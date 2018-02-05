// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runes

import (
	"strings"
	"testing"
	"unicode/utf8"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/transform"
)

type transformTest struct ***REMOVED***
	desc    string
	szDst   int
	atEOF   bool
	repl    string
	in      string
	out     string // result string of first call to Transform
	outFull string // transform of entire input string
	err     error
	errSpan error
	nSpan   int

	t transform.SpanningTransformer
***REMOVED***

const large = 10240

func (tt *transformTest) check(t *testing.T, i int) ***REMOVED***
	if tt.t == nil ***REMOVED***
		return
	***REMOVED***
	dst := make([]byte, tt.szDst)
	src := []byte(tt.in)
	nDst, nSrc, err := tt.t.Transform(dst, src, tt.atEOF)
	if err != tt.err ***REMOVED***
		t.Errorf("%d:%s:error: got %v; want %v", i, tt.desc, err, tt.err)
	***REMOVED***
	if got := string(dst[:nDst]); got != tt.out ***REMOVED***
		t.Errorf("%d:%s:out: got %q; want %q", i, tt.desc, got, tt.out)
	***REMOVED***

	// Calls tt.t.Transform for the remainder of the input. We use this to test
	// the nSrc return value.
	out := make([]byte, large)
	n := copy(out, dst[:nDst])
	nDst, _, _ = tt.t.Transform(out[n:], src[nSrc:], true)
	if got, want := string(out[:n+nDst]), tt.outFull; got != want ***REMOVED***
		t.Errorf("%d:%s:outFull: got %q; want %q", i, tt.desc, got, want)
	***REMOVED***

	tt.t.Reset()
	p := 0
	for ; p < len(tt.in) && p < len(tt.outFull) && tt.in[p] == tt.outFull[p]; p++ ***REMOVED***
	***REMOVED***
	if tt.nSpan != 0 ***REMOVED***
		p = tt.nSpan
	***REMOVED***
	if n, err = tt.t.Span([]byte(tt.in), tt.atEOF); n != p || err != tt.errSpan ***REMOVED***
		t.Errorf("%d:%s:span: got %d, %v; want %d, %v", i, tt.desc, n, err, p, tt.errSpan)
	***REMOVED***
***REMOVED***

func idem(r rune) rune ***REMOVED*** return r ***REMOVED***

