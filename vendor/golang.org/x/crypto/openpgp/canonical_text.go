// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openpgp

import "hash"

// NewCanonicalTextHash reformats text written to it into the canonical
// form and then applies the hash h.  See RFC 4880, section 5.2.1.
func NewCanonicalTextHash(h hash.Hash) hash.Hash ***REMOVED***
	return &canonicalTextHash***REMOVED***h, 0***REMOVED***
***REMOVED***

type canonicalTextHash struct ***REMOVED***
	h hash.Hash
	s int
***REMOVED***

var newline = []byte***REMOVED***'\r', '\n'***REMOVED***

func (cth *canonicalTextHash) Write(buf []byte) (int, error) ***REMOVED***
	start := 0

	for i, c := range buf ***REMOVED***
		switch cth.s ***REMOVED***
		case 0:
			if c == '\r' ***REMOVED***
				cth.s = 1
			***REMOVED*** else if c == '\n' ***REMOVED***
				cth.h.Write(buf[start:i])
				cth.h.Write(newline)
				start = i + 1
			***REMOVED***
		case 1:
			cth.s = 0
		***REMOVED***
	***REMOVED***

	cth.h.Write(buf[start:])
	return len(buf), nil
***REMOVED***

func (cth *canonicalTextHash) Sum(in []byte) []byte ***REMOVED***
	return cth.h.Sum(in)
***REMOVED***

func (cth *canonicalTextHash) Reset() ***REMOVED***
	cth.h.Reset()
	cth.s = 0
***REMOVED***

func (cth *canonicalTextHash) Size() int ***REMOVED***
	return cth.h.Size()
***REMOVED***

func (cth *canonicalTextHash) BlockSize() int ***REMOVED***
	return cth.h.BlockSize()
***REMOVED***
