// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/internal"
	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/number"
	"golang.org/x/text/internal/stringset"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

var (
	test = flag.Bool("test", false,
		"test existing tables; can be used to compare web data with package data.")
	outputFile     = flag.String("output", "tables.go", "output file")
	outputTestFile = flag.String("testoutput", "data_test.go", "output file")

	draft = flag.String("draft",
		"contributed",
		`Minimal draft requirements (approved, contributed, provisional, unconfirmed).`)
)

func main() ***REMOVED***
	gen.Init()

	const pkg = "number"

	gen.Repackage("gen_common.go", "common.go", pkg)
	// Read the CLDR zip file.
	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder***REMOVED******REMOVED***
	d.SetDirFilter("supplemental", "main")
	d.SetSectionFilter("numbers", "numberingSystem")
	data, err := d.DecodeZip(r)
	if err != nil ***REMOVED***
		log.Fatalf("DecodeZip: %v", err)
	***REMOVED***

	w := gen.NewCodeWriter()
	defer w.WriteGoFile(*outputFile, pkg)

	fmt.Fprintln(w, `import "golang.org/x/text/internal/stringset"`)

	gen.WriteCLDRVersion(w)

	genNumSystem(w, data)
	genSymbols(w, data)
	genFormats(w, data)
***REMOVED***

var systemMap = map[string]system***REMOVED***"latn": 0***REMOVED***

func getNumberSystem(str string) system ***REMOVED***
	ns, ok := systemMap[str]
	if !ok ***REMOVED***
		log.Fatalf("No index for numbering system %q", str)
	***REMOVED***
	return ns
***REMOVED***

func genNumSystem(w *gen.CodeWriter, data *cldr.CLDR) ***REMOVED***
	numSysData := []systemData***REMOVED***
		***REMOVED***digitSize: 1, zero: [4]byte***REMOVED***'0'***REMOVED******REMOVED***,
	***REMOVED***

	for _, ns := range data.Supplemental().NumberingSystems.NumberingSystem ***REMOVED***
		if len(ns.Digits) == 0 ***REMOVED***
			continue
		***REMOVED***
		switch ns.Id ***REMOVED***
		case "latn":
			// hard-wired
			continue
		case "hanidec":
			// non-consecutive digits: treat as "algorithmic"
			continue
		***REMOVED***

		zero, sz := utf8.DecodeRuneInString(ns.Digits)
		if ns.Digits[sz-1]+9 > 0xBF ***REMOVED*** // 1011 1111: highest continuation byte
			log.Fatalf("Last byte of zero value overflows for %s", ns.Id)
		***REMOVED***

		i := rune(0)
		for _, r := range ns.Digits ***REMOVED***
			// Verify that we can do simple math on the UTF-8 byte sequence
			// of zero to get the digit.
			if zero+i != r ***REMOVED***
				// Runes not consecutive.
				log.Fatalf("Digit %d of %s (%U) is not offset correctly from zero value", i, ns.Id, r)
			***REMOVED***
			i++
		***REMOVED***
		var x [utf8.UTFMax]byte
		utf8.EncodeRune(x[:], zero)
		id := system(len(numSysData))
		systemMap[ns.Id] = id
		numSysData = append(numSysData, systemData***REMOVED***
			id:        id,
			digitSize: byte(sz),
			zero:      x,
		***REMOVED***)
	***REMOVED***
	w.WriteVar("numSysData", numSysData)

	algoID := system(len(numSysData))
	fmt.Fprintln(w, "const (")
	for _, ns := range data.Supplemental().NumberingSystems.NumberingSystem ***REMOVED***
		id, ok := systemMap[ns.Id]
		if !ok ***REMOVED***
			id = algoID
			systemMap[ns.Id] = id
			algoID++
		***REMOVED***
		fmt.Fprintf(w, "num%s = %#x\n", strings.Title(ns.Id), id)
	***REMOVED***
	fmt.Fprintln(w, "numNumberSystems")
	fmt.Fprintln(w, ")")

	fmt.Fprintln(w, "var systemMap = map[string]system***REMOVED***")
	for _, ns := range data.Supplemental().NumberingSystems.NumberingSystem ***REMOVED***
		fmt.Fprintf(w, "%q: num%s,\n", ns.Id, strings.Title(ns.Id))
		w.Size += len(ns.Id) + 16 + 1 // very coarse approximation
	***REMOVED***
	fmt.Fprintln(w, "***REMOVED***")
