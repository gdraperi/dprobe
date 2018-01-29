// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tea

import (
	"bytes"
	"testing"
)

// A sample test key for when we just want to initialize a cipher
var testKey = []byte***REMOVED***0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF***REMOVED***

// Test that the block size for tea is correct
func TestBlocksize(t *testing.T) ***REMOVED***
	c, err := NewCipher(testKey)
	if err != nil ***REMOVED***
		t.Fatalf("NewCipher returned error: %s", err)
	***REMOVED***

	if result := c.BlockSize(); result != BlockSize ***REMOVED***
		t.Errorf("cipher.BlockSize returned %d, but expected %d", result, BlockSize)
	***REMOVED***
***REMOVED***

// Test that invalid key sizes return an error
func TestInvalidKeySize(t *testing.T) ***REMOVED***
	var key [KeySize + 1]byte

	if _, err := NewCipher(key[:]); err == nil ***REMOVED***
		t.Errorf("invalid key size %d didn't result in an error.", len(key))
	***REMOVED***

	if _, err := NewCipher(key[:KeySize-1]); err == nil ***REMOVED***
		t.Errorf("invalid key size %d didn't result in an error.", KeySize-1)
	***REMOVED***
***REMOVED***

// Test Vectors
type teaTest struct ***REMOVED***
	rounds     int
	key        []byte
	plaintext  []byte
	ciphertext []byte
***REMOVED***

var teaTests = []teaTest***REMOVED***
	// These were sourced from https://github.com/froydnj/ironclad/blob/master/testing/test-vectors/tea.testvec
	***REMOVED***
		numRounds,
		[]byte***REMOVED***0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00***REMOVED***,
		[]byte***REMOVED***0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00***REMOVED***,
		[]byte***REMOVED***0x41, 0xea, 0x3a, 0x0a, 0x94, 0xba, 0xa9, 0x40***REMOVED***,
	***REMOVED***,
	***REMOVED***
		numRounds,
		[]byte***REMOVED***0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff***REMOVED***,
		[]byte***REMOVED***0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff***REMOVED***,
		[]byte***REMOVED***0x31, 0x9b, 0xbe, 0xfb, 0x01, 0x6a, 0xbd, 0xb2***REMOVED***,
	***REMOVED***,
	***REMOVED***
		16,
		[]byte***REMOVED***0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00***REMOVED***,
		[]byte***REMOVED***0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00***REMOVED***,
		[]byte***REMOVED***0xed, 0x28, 0x5d, 0xa1, 0x45, 0x5b, 0x33, 0xc1***REMOVED***,
	***REMOVED***,
***REMOVED***

// Test encryption
func TestCipherEncrypt(t *testing.T) ***REMOVED***
	// Test encryption with standard 64 rounds
	for i, test := range teaTests ***REMOVED***
		c, err := NewCipherWithRounds(test.key, test.rounds)
		if err != nil ***REMOVED***
			t.Fatalf("#%d: NewCipher returned error: %s", i, err)
		***REMOVED***

		var ciphertext [BlockSize]byte
		c.Encrypt(ciphertext[:], test.plaintext)

		if !bytes.Equal(ciphertext[:], test.ciphertext) ***REMOVED***
			t.Errorf("#%d: incorrect ciphertext. Got %x, wanted %x", i, ciphertext, test.ciphertext)
		***REMOVED***

		var plaintext2 [BlockSize]byte
		c.Decrypt(plaintext2[:], ciphertext[:])

		if !bytes.Equal(plaintext2[:], test.plaintext) ***REMOVED***
			t.Errorf("#%d: incorrect plaintext. Got %x, wanted %x", i, plaintext2, test.plaintext)
		***REMOVED***
	***REMOVED***
***REMOVED***
