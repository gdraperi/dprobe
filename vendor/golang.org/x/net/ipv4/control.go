// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

import (
	"fmt"
	"net"
	"sync"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/socket"
)

type rawOpt struct ***REMOVED***
	sync.RWMutex
	cflags ControlFlags
***REMOVED***

func (c *rawOpt) set(f ControlFlags)        ***REMOVED*** c.cflags |= f ***REMOVED***
func (c *rawOpt) clear(f ControlFlags)      ***REMOVED*** c.cflags &^= f ***REMOVED***
func (c *rawOpt) isset(f ControlFlags) bool ***REMOVED*** return c.cflags&f != 0 ***REMOVED***

type ControlFlags uint

const (
	FlagTTL       ControlFlags = 1 << iota // pass the TTL on the received packet
	FlagSrc                                // pass the source address on the received packet
	FlagDst                                // pass the destination address on the received packet
	FlagInterface                          // pass the interface index on the received packet
)

// A ControlMessage represents per packet basis IP-level socket options.
type ControlMessage struct ***REMOVED***
	// Receiving socket options: SetControlMessage allows to
	// receive the options from the protocol stack using ReadFrom
	// method of PacketConn or RawConn.
	//
	// Specifying socket options: ControlMessage for WriteTo
	// method of PacketConn or RawConn allows to send the options
	// to the protocol stack.
	//
	TTL     int    // time-to-live, receiving only
	Src     net.IP // source address, specifying only
	Dst     net.IP // destination address, receiving only
	IfIndex int    // interface index, must be 1 <= value when specifying
***REMOVED***

func (cm *ControlMessage) String() string ***REMOVED***
	if cm == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	return fmt.Sprintf("ttl=%d src=%v dst=%v ifindex=%d", cm.TTL, cm.Src, cm.Dst, cm.IfIndex)
***REMOVED***

// Marshal returns the binary encoding of cm.
func (cm *ControlMessage) Marshal() []byte ***REMOVED***
	if cm == nil ***REMOVED***
		return nil
	***REMOVED***
	var m socket.ControlMessage
	if ctlOpts[ctlPacketInfo].name > 0 && (cm.Src.To4() != nil || cm.IfIndex > 0) ***REMOVED***
		m = socket.NewControlMessage([]int***REMOVED***ctlOpts[ctlPacketInfo].length***REMOVED***)
	***REMOVED***
	if len(m) > 0 ***REMOVED***
		ctlOpts[ctlPacketInfo].marshal(m, cm)
	***REMOVED***
	return m
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
		if lvl != iana.ProtocolIP ***REMOVED***
			continue
		***REMOVED***
		switch ***REMOVED***
		case typ == ctlOpts[ctlTTL].name && l >= ctlOpts[ctlTTL].length:
			ctlOpts[ctlTTL].parse(cm, m.Data(l))
		case typ == ctlOpts[ctlDst].name && l >= ctlOpts[ctlDst].length:
			ctlOpts[ctlDst].parse(cm, m.Data(l))
		case typ == ctlOpts[ctlInterface].name && l >= ctlOpts[ctlInterface].length:
			ctlOpts[ctlInterface].parse(cm, m.Data(l))
		case typ == ctlOpts[ctlPacketInfo].name && l >= ctlOpts[ctlPacketInfo].length:
			ctlOpts[ctlPacketInfo].parse(cm, m.Data(l))
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
	if opt.isset(FlagTTL) && ctlOpts[ctlTTL].name > 0 ***REMOVED***
		l += socket.ControlMessageSpace(ctlOpts[ctlTTL].length)
	***REMOVED***
	if ctlOpts[ctlPacketInfo].name > 0 ***REMOVED***
		if opt.isset(FlagSrc | FlagDst | FlagInterface) ***REMOVED***
			l += socket.ControlMessageSpace(ctlOpts[ctlPacketInfo].length)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if opt.isset(FlagDst) && ctlOpts[ctlDst].name > 0 ***REMOVED***
			l += socket.ControlMessageSpace(ctlOpts[ctlDst].length)
		***REMOVED***
		if opt.isset(FlagInterface) && ctlOpts[ctlInterface].name > 0 ***REMOVED***
			l += socket.ControlMessageSpace(ctlOpts[ctlInterface].length)
		***REMOVED***
	***REMOVED***
	var b []byte
	if l > 0 ***REMOVED***
		b = make([]byte, l)
	***REMOVED***
	return b
***REMOVED***

// Ancillary data socket options
const (
	ctlTTL        = iota // header field
	ctlSrc               // header field
	ctlDst               // header field
	ctlInterface         // inbound or outbound interface
	ctlPacketInfo        // inbound or outbound packet path
	ctlMax
)

// A ctlOpt represents a binding for ancillary data socket option.
type ctlOpt struct ***REMOVED***
	name    int // option name, must be equal or greater than 1
	length  int // option length
	marshal func([]byte, *ControlMessage) []byte
	parse   func(*ControlMessage, []byte)
***REMOVED***
