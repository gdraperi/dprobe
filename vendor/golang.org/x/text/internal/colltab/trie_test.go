// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colltab

import (
	"testing"
)

// We take the smallest, largest and an arbitrary value for each
// of the UTF-8 sequence lengths.
var testRunes = []rune***REMOVED***
	0x01, 0x0C, 0x7F, // 1-byte sequences
	0x80, 0x100, 0x7FF, // 2-byte sequences
	0x800, 0x999, 0xFFFF, // 3-byte sequences
	0x10000, 0x10101, 0x10FFFF, // 4-byte sequences
	0x200, 0x201, 0x202, 0x210, 0x215, // five entries in one sparse block
***REMOVED***

// Test cases for illegal runes.
type trietest struct ***REMOVED***
	size  int
	bytes []byte
***REMOVED***

var tests = []trietest***REMOVED***
	// illegal runes
	***REMOVED***1, []byte***REMOVED***0x80***REMOVED******REMOVED***,
	***REMOVED***1, []byte***REMOVED***0xFF***REMOVED******REMOVED***,
	***REMOVED***1, []byte***REMOVED***t2, tx - 1***REMOVED******REMOVED***,
	***REMOVED***1, []byte***REMOVED***t2, t2***REMOVED******REMOVED***,
	***REMOVED***2, []byte***REMOVED***t3, tx, tx - 1***REMOVED******REMOVED***,
	***REMOVED***2, []byte***REMOVED***t3, tx, t2***REMOVED******REMOVED***,
	***REMOVED***1, []byte***REMOVED***t3, tx - 1, tx***REMOVED******REMOVED***,
	***REMOVED***3, []byte***REMOVED***t4, tx, tx, tx - 1***REMOVED******REMOVED***,
	***REMOVED***3, []byte***REMOVED***t4, tx, tx, t2***REMOVED******REMOVED***,
	***REMOVED***1, []byte***REMOVED***t4, t2, tx, tx - 1***REMOVED******REMOVED***,
	***REMOVED***2, []byte***REMOVED***t4, tx, t2, tx - 1***REMOVED******REMOVED***,

	// short runes
	***REMOVED***0, []byte***REMOVED***t2***REMOVED******REMOVED***,
	***REMOVED***0, []byte***REMOVED***t3, tx***REMOVED******REMOVED***,
	***REMOVED***0, []byte***REMOVED***t4, tx, tx***REMOVED******REMOVED***,

	// we only support UTF-8 up to utf8.UTFMax bytes (4 bytes)
	***REMOVED***1, []byte***REMOVED***t5, tx, tx, tx, tx***REMOVED******REMOVED***,
	***REMOVED***1, []byte***REMOVED***t6, tx, tx, tx, tx, tx***REMOVED******REMOVED***,
***REMOVED***

func TestLookupTrie(t *testing.T) ***REMOVED***
	for i, r := range testRunes ***REMOVED***
		b := []byte(string(r))
		v, sz := testTrie.lookup(b)
		if int(v) != i ***REMOVED***
			t.Errorf("lookup(%U): found value %#x, expected %#x", r, v, i)
		***REMOVED***
		if sz != len(b) ***REMOVED***
			t.Errorf("lookup(%U): found size %d, expected %d", r, sz, len(b))
		***REMOVED***
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		v, sz := testTrie.lookup(tt.bytes)
		if int(v) != 0 ***REMOVED***
			t.Errorf("lookup of illegal rune, case %d: found value %#x, expected 0", i, v)
		***REMOVED***
		if sz != tt.size ***REMOVED***
			t.Errorf("lookup of illegal rune, case %d: found size %d, expected %d", i, sz, tt.size)
		***REMOVED***
	***REMOVED***
***REMOVED***

// test data is taken from exp/collate/locale/build/trie_test.go
var testValues = [832]uint32***REMOVED***
	0x000c: 0x00000001,
	0x007f: 0x00000002,
	0x00c0: 0x00000003,
	0x0100: 0x00000004,
	0x0140: 0x0000000c, 0x0141: 0x0000000d, 0x0142: 0x0000000e,
	0x0150: 0x0000000f,
	0x0155: 0x00000010,
	0x01bf: 0x00000005,
	0x01c0: 0x00000006,
	0x0219: 0x00000007,
	0x027f: 0x00000008,
	0x0280: 0x00000009,
	0x02c1: 0x0000000a,
	0x033f: 0x0000000b,
***REMOVED***

var testLookup = [640]uint16***REMOVED***
	0x0e0: 0x05, 0x0e6: 0x06,
	0x13f: 0x07,
	0x140: 0x08, 0x144: 0x09,
	0x190: 0x03,
	0x1ff: 0x0a,
	0x20f: 0x05,
	0x242: 0x01, 0x244: 0x02,
	0x248: 0x03,
	0x25f: 0x04,
	0x260: 0x01,
	0x26f: 0x02,
	0x270: 0x04, 0x274: 0x06,
***REMOVED***

var testTrie = Trie***REMOVED***testLookup[6*blockSize:], testValues[:], testLookup[:], testValues[:]***REMOVED***
