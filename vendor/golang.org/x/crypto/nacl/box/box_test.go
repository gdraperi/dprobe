// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package box

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"golang.org/x/crypto/curve25519"
)

func TestSealOpen(t *testing.T) ***REMOVED***
	publicKey1, privateKey1, _ := GenerateKey(rand.Reader)
	publicKey2, privateKey2, _ := GenerateKey(rand.Reader)

	if *privateKey1 == *privateKey2 ***REMOVED***
		t.Fatalf("private keys are equal!")
	***REMOVED***
	if *publicKey1 == *publicKey2 ***REMOVED***
		t.Fatalf("public keys are equal!")
	***REMOVED***
	message := []byte("test message")
	var nonce [24]byte

	box := Seal(nil, message, &nonce, publicKey1, privateKey2)
	opened, ok := Open(nil, box, &nonce, publicKey2, privateKey1)
	if !ok ***REMOVED***
		t.Fatalf("failed to open box")
	***REMOVED***

	if !bytes.Equal(opened, message) ***REMOVED***
		t.Fatalf("got %x, want %x", opened, message)
	***REMOVED***

	for i := range box ***REMOVED***
		box[i] ^= 0x40
		_, ok := Open(nil, box, &nonce, publicKey2, privateKey1)
		if ok ***REMOVED***
			t.Fatalf("opened box with byte %d corrupted", i)
		***REMOVED***
		box[i] ^= 0x40
	***REMOVED***
***REMOVED***

func TestBox(t *testing.T) ***REMOVED***
	var privateKey1, privateKey2 [32]byte
	for i := range privateKey1[:] ***REMOVED***
		privateKey1[i] = 1
	***REMOVED***
	for i := range privateKey2[:] ***REMOVED***
		privateKey2[i] = 2
	***REMOVED***

	var publicKey1 [32]byte
	curve25519.ScalarBaseMult(&publicKey1, &privateKey1)
	var message [64]byte
	for i := range message[:] ***REMOVED***
		message[i] = 3
	***REMOVED***

	var nonce [24]byte
	for i := range nonce[:] ***REMOVED***
		nonce[i] = 4
	***REMOVED***

	box := Seal(nil, message[:], &nonce, &publicKey1, &privateKey2)

	// expected was generated using the C implementation of NaCl.
	expected, _ := hex.DecodeString("78ea30b19d2341ebbdba54180f821eec265cf86312549bea8a37652a8bb94f07b78a73ed1708085e6ddd0e943bbdeb8755079a37eb31d86163ce241164a47629c0539f330b4914cd135b3855bc2a2dfc")

	if !bytes.Equal(box, expected) ***REMOVED***
		t.Fatalf("box didn't match, got\n%x\n, expected\n%x", box, expected)
	***REMOVED***
***REMOVED***
