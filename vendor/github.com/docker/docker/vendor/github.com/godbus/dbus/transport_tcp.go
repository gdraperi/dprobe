//+build !windows

package dbus

import (
	"errors"
	"net"
)

func init() ***REMOVED***
	transports["tcp"] = newTcpTransport
***REMOVED***

func tcpFamily(keys string) (string, error) ***REMOVED***
	switch getKey(keys, "family") ***REMOVED***
	case "":
		return "tcp", nil
	case "ipv4":
		return "tcp4", nil
	case "ipv6":
		return "tcp6", nil
	default:
		return "", errors.New("dbus: invalid tcp family (must be ipv4 or ipv6)")
	***REMOVED***
***REMOVED***

func newTcpTransport(keys string) (transport, error) ***REMOVED***
	host := getKey(keys, "host")
	port := getKey(keys, "port")
	if host == "" || port == "" ***REMOVED***
		return nil, errors.New("dbus: unsupported address (must set host and port)")
	***REMOVED***

	protocol, err := tcpFamily(keys)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	socket, err := net.Dial(protocol, net.JoinHostPort(host, port))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewConn(socket)
***REMOVED***
