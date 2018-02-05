// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collate

import (
	"testing"

	"golang.org/x/text/collate/build"
	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/unicode/norm"
)

type ColElems []Weights

type input struct ***REMOVED***
	str string
	ces [][]int
***REMOVED***

type check struct ***REMOVED***
	in  string
	n   int
	out ColElems
***REMOVED***

type tableTest struct ***REMOVED***
	in  []input
	chk []check
***REMOVED***

func w(ce ...int) Weights ***REMOVED***
	return W(ce...)
***REMOVED***

var defaults = w(0)

func pt(p, t int) []int ***REMOVED***
	return []int***REMOVED***p, defaults.Secondary, t***REMOVED***
***REMOVED***

func makeTable(in []input) (*Collator, error) ***REMOVED***
	b := build.NewBuilder()
	for _, r := range in ***REMOVED***
		if e := b.Add([]rune(r.str), r.ces, nil); e != nil ***REMOVED***
			panic(e)
		***REMOVED***
	***REMOVED***
	t, err := b.Build()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewFromTable(t), nil
***REMOVED***

// modSeq holds a seqeunce of modifiers in increasing order of CCC long enough
// to cause a segment overflow if not handled correctly. The last rune in this
// list has a CCC of 214.
var modSeq = []rune***REMOVED***
	0x05B1, 0x05B2, 0x05B3, 0x05B4, 0x05B5, 0x05B6, 0x05B7, 0x05B8, 0x05B9, 0x05BB,
	0x05BC, 0x05BD, 0x05BF, 0x05C1, 0x05C2, 0xFB1E, 0x064B, 0x064C, 0x064D, 0x064E,
	0x064F, 0x0650, 0x0651, 0x0652, 0x0670, 0x0711, 0x0C55, 0x0C56, 0x0E38, 0x0E48,
	0x0EB8, 0x0EC8, 0x0F71, 0x0F72, 0x0F74, 0x0321, 0x1DCE,
***REMOVED***

var mods []input
var modW = func() ColElems ***REMOVED***
	ws := ColElems***REMOVED******REMOVED***
	for _, r := range modSeq ***REMOVED***
		rune := norm.NFC.PropertiesString(string(r))
		ws = append(ws, w(0, int(rune.CCC())))
		mods = append(mods, input***REMOVED***string(r), [][]int***REMOVED******REMOVED***0, int(rune.CCC())***REMOVED******REMOVED******REMOVED***)
	***REMOVED***
	return ws
***REMOVED***()

