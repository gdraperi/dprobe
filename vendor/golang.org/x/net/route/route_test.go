// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package route

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

func (m *RouteMessage) String() string ***REMOVED***
	return fmt.Sprintf("%s", addrAttrs(nativeEndian.Uint32(m.raw[12:16])))
***REMOVED***

func (m *InterfaceMessage) String() string ***REMOVED***
	var attrs addrAttrs
	if runtime.GOOS == "openbsd" ***REMOVED***
		attrs = addrAttrs(nativeEndian.Uint32(m.raw[12:16]))
	***REMOVED*** else ***REMOVED***
		attrs = addrAttrs(nativeEndian.Uint32(m.raw[4:8]))
	***REMOVED***
	return fmt.Sprintf("%s", attrs)
***REMOVED***

func (m *InterfaceAddrMessage) String() string ***REMOVED***
	var attrs addrAttrs
	if runtime.GOOS == "openbsd" ***REMOVED***
		attrs = addrAttrs(nativeEndian.Uint32(m.raw[12:16]))
	***REMOVED*** else ***REMOVED***
		attrs = addrAttrs(nativeEndian.Uint32(m.raw[4:8]))
	***REMOVED***
	return fmt.Sprintf("%s", attrs)
***REMOVED***

func (m *InterfaceMulticastAddrMessage) String() string ***REMOVED***
	return fmt.Sprintf("%s", addrAttrs(nativeEndian.Uint32(m.raw[4:8])))
***REMOVED***

func (m *InterfaceAnnounceMessage) String() string ***REMOVED***
	what := "<nil>"
	switch m.What ***REMOVED***
	case 0:
		what = "arrival"
	case 1:
		what = "departure"
	***REMOVED***
	return fmt.Sprintf("(%d %s %s)", m.Index, m.Name, what)
***REMOVED***

func (m *InterfaceMetrics) String() string ***REMOVED***
	return fmt.Sprintf("(type=%d mtu=%d)", m.Type, m.MTU)
***REMOVED***

func (m *RouteMetrics) String() string ***REMOVED***
	return fmt.Sprintf("(pmtu=%d)", m.PathMTU)
***REMOVED***

type addrAttrs uint

var addrAttrNames = [...]string***REMOVED***
	"dst",
	"gateway",
	"netmask",
	"genmask",
	"ifp",
	"ifa",
	"author",
	"brd",
	"df:mpls1-n:tag-o:src", // mpls1 for dragonfly, tag for netbsd, src for openbsd
	"df:mpls2-o:srcmask",   // mpls2 for dragonfly, srcmask for openbsd
	"df:mpls3-o:label",     // mpls3 for dragonfly, label for openbsd
	"o:bfd",                // bfd for openbsd
	"o:dns",                // dns for openbsd
	"o:static",             // static for openbsd
	"o:search",             // search for openbsd
***REMOVED***

func (attrs addrAttrs) String() string ***REMOVED***
	var s string
	for i, name := range addrAttrNames ***REMOVED***
		if attrs&(1<<uint(i)) != 0 ***REMOVED***
			if s != "" ***REMOVED***
				s += "|"
			***REMOVED***
			s += name
		***REMOVED***
	***REMOVED***
	if s == "" ***REMOVED***
		return "<nil>"
	***REMOVED***
	return s
***REMOVED***

type msgs []Message

func (ms msgs) validate() ([]string, error) ***REMOVED***
	var ss []string
	for _, m := range ms ***REMOVED***
		switch m := m.(type) ***REMOVED***
		case *RouteMessage:
			if err := addrs(m.Addrs).match(addrAttrs(nativeEndian.Uint32(m.raw[12:16]))); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			sys := m.Sys()
			if sys == nil ***REMOVED***
				return nil, fmt.Errorf("no sys for %s", m.String())
			***REMOVED***
			ss = append(ss, m.String()+" "+syss(sys).String()+" "+addrs(m.Addrs).String())
		case *InterfaceMessage:
			var attrs addrAttrs
			if runtime.GOOS == "openbsd" ***REMOVED***
				attrs = addrAttrs(nativeEndian.Uint32(m.raw[12:16]))
			***REMOVED*** else ***REMOVED***
				attrs = addrAttrs(nativeEndian.Uint32(m.raw[4:8]))
			***REMOVED***
			if err := addrs(m.Addrs).match(attrs); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			sys := m.Sys()
			if sys == nil ***REMOVED***
				return nil, fmt.Errorf("no sys for %s", m.String())
			***REMOVED***
			ss = append(ss, m.String()+" "+syss(sys).String()+" "+addrs(m.Addrs).String())
		case *InterfaceAddrMessage:
			var attrs addrAttrs
			if runtime.GOOS == "openbsd" ***REMOVED***
				attrs = addrAttrs(nativeEndian.Uint32(m.raw[12:16]))
			***REMOVED*** else ***REMOVED***
				attrs = addrAttrs(nativeEndian.Uint32(m.raw[4:8]))
			***REMOVED***
			if err := addrs(m.Addrs).match(attrs); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			ss = append(ss, m.String()+" "+addrs(m.Addrs).String())
		case *InterfaceMulticastAddrMessage:
			if err := addrs(m.Addrs).match(addrAttrs(nativeEndian.Uint32(m.raw[4:8]))); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			ss = append(ss, m.String()+" "+addrs(m.Addrs).String())
		case *InterfaceAnnounceMessage:
			ss = append(ss, m.String())
		default:
			ss = append(ss, fmt.Sprintf("%+v", m))
		***REMOVED***
	***REMOVED***
	return ss, nil
