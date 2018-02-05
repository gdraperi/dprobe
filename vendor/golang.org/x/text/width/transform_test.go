// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package width

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/transform"
)

func foldRune(r rune) (folded rune, ok bool) ***REMOVED***
	alt, ok := mapRunes[r]
	if ok && alt.e&tagNeedsFold != 0 ***REMOVED***
		return alt.r, true
	***REMOVED***
	return r, false
***REMOVED***

func widenRune(r rune) (wide rune, ok bool) ***REMOVED***
	alt, ok := mapRunes[r]
	if k := alt.e.kind(); k == EastAsianHalfwidth || k == EastAsianNarrow ***REMOVED***
		return alt.r, true
	***REMOVED***
	return r, false
***REMOVED***

func narrowRune(r rune) (narrow rune, ok bool) ***REMOVED***
	alt, ok := mapRunes[r]
	if k := alt.e.kind(); k == EastAsianFullwidth || k == EastAsianWide || k == EastAsianAmbiguous ***REMOVED***
		return alt.r, true
	***REMOVED***
	return r, false
***REMOVED***

func TestFoldSingleRunes(t *testing.T) ***REMOVED***
	for r := rune(0); r < 0x1FFFF; r++ ***REMOVED***
		if loSurrogate <= r && r <= hiSurrogate ***REMOVED***
			continue
		***REMOVED***
		x, _ := foldRune(r)
		want := string(x)
		got := Fold.String(string(r))
		if got != want ***REMOVED***
			t.Errorf("Fold().String(%U) = %+q; want %+q", r, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

type transformTest struct ***REMOVED***
	desc    string
	src     string
	nBuf    int
	nDst    int
	atEOF   bool
	dst     string
	nSrc    int
	err     error
	nSpan   int
	errSpan error
***REMOVED***

func (tc *transformTest) doTest(t *testing.T, tr Transformer) ***REMOVED***
	testtext.Run(t, tc.desc, func(t *testing.T) ***REMOVED***
		b := make([]byte, tc.nBuf)
		nDst, nSrc, err := tr.Transform(b, []byte(tc.src), tc.atEOF)
		if got := string(b[:nDst]); got != tc.dst[:nDst] ***REMOVED***
			t.Errorf("dst was %+q; want %+q", got, tc.dst)
		***REMOVED***
		if nDst != tc.nDst ***REMOVED***
			t.Errorf("nDst was %d; want %d", nDst, tc.nDst)
		***REMOVED***
		if nSrc != tc.nSrc ***REMOVED***
			t.Errorf("nSrc was %d; want %d", nSrc, tc.nSrc)
		***REMOVED***
		if err != tc.err ***REMOVED***
			t.Errorf("error was %v; want %v", err, tc.err)
		***REMOVED***
		if got := tr.String(tc.src); got != tc.dst ***REMOVED***
			t.Errorf("String(%q) = %q; want %q", tc.src, got, tc.dst)
		***REMOVED***
		n, err := tr.Span([]byte(tc.src), tc.atEOF)
		if n != tc.nSpan || err != tc.errSpan ***REMOVED***
			t.Errorf("Span: got %d, %v; want %d, %v", n, err, tc.nSpan, tc.errSpan)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestFold(t *testing.T) ***REMOVED***
	for _, tc := range []transformTest***REMOVED******REMOVED***
		desc:    "empty",
		src:     "",
		nBuf:    10,
		dst:     "",
		nDst:    0,
		nSrc:    0,
		atEOF:   false,
		err:     nil,
		nSpan:   0,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "short source 1",
		src:     "a\xc2",
		nBuf:    10,
		dst:     "a\xc2",
		nDst:    1,
		nSrc:    1,
		atEOF:   false,
		err:     transform.ErrShortSrc,
		nSpan:   1,
		errSpan: transform.ErrShortSrc,
	***REMOVED***, ***REMOVED***
		desc:    "short source 2",
		src:     "a\xe0\x80",
		nBuf:    10,
		dst:     "a\xe0\x80",
		nDst:    1,
		nSrc:    1,
		atEOF:   false,
		err:     transform.ErrShortSrc,
		nSpan:   1,
		errSpan: transform.ErrShortSrc,
	***REMOVED***, ***REMOVED***
		desc:    "incomplete but terminated source 1",
		src:     "a\xc2",
		nBuf:    10,
		dst:     "a\xc2",
		nDst:    2,
		nSrc:    2,
		atEOF:   true,
		err:     nil,
		nSpan:   2,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "incomplete but terminated source 2",
		src:     "a\xe0\x80",
		nBuf:    10,
		dst:     "a\xe0\x80",
		nDst:    3,
		nSrc:    3,
		atEOF:   true,
		err:     nil,
		nSpan:   3,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "exact fit dst",
		src:     "a\uff01",
		nBuf:    2,
		dst:     "a!",
		nDst:    2,
		nSrc:    4,
		atEOF:   false,
		err:     nil,
		nSpan:   1,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "exact fit dst and src ascii",
		src:     "ab",
		nBuf:    2,
		dst:     "ab",
		nDst:    2,
		nSrc:    2,
		atEOF:   true,
		err:     nil,
		nSpan:   2,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "empty dst",
		src:     "\u0300",
		nBuf:    0,
		dst:     "\u0300",
		nDst:    0,
		nSrc:    0,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   2,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "empty dst ascii",
		src:     "a",
		nBuf:    0,
		dst:     "a",
		nDst:    0,
		nSrc:    0,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   1,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "short dst 1",
		src:     "a\uffe0", // ￠
		nBuf:    2,
		dst:     "a\u00a2", // ¢
		nDst:    1,
		nSrc:    1,
		atEOF:   false,
		err:     transform.ErrShortDst,
		nSpan:   1,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "short dst 2",
		src:     "不夠",
		nBuf:    3,
		dst:     "不夠",
		nDst:    3,
		nSrc:    3,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   6,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "short dst fast path",
		src:     "fast",
		nDst:    3,
		dst:     "fast",
		nBuf:    3,
		nSrc:    3,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   4,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "short dst larger buffer",
		src:     "\uff21" + strings.Repeat("0", 127) + "B",
		nBuf:    128,
		dst:     "A" + strings.Repeat("0", 127) + "B",
		nDst:    128,
		nSrc:    130,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "fast path alternation",
		src:     "fast路徑fast路徑",
		nBuf:    20,
		dst:     "fast路徑fast路徑",
		nDst:    20,
		nSrc:    20,
		atEOF:   true,
		err:     nil,
		nSpan:   20,
		errSpan: nil,
	***REMOVED******REMOVED*** ***REMOVED***
		tc.doTest(t, Fold)
	***REMOVED***
***REMOVED***

func TestWidenSingleRunes(t *testing.T) ***REMOVED***
	for r := rune(0); r < 0x1FFFF; r++ ***REMOVED***
		if loSurrogate <= r && r <= hiSurrogate ***REMOVED***
			continue
		***REMOVED***
		alt, _ := widenRune(r)
		want := string(alt)
		got := Widen.String(string(r))
		if got != want ***REMOVED***
			t.Errorf("Widen().String(%U) = %+q; want %+q", r, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWiden(t *testing.T) ***REMOVED***
	for _, tc := range []transformTest***REMOVED******REMOVED***
		desc:    "empty",
		src:     "",
		nBuf:    10,
		dst:     "",
		nDst:    0,
		nSrc:    0,
		atEOF:   false,
		err:     nil,
		nSpan:   0,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "short source 1",
		src:     "a\xc2",
		nBuf:    10,
		dst:     "ａ\xc2",
		nDst:    3,
		nSrc:    1,
		atEOF:   false,
		err:     transform.ErrShortSrc,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "short source 2",
		src:     "a\xe0\x80",
		nBuf:    10,
		dst:     "ａ\xe0\x80",
		nDst:    3,
		nSrc:    1,
		atEOF:   false,
		err:     transform.ErrShortSrc,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "incomplete but terminated source 1",
		src:     "a\xc2",
		nBuf:    10,
		dst:     "ａ\xc2",
		nDst:    4,
		nSrc:    2,
		atEOF:   true,
		err:     nil,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "incomplete but terminated source 2",
		src:     "a\xe0\x80",
		nBuf:    10,
		dst:     "ａ\xe0\x80",
		nDst:    5,
		nSrc:    3,
		atEOF:   true,
		err:     nil,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "short source 1 some span",
		src:     "ａ\xc2",
		nBuf:    10,
		dst:     "ａ\xc2",
		nDst:    3,
		nSrc:    3,
		atEOF:   false,
		err:     transform.ErrShortSrc,
		nSpan:   3,
		errSpan: transform.ErrShortSrc,
	***REMOVED***, ***REMOVED***
		desc:    "short source 2 some span",
		src:     "ａ\xe0\x80",
		nBuf:    10,
		dst:     "ａ\xe0\x80",
		nDst:    3,
		nSrc:    3,
		atEOF:   false,
		err:     transform.ErrShortSrc,
		nSpan:   3,
		errSpan: transform.ErrShortSrc,
	***REMOVED***, ***REMOVED***
		desc:    "incomplete but terminated source 1 some span",
		src:     "ａ\xc2",
		nBuf:    10,
		dst:     "ａ\xc2",
		nDst:    4,
		nSrc:    4,
		atEOF:   true,
		err:     nil,
		nSpan:   4,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "incomplete but terminated source 2 some span",
		src:     "ａ\xe0\x80",
		nBuf:    10,
		dst:     "ａ\xe0\x80",
		nDst:    5,
		nSrc:    5,
		atEOF:   true,
		err:     nil,
		nSpan:   5,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "exact fit dst",
		src:     "a!",
		nBuf:    6,
		dst:     "ａ\uff01",
		nDst:    6,
		nSrc:    2,
		atEOF:   false,
		err:     nil,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "empty dst",
		src:     "\u0300",
		nBuf:    0,
		dst:     "\u0300",
		nDst:    0,
		nSrc:    0,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   2,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "empty dst ascii",
		src:     "a",
		nBuf:    0,
		dst:     "ａ",
		nDst:    0,
		nSrc:    0,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "short dst 1",
		src:     "a\uffe0",
		nBuf:    4,
		dst:     "ａ\uffe0",
		nDst:    3,
		nSrc:    1,
		atEOF:   false,
		err:     transform.ErrShortDst,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "short dst 2",
		src:     "不夠",
		nBuf:    3,
		dst:     "不夠",
		nDst:    3,
		nSrc:    3,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   6,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "short dst ascii",
		src:     "ascii",
		nBuf:    3,
		dst:     "ａｓｃｉｉ", // U+ff41, ...
		nDst:    3,
		nSrc:    1,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "ambiguous",
		src:     "\uffe9",
		nBuf:    4,
		dst:     "\u2190",
		nDst:    3,
		nSrc:    3,
		atEOF:   false,
		err:     nil,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED******REMOVED*** ***REMOVED***
		tc.doTest(t, Widen)
	***REMOVED***
***REMOVED***

func TestNarrowSingleRunes(t *testing.T) ***REMOVED***
	for r := rune(0); r < 0x1FFFF; r++ ***REMOVED***
		if loSurrogate <= r && r <= hiSurrogate ***REMOVED***
			continue
		***REMOVED***
		alt, _ := narrowRune(r)
		want := string(alt)
		got := Narrow.String(string(r))
		if got != want ***REMOVED***
			t.Errorf("Narrow().String(%U) = %+q; want %+q", r, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNarrow(t *testing.T) ***REMOVED***
	for _, tc := range []transformTest***REMOVED******REMOVED***
		desc:    "empty",
		src:     "",
		nBuf:    10,
		dst:     "",
		nDst:    0,
		nSrc:    0,
		atEOF:   false,
		err:     nil,
		nSpan:   0,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "short source 1",
		src:     "a\xc2",
		nBuf:    10,
		dst:     "a\xc2",
		nDst:    1,
		nSrc:    1,
		atEOF:   false,
		err:     transform.ErrShortSrc,
		nSpan:   1,
		errSpan: transform.ErrShortSrc,
	***REMOVED***, ***REMOVED***
		desc:    "short source 2",
		src:     "ａ\xe0\x80",
		nBuf:    10,
		dst:     "a\xe0\x80",
		nDst:    1,
		nSrc:    3,
		atEOF:   false,
		err:     transform.ErrShortSrc,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "incomplete but terminated source 1",
		src:     "ａ\xc2",
		nBuf:    10,
		dst:     "a\xc2",
		nDst:    2,
		nSrc:    4,
		atEOF:   true,
		err:     nil,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "incomplete but terminated source 2",
		src:     "ａ\xe0\x80",
		nBuf:    10,
		dst:     "a\xe0\x80",
		nDst:    3,
		nSrc:    5,
		atEOF:   true,
		err:     nil,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "exact fit dst",
		src:     "ａ\uff01",
		nBuf:    2,
		dst:     "a!",
		nDst:    2,
		nSrc:    6,
		atEOF:   false,
		err:     nil,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "exact fit dst some span",
		src:     "a\uff01",
		nBuf:    2,
		dst:     "a!",
		nDst:    2,
		nSrc:    4,
		atEOF:   false,
		err:     nil,
		nSpan:   1,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "empty dst",
		src:     "\u0300",
		nBuf:    0,
		dst:     "\u0300",
		nDst:    0,
		nSrc:    0,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   2,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "empty dst ascii",
		src:     "a",
		nBuf:    0,
		dst:     "a",
		nDst:    0,
		nSrc:    0,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   1,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "short dst 1",
		src:     "ａ\uffe0", // ￠
		nBuf:    2,
		dst:     "a\u00a2", // ¢
		nDst:    1,
		nSrc:    3,
		atEOF:   false,
		err:     transform.ErrShortDst,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "short dst 2",
		src:     "不夠",
		nBuf:    3,
		dst:     "不夠",
		nDst:    3,
		nSrc:    3,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   6,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		// Create a narrow variant of ambiguous runes, if they exist.
		desc:    "ambiguous",
		src:     "\u2190",
		nBuf:    4,
		dst:     "\uffe9",
		nDst:    3,
		nSrc:    3,
		atEOF:   false,
		err:     nil,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "short dst fast path",
		src:     "fast",
		nBuf:    3,
		dst:     "fast",
		nDst:    3,
		nSrc:    3,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   4,
		errSpan: nil,
	***REMOVED***, ***REMOVED***
		desc:    "short dst larger buffer",
		src:     "\uff21" + strings.Repeat("0", 127) + "B",
		nBuf:    128,
		dst:     "A" + strings.Repeat("0", 127) + "B",
		nDst:    128,
		nSrc:    130,
		atEOF:   true,
		err:     transform.ErrShortDst,
		nSpan:   0,
		errSpan: transform.ErrEndOfSpan,
	***REMOVED***, ***REMOVED***
		desc:    "fast path alternation",
		src:     "fast路徑fast路徑",
		nBuf:    20,
		dst:     "fast路徑fast路徑",
		nDst:    20,
		nSrc:    20,
		atEOF:   true,
		err:     nil,
		nSpan:   20,
		errSpan: nil,
	***REMOVED******REMOVED*** ***REMOVED***
		tc.doTest(t, Narrow)
	***REMOVED***
***REMOVED***

func bench(b *testing.B, t Transformer, s string) ***REMOVED***
	dst := make([]byte, 1024)
	src := []byte(s)
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		t.Transform(dst, src, true)
	***REMOVED***
***REMOVED***

func changingRunes(f func(r rune) (rune, bool)) string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***
	for r := rune(0); r <= 0xFFFF; r++ ***REMOVED***
		if _, ok := foldRune(r); ok ***REMOVED***
			buf.WriteRune(r)
		***REMOVED***
	***REMOVED***
	return buf.String()
***REMOVED***

func BenchmarkFoldASCII(b *testing.B) ***REMOVED***
	bench(b, Fold, testtext.ASCII)
***REMOVED***

func BenchmarkFoldCJK(b *testing.B) ***REMOVED***
	bench(b, Fold, testtext.CJK)
***REMOVED***

func BenchmarkFoldNonCanonical(b *testing.B) ***REMOVED***
	bench(b, Fold, changingRunes(foldRune))
***REMOVED***

func BenchmarkFoldOther(b *testing.B) ***REMOVED***
	bench(b, Fold, testtext.TwoByteUTF8+testtext.ThreeByteUTF8)
***REMOVED***

func BenchmarkWideASCII(b *testing.B) ***REMOVED***
	bench(b, Widen, testtext.ASCII)
***REMOVED***

func BenchmarkWideCJK(b *testing.B) ***REMOVED***
	bench(b, Widen, testtext.CJK)
***REMOVED***

func BenchmarkWideNonCanonical(b *testing.B) ***REMOVED***
	bench(b, Widen, changingRunes(widenRune))
***REMOVED***

func BenchmarkWideOther(b *testing.B) ***REMOVED***
	bench(b, Widen, testtext.TwoByteUTF8+testtext.ThreeByteUTF8)
***REMOVED***

func BenchmarkNarrowASCII(b *testing.B) ***REMOVED***
	bench(b, Narrow, testtext.ASCII)
***REMOVED***

func BenchmarkNarrowCJK(b *testing.B) ***REMOVED***
	bench(b, Narrow, testtext.CJK)
***REMOVED***

func BenchmarkNarrowNonCanonical(b *testing.B) ***REMOVED***
	bench(b, Narrow, changingRunes(narrowRune))
***REMOVED***

func BenchmarkNarrowOther(b *testing.B) ***REMOVED***
	bench(b, Narrow, testtext.TwoByteUTF8+testtext.ThreeByteUTF8)
***REMOVED***
