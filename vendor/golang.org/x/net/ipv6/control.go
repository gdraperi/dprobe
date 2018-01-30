// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

import (
	"fmt"
	"net"
	"sync"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/socket"
)

// Note that RFC 3542 obsoletes RFC 2292 but OS X Snow Leopard and the
// former still support RFC 2292 only. Please be aware that almost
// all protocol implementations prohibit using a combination of RFC
// 2292 and RFC 3542 for some practical reasons.

type rawOpt struct ***REMOVED***
	sync.RWMutex
	cflags ControlFlags
***REMOVED***

func (c *rawOpt) set(f ControlFlags)        ***REMOVED*** c.cflags |= f ***REMOVED***
func (c *rawOpt) clear(f ControlFlags)      ***REMOVED*** c.cflags &^= f ***REMOVED***
func (c *rawOpt) isset(f ControlFlags) bool ***REMOVED*** return c.cflags&f != 0 ***REMOVED***

// A ControlFlags represents per packet basis IP-level socket option
// control flags.
type ControlFlags uint

const (
	FlagTrafficClass ControlFlags = 1 << iota // pass the traffic class on the received packet
	FlagHopLimit                              // pass the hop limit on the received packet
	FlagSrc                                   // pass the source address on the received packet
	FlagDst                                   // pass the destination address on the received packet
	FlagInterface                             // pass the interface index on the received packet
	FlagPathMTU                               // pass the path MTU on the received packet path
)

const flagPacketInfo = FlagDst | FlagInterface

// A ControlMessage represents per packet basis IP-level socket
// options.
type ControlMessage struct ***REMOVED***
	// Receiving socket options: SetControlMessage allows to
	// receive the options from the protocol stack using ReadFrom
	// method of PacketConn.
	//
	// Specifying socket options: ControlMessage for WriteTo
	// method of PacketConn allows to send the options to the
	// protocol stack.
	//
	TrafficClass int    // traffic class, must be 1 <= value <= 255 when specifying
	HopLimit     int    // hop limit, must be 1 <= value <= 255 when specifying
	Src          net.IP // source address, specifying only
	Dst          net.IP // destination address, receiving only
	IfIndex      int    // interface index, must be 1 <= value when specifying
	NextHop      net.IP // next hop address, specifying only
	MTU          int    // path MTU, receiving only
***REMOVED***

func (cm *ControlMessage) String() string ***REMOVED***
	if cm == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	return fmt.Sprintf("tclass=%#x hoplim=%d src=%v dst=%v ifindex=%d nexthop=%v mtu=%d", cm.TrafficClass, cm.HopLimit, cm.Src, cm.Dst, cm.IfIndex, cm.NextHop, cm.MTU)
***REMOVED***

// Marshal returns the binary encoding of cm.
func (cm *ControlMessage) Marshal() []byte ***REMOVED***
	if cm == nil ***REMOVED***
		return nil
	***REMOVED***
	var l int
	tclass := false
	if ctlOpts[ctlTrafficClass].name > 0 && cm.TrafficClass > 0 ***REMOVED***
		tclass = true
		l += socket.ControlMessageSpace(ctlOpts[ctlTrafficClass].length)
	***REMOVED***
	hoplimit := false
	if ctlOpts[ctlHopLimit].name > 0 && cm.HopLimit > 0 ***REMOVED***
		hoplimit = true
		l += socket.ControlMessageSpace(ctlOpts[ctlHopLimit].length)
	***REMOVED***
	pktinfo := false
	if ctlOpts[ctlPacketInfo].name > 0 && (cm.Src.To16() != nil && cm.Src.To4() == nil || cm.IfIndex > 0) ***REMOVED***
		pktinfo = true
		l += socket.ControlMessageSpace(ctlOpts[ctlPacketInfo].length)
	***REMOVED***
	nexthop := false
	if ctlOpts[ctlNextHop].name > 0 && cm.NextHop.To16() != nil && cm.NextHop.To4() == nil ***REMOVED***
		nexthop = true
		l += socket.ControlMessageSpace(ctlOpts[ctlNextHop].length)
	***REMOVED***
	var b []byte
	if l > 0 ***REMOVED***
		b = make([]byte, l)
		bb := b
		if tclass ***REMOVED***
			bb = ctlOpts[ctlTrafficClass].marshal(bb, cm)
		***REMOVED***
		if hoplimit ***REMOVED***
			bb = ctlOpts[ctlHopLimit].marshal(bb, cm)
		***REMOVED***
		if pktinfo ***REMOVED***
			bb = ctlOpts[ctlPacketInfo].marshal(bb, cm)
		***REMOVED***
		if nexthop ***REMOVED***
			bb = ctlOpts[ctlNextHop].marshal(bb, cm)
		***REMOVED***
	***REMOVED***
	return b
