// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package width

import (
	"testing"

	"golang.org/x/text/internal/testtext"
)

const (
	loSurrogate = 0xD800
	hiSurrogate = 0xDFFF
)

func TestTables(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	runes := map[rune]Kind***REMOVED******REMOVED***
	getWidthData(func(r rune, tag elem, _ rune) ***REMOVED***
		runes[r] = tag.kind()
	***REMOVED***)
	for r := rune(0); r < 0x10FFFF; r++ ***REMOVED***
		if loSurrogate <= r && r <= hiSurrogate ***REMOVED***
			continue
		***REMOVED***
		p := LookupRune(r)
		if got, want := p.Kind(), runes[r]; got != want ***REMOVED***
			t.Errorf("Kind of %U was %s; want %s.", r, got, want)
		***REMOVED***
		want, mapped := foldRune(r)
		if got := p.Folded(); (got == 0) == mapped || got != 0 && got != want ***REMOVED***
			t.Errorf("Folded(%U) = %U; want %U", r, got, want)
		***REMOVED***
		want, mapped = widenRune(r)
		if got := p.Wide(); (got == 0) == mapped || got != 0 && got != want ***REMOVED***
			t.Errorf("Wide(%U) = %U; want %U", r, got, want)
		***REMOVED***
		want, mapped = narrowRune(r)
		if got := p.Narrow(); (got == 0) == mapped || got != 0 && got != want ***REMOVED***
			t.Errorf("Narrow(%U) = %U; want %U", r, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestAmbiguous verifies that that ambiguous runes with a mapping always map to
// a halfwidth rune.
func TestAmbiguous(t *testing.T) ***REMOVED***
	for r, m := range mapRunes ***REMOVED***
		if m.e != tagAmbiguous ***REMOVED***
			continue
		***REMOVED***
		if k := mapRunes[m.r].e.kind(); k != EastAsianHalfwidth ***REMOVED***
			t.Errorf("Rune %U is ambiguous and maps to a rune of type %v", r, k)
		***REMOVED***
	***REMOVED***
***REMOVED***
