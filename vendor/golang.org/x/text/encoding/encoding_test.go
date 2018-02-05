// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encoding_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func TestEncodeInvalidUTF8(t *testing.T) ***REMOVED***
	inputs := []string***REMOVED***
		"hello.",
		"wo\ufffdld.",
		"ABC\xff\x80\x80", // Invalid UTF-8.
		"\x80\x80\x80\x80\x80",
		"\x80\x80D\x80\x80",          // Valid rune at "D".
		"E\xed\xa0\x80\xed\xbf\xbfF", // Two invalid UTF-8 runes (surrogates).
		"G",
		"H\xe2\x82",     // U+20AC in UTF-8 is "\xe2\x82\xac", which we split over two
		"\xacI\xe2\x82", // input lines. It maps to 0x80 in the Windows-1252 encoding.
	***REMOVED***
	// Each invalid source byte becomes '\x1a'.
	want := strings.Replace("hello.wo?ld.ABC??????????D??E??????FGH\x80I??", "?", "\x1a", -1)

	transformer := encoding.ReplaceUnsupported(charmap.Windows1252.NewEncoder())
	gotBuf := make([]byte, 0, 1024)
	src := make([]byte, 0, 1024)
	for i, input := range inputs ***REMOVED***
		dst := make([]byte, 1024)
		src = append(src, input...)
		atEOF := i == len(inputs)-1
		nDst, nSrc, err := transformer.Transform(dst, src, atEOF)
		gotBuf = append(gotBuf, dst[:nDst]...)
		src = src[nSrc:]
		if err != nil && err != transform.ErrShortSrc ***REMOVED***
			t.Fatalf("i=%d: %v", i, err)
		***REMOVED***
		if atEOF && err != nil ***REMOVED***
			t.Fatalf("i=%d: atEOF: %v", i, err)
		***REMOVED***
	***REMOVED***
	if got := string(gotBuf); got != want ***REMOVED***
		t.Fatalf("\ngot  %+q\nwant %+q", got, want)
	***REMOVED***
***REMOVED***