func TestMap(t *testing.T) ***REMOVED***
	runes := []rune***REMOVED***'a', 'ç', '中', '\U00012345', 'a'***REMOVED***
	// Default mapper used for this test.
	rotate := Map(func(r rune) rune ***REMOVED***
		for i, m := range runes ***REMOVED***
			if m == r ***REMOVED***
				return runes[i+1]
			***REMOVED***
		***REMOVED***
		return r
	***REMOVED***)

	for i, tt := range []transformTest***REMOVED******REMOVED***
		desc:    "empty",
		szDst:   large,
		atEOF:   true,
		in:      "",
		out:     "",
		outFull: "",
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "no change",
		szDst:   1,
		atEOF:   true,
		in:      "b",
		out:     "b",
		outFull: "b",
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "short dst",
		szDst:   2,
		atEOF:   true,
		in:      "aaaa",
		out:     "ç",
		outFull: "çççç",
		err:     transform.ErrShortDst,
		errSpan: transform.ErrEndOfSpan,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "short dst ascii, no change",
		szDst:   2,
		atEOF:   true,
		in:      "bbb",
		out:     "bb",
		outFull: "bbb",
		err:     transform.ErrShortDst,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "short dst writing error",
		szDst:   2,
		atEOF:   false,
		in:      "a\x80",
		out:     "ç",
		outFull: "ç\ufffd",
		err:     transform.ErrShortDst,
		errSpan: transform.ErrEndOfSpan,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "short dst writing incomplete rune",
		szDst:   2,
		atEOF:   true,
		in:      "a\xc0",
		out:     "ç",
		outFull: "ç\ufffd",
		err:     transform.ErrShortDst,
		errSpan: transform.ErrEndOfSpan,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "short dst, longer",
		szDst:   5,
		atEOF:   true,
		in:      "Hellø",
		out:     "Hell",
		outFull: "Hellø",
		err:     transform.ErrShortDst,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "short dst, single",
		szDst:   1,
		atEOF:   false,
		in:      "ø",
		out:     "",
		outFull: "ø",
		err:     transform.ErrShortDst,
		t:       Map(idem),
	***REMOVED***, ***REMOVED***
		desc:    "short dst, longer, writing error",
		szDst:   8,
		atEOF:   false,
		in:      "\x80Hello\x80",
		out:     "\ufffdHello",
		outFull: "\ufffdHello\ufffd",
		err:     transform.ErrShortDst,
		errSpan: transform.ErrEndOfSpan,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "short src",
		szDst:   2,
		atEOF:   false,
		in:      "a\xc2",
		out:     "ç",
		outFull: "ç\ufffd",
		err:     transform.ErrShortSrc,
		errSpan: transform.ErrEndOfSpan,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "invalid input, atEOF",
		szDst:   large,
		atEOF:   true,
		in:      "\x80",
		out:     "\ufffd",
		outFull: "\ufffd",
		errSpan: transform.ErrEndOfSpan,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "invalid input, !atEOF",
		szDst:   large,
		atEOF:   false,
		in:      "\x80",
		out:     "\ufffd",
		outFull: "\ufffd",
		errSpan: transform.ErrEndOfSpan,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "incomplete rune !atEOF",
		szDst:   large,
		atEOF:   false,
		in:      "\xc2",
		out:     "",
		outFull: "\ufffd",
		err:     transform.ErrShortSrc,
		errSpan: transform.ErrShortSrc,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "invalid input, incomplete rune atEOF",
		szDst:   large,
		atEOF:   true,
		in:      "\xc2",
		out:     "\ufffd",
		outFull: "\ufffd",
		errSpan: transform.ErrEndOfSpan,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "misc correct",
		szDst:   large,
		atEOF:   true,
		in:      "a\U00012345 ç!",
		out:     "ça 中!",
		outFull: "ça 中!",
		errSpan: transform.ErrEndOfSpan,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "misc correct and invalid",
		szDst:   large,
		atEOF:   true,
		in:      "Hello\x80 w\x80orl\xc0d!\xc0",
		out:     "Hello\ufffd w\ufffdorl\ufffdd!\ufffd",
		outFull: "Hello\ufffd w\ufffdorl\ufffdd!\ufffd",
		errSpan: transform.ErrEndOfSpan,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "misc correct and invalid, short src",
		szDst:   large,
		atEOF:   false,
		in:      "Hello\x80 w\x80orl\xc0d!\xc2",
		out:     "Hello\ufffd w\ufffdorl\ufffdd!",
		outFull: "Hello\ufffd w\ufffdorl\ufffdd!\ufffd",
		err:     transform.ErrShortSrc,
		errSpan: transform.ErrEndOfSpan,
		t:       rotate,
	***REMOVED***, ***REMOVED***
		desc:    "misc correct and invalid, short src, replacing RuneError",
		szDst:   large,
		atEOF:   false,
		in:      "Hel\ufffdlo\x80 w\x80orl\xc0d!\xc2",
		out:     "Hel?lo? w?orl?d!",
		outFull: "Hel?lo? w?orl?d!?",
		errSpan: transform.ErrEndOfSpan,
		err:     transform.ErrShortSrc,
		t: Map(func(r rune) rune ***REMOVED***
			if r == utf8.RuneError ***REMOVED***
				return '?'
			***REMOVED***
			return r
		***REMOVED***),
	***REMOVED******REMOVED*** ***REMOVED***
		tt.check(t, i)
	***REMOVED***
***REMOVED***

