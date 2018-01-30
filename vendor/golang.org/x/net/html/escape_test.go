// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import "testing"

type unescapeTest struct ***REMOVED***
	// A short description of the test case.
	desc string
	// The HTML text.
	html string
	// The unescaped text.
	unescaped string
***REMOVED***

var unescapeTests = []unescapeTest***REMOVED***
	// Handle no entities.
	***REMOVED***
		"copy",
		"A\ttext\nstring",
		"A\ttext\nstring",
	***REMOVED***,
	// Handle simple named entities.
	***REMOVED***
		"simple",
		"&amp; &gt; &lt;",
		"& > <",
	***REMOVED***,
	// Handle hitting the end of the string.
	***REMOVED***
		"stringEnd",
		"&amp &amp",
		"& &",
	***REMOVED***,
	// Handle entities with two codepoints.
	***REMOVED***
		"multiCodepoint",
		"text &gesl; blah",
		"text \u22db\ufe00 blah",
	***REMOVED***,
	// Handle decimal numeric entities.
	***REMOVED***
		"decimalEntity",
		"Delta = &#916; ",
		"Delta = Δ ",
	***REMOVED***,
	// Handle hexadecimal numeric entities.
	***REMOVED***
		"hexadecimalEntity",
		"Lambda = &#x3bb; = &#X3Bb ",
		"Lambda = λ = λ ",
	***REMOVED***,
	// Handle numeric early termination.
	***REMOVED***
		"numericEnds",
		"&# &#x &#128;43 &copy = &#169f = &#xa9",
		"&# &#x €43 © = ©f = ©",
	***REMOVED***,
	// Handle numeric ISO-8859-1 entity replacements.
	***REMOVED***
		"numericReplacements",
		"Footnote&#x87;",
		"Footnote‡",
	***REMOVED***,
***REMOVED***

func TestUnescape(t *testing.T) ***REMOVED***
	for _, tt := range unescapeTests ***REMOVED***
		unescaped := UnescapeString(tt.html)
		if unescaped != tt.unescaped ***REMOVED***
			t.Errorf("TestUnescape %s: want %q, got %q", tt.desc, tt.unescaped, unescaped)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUnescapeEscape(t *testing.T) ***REMOVED***
	ss := []string***REMOVED***
		``,
		`abc def`,
		`a & b`,
		`a&amp;b`,
		`a &amp b`,
		`&quot;`,
		`"`,
		`"<&>"`,
		`&quot;&lt;&amp;&gt;&quot;`,
		`3&5==1 && 0<1, "0&lt;1", a+acute=&aacute;`,
		`The special characters are: <, >, &, ' and "`,
	***REMOVED***
	for _, s := range ss ***REMOVED***
		if got := UnescapeString(EscapeString(s)); got != s ***REMOVED***
			t.Errorf("got %q want %q", got, s)
		***REMOVED***
	***REMOVED***
***REMOVED***
