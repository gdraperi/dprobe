// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package netutil provides network utility functions, complementing the more
// common ones in the net package.
package netutil // import "golang.org/x/net/netutil"

import (
	"net"
	"sync"
)

// LimitListener returns a Listener that accepts at most n simultaneous
// connections from the provided Listener.
func LimitListener(l net.Listener, n int) net.Listener ***REMOVED***
	return &limitListener***REMOVED***l, make(chan struct***REMOVED******REMOVED***, n)***REMOVED***
***REMOVED***

type limitListener struct ***REMOVED***
	net.Listener
	sem chan struct***REMOVED******REMOVED***
***REMOVED***

func (l *limitListener) acquire() ***REMOVED*** l.sem <- struct***REMOVED******REMOVED******REMOVED******REMOVED*** ***REMOVED***
func (l *limitListener) release() ***REMOVED*** <-l.sem ***REMOVED***

func (l *limitListener) Accept() (net.Conn, error) ***REMOVED***
	l.acquire()
	c, err := l.Listener.Accept()
	if err != nil ***REMOVED***
		l.release()
		return nil, err
	***REMOVED***
	return &limitListenerConn***REMOVED***Conn: c, release: l.release***REMOVED***, nil
***REMOVED***

type limitListenerConn struct ***REMOVED***
	net.Conn
	releaseOnce sync.Once
	release     func()
***REMOVED***

func (l *limitListenerConn) Close() error ***REMOVED***
	err := l.Conn.Close()
	l.releaseOnce.Do(l.release)
	return err
***REMOVED***
