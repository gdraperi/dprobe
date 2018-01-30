// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

import (
	"net"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/socket"
)

var (
	ctlOpts = [ctlMax]ctlOpt***REMOVED***
		ctlTrafficClass: ***REMOVED***sysIPV6_TCLASS, 4, marshalTrafficClass, parseTrafficClass***REMOVED***,
		ctlHopLimit:     ***REMOVED***sysIPV6_HOPLIMIT, 4, marshalHopLimit, parseHopLimit***REMOVED***,
		ctlPacketInfo:   ***REMOVED***sysIPV6_PKTINFO, sizeofInet6Pktinfo, marshalPacketInfo, parsePacketInfo***REMOVED***,
		ctlNextHop:      ***REMOVED***sysIPV6_NEXTHOP, sizeofSockaddrInet6, marshalNextHop, parseNextHop***REMOVED***,
		ctlPathMTU:      ***REMOVED***sysIPV6_PATHMTU, sizeofIPv6Mtuinfo, marshalPathMTU, parsePathMTU***REMOVED***,
	***REMOVED***

	sockOpts = map[int]sockOpt***REMOVED***
		ssoTrafficClass:        ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_TCLASS, Len: 4***REMOVED******REMOVED***,
		ssoHopLimit:            ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_UNICAST_HOPS, Len: 4***REMOVED******REMOVED***,
		ssoMulticastInterface:  ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_MULTICAST_IF, Len: 4***REMOVED******REMOVED***,
		ssoMulticastHopLimit:   ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_MULTICAST_HOPS, Len: 4***REMOVED******REMOVED***,
		ssoMulticastLoopback:   ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_MULTICAST_LOOP, Len: 4***REMOVED******REMOVED***,
		ssoReceiveTrafficClass: ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_RECVTCLASS, Len: 4***REMOVED******REMOVED***,
		ssoReceiveHopLimit:     ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_RECVHOPLIMIT, Len: 4***REMOVED******REMOVED***,
		ssoReceivePacketInfo:   ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_RECVPKTINFO, Len: 4***REMOVED******REMOVED***,
		ssoReceivePathMTU:      ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_RECVPATHMTU, Len: 4***REMOVED******REMOVED***,
		ssoPathMTU:             ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_PATHMTU, Len: sizeofIPv6Mtuinfo***REMOVED******REMOVED***,
		ssoChecksum:            ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysIPV6_CHECKSUM, Len: 4***REMOVED******REMOVED***,
		ssoICMPFilter:          ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6ICMP, Name: sysICMP6_FILTER, Len: sizeofICMPv6Filter***REMOVED******REMOVED***,
		ssoJoinGroup:           ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_JOIN_GROUP, Len: sizeofGroupReq***REMOVED***, typ: ssoTypeGroupReq***REMOVED***,
		ssoLeaveGroup:          ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_LEAVE_GROUP, Len: sizeofGroupReq***REMOVED***, typ: ssoTypeGroupReq***REMOVED***,
		ssoJoinSourceGroup:     ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_JOIN_SOURCE_GROUP, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
		ssoLeaveSourceGroup:    ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_LEAVE_SOURCE_GROUP, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
		ssoBlockSourceGroup:    ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_BLOCK_SOURCE, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
		ssoUnblockSourceGroup:  ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIPv6, Name: sysMCAST_UNBLOCK_SOURCE, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
	***REMOVED***
)

func init() ***REMOVED***
	if runtime.GOOS == "freebsd" && runtime.GOARCH == "386" ***REMOVED***
		archs, _ := syscall.Sysctl("kern.supported_archs")
		for _, s := range strings.Fields(archs) ***REMOVED***
			if s == "amd64" ***REMOVED***
				freebsd32o64 = true
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
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
	sa := (*sockaddrInet6)(unsafe.Pointer(&gr.Group))
	sa.Len = sizeofSockaddrInet6
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], grp)
***REMOVED***

func (gsr *groupSourceReq) setSourceGroup(grp, src net.IP) ***REMOVED***
	sa := (*sockaddrInet6)(unsafe.Pointer(&gsr.Group))
	sa.Len = sizeofSockaddrInet6
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], grp)
	sa = (*sockaddrInet6)(unsafe.Pointer(&gsr.Source))
	sa.Len = sizeofSockaddrInet6
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], src)
***REMOVED***
