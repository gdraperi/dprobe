// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows
// +build go1.3

package svc

import "unsafe"

const ptrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const

// Should be a built-in for unsafe.Pointer?
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer ***REMOVED***
	return unsafe.Pointer(uintptr(p) + x)
***REMOVED***

// funcPC returns the entry PC of the function f.
// It assumes that f is a func value. Otherwise the behavior is undefined.
func funcPC(f interface***REMOVED******REMOVED***) uintptr ***REMOVED***
	return **(**uintptr)(add(unsafe.Pointer(&f), ptrSize))
***REMOVED***

// from sys_386.s and sys_amd64.s
func servicectlhandler(ctl uint32) uintptr
func servicemain(argc uint32, argv **uint16)

func getServiceMain(r *uintptr) ***REMOVED***
	*r = funcPC(servicemain)
***REMOVED***
