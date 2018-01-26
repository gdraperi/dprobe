// +build !linux,!windows

package client

import (
	"os/exec"
	"syscall"
)

func getSysProcAttr() *syscall.SysProcAttr ***REMOVED***
	return &syscall.SysProcAttr***REMOVED***
		Setpgid: true,
	***REMOVED***
***REMOVED***

func setCgroup(cgroupPath string, cmd *exec.Cmd) error ***REMOVED***
	return nil
***REMOVED***
