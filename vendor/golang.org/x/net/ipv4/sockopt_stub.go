// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package ipv4

import (
	"net"

	"golang.org/x/net/bpf"
	"golang.org/x/net/internal/socket"
)

func (so *sockOpt) getMulticastInterface(c *socket.Conn) (*net.Interface, error) ***REMOVED***
	return nil, errOpNoSupport
***REMOVED***

func (so *sockOpt) setMulticastInterface(c *socket.Conn, ifi *net.Interface) error ***REMOVED***
	return errOpNoSupport
***REMOVED***

func (so *sockOpt) getICMPFilter(c *socket.Conn) (*ICMPFilter, error) ***REMOVED***
	return nil, errOpNoSupport
***REMOVED***

func (so *sockOpt) setICMPFilter(c *socket.Conn, f *ICMPFilter) error ***REMOVED***
	return errOpNoSupport
***REMOVED***

func (so *sockOpt) setGroup(c *socket.Conn, ifi *net.Interface, grp net.IP) error ***REMOVED***
	return errOpNoSupport
***REMOVED***

func (so *sockOpt) setSourceGroup(c *socket.Conn, ifi *net.Interface, grp, src net.IP) error ***REMOVED***
	return errOpNoSupport
***REMOVED***

func (so *sockOpt) setBPF(c *socket.Conn, f []bpf.RawInstruction) error ***REMOVED***
	return errOpNoSupport
***REMOVED***
