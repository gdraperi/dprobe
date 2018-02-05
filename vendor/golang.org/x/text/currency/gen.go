// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// Generator for currency-related data.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/internal"
	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/tag"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

var (
	test = flag.Bool("test", false,
		"test existing tables; can be used to compare web data with package data.")
	outputFile = flag.String("output", "tables.go", "output file")

	draft = flag.String("draft",
		"contributed",
		`Minimal draft requirements (approved, contributed, provisional, unconfirmed).`)
)

func main() ***REMOVED***
	gen.Init()

	gen.Repackage("gen_common.go", "common.go", "currency")

	// Read the CLDR zip file.
	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder***REMOVED******REMOVED***
	d.SetDirFilter("supplemental", "main")
	d.SetSectionFilter("numbers")
	data, err := d.DecodeZip(r)
	if err != nil ***REMOVED***
		log.Fatalf("DecodeZip: %v", err)
	***REMOVED***

	w := gen.NewCodeWriter()
	defer w.WriteGoFile(*outputFile, "currency")

	fmt.Fprintln(w, `import "golang.org/x/text/internal/tag"`)

	gen.WriteCLDRVersion(w)
	b := &builder***REMOVED******REMOVED***
	b.genCurrencies(w, data.Supplemental())
	b.genSymbols(w, data)
***REMOVED***

var constants = []string***REMOVED***
	// Undefined and testing.
	"XXX", "XTS",
	// G11 currencies https://en.wikipedia.org/wiki/G10_currencies.
	"USD", "EUR", "JPY", "GBP", "CHF", "AUD", "NZD", "CAD", "SEK", "NOK", "DKK",
	// Precious metals.
	"XAG", "XAU", "XPT", "XPD",

	// Additional common currencies as defined by CLDR.
	"BRL", "CNY", "INR", "RUB", "HKD", "IDR", "KRW", "MXN", "PLN", "SAR",
	"THB", "TRY", "TWD", "ZAR",
***REMOVED***

type builder struct ***REMOVED***
	currencies    tag.Index
	numCurrencies int
***REMOVED***

