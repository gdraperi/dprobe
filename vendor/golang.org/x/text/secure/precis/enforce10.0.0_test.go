// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.10

package precis

import (
	"golang.org/x/text/secure/bidirule"
)

var enforceTestCases = []struct ***REMOVED***
	name  string
	p     *Profile
	cases []testCase
***REMOVED******REMOVED***
	***REMOVED***"Basic", NewFreeform(), []testCase***REMOVED***
		***REMOVED***"e\u0301\u031f", "\u00e9\u031f", nil***REMOVED***, // normalize
	***REMOVED******REMOVED***,

	***REMOVED***"Context Rule 1", NewFreeform(), []testCase***REMOVED***
		// Rule 1: zero-width non-joiner (U+200C)
		// From RFC:
		//   False
		//   If Canonical_Combining_Class(Before(cp)) .eq.  Virama Then True;
		//   If RegExpMatch((Joining_Type:***REMOVED***L,D***REMOVED***)(Joining_Type:T)*\u200C
		//          (Joining_Type:T)*(Joining_Type:***REMOVED***R,D***REMOVED***)) Then True;
		//
		// Example runes for different joining types:
		// Join L: U+A872; PHAGS-PA SUPERFIXED LETTER RA
		// Join D: U+062C; HAH WITH DOT BELOW
		// Join T: U+0610; ARABIC SIGN SALLALLAHOU ALAYHE WASSALLAM
		// Join R: U+0627; ALEF
		// Virama: U+0A4D; GURMUKHI SIGN VIRAMA
		// Virama and Join T: U+0ACD; GUJARATI SIGN VIRAMA
		***REMOVED***"\u200c", "", errContext***REMOVED***,
		***REMOVED***"\u200ca", "", errContext***REMOVED***,
		***REMOVED***"a\u200c", "", errContext***REMOVED***,
		***REMOVED***"\u200c\u0627", "", errContext***REMOVED***,             // missing JoinStart
		***REMOVED***"\u062c\u200c", "", errContext***REMOVED***,             // missing JoinEnd
		***REMOVED***"\u0610\u200c\u0610\u0627", "", errContext***REMOVED***, // missing JoinStart
		***REMOVED***"\u062c\u0610\u200c\u0610", "", errContext***REMOVED***, // missing JoinEnd

		// Variants of: D T* U+200c T* R
		***REMOVED***"\u062c\u200c\u0627", "\u062c\u200c\u0627", nil***REMOVED***,
		***REMOVED***"\u062c\u0610\u200c\u0610\u0627", "\u062c\u0610\u200c\u0610\u0627", nil***REMOVED***,
		***REMOVED***"\u062c\u0610\u0610\u200c\u0610\u0610\u0627", "\u062c\u0610\u0610\u200c\u0610\u0610\u0627", nil***REMOVED***,
		***REMOVED***"\u062c\u0610\u200c\u0627", "\u062c\u0610\u200c\u0627", nil***REMOVED***,
		***REMOVED***"\u062c\u200c\u0610\u0627", "\u062c\u200c\u0610\u0627", nil***REMOVED***,

		// Variants of: L T* U+200c T* D
		***REMOVED***"\ua872\u200c\u062c", "\ua872\u200c\u062c", nil***REMOVED***,
		***REMOVED***"\ua872\u0610\u200c\u0610\u062c", "\ua872\u0610\u200c\u0610\u062c", nil***REMOVED***,
		***REMOVED***"\ua872\u0610\u0610\u200c\u0610\u0610\u062c", "\ua872\u0610\u0610\u200c\u0610\u0610\u062c", nil***REMOVED***,
		***REMOVED***"\ua872\u0610\u200c\u062c", "\ua872\u0610\u200c\u062c", nil***REMOVED***,
		***REMOVED***"\ua872\u200c\u0610\u062c", "\ua872\u200c\u0610\u062c", nil***REMOVED***,

		// Virama
		***REMOVED***"\u0a4d\u200c", "\u0a4d\u200c", nil***REMOVED***,
		***REMOVED***"\ua872\u0a4d\u200c", "\ua872\u0a4d\u200c", nil***REMOVED***,
		***REMOVED***"\ua872\u0a4d\u0610\u200c", "", errContext***REMOVED***,
		***REMOVED***"\ua872\u0a4d\u0610\u200c", "", errContext***REMOVED***,

		***REMOVED***"\u0acd\u200c", "\u0acd\u200c", nil***REMOVED***,
		***REMOVED***"\ua872\u0acd\u200c", "\ua872\u0acd\u200c", nil***REMOVED***,
		***REMOVED***"\ua872\u0acd\u0610\u200c", "", errContext***REMOVED***,
		***REMOVED***"\ua872\u0acd\u0610\u200c", "", errContext***REMOVED***,

		// Using Virama as join T
		***REMOVED***"\ua872\u0acd\u200c\u062c", "\ua872\u0acd\u200c\u062c", nil***REMOVED***,
		***REMOVED***"\ua872\u200c\u0acd\u062c", "\ua872\u200c\u0acd\u062c", nil***REMOVED***,
	***REMOVED******REMOVED***,

	***REMOVED***"Context Rule 2", NewFreeform(), []testCase***REMOVED***
		// Rule 2: zero-width joiner (U+200D)
		***REMOVED***"\u200d", "", errContext***REMOVED***,
		***REMOVED***"\u200da", "", errContext***REMOVED***,
		***REMOVED***"a\u200d", "", errContext***REMOVED***,

		***REMOVED***"\u0a4d\u200d", "\u0a4d\u200d", nil***REMOVED***,
		***REMOVED***"\ua872\u0a4d\u200d", "\ua872\u0a4d\u200d", nil***REMOVED***,
		***REMOVED***"\u0a4da\u200d", "", errContext***REMOVED***,
	***REMOVED******REMOVED***,

	***REMOVED***"Context Rule 3", NewFreeform(), []testCase***REMOVED***
		// Rule 3: middle dot
		***REMOVED***"·", "", errContext***REMOVED***,
		***REMOVED***"l·", "", errContext***REMOVED***,
		***REMOVED***"·l", "", errContext***REMOVED***,
		***REMOVED***"a·", "", errContext***REMOVED***,
		***REMOVED***"l·a", "", errContext***REMOVED***,
		***REMOVED***"a·a", "", errContext***REMOVED***,
		***REMOVED***"l·l", "l·l", nil***REMOVED***,
		***REMOVED***"al·la", "al·la", nil***REMOVED***,
	***REMOVED******REMOVED***,

	***REMOVED***"Context Rule 4", NewFreeform(), []testCase***REMOVED***
		// Rule 4: Greek lower numeral U+0375
		***REMOVED***"͵", "", errContext***REMOVED***,
		***REMOVED***"͵a", "", errContext***REMOVED***,
		***REMOVED***"α͵", "", errContext***REMOVED***,
		***REMOVED***"͵α", "͵α", nil***REMOVED***,
		***REMOVED***"α͵α", "α͵α", nil***REMOVED***,
		***REMOVED***"͵͵α", "͵͵α", nil***REMOVED***, // The numeric sign is itself Greek.
		***REMOVED***"α͵͵α", "α͵͵α", nil***REMOVED***,
		***REMOVED***"α͵͵", "", errContext***REMOVED***,
		***REMOVED***"α͵͵a", "", errContext***REMOVED***,
	***REMOVED******REMOVED***,

	***REMOVED***"Context Rule 5+6", NewFreeform(), []testCase***REMOVED***
		// Rule 5+6: Hebrew preceding
		// U+05f3: Geresh
		***REMOVED***"׳", "", errContext***REMOVED***,
		***REMOVED***"׳ה", "", errContext***REMOVED***,
		***REMOVED***"a׳b", "", errContext***REMOVED***,
		***REMOVED***"ש׳", "ש׳", nil***REMOVED***,     // U+05e9 U+05f3
		***REMOVED***"ש׳׳׳", "ש׳׳׳", nil***REMOVED***, // U+05e9 U+05f3

		// U+05f4: Gershayim
		***REMOVED***"״", "", errContext***REMOVED***,
		***REMOVED***"״ה", "", errContext***REMOVED***,
		***REMOVED***"a״b", "", errContext***REMOVED***,
		***REMOVED***"ש״", "ש״", nil***REMOVED***,       // U+05e9 U+05f4
		***REMOVED***"ש״״״", "ש״״״", nil***REMOVED***,   // U+05e9 U+05f4
		***REMOVED***"aש״״״", "aש״״״", nil***REMOVED***, // U+05e9 U+05f4
	***REMOVED******REMOVED***,

	***REMOVED***"Context Rule 7", NewFreeform(), []testCase***REMOVED***
		// Rule 7: Katakana middle Dot
		***REMOVED***"・", "", errContext***REMOVED***,
		***REMOVED***"abc・", "", errContext***REMOVED***,
		***REMOVED***"・def", "", errContext***REMOVED***,
		***REMOVED***"abc・def", "", errContext***REMOVED***,
		***REMOVED***"aヅc・def", "aヅc・def", nil***REMOVED***,
		***REMOVED***"abc・dぶf", "abc・dぶf", nil***REMOVED***,
		***REMOVED***"⺐bc・def", "⺐bc・def", nil***REMOVED***,
	***REMOVED******REMOVED***,

	***REMOVED***"Context Rule 8+9", NewFreeform(), []testCase***REMOVED***
		// Rule 8+9: Arabic Indic Digit
		***REMOVED***"١٢٣٤٥۶", "", errContext***REMOVED***,
		***REMOVED***"۱۲۳۴۵٦", "", errContext***REMOVED***,
		***REMOVED***"١٢٣٤٥", "١٢٣٤٥", nil***REMOVED***,
		***REMOVED***"۱۲۳۴۵", "۱۲۳۴۵", nil***REMOVED***,
	***REMOVED******REMOVED***,

	***REMOVED***"Nickname", Nickname, []testCase***REMOVED***
		***REMOVED***"  Swan  of   Avon   ", "Swan of Avon", nil***REMOVED***,
		***REMOVED***"", "", errEmptyString***REMOVED***,
		***REMOVED***" ", "", errEmptyString***REMOVED***,
		***REMOVED***"  ", "", errEmptyString***REMOVED***,
		***REMOVED***"a\u00A0a\u1680a\u2000a\u2001a\u2002a\u2003a\u2004a\u2005a\u2006a\u2007a\u2008a\u2009a\u200Aa\u202Fa\u205Fa\u3000a", "a a a a a a a a a a a a a a a a a", nil***REMOVED***,
		***REMOVED***"Foo", "Foo", nil***REMOVED***,
		***REMOVED***"foo", "foo", nil***REMOVED***,
		***REMOVED***"Foo Bar", "Foo Bar", nil***REMOVED***,
		***REMOVED***"foo bar", "foo bar", nil***REMOVED***,
		***REMOVED***"\u03A3", "\u03A3", nil***REMOVED***,
		***REMOVED***"\u03C3", "\u03C3", nil***REMOVED***,
		// Greek final sigma is left as is (do not fold!)
		***REMOVED***"\u03C2", "\u03C2", nil***REMOVED***,
		***REMOVED***"\u265A", "♚", nil***REMOVED***,
		***REMOVED***"Richard \u2163", "Richard IV", nil***REMOVED***,
		***REMOVED***"\u212B", "Å", nil***REMOVED***,
		***REMOVED***"\uFB00", "ff", nil***REMOVED***, // because of NFKC
		***REMOVED***"שa", "שa", nil***REMOVED***,     // no bidi rule
		***REMOVED***"동일조건변경허락", "동일조건변경허락", nil***REMOVED***,
	***REMOVED******REMOVED***,
	***REMOVED***"OpaqueString", OpaqueString, []testCase***REMOVED***
		***REMOVED***"  Swan  of   Avon   ", "  Swan  of   Avon   ", nil***REMOVED***,
		***REMOVED***"", "", errEmptyString***REMOVED***,
		***REMOVED***" ", " ", nil***REMOVED***,
		***REMOVED***"  ", "  ", nil***REMOVED***,
		***REMOVED***"a\u00A0a\u1680a\u2000a\u2001a\u2002a\u2003a\u2004a\u2005a\u2006a\u2007a\u2008a\u2009a\u200Aa\u202Fa\u205Fa\u3000a", "a a a a a a a a a a a a a a a a a", nil***REMOVED***,
		***REMOVED***"Foo", "Foo", nil***REMOVED***,
		***REMOVED***"foo", "foo", nil***REMOVED***,
		***REMOVED***"Foo Bar", "Foo Bar", nil***REMOVED***,
		***REMOVED***"foo bar", "foo bar", nil***REMOVED***,
		***REMOVED***"\u03C3", "\u03C3", nil***REMOVED***,
		***REMOVED***"Richard \u2163", "Richard \u2163", nil***REMOVED***,
		***REMOVED***"\u212B", "Å", nil***REMOVED***,
		***REMOVED***"Jack of \u2666s", "Jack of \u2666s", nil***REMOVED***,
		***REMOVED***"my cat is a \u0009by", "", errDisallowedRune***REMOVED***,
		***REMOVED***"שa", "שa", nil***REMOVED***, // no bidi rule
	***REMOVED******REMOVED***,
	***REMOVED***"UsernameCaseMapped", UsernameCaseMapped, []testCase***REMOVED***
		// TODO: Should this work?
		// ***REMOVED***UsernameCaseMapped, "", "", errDisallowedRune***REMOVED***,
		***REMOVED***"juliet@example.com", "juliet@example.com", nil***REMOVED***,
		***REMOVED***"fussball", "fussball", nil***REMOVED***,
		***REMOVED***"fu\u00DFball", "fu\u00DFball", nil***REMOVED***,
		***REMOVED***"\u03C0", "\u03C0", nil***REMOVED***,
		***REMOVED***"\u03A3", "\u03C3", nil***REMOVED***,
		***REMOVED***"\u03C3", "\u03C3", nil***REMOVED***,
		// Greek final sigma is left as is (do not fold!)
		***REMOVED***"\u03C2", "\u03C2", nil***REMOVED***,
		***REMOVED***"\u0049", "\u0069", nil***REMOVED***,
		***REMOVED***"\u0049", "\u0069", nil***REMOVED***,
		***REMOVED***"\u03D2", "", errDisallowedRune***REMOVED***,
		***REMOVED***"\u03B0", "\u03B0", nil***REMOVED***,
		***REMOVED***"foo bar", "", errDisallowedRune***REMOVED***,
		***REMOVED***"♚", "", bidirule.ErrInvalid***REMOVED***,
		***REMOVED***"\u007E", "~", nil***REMOVED***,
		***REMOVED***"a", "a", nil***REMOVED***,
		***REMOVED***"!", "!", nil***REMOVED***,
		***REMOVED***"²", "", bidirule.ErrInvalid***REMOVED***,
		***REMOVED***"\t", "", errDisallowedRune***REMOVED***,
		***REMOVED***"\n", "", errDisallowedRune***REMOVED***,
		***REMOVED***"\u26D6", "", bidirule.ErrInvalid***REMOVED***,
		***REMOVED***"\u26FF", "", bidirule.ErrInvalid***REMOVED***,
		***REMOVED***"\uFB00", "", errDisallowedRune***REMOVED***,
		***REMOVED***"\u1680", "", bidirule.ErrInvalid***REMOVED***,
		***REMOVED***" ", "", errDisallowedRune***REMOVED***,
		***REMOVED***"  ", "", errDisallowedRune***REMOVED***,
		***REMOVED***"\u01C5", "", errDisallowedRune***REMOVED***,
		***REMOVED***"\u16EE", "", errDisallowedRune***REMOVED***,   // Nl RUNIC ARLAUG SYMBOL
		***REMOVED***"\u0488", "", bidirule.ErrInvalid***REMOVED***, // Me COMBINING CYRILLIC HUNDRED THOUSANDS SIGN
		***REMOVED***"\u212B", "\u00e5", nil***REMOVED***,           // Angstrom sign, NFC -> U+00E5
		***REMOVED***"A\u030A", "å", nil***REMOVED***,               // A + ring
		***REMOVED***"\u00C5", "å", nil***REMOVED***,                // A with ring
		***REMOVED***"\u00E7", "ç", nil***REMOVED***,                // c cedille
		***REMOVED***"\u0063\u0327", "ç", nil***REMOVED***,          // c + cedille
		***REMOVED***"\u0158", "ř", nil***REMOVED***,
		***REMOVED***"\u0052\u030C", "ř", nil***REMOVED***,

		***REMOVED***"\u1E61", "\u1E61", nil***REMOVED***, // LATIN SMALL LETTER S WITH DOT ABOVE

		// Confusable characters ARE allowed and should NOT be mapped.
		***REMOVED***"\u0410", "\u0430", nil***REMOVED***, // CYRILLIC CAPITAL LETTER A

		// Full width should be mapped to the canonical decomposition.
		***REMOVED***"ＡＢ", "ab", nil***REMOVED***,
		***REMOVED***"שc", "", bidirule.ErrInvalid***REMOVED***, // bidi rule

	***REMOVED******REMOVED***,
	***REMOVED***"UsernameCasePreserved", UsernameCasePreserved, []testCase***REMOVED***
		***REMOVED***"ABC", "ABC", nil***REMOVED***,
		***REMOVED***"ＡＢ", "AB", nil***REMOVED***,
		***REMOVED***"שc", "", bidirule.ErrInvalid***REMOVED***, // bidi rule
		***REMOVED***"\uFB00", "", errDisallowedRune***REMOVED***,
		***REMOVED***"\u212B", "\u00c5", nil***REMOVED***,    // Angstrom sign, NFC -> U+00E5
		***REMOVED***"ẛ", "", errDisallowedRune***REMOVED***, // LATIN SMALL LETTER LONG S WITH DOT ABOVE
	***REMOVED******REMOVED***,
***REMOVED***
