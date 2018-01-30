// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd linux solaris

package ipv4

import (
	"net"
	"unsafe"

	"golang.org/x/net/internal/socket"
)

var freebsd32o64 bool

func (so *sockOpt) setGroupReq(c *socket.Conn, ifi *net.Interface, grp net.IP) error ***REMOVED***
	var gr groupReq
	if ifi != nil ***REMOVED***
		gr.Interface = uint32(ifi.Index)
	***REMOVED***
	gr.setGroup(grp)
	var b []byte
	if freebsd32o64 ***REMOVED***
		var d [sizeofGroupReq + 4]byte
		s := (*[sizeofGroupReq]byte)(unsafe.Pointer(&gr))
		copy(d[:4], s[:4])
		copy(d[8:], s[4:])
		b = d[:]
	***REMOVED*** else ***REMOVED***
		b = (*[sizeofGroupReq]byte)(unsafe.Pointer(&gr))[:sizeofGroupReq]
	***REMOVED***
	return so.Set(c, b)
***REMOVED***

func (so *sockOpt) setGroupSourceReq(c *socket.Conn, ifi *net.Interface, grp, src net.IP) error ***REMOVED***
	var gsr groupSourceReq
	if ifi != nil ***REMOVED***
		gsr.Interface = uint32(ifi.Index)
	***REMOVED***
	gsr.setSourceGroup(grp, src)
	var b []byte
	if freebsd32o64 ***REMOVED***
		var d [sizeofGroupSourceReq + 4]byte
		s := (*[sizeofGroupSourceReq]byte)(unsafe.Pointer(&gsr))
		copy(d[:4], s[:4])
		copy(d[8:], s[4:])
		b = d[:]
	***REMOVED*** else ***REMOVED***
		b = (*[sizeofGroupSourceReq]byte)(unsafe.Pointer(&gsr))[:sizeofGroupSourceReq]
	***REMOVED***
	return so.Set(c, b)
***REMOVED***
