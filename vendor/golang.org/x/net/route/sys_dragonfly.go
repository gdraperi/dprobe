// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import "unsafe"

func (typ RIBType) parseable() bool ***REMOVED*** return true ***REMOVED***

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
			PathMTU: int(nativeEndian.Uint64(m.raw[m.extOff+8 : m.extOff+16])),
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
	var p uintptr
	rtm := &wireFormat***REMOVED***extOff: 40, bodyOff: sizeofRtMsghdrDragonFlyBSD4***REMOVED***
	rtm.parse = rtm.parseRouteMessage
	ifm := &wireFormat***REMOVED***extOff: 16, bodyOff: sizeofIfMsghdrDragonFlyBSD4***REMOVED***
	ifm.parse = ifm.parseInterfaceMessage
	ifam := &wireFormat***REMOVED***extOff: sizeofIfaMsghdrDragonFlyBSD4, bodyOff: sizeofIfaMsghdrDragonFlyBSD4***REMOVED***
	ifam.parse = ifam.parseInterfaceAddrMessage
	ifmam := &wireFormat***REMOVED***extOff: sizeofIfmaMsghdrDragonFlyBSD4, bodyOff: sizeofIfmaMsghdrDragonFlyBSD4***REMOVED***
	ifmam.parse = ifmam.parseInterfaceMulticastAddrMessage
	ifanm := &wireFormat***REMOVED***extOff: sizeofIfAnnouncemsghdrDragonFlyBSD4, bodyOff: sizeofIfAnnouncemsghdrDragonFlyBSD4***REMOVED***
	ifanm.parse = ifanm.parseInterfaceAnnounceMessage
	return int(unsafe.Sizeof(p)), map[int]*wireFormat***REMOVED***
		sysRTM_ADD:        rtm,
		sysRTM_DELETE:     rtm,
		sysRTM_CHANGE:     rtm,
		sysRTM_GET:        rtm,
		sysRTM_LOSING:     rtm,
		sysRTM_REDIRECT:   rtm,
		sysRTM_MISS:       rtm,
		sysRTM_LOCK:       rtm,
		sysRTM_RESOLVE:    rtm,
		sysRTM_NEWADDR:    ifam,
		sysRTM_DELADDR:    ifam,
		sysRTM_IFINFO:     ifm,
		sysRTM_NEWMADDR:   ifmam,
		sysRTM_DELMADDR:   ifmam,
		sysRTM_IFANNOUNCE: ifanm,
	***REMOVED***
***REMOVED***
