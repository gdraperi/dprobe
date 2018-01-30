// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

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
		ctlTTL:       ***REMOVED***sysIP_RECVTTL, 1, marshalTTL, parseTTL***REMOVED***,
		ctlDst:       ***REMOVED***sysIP_RECVDSTADDR, net.IPv4len, marshalDst, parseDst***REMOVED***,
		ctlInterface: ***REMOVED***sysIP_RECVIF, syscall.SizeofSockaddrDatalink, marshalInterface, parseInterface***REMOVED***,
	***REMOVED***

	sockOpts = map[int]*sockOpt***REMOVED***
		ssoTOS:                ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_TOS, Len: 4***REMOVED******REMOVED***,
		ssoTTL:                ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_TTL, Len: 4***REMOVED******REMOVED***,
		ssoMulticastTTL:       ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_MULTICAST_TTL, Len: 1***REMOVED******REMOVED***,
		ssoMulticastInterface: ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_MULTICAST_IF, Len: 4***REMOVED******REMOVED***,
		ssoMulticastLoopback:  ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_MULTICAST_LOOP, Len: 4***REMOVED******REMOVED***,
		ssoReceiveTTL:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_RECVTTL, Len: 4***REMOVED******REMOVED***,
		ssoReceiveDst:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_RECVDSTADDR, Len: 4***REMOVED******REMOVED***,
		ssoReceiveInterface:   ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_RECVIF, Len: 4***REMOVED******REMOVED***,
		ssoHeaderPrepend:      ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_HDRINCL, Len: 4***REMOVED******REMOVED***,
		ssoStripHeader:        ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_STRIPHDR, Len: 4***REMOVED******REMOVED***,
		ssoJoinGroup:          ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_ADD_MEMBERSHIP, Len: sizeofIPMreq***REMOVED***, typ: ssoTypeIPMreq***REMOVED***,
		ssoLeaveGroup:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_DROP_MEMBERSHIP, Len: sizeofIPMreq***REMOVED***, typ: ssoTypeIPMreq***REMOVED***,
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
	ctlOpts[ctlPacketInfo].name = sysIP_PKTINFO
	ctlOpts[ctlPacketInfo].length = sizeofInetPktinfo
	ctlOpts[ctlPacketInfo].marshal = marshalPacketInfo
	ctlOpts[ctlPacketInfo].parse = parsePacketInfo
	sockOpts[ssoPacketInfo] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_RECVPKTINFO, Len: 4***REMOVED******REMOVED***
	sockOpts[ssoMulticastInterface] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_MULTICAST_IF, Len: sizeofIPMreqn***REMOVED***, typ: ssoTypeIPMreqn***REMOVED***
	sockOpts[ssoJoinGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_JOIN_GROUP, Len: sizeofGroupReq***REMOVED***, typ: ssoTypeGroupReq***REMOVED***
	sockOpts[ssoLeaveGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_LEAVE_GROUP, Len: sizeofGroupReq***REMOVED***, typ: ssoTypeGroupReq***REMOVED***
	sockOpts[ssoJoinSourceGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_JOIN_SOURCE_GROUP, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***
	sockOpts[ssoLeaveSourceGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_LEAVE_SOURCE_GROUP, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***
	sockOpts[ssoBlockSourceGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_BLOCK_SOURCE, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***
	sockOpts[ssoUnblockSourceGroup] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_UNBLOCK_SOURCE, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***
***REMOVED***

func (pi *inetPktinfo) setIfindex(i int) ***REMOVED***
	pi.Ifindex = uint32(i)
***REMOVED***

func (gr *groupReq) setGroup(grp net.IP) ***REMOVED***
	sa := (*sockaddrInet)(unsafe.Pointer(uintptr(unsafe.Pointer(gr)) + 4))
	sa.Len = sizeofSockaddrInet
	sa.Family = syscall.AF_INET
	copy(sa.Addr[:], grp)
***REMOVED***

func (gsr *groupSourceReq) setSourceGroup(grp, src net.IP) ***REMOVED***
	sa := (*sockaddrInet)(unsafe.Pointer(uintptr(unsafe.Pointer(gsr)) + 4))
	sa.Len = sizeofSockaddrInet
	sa.Family = syscall.AF_INET
	copy(sa.Addr[:], grp)
	sa = (*sockaddrInet)(unsafe.Pointer(uintptr(unsafe.Pointer(gsr)) + 132))
	sa.Len = sizeofSockaddrInet
	sa.Family = syscall.AF_INET
	copy(sa.Addr[:], src)
***REMOVED***
