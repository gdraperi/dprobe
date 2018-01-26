// +build !appengine,!plan9

package request

import (
	"net"
	"os"
	"syscall"
)

func isErrConnectionReset(err error) bool ***REMOVED***
	if opErr, ok := err.(*net.OpError); ok ***REMOVED***
		if sysErr, ok := opErr.Err.(*os.SyscallError); ok ***REMOVED***
			return sysErr.Err == syscall.ECONNRESET
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***