var appendNextTests = []tableTest***REMOVED***
	***REMOVED*** // test getWeights
		[]input***REMOVED***
			***REMOVED***"a", [][]int***REMOVED******REMOVED***100***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"b", [][]int***REMOVED******REMOVED***105***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"c", [][]int***REMOVED******REMOVED***110***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"ß", [][]int***REMOVED******REMOVED***120***REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
		[]check***REMOVED***
			***REMOVED***"a", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			***REMOVED***"b", 1, ColElems***REMOVED***w(105)***REMOVED******REMOVED***,
			***REMOVED***"c", 1, ColElems***REMOVED***w(110)***REMOVED******REMOVED***,
			***REMOVED***"d", 1, ColElems***REMOVED***w(0x50064)***REMOVED******REMOVED***,
			***REMOVED***"ab", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			***REMOVED***"bc", 1, ColElems***REMOVED***w(105)***REMOVED******REMOVED***,
			***REMOVED***"dd", 1, ColElems***REMOVED***w(0x50064)***REMOVED******REMOVED***,
			***REMOVED***"ß", 2, ColElems***REMOVED***w(120)***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED*** // test expansion
		[]input***REMOVED***
			***REMOVED***"u", [][]int***REMOVED******REMOVED***100***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"U", [][]int***REMOVED******REMOVED***100***REMOVED***, ***REMOVED***0, 25***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"w", [][]int***REMOVED******REMOVED***100***REMOVED***, ***REMOVED***100***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"W", [][]int***REMOVED******REMOVED***100***REMOVED***, ***REMOVED***0, 25***REMOVED***, ***REMOVED***100***REMOVED***, ***REMOVED***0, 25***REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
		[]check***REMOVED***
			***REMOVED***"u", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			***REMOVED***"U", 1, ColElems***REMOVED***w(100), w(0, 25)***REMOVED******REMOVED***,
			***REMOVED***"w", 1, ColElems***REMOVED***w(100), w(100)***REMOVED******REMOVED***,
			***REMOVED***"W", 1, ColElems***REMOVED***w(100), w(0, 25), w(100), w(0, 25)***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED*** // test decompose
		[]input***REMOVED***
			***REMOVED***"D", [][]int***REMOVED***pt(104, 8)***REMOVED******REMOVED***,
			***REMOVED***"z", [][]int***REMOVED***pt(130, 8)***REMOVED******REMOVED***,
			***REMOVED***"\u030C", [][]int***REMOVED******REMOVED***0, 40***REMOVED******REMOVED******REMOVED***,                               // Caron
			***REMOVED***"\u01C5", [][]int***REMOVED***pt(104, 9), pt(130, 4), ***REMOVED***0, 40, 0x1F***REMOVED******REMOVED******REMOVED***, // ǅ = D+z+caron
		***REMOVED***,
		[]check***REMOVED***
			***REMOVED***"\u01C5", 2, ColElems***REMOVED***w(pt(104, 9)...), w(pt(130, 4)...), w(0, 40, 0x1F)***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED*** // test basic contraction
		[]input***REMOVED***
			***REMOVED***"a", [][]int***REMOVED******REMOVED***100***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"ab", [][]int***REMOVED******REMOVED***101***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"aab", [][]int***REMOVED******REMOVED***101***REMOVED***, ***REMOVED***101***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"abc", [][]int***REMOVED******REMOVED***102***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"b", [][]int***REMOVED******REMOVED***200***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"c", [][]int***REMOVED******REMOVED***300***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"d", [][]int***REMOVED******REMOVED***400***REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
		[]check***REMOVED***
			***REMOVED***"a", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			***REMOVED***"aa", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			***REMOVED***"aac", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			***REMOVED***"d", 1, ColElems***REMOVED***w(400)***REMOVED******REMOVED***,
			***REMOVED***"ab", 2, ColElems***REMOVED***w(101)***REMOVED******REMOVED***,
			***REMOVED***"abb", 2, ColElems***REMOVED***w(101)***REMOVED******REMOVED***,
			***REMOVED***"aab", 3, ColElems***REMOVED***w(101), w(101)***REMOVED******REMOVED***,
			***REMOVED***"aaba", 3, ColElems***REMOVED***w(101), w(101)***REMOVED******REMOVED***,
			***REMOVED***"abc", 3, ColElems***REMOVED***w(102)***REMOVED******REMOVED***,
			***REMOVED***"abcd", 3, ColElems***REMOVED***w(102)***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED*** // test discontinuous contraction
		append(mods, []input***REMOVED***
			// modifiers; secondary weight equals ccc
			***REMOVED***"\u0316", [][]int***REMOVED******REMOVED***0, 220***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u0317", [][]int***REMOVED******REMOVED***0, 220***REMOVED***, ***REMOVED***0, 220***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u302D", [][]int***REMOVED******REMOVED***0, 222***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u302E", [][]int***REMOVED******REMOVED***0, 225***REMOVED******REMOVED******REMOVED***, // used as starter
			***REMOVED***"\u302F", [][]int***REMOVED******REMOVED***0, 224***REMOVED******REMOVED******REMOVED***, // used as starter
			***REMOVED***"\u18A9", [][]int***REMOVED******REMOVED***0, 228***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u0300", [][]int***REMOVED******REMOVED***0, 230***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u0301", [][]int***REMOVED******REMOVED***0, 230***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u0315", [][]int***REMOVED******REMOVED***0, 232***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u031A", [][]int***REMOVED******REMOVED***0, 232***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u035C", [][]int***REMOVED******REMOVED***0, 233***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u035F", [][]int***REMOVED******REMOVED***0, 233***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u035D", [][]int***REMOVED******REMOVED***0, 234***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u035E", [][]int***REMOVED******REMOVED***0, 234***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u0345", [][]int***REMOVED******REMOVED***0, 240***REMOVED******REMOVED******REMOVED***,

			// starters
			***REMOVED***"a", [][]int***REMOVED******REMOVED***100***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"b", [][]int***REMOVED******REMOVED***200***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"c", [][]int***REMOVED******REMOVED***300***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u03B1", [][]int***REMOVED******REMOVED***900***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\x01", [][]int***REMOVED******REMOVED***0, 0, 0, 0***REMOVED******REMOVED******REMOVED***,

			// contractions
			***REMOVED***"a\u0300", [][]int***REMOVED******REMOVED***101***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"a\u0301", [][]int***REMOVED******REMOVED***102***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"a\u035E", [][]int***REMOVED******REMOVED***110***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"a\u035Eb\u035E", [][]int***REMOVED******REMOVED***115***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"ac\u035Eaca\u035E", [][]int***REMOVED******REMOVED***116***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"a\u035Db\u035D", [][]int***REMOVED******REMOVED***117***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"a\u0301\u035Db", [][]int***REMOVED******REMOVED***120***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"a\u0301\u035F", [][]int***REMOVED******REMOVED***121***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"a\u0301\u035Fb", [][]int***REMOVED******REMOVED***119***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u03B1\u0345", [][]int***REMOVED******REMOVED***901***REMOVED***, ***REMOVED***902***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u302E\u302F", [][]int***REMOVED******REMOVED***0, 131***REMOVED***, ***REMOVED***0, 131***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"\u302F\u18A9", [][]int***REMOVED******REMOVED***0, 130***REMOVED******REMOVED******REMOVED***,
		***REMOVED***...),
		[]check***REMOVED***
			***REMOVED***"a\x01\u0300", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			***REMOVED***"ab", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,                              // closing segment
			***REMOVED***"a\u0316\u0300b", 5, ColElems***REMOVED***w(101), w(0, 220)***REMOVED******REMOVED***,       // closing segment
			***REMOVED***"a\u0316\u0300", 5, ColElems***REMOVED***w(101), w(0, 220)***REMOVED******REMOVED***,        // no closing segment
			***REMOVED***"a\u0316\u0300\u035Cb", 5, ColElems***REMOVED***w(101), w(0, 220)***REMOVED******REMOVED***, // completes before segment end
			***REMOVED***"a\u0316\u0300\u035C", 5, ColElems***REMOVED***w(101), w(0, 220)***REMOVED******REMOVED***,  // completes before segment end

			***REMOVED***"a\u0316\u0301b", 5, ColElems***REMOVED***w(102), w(0, 220)***REMOVED******REMOVED***,       // closing segment
			***REMOVED***"a\u0316\u0301", 5, ColElems***REMOVED***w(102), w(0, 220)***REMOVED******REMOVED***,        // no closing segment
			***REMOVED***"a\u0316\u0301\u035Cb", 5, ColElems***REMOVED***w(102), w(0, 220)***REMOVED******REMOVED***, // completes before segment end
			***REMOVED***"a\u0316\u0301\u035C", 5, ColElems***REMOVED***w(102), w(0, 220)***REMOVED******REMOVED***,  // completes before segment end

			// match blocked by modifier with same ccc
			***REMOVED***"a\u0301\u0315\u031A\u035Fb", 3, ColElems***REMOVED***w(102)***REMOVED******REMOVED***,

			// multiple gaps
			***REMOVED***"a\u0301\u035Db", 6, ColElems***REMOVED***w(120)***REMOVED******REMOVED***,
			***REMOVED***"a\u0301\u035F", 5, ColElems***REMOVED***w(121)***REMOVED******REMOVED***,
			***REMOVED***"a\u0301\u035Fb", 6, ColElems***REMOVED***w(119)***REMOVED******REMOVED***,
			***REMOVED***"a\u0316\u0301\u035F", 7, ColElems***REMOVED***w(121), w(0, 220)***REMOVED******REMOVED***,
			***REMOVED***"a\u0301\u0315\u035Fb", 7, ColElems***REMOVED***w(121), w(0, 232)***REMOVED******REMOVED***,
			***REMOVED***"a\u0316\u0301\u0315\u035Db", 5, ColElems***REMOVED***w(102), w(0, 220)***REMOVED******REMOVED***,
			***REMOVED***"a\u0316\u0301\u0315\u035F", 9, ColElems***REMOVED***w(121), w(0, 220), w(0, 232)***REMOVED******REMOVED***,
			***REMOVED***"a\u0316\u0301\u0315\u035Fb", 9, ColElems***REMOVED***w(121), w(0, 220), w(0, 232)***REMOVED******REMOVED***,
			***REMOVED***"a\u0316\u0301\u0315\u035F\u035D", 9, ColElems***REMOVED***w(121), w(0, 220), w(0, 232)***REMOVED******REMOVED***,
			***REMOVED***"a\u0316\u0301\u0315\u035F\u035Db", 9, ColElems***REMOVED***w(121), w(0, 220), w(0, 232)***REMOVED******REMOVED***,

			// handling of segment overflow
			***REMOVED*** // just fits within segment
				"a" + string(modSeq[:30]) + "\u0301",
				3 + len(string(modSeq[:30])),
				append(ColElems***REMOVED***w(102)***REMOVED***, modW[:30]...),
			***REMOVED***,
			***REMOVED***"a" + string(modSeq[:31]) + "\u0301", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***, // overflow
			***REMOVED***"a" + string(modSeq) + "\u0301", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			***REMOVED*** // just fits within segment with two interstitial runes
				"a" + string(modSeq[:28]) + "\u0301\u0315\u035F",
				7 + len(string(modSeq[:28])),
				append(append(ColElems***REMOVED***w(121)***REMOVED***, modW[:28]...), w(0, 232)),
			***REMOVED***,
			***REMOVED*** // second half does not fit within segment
				"a" + string(modSeq[:29]) + "\u0301\u0315\u035F",
				3 + len(string(modSeq[:29])),
				append(ColElems***REMOVED***w(102)***REMOVED***, modW[:29]...),
			***REMOVED***,

			// discontinuity can only occur in last normalization segment
			***REMOVED***"a\u035Eb\u035E", 6, ColElems***REMOVED***w(115)***REMOVED******REMOVED***,
			***REMOVED***"a\u0316\u035Eb\u035E", 5, ColElems***REMOVED***w(110), w(0, 220)***REMOVED******REMOVED***,
			***REMOVED***"a\u035Db\u035D", 6, ColElems***REMOVED***w(117)***REMOVED******REMOVED***,
			***REMOVED***"a\u0316\u035Db\u035D", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			***REMOVED***"a\u035Eb\u0316\u035E", 8, ColElems***REMOVED***w(115), w(0, 220)***REMOVED******REMOVED***,
			***REMOVED***"a\u035Db\u0316\u035D", 8, ColElems***REMOVED***w(117), w(0, 220)***REMOVED******REMOVED***,
			***REMOVED***"ac\u035Eaca\u035E", 9, ColElems***REMOVED***w(116)***REMOVED******REMOVED***,
			***REMOVED***"a\u0316c\u035Eaca\u035E", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			***REMOVED***"ac\u035Eac\u0316a\u035E", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,

			// expanding contraction
			***REMOVED***"\u03B1\u0345", 4, ColElems***REMOVED***w(901), w(902)***REMOVED******REMOVED***,

			// Theoretical possibilities
			// contraction within a gap
			***REMOVED***"a\u302F\u18A9\u0301", 9, ColElems***REMOVED***w(102), w(0, 130)***REMOVED******REMOVED***,
			// expansion within a gap
			***REMOVED***"a\u0317\u0301", 5, ColElems***REMOVED***w(102), w(0, 220), w(0, 220)***REMOVED******REMOVED***,
			// repeating CCC blocks last modifier
			***REMOVED***"a\u302E\u302F\u0301", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			// The trailing combining characters (with lower CCC) should block the first one.
			// TODO: make the following pass.
			// ***REMOVED***"a\u035E\u0316\u0316", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
			***REMOVED***"a\u035F\u035Eb", 5, ColElems***REMOVED***w(110), w(0, 233)***REMOVED******REMOVED***,
			// Last combiner should match after normalization.
			// TODO: make the following pass.
			// ***REMOVED***"a\u035D\u0301", 3, ColElems***REMOVED***w(102), w(0, 234)***REMOVED******REMOVED***,
			// The first combiner is blocking the second one as they have the same CCC.
			***REMOVED***"a\u035D\u035Eb", 1, ColElems***REMOVED***w(100)***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestAppendNext(t *testing.T) ***REMOVED***
	for i, tt := range appendNextTests ***REMOVED***
		c, err := makeTable(tt.in)
		if err != nil ***REMOVED***
			t.Errorf("%d: error creating table: %v", i, err)
			continue
		***REMOVED***
		for j, chk := range tt.chk ***REMOVED***
			ws, n := c.t.AppendNext(nil, []byte(chk.in))
			if n != chk.n ***REMOVED***
				t.Errorf("%d:%d: bytes consumed was %d; want %d", i, j, n, chk.n)
			***REMOVED***
			out := convertFromWeights(chk.out)
			if len(ws) != len(out) ***REMOVED***
				t.Errorf("%d:%d: len(ws) was %d; want %d (%X vs %X)\n%X", i, j, len(ws), len(out), ws, out, chk.in)
				continue
			***REMOVED***
			for k, w := range ws ***REMOVED***
				w, _ = colltab.MakeElem(w.Primary(), w.Secondary(), int(w.Tertiary()), 0)
				if w != out[k] ***REMOVED***
					t.Errorf("%d:%d: Weights %d was %X; want %X", i, j, k, w, out[k])
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
