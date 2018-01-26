// +build !windows

package xnet

import (
	"net"
	"time"
)

// ListenLocal opens a local socket for control communication
func ListenLocal(socket string) (net.Listener, error) ***REMOVED***
	// on unix it's just a unix socket
	return net.Listen("unix", socket)
***REMOVED***

// DialTimeoutLocal is a DialTimeout function for local sockets
func DialTimeoutLocal(socket string, timeout time.Duration) (net.Conn, error) ***REMOVED***
	// on unix, we dial a unix socket
	return net.DialTimeout("unix", socket, timeout)
***REMOVED***
