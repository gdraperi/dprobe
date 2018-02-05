// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
	"unicode/utf8"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/transform"
)

var (
	testn = flag.Int("testn", -1, "specific test number to run or -1 for all")
)

// pc replaces any rune r that is repeated n times, for n > 1, with r***REMOVED***n***REMOVED***.
func pc(s string) []byte ***REMOVED***
	b := bytes.NewBuffer(make([]byte, 0, len(s)))
	for i := 0; i < len(s); ***REMOVED***
		r, sz := utf8.DecodeRuneInString(s[i:])
		n := 0
		if sz == 1 ***REMOVED***
			// Special-case one-byte case to handle repetition for invalid UTF-8.
			for c := s[i]; i+n < len(s) && s[i+n] == c; n++ ***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for _, r2 := range s[i:] ***REMOVED***
				if r2 != r ***REMOVED***
					break
				***REMOVED***
				n++
			***REMOVED***
		***REMOVED***
		b.WriteString(s[i : i+sz])
		if n > 1 ***REMOVED***
			fmt.Fprintf(b, "***REMOVED***%d***REMOVED***", n)
		***REMOVED***
		i += sz * n
	***REMOVED***
	return b.Bytes()
***REMOVED***

// pidx finds the index from which two strings start to differ, plus context.
// It returns the index and ellipsis if the index is greater than 0.
func pidx(a, b string) (i int, prefix string) ***REMOVED***
	for ; i < len(a) && i < len(b) && a[i] == b[i]; i++ ***REMOVED***
	***REMOVED***
	if i < 8 ***REMOVED***
		return 0, ""
	***REMOVED***
	i -= 3 // ensure taking at least one full rune before the difference.
	for k := i - 7; i > k && !utf8.RuneStart(a[i]); i-- ***REMOVED***
	***REMOVED***
	return i, "..."
***REMOVED***

type PositionTest struct ***REMOVED***
	input  string
	pos    int
	buffer string // expected contents of reorderBuffer, if applicable
***REMOVED***

type positionFunc func(rb *reorderBuffer, s string) (int, []byte)

func runPosTests(t *testing.T, name string, f Form, fn positionFunc, tests []PositionTest) ***REMOVED***
	rb := reorderBuffer***REMOVED******REMOVED***
	rb.init(f, nil)
	for i, test := range tests ***REMOVED***
		rb.reset()
		rb.src = inputString(test.input)
		rb.nsrc = len(test.input)
		pos, out := fn(&rb, test.input)
		if pos != test.pos ***REMOVED***
			t.Errorf("%s:%d: position is %d; want %d", name, i, pos, test.pos)
		***REMOVED***
		if outs := string(out); outs != test.buffer ***REMOVED***
			k, pfx := pidx(outs, test.buffer)
			t.Errorf("%s:%d: buffer \nwas  %s%+q; \nwant %s%+q", name, i, pfx, pc(outs[k:]), pfx, pc(test.buffer[k:]))
		***REMOVED***
	***REMOVED***
***REMOVED***

func grave(n int) string ***REMOVED***
	return rep(0x0300, n)
***REMOVED***

func rep(r rune, n int) string ***REMOVED***
	return strings.Repeat(string(r), n)
***REMOVED***

const segSize = maxByteBufferSize

var cgj = GraphemeJoiner

var decomposeSegmentTests = []PositionTest***REMOVED***
	// illegal runes
	***REMOVED***"\xC2", 0, ""***REMOVED***,
	***REMOVED***"\xC0", 1, "\xC0"***REMOVED***,
	***REMOVED***"\u00E0\x80", 2, "\u0061\u0300"***REMOVED***,
	// starter
	***REMOVED***"a", 1, "a"***REMOVED***,
	***REMOVED***"ab", 1, "a"***REMOVED***,
	// starter + composing
	***REMOVED***"a\u0300", 3, "a\u0300"***REMOVED***,
	***REMOVED***"a\u0300b", 3, "a\u0300"***REMOVED***,
	// with decomposition
	***REMOVED***"\u00C0", 2, "A\u0300"***REMOVED***,
	***REMOVED***"\u00C0b", 2, "A\u0300"***REMOVED***,
	// long
	***REMOVED***grave(31), 60, grave(30) + cgj***REMOVED***,
	***REMOVED***"a" + grave(31), 61, "a" + grave(30) + cgj***REMOVED***,

	// Stability tests: see http://www.unicode.org/review/pr-29.html.
	// U+0300 COMBINING GRAVE ACCENT;Mn;230;NSM;;;;;N;NON-SPACING GRAVE;;;;
	// U+0B47 ORIYA VOWEL SIGN E;Mc;0;L;;;;;N;;;;;
	// U+0B3E ORIYA VOWEL SIGN AA;Mc;0;L;;;;;N;;;;;
	// U+1100 HANGUL CHOSEONG KIYEOK;Lo;0;L;;;;;N;;;;;
	// U+1161 HANGUL JUNGSEONG A;Lo;0;L;;;;;N;;;;;
	***REMOVED***"\u0B47\u0300\u0B3E", 8, "\u0B47\u0300\u0B3E"***REMOVED***,
	***REMOVED***"\u1100\u0300\u1161", 8, "\u1100\u0300\u1161"***REMOVED***,
	***REMOVED***"\u0B47\u0B3E", 6, "\u0B47\u0B3E"***REMOVED***,
	***REMOVED***"\u1100\u1161", 6, "\u1100\u1161"***REMOVED***,

	// U+04DA MALAYALAM VOWEL SIGN O;Mc;0;L;0D46 0D3E;;;;N;;;;;
	// Sequence of decomposing characters that are starters and modifiers.
	***REMOVED***"\u0d4a" + strings.Repeat("\u0d3e", 31), 90, "\u0d46" + strings.Repeat("\u0d3e", 30) + cgj***REMOVED***,

	***REMOVED***grave(30), 60, grave(30)***REMOVED***,
	// U+FF9E is a starter, but decomposes to U+3099, which is not.
	***REMOVED***grave(30) + "\uff9e", 60, grave(30) + cgj***REMOVED***,
	// ends with incomplete UTF-8 encoding
	***REMOVED***"\xCC", 0, ""***REMOVED***,
	***REMOVED***"\u0300\xCC", 2, "\u0300"***REMOVED***,
***REMOVED***

func decomposeSegmentF(rb *reorderBuffer, s string) (int, []byte) ***REMOVED***
	rb.initString(NFD, s)
	rb.setFlusher(nil, appendFlush)
	p := decomposeSegment(rb, 0, true)
	return p, rb.out
***REMOVED***

func TestDecomposeSegment(t *testing.T) ***REMOVED***
	runPosTests(t, "TestDecomposeSegment", NFC, decomposeSegmentF, decomposeSegmentTests)
***REMOVED***

var firstBoundaryTests = []PositionTest***REMOVED***
	// no boundary
	***REMOVED***"", -1, ""***REMOVED***,
	***REMOVED***"\u0300", -1, ""***REMOVED***,
	***REMOVED***"\x80\x80", -1, ""***REMOVED***,
	// illegal runes
	***REMOVED***"\xff", 0, ""***REMOVED***,
	***REMOVED***"\u0300\xff", 2, ""***REMOVED***,
	***REMOVED***"\u0300\xc0\x80\x80", 2, ""***REMOVED***,
	// boundaries
	***REMOVED***"a", 0, ""***REMOVED***,
	***REMOVED***"\u0300a", 2, ""***REMOVED***,
	// Hangul
	***REMOVED***"\u1103\u1161", 0, ""***REMOVED***,
	***REMOVED***"\u110B\u1173\u11B7", 0, ""***REMOVED***,
	***REMOVED***"\u1161\u110B\u1173\u11B7", 3, ""***REMOVED***,
	***REMOVED***"\u1173\u11B7\u1103\u1161", 6, ""***REMOVED***,
	// too many combining characters.
	***REMOVED***grave(maxNonStarters - 1), -1, ""***REMOVED***,
	***REMOVED***grave(maxNonStarters), 60, ""***REMOVED***,
	***REMOVED***grave(maxNonStarters + 1), 60, ""***REMOVED***,
***REMOVED***

func firstBoundaryF(rb *reorderBuffer, s string) (int, []byte) ***REMOVED***
	return rb.f.form.FirstBoundary([]byte(s)), nil
***REMOVED***

func firstBoundaryStringF(rb *reorderBuffer, s string) (int, []byte) ***REMOVED***
	return rb.f.form.FirstBoundaryInString(s), nil
***REMOVED***

func TestFirstBoundary(t *testing.T) ***REMOVED***
	runPosTests(t, "TestFirstBoundary", NFC, firstBoundaryF, firstBoundaryTests)
	runPosTests(t, "TestFirstBoundaryInString", NFC, firstBoundaryStringF, firstBoundaryTests)
***REMOVED***

