// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build netbsd openbsd

package ipv4

import (
	"net"
	"syscall"

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
		ssoMulticastLoopback:  ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_MULTICAST_LOOP, Len: 1***REMOVED******REMOVED***,
		ssoReceiveTTL:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_RECVTTL, Len: 4***REMOVED******REMOVED***,
		ssoReceiveDst:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_RECVDSTADDR, Len: 4***REMOVED******REMOVED***,
		ssoReceiveInterface:   ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_RECVIF, Len: 4***REMOVED******REMOVED***,
		ssoHeaderPrepend:      ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_HDRINCL, Len: 4***REMOVED******REMOVED***,
		ssoJoinGroup:          ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_ADD_MEMBERSHIP, Len: sizeofIPMreq***REMOVED***, typ: ssoTypeIPMreq***REMOVED***,
		ssoLeaveGroup:         ***REMOVED***Option: socket.Option***REMOVED***Level: iana.ProtocolIP, Name: sysIP_DROP_MEMBERSHIP, Len: sizeofIPMreq***REMOVED***, typ: ssoTypeIPMreq***REMOVED***,
	***REMOVED***
)
