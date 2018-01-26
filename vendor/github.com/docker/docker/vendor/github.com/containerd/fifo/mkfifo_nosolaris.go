// +build !solaris

package fifo

import "syscall"

func mkfifo(path string, mode uint32) (err error) ***REMOVED***
	return syscall.Mkfifo(path, mode)
***REMOVED***
