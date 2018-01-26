// +build linux freebsd darwin

package system

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// IsProcessAlive returns true if process with a given pid is running.
func IsProcessAlive(pid int) bool ***REMOVED***
	err := unix.Kill(pid, syscall.Signal(0))
	if err == nil || err == unix.EPERM ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

// KillProcess force-stops a process.
func KillProcess(pid int) ***REMOVED***
	unix.Kill(pid, unix.SIGKILL)
***REMOVED***
