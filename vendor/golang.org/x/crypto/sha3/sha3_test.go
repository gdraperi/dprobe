// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sha3

// Tests include all the ShortMsgKATs provided by the Keccak team at
// https://github.com/gvanas/KeccakCodePackage
//
// They only include the zero-bit case of the bitwise testvectors
// published by NIST in the draft of FIPS-202.

import (
	"bytes"
	"compress/flate"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"os"
	"strings"
	"testing"
)

const (
	testString  = "brekeccakkeccak koax koax"
	katFilename = "testdata/keccakKats.json.deflate"
)

// Internal-use instances of SHAKE used to test against KATs.
func newHashShake128() hash.Hash ***REMOVED***
	return &state***REMOVED***rate: 168, dsbyte: 0x1f, outputLen: 512***REMOVED***
***REMOVED***
func newHashShake256() hash.Hash ***REMOVED***
	return &state***REMOVED***rate: 136, dsbyte: 0x1f, outputLen: 512***REMOVED***
***REMOVED***

// testDigests contains functions returning hash.Hash instances
// with output-length equal to the KAT length for both SHA-3 and
// SHAKE instances.
var testDigests = map[string]func() hash.Hash***REMOVED***
	"SHA3-224": New224,
	"SHA3-256": New256,
	"SHA3-384": New384,
	"SHA3-512": New512,
	"SHAKE128": newHashShake128,
	"SHAKE256": newHashShake256,
***REMOVED***

// testShakes contains functions that return ShakeHash instances for
// testing the ShakeHash-specific interface.
var testShakes = map[string]func() ShakeHash***REMOVED***
	"SHAKE128": NewShake128,
	"SHAKE256": NewShake256,
***REMOVED***

