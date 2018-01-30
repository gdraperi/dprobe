// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

import (
	"errors"
	"net"
)

var (
	errMissingAddress  = errors.New("missing address")
	errHeaderTooShort  = errors.New("header too short")
	errInvalidConnType = errors.New("invalid conn type")
	errOpNoSupport     = errors.New("operation not supported")
	errNoSuchInterface = errors.New("no such interface")
)

func boolint(b bool) int ***REMOVED***
	if b ***REMOVED***
		return 1
	***REMOVED***
	return 0
***REMOVED***

func netAddrToIP16(a net.Addr) net.IP ***REMOVED***
	switch v := a.(type) ***REMOVED***
	case *net.UDPAddr:
		if ip := v.IP.To16(); ip != nil && ip.To4() == nil ***REMOVED***
			return ip
		***REMOVED***
	case *net.IPAddr:
		if ip := v.IP.To16(); ip != nil && ip.To4() == nil ***REMOVED***
			return ip
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func opAddr(a net.Addr) net.Addr ***REMOVED***
	switch a.(type) ***REMOVED***
	case *net.TCPAddr:
		if a == nil ***REMOVED***
			return nil
		***REMOVED***
	case *net.UDPAddr:
		if a == nil ***REMOVED***
			return nil
		***REMOVED***
	case *net.IPAddr:
		if a == nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return a
***REMOVED***
