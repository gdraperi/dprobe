// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package idna

import (
	"testing"
)

var idnaTestCases = [...]struct ***REMOVED***
	ascii, unicode string
***REMOVED******REMOVED***
	// Labels.
	***REMOVED***"books", "books"***REMOVED***,
	***REMOVED***"xn--bcher-kva", "bücher"***REMOVED***,

	// Domains.
	***REMOVED***"foo--xn--bar.org", "foo--xn--bar.org"***REMOVED***,
	***REMOVED***"golang.org", "golang.org"***REMOVED***,
	***REMOVED***"example.xn--p1ai", "example.рф"***REMOVED***,
	***REMOVED***"xn--czrw28b.tw", "商業.tw"***REMOVED***,
	***REMOVED***"www.xn--mller-kva.de", "www.müller.de"***REMOVED***,
***REMOVED***

func TestIDNA(t *testing.T) ***REMOVED***
	for _, tc := range idnaTestCases ***REMOVED***
		if a, err := ToASCII(tc.unicode); err != nil ***REMOVED***
			t.Errorf("ToASCII(%q): %v", tc.unicode, err)
		***REMOVED*** else if a != tc.ascii ***REMOVED***
			t.Errorf("ToASCII(%q): got %q, want %q", tc.unicode, a, tc.ascii)
		***REMOVED***

		if u, err := ToUnicode(tc.ascii); err != nil ***REMOVED***
			t.Errorf("ToUnicode(%q): %v", tc.ascii, err)
		***REMOVED*** else if u != tc.unicode ***REMOVED***
			t.Errorf("ToUnicode(%q): got %q, want %q", tc.ascii, u, tc.unicode)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIDNASeparators(t *testing.T) ***REMOVED***
	type subCase struct ***REMOVED***
		unicode   string
		wantASCII string
		wantErr   bool
	***REMOVED***

	testCases := []struct ***REMOVED***
		name     string
		profile  *Profile
		subCases []subCase
	***REMOVED******REMOVED***
		***REMOVED***
			name: "Punycode", profile: Punycode,
			subCases: []subCase***REMOVED***
				***REMOVED***"example\u3002jp", "xn--examplejp-ck3h", false***REMOVED***,
				***REMOVED***"東京\uFF0Ejp", "xn--jp-l92cn98g071o", false***REMOVED***,
				***REMOVED***"大阪\uFF61jp", "xn--jp-ku9cz72u463f", false***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "Lookup", profile: Lookup,
			subCases: []subCase***REMOVED***
				***REMOVED***"example\u3002jp", "example.jp", false***REMOVED***,
				***REMOVED***"東京\uFF0Ejp", "xn--1lqs71d.jp", false***REMOVED***,
				***REMOVED***"大阪\uFF61jp", "xn--pssu33l.jp", false***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "Display", profile: Display,
			subCases: []subCase***REMOVED***
				***REMOVED***"example\u3002jp", "example.jp", false***REMOVED***,
				***REMOVED***"東京\uFF0Ejp", "xn--1lqs71d.jp", false***REMOVED***,
				***REMOVED***"大阪\uFF61jp", "xn--pssu33l.jp", false***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "Registration", profile: Registration,
			subCases: []subCase***REMOVED***
				***REMOVED***"example\u3002jp", "", true***REMOVED***,
				***REMOVED***"東京\uFF0Ejp", "", true***REMOVED***,
				***REMOVED***"大阪\uFF61jp", "", true***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(tc.name, func(t *testing.T) ***REMOVED***
			for _, c := range tc.subCases ***REMOVED***
				gotA, err := tc.profile.ToASCII(c.unicode)
				if c.wantErr ***REMOVED***
					if err == nil ***REMOVED***
						t.Errorf("ToASCII(%q): got no error, but an error expected", c.unicode)
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if err != nil ***REMOVED***
						t.Errorf("ToASCII(%q): got err=%v, but no error expected", c.unicode, err)
					***REMOVED*** else if gotA != c.wantASCII ***REMOVED***
						t.Errorf("ToASCII(%q): got %q, want %q", c.unicode, gotA, c.wantASCII)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

// TODO(nigeltao): test errors, once we've specified when ToASCII and ToUnicode
// return errors.
