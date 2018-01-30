// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9
// +build !nacl,!plan9,!windows

package ipv6

import (
	"net"

	"golang.org/x/net/internal/socket"
)

func (c *payloadHandler) readFrom(b []byte) (int, *ControlMessage, net.Addr, error) ***REMOVED***
	c.rawOpt.RLock()
	m := socket.Message***REMOVED***
		Buffers: [][]byte***REMOVED***b***REMOVED***,
		OOB:     NewControlMessage(c.rawOpt.cflags),
	***REMOVED***
	c.rawOpt.RUnlock()
	switch c.PacketConn.(type) ***REMOVED***
	case *net.UDPConn:
		if err := c.RecvMsg(&m, 0); err != nil ***REMOVED***
			return 0, nil, nil, &net.OpError***REMOVED***Op: "read", Net: c.PacketConn.LocalAddr().Network(), Source: c.PacketConn.LocalAddr(), Err: err***REMOVED***
		***REMOVED***
	case *net.IPConn:
		if err := c.RecvMsg(&m, 0); err != nil ***REMOVED***
			return 0, nil, nil, &net.OpError***REMOVED***Op: "read", Net: c.PacketConn.LocalAddr().Network(), Source: c.PacketConn.LocalAddr(), Err: err***REMOVED***
		***REMOVED***
	default:
		return 0, nil, nil, &net.OpError***REMOVED***Op: "read", Net: c.PacketConn.LocalAddr().Network(), Source: c.PacketConn.LocalAddr(), Err: errInvalidConnType***REMOVED***
	***REMOVED***
	var cm *ControlMessage
	if m.NN > 0 ***REMOVED***
		cm = new(ControlMessage)
		if err := cm.Parse(m.OOB[:m.NN]); err != nil ***REMOVED***
			return 0, nil, nil, &net.OpError***REMOVED***Op: "read", Net: c.PacketConn.LocalAddr().Network(), Source: c.PacketConn.LocalAddr(), Err: err***REMOVED***
		***REMOVED***
		cm.Src = netAddrToIP16(m.Addr)
	***REMOVED***
	return m.N, cm, m.Addr, nil
***REMOVED***

func (c *payloadHandler) writeTo(b []byte, cm *ControlMessage, dst net.Addr) (int, error) ***REMOVED***
	m := socket.Message***REMOVED***
		Buffers: [][]byte***REMOVED***b***REMOVED***,
		OOB:     cm.Marshal(),
		Addr:    dst,
	***REMOVED***
	err := c.SendMsg(&m, 0)
	if err != nil ***REMOVED***
		err = &net.OpError***REMOVED***Op: "write", Net: c.PacketConn.LocalAddr().Network(), Source: c.PacketConn.LocalAddr(), Addr: opAddr(dst), Err: err***REMOVED***
	***REMOVED***
	return m.N, err
***REMOVED***