***REMOVED***

type syss []Sys

func (sys syss) String() string ***REMOVED***
	var s string
	for _, sy := range sys ***REMOVED***
		switch sy := sy.(type) ***REMOVED***
		case *InterfaceMetrics:
			if len(s) > 0 ***REMOVED***
				s += " "
			***REMOVED***
			s += sy.String()
		case *RouteMetrics:
			if len(s) > 0 ***REMOVED***
				s += " "
			***REMOVED***
			s += sy.String()
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***

type addrFamily int

func (af addrFamily) String() string ***REMOVED***
	switch af ***REMOVED***
	case sysAF_UNSPEC:
		return "unspec"
	case sysAF_LINK:
		return "link"
	case sysAF_INET:
		return "inet4"
	case sysAF_INET6:
		return "inet6"
	default:
		return fmt.Sprintf("%d", af)
	***REMOVED***
***REMOVED***

const hexDigit = "0123456789abcdef"

type llAddr []byte

func (a llAddr) String() string ***REMOVED***
	if len(a) == 0 ***REMOVED***
		return ""
	***REMOVED***
	buf := make([]byte, 0, len(a)*3-1)
	for i, b := range a ***REMOVED***
		if i > 0 ***REMOVED***
			buf = append(buf, ':')
		***REMOVED***
		buf = append(buf, hexDigit[b>>4])
		buf = append(buf, hexDigit[b&0xF])
	***REMOVED***
	return string(buf)
***REMOVED***

type ipAddr []byte

func (a ipAddr) String() string ***REMOVED***
	if len(a) == 0 ***REMOVED***
		return "<nil>"
	***REMOVED***
	if len(a) == 4 ***REMOVED***
		return fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
	***REMOVED***
	if len(a) == 16 ***REMOVED***
		return fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x", a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], a[14], a[15])
	***REMOVED***
	s := make([]byte, len(a)*2)
	for i, tn := range a ***REMOVED***
		s[i*2], s[i*2+1] = hexDigit[tn>>4], hexDigit[tn&0xf]
	***REMOVED***
	return string(s)
***REMOVED***

func (a *LinkAddr) String() string ***REMOVED***
	name := a.Name
	if name == "" ***REMOVED***
		name = "<nil>"
	***REMOVED***
	lla := llAddr(a.Addr).String()
	if lla == "" ***REMOVED***
		lla = "<nil>"
	***REMOVED***
	return fmt.Sprintf("(%v %d %s %s)", addrFamily(a.Family()), a.Index, name, lla)
***REMOVED***

func (a *Inet4Addr) String() string ***REMOVED***
	return fmt.Sprintf("(%v %v)", addrFamily(a.Family()), ipAddr(a.IP[:]))
***REMOVED***

func (a *Inet6Addr) String() string ***REMOVED***
	return fmt.Sprintf("(%v %v %d)", addrFamily(a.Family()), ipAddr(a.IP[:]), a.ZoneID)
***REMOVED***

func (a *DefaultAddr) String() string ***REMOVED***
	return fmt.Sprintf("(%v %s)", addrFamily(a.Family()), ipAddr(a.Raw[2:]).String())
***REMOVED***

type addrs []Addr

