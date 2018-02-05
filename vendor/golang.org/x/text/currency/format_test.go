// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package currency

import (
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	en    = language.English
	fr    = language.French
	en_US = language.AmericanEnglish
	en_GB = language.BritishEnglish
	en_AU = language.MustParse("en-AU")
	und   = language.Und
)

func TestFormatting(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		tag    language.Tag
		value  interface***REMOVED******REMOVED***
		format Formatter
		want   string
	***REMOVED******REMOVED***
		0: ***REMOVED***en, USD.Amount(0.1), nil, "USD 0.10"***REMOVED***,
		1: ***REMOVED***en, XPT.Amount(1.0), Symbol, "XPT 1.00"***REMOVED***,

		2: ***REMOVED***en, USD.Amount(2.0), ISO, "USD 2.00"***REMOVED***,
		3: ***REMOVED***und, USD.Amount(3.0), Symbol, "US$ 3.00"***REMOVED***,
		4: ***REMOVED***en, USD.Amount(4.0), Symbol, "$ 4.00"***REMOVED***,

		5: ***REMOVED***en, USD.Amount(5.20), NarrowSymbol, "$ 5.20"***REMOVED***,
		6: ***REMOVED***en, AUD.Amount(6.20), Symbol, "A$ 6.20"***REMOVED***,

		7: ***REMOVED***en_AU, AUD.Amount(7.20), Symbol, "$ 7.20"***REMOVED***,
		8: ***REMOVED***en_GB, USD.Amount(8.20), Symbol, "US$ 8.20"***REMOVED***,

		9:  ***REMOVED***en, 9.0, Symbol.Default(EUR), "€ 9.00"***REMOVED***,
		10: ***REMOVED***en, 10.123, Symbol.Default(KRW), "₩ 10"***REMOVED***,
		11: ***REMOVED***fr, 11.52, Symbol.Default(TWD), "TWD 11.52"***REMOVED***,
		12: ***REMOVED***en, 12.123, Symbol.Default(czk), "CZK 12.12"***REMOVED***,
		13: ***REMOVED***en, 13.123, Symbol.Default(czk).Kind(Cash), "CZK 13"***REMOVED***,
		14: ***REMOVED***en, 14.12345, ISO.Default(MustParseISO("CLF")), "CLF 14.1235"***REMOVED***,
		15: ***REMOVED***en, USD.Amount(15.00), ISO.Default(TWD), "USD 15.00"***REMOVED***,
		16: ***REMOVED***en, KRW.Amount(16.00), ISO.Kind(Cash), "KRW 16"***REMOVED***,

		// TODO: support integers as well.

		17: ***REMOVED***en, USD, nil, "USD"***REMOVED***,
		18: ***REMOVED***en, USD, ISO, "USD"***REMOVED***,
		19: ***REMOVED***en, USD, Symbol, "$"***REMOVED***,
		20: ***REMOVED***en_GB, USD, Symbol, "US$"***REMOVED***,
		21: ***REMOVED***en_AU, USD, NarrowSymbol, "$"***REMOVED***,
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		p := message.NewPrinter(tc.tag)
		v := tc.value
		if tc.format != nil ***REMOVED***
			v = tc.format(v)
		***REMOVED***
		if got := p.Sprint(v); got != tc.want ***REMOVED***
			t.Errorf("%d: got %q; want %q", i, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***
