package request

import (
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

func npipeDial(path string, timeout time.Duration) (net.Conn, error) ***REMOVED***
	return winio.DialPipe(path, &timeout)
***REMOVED***
