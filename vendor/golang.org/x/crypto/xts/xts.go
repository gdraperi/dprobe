// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xts implements the XTS cipher mode as specified in IEEE P1619/D16.
//
// XTS mode is typically used for disk encryption, which presents a number of
// novel problems that make more common modes inapplicable. The disk is
// conceptually an array of sectors and we must be able to encrypt and decrypt
// a sector in isolation. However, an attacker must not be able to transpose
// two sectors of plaintext by transposing their ciphertext.
//
// XTS wraps a block cipher with Rogaway's XEX mode in order to build a
// tweakable block cipher. This allows each sector to have a unique tweak and
// effectively create a unique key for each sector.
//
// XTS does not provide any authentication. An attacker can manipulate the
// ciphertext and randomise a block (16 bytes) of the plaintext.
//
// (Note: this package does not implement ciphertext-stealing so sectors must
// be a multiple of 16 bytes.)
package xts // import "golang.org/x/crypto/xts"

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"
)

// Cipher contains an expanded key structure. It doesn't contain mutable state
// and therefore can be used concurrently.
type Cipher struct ***REMOVED***
	k1, k2 cipher.Block
***REMOVED***

// blockSize is the block size that the underlying cipher must have. XTS is
// only defined for 16-byte ciphers.
const blockSize = 16

// NewCipher creates a Cipher given a function for creating the underlying
// block cipher (which must have a block size of 16 bytes). The key must be
// twice the length of the underlying cipher's key.
func NewCipher(cipherFunc func([]byte) (cipher.Block, error), key []byte) (c *Cipher, err error) ***REMOVED***
	c = new(Cipher)
	if c.k1, err = cipherFunc(key[:len(key)/2]); err != nil ***REMOVED***
		return
	***REMOVED***
	c.k2, err = cipherFunc(key[len(key)/2:])

	if c.k1.BlockSize() != blockSize ***REMOVED***
		err = errors.New("xts: cipher does not have a block size of 16")
	***REMOVED***

	return
***REMOVED***

// Encrypt encrypts a sector of plaintext and puts the result into ciphertext.
// Plaintext and ciphertext must overlap entirely or not at all.
// Sectors must be a multiple of 16 bytes and less than 2²⁴ bytes.
func (c *Cipher) Encrypt(ciphertext, plaintext []byte, sectorNum uint64) ***REMOVED***
	if len(ciphertext) < len(plaintext) ***REMOVED***
		panic("xts: ciphertext is smaller than plaintext")
	***REMOVED***
	if len(plaintext)%blockSize != 0 ***REMOVED***
		panic("xts: plaintext is not a multiple of the block size")
	***REMOVED***

	var tweak [blockSize]byte
	binary.LittleEndian.PutUint64(tweak[:8], sectorNum)

	c.k2.Encrypt(tweak[:], tweak[:])

	for len(plaintext) > 0 ***REMOVED***
		for j := range tweak ***REMOVED***
			ciphertext[j] = plaintext[j] ^ tweak[j]
		***REMOVED***
		c.k1.Encrypt(ciphertext, ciphertext)
		for j := range tweak ***REMOVED***
			ciphertext[j] ^= tweak[j]
		***REMOVED***
		plaintext = plaintext[blockSize:]
		ciphertext = ciphertext[blockSize:]

		mul2(&tweak)
	***REMOVED***
***REMOVED***

// Decrypt decrypts a sector of ciphertext and puts the result into plaintext.
// Plaintext and ciphertext must overlap entirely or not at all.
// Sectors must be a multiple of 16 bytes and less than 2²⁴ bytes.
func (c *Cipher) Decrypt(plaintext, ciphertext []byte, sectorNum uint64) ***REMOVED***
	if len(plaintext) < len(ciphertext) ***REMOVED***
		panic("xts: plaintext is smaller than ciphertext")
	***REMOVED***
	if len(ciphertext)%blockSize != 0 ***REMOVED***
		panic("xts: ciphertext is not a multiple of the block size")
	***REMOVED***

	var tweak [blockSize]byte
	binary.LittleEndian.PutUint64(tweak[:8], sectorNum)

	c.k2.Encrypt(tweak[:], tweak[:])

	for len(ciphertext) > 0 ***REMOVED***
		for j := range tweak ***REMOVED***
			plaintext[j] = ciphertext[j] ^ tweak[j]
		***REMOVED***
		c.k1.Decrypt(plaintext, plaintext)
		for j := range tweak ***REMOVED***
			plaintext[j] ^= tweak[j]
		***REMOVED***
		plaintext = plaintext[blockSize:]
		ciphertext = ciphertext[blockSize:]

		mul2(&tweak)
	***REMOVED***
***REMOVED***

// mul2 multiplies tweak by 2 in GF(2¹²⁸) with an irreducible polynomial of
// x¹²⁸ + x⁷ + x² + x + 1.
func mul2(tweak *[blockSize]byte) ***REMOVED***
	var carryIn byte
	for j := range tweak ***REMOVED***
		carryOut := tweak[j] >> 7
		tweak[j] = (tweak[j] << 1) + carryIn
		carryIn = carryOut
	***REMOVED***
	if carryIn != 0 ***REMOVED***
		// If we have a carry bit then we need to subtract a multiple
		// of the irreducible polynomial (x¹²⁸ + x⁷ + x² + x + 1).
		// By dropping the carry bit, we're subtracting the x^128 term
		// so all that remains is to subtract x⁷ + x² + x + 1.
		// Subtraction (and addition) in this representation is just
		// XOR.
		tweak[0] ^= 1<<7 | 1<<2 | 1<<1 | 1
	***REMOVED***
***REMOVED***
