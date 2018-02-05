// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"bytes"
	"sort"
	"testing"

	"golang.org/x/text/internal/colltab"
)

var largetosmall = []stridx***REMOVED***
	***REMOVED***"a", 5***REMOVED***,
	***REMOVED***"ab", 4***REMOVED***,
	***REMOVED***"abc", 3***REMOVED***,
	***REMOVED***"abcd", 2***REMOVED***,
	***REMOVED***"abcde", 1***REMOVED***,
	***REMOVED***"abcdef", 0***REMOVED***,
***REMOVED***

var offsetSortTests = [][]stridx***REMOVED***
	***REMOVED***
		***REMOVED***"bcde", 1***REMOVED***,
		***REMOVED***"bc", 5***REMOVED***,
		***REMOVED***"ab", 4***REMOVED***,
		***REMOVED***"bcd", 3***REMOVED***,
		***REMOVED***"abcd", 0***REMOVED***,
		***REMOVED***"abc", 2***REMOVED***,
	***REMOVED***,
	largetosmall,
***REMOVED***

func TestOffsetSort(t *testing.T) ***REMOVED***
	for i, st := range offsetSortTests ***REMOVED***
		sort.Sort(offsetSort(st))
		for j, si := range st ***REMOVED***
			if j != si.index ***REMOVED***
				t.Errorf("%d: failed: %v", i, st)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for i, tt := range genStateTests ***REMOVED***
		// ensure input is well-formed
		sort.Sort(offsetSort(tt.in))
		for j, si := range tt.in ***REMOVED***
			if si.index != j+1 ***REMOVED***
				t.Errorf("%dth sort failed: %v", i, tt.in)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var genidxtest1 = []stridx***REMOVED***
	***REMOVED***"bcde", 3***REMOVED***,
	***REMOVED***"bc", 6***REMOVED***,
	***REMOVED***"ab", 2***REMOVED***,
	***REMOVED***"bcd", 5***REMOVED***,
	***REMOVED***"abcd", 0***REMOVED***,
	***REMOVED***"abc", 1***REMOVED***,
	***REMOVED***"bcdf", 4***REMOVED***,
***REMOVED***

var genidxSortTests = [][]stridx***REMOVED***
	genidxtest1,
	largetosmall,
***REMOVED***

func TestGenIdxSort(t *testing.T) ***REMOVED***
	for i, st := range genidxSortTests ***REMOVED***
		sort.Sort(genidxSort(st))
		for j, si := range st ***REMOVED***
			if j != si.index ***REMOVED***
				t.Errorf("%dth sort failed %v", i, st)
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var entrySortTests = []colltab.ContractTrieSet***REMOVED***
	***REMOVED***
		***REMOVED***10, 0, 1, 3***REMOVED***,
		***REMOVED***99, 0, 1, 0***REMOVED***,
		***REMOVED***20, 50, 0, 2***REMOVED***,
		***REMOVED***30, 0, 1, 1***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestEntrySort(t *testing.T) ***REMOVED***
	for i, et := range entrySortTests ***REMOVED***
		sort.Sort(entrySort(et))
		for j, fe := range et ***REMOVED***
			if j != int(fe.I) ***REMOVED***
				t.Errorf("%dth sort failed %v", i, et)
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type GenStateTest struct ***REMOVED***
	in            []stridx
	firstBlockLen int
	out           colltab.ContractTrieSet
***REMOVED***

