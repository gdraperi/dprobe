// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import (
	"errors"
	"syscall"
	"unsafe"
)

func probeProtocolStack() int ***REMOVED***
	var p uintptr
	return int(unsafe.Sizeof(p))
***REMOVED***

const (
	sysAF_UNSPEC = 0x0
	sysAF_INET   = 0x2
	sysAF_INET6  = 0x17

	sysSOCK_RAW = 0x3
)

type sockaddrInet struct ***REMOVED***
	Family uint16
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]uint8
***REMOVED***

type sockaddrInet6 struct ***REMOVED***
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
***REMOVED***

const (
	sizeofSockaddrInet  = 0x10
	sizeofSockaddrInet6 = 0x1c
)

func getsockopt(s uintptr, level, name int, b []byte) (int, error) ***REMOVED***
	l := uint32(len(b))
	err := syscall.Getsockopt(syscall.Handle(s), int32(level), int32(name), (*byte)(unsafe.Pointer(&b[0])), (*int32)(unsafe.Pointer(&l)))
	return int(l), err
***REMOVED***

func setsockopt(s uintptr, level, name int, b []byte) error ***REMOVED***
	return syscall.Setsockopt(syscall.Handle(s), int32(level), int32(name), (*byte)(unsafe.Pointer(&b[0])), int32(len(b)))
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
