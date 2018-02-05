// Copyright 2016 The Gorilla WebSocket Authors. All rights reserved.  Use of
// this source code is governed by a BSD-style license that can be found in the
// LICENSE file.

// Require 1.7 for sub-bencmarks
// +build go1.7,!appengine

package websocket

import (
	"fmt"
	"testing"
)

func maskBytesByByte(key [4]byte, pos int, b []byte) int ***REMOVED***
	for i := range b ***REMOVED***
		b[i] ^= key[pos&3]
		pos++
	***REMOVED***
	return pos & 3
***REMOVED***

func notzero(b []byte) int ***REMOVED***
	for i := range b ***REMOVED***
		if b[i] != 0 ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

func TestMaskBytes(t *testing.T) ***REMOVED***
	key := [4]byte***REMOVED***1, 2, 3, 4***REMOVED***
	for size := 1; size <= 1024; size++ ***REMOVED***
		for align := 0; align < wordSize; align++ ***REMOVED***
			for pos := 0; pos < 4; pos++ ***REMOVED***
				b := make([]byte, size+align)[align:]
				maskBytes(key, pos, b)
				maskBytesByByte(key, pos, b)
				if i := notzero(b); i >= 0 ***REMOVED***
					t.Errorf("size:%d, align:%d, pos:%d, offset:%d", size, align, pos, i)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkMaskBytes(b *testing.B) ***REMOVED***
	for _, size := range []int***REMOVED***2, 4, 8, 16, 32, 512, 1024***REMOVED*** ***REMOVED***
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) ***REMOVED***
			for _, align := range []int***REMOVED***wordSize / 2***REMOVED*** ***REMOVED***
				b.Run(fmt.Sprintf("align-%d", align), func(b *testing.B) ***REMOVED***
					for _, fn := range []struct ***REMOVED***
						name string
						fn   func(key [4]byte, pos int, b []byte) int
					***REMOVED******REMOVED***
						***REMOVED***"byte", maskBytesByByte***REMOVED***,
						***REMOVED***"word", maskBytes***REMOVED***,
					***REMOVED*** ***REMOVED***
						b.Run(fn.name, func(b *testing.B) ***REMOVED***
							key := newMaskKey()
							data := make([]byte, size+align)[align:]
							for i := 0; i < b.N; i++ ***REMOVED***
								fn.fn(key, 0, data)
							***REMOVED***
							b.SetBytes(int64(len(data)))
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
