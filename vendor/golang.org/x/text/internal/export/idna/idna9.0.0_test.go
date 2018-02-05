// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.10

package idna

import "testing"

// TestLabelErrors tests strings returned in case of error. All results should
// be identical to the reference implementation and can be verified at
// http://unicode.org/cldr/utility/idna.jsp. The reference implementation,
// however, seems to not display Bidi and ContextJ errors.
//
// In some cases the behavior of browsers is added as a comment. In all cases,
// whenever a resolve search returns an error here, Chrome will treat the input
// string as a search string (including those for Bidi and Context J errors),
// unless noted otherwise.
func TestLabelErrors(t *testing.T) ***REMOVED***
	encode := func(s string) string ***REMOVED*** s, _ = encode(acePrefix, s); return s ***REMOVED***
	type kind struct ***REMOVED***
		name string
		f    func(string) (string, error)
	***REMOVED***
	punyA := kind***REMOVED***"PunycodeA", punycode.ToASCII***REMOVED***
	resolve := kind***REMOVED***"ResolveA", Lookup.ToASCII***REMOVED***
	display := kind***REMOVED***"ToUnicode", Display.ToUnicode***REMOVED***
	p := New(VerifyDNSLength(true), MapForLookup(), BidiRule())
	lengthU := kind***REMOVED***"CheckLengthU", p.ToUnicode***REMOVED***
	lengthA := kind***REMOVED***"CheckLengthA", p.ToASCII***REMOVED***
	p = New(MapForLookup(), StrictDomainName(false))
	std3 := kind***REMOVED***"STD3", p.ToASCII***REMOVED***

	testCases := []struct ***REMOVED***
		kind
		input   string
		want    string
		wantErr string
	***REMOVED******REMOVED***
		***REMOVED***lengthU, "", "", "A4"***REMOVED***, // From UTS 46 conformance test.
		***REMOVED***lengthA, "", "", "A4"***REMOVED***,

		***REMOVED***lengthU, "xn--", "", "A4"***REMOVED***,
		***REMOVED***lengthU, "foo.xn--", "foo.", "A4"***REMOVED***, // TODO: is dropping xn-- correct?
		***REMOVED***lengthU, "xn--.foo", ".foo", "A4"***REMOVED***,
		***REMOVED***lengthU, "foo.xn--.bar", "foo..bar", "A4"***REMOVED***,

		***REMOVED***display, "xn--", "", ""***REMOVED***,
		***REMOVED***display, "foo.xn--", "foo.", ""***REMOVED***, // TODO: is dropping xn-- correct?
		***REMOVED***display, "xn--.foo", ".foo", ""***REMOVED***,
		***REMOVED***display, "foo.xn--.bar", "foo..bar", ""***REMOVED***,

		***REMOVED***lengthA, "a..b", "a..b", "A4"***REMOVED***,
		***REMOVED***punyA, ".b", ".b", ""***REMOVED***,
		// For backwards compatibility, the Punycode profile does not map runes.
		***REMOVED***punyA, "\u3002b", "xn--b-83t", ""***REMOVED***,
		***REMOVED***punyA, "..b", "..b", ""***REMOVED***,
		// Only strip leading empty labels for certain profiles. Stripping
		// leading empty labels here but not for "empty" punycode above seems
		// inconsistent, but seems to be applied by both the conformance test
		// and Chrome. So we turn it off by default, support it as an option,
		// and enable it in profiles where it seems commonplace.
		***REMOVED***lengthA, ".b", "b", ""***REMOVED***,
		***REMOVED***lengthA, "\u3002b", "b", ""***REMOVED***,
		***REMOVED***lengthA, "..b", "b", ""***REMOVED***,
		***REMOVED***lengthA, "b..", "b..", ""***REMOVED***,

		***REMOVED***resolve, "a..b", "a..b", ""***REMOVED***,
		***REMOVED***resolve, ".b", "b", ""***REMOVED***,
		***REMOVED***resolve, "\u3002b", "b", ""***REMOVED***,
		***REMOVED***resolve, "..b", "b", ""***REMOVED***,
		***REMOVED***resolve, "b..", "b..", ""***REMOVED***,

		// Raw punycode
		***REMOVED***punyA, "", "", ""***REMOVED***,
		***REMOVED***punyA, "*.foo.com", "*.foo.com", ""***REMOVED***,
		***REMOVED***punyA, "Foo.com", "Foo.com", ""***REMOVED***,

		// STD3 rules
		***REMOVED***display, "*.foo.com", "*.foo.com", "P1"***REMOVED***,
		***REMOVED***std3, "*.foo.com", "*.foo.com", ""***REMOVED***,

		// Don't map U+2490 (DIGIT NINE FULL STOP). This is the behavior of
		// Chrome, Safari, and IE. Firefox will first map ⒐ to 9. and return
		// lab9.be.
		***REMOVED***resolve, "lab⒐be", "xn--labbe-zh9b", "P1"***REMOVED***, // encode("lab⒐be")
		***REMOVED***display, "lab⒐be", "lab⒐be", "P1"***REMOVED***,

		***REMOVED***resolve, "plan⒐faß.de", "xn--planfass-c31e.de", "P1"***REMOVED***, // encode("plan⒐fass") + ".de"
		***REMOVED***display, "Plan⒐faß.de", "plan⒐faß.de", "P1"***REMOVED***,

		// Chrome 54.0 recognizes the error and treats this input verbatim as a
		// search string.
		// Safari 10.0 (non-conform spec) decomposes "⒈" and computes the
		// punycode on the result using transitional mapping.
		// Firefox 49.0.1 goes haywire on this string and prints a bunch of what
		// seems to be nested punycode encodings.
		***REMOVED***resolve, "日本⒈co.ßßß.de", "xn--co-wuw5954azlb.ssssss.de", "P1"***REMOVED***,
		***REMOVED***display, "日本⒈co.ßßß.de", "日本⒈co.ßßß.de", "P1"***REMOVED***,

		***REMOVED***resolve, "a\u200Cb", "ab", ""***REMOVED***,
		***REMOVED***display, "a\u200Cb", "a\u200Cb", "C"***REMOVED***,

		***REMOVED***resolve, encode("a\u200Cb"), encode("a\u200Cb"), "C"***REMOVED***,
		***REMOVED***display, "a\u200Cb", "a\u200Cb", "C"***REMOVED***,

		***REMOVED***resolve, "grﻋﺮﺑﻲ.de", "xn--gr-gtd9a1b0g.de", "B"***REMOVED***,
		***REMOVED***
			// Notice how the string gets transformed, even with an error.
			// Chrome will use the original string if it finds an error, so not
			// the transformed one.
			display,
			"gr\ufecb\ufeae\ufe91\ufef2.de",
			"gr\u0639\u0631\u0628\u064a.de",
			"B",
		***REMOVED***,

		***REMOVED***resolve, "\u0671.\u03c3\u07dc", "xn--qib.xn--4xa21s", "B"***REMOVED***, // ٱ.σߜ
		***REMOVED***display, "\u0671.\u03c3\u07dc", "\u0671.\u03c3\u07dc", "B"***REMOVED***,

		// normalize input
		***REMOVED***resolve, "a\u0323\u0322", "xn--jta191l", ""***REMOVED***, // ạ̢
		***REMOVED***display, "a\u0323\u0322", "\u1ea1\u0322", ""***REMOVED***,

		// Non-normalized strings are not normalized when they originate from
		// punycode. Despite the error, Chrome, Safari and Firefox will attempt
		// to look up the input punycode.
		***REMOVED***resolve, encode("a\u0323\u0322") + ".com", "xn--a-tdbc.com", "V1"***REMOVED***,
		***REMOVED***display, encode("a\u0323\u0322") + ".com", "a\u0323\u0322.com", "V1"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		doTest(t, tc.f, tc.name, tc.input, tc.want, tc.wantErr)
	***REMOVED***
***REMOVED***
