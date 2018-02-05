// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package idna

// appendMapping appends the mapping for the respective rune. isMapped must be
// true. A mapping is a categorization of a rune as defined in UTS #46.
func (c info) appendMapping(b []byte, s string) []byte ***REMOVED***
	index := int(c >> indexShift)
	if c&xorBit == 0 ***REMOVED***
		s := mappings[index:]
		return append(b, s[1:s[0]+1]...)
	***REMOVED***
	b = append(b, s...)
	if c&inlineXOR == inlineXOR ***REMOVED***
		// TODO: support and handle two-byte inline masks
		b[len(b)-1] ^= byte(index)
	***REMOVED*** else ***REMOVED***
		for p := len(b) - int(xorData[index]); p < len(b); p++ ***REMOVED***
			index++
			b[p] ^= xorData[index]
		***REMOVED***
	***REMOVED***
	return b
***REMOVED***

// Sparse block handling code.

type valueRange struct ***REMOVED***
	value  uint16 // header: value:stride
	lo, hi byte   // header: lo:n
***REMOVED***

type sparseBlocks struct ***REMOVED***
	values []valueRange
	offset []uint16
***REMOVED***

var idnaSparse = sparseBlocks***REMOVED***
	values: idnaSparseValues[:],
	offset: idnaSparseOffset[:],
***REMOVED***

// Don't use newIdnaTrie to avoid unconditional linking in of the table.
var trie = &idnaTrie***REMOVED******REMOVED***

// lookup determines the type of block n and looks up the value for b.
// For n < t.cutoff, the block is a simple lookup table. Otherwise, the block
// is a list of ranges with an accompanying value. Given a matching range r,
// the value for b is by r.value + (b - r.lo) * stride.
func (t *sparseBlocks) lookup(n uint32, b byte) uint16 ***REMOVED***
	offset := t.offset[n]
	header := t.values[offset]
	lo := offset + 1
	hi := lo + uint16(header.lo)
	for lo < hi ***REMOVED***
		m := lo + (hi-lo)/2
		r := t.values[m]
		if r.lo <= b && b <= r.hi ***REMOVED***
			return r.value + uint16(b-r.lo)*header.value
		***REMOVED***
		if b < r.lo ***REMOVED***
			hi = m
		***REMOVED*** else ***REMOVED***
			lo = m + 1
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***
