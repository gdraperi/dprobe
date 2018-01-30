// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

package lif

import (
	"errors"
	"unsafe"
)

// An Addr represents an address associated with packet routing.
type Addr interface ***REMOVED***
	// Family returns an address family.
	Family() int
***REMOVED***

// An Inet4Addr represents an internet address for IPv4.
type Inet4Addr struct ***REMOVED***
	IP        [4]byte // IP address
	PrefixLen int     // address prefix length
***REMOVED***

// Family implements the Family method of Addr interface.
func (a *Inet4Addr) Family() int ***REMOVED*** return sysAF_INET ***REMOVED***

// An Inet6Addr represents an internet address for IPv6.
type Inet6Addr struct ***REMOVED***
	IP        [16]byte // IP address
	PrefixLen int      // address prefix length
	ZoneID    int      // zone identifier
***REMOVED***

// Family implements the Family method of Addr interface.
func (a *Inet6Addr) Family() int ***REMOVED*** return sysAF_INET6 ***REMOVED***

// Addrs returns a list of interface addresses.
//
// The provided af must be an address family and name must be a data
// link name. The zero value of af or name means a wildcard.
func Addrs(af int, name string) ([]Addr, error) ***REMOVED***
	eps, err := newEndpoints(af)
	if len(eps) == 0 ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		for _, ep := range eps ***REMOVED***
			ep.close()
		***REMOVED***
	***REMOVED***()
	lls, err := links(eps, name)
	if len(lls) == 0 ***REMOVED***
		return nil, err
	***REMOVED***
	var as []Addr
	for _, ll := range lls ***REMOVED***
		var lifr lifreq
		for i := 0; i < len(ll.Name); i++ ***REMOVED***
			lifr.Name[i] = int8(ll.Name[i])
		***REMOVED***
		for _, ep := range eps ***REMOVED***
			ioc := int64(sysSIOCGLIFADDR)
			err := ioctl(ep.s, uintptr(ioc), unsafe.Pointer(&lifr))
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			sa := (*sockaddrStorage)(unsafe.Pointer(&lifr.Lifru[0]))
			l := int(nativeEndian.Uint32(lifr.Lifru1[:4]))
			if l == 0 ***REMOVED***
				continue
			***REMOVED***
			switch sa.Family ***REMOVED***
			case sysAF_INET:
				a := &Inet4Addr***REMOVED***PrefixLen: l***REMOVED***
				copy(a.IP[:], lifr.Lifru[4:8])
				as = append(as, a)
			case sysAF_INET6:
				a := &Inet6Addr***REMOVED***PrefixLen: l, ZoneID: int(nativeEndian.Uint32(lifr.Lifru[24:28]))***REMOVED***
				copy(a.IP[:], lifr.Lifru[8:24])
				as = append(as, a)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return as, nil
***REMOVED***

func parseLinkAddr(b []byte) ([]byte, error) ***REMOVED***
	nlen, alen, slen := int(b[1]), int(b[2]), int(b[3])
	l := 4 + nlen + alen + slen
	if len(b) < l ***REMOVED***
		return nil, errors.New("invalid address")
	***REMOVED***
	b = b[4:]
	var addr []byte
	if nlen > 0 ***REMOVED***
		b = b[nlen:]
	***REMOVED***
	if alen > 0 ***REMOVED***
		addr = make([]byte, alen)
		copy(addr, b[:alen])
	***REMOVED***
	return addr, nil
***REMOVED***
