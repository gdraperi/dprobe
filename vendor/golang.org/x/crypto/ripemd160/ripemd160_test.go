// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ripemd160

// Test vectors are from:
// http://homes.esat.kuleuven.be/~bosselae/ripemd160.html

import (
	"fmt"
	"io"
	"testing"
)

type mdTest struct ***REMOVED***
	out string
	in  string
***REMOVED***

var vectors = [...]mdTest***REMOVED***
	***REMOVED***"9c1185a5c5e9fc54612808977ee8f548b2258d31", ""***REMOVED***,
	***REMOVED***"0bdc9d2d256b3ee9daae347be6f4dc835a467ffe", "a"***REMOVED***,
	***REMOVED***"8eb208f7e05d987a9b044a8e98c6b087f15a0bfc", "abc"***REMOVED***,
	***REMOVED***"5d0689ef49d2fae572b881b123a85ffa21595f36", "message digest"***REMOVED***,
	***REMOVED***"f71c27109c692c1b56bbdceb5b9d2865b3708dbc", "abcdefghijklmnopqrstuvwxyz"***REMOVED***,
	***REMOVED***"12a053384a9c0c88e405a06c27dcf49ada62eb2b", "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq"***REMOVED***,
	***REMOVED***"b0e20b6e3116640286ed3a87a5713079b21f5189", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"***REMOVED***,
	***REMOVED***"9b752e45573d4b39f4dbd3323cab82bf63326bfb", "12345678901234567890123456789012345678901234567890123456789012345678901234567890"***REMOVED***,
***REMOVED***

func TestVectors(t *testing.T) ***REMOVED***
	for i := 0; i < len(vectors); i++ ***REMOVED***
		tv := vectors[i]
		md := New()
		for j := 0; j < 3; j++ ***REMOVED***
			if j < 2 ***REMOVED***
				io.WriteString(md, tv.in)
			***REMOVED*** else ***REMOVED***
				io.WriteString(md, tv.in[0:len(tv.in)/2])
				md.Sum(nil)
				io.WriteString(md, tv.in[len(tv.in)/2:])
			***REMOVED***
			s := fmt.Sprintf("%x", md.Sum(nil))
			if s != tv.out ***REMOVED***
				t.Fatalf("RIPEMD-160[%d](%s) = %s, expected %s", j, tv.in, s, tv.out)
			***REMOVED***
			md.Reset()
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMillionA(t *testing.T) ***REMOVED***
	md := New()
	for i := 0; i < 100000; i++ ***REMOVED***
		io.WriteString(md, "aaaaaaaaaa")
	***REMOVED***
	out := "52783243c1697bdbe16d37f97f68f08325dc1528"
	s := fmt.Sprintf("%x", md.Sum(nil))
	if s != out ***REMOVED***
		t.Fatalf("RIPEMD-160 (1 million 'a') = %s, expected %s", s, out)
	***REMOVED***
	md.Reset()
***REMOVED***
