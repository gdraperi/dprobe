// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runenames

import (
	"strings"
	"testing"
	"unicode"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/internal/ucd"
)

func TestName(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	wants := make([]string, 1+unicode.MaxRune)
	ucd.Parse(gen.OpenUCDFile("UnicodeData.txt"), func(p *ucd.Parser) ***REMOVED***
		r, s := p.Rune(0), p.String(ucd.Name)
		if s == "" ***REMOVED***
			return
		***REMOVED***
		if s[0] == '<' ***REMOVED***
			const first = ", First>"
			if i := strings.Index(s, first); i >= 0 ***REMOVED***
				s = s[:i] + ">"
			***REMOVED***
		***REMOVED***
		wants[r] = s
	***REMOVED***)

	nErrors := 0
	for r, want := range wants ***REMOVED***
		got := Name(rune(r))
		if got != want ***REMOVED***
			t.Errorf("r=%#08x: got %q, want %q", r, got, want)
			nErrors++
			if nErrors == 100 ***REMOVED***
				t.Fatal("too many errors")
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
