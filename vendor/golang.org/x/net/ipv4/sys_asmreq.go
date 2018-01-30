// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd solaris windows

package ipv4

import (
	"net"
	"unsafe"

	"golang.org/x/net/internal/socket"
)

func (so *sockOpt) setIPMreq(c *socket.Conn, ifi *net.Interface, grp net.IP) error ***REMOVED***
	mreq := ipMreq***REMOVED***Multiaddr: [4]byte***REMOVED***grp[0], grp[1], grp[2], grp[3]***REMOVED******REMOVED***
	if err := setIPMreqInterface(&mreq, ifi); err != nil ***REMOVED***
		return err
	***REMOVED***
	b := (*[sizeofIPMreq]byte)(unsafe.Pointer(&mreq))[:sizeofIPMreq]
	return so.Set(c, b)
***REMOVED***

func (so *sockOpt) getMulticastIf(c *socket.Conn) (*net.Interface, error) ***REMOVED***
	var b [4]byte
	if _, err := so.Get(c, b[:]); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ifi, err := netIP4ToInterface(net.IPv4(b[0], b[1], b[2], b[3]))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return ifi, nil
***REMOVED***

func (so *sockOpt) setMulticastIf(c *socket.Conn, ifi *net.Interface) error ***REMOVED***
	ip, err := netInterfaceToIP4(ifi)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	var b [4]byte
	copy(b[:], ip)
	return so.Set(c, b[:])
***REMOVED***

func setIPMreqInterface(mreq *ipMreq, ifi *net.Interface) error ***REMOVED***
	if ifi == nil ***REMOVED***
		return nil
	***REMOVED***
	ifat, err := ifi.Addrs()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, ifa := range ifat ***REMOVED***
		switch ifa := ifa.(type) ***REMOVED***
		case *net.IPAddr:
			if ip := ifa.IP.To4(); ip != nil ***REMOVED***
				copy(mreq.Interface[:], ip)
				return nil
			***REMOVED***
		case *net.IPNet:
			if ip := ifa.IP.To4(); ip != nil ***REMOVED***
				copy(mreq.Interface[:], ip)
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return errNoSuchInterface
***REMOVED***

func netIP4ToInterface(ip net.IP) (*net.Interface, error) ***REMOVED***
	ift, err := net.Interfaces()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, ifi := range ift ***REMOVED***
		ifat, err := ifi.Addrs()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		for _, ifa := range ifat ***REMOVED***
			switch ifa := ifa.(type) ***REMOVED***
			case *net.IPAddr:
				if ip.Equal(ifa.IP) ***REMOVED***
					return &ifi, nil
				***REMOVED***
			case *net.IPNet:
				if ip.Equal(ifa.IP) ***REMOVED***
					return &ifi, nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil, errNoSuchInterface
***REMOVED***

func netInterfaceToIP4(ifi *net.Interface) (net.IP, error) ***REMOVED***
	if ifi == nil ***REMOVED***
		return net.IPv4zero.To4(), nil
	***REMOVED***
	ifat, err := ifi.Addrs()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, ifa := range ifat ***REMOVED***
		switch ifa := ifa.(type) ***REMOVED***
		case *net.IPAddr:
			if ip := ifa.IP.To4(); ip != nil ***REMOVED***
				return ip, nil
			***REMOVED***
		case *net.IPNet:
			if ip := ifa.IP.To4(); ip != nil ***REMOVED***
				return ip, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil, errNoSuchInterface
***REMOVED***
