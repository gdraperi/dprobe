// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package ipv6

import (
	"unsafe"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/socket"
)

func marshal2292HopLimit(b []byte, cm *ControlMessage) []byte ***REMOVED***
	m := socket.ControlMessage(b)
	m.MarshalHeader(iana.ProtocolIPv6, sysIPV6_2292HOPLIMIT, 4)
	if cm != nil ***REMOVED***
		socket.NativeEndian.PutUint32(m.Data(4), uint32(cm.HopLimit))
	***REMOVED***
	return m.Next(4)
***REMOVED***

func marshal2292PacketInfo(b []byte, cm *ControlMessage) []byte ***REMOVED***
	m := socket.ControlMessage(b)
	m.MarshalHeader(iana.ProtocolIPv6, sysIPV6_2292PKTINFO, sizeofInet6Pktinfo)
	if cm != nil ***REMOVED***
		pi := (*inet6Pktinfo)(unsafe.Pointer(&m.Data(sizeofInet6Pktinfo)[0]))
		if ip := cm.Src.To16(); ip != nil && ip.To4() == nil ***REMOVED***
			copy(pi.Addr[:], ip)
		***REMOVED***
		if cm.IfIndex > 0 ***REMOVED***
			pi.setIfindex(cm.IfIndex)
		***REMOVED***
	***REMOVED***
	return m.Next(sizeofInet6Pktinfo)
***REMOVED***

func marshal2292NextHop(b []byte, cm *ControlMessage) []byte ***REMOVED***
	m := socket.ControlMessage(b)
	m.MarshalHeader(iana.ProtocolIPv6, sysIPV6_2292NEXTHOP, sizeofSockaddrInet6)
	if cm != nil ***REMOVED***
		sa := (*sockaddrInet6)(unsafe.Pointer(&m.Data(sizeofSockaddrInet6)[0]))
		sa.setSockaddr(cm.NextHop, cm.IfIndex)
	***REMOVED***
	return m.Next(sizeofSockaddrInet6)
***REMOVED***
