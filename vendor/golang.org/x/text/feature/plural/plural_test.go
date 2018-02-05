// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plural

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/text/language"
)

func TestGetIntApprox(t *testing.T) ***REMOVED***
	const big = 1234567890
	testCases := []struct ***REMOVED***
		digits string
		start  int
		end    int
		nMod   int
		want   int
	***REMOVED******REMOVED***
		***REMOVED***"123", 0, 1, 1, 1***REMOVED***,
		***REMOVED***"123", 0, 2, 1, big***REMOVED***,
		***REMOVED***"123", 0, 2, 2, 12***REMOVED***,
		***REMOVED***"123", 3, 4, 2, 0***REMOVED***,
		***REMOVED***"12345", 3, 4, 2, 4***REMOVED***,
		***REMOVED***"40", 0, 1, 2, 4***REMOVED***,
		***REMOVED***"1", 0, 7, 2, big***REMOVED***,

		***REMOVED***"123", 0, 5, 2, big***REMOVED***,
		***REMOVED***"123", 0, 5, 3, big***REMOVED***,
		***REMOVED***"123", 0, 5, 4, big***REMOVED***,
		***REMOVED***"123", 0, 5, 5, 12300***REMOVED***,
		***REMOVED***"123", 0, 5, 6, 12300***REMOVED***,
		***REMOVED***"123", 0, 5, 7, 12300***REMOVED***,

		// Translation of examples in MatchDigits.
		// Integer parts
		***REMOVED***"123", 0, 3, 3, 123***REMOVED***,  // 123
		***REMOVED***"1234", 0, 3, 3, 123***REMOVED***, // 123.4
		***REMOVED***"1", 0, 6, 8, 100000***REMOVED***, // 100000

		// Fraction parts
		***REMOVED***"123", 3, 3, 3, 0***REMOVED***,   // 123
		***REMOVED***"1234", 3, 4, 3, 4***REMOVED***,  // 123.4
		***REMOVED***"1234", 3, 5, 3, 40***REMOVED***, // 123.40
		***REMOVED***"1", 6, 8, 8, 0***REMOVED***,     // 100000.00
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(fmt.Sprintf("%s:%d:%d/%d", tc.digits, tc.start, tc.end, tc.nMod), func(t *testing.T) ***REMOVED***
			got := getIntApprox(mkDigits(tc.digits), tc.start, tc.end, tc.nMod, big)
			if got != tc.want ***REMOVED***
				t.Errorf("got %d; want %d", got, tc.want)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func mkDigits(s string) []byte ***REMOVED***
	b := []byte(s)
	for i := range b ***REMOVED***
		b[i] -= '0'
	***REMOVED***
	return b
***REMOVED***

func TestValidForms(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		tag  language.Tag
		want []Form
	***REMOVED******REMOVED***
		***REMOVED***language.AmericanEnglish, []Form***REMOVED***Other, One***REMOVED******REMOVED***,
		***REMOVED***language.Portuguese, []Form***REMOVED***Other, One***REMOVED******REMOVED***,
		***REMOVED***language.Latvian, []Form***REMOVED***Other, Zero, One***REMOVED******REMOVED***,
		***REMOVED***language.Arabic, []Form***REMOVED***Other, Zero, One, Two, Few, Many***REMOVED******REMOVED***,
		***REMOVED***language.Russian, []Form***REMOVED***Other, One, Few, Many***REMOVED******REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		got := validForms(cardinal, tc.tag)
		if !reflect.DeepEqual(got, tc.want) ***REMOVED***
			t.Errorf("validForms(%v): got %v; want %v", tc.tag, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestOrdinal(t *testing.T) ***REMOVED***
	testPlurals(t, Ordinal, ordinalTests)
***REMOVED***

func TestCardinal(t *testing.T) ***REMOVED***
	testPlurals(t, Cardinal, cardinalTests)
***REMOVED***

func testPlurals(t *testing.T, p *Rules, testCases []pluralTest) ***REMOVED***
	for _, tc := range testCases ***REMOVED***
		for _, loc := range strings.Split(tc.locales, " ") ***REMOVED***
			tag := language.MustParse(loc)
			// Test integers
			for _, s := range tc.integer ***REMOVED***
				a := strings.Split(s, "~")
				from := parseUint(t, a[0])
				to := from
				if len(a) > 1 ***REMOVED***
					to = parseUint(t, a[1])
				***REMOVED***
				for n := from; n <= to; n++ ***REMOVED***
					t.Run(fmt.Sprintf("%s/int(%d)", loc, n), func(t *testing.T) ***REMOVED***
						if f := p.matchComponents(tag, n, 0, 0); f != Form(tc.form) ***REMOVED***
							t.Errorf("matchComponents: got %v; want %v", f, Form(tc.form))
						***REMOVED***
						digits := []byte(fmt.Sprint(n))
						for i := range digits ***REMOVED***
							digits[i] -= '0'
						***REMOVED***
						if f := p.MatchDigits(tag, digits, len(digits), 0); f != Form(tc.form) ***REMOVED***
							t.Errorf("MatchDigits: got %v; want %v", f, Form(tc.form))
						***REMOVED***
					***REMOVED***)
				***REMOVED***
			***REMOVED***
			// Test decimals
			for _, s := range tc.decimal ***REMOVED***
				a := strings.Split(s, "~")
				from, scale := parseFixedPoint(t, a[0])
				to := from
				if len(a) > 1 ***REMOVED***
					var toScale int
					if to, toScale = parseFixedPoint(t, a[1]); toScale != scale ***REMOVED***
						t.Fatalf("%s:%s: non-matching scales %d versus %d", loc, s, scale, toScale)
					***REMOVED***
				***REMOVED***
				m := 1
				for i := 0; i < scale; i++ ***REMOVED***
					m *= 10
				***REMOVED***
				for n := from; n <= to; n++ ***REMOVED***
					num := fmt.Sprintf("%[1]d.%0[3]*[2]d", n/m, n%m, scale)
					name := fmt.Sprintf("%s:dec(%s)", loc, num)
					t.Run(name, func(t *testing.T) ***REMOVED***
						ff := n % m
						tt := ff
						w := scale
						for tt > 0 && tt%10 == 0 ***REMOVED***
							w--
							tt /= 10
						***REMOVED***
						if f := p.MatchPlural(tag, n/m, scale, w, ff, tt); f != Form(tc.form) ***REMOVED***
							t.Errorf("MatchPlural: got %v; want %v", f, Form(tc.form))
						***REMOVED***
						if f := p.matchComponents(tag, n/m, n%m, scale); f != Form(tc.form) ***REMOVED***
							t.Errorf("matchComponents: got %v; want %v", f, Form(tc.form))
						***REMOVED***
						exp := strings.IndexByte(num, '.')
						digits := []byte(strings.Replace(num, ".", "", 1))
						for i := range digits ***REMOVED***
							digits[i] -= '0'
						***REMOVED***
						if f := p.MatchDigits(tag, digits, exp, scale); f != Form(tc.form) ***REMOVED***
							t.Errorf("MatchDigits: got %v; want %v", f, Form(tc.form))
						***REMOVED***
					***REMOVED***)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func parseUint(t *testing.T, s string) int ***REMOVED***
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	return int(val)
***REMOVED***

func parseFixedPoint(t *testing.T, s string) (val, scale int) ***REMOVED***
	p := strings.Index(s, ".")
	s = strings.Replace(s, ".", "", 1)
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	return int(v), len(s) - p
***REMOVED***

func BenchmarkPluralSimpleCases(b *testing.B) ***REMOVED***
	p := Cardinal
	en, _ := language.CompactIndex(language.English)
	zh, _ := language.CompactIndex(language.Chinese)
	for i := 0; i < b.N; i++ ***REMOVED***
		matchPlural(p, en, 0, 0, 0)  // 0
		matchPlural(p, en, 1, 0, 0)  // 1
		matchPlural(p, en, 2, 12, 3) // 2.120
		matchPlural(p, zh, 0, 0, 0)  // 0
		matchPlural(p, zh, 1, 0, 0)  // 1
		matchPlural(p, zh, 2, 12, 3) // 2.120
	***REMOVED***
***REMOVED***

func BenchmarkPluralComplexCases(b *testing.B) ***REMOVED***
	p := Cardinal
	ar, _ := language.CompactIndex(language.Arabic)
	lv, _ := language.CompactIndex(language.Latvian)
	for i := 0; i < b.N; i++ ***REMOVED***
		matchPlural(p, lv, 0, 19, 2)    // 0.19
		matchPlural(p, lv, 11, 0, 3)    // 11.000
		matchPlural(p, lv, 100, 123, 4) // 0.1230
		matchPlural(p, ar, 0, 0, 0)     // 0
		matchPlural(p, ar, 110, 0, 0)   // 110
		matchPlural(p, ar, 99, 99, 2)   // 99.99
	***REMOVED***
***REMOVED***
