// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"golang.org/x/text/internal/testtext"
)

type lowerCaseASCII struct***REMOVED*** NopResetter ***REMOVED***

func (lowerCaseASCII) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	n := len(src)
	if n > len(dst) ***REMOVED***
		n, err = len(dst), ErrShortDst
	***REMOVED***
	for i, c := range src[:n] ***REMOVED***
		if 'A' <= c && c <= 'Z' ***REMOVED***
			c += 'a' - 'A'
		***REMOVED***
		dst[i] = c
	***REMOVED***
	return n, n, err
***REMOVED***

// lowerCaseASCIILookahead lowercases the string and reports ErrShortSrc as long
// as the input is not atEOF.
type lowerCaseASCIILookahead struct***REMOVED*** NopResetter ***REMOVED***

func (lowerCaseASCIILookahead) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	n := len(src)
	if n > len(dst) ***REMOVED***
		n, err = len(dst), ErrShortDst
	***REMOVED***
	for i, c := range src[:n] ***REMOVED***
		if 'A' <= c && c <= 'Z' ***REMOVED***
			c += 'a' - 'A'
		***REMOVED***
		dst[i] = c
	***REMOVED***
	if !atEOF ***REMOVED***
		err = ErrShortSrc
	***REMOVED***
	return n, n, err
***REMOVED***

var errYouMentionedX = errors.New("you mentioned X")

type dontMentionX struct***REMOVED*** NopResetter ***REMOVED***

func (dontMentionX) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	n := len(src)
	if n > len(dst) ***REMOVED***
		n, err = len(dst), ErrShortDst
	***REMOVED***
	for i, c := range src[:n] ***REMOVED***
		if c == 'X' ***REMOVED***
			return i, i, errYouMentionedX
		***REMOVED***
		dst[i] = c
	***REMOVED***
	return n, n, err
***REMOVED***

var errAtEnd = errors.New("error after all text")

type errorAtEnd struct***REMOVED*** NopResetter ***REMOVED***

func (errorAtEnd) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	n := copy(dst, src)
	if n < len(src) ***REMOVED***
		return n, n, ErrShortDst
	***REMOVED***
	if atEOF ***REMOVED***
		return n, n, errAtEnd
	***REMOVED***
	return n, n, nil
***REMOVED***

type replaceWithConstant struct ***REMOVED***
	replacement string
	written     int
***REMOVED***

func (t *replaceWithConstant) Reset() ***REMOVED***
	t.written = 0
***REMOVED***

func (t *replaceWithConstant) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	if atEOF ***REMOVED***
		nDst = copy(dst, t.replacement[t.written:])
		t.written += nDst
		if t.written < len(t.replacement) ***REMOVED***
			err = ErrShortDst
		***REMOVED***
	***REMOVED***
	return nDst, len(src), err
***REMOVED***

type addAnXAtTheEnd struct***REMOVED*** NopResetter ***REMOVED***

func (addAnXAtTheEnd) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	n := copy(dst, src)
	if n < len(src) ***REMOVED***
		return n, n, ErrShortDst
	***REMOVED***
	if !atEOF ***REMOVED***
		return n, n, nil
	***REMOVED***
	if len(dst) == n ***REMOVED***
		return n, n, ErrShortDst
	***REMOVED***
	dst[n] = 'X'
	return n + 1, n, nil
***REMOVED***

// doublerAtEOF is a strange Transformer that transforms "this" to "tthhiiss",
// but only if atEOF is true.
type doublerAtEOF struct***REMOVED*** NopResetter ***REMOVED***

func (doublerAtEOF) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	if !atEOF ***REMOVED***
		return 0, 0, ErrShortSrc
	***REMOVED***
	for i, c := range src ***REMOVED***
		if 2*i+2 >= len(dst) ***REMOVED***
			return 2 * i, i, ErrShortDst
		***REMOVED***
		dst[2*i+0] = c
		dst[2*i+1] = c
	***REMOVED***
	return 2 * len(src), len(src), nil
***REMOVED***

// rleDecode and rleEncode implement a toy run-length encoding: "aabbbbbbbbbb"
// is encoded as "2a10b". The decoding is assumed to not contain any numbers.

type rleDecode struct***REMOVED*** NopResetter ***REMOVED***

func (rleDecode) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
loop:
	for len(src) > 0 ***REMOVED***
		n := 0
		for i, c := range src ***REMOVED***
			if '0' <= c && c <= '9' ***REMOVED***
				n = 10*n + int(c-'0')
				continue
			***REMOVED***
			if i == 0 ***REMOVED***
				return nDst, nSrc, errors.New("rleDecode: bad input")
			***REMOVED***
			if n > len(dst) ***REMOVED***
				return nDst, nSrc, ErrShortDst
			***REMOVED***
			for j := 0; j < n; j++ ***REMOVED***
				dst[j] = c
			***REMOVED***
			dst, src = dst[n:], src[i+1:]
			nDst, nSrc = nDst+n, nSrc+i+1
			continue loop
		***REMOVED***
		if atEOF ***REMOVED***
			return nDst, nSrc, errors.New("rleDecode: bad input")
		***REMOVED***
		return nDst, nSrc, ErrShortSrc
	***REMOVED***
	return nDst, nSrc, nil
***REMOVED***

type rleEncode struct ***REMOVED***
	NopResetter

	// allowStutter means that "xxxxxxxx" can be encoded as "5x3x"
	// instead of always as "8x".
	allowStutter bool
***REMOVED***

func (e rleEncode) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	for len(src) > 0 ***REMOVED***
		n, c0 := len(src), src[0]
		for i, c := range src[1:] ***REMOVED***
			if c != c0 ***REMOVED***
				n = i + 1
				break
			***REMOVED***
		***REMOVED***
		if n == len(src) && !atEOF && !e.allowStutter ***REMOVED***
			return nDst, nSrc, ErrShortSrc
		***REMOVED***
		s := strconv.Itoa(n)
		if len(s) >= len(dst) ***REMOVED***
			return nDst, nSrc, ErrShortDst
		***REMOVED***
		copy(dst, s)
		dst[len(s)] = c0
		dst, src = dst[len(s)+1:], src[n:]
		nDst, nSrc = nDst+len(s)+1, nSrc+n
	***REMOVED***
	return nDst, nSrc, nil
