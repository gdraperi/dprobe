// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package socket

import (
	"errors"
	"net"
	"runtime"
	"unsafe"
)

const (
	sysAF_UNSPEC = 0x0
	sysAF_INET   = 0x2
	sysAF_INET6  = 0xa

	sysSOCK_RAW = 0x3
)

func probeProtocolStack() int ***REMOVED***
	switch runtime.GOARCH ***REMOVED***
	case "amd64p32", "mips64p32":
		return 4
	default:
		var p uintptr
		return int(unsafe.Sizeof(p))
	***REMOVED***
***REMOVED***

func marshalInetAddr(ip net.IP, port int, zone string) []byte ***REMOVED***
	return nil
***REMOVED***

func parseInetAddr(b []byte, network string) (net.Addr, error) ***REMOVED***
	return nil, errors.New("not implemented")
***REMOVED***

func getsockopt(s uintptr, level, name int, b []byte) (int, error) ***REMOVED***
	return 0, errors.New("not implemented")
***REMOVED***

func setsockopt(s uintptr, level, name int, b []byte) error ***REMOVED***
	return errors.New("not implemented")
***REMOVED***

func recvmsg(s uintptr, h *msghdr, flags int) (int, error) ***REMOVED***
	return 0, errors.New("not implemented")
***REMOVED***

func sendmsg(s uintptr, h *msghdr, flags int) (int, error) ***REMOVED***
	return 0, errors.New("not implemented")
***REMOVED***

func recvmmsg(s uintptr, hs []mmsghdr, flags int) (int, error) ***REMOVED***
	return 0, errors.New("not implemented")
***REMOVED***

func sendmmsg(s uintptr, hs []mmsghdr, flags int) (int, error) ***REMOVED***
	return 0, errors.New("not implemented")
***REMOVED***
