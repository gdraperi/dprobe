// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package number

import (
	"reflect"
	"testing"
	"unsafe"
)

var testCases = []struct ***REMOVED***
	pat  string
	want *Pattern
***REMOVED******REMOVED******REMOVED***
	"#",
	&Pattern***REMOVED***
		FormatWidth: 1,
		// TODO: Should MinIntegerDigits be 1?
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"0",
	&Pattern***REMOVED***
		FormatWidth: 1,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"+0",
	&Pattern***REMOVED***
		Affix:       "\x01+\x00",
		FormatWidth: 2,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"0+",
	&Pattern***REMOVED***
		Affix:       "\x00\x01+",
		FormatWidth: 2,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"0000",
	&Pattern***REMOVED***
		FormatWidth: 4,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits: 4,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	".#",
	&Pattern***REMOVED***
		FormatWidth: 2,
		RoundingContext: RoundingContext***REMOVED***
			MaxFractionDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"#0.###",
	&Pattern***REMOVED***
		FormatWidth: 6,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MaxFractionDigits: 3,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"#0.######",
	&Pattern***REMOVED***
		FormatWidth: 9,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MaxFractionDigits: 6,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"#,0",
	&Pattern***REMOVED***
		FormatWidth:  3,
		GroupingSize: [2]uint8***REMOVED***1, 0***REMOVED***,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"#,0.00",
	&Pattern***REMOVED***
		FormatWidth:  6,
		GroupingSize: [2]uint8***REMOVED***1, 0***REMOVED***,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MinFractionDigits: 2,
			MaxFractionDigits: 2,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"#,##0.###",
	&Pattern***REMOVED***
		FormatWidth:  9,
		GroupingSize: [2]uint8***REMOVED***3, 0***REMOVED***,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MaxFractionDigits: 3,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"#,##,##0.###",
	&Pattern***REMOVED***
		FormatWidth:  12,
		GroupingSize: [2]uint8***REMOVED***3, 2***REMOVED***,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MaxFractionDigits: 3,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	// Ignore additional separators.
	"#,####,##,##0.###",
	&Pattern***REMOVED***
		FormatWidth:  17,
		GroupingSize: [2]uint8***REMOVED***3, 2***REMOVED***,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MaxFractionDigits: 3,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"#E0",
	&Pattern***REMOVED***
		FormatWidth: 3,
		RoundingContext: RoundingContext***REMOVED***
			MaxIntegerDigits:  1,
			MinExponentDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	// At least one exponent digit is required. As long as this is true, one can
	// determine that scientific rendering is needed if MinExponentDigits > 0.
	"#E#",
	nil,
***REMOVED***, ***REMOVED***
	"0E0",
	&Pattern***REMOVED***
		FormatWidth: 3,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MinExponentDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"##0.###E00",
	&Pattern***REMOVED***
		FormatWidth: 10,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MaxIntegerDigits:  3,
			MaxFractionDigits: 3,
			MinExponentDigits: 2,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"##00.0#E0",
	&Pattern***REMOVED***
		FormatWidth: 9,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  2,
			MaxIntegerDigits:  4,
			MinFractionDigits: 1,
			MaxFractionDigits: 2,
			MinExponentDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"#00.0E+0",
	&Pattern***REMOVED***
		FormatWidth: 8,
		Flags:       AlwaysExpSign,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  2,
			MaxIntegerDigits:  3,
			MinFractionDigits: 1,
			MaxFractionDigits: 1,
			MinExponentDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"0.0E++0",
	nil,
***REMOVED***, ***REMOVED***
	"#0E+",
	nil,
***REMOVED***, ***REMOVED***
	// significant digits
	"@",
	&Pattern***REMOVED***
		FormatWidth: 1,
		RoundingContext: RoundingContext***REMOVED***
			MinSignificantDigits: 1,
			MaxSignificantDigits: 1,
			MaxFractionDigits:    -1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	// significant digits
	"@@@@",
	&Pattern***REMOVED***
		FormatWidth: 4,
		RoundingContext: RoundingContext***REMOVED***
			MinSignificantDigits: 4,
			MaxSignificantDigits: 4,
			MaxFractionDigits:    -1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"@###",
	&Pattern***REMOVED***
		FormatWidth: 4,
		RoundingContext: RoundingContext***REMOVED***
			MinSignificantDigits: 1,
			MaxSignificantDigits: 4,
			MaxFractionDigits:    -1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	// Exponents in significant digits mode gets normalized.
	"@@E0",
	&Pattern***REMOVED***
		FormatWidth: 4,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MaxIntegerDigits:  1,
			MinFractionDigits: 1,
			MaxFractionDigits: 1,
			MinExponentDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"@###E00",
	&Pattern***REMOVED***
		FormatWidth: 7,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MaxIntegerDigits:  1,
			MinFractionDigits: 0,
			MaxFractionDigits: 3,
			MinExponentDigits: 2,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	// The significant digits mode does not allow fractions.
	"@###.#E0",
	nil,