func TestRemove(t *testing.T) ***REMOVED***
	remove := Remove(Predicate(func(r rune) bool ***REMOVED***
		return strings.ContainsRune("aeiou\u0300\uFF24\U00012345", r)
	***REMOVED***))

	for i, tt := range []transformTest***REMOVED***
		0: ***REMOVED***
			szDst:   large,
			atEOF:   true,
			in:      "",
			out:     "",
			outFull: "",
			t:       remove,
		***REMOVED***,
		1: ***REMOVED***
			szDst:   0,
			atEOF:   true,
			in:      "aaaa",
			out:     "",
			outFull: "",
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		2: ***REMOVED***
			szDst:   1,
			atEOF:   true,
			in:      "aaaa",
			out:     "",
			outFull: "",
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		3: ***REMOVED***
			szDst:   1,
			atEOF:   true,
			in:      "baaaa",
			out:     "b",
			outFull: "b",
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		4: ***REMOVED***
			szDst:   2,
			atEOF:   true,
			in:      "açaaa",
			out:     "ç",
			outFull: "ç",
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		5: ***REMOVED***
			szDst:   2,
			atEOF:   true,
			in:      "aaaç",
			out:     "ç",
			outFull: "ç",
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		6: ***REMOVED***
			szDst:   2,
			atEOF:   false,
			in:      "a\x80",
			out:     "",
			outFull: "\ufffd",
			err:     transform.ErrShortDst,
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		7: ***REMOVED***
			szDst:   1,
			atEOF:   true,
			in:      "a\xc0",
			out:     "",
			outFull: "\ufffd",
			err:     transform.ErrShortDst,
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		8: ***REMOVED***
			szDst:   1,
			atEOF:   false,
			in:      "a\xc2",
			out:     "",
			outFull: "\ufffd",
			err:     transform.ErrShortSrc,
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		9: ***REMOVED***
			szDst:   large,
			atEOF:   true,
			in:      "\x80",
			out:     "\ufffd",
			outFull: "\ufffd",
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		10: ***REMOVED***
			szDst:   large,
			atEOF:   false,
			in:      "\x80",
			out:     "\ufffd",
			outFull: "\ufffd",
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		11: ***REMOVED***
			szDst:   large,
			atEOF:   true,
			in:      "\xc2",
			out:     "\ufffd",
			outFull: "\ufffd",
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		12: ***REMOVED***
			szDst:   large,
			atEOF:   false,
			in:      "\xc2",
			out:     "",
			outFull: "\ufffd",
			err:     transform.ErrShortSrc,
			errSpan: transform.ErrShortSrc,
			t:       remove,
		***REMOVED***,
		13: ***REMOVED***
			szDst:   large,
			atEOF:   true,
			in:      "Hello \U00012345world!",
			out:     "Hll wrld!",
			outFull: "Hll wrld!",
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		14: ***REMOVED***
			szDst:   large,
			atEOF:   true,
			in:      "Hello\x80 w\x80orl\xc0d!\xc0",
			out:     "Hll\ufffd w\ufffdrl\ufffdd!\ufffd",
			outFull: "Hll\ufffd w\ufffdrl\ufffdd!\ufffd",
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		15: ***REMOVED***
			szDst:   large,
			atEOF:   false,
			in:      "Hello\x80 w\x80orl\xc0d!\xc2",
			out:     "Hll\ufffd w\ufffdrl\ufffdd!",
			outFull: "Hll\ufffd w\ufffdrl\ufffdd!\ufffd",
			err:     transform.ErrShortSrc,
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		16: ***REMOVED***
			szDst:   large,
			atEOF:   false,
			in:      "Hel\ufffdlo\x80 w\x80orl\xc0d!\xc2",
			out:     "Hello world!",
			outFull: "Hello world!",
			err:     transform.ErrShortSrc,
			errSpan: transform.ErrEndOfSpan,
			t:       Remove(Predicate(func(r rune) bool ***REMOVED*** return r == utf8.RuneError ***REMOVED***)),
		***REMOVED***,
		17: ***REMOVED***
			szDst:   4,
			atEOF:   true,
			in:      "Hellø",
			out:     "Hll",
			outFull: "Hllø",
			err:     transform.ErrShortDst,
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		18: ***REMOVED***
			szDst:   4,
			atEOF:   false,
			in:      "Hellø",
			out:     "Hll",
			outFull: "Hllø",
			err:     transform.ErrShortDst,
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		19: ***REMOVED***
			szDst:   8,
			atEOF:   false,
			in:      "\x80Hello\uFF24\x80",
			out:     "\ufffdHll",
			outFull: "\ufffdHll\ufffd",
			err:     transform.ErrShortDst,
			errSpan: transform.ErrEndOfSpan,
			t:       remove,
		***REMOVED***,
		20: ***REMOVED***
			szDst:   8,
			atEOF:   false,
			in:      "Hllll",
			out:     "Hllll",
			outFull: "Hllll",
			t:       remove,
		***REMOVED******REMOVED*** ***REMOVED***
		tt.check(t, i)
	***REMOVED***
