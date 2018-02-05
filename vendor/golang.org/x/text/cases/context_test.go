// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cases

import (
	"strings"
	"testing"
	"unicode"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/language"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/rangetable"
)

// The following definitions are taken directly from Chapter 3 of The Unicode
// Standard.

func propCased(r rune) bool ***REMOVED***
	return propLower(r) || propUpper(r) || unicode.IsTitle(r)
***REMOVED***

func propLower(r rune) bool ***REMOVED***
	return unicode.IsLower(r) || unicode.Is(unicode.Other_Lowercase, r)
***REMOVED***

func propUpper(r rune) bool ***REMOVED***
	return unicode.IsUpper(r) || unicode.Is(unicode.Other_Uppercase, r)
***REMOVED***

func propIgnore(r rune) bool ***REMOVED***
	if unicode.In(r, unicode.Mn, unicode.Me, unicode.Cf, unicode.Lm, unicode.Sk) ***REMOVED***
		return true
	***REMOVED***
	return caseIgnorable[r]
***REMOVED***

func hasBreakProp(r rune) bool ***REMOVED***
	// binary search over ranges
	lo := 0
	hi := len(breakProp)
	for lo < hi ***REMOVED***
		m := lo + (hi-lo)/2
		bp := &breakProp[m]
		if bp.lo <= r && r <= bp.hi ***REMOVED***
			return true
		***REMOVED***
		if r < bp.lo ***REMOVED***
			hi = m
		***REMOVED*** else ***REMOVED***
			lo = m + 1
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func contextFromRune(r rune) *context ***REMOVED***
	c := context***REMOVED***dst: make([]byte, 128), src: []byte(string(r)), atEOF: true***REMOVED***
	c.next()
	return &c
***REMOVED***

func TestCaseProperties(t *testing.T) ***REMOVED***
	if unicode.Version != UnicodeVersion ***REMOVED***
		// Properties of existing code points may change by Unicode version, so
		// we need to skip.
		t.Skipf("Skipping as core Unicode version %s different than %s", unicode.Version, UnicodeVersion)
	***REMOVED***
	assigned := rangetable.Assigned(UnicodeVersion)
	coreVersion := rangetable.Assigned(unicode.Version)
	for r := rune(0); r <= lastRuneForTesting; r++ ***REMOVED***
		if !unicode.In(r, assigned) || !unicode.In(r, coreVersion) ***REMOVED***
			continue
		***REMOVED***
		c := contextFromRune(r)
		if got, want := c.info.isCaseIgnorable(), propIgnore(r); got != want ***REMOVED***
			t.Errorf("caseIgnorable(%U): got %v; want %v (%x)", r, got, want, c.info)
		***REMOVED***
		// New letters may change case types, but existing case pairings should
		// not change. See Case Pair Stability in
		// http://unicode.org/policies/stability_policy.html.
		if rf := unicode.SimpleFold(r); rf != r && unicode.In(rf, assigned) ***REMOVED***
			if got, want := c.info.isCased(), propCased(r); got != want ***REMOVED***
				t.Errorf("cased(%U): got %v; want %v (%x)", r, got, want, c.info)
			***REMOVED***
			if got, want := c.caseType() == cUpper, propUpper(r); got != want ***REMOVED***
				t.Errorf("upper(%U): got %v; want %v (%x)", r, got, want, c.info)
			***REMOVED***
			if got, want := c.caseType() == cLower, propLower(r); got != want ***REMOVED***
				t.Errorf("lower(%U): got %v; want %v (%x)", r, got, want, c.info)
			***REMOVED***
		***REMOVED***
		if got, want := c.info.isBreak(), hasBreakProp(r); got != want ***REMOVED***
			t.Errorf("isBreak(%U): got %v; want %v (%x)", r, got, want, c.info)
		***REMOVED***
	***REMOVED***
	// TODO: get title case from unicode file.
***REMOVED***

