// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hkdf implements the HMAC-based Extract-and-Expand Key Derivation
// Function (HKDF) as defined in RFC 5869.
//
// HKDF is a cryptographic key derivation function (KDF) with the goal of
// expanding limited input keying material into one or more cryptographically
// strong secret keys.
//
// RFC 5869: https://tools.ietf.org/html/rfc5869
package hkdf // import "golang.org/x/crypto/hkdf"

import (
	"crypto/hmac"
	"errors"
	"hash"
	"io"
)

type hkdf struct ***REMOVED***
	expander hash.Hash
	size     int

	info    []byte
	counter byte

	prev  []byte
	cache []byte
***REMOVED***

func (f *hkdf) Read(p []byte) (int, error) ***REMOVED***
	// Check whether enough data can be generated
	need := len(p)
	remains := len(f.cache) + int(255-f.counter+1)*f.size
	if remains < need ***REMOVED***
		return 0, errors.New("hkdf: entropy limit reached")
	***REMOVED***
	// Read from the cache, if enough data is present
	n := copy(p, f.cache)
	p = p[n:]

	// Fill the buffer
	for len(p) > 0 ***REMOVED***
		f.expander.Reset()
		f.expander.Write(f.prev)
		f.expander.Write(f.info)
		f.expander.Write([]byte***REMOVED***f.counter***REMOVED***)
		f.prev = f.expander.Sum(f.prev[:0])
		f.counter++

		// Copy the new batch into p
		f.cache = f.prev
		n = copy(p, f.cache)
		p = p[n:]
	***REMOVED***
	// Save leftovers for next run
	f.cache = f.cache[n:]

	return need, nil
***REMOVED***

// New returns a new HKDF using the given hash, the secret keying material to expand
// and optional salt and info fields.
func New(hash func() hash.Hash, secret, salt, info []byte) io.Reader ***REMOVED***
	if salt == nil ***REMOVED***
		salt = make([]byte, hash().Size())
	***REMOVED***
	extractor := hmac.New(hash, salt)
	extractor.Write(secret)
	prk := extractor.Sum(nil)

	return &hkdf***REMOVED***hmac.New(hash, prk), extractor.Size(), info, 1, nil, nil***REMOVED***
***REMOVED***
