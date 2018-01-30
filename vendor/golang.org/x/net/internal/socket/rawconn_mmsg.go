// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9
// +build linux

package socket

import (
	"net"
	"os"
	"syscall"
)

func (c *Conn) recvMsgs(ms []Message, flags int) (int, error) ***REMOVED***
	hs := make(mmsghdrs, len(ms))
	var parseFn func([]byte, string) (net.Addr, error)
	if c.network != "tcp" ***REMOVED***
		parseFn = parseInetAddr
	***REMOVED***
	if err := hs.pack(ms, parseFn, nil); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	var operr error
	var n int
	fn := func(s uintptr) bool ***REMOVED***
		n, operr = recvmmsg(s, hs, flags)
		if operr == syscall.EAGAIN ***REMOVED***
			return false
		***REMOVED***
		return true
	***REMOVED***
	if err := c.c.Read(fn); err != nil ***REMOVED***
		return n, err
	***REMOVED***
	if operr != nil ***REMOVED***
		return n, os.NewSyscallError("recvmmsg", operr)
	***REMOVED***
	if err := hs[:n].unpack(ms[:n], parseFn, c.network); err != nil ***REMOVED***
		return n, err
	***REMOVED***
	return n, nil
***REMOVED***

func (c *Conn) sendMsgs(ms []Message, flags int) (int, error) ***REMOVED***
	hs := make(mmsghdrs, len(ms))
	var marshalFn func(net.Addr) []byte
	if c.network != "tcp" ***REMOVED***
		marshalFn = marshalInetAddr
	***REMOVED***
	if err := hs.pack(ms, nil, marshalFn); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	var operr error
	var n int
	fn := func(s uintptr) bool ***REMOVED***
		n, operr = sendmmsg(s, hs, flags)
		if operr == syscall.EAGAIN ***REMOVED***
			return false
		***REMOVED***
		return true
	***REMOVED***
	if err := c.c.Write(fn); err != nil ***REMOVED***
		return n, err
	***REMOVED***
	if operr != nil ***REMOVED***
		return n, os.NewSyscallError("sendmmsg", operr)
	***REMOVED***
	if err := hs[:n].unpack(ms[:n], nil, ""); err != nil ***REMOVED***
		return n, err
	***REMOVED***
	return n, nil
***REMOVED***