func TestMapping(t *testing.T) ***REMOVED***
	assigned := rangetable.Assigned(UnicodeVersion)
	coreVersion := rangetable.Assigned(unicode.Version)
	if coreVersion == nil ***REMOVED***
		coreVersion = assigned
	***REMOVED***
	apply := func(r rune, f func(c *context) bool) string ***REMOVED***
		c := contextFromRune(r)
		f(c)
		return string(c.dst[:c.pDst])
	***REMOVED***

	for r, tt := range special ***REMOVED***
		if got, want := apply(r, lower), tt.toLower; got != want ***REMOVED***
			t.Errorf("lowerSpecial:(%U): got %+q; want %+q", r, got, want)
		***REMOVED***
		if got, want := apply(r, title), tt.toTitle; got != want ***REMOVED***
			t.Errorf("titleSpecial:(%U): got %+q; want %+q", r, got, want)
		***REMOVED***
		if got, want := apply(r, upper), tt.toUpper; got != want ***REMOVED***
			t.Errorf("upperSpecial:(%U): got %+q; want %+q", r, got, want)
		***REMOVED***
	***REMOVED***

	for r := rune(0); r <= lastRuneForTesting; r++ ***REMOVED***
		if !unicode.In(r, assigned) || !unicode.In(r, coreVersion) ***REMOVED***
			continue
		***REMOVED***
		if rf := unicode.SimpleFold(r); rf == r || !unicode.In(rf, assigned) ***REMOVED***
			continue
		***REMOVED***
		if _, ok := special[r]; ok ***REMOVED***
			continue
		***REMOVED***
		want := string(unicode.ToLower(r))
		if got := apply(r, lower); got != want ***REMOVED***
			t.Errorf("lower:%q (%U): got %q %U; want %q %U", r, r, got, []rune(got), want, []rune(want))
		***REMOVED***

		want = string(unicode.ToUpper(r))
		if got := apply(r, upper); got != want ***REMOVED***
			t.Errorf("upper:%q (%U): got %q %U; want %q %U", r, r, got, []rune(got), want, []rune(want))
		***REMOVED***

		want = string(unicode.ToTitle(r))
		if got := apply(r, title); got != want ***REMOVED***
			t.Errorf("title:%q (%U): got %q %U; want %q %U", r, r, got, []rune(got), want, []rune(want))
		***REMOVED***
	***REMOVED***
***REMOVED***

func runeFoldData(r rune) (x struct***REMOVED*** simple, full, special string ***REMOVED***) ***REMOVED***
	x = foldMap[r]
	if x.simple == "" ***REMOVED***
		x.simple = string(unicode.ToLower(r))
	***REMOVED***
	if x.full == "" ***REMOVED***
		x.full = string(unicode.ToLower(r))
	***REMOVED***
	if x.special == "" ***REMOVED***
		x.special = x.full
	***REMOVED***
	return
***REMOVED***

func TestFoldData(t *testing.T) ***REMOVED***
	assigned := rangetable.Assigned(UnicodeVersion)
	coreVersion := rangetable.Assigned(unicode.Version)
	if coreVersion == nil ***REMOVED***
		coreVersion = assigned
	***REMOVED***
	apply := func(r rune, f func(c *context) bool) (string, info) ***REMOVED***
		c := contextFromRune(r)
		f(c)
		return string(c.dst[:c.pDst]), c.info.cccType()
	***REMOVED***
	for r := rune(0); r <= lastRuneForTesting; r++ ***REMOVED***
		if !unicode.In(r, assigned) || !unicode.In(r, coreVersion) ***REMOVED***
			continue
		***REMOVED***
		x := runeFoldData(r)
		if got, info := apply(r, foldFull); got != x.full ***REMOVED***
			t.Errorf("full:%q (%U): got %q %U; want %q %U (ccc=%x)", r, r, got, []rune(got), x.full, []rune(x.full), info)
		***REMOVED***
		// TODO: special and simple.
	***REMOVED***
***REMOVED***

