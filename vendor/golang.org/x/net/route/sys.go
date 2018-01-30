// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package route

import "unsafe"

var (
	nativeEndian binaryByteOrder
	kernelAlign  int
	wireFormats  map[int]*wireFormat
)

func init() ***REMOVED***
	i := uint32(1)
	b := (*[4]byte)(unsafe.Pointer(&i))
	if b[0] == 1 ***REMOVED***
		nativeEndian = littleEndian
	***REMOVED*** else ***REMOVED***
		nativeEndian = bigEndian
	***REMOVED***
	kernelAlign, wireFormats = probeRoutingStack()
***REMOVED***

func roundup(l int) int ***REMOVED***
	if l == 0 ***REMOVED***
		return kernelAlign
	***REMOVED***
	return (l + kernelAlign - 1) & ^(kernelAlign - 1)
***REMOVED***

type wireFormat struct ***REMOVED***
	extOff  int // offset of header extension
	bodyOff int // offset of message body
	parse   func(RIBType, []byte) (Message, error)
***REMOVED***
