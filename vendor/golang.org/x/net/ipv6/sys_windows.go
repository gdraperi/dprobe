// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

import (
	"net"
	"syscall"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/socket"
)

const (
	// See ws2tcpip.h.
	sysIPV6_UNICAST_HOPS   = 0x4
	sysIPV6_MULTICAST_IF   = 0x9
	sysIPV6_MULTICAST_HOPS = 0xa
	sysIPV6_MULTICAST_LOOP = 0xb
	sysIPV6_JOIN_GROUP     = 0xc
	sysIPV6_LEAVE_GROUP    = 0xd
	sysIPV6_PKTINFO        = 0x13

	sizeofSockaddrInet6 = 0x1c

	sizeofIPv6Mreq     = 0x14
	sizeofIPv6Mtuinfo  = 0x20
	sizeofICMPv6Filter = 0
)

type sockaddrInet6 struct ***REMOVED***
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
***REMOVED***

type ipv6Mreq struct ***REMOVED***
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
***REMOVED***

type ipv6Mtuinfo struct ***REMOVED***
	Addr sockaddrInet6
	Mtu  uint32
***REMOVED***

type icmpv6Filter struct ***REMOVED***
	// TODO(mikio): implement this
***REMOVED***

var (
	ctlOpts = [ctlMax]ctlOpt***REMOVED******REMOVED***

	sockOpts = map[int]*sockOpt***REMOVED***
		ssoHopLimit:           ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_UNICAST_HOPS, Len: 4***REMOVED******REMOVED***,
		ssoMulticastInterface: ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_MULTICAST_IF, Len: 4***REMOVED******REMOVED***,
		ssoMulticastHopLimit:  ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_MULTICAST_HOPS, Len: 4***REMOVED******REMOVED***,
		ssoMulticastLoopback:  ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_MULTICAST_LOOP, Len: 4***REMOVED******REMOVED***,
		ssoJoinGroup:          ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_JOIN_GROUP, Len: sizeofIPv6Mreq***REMOVED***, typ: ssoTypeIPMreq***REMOVED***,
		ssoLeaveGroup:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_LEAVE_GROUP, Len: sizeofIPv6Mreq***REMOVED***, typ: ssoTypeIPMreq***REMOVED***,
	***REMOVED***
)

func (sa *sockaddrInet6) setSockaddr(ip net.IP, i int) ***REMOVED***
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], ip)
	sa.Scope_id = uint32(i)
***REMOVED***

func (mreq *ipv6Mreq) setIfindex(i int) ***REMOVED***
	mreq.Interface = uint32(i)
***REMOVED***