// decodeHex converts a hex-encoded string into a raw byte string.
func decodeHex(s string) []byte ***REMOVED***
	b, err := hex.DecodeString(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return b
***REMOVED***

// structs used to marshal JSON test-cases.
type KeccakKats struct ***REMOVED***
	Kats map[string][]struct ***REMOVED***
		Digest  string `json:"digest"`
		Length  int64  `json:"length"`
		Message string `json:"message"`
	***REMOVED***
***REMOVED***

func testUnalignedAndGeneric(t *testing.T, testf func(impl string)) ***REMOVED***
	xorInOrig, copyOutOrig := xorIn, copyOut
	xorIn, copyOut = xorInGeneric, copyOutGeneric
	testf("generic")
	if xorImplementationUnaligned != "generic" ***REMOVED***
		xorIn, copyOut = xorInUnaligned, copyOutUnaligned
		testf("unaligned")
	***REMOVED***
	xorIn, copyOut = xorInOrig, copyOutOrig
***REMOVED***

// TestKeccakKats tests the SHA-3 and Shake implementations against all the
// ShortMsgKATs from https://github.com/gvanas/KeccakCodePackage
// (The testvectors are stored in keccakKats.json.deflate due to their length.)
func TestKeccakKats(t *testing.T) ***REMOVED***
	testUnalignedAndGeneric(t, func(impl string) ***REMOVED***
		// Read the KATs.
		deflated, err := os.Open(katFilename)
		if err != nil ***REMOVED***
			t.Errorf("error opening %s: %s", katFilename, err)
		***REMOVED***
		file := flate.NewReader(deflated)
		dec := json.NewDecoder(file)
		var katSet KeccakKats
		err = dec.Decode(&katSet)
		if err != nil ***REMOVED***
			t.Errorf("error decoding KATs: %s", err)
		***REMOVED***

		// Do the KATs.
		for functionName, kats := range katSet.Kats ***REMOVED***
			d := testDigests[functionName]()
			for _, kat := range kats ***REMOVED***
				d.Reset()
				in, err := hex.DecodeString(kat.Message)
				if err != nil ***REMOVED***
					t.Errorf("error decoding KAT: %s", err)
				***REMOVED***
				d.Write(in[:kat.Length/8])
				got := strings.ToUpper(hex.EncodeToString(d.Sum(nil)))
				if got != kat.Digest ***REMOVED***
					t.Errorf("function=%s, implementation=%s, length=%d\nmessage:\n  %s\ngot:\n  %s\nwanted:\n %s",
						functionName, impl, kat.Length, kat.Message, got, kat.Digest)
					t.Logf("wanted %+v", kat)
					t.FailNow()
				***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

// TestUnalignedWrite tests that writing data in an arbitrary pattern with
// small input buffers.
func testUnalignedWrite(t *testing.T) ***REMOVED***
	testUnalignedAndGeneric(t, func(impl string) ***REMOVED***
		buf := sequentialBytes(0x10000)
		for alg, df := range testDigests ***REMOVED***
			d := df()
			d.Reset()
			d.Write(buf)
			want := d.Sum(nil)
			d.Reset()
			for i := 0; i < len(buf); ***REMOVED***
				// Cycle through offsets which make a 137 byte sequence.
				// Because 137 is prime this sequence should exercise all corner cases.
				offsets := [17]int***REMOVED***1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 1***REMOVED***
				for _, j := range offsets ***REMOVED***
					if v := len(buf) - i; v < j ***REMOVED***
						j = v
					***REMOVED***
					d.Write(buf[i : i+j])
					i += j
				***REMOVED***
			***REMOVED***
			got := d.Sum(nil)
			if !bytes.Equal(got, want) ***REMOVED***
				t.Errorf("Unaligned writes, implementation=%s, alg=%s\ngot %q, want %q", impl, alg, got, want)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

// TestAppend checks that appending works when reallocation is necessary.
func TestAppend(t *testing.T) ***REMOVED***
	testUnalignedAndGeneric(t, func(impl string) ***REMOVED***
		d := New224()

		for capacity := 2; capacity <= 66; capacity += 64 ***REMOVED***
			// The first time around the loop, Sum will have to reallocate.
			// The second time, it will not.
			buf := make([]byte, 2, capacity)
			d.Reset()
			d.Write([]byte***REMOVED***0xcc***REMOVED***)
			buf = d.Sum(buf)
			expected := "0000DF70ADC49B2E76EEE3A6931B93FA41841C3AF2CDF5B32A18B5478C39"
			if got := strings.ToUpper(hex.EncodeToString(buf)); got != expected ***REMOVED***
				t.Errorf("got %s, want %s", got, expected)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

// TestAppendNoRealloc tests that appending works when no reallocation is necessary.
func TestAppendNoRealloc(t *testing.T) ***REMOVED***
	testUnalignedAndGeneric(t, func(impl string) ***REMOVED***
		buf := make([]byte, 1, 200)
		d := New224()
		d.Write([]byte***REMOVED***0xcc***REMOVED***)
		buf = d.Sum(buf)
		expected := "00DF70ADC49B2E76EEE3A6931B93FA41841C3AF2CDF5B32A18B5478C39"
		if got := strings.ToUpper(hex.EncodeToString(buf)); got != expected ***REMOVED***
			t.Errorf("%s: got %s, want %s", impl, got, expected)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// TestSqueezing checks that squeezing the full output a single time produces
// the same output as repeatedly squeezing the instance.
func TestSqueezing(t *testing.T) ***REMOVED***
	testUnalignedAndGeneric(t, func(impl string) ***REMOVED***
		for functionName, newShakeHash := range testShakes ***REMOVED***
			d0 := newShakeHash()
			d0.Write([]byte(testString))
			ref := make([]byte, 32)
			d0.Read(ref)

			d1 := newShakeHash()
			d1.Write([]byte(testString))
			var multiple []byte
			for range ref ***REMOVED***
				one := make([]byte, 1)
				d1.Read(one)
				multiple = append(multiple, one...)
			***REMOVED***
			if !bytes.Equal(ref, multiple) ***REMOVED***
				t.Errorf("%s (%s): squeezing %d bytes one at a time failed", functionName, impl, len(ref))
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

// sequentialBytes produces a buffer of size consecutive bytes 0x00, 0x01, ..., used for testing.
func sequentialBytes(size int) []byte ***REMOVED***
	result := make([]byte, size)
	for i := range result ***REMOVED***
		result[i] = byte(i)
	***REMOVED***
	return result
***REMOVED***

// BenchmarkPermutationFunction measures the speed of the permutation function
// with no input data.
func BenchmarkPermutationFunction(b *testing.B) ***REMOVED***
	b.SetBytes(int64(200))
	var lanes [25]uint64
	for i := 0; i < b.N; i++ ***REMOVED***
		keccakF1600(&lanes)
	***REMOVED***
***REMOVED***

// benchmarkHash tests the speed to hash num buffers of buflen each.
func benchmarkHash(b *testing.B, h hash.Hash, size, num int) ***REMOVED***
	b.StopTimer()
	h.Reset()
	data := sequentialBytes(size)
	b.SetBytes(int64(size * num))
	b.StartTimer()

	var state []byte
	for i := 0; i < b.N; i++ ***REMOVED***
		for j := 0; j < num; j++ ***REMOVED***
			h.Write(data)
		***REMOVED***
		state = h.Sum(state[:0])
	***REMOVED***
	b.StopTimer()
	h.Reset()
***REMOVED***

// benchmarkShake is specialized to the Shake instances, which don't
// require a copy on reading output.
func benchmarkShake(b *testing.B, h ShakeHash, size, num int) ***REMOVED***
	b.StopTimer()
	h.Reset()
	data := sequentialBytes(size)
	d := make([]byte, 32)

	b.SetBytes(int64(size * num))
	b.StartTimer()

	for i := 0; i < b.N; i++ ***REMOVED***
		h.Reset()
		for j := 0; j < num; j++ ***REMOVED***
			h.Write(data)
		***REMOVED***
		h.Read(d)
	***REMOVED***
***REMOVED***

func BenchmarkSha3_512_MTU(b *testing.B) ***REMOVED*** benchmarkHash(b, New512(), 1350, 1) ***REMOVED***
func BenchmarkSha3_384_MTU(b *testing.B) ***REMOVED*** benchmarkHash(b, New384(), 1350, 1) ***REMOVED***
func BenchmarkSha3_256_MTU(b *testing.B) ***REMOVED*** benchmarkHash(b, New256(), 1350, 1) ***REMOVED***
func BenchmarkSha3_224_MTU(b *testing.B) ***REMOVED*** benchmarkHash(b, New224(), 1350, 1) ***REMOVED***

func BenchmarkShake128_MTU(b *testing.B)  ***REMOVED*** benchmarkShake(b, NewShake128(), 1350, 1) ***REMOVED***
func BenchmarkShake256_MTU(b *testing.B)  ***REMOVED*** benchmarkShake(b, NewShake256(), 1350, 1) ***REMOVED***
func BenchmarkShake256_16x(b *testing.B)  ***REMOVED*** benchmarkShake(b, NewShake256(), 16, 1024) ***REMOVED***
func BenchmarkShake256_1MiB(b *testing.B) ***REMOVED*** benchmarkShake(b, NewShake256(), 1024, 1024) ***REMOVED***

func BenchmarkSha3_512_1MiB(b *testing.B) ***REMOVED*** benchmarkHash(b, New512(), 1024, 1024) ***REMOVED***

func Example_sum() ***REMOVED***
	buf := []byte("some data to hash")
	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	h := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	ShakeSum256(h, buf)
	fmt.Printf("%x\n", h)
	// Output: 0f65fe41fc353e52c55667bb9e2b27bfcc8476f2c413e9437d272ee3194a4e3146d05ec04a25d16b8f577c19b82d16b1424c3e022e783d2b4da98de3658d363d
***REMOVED***

func Example_mac() ***REMOVED***
	k := []byte("this is a secret key; you should generate a strong random key that's at least 32 bytes long")
	buf := []byte("and this is some data to authenticate")
	// A MAC with 32 bytes of output has 256-bit security strength -- if you use at least a 32-byte-long key.
	h := make([]byte, 32)
	d := NewShake256()
	// Write the key into the hash.
	d.Write(k)
	// Now write the data.
	d.Write(buf)
	// Read 32 bytes of output from the hash into h.
	d.Read(h)
	fmt.Printf("%x\n", h)
	// Output: 78de2974bd2711d5549ffd32b753ef0f5fa80a0db2556db60f0987eb8a9218ff
***REMOVED***
