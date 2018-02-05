// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bidirule

import (
	"fmt"
	"testing"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/unicode/bidi"
)

const (
	strL   = "ABC"    // Left to right - most letters in LTR scripts
	strR   = "עברית"  // Right to left - most letters in non-Arabic RTL scripts
	strAL  = "دبي"    // Arabic letters - most letters in the Arabic script
	strEN  = "123"    // European Number (0-9, and Extended Arabic-Indic numbers)
	strES  = "+-"     // European Number Separator (+ and -)
	strET  = "$"      // European Number Terminator (currency symbols, the hash sign, the percent sign and so on)
	strAN  = "\u0660" // Arabic Number; this encompasses the Arabic-Indic numbers, but not the Extended Arabic-Indic numbers
	strCS  = ","      // Common Number Separator (. , / : et al)
	strNSM = "\u0300" // Nonspacing Mark - most combining accents
	strBN  = "\u200d" // Boundary Neutral - control characters (ZWNJ, ZWJ, and others)
	strB   = "\u2029" // Paragraph Separator
	strS   = "\u0009" // Segment Separator
	strWS  = " "      // Whitespace, including the SPACE character
	strON  = "@"      // Other Neutrals, including @, &, parentheses, MIDDLE DOT
)

type ruleTest struct ***REMOVED***
	in  string
	dir bidi.Direction
	n   int // position at which the rule fails
	err error

	// For tests that split the string in two.
	pSrc  int   // number of source bytes to consume first
	szDst int   // size of destination buffer
	nSrc  int   // source bytes consumed and bytes written
	err0  error // error after first run
***REMOVED***

func init() ***REMOVED***
	for rule, cases := range testCases ***REMOVED***
		for i, tc := range cases ***REMOVED***
			if tc.err == nil ***REMOVED***
				testCases[rule][i].n = len(tc.in)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func doTests(t *testing.T, fn func(t *testing.T, tc ruleTest)) ***REMOVED***
	for rule, cases := range testCases ***REMOVED***
		for i, tc := range cases ***REMOVED***
			name := fmt.Sprintf("%d/%d:%+q:%s", rule, i, tc.in, tc.in)
			testtext.Run(t, name, func(t *testing.T) ***REMOVED***
				fn(t, tc)
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDirection(t *testing.T) ***REMOVED***
	doTests(t, func(t *testing.T, tc ruleTest) ***REMOVED***
		dir := Direction([]byte(tc.in))
		if dir != tc.dir ***REMOVED***
			t.Errorf("dir was %v; want %v", dir, tc.dir)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestDirectionString(t *testing.T) ***REMOVED***
	doTests(t, func(t *testing.T, tc ruleTest) ***REMOVED***
		dir := DirectionString(tc.in)
		if dir != tc.dir ***REMOVED***
			t.Errorf("dir was %v; want %v", dir, tc.dir)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestValid(t *testing.T) ***REMOVED***
	doTests(t, func(t *testing.T, tc ruleTest) ***REMOVED***
		got := Valid([]byte(tc.in))
		want := tc.err == nil
		if got != want ***REMOVED***
			t.Fatalf("Valid: got %v; want %v", got, want)
		***REMOVED***

		got = ValidString(tc.in)
		want = tc.err == nil
		if got != want ***REMOVED***
			t.Fatalf("Valid: got %v; want %v", got, want)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestSpan(t *testing.T) ***REMOVED***
	doTests(t, func(t *testing.T, tc ruleTest) ***REMOVED***
		// Skip tests that test for limited destination buffer size.
		if tc.szDst > 0 ***REMOVED***
			return
		***REMOVED***

		r := New()
		src := []byte(tc.in)

		n, err := r.Span(src[:tc.pSrc], tc.pSrc == len(tc.in))
		if err != tc.err0 ***REMOVED***
			t.Errorf("err0 was %v; want %v", err, tc.err0)
		***REMOVED***
		if n != tc.nSrc ***REMOVED***
			t.Fatalf("nSrc was %d; want %d", n, tc.nSrc)
		***REMOVED***

		n, err = r.Span(src[n:], true)
		if err != tc.err ***REMOVED***
			t.Errorf("error was %v; want %v", err, tc.err)
		***REMOVED***
		if got := n + tc.nSrc; got != tc.n ***REMOVED***
			t.Errorf("n was %d; want %d", got, tc.n)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestTransform(t *testing.T) ***REMOVED***
	doTests(t, func(t *testing.T, tc ruleTest) ***REMOVED***
		r := New()

		src := []byte(tc.in)
		dst := make([]byte, len(tc.in))
		if tc.szDst > 0 ***REMOVED***
			dst = make([]byte, tc.szDst)
		***REMOVED***

		// First transform operates on a zero-length string for most tests.
		nDst, nSrc, err := r.Transform(dst, src[:tc.pSrc], tc.pSrc == len(tc.in))
		if err != tc.err0 ***REMOVED***
			t.Errorf("err0 was %v; want %v", err, tc.err0)
		***REMOVED***
		if nDst != nSrc ***REMOVED***
			t.Fatalf("nDst (%d) and nSrc (%d) should match", nDst, nSrc)
		***REMOVED***
		if nSrc != tc.nSrc ***REMOVED***
			t.Fatalf("nSrc was %d; want %d", nSrc, tc.nSrc)
		***REMOVED***

		dst1 := make([]byte, len(tc.in))
		copy(dst1, dst[:nDst])

		nDst, nSrc, err = r.Transform(dst1[nDst:], src[nSrc:], true)
		if err != tc.err ***REMOVED***
			t.Errorf("error was %v; want %v", err, tc.err)
		***REMOVED***
		if nDst != nSrc ***REMOVED***
			t.Fatalf("nDst (%d) and nSrc (%d) should match", nDst, nSrc)
		***REMOVED***
		n := nSrc + tc.nSrc
		if n != tc.n ***REMOVED***
			t.Fatalf("n was %d; want %d", n, tc.n)
		***REMOVED***
		if got, want := string(dst1[:n]), tc.in[:tc.n]; got != want ***REMOVED***
			t.Errorf("got %+q; want %+q", got, want)
		***REMOVED***
	***REMOVED***)
***REMOVED***
