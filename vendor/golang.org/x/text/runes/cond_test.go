// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runes

import (
	"strings"
	"testing"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/transform"
)

var (
	toUpper = cases.Upper(language.Und)
	toLower = cases.Lower(language.Und)
)

type spanformer interface ***REMOVED***
	transform.SpanningTransformer
***REMOVED***

func TestPredicate(t *testing.T) ***REMOVED***
	testConditional(t, func(rt *unicode.RangeTable, t, f spanformer) spanformer ***REMOVED***
		return If(Predicate(func(r rune) bool ***REMOVED***
			return unicode.Is(rt, r)
		***REMOVED***), t, f)
	***REMOVED***)
***REMOVED***

func TestIn(t *testing.T) ***REMOVED***
	testConditional(t, func(rt *unicode.RangeTable, t, f spanformer) spanformer ***REMOVED***
		return If(In(rt), t, f)
	***REMOVED***)
***REMOVED***

func TestNotIn(t *testing.T) ***REMOVED***
	testConditional(t, func(rt *unicode.RangeTable, t, f spanformer) spanformer ***REMOVED***
		return If(NotIn(rt), f, t)
	***REMOVED***)
***REMOVED***

func testConditional(t *testing.T, f func(rt *unicode.RangeTable, t, f spanformer) spanformer) ***REMOVED***
	lower := f(unicode.Latin, toLower, toLower)

	for i, tt := range []transformTest***REMOVED******REMOVED***
		desc:    "empty",
		szDst:   large,
		atEOF:   true,
		in:      "",
		out:     "",
		outFull: "",
		t:       lower,
	***REMOVED***, ***REMOVED***
		desc:    "small",
		szDst:   1,
		atEOF:   true,
		in:      "B",
		out:     "b",
		outFull: "b",
		errSpan: transform.ErrEndOfSpan,
		t:       lower,
	***REMOVED***, ***REMOVED***
		desc:    "short dst",
		szDst:   2,
		atEOF:   true,
		in:      "AAA",
		out:     "aa",
		outFull: "aaa",
		err:     transform.ErrShortDst,
		errSpan: transform.ErrEndOfSpan,
		t:       lower,
	***REMOVED***, ***REMOVED***
		desc:    "short dst writing error",
		szDst:   1,
		atEOF:   false,
		in:      "A\x80",
		out:     "a",
		outFull: "a\x80",
		err:     transform.ErrShortDst,
		errSpan: transform.ErrEndOfSpan,
		t:       lower,
	***REMOVED***, ***REMOVED***
		desc:    "short dst writing incomplete rune",
		szDst:   2,
		atEOF:   true,
		in:      "Σ\xc2",
		out:     "Σ",
		outFull: "Σ\xc2",
		err:     transform.ErrShortDst,
		t:       f(unicode.Latin, toLower, nil),
	***REMOVED***, ***REMOVED***
		desc:    "short dst, longer",
		szDst:   5,
		atEOF:   true,
		in:      "Hellø",
		out:     "Hell",
		outFull: "Hellø",
		err:     transform.ErrShortDst,
		// idem is used to test short buffers by forcing processing of full-rune increments.
		t: f(unicode.Latin, Map(idem), nil),
	***REMOVED***, ***REMOVED***
		desc:    "short dst, longer, writing error",
		szDst:   6,
		atEOF:   false,
		in:      "\x80Hello\x80",
		out:     "\x80Hello",
		outFull: "\x80Hello\x80",
		err:     transform.ErrShortDst,
		t:       f(unicode.Latin, Map(idem), nil),
	***REMOVED***, ***REMOVED***
		desc:    "short src",
		szDst:   2,
		atEOF:   false,
		in:      "A\xc2",
		out:     "a",
		outFull: "a\xc2",
		err:     transform.ErrShortSrc,
		errSpan: transform.ErrEndOfSpan,
		t:       lower,
	***REMOVED***, ***REMOVED***
		desc:    "short src no change",
		szDst:   2,
		atEOF:   false,
		in:      "a\xc2",
		out:     "a",
		outFull: "a\xc2",
		err:     transform.ErrShortSrc,
		errSpan: transform.ErrShortSrc,
		nSpan:   1,
		t:       lower,
	***REMOVED***, ***REMOVED***
		desc:    "invalid input, atEOF",
		szDst:   large,
		atEOF:   true,
		in:      "\x80",
		out:     "\x80",
		outFull: "\x80",
		t:       lower,
	***REMOVED***, ***REMOVED***
		desc:    "invalid input, !atEOF",
		szDst:   large,
		atEOF:   false,
		in:      "\x80",
		out:     "\x80",
		outFull: "\x80",
		t:       lower,
	***REMOVED***, ***REMOVED***
		desc:    "invalid input, incomplete rune atEOF",
		szDst:   large,
		atEOF:   true,
		in:      "\xc2",
		out:     "\xc2",
		outFull: "\xc2",
		t:       lower,
	***REMOVED***, ***REMOVED***
		desc:    "nop",
		szDst:   large,
		atEOF:   true,
		in:      "Hello World!",
		out:     "Hello World!",
		outFull: "Hello World!",
		t:       f(unicode.Latin, nil, nil),
	***REMOVED***, ***REMOVED***
		desc:    "nop in",
		szDst:   large,
		atEOF:   true,
		in:      "THIS IS α ΤΕΣΤ",
		out:     "this is α ΤΕΣΤ",
		outFull: "this is α ΤΕΣΤ",
		errSpan: transform.ErrEndOfSpan,
		t:       f(unicode.Greek, nil, toLower),
	***REMOVED***, ***REMOVED***
		desc:    "nop in latin",
		szDst:   large,
		atEOF:   true,
		in:      "THIS IS α ΤΕΣΤ",
		out:     "THIS IS α τεστ",
		outFull: "THIS IS α τεστ",
		errSpan: transform.ErrEndOfSpan,
		t:       f(unicode.Latin, nil, toLower),
	***REMOVED***, ***REMOVED***
		desc:    "nop not in",
		szDst:   large,
		atEOF:   true,
		in:      "THIS IS α ΤΕΣΤ",
		out:     "this is α ΤΕΣΤ",
		outFull: "this is α ΤΕΣΤ",
		errSpan: transform.ErrEndOfSpan,
		t:       f(unicode.Latin, toLower, nil),
	***REMOVED***, ***REMOVED***
		desc:    "pass atEOF is true when at end",
		szDst:   large,
		atEOF:   true,
		in:      "hello",
		out:     "HELLO",
		outFull: "HELLO",
		errSpan: transform.ErrEndOfSpan,
		t:       f(unicode.Latin, upperAtEOF***REMOVED******REMOVED***, nil),
	***REMOVED***, ***REMOVED***
		desc:    "pass atEOF is true when at end of segment",
		szDst:   large,
		atEOF:   true,
		in:      "hello ",
		out:     "HELLO ",
		outFull: "HELLO ",
		errSpan: transform.ErrEndOfSpan,
		t:       f(unicode.Latin, upperAtEOF***REMOVED******REMOVED***, nil),
	***REMOVED***, ***REMOVED***
		desc:    "don't pass atEOF is true when atEOF is false",
		szDst:   large,
		atEOF:   false,
		in:      "hello",
		out:     "",
		outFull: "HELLO",
		err:     transform.ErrShortSrc,
		errSpan: transform.ErrShortSrc,
		t:       f(unicode.Latin, upperAtEOF***REMOVED******REMOVED***, nil),
	***REMOVED***, ***REMOVED***
		desc:    "pass atEOF is true when at end, no change",
		szDst:   large,
		atEOF:   true,
		in:      "HELLO",
		out:     "HELLO",
		outFull: "HELLO",
		t:       f(unicode.Latin, upperAtEOF***REMOVED******REMOVED***, nil),
	***REMOVED***, ***REMOVED***
		desc:    "pass atEOF is true when at end of segment, no change",
		szDst:   large,
		atEOF:   true,
		in:      "HELLO ",
		out:     "HELLO ",
		outFull: "HELLO ",
		t:       f(unicode.Latin, upperAtEOF***REMOVED******REMOVED***, nil),
	***REMOVED***, ***REMOVED***
		desc:    "large input ASCII",
		szDst:   12000,
		atEOF:   false,
		in:      strings.Repeat("HELLO", 2000),
		out:     strings.Repeat("hello", 2000),
		outFull: strings.Repeat("hello", 2000),
		errSpan: transform.ErrEndOfSpan,
		err:     nil,
		t:       lower,
	***REMOVED***, ***REMOVED***
		desc:    "large input non-ASCII",
		szDst:   12000,
		atEOF:   false,
		in:      strings.Repeat("\u3333", 2000),
		out:     strings.Repeat("\u3333", 2000),
		outFull: strings.Repeat("\u3333", 2000),
		err:     nil,
		t:       lower,
	***REMOVED******REMOVED*** ***REMOVED***
		tt.check(t, i)
	***REMOVED***
***REMOVED***

// upperAtEOF is a strange Transformer that converts text to uppercase, but only
// if atEOF is true.
type upperAtEOF struct***REMOVED*** transform.NopResetter ***REMOVED***

func (upperAtEOF) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	if !atEOF ***REMOVED***
		return 0, 0, transform.ErrShortSrc
	***REMOVED***
	return toUpper.Transform(dst, src, atEOF)
***REMOVED***

func (upperAtEOF) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	if !atEOF ***REMOVED***
		return 0, transform.ErrShortSrc
	***REMOVED***
	return toUpper.Span(src, atEOF)
***REMOVED***

func BenchmarkConditional(b *testing.B) ***REMOVED***
	doBench(b, If(In(unicode.Hangul), transform.Nop, transform.Nop))
***REMOVED***
