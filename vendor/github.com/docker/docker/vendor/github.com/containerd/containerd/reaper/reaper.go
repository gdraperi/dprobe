// +build !windows

package reaper

import (
	"os/exec"
	"sync"
	"time"

	"github.com/containerd/containerd/sys"
	runc "github.com/containerd/go-runc"
	"github.com/pkg/errors"
)

// ErrNoSuchProcess is returned when the process no longer exists
var ErrNoSuchProcess = errors.New("no such process")

const bufferSize = 1024

// Reap should be called when the process receives an SIGCHLD.  Reap will reap
// all exited processes and close their wait channels
func Reap() error ***REMOVED***
	now := time.Now()
	exits, err := sys.Reap(false)
	Default.Lock()
	for c := range Default.subscribers ***REMOVED***
		for _, e := range exits ***REMOVED***
			c <- runc.Exit***REMOVED***
				Timestamp: now,
				Pid:       e.Pid,
				Status:    e.Status,
			***REMOVED***
		***REMOVED***

	***REMOVED***
	Default.Unlock()
	return err
***REMOVED***

// Default is the default monitor initialized for the package
var Default = &Monitor***REMOVED***
	subscribers: make(map[chan runc.Exit]struct***REMOVED******REMOVED***),
***REMOVED***

// Monitor monitors the underlying system for process status changes
type Monitor struct ***REMOVED***
	sync.Mutex

	subscribers map[chan runc.Exit]struct***REMOVED******REMOVED***
***REMOVED***

// Start starts the command a registers the process with the reaper
func (m *Monitor) Start(c *exec.Cmd) (chan runc.Exit, error) ***REMOVED***
	ec := m.Subscribe()
	if err := c.Start(); err != nil ***REMOVED***
		m.Unsubscribe(ec)
		return nil, err
	***REMOVED***
	return ec, nil
***REMOVED***

// Wait blocks until a process is signal as dead.
// User should rely on the value of the exit status to determine if the
// command was successful or not.
func (m *Monitor) Wait(c *exec.Cmd, ec chan runc.Exit) (int, error) ***REMOVED***
	for e := range ec ***REMOVED***
		if e.Pid == c.Process.Pid ***REMOVED***
			// make sure we flush all IO
			c.Wait()
			m.Unsubscribe(ec)
			return e.Status, nil
		***REMOVED***
	***REMOVED***
	// return no such process if the ec channel is closed and no more exit
	// events will be sent
	return -1, ErrNoSuchProcess
***REMOVED***

// Subscribe to process exit changes
func (m *Monitor) Subscribe() chan runc.Exit ***REMOVED***
	c := make(chan runc.Exit, bufferSize)
	m.Lock()
	m.subscribers[c] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	m.Unlock()
	return c
***REMOVED***

// Unsubscribe to process exit changes
func (m *Monitor) Unsubscribe(c chan runc.Exit) ***REMOVED***
	m.Lock()
	delete(m.subscribers, c)
	close(c)
	m.Unlock()
***REMOVED***
