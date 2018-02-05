// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package precis

import (
	"fmt"
	"math/rand"
	"testing"
	"unicode"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/transform"
)

// copyOrbit is a Transformer for the sole purpose of testing the apply method,
// testing that apply will always call Span for the prefix of the input that
// remains identical and then call Transform for the remainder. It will produce
// inconsistent output for other usage patterns.
// Provided that copyOrbit is used this way, the first t bytes of the output
// will be identical to the input and the remaining output will be the result
// of calling caseOrbit on the remaining input bytes.
type copyOrbit int

func (t copyOrbit) Reset() ***REMOVED******REMOVED***
func (t copyOrbit) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	if int(t) == len(src) ***REMOVED***
		return int(t), nil
	***REMOVED***
	return int(t), transform.ErrEndOfSpan
***REMOVED***

// Transform implements transform.Transformer specifically for testing the apply method.
// See documentation of copyOrbit before using this method.
func (t copyOrbit) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	n := copy(dst, src)
	for i, c := range dst[:n] ***REMOVED***
		dst[i] = orbitCase(c)
	***REMOVED***
	return n, n, nil
***REMOVED***

func orbitCase(c byte) byte ***REMOVED***
	if unicode.IsLower(rune(c)) ***REMOVED***
		return byte(unicode.ToUpper(rune(c)))
	***REMOVED*** else ***REMOVED***
		return byte(unicode.ToLower(rune(c)))
	***REMOVED***
***REMOVED***

func TestBuffers(t *testing.T) ***REMOVED***
	want := "Those who cannot remember the past are condemned to compute it."

	spans := rand.Perm(len(want) + 1)

	// Compute the result of applying copyOrbit(span) transforms in reverse.
	input := []byte(want)
	for i := len(spans) - 1; i >= 0; i-- ***REMOVED***
		for j := spans[i]; j < len(input); j++ ***REMOVED***
			input[j] = orbitCase(input[j])
		***REMOVED***
	***REMOVED***

	// Apply the copyOrbit(span) transforms.
	b := buffers***REMOVED***src: input***REMOVED***
	for _, n := range spans ***REMOVED***
		b.apply(copyOrbit(n))
		if n%11 == 0 ***REMOVED***
			b.apply(transform.Nop)
		***REMOVED***
	***REMOVED***
	if got := string(b.src); got != want ***REMOVED***
		t.Errorf("got %q; want %q", got, want)
	***REMOVED***
***REMOVED***

type compareTestCase struct ***REMOVED***
	a      string
	b      string
	result bool
***REMOVED***

var compareTestCases = []struct ***REMOVED***
	name  string
	p     *Profile
	cases []compareTestCase
***REMOVED******REMOVED***
	***REMOVED***"Nickname", Nickname, []compareTestCase***REMOVED***
		***REMOVED***"a", "b", false***REMOVED***,
		***REMOVED***"  Swan  of   Avon   ", "swan of avon", true***REMOVED***,
		***REMOVED***"Foo", "foo", true***REMOVED***,
		***REMOVED***"foo", "foo", true***REMOVED***,
		***REMOVED***"Foo Bar", "foo bar", true***REMOVED***,
		***REMOVED***"foo bar", "foo bar", true***REMOVED***,
		***REMOVED***"\u03A3", "\u03C3", true***REMOVED***,
		***REMOVED***"\u03A3", "\u03C2", false***REMOVED***,
		***REMOVED***"\u03C3", "\u03C2", false***REMOVED***,
		***REMOVED***"Richard \u2163", "richard iv", true***REMOVED***,
		***REMOVED***"Å", "å", true***REMOVED***,
		***REMOVED***"ﬀ", "ff", true***REMOVED***, // because of NFKC
		***REMOVED***"ß", "sS", false***REMOVED***,

		// After applying the Nickname profile, \u00a8  becomes \u0020\u0308,
		// however because the nickname profile is not idempotent, applying it again
		// to \u0020\u0308 results in \u0308.
		***REMOVED***"\u00a8", "\u0020\u0308", true***REMOVED***,
		***REMOVED***"\u00a8", "\u0308", true***REMOVED***,
		***REMOVED***"\u0020\u0308", "\u0308", true***REMOVED***,
	***REMOVED******REMOVED***,
***REMOVED***

func doCompareTests(t *testing.T, fn func(t *testing.T, p *Profile, tc compareTestCase)) ***REMOVED***
	for _, g := range compareTestCases ***REMOVED***
		for i, tc := range g.cases ***REMOVED***
			name := fmt.Sprintf("%s:%d:%+q", g.name, i, tc.a)
			testtext.Run(t, name, func(t *testing.T) ***REMOVED***
				fn(t, g.p, tc)
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCompare(t *testing.T) ***REMOVED***
	doCompareTests(t, func(t *testing.T, p *Profile, tc compareTestCase) ***REMOVED***
		if result := p.Compare(tc.a, tc.b); result != tc.result ***REMOVED***
			t.Errorf("got %v; want %v", result, tc.result)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestCompareString(t *testing.T) ***REMOVED***
	doCompareTests(t, func(t *testing.T, p *Profile, tc compareTestCase) ***REMOVED***
		a, err := p.CompareKey(tc.a)
		if err != nil ***REMOVED***
			t.Errorf("Unexpected error when creating key: %v", err)
			return
		***REMOVED***
		b, err := p.CompareKey(tc.b)
		if err != nil ***REMOVED***
			t.Errorf("Unexpected error when creating key: %v", err)
			return
		***REMOVED***

		if result := (a == b); result != tc.result ***REMOVED***
			t.Errorf("got %v; want %v", result, tc.result)
		***REMOVED***
	***REMOVED***)
***REMOVED***
