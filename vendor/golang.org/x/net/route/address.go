// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package route

import "runtime"

// An Addr represents an address associated with packet routing.
type Addr interface ***REMOVED***
	// Family returns an address family.
	Family() int
***REMOVED***

// A LinkAddr represents a link-layer address.
type LinkAddr struct ***REMOVED***
	Index int    // interface index when attached
	Name  string // interface name when attached
	Addr  []byte // link-layer address when attached
***REMOVED***

// Family implements the Family method of Addr interface.
func (a *LinkAddr) Family() int ***REMOVED*** return sysAF_LINK ***REMOVED***

func (a *LinkAddr) lenAndSpace() (int, int) ***REMOVED***
	l := 8 + len(a.Name) + len(a.Addr)
	return l, roundup(l)
***REMOVED***

func (a *LinkAddr) marshal(b []byte) (int, error) ***REMOVED***
	l, ll := a.lenAndSpace()
	if len(b) < ll ***REMOVED***
		return 0, errShortBuffer
	***REMOVED***
	nlen, alen := len(a.Name), len(a.Addr)
	if nlen > 255 || alen > 255 ***REMOVED***
		return 0, errInvalidAddr
	***REMOVED***
	b[0] = byte(l)
	b[1] = sysAF_LINK
	if a.Index > 0 ***REMOVED***
		nativeEndian.PutUint16(b[2:4], uint16(a.Index))
	***REMOVED***
	data := b[8:]
	if nlen > 0 ***REMOVED***
		b[5] = byte(nlen)
		copy(data[:nlen], a.Addr)
		data = data[nlen:]
	***REMOVED***
	if alen > 0 ***REMOVED***
		b[6] = byte(alen)
		copy(data[:alen], a.Name)
		data = data[alen:]
	***REMOVED***
	return ll, nil
***REMOVED***