func (b *builder) genCurrencies(w *gen.CodeWriter, data *cldr.SupplementalData) ***REMOVED***
	// 3-letter ISO currency codes
	// Start with dummy to let index start at 1.
	currencies := []string***REMOVED***"\x00\x00\x00\x00"***REMOVED***

	// currency codes
	for _, reg := range data.CurrencyData.Region ***REMOVED***
		for _, cur := range reg.Currency ***REMOVED***
			currencies = append(currencies, cur.Iso4217)
		***REMOVED***
	***REMOVED***
	// Not included in the list for some reasons:
	currencies = append(currencies, "MVP")

	sort.Strings(currencies)
	// Unique the elements.
	k := 0
	for i := 1; i < len(currencies); i++ ***REMOVED***
		if currencies[k] != currencies[i] ***REMOVED***
			currencies[k+1] = currencies[i]
			k++
		***REMOVED***
	***REMOVED***
	currencies = currencies[:k+1]

	// Close with dummy for simpler and faster searching.
	currencies = append(currencies, "\xff\xff\xff\xff")

	// Write currency values.
	fmt.Fprintln(w, "const (")
	for _, c := range constants ***REMOVED***
		index := sort.SearchStrings(currencies, c)
		fmt.Fprintf(w, "\t%s = %d\n", strings.ToLower(c), index)
	***REMOVED***
	fmt.Fprint(w, ")")

	// Compute currency-related data that we merge into the table.
	for _, info := range data.CurrencyData.Fractions[0].Info ***REMOVED***
		if info.Iso4217 == "DEFAULT" ***REMOVED***
			continue
		***REMOVED***
		standard := getRoundingIndex(info.Digits, info.Rounding, 0)
		cash := getRoundingIndex(info.CashDigits, info.CashRounding, standard)

		index := sort.SearchStrings(currencies, info.Iso4217)
		currencies[index] += mkCurrencyInfo(standard, cash)
	***REMOVED***

	// Set default values for entries that weren't touched.
	for i, c := range currencies ***REMOVED***
		if len(c) == 3 ***REMOVED***
			currencies[i] += mkCurrencyInfo(0, 0)
		***REMOVED***
	***REMOVED***

	b.currencies = tag.Index(strings.Join(currencies, ""))
	w.WriteComment(`
	currency holds an alphabetically sorted list of canonical 3-letter currency
	identifiers. Each identifier is followed by a byte of type currencyInfo,
	defined in gen_common.go.`)
	w.WriteConst("currency", b.currencies)

	// Hack alert: gofmt indents a trailing comment after an indented string.
	// Ensure that the next thing written is not a comment.
	b.numCurrencies = (len(b.currencies) / 4) - 2
	w.WriteConst("numCurrencies", b.numCurrencies)

	// Create a table that maps regions to currencies.
	regionToCurrency := []toCurrency***REMOVED******REMOVED***

	for _, reg := range data.CurrencyData.Region ***REMOVED***
		if len(reg.Iso3166) != 2 ***REMOVED***
			log.Fatalf("Unexpected group %q in region data", reg.Iso3166)
		***REMOVED***
		if len(reg.Currency) == 0 ***REMOVED***
			continue
		***REMOVED***
		cur := reg.Currency[0]
		if cur.To != "" || cur.Tender == "false" ***REMOVED***
			continue
		***REMOVED***
		regionToCurrency = append(regionToCurrency, toCurrency***REMOVED***
			region: regionToCode(language.MustParseRegion(reg.Iso3166)),
			code:   uint16(b.currencies.Index([]byte(cur.Iso4217))),
		***REMOVED***)
	***REMOVED***
	sort.Sort(byRegion(regionToCurrency))

	w.WriteType(toCurrency***REMOVED******REMOVED***)
	w.WriteVar("regionToCurrency", regionToCurrency)

	// Create a table that maps regions to currencies.
	regionData := []regionInfo***REMOVED******REMOVED***

	for _, reg := range data.CurrencyData.Region ***REMOVED***
		if len(reg.Iso3166) != 2 ***REMOVED***
			log.Fatalf("Unexpected group %q in region data", reg.Iso3166)
		***REMOVED***
		for _, cur := range reg.Currency ***REMOVED***
			from, _ := time.Parse("2006-01-02", cur.From)
			to, _ := time.Parse("2006-01-02", cur.To)
			code := uint16(b.currencies.Index([]byte(cur.Iso4217)))
			if cur.Tender == "false" ***REMOVED***
				code |= nonTenderBit
			***REMOVED***
			regionData = append(regionData, regionInfo***REMOVED***
				region: regionToCode(language.MustParseRegion(reg.Iso3166)),
				code:   code,
				from:   toDate(from),
				to:     toDate(to),
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	sort.Stable(byRegionCode(regionData))

	w.WriteType(regionInfo***REMOVED******REMOVED***)
	w.WriteVar("regionData", regionData)
***REMOVED***

type regionInfo struct ***REMOVED***
	region uint16
	code   uint16 // 0x8000 not legal tender
	from   uint32
	to     uint32
***REMOVED***

type byRegionCode []regionInfo

func (a byRegionCode) Len() int           ***REMOVED*** return len(a) ***REMOVED***
func (a byRegionCode) Swap(i, j int)      ***REMOVED*** a[i], a[j] = a[j], a[i] ***REMOVED***
func (a byRegionCode) Less(i, j int) bool ***REMOVED*** return a[i].region < a[j].region ***REMOVED***

type toCurrency struct ***REMOVED***
	region uint16
	code   uint16
***REMOVED***

type byRegion []toCurrency

func (a byRegion) Len() int           ***REMOVED*** return len(a) ***REMOVED***
func (a byRegion) Swap(i, j int)      ***REMOVED*** a[i], a[j] = a[j], a[i] ***REMOVED***
func (a byRegion) Less(i, j int) bool ***REMOVED*** return a[i].region < a[j].region ***REMOVED***

func mkCurrencyInfo(standard, cash int) string ***REMOVED***
	return string([]byte***REMOVED***byte(cash<<cashShift | standard)***REMOVED***)
***REMOVED***

func getRoundingIndex(digits, rounding string, defIndex int) int ***REMOVED***
	round := roundings[defIndex] // default

	if digits != "" ***REMOVED***
		round.scale = parseUint8(digits)
	***REMOVED***
	if rounding != "" && rounding != "0" ***REMOVED*** // 0 means 1 here in CLDR
		round.increment = parseUint8(rounding)
	***REMOVED***

	// Will panic if the entry doesn't exist:
	for i, r := range roundings ***REMOVED***
		if r == round ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	log.Fatalf("Rounding entry %#v does not exist.", round)
	panic("unreachable")
***REMOVED***

// genSymbols generates the symbols used for currencies. Most symbols are
// defined in root and there is only very small variation per language.
// The following rules apply:
// - A symbol can be requested as normal or narrow.
// - If a symbol is not defined for a currency, it defaults to its ISO code.
func (b *builder) genSymbols(w *gen.CodeWriter, data *cldr.CLDR) ***REMOVED***
	d, err := cldr.ParseDraft(*draft)
	if err != nil ***REMOVED***
		log.Fatalf("filter: %v", err)
	***REMOVED***

	const (
		normal = iota
		narrow
		numTypes
	)
	// language -> currency -> type ->  symbol
	var symbols [language.NumCompactTags][][numTypes]*string

	// Collect symbol information per language.
	for _, lang := range data.Locales() ***REMOVED***
		ldml := data.RawLDML(lang)
		if ldml.Numbers == nil || ldml.Numbers.Currencies == nil ***REMOVED***
			continue
		***REMOVED***

		langIndex, ok := language.CompactIndex(language.MustParse(lang))
		if !ok ***REMOVED***
			log.Fatalf("No compact index for language %s", lang)
		***REMOVED***

		symbols[langIndex] = make([][numTypes]*string, b.numCurrencies+1)

		for _, c := range ldml.Numbers.Currencies.Currency ***REMOVED***
			syms := cldr.MakeSlice(&c.Symbol)
			syms.SelectDraft(d)

			for _, sym := range c.Symbol ***REMOVED***
				v := sym.Data()
				if v == c.Type ***REMOVED***
					// We define "" to mean the ISO symbol.
					v = ""
				***REMOVED***
				cur := b.currencies.Index([]byte(c.Type))
				// XXX gets reassigned to 0 in the package's code.
				if c.Type == "XXX" ***REMOVED***
					cur = 0
				***REMOVED***
				if cur == -1 ***REMOVED***
					fmt.Println("Unsupported:", c.Type)
					continue
				***REMOVED***

				switch sym.Alt ***REMOVED***
				case "":
					symbols[langIndex][cur][normal] = &v
				case "narrow":
					symbols[langIndex][cur][narrow] = &v
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Remove values identical to the parent.
	for langIndex, data := range symbols ***REMOVED***
		for curIndex, curs := range data ***REMOVED***
			for typ, sym := range curs ***REMOVED***
				if sym == nil ***REMOVED***
					continue
				***REMOVED***
				for p := uint16(langIndex); p != 0; ***REMOVED***
					p = internal.Parent[p]
					x := symbols[p]
					if x == nil ***REMOVED***
						continue
					***REMOVED***
					if v := x[curIndex][typ]; v != nil || p == 0 ***REMOVED***
						// Value is equal to the default value root value is undefined.
						parentSym := ""
						if v != nil ***REMOVED***
							parentSym = *v
						***REMOVED***
						if parentSym == *sym ***REMOVED***
							// Value is the same as parent.
							data[curIndex][typ] = nil
						***REMOVED***
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Create symbol index.
	symbolData := []byte***REMOVED***0***REMOVED***
	symbolLookup := map[string]uint16***REMOVED***"": 0***REMOVED*** // 0 means default, so block that value.
	for _, data := range symbols ***REMOVED***
		for _, curs := range data ***REMOVED***
			for _, sym := range curs ***REMOVED***
				if sym == nil ***REMOVED***
					continue
				***REMOVED***
				if _, ok := symbolLookup[*sym]; !ok ***REMOVED***
					symbolLookup[*sym] = uint16(len(symbolData))
					symbolData = append(symbolData, byte(len(*sym)))
					symbolData = append(symbolData, *sym...)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	w.WriteComment(`
	symbols holds symbol data of the form <n> <str>, where n is the length of
	the symbol string str.`)
	w.WriteConst("symbols", string(symbolData))

	// Create index from language to currency lookup to symbol.
	type curToIndex struct***REMOVED*** cur, idx uint16 ***REMOVED***
	w.WriteType(curToIndex***REMOVED******REMOVED***)

	prefix := []string***REMOVED***"normal", "narrow"***REMOVED***
	// Create data for regular and narrow symbol data.
	for typ := normal; typ <= narrow; typ++ ***REMOVED***

		indexes := []curToIndex***REMOVED******REMOVED*** // maps currency to symbol index
		languages := []uint16***REMOVED******REMOVED***

		for _, data := range symbols ***REMOVED***
			languages = append(languages, uint16(len(indexes)))
			for curIndex, curs := range data ***REMOVED***

				if sym := curs[typ]; sym != nil ***REMOVED***
					indexes = append(indexes, curToIndex***REMOVED***uint16(curIndex), symbolLookup[*sym]***REMOVED***)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		languages = append(languages, uint16(len(indexes)))

		w.WriteVar(prefix[typ]+"LangIndex", languages)
		w.WriteVar(prefix[typ]+"SymIndex", indexes)
	***REMOVED***
***REMOVED***
func parseUint8(str string) uint8 ***REMOVED***
	x, err := strconv.ParseUint(str, 10, 8)
	if err != nil ***REMOVED***
		// Show line number of where this function was called.
		log.New(os.Stderr, "", log.Lshortfile).Output(2, err.Error())
		os.Exit(1)
	***REMOVED***
	return uint8(x)
***REMOVED***
