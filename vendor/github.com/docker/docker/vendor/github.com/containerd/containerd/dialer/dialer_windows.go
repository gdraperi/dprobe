package dialer

import (
	"net"
	"os"
	"syscall"
	"time"

	winio "github.com/Microsoft/go-winio"
)

func isNoent(err error) bool ***REMOVED***
	if err != nil ***REMOVED***
		if oerr, ok := err.(*os.PathError); ok ***REMOVED***
			if oerr.Err == syscall.ENOENT ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func dialer(address string, timeout time.Duration) (net.Conn, error) ***REMOVED***
	return winio.DialPipe(address, &timeout)
***REMOVED***

// DialAddress returns the dial address
func DialAddress(address string) string ***REMOVED***
	return address
***REMOVED***