***REMOVED***

func TestReplaceIllFormed(t *testing.T) ***REMOVED***
	replace := ReplaceIllFormed()

	for i, tt := range []transformTest***REMOVED***
		0: ***REMOVED***
			szDst:   large,
			atEOF:   true,
			in:      "",
			out:     "",
			outFull: "",
			t:       replace,
		***REMOVED***,
		1: ***REMOVED***
			szDst:   1,
			atEOF:   true,
			in:      "aa",
			out:     "a",
			outFull: "aa",
			err:     transform.ErrShortDst,
			t:       replace,
		***REMOVED***,
		2: ***REMOVED***
			szDst:   1,
			atEOF:   true,
			in:      "a\x80",
			out:     "a",
			outFull: "a\ufffd",
			err:     transform.ErrShortDst,
			errSpan: transform.ErrEndOfSpan,
			t:       replace,
		***REMOVED***,
		3: ***REMOVED***
			szDst:   1,
			atEOF:   true,
			in:      "a\xc2",
			out:     "a",
			outFull: "a\ufffd",
			err:     transform.ErrShortDst,
			errSpan: transform.ErrEndOfSpan,
			t:       replace,
		***REMOVED***,
		4: ***REMOVED***
			szDst:   large,
			atEOF:   true,
			in:      "\x80",
			out:     "\ufffd",
			outFull: "\ufffd",
			errSpan: transform.ErrEndOfSpan,
			t:       replace,
		***REMOVED***,
		5: ***REMOVED***
			szDst:   large,
			atEOF:   false,
			in:      "\x80",
			out:     "\ufffd",
			outFull: "\ufffd",
			errSpan: transform.ErrEndOfSpan,
			t:       replace,
		***REMOVED***,
		6: ***REMOVED***
			szDst:   large,
			atEOF:   true,
			in:      "\xc2",
			out:     "\ufffd",
			outFull: "\ufffd",
			errSpan: transform.ErrEndOfSpan,
			t:       replace,
		***REMOVED***,
		7: ***REMOVED***
			szDst:   large,
			atEOF:   false,
			in:      "\xc2",
			out:     "",
			outFull: "\ufffd",
			err:     transform.ErrShortSrc,
			errSpan: transform.ErrShortSrc,
			t:       replace,
		***REMOVED***,
		8: ***REMOVED***
			szDst:   large,
			atEOF:   true,
			in:      "Hello world!",
			out:     "Hello world!",
			outFull: "Hello world!",
			t:       replace,
		***REMOVED***,
		9: ***REMOVED***
			szDst:   large,
			atEOF:   true,
			in:      "Hello\x80 w\x80orl\xc2d!\xc2",
			out:     "Hello\ufffd w\ufffdorl\ufffdd!\ufffd",
			outFull: "Hello\ufffd w\ufffdorl\ufffdd!\ufffd",
			errSpan: transform.ErrEndOfSpan,
			t:       replace,
		***REMOVED***,
		10: ***REMOVED***
			szDst:   large,
			atEOF:   false,
			in:      "Hello\x80 w\x80orl\xc2d!\xc2",
			out:     "Hello\ufffd w\ufffdorl\ufffdd!",
			outFull: "Hello\ufffd w\ufffdorl\ufffdd!\ufffd",
			err:     transform.ErrShortSrc,
			errSpan: transform.ErrEndOfSpan,
			t:       replace,
		***REMOVED***,
		16: ***REMOVED***
			szDst:   10,
			atEOF:   false,
			in:      "\x80Hello\x80",
			out:     "\ufffdHello",
			outFull: "\ufffdHello\ufffd",
			err:     transform.ErrShortDst,
			errSpan: transform.ErrEndOfSpan,
			t:       replace,
		***REMOVED***,
		17: ***REMOVED***
			szDst:   10,
			atEOF:   false,
			in:      "\ufffdHello\ufffd",
			out:     "\ufffdHello",
			outFull: "\ufffdHello\ufffd",
			err:     transform.ErrShortDst,
			t:       replace,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		tt.check(t, i)
	***REMOVED***
