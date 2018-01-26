package memberlist

import (
	"fmt"
	"net"
)

func LogAddress(addr net.Addr) string ***REMOVED***
	if addr == nil ***REMOVED***
		return "from=<unknown address>"
	***REMOVED***

	return fmt.Sprintf("from=%s", addr.String())
***REMOVED***

func LogConn(conn net.Conn) string ***REMOVED***
	if conn == nil ***REMOVED***
		return LogAddress(nil)
	***REMOVED***

	return LogAddress(conn.RemoteAddr())
***REMOVED***
