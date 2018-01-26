// +build !windows

package request

import (
	"net"
	"time"
)

func npipeDial(path string, timeout time.Duration) (net.Conn, error) ***REMOVED***
	panic("npipe protocol only supported on Windows")
***REMOVED***
