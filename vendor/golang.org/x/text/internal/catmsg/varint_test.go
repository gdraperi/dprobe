// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package catmsg

import (
	"fmt"
	"testing"
)

func TestEncodeUint(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		x   uint64
		enc string
	***REMOVED******REMOVED***
		***REMOVED***0, "\x00"***REMOVED***,
		***REMOVED***1, "\x01"***REMOVED***,
		***REMOVED***2, "\x02"***REMOVED***,
		***REMOVED***0x7f, "\x7f"***REMOVED***,
		***REMOVED***0x80, "\x80\x01"***REMOVED***,
		***REMOVED***1 << 14, "\x80\x80\x01"***REMOVED***,
		***REMOVED***0xffffffff, "\xff\xff\xff\xff\x0f"***REMOVED***,
		***REMOVED***0xffffffffffffffff, "\xff\xff\xff\xff\xff\xff\xff\xff\xff\x01"***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		buf := [maxVarintBytes]byte***REMOVED******REMOVED***
		got := string(buf[:encodeUint(buf[:], tc.x)])
		if got != tc.enc ***REMOVED***
			t.Errorf("EncodeUint(%#x) = %q; want %q", tc.x, got, tc.enc)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDecodeUint(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		x    uint64
		size int
		enc  string
		err  error
	***REMOVED******REMOVED******REMOVED***
		x:    0,
		size: 0,
		enc:  "",
		err:  errIllegalVarint,
	***REMOVED***, ***REMOVED***
		x:    0,
		size: 1,
		enc:  "\x80",
		err:  errIllegalVarint,
	***REMOVED***, ***REMOVED***
		x:    0,
		size: 3,
		enc:  "\x80\x80\x80",
		err:  errIllegalVarint,
	***REMOVED***, ***REMOVED***
		x:    0,
		size: 1,
		enc:  "\x00",
	***REMOVED***, ***REMOVED***
		x:    1,
		size: 1,
		enc:  "\x01",
	***REMOVED***, ***REMOVED***
		x:    2,
		size: 1,
		enc:  "\x02",
	***REMOVED***, ***REMOVED***
		x:    0x7f,
		size: 1,
		enc:  "\x7f",
	***REMOVED***, ***REMOVED***
		x:    0x80,
		size: 2,
		enc:  "\x80\x01",
	***REMOVED***, ***REMOVED***
		x:    1 << 14,
		size: 3,
		enc:  "\x80\x80\x01",
	***REMOVED***, ***REMOVED***
		x:    0xffffffff,
		size: 5,
		enc:  "\xff\xff\xff\xff\x0f",
	***REMOVED***, ***REMOVED***
		x:    0xffffffffffffffff,
		size: 10,
		enc:  "\xff\xff\xff\xff\xff\xff\xff\xff\xff\x01",
	***REMOVED***, ***REMOVED***
		x:    0xffffffffffffffff,
		size: 10,
		enc:  "\xff\xff\xff\xff\xff\xff\xff\xff\xff\x01\x00",
	***REMOVED***, ***REMOVED***
		x:    0,
		size: 10,
		enc:  "\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\x01",
		err:  errVarintTooLarge,
	***REMOVED******REMOVED***
	forms := []struct ***REMOVED***
		name   string
		decode func(s string) (x uint64, size int, err error)
	***REMOVED******REMOVED***
		***REMOVED***"decode", func(s string) (x uint64, size int, err error) ***REMOVED***
			return decodeUint([]byte(s))
		***REMOVED******REMOVED***,
		***REMOVED***"decodeString", decodeUintString***REMOVED***,
	***REMOVED***
	for _, f := range forms ***REMOVED***
		for _, tc := range testCases ***REMOVED***
			t.Run(fmt.Sprintf("%s:%q", f.name, tc.enc), func(t *testing.T) ***REMOVED***
				x, size, err := f.decode(tc.enc)
				if err != tc.err ***REMOVED***
					t.Errorf("err = %q; want %q", err, tc.err)
				***REMOVED***
				if size != tc.size ***REMOVED***
					t.Errorf("size = %d; want %d", size, tc.size)
				***REMOVED***
				if x != tc.x ***REMOVED***
					t.Errorf("decode = %#x; want %#x", x, tc.x)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***
