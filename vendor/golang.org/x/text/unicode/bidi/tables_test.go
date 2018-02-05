// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bidi

import (
	"testing"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/internal/ucd"
)

var labels = []string***REMOVED***
	AL:  "AL",
	AN:  "AN",
	B:   "B",
	BN:  "BN",
	CS:  "CS",
	EN:  "EN",
	ES:  "ES",
	ET:  "ET",
	L:   "L",
	NSM: "NSM",
	ON:  "ON",
	R:   "R",
	S:   "S",
	WS:  "WS",

	LRO: "LRO",
	RLO: "RLO",
	LRE: "LRE",
	RLE: "RLE",
	PDF: "PDF",
	LRI: "LRI",
	RLI: "RLI",
	FSI: "FSI",
	PDI: "PDI",
***REMOVED***

func TestTables(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	ucd.Parse(gen.OpenUCDFile("BidiBrackets.txt"), func(p *ucd.Parser) ***REMOVED***
		r1 := p.Rune(0)
		want := p.Rune(1)

		e, _ := LookupRune(r1)
		if got := e.reverseBracket(r1); got != want ***REMOVED***
			t.Errorf("Reverse(%U) = %U; want %U", r1, got, want)
		***REMOVED***
	***REMOVED***)

	done := map[rune]bool***REMOVED******REMOVED***
	test := func(name string, r rune, want string) ***REMOVED***
		str := string(r)
		e, _ := LookupString(str)
		if got := labels[e.Class()]; got != want ***REMOVED***
			t.Errorf("%s:%U: got %s; want %s", name, r, got, want)
		***REMOVED***
		if e2, sz := LookupRune(r); e != e2 || sz != len(str) ***REMOVED***
			t.Errorf("LookupRune(%U) = %v, %d; want %v, %d", r, e2, e, sz, len(str))
		***REMOVED***
		if e2, sz := Lookup([]byte(str)); e != e2 || sz != len(str) ***REMOVED***
			t.Errorf("Lookup(%U) = %v, %d; want %v, %d", r, e2, e, sz, len(str))
		***REMOVED***
		done[r] = true
	***REMOVED***

	// Insert the derived BiDi properties.
	ucd.Parse(gen.OpenUCDFile("extracted/DerivedBidiClass.txt"), func(p *ucd.Parser) ***REMOVED***
		r := p.Rune(0)
		test("derived", r, p.String(1))
	***REMOVED***)
	visitDefaults(func(r rune, c Class) ***REMOVED***
		if !done[r] ***REMOVED***
			test("default", r, labels[c])
		***REMOVED***
	***REMOVED***)

***REMOVED***
