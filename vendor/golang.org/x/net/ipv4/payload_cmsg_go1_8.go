// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.9
// +build !nacl,!plan9,!windows

package ipv4

import "net"

func (c *payloadHandler) readFrom(b []byte) (n int, cm *ControlMessage, src net.Addr, err error) ***REMOVED***
	c.rawOpt.RLock()
	oob := NewControlMessage(c.rawOpt.cflags)
	c.rawOpt.RUnlock()
	var nn int
	switch c := c.PacketConn.(type) ***REMOVED***
	case *net.UDPConn:
		if n, nn, _, src, err = c.ReadMsgUDP(b, oob); err != nil ***REMOVED***
			return 0, nil, nil, err
		***REMOVED***
	case *net.IPConn:
		nb := make([]byte, maxHeaderLen+len(b))
		if n, nn, _, src, err = c.ReadMsgIP(nb, oob); err != nil ***REMOVED***
			return 0, nil, nil, err
		***REMOVED***
		hdrlen := int(nb[0]&0x0f) << 2
		copy(b, nb[hdrlen:])
		n -= hdrlen
	default:
		return 0, nil, nil, &net.OpError***REMOVED***Op: "read", Net: c.LocalAddr().Network(), Source: c.LocalAddr(), Err: errInvalidConnType***REMOVED***
	***REMOVED***
	if nn > 0 ***REMOVED***
		cm = new(ControlMessage)
		if err = cm.Parse(oob[:nn]); err != nil ***REMOVED***
			return 0, nil, nil, &net.OpError***REMOVED***Op: "read", Net: c.PacketConn.LocalAddr().Network(), Source: c.PacketConn.LocalAddr(), Err: err***REMOVED***
		***REMOVED***
	***REMOVED***
	if cm != nil ***REMOVED***
		cm.Src = netAddrToIP4(src)
	***REMOVED***
	return
***REMOVED***

func (c *payloadHandler) writeTo(b []byte, cm *ControlMessage, dst net.Addr) (n int, err error) ***REMOVED***
	oob := cm.Marshal()
	if dst == nil ***REMOVED***
		return 0, &net.OpError***REMOVED***Op: "write", Net: c.PacketConn.LocalAddr().Network(), Source: c.PacketConn.LocalAddr(), Err: errMissingAddress***REMOVED***
	***REMOVED***
	switch c := c.PacketConn.(type) ***REMOVED***
	case *net.UDPConn:
		n, _, err = c.WriteMsgUDP(b, oob, dst.(*net.UDPAddr))
	case *net.IPConn:
		n, _, err = c.WriteMsgIP(b, oob, dst.(*net.IPAddr))
	default:
		return 0, &net.OpError***REMOVED***Op: "write", Net: c.LocalAddr().Network(), Source: c.LocalAddr(), Addr: opAddr(dst), Err: errInvalidConnType***REMOVED***
	***REMOVED***
	return
***REMOVED***
