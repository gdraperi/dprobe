// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import (
	"encoding/binary"
	"unsafe"
)

var (
	// NativeEndian is the machine native endian implementation of
	// ByteOrder.
	NativeEndian binary.ByteOrder

	kernelAlign int
)

func init() ***REMOVED***
	i := uint32(1)
	b := (*[4]byte)(unsafe.Pointer(&i))
	if b[0] == 1 ***REMOVED***
		NativeEndian = binary.LittleEndian
	***REMOVED*** else ***REMOVED***
		NativeEndian = binary.BigEndian
	***REMOVED***
	kernelAlign = probeProtocolStack()
***REMOVED***

func roundup(l int) int ***REMOVED***
	return (l + kernelAlign - 1) & ^(kernelAlign - 1)
***REMOVED***
