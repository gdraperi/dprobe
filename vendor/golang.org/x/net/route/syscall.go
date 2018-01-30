// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package route

import (
	"syscall"
	"unsafe"
)

var zero uintptr

func sysctl(mib []int32, old *byte, oldlen *uintptr, new *byte, newlen uintptr) error ***REMOVED***
	var p unsafe.Pointer
	if len(mib) > 0 ***REMOVED***
		p = unsafe.Pointer(&mib[0])
	***REMOVED*** else ***REMOVED***
		p = unsafe.Pointer(&zero)
	***REMOVED***
	_, _, errno := syscall.Syscall6(syscall.SYS___SYSCTL, uintptr(p), uintptr(len(mib)), uintptr(unsafe.Pointer(old)), uintptr(unsafe.Pointer(oldlen)), uintptr(unsafe.Pointer(new)), uintptr(newlen))
	if errno != 0 ***REMOVED***
		return error(errno)
	***REMOVED***
	return nil
***REMOVED***
