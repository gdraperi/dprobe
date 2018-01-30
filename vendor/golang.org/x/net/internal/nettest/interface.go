// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nettest

import "net"

// IsMulticastCapable reports whether ifi is an IP multicast-capable
// network interface. Network must be "ip", "ip4" or "ip6".
func IsMulticastCapable(network string, ifi *net.Interface) (net.IP, bool) ***REMOVED***
	switch network ***REMOVED***
	case "ip", "ip4", "ip6":
	default:
		return nil, false
	***REMOVED***
	if ifi == nil || ifi.Flags&net.FlagUp == 0 || ifi.Flags&net.FlagMulticast == 0 ***REMOVED***
		return nil, false
	***REMOVED***
	return hasRoutableIP(network, ifi)
***REMOVED***

// RoutedInterface returns a network interface that can route IP
// traffic and satisfies flags. It returns nil when an appropriate
// network interface is not found. Network must be "ip", "ip4" or
// "ip6".
func RoutedInterface(network string, flags net.Flags) *net.Interface ***REMOVED***
	switch network ***REMOVED***
	case "ip", "ip4", "ip6":
	default:
		return nil
	***REMOVED***
	ift, err := net.Interfaces()
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	for _, ifi := range ift ***REMOVED***
		if ifi.Flags&flags != flags ***REMOVED***
			continue
		***REMOVED***
		if _, ok := hasRoutableIP(network, &ifi); !ok ***REMOVED***
			continue
		***REMOVED***
		return &ifi
	***REMOVED***
	return nil
***REMOVED***

func hasRoutableIP(network string, ifi *net.Interface) (net.IP, bool) ***REMOVED***
	ifat, err := ifi.Addrs()
	if err != nil ***REMOVED***
		return nil, false
	***REMOVED***
	for _, ifa := range ifat ***REMOVED***
		switch ifa := ifa.(type) ***REMOVED***
		case *net.IPAddr:
			if ip := routableIP(network, ifa.IP); ip != nil ***REMOVED***
				return ip, true
			***REMOVED***
		case *net.IPNet:
			if ip := routableIP(network, ifa.IP); ip != nil ***REMOVED***
				return ip, true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

func routableIP(network string, ip net.IP) net.IP ***REMOVED***
	if !ip.IsLoopback() && !ip.IsLinkLocalUnicast() && !ip.IsGlobalUnicast() ***REMOVED***
		return nil
	***REMOVED***
	switch network ***REMOVED***
	case "ip4":
		if ip := ip.To4(); ip != nil ***REMOVED***
			return ip
		***REMOVED***
	case "ip6":
		if ip.IsLoopback() ***REMOVED*** // addressing scope of the loopback address depends on each implementation
			return nil
		***REMOVED***
		if ip := ip.To16(); ip != nil && ip.To4() == nil ***REMOVED***
			return ip
		***REMOVED***
	default:
		if ip := ip.To4(); ip != nil ***REMOVED***
			return ip
		***REMOVED***
		if ip := ip.To16(); ip != nil ***REMOVED***
			return ip
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
