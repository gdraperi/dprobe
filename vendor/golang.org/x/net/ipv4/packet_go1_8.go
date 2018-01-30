// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.9

package ipv4

import "net"

func (c *packetHandler) readFrom(b []byte) (h *Header, p []byte, cm *ControlMessage, err error) ***REMOVED***
	c.rawOpt.RLock()
	oob := NewControlMessage(c.rawOpt.cflags)
	c.rawOpt.RUnlock()
	n, nn, _, src, err := c.ReadMsgIP(b, oob)
	if err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***
	var hs []byte
	if hs, p, err = slicePacket(b[:n]); err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***
	if h, err = ParseHeader(hs); err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***
	if nn > 0 ***REMOVED***
		cm = new(ControlMessage)
		if err := cm.Parse(oob[:nn]); err != nil ***REMOVED***
			return nil, nil, nil, err
		***REMOVED***
	***REMOVED***
	if src != nil && cm != nil ***REMOVED***
		cm.Src = src.IP
	***REMOVED***
	return
***REMOVED***

func (c *packetHandler) writeTo(h *Header, p []byte, cm *ControlMessage) error ***REMOVED***
	oob := cm.Marshal()
	wh, err := h.Marshal()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	dst := new(net.IPAddr)
	if cm != nil ***REMOVED***
		if ip := cm.Dst.To4(); ip != nil ***REMOVED***
			dst.IP = ip
		***REMOVED***
	***REMOVED***
	if dst.IP == nil ***REMOVED***
		dst.IP = h.Dst
	***REMOVED***
	wh = append(wh, p...)
	_, _, err = c.WriteMsgIP(wh, oob, dst)
	return err
***REMOVED***