***REMOVED***

// trickler consumes all input bytes, but writes a single byte at a time to dst.
type trickler []byte

func (t *trickler) Reset() ***REMOVED***
	*t = nil
***REMOVED***

func (t *trickler) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	*t = append(*t, src...)
	if len(*t) == 0 ***REMOVED***
		return 0, 0, nil
	***REMOVED***
	if len(dst) == 0 ***REMOVED***
		return 0, len(src), ErrShortDst
	***REMOVED***
	dst[0] = (*t)[0]
	*t = (*t)[1:]
	if len(*t) > 0 ***REMOVED***
		err = ErrShortDst
	***REMOVED***
	return 1, len(src), err
***REMOVED***

// delayedTrickler is like trickler, but delays writing output to dst. This is
// highly unlikely to be relevant in practice, but it seems like a good idea
// to have some tolerance as long as progress can be detected.
type delayedTrickler []byte

func (t *delayedTrickler) Reset() ***REMOVED***
	*t = nil
***REMOVED***

func (t *delayedTrickler) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	if len(*t) > 0 && len(dst) > 0 ***REMOVED***
		dst[0] = (*t)[0]
		*t = (*t)[1:]
		nDst = 1
	***REMOVED***
	*t = append(*t, src...)
	if len(*t) > 0 ***REMOVED***
		err = ErrShortDst
	***REMOVED***
	return nDst, len(src), err
***REMOVED***

type testCase struct ***REMOVED***
	desc     string
	t        Transformer
	src      string
	dstSize  int
	srcSize  int
	ioSize   int
	wantStr  string
	wantErr  error
	wantIter int // number of iterations taken; 0 means we don't care.
***REMOVED***

func (t testCase) String() string ***REMOVED***
	return tstr(t.t) + "; " + t.desc
***REMOVED***

func tstr(t Transformer) string ***REMOVED***
	if stringer, ok := t.(fmt.Stringer); ok ***REMOVED***
		return stringer.String()
	***REMOVED***
	s := fmt.Sprintf("%T", t)
	return s[1+strings.Index(s, "."):]
***REMOVED***

func (c chain) String() string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***
	buf.WriteString("Chain(")
	for i, l := range c.link[:len(c.link)-1] ***REMOVED***
		if i != 0 ***REMOVED***
			fmt.Fprint(buf, ", ")
		***REMOVED***
		buf.WriteString(tstr(l.t))
	***REMOVED***
	buf.WriteString(")")
	return buf.String()
***REMOVED***

