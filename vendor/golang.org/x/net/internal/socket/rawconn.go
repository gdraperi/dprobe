// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9

package socket

import (
	"errors"
	"net"
	"os"
	"syscall"
)

// A Conn represents a raw connection.
type Conn struct ***REMOVED***
	network string
	c       syscall.RawConn
***REMOVED***

// NewConn returns a new raw connection.
func NewConn(c net.Conn) (*Conn, error) ***REMOVED***
	var err error
	var cc Conn
	switch c := c.(type) ***REMOVED***
	case *net.TCPConn:
		cc.network = "tcp"
		cc.c, err = c.SyscallConn()
	case *net.UDPConn:
		cc.network = "udp"
		cc.c, err = c.SyscallConn()
	case *net.IPConn:
		cc.network = "ip"
		cc.c, err = c.SyscallConn()
	default:
		return nil, errors.New("unknown connection type")
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &cc, nil
***REMOVED***

func (o *Option) get(c *Conn, b []byte) (int, error) ***REMOVED***
	var operr error
	var n int
	fn := func(s uintptr) ***REMOVED***
		n, operr = getsockopt(s, o.Level, o.Name, b)
	***REMOVED***
	if err := c.c.Control(fn); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return n, os.NewSyscallError("getsockopt", operr)
***REMOVED***

func (o *Option) set(c *Conn, b []byte) error ***REMOVED***
	var operr error
	fn := func(s uintptr) ***REMOVED***
		operr = setsockopt(s, o.Level, o.Name, b)
	***REMOVED***
	if err := c.c.Control(fn); err != nil ***REMOVED***
		return err
	***REMOVED***
	return os.NewSyscallError("setsockopt", operr)
***REMOVED***
