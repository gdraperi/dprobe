// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package icmp

import (
	"net"
	"strconv"
	"syscall"
)

func sockaddr(family int, address string) (syscall.Sockaddr, error) ***REMOVED***
	switch family ***REMOVED***
	case syscall.AF_INET:
		a, err := net.ResolveIPAddr("ip4", address)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if len(a.IP) == 0 ***REMOVED***
			a.IP = net.IPv4zero
		***REMOVED***
		if a.IP = a.IP.To4(); a.IP == nil ***REMOVED***
			return nil, net.InvalidAddrError("non-ipv4 address")
		***REMOVED***
		sa := &syscall.SockaddrInet4***REMOVED******REMOVED***
		copy(sa.Addr[:], a.IP)
		return sa, nil
	case syscall.AF_INET6:
		a, err := net.ResolveIPAddr("ip6", address)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if len(a.IP) == 0 ***REMOVED***
			a.IP = net.IPv6unspecified
		***REMOVED***
		if a.IP.Equal(net.IPv4zero) ***REMOVED***
			a.IP = net.IPv6unspecified
		***REMOVED***
		if a.IP = a.IP.To16(); a.IP == nil || a.IP.To4() != nil ***REMOVED***
			return nil, net.InvalidAddrError("non-ipv6 address")
		***REMOVED***
		sa := &syscall.SockaddrInet6***REMOVED***ZoneId: zoneToUint32(a.Zone)***REMOVED***
		copy(sa.Addr[:], a.IP)
		return sa, nil
	default:
		return nil, net.InvalidAddrError("unexpected family")
	***REMOVED***
***REMOVED***

func zoneToUint32(zone string) uint32 ***REMOVED***
	if zone == "" ***REMOVED***
		return 0
	***REMOVED***
	if ifi, err := net.InterfaceByName(zone); err == nil ***REMOVED***
		return uint32(ifi.Index)
	***REMOVED***
	n, err := strconv.Atoi(zone)
	if err != nil ***REMOVED***
		return 0
	***REMOVED***
	return uint32(n)
***REMOVED***

func last(s string, b byte) int ***REMOVED***
	i := len(s)
	for i--; i >= 0; i-- ***REMOVED***
		if s[i] == b ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return i
***REMOVED***
