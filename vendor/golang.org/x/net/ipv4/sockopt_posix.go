// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package ipv4

import (
	"net"
	"unsafe"

	"golang.org/x/net/bpf"
	"golang.org/x/net/internal/socket"
)

func (so *sockOpt) getMulticastInterface(c *socket.Conn) (*net.Interface, error) ***REMOVED***
	switch so.typ ***REMOVED***
	case ssoTypeIPMreqn:
		return so.getIPMreqn(c)
	default:
		return so.getMulticastIf(c)
	***REMOVED***
***REMOVED***

func (so *sockOpt) setMulticastInterface(c *socket.Conn, ifi *net.Interface) error ***REMOVED***
	switch so.typ ***REMOVED***
	case ssoTypeIPMreqn:
		return so.setIPMreqn(c, ifi, nil)
	default:
		return so.setMulticastIf(c, ifi)
	***REMOVED***
***REMOVED***

func (so *sockOpt) getICMPFilter(c *socket.Conn) (*ICMPFilter, error) ***REMOVED***
	b := make([]byte, so.Len)
	n, err := so.Get(c, b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if n != sizeofICMPFilter ***REMOVED***
		return nil, errOpNoSupport
	***REMOVED***
	return (*ICMPFilter)(unsafe.Pointer(&b[0])), nil
***REMOVED***

func (so *sockOpt) setICMPFilter(c *socket.Conn, f *ICMPFilter) error ***REMOVED***
	b := (*[sizeofICMPFilter]byte)(unsafe.Pointer(f))[:sizeofICMPFilter]
	return so.Set(c, b)
***REMOVED***

func (so *sockOpt) setGroup(c *socket.Conn, ifi *net.Interface, grp net.IP) error ***REMOVED***
	switch so.typ ***REMOVED***
	case ssoTypeIPMreq:
		return so.setIPMreq(c, ifi, grp)
	case ssoTypeIPMreqn:
		return so.setIPMreqn(c, ifi, grp)
	case ssoTypeGroupReq:
		return so.setGroupReq(c, ifi, grp)
	default:
		return errOpNoSupport
	***REMOVED***
***REMOVED***

func (so *sockOpt) setSourceGroup(c *socket.Conn, ifi *net.Interface, grp, src net.IP) error ***REMOVED***
	return so.setGroupSourceReq(c, ifi, grp, src)
***REMOVED***

func (so *sockOpt) setBPF(c *socket.Conn, f []bpf.RawInstruction) error ***REMOVED***
	return so.setAttachFilter(c, f)
***REMOVED***