func TestNextBoundary(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		input string
		atEOF bool
		want  int
	***REMOVED******REMOVED***
		// no boundary
		***REMOVED***"", true, 0***REMOVED***,
		***REMOVED***"", false, -1***REMOVED***,
		***REMOVED***"\u0300", true, 2***REMOVED***,
		***REMOVED***"\u0300", false, -1***REMOVED***,
		***REMOVED***"\x80\x80", true, 1***REMOVED***,
		***REMOVED***"\x80\x80", false, 1***REMOVED***,
		// illegal runes
		***REMOVED***"\xff", false, 1***REMOVED***,
		***REMOVED***"\u0300\xff", false, 2***REMOVED***,
		***REMOVED***"\u0300\xc0\x80\x80", false, 2***REMOVED***,
		***REMOVED***"\xc2\x80\x80", false, 2***REMOVED***,
		***REMOVED***"\xc2", false, -1***REMOVED***,
		***REMOVED***"\xc2", true, 1***REMOVED***,
		***REMOVED***"a\u0300\xc2", false, -1***REMOVED***,
		***REMOVED***"a\u0300\xc2", true, 3***REMOVED***,
		// boundaries
		***REMOVED***"a", true, 1***REMOVED***,
		***REMOVED***"a", false, -1***REMOVED***,
		***REMOVED***"aa", false, 1***REMOVED***,
		***REMOVED***"\u0300", true, 2***REMOVED***,
		***REMOVED***"\u0300", false, -1***REMOVED***,
		***REMOVED***"\u0300a", false, 2***REMOVED***,
		// Hangul
		***REMOVED***"\u1103\u1161", true, 6***REMOVED***,
		***REMOVED***"\u1103\u1161", false, -1***REMOVED***,
		***REMOVED***"\u110B\u1173\u11B7", false, -1***REMOVED***,
		***REMOVED***"\u110B\u1173\u11B7\u110B\u1173\u11B7", false, 9***REMOVED***,
		***REMOVED***"\u1161\u110B\u1173\u11B7", false, 3***REMOVED***,
		***REMOVED***"\u1173\u11B7\u1103\u1161", false, 6***REMOVED***,
		// too many combining characters.
		***REMOVED***grave(maxNonStarters - 1), false, -1***REMOVED***,
		***REMOVED***grave(maxNonStarters), false, 60***REMOVED***,
		***REMOVED***grave(maxNonStarters + 1), false, 60***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		if got := NFC.NextBoundary([]byte(tc.input), tc.atEOF); got != tc.want ***REMOVED***
			t.Errorf("NextBoundary(%+q, %v) = %d; want %d", tc.input, tc.atEOF, got, tc.want)
		***REMOVED***
		if got := NFC.NextBoundaryInString(tc.input, tc.atEOF); got != tc.want ***REMOVED***
			t.Errorf("NextBoundaryInString(%+q, %v) = %d; want %d", tc.input, tc.atEOF, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

var decomposeToLastTests = []PositionTest***REMOVED***
	// ends with inert character
	***REMOVED***"Hello!", 6, ""***REMOVED***,
	***REMOVED***"\u0632", 2, ""***REMOVED***,
	***REMOVED***"a\u0301\u0635", 5, ""***REMOVED***,
	// ends with non-inert starter
	***REMOVED***"a", 0, "a"***REMOVED***,
	***REMOVED***"a\u0301a", 3, "a"***REMOVED***,
	***REMOVED***"a\u0301\u03B9", 3, "\u03B9"***REMOVED***,
	***REMOVED***"a\u0327", 0, "a\u0327"***REMOVED***,
	// illegal runes
	***REMOVED***"\xFF", 1, ""***REMOVED***,
	***REMOVED***"aa\xFF", 3, ""***REMOVED***,
	***REMOVED***"\xC0\x80\x80", 3, ""***REMOVED***,
	***REMOVED***"\xCC\x80\x80", 3, ""***REMOVED***,
	// ends with incomplete UTF-8 encoding
	***REMOVED***"a\xCC", 2, ""***REMOVED***,
	// ends with combining characters
	***REMOVED***"\u0300\u0301", 0, "\u0300\u0301"***REMOVED***,
	***REMOVED***"a\u0300\u0301", 0, "a\u0300\u0301"***REMOVED***,
	***REMOVED***"a\u0301\u0308", 0, "a\u0301\u0308"***REMOVED***,
	***REMOVED***"a\u0308\u0301", 0, "a\u0308\u0301"***REMOVED***,
	***REMOVED***"aaaa\u0300\u0301", 3, "a\u0300\u0301"***REMOVED***,
	***REMOVED***"\u0300a\u0300\u0301", 2, "a\u0300\u0301"***REMOVED***,
	***REMOVED***"\u00C0", 0, "A\u0300"***REMOVED***,
	***REMOVED***"a\u00C0", 1, "A\u0300"***REMOVED***,
	// decomposing
	***REMOVED***"a\u0300\u00E0", 3, "a\u0300"***REMOVED***,
	// multisegment decompositions (flushes leading segments)
	***REMOVED***"a\u0300\uFDC0", 7, "\u064A"***REMOVED***,
	***REMOVED***"\uFDC0" + grave(29), 4, "\u064A" + grave(29)***REMOVED***,
	***REMOVED***"\uFDC0" + grave(30), 4, "\u064A" + grave(30)***REMOVED***,
	***REMOVED***"\uFDC0" + grave(31), 5, grave(30)***REMOVED***,
	***REMOVED***"\uFDFA" + grave(14), 31, "\u0645" + grave(14)***REMOVED***,
	// Overflow
	***REMOVED***"\u00E0" + grave(29), 0, "a" + grave(30)***REMOVED***,
	***REMOVED***"\u00E0" + grave(30), 2, grave(30)***REMOVED***,
	// Hangul
	***REMOVED***"a\u1103", 1, "\u1103"***REMOVED***,
	***REMOVED***"a\u110B", 1, "\u110B"***REMOVED***,
	***REMOVED***"a\u110B\u1173", 1, "\u110B\u1173"***REMOVED***,
	// See comment in composition.go:compBoundaryAfter.
	***REMOVED***"a\u110B\u1173\u11B7", 1, "\u110B\u1173\u11B7"***REMOVED***,
	***REMOVED***"a\uC73C", 1, "\u110B\u1173"***REMOVED***,
	***REMOVED***"다음", 3, "\u110B\u1173\u11B7"***REMOVED***,
	***REMOVED***"다", 0, "\u1103\u1161"***REMOVED***,
	***REMOVED***"\u1103\u1161\u110B\u1173\u11B7", 6, "\u110B\u1173\u11B7"***REMOVED***,
	***REMOVED***"\u110B\u1173\u11B7\u1103\u1161", 9, "\u1103\u1161"***REMOVED***,
	***REMOVED***"다음음", 6, "\u110B\u1173\u11B7"***REMOVED***,
	***REMOVED***"음다다", 6, "\u1103\u1161"***REMOVED***,
	// maximized buffer
	***REMOVED***"a" + grave(30), 0, "a" + grave(30)***REMOVED***,
	// Buffer overflow
	***REMOVED***"a" + grave(31), 3, grave(30)***REMOVED***,
	// weird UTF-8
	***REMOVED***"a\u0300\u11B7", 0, "a\u0300\u11B7"***REMOVED***,
***REMOVED***

func decomposeToLast(rb *reorderBuffer, s string) (int, []byte) ***REMOVED***
	rb.setFlusher([]byte(s), appendFlush)
	decomposeToLastBoundary(rb)
	buf := rb.flush(nil)
	return len(rb.out), buf
***REMOVED***

func TestDecomposeToLastBoundary(t *testing.T) ***REMOVED***
	runPosTests(t, "TestDecomposeToLastBoundary", NFKC, decomposeToLast, decomposeToLastTests)
***REMOVED***

var lastBoundaryTests = []PositionTest***REMOVED***
	// ends with inert character
	***REMOVED***"Hello!", 6, ""***REMOVED***,
	***REMOVED***"\u0632", 2, ""***REMOVED***,
	// ends with non-inert starter
	***REMOVED***"a", 0, ""***REMOVED***,
	// illegal runes
	***REMOVED***"\xff", 1, ""***REMOVED***,
	***REMOVED***"aa\xff", 3, ""***REMOVED***,
	***REMOVED***"a\xff\u0300", 1, ""***REMOVED***, // TODO: should probably be 2.
	***REMOVED***"\xc0\x80\x80", 3, ""***REMOVED***,
	***REMOVED***"\xc0\x80\x80\u0300", 3, ""***REMOVED***,
	// ends with incomplete UTF-8 encoding
	***REMOVED***"\xCC", -1, ""***REMOVED***,
	***REMOVED***"\xE0\x80", -1, ""***REMOVED***,
	***REMOVED***"\xF0\x80\x80", -1, ""***REMOVED***,
	***REMOVED***"a\xCC", 0, ""***REMOVED***,
	***REMOVED***"\x80\xCC", 1, ""***REMOVED***,
	***REMOVED***"\xCC\xCC", 1, ""***REMOVED***,
	// ends with combining characters
	***REMOVED***"a\u0300\u0301", 0, ""***REMOVED***,
	***REMOVED***"aaaa\u0300\u0301", 3, ""***REMOVED***,
	***REMOVED***"\u0300a\u0300\u0301", 2, ""***REMOVED***,
	***REMOVED***"\u00C2", 0, ""***REMOVED***,
	***REMOVED***"a\u00C2", 1, ""***REMOVED***,
	// decomposition may recombine
	***REMOVED***"\u0226", 0, ""***REMOVED***,
	// no boundary
	***REMOVED***"", -1, ""***REMOVED***,
	***REMOVED***"\u0300\u0301", -1, ""***REMOVED***,
	***REMOVED***"\u0300", -1, ""***REMOVED***,
	***REMOVED***"\x80\x80", -1, ""***REMOVED***,
	***REMOVED***"\x80\x80\u0301", -1, ""***REMOVED***,
	// Hangul
	***REMOVED***"다음", 3, ""***REMOVED***,
	***REMOVED***"다", 0, ""***REMOVED***,
	***REMOVED***"\u1103\u1161\u110B\u1173\u11B7", 6, ""***REMOVED***,
	***REMOVED***"\u110B\u1173\u11B7\u1103\u1161", 9, ""***REMOVED***,
	// too many combining characters.
	***REMOVED***grave(maxNonStarters - 1), -1, ""***REMOVED***,
	// May still be preceded with a non-starter.
	***REMOVED***grave(maxNonStarters), -1, ""***REMOVED***,
	// May still need to insert a cgj after the last combiner.
	***REMOVED***grave(maxNonStarters + 1), 2, ""***REMOVED***,
	***REMOVED***grave(maxNonStarters + 2), 4, ""***REMOVED***,

	***REMOVED***"a" + grave(maxNonStarters-1), 0, ""***REMOVED***,
	***REMOVED***"a" + grave(maxNonStarters), 0, ""***REMOVED***,
	// May still need to insert a cgj after the last combiner.
	***REMOVED***"a" + grave(maxNonStarters+1), 3, ""***REMOVED***,
	***REMOVED***"a" + grave(maxNonStarters+2), 5, ""***REMOVED***,
***REMOVED***

func lastBoundaryF(rb *reorderBuffer, s string) (int, []byte) ***REMOVED***
	return rb.f.form.LastBoundary([]byte(s)), nil
***REMOVED***

func TestLastBoundary(t *testing.T) ***REMOVED***
	runPosTests(t, "TestLastBoundary", NFC, lastBoundaryF, lastBoundaryTests)
***REMOVED***

type spanTest struct ***REMOVED***
	input string
	atEOF bool
	n     int
	err   error
***REMOVED***

var quickSpanTests = []spanTest***REMOVED***
	***REMOVED***"", true, 0, nil***REMOVED***,
	// starters
	***REMOVED***"a", true, 1, nil***REMOVED***,
	***REMOVED***"abc", true, 3, nil***REMOVED***,
	***REMOVED***"\u043Eb", true, 3, nil***REMOVED***,
	// incomplete last rune.
	***REMOVED***"\xCC", true, 1, nil***REMOVED***,
	***REMOVED***"\xCC", false, 0, transform.ErrShortSrc***REMOVED***,
	***REMOVED***"a\xCC", true, 2, nil***REMOVED***,
	***REMOVED***"a\xCC", false, 0, transform.ErrShortSrc***REMOVED***, // TODO: could be 1 for NFD
	// incorrectly ordered combining characters
	***REMOVED***"\u0300\u0316", true, 0, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"\u0300\u0316", false, 0, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"\u0300\u0316cd", true, 0, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"\u0300\u0316cd", false, 0, transform.ErrEndOfSpan***REMOVED***,
	// have a maximum number of combining characters.
	***REMOVED***rep(0x035D, 30) + "\u035B", true, 0, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"a" + rep(0x035D, 30) + "\u035B", true, 0, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"Ɵ" + rep(0x035D, 30) + "\u035B", true, 0, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"aa" + rep(0x035D, 30) + "\u035B", true, 1, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***rep(0x035D, 30) + cgj + "\u035B", true, 64, nil***REMOVED***,
	***REMOVED***"a" + rep(0x035D, 30) + cgj + "\u035B", true, 65, nil***REMOVED***,
	***REMOVED***"Ɵ" + rep(0x035D, 30) + cgj + "\u035B", true, 66, nil***REMOVED***,
	***REMOVED***"aa" + rep(0x035D, 30) + cgj + "\u035B", true, 66, nil***REMOVED***,

	***REMOVED***"a" + rep(0x035D, 30) + cgj + "\u035B", false, 61, transform.ErrShortSrc***REMOVED***,
	***REMOVED***"Ɵ" + rep(0x035D, 30) + cgj + "\u035B", false, 62, transform.ErrShortSrc***REMOVED***,
	***REMOVED***"aa" + rep(0x035D, 30) + cgj + "\u035B", false, 62, transform.ErrShortSrc***REMOVED***,
***REMOVED***

var quickSpanNFDTests = []spanTest***REMOVED***
	// needs decomposing
	***REMOVED***"\u00C0", true, 0, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"abc\u00C0", true, 3, transform.ErrEndOfSpan***REMOVED***,
	// correctly ordered combining characters
	***REMOVED***"\u0300", true, 2, nil***REMOVED***,
	***REMOVED***"ab\u0300", true, 4, nil***REMOVED***,
	***REMOVED***"ab\u0300cd", true, 6, nil***REMOVED***,
	***REMOVED***"\u0300cd", true, 4, nil***REMOVED***,
	***REMOVED***"\u0316\u0300", true, 4, nil***REMOVED***,
	***REMOVED***"ab\u0316\u0300", true, 6, nil***REMOVED***,
	***REMOVED***"ab\u0316\u0300cd", true, 8, nil***REMOVED***,
	***REMOVED***"ab\u0316\u0300\u00C0", true, 6, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"\u0316\u0300cd", true, 6, nil***REMOVED***,
	***REMOVED***"\u043E\u0308b", true, 5, nil***REMOVED***,
	// incorrectly ordered combining characters
	***REMOVED***"ab\u0300\u0316", true, 1, transform.ErrEndOfSpan***REMOVED***, // TODO: we could skip 'b' as well.
	***REMOVED***"ab\u0300\u0316cd", true, 1, transform.ErrEndOfSpan***REMOVED***,
	// Hangul
	***REMOVED***"같은", true, 0, transform.ErrEndOfSpan***REMOVED***,
***REMOVED***

var quickSpanNFCTests = []spanTest***REMOVED***
	// okay composed
	***REMOVED***"\u00C0", true, 2, nil***REMOVED***,
	***REMOVED***"abc\u00C0", true, 5, nil***REMOVED***,
	// correctly ordered combining characters
	// TODO: b may combine with modifiers, which is why this fails. We could
	// make a more precise test that that actually checks whether last
	// characters combines. Probably not worth it.
	***REMOVED***"ab\u0300", true, 1, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"ab\u0300cd", true, 1, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"ab\u0316\u0300", true, 1, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"ab\u0316\u0300cd", true, 1, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"\u00C0\u035D", true, 4, nil***REMOVED***,
	// we do not special case leading combining characters
	***REMOVED***"\u0300cd", true, 0, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"\u0300", true, 0, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"\u0316\u0300", true, 0, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"\u0316\u0300cd", true, 0, transform.ErrEndOfSpan***REMOVED***,
	// incorrectly ordered combining characters
	***REMOVED***"ab\u0300\u0316", true, 1, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***"ab\u0300\u0316cd", true, 1, transform.ErrEndOfSpan***REMOVED***,
	// Hangul
	***REMOVED***"같은", true, 6, nil***REMOVED***,
	***REMOVED***"같은", false, 3, transform.ErrShortSrc***REMOVED***,
	// We return the start of the violating segment in case of overflow.
	***REMOVED***grave(30) + "\uff9e", true, 0, transform.ErrEndOfSpan***REMOVED***,
	***REMOVED***grave(30), true, 0, transform.ErrEndOfSpan***REMOVED***,
***REMOVED***

func runSpanTests(t *testing.T, name string, f Form, testCases []spanTest) ***REMOVED***
	for i, tc := range testCases ***REMOVED***
		s := fmt.Sprintf("Bytes/%s/%d=%+q/atEOF=%v", name, i, pc(tc.input), tc.atEOF)
		ok := testtext.Run(t, s, func(t *testing.T) ***REMOVED***
			n, err := f.Span([]byte(tc.input), tc.atEOF)
			if n != tc.n || err != tc.err ***REMOVED***
				t.Errorf("\n got %d, %v;\nwant %d, %v", n, err, tc.n, tc.err)
			***REMOVED***
		***REMOVED***)
		if !ok ***REMOVED***
			continue // Don't do the String variant if the Bytes variant failed.
		***REMOVED***
		s = fmt.Sprintf("String/%s/%d=%+q/atEOF=%v", name, i, pc(tc.input), tc.atEOF)
		testtext.Run(t, s, func(t *testing.T) ***REMOVED***
			n, err := f.SpanString(tc.input, tc.atEOF)
			if n != tc.n || err != tc.err ***REMOVED***
				t.Errorf("\n got %d, %v;\nwant %d, %v", n, err, tc.n, tc.err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSpan(t *testing.T) ***REMOVED***
	runSpanTests(t, "NFD", NFD, quickSpanTests)
	runSpanTests(t, "NFD", NFD, quickSpanNFDTests)
	runSpanTests(t, "NFC", NFC, quickSpanTests)
	runSpanTests(t, "NFC", NFC, quickSpanNFCTests)
***REMOVED***

var isNormalTests = []PositionTest***REMOVED***
	***REMOVED***"", 1, ""***REMOVED***,
	// illegal runes
	***REMOVED***"\xff", 1, ""***REMOVED***,
	// starters
	***REMOVED***"a", 1, ""***REMOVED***,
	***REMOVED***"abc", 1, ""***REMOVED***,
	***REMOVED***"\u043Eb", 1, ""***REMOVED***,
	// incorrectly ordered combining characters
	***REMOVED***"\u0300\u0316", 0, ""***REMOVED***,
	***REMOVED***"ab\u0300\u0316", 0, ""***REMOVED***,
	***REMOVED***"ab\u0300\u0316cd", 0, ""***REMOVED***,
	***REMOVED***"\u0300\u0316cd", 0, ""***REMOVED***,
***REMOVED***
var isNormalNFDTests = []PositionTest***REMOVED***
	// needs decomposing
	***REMOVED***"\u00C0", 0, ""***REMOVED***,
	***REMOVED***"abc\u00C0", 0, ""***REMOVED***,
	// correctly ordered combining characters
	***REMOVED***"\u0300", 1, ""***REMOVED***,
	***REMOVED***"ab\u0300", 1, ""***REMOVED***,
	***REMOVED***"ab\u0300cd", 1, ""***REMOVED***,
	***REMOVED***"\u0300cd", 1, ""***REMOVED***,
	***REMOVED***"\u0316\u0300", 1, ""***REMOVED***,
	***REMOVED***"ab\u0316\u0300", 1, ""***REMOVED***,
	***REMOVED***"ab\u0316\u0300cd", 1, ""***REMOVED***,
	***REMOVED***"\u0316\u0300cd", 1, ""***REMOVED***,
	***REMOVED***"\u043E\u0308b", 1, ""***REMOVED***,
	// Hangul
	***REMOVED***"같은", 0, ""***REMOVED***,
***REMOVED***
var isNormalNFCTests = []PositionTest***REMOVED***
	// okay composed
	***REMOVED***"\u00C0", 1, ""***REMOVED***,
	***REMOVED***"abc\u00C0", 1, ""***REMOVED***,
	// need reordering
	***REMOVED***"a\u0300", 0, ""***REMOVED***,
	***REMOVED***"a\u0300cd", 0, ""***REMOVED***,
	***REMOVED***"a\u0316\u0300", 0, ""***REMOVED***,
	***REMOVED***"a\u0316\u0300cd", 0, ""***REMOVED***,
	// correctly ordered combining characters
	***REMOVED***"ab\u0300", 1, ""***REMOVED***,
	***REMOVED***"ab\u0300cd", 1, ""***REMOVED***,
	***REMOVED***"ab\u0316\u0300", 1, ""***REMOVED***,
	***REMOVED***"ab\u0316\u0300cd", 1, ""***REMOVED***,
	***REMOVED***"\u00C0\u035D", 1, ""***REMOVED***,
	***REMOVED***"\u0300", 1, ""***REMOVED***,
	***REMOVED***"\u0316\u0300cd", 1, ""***REMOVED***,
	// Hangul
	***REMOVED***"같은", 1, ""***REMOVED***,
***REMOVED***

var isNormalNFKXTests = []PositionTest***REMOVED***
	// Special case.
	***REMOVED***"\u00BC", 0, ""***REMOVED***,
***REMOVED***

func isNormalF(rb *reorderBuffer, s string) (int, []byte) ***REMOVED***
	if rb.f.form.IsNormal([]byte(s)) ***REMOVED***
		return 1, nil
	***REMOVED***
	return 0, nil
***REMOVED***

func isNormalStringF(rb *reorderBuffer, s string) (int, []byte) ***REMOVED***
	if rb.f.form.IsNormalString(s) ***REMOVED***
		return 1, nil
	***REMOVED***
	return 0, nil
***REMOVED***

func TestIsNormal(t *testing.T) ***REMOVED***
	runPosTests(t, "TestIsNormalNFD1", NFD, isNormalF, isNormalTests)
	runPosTests(t, "TestIsNormalNFD2", NFD, isNormalF, isNormalNFDTests)
	runPosTests(t, "TestIsNormalNFC1", NFC, isNormalF, isNormalTests)
	runPosTests(t, "TestIsNormalNFC2", NFC, isNormalF, isNormalNFCTests)
	runPosTests(t, "TestIsNormalNFKD1", NFKD, isNormalF, isNormalTests)
	runPosTests(t, "TestIsNormalNFKD2", NFKD, isNormalF, isNormalNFDTests)
	runPosTests(t, "TestIsNormalNFKD3", NFKD, isNormalF, isNormalNFKXTests)
	runPosTests(t, "TestIsNormalNFKC1", NFKC, isNormalF, isNormalTests)
	runPosTests(t, "TestIsNormalNFKC2", NFKC, isNormalF, isNormalNFCTests)
	runPosTests(t, "TestIsNormalNFKC3", NFKC, isNormalF, isNormalNFKXTests)
***REMOVED***

func TestIsNormalString(t *testing.T) ***REMOVED***
	runPosTests(t, "TestIsNormalNFD1", NFD, isNormalStringF, isNormalTests)
	runPosTests(t, "TestIsNormalNFD2", NFD, isNormalStringF, isNormalNFDTests)
	runPosTests(t, "TestIsNormalNFC1", NFC, isNormalStringF, isNormalTests)
	runPosTests(t, "TestIsNormalNFC2", NFC, isNormalStringF, isNormalNFCTests)
***REMOVED***

type AppendTest struct ***REMOVED***
	left  string
	right string
	out   string
***REMOVED***

type appendFunc func(f Form, out []byte, s string) []byte

var fstr = []string***REMOVED***"NFC", "NFD", "NFKC", "NFKD"***REMOVED***

func runNormTests(t *testing.T, name string, fn appendFunc) ***REMOVED***
	for f := NFC; f <= NFKD; f++ ***REMOVED***
		runAppendTests(t, name, f, fn, normTests[f])
	***REMOVED***
***REMOVED***

func runAppendTests(t *testing.T, name string, f Form, fn appendFunc, tests []AppendTest) ***REMOVED***
	for i, test := range tests ***REMOVED***
		t.Run(fmt.Sprintf("%s/%d", fstr[f], i), func(t *testing.T) ***REMOVED***
			id := pc(test.left + test.right)
			if *testn >= 0 && i != *testn ***REMOVED***
				return
			***REMOVED***
			t.Run("fn", func(t *testing.T) ***REMOVED***
				out := []byte(test.left)
				have := string(fn(f, out, test.right))
				if len(have) != len(test.out) ***REMOVED***
					t.Errorf("%+q: length is %d; want %d (%+q vs %+q)", id, len(have), len(test.out), pc(have), pc(test.out))
				***REMOVED***
				if have != test.out ***REMOVED***
					k, pf := pidx(have, test.out)
					t.Errorf("%+q:\nwas  %s%+q; \nwant %s%+q", id, pf, pc(have[k:]), pf, pc(test.out[k:]))
				***REMOVED***
			***REMOVED***)

			// Bootstrap by normalizing input. Ensures that the various variants
			// behave the same.
			for g := NFC; g <= NFKD; g++ ***REMOVED***
				if f == g ***REMOVED***
					continue
				***REMOVED***
				t.Run(fstr[g], func(t *testing.T) ***REMOVED***
					want := g.String(test.left + test.right)
					have := string(fn(g, g.AppendString(nil, test.left), test.right))
					if len(have) != len(want) ***REMOVED***
						t.Errorf("%+q: length is %d; want %d (%+q vs %+q)", id, len(have), len(want), pc(have), pc(want))
					***REMOVED***
					if have != want ***REMOVED***
						k, pf := pidx(have, want)
						t.Errorf("%+q:\nwas  %s%+q; \nwant %s%+q", id, pf, pc(have[k:]), pf, pc(want[k:]))
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

var normTests = [][]AppendTest***REMOVED***
	appendTestsNFC,
	appendTestsNFD,
	appendTestsNFKC,
	appendTestsNFKD,
***REMOVED***

var appendTestsNFC = []AppendTest***REMOVED***
	***REMOVED***"", ascii, ascii***REMOVED***,
	***REMOVED***"", txt_all, txt_all***REMOVED***,
	***REMOVED***"\uff9e", grave(30), "\uff9e" + grave(29) + cgj + grave(1)***REMOVED***,
	***REMOVED***grave(30), "\uff9e", grave(30) + cgj + "\uff9e"***REMOVED***,

	// Tests designed for Iter.
	***REMOVED*** // ordering of non-composing combining characters
		"",
		"\u0305\u0316",
		"\u0316\u0305",
	***REMOVED***,
	***REMOVED*** // segment overflow
		"",
		"a" + rep(0x0305, maxNonStarters+4) + "\u0316",
		"a" + rep(0x0305, maxNonStarters) + cgj + "\u0316" + rep(0x305, 4),
	***REMOVED***,

	***REMOVED*** // Combine across non-blocking non-starters.
		// U+0327 COMBINING CEDILLA;Mn;202;NSM;;;;;N;NON-SPACING CEDILLA;;;;
		// U+0325 COMBINING RING BELOW;Mn;220;NSM;;;;;N;NON-SPACING RING BELOW;;;;
		"", "a\u0327\u0325", "\u1e01\u0327",
	***REMOVED***,

	***REMOVED*** // Jamo V+T does not combine.
		"",
		"\u1161\u11a8",
		"\u1161\u11a8",
	***REMOVED***,

	// Stability tests: see http://www.unicode.org/review/pr-29.html.
	***REMOVED***"", "\u0b47\u0300\u0b3e", "\u0b47\u0300\u0b3e"***REMOVED***,
	***REMOVED***"", "\u1100\u0300\u1161", "\u1100\u0300\u1161"***REMOVED***,
	***REMOVED***"", "\u0b47\u0b3e", "\u0b4b"***REMOVED***,
	***REMOVED***"", "\u1100\u1161", "\uac00"***REMOVED***,

	// U+04DA MALAYALAM VOWEL SIGN O;Mc;0;L;0D46 0D3E;;;;N;;;;;
	***REMOVED*** // 0d4a starts a new segment.
		"",
		"\u0d4a" + strings.Repeat("\u0d3e", 15) + "\u0d4a" + strings.Repeat("\u0d3e", 15),
		"\u0d4a" + strings.Repeat("\u0d3e", 15) + "\u0d4a" + strings.Repeat("\u0d3e", 15),
	***REMOVED***,

	***REMOVED*** // Split combining characters.
		// TODO: don't insert CGJ before starters.
		"",
		"\u0d46" + strings.Repeat("\u0d3e", 31),
		"\u0d4a" + strings.Repeat("\u0d3e", 29) + cgj + "\u0d3e",
	***REMOVED***,

	***REMOVED*** // Split combining characters.
		"",
		"\u0d4a" + strings.Repeat("\u0d3e", 30),
		"\u0d4a" + strings.Repeat("\u0d3e", 29) + cgj + "\u0d3e",
	***REMOVED***,

	***REMOVED*** //  https://golang.org/issues/20079
		"",
		"\xeb\u0344",
		"\xeb\u0308\u0301",
	***REMOVED***,

	***REMOVED*** //  https://golang.org/issues/20079
		"",
		"\uac00" + strings.Repeat("\u0300", 30),
		"\uac00" + strings.Repeat("\u0300", 29) + "\u034f\u0300",
	***REMOVED***,

	***REMOVED*** //  https://golang.org/issues/20079
		"",
		"\xeb" + strings.Repeat("\u0300", 31),
		"\xeb" + strings.Repeat("\u0300", 30) + "\u034f\u0300",
	***REMOVED***,
***REMOVED***

var appendTestsNFD = []AppendTest***REMOVED***
	// TODO: Move some of the tests here.
***REMOVED***

var appendTestsNFKC = []AppendTest***REMOVED***
	// empty buffers
	***REMOVED***"", "", ""***REMOVED***,
	***REMOVED***"a", "", "a"***REMOVED***,
	***REMOVED***"", "a", "a"***REMOVED***,
	***REMOVED***"", "\u0041\u0307\u0304", "\u01E0"***REMOVED***,
	// segment split across buffers
	***REMOVED***"", "a\u0300b", "\u00E0b"***REMOVED***,
	***REMOVED***"a", "\u0300b", "\u00E0b"***REMOVED***,
	***REMOVED***"a", "\u0300\u0316", "\u00E0\u0316"***REMOVED***,
	***REMOVED***"a", "\u0316\u0300", "\u00E0\u0316"***REMOVED***,
	***REMOVED***"a", "\u0300a\u0300", "\u00E0\u00E0"***REMOVED***,
	***REMOVED***"a", "\u0300a\u0300a\u0300", "\u00E0\u00E0\u00E0"***REMOVED***,
	***REMOVED***"a", "\u0300aaa\u0300aaa\u0300", "\u00E0aa\u00E0aa\u00E0"***REMOVED***,
	***REMOVED***"a\u0300", "\u0327", "\u00E0\u0327"***REMOVED***,
	***REMOVED***"a\u0327", "\u0300", "\u00E0\u0327"***REMOVED***,
	***REMOVED***"a\u0316", "\u0300", "\u00E0\u0316"***REMOVED***,
	***REMOVED***"\u0041\u0307", "\u0304", "\u01E0"***REMOVED***,
	// Hangul
	***REMOVED***"", "\u110B\u1173", "\uC73C"***REMOVED***,
	***REMOVED***"", "\u1103\u1161", "\uB2E4"***REMOVED***,
	***REMOVED***"", "\u110B\u1173\u11B7", "\uC74C"***REMOVED***,
	***REMOVED***"", "\u320E", "\x28\uAC00\x29"***REMOVED***,
	***REMOVED***"", "\x28\u1100\u1161\x29", "\x28\uAC00\x29"***REMOVED***,
	***REMOVED***"\u1103", "\u1161", "\uB2E4"***REMOVED***,
	***REMOVED***"\u110B", "\u1173\u11B7", "\uC74C"***REMOVED***,
	***REMOVED***"\u110B\u1173", "\u11B7", "\uC74C"***REMOVED***,
	***REMOVED***"\uC73C", "\u11B7", "\uC74C"***REMOVED***,
	// UTF-8 encoding split across buffers
	***REMOVED***"a\xCC", "\x80", "\u00E0"***REMOVED***,
	***REMOVED***"a\xCC", "\x80b", "\u00E0b"***REMOVED***,
	***REMOVED***"a\xCC", "\x80a\u0300", "\u00E0\u00E0"***REMOVED***,
	***REMOVED***"a\xCC", "\x80\x80", "\u00E0\x80"***REMOVED***,
	***REMOVED***"a\xCC", "\x80\xCC", "\u00E0\xCC"***REMOVED***,
	***REMOVED***"a\u0316\xCC", "\x80a\u0316\u0300", "\u00E0\u0316\u00E0\u0316"***REMOVED***,
	// ending in incomplete UTF-8 encoding
	***REMOVED***"", "\xCC", "\xCC"***REMOVED***,
	***REMOVED***"a", "\xCC", "a\xCC"***REMOVED***,
	***REMOVED***"a", "b\xCC", "ab\xCC"***REMOVED***,
	***REMOVED***"\u0226", "\xCC", "\u0226\xCC"***REMOVED***,
	// illegal runes
	***REMOVED***"", "\x80", "\x80"***REMOVED***,
	***REMOVED***"", "\x80\x80\x80", "\x80\x80\x80"***REMOVED***,
	***REMOVED***"", "\xCC\x80\x80\x80", "\xCC\x80\x80\x80"***REMOVED***,
	***REMOVED***"", "a\x80", "a\x80"***REMOVED***,
	***REMOVED***"", "a\x80\x80\x80", "a\x80\x80\x80"***REMOVED***,
	***REMOVED***"", "a\x80\x80\x80\x80\x80\x80", "a\x80\x80\x80\x80\x80\x80"***REMOVED***,
	***REMOVED***"a", "\x80\x80\x80", "a\x80\x80\x80"***REMOVED***,
	// overflow
	***REMOVED***"", strings.Repeat("\x80", 33), strings.Repeat("\x80", 33)***REMOVED***,
	***REMOVED***strings.Repeat("\x80", 33), "", strings.Repeat("\x80", 33)***REMOVED***,
	***REMOVED***strings.Repeat("\x80", 33), strings.Repeat("\x80", 33), strings.Repeat("\x80", 66)***REMOVED***,
	// overflow of combining characters
	***REMOVED***"", grave(34), grave(30) + cgj + grave(4)***REMOVED***,
	***REMOVED***"", grave(36), grave(30) + cgj + grave(6)***REMOVED***,
	***REMOVED***grave(29), grave(5), grave(30) + cgj + grave(4)***REMOVED***,
	***REMOVED***grave(30), grave(4), grave(30) + cgj + grave(4)***REMOVED***,
	***REMOVED***grave(30), grave(3), grave(30) + cgj + grave(3)***REMOVED***,
	***REMOVED***grave(30) + "\xCC", "\x80", grave(30) + cgj + grave(1)***REMOVED***,
	***REMOVED***"", "\uFDFA" + grave(14), "\u0635\u0644\u0649 \u0627\u0644\u0644\u0647 \u0639\u0644\u064a\u0647 \u0648\u0633\u0644\u0645" + grave(14)***REMOVED***,
	***REMOVED***"", "\uFDFA" + grave(28) + "\u0316", "\u0635\u0644\u0649 \u0627\u0644\u0644\u0647 \u0639\u0644\u064a\u0647 \u0648\u0633\u0644\u0645\u0316" + grave(28)***REMOVED***,
	// - First rune has a trailing non-starter.
	***REMOVED***"\u00d5", grave(30), "\u00d5" + grave(29) + cgj + grave(1)***REMOVED***,
	// - U+FF9E decomposes into a non-starter in compatibility mode. A CGJ must be
	//   inserted even when FF9E starts a new segment.
	***REMOVED***"\uff9e", grave(30), "\u3099" + grave(29) + cgj + grave(1)***REMOVED***,
	***REMOVED***grave(30), "\uff9e", grave(30) + cgj + "\u3099"***REMOVED***,
	// - Many non-starter decompositions in a row causing overflow.
	***REMOVED***"", rep(0x340, 31), rep(0x300, 30) + cgj + "\u0300"***REMOVED***,
	***REMOVED***"", rep(0xFF9E, 31), rep(0x3099, 30) + cgj + "\u3099"***REMOVED***,

	***REMOVED***"", "\u0644\u0625" + rep(0x300, 31), "\u0644\u0625" + rep(0x300, 29) + cgj + "\u0300\u0300"***REMOVED***,
	***REMOVED***"", "\ufef9" + rep(0x300, 31), "\u0644\u0625" + rep(0x300, 29) + cgj + rep(0x0300, 2)***REMOVED***,
	***REMOVED***"", "\ufef9" + rep(0x300, 31), "\u0644\u0625" + rep(0x300, 29) + cgj + rep(0x0300, 2)***REMOVED***,

	// U+0F81 TIBETAN VOWEL SIGN REVERSED II splits into two modifiers.
	***REMOVED***"", "\u0f7f" + rep(0xf71, 29) + "\u0f81", "\u0f7f" + rep(0xf71, 29) + cgj + "\u0f71\u0f80"***REMOVED***,
	***REMOVED***"", "\u0f7f" + rep(0xf71, 28) + "\u0f81", "\u0f7f" + rep(0xf71, 29) + "\u0f80"***REMOVED***,
	***REMOVED***"", "\u0f7f" + rep(0xf81, 16), "\u0f7f" + rep(0xf71, 15) + rep(0xf80, 15) + cgj + "\u0f71\u0f80"***REMOVED***,

	// weird UTF-8
	***REMOVED***"\u00E0\xE1", "\x86", "\u00E0\xE1\x86"***REMOVED***,
	***REMOVED***"a\u0300\u11B7", "\u0300", "\u00E0\u11B7\u0300"***REMOVED***,
	***REMOVED***"a\u0300\u11B7\u0300", "\u0300", "\u00E0\u11B7\u0300\u0300"***REMOVED***,
	***REMOVED***"\u0300", "\xF8\x80\x80\x80\x80\u0300", "\u0300\xF8\x80\x80\x80\x80\u0300"***REMOVED***,
	***REMOVED***"\u0300", "\xFC\x80\x80\x80\x80\x80\u0300", "\u0300\xFC\x80\x80\x80\x80\x80\u0300"***REMOVED***,
	***REMOVED***"\xF8\x80\x80\x80\x80\u0300", "\u0300", "\xF8\x80\x80\x80\x80\u0300\u0300"***REMOVED***,
	***REMOVED***"\xFC\x80\x80\x80\x80\x80\u0300", "\u0300", "\xFC\x80\x80\x80\x80\x80\u0300\u0300"***REMOVED***,
	***REMOVED***"\xF8\x80\x80\x80", "\x80\u0300\u0300", "\xF8\x80\x80\x80\x80\u0300\u0300"***REMOVED***,

	***REMOVED***"", strings.Repeat("a\u0316\u0300", 6), strings.Repeat("\u00E0\u0316", 6)***REMOVED***,
	// large input.
	***REMOVED***"", strings.Repeat("a\u0300\u0316", 31), strings.Repeat("\u00E0\u0316", 31)***REMOVED***,
	***REMOVED***"", strings.Repeat("a\u0300\u0316", 4000), strings.Repeat("\u00E0\u0316", 4000)***REMOVED***,
	***REMOVED***"", strings.Repeat("\x80\x80", 4000), strings.Repeat("\x80\x80", 4000)***REMOVED***,
	***REMOVED***"", "\u0041\u0307\u0304", "\u01E0"***REMOVED***,
***REMOVED***

var appendTestsNFKD = []AppendTest***REMOVED***
	***REMOVED***"", "a" + grave(64), "a" + grave(30) + cgj + grave(30) + cgj + grave(4)***REMOVED***,

	***REMOVED*** // segment overflow on unchanged character
		"",
		"a" + grave(64) + "\u0316",
		"a" + grave(30) + cgj + grave(30) + cgj + "\u0316" + grave(4),
	***REMOVED***,
	***REMOVED*** // segment overflow on unchanged character + start value
		"",
		"a" + grave(98) + "\u0316",
		"a" + grave(30) + cgj + grave(30) + cgj + grave(30) + cgj + "\u0316" + grave(8),
	***REMOVED***,
	***REMOVED*** // segment overflow on decomposition. (U+0340 decomposes to U+0300.)
		"",
		"a" + grave(59) + "\u0340",
		"a" + grave(30) + cgj + grave(30),
	***REMOVED***,
	***REMOVED*** // segment overflow on non-starter decomposition
		"",
		"a" + grave(33) + "\u0340" + grave(30) + "\u0320",
		"a" + grave(30) + cgj + grave(30) + cgj + "\u0320" + grave(4),
	***REMOVED***,
	***REMOVED*** // start value after ASCII overflow
		"",
		rep('a', segSize) + grave(32) + "\u0320",
		rep('a', segSize) + grave(30) + cgj + "\u0320" + grave(2),
	***REMOVED***,
	***REMOVED*** // Jamo overflow
		"",
		"\u1100\u1161" + grave(30) + "\u0320" + grave(2),
		"\u1100\u1161" + grave(29) + cgj + "\u0320" + grave(3),
	***REMOVED***,
	***REMOVED*** // Hangul
		"",
		"\uac00",
		"\u1100\u1161",
	***REMOVED***,
	***REMOVED*** // Hangul overflow
		"",
		"\uac00" + grave(32) + "\u0320",
		"\u1100\u1161" + grave(29) + cgj + "\u0320" + grave(3),
	***REMOVED***,
	***REMOVED*** // Hangul overflow in Hangul mode.
		"",
		"\uac00\uac00" + grave(32) + "\u0320",
		"\u1100\u1161\u1100\u1161" + grave(29) + cgj + "\u0320" + grave(3),
	***REMOVED***,
	***REMOVED*** // Hangul overflow in Hangul mode.
		"",
		strings.Repeat("\uac00", 3) + grave(32) + "\u0320",
		strings.Repeat("\u1100\u1161", 3) + grave(29) + cgj + "\u0320" + grave(3),
	***REMOVED***,
	***REMOVED*** // start value after cc=0
		"",
		"您您" + grave(34) + "\u0320",
		"您您" + grave(30) + cgj + "\u0320" + grave(4),
	***REMOVED***,
	***REMOVED*** // start value after normalization
		"",
		"\u0300\u0320a" + grave(34) + "\u0320",
		"\u0320\u0300a" + grave(30) + cgj + "\u0320" + grave(4),
	***REMOVED***,
	***REMOVED***
		// U+0F81 TIBETAN VOWEL SIGN REVERSED II splits into two modifiers.
		"",
		"a\u0f7f" + rep(0xf71, 29) + "\u0f81",
		"a\u0f7f" + rep(0xf71, 29) + cgj + "\u0f71\u0f80",
	***REMOVED***,
***REMOVED***

func TestAppend(t *testing.T) ***REMOVED***
	runNormTests(t, "Append", func(f Form, out []byte, s string) []byte ***REMOVED***
		return f.Append(out, []byte(s)...)
	***REMOVED***)
***REMOVED***

func TestAppendString(t *testing.T) ***REMOVED***
	runNormTests(t, "AppendString", func(f Form, out []byte, s string) []byte ***REMOVED***
		return f.AppendString(out, s)
	***REMOVED***)
***REMOVED***

func TestBytes(t *testing.T) ***REMOVED***
	runNormTests(t, "Bytes", func(f Form, out []byte, s string) []byte ***REMOVED***
		buf := []byte***REMOVED******REMOVED***
		buf = append(buf, out...)
		buf = append(buf, s...)
		return f.Bytes(buf)
	***REMOVED***)
***REMOVED***

func TestString(t *testing.T) ***REMOVED***
	runNormTests(t, "String", func(f Form, out []byte, s string) []byte ***REMOVED***
		outs := string(out) + s
		return []byte(f.String(outs))
	***REMOVED***)
***REMOVED***

func TestLinking(t *testing.T) ***REMOVED***
	const prog = `
	package main
	import "fmt"
	import "golang.org/x/text/unicode/norm"
	func main() ***REMOVED*** fmt.Println(norm.%s) ***REMOVED***
	`
	baseline, errB := testtext.CodeSize(fmt.Sprintf(prog, "MaxSegmentSize"))
	withTables, errT := testtext.CodeSize(fmt.Sprintf(prog, `NFC.String("")`))
	if errB != nil || errT != nil ***REMOVED***
		t.Skipf("code size failed: %v and %v", errB, errT)
	***REMOVED***
	// Tables are at least 50K
	if d := withTables - baseline; d < 50*1024 ***REMOVED***
		t.Errorf("tables appear not to be dropped: %d - %d = %d",
			withTables, baseline, d)
	***REMOVED***
***REMOVED***

func appendBench(f Form, in []byte) func() ***REMOVED***
	buf := make([]byte, 0, 4*len(in))
	return func() ***REMOVED***
		f.Append(buf, in...)
	***REMOVED***
***REMOVED***

func bytesBench(f Form, in []byte) func() ***REMOVED***
	return func() ***REMOVED***
		f.Bytes(in)
	***REMOVED***
***REMOVED***

func iterBench(f Form, in []byte) func() ***REMOVED***
	iter := Iter***REMOVED******REMOVED***
	return func() ***REMOVED***
		iter.Init(f, in)
		for !iter.Done() ***REMOVED***
			iter.Next()
		***REMOVED***
	***REMOVED***
***REMOVED***

func transformBench(f Form, in []byte) func() ***REMOVED***
	buf := make([]byte, 4*len(in))
	return func() ***REMOVED***
		if _, n, err := f.Transform(buf, in, true); err != nil || len(in) != n ***REMOVED***
			log.Panic(n, len(in), err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func readerBench(f Form, in []byte) func() ***REMOVED***
	buf := make([]byte, 4*len(in))
	return func() ***REMOVED***
		r := f.Reader(bytes.NewReader(in))
		var err error
		for err == nil ***REMOVED***
			_, err = r.Read(buf)
		***REMOVED***
		if err != io.EOF ***REMOVED***
			panic("")
		***REMOVED***
	***REMOVED***
***REMOVED***

func writerBench(f Form, in []byte) func() ***REMOVED***
	buf := make([]byte, 0, 4*len(in))
	return func() ***REMOVED***
		r := f.Writer(bytes.NewBuffer(buf))
		if _, err := r.Write(in); err != nil ***REMOVED***
			panic("")
		***REMOVED***
	***REMOVED***
***REMOVED***

func appendBenchmarks(bm []func(), f Form, in []byte) []func() ***REMOVED***
	bm = append(bm, appendBench(f, in))
	bm = append(bm, iterBench(f, in))
	bm = append(bm, transformBench(f, in))
	bm = append(bm, readerBench(f, in))
	bm = append(bm, writerBench(f, in))
	return bm
***REMOVED***

func doFormBenchmark(b *testing.B, inf, f Form, s string) ***REMOVED***
	b.StopTimer()
	in := inf.Bytes([]byte(s))
	bm := appendBenchmarks(nil, f, in)
	b.SetBytes(int64(len(in) * len(bm)))
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		for _, fn := range bm ***REMOVED***
			fn()
		***REMOVED***
	***REMOVED***
***REMOVED***

func doSingle(b *testing.B, f func(Form, []byte) func(), s []byte) ***REMOVED***
	b.StopTimer()
	fn := f(NFC, s)
	b.SetBytes(int64(len(s)))
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		fn()
	***REMOVED***
***REMOVED***

var (
	smallNoChange = []byte("nörmalization")
	smallChange   = []byte("No\u0308rmalization")
	ascii         = strings.Repeat("There is nothing to change here! ", 500)
)

func lowerBench(f Form, in []byte) func() ***REMOVED***
	// Use package strings instead of bytes as it doesn't allocate memory
	// if there aren't any changes.
	s := string(in)
	return func() ***REMOVED***
		strings.ToLower(s)
	***REMOVED***
***REMOVED***

func BenchmarkLowerCaseNoChange(b *testing.B) ***REMOVED***
	doSingle(b, lowerBench, smallNoChange)
***REMOVED***
func BenchmarkLowerCaseChange(b *testing.B) ***REMOVED***
	doSingle(b, lowerBench, smallChange)
***REMOVED***

func quickSpanBench(f Form, in []byte) func() ***REMOVED***
	return func() ***REMOVED***
		f.QuickSpan(in)
	***REMOVED***
***REMOVED***

func BenchmarkQuickSpanChangeNFC(b *testing.B) ***REMOVED***
	doSingle(b, quickSpanBench, smallNoChange)
***REMOVED***

func BenchmarkBytesNoChangeNFC(b *testing.B) ***REMOVED***
	doSingle(b, bytesBench, smallNoChange)
***REMOVED***
func BenchmarkBytesChangeNFC(b *testing.B) ***REMOVED***
	doSingle(b, bytesBench, smallChange)
***REMOVED***

func BenchmarkAppendNoChangeNFC(b *testing.B) ***REMOVED***
	doSingle(b, appendBench, smallNoChange)
***REMOVED***
func BenchmarkAppendChangeNFC(b *testing.B) ***REMOVED***
	doSingle(b, appendBench, smallChange)
***REMOVED***
func BenchmarkAppendLargeNFC(b *testing.B) ***REMOVED***
	doSingle(b, appendBench, txt_all_bytes)
***REMOVED***

func BenchmarkIterNoChangeNFC(b *testing.B) ***REMOVED***
	doSingle(b, iterBench, smallNoChange)
***REMOVED***
func BenchmarkIterChangeNFC(b *testing.B) ***REMOVED***
	doSingle(b, iterBench, smallChange)
***REMOVED***
func BenchmarkIterLargeNFC(b *testing.B) ***REMOVED***
	doSingle(b, iterBench, txt_all_bytes)
***REMOVED***

func BenchmarkTransformNoChangeNFC(b *testing.B) ***REMOVED***
	doSingle(b, transformBench, smallNoChange)
***REMOVED***
func BenchmarkTransformChangeNFC(b *testing.B) ***REMOVED***
	doSingle(b, transformBench, smallChange)
***REMOVED***
func BenchmarkTransformLargeNFC(b *testing.B) ***REMOVED***
	doSingle(b, transformBench, txt_all_bytes)
***REMOVED***

func BenchmarkNormalizeAsciiNFC(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFC, NFC, ascii)
***REMOVED***
func BenchmarkNormalizeAsciiNFD(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFC, NFD, ascii)
***REMOVED***
func BenchmarkNormalizeAsciiNFKC(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFC, NFKC, ascii)
***REMOVED***
func BenchmarkNormalizeAsciiNFKD(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFC, NFKD, ascii)
***REMOVED***

func BenchmarkNormalizeNFC2NFC(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFC, NFC, txt_all)
***REMOVED***
func BenchmarkNormalizeNFC2NFD(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFC, NFD, txt_all)
***REMOVED***
func BenchmarkNormalizeNFD2NFC(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFD, NFC, txt_all)
***REMOVED***
func BenchmarkNormalizeNFD2NFD(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFD, NFD, txt_all)
***REMOVED***

// Hangul is often special-cased, so we test it separately.
func BenchmarkNormalizeHangulNFC2NFC(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFC, NFC, txt_kr)
***REMOVED***
func BenchmarkNormalizeHangulNFC2NFD(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFC, NFD, txt_kr)
***REMOVED***
func BenchmarkNormalizeHangulNFD2NFC(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFD, NFC, txt_kr)
***REMOVED***
func BenchmarkNormalizeHangulNFD2NFD(b *testing.B) ***REMOVED***
	doFormBenchmark(b, NFD, NFD, txt_kr)
***REMOVED***

var forms = []Form***REMOVED***NFC, NFD, NFKC, NFKD***REMOVED***

func doTextBenchmark(b *testing.B, s string) ***REMOVED***
	b.StopTimer()
	in := []byte(s)
	bm := []func()***REMOVED******REMOVED***
	for _, f := range forms ***REMOVED***
		bm = appendBenchmarks(bm, f, in)
	***REMOVED***
	b.SetBytes(int64(len(s) * len(bm)))
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		for _, f := range bm ***REMOVED***
			f()
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkCanonicalOrdering(b *testing.B) ***REMOVED***
	doTextBenchmark(b, txt_canon)
***REMOVED***
func BenchmarkExtendedLatin(b *testing.B) ***REMOVED***
	doTextBenchmark(b, txt_vn)
***REMOVED***
func BenchmarkMiscTwoByteUtf8(b *testing.B) ***REMOVED***
	doTextBenchmark(b, twoByteUtf8)
***REMOVED***
func BenchmarkMiscThreeByteUtf8(b *testing.B) ***REMOVED***
	doTextBenchmark(b, threeByteUtf8)
***REMOVED***
func BenchmarkHangul(b *testing.B) ***REMOVED***
	doTextBenchmark(b, txt_kr)
***REMOVED***
func BenchmarkJapanese(b *testing.B) ***REMOVED***
	doTextBenchmark(b, txt_jp)
***REMOVED***
func BenchmarkChinese(b *testing.B) ***REMOVED***
	doTextBenchmark(b, txt_cn)
***REMOVED***
func BenchmarkOverflow(b *testing.B) ***REMOVED***
	doTextBenchmark(b, overflow)
***REMOVED***

var overflow = string(bytes.Repeat([]byte("\u035D"), 4096)) + "\u035B"

// Tests sampled from the Canonical ordering tests (Part 2) of
// http://unicode.org/Public/UNIDATA/NormalizationTest.txt
const txt_canon = `\u0061\u0315\u0300\u05AE\u0300\u0062 \u0061\u0300\u0315\u0300\u05AE\u0062
\u0061\u0302\u0315\u0300\u05AE\u0062 \u0061\u0307\u0315\u0300\u05AE\u0062
\u0061\u0315\u0300\u05AE\u030A\u0062 \u0061\u059A\u0316\u302A\u031C\u0062
\u0061\u032E\u059A\u0316\u302A\u0062 \u0061\u0338\u093C\u0334\u0062 
\u0061\u059A\u0316\u302A\u0339       \u0061\u0341\u0315\u0300\u05AE\u0062
\u0061\u0348\u059A\u0316\u302A\u0062 \u0061\u0361\u0345\u035D\u035C\u0062
\u0061\u0366\u0315\u0300\u05AE\u0062 \u0061\u0315\u0300\u05AE\u0486\u0062
\u0061\u05A4\u059A\u0316\u302A\u0062 \u0061\u0315\u0300\u05AE\u0613\u0062
\u0061\u0315\u0300\u05AE\u0615\u0062 \u0061\u0617\u0315\u0300\u05AE\u0062
\u0061\u0619\u0618\u064D\u064E\u0062 \u0061\u0315\u0300\u05AE\u0654\u0062
\u0061\u0315\u0300\u05AE\u06DC\u0062 \u0061\u0733\u0315\u0300\u05AE\u0062
\u0061\u0744\u059A\u0316\u302A\u0062 \u0061\u0315\u0300\u05AE\u0745\u0062
\u0061\u09CD\u05B0\u094D\u3099\u0062 \u0061\u0E38\u0E48\u0E38\u0C56\u0062
\u0061\u0EB8\u0E48\u0E38\u0E49\u0062 \u0061\u0F72\u0F71\u0EC8\u0F71\u0062
\u0061\u1039\u05B0\u094D\u3099\u0062 \u0061\u05B0\u094D\u3099\u1A60\u0062
\u0061\u3099\u093C\u0334\u1BE6\u0062 \u0061\u3099\u093C\u0334\u1C37\u0062
\u0061\u1CD9\u059A\u0316\u302A\u0062 \u0061\u2DED\u0315\u0300\u05AE\u0062
\u0061\u2DEF\u0315\u0300\u05AE\u0062 \u0061\u302D\u302E\u059A\u0316\u0062`

// Taken from http://creativecommons.org/licenses/by-sa/3.0/vn/
const txt_vn = `Với các điều kiện sau: Ghi nhận công của tác giả. 
Nếu bạn sử dụng, chuyển đổi, hoặc xây dựng dự án từ 
nội dung được chia sẻ này, bạn phải áp dụng giấy phép này hoặc 
một giấy phép khác có các điều khoản tương tự như giấy phép này
cho dự án của bạn. Hiểu rằng: Miễn — Bất kỳ các điều kiện nào
trên đây cũng có thể được miễn bỏ nếu bạn được sự cho phép của
người sở hữu bản quyền. Phạm vi công chúng — Khi tác phẩm hoặc
bất kỳ chương nào của tác phẩm đã trong vùng dành cho công
chúng theo quy định của pháp luật thì tình trạng của nó không 
bị ảnh hưởng bởi giấy phép trong bất kỳ trường hợp nào.`

// Taken from http://creativecommons.org/licenses/by-sa/1.0/deed.ru
const txt_ru = `При обязательном соблюдении следующих условий:
Attribution — Вы должны атрибутировать произведение (указывать
автора и источник) в порядке, предусмотренном автором или
лицензиаром (но только так, чтобы никоим образом не подразумевалось,
что они поддерживают вас или использование вами данного произведения).
Υπό τις ακόλουθες προϋποθέσεις:`

// Taken from http://creativecommons.org/licenses/by-sa/3.0/gr/
const txt_gr = `Αναφορά Δημιουργού — Θα πρέπει να κάνετε την αναφορά στο έργο με τον
τρόπο που έχει οριστεί από το δημιουργό ή το χορηγούντο την άδεια
(χωρίς όμως να εννοείται με οποιονδήποτε τρόπο ότι εγκρίνουν εσάς ή
τη χρήση του έργου από εσάς). Παρόμοια Διανομή — Εάν αλλοιώσετε,
τροποποιήσετε ή δημιουργήσετε περαιτέρω βασισμένοι στο έργο θα
μπορείτε να διανέμετε το έργο που θα προκύψει μόνο με την ίδια ή
παρόμοια άδεια.`

// Taken from http://creativecommons.org/licenses/by-sa/3.0/deed.ar
const txt_ar = `بموجب الشروط التالية نسب المصنف — يجب عليك أن
تنسب العمل بالطريقة التي تحددها المؤلف أو المرخص (ولكن ليس بأي حال من
الأحوال أن توحي وتقترح بتحول أو استخدامك للعمل).
المشاركة على قدم المساواة — إذا كنت يعدل ، والتغيير ، أو الاستفادة
من هذا العمل ، قد ينتج عن توزيع العمل إلا في ظل تشابه او تطابق فى واحد
لهذا الترخيص.`

// Taken from http://creativecommons.org/licenses/by-sa/1.0/il/
const txt_il = `בכפוף לתנאים הבאים: ייחוס — עליך לייחס את היצירה (לתת קרדיט) באופן
המצויין על-ידי היוצר או מעניק הרישיון (אך לא בשום אופן המרמז על כך
שהם תומכים בך או בשימוש שלך ביצירה). שיתוף זהה — אם תחליט/י לשנות,
לעבד או ליצור יצירה נגזרת בהסתמך על יצירה זו, תוכל/י להפיץ את יצירתך
החדשה רק תחת אותו הרישיון או רישיון דומה לרישיון זה.`

const twoByteUtf8 = txt_ru + txt_gr + txt_ar + txt_il

// Taken from http://creativecommons.org/licenses/by-sa/2.0/kr/
const txt_kr = `다음과 같은 조건을 따라야 합니다: 저작자표시
(Attribution) — 저작자나 이용허락자가 정한 방법으로 저작물의
원저작자를 표시하여야 합니다(그러나 원저작자가 이용자나 이용자의
이용을 보증하거나 추천한다는 의미로 표시해서는 안됩니다). 
동일조건변경허락 — 이 저작물을 이용하여 만든 이차적 저작물에는 본
라이선스와 동일한 라이선스를 적용해야 합니다.`

// Taken from http://creativecommons.org/licenses/by-sa/3.0/th/
const txt_th = `ภายใต้เงื่อนไข ดังต่อไปนี้ : แสดงที่มา — คุณต้องแสดงที่
มาของงานดังกล่าว ตามรูปแบบที่ผู้สร้างสรรค์หรือผู้อนุญาตกำหนด (แต่
ไม่ใช่ในลักษณะที่ว่า พวกเขาสนับสนุนคุณหรือสนับสนุนการที่
คุณนำงานไปใช้) อนุญาตแบบเดียวกัน — หากคุณดัดแปลง เปลี่ยนรูป หรื
อต่อเติมงานนี้ คุณต้องใช้สัญญาอนุญาตแบบเดียวกันหรือแบบที่เหมื
อนกับสัญญาอนุญาตที่ใช้กับงานนี้เท่านั้น`

const threeByteUtf8 = txt_th

// Taken from http://creativecommons.org/licenses/by-sa/2.0/jp/
const txt_jp = `あなたの従うべき条件は以下の通りです。
表示 — あなたは原著作者のクレジットを表示しなければなりません。
継承 — もしあなたがこの作品を改変、変形または加工した場合、
あなたはその結果生じた作品をこの作品と同一の許諾条件の下でのみ
頒布することができます。`

// http://creativecommons.org/licenses/by-sa/2.5/cn/
const txt_cn = `您可以自由： 复制、发行、展览、表演、放映、
广播或通过信息网络传播本作品 创作演绎作品
对本作品进行商业性使用 惟须遵守下列条件：
署名 — 您必须按照作者或者许可人指定的方式对作品进行署名。
相同方式共享 — 如果您改变、转换本作品或者以本作品为基础进行创作，
您只能采用与本协议相同的许可协议发布基于本作品的演绎作品。`

const txt_cjk = txt_cn + txt_jp + txt_kr
const txt_all = txt_vn + twoByteUtf8 + threeByteUtf8 + txt_cjk

var txt_all_bytes = []byte(txt_all)