***REMOVED***

func genSymbols(w *gen.CodeWriter, data *cldr.CLDR) ***REMOVED***
	d, err := cldr.ParseDraft(*draft)
	if err != nil ***REMOVED***
		log.Fatalf("invalid draft level: %v", err)
	***REMOVED***

	nNumberSystems := system(len(systemMap))

	type symbols [NumSymbolTypes]string

	type key struct ***REMOVED***
		tag    int // from language.CompactIndex
		system system
	***REMOVED***
	symbolMap := map[key]*symbols***REMOVED******REMOVED***

	defaults := map[int]system***REMOVED******REMOVED***

	for _, lang := range data.Locales() ***REMOVED***
		ldml := data.RawLDML(lang)
		if ldml.Numbers == nil ***REMOVED***
			continue
		***REMOVED***
		langIndex, ok := language.CompactIndex(language.MustParse(lang))
		if !ok ***REMOVED***
			log.Fatalf("No compact index for language %s", lang)
		***REMOVED***
		if d := ldml.Numbers.DefaultNumberingSystem; len(d) > 0 ***REMOVED***
			defaults[langIndex] = getNumberSystem(d[0].Data())
		***REMOVED***

		syms := cldr.MakeSlice(&ldml.Numbers.Symbols)
		syms.SelectDraft(d)

		getFirst := func(name string, x interface***REMOVED******REMOVED***) string ***REMOVED***
			v := reflect.ValueOf(x)
			slice := cldr.MakeSlice(x)
			slice.SelectAnyOf("alt", "", "alt")
			if reflect.Indirect(v).Len() == 0 ***REMOVED***
				return ""
			***REMOVED*** else if reflect.Indirect(v).Len() > 1 ***REMOVED***
				log.Fatalf("%s: multiple values of %q within single symbol not supported.", lang, name)
			***REMOVED***
			return reflect.Indirect(v).Index(0).MethodByName("Data").Call(nil)[0].String()
		***REMOVED***

		for _, sym := range ldml.Numbers.Symbols ***REMOVED***
			if sym.NumberSystem == "" ***REMOVED***
				// This is just linking the default of root to "latn".
				continue
			***REMOVED***
			symbolMap[key***REMOVED***langIndex, getNumberSystem(sym.NumberSystem)***REMOVED***] = &symbols***REMOVED***
				SymDecimal:                getFirst("decimal", &sym.Decimal),
				SymGroup:                  getFirst("group", &sym.Group),
				SymList:                   getFirst("list", &sym.List),
				SymPercentSign:            getFirst("percentSign", &sym.PercentSign),
				SymPlusSign:               getFirst("plusSign", &sym.PlusSign),
				SymMinusSign:              getFirst("minusSign", &sym.MinusSign),
				SymExponential:            getFirst("exponential", &sym.Exponential),
				SymSuperscriptingExponent: getFirst("superscriptingExponent", &sym.SuperscriptingExponent),
				SymPerMille:               getFirst("perMille", &sym.PerMille),
				SymInfinity:               getFirst("infinity", &sym.Infinity),
				SymNan:                    getFirst("nan", &sym.Nan),
				SymTimeSeparator:          getFirst("timeSeparator", &sym.TimeSeparator),
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Expand all values.
	for k, syms := range symbolMap ***REMOVED***
		for t := SymDecimal; t < NumSymbolTypes; t++ ***REMOVED***
			p := k.tag
			for syms[t] == "" ***REMOVED***
				p = int(internal.Parent[p])
				if pSyms, ok := symbolMap[key***REMOVED***p, k.system***REMOVED***]; ok && (*pSyms)[t] != "" ***REMOVED***
					syms[t] = (*pSyms)[t]
					break
				***REMOVED***
				if p == 0 /* und */ ***REMOVED***
					// Default to root, latn.
					syms[t] = (*symbolMap[key***REMOVED******REMOVED***])[t]
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Unique the symbol sets and write the string data.
	m := map[symbols]int***REMOVED******REMOVED***
	sb := stringset.NewBuilder()

	symIndex := [][NumSymbolTypes]byte***REMOVED******REMOVED***

	for ns := system(0); ns < nNumberSystems; ns++ ***REMOVED***
		for _, l := range data.Locales() ***REMOVED***
			langIndex, _ := language.CompactIndex(language.MustParse(l))
			s := symbolMap[key***REMOVED***langIndex, ns***REMOVED***]
			if s == nil ***REMOVED***
				continue
			***REMOVED***
			if _, ok := m[*s]; !ok ***REMOVED***
				m[*s] = len(symIndex)
				sb.Add(s[:]...)
				var x [NumSymbolTypes]byte
				for i := SymDecimal; i < NumSymbolTypes; i++ ***REMOVED***
					x[i] = byte(sb.Index((*s)[i]))
				***REMOVED***
				symIndex = append(symIndex, x)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	w.WriteVar("symIndex", symIndex)
	w.WriteVar("symData", sb.Set())

	// resolveSymbolIndex gets the index from the closest matching locale,
	// including the locale itself.
	resolveSymbolIndex := func(langIndex int, ns system) symOffset ***REMOVED***
		for ***REMOVED***
			if sym := symbolMap[key***REMOVED***langIndex, ns***REMOVED***]; sym != nil ***REMOVED***
				return symOffset(m[*sym])
			***REMOVED***
			if langIndex == 0 ***REMOVED***
				return 0 // und, latn
			***REMOVED***
			langIndex = int(internal.Parent[langIndex])
		***REMOVED***
	***REMOVED***

	// Create an index with the symbols for each locale for the latn numbering
	// system. If this is not the default, or the only one, for a locale, we
	// will overwrite the value later.
	var langToDefaults [language.NumCompactTags]symOffset
	for _, l := range data.Locales() ***REMOVED***
		langIndex, _ := language.CompactIndex(language.MustParse(l))
		langToDefaults[langIndex] = resolveSymbolIndex(langIndex, 0)
	***REMOVED***

	// Delete redundant entries.
	for _, l := range data.Locales() ***REMOVED***
		langIndex, _ := language.CompactIndex(language.MustParse(l))
		def := defaults[langIndex]
		syms := symbolMap[key***REMOVED***langIndex, def***REMOVED***]
		if syms == nil ***REMOVED***
			continue
		***REMOVED***
		for ns := system(0); ns < nNumberSystems; ns++ ***REMOVED***
			if ns == def ***REMOVED***
				continue
			***REMOVED***
			if altSyms, ok := symbolMap[key***REMOVED***langIndex, ns***REMOVED***]; ok && *altSyms == *syms ***REMOVED***
				delete(symbolMap, key***REMOVED***langIndex, ns***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Create a sorted list of alternatives per language. This will only need to
	// be referenced if a user specified an alternative numbering system.
	var langToAlt []altSymData
	for _, l := range data.Locales() ***REMOVED***
		langIndex, _ := language.CompactIndex(language.MustParse(l))
		start := len(langToAlt)
		if start >= hasNonLatnMask ***REMOVED***
			log.Fatalf("Number of alternative assignments >= %x", hasNonLatnMask)
		***REMOVED***
		// Create the entry for the default value.
		def := defaults[langIndex]
		langToAlt = append(langToAlt, altSymData***REMOVED***
			compactTag: uint16(langIndex),
			system:     def,
			symIndex:   resolveSymbolIndex(langIndex, def),
		***REMOVED***)

		for ns := system(0); ns < nNumberSystems; ns++ ***REMOVED***
			if def == ns ***REMOVED***
				continue
			***REMOVED***
			if sym := symbolMap[key***REMOVED***langIndex, ns***REMOVED***]; sym != nil ***REMOVED***
				langToAlt = append(langToAlt, altSymData***REMOVED***
					compactTag: uint16(langIndex),
					system:     ns,
					symIndex:   resolveSymbolIndex(langIndex, ns),
				***REMOVED***)
			***REMOVED***
		***REMOVED***
		if def == 0 && len(langToAlt) == start+1 ***REMOVED***
			// No additional data: erase the entry.
			langToAlt = langToAlt[:start]
		***REMOVED*** else ***REMOVED***
			// Overwrite the entry in langToDefaults.
			langToDefaults[langIndex] = hasNonLatnMask | symOffset(start)
		***REMOVED***
	***REMOVED***
	w.WriteComment(`
langToDefaults maps a compact language index to the default numbering system
and default symbol set`)
	w.WriteVar("langToDefaults", langToDefaults)

	w.WriteComment(`
langToAlt is a list of numbering system and symbol set pairs, sorted and
marked by compact language index.`)
	w.WriteVar("langToAlt", langToAlt)
***REMOVED***

// genFormats generates the lookup table for decimal, scientific and percent
// patterns.
//
// CLDR allows for patterns to be different per language for different numbering
// systems. In practice the patterns are set to be consistent for a language
// independent of the numbering system. genFormats verifies that no language
// deviates from this.
func genFormats(w *gen.CodeWriter, data *cldr.CLDR) ***REMOVED***
	d, err := cldr.ParseDraft(*draft)
	if err != nil ***REMOVED***
		log.Fatalf("invalid draft level: %v", err)
	***REMOVED***

	// Fill the first slot with a dummy so we can identify unspecified tags.
	formats := []number.Pattern***REMOVED******REMOVED******REMOVED******REMOVED***
	patterns := map[string]int***REMOVED******REMOVED***

	// TODO: It would be possible to eliminate two of these slices by having
	// another indirection and store a reference to the combination of patterns.
	decimal := make([]byte, language.NumCompactTags)
	scientific := make([]byte, language.NumCompactTags)
	percent := make([]byte, language.NumCompactTags)

	for _, lang := range data.Locales() ***REMOVED***
		ldml := data.RawLDML(lang)
		if ldml.Numbers == nil ***REMOVED***
			continue
		***REMOVED***
		langIndex, ok := language.CompactIndex(language.MustParse(lang))
		if !ok ***REMOVED***
			log.Fatalf("No compact index for language %s", lang)
		***REMOVED***
		type patternSlice []*struct ***REMOVED***
			cldr.Common
			Numbers string `xml:"numbers,attr"`
			Count   string `xml:"count,attr"`
		***REMOVED***

		add := func(name string, tags []byte, ps patternSlice) ***REMOVED***
			sl := cldr.MakeSlice(&ps)
			sl.SelectDraft(d)
			if len(ps) == 0 ***REMOVED***
				return
			***REMOVED***
			if len(ps) > 2 || len(ps) == 2 && ps[0] != ps[1] ***REMOVED***
				log.Fatalf("Inconsistent %d patterns for language %s", name, lang)
			***REMOVED***
			s := ps[0].Data()

			index, ok := patterns[s]
			if !ok ***REMOVED***
				nf, err := number.ParsePattern(s)
				if err != nil ***REMOVED***
					log.Fatal(err)
				***REMOVED***
				index = len(formats)
				patterns[s] = index
				formats = append(formats, *nf)
			***REMOVED***
			tags[langIndex] = byte(index)
		***REMOVED***

		for _, df := range ldml.Numbers.DecimalFormats ***REMOVED***
			for _, l := range df.DecimalFormatLength ***REMOVED***
				if l.Type != "" ***REMOVED***
					continue
				***REMOVED***
				for _, f := range l.DecimalFormat ***REMOVED***
					add("decimal", decimal, f.Pattern)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		for _, df := range ldml.Numbers.ScientificFormats ***REMOVED***
			for _, l := range df.ScientificFormatLength ***REMOVED***
				if l.Type != "" ***REMOVED***
					continue
				***REMOVED***
				for _, f := range l.ScientificFormat ***REMOVED***
					add("scientific", scientific, f.Pattern)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		for _, df := range ldml.Numbers.PercentFormats ***REMOVED***
			for _, l := range df.PercentFormatLength ***REMOVED***
				if l.Type != "" ***REMOVED***
					continue
				***REMOVED***
				for _, f := range l.PercentFormat ***REMOVED***
					add("percent", percent, f.Pattern)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Complete the parent tag array to reflect inheritance. An index of 0
	// indicates an unspecified value.
	for _, data := range [][]byte***REMOVED***decimal, scientific, percent***REMOVED*** ***REMOVED***
		for i := range data ***REMOVED***
			p := uint16(i)
			for ; data[p] == 0; p = internal.Parent[p] ***REMOVED***
			***REMOVED***
			data[i] = data[p]
		***REMOVED***
	***REMOVED***
	w.WriteVar("tagToDecimal", decimal)
	w.WriteVar("tagToScientific", scientific)
	w.WriteVar("tagToPercent", percent)

	value := strings.Replace(fmt.Sprintf("%#v", formats), "number.", "", -1)
	// Break up the lines. This won't give ideal perfect formatting, but it is
	// better than one huge line.
	value = strings.Replace(value, ", ", ",\n", -1)
	fmt.Fprintf(w, "var formats = %s\n", value)
***REMOVED***