***REMOVED***

func TestMapAlloc(t *testing.T) ***REMOVED***
	if n := testtext.AllocsPerRun(3, func() ***REMOVED***
		Map(idem).Transform(nil, nil, false)
	***REMOVED***); n > 0 ***REMOVED***
		t.Errorf("got %f; want 0", n)
	***REMOVED***
***REMOVED***

func rmNop(r rune) bool ***REMOVED*** return false ***REMOVED***

func TestRemoveAlloc(t *testing.T) ***REMOVED***
	if n := testtext.AllocsPerRun(3, func() ***REMOVED***
		Remove(Predicate(rmNop)).Transform(nil, nil, false)
	***REMOVED***); n > 0 ***REMOVED***
		t.Errorf("got %f; want 0", n)
	***REMOVED***
***REMOVED***

func TestReplaceIllFormedAlloc(t *testing.T) ***REMOVED***
	if n := testtext.AllocsPerRun(3, func() ***REMOVED***
		ReplaceIllFormed().Transform(nil, nil, false)
	***REMOVED***); n > 0 ***REMOVED***
		t.Errorf("got %f; want 0", n)
	***REMOVED***
***REMOVED***

func doBench(b *testing.B, t Transformer) ***REMOVED***
	for _, bc := range []struct***REMOVED*** name, data string ***REMOVED******REMOVED***
		***REMOVED***"ascii", testtext.ASCII***REMOVED***,
		***REMOVED***"3byte", testtext.ThreeByteUTF8***REMOVED***,
	***REMOVED*** ***REMOVED***
		dst := make([]byte, 2*len(bc.data))
		src := []byte(bc.data)

		testtext.Bench(b, bc.name+"/transform", func(b *testing.B) ***REMOVED***
			b.SetBytes(int64(len(src)))
			for i := 0; i < b.N; i++ ***REMOVED***
				t.Transform(dst, src, true)
			***REMOVED***
		***REMOVED***)
		src = t.Bytes(src)
		t.Reset()
		testtext.Bench(b, bc.name+"/span", func(b *testing.B) ***REMOVED***
			b.SetBytes(int64(len(src)))
			for i := 0; i < b.N; i++ ***REMOVED***
				t.Span(src, true)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkRemove(b *testing.B) ***REMOVED***
	doBench(b, Remove(Predicate(func(r rune) bool ***REMOVED*** return r == 'e' ***REMOVED***)))
***REMOVED***

func BenchmarkMapAll(b *testing.B) ***REMOVED***
	doBench(b, Map(func(r rune) rune ***REMOVED*** return 'a' ***REMOVED***))
***REMOVED***

func BenchmarkMapNone(b *testing.B) ***REMOVED***
	doBench(b, Map(func(r rune) rune ***REMOVED*** return r ***REMOVED***))
***REMOVED***

func BenchmarkReplaceIllFormed(b *testing.B) ***REMOVED***
	doBench(b, ReplaceIllFormed())
***REMOVED***

var (
	input = strings.Repeat("Thé qüick brøwn føx jumps øver the lazy døg. ", 100)
)
