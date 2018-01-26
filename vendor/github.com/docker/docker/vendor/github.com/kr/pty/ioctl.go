package pty

import "syscall"

func ioctl(fd, cmd, ptr uintptr) error ***REMOVED***
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, ptr)
	if e != 0 ***REMOVED***
		return e
	***REMOVED***
	return nil
***REMOVED***