var testCases = []testCase***REMOVED***
	***REMOVED***
		desc:    "empty",
		t:       lowerCaseASCII***REMOVED******REMOVED***,
		src:     "",
		dstSize: 100,
		srcSize: 100,
		wantStr: "",
	***REMOVED***,

	***REMOVED***
		desc:    "basic",
		t:       lowerCaseASCII***REMOVED******REMOVED***,
		src:     "Hello WORLD.",
		dstSize: 100,
		srcSize: 100,
		wantStr: "hello world.",
	***REMOVED***,

	***REMOVED***
		desc:    "small dst",
		t:       lowerCaseASCII***REMOVED******REMOVED***,
		src:     "Hello WORLD.",
		dstSize: 3,
		srcSize: 100,
		wantStr: "hello world.",
	***REMOVED***,

	***REMOVED***
		desc:    "small src",
		t:       lowerCaseASCII***REMOVED******REMOVED***,
		src:     "Hello WORLD.",
		dstSize: 100,
		srcSize: 4,
		wantStr: "hello world.",
	***REMOVED***,

	***REMOVED***
		desc:    "small buffers",
		t:       lowerCaseASCII***REMOVED******REMOVED***,
		src:     "Hello WORLD.",
		dstSize: 3,
		srcSize: 4,
		wantStr: "hello world.",
	***REMOVED***,

	***REMOVED***
		desc:    "very small buffers",
		t:       lowerCaseASCII***REMOVED******REMOVED***,
		src:     "Hello WORLD.",
		dstSize: 1,
		srcSize: 1,
		wantStr: "hello world.",
	***REMOVED***,

	***REMOVED***
		desc:    "small dst with lookahead",
		t:       lowerCaseASCIILookahead***REMOVED******REMOVED***,
		src:     "Hello WORLD.",
		dstSize: 3,
		srcSize: 100,
		wantStr: "hello world.",
	***REMOVED***,

	***REMOVED***
		desc:    "small src with lookahead",
		t:       lowerCaseASCIILookahead***REMOVED******REMOVED***,
		src:     "Hello WORLD.",
		dstSize: 100,
		srcSize: 4,
		wantStr: "hello world.",
	***REMOVED***,

	***REMOVED***
		desc:    "small buffers with lookahead",
		t:       lowerCaseASCIILookahead***REMOVED******REMOVED***,
		src:     "Hello WORLD.",
		dstSize: 3,
		srcSize: 4,
		wantStr: "hello world.",
	***REMOVED***,

	***REMOVED***
		desc:    "very small buffers with lookahead",
		t:       lowerCaseASCIILookahead***REMOVED******REMOVED***,
		src:     "Hello WORLD.",
		dstSize: 1,
		srcSize: 2,
		wantStr: "hello world.",
	***REMOVED***,

	***REMOVED***
		desc:    "user error",
		t:       dontMentionX***REMOVED******REMOVED***,
		src:     "The First Rule of Transform Club: don't mention Mister X, ever.",
		dstSize: 100,
		srcSize: 100,
		wantStr: "The First Rule of Transform Club: don't mention Mister ",
		wantErr: errYouMentionedX,
	***REMOVED***,

	***REMOVED***
		desc:    "user error at end",
		t:       errorAtEnd***REMOVED******REMOVED***,
		src:     "All goes well until it doesn't.",
		dstSize: 100,
		srcSize: 100,
		wantStr: "All goes well until it doesn't.",
		wantErr: errAtEnd,
	***REMOVED***,

	***REMOVED***
		desc:    "user error at end, incremental",
		t:       errorAtEnd***REMOVED******REMOVED***,
		src:     "All goes well until it doesn't.",
		dstSize: 10,
		srcSize: 10,
		wantStr: "All goes well until it doesn't.",
		wantErr: errAtEnd,
	***REMOVED***,

	***REMOVED***
		desc:    "replace entire non-empty string with one byte",
		t:       &replaceWithConstant***REMOVED***replacement: "X"***REMOVED***,
		src:     "none of this will be copied",
		dstSize: 1,
		srcSize: 10,
		wantStr: "X",
	***REMOVED***,

	***REMOVED***
		desc:    "replace entire empty string with one byte",
		t:       &replaceWithConstant***REMOVED***replacement: "X"***REMOVED***,
		src:     "",
		dstSize: 1,
		srcSize: 10,
		wantStr: "X",
	***REMOVED***,

	***REMOVED***
		desc:    "replace entire empty string with seven bytes",
		t:       &replaceWithConstant***REMOVED***replacement: "ABCDEFG"***REMOVED***,
		src:     "",
		dstSize: 3,
		srcSize: 10,
		wantStr: "ABCDEFG",
	***REMOVED***,

	***REMOVED***
		desc:    "add an X (initialBufSize-1)",
		t:       addAnXAtTheEnd***REMOVED******REMOVED***,
		src:     aaa[:initialBufSize-1],
		dstSize: 10,
		srcSize: 10,
		wantStr: aaa[:initialBufSize-1] + "X",
	***REMOVED***,

	***REMOVED***
		desc:    "add an X (initialBufSize+0)",
		t:       addAnXAtTheEnd***REMOVED******REMOVED***,
		src:     aaa[:initialBufSize+0],
		dstSize: 10,
		srcSize: 10,
		wantStr: aaa[:initialBufSize+0] + "X",
	***REMOVED***,

	***REMOVED***
		desc:    "add an X (initialBufSize+1)",
		t:       addAnXAtTheEnd***REMOVED******REMOVED***,
		src:     aaa[:initialBufSize+1],
		dstSize: 10,
		srcSize: 10,
		wantStr: aaa[:initialBufSize+1] + "X",
	***REMOVED***,

	***REMOVED***
		desc:    "small buffers",
		t:       dontMentionX***REMOVED******REMOVED***,
		src:     "The First Rule of Transform Club: don't mention Mister X, ever.",
		dstSize: 10,
		srcSize: 10,
		wantStr: "The First Rule of Transform Club: don't mention Mister ",
		wantErr: errYouMentionedX,
	***REMOVED***,

	***REMOVED***
		desc:    "very small buffers",
		t:       dontMentionX***REMOVED******REMOVED***,
		src:     "The First Rule of Transform Club: don't mention Mister X, ever.",
		dstSize: 1,
		srcSize: 1,
		wantStr: "The First Rule of Transform Club: don't mention Mister ",
		wantErr: errYouMentionedX,
	***REMOVED***,

	***REMOVED***
		desc:    "only transform at EOF",
		t:       doublerAtEOF***REMOVED******REMOVED***,
		src:     "this",
		dstSize: 100,
		srcSize: 100,
		wantStr: "tthhiiss",
	***REMOVED***,

	***REMOVED***
		desc:    "basic",
		t:       rleDecode***REMOVED******REMOVED***,
		src:     "1a2b3c10d11e0f1g",
		dstSize: 100,
		srcSize: 100,
		wantStr: "abbcccddddddddddeeeeeeeeeeeg",
	***REMOVED***,

	***REMOVED***
		desc:    "long",
		t:       rleDecode***REMOVED******REMOVED***,
		src:     "12a23b34c45d56e99z",
		dstSize: 100,
		srcSize: 100,
		wantStr: strings.Repeat("a", 12) +
			strings.Repeat("b", 23) +
			strings.Repeat("c", 34) +
			strings.Repeat("d", 45) +
			strings.Repeat("e", 56) +
			strings.Repeat("z", 99),
	***REMOVED***,

	***REMOVED***
		desc:    "tight buffers",
		t:       rleDecode***REMOVED******REMOVED***,
		src:     "1a2b3c10d11e0f1g",
		dstSize: 11,
		srcSize: 3,
		wantStr: "abbcccddddddddddeeeeeeeeeeeg",
	***REMOVED***,

	***REMOVED***
		desc:    "short dst",
		t:       rleDecode***REMOVED******REMOVED***,
		src:     "1a2b3c10d11e0f1g",
		dstSize: 10,
		srcSize: 3,
		wantStr: "abbcccdddddddddd",
		wantErr: ErrShortDst,
	***REMOVED***,

	***REMOVED***
		desc:    "short src",
		t:       rleDecode***REMOVED******REMOVED***,
		src:     "1a2b3c10d11e0f1g",
		dstSize: 11,
		srcSize: 2,
		ioSize:  2,
		wantStr: "abbccc",
		wantErr: ErrShortSrc,
	***REMOVED***,

	***REMOVED***
		desc:    "basic",
		t:       rleEncode***REMOVED******REMOVED***,
		src:     "abbcccddddddddddeeeeeeeeeeeg",
		dstSize: 100,
		srcSize: 100,
		wantStr: "1a2b3c10d11e1g",
	***REMOVED***,

	***REMOVED***
		desc: "long",
		t:    rleEncode***REMOVED******REMOVED***,
		src: strings.Repeat("a", 12) +
			strings.Repeat("b", 23) +
			strings.Repeat("c", 34) +
			strings.Repeat("d", 45) +
			strings.Repeat("e", 56) +
			strings.Repeat("z", 99),
		dstSize: 100,
		srcSize: 100,
		wantStr: "12a23b34c45d56e99z",
	***REMOVED***,

	***REMOVED***
		desc:    "tight buffers",
		t:       rleEncode***REMOVED******REMOVED***,
		src:     "abbcccddddddddddeeeeeeeeeeeg",
		dstSize: 3,
		srcSize: 12,
		wantStr: "1a2b3c10d11e1g",
	***REMOVED***,

	***REMOVED***
		desc:    "short dst",
		t:       rleEncode***REMOVED******REMOVED***,
		src:     "abbcccddddddddddeeeeeeeeeeeg",
		dstSize: 2,
		srcSize: 12,
		wantStr: "1a2b3c",
		wantErr: ErrShortDst,
	***REMOVED***,

	***REMOVED***
		desc:    "short src",
		t:       rleEncode***REMOVED******REMOVED***,
		src:     "abbcccddddddddddeeeeeeeeeeeg",
		dstSize: 3,
		srcSize: 11,
		ioSize:  11,
		wantStr: "1a2b3c10d",
		wantErr: ErrShortSrc,
	***REMOVED***,

	***REMOVED***
		desc:    "allowStutter = false",
		t:       rleEncode***REMOVED***allowStutter: false***REMOVED***,
		src:     "aaaabbbbbbbbccccddddd",
		dstSize: 10,
		srcSize: 10,
		wantStr: "4a8b4c5d",
	***REMOVED***,

	***REMOVED***
		desc:    "allowStutter = true",
		t:       rleEncode***REMOVED***allowStutter: true***REMOVED***,
		src:     "aaaabbbbbbbbccccddddd",
		dstSize: 10,
		srcSize: 10,
		ioSize:  10,
		wantStr: "4a6b2b4c4d1d",
	***REMOVED***,

	***REMOVED***
		desc:    "trickler",
		t:       &trickler***REMOVED******REMOVED***,
		src:     "abcdefghijklm",
		dstSize: 3,
		srcSize: 15,
		wantStr: "abcdefghijklm",
	***REMOVED***,

	***REMOVED***
		desc:    "delayedTrickler",
		t:       &delayedTrickler***REMOVED******REMOVED***,
		src:     "abcdefghijklm",
		dstSize: 3,
		srcSize: 15,
		wantStr: "abcdefghijklm",
	***REMOVED***,