func TestCCC(t *testing.T) ***REMOVED***
	assigned := rangetable.Assigned(UnicodeVersion)
	normVersion := rangetable.Assigned(norm.Version)
	for r := rune(0); r <= lastRuneForTesting; r++ ***REMOVED***
		if !unicode.In(r, assigned) || !unicode.In(r, normVersion) ***REMOVED***
			continue
		***REMOVED***
		c := contextFromRune(r)

		p := norm.NFC.PropertiesString(string(r))
		want := cccOther
		switch p.CCC() ***REMOVED***
		case 0:
			want = cccZero
		case above:
			want = cccAbove
		***REMOVED***
		if got := c.info.cccType(); got != want ***REMOVED***
			t.Errorf("%U: got %x; want %x", r, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWordBreaks(t *testing.T) ***REMOVED***
	for _, tt := range breakTest ***REMOVED***
		testtext.Run(t, tt, func(t *testing.T) ***REMOVED***
			parts := strings.Split(tt, "|")
			want := ""
			for _, s := range parts ***REMOVED***
				found := false
				// This algorithm implements title casing given word breaks
				// as defined in the Unicode standard 3.13 R3.
				for _, r := range s ***REMOVED***
					title := unicode.ToTitle(r)
					lower := unicode.ToLower(r)
					if !found && title != lower ***REMOVED***
						found = true
						want += string(title)
					***REMOVED*** else ***REMOVED***
						want += string(lower)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			src := strings.Join(parts, "")
			got := Title(language.Und).String(src)
			if got != want ***REMOVED***
				t.Errorf("got %q; want %q", got, want)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestContext(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		desc       string
		dstSize    int
		atEOF      bool
		src        string
		out        string
		nSrc       int
		err        error
		ops        string
		prefixArg  string
		prefixWant bool
	***REMOVED******REMOVED******REMOVED***
		desc:    "next: past end, atEOF, no checkpoint",
		dstSize: 10,
		atEOF:   true,
		src:     "12",
		out:     "",
		nSrc:    2,
		ops:     "next;next;next",
		// Test that calling prefix with a non-empty argument when the buffer
		// is depleted returns false.
		prefixArg:  "x",
		prefixWant: false,
	***REMOVED***, ***REMOVED***
		desc:       "next: not at end, atEOF, no checkpoint",
		dstSize:    10,
		atEOF:      false,
		src:        "12",
		out:        "",
		nSrc:       0,
		err:        transform.ErrShortSrc,
		ops:        "next;next",
		prefixArg:  "",
		prefixWant: true,
	***REMOVED***, ***REMOVED***
		desc:       "next: past end, !atEOF, no checkpoint",
		dstSize:    10,
		atEOF:      false,
		src:        "12",
		out:        "",
		nSrc:       0,
		err:        transform.ErrShortSrc,
		ops:        "next;next;next",
		prefixArg:  "",
		prefixWant: true,
	***REMOVED***, ***REMOVED***
		desc:       "next: past end, !atEOF, checkpoint",
		dstSize:    10,
		atEOF:      false,
		src:        "12",
		out:        "",
		nSrc:       2,
		ops:        "next;next;checkpoint;next",
		prefixArg:  "",
		prefixWant: true,
	***REMOVED***, ***REMOVED***
		desc:       "copy: exact count, atEOF, no checkpoint",
		dstSize:    2,
		atEOF:      true,
		src:        "12",
		out:        "12",
		nSrc:       2,
		ops:        "next;copy;next;copy;next",
		prefixArg:  "",
		prefixWant: true,
	***REMOVED***, ***REMOVED***
		desc:       "copy: past end, !atEOF, no checkpoint",
		dstSize:    2,
		atEOF:      false,
		src:        "12",
		out:        "",
		nSrc:       0,
		err:        transform.ErrShortSrc,
		ops:        "next;copy;next;copy;next",
		prefixArg:  "",
		prefixWant: true,
	***REMOVED***, ***REMOVED***
		desc:       "copy: past end, !atEOF, checkpoint",
		dstSize:    2,
		atEOF:      false,
		src:        "12",
		out:        "12",
		nSrc:       2,
		ops:        "next;copy;next;copy;checkpoint;next",
		prefixArg:  "",
		prefixWant: true,
	***REMOVED***, ***REMOVED***
		desc:       "copy: short dst",
		dstSize:    1,
		atEOF:      false,
		src:        "12",
		out:        "",
		nSrc:       0,
		err:        transform.ErrShortDst,
		ops:        "next;copy;next;copy;checkpoint;next",
		prefixArg:  "12",
		prefixWant: false,
	***REMOVED***, ***REMOVED***
		desc:       "copy: short dst, checkpointed",
		dstSize:    1,
		atEOF:      false,
		src:        "12",
		out:        "1",
		nSrc:       1,
		err:        transform.ErrShortDst,
		ops:        "next;copy;checkpoint;next;copy;next",
		prefixArg:  "",
		prefixWant: true,
	***REMOVED***, ***REMOVED***
		desc:       "writeString: simple",
		dstSize:    3,
		atEOF:      true,
		src:        "1",
		out:        "1ab",
		nSrc:       1,
		ops:        "next;copy;writeab;next",
		prefixArg:  "",
		prefixWant: true,
	***REMOVED***, ***REMOVED***
		desc:       "writeString: short dst",
		dstSize:    2,
		atEOF:      true,
		src:        "12",
		out:        "",
		nSrc:       0,
		err:        transform.ErrShortDst,
		ops:        "next;copy;writeab;next",
		prefixArg:  "2",
		prefixWant: true,
	***REMOVED***, ***REMOVED***
		desc:       "writeString: simple",
		dstSize:    3,
		atEOF:      true,
		src:        "12",
		out:        "1ab",
		nSrc:       2,
		ops:        "next;copy;next;writeab;next",
		prefixArg:  "",
		prefixWant: true,
	***REMOVED***, ***REMOVED***
		desc:       "writeString: short dst",
		dstSize:    2,
		atEOF:      true,
		src:        "12",
		out:        "",
		nSrc:       0,
		err:        transform.ErrShortDst,
		ops:        "next;copy;next;writeab;next",
		prefixArg:  "1",
		prefixWant: false,
	***REMOVED***, ***REMOVED***
		desc:    "prefix",
		dstSize: 2,
		atEOF:   true,
		src:     "12",
		out:     "",
		nSrc:    0,
		// Context will assign an ErrShortSrc if the input wasn't exhausted.
		err:        transform.ErrShortSrc,
		prefixArg:  "12",
		prefixWant: true,
	***REMOVED******REMOVED***
	for _, tt := range tests ***REMOVED***
		c := context***REMOVED***dst: make([]byte, tt.dstSize), src: []byte(tt.src), atEOF: tt.atEOF***REMOVED***

		for _, op := range strings.Split(tt.ops, ";") ***REMOVED***
			switch op ***REMOVED***
			case "next":
				c.next()
			case "checkpoint":
				c.checkpoint()
			case "writeab":
				c.writeString("ab")
			case "copy":
				c.copy()
			case "":
			default:
				t.Fatalf("unknown op %q", op)
			***REMOVED***
		***REMOVED***
		if got := c.hasPrefix(tt.prefixArg); got != tt.prefixWant ***REMOVED***
			t.Errorf("%s:\nprefix was %v; want %v", tt.desc, got, tt.prefixWant)
		***REMOVED***
		nDst, nSrc, err := c.ret()
		if err != tt.err ***REMOVED***
			t.Errorf("%s:\nerror was %v; want %v", tt.desc, err, tt.err)
		***REMOVED***
		if out := string(c.dst[:nDst]); out != tt.out ***REMOVED***
			t.Errorf("%s:\nout was %q; want %q", tt.desc, out, tt.out)
		***REMOVED***
		if nSrc != tt.nSrc ***REMOVED***
			t.Errorf("%s:\nnSrc was %d; want %d", tt.desc, nSrc, tt.nSrc)
		***REMOVED***
	***REMOVED***
***REMOVED***
