// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

import (
	"net"
	"syscall"
	"unsafe"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/socket"
)

var (
	ctlOpts = [ctlMax]ctlOpt***REMOVED***
		ctlTTL:        ***REMOVED***sysIP_TTL, 1, marshalTTL, parseTTL***REMOVED***,
		ctlPacketInfo: ***REMOVED***sysIP_PKTINFO, sizeofInetPktinfo, marshalPacketInfo, parsePacketInfo***REMOVED***,
	***REMOVED***

	sockOpts = map[int]*sockOpt***REMOVED***
		ssoTOS:                ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_TOS, Len: 4***REMOVED******REMOVED***,
		ssoTTL:                ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_TTL, Len: 4***REMOVED******REMOVED***,
		ssoMulticastTTL:       ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_MULTICAST_TTL, Len: 4***REMOVED******REMOVED***,
		ssoMulticastInterface: ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_MULTICAST_IF, Len: sizeofIPMreqn***REMOVED***, typ: ssoTypeIPMreqn***REMOVED***,
		ssoMulticastLoopback:  ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_MULTICAST_LOOP, Len: 4***REMOVED******REMOVED***,
		ssoReceiveTTL:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_RECVTTL, Len: 4***REMOVED******REMOVED***,
		ssoPacketInfo:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_PKTINFO, Len: 4***REMOVED******REMOVED***,
		ssoHeaderPrepend:      ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_HDRINCL, Len: 4***REMOVED******REMOVED***,
		ssoICMPFilter:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolReserved, Name: sysICMP_FILTER, Len: sizeofICMPFilter***REMOVED******REMOVED***,
		ssoJoinGroup:          ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_JOIN_GROUP, Len: sizeofGroupReq***REMOVED***, typ: ssoTypeGroupReq***REMOVED***,
		ssoLeaveGroup:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_LEAVE_GROUP, Len: sizeofGroupReq***REMOVED***, typ: ssoTypeGroupReq***REMOVED***,
		ssoJoinSourceGroup:    ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_JOIN_SOURCE_GROUP, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
		ssoLeaveSourceGroup:   ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_LEAVE_SOURCE_GROUP, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
		ssoBlockSourceGroup:   ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_BLOCK_SOURCE, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
		ssoUnblockSourceGroup: ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_UNBLOCK_SOURCE, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
		ssoAttachFilter:       ***REMOVED***Option: socket.Option***REMOVED***Level: sysSOL_SOCKET, Name: sysSO_ATTACH_FILTER, Len: sizeofSockFprog***REMOVED******REMOVED***,
	***REMOVED***
)

func (pi *inetPktinfo) setIfindex(i int) ***REMOVED***
	pi.Ifindex = int32(i)
***REMOVED***

func (gr *groupReq) setGroup(grp net.IP) ***REMOVED***
	sa := (*sockaddrInet)(unsafe.Pointer(&gr.Group))
	sa.Family = syscall.AF_INET
	copy(sa.Addr[:], grp)
***REMOVED***

func (gsr *groupSourceReq) setSourceGroup(grp, src net.IP) ***REMOVED***
	sa := (*sockaddrInet)(unsafe.Pointer(&gsr.Group))
	sa.Family = syscall.AF_INET
	copy(sa.Addr[:], grp)
	sa = (*sockaddrInet)(unsafe.Pointer(&gsr.Source))
	sa.Family = syscall.AF_INET
	copy(sa.Addr[:], src)
***REMOVED***
