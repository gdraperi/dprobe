// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package idna

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/internal/ucd"
)

func TestAllocToUnicode(t *testing.T) ***REMOVED***
	avg := testtext.AllocsPerRun(1000, func() ***REMOVED***
		ToUnicode("www.golang.org")
	***REMOVED***)
	if avg > 0 ***REMOVED***
		t.Errorf("got %f; want 0", avg)
	***REMOVED***
***REMOVED***

func TestAllocToASCII(t *testing.T) ***REMOVED***
	avg := testtext.AllocsPerRun(1000, func() ***REMOVED***
		ToASCII("www.golang.org")
	***REMOVED***)
	if avg > 0 ***REMOVED***
		t.Errorf("got %f; want 0", avg)
	***REMOVED***
***REMOVED***

func TestProfiles(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		name      string
		want, got *Profile
	***REMOVED******REMOVED***
		***REMOVED***"Punycode", punycode, New()***REMOVED***,
		***REMOVED***"Registration", registration, New(ValidateForRegistration())***REMOVED***,
		***REMOVED***"Registration", registration, New(
			ValidateForRegistration(),
			VerifyDNSLength(true),
			BidiRule(),
		)***REMOVED***,
		***REMOVED***"Lookup", lookup, New(MapForLookup(), BidiRule(), Transitional(true))***REMOVED***,
		***REMOVED***"Display", display, New(MapForLookup(), BidiRule())***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		// Functions are not comparable, but the printed version will include
		// their pointers.
		got := fmt.Sprintf("%#v", tc.got)
		want := fmt.Sprintf("%#v", tc.want)
		if got != want ***REMOVED***
			t.Errorf("%s: \ngot  %#v,\nwant %#v", tc.name, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// doTest performs a single test f(input) and verifies that the output matches
// out and that the returned error is expected. The errors string contains
// all allowed error codes as categorized in
// http://www.unicode.org/Public/idna/9.0.0/IdnaTest.txt:
// P: Processing
// V: Validity
// A: to ASCII
// B: Bidi
// C: Context J
func doTest(t *testing.T, f func(string) (string, error), name, input, want, errors string) ***REMOVED***
	errors = strings.Trim(errors, "[]")
	test := "ok"
	if errors != "" ***REMOVED***
		test = "err:" + errors
	***REMOVED***
	// Replace some of the escape sequences to make it easier to single out
	// tests on the command name.
	in := strings.Trim(strconv.QuoteToASCII(input), `"`)
	in = strings.Replace(in, `\u`, "#", -1)
	in = strings.Replace(in, `\U`, "#", -1)
	name = fmt.Sprintf("%s/%s/%s", name, in, test)

	testtext.Run(t, name, func(t *testing.T) ***REMOVED***
		got, err := f(input)

		if err != nil ***REMOVED***
			code := err.(interface ***REMOVED***
				code() string
			***REMOVED***).code()
			if strings.Index(errors, code) == -1 ***REMOVED***
				t.Errorf("error %q not in set of expected errors ***REMOVED***%v***REMOVED***", code, errors)
			***REMOVED***
		***REMOVED*** else if errors != "" ***REMOVED***
			t.Errorf("no errors; want error in ***REMOVED***%v***REMOVED***", errors)
		***REMOVED***

		if want != "" && got != want ***REMOVED***
			t.Errorf(`string: got %+q; want %+q`, got, want)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestConformance(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	r := gen.OpenUnicodeFile("idna", "", "IdnaTest.txt")
	defer r.Close()

	section := "main"
	started := false
	p := ucd.New(r, ucd.CommentHandler(func(s string) ***REMOVED***
		if started ***REMOVED***
			section = strings.ToLower(strings.Split(s, " ")[0])
		***REMOVED***
	***REMOVED***))
	transitional := New(Transitional(true), VerifyDNSLength(true), BidiRule(), MapForLookup())
	nonTransitional := New(VerifyDNSLength(true), BidiRule(), MapForLookup())
	for p.Next() ***REMOVED***
		started = true

		// What to test
		profiles := []*Profile***REMOVED******REMOVED***
		switch p.String(0) ***REMOVED***
		case "T":
			profiles = append(profiles, transitional)
		case "N":
			profiles = append(profiles, nonTransitional)
		case "B":
			profiles = append(profiles, transitional)
			profiles = append(profiles, nonTransitional)
		***REMOVED***

		src := unescape(p.String(1))

		wantToUnicode := unescape(p.String(2))
		if wantToUnicode == "" ***REMOVED***
			wantToUnicode = src
		***REMOVED***
		wantToASCII := unescape(p.String(3))
		if wantToASCII == "" ***REMOVED***
			wantToASCII = wantToUnicode
		***REMOVED***
		wantErrToUnicode := ""
		if strings.HasPrefix(wantToUnicode, "[") ***REMOVED***
			wantErrToUnicode = wantToUnicode
			wantToUnicode = ""
		***REMOVED***
		wantErrToASCII := ""
		if strings.HasPrefix(wantToASCII, "[") ***REMOVED***
			wantErrToASCII = wantToASCII
			wantToASCII = ""
		***REMOVED***

		// TODO: also do IDNA tests.
		// invalidInIDNA2008 := p.String(4) == "NV8"

		for _, p := range profiles ***REMOVED***
			name := fmt.Sprintf("%s:%s", section, p)
			doTest(t, p.ToUnicode, name+":ToUnicode", src, wantToUnicode, wantErrToUnicode)
			doTest(t, p.ToASCII, name+":ToASCII", src, wantToASCII, wantErrToASCII)
		***REMOVED***
	***REMOVED***
***REMOVED***

func unescape(s string) string ***REMOVED***
	s, err := strconv.Unquote(`"` + s + `"`)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return s
***REMOVED***

func BenchmarkProfile(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		Lookup.ToASCII("www.yahoogle.com")
	***REMOVED***
***REMOVED***
