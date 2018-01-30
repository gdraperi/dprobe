// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9
// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package socket

import (
	"os"
	"syscall"
)

func (c *Conn) recvMsg(m *Message, flags int) error ***REMOVED***
	var h msghdr
	vs := make([]iovec, len(m.Buffers))
	var sa []byte
	if c.network != "tcp" ***REMOVED***
		sa = make([]byte, sizeofSockaddrInet6)
	***REMOVED***
	h.pack(vs, m.Buffers, m.OOB, sa)
	var operr error
	var n int
	fn := func(s uintptr) bool ***REMOVED***
		n, operr = recvmsg(s, &h, flags)
		if operr == syscall.EAGAIN ***REMOVED***
			return false
		***REMOVED***
		return true
	***REMOVED***
	if err := c.c.Read(fn); err != nil ***REMOVED***
		return err
	***REMOVED***
	if operr != nil ***REMOVED***
		return os.NewSyscallError("recvmsg", operr)
	***REMOVED***
	if c.network != "tcp" ***REMOVED***
		var err error
		m.Addr, err = parseInetAddr(sa[:], c.network)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	m.N = n
	m.NN = h.controllen()
	m.Flags = h.flags()
	return nil
***REMOVED***

func (c *Conn) sendMsg(m *Message, flags int) error ***REMOVED***
	var h msghdr
	vs := make([]iovec, len(m.Buffers))
	var sa []byte
	if m.Addr != nil ***REMOVED***
		sa = marshalInetAddr(m.Addr)
	***REMOVED***
	h.pack(vs, m.Buffers, m.OOB, sa)
	var operr error
	var n int
	fn := func(s uintptr) bool ***REMOVED***
		n, operr = sendmsg(s, &h, flags)
		if operr == syscall.EAGAIN ***REMOVED***
			return false
		***REMOVED***
		return true
	***REMOVED***
	if err := c.c.Write(fn); err != nil ***REMOVED***
		return err
	***REMOVED***
	if operr != nil ***REMOVED***
		return os.NewSyscallError("sendmsg", operr)
	***REMOVED***
	m.N = n
	m.NN = len(m.OOB)
	return nil
***REMOVED***
