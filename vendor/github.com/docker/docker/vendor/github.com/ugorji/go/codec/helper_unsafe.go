//+build unsafe

// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"unsafe"
)

// This file has unsafe variants of some helper methods.

type unsafeString struct ***REMOVED***
	Data uintptr
	Len  int
***REMOVED***

type unsafeBytes struct ***REMOVED***
	Data uintptr
	Len  int
	Cap  int
***REMOVED***

// stringView returns a view of the []byte as a string.
// In unsafe mode, it doesn't incur allocation and copying caused by conversion.
// In regular safe mode, it is an allocation and copy.
func stringView(v []byte) string ***REMOVED***
	if len(v) == 0 ***REMOVED***
		return ""
	***REMOVED***
	x := unsafeString***REMOVED***uintptr(unsafe.Pointer(&v[0])), len(v)***REMOVED***
	return *(*string)(unsafe.Pointer(&x))
***REMOVED***

// bytesView returns a view of the string as a []byte.
// In unsafe mode, it doesn't incur allocation and copying caused by conversion.
// In regular safe mode, it is an allocation and copy.
func bytesView(v string) []byte ***REMOVED***
	if len(v) == 0 ***REMOVED***
		return zeroByteSlice
	***REMOVED***
	x := unsafeBytes***REMOVED***uintptr(unsafe.Pointer(&v)), len(v), len(v)***REMOVED***
	return *(*[]byte)(unsafe.Pointer(&x))
***REMOVED***
