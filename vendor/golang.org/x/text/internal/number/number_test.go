// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package number

import (
	"fmt"
	"testing"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/language"
)

func TestInfo(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		lang     string
		sym      SymbolType
		wantSym  string
		wantNine rune
	***REMOVED******REMOVED***
		***REMOVED***"und", SymDecimal, ".", '9'***REMOVED***,
		***REMOVED***"de", SymGroup, ".", '9'***REMOVED***,
		***REMOVED***"de-BE", SymGroup, ".", '9'***REMOVED***,          // inherits from de (no number data in CLDR)
		***REMOVED***"de-BE-oxendict", SymGroup, ".", '9'***REMOVED***, // inherits from de (no compact index)

		// U+096F DEVANAGARI DIGIT NINE ('९')
		***REMOVED***"de-BE-u-nu-deva", SymGroup, ".", '\u096f'***REMOVED***, // miss -> latn -> de
		***REMOVED***"de-Cyrl-BE", SymGroup, ",", '9'***REMOVED***,           // inherits from root
		***REMOVED***"de-CH", SymGroup, "’", '9'***REMOVED***,                // overrides values in de
		***REMOVED***"de-CH-oxendict", SymGroup, "’", '9'***REMOVED***,       // inherits from de-CH (no compact index)
		***REMOVED***"de-CH-u-nu-deva", SymGroup, "’", '\u096f'***REMOVED***, // miss -> latn -> de-CH

		***REMOVED***"bn-u-nu-beng", SymGroup, ",", '\u09ef'***REMOVED***,
		***REMOVED***"bn-u-nu-deva", SymGroup, ",", '\u096f'***REMOVED***,
		***REMOVED***"bn-u-nu-latn", SymGroup, ",", '9'***REMOVED***,

		***REMOVED***"pa", SymExponential, "E", '9'***REMOVED***,

		// "×۱۰^" -> U+00d7 U+06f1 U+06f0^"
		// U+06F0 EXTENDED ARABIC-INDIC DIGIT ZERO
		// U+06F1 EXTENDED ARABIC-INDIC DIGIT ONE
		// U+06F9 EXTENDED ARABIC-INDIC DIGIT NINE
		***REMOVED***"pa-u-nu-arabext", SymExponential, "\u00d7\u06f1\u06f0^", '\u06f9'***REMOVED***,

		//  "གྲངས་མེད" - > U+0f42 U+0fb2 U+0f44 U+0f66 U+0f0b U+0f58 U+0f7a U+0f51
		// Examples:
		// U+0F29 TIBETAN DIGIT NINE (༩)
		***REMOVED***"dz", SymInfinity, "\u0f42\u0fb2\u0f44\u0f66\u0f0b\u0f58\u0f7a\u0f51", '\u0f29'***REMOVED***, // defaults to tibt
		***REMOVED***"dz-u-nu-latn", SymInfinity, "∞", '9'***REMOVED***,                                           // select alternative
		***REMOVED***"dz-u-nu-tibt", SymInfinity, "\u0f42\u0fb2\u0f44\u0f66\u0f0b\u0f58\u0f7a\u0f51", '\u0f29'***REMOVED***,
		***REMOVED***"en-u-nu-tibt", SymInfinity, "∞", '\u0f29'***REMOVED***,

		// algorithmic number systems fall back to ASCII if Digits is used.
		***REMOVED***"en-u-nu-hanidec", SymPlusSign, "+", '9'***REMOVED***,
		***REMOVED***"en-u-nu-roman", SymPlusSign, "+", '9'***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(fmt.Sprintf("%s:%v", tc.lang, tc.sym), func(t *testing.T) ***REMOVED***
			info := InfoFromTag(language.MustParse(tc.lang))
			if got := info.Symbol(tc.sym); got != tc.wantSym ***REMOVED***
				t.Errorf("sym: got %q; want %q", got, tc.wantSym)
			***REMOVED***
			if got := info.Digit('9'); got != tc.wantNine ***REMOVED***
				t.Errorf("Digit(9): got %+q; want %+q", got, tc.wantNine)
			***REMOVED***
			var buf [4]byte
			if got := string(buf[:info.WriteDigit(buf[:], '9')]); got != string(tc.wantNine) ***REMOVED***
				t.Errorf("WriteDigit(9): got %+q; want %+q", got, tc.wantNine)
			***REMOVED***
			if got := string(info.AppendDigit([]byte***REMOVED******REMOVED***, 9)); got != string(tc.wantNine) ***REMOVED***
				t.Errorf("AppendDigit(9): got %+q; want %+q", got, tc.wantNine)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestFormats(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		lang    string
		pattern string
		index   []byte
	***REMOVED******REMOVED***
		***REMOVED***"en", "#,##0.###", tagToDecimal***REMOVED***,
		***REMOVED***"de", "#,##0.###", tagToDecimal***REMOVED***,
		***REMOVED***"de-CH", "#,##0.###", tagToDecimal***REMOVED***,
		***REMOVED***"pa", "#,##,##0.###", tagToDecimal***REMOVED***,
		***REMOVED***"pa-Arab", "#,##0.###", tagToDecimal***REMOVED***, // Does NOT inherit from pa!
		***REMOVED***"mr", "#,##,##0.###", tagToDecimal***REMOVED***,
		***REMOVED***"mr-IN", "#,##,##0.###", tagToDecimal***REMOVED***, // Inherits from mr.
		***REMOVED***"nl", "#E0", tagToScientific***REMOVED***,
		***REMOVED***"nl-MX", "#E0", tagToScientific***REMOVED***, // Inherits through Tag.Parent.
		***REMOVED***"zgh", "#,##0 %", tagToPercent***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		testtext.Run(t, tc.lang, func(t *testing.T) ***REMOVED***
			got := formatForLang(language.MustParse(tc.lang), tc.index)
			want, _ := ParsePattern(tc.pattern)
			if *got != *want ***REMOVED***
				t.Errorf("\ngot  %#v;\nwant %#v", got, want)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
