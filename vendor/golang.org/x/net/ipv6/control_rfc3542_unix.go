// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package ipv6

import (
	"net"
	"unsafe"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/socket"
)

func marshalTrafficClass(b []byte, cm *ControlMessage) []byte ***REMOVED***
	m := socket.ControlMessage(b)
	m.MarshalHeader(iana.ProtocolIPv6, sysIPV6_TCLASS, 4)
	if cm != nil ***REMOVED***
		socket.NativeEndian.PutUint32(m.Data(4), uint32(cm.TrafficClass))
	***REMOVED***
	return m.Next(4)
***REMOVED***

func parseTrafficClass(cm *ControlMessage, b []byte) ***REMOVED***
	cm.TrafficClass = int(socket.NativeEndian.Uint32(b[:4]))
***REMOVED***

func marshalHopLimit(b []byte, cm *ControlMessage) []byte ***REMOVED***
	m := socket.ControlMessage(b)
	m.MarshalHeader(iana.ProtocolIPv6, sysIPV6_HOPLIMIT, 4)
	if cm != nil ***REMOVED***
		socket.NativeEndian.PutUint32(m.Data(4), uint32(cm.HopLimit))
	***REMOVED***
	return m.Next(4)
***REMOVED***

func parseHopLimit(cm *ControlMessage, b []byte) ***REMOVED***
	cm.HopLimit = int(socket.NativeEndian.Uint32(b[:4]))
***REMOVED***

func marshalPacketInfo(b []byte, cm *ControlMessage) []byte ***REMOVED***
	m := socket.ControlMessage(b)
	m.MarshalHeader(iana.ProtocolIPv6, sysIPV6_PKTINFO, sizeofInet6Pktinfo)
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

func parsePacketInfo(cm *ControlMessage, b []byte) ***REMOVED***
	pi := (*inet6Pktinfo)(unsafe.Pointer(&b[0]))
	if len(cm.Dst) < net.IPv6len ***REMOVED***
		cm.Dst = make(net.IP, net.IPv6len)
	***REMOVED***
	copy(cm.Dst, pi.Addr[:])
	cm.IfIndex = int(pi.Ifindex)
***REMOVED***

func marshalNextHop(b []byte, cm *ControlMessage) []byte ***REMOVED***
	m := socket.ControlMessage(b)
	m.MarshalHeader(iana.ProtocolIPv6, sysIPV6_NEXTHOP, sizeofSockaddrInet6)
	if cm != nil ***REMOVED***
		sa := (*sockaddrInet6)(unsafe.Pointer(&m.Data(sizeofSockaddrInet6)[0]))
		sa.setSockaddr(cm.NextHop, cm.IfIndex)
	***REMOVED***
	return m.Next(sizeofSockaddrInet6)
***REMOVED***

func parseNextHop(cm *ControlMessage, b []byte) ***REMOVED***
***REMOVED***

func marshalPathMTU(b []byte, cm *ControlMessage) []byte ***REMOVED***
	m := socket.ControlMessage(b)
	m.MarshalHeader(iana.ProtocolIPv6, sysIPV6_PATHMTU, sizeofIPv6Mtuinfo)
	return m.Next(sizeofIPv6Mtuinfo)
***REMOVED***

func parsePathMTU(cm *ControlMessage, b []byte) ***REMOVED***
	mi := (*ipv6Mtuinfo)(unsafe.Pointer(&b[0]))
	if len(cm.Dst) < net.IPv6len ***REMOVED***
		cm.Dst = make(net.IP, net.IPv6len)
	***REMOVED***
	copy(cm.Dst, mi.Addr.Addr[:])
	cm.IfIndex = int(mi.Addr.Scope_id)
	cm.MTU = int(mi.Mtu)
***REMOVED***
