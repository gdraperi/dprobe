// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package secretbox

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"
)

func TestSealOpen(t *testing.T) ***REMOVED***
	var key [32]byte
	var nonce [24]byte

	rand.Reader.Read(key[:])
	rand.Reader.Read(nonce[:])

	var box, opened []byte

	for msgLen := 0; msgLen < 128; msgLen += 17 ***REMOVED***
		message := make([]byte, msgLen)
		rand.Reader.Read(message)

		box = Seal(box[:0], message, &nonce, &key)
		var ok bool
		opened, ok = Open(opened[:0], box, &nonce, &key)
		if !ok ***REMOVED***
			t.Errorf("%d: failed to open box", msgLen)
			continue
		***REMOVED***

		if !bytes.Equal(opened, message) ***REMOVED***
			t.Errorf("%d: got %x, expected %x", msgLen, opened, message)
			continue
		***REMOVED***
	***REMOVED***

	for i := range box ***REMOVED***
		box[i] ^= 0x20
		_, ok := Open(opened[:0], box, &nonce, &key)
		if ok ***REMOVED***
			t.Errorf("box was opened after corrupting byte %d", i)
		***REMOVED***
		box[i] ^= 0x20
	***REMOVED***
***REMOVED***

func TestSecretBox(t *testing.T) ***REMOVED***
	var key [32]byte
	var nonce [24]byte
	var message [64]byte

	for i := range key[:] ***REMOVED***
		key[i] = 1
	***REMOVED***
	for i := range nonce[:] ***REMOVED***
		nonce[i] = 2
	***REMOVED***
	for i := range message[:] ***REMOVED***
		message[i] = 3
	***REMOVED***

	box := Seal(nil, message[:], &nonce, &key)
	// expected was generated using the C implementation of NaCl.
	expected, _ := hex.DecodeString("8442bc313f4626f1359e3b50122b6ce6fe66ddfe7d39d14e637eb4fd5b45beadab55198df6ab5368439792a23c87db70acb6156dc5ef957ac04f6276cf6093b84be77ff0849cc33e34b7254d5a8f65ad")

	if !bytes.Equal(box, expected) ***REMOVED***
		t.Fatalf("box didn't match, got\n%x\n, expected\n%x", box, expected)
	***REMOVED***
***REMOVED***

func TestAppend(t *testing.T) ***REMOVED***
	var key [32]byte
	var nonce [24]byte
	var message [8]byte

	out := make([]byte, 4)
	box := Seal(out, message[:], &nonce, &key)
	if !bytes.Equal(box[:4], out[:4]) ***REMOVED***
		t.Fatalf("Seal didn't correctly append")
	***REMOVED***

	out = make([]byte, 4, 100)
	box = Seal(out, message[:], &nonce, &key)
	if !bytes.Equal(box[:4], out[:4]) ***REMOVED***
		t.Fatalf("Seal didn't correctly append with sufficient capacity.")
	***REMOVED***
***REMOVED***

func benchmarkSealSize(b *testing.B, size int) ***REMOVED***
	message := make([]byte, size)
	out := make([]byte, size+Overhead)
	var nonce [24]byte
	var key [32]byte

	b.SetBytes(int64(size))
	b.ResetTimer()

	for i := 0; i < b.N; i++ ***REMOVED***
		out = Seal(out[:0], message, &nonce, &key)
	***REMOVED***
***REMOVED***

func BenchmarkSeal8Bytes(b *testing.B) ***REMOVED***
	benchmarkSealSize(b, 8)
***REMOVED***

func BenchmarkSeal100Bytes(b *testing.B) ***REMOVED***
	benchmarkSealSize(b, 100)
***REMOVED***

func BenchmarkSeal1K(b *testing.B) ***REMOVED***
	benchmarkSealSize(b, 1024)
***REMOVED***

func BenchmarkSeal8K(b *testing.B) ***REMOVED***
	benchmarkSealSize(b, 8192)
***REMOVED***

func benchmarkOpenSize(b *testing.B, size int) ***REMOVED***
	msg := make([]byte, size)
	result := make([]byte, size)
	var nonce [24]byte
	var key [32]byte
	box := Seal(nil, msg, &nonce, &key)

	b.SetBytes(int64(size))
	b.ResetTimer()

	for i := 0; i < b.N; i++ ***REMOVED***
		if _, ok := Open(result[:0], box, &nonce, &key); !ok ***REMOVED***
			panic("Open failed")
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkOpen8Bytes(b *testing.B) ***REMOVED***
	benchmarkOpenSize(b, 8)
***REMOVED***

func BenchmarkOpen100Bytes(b *testing.B) ***REMOVED***
	benchmarkOpenSize(b, 100)
***REMOVED***

func BenchmarkOpen1K(b *testing.B) ***REMOVED***
	benchmarkOpenSize(b, 1024)
***REMOVED***

func BenchmarkOpen8K(b *testing.B) ***REMOVED***
	benchmarkOpenSize(b, 8192)
***REMOVED***