***REMOVED***, ***REMOVED***
	//alternative negative pattern
	"#0.###;(#0.###)",
	&Pattern***REMOVED***
		Affix:       "\x00\x00\x01(\x01)",
		NegOffset:   2,
		FormatWidth: 6,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MaxFractionDigits: 3,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	// Rounding increment
	"1.05",
	&Pattern***REMOVED***
		FormatWidth: 4,
		RoundingContext: RoundingContext***REMOVED***
			Increment:         105,
			IncrementScale:    2,
			MinIntegerDigits:  1,
			MinFractionDigits: 2,
			MaxFractionDigits: 2,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	// Rounding increment with grouping
	"1,05",
	&Pattern***REMOVED***
		FormatWidth:  4,
		GroupingSize: [2]uint8***REMOVED***2, 0***REMOVED***,
		RoundingContext: RoundingContext***REMOVED***
			Increment:         105,
			IncrementScale:    0,
			MinIntegerDigits:  3,
			MinFractionDigits: 0,
			MaxFractionDigits: 0,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"0.0%",
	&Pattern***REMOVED***
		Affix:       "\x00\x01%",
		FormatWidth: 4,
		RoundingContext: RoundingContext***REMOVED***
			DigitShift:        2,
			MinIntegerDigits:  1,
			MinFractionDigits: 1,
			MaxFractionDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"0.0‰",
	&Pattern***REMOVED***
		Affix:       "\x00\x03‰",
		FormatWidth: 4,
		RoundingContext: RoundingContext***REMOVED***
			DigitShift:        3,
			MinIntegerDigits:  1,
			MinFractionDigits: 1,
			MaxFractionDigits: 1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"#,##0.00¤",
	&Pattern***REMOVED***
		Affix:        "\x00\x02¤",
		FormatWidth:  9,
		GroupingSize: [2]uint8***REMOVED***3, 0***REMOVED***,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits:  1,
			MinFractionDigits: 2,
			MaxFractionDigits: 2,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"#,##0.00 ¤;(#,##0.00 ¤)",
	&Pattern***REMOVED***Affix: "\x00\x04\u00a0¤\x01(\x05\u00a0¤)",
		NegOffset:    6,
		FormatWidth:  10,
		GroupingSize: [2]uint8***REMOVED***3, 0***REMOVED***,
		RoundingContext: RoundingContext***REMOVED***
			DigitShift:        0,
			MinIntegerDigits:  1,
			MinFractionDigits: 2,
			MaxFractionDigits: 2,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	// padding
	"*x#",
	&Pattern***REMOVED***
		PadRune:     'x',
		FormatWidth: 1,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	// padding
	"#*x",
	&Pattern***REMOVED***
		PadRune:     'x',
		FormatWidth: 1,
		Flags:       PadBeforeSuffix,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"*xpre#suf",
	&Pattern***REMOVED***
		Affix:       "\x03pre\x03suf",
		PadRune:     'x',
		FormatWidth: 7,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"pre*x#suf",
	&Pattern***REMOVED***
		Affix:       "\x03pre\x03suf",
		PadRune:     'x',
		FormatWidth: 7,
		Flags:       PadAfterPrefix,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"pre#*xsuf",
	&Pattern***REMOVED***
		Affix:       "\x03pre\x03suf",
		PadRune:     'x',
		FormatWidth: 7,
		Flags:       PadBeforeSuffix,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	"pre#suf*x",
	&Pattern***REMOVED***
		Affix:       "\x03pre\x03suf",
		PadRune:     'x',
		FormatWidth: 7,
		Flags:       PadAfterSuffix,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	`* #0 o''clock`,
	&Pattern***REMOVED***Affix: "\x00\x09 o\\'clock",
		FormatWidth: 10,
		PadRune:     32,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits: 0x1,
		***REMOVED***,
	***REMOVED***,
***REMOVED***, ***REMOVED***
	`'123'* #0'456'`,
	&Pattern***REMOVED***Affix: "\x05'123'\x05'456'",
		FormatWidth: 8,
		PadRune:     32,
		RoundingContext: RoundingContext***REMOVED***
			MinIntegerDigits: 0x1,
		***REMOVED***,
		Flags: PadAfterPrefix***REMOVED***,
***REMOVED***, ***REMOVED***
	// no duplicate padding
	"*xpre#suf*x", nil,
***REMOVED***, ***REMOVED***
	// no duplicate padding
	"*xpre#suf*x", nil,
***REMOVED******REMOVED***

func TestParsePattern(t *testing.T) ***REMOVED***
	for i, tc := range testCases ***REMOVED***
		t.Run(tc.pat, func(t *testing.T) ***REMOVED***
			f, err := ParsePattern(tc.pat)
			if !reflect.DeepEqual(f, tc.want) ***REMOVED***
				t.Errorf("%d:%s:\ngot %#v;\nwant %#v", i, tc.pat, f, tc.want)
			***REMOVED***
			if got, want := err != nil, tc.want == nil; got != want ***REMOVED***
				t.Errorf("%d:%s:error: got %v; want %v", i, tc.pat, err, want)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestPatternSize(t *testing.T) ***REMOVED***
	if sz := unsafe.Sizeof(Pattern***REMOVED******REMOVED***); sz > 56 ***REMOVED***
		t.Errorf("got %d; want <= 56", sz)
	***REMOVED***

***REMOVED***
