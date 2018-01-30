// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp

import (
	"encoding/binary"
	"net"
	"strings"

	"golang.org/x/net/internal/iana"
)

const (
	classInterfaceInfo = 2

	afiIPv4 = 1
	afiIPv6 = 2
)

const (
	attrMTU = 1 << iota
	attrName
	attrIPAddr
	attrIfIndex
)

// An InterfaceInfo represents interface and next-hop identification.
type InterfaceInfo struct ***REMOVED***
	Class     int // extension object class number
	Type      int // extension object sub-type
	Interface *net.Interface
	Addr      *net.IPAddr
***REMOVED***

func (ifi *InterfaceInfo) nameLen() int ***REMOVED***
	if len(ifi.Interface.Name) > 63 ***REMOVED***
		return 64
	***REMOVED***
	l := 1 + len(ifi.Interface.Name)
	return (l + 3) &^ 3
***REMOVED***

func (ifi *InterfaceInfo) attrsAndLen(proto int) (attrs, l int) ***REMOVED***
	l = 4
	if ifi.Interface != nil && ifi.Interface.Index > 0 ***REMOVED***
		attrs |= attrIfIndex
		l += 4
		if len(ifi.Interface.Name) > 0 ***REMOVED***
			attrs |= attrName
			l += ifi.nameLen()
		***REMOVED***
		if ifi.Interface.MTU > 0 ***REMOVED***
			attrs |= attrMTU
			l += 4
		***REMOVED***
	***REMOVED***
	if ifi.Addr != nil ***REMOVED***
		switch proto ***REMOVED***
		case iana.ProtocolICMP:
			if ifi.Addr.IP.To4() != nil ***REMOVED***
				attrs |= attrIPAddr
				l += 4 + net.IPv4len
			***REMOVED***
		case iana.ProtocolIPv6ICMP:
			if ifi.Addr.IP.To16() != nil && ifi.Addr.IP.To4() == nil ***REMOVED***
				attrs |= attrIPAddr
				l += 4 + net.IPv6len
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// Len implements the Len method of Extension interface.
func (ifi *InterfaceInfo) Len(proto int) int ***REMOVED***
	_, l := ifi.attrsAndLen(proto)
	return l
***REMOVED***

// Marshal implements the Marshal method of Extension interface.
func (ifi *InterfaceInfo) Marshal(proto int) ([]byte, error) ***REMOVED***
	attrs, l := ifi.attrsAndLen(proto)
	b := make([]byte, l)
	if err := ifi.marshal(proto, b, attrs, l); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b, nil
***REMOVED***

func (ifi *InterfaceInfo) marshal(proto int, b []byte, attrs, l int) error ***REMOVED***
	binary.BigEndian.PutUint16(b[:2], uint16(l))
	b[2], b[3] = classInterfaceInfo, byte(ifi.Type)
	for b = b[4:]; len(b) > 0 && attrs != 0; ***REMOVED***
		switch ***REMOVED***
		case attrs&attrIfIndex != 0:
			b = ifi.marshalIfIndex(proto, b)
			attrs &^= attrIfIndex
		case attrs&attrIPAddr != 0:
			b = ifi.marshalIPAddr(proto, b)
			attrs &^= attrIPAddr
		case attrs&attrName != 0:
			b = ifi.marshalName(proto, b)
			attrs &^= attrName
		case attrs&attrMTU != 0:
			b = ifi.marshalMTU(proto, b)
			attrs &^= attrMTU
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (ifi *InterfaceInfo) marshalIfIndex(proto int, b []byte) []byte ***REMOVED***
	binary.BigEndian.PutUint32(b[:4], uint32(ifi.Interface.Index))
	return b[4:]
***REMOVED***

func (ifi *InterfaceInfo) parseIfIndex(b []byte) ([]byte, error) ***REMOVED***
	if len(b) < 4 ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	ifi.Interface.Index = int(binary.BigEndian.Uint32(b[:4]))
	return b[4:], nil
***REMOVED***

