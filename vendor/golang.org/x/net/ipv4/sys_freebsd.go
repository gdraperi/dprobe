// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

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
		ssoJoinGroup:          ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_JOIN_GROUP, Len: sizeofGroupReq***REMOVED***, typ: ssoTypeGroupReq***REMOVED***,
		ssoLeaveGroup:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_LEAVE_GROUP, Len: sizeofGroupReq***REMOVED***, typ: ssoTypeGroupReq***REMOVED***,
		ssoJoinSourceGroup:    ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_JOIN_SOURCE_GROUP, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
		ssoLeaveSourceGroup:   ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_LEAVE_SOURCE_GROUP, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
		ssoBlockSourceGroup:   ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_BLOCK_SOURCE, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
		ssoUnblockSourceGroup: ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysMCAST_UNBLOCK_SOURCE, Len: sizeofGroupSourceReq***REMOVED***, typ: ssoTypeGroupSourceReq***REMOVED***,
	***REMOVED***
)

func init() ***REMOVED***
	freebsdVersion, _ = syscall.SysctlUint32("kern.osreldate")
	if freebsdVersion >= 1000000 ***REMOVED***
		sockOpts[ssoMulticastInterface] = &sockOpt***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_MULTICAST_IF, Len: sizeofIPMreqn***REMOVED***, typ: ssoTypeIPMreqn***REMOVED***
	***REMOVED***
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

func (gr *groupReq) setGroup(grp net.IP) ***REMOVED***
	sa := (*sockaddrInet)(unsafe.Pointer(&gr.Group))
	sa.Len = sizeofSockaddrInet
	sa.Family = syscall.AF_INET
	copy(sa.Addr[:], grp)
***REMOVED***

func (gsr *groupSourceReq) setSourceGroup(grp, src net.IP) ***REMOVED***
	sa := (*sockaddrInet)(unsafe.Pointer(&gsr.Group))
	sa.Len = sizeofSockaddrInet
	sa.Family = syscall.AF_INET
	copy(sa.Addr[:], grp)
	sa = (*sockaddrInet)(unsafe.Pointer(&gsr.Source))
	sa.Len = sizeofSockaddrInet
	sa.Family = syscall.AF_INET
	copy(sa.Addr[:], src)
***REMOVED***
