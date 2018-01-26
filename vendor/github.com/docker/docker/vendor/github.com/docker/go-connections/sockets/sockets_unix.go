// +build !windows

package sockets

import (
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"
)

const maxUnixSocketPathSize = len(syscall.RawSockaddrUnix***REMOVED******REMOVED***.Path)

func configureUnixTransport(tr *http.Transport, proto, addr string) error ***REMOVED***
	if len(addr) > maxUnixSocketPathSize ***REMOVED***
		return fmt.Errorf("Unix socket path %q is too long", addr)
	***REMOVED***
	// No need for compression in local communications.
	tr.DisableCompression = true
	tr.Dial = func(_, _ string) (net.Conn, error) ***REMOVED***
		return net.DialTimeout(proto, addr, defaultTimeout)
	***REMOVED***
	return nil
***REMOVED***

func configureNpipeTransport(tr *http.Transport, proto, addr string) error ***REMOVED***
	return ErrProtocolNotAvailable
***REMOVED***

// DialPipe connects to a Windows named pipe.
// This is not supported on other OSes.
func DialPipe(_ string, _ time.Duration) (net.Conn, error) ***REMOVED***
	return nil, syscall.EAFNOSUPPORT
***REMOVED***