func (ifi *InterfaceInfo) marshalIPAddr(proto int, b []byte) []byte ***REMOVED***
	switch proto ***REMOVED***
	case iana.ProtocolICMP:
		binary.BigEndian.PutUint16(b[:2], uint16(afiIPv4))
		copy(b[4:4+net.IPv4len], ifi.Addr.IP.To4())
		b = b[4+net.IPv4len:]
	case iana.ProtocolIPv6ICMP:
		binary.BigEndian.PutUint16(b[:2], uint16(afiIPv6))
		copy(b[4:4+net.IPv6len], ifi.Addr.IP.To16())
		b = b[4+net.IPv6len:]
	***REMOVED***
	return b
***REMOVED***

func (ifi *InterfaceInfo) parseIPAddr(b []byte) ([]byte, error) ***REMOVED***
	if len(b) < 4 ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	afi := int(binary.BigEndian.Uint16(b[:2]))
	b = b[4:]
	switch afi ***REMOVED***
	case afiIPv4:
		if len(b) < net.IPv4len ***REMOVED***
			return nil, errMessageTooShort
		***REMOVED***
		ifi.Addr.IP = make(net.IP, net.IPv4len)
		copy(ifi.Addr.IP, b[:net.IPv4len])
		b = b[net.IPv4len:]
	case afiIPv6:
		if len(b) < net.IPv6len ***REMOVED***
			return nil, errMessageTooShort
		***REMOVED***
		ifi.Addr.IP = make(net.IP, net.IPv6len)
		copy(ifi.Addr.IP, b[:net.IPv6len])
		b = b[net.IPv6len:]
	***REMOVED***
	return b, nil
***REMOVED***

func (ifi *InterfaceInfo) marshalName(proto int, b []byte) []byte ***REMOVED***
	l := byte(ifi.nameLen())
	b[0] = l
	copy(b[1:], []byte(ifi.Interface.Name))
	return b[l:]
***REMOVED***

func (ifi *InterfaceInfo) parseName(b []byte) ([]byte, error) ***REMOVED***
	if 4 > len(b) || len(b) < int(b[0]) ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	l := int(b[0])
	if l%4 != 0 || 4 > l || l > 64 ***REMOVED***
		return nil, errInvalidExtension
	***REMOVED***
	var name [63]byte
	copy(name[:], b[1:l])
	ifi.Interface.Name = strings.Trim(string(name[:]), "\000")
	return b[l:], nil
***REMOVED***

func (ifi *InterfaceInfo) marshalMTU(proto int, b []byte) []byte ***REMOVED***
	binary.BigEndian.PutUint32(b[:4], uint32(ifi.Interface.MTU))
	return b[4:]
***REMOVED***

func (ifi *InterfaceInfo) parseMTU(b []byte) ([]byte, error) ***REMOVED***
	if len(b) < 4 ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	ifi.Interface.MTU = int(binary.BigEndian.Uint32(b[:4]))
	return b[4:], nil
***REMOVED***

func parseInterfaceInfo(b []byte) (Extension, error) ***REMOVED***
	ifi := &InterfaceInfo***REMOVED***
		Class: int(b[2]),
		Type:  int(b[3]),
	***REMOVED***
	if ifi.Type&(attrIfIndex|attrName|attrMTU) != 0 ***REMOVED***
		ifi.Interface = &net.Interface***REMOVED******REMOVED***
	***REMOVED***
	if ifi.Type&attrIPAddr != 0 ***REMOVED***
		ifi.Addr = &net.IPAddr***REMOVED******REMOVED***
	***REMOVED***
	attrs := ifi.Type & (attrIfIndex | attrIPAddr | attrName | attrMTU)
	for b = b[4:]; len(b) > 0 && attrs != 0; ***REMOVED***
		var err error
		switch ***REMOVED***
		case attrs&attrIfIndex != 0:
			b, err = ifi.parseIfIndex(b)
			attrs &^= attrIfIndex
		case attrs&attrIPAddr != 0:
			b, err = ifi.parseIPAddr(b)
			attrs &^= attrIPAddr
		case attrs&attrName != 0:
			b, err = ifi.parseName(b)
			attrs &^= attrName
		case attrs&attrMTU != 0:
			b, err = ifi.parseMTU(b)
			attrs &^= attrMTU
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if ifi.Interface != nil && ifi.Interface.Name != "" && ifi.Addr != nil && ifi.Addr.IP.To16() != nil && ifi.Addr.IP.To4() == nil ***REMOVED***
		ifi.Addr.Zone = ifi.Interface.Name
	***REMOVED***
	return ifi, nil
***REMOVED***
