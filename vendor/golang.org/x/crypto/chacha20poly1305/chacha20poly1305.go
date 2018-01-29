// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package chacha20poly1305 implements the ChaCha20-Poly1305 AEAD as specified in RFC 7539.
package chacha20poly1305 // import "golang.org/x/crypto/chacha20poly1305"

import (
	"crypto/cipher"
	"errors"
)

const (
	// KeySize is the size of the key used by this AEAD, in bytes.
	KeySize = 32
	// NonceSize is the size of the nonce used with this AEAD, in bytes.
	NonceSize = 12
)

type chacha20poly1305 struct ***REMOVED***
	key [32]byte
***REMOVED***

// New returns a ChaCha20-Poly1305 AEAD that uses the given, 256-bit key.
func New(key []byte) (cipher.AEAD, error) ***REMOVED***
	if len(key) != KeySize ***REMOVED***
		return nil, errors.New("chacha20poly1305: bad key length")
	***REMOVED***
	ret := new(chacha20poly1305)
	copy(ret.key[:], key)
	return ret, nil
***REMOVED***

func (c *chacha20poly1305) NonceSize() int ***REMOVED***
	return NonceSize
***REMOVED***

func (c *chacha20poly1305) Overhead() int ***REMOVED***
	return 16
***REMOVED***

func (c *chacha20poly1305) Seal(dst, nonce, plaintext, additionalData []byte) []byte ***REMOVED***
	if len(nonce) != NonceSize ***REMOVED***
		panic("chacha20poly1305: bad nonce length passed to Seal")
	***REMOVED***

	if uint64(len(plaintext)) > (1<<38)-64 ***REMOVED***
		panic("chacha20poly1305: plaintext too large")
	***REMOVED***

	return c.seal(dst, nonce, plaintext, additionalData)
***REMOVED***

var errOpen = errors.New("chacha20poly1305: message authentication failed")

func (c *chacha20poly1305) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) ***REMOVED***
	if len(nonce) != NonceSize ***REMOVED***
		panic("chacha20poly1305: bad nonce length passed to Open")
	***REMOVED***
	if len(ciphertext) < 16 ***REMOVED***
		return nil, errOpen
	***REMOVED***
	if uint64(len(ciphertext)) > (1<<38)-48 ***REMOVED***
		panic("chacha20poly1305: ciphertext too large")
	***REMOVED***

	return c.open(dst, nonce, ciphertext, additionalData)
***REMOVED***

// sliceForAppend takes a slice and a requested number of bytes. It returns a
// slice with the contents of the given slice followed by that many bytes and a
// second slice that aliases into it and contains only the extra bytes. If the
// original slice has sufficient capacity then no allocation is performed.
func sliceForAppend(in []byte, n int) (head, tail []byte) ***REMOVED***
	if total := len(in) + n; cap(in) >= total ***REMOVED***
		head = in[:total]
	***REMOVED*** else ***REMOVED***
		head = make([]byte, total)
		copy(head, in)
	***REMOVED***
	tail = head[len(in):]
	return
***REMOVED***
