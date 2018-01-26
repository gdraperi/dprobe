package container

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-units"
)

// State holds the current container state, and has methods to get and
// set the state. Container has an embed, which allows all of the
// functions defined against State to run against Container.
type State struct ***REMOVED***
	sync.Mutex
	// Note that `Running` and `Paused` are not mutually exclusive:
	// When pausing a container (on Linux), the cgroups freezer is used to suspend
	// all processes in the container. Freezing the process requires the process to
	// be running. As a result, paused containers are both `Running` _and_ `Paused`.
	Running           bool
	Paused            bool
	Restarting        bool
	OOMKilled         bool
	RemovalInProgress bool // Not need for this to be persistent on disk.
	Dead              bool
	Pid               int
	ExitCodeValue     int    `json:"ExitCode"`
	ErrorMsg          string `json:"Error"` // contains last known error during container start or remove
	StartedAt         time.Time
	FinishedAt        time.Time
	Health            *Health

	waitStop   chan struct***REMOVED******REMOVED***
	waitRemove chan struct***REMOVED******REMOVED***
***REMOVED***

// StateStatus is used to return container wait results.
// Implements exec.ExitCode interface.
// This type is needed as State include a sync.Mutex field which make
// copying it unsafe.
type StateStatus struct ***REMOVED***
	exitCode int
	err      error
***REMOVED***

// ExitCode returns current exitcode for the state.
func (s StateStatus) ExitCode() int ***REMOVED***
	return s.exitCode
***REMOVED***

// Err returns current error for the state. Returns nil if the container had
// exited on its own.
func (s StateStatus) Err() error ***REMOVED***
	return s.err
***REMOVED***

// NewState creates a default state object with a fresh channel for state changes.
func NewState() *State ***REMOVED***
	return &State***REMOVED***
		waitStop:   make(chan struct***REMOVED******REMOVED***),
		waitRemove: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// String returns a human-readable description of the state
func (s *State) String() string ***REMOVED***
	if s.Running ***REMOVED***
		if s.Paused ***REMOVED***
			return fmt.Sprintf("Up %s (Paused)", units.HumanDuration(time.Now().UTC().Sub(s.StartedAt)))
		***REMOVED***
		if s.Restarting ***REMOVED***
			return fmt.Sprintf("Restarting (%d) %s ago", s.ExitCodeValue, units.HumanDuration(time.Now().UTC().Sub(s.FinishedAt)))
		***REMOVED***

		if h := s.Health; h != nil ***REMOVED***
			return fmt.Sprintf("Up %s (%s)", units.HumanDuration(time.Now().UTC().Sub(s.StartedAt)), h.String())
		***REMOVED***

		return fmt.Sprintf("Up %s", units.HumanDuration(time.Now().UTC().Sub(s.StartedAt)))
	***REMOVED***

	if s.RemovalInProgress ***REMOVED***
		return "Removal In Progress"
	***REMOVED***

	if s.Dead ***REMOVED***
		return "Dead"
	***REMOVED***

	if s.StartedAt.IsZero() ***REMOVED***
		return "Created"
	***REMOVED***

	if s.FinishedAt.IsZero() ***REMOVED***
		return ""
	***REMOVED***

	return fmt.Sprintf("Exited (%d) %s ago", s.ExitCodeValue, units.HumanDuration(time.Now().UTC().Sub(s.FinishedAt)))
***REMOVED***

// IsValidHealthString checks if the provided string is a valid container health status or not.
func IsValidHealthString(s string) bool ***REMOVED***
	return s == types.Starting ||
		s == types.Healthy ||
		s == types.Unhealthy ||
		s == types.NoHealthcheck
***REMOVED***

// StateString returns a single string to describe state
func (s *State) StateString() string ***REMOVED***
	if s.Running ***REMOVED***
		if s.Paused ***REMOVED***
			return "paused"
		***REMOVED***
		if s.Restarting ***REMOVED***
			return "restarting"
		***REMOVED***
		return "running"
	***REMOVED***

	if s.RemovalInProgress ***REMOVED***
		return "removing"
	***REMOVED***

	if s.Dead ***REMOVED***
		return "dead"
	***REMOVED***

	if s.StartedAt.IsZero() ***REMOVED***
		return "created"
	***REMOVED***

	return "exited"
***REMOVED***