func parseLinkAddr(b []byte) (Addr, error) ***REMOVED***
	if len(b) < 8 ***REMOVED***
		return nil, errInvalidAddr
	***REMOVED***
	_, a, err := parseKernelLinkAddr(sysAF_LINK, b[4:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	a.(*LinkAddr).Index = int(nativeEndian.Uint16(b[2:4]))
	return a, nil
***REMOVED***

// parseKernelLinkAddr parses b as a link-layer address in
// conventional BSD kernel form.
func parseKernelLinkAddr(_ int, b []byte) (int, Addr, error) ***REMOVED***
	// The encoding looks like the following:
	// +----------------------------+
	// | Type             (1 octet) |
	// +----------------------------+
	// | Name length      (1 octet) |
	// +----------------------------+
	// | Address length   (1 octet) |
	// +----------------------------+
	// | Selector length  (1 octet) |
	// +----------------------------+
	// | Data            (variable) |
	// +----------------------------+
	//
	// On some platforms, all-bit-one of length field means "don't
	// care".
	nlen, alen, slen := int(b[1]), int(b[2]), int(b[3])
	if nlen == 0xff ***REMOVED***
		nlen = 0
	***REMOVED***
	if alen == 0xff ***REMOVED***
		alen = 0
	***REMOVED***
	if slen == 0xff ***REMOVED***
		slen = 0
	***REMOVED***
	l := 4 + nlen + alen + slen
	if len(b) < l ***REMOVED***
		return 0, nil, errInvalidAddr
	***REMOVED***
	data := b[4:]
	var name string
	var addr []byte
	if nlen > 0 ***REMOVED***
		name = string(data[:nlen])
		data = data[nlen:]
	***REMOVED***
	if alen > 0 ***REMOVED***
		addr = data[:alen]
		data = data[alen:]
	***REMOVED***
	return l, &LinkAddr***REMOVED***Name: name, Addr: addr***REMOVED***, nil
***REMOVED***

// An Inet4Addr represents an internet address for IPv4.
type Inet4Addr struct ***REMOVED***
	IP [4]byte // IP address
***REMOVED***

// Family implements the Family method of Addr interface.
func (a *Inet4Addr) Family() int ***REMOVED*** return sysAF_INET ***REMOVED***

func (a *Inet4Addr) lenAndSpace() (int, int) ***REMOVED***
	return sizeofSockaddrInet, roundup(sizeofSockaddrInet)
***REMOVED***

func (a *Inet4Addr) marshal(b []byte) (int, error) ***REMOVED***
	l, ll := a.lenAndSpace()
	if len(b) < ll ***REMOVED***
		return 0, errShortBuffer
	***REMOVED***
	b[0] = byte(l)
	b[1] = sysAF_INET
	copy(b[4:8], a.IP[:])
	return ll, nil
***REMOVED***

// An Inet6Addr represents an internet address for IPv6.
type Inet6Addr struct ***REMOVED***
	IP     [16]byte // IP address
	ZoneID int      // zone identifier
***REMOVED***

// Family implements the Family method of Addr interface.
func (a *Inet6Addr) Family() int ***REMOVED*** return sysAF_INET6 ***REMOVED***

func (a *Inet6Addr) lenAndSpace() (int, int) ***REMOVED***
	return sizeofSockaddrInet6, roundup(sizeofSockaddrInet6)
***REMOVED***

func (a *Inet6Addr) marshal(b []byte) (int, error) ***REMOVED***
	l, ll := a.lenAndSpace()
	if len(b) < ll ***REMOVED***
		return 0, errShortBuffer
	***REMOVED***
	b[0] = byte(l)
	b[1] = sysAF_INET6
	copy(b[8:24], a.IP[:])
	if a.ZoneID > 0 ***REMOVED***
		nativeEndian.PutUint32(b[24:28], uint32(a.ZoneID))
	***REMOVED***
	return ll, nil
***REMOVED***

// parseInetAddr parses b as an internet address for IPv4 or IPv6.
func parseInetAddr(af int, b []byte) (Addr, error) ***REMOVED***
	switch af ***REMOVED***
	case sysAF_INET:
		if len(b) < sizeofSockaddrInet ***REMOVED***
			return nil, errInvalidAddr
		***REMOVED***
		a := &Inet4Addr***REMOVED******REMOVED***
		copy(a.IP[:], b[4:8])
		return a, nil
	case sysAF_INET6:
		if len(b) < sizeofSockaddrInet6 ***REMOVED***
			return nil, errInvalidAddr
		***REMOVED***
		a := &Inet6Addr***REMOVED***ZoneID: int(nativeEndian.Uint32(b[24:28]))***REMOVED***
		copy(a.IP[:], b[8:24])
		if a.IP[0] == 0xfe && a.IP[1]&0xc0 == 0x80 || a.IP[0] == 0xff && (a.IP[1]&0x0f == 0x01 || a.IP[1]&0x0f == 0x02) ***REMOVED***
			// KAME based IPv6 protocol stack usually
			// embeds the interface index in the
			// interface-local or link-local address as
			// the kernel-internal form.
			id := int(bigEndian.Uint16(a.IP[2:4]))
			if id != 0 ***REMOVED***
				a.ZoneID = id
				a.IP[2], a.IP[3] = 0, 0
			***REMOVED***
		***REMOVED***
		return a, nil
	default:
		return nil, errInvalidAddr
	***REMOVED***
***REMOVED***

// parseKernelInetAddr parses b as an internet address in conventional
// BSD kernel form.
func parseKernelInetAddr(af int, b []byte) (int, Addr, error) ***REMOVED***
	// The encoding looks similar to the NLRI encoding.
	// +----------------------------+
	// | Length           (1 octet) |
	// +----------------------------+
	// | Address prefix  (variable) |
	// +----------------------------+
	//
	// The differences between the kernel form and the NLRI
	// encoding are:
	//
	// - The length field of the kernel form indicates the prefix
	//   length in bytes, not in bits
	//
	// - In the kernel form, zero value of the length field
	//   doesn't mean 0.0.0.0/0 or ::/0
	//
	// - The kernel form appends leading bytes to the prefix field
	//   to make the <length, prefix> tuple to be conformed with
	//   the routing message boundary
	l := int(b[0])
	if runtime.GOOS == "darwin" ***REMOVED***
		// On Darwn, an address in the kernel form is also
		// used as a message filler.
		if l == 0 || len(b) > roundup(l) ***REMOVED***
			l = roundup(l)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		l = roundup(l)
	***REMOVED***
	if len(b) < l ***REMOVED***
		return 0, nil, errInvalidAddr
	***REMOVED***
	// Don't reorder case expressions.
	// The case expressions for IPv6 must come first.
	const (
		off4 = 4 // offset of in_addr
		off6 = 8 // offset of in6_addr
	)
	switch ***REMOVED***
	case b[0] == sizeofSockaddrInet6:
		a := &Inet6Addr***REMOVED******REMOVED***
		copy(a.IP[:], b[off6:off6+16])
		return int(b[0]), a, nil
	case af == sysAF_INET6:
		a := &Inet6Addr***REMOVED******REMOVED***
		if l-1 < off6 ***REMOVED***
			copy(a.IP[:], b[1:l])
		***REMOVED*** else ***REMOVED***
			copy(a.IP[:], b[l-off6:l])
		***REMOVED***
		return int(b[0]), a, nil
	case b[0] == sizeofSockaddrInet:
		a := &Inet4Addr***REMOVED******REMOVED***
		copy(a.IP[:], b[off4:off4+4])
		return int(b[0]), a, nil
	default: // an old fashion, AF_UNSPEC or unknown means AF_INET
		a := &Inet4Addr***REMOVED******REMOVED***
		if l-1 < off4 ***REMOVED***
			copy(a.IP[:], b[1:l])
		***REMOVED*** else ***REMOVED***
			copy(a.IP[:], b[l-off4:l])
		***REMOVED***
		return int(b[0]), a, nil
	***REMOVED***
***REMOVED***

// A DefaultAddr represents an address of various operating
// system-specific features.
type DefaultAddr struct ***REMOVED***
	af  int
	Raw []byte // raw format of address
***REMOVED***

// Family implements the Family method of Addr interface.
func (a *DefaultAddr) Family() int ***REMOVED*** return a.af ***REMOVED***

func (a *DefaultAddr) lenAndSpace() (int, int) ***REMOVED***
	l := len(a.Raw)
	return l, roundup(l)
***REMOVED***

func (a *DefaultAddr) marshal(b []byte) (int, error) ***REMOVED***
	l, ll := a.lenAndSpace()
	if len(b) < ll ***REMOVED***
		return 0, errShortBuffer
	***REMOVED***
	if l > 255 ***REMOVED***
		return 0, errInvalidAddr
	***REMOVED***
	b[1] = byte(l)
	copy(b[:l], a.Raw)
	return ll, nil
***REMOVED***

func parseDefaultAddr(b []byte) (Addr, error) ***REMOVED***
	if len(b) < 2 || len(b) < int(b[0]) ***REMOVED***
		return nil, errInvalidAddr
	***REMOVED***
	a := &DefaultAddr***REMOVED***af: int(b[1]), Raw: b[:b[0]]***REMOVED***
	return a, nil
***REMOVED***

func addrsSpace(as []Addr) int ***REMOVED***
	var l int
	for _, a := range as ***REMOVED***
		switch a := a.(type) ***REMOVED***
		case *LinkAddr:
			_, ll := a.lenAndSpace()
			l += ll
		case *Inet4Addr:
			_, ll := a.lenAndSpace()
			l += ll
		case *Inet6Addr:
			_, ll := a.lenAndSpace()
			l += ll
		case *DefaultAddr:
			_, ll := a.lenAndSpace()
			l += ll
		***REMOVED***
	***REMOVED***
	return l
***REMOVED***

// marshalAddrs marshals as and returns a bitmap indicating which
// address is stored in b.
func marshalAddrs(b []byte, as []Addr) (uint, error) ***REMOVED***
	var attrs uint
	for i, a := range as ***REMOVED***
		switch a := a.(type) ***REMOVED***
		case *LinkAddr:
			l, err := a.marshal(b)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			b = b[l:]
			attrs |= 1 << uint(i)
		case *Inet4Addr:
			l, err := a.marshal(b)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			b = b[l:]
			attrs |= 1 << uint(i)
		case *Inet6Addr:
			l, err := a.marshal(b)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			b = b[l:]
			attrs |= 1 << uint(i)
		case *DefaultAddr:
			l, err := a.marshal(b)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			b = b[l:]
			attrs |= 1 << uint(i)
		***REMOVED***
	***REMOVED***
	return attrs, nil
***REMOVED***

func parseAddrs(attrs uint, fn func(int, []byte) (int, Addr, error), b []byte) ([]Addr, error) ***REMOVED***
	var as [sysRTAX_MAX]Addr
	af := int(sysAF_UNSPEC)
	for i := uint(0); i < sysRTAX_MAX && len(b) >= roundup(0); i++ ***REMOVED***
		if attrs&(1<<i) == 0 ***REMOVED***
			continue
		***REMOVED***
		if i <= sysRTAX_BRD ***REMOVED***
			switch b[1] ***REMOVED***
			case sysAF_LINK:
				a, err := parseLinkAddr(b)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				as[i] = a
				l := roundup(int(b[0]))
				if len(b) < l ***REMOVED***
					return nil, errMessageTooShort
				***REMOVED***
				b = b[l:]
			case sysAF_INET, sysAF_INET6:
				af = int(b[1])
				a, err := parseInetAddr(af, b)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				as[i] = a
				l := roundup(int(b[0]))
				if len(b) < l ***REMOVED***
					return nil, errMessageTooShort
				***REMOVED***
				b = b[l:]
			default:
				l, a, err := fn(af, b)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				as[i] = a
				ll := roundup(l)
				if len(b) < ll ***REMOVED***
					b = b[l:]
				***REMOVED*** else ***REMOVED***
					b = b[ll:]
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			a, err := parseDefaultAddr(b)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			as[i] = a
			l := roundup(int(b[0]))
			if len(b) < l ***REMOVED***
				return nil, errMessageTooShort
			***REMOVED***
			b = b[l:]
		***REMOVED***
	***REMOVED***
	return as[:], nil
***REMOVED***