func TestReplacement(t *testing.T) ***REMOVED***
	for _, direction := range []string***REMOVED***"Decode", "Encode"***REMOVED*** ***REMOVED***
		enc, want := (transform.Transformer)(nil), ""
		if direction == "Decode" ***REMOVED***
			enc = encoding.Replacement.NewDecoder()
			want = "\ufffd"
		***REMOVED*** else ***REMOVED***
			enc = encoding.Replacement.NewEncoder()
			want = "AB\x00CD\ufffdYZ"
		***REMOVED***
		sr := strings.NewReader("AB\x00CD\x80YZ")
		g, err := ioutil.ReadAll(transform.NewReader(sr, enc))
		if err != nil ***REMOVED***
			t.Errorf("%s: ReadAll: %v", direction, err)
			continue
		***REMOVED***
		if got := string(g); got != want ***REMOVED***
			t.Errorf("%s:\ngot  %q\nwant %q", direction, got, want)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUTF8Validator(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		desc    string
		dstSize int
		src     string
		atEOF   bool
		want    string
		wantErr error
	***REMOVED******REMOVED***
		***REMOVED***
			"empty input",
			100,
			"",
			false,
			"",
			nil,
		***REMOVED***,
		***REMOVED***
			"valid 1-byte 1-rune input",
			100,
			"a",
			false,
			"a",
			nil,
		***REMOVED***,
		***REMOVED***
			"valid 3-byte 1-rune input",
			100,
			"\u1234",
			false,
			"\u1234",
			nil,
		***REMOVED***,
		***REMOVED***
			"valid 5-byte 3-rune input",
			100,
			"a\u0100\u0101",
			false,
			"a\u0100\u0101",
			nil,
		***REMOVED***,
		***REMOVED***
			"perfectly sized dst (non-ASCII)",
			5,
			"a\u0100\u0101",
			false,
			"a\u0100\u0101",
			nil,
		***REMOVED***,
		***REMOVED***
			"short dst (non-ASCII)",
			4,
			"a\u0100\u0101",
			false,
			"a\u0100",
			transform.ErrShortDst,
		***REMOVED***,
		***REMOVED***
			"perfectly sized dst (ASCII)",
			5,
			"abcde",
			false,
			"abcde",
			nil,
		***REMOVED***,
		***REMOVED***
			"short dst (ASCII)",
			4,
			"abcde",
			false,
			"abcd",
			transform.ErrShortDst,
		***REMOVED***,
		***REMOVED***
			"partial input (!EOF)",
			100,
			"a\u0100\xf1",
			false,
			"a\u0100",
			transform.ErrShortSrc,
		***REMOVED***,
		***REMOVED***
			"invalid input (EOF)",
			100,
			"a\u0100\xf1",
			true,
			"a\u0100",
			encoding.ErrInvalidUTF8,
		***REMOVED***,
		***REMOVED***
			"invalid input (!EOF)",
			100,
			"a\u0100\x80",
			false,
			"a\u0100",
			encoding.ErrInvalidUTF8,
		***REMOVED***,
		***REMOVED***
			"invalid input (above U+10FFFF)",
			100,
			"a\u0100\xf7\xbf\xbf\xbf",
			false,
			"a\u0100",
			encoding.ErrInvalidUTF8,
		***REMOVED***,
		***REMOVED***
			"invalid input (surrogate half)",
			100,
			"a\u0100\xed\xa0\x80",
			false,
			"a\u0100",
			encoding.ErrInvalidUTF8,
		***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		dst := make([]byte, tc.dstSize)
		nDst, nSrc, err := encoding.UTF8Validator.Transform(dst, []byte(tc.src), tc.atEOF)
		if nDst < 0 || len(dst) < nDst ***REMOVED***
			t.Errorf("%s: nDst=%d out of range", tc.desc, nDst)
			continue
		***REMOVED***
		got := string(dst[:nDst])
		if got != tc.want || nSrc != len(tc.want) || err != tc.wantErr ***REMOVED***
			t.Errorf("%s:\ngot  %+q, %d, %v\nwant %+q, %d, %v",
				tc.desc, got, nSrc, err, tc.want, len(tc.want), tc.wantErr)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestErrorHandler(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		desc      string
		handler   func(*encoding.Encoder) *encoding.Encoder
		sizeDst   int
		src, want string
		nSrc      int
		err       error
	***REMOVED******REMOVED***
		***REMOVED***
			desc:    "one rune replacement",
			handler: encoding.ReplaceUnsupported,
			sizeDst: 100,
			src:     "\uAC00",
			want:    "\x1a",
			nSrc:    3,
		***REMOVED***,
		***REMOVED***
			desc:    "mid-stream rune replacement",
			handler: encoding.ReplaceUnsupported,
			sizeDst: 100,
			src:     "a\uAC00bcd\u00e9",
			want:    "a\x1abcd\xe9",
			nSrc:    9,
		***REMOVED***,
		***REMOVED***
			desc:    "at end rune replacement",
			handler: encoding.ReplaceUnsupported,
			sizeDst: 10,
			src:     "\u00e9\uAC00",
			want:    "\xe9\x1a",
			nSrc:    5,
		***REMOVED***,
		***REMOVED***
			desc:    "short buffer replacement",
			handler: encoding.ReplaceUnsupported,
			sizeDst: 1,
			src:     "\u00e9\uAC00",
			want:    "\xe9",
			nSrc:    2,
			err:     transform.ErrShortDst,
		***REMOVED***,
		***REMOVED***
			desc:    "one rune html escape",
			handler: encoding.HTMLEscapeUnsupported,
			sizeDst: 100,
			src:     "\uAC00",
			want:    "&#44032;",
			nSrc:    3,
		***REMOVED***,
		***REMOVED***
			desc:    "mid-stream html escape",
			handler: encoding.HTMLEscapeUnsupported,
			sizeDst: 100,
			src:     "\u00e9\uAC00dcba",
			want:    "\xe9&#44032;dcba",
			nSrc:    9,
		***REMOVED***,
		***REMOVED***
			desc:    "short buffer html escape",
			handler: encoding.HTMLEscapeUnsupported,
			sizeDst: 9,
			src:     "ab\uAC01",
			want:    "ab",
			nSrc:    2,
			err:     transform.ErrShortDst,
		***REMOVED***,
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		tr := tc.handler(charmap.Windows1250.NewEncoder())
		b := make([]byte, tc.sizeDst)
		nDst, nSrc, err := tr.Transform(b, []byte(tc.src), true)
		if err != tc.err ***REMOVED***
			t.Errorf("%d:%s: error was %v; want %v", i, tc.desc, err, tc.err)
		***REMOVED***
		if got := string(b[:nDst]); got != tc.want ***REMOVED***
			t.Errorf("%d:%s: result was %q: want %q", i, tc.desc, got, tc.want)
		***REMOVED***
		if nSrc != tc.nSrc ***REMOVED***
			t.Errorf("%d:%s: nSrc was %d; want %d", i, tc.desc, nSrc, tc.nSrc)
		***REMOVED***

	***REMOVED***
***REMOVED***
