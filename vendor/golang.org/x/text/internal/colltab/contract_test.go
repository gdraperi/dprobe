// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colltab

import (
	"testing"
)

type lookupStrings struct ***REMOVED***
	str    string
	offset int
	n      int // bytes consumed from input
***REMOVED***

type LookupTest struct ***REMOVED***
	lookup []lookupStrings
	n      int
	tries  ContractTrieSet
***REMOVED***

var lookupTests = []LookupTest***REMOVED******REMOVED***
	[]lookupStrings***REMOVED***
		***REMOVED***"abc", 1, 3***REMOVED***,
		***REMOVED***"a", 0, 0***REMOVED***,
		***REMOVED***"b", 0, 0***REMOVED***,
		***REMOVED***"c", 0, 0***REMOVED***,
		***REMOVED***"d", 0, 0***REMOVED***,
	***REMOVED***,
	1,
	ContractTrieSet***REMOVED***
		***REMOVED***'a', 0, 1, 0xFF***REMOVED***,
		***REMOVED***'b', 0, 1, 0xFF***REMOVED***,
		***REMOVED***'c', 'c', 0, 1***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	[]lookupStrings***REMOVED***
		***REMOVED***"abc", 1, 3***REMOVED***,
		***REMOVED***"abd", 2, 3***REMOVED***,
		***REMOVED***"abe", 3, 3***REMOVED***,
		***REMOVED***"a", 0, 0***REMOVED***,
		***REMOVED***"ab", 0, 0***REMOVED***,
		***REMOVED***"d", 0, 0***REMOVED***,
		***REMOVED***"f", 0, 0***REMOVED***,
	***REMOVED***,
	1,
	ContractTrieSet***REMOVED***
		***REMOVED***'a', 0, 1, 0xFF***REMOVED***,
		***REMOVED***'b', 0, 1, 0xFF***REMOVED***,
		***REMOVED***'c', 'e', 0, 1***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	[]lookupStrings***REMOVED***
		***REMOVED***"abc", 1, 3***REMOVED***,
		***REMOVED***"ab", 2, 2***REMOVED***,
		***REMOVED***"a", 3, 1***REMOVED***,
		***REMOVED***"abcd", 1, 3***REMOVED***,
		***REMOVED***"abe", 2, 2***REMOVED***,
	***REMOVED***,
	1,
	ContractTrieSet***REMOVED***
		***REMOVED***'a', 0, 1, 3***REMOVED***,
		***REMOVED***'b', 0, 1, 2***REMOVED***,
		***REMOVED***'c', 'c', 0, 1***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	[]lookupStrings***REMOVED***
		***REMOVED***"abc", 1, 3***REMOVED***,
		***REMOVED***"abd", 2, 3***REMOVED***,
		***REMOVED***"ab", 3, 2***REMOVED***,
		***REMOVED***"ac", 4, 2***REMOVED***,
		***REMOVED***"a", 5, 1***REMOVED***,
		***REMOVED***"b", 6, 1***REMOVED***,
		***REMOVED***"ba", 6, 1***REMOVED***,
	***REMOVED***,
	2,
	ContractTrieSet***REMOVED***
		***REMOVED***'b', 'b', 0, 6***REMOVED***,
		***REMOVED***'a', 0, 2, 5***REMOVED***,
		***REMOVED***'c', 'c', 0, 4***REMOVED***,
		***REMOVED***'b', 0, 1, 3***REMOVED***,
		***REMOVED***'c', 'd', 0, 1***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	[]lookupStrings***REMOVED***
		***REMOVED***"bcde", 2, 4***REMOVED***,
		***REMOVED***"bc", 7, 2***REMOVED***,
		***REMOVED***"ab", 6, 2***REMOVED***,
		***REMOVED***"bcd", 5, 3***REMOVED***,
		***REMOVED***"abcd", 1, 4***REMOVED***,
		***REMOVED***"abc", 4, 3***REMOVED***,
		***REMOVED***"bcdf", 3, 4***REMOVED***,
	***REMOVED***,
	2,
	ContractTrieSet***REMOVED***
		***REMOVED***'b', 3, 1, 0xFF***REMOVED***,
		***REMOVED***'a', 0, 1, 0xFF***REMOVED***,
		***REMOVED***'b', 0, 1, 6***REMOVED***,
		***REMOVED***'c', 0, 1, 4***REMOVED***,
		***REMOVED***'d', 'd', 0, 1***REMOVED***,
		***REMOVED***'c', 0, 1, 7***REMOVED***,
		***REMOVED***'d', 0, 1, 5***REMOVED***,
		***REMOVED***'e', 'f', 0, 2***REMOVED***,
	***REMOVED***,
***REMOVED******REMOVED***

func lookup(c *ContractTrieSet, nnode int, s []uint8) (i, n int) ***REMOVED***
	scan := c.scanner(0, nnode, s)
	scan.scan(0)
	return scan.result()
***REMOVED***

func TestLookupContraction(t *testing.T) ***REMOVED***
	for i, tt := range lookupTests ***REMOVED***
		cts := ContractTrieSet(tt.tries)
		for j, lu := range tt.lookup ***REMOVED***
			str := lu.str
			for _, s := range []string***REMOVED***str, str + "X"***REMOVED*** ***REMOVED***
				const msg = `%d:%d: %s of "%s" %v; want %v`
				offset, n := lookup(&cts, tt.n, []byte(s))
				if offset != lu.offset ***REMOVED***
					t.Errorf(msg, i, j, "offset", s, offset, lu.offset)
				***REMOVED***
				if n != lu.n ***REMOVED***
					t.Errorf(msg, i, j, "bytes consumed", s, n, len(str))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
