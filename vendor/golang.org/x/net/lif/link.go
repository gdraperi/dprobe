// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

package lif

import "unsafe"

// A Link represents logical data link information.
//
// It also represents base information for logical network interface.
// On Solaris, each logical network interface represents network layer
// adjacency information and the interface has a only single network
// address or address pair for tunneling. It's usual that multiple
// logical network interfaces share the same logical data link.
type Link struct ***REMOVED***
	Name  string // name, equivalent to IP interface name
	Index int    // index, equivalent to IP interface index
	Type  int    // type
	Flags int    // flags
	MTU   int    // maximum transmission unit, basically link MTU but may differ between IP address families
	Addr  []byte // address
***REMOVED***

func (ll *Link) fetch(s uintptr) ***REMOVED***
	var lifr lifreq
	for i := 0; i < len(ll.Name); i++ ***REMOVED***
		lifr.Name[i] = int8(ll.Name[i])
	***REMOVED***
	ioc := int64(sysSIOCGLIFINDEX)
	if err := ioctl(s, uintptr(ioc), unsafe.Pointer(&lifr)); err == nil ***REMOVED***
		ll.Index = int(nativeEndian.Uint32(lifr.Lifru[:4]))
	***REMOVED***
	ioc = int64(sysSIOCGLIFFLAGS)
	if err := ioctl(s, uintptr(ioc), unsafe.Pointer(&lifr)); err == nil ***REMOVED***
		ll.Flags = int(nativeEndian.Uint64(lifr.Lifru[:8]))
	***REMOVED***
	ioc = int64(sysSIOCGLIFMTU)
	if err := ioctl(s, uintptr(ioc), unsafe.Pointer(&lifr)); err == nil ***REMOVED***
		ll.MTU = int(nativeEndian.Uint32(lifr.Lifru[:4]))
	***REMOVED***
	switch ll.Type ***REMOVED***
	case sysIFT_IPV4, sysIFT_IPV6, sysIFT_6TO4:
	default:
		ioc = int64(sysSIOCGLIFHWADDR)
		if err := ioctl(s, uintptr(ioc), unsafe.Pointer(&lifr)); err == nil ***REMOVED***
			ll.Addr, _ = parseLinkAddr(lifr.Lifru[4:])
		***REMOVED***
	***REMOVED***
***REMOVED***

// Links returns a list of logical data links.
//
// The provided af must be an address family and name must be a data
// link name. The zero value of af or name means a wildcard.
func Links(af int, name string) ([]Link, error) ***REMOVED***
	eps, err := newEndpoints(af)
	if len(eps) == 0 ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		for _, ep := range eps ***REMOVED***
			ep.close()
		***REMOVED***
	***REMOVED***()
	return links(eps, name)
***REMOVED***

func links(eps []endpoint, name string) ([]Link, error) ***REMOVED***
	var lls []Link
	lifn := lifnum***REMOVED***Flags: sysLIFC_NOXMIT | sysLIFC_TEMPORARY | sysLIFC_ALLZONES | sysLIFC_UNDER_IPMP***REMOVED***
	lifc := lifconf***REMOVED***Flags: sysLIFC_NOXMIT | sysLIFC_TEMPORARY | sysLIFC_ALLZONES | sysLIFC_UNDER_IPMP***REMOVED***
	for _, ep := range eps ***REMOVED***
		lifn.Family = uint16(ep.af)
		ioc := int64(sysSIOCGLIFNUM)
		if err := ioctl(ep.s, uintptr(ioc), unsafe.Pointer(&lifn)); err != nil ***REMOVED***
			continue
		***REMOVED***
		if lifn.Count == 0 ***REMOVED***
			continue
		***REMOVED***
		b := make([]byte, lifn.Count*sizeofLifreq)
		lifc.Family = uint16(ep.af)
		lifc.Len = lifn.Count * sizeofLifreq
		if len(lifc.Lifcu) == 8 ***REMOVED***
			nativeEndian.PutUint64(lifc.Lifcu[:], uint64(uintptr(unsafe.Pointer(&b[0]))))
		***REMOVED*** else ***REMOVED***
			nativeEndian.PutUint32(lifc.Lifcu[:], uint32(uintptr(unsafe.Pointer(&b[0]))))
		***REMOVED***
		ioc = int64(sysSIOCGLIFCONF)
		if err := ioctl(ep.s, uintptr(ioc), unsafe.Pointer(&lifc)); err != nil ***REMOVED***
			continue
		***REMOVED***
		nb := make([]byte, 32) // see LIFNAMSIZ in net/if.h
		for i := 0; i < int(lifn.Count); i++ ***REMOVED***
			lifr := (*lifreq)(unsafe.Pointer(&b[i*sizeofLifreq]))
			for i := 0; i < 32; i++ ***REMOVED***
				if lifr.Name[i] == 0 ***REMOVED***
					nb = nb[:i]
					break
				***REMOVED***
				nb[i] = byte(lifr.Name[i])
			***REMOVED***
			llname := string(nb)
			nb = nb[:32]
			if isDupLink(lls, llname) || name != "" && name != llname ***REMOVED***
				continue
			***REMOVED***
			ll := Link***REMOVED***Name: llname, Type: int(lifr.Type)***REMOVED***
			ll.fetch(ep.s)
			lls = append(lls, ll)
		***REMOVED***
	***REMOVED***
	return lls, nil
***REMOVED***

func isDupLink(lls []Link, name string) bool ***REMOVED***
	for _, ll := range lls ***REMOVED***
		if ll.Name == name ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
