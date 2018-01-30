// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9

package ipv4

import (
	"net"

	"golang.org/x/net/internal/socket"
)

func (c *packetHandler) readFrom(b []byte) (h *Header, p []byte, cm *ControlMessage, err error) ***REMOVED***
	c.rawOpt.RLock()
	m := socket.Message***REMOVED***
		Buffers: [][]byte***REMOVED***b***REMOVED***,
		OOB:     NewControlMessage(c.rawOpt.cflags),
	***REMOVED***
	c.rawOpt.RUnlock()
	if err := c.RecvMsg(&m, 0); err != nil ***REMOVED***
		return nil, nil, nil, &net.OpError***REMOVED***Op: "read", Net: c.IPConn.LocalAddr().Network(), Source: c.IPConn.LocalAddr(), Err: err***REMOVED***
	***REMOVED***
	var hs []byte
	if hs, p, err = slicePacket(b[:m.N]); err != nil ***REMOVED***
		return nil, nil, nil, &net.OpError***REMOVED***Op: "read", Net: c.IPConn.LocalAddr().Network(), Source: c.IPConn.LocalAddr(), Err: err***REMOVED***
	***REMOVED***
	if h, err = ParseHeader(hs); err != nil ***REMOVED***
		return nil, nil, nil, &net.OpError***REMOVED***Op: "read", Net: c.IPConn.LocalAddr().Network(), Source: c.IPConn.LocalAddr(), Err: err***REMOVED***
	***REMOVED***
	if m.NN > 0 ***REMOVED***
		cm = new(ControlMessage)
		if err := cm.Parse(m.OOB[:m.NN]); err != nil ***REMOVED***
			return nil, nil, nil, &net.OpError***REMOVED***Op: "read", Net: c.IPConn.LocalAddr().Network(), Source: c.IPConn.LocalAddr(), Err: err***REMOVED***
		***REMOVED***
	***REMOVED***
	if src, ok := m.Addr.(*net.IPAddr); ok && cm != nil ***REMOVED***
		cm.Src = src.IP
	***REMOVED***
	return
***REMOVED***

func (c *packetHandler) writeTo(h *Header, p []byte, cm *ControlMessage) error ***REMOVED***
	m := socket.Message***REMOVED***
		OOB: cm.Marshal(),
	***REMOVED***
	wh, err := h.Marshal()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	m.Buffers = [][]byte***REMOVED***wh, p***REMOVED***
	dst := new(net.IPAddr)
	if cm != nil ***REMOVED***
		if ip := cm.Dst.To4(); ip != nil ***REMOVED***
			dst.IP = ip
		***REMOVED***
	***REMOVED***
	if dst.IP == nil ***REMOVED***
		dst.IP = h.Dst
	***REMOVED***
	m.Addr = dst
	if err := c.SendMsg(&m, 0); err != nil ***REMOVED***
		return &net.OpError***REMOVED***Op: "write", Net: c.IPConn.LocalAddr().Network(), Source: c.IPConn.LocalAddr(), Addr: opAddr(dst), Err: err***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
