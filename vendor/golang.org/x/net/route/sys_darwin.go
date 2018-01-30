// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

func (typ RIBType) parseable() bool ***REMOVED***
	switch typ ***REMOVED***
	case sysNET_RT_STAT, sysNET_RT_TRASH:
		return false
	default:
		return true
	***REMOVED***
***REMOVED***

// RouteMetrics represents route metrics.
type RouteMetrics struct ***REMOVED***
	PathMTU int // path maximum transmission unit
***REMOVED***

// SysType implements the SysType method of Sys interface.
func (rmx *RouteMetrics) SysType() SysType ***REMOVED*** return SysMetrics ***REMOVED***

// Sys implements the Sys method of Message interface.
func (m *RouteMessage) Sys() []Sys ***REMOVED***
	return []Sys***REMOVED***
		&RouteMetrics***REMOVED***
			PathMTU: int(nativeEndian.Uint32(m.raw[m.extOff+4 : m.extOff+8])),
		***REMOVED***,
	***REMOVED***
***REMOVED***

// InterfaceMetrics represents interface metrics.
type InterfaceMetrics struct ***REMOVED***
	Type int // interface type
	MTU  int // maximum transmission unit
***REMOVED***

// SysType implements the SysType method of Sys interface.
func (imx *InterfaceMetrics) SysType() SysType ***REMOVED*** return SysMetrics ***REMOVED***

// Sys implements the Sys method of Message interface.
func (m *InterfaceMessage) Sys() []Sys ***REMOVED***
	return []Sys***REMOVED***
		&InterfaceMetrics***REMOVED***
			Type: int(m.raw[m.extOff]),
			MTU:  int(nativeEndian.Uint32(m.raw[m.extOff+8 : m.extOff+12])),
		***REMOVED***,
	***REMOVED***
***REMOVED***

func probeRoutingStack() (int, map[int]*wireFormat) ***REMOVED***
	rtm := &wireFormat***REMOVED***extOff: 36, bodyOff: sizeofRtMsghdrDarwin15***REMOVED***
	rtm.parse = rtm.parseRouteMessage
	rtm2 := &wireFormat***REMOVED***extOff: 36, bodyOff: sizeofRtMsghdr2Darwin15***REMOVED***
	rtm2.parse = rtm2.parseRouteMessage
	ifm := &wireFormat***REMOVED***extOff: 16, bodyOff: sizeofIfMsghdrDarwin15***REMOVED***
	ifm.parse = ifm.parseInterfaceMessage
	ifm2 := &wireFormat***REMOVED***extOff: 32, bodyOff: sizeofIfMsghdr2Darwin15***REMOVED***
	ifm2.parse = ifm2.parseInterfaceMessage
	ifam := &wireFormat***REMOVED***extOff: sizeofIfaMsghdrDarwin15, bodyOff: sizeofIfaMsghdrDarwin15***REMOVED***
	ifam.parse = ifam.parseInterfaceAddrMessage
	ifmam := &wireFormat***REMOVED***extOff: sizeofIfmaMsghdrDarwin15, bodyOff: sizeofIfmaMsghdrDarwin15***REMOVED***
	ifmam.parse = ifmam.parseInterfaceMulticastAddrMessage
	ifmam2 := &wireFormat***REMOVED***extOff: sizeofIfmaMsghdr2Darwin15, bodyOff: sizeofIfmaMsghdr2Darwin15***REMOVED***
	ifmam2.parse = ifmam2.parseInterfaceMulticastAddrMessage
	// Darwin kernels require 32-bit aligned access to routing facilities.
	return 4, map[int]*wireFormat***REMOVED***
		sysRTM_ADD:       rtm,
		sysRTM_DELETE:    rtm,
		sysRTM_CHANGE:    rtm,
		sysRTM_GET:       rtm,
		sysRTM_LOSING:    rtm,
		sysRTM_REDIRECT:  rtm,
		sysRTM_MISS:      rtm,
		sysRTM_LOCK:      rtm,
		sysRTM_RESOLVE:   rtm,
		sysRTM_NEWADDR:   ifam,
		sysRTM_DELADDR:   ifam,
		sysRTM_IFINFO:    ifm,
		sysRTM_NEWMADDR:  ifmam,
		sysRTM_DELMADDR:  ifmam,
		sysRTM_IFINFO2:   ifm2,
		sysRTM_NEWMADDR2: ifmam2,
		sysRTM_GET2:      rtm2,
	***REMOVED***
***REMOVED***
