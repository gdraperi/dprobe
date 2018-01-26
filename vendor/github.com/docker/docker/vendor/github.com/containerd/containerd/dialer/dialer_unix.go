// +build !windows

package dialer

import (
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

// DialAddress returns the address with unix:// prepended to the
// provided address
func DialAddress(address string) string ***REMOVED***
	return fmt.Sprintf("unix://%s", address)
***REMOVED***

func isNoent(err error) bool ***REMOVED***
	if err != nil ***REMOVED***
		if nerr, ok := err.(*net.OpError); ok ***REMOVED***
			if serr, ok := nerr.Err.(*os.SyscallError); ok ***REMOVED***
				if serr.Err == syscall.ENOENT ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func dialer(address string, timeout time.Duration) (net.Conn, error) ***REMOVED***
	address = strings.TrimPrefix(address, "unix://")
	return net.DialTimeout("unix", address, timeout)
***REMOVED***