func (as addrs) String() string ***REMOVED***
	var s string
	for _, a := range as ***REMOVED***
		if a == nil ***REMOVED***
			continue
		***REMOVED***
		if len(s) > 0 ***REMOVED***
			s += " "
		***REMOVED***
		switch a := a.(type) ***REMOVED***
		case *LinkAddr:
			s += a.String()
		case *Inet4Addr:
			s += a.String()
		case *Inet6Addr:
			s += a.String()
		case *DefaultAddr:
			s += a.String()
		***REMOVED***
	***REMOVED***
	if s == "" ***REMOVED***
		return "<nil>"
	***REMOVED***
	return s
***REMOVED***

func (as addrs) match(attrs addrAttrs) error ***REMOVED***
	var ts addrAttrs
	af := sysAF_UNSPEC
	for i := range as ***REMOVED***
		if as[i] != nil ***REMOVED***
			ts |= 1 << uint(i)
		***REMOVED***
		switch as[i].(type) ***REMOVED***
		case *Inet4Addr:
			if af == sysAF_UNSPEC ***REMOVED***
				af = sysAF_INET
			***REMOVED***
			if af != sysAF_INET ***REMOVED***
				return fmt.Errorf("got %v; want %v", addrs(as), addrFamily(af))
			***REMOVED***
		case *Inet6Addr:
			if af == sysAF_UNSPEC ***REMOVED***
				af = sysAF_INET6
			***REMOVED***
			if af != sysAF_INET6 ***REMOVED***
				return fmt.Errorf("got %v; want %v", addrs(as), addrFamily(af))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if ts != attrs && ts > attrs ***REMOVED***
		return fmt.Errorf("%v not included in %v", ts, attrs)
	***REMOVED***
	return nil
***REMOVED***

func fetchAndParseRIB(af int, typ RIBType) ([]Message, error) ***REMOVED***
	var err error
	var b []byte
	for i := 0; i < 3; i++ ***REMOVED***
		if b, err = FetchRIB(af, typ, 0); err != nil ***REMOVED***
			time.Sleep(10 * time.Millisecond)
			continue
		***REMOVED***
		break
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("%v %d %v", addrFamily(af), typ, err)
	***REMOVED***
	ms, err := ParseRIB(typ, b)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("%v %d %v", addrFamily(af), typ, err)
	***REMOVED***
	return ms, nil
***REMOVED***

// propVirtual is a proprietary virtual network interface.
type propVirtual struct ***REMOVED***
	name         string
	addr, mask   string
	setupCmds    []*exec.Cmd
	teardownCmds []*exec.Cmd
***REMOVED***

func (pv *propVirtual) setup() error ***REMOVED***
	for _, cmd := range pv.setupCmds ***REMOVED***
		if err := cmd.Run(); err != nil ***REMOVED***
			pv.teardown()
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (pv *propVirtual) teardown() error ***REMOVED***
	for _, cmd := range pv.teardownCmds ***REMOVED***
		if err := cmd.Run(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (pv *propVirtual) configure(suffix int) error ***REMOVED***
	if runtime.GOOS == "openbsd" ***REMOVED***
		pv.name = fmt.Sprintf("vether%d", suffix)
	***REMOVED*** else ***REMOVED***
		pv.name = fmt.Sprintf("vlan%d", suffix)
	***REMOVED***
	xname, err := exec.LookPath("ifconfig")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	pv.setupCmds = append(pv.setupCmds, &exec.Cmd***REMOVED***
		Path: xname,
		Args: []string***REMOVED***"ifconfig", pv.name, "create"***REMOVED***,
	***REMOVED***)
	if runtime.GOOS == "netbsd" ***REMOVED***
		// NetBSD requires an underlying dot1Q-capable network
		// interface.
		pv.setupCmds = append(pv.setupCmds, &exec.Cmd***REMOVED***
			Path: xname,
			Args: []string***REMOVED***"ifconfig", pv.name, "vlan", fmt.Sprintf("%d", suffix&0xfff), "vlanif", "wm0"***REMOVED***,
		***REMOVED***)
	***REMOVED***
	pv.setupCmds = append(pv.setupCmds, &exec.Cmd***REMOVED***
		Path: xname,
		Args: []string***REMOVED***"ifconfig", pv.name, "inet", pv.addr, "netmask", pv.mask***REMOVED***,
	***REMOVED***)
	pv.teardownCmds = append(pv.teardownCmds, &exec.Cmd***REMOVED***
		Path: xname,
		Args: []string***REMOVED***"ifconfig", pv.name, "destroy"***REMOVED***,
	***REMOVED***)
	return nil
***REMOVED***
