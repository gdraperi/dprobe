package restartmanager

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
)

const (
	backoffMultiplier = 2
	defaultTimeout    = 100 * time.Millisecond
	maxRestartTimeout = 1 * time.Minute
)

// ErrRestartCanceled is returned when the restart manager has been
// canceled and will no longer restart the container.
var ErrRestartCanceled = errors.New("restart canceled")

// RestartManager defines object that controls container restarting rules.
type RestartManager interface ***REMOVED***
	Cancel() error
	ShouldRestart(exitCode uint32, hasBeenManuallyStopped bool, executionDuration time.Duration) (bool, chan error, error)
***REMOVED***

type restartManager struct ***REMOVED***
	sync.Mutex
	sync.Once
	policy       container.RestartPolicy
	restartCount int
	timeout      time.Duration
	active       bool
	cancel       chan struct***REMOVED******REMOVED***
	canceled     bool
***REMOVED***

// New returns a new restartManager based on a policy.
func New(policy container.RestartPolicy, restartCount int) RestartManager ***REMOVED***
	return &restartManager***REMOVED***policy: policy, restartCount: restartCount, cancel: make(chan struct***REMOVED******REMOVED***)***REMOVED***
***REMOVED***

func (rm *restartManager) SetPolicy(policy container.RestartPolicy) ***REMOVED***
	rm.Lock()
	rm.policy = policy
	rm.Unlock()
***REMOVED***

func (rm *restartManager) ShouldRestart(exitCode uint32, hasBeenManuallyStopped bool, executionDuration time.Duration) (bool, chan error, error) ***REMOVED***
	if rm.policy.IsNone() ***REMOVED***
		return false, nil, nil
	***REMOVED***
	rm.Lock()
	unlockOnExit := true
	defer func() ***REMOVED***
		if unlockOnExit ***REMOVED***
			rm.Unlock()
		***REMOVED***
	***REMOVED***()

	if rm.canceled ***REMOVED***
		return false, nil, ErrRestartCanceled
	***REMOVED***

	if rm.active ***REMOVED***
		return false, nil, fmt.Errorf("invalid call on an active restart manager")
	***REMOVED***
	// if the container ran for more than 10s, regardless of status and policy reset the
	// the timeout back to the default.
	if executionDuration.Seconds() >= 10 ***REMOVED***
		rm.timeout = 0
	***REMOVED***
	switch ***REMOVED***
	case rm.timeout == 0:
		rm.timeout = defaultTimeout
	case rm.timeout < maxRestartTimeout:
		rm.timeout *= backoffMultiplier
	***REMOVED***
	if rm.timeout > maxRestartTimeout ***REMOVED***
		rm.timeout = maxRestartTimeout
	***REMOVED***

	var restart bool
	switch ***REMOVED***
	case rm.policy.IsAlways():
		restart = true
	case rm.policy.IsUnlessStopped() && !hasBeenManuallyStopped:
		restart = true
	case rm.policy.IsOnFailure():
		// the default value of 0 for MaximumRetryCount means that we will not enforce a maximum count
		if max := rm.policy.MaximumRetryCount; max == 0 || rm.restartCount < max ***REMOVED***
			restart = exitCode != 0
		***REMOVED***
	***REMOVED***

	if !restart ***REMOVED***
		rm.active = false
		return false, nil, nil
	***REMOVED***

	rm.restartCount++

	unlockOnExit = false
	rm.active = true
	rm.Unlock()

	ch := make(chan error)
	go func() ***REMOVED***
		select ***REMOVED***
		case <-rm.cancel:
			ch <- ErrRestartCanceled
			close(ch)
		case <-time.After(rm.timeout):
			rm.Lock()
			close(ch)
			rm.active = false
			rm.Unlock()
		***REMOVED***
	***REMOVED***()

	return true, ch, nil
***REMOVED***

func (rm *restartManager) Cancel() error ***REMOVED***
	rm.Do(func() ***REMOVED***
		rm.Lock()
		rm.canceled = true
		close(rm.cancel)
		rm.Unlock()
	***REMOVED***)
	return nil
***REMOVED***
