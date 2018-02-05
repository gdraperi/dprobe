// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package number

import (
	"flag"
	"log"
	"reflect"
	"testing"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

var draft = flag.String("draft",
	"contributed",
	`Minimal draft requirements (approved, contributed, provisional, unconfirmed).`)

func TestNumberSystems(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder***REMOVED******REMOVED***
	d.SetDirFilter("supplemental")
	d.SetSectionFilter("numberingSystem")
	data, err := d.DecodeZip(r)
	if err != nil ***REMOVED***
		t.Fatalf("DecodeZip: %v", err)
	***REMOVED***

	for _, ns := range data.Supplemental().NumberingSystems.NumberingSystem ***REMOVED***
		n := systemMap[ns.Id]
		if int(n) >= len(numSysData) ***REMOVED***
			continue
		***REMOVED***
		info := InfoFromLangID(0, ns.Id)
		val := '0'
		for _, rWant := range ns.Digits ***REMOVED***
			if rGot := info.Digit(val); rGot != rWant ***REMOVED***
				t.Errorf("%s:%d: got %U; want %U", ns.Id, val, rGot, rWant)
			***REMOVED***
			val++
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSymbols(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	draft, err := cldr.ParseDraft(*draft)
	if err != nil ***REMOVED***
		log.Fatalf("invalid draft level: %v", err)
	***REMOVED***

	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder***REMOVED******REMOVED***
	d.SetDirFilter("main")
	d.SetSectionFilter("numbers")
	data, err := d.DecodeZip(r)
	if err != nil ***REMOVED***
		t.Fatalf("DecodeZip: %v", err)
	***REMOVED***

	for _, lang := range data.Locales() ***REMOVED***
		ldml := data.RawLDML(lang)
		if ldml.Numbers == nil ***REMOVED***
			continue
		***REMOVED***
		langIndex, ok := language.CompactIndex(language.MustParse(lang))
		if !ok ***REMOVED***
			t.Fatalf("No compact index for language %s", lang)
		***REMOVED***

		syms := cldr.MakeSlice(&ldml.Numbers.Symbols)
		syms.SelectDraft(draft)

		for _, sym := range ldml.Numbers.Symbols ***REMOVED***
			if sym.NumberSystem == "" ***REMOVED***
				continue
			***REMOVED***
			testCases := []struct ***REMOVED***
				name string
				st   SymbolType
				x    interface***REMOVED******REMOVED***
			***REMOVED******REMOVED***
				***REMOVED***"Decimal", SymDecimal, sym.Decimal***REMOVED***,
				***REMOVED***"Group", SymGroup, sym.Group***REMOVED***,
				***REMOVED***"List", SymList, sym.List***REMOVED***,
				***REMOVED***"PercentSign", SymPercentSign, sym.PercentSign***REMOVED***,
				***REMOVED***"PlusSign", SymPlusSign, sym.PlusSign***REMOVED***,
				***REMOVED***"MinusSign", SymMinusSign, sym.MinusSign***REMOVED***,
				***REMOVED***"Exponential", SymExponential, sym.Exponential***REMOVED***,
				***REMOVED***"SuperscriptingExponent", SymSuperscriptingExponent, sym.SuperscriptingExponent***REMOVED***,
				***REMOVED***"PerMille", SymPerMille, sym.PerMille***REMOVED***,
				***REMOVED***"Infinity", SymInfinity, sym.Infinity***REMOVED***,
				***REMOVED***"NaN", SymNan, sym.Nan***REMOVED***,
				***REMOVED***"TimeSeparator", SymTimeSeparator, sym.TimeSeparator***REMOVED***,
			***REMOVED***
			info := InfoFromLangID(langIndex, sym.NumberSystem)
			for _, tc := range testCases ***REMOVED***
				// Extract the wanted value.
				v := reflect.ValueOf(tc.x)
				if v.Len() == 0 ***REMOVED***
					return
				***REMOVED***
				if v.Len() > 1 ***REMOVED***
					t.Fatalf("Multiple values of %q within single symbol not supported.", tc.name)
				***REMOVED***
				want := v.Index(0).MethodByName("Data").Call(nil)[0].String()
				got := info.Symbol(tc.st)
				if got != want ***REMOVED***
					t.Errorf("%s:%s:%s: got %q; want %q", lang, sym.NumberSystem, tc.name, got, want)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
