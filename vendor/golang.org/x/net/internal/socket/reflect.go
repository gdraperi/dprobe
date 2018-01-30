// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.9

package socket

import (
	"errors"
	"net"
	"os"
	"reflect"
	"runtime"
)

// A Conn represents a raw connection.
type Conn struct ***REMOVED***
	c net.Conn
***REMOVED***

// NewConn returns a new raw connection.
func NewConn(c net.Conn) (*Conn, error) ***REMOVED***
	return &Conn***REMOVED***c: c***REMOVED***, nil
***REMOVED***

func (o *Option) get(c *Conn, b []byte) (int, error) ***REMOVED***
	s, err := socketOf(c.c)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	n, err := getsockopt(s, o.Level, o.Name, b)
	return n, os.NewSyscallError("getsockopt", err)
***REMOVED***

func (o *Option) set(c *Conn, b []byte) error ***REMOVED***
	s, err := socketOf(c.c)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return os.NewSyscallError("setsockopt", setsockopt(s, o.Level, o.Name, b))
***REMOVED***

func socketOf(c net.Conn) (uintptr, error) ***REMOVED***
	switch c.(type) ***REMOVED***
	case *net.TCPConn, *net.UDPConn, *net.IPConn:
		v := reflect.ValueOf(c)
		switch e := v.Elem(); e.Kind() ***REMOVED***
		case reflect.Struct:
			fd := e.FieldByName("conn").FieldByName("fd")
			switch e := fd.Elem(); e.Kind() ***REMOVED***
			case reflect.Struct:
				sysfd := e.FieldByName("sysfd")
				if runtime.GOOS == "windows" ***REMOVED***
					return uintptr(sysfd.Uint()), nil
				***REMOVED***
				return uintptr(sysfd.Int()), nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return 0, errors.New("invalid type")
***REMOVED***
