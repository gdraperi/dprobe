package srslog

import (
	"crypto/tls"
	"net"
)

// dialerFunctionWrapper is a simple object that consists of a dialer function
// and its name. This is primarily for testing, so we can make sure that the
// getDialer method returns the correct dialer function. However, if you ever
// find that you need to check which dialer function you have, this would also
// be useful for you without having to use reflection.
type dialerFunctionWrapper struct ***REMOVED***
	Name   string
	Dialer func() (serverConn, string, error)
***REMOVED***

// Call the wrapped dialer function and return its return values.
func (df dialerFunctionWrapper) Call() (serverConn, string, error) ***REMOVED***
	return df.Dialer()
***REMOVED***

// getDialer returns a "dialer" function that can be called to connect to a
// syslog server.
//
// Each dialer function is responsible for dialing the remote host and returns
// a serverConn, the hostname (or a default if the Writer has not specified a
// hostname), and an error in case dialing fails.
//
// The reason for separate dialers is that different network types may need
// to dial their connection differently, yet still provide a net.Conn interface
// that you can use once they have dialed. Rather than an increasingly long
// conditional, we have a map of network -> dialer function (with a sane default
// value), and adding a new network type is as easy as writing the dialer
// function and adding it to the map.
func (w *Writer) getDialer() dialerFunctionWrapper ***REMOVED***
	dialers := map[string]dialerFunctionWrapper***REMOVED***
		"":        dialerFunctionWrapper***REMOVED***"unixDialer", w.unixDialer***REMOVED***,
		"tcp+tls": dialerFunctionWrapper***REMOVED***"tlsDialer", w.tlsDialer***REMOVED***,
	***REMOVED***
	dialer, ok := dialers[w.network]
	if !ok ***REMOVED***
		dialer = dialerFunctionWrapper***REMOVED***"basicDialer", w.basicDialer***REMOVED***
	***REMOVED***
	return dialer
***REMOVED***

// unixDialer uses the unixSyslog method to open a connection to the syslog
// daemon running on the local machine.
func (w *Writer) unixDialer() (serverConn, string, error) ***REMOVED***
	sc, err := unixSyslog()
	hostname := w.hostname
	if hostname == "" ***REMOVED***
		hostname = "localhost"
	***REMOVED***
	return sc, hostname, err
***REMOVED***

// tlsDialer connects to TLS over TCP, and is used for the "tcp+tls" network
// type.
func (w *Writer) tlsDialer() (serverConn, string, error) ***REMOVED***
	c, err := tls.Dial("tcp", w.raddr, w.tlsConfig)
	var sc serverConn
	hostname := w.hostname
	if err == nil ***REMOVED***
		sc = &netConn***REMOVED***conn: c***REMOVED***
		if hostname == "" ***REMOVED***
			hostname = c.LocalAddr().String()
		***REMOVED***
	***REMOVED***
	return sc, hostname, err
***REMOVED***

// basicDialer is the most common dialer for syslog, and supports both TCP and
// UDP connections.
func (w *Writer) basicDialer() (serverConn, string, error) ***REMOVED***
	c, err := net.Dial(w.network, w.raddr)
	var sc serverConn
	hostname := w.hostname
	if err == nil ***REMOVED***
		sc = &netConn***REMOVED***conn: c***REMOVED***
		if hostname == "" ***REMOVED***
			hostname = c.LocalAddr().String()
		***REMOVED***
	***REMOVED***
	return sc, hostname, err
***REMOVED***