***REMOVED***

func TestReader(t *testing.T) ***REMOVED***
	for _, tc := range testCases ***REMOVED***
		testtext.Run(t, tc.desc, func(t *testing.T) ***REMOVED***
			r := NewReader(strings.NewReader(tc.src), tc.t)
			// Differently sized dst and src buffers are not part of the
			// exported API. We override them manually.
			r.dst = make([]byte, tc.dstSize)
			r.src = make([]byte, tc.srcSize)
			got, err := ioutil.ReadAll(r)
			str := string(got)
			if str != tc.wantStr || err != tc.wantErr ***REMOVED***
				t.Errorf("\ngot  %q, %v\nwant %q, %v", str, err, tc.wantStr, tc.wantErr)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestWriter(t *testing.T) ***REMOVED***
	tests := append(testCases, chainTests()...)
	for _, tc := range tests ***REMOVED***
		sizes := []int***REMOVED***1, 2, 3, 4, 5, 10, 100, 1000***REMOVED***
		if tc.ioSize > 0 ***REMOVED***
			sizes = []int***REMOVED***tc.ioSize***REMOVED***
		***REMOVED***
		for _, sz := range sizes ***REMOVED***
			testtext.Run(t, fmt.Sprintf("%s/%d", tc.desc, sz), func(t *testing.T) ***REMOVED***
				bb := &bytes.Buffer***REMOVED******REMOVED***
				w := NewWriter(bb, tc.t)
				// Differently sized dst and src buffers are not part of the
				// exported API. We override them manually.
				w.dst = make([]byte, tc.dstSize)
				w.src = make([]byte, tc.srcSize)
				src := make([]byte, sz)
				var err error
				for b := tc.src; len(b) > 0 && err == nil; ***REMOVED***
					n := copy(src, b)
					b = b[n:]
					m := 0
					m, err = w.Write(src[:n])
					if m != n && err == nil ***REMOVED***
						t.Errorf("did not consume all bytes %d < %d", m, n)
					***REMOVED***
				***REMOVED***
				if err == nil ***REMOVED***
					err = w.Close()
				***REMOVED***
				str := bb.String()
				if str != tc.wantStr || err != tc.wantErr ***REMOVED***
					t.Errorf("\ngot  %q, %v\nwant %q, %v", str, err, tc.wantStr, tc.wantErr)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNop(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		str     string
		dstSize int
		err     error
	***REMOVED******REMOVED***
		***REMOVED***"", 0, nil***REMOVED***,
		***REMOVED***"", 10, nil***REMOVED***,
		***REMOVED***"a", 0, ErrShortDst***REMOVED***,
		***REMOVED***"a", 1, nil***REMOVED***,
		***REMOVED***"a", 10, nil***REMOVED***,
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		dst := make([]byte, tc.dstSize)
		nDst, nSrc, err := Nop.Transform(dst, []byte(tc.str), true)
		want := tc.str
		if tc.dstSize < len(want) ***REMOVED***
			want = want[:tc.dstSize]
		***REMOVED***
		if got := string(dst[:nDst]); got != want || err != tc.err || nSrc != nDst ***REMOVED***
			t.Errorf("%d:\ngot %q, %d, %v\nwant %q, %d, %v", i, got, nSrc, err, want, nDst, tc.err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDiscard(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		str     string
		dstSize int
	***REMOVED******REMOVED***
		***REMOVED***"", 0***REMOVED***,
		***REMOVED***"", 10***REMOVED***,
		***REMOVED***"a", 0***REMOVED***,
		***REMOVED***"ab", 10***REMOVED***,
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		nDst, nSrc, err := Discard.Transform(make([]byte, tc.dstSize), []byte(tc.str), true)
		if nDst != 0 || nSrc != len(tc.str) || err != nil ***REMOVED***
			t.Errorf("%d:\ngot %q, %d, %v\nwant 0, %d, nil", i, nDst, nSrc, err, len(tc.str))
		***REMOVED***
	***REMOVED***
***REMOVED***

// mkChain creates a Chain transformer. x must be alternating between transformer
// and bufSize, like T, (sz, T)*
func mkChain(x ...interface***REMOVED******REMOVED***) *chain ***REMOVED***
	t := []Transformer***REMOVED******REMOVED***
	for i := 0; i < len(x); i += 2 ***REMOVED***
		t = append(t, x[i].(Transformer))
	***REMOVED***
	c := Chain(t...).(*chain)
	for i, j := 1, 1; i < len(x); i, j = i+2, j+1 ***REMOVED***
		c.link[j].b = make([]byte, x[i].(int))
	***REMOVED***
	return c
***REMOVED***

func chainTests() []testCase ***REMOVED***
	return []testCase***REMOVED***
		***REMOVED***
			desc:     "nil error",
			t:        mkChain(rleEncode***REMOVED******REMOVED***, 100, lowerCaseASCII***REMOVED******REMOVED***),
			src:      "ABB",
			dstSize:  100,
			srcSize:  100,
			wantStr:  "1a2b",
			wantErr:  nil,
			wantIter: 1,
		***REMOVED***,

		***REMOVED***
			desc:    "short dst buffer",
			t:       mkChain(lowerCaseASCII***REMOVED******REMOVED***, 3, rleDecode***REMOVED******REMOVED***),
			src:     "1a2b3c10d11e0f1g",
			dstSize: 10,
			srcSize: 3,
			wantStr: "abbcccdddddddddd",
			wantErr: ErrShortDst,
		***REMOVED***,

		***REMOVED***
			desc:    "short internal dst buffer",
			t:       mkChain(lowerCaseASCII***REMOVED******REMOVED***, 3, rleDecode***REMOVED******REMOVED***, 10, Nop),
			src:     "1a2b3c10d11e0f1g",
			dstSize: 100,
			srcSize: 3,
			wantStr: "abbcccdddddddddd",
			wantErr: errShortInternal,
		***REMOVED***,

		***REMOVED***
			desc:    "short internal dst buffer from input",
			t:       mkChain(rleDecode***REMOVED******REMOVED***, 10, Nop),
			src:     "1a2b3c10d11e0f1g",
			dstSize: 100,
			srcSize: 3,
			wantStr: "abbcccdddddddddd",
			wantErr: errShortInternal,
		***REMOVED***,

		***REMOVED***
			desc:    "empty short internal dst buffer",
			t:       mkChain(lowerCaseASCII***REMOVED******REMOVED***, 3, rleDecode***REMOVED******REMOVED***, 10, Nop),
			src:     "4a7b11e0f1g",
			dstSize: 100,
			srcSize: 3,
			wantStr: "aaaabbbbbbb",
			wantErr: errShortInternal,
		***REMOVED***,

		***REMOVED***
			desc:    "empty short internal dst buffer from input",
			t:       mkChain(rleDecode***REMOVED******REMOVED***, 10, Nop),
			src:     "4a7b11e0f1g",
			dstSize: 100,
			srcSize: 3,
			wantStr: "aaaabbbbbbb",
			wantErr: errShortInternal,
		***REMOVED***,

		***REMOVED***
			desc:     "short internal src buffer after full dst buffer",
			t:        mkChain(Nop, 5, rleEncode***REMOVED******REMOVED***, 10, Nop),
			src:      "cccccddddd",
			dstSize:  100,
			srcSize:  100,
			wantStr:  "",
			wantErr:  errShortInternal,
			wantIter: 1,
		***REMOVED***,

		***REMOVED***
			desc:    "short internal src buffer after short dst buffer; test lastFull",
			t:       mkChain(rleDecode***REMOVED******REMOVED***, 5, rleEncode***REMOVED******REMOVED***, 4, Nop),
			src:     "2a1b4c6d",
			dstSize: 100,
			srcSize: 100,
			wantStr: "2a1b",
			wantErr: errShortInternal,
		***REMOVED***,

		***REMOVED***
			desc:     "short internal src buffer after successful complete fill",
			t:        mkChain(Nop, 3, rleDecode***REMOVED******REMOVED***),
			src:      "123a4b",
			dstSize:  4,
			srcSize:  3,
			wantStr:  "",
			wantErr:  errShortInternal,
			wantIter: 1,
		***REMOVED***,

		***REMOVED***
			desc:    "short internal src buffer after short dst buffer; test lastFull",
			t:       mkChain(rleDecode***REMOVED******REMOVED***, 5, rleEncode***REMOVED******REMOVED***),
			src:     "2a1b4c6d",
			dstSize: 4,
			srcSize: 100,
			wantStr: "2a1b",
			wantErr: errShortInternal,
		***REMOVED***,

		***REMOVED***
			desc:    "short src buffer",
			t:       mkChain(rleEncode***REMOVED******REMOVED***, 5, Nop),
			src:     "abbcccddddeeeee",
			dstSize: 4,
			srcSize: 4,
			ioSize:  4,
			wantStr: "1a2b3c",
			wantErr: ErrShortSrc,
		***REMOVED***,

		***REMOVED***
			desc:     "process all in one go",
			t:        mkChain(rleEncode***REMOVED******REMOVED***, 5, Nop),
			src:      "abbcccddddeeeeeffffff",
			dstSize:  100,
			srcSize:  100,
			wantStr:  "1a2b3c4d5e6f",
			wantErr:  nil,
			wantIter: 1,
		***REMOVED***,

		***REMOVED***
			desc:    "complete processing downstream after error",
			t:       mkChain(dontMentionX***REMOVED******REMOVED***, 2, rleDecode***REMOVED******REMOVED***, 5, Nop),
			src:     "3a4b5eX",
			dstSize: 100,
			srcSize: 100,
			ioSize:  100,
			wantStr: "aaabbbbeeeee",
			wantErr: errYouMentionedX,
		***REMOVED***,

		***REMOVED***
			desc:    "return downstream fatal errors first (followed by short dst)",
			t:       mkChain(dontMentionX***REMOVED******REMOVED***, 8, rleDecode***REMOVED******REMOVED***, 4, Nop),
			src:     "3a4b5eX",
			dstSize: 100,
			srcSize: 100,
			ioSize:  100,
			wantStr: "aaabbbb",
			wantErr: errShortInternal,
		***REMOVED***,

		***REMOVED***
			desc:    "return downstream fatal errors first (followed by short src)",
			t:       mkChain(dontMentionX***REMOVED******REMOVED***, 5, Nop, 1, rleDecode***REMOVED******REMOVED***),
			src:     "1a5bX",
			dstSize: 100,
			srcSize: 100,
			ioSize:  100,
			wantStr: "",
			wantErr: errShortInternal,
		***REMOVED***,

		***REMOVED***
			desc:    "short internal",
			t:       mkChain(Nop, 11, rleEncode***REMOVED******REMOVED***, 3, Nop),
			src:     "abbcccddddddddddeeeeeeeeeeeg",
			dstSize: 3,
			srcSize: 100,
			wantStr: "1a2b3c10d",
			wantErr: errShortInternal,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func doTransform(tc testCase) (res string, iter int, err error) ***REMOVED***
	tc.t.Reset()
	dst := make([]byte, tc.dstSize)
	out, in := make([]byte, 0, 2*len(tc.src)), []byte(tc.src)
	for ***REMOVED***
		iter++
		src, atEOF := in, true
		if len(src) > tc.srcSize ***REMOVED***
			src, atEOF = src[:tc.srcSize], false
		***REMOVED***
		nDst, nSrc, err := tc.t.Transform(dst, src, atEOF)
		out = append(out, dst[:nDst]...)
		in = in[nSrc:]
		switch ***REMOVED***
		case err == nil && len(in) != 0:
		case err == ErrShortSrc && nSrc > 0:
		case err == ErrShortDst && (nDst > 0 || nSrc > 0):
		default:
			return string(out), iter, err
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestChain(t *testing.T) ***REMOVED***
	if c, ok := Chain().(nop); !ok ***REMOVED***
		t.Errorf("empty chain: %v; want Nop", c)
	***REMOVED***

	// Test Chain for a single Transformer.
	for _, tc := range testCases ***REMOVED***
		tc.t = Chain(tc.t)
		str, _, err := doTransform(tc)
		if str != tc.wantStr || err != tc.wantErr ***REMOVED***
			t.Errorf("%s:\ngot  %q, %v\nwant %q, %v", tc, str, err, tc.wantStr, tc.wantErr)
		***REMOVED***
	***REMOVED***

	tests := chainTests()
	sizes := []int***REMOVED***1, 2, 3, 4, 5, 7, 10, 100, 1000***REMOVED***
	addTest := func(tc testCase, t *chain) ***REMOVED***
		if t.link[0].t != tc.t && tc.wantErr == ErrShortSrc ***REMOVED***
			tc.wantErr = errShortInternal
		***REMOVED***
		if t.link[len(t.link)-2].t != tc.t && tc.wantErr == ErrShortDst ***REMOVED***
			tc.wantErr = errShortInternal
		***REMOVED***
		tc.t = t
		tests = append(tests, tc)
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		for _, sz := range sizes ***REMOVED***
			tt := tc
			tt.dstSize = sz
			addTest(tt, mkChain(tc.t, tc.dstSize, Nop))
			addTest(tt, mkChain(tc.t, tc.dstSize, Nop, 2, Nop))
			addTest(tt, mkChain(Nop, tc.srcSize, tc.t, tc.dstSize, Nop))
			if sz >= tc.dstSize && (tc.wantErr != ErrShortDst || sz == tc.dstSize) ***REMOVED***
				addTest(tt, mkChain(Nop, tc.srcSize, tc.t))
				addTest(tt, mkChain(Nop, 100, Nop, tc.srcSize, tc.t))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		tt := tc
		tt.dstSize = 1
		tt.wantStr = ""
		addTest(tt, mkChain(tc.t, tc.dstSize, Discard))
		addTest(tt, mkChain(Nop, tc.srcSize, tc.t, tc.dstSize, Discard))
		addTest(tt, mkChain(Nop, tc.srcSize, tc.t, tc.dstSize, Nop, tc.dstSize, Discard))
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		tt := tc
		tt.dstSize = 100
		tt.wantStr = strings.Replace(tc.src, "0f", "", -1)
		// Chain encoders and decoders.
		if _, ok := tc.t.(rleEncode); ok && tc.wantErr == nil ***REMOVED***
			addTest(tt, mkChain(tc.t, tc.dstSize, Nop, 1000, rleDecode***REMOVED******REMOVED***))
			addTest(tt, mkChain(tc.t, tc.dstSize, Nop, tc.dstSize, rleDecode***REMOVED******REMOVED***))
			addTest(tt, mkChain(Nop, tc.srcSize, tc.t, tc.dstSize, Nop, 100, rleDecode***REMOVED******REMOVED***))
			// decoding needs larger destinations
			addTest(tt, mkChain(Nop, tc.srcSize, tc.t, tc.dstSize, rleDecode***REMOVED******REMOVED***, 100, Nop))
			addTest(tt, mkChain(Nop, tc.srcSize, tc.t, tc.dstSize, Nop, 100, rleDecode***REMOVED******REMOVED***, 100, Nop))
		***REMOVED*** else if _, ok := tc.t.(rleDecode); ok && tc.wantErr == nil ***REMOVED***
			// The internal buffer size may need to be the sum of the maximum segment
			// size of the two encoders!
			addTest(tt, mkChain(tc.t, 2*tc.dstSize, rleEncode***REMOVED******REMOVED***))
			addTest(tt, mkChain(tc.t, tc.dstSize, Nop, 101, rleEncode***REMOVED******REMOVED***))
			addTest(tt, mkChain(Nop, tc.srcSize, tc.t, tc.dstSize, Nop, 100, rleEncode***REMOVED******REMOVED***))
			addTest(tt, mkChain(Nop, tc.srcSize, tc.t, tc.dstSize, Nop, 200, rleEncode***REMOVED******REMOVED***, 100, Nop))
		***REMOVED***
	***REMOVED***
	for _, tc := range tests ***REMOVED***
		str, iter, err := doTransform(tc)
		mi := tc.wantIter != 0 && tc.wantIter != iter
		if str != tc.wantStr || err != tc.wantErr || mi ***REMOVED***
			t.Errorf("%s:\ngot  iter:%d, %q, %v\nwant iter:%d, %q, %v", tc, iter, str, err, tc.wantIter, tc.wantStr, tc.wantErr)
		***REMOVED***
		break
	***REMOVED***
***REMOVED***

func TestRemoveFunc(t *testing.T) ***REMOVED***
	filter := RemoveFunc(func(r rune) bool ***REMOVED***
		return strings.IndexRune("ab\u0300\u1234,", r) != -1
	***REMOVED***)
	tests := []testCase***REMOVED***
		***REMOVED***
			src:     ",",
			wantStr: "",
		***REMOVED***,

		***REMOVED***
			src:     "c",
			wantStr: "c",
		***REMOVED***,

		***REMOVED***
			src:     "\u2345",
			wantStr: "\u2345",
		***REMOVED***,

		***REMOVED***
			src:     "tschüß",
			wantStr: "tschüß",
		***REMOVED***,

		***REMOVED***
			src:     ",до,свидания,",
			wantStr: "досвидания",
		***REMOVED***,

		***REMOVED***
			src:     "a\xbd\xb2=\xbc ⌘",
			wantStr: "\uFFFD\uFFFD=\uFFFD ⌘",
		***REMOVED***,

		***REMOVED***
			// If we didn't replace illegal bytes with RuneError, the result
			// would be \u0300 or the code would need to be more complex.
			src:     "\xcc\u0300\x80",
			wantStr: "\uFFFD\uFFFD",
		***REMOVED***,

		***REMOVED***
			src:      "\xcc\u0300\x80",
			dstSize:  3,
			wantStr:  "\uFFFD\uFFFD",
			wantIter: 2,
		***REMOVED***,

		***REMOVED***
			// Test a long buffer greater than the internal buffer size
			src:      "hello\xcc\xcc\xccworld",
			srcSize:  13,
			wantStr:  "hello\uFFFD\uFFFD\uFFFDworld",
			wantIter: 1,
		***REMOVED***,

		***REMOVED***
			src:     "\u2345",
			dstSize: 2,
			wantStr: "",
			wantErr: ErrShortDst,
		***REMOVED***,

		***REMOVED***
			src:     "\xcc",
			dstSize: 2,
			wantStr: "",
			wantErr: ErrShortDst,
		***REMOVED***,

		***REMOVED***
			src:     "\u0300",
			dstSize: 2,
			srcSize: 1,
			wantStr: "",
			wantErr: ErrShortSrc,
		***REMOVED***,

		***REMOVED***
			t: RemoveFunc(func(r rune) bool ***REMOVED***
				return r == utf8.RuneError
			***REMOVED***),
			src:     "\xcc\u0300\x80",
			wantStr: "\u0300",
		***REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		tc.desc = tc.src
		if tc.t == nil ***REMOVED***
			tc.t = filter
		***REMOVED***
		if tc.dstSize == 0 ***REMOVED***
			tc.dstSize = 100
		***REMOVED***
		if tc.srcSize == 0 ***REMOVED***
			tc.srcSize = 100
		***REMOVED***
		str, iter, err := doTransform(tc)
		mi := tc.wantIter != 0 && tc.wantIter != iter
		if str != tc.wantStr || err != tc.wantErr || mi ***REMOVED***
			t.Errorf("%+q:\ngot  iter:%d, %+q, %v\nwant iter:%d, %+q, %v", tc.src, iter, str, err, tc.wantIter, tc.wantStr, tc.wantErr)
		***REMOVED***

		tc.src = str
		idem, _, _ := doTransform(tc)
		if str != idem ***REMOVED***
			t.Errorf("%+q: found %+q; want %+q", tc.src, idem, str)
		***REMOVED***
	***REMOVED***
***REMOVED***

func testString(t *testing.T, f func(Transformer, string) (string, int, error)) ***REMOVED***
	for _, tt := range append(testCases, chainTests()...) ***REMOVED***
		if tt.desc == "allowStutter = true" ***REMOVED***
			// We don't have control over the buffer size, so we eliminate tests
			// that depend on a specific buffer size being set.
			continue
		***REMOVED***
		if tt.wantErr == ErrShortDst || tt.wantErr == ErrShortSrc ***REMOVED***
			// The result string will be different.
			continue
		***REMOVED***
		testtext.Run(t, tt.desc, func(t *testing.T) ***REMOVED***
			got, n, err := f(tt.t, tt.src)
			if tt.wantErr != err ***REMOVED***
				t.Errorf("error: got %v; want %v", err, tt.wantErr)
			***REMOVED***
			// Check that err == nil implies that n == len(tt.src). Note that vice
			// versa isn't necessarily true.
			if err == nil && n != len(tt.src) ***REMOVED***
				t.Errorf("err == nil: got %d bytes, want %d", n, err)
			***REMOVED***
			if got != tt.wantStr ***REMOVED***
				t.Errorf("string: got %q; want %q", got, tt.wantStr)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestBytes(t *testing.T) ***REMOVED***
	testString(t, func(z Transformer, s string) (string, int, error) ***REMOVED***
		b, n, err := Bytes(z, []byte(s))
		return string(b), n, err
	***REMOVED***)
***REMOVED***

func TestAppend(t *testing.T) ***REMOVED***
	// Create a bunch of subtests for different buffer sizes.
	testCases := [][]byte***REMOVED***
		nil,
		make([]byte, 0, 0),
		make([]byte, 0, 1),
		make([]byte, 1, 1),
		make([]byte, 1, 5),
		make([]byte, 100, 100),
		make([]byte, 100, 200),
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		testString(t, func(z Transformer, s string) (string, int, error) ***REMOVED***
			b, n, err := Append(z, tc, []byte(s))
			return string(b[len(tc):]), n, err
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestString(t *testing.T) ***REMOVED***
	testtext.Run(t, "transform", func(t *testing.T) ***REMOVED*** testString(t, String) ***REMOVED***)

	// Overrun the internal destination buffer.
	for i, s := range []string***REMOVED***
		aaa[:1*initialBufSize-1],
		aaa[:1*initialBufSize+0],
		aaa[:1*initialBufSize+1],
		AAA[:1*initialBufSize-1],
		AAA[:1*initialBufSize+0],
		AAA[:1*initialBufSize+1],
		AAA[:2*initialBufSize-1],
		AAA[:2*initialBufSize+0],
		AAA[:2*initialBufSize+1],
		aaa[:1*initialBufSize-2] + "A",
		aaa[:1*initialBufSize-1] + "A",
		aaa[:1*initialBufSize+0] + "A",
		aaa[:1*initialBufSize+1] + "A",
	***REMOVED*** ***REMOVED***
		testtext.Run(t, fmt.Sprint("dst buffer test using lower/", i), func(t *testing.T) ***REMOVED***
			got, _, _ := String(lowerCaseASCII***REMOVED******REMOVED***, s)
			if want := strings.ToLower(s); got != want ***REMOVED***
				t.Errorf("got %s (%d); want %s (%d)", got, len(got), want, len(want))
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	// Overrun the internal source buffer.
	for i, s := range []string***REMOVED***
		aaa[:1*initialBufSize-1],
		aaa[:1*initialBufSize+0],
		aaa[:1*initialBufSize+1],
		aaa[:2*initialBufSize+1],
		aaa[:2*initialBufSize+0],
		aaa[:2*initialBufSize+1],
	***REMOVED*** ***REMOVED***
		testtext.Run(t, fmt.Sprint("src buffer test using rleEncode/", i), func(t *testing.T) ***REMOVED***
			got, _, _ := String(rleEncode***REMOVED******REMOVED***, s)
			if want := fmt.Sprintf("%da", len(s)); got != want ***REMOVED***
				t.Errorf("got %s (%d); want %s (%d)", got, len(got), want, len(want))
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	// Test allocations for non-changing strings.
	// Note we still need to allocate a single buffer.
	for i, s := range []string***REMOVED***
		"",
		"123456789",
		aaa[:initialBufSize-1],
		aaa[:initialBufSize+0],
		aaa[:initialBufSize+1],
		aaa[:10*initialBufSize],
	***REMOVED*** ***REMOVED***
		testtext.Run(t, fmt.Sprint("alloc/", i), func(t *testing.T) ***REMOVED***
			if n := testtext.AllocsPerRun(5, func() ***REMOVED*** String(&lowerCaseASCIILookahead***REMOVED******REMOVED***, s) ***REMOVED***); n > 1 ***REMOVED***
				t.Errorf("#allocs was %f; want 1", n)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

// TestBytesAllocation tests that buffer growth stays limited with the trickler
// transformer, which behaves oddly but within spec. In case buffer growth is
// not correctly handled, the test will either panic with a failed allocation or
// thrash. To ensure the tests terminate under the last condition, we time out
// after some sufficiently long period of time.
func TestBytesAllocation(t *testing.T) ***REMOVED***
	done := make(chan bool)
	go func() ***REMOVED***
		in := bytes.Repeat([]byte***REMOVED***'a'***REMOVED***, 1000)
		tr := trickler(make([]byte, 1))
		Bytes(&tr, in)
		done <- true
	***REMOVED***()
	select ***REMOVED***
	case <-done:
	case <-time.After(3 * time.Second):
		t.Error("time out, likely due to excessive allocation")
	***REMOVED***
***REMOVED***

// TestStringAllocation tests that buffer growth stays limited with the trickler
// transformer, which behaves oddly but within spec. In case buffer growth is
// not correctly handled, the test will either panic with a failed allocation or
// thrash. To ensure the tests terminate under the last condition, we time out
// after some sufficiently long period of time.
func TestStringAllocation(t *testing.T) ***REMOVED***
	done := make(chan bool)
	go func() ***REMOVED***
		tr := trickler(make([]byte, 1))
		String(&tr, aaa[:1000])
		done <- true
	***REMOVED***()
	select ***REMOVED***
	case <-done:
	case <-time.After(3 * time.Second):
		t.Error("time out, likely due to excessive allocation")
	***REMOVED***
***REMOVED***

func BenchmarkStringLowerEmpty(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		String(&lowerCaseASCIILookahead***REMOVED******REMOVED***, "")
	***REMOVED***
***REMOVED***

func BenchmarkStringLowerIdentical(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		String(&lowerCaseASCIILookahead***REMOVED******REMOVED***, aaa[:4096])
	***REMOVED***
***REMOVED***

func BenchmarkStringLowerChanged(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		String(&lowerCaseASCIILookahead***REMOVED******REMOVED***, AAA[:4096])
	***REMOVED***
***REMOVED***

var (
	aaa = strings.Repeat("a", 4096)
	AAA = strings.Repeat("A", 4096)
)
