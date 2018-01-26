// Package sockets provides helper functions to create and configure Unix or TCP sockets.
package sockets

import (
	"crypto/tls"
	"net"
)

// NewTCPSocket creates a TCP socket listener with the specified address and
// the specified tls configuration. If TLSConfig is set, will encapsulate the
// TCP listener inside a TLS one.
func NewTCPSocket(addr string, tlsConfig *tls.Config) (net.Listener, error) ***REMOVED***
	l, err := net.Listen("tcp", addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if tlsConfig != nil ***REMOVED***
		tlsConfig.NextProtos = []string***REMOVED***"http/1.1"***REMOVED***
		l = tls.NewListener(l, tlsConfig)
	***REMOVED***
	return l, nil
***REMOVED***
