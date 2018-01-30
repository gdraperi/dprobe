// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build nacl plan9 windows

package ipv4

import (
	"net"
	"syscall"
)

// ReadFrom reads a payload of the received IPv4 datagram, from the
// endpoint c, copying the payload into b. It returns the number of
// bytes copied into b, the control message cm and the source address
// src of the received datagram.
func (c *payloadHandler) ReadFrom(b []byte) (n int, cm *ControlMessage, src net.Addr, err error) ***REMOVED***
	if !c.ok() ***REMOVED***
		return 0, nil, nil, syscall.EINVAL
	***REMOVED***
	if n, src, err = c.PacketConn.ReadFrom(b); err != nil ***REMOVED***
		return 0, nil, nil, err
	***REMOVED***
	return
***REMOVED***

// WriteTo writes a payload of the IPv4 datagram, to the destination
// address dst through the endpoint c, copying the payload from b. It
// returns the number of bytes written. The control message cm allows
// the datagram path and the outgoing interface to be specified.
// Currently only Darwin and Linux support this. The cm may be nil if
// control of the outgoing datagram is not required.
func (c *payloadHandler) WriteTo(b []byte, cm *ControlMessage, dst net.Addr) (n int, err error) ***REMOVED***
	if !c.ok() ***REMOVED***
		return 0, syscall.EINVAL
	***REMOVED***
	if dst == nil ***REMOVED***
		return 0, errMissingAddress
	***REMOVED***
	return c.PacketConn.WriteTo(b, dst)
***REMOVED***
