package libcontainerd

import "syscall"

// containerdSysProcAttr returns the SysProcAttr to use when exec'ing
// containerd
func containerdSysProcAttr() *syscall.SysProcAttr ***REMOVED***
	return &syscall.SysProcAttr***REMOVED***
		Setsid:    true,
		Pdeathsig: syscall.SIGKILL,
	***REMOVED***
***REMOVED***
