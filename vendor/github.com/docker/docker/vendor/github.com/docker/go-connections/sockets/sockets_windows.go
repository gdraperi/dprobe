package sockets

import (
	"net"
	"net/http"
	"time"

	"github.com/Microsoft/go-winio"
)

func configureUnixTransport(tr *http.Transport, proto, addr string) error ***REMOVED***
	return ErrProtocolNotAvailable
***REMOVED***

func configureNpipeTransport(tr *http.Transport, proto, addr string) error ***REMOVED***
	// No need for compression in local communications.
	tr.DisableCompression = true
	tr.Dial = func(_, _ string) (net.Conn, error) ***REMOVED***
		return DialPipe(addr, defaultTimeout)
	***REMOVED***
	return nil
***REMOVED***

// DialPipe connects to a Windows named pipe.
func DialPipe(addr string, timeout time.Duration) (net.Conn, error) ***REMOVED***
	return winio.DialPipe(addr, &timeout)
***REMOVED***
