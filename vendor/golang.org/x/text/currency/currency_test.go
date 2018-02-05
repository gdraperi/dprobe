// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package currency

import (
	"fmt"
	"testing"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/language"
)

var (
	cup = MustParseISO("CUP")
	czk = MustParseISO("CZK")
	xcd = MustParseISO("XCD")
	zwr = MustParseISO("ZWR")
)

func TestParseISO(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		in  string
		out Unit
		ok  bool
	***REMOVED******REMOVED***
		***REMOVED***"USD", USD, true***REMOVED***,
		***REMOVED***"xxx", XXX, true***REMOVED***,
		***REMOVED***"xts", XTS, true***REMOVED***,
		***REMOVED***"XX", XXX, false***REMOVED***,
		***REMOVED***"XXXX", XXX, false***REMOVED***,
		***REMOVED***"", XXX, false***REMOVED***,       // not well-formed
		***REMOVED***"UUU", XXX, false***REMOVED***,    // unknown
		***REMOVED***"\u22A9", XXX, false***REMOVED***, // non-ASCII, printable

		***REMOVED***"aaa", XXX, false***REMOVED***,
		***REMOVED***"zzz", XXX, false***REMOVED***,
		***REMOVED***"000", XXX, false***REMOVED***,
		***REMOVED***"999", XXX, false***REMOVED***,
		***REMOVED***"---", XXX, false***REMOVED***,
		***REMOVED***"\x00\x00\x00", XXX, false***REMOVED***,
		***REMOVED***"\xff\xff\xff", XXX, false***REMOVED***,
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		if x, err := ParseISO(tc.in); x != tc.out || err == nil != tc.ok ***REMOVED***
			t.Errorf("%d:%s: was %s, %v; want %s, %v", i, tc.in, x, err == nil, tc.out, tc.ok)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFromRegion(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		region   string
		currency Unit
		ok       bool
	***REMOVED******REMOVED***
		***REMOVED***"NL", EUR, true***REMOVED***,
		***REMOVED***"BE", EUR, true***REMOVED***,
		***REMOVED***"AG", xcd, true***REMOVED***,
		***REMOVED***"CH", CHF, true***REMOVED***,
		***REMOVED***"CU", cup, true***REMOVED***,   // first of multiple
		***REMOVED***"DG", USD, true***REMOVED***,   // does not have M49 code
		***REMOVED***"150", XXX, false***REMOVED***, // implicit false
		***REMOVED***"CP", XXX, false***REMOVED***,  // explicit false in CLDR
		***REMOVED***"CS", XXX, false***REMOVED***,  // all expired
		***REMOVED***"ZZ", XXX, false***REMOVED***,  // none match
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		cur, ok := FromRegion(language.MustParseRegion(tc.region))
		if cur != tc.currency || ok != tc.ok ***REMOVED***
			t.Errorf("%s: got %v, %v; want %v, %v", tc.region, cur, ok, tc.currency, tc.ok)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFromTag(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		tag      string
		currency Unit
		conf     language.Confidence
	***REMOVED******REMOVED***
		***REMOVED***"nl", EUR, language.Low***REMOVED***,      // nl also spoken outside Euro land.
		***REMOVED***"nl-BE", EUR, language.Exact***REMOVED***, // region is known
		***REMOVED***"pt", BRL, language.Low***REMOVED***,
		***REMOVED***"en", USD, language.Low***REMOVED***,
		***REMOVED***"en-u-cu-eur", EUR, language.Exact***REMOVED***,
		***REMOVED***"tlh", XXX, language.No***REMOVED***, // Klingon has no country.
		***REMOVED***"es-419", XXX, language.No***REMOVED***,
		***REMOVED***"und", USD, language.Low***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		cur, conf := FromTag(language.MustParse(tc.tag))
		if cur != tc.currency || conf != tc.conf ***REMOVED***
			t.Errorf("%s: got %v, %v; want %v, %v", tc.tag, cur, conf, tc.currency, tc.conf)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTable(t *testing.T) ***REMOVED***
	for i := 4; i < len(currency); i += 4 ***REMOVED***
		if a, b := currency[i-4:i-1], currency[i:i+3]; a >= b ***REMOVED***
			t.Errorf("currency unordered at element %d: %s >= %s", i, a, b)
		***REMOVED***
	***REMOVED***
	// First currency has index 1, last is numCurrencies.
	if c := currency.Elem(1)[:3]; c != "ADP" ***REMOVED***
		t.Errorf("first was %q; want ADP", c)
	***REMOVED***
	if c := currency.Elem(numCurrencies)[:3]; c != "ZWR" ***REMOVED***
		t.Errorf("last was %q; want ZWR", c)
	***REMOVED***
***REMOVED***

func TestKindRounding(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		kind  Kind
		cur   Unit
		scale int
		inc   int
	***REMOVED******REMOVED***
		***REMOVED***Standard, USD, 2, 1***REMOVED***,
		***REMOVED***Standard, CHF, 2, 1***REMOVED***,
		***REMOVED***Cash, CHF, 2, 5***REMOVED***,
		***REMOVED***Standard, TWD, 2, 1***REMOVED***,
		***REMOVED***Cash, TWD, 0, 1***REMOVED***,
		***REMOVED***Standard, czk, 2, 1***REMOVED***,
		***REMOVED***Cash, czk, 0, 1***REMOVED***,
		***REMOVED***Standard, zwr, 2, 1***REMOVED***,
		***REMOVED***Cash, zwr, 0, 1***REMOVED***,
		***REMOVED***Standard, KRW, 0, 1***REMOVED***,
		***REMOVED***Cash, KRW, 0, 1***REMOVED***, // Cash defaults to standard.
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		if scale, inc := tc.kind.Rounding(tc.cur); scale != tc.scale && inc != tc.inc ***REMOVED***
			t.Errorf("%d: got %d, %d; want %d, %d", i, scale, inc, tc.scale, tc.inc)
		***REMOVED***
	***REMOVED***
***REMOVED***

const body = `package main
import (
	"fmt"
	"golang.org/x/text/currency"
)
func main() ***REMOVED***
	%s
***REMOVED***
`

func TestLinking(t *testing.T) ***REMOVED***
	base := getSize(t, `fmt.Print(currency.CLDRVersion)`)
	symbols := getSize(t, `fmt.Print(currency.Symbol(currency.USD))`)
	if d := symbols - base; d < 2*1024 ***REMOVED***
		t.Errorf("size(symbols)-size(base) was %d; want > 2K", d)
	***REMOVED***
***REMOVED***

func getSize(t *testing.T, main string) int ***REMOVED***
	size, err := testtext.CodeSize(fmt.Sprintf(body, main))
	if err != nil ***REMOVED***
		t.Skipf("skipping link size test; binary size could not be determined: %v", err)
	***REMOVED***
	return size
***REMOVED***

func BenchmarkString(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		USD.String()
	***REMOVED***
***REMOVED***