***REMOVED***

// Parse parses b as a control message and stores the result in cm.
func (cm *ControlMessage) Parse(b []byte) error ***REMOVED***
	ms, err := socket.ControlMessage(b).Parse()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, m := range ms ***REMOVED***
		lvl, typ, l, err := m.ParseHeader()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if lvl != iana.ProtocolIPv6 ***REMOVED***
			continue
		***REMOVED***
		switch ***REMOVED***
		case typ == ctlOpts[ctlTrafficClass].name && l >= ctlOpts[ctlTrafficClass].length:
			ctlOpts[ctlTrafficClass].parse(cm, m.Data(l))
		case typ == ctlOpts[ctlHopLimit].name && l >= ctlOpts[ctlHopLimit].length:
			ctlOpts[ctlHopLimit].parse(cm, m.Data(l))
		case typ == ctlOpts[ctlPacketInfo].name && l >= ctlOpts[ctlPacketInfo].length:
			ctlOpts[ctlPacketInfo].parse(cm, m.Data(l))
		case typ == ctlOpts[ctlPathMTU].name && l >= ctlOpts[ctlPathMTU].length:
			ctlOpts[ctlPathMTU].parse(cm, m.Data(l))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// NewControlMessage returns a new control message.
//
// The returned message is large enough for options specified by cf.
func NewControlMessage(cf ControlFlags) []byte ***REMOVED***
	opt := rawOpt***REMOVED***cflags: cf***REMOVED***
	var l int
	if opt.isset(FlagTrafficClass) && ctlOpts[ctlTrafficClass].name > 0 ***REMOVED***
		l += socket.ControlMessageSpace(ctlOpts[ctlTrafficClass].length)
	***REMOVED***
	if opt.isset(FlagHopLimit) && ctlOpts[ctlHopLimit].name > 0 ***REMOVED***
		l += socket.ControlMessageSpace(ctlOpts[ctlHopLimit].length)
	***REMOVED***
	if opt.isset(flagPacketInfo) && ctlOpts[ctlPacketInfo].name > 0 ***REMOVED***
		l += socket.ControlMessageSpace(ctlOpts[ctlPacketInfo].length)
	***REMOVED***
	if opt.isset(FlagPathMTU) && ctlOpts[ctlPathMTU].name > 0 ***REMOVED***
		l += socket.ControlMessageSpace(ctlOpts[ctlPathMTU].length)
	***REMOVED***
	var b []byte
	if l > 0 ***REMOVED***
		b = make([]byte, l)
	***REMOVED***
	return b
***REMOVED***

// Ancillary data socket options
const (
	ctlTrafficClass = iota // header field
	ctlHopLimit            // header field
	ctlPacketInfo          // inbound or outbound packet path
	ctlNextHop             // nexthop
	ctlPathMTU             // path mtu
	ctlMax
)

// A ctlOpt represents a binding for ancillary data socket option.
type ctlOpt struct ***REMOVED***
	name    int // option name, must be equal or greater than 1
	length  int // option length
	marshal func([]byte, *ControlMessage) []byte
	parse   func(*ControlMessage, []byte)
***REMOVED***