// IsValidStateString checks if the provided string is a valid container state or not.
func IsValidStateString(s string) bool ***REMOVED***
	if s != "paused" &&
		s != "restarting" &&
		s != "removing" &&
		s != "running" &&
		s != "dead" &&
		s != "created" &&
		s != "exited" ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// WaitCondition is an enum type for different states to wait for.
type WaitCondition int

// Possible WaitCondition Values.
//
// WaitConditionNotRunning (default) is used to wait for any of the non-running
// states: "created", "exited", "dead", "removing", or "removed".
//
// WaitConditionNextExit is used to wait for the next time the state changes
// to a non-running state. If the state is currently "created" or "exited",
// this would cause Wait() to block until either the container runs and exits
// or is removed.
//
// WaitConditionRemoved is used to wait for the container to be removed.
const (
	WaitConditionNotRunning WaitCondition = iota
	WaitConditionNextExit
	WaitConditionRemoved
)

// Wait waits until the container is in a certain state indicated by the given
// condition. A context must be used for cancelling the request, controlling
// timeouts, and avoiding goroutine leaks. Wait must be called without holding
// the state lock. Returns a channel from which the caller will receive the
// result. If the container exited on its own, the result's Err() method will
// be nil and its ExitCode() method will return the container's exit code,
// otherwise, the results Err() method will return an error indicating why the
// wait operation failed.
func (s *State) Wait(ctx context.Context, condition WaitCondition) <-chan StateStatus ***REMOVED***
	s.Lock()
	defer s.Unlock()

	if condition == WaitConditionNotRunning && !s.Running ***REMOVED***
		// Buffer so we can put it in the channel now.
		resultC := make(chan StateStatus, 1)

		// Send the current status.
		resultC <- StateStatus***REMOVED***
			exitCode: s.ExitCode(),
			err:      s.Err(),
		***REMOVED***

		return resultC
	***REMOVED***

	// If we are waiting only for removal, the waitStop channel should
	// remain nil and block forever.
	var waitStop chan struct***REMOVED******REMOVED***
	if condition < WaitConditionRemoved ***REMOVED***
		waitStop = s.waitStop
	***REMOVED***

	// Always wait for removal, just in case the container gets removed
	// while it is still in a "created" state, in which case it is never
	// actually stopped.
	waitRemove := s.waitRemove

	resultC := make(chan StateStatus)

	go func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			// Context timeout or cancellation.
			resultC <- StateStatus***REMOVED***
				exitCode: -1,
				err:      ctx.Err(),
			***REMOVED***
			return
		case <-waitStop:
		case <-waitRemove:
		***REMOVED***

		s.Lock()
		result := StateStatus***REMOVED***
			exitCode: s.ExitCode(),
			err:      s.Err(),
		***REMOVED***
		s.Unlock()

		resultC <- result
	***REMOVED***()

	return resultC
***REMOVED***

// IsRunning returns whether the running flag is set. Used by Container to check whether a container is running.
func (s *State) IsRunning() bool ***REMOVED***
	s.Lock()
	res := s.Running
	s.Unlock()
	return res
***REMOVED***

// GetPID holds the process id of a container.
func (s *State) GetPID() int ***REMOVED***
	s.Lock()
	res := s.Pid
	s.Unlock()
	return res
***REMOVED***

// ExitCode returns current exitcode for the state. Take lock before if state
// may be shared.
func (s *State) ExitCode() int ***REMOVED***
	return s.ExitCodeValue
***REMOVED***

// SetExitCode sets current exitcode for the state. Take lock before if state
// may be shared.
func (s *State) SetExitCode(ec int) ***REMOVED***
	s.ExitCodeValue = ec
***REMOVED***

// SetRunning sets the state of the container to "running".
func (s *State) SetRunning(pid int, initial bool) ***REMOVED***
	s.ErrorMsg = ""
	s.Paused = false
	s.Running = true
	s.Restarting = false
	if initial ***REMOVED***
		s.Paused = false
	***REMOVED***
	s.ExitCodeValue = 0
	s.Pid = pid
	if initial ***REMOVED***
		s.StartedAt = time.Now().UTC()
	***REMOVED***
***REMOVED***

