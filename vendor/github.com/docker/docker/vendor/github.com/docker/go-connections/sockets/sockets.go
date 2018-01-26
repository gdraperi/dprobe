// Package sockets provides helper functions to create and configure Unix or TCP sockets.
package sockets

import (
	"errors"
	"net"
	"net/http"
	"time"
)

// Why 32? See https://github.com/docker/docker/pull/8035.
const defaultTimeout = 32 * time.Second

// ErrProtocolNotAvailable is returned when a given transport protocol is not provided by the operating system.
var ErrProtocolNotAvailable = errors.New("protocol not available")

// ConfigureTransport configures the specified Transport according to the
// specified proto and addr.
// If the proto is unix (using a unix socket to communicate) or npipe the
// compression is disabled.
func ConfigureTransport(tr *http.Transport, proto, addr string) error ***REMOVED***
	switch proto ***REMOVED***
	case "unix":
		return configureUnixTransport(tr, proto, addr)
	case "npipe":
		return configureNpipeTransport(tr, proto, addr)
	default:
		tr.Proxy = http.ProxyFromEnvironment
		dialer, err := DialerFromEnvironment(&net.Dialer***REMOVED***
			Timeout: defaultTimeout,
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		tr.Dial = dialer.Dial
	***REMOVED***
	return nil
***REMOVED***
