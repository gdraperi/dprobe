package currency

import (
	"flag"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/unicode/cldr"
)

var draft = flag.String("draft",
	"contributed",
	`Minimal draft requirements (approved, contributed, provisional, unconfirmed).`)

func TestTables(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	// Read the CLDR zip file.
	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder***REMOVED******REMOVED***
	d.SetDirFilter("supplemental", "main")
	d.SetSectionFilter("numbers")
	data, err := d.DecodeZip(r)
	if err != nil ***REMOVED***
		t.Fatalf("DecodeZip: %v", err)
	***REMOVED***

	dr, err := cldr.ParseDraft(*draft)
	if err != nil ***REMOVED***
		t.Fatalf("filter: %v", err)
	***REMOVED***

	for _, lang := range data.Locales() ***REMOVED***
		p := message.NewPrinter(language.MustParse(lang))

		ldml := data.RawLDML(lang)
		if ldml.Numbers == nil || ldml.Numbers.Currencies == nil ***REMOVED***
			continue
		***REMOVED***
		for _, c := range ldml.Numbers.Currencies.Currency ***REMOVED***
			syms := cldr.MakeSlice(&c.Symbol)
			syms.SelectDraft(dr)

			for _, sym := range c.Symbol ***REMOVED***
				cur, err := ParseISO(c.Type)
				if err != nil ***REMOVED***
					continue
				***REMOVED***
				formatter := Symbol
				switch sym.Alt ***REMOVED***
				case "":
				case "narrow":
					formatter = NarrowSymbol
				default:
					continue
				***REMOVED***
				want := sym.Data()
				if got := p.Sprint(formatter(cur)); got != want ***REMOVED***
					t.Errorf("%s:%sSymbol(%s) = %s; want %s", lang, strings.Title(sym.Alt), c.Type, got, want)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, reg := range data.Supplemental().CurrencyData.Region ***REMOVED***
		i := 0
		for ; regionData[i].Region().String() != reg.Iso3166; i++ ***REMOVED***
		***REMOVED***
		it := Query(Historical, NonTender, Region(language.MustParseRegion(reg.Iso3166)))
		for _, cur := range reg.Currency ***REMOVED***
			from, _ := time.Parse("2006-01-02", cur.From)
			to, _ := time.Parse("2006-01-02", cur.To)

			it.Next()
			for j, r := range []QueryIter***REMOVED***&iter***REMOVED***regionInfo: &regionData[i]***REMOVED***, it***REMOVED*** ***REMOVED***
				if got, _ := r.From(); from != got ***REMOVED***
					t.Errorf("%d:%s:%s:from: got %v; want %v", j, reg.Iso3166, cur.Iso4217, got, from)
				***REMOVED***
				if got, _ := r.To(); to != got ***REMOVED***
					t.Errorf("%d:%s:%s:to: got %v; want %v", j, reg.Iso3166, cur.Iso4217, got, to)
				***REMOVED***
			***REMOVED***
			i++
		***REMOVED***
	***REMOVED***
***REMOVED***