// SetStopped sets the container state to "stopped" without locking.
func (s *State) SetStopped(exitStatus *ExitStatus) ***REMOVED***
	s.Running = false
	s.Paused = false
	s.Restarting = false
	s.Pid = 0
	if exitStatus.ExitedAt.IsZero() ***REMOVED***
		s.FinishedAt = time.Now().UTC()
	***REMOVED*** else ***REMOVED***
		s.FinishedAt = exitStatus.ExitedAt
	***REMOVED***
	s.ExitCodeValue = exitStatus.ExitCode
	s.OOMKilled = exitStatus.OOMKilled
	close(s.waitStop) // fire waiters for stop
	s.waitStop = make(chan struct***REMOVED******REMOVED***)
***REMOVED***

// SetRestarting sets the container state to "restarting" without locking.
// It also sets the container PID to 0.
func (s *State) SetRestarting(exitStatus *ExitStatus) ***REMOVED***
	// we should consider the container running when it is restarting because of
	// all the checks in docker around rm/stop/etc
	s.Running = true
	s.Restarting = true
	s.Paused = false
	s.Pid = 0
	s.FinishedAt = time.Now().UTC()
	s.ExitCodeValue = exitStatus.ExitCode
	s.OOMKilled = exitStatus.OOMKilled
	close(s.waitStop) // fire waiters for stop
	s.waitStop = make(chan struct***REMOVED******REMOVED***)
***REMOVED***

// SetError sets the container's error state. This is useful when we want to
// know the error that occurred when container transits to another state
// when inspecting it
func (s *State) SetError(err error) ***REMOVED***
	s.ErrorMsg = ""
	if err != nil ***REMOVED***
		s.ErrorMsg = err.Error()
	***REMOVED***
***REMOVED***

// IsPaused returns whether the container is paused or not.
func (s *State) IsPaused() bool ***REMOVED***
	s.Lock()
	res := s.Paused
	s.Unlock()
	return res
***REMOVED***

// IsRestarting returns whether the container is restarting or not.
func (s *State) IsRestarting() bool ***REMOVED***
	s.Lock()
	res := s.Restarting
	s.Unlock()
	return res
***REMOVED***

// SetRemovalInProgress sets the container state as being removed.
// It returns true if the container was already in that state.
func (s *State) SetRemovalInProgress() bool ***REMOVED***
	s.Lock()
	defer s.Unlock()
	if s.RemovalInProgress ***REMOVED***
		return true
	***REMOVED***
	s.RemovalInProgress = true
	return false
***REMOVED***

// ResetRemovalInProgress makes the RemovalInProgress state to false.
func (s *State) ResetRemovalInProgress() ***REMOVED***
	s.Lock()
	s.RemovalInProgress = false
	s.Unlock()
***REMOVED***

// IsRemovalInProgress returns whether the RemovalInProgress flag is set.
// Used by Container to check whether a container is being removed.
func (s *State) IsRemovalInProgress() bool ***REMOVED***
	s.Lock()
	res := s.RemovalInProgress
	s.Unlock()
	return res
***REMOVED***

// SetDead sets the container state to "dead"
func (s *State) SetDead() ***REMOVED***
	s.Lock()
	s.Dead = true
	s.Unlock()
***REMOVED***

// IsDead returns whether the Dead flag is set. Used by Container to check whether a container is dead.
func (s *State) IsDead() bool ***REMOVED***
	s.Lock()
	res := s.Dead
	s.Unlock()
	return res
***REMOVED***

// SetRemoved assumes this container is already in the "dead" state and
// closes the internal waitRemove channel to unblock callers waiting for a
// container to be removed.
func (s *State) SetRemoved() ***REMOVED***
	s.SetRemovalError(nil)
***REMOVED***

// SetRemovalError is to be called in case a container remove failed.
// It sets an error and closes the internal waitRemove channel to unblock
// callers waiting for the container to be removed.
func (s *State) SetRemovalError(err error) ***REMOVED***
	s.SetError(err)
	s.Lock()
	close(s.waitRemove) // Unblock those waiting on remove.
	// Recreate the channel so next ContainerWait will work
	s.waitRemove = make(chan struct***REMOVED******REMOVED***)
	s.Unlock()
***REMOVED***

// Err returns an error if there is one.
func (s *State) Err() error ***REMOVED***
	if s.ErrorMsg != "" ***REMOVED***
		return errors.New(s.ErrorMsg)
	***REMOVED***
	return nil
***REMOVED***
