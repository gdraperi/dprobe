package sockets

import (
	"errors"
	"net"
	"sync"
)

var errClosed = errors.New("use of closed network connection")

// InmemSocket implements net.Listener using in-memory only connections.
type InmemSocket struct ***REMOVED***
	chConn  chan net.Conn
	chClose chan struct***REMOVED******REMOVED***
	addr    string
	mu      sync.Mutex
***REMOVED***

// dummyAddr is used to satisfy net.Addr for the in-mem socket
// it is just stored as a string and returns the string for all calls
type dummyAddr string

// NewInmemSocket creates an in-memory only net.Listener
// The addr argument can be any string, but is used to satisfy the `Addr()` part
// of the net.Listener interface
func NewInmemSocket(addr string, bufSize int) *InmemSocket ***REMOVED***
	return &InmemSocket***REMOVED***
		chConn:  make(chan net.Conn, bufSize),
		chClose: make(chan struct***REMOVED******REMOVED***),
		addr:    addr,
	***REMOVED***
***REMOVED***

// Addr returns the socket's addr string to satisfy net.Listener
func (s *InmemSocket) Addr() net.Addr ***REMOVED***
	return dummyAddr(s.addr)
***REMOVED***

// Accept implements the Accept method in the Listener interface; it waits for the next call and returns a generic Conn.
func (s *InmemSocket) Accept() (net.Conn, error) ***REMOVED***
	select ***REMOVED***
	case conn := <-s.chConn:
		return conn, nil
	case <-s.chClose:
		return nil, errClosed
	***REMOVED***
***REMOVED***

// Close closes the listener. It will be unavailable for use once closed.
func (s *InmemSocket) Close() error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	select ***REMOVED***
	case <-s.chClose:
	default:
		close(s.chClose)
	***REMOVED***
	return nil
***REMOVED***

// Dial is used to establish a connection with the in-mem server
func (s *InmemSocket) Dial(network, addr string) (net.Conn, error) ***REMOVED***
	srvConn, clientConn := net.Pipe()
	select ***REMOVED***
	case s.chConn <- srvConn:
	case <-s.chClose:
		return nil, errClosed
	***REMOVED***

	return clientConn, nil
***REMOVED***

// Network returns the addr string, satisfies net.Addr
func (a dummyAddr) Network() string ***REMOVED***
	return string(a)
***REMOVED***

// String returns the string form
func (a dummyAddr) String() string ***REMOVED***
	return string(a)
***REMOVED***
