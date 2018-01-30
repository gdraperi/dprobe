// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9
// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package socket

import (
	"encoding/binary"
	"errors"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func marshalInetAddr(a net.Addr) []byte ***REMOVED***
	switch a := a.(type) ***REMOVED***
	case *net.TCPAddr:
		return marshalSockaddr(a.IP, a.Port, a.Zone)
	case *net.UDPAddr:
		return marshalSockaddr(a.IP, a.Port, a.Zone)
	case *net.IPAddr:
		return marshalSockaddr(a.IP, 0, a.Zone)
	default:
		return nil
	***REMOVED***
***REMOVED***

func marshalSockaddr(ip net.IP, port int, zone string) []byte ***REMOVED***
	if ip4 := ip.To4(); ip4 != nil ***REMOVED***
		b := make([]byte, sizeofSockaddrInet)
		switch runtime.GOOS ***REMOVED***
		case "android", "linux", "solaris", "windows":
			NativeEndian.PutUint16(b[:2], uint16(sysAF_INET))
		default:
			b[0] = sizeofSockaddrInet
			b[1] = sysAF_INET
		***REMOVED***
		binary.BigEndian.PutUint16(b[2:4], uint16(port))
		copy(b[4:8], ip4)
		return b
	***REMOVED***
	if ip6 := ip.To16(); ip6 != nil && ip.To4() == nil ***REMOVED***
		b := make([]byte, sizeofSockaddrInet6)
		switch runtime.GOOS ***REMOVED***
		case "android", "linux", "solaris", "windows":
			NativeEndian.PutUint16(b[:2], uint16(sysAF_INET6))
		default:
			b[0] = sizeofSockaddrInet6
			b[1] = sysAF_INET6
		***REMOVED***
		binary.BigEndian.PutUint16(b[2:4], uint16(port))
		copy(b[8:24], ip6)
		if zone != "" ***REMOVED***
			NativeEndian.PutUint32(b[24:28], uint32(zoneCache.index(zone)))
		***REMOVED***
		return b
	***REMOVED***
	return nil
***REMOVED***

func parseInetAddr(b []byte, network string) (net.Addr, error) ***REMOVED***
	if len(b) < 2 ***REMOVED***
		return nil, errors.New("invalid address")
	***REMOVED***
	var af int
	switch runtime.GOOS ***REMOVED***
	case "android", "linux", "solaris", "windows":
		af = int(NativeEndian.Uint16(b[:2]))
	default:
		af = int(b[1])
	***REMOVED***
	var ip net.IP
	var zone string
	if af == sysAF_INET ***REMOVED***
		if len(b) < sizeofSockaddrInet ***REMOVED***
			return nil, errors.New("short address")
		***REMOVED***
		ip = make(net.IP, net.IPv4len)
		copy(ip, b[4:8])
	***REMOVED***
	if af == sysAF_INET6 ***REMOVED***
		if len(b) < sizeofSockaddrInet6 ***REMOVED***
			return nil, errors.New("short address")
		***REMOVED***
		ip = make(net.IP, net.IPv6len)
		copy(ip, b[8:24])
		if id := int(NativeEndian.Uint32(b[24:28])); id > 0 ***REMOVED***
			zone = zoneCache.name(id)
		***REMOVED***
	***REMOVED***
	switch network ***REMOVED***
	case "tcp", "tcp4", "tcp6":
		return &net.TCPAddr***REMOVED***IP: ip, Port: int(binary.BigEndian.Uint16(b[2:4])), Zone: zone***REMOVED***, nil
	case "udp", "udp4", "udp6":
		return &net.UDPAddr***REMOVED***IP: ip, Port: int(binary.BigEndian.Uint16(b[2:4])), Zone: zone***REMOVED***, nil
	default:
		return &net.IPAddr***REMOVED***IP: ip, Zone: zone***REMOVED***, nil
	***REMOVED***
***REMOVED***

// An ipv6ZoneCache represents a cache holding partial network
// interface information. It is used for reducing the cost of IPv6
// addressing scope zone resolution.
//
// Multiple names sharing the index are managed by first-come
// first-served basis for consistency.
type ipv6ZoneCache struct ***REMOVED***
	sync.RWMutex                // guard the following
	lastFetched  time.Time      // last time routing information was fetched
	toIndex      map[string]int // interface name to its index
	toName       map[int]string // interface index to its name
***REMOVED***

var zoneCache = ipv6ZoneCache***REMOVED***
	toIndex: make(map[string]int),
	toName:  make(map[int]string),
***REMOVED***

func (zc *ipv6ZoneCache) update(ift []net.Interface) ***REMOVED***
	zc.Lock()
	defer zc.Unlock()
	now := time.Now()
	if zc.lastFetched.After(now.Add(-60 * time.Second)) ***REMOVED***
		return
	***REMOVED***
	zc.lastFetched = now
	if len(ift) == 0 ***REMOVED***
		var err error
		if ift, err = net.Interfaces(); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	zc.toIndex = make(map[string]int, len(ift))
	zc.toName = make(map[int]string, len(ift))
	for _, ifi := range ift ***REMOVED***
		zc.toIndex[ifi.Name] = ifi.Index
		if _, ok := zc.toName[ifi.Index]; !ok ***REMOVED***
			zc.toName[ifi.Index] = ifi.Name
		***REMOVED***
	***REMOVED***
***REMOVED***

func (zc *ipv6ZoneCache) name(zone int) string ***REMOVED***
	zoneCache.update(nil)
	zoneCache.RLock()
	defer zoneCache.RUnlock()
	name, ok := zoneCache.toName[zone]
	if !ok ***REMOVED***
		name = strconv.Itoa(zone)
	***REMOVED***
	return name
***REMOVED***

func (zc *ipv6ZoneCache) index(zone string) int ***REMOVED***
	zoneCache.update(nil)
	zoneCache.RLock()
	defer zoneCache.RUnlock()
	index, ok := zoneCache.toIndex[zone]
	if !ok ***REMOVED***
		index, _ = strconv.Atoi(zone)
	***REMOVED***
	return index
***REMOVED***
