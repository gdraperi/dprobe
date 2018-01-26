// +build !linux

package dns

import (
	"net"
	"syscall"
)

// These do nothing. See udp_linux.go for an example of how to implement this.

// We tried to adhire to some kind of naming scheme.

func setUDPSocketOptions4(conn *net.UDPConn) error                 ***REMOVED*** return nil ***REMOVED***
func setUDPSocketOptions6(conn *net.UDPConn) error                 ***REMOVED*** return nil ***REMOVED***
func getUDPSocketOptions6Only(conn *net.UDPConn) (bool, error)     ***REMOVED*** return false, nil ***REMOVED***
func getUDPSocketName(conn *net.UDPConn) (syscall.Sockaddr, error) ***REMOVED*** return nil, nil ***REMOVED***
