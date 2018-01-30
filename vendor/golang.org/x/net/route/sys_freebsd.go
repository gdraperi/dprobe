// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"syscall"
	"unsafe"
)

func (typ RIBType) parseable() bool ***REMOVED*** return true ***REMOVED***

// RouteMetrics represents route metrics.
type RouteMetrics struct ***REMOVED***
	PathMTU int // path maximum transmission unit
***REMOVED***

// SysType implements the SysType method of Sys interface.
func (rmx *RouteMetrics) SysType() SysType ***REMOVED*** return SysMetrics ***REMOVED***

// Sys implements the Sys method of Message interface.
func (m *RouteMessage) Sys() []Sys ***REMOVED***
	if kernelAlign == 8 ***REMOVED***
		return []Sys***REMOVED***
			&RouteMetrics***REMOVED***
				PathMTU: int(nativeEndian.Uint64(m.raw[m.extOff+8 : m.extOff+16])),
			***REMOVED***,
		***REMOVED***
	***REMOVED***
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
	var p uintptr
	wordSize := int(unsafe.Sizeof(p))
	align := int(unsafe.Sizeof(p))
	// In the case of kern.supported_archs="amd64 i386", we need
	// to know the underlying kernel's architecture because the
	// alignment for routing facilities are set at the build time
	// of the kernel.
	conf, _ := syscall.Sysctl("kern.conftxt")
	for i, j := 0, 0; j < len(conf); j++ ***REMOVED***
		if conf[j] != '\n' ***REMOVED***
			continue
		***REMOVED***
		s := conf[i:j]
		i = j + 1
		if len(s) > len("machine") && s[:len("machine")] == "machine" ***REMOVED***
			s = s[len("machine"):]
			for k := 0; k < len(s); k++ ***REMOVED***
				if s[k] == ' ' || s[k] == '\t' ***REMOVED***
					s = s[1:]
				***REMOVED***
				break
			***REMOVED***
			if s == "amd64" ***REMOVED***
				align = 8
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	var rtm, ifm, ifam, ifmam, ifanm *wireFormat
	if align != wordSize ***REMOVED*** // 386 emulation on amd64
		rtm = &wireFormat***REMOVED***extOff: sizeofRtMsghdrFreeBSD10Emu - sizeofRtMetricsFreeBSD10Emu, bodyOff: sizeofRtMsghdrFreeBSD10Emu***REMOVED***
		ifm = &wireFormat***REMOVED***extOff: 16***REMOVED***
		ifam = &wireFormat***REMOVED***extOff: sizeofIfaMsghdrFreeBSD10Emu, bodyOff: sizeofIfaMsghdrFreeBSD10Emu***REMOVED***
		ifmam = &wireFormat***REMOVED***extOff: sizeofIfmaMsghdrFreeBSD10Emu, bodyOff: sizeofIfmaMsghdrFreeBSD10Emu***REMOVED***
		ifanm = &wireFormat***REMOVED***extOff: sizeofIfAnnouncemsghdrFreeBSD10Emu, bodyOff: sizeofIfAnnouncemsghdrFreeBSD10Emu***REMOVED***
	***REMOVED*** else ***REMOVED***
		rtm = &wireFormat***REMOVED***extOff: sizeofRtMsghdrFreeBSD10 - sizeofRtMetricsFreeBSD10, bodyOff: sizeofRtMsghdrFreeBSD10***REMOVED***
		ifm = &wireFormat***REMOVED***extOff: 16***REMOVED***
		ifam = &wireFormat***REMOVED***extOff: sizeofIfaMsghdrFreeBSD10, bodyOff: sizeofIfaMsghdrFreeBSD10***REMOVED***
		ifmam = &wireFormat***REMOVED***extOff: sizeofIfmaMsghdrFreeBSD10, bodyOff: sizeofIfmaMsghdrFreeBSD10***REMOVED***
		ifanm = &wireFormat***REMOVED***extOff: sizeofIfAnnouncemsghdrFreeBSD10, bodyOff: sizeofIfAnnouncemsghdrFreeBSD10***REMOVED***
	***REMOVED***
	rel, _ := syscall.SysctlUint32("kern.osreldate")
	switch ***REMOVED***
	case rel < 800000:
		if align != wordSize ***REMOVED*** // 386 emulation on amd64
			ifm.bodyOff = sizeofIfMsghdrFreeBSD7Emu
		***REMOVED*** else ***REMOVED***
			ifm.bodyOff = sizeofIfMsghdrFreeBSD7
		***REMOVED***
	case 800000 <= rel && rel < 900000:
		if align != wordSize ***REMOVED*** // 386 emulation on amd64
			ifm.bodyOff = sizeofIfMsghdrFreeBSD8Emu
		***REMOVED*** else ***REMOVED***
			ifm.bodyOff = sizeofIfMsghdrFreeBSD8
		***REMOVED***
	case 900000 <= rel && rel < 1000000:
		if align != wordSize ***REMOVED*** // 386 emulation on amd64
			ifm.bodyOff = sizeofIfMsghdrFreeBSD9Emu
		***REMOVED*** else ***REMOVED***
			ifm.bodyOff = sizeofIfMsghdrFreeBSD9
		***REMOVED***
	case 1000000 <= rel && rel < 1100000:
		if align != wordSize ***REMOVED*** // 386 emulation on amd64
			ifm.bodyOff = sizeofIfMsghdrFreeBSD10Emu
		***REMOVED*** else ***REMOVED***
			ifm.bodyOff = sizeofIfMsghdrFreeBSD10
		***REMOVED***
	default:
		if align != wordSize ***REMOVED*** // 386 emulation on amd64
			ifm.bodyOff = sizeofIfMsghdrFreeBSD11Emu
		***REMOVED*** else ***REMOVED***
			ifm.bodyOff = sizeofIfMsghdrFreeBSD11
		***REMOVED***
	***REMOVED***
	rtm.parse = rtm.parseRouteMessage
	ifm.parse = ifm.parseInterfaceMessage
	ifam.parse = ifam.parseInterfaceAddrMessage
	ifmam.parse = ifmam.parseInterfaceMulticastAddrMessage
	ifanm.parse = ifanm.parseInterfaceAnnounceMessage
	return align, map[int]*wireFormat***REMOVED***
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
