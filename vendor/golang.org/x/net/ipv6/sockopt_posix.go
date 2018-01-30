// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package ipv6

import (
	"net"
	"unsafe"

	"golang.org/x/net/bpf"
	"golang.org/x/net/internal/socket"
)

func (so *sockOpt) getMulticastInterface(c *socket.Conn) (*net.Interface, error) ***REMOVED***
	n, err := so.GetInt(c)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return net.InterfaceByIndex(n)
***REMOVED***

func (so *sockOpt) setMulticastInterface(c *socket.Conn, ifi *net.Interface) error ***REMOVED***
	var n int
	if ifi != nil ***REMOVED***
		n = ifi.Index
	***REMOVED***
	return so.SetInt(c, n)
***REMOVED***

func (so *sockOpt) getICMPFilter(c *socket.Conn) (*ICMPFilter, error) ***REMOVED***
	b := make([]byte, so.Len)
	n, err := so.Get(c, b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if n != sizeofICMPv6Filter ***REMOVED***
		return nil, errOpNoSupport
	***REMOVED***
	return (*ICMPFilter)(unsafe.Pointer(&b[0])), nil
***REMOVED***

func (so *sockOpt) setICMPFilter(c *socket.Conn, f *ICMPFilter) error ***REMOVED***
	b := (*[sizeofICMPv6Filter]byte)(unsafe.Pointer(f))[:sizeofICMPv6Filter]
	return so.Set(c, b)
***REMOVED***

func (so *sockOpt) getMTUInfo(c *socket.Conn) (*net.Interface, int, error) ***REMOVED***
	b := make([]byte, so.Len)
	n, err := so.Get(c, b)
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***
	if n != sizeofIPv6Mtuinfo ***REMOVED***
		return nil, 0, errOpNoSupport
	***REMOVED***
	mi := (*ipv6Mtuinfo)(unsafe.Pointer(&b[0]))
	if mi.Addr.Scope_id == 0 ***REMOVED***
		return nil, int(mi.Mtu), nil
	***REMOVED***
	ifi, err := net.InterfaceByIndex(int(mi.Addr.Scope_id))
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***
	return ifi, int(mi.Mtu), nil
***REMOVED***

func (so *sockOpt) setGroup(c *socket.Conn, ifi *net.Interface, grp net.IP) error ***REMOVED***
	switch so.typ ***REMOVED***
	case ssoTypeIPMreq:
		return so.setIPMreq(c, ifi, grp)
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
