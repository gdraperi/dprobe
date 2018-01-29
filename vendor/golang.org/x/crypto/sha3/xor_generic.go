// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sha3

import "encoding/binary"

// xorInGeneric xors the bytes in buf into the state; it
// makes no non-portable assumptions about memory layout
// or alignment.
func xorInGeneric(d *state, buf []byte) ***REMOVED***
	n := len(buf) / 8

	for i := 0; i < n; i++ ***REMOVED***
		a := binary.LittleEndian.Uint64(buf)
		d.a[i] ^= a
		buf = buf[8:]
	***REMOVED***
***REMOVED***

// copyOutGeneric copies ulint64s to a byte buffer.
func copyOutGeneric(d *state, b []byte) ***REMOVED***
	for i := 0; len(b) >= 8; i++ ***REMOVED***
		binary.LittleEndian.PutUint64(b, d.a[i])
		b = b[8:]
	***REMOVED***
***REMOVED***
