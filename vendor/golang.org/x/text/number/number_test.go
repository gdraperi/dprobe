// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package number

import (
	"strings"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestFormatter(t *testing.T) ***REMOVED***
	overrides := map[string]string***REMOVED***
		"en": "*e#######0",
		"nl": "*n#######0",
	***REMOVED***
	testCases := []struct ***REMOVED***
		desc string
		tag  string
		f    Formatter
		want string
	***REMOVED******REMOVED******REMOVED***
		desc: "decimal",
		f:    Decimal(3),
		want: "3",
	***REMOVED***, ***REMOVED***
		desc: "decimal fraction",
		f:    Decimal(0.123),
		want: "0.123",
	***REMOVED***, ***REMOVED***
		desc: "separators",
		f:    Decimal(1234.567),
		want: "1,234.567",
	***REMOVED***, ***REMOVED***
		desc: "no separators",
		f:    Decimal(1234.567, NoSeparator()),
		want: "1234.567",
	***REMOVED***, ***REMOVED***
		desc: "max integer",
		f:    Decimal(1973, MaxIntegerDigits(2)),
		want: "73",
	***REMOVED***, ***REMOVED***
		desc: "max integer overflow",
		f:    Decimal(1973, MaxIntegerDigits(1000)),
		want: "1,973",
	***REMOVED***, ***REMOVED***
		desc: "min integer",
		f:    Decimal(12, MinIntegerDigits(5)),
		want: "00,012",
	***REMOVED***, ***REMOVED***
		desc: "max fraction zero",
		f:    Decimal(0.12345, MaxFractionDigits(0)),
		want: "0",
	***REMOVED***, ***REMOVED***
		desc: "max fraction 2",
		f:    Decimal(0.12, MaxFractionDigits(2)),
		want: "0.12",
	***REMOVED***, ***REMOVED***
		desc: "min fraction 2",
		f:    Decimal(0.12, MaxFractionDigits(2)),
		want: "0.12",
	***REMOVED***, ***REMOVED***
		desc: "max fraction overflow",
		f:    Decimal(0.125, MaxFractionDigits(1e6)),
		want: "0.125",
	***REMOVED***, ***REMOVED***
		desc: "min integer overflow",
		f:    Decimal(0, MinIntegerDigits(1e6)),
		want: strings.Repeat("000,", 255/3-1) + "000",
	***REMOVED***, ***REMOVED***
		desc: "min fraction overflow",
		f:    Decimal(0, MinFractionDigits(1e6)),
		want: "0." + strings.Repeat("0", 255), // TODO: fraction separators
	***REMOVED***, ***REMOVED***
		desc: "format width",
		f:    Decimal(123, FormatWidth(10)),
		want: "       123",
	***REMOVED***, ***REMOVED***
		desc: "format width pad option before",
		f:    Decimal(123, Pad('*'), FormatWidth(10)),
		want: "*******123",
	***REMOVED***, ***REMOVED***
		desc: "format width pad option after",
		f:    Decimal(123, FormatWidth(10), Pad('*')),
		want: "*******123",
	***REMOVED***, ***REMOVED***
		desc: "format width illegal",
		f:    Decimal(123, FormatWidth(-1)),
		want: "123",
	***REMOVED***, ***REMOVED***
		desc: "increment",
		f:    Decimal(10.33, IncrementString("0.5")),
		want: "10.5",
	***REMOVED***, ***REMOVED***
		desc: "increment",
		f:    Decimal(10, IncrementString("ppp")),
		want: "10",
	***REMOVED***, ***REMOVED***
		desc: "increment and scale",
		f:    Decimal(10.33, IncrementString("0.5"), Scale(2)),
		want: "10.50",
	***REMOVED***, ***REMOVED***
		desc: "pattern overrides en",
		tag:  "en",
		f:    Decimal(101, PatternOverrides(overrides)),
		want: "eeeee101",
	***REMOVED***, ***REMOVED***
		desc: "pattern overrides nl",
		tag:  "nl",
		f:    Decimal(101, PatternOverrides(overrides)),
		want: "nnnnn101",
	***REMOVED***, ***REMOVED***
		desc: "pattern overrides de",
		tag:  "de",
		f:    Decimal(101, PatternOverrides(overrides)),
		want: "101",
	***REMOVED***, ***REMOVED***
		desc: "language selection",
		tag:  "bn",
		f:    Decimal(123456.78, Scale(2)),
		want: "১,২৩,৪৫৬.৭৮",
	***REMOVED***, ***REMOVED***
		desc: "scale",
		f:    Decimal(1234.567, Scale(2)),
		want: "1,234.57",
	***REMOVED***, ***REMOVED***
		desc: "scientific",
		f:    Scientific(3.00),
		want: "3\u202f×\u202f10⁰",
	***REMOVED***, ***REMOVED***
		desc: "scientific",
		f:    Scientific(1234),
		want: "1.234\u202f×\u202f10³",
	***REMOVED***, ***REMOVED***
		desc: "scientific",
		f:    Scientific(1234, Scale(2)),
		want: "1.23\u202f×\u202f10³",
	***REMOVED***, ***REMOVED***
		desc: "engineering",
		f:    Engineering(12345),
		want: "12.345\u202f×\u202f10³",
	***REMOVED***, ***REMOVED***
		desc: "engineering scale",
		f:    Engineering(12345, Scale(2)),
		want: "12.34\u202f×\u202f10³",
	***REMOVED***, ***REMOVED***
		desc: "engineering precision(4)",
		f:    Engineering(12345, Precision(4)),
		want: "12.34\u202f×\u202f10³",
	***REMOVED***, ***REMOVED***
		desc: "engineering precision(2)",
		f:    Engineering(1234.5, Precision(2)),
		want: "1.2\u202f×\u202f10³",
	***REMOVED***, ***REMOVED***
		desc: "percent",
		f:    Percent(0.12),
		want: "12%",
	***REMOVED***, ***REMOVED***
		desc: "permille",
		f:    PerMille(0.123),
		want: "123‰",
	***REMOVED***, ***REMOVED***
		desc: "percent rounding",
		f:    PerMille(0.12345),
		want: "123‰",
	***REMOVED***, ***REMOVED***
		desc: "percent fraction",
		f:    PerMille(0.12345, Scale(2)),
		want: "123.45‰",
	***REMOVED***, ***REMOVED***
		desc: "percent fraction",
		f:    PerMille(0.12344, Scale(1)),
		want: "123.4‰",
	***REMOVED******REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(tc.desc, func(t *testing.T) ***REMOVED***
			tag := language.Und
			if tc.tag != "" ***REMOVED***
				tag = language.MustParse(tc.tag)
			***REMOVED***
			got := message.NewPrinter(tag).Sprint(tc.f)
			if got != tc.want ***REMOVED***
				t.Errorf("got %q; want %q", got, tc.want)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
