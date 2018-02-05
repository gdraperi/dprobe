// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.10

package idna

import (
	"testing"
	"unicode"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/internal/ucd"
)

func TestTables(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	lookup := func(r rune) info ***REMOVED***
		v, _ := trie.lookupString(string(r))
		return info(v)
	***REMOVED***

	ucd.Parse(gen.OpenUnicodeFile("idna", "", "IdnaMappingTable.txt"), func(p *ucd.Parser) ***REMOVED***
		r := p.Rune(0)
		x := lookup(r)
		if got, want := x.category(), catFromEntry(p); got != want ***REMOVED***
			t.Errorf("%U:category: got %x; want %x", r, got, want)
		***REMOVED***

		mapped := false
		switch p.String(1) ***REMOVED***
		case "mapped", "disallowed_STD3_mapped", "deviation":
			mapped = true
		***REMOVED***
		if x.isMapped() != mapped ***REMOVED***
			t.Errorf("%U:isMapped: got %v; want %v", r, x.isMapped(), mapped)
		***REMOVED***
		if !mapped ***REMOVED***
			return
		***REMOVED***
		want := string(p.Runes(2))
		got := string(x.appendMapping(nil, string(r)))
		if got != want ***REMOVED***
			t.Errorf("%U:mapping: got %+q; want %+q", r, got, want)
		***REMOVED***

		if x.isMapped() ***REMOVED***
			return
		***REMOVED***
		wantMark := unicode.In(r, unicode.Mark)
		gotMark := x.isModifier()
		if gotMark != wantMark ***REMOVED***
			t.Errorf("IsMark(%U) = %v; want %v", r, gotMark, wantMark)
		***REMOVED***
	***REMOVED***)

	ucd.Parse(gen.OpenUCDFile("UnicodeData.txt"), func(p *ucd.Parser) ***REMOVED***
		r := p.Rune(0)
		x := lookup(r)
		got := x.isViramaModifier()

		const cccVirama = 9
		want := p.Int(ucd.CanonicalCombiningClass) == cccVirama
		if got != want ***REMOVED***
			t.Errorf("IsVirama(%U) = %v; want %v", r, got, want)
		***REMOVED***

		rtl := false
		switch p.String(ucd.BidiClass) ***REMOVED***
		case "R", "AL", "AN":
			rtl = true
		***REMOVED***
		if got := x.isBidi("A"); got != rtl && !x.isMapped() ***REMOVED***
			t.Errorf("IsBidi(%U) = %v; want %v", r, got, rtl)
		***REMOVED***
	***REMOVED***)

	ucd.Parse(gen.OpenUCDFile("extracted/DerivedJoiningType.txt"), func(p *ucd.Parser) ***REMOVED***
		r := p.Rune(0)
		x := lookup(r)
		if x.isMapped() ***REMOVED***
			return
		***REMOVED***
		got := x.joinType()
		want := joinType[p.String(1)]
		if got != want ***REMOVED***
			t.Errorf("JoinType(%U) = %x; want %x", r, got, want)
		***REMOVED***
	***REMOVED***)
***REMOVED***
