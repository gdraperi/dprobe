// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sha3

// This file provides functions for creating instances of the SHA-3
// and SHAKE hash functions, as well as utility functions for hashing
// bytes.

import (
	"hash"
)

// New224 creates a new SHA3-224 hash.
// Its generic security strength is 224 bits against preimage attacks,
// and 112 bits against collision attacks.
func New224() hash.Hash ***REMOVED*** return &state***REMOVED***rate: 144, outputLen: 28, dsbyte: 0x06***REMOVED*** ***REMOVED***

// New256 creates a new SHA3-256 hash.
// Its generic security strength is 256 bits against preimage attacks,
// and 128 bits against collision attacks.
func New256() hash.Hash ***REMOVED*** return &state***REMOVED***rate: 136, outputLen: 32, dsbyte: 0x06***REMOVED*** ***REMOVED***

// New384 creates a new SHA3-384 hash.
// Its generic security strength is 384 bits against preimage attacks,
// and 192 bits against collision attacks.
func New384() hash.Hash ***REMOVED*** return &state***REMOVED***rate: 104, outputLen: 48, dsbyte: 0x06***REMOVED*** ***REMOVED***

// New512 creates a new SHA3-512 hash.
// Its generic security strength is 512 bits against preimage attacks,
// and 256 bits against collision attacks.
func New512() hash.Hash ***REMOVED*** return &state***REMOVED***rate: 72, outputLen: 64, dsbyte: 0x06***REMOVED*** ***REMOVED***

// Sum224 returns the SHA3-224 digest of the data.
func Sum224(data []byte) (digest [28]byte) ***REMOVED***
	h := New224()
	h.Write(data)
	h.Sum(digest[:0])
	return
***REMOVED***

// Sum256 returns the SHA3-256 digest of the data.
func Sum256(data []byte) (digest [32]byte) ***REMOVED***
	h := New256()
	h.Write(data)
	h.Sum(digest[:0])
	return
***REMOVED***

// Sum384 returns the SHA3-384 digest of the data.
func Sum384(data []byte) (digest [48]byte) ***REMOVED***
	h := New384()
	h.Write(data)
	h.Sum(digest[:0])
	return
***REMOVED***

// Sum512 returns the SHA3-512 digest of the data.
func Sum512(data []byte) (digest [64]byte) ***REMOVED***
	h := New512()
	h.Write(data)
	h.Sum(digest[:0])
	return
***REMOVED***
