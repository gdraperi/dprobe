// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

import (
	"net"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/socket"
)

var (
	ctlOpts = [ctlMax]ctlOpt***REMOVED***
		ctlHopLimit:   ***REMOVED***sysIPV6_2292HOPLIMIT, 4, marshal2292HopLimit, parseHopLimit***REMOVED***,
		ctlPacketInfo: ***REMOVED***sysIPV6_2292PKTINFO, sizeofInet6Pktinfo, marshal2292PacketInfo, parsePacketInfo***REMOVED***,
	***REMOVED***

	sockOpts = map[int]*sockOpt***REMOVED***
		ssoHopLimit:           ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_UNICAST_HOPS, Len: 4***REMOVED******REMOVED***,
		ssoMulticastInterface: ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_MULTICAST_IF, Len: 4***REMOVED******REMOVED***,
		ssoMulticastHopLimit:  ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_MULTICAST_HOPS, Len: 4***REMOVED******REMOVED***,
		ssoMulticastLoopback:  ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_MULTICAST_LOOP, Len: 4***REMOVED******REMOVED***,
		ssoReceiveHopLimit:    ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_2292HOPLIMIT, Len: 4***REMOVED******REMOVED***,
		ssoReceivePacketInfo:  ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_2292PKTINFO, Len: 4***REMOVED******REMOVED***,
		ssoChecksum:           ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_CHECKSUM, Len: 4***REMOVED******REMOVED***,
		ssoICMPFilter:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6ICMP, Name: sysICMP6_FILTER, Len: sizeofICMPv6Filter***REMOVED******REMOVED***,
		ssoJoinGroup:          ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_JOIN_GROUP, Len: sizeofIPv6Mreq***REMOVED***, typ: ssoTypeIPMreq***REMOVED***,
		ssoLeaveGroup:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_LEAVE_GROUP, Len: sizeofIPv6Mreq***REMOVED***, typ: ssoTypeIPMreq***REMOVED***,
	***REMOVED***
)

func init() ***REMOVED***
	// Seems like kern.osreldate is veiled on latest OS X. We use
	// kern.osrelease instead.
	s, err := syscall.Sysctl("kern.osrelease")
	if err != nil ***REMOVED***
		return
	***REMOVED***
	ss := strings.Split(s, ".")
	if len(ss) == 0 ***REMOVED***
		return
	***REMOVED***
	// The IP_PKTINFO and protocol-independent multicast API were
	// introduced in OS X 10.7 (Darwin 11). But it looks like
	// those features require OS X 10.8 (Darwin 12) or above.
	// See http://support.apple.com/kb/HT1633.
	if mjver, err := strconv.Atoi(ss[0]); err != nil || mjver < 12 ***REMOVED***
		return
	***REMOVED***
	ctlOpts[ctlTrafficClass] = ctlOpt***REMOVED***sysIPV6_TCLASS, 4, marshalTrafficClass, parseTrafficClass***REMOVED***
	ctlOpts[ctlHopLimit] = ctlOpt***REMOVED***sysIPV6_HOPLIMIT, 4, marshalHopLimit, parseHopLimit***REMOVED***
	ctlOpts[ctlPacketInfo] = ctlOpt***REMOVED***sysIPV6_PKTINFO, sizeofInet6Pktinfo, marshalPacketInfo, parsePacketInfo***REMOVED***
	ctlOpts[ctlNextHop] = ctlOpt***REMOVED***sysIPV6_NEXTHOP, sizeofSockaddrInet6, marshalNextHop, parseNextHop***REMOVED***
	ctlOpts[ctlPathMTU] = ctlOpt***REMOVED***sysIPV6_PATHMTU, sizeofIPv6Mtuinfo, marshalPathMTU, parsePathMTU***REMOVED***
	sockOpts[ssoTrafficClass] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_TCLASS, Len: 4***REMOVED******REMOVED***
	sockOpts[ssoReceiveTrafficClass] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_RECVTCLASS, Len: 4***REMOVED******REMOVED***
	sockOpts[ssoReceiveHopLimit] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_RECVHOPLIMIT, Len: 4***REMOVED******REMOVED***
	sockOpts[ssoReceivePacketInfo] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_RECVPKTINFO, Len: 4***REMOVED******REMOVED***
	sockOpts[ssoReceivePathMTU] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_RECVPATHMTU, Len: 4***REMOVED******REMOVED***
	sockOpts[ssoPathMTU] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_PATHMTU, Len: sizeofIPv6Mtuinfo***REMOVED******REMOVED***
	sockOpts[ssoJoinGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_JOIN_GROUP, Len: sizeofGroupReq***REMOVED***, typ: ssoTypeGroupReq***REMOVED***
	sockOpts[ssoLeaveGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_LEAVE_GROUP, Len: sizeofGroupReq***REMOVED***, typ: ssoTypeGroupReq***REMOVED***
	sockOpts[ssoJoinSourceGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_JOIN_SOURCE_GROUP, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***
	sockOpts[ssoLeaveSourceGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_LEAVE_SOURCE_GROUP, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***
	sockOpts[ssoBlockSourceGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_BLOCK_SOURCE, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***
	sockOpts[ssoUnblockSourceGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_UNBLOCK_SOURCE, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***
***REMOVED***

func (sa *sockaddrInet6) setSockaddr(ip net.IP, i int) ***REMOVED***
	sa.Len = sizeofSockaddrInet6
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], ip)
	sa.Scope_id = uint32(i)
***REMOVED***

func (pi *inet6Pktinfo) setIfindex(i int) ***REMOVED***
	pi.Ifindex = uint32(i)
***REMOVED***

func (mreq *ipv6Mreq) setIfindex(i int) ***REMOVED***
	mreq.Interface = uint32(i)
***REMOVED***

func (gr *groupReq) setGroup(grp net.IP) ***REMOVED***
	sa := (*sockaddrInet6)(unsafe.Pointer(uintptr(unsafe.Pointer(gr)) + 4))
	sa.Len = sizeofSockaddrInet6
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], grp)
***REMOVED***

func (gsr *groupSourceReq) setSourceGroup(grp, src net.IP) ***REMOVED***
	sa := (*sockaddrInet6)(unsafe.Pointer(uintptr(unsafe.Pointer(gsr)) + 4))
	sa.Len = sizeofSockaddrInet6
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], grp)
	sa = (*sockaddrInet6)(unsafe.Pointer(uintptr(unsafe.Pointer(gsr)) + 132))
	sa.Len = sizeofSockaddrInet6
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], src)
***REMOVED***
