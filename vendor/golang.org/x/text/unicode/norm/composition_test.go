// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import "testing"

// TestCase is used for most tests.
type TestCase struct ***REMOVED***
	in  []rune
	out []rune
***REMOVED***

func runTests(t *testing.T, name string, fm Form, tests []TestCase) ***REMOVED***
	rb := reorderBuffer***REMOVED******REMOVED***
	rb.init(fm, nil)
	for i, test := range tests ***REMOVED***
		rb.setFlusher(nil, appendFlush)
		for j, rune := range test.in ***REMOVED***
			b := []byte(string(rune))
			src := inputBytes(b)
			info := rb.f.info(src, 0)
			if j == 0 ***REMOVED***
				rb.ss.first(info)
			***REMOVED*** else ***REMOVED***
				rb.ss.next(info)
			***REMOVED***
			if rb.insertFlush(src, 0, info) < 0 ***REMOVED***
				t.Errorf("%s:%d: insert failed for rune %d", name, i, j)
			***REMOVED***
		***REMOVED***
		rb.doFlush()
		was := string(rb.out)
		want := string(test.out)
		if len(was) != len(want) ***REMOVED***
			t.Errorf("%s:%d: length = %d; want %d", name, i, len(was), len(want))
		***REMOVED***
		if was != want ***REMOVED***
			k, pfx := pidx(was, want)
			t.Errorf("%s:%d: \nwas  %s%+q; \nwant %s%+q", name, i, pfx, was[k:], pfx, want[k:])
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFlush(t *testing.T) ***REMOVED***
	const (
		hello = "Hello "
		world = "world!"
	)
	buf := make([]byte, maxByteBufferSize)
	p := copy(buf, hello)
	out := buf[p:]
	rb := reorderBuffer***REMOVED******REMOVED***
	rb.initString(NFC, world)
	if i := rb.flushCopy(out); i != 0 ***REMOVED***
		t.Errorf("wrote bytes on flush of empty buffer. (len(out) = %d)", i)
	***REMOVED***

	for i := range world ***REMOVED***
		// No need to set streamSafe values for this test.
		rb.insertFlush(rb.src, i, rb.f.info(rb.src, i))
		n := rb.flushCopy(out)
		out = out[n:]
		p += n
	***REMOVED***

	was := buf[:p]
	want := hello + world
	if string(was) != want ***REMOVED***
		t.Errorf(`output after flush was "%s"; want "%s"`, string(was), want)
	***REMOVED***
	if rb.nrune != 0 ***REMOVED***
		t.Errorf("non-null size of info buffer (rb.nrune == %d)", rb.nrune)
	***REMOVED***
	if rb.nbyte != 0 ***REMOVED***
		t.Errorf("non-null size of byte buffer (rb.nbyte == %d)", rb.nbyte)
	***REMOVED***
***REMOVED***

var insertTests = []TestCase***REMOVED***
	***REMOVED***[]rune***REMOVED***'a'***REMOVED***, []rune***REMOVED***'a'***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x300***REMOVED***, []rune***REMOVED***0x300***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x300, 0x316***REMOVED***, []rune***REMOVED***0x316, 0x300***REMOVED******REMOVED***, // CCC(0x300)==230; CCC(0x316)==220
	***REMOVED***[]rune***REMOVED***0x316, 0x300***REMOVED***, []rune***REMOVED***0x316, 0x300***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x41, 0x316, 0x300***REMOVED***, []rune***REMOVED***0x41, 0x316, 0x300***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x41, 0x300, 0x316***REMOVED***, []rune***REMOVED***0x41, 0x316, 0x300***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x300, 0x316, 0x41***REMOVED***, []rune***REMOVED***0x316, 0x300, 0x41***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x41, 0x300, 0x40, 0x316***REMOVED***, []rune***REMOVED***0x41, 0x300, 0x40, 0x316***REMOVED******REMOVED***,
***REMOVED***

func TestInsert(t *testing.T) ***REMOVED***
	runTests(t, "TestInsert", NFD, insertTests)
***REMOVED***

var decompositionNFDTest = []TestCase***REMOVED***
	***REMOVED***[]rune***REMOVED***0xC0***REMOVED***, []rune***REMOVED***0x41, 0x300***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0xAC00***REMOVED***, []rune***REMOVED***0x1100, 0x1161***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x01C4***REMOVED***, []rune***REMOVED***0x01C4***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x320E***REMOVED***, []rune***REMOVED***0x320E***REMOVED******REMOVED***,
	***REMOVED***[]rune("음ẻ과"), []rune***REMOVED***0x110B, 0x1173, 0x11B7, 0x65, 0x309, 0x1100, 0x116A***REMOVED******REMOVED***,
***REMOVED***

var decompositionNFKDTest = []TestCase***REMOVED***
	***REMOVED***[]rune***REMOVED***0xC0***REMOVED***, []rune***REMOVED***0x41, 0x300***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0xAC00***REMOVED***, []rune***REMOVED***0x1100, 0x1161***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x01C4***REMOVED***, []rune***REMOVED***0x44, 0x5A, 0x030C***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x320E***REMOVED***, []rune***REMOVED***0x28, 0x1100, 0x1161, 0x29***REMOVED******REMOVED***,
***REMOVED***

func TestDecomposition(t *testing.T) ***REMOVED***
	runTests(t, "TestDecompositionNFD", NFD, decompositionNFDTest)
	runTests(t, "TestDecompositionNFKD", NFKD, decompositionNFKDTest)
***REMOVED***

var compositionTest = []TestCase***REMOVED***
	***REMOVED***[]rune***REMOVED***0x41, 0x300***REMOVED***, []rune***REMOVED***0xC0***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x41, 0x316***REMOVED***, []rune***REMOVED***0x41, 0x316***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x41, 0x300, 0x35D***REMOVED***, []rune***REMOVED***0xC0, 0x35D***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x41, 0x316, 0x300***REMOVED***, []rune***REMOVED***0xC0, 0x316***REMOVED******REMOVED***,
	// blocking starter
	***REMOVED***[]rune***REMOVED***0x41, 0x316, 0x40, 0x300***REMOVED***, []rune***REMOVED***0x41, 0x316, 0x40, 0x300***REMOVED******REMOVED***,
	***REMOVED***[]rune***REMOVED***0x1100, 0x1161***REMOVED***, []rune***REMOVED***0xAC00***REMOVED******REMOVED***,
	// parenthesized Hangul, alternate between ASCII and Hangul.
	***REMOVED***[]rune***REMOVED***0x28, 0x1100, 0x1161, 0x29***REMOVED***, []rune***REMOVED***0x28, 0xAC00, 0x29***REMOVED******REMOVED***,
***REMOVED***

func TestComposition(t *testing.T) ***REMOVED***
	runTests(t, "TestComposition", NFC, compositionTest)
***REMOVED***
