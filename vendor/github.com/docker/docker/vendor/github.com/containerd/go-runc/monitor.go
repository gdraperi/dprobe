package runc

import (
	"os/exec"
	"syscall"
	"time"
)

var Monitor ProcessMonitor = &defaultMonitor***REMOVED******REMOVED***

type Exit struct ***REMOVED***
	Timestamp time.Time
	Pid       int
	Status    int
***REMOVED***

// ProcessMonitor is an interface for process monitoring
//
// It allows daemons using go-runc to have a SIGCHLD handler
// to handle exits without introducing races between the handler
// and go's exec.Cmd
// These methods should match the methods exposed by exec.Cmd to provide
// a consistent experience for the caller
type ProcessMonitor interface ***REMOVED***
	Start(*exec.Cmd) (chan Exit, error)
	Wait(*exec.Cmd, chan Exit) (int, error)
***REMOVED***

type defaultMonitor struct ***REMOVED***
***REMOVED***

func (m *defaultMonitor) Start(c *exec.Cmd) (chan Exit, error) ***REMOVED***
	if err := c.Start(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ec := make(chan Exit, 1)
	go func() ***REMOVED***
		var status int
		if err := c.Wait(); err != nil ***REMOVED***
			status = 255
			if exitErr, ok := err.(*exec.ExitError); ok ***REMOVED***
				if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok ***REMOVED***
					status = ws.ExitStatus()
				***REMOVED***
			***REMOVED***
		***REMOVED***
		ec <- Exit***REMOVED***
			Timestamp: time.Now(),
			Pid:       c.Process.Pid,
			Status:    status,
		***REMOVED***
		close(ec)
	***REMOVED***()
	return ec, nil
***REMOVED***

func (m *defaultMonitor) Wait(c *exec.Cmd, ec chan Exit) (int, error) ***REMOVED***
	e := <-ec
	return e.Status, nil
***REMOVED***