var genStateTests = []GenStateTest***REMOVED***
	***REMOVED***[]stridx***REMOVED***
		***REMOVED***"abc", 1***REMOVED***,
	***REMOVED***,
		1,
		colltab.ContractTrieSet***REMOVED***
			***REMOVED***'a', 0, 1, noIndex***REMOVED***,
			***REMOVED***'b', 0, 1, noIndex***REMOVED***,
			***REMOVED***'c', 'c', final, 1***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***[]stridx***REMOVED***
		***REMOVED***"abc", 1***REMOVED***,
		***REMOVED***"abd", 2***REMOVED***,
		***REMOVED***"abe", 3***REMOVED***,
	***REMOVED***,
		1,
		colltab.ContractTrieSet***REMOVED***
			***REMOVED***'a', 0, 1, noIndex***REMOVED***,
			***REMOVED***'b', 0, 1, noIndex***REMOVED***,
			***REMOVED***'c', 'e', final, 1***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***[]stridx***REMOVED***
		***REMOVED***"abc", 1***REMOVED***,
		***REMOVED***"ab", 2***REMOVED***,
		***REMOVED***"a", 3***REMOVED***,
	***REMOVED***,
		1,
		colltab.ContractTrieSet***REMOVED***
			***REMOVED***'a', 0, 1, 3***REMOVED***,
			***REMOVED***'b', 0, 1, 2***REMOVED***,
			***REMOVED***'c', 'c', final, 1***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***[]stridx***REMOVED***
		***REMOVED***"abc", 1***REMOVED***,
		***REMOVED***"abd", 2***REMOVED***,
		***REMOVED***"ab", 3***REMOVED***,
		***REMOVED***"ac", 4***REMOVED***,
		***REMOVED***"a", 5***REMOVED***,
		***REMOVED***"b", 6***REMOVED***,
	***REMOVED***,
		2,
		colltab.ContractTrieSet***REMOVED***
			***REMOVED***'b', 'b', final, 6***REMOVED***,
			***REMOVED***'a', 0, 2, 5***REMOVED***,
			***REMOVED***'c', 'c', final, 4***REMOVED***,
			***REMOVED***'b', 0, 1, 3***REMOVED***,
			***REMOVED***'c', 'd', final, 1***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***[]stridx***REMOVED***
		***REMOVED***"bcde", 2***REMOVED***,
		***REMOVED***"bc", 7***REMOVED***,
		***REMOVED***"ab", 6***REMOVED***,
		***REMOVED***"bcd", 5***REMOVED***,
		***REMOVED***"abcd", 1***REMOVED***,
		***REMOVED***"abc", 4***REMOVED***,
		***REMOVED***"bcdf", 3***REMOVED***,
	***REMOVED***,
		2,
		colltab.ContractTrieSet***REMOVED***
			***REMOVED***'b', 3, 1, noIndex***REMOVED***,
			***REMOVED***'a', 0, 1, noIndex***REMOVED***,
			***REMOVED***'b', 0, 1, 6***REMOVED***,
			***REMOVED***'c', 0, 1, 4***REMOVED***,
			***REMOVED***'d', 'd', final, 1***REMOVED***,
			***REMOVED***'c', 0, 1, 7***REMOVED***,
			***REMOVED***'d', 0, 1, 5***REMOVED***,
			***REMOVED***'e', 'f', final, 2***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestGenStates(t *testing.T) ***REMOVED***
	for i, tt := range genStateTests ***REMOVED***
		si := []stridx***REMOVED******REMOVED***
		for _, e := range tt.in ***REMOVED***
			si = append(si, e)
		***REMOVED***
		// ensure input is well-formed
		sort.Sort(genidxSort(si))
		ct := colltab.ContractTrieSet***REMOVED******REMOVED***
		n, _ := genStates(&ct, si)
		if nn := tt.firstBlockLen; nn != n ***REMOVED***
			t.Errorf("%d: block len %v; want %v", i, n, nn)
		***REMOVED***
		if lv, lw := len(ct), len(tt.out); lv != lw ***REMOVED***
			t.Errorf("%d: len %v; want %v", i, lv, lw)
			continue
		***REMOVED***
		for j, fe := range tt.out ***REMOVED***
			const msg = "%d:%d: value %s=%v; want %v"
			if fe.L != ct[j].L ***REMOVED***
				t.Errorf(msg, i, j, "l", ct[j].L, fe.L)
			***REMOVED***
			if fe.H != ct[j].H ***REMOVED***
				t.Errorf(msg, i, j, "h", ct[j].H, fe.H)
			***REMOVED***
			if fe.N != ct[j].N ***REMOVED***
				t.Errorf(msg, i, j, "n", ct[j].N, fe.N)
			***REMOVED***
			if fe.I != ct[j].I ***REMOVED***
				t.Errorf(msg, i, j, "i", ct[j].I, fe.I)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestLookupContraction(t *testing.T) ***REMOVED***
	for i, tt := range genStateTests ***REMOVED***
		input := []string***REMOVED******REMOVED***
		for _, e := range tt.in ***REMOVED***
			input = append(input, e.str)
		***REMOVED***
		cts := colltab.ContractTrieSet***REMOVED******REMOVED***
		h, _ := appendTrie(&cts, input)
		for j, si := range tt.in ***REMOVED***
			str := si.str
			for _, s := range []string***REMOVED***str, str + "X"***REMOVED*** ***REMOVED***
				msg := "%d:%d: %s(%s) %v; want %v"
				idx, sn := lookup(&cts, h, []byte(s))
				if idx != si.index ***REMOVED***
					t.Errorf(msg, i, j, "index", s, idx, si.index)
				***REMOVED***
				if sn != len(str) ***REMOVED***
					t.Errorf(msg, i, j, "sn", s, sn, len(str))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestPrintContractionTrieSet(t *testing.T) ***REMOVED***
	testdata := colltab.ContractTrieSet(genStateTests[4].out)
	buf := &bytes.Buffer***REMOVED******REMOVED***
	print(&testdata, buf, "test")
	if contractTrieOutput != buf.String() ***REMOVED***
		t.Errorf("output differs; found\n%s", buf.String())
		println(string(buf.Bytes()))
	***REMOVED***
***REMOVED***

const contractTrieOutput = `// testCTEntries: 8 entries, 32 bytes
var testCTEntries = [8]struct***REMOVED***L,H,N,I uint8***REMOVED******REMOVED***
	***REMOVED***0x62, 0x3, 1, 255***REMOVED***,
	***REMOVED***0x61, 0x0, 1, 255***REMOVED***,
	***REMOVED***0x62, 0x0, 1, 6***REMOVED***,
	***REMOVED***0x63, 0x0, 1, 4***REMOVED***,
	***REMOVED***0x64, 0x64, 0, 1***REMOVED***,
	***REMOVED***0x63, 0x0, 1, 7***REMOVED***,
	***REMOVED***0x64, 0x0, 1, 5***REMOVED***,
	***REMOVED***0x65, 0x66, 0, 2***REMOVED***,
***REMOVED***
var testContractTrieSet = colltab.ContractTrieSet( testCTEntries[:] )
`
