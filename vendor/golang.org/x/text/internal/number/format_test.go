// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package number

import (
	"fmt"
	"log"
	"testing"

	"golang.org/x/text/language"
)

func TestAppendDecimal(t *testing.T) ***REMOVED***
	type pairs map[string]string // alternates with decimal input and result

	testCases := []struct ***REMOVED***
		pattern string
		// We want to be able to test some forms of patterns that cannot be
		// represented as a string.
		pat *Pattern

		test pairs
	***REMOVED******REMOVED******REMOVED***
		pattern: "0",
		test: pairs***REMOVED***
			"0":    "0",
			"1":    "1",
			"-1":   "-1",
			".00":  "0",
			"10.":  "10",
			"12":   "12",
			"1.2":  "1",
			"NaN":  "NaN",
			"-Inf": "-∞",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "+0;+0",
		test: pairs***REMOVED***
			"0":    "+0",
			"1":    "+1",
			"-1":   "-1",
			".00":  "+0",
			"10.":  "+10",
			"12":   "+12",
			"1.2":  "+1",
			"NaN":  "NaN",
			"-Inf": "-∞",
			"Inf":  "+∞",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "0 +;0 +",
		test: pairs***REMOVED***
			"0":   "0 +",
			"1":   "1 +",
			"-1":  "1 -",
			".00": "0 +",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "0;0-",
		test: pairs***REMOVED***
			"-1":   "1-",
			"NaN":  "NaN",
			"-Inf": "∞-",
			"Inf":  "∞",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "0000",
		test: pairs***REMOVED***
			"0":     "0000",
			"1":     "0001",
			"12":    "0012",
			"12345": "12345",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: ".0",
		test: pairs***REMOVED***
			"0":      ".0",
			"1":      "1.0",
			"1.2":    "1.2",
			"1.2345": "1.2",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#.0",
		test: pairs***REMOVED***
			"0": ".0",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#.0#",
		test: pairs***REMOVED***
			"0": ".0",
			"1": "1.0",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "0.0#",
		test: pairs***REMOVED***
			"0": "0.0",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#0.###",
		test: pairs***REMOVED***
			"0":        "0",
			"1":        "1",
			"1.2":      "1.2",
			"1.2345":   "1.234", // rounding should have been done earlier
			"1234.5":   "1234.5",
			"1234.567": "1234.567",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#0.######",
		test: pairs***REMOVED***
			"0":           "0",
			"1234.5678":   "1234.5678",
			"0.123456789": "0.123457",
			"NaN":         "NaN",
			"Inf":         "∞",
		***REMOVED***,

		// Test separators.
	***REMOVED***, ***REMOVED***
		pattern: "#,#.00",
		test: pairs***REMOVED***
			"100": "1,0,0.00",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#,0.##",
		test: pairs***REMOVED***
			"10": "1,0",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#,0",
		test: pairs***REMOVED***
			"10": "1,0",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#,##,#.00",
		test: pairs***REMOVED***
			"1000": "1,00,0.00",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#,##0.###",
		test: pairs***REMOVED***
			"0":           "0",
			"1234.5678":   "1,234.568",
			"0.123456789": "0.123",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#,##,##0.###",
		test: pairs***REMOVED***
			"0":            "0",
			"123456789012": "1,23,45,67,89,012",
			"0.123456789":  "0.123",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "0,00,000.###",
		test: pairs***REMOVED***
			"0":            "0,00,000",
			"123456789012": "1,23,45,67,89,012",
			"12.3456789":   "0,00,012.346",
			"0.123456789":  "0,00,000.123",
		***REMOVED***,

		// Support for ill-formed patterns.
	***REMOVED***, ***REMOVED***
		pattern: "#",
		test: pairs***REMOVED***
			".00": "", // This is the behavior of fmt.
			"0":   "", // This is the behavior of fmt.
			"1":   "1",
			"10.": "10",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: ".#",
		test: pairs***REMOVED***
			"0":      "", // This is the behavior of fmt.
			"1":      "1",
			"1.2":    "1.2",
			"1.2345": "1.2",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#,#.##",
		test: pairs***REMOVED***
			"10": "1,0",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#,#",
		test: pairs***REMOVED***
			"10": "1,0",
		***REMOVED***,

		// Special patterns
	***REMOVED***, ***REMOVED***
		pattern: "#,max_int=2",
		pat: &Pattern***REMOVED***
			RoundingContext: RoundingContext***REMOVED***
				MaxIntegerDigits: 2,
			***REMOVED***,
		***REMOVED***,
		test: pairs***REMOVED***
			"2017": "17",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "0,max_int=2",
		pat: &Pattern***REMOVED***
			RoundingContext: RoundingContext***REMOVED***
				MaxIntegerDigits: 2,
				MinIntegerDigits: 1,
			***REMOVED***,
		***REMOVED***,
		test: pairs***REMOVED***
			"2000": "0",
			"2001": "1",
			"2017": "17",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "00,max_int=2",
		pat: &Pattern***REMOVED***
			RoundingContext: RoundingContext***REMOVED***
				MaxIntegerDigits: 2,
				MinIntegerDigits: 2,
			***REMOVED***,
		***REMOVED***,
		test: pairs***REMOVED***
			"2000": "00",
			"2001": "01",
			"2017": "17",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "@@@@,max_int=2",
		pat: &Pattern***REMOVED***
			RoundingContext: RoundingContext***REMOVED***
				MaxIntegerDigits:     2,
				MinSignificantDigits: 4,
			***REMOVED***,
		***REMOVED***,
		test: pairs***REMOVED***
			"2017": "17.00",
			"2000": "0.000",
			"2001": "1.000",
		***REMOVED***,

		// Significant digits
	***REMOVED***, ***REMOVED***
		pattern: "@@##",
		test: pairs***REMOVED***
			"1":     "1.0",
			"0.1":   "0.10", // leading zero does not count as significant digit
			"123":   "123",
			"1234":  "1234",
			"12345": "12340",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "@@@@",
		test: pairs***REMOVED***
			"1":     "1.000",
			".1":    "0.1000",
			".001":  "0.001000",
			"123":   "123.0",
			"1234":  "1234",
			"12345": "12340", // rounding down
			"NaN":   "NaN",
			"-Inf":  "-∞",
		***REMOVED***,

		// TODO: rounding
		// ***REMOVED***"@@@@": "23456": "23460"***REMOVED***, // rounding up
		// TODO: padding

		// Scientific and Engineering notation
	***REMOVED***, ***REMOVED***
		pattern: "#E0",
		test: pairs***REMOVED***
			"0":       "0\u202f×\u202f10⁰",
			"1":       "1\u202f×\u202f10⁰",
			"123.456": "1\u202f×\u202f10²",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "#E+0",
		test: pairs***REMOVED***
			"0":      "0\u202f×\u202f10⁺⁰",
			"1000":   "1\u202f×\u202f10⁺³",
			"1E100":  "1\u202f×\u202f10⁺¹⁰⁰",
			"1E-100": "1\u202f×\u202f10⁻¹⁰⁰",
			"NaN":    "NaN",
			"-Inf":   "-∞",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "##0E00",
		test: pairs***REMOVED***
			"100":     "100\u202f×\u202f10⁰⁰",
			"12345":   "12\u202f×\u202f10⁰³",
			"123.456": "123\u202f×\u202f10⁰⁰",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "##0.###E00",
		test: pairs***REMOVED***
			"100":      "100\u202f×\u202f10⁰⁰",
			"12345":    "12.345\u202f×\u202f10⁰³",
			"123456":   "123.456\u202f×\u202f10⁰³",
			"123.456":  "123.456\u202f×\u202f10⁰⁰",
			"123.4567": "123.457\u202f×\u202f10⁰⁰",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "##0.000E00",
		test: pairs***REMOVED***
			"100":     "100.000\u202f×\u202f10⁰⁰",
			"12345":   "12.345\u202f×\u202f10⁰³",
			"123.456": "123.456\u202f×\u202f10⁰⁰",
			"12.3456": "12.346\u202f×\u202f10⁰⁰",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "@@E0",
		test: pairs***REMOVED***
			"0":    "0.0\u202f×\u202f10⁰",
			"99":   "9.9\u202f×\u202f10¹",
			"0.99": "9.9\u202f×\u202f10⁻¹",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "@###E00",
		test: pairs***REMOVED***
			"0":     "0\u202f×\u202f10⁰⁰",
			"1":     "1\u202f×\u202f10⁰⁰",
			"11":    "1.1\u202f×\u202f10⁰¹",
			"111":   "1.11\u202f×\u202f10⁰²",
			"1111":  "1.111\u202f×\u202f10⁰³",
			"11111": "1.111\u202f×\u202f10⁰⁴",
			"0.1":   "1\u202f×\u202f10⁻⁰¹",
			"0.11":  "1.1\u202f×\u202f10⁻⁰¹",
			"0.001": "1\u202f×\u202f10⁻⁰³",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "*x##0",
		test: pairs***REMOVED***
			"0":    "xx0",
			"10":   "x10",
			"100":  "100",
			"1000": "1000",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "##0*x",
		test: pairs***REMOVED***
			"0":    "0xx",
			"10":   "10x",
			"100":  "100",
			"1000": "1000",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "* ###0.000",
		test: pairs***REMOVED***
			"0":        "   0.000",
			"123":      " 123.000",
			"123.456":  " 123.456",
			"1234.567": "1234.567",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "**0.0#######E00",
		test: pairs***REMOVED***
			"0":     "***0.0\u202f×\u202f10⁰⁰",
			"10":    "***1.0\u202f×\u202f10⁰¹",
			"11":    "***1.1\u202f×\u202f10⁰¹",
			"111":   "**1.11\u202f×\u202f10⁰²",
			"1111":  "*1.111\u202f×\u202f10⁰³",
			"11111": "1.1111\u202f×\u202f10⁰⁴",
			"11110": "*1.111\u202f×\u202f10⁰⁴",
			"11100": "**1.11\u202f×\u202f10⁰⁴",
			"11000": "***1.1\u202f×\u202f10⁰⁴",
			"10000": "***1.0\u202f×\u202f10⁰⁴",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "*xpre0suf",
		test: pairs***REMOVED***
			"0":  "pre0suf",
			"10": "pre10suf",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "*∞ pre ###0 suf",
		test: pairs***REMOVED***
			"0":    "∞∞∞ pre 0 suf",
			"10":   "∞∞ pre 10 suf",
			"100":  "∞ pre 100 suf",
			"1000": " pre 1000 suf",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "pre *∞###0 suf",
		test: pairs***REMOVED***
			"0":    "pre ∞∞∞0 suf",
			"10":   "pre ∞∞10 suf",
			"100":  "pre ∞100 suf",
			"1000": "pre 1000 suf",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "pre ###0*∞ suf",
		test: pairs***REMOVED***
			"0":    "pre 0∞∞∞ suf",
			"10":   "pre 10∞∞ suf",
			"100":  "pre 100∞ suf",
			"1000": "pre 1000 suf",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "pre ###0 suf *∞",
		test: pairs***REMOVED***
			"0":    "pre 0 suf ∞∞∞",
			"10":   "pre 10 suf ∞∞",
			"100":  "pre 100 suf ∞",
			"1000": "pre 1000 suf ",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		// Take width of positive pattern.
		pattern: "**###0;**-#####0x",
		test: pairs***REMOVED***
			"0":  "***0",
			"-1": "*-1x",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "0.00%",
		test: pairs***REMOVED***
			"0.1": "10.00%",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "0.##%",
		test: pairs***REMOVED***
			"0.1":     "10%",
			"0.11":    "11%",
			"0.111":   "11.1%",
			"0.1111":  "11.11%",
			"0.11111": "11.11%",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		pattern: "‰ 0.0#",
		test: pairs***REMOVED***
			"0.1":      "‰ 100.0",
			"0.11":     "‰ 110.0",
			"0.111":    "‰ 111.0",
			"0.1111":   "‰ 111.1",
			"0.11111":  "‰ 111.11",
			"0.111111": "‰ 111.11",
		***REMOVED***,
	***REMOVED******REMOVED***

	// TODO:
	// 	"#,##0.00¤",
	// 	"#,##0.00 ¤;(#,##0.00 ¤)",

	for _, tc := range testCases ***REMOVED***
		pat := tc.pat
		if pat == nil ***REMOVED***
			var err error
			if pat, err = ParsePattern(tc.pattern); err != nil ***REMOVED***
				log.Fatal(err)
			***REMOVED***
		***REMOVED***
		var f Formatter
		f.InitPattern(language.English, pat)
		for num, want := range tc.test ***REMOVED***
			buf := make([]byte, 100)
			t.Run(tc.pattern+"/"+num, func(t *testing.T) ***REMOVED***
				var d Decimal
				d.Convert(f.RoundingContext, dec(num))
				buf = f.Format(buf[:0], &d)
				if got := string(buf); got != want ***REMOVED***
					t.Errorf("\n got %[1]q (%[1]s)\nwant %[2]q (%[2]s)", got, want)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestLocales(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		tag  language.Tag
		num  string
		want string
	***REMOVED******REMOVED***
		***REMOVED***language.Make("en"), "123456.78", "123,456.78"***REMOVED***,
		***REMOVED***language.Make("de"), "123456.78", "123.456,78"***REMOVED***,
		***REMOVED***language.Make("de-CH"), "123456.78", "123’456.78"***REMOVED***,
		***REMOVED***language.Make("fr"), "123456.78", "123 456,78"***REMOVED***,
		***REMOVED***language.Make("bn"), "123456.78", "১,২৩,৪৫৬.৭৮"***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(fmt.Sprint(tc.tag, "/", tc.num), func(t *testing.T) ***REMOVED***
			var f Formatter
			f.InitDecimal(tc.tag)
			var d Decimal
			d.Convert(f.RoundingContext, dec(tc.num))
			b := f.Format(nil, &d)
			if got := string(b); got != tc.want ***REMOVED***
				t.Errorf("got %[1]q (%[1]s); want %[2]q (%[2]s)", got, tc.want)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestFormatters(t *testing.T) ***REMOVED***
	var f Formatter
	testCases := []struct ***REMOVED***
		init func(t language.Tag)
		num  string
		want string
	***REMOVED******REMOVED***
		***REMOVED***f.InitDecimal, "123456.78", "123,456.78"***REMOVED***,
		***REMOVED***f.InitScientific, "123456.78", "1.23\u202f×\u202f10⁵"***REMOVED***,
		***REMOVED***f.InitEngineering, "123456.78", "123.46\u202f×\u202f10³"***REMOVED***,
		***REMOVED***f.InitEngineering, "1234", "1.23\u202f×\u202f10³"***REMOVED***,

		***REMOVED***f.InitPercent, "0.1234", "12.34%"***REMOVED***,
		***REMOVED***f.InitPerMille, "0.1234", "123.40‰"***REMOVED***,
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		t.Run(fmt.Sprint(i, "/", tc.num), func(t *testing.T) ***REMOVED***
			tc.init(language.English)
			f.SetScale(2)
			var d Decimal
			d.Convert(f.RoundingContext, dec(tc.num))
			b := f.Format(nil, &d)
			if got := string(b); got != tc.want ***REMOVED***
				t.Errorf("got %[1]q (%[1]s); want %[2]q (%[2]s)", got, tc.want)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
