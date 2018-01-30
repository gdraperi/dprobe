// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package atom

import (
	"sort"
	"testing"
)

func TestKnown(t *testing.T) ***REMOVED***
	for _, s := range testAtomList ***REMOVED***
		if atom := Lookup([]byte(s)); atom.String() != s ***REMOVED***
			t.Errorf("Lookup(%q) = %#x (%q)", s, uint32(atom), atom.String())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestHits(t *testing.T) ***REMOVED***
	for _, a := range table ***REMOVED***
		if a == 0 ***REMOVED***
			continue
		***REMOVED***
		got := Lookup([]byte(a.String()))
		if got != a ***REMOVED***
			t.Errorf("Lookup(%q) = %#x, want %#x", a.String(), uint32(got), uint32(a))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMisses(t *testing.T) ***REMOVED***
	testCases := []string***REMOVED***
		"",
		"\x00",
		"\xff",
		"A",
		"DIV",
		"Div",
		"dIV",
		"aa",
		"a\x00",
		"ab",
		"abb",
		"abbr0",
		"abbr ",
		" abbr",
		" a",
		"acceptcharset",
		"acceptCharset",
		"accept_charset",
		"h0",
		"h1h2",
		"h7",
		"onClick",
		"λ",
		// The following string has the same hash (0xa1d7fab7) as "onmouseover".
		"\x00\x00\x00\x00\x00\x50\x18\xae\x38\xd0\xb7",
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		got := Lookup([]byte(tc))
		if got != 0 ***REMOVED***
			t.Errorf("Lookup(%q): got %d, want 0", tc, got)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestForeignObject(t *testing.T) ***REMOVED***
	const (
		afo = Foreignobject
		afO = ForeignObject
		sfo = "foreignobject"
		sfO = "foreignObject"
	)
	if got := Lookup([]byte(sfo)); got != afo ***REMOVED***
		t.Errorf("Lookup(%q): got %#v, want %#v", sfo, got, afo)
	***REMOVED***
	if got := Lookup([]byte(sfO)); got != afO ***REMOVED***
		t.Errorf("Lookup(%q): got %#v, want %#v", sfO, got, afO)
	***REMOVED***
	if got := afo.String(); got != sfo ***REMOVED***
		t.Errorf("Atom(%#v).String(): got %q, want %q", afo, got, sfo)
	***REMOVED***
	if got := afO.String(); got != sfO ***REMOVED***
		t.Errorf("Atom(%#v).String(): got %q, want %q", afO, got, sfO)
	***REMOVED***
***REMOVED***

func BenchmarkLookup(b *testing.B) ***REMOVED***
	sortedTable := make([]string, 0, len(table))
	for _, a := range table ***REMOVED***
		if a != 0 ***REMOVED***
			sortedTable = append(sortedTable, a.String())
		***REMOVED***
	***REMOVED***
	sort.Strings(sortedTable)

	x := make([][]byte, 1000)
	for i := range x ***REMOVED***
		x[i] = []byte(sortedTable[i%len(sortedTable)])
	***REMOVED***

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		for _, s := range x ***REMOVED***
			Lookup(s)
		***REMOVED***
	***REMOVED***
***REMOVED***
