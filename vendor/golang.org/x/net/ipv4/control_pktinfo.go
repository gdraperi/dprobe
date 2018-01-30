// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux solaris

package ipv4

import (
	"net"
	"unsafe"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/socket"
)

func marshalPacketInfo(b []byte, cm *ControlMessage) []byte ***REMOVED***
	m := socket.ControlMessage(b)
	m.MarshalHeader(iana.ProtocolIP, sysIP_PKTINFO, sizeofInetPktinfo)
	if cm != nil ***REMOVED***
		pi := (*inetPktinfo)(unsafe.Pointer(&m.Data(sizeofInetPktinfo)[0]))
		if ip := cm.Src.To4(); ip != nil ***REMOVED***
			copy(pi.Spec_dst[:], ip)
		***REMOVED***
		if cm.IfIndex > 0 ***REMOVED***
			pi.setIfindex(cm.IfIndex)
		***REMOVED***
	***REMOVED***
	return m.Next(sizeofInetPktinfo)
***REMOVED***

func parsePacketInfo(cm *ControlMessage, b []byte) ***REMOVED***
	pi := (*inetPktinfo)(unsafe.Pointer(&b[0]))
	cm.IfIndex = int(pi.Ifindex)
	if len(cm.Dst) < net.IPv4len ***REMOVED***
		cm.Dst = make(net.IP, net.IPv4len)
	***REMOVED***
	copy(cm.Dst, pi.Addr[:])
***REMOVED***
