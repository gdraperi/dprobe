// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package poly1305

import (
	"bytes"
	"encoding/hex"
	"flag"
	"testing"
	"unsafe"
)

var stressFlag = flag.Bool("stress", false, "run slow stress tests")

var testData = []struct ***REMOVED***
	in, k, correct []byte
***REMOVED******REMOVED***
	***REMOVED***
		[]byte("Hello world!"),
		[]byte("this is 32-byte key for Poly1305"),
		[]byte***REMOVED***0xa6, 0xf7, 0x45, 0x00, 0x8f, 0x81, 0xc9, 0x16, 0xa2, 0x0d, 0xcc, 0x74, 0xee, 0xf2, 0xb2, 0xf0***REMOVED***,
	***REMOVED***,
	***REMOVED***
		make([]byte, 32),
		[]byte("this is 32-byte key for Poly1305"),
		[]byte***REMOVED***0x49, 0xec, 0x78, 0x09, 0x0e, 0x48, 0x1e, 0xc6, 0xc2, 0x6b, 0x33, 0xb9, 0x1c, 0xcc, 0x03, 0x07***REMOVED***,
	***REMOVED***,
	***REMOVED***
		make([]byte, 2007),
		[]byte("this is 32-byte key for Poly1305"),
		[]byte***REMOVED***0xda, 0x84, 0xbc, 0xab, 0x02, 0x67, 0x6c, 0x38, 0xcd, 0xb0, 0x15, 0x60, 0x42, 0x74, 0xc2, 0xaa***REMOVED***,
	***REMOVED***,
	***REMOVED***
		make([]byte, 2007),
		make([]byte, 32),
		make([]byte, 16),
	***REMOVED***,
	***REMOVED***
		// This test triggers an edge-case. See https://go-review.googlesource.com/#/c/30101/.
		[]byte***REMOVED***0x81, 0xd8, 0xb2, 0xe4, 0x6a, 0x25, 0x21, 0x3b, 0x58, 0xfe, 0xe4, 0x21, 0x3a, 0x2a, 0x28, 0xe9, 0x21, 0xc1, 0x2a, 0x96, 0x32, 0x51, 0x6d, 0x3b, 0x73, 0x27, 0x27, 0x27, 0xbe, 0xcf, 0x21, 0x29***REMOVED***,
		[]byte***REMOVED***0x3b, 0x3a, 0x29, 0xe9, 0x3b, 0x21, 0x3a, 0x5c, 0x5c, 0x3b, 0x3b, 0x05, 0x3a, 0x3a, 0x8c, 0x0d***REMOVED***,
		[]byte***REMOVED***0x6d, 0xc1, 0x8b, 0x8c, 0x34, 0x4c, 0xd7, 0x99, 0x27, 0x11, 0x8b, 0xbe, 0x84, 0xb7, 0xf3, 0x14***REMOVED***,
	***REMOVED***,
	***REMOVED***
		// This test generates a result of (2^130-1) % (2^130-5).
		[]byte***REMOVED***
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		***REMOVED***,
		[]byte***REMOVED***1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0***REMOVED***,
		[]byte***REMOVED***4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0***REMOVED***,
	***REMOVED***,
	***REMOVED***
		// This test generates a result of (2^130-6) % (2^130-5).
		[]byte***REMOVED***
			0xfa, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		***REMOVED***,
		[]byte***REMOVED***1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0***REMOVED***,
		[]byte***REMOVED***0xfa, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff***REMOVED***,
	***REMOVED***,
	***REMOVED***
		// This test generates a result of (2^130-5) % (2^130-5).
		[]byte***REMOVED***
			0xfb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		***REMOVED***,
		[]byte***REMOVED***1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0***REMOVED***,
		[]byte***REMOVED***0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0***REMOVED***,
	***REMOVED***,
***REMOVED***

func testSum(t *testing.T, unaligned bool) ***REMOVED***
	var out [16]byte
	var key [32]byte

	for i, v := range testData ***REMOVED***
		in := v.in
		if unaligned ***REMOVED***
			in = unalignBytes(in)
		***REMOVED***
		copy(key[:], v.k)
		Sum(&out, in, &key)
		if !bytes.Equal(out[:], v.correct) ***REMOVED***
			t.Errorf("%d: expected %x, got %x", i, v.correct, out[:])
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBurnin(t *testing.T) ***REMOVED***
	// This test can be used to sanity-check significant changes. It can
	// take about many minutes to run, even on fast machines. It's disabled
	// by default.
	if !*stressFlag ***REMOVED***
		t.Skip("skipping without -stress")
	***REMOVED***

	var key [32]byte
	var input [25]byte
	var output [16]byte

	for i := range key ***REMOVED***
		key[i] = 1
	***REMOVED***
	for i := range input ***REMOVED***
		input[i] = 2
	***REMOVED***

	for i := uint64(0); i < 1e10; i++ ***REMOVED***
		Sum(&output, input[:], &key)
		copy(key[0:], output[:])
		copy(key[16:], output[:])
		copy(input[:], output[:])
		copy(input[16:], output[:])
	***REMOVED***

	const expected = "5e3b866aea0b636d240c83c428f84bfa"
	if got := hex.EncodeToString(output[:]); got != expected ***REMOVED***
		t.Errorf("expected %s, got %s", expected, got)
	***REMOVED***
***REMOVED***

func TestSum(t *testing.T)          ***REMOVED*** testSum(t, false) ***REMOVED***
func TestSumUnaligned(t *testing.T) ***REMOVED*** testSum(t, true) ***REMOVED***

func benchmark(b *testing.B, size int, unaligned bool) ***REMOVED***
	var out [16]byte
	var key [32]byte
	in := make([]byte, size)
	if unaligned ***REMOVED***
		in = unalignBytes(in)
	***REMOVED***
	b.SetBytes(int64(len(in)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		Sum(&out, in, &key)
	***REMOVED***
***REMOVED***

func Benchmark64(b *testing.B)          ***REMOVED*** benchmark(b, 64, false) ***REMOVED***
func Benchmark1K(b *testing.B)          ***REMOVED*** benchmark(b, 1024, false) ***REMOVED***
func Benchmark64Unaligned(b *testing.B) ***REMOVED*** benchmark(b, 64, true) ***REMOVED***
func Benchmark1KUnaligned(b *testing.B) ***REMOVED*** benchmark(b, 1024, true) ***REMOVED***

func unalignBytes(in []byte) []byte ***REMOVED***
	out := make([]byte, len(in)+1)
	if uintptr(unsafe.Pointer(&out[0]))&(unsafe.Alignof(uint32(0))-1) == 0 ***REMOVED***
		out = out[1:]
	***REMOVED*** else ***REMOVED***
		out = out[:len(in)]
	***REMOVED***
	copy(out, in)
	return out
***REMOVED***
