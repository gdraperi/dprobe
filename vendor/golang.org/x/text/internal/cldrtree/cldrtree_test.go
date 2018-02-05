// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldrtree

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

var genOutput = flag.Bool("gen", false, "generate output files")

func TestAliasRegexp(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		alias string
		want  []string
	***REMOVED******REMOVED******REMOVED***
		alias: "miscPatterns[@numberSystem='latn']",
		want: []string***REMOVED***
			"miscPatterns[@numberSystem='latn']",
			"miscPatterns",
			"[@numberSystem='latn']",
			"numberSystem",
			"latn",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		alias: `calendar[@type='greg-foo']/days/`,
		want: []string***REMOVED***
			"calendar[@type='greg-foo']",
			"calendar",
			"[@type='greg-foo']",
			"type",
			"greg-foo",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		alias: "eraAbbr",
		want: []string***REMOVED***
			"eraAbbr",
			"eraAbbr",
			"",
			"",
			"",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		// match must be anchored at beginning.
		alias: `../calendar[@type='gregorian']/days/`,
	***REMOVED******REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(tc.alias, func(t *testing.T) ***REMOVED***
			got := aliasRe.FindStringSubmatch(tc.alias)
			if !reflect.DeepEqual(got, tc.want) ***REMOVED***
				t.Errorf("got %v; want %v", got, tc.want)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestBuild(t *testing.T) ***REMOVED***
	tree1, _ := loadTestdata(t, "test1")
	tree2, _ := loadTestdata(t, "test2")

	// Constants for second test test
	const (
		calendar = iota
		field
	)
	const (
		month = iota
		era
		filler
		cyclicNameSet
	)
	const (
		abbreviated = iota
		narrow
		wide
	)

	testCases := []struct ***REMOVED***
		desc      string
		tree      *Tree
		locale    string
		path      []uint16
		isFeature bool
		result    string
	***REMOVED******REMOVED******REMOVED***
		desc:   "und/chinese month format wide m1",
		tree:   tree1,
		locale: "und",
		path:   path(calendar, 0, month, 0, wide, 1),
		result: "cM01",
	***REMOVED***, ***REMOVED***
		desc:   "und/chinese month format wide m12",
		tree:   tree1,
		locale: "und",
		path:   path(calendar, 0, month, 0, wide, 12),
		result: "cM12",
	***REMOVED***, ***REMOVED***
		desc:   "und/non-existing value",
		tree:   tree1,
		locale: "und",
		path:   path(calendar, 0, month, 0, wide, 13),
		result: "",
	***REMOVED***, ***REMOVED***
		desc:   "und/dangi:chinese month format wide",
		tree:   tree1,
		locale: "und",
		path:   path(calendar, 1, month, 0, wide, 1),
		result: "cM01",
	***REMOVED***, ***REMOVED***
		desc:   "und/chinese month format abbreviated:wide",
		tree:   tree1,
		locale: "und",
		path:   path(calendar, 0, month, 0, abbreviated, 1),
		result: "cM01",
	***REMOVED***, ***REMOVED***
		desc:   "und/chinese month format narrow:wide",
		tree:   tree1,
		locale: "und",
		path:   path(calendar, 0, month, 0, narrow, 1),
		result: "cM01",
	***REMOVED***, ***REMOVED***
		desc:   "und/gregorian month format wide",
		tree:   tree1,
		locale: "und",
		path:   path(calendar, 2, month, 0, wide, 2),
		result: "gM02",
	***REMOVED***, ***REMOVED***
		desc:   "und/gregorian month format:stand-alone narrow",
		tree:   tree1,
		locale: "und",
		path:   path(calendar, 2, month, 0, narrow, 1),
		result: "1",
	***REMOVED***, ***REMOVED***
		desc:   "und/gregorian month stand-alone:format abbreviated",
		tree:   tree1,
		locale: "und",
		path:   path(calendar, 2, month, 1, abbreviated, 1),
		result: "gM01",
	***REMOVED***, ***REMOVED***
		desc:   "und/gregorian month stand-alone:format wide ",
		tree:   tree1,
		locale: "und",
		path:   path(calendar, 2, month, 1, abbreviated, 1),
		result: "gM01",
	***REMOVED***, ***REMOVED***
		desc:   "und/dangi:chinese month format narrow:wide ",
		tree:   tree1,
		locale: "und",
		path:   path(calendar, 1, month, 0, narrow, 4),
		result: "cM04",
	***REMOVED***, ***REMOVED***
		desc:   "und/field era displayname 0",
		tree:   tree2,
		locale: "und",
		path:   path(field, 0, 0, 0),
		result: "Era",
	***REMOVED***, ***REMOVED***
		desc:   "en/field era displayname 0",
		tree:   tree2,
		locale: "en",
		path:   path(field, 0, 0, 0),
		result: "era",
	***REMOVED***, ***REMOVED***
		desc:   "und/calendar hebrew format wide 7-leap",
		tree:   tree2,
		locale: "und",
		path:   path(calendar, 7, month, 0, wide, 0),
		result: "Adar II",
	***REMOVED***, ***REMOVED***
		desc:   "en-GB:en-001:en:und/calendar hebrew format wide 7-leap",
		tree:   tree2,
		locale: "en-GB",
		path:   path(calendar, 7, month, 0, wide, 0),
		result: "Adar II",
	***REMOVED***, ***REMOVED***
		desc:   "und/buddhist month format wide 11",
		tree:   tree2,
		locale: "und",
		path:   path(calendar, 0, month, 0, wide, 12),
		result: "genWideM12",
	***REMOVED***, ***REMOVED***
		desc:   "en-GB/gregorian month stand-alone narrow 2",
		tree:   tree2,
		locale: "en-GB",
		path:   path(calendar, 6, month, 1, narrow, 3),
		result: "gbNarrowM3",
	***REMOVED***, ***REMOVED***
		desc:   "en-GB/gregorian month format narrow 3/missing in en-GB",
		tree:   tree2,
		locale: "en-GB",
		path:   path(calendar, 6, month, 0, narrow, 4),
		result: "enNarrowM4",
	***REMOVED***, ***REMOVED***
		desc:   "en-GB/gregorian month format narrow 3/missing in en and en-GB",
		tree:   tree2,
		locale: "en-GB",
		path:   path(calendar, 6, month, 0, narrow, 7),
		result: "gregNarrowM7",
	***REMOVED***, ***REMOVED***
		desc:   "en-GB/gregorian month format narrow 3/missing in en and en-GB",
		tree:   tree2,
		locale: "en-GB",
		path:   path(calendar, 6, month, 0, narrow, 7),
		result: "gregNarrowM7",
	***REMOVED***, ***REMOVED***
		desc:      "en-GB/gregorian era narrow",
		tree:      tree2,
		locale:    "en-GB",
		path:      path(calendar, 6, era, abbreviated, 0, 1),
		isFeature: true,
		result:    "AD",
	***REMOVED***, ***REMOVED***
		desc:      "en-GB/gregorian era narrow",
		tree:      tree2,
		locale:    "en-GB",
		path:      path(calendar, 6, era, narrow, 0, 0),
		isFeature: true,
		result:    "BC",
	***REMOVED***, ***REMOVED***
		desc:      "en-GB/gregorian era narrow",
		tree:      tree2,
		locale:    "en-GB",
		path:      path(calendar, 6, era, wide, 1, 0),
		isFeature: true,
		result:    "Before Common Era",
	***REMOVED***, ***REMOVED***
		desc:      "en-GB/dangi:chinese cyclicName, months, format, narrow:abbreviated 2",
		tree:      tree2,
		locale:    "en-GB",
		path:      path(calendar, 1, cyclicNameSet, 3, 0, 1, 2),
		isFeature: true,
		result:    "year2",
	***REMOVED***, ***REMOVED***
		desc:   "en-GB/field era-narrow ",
		tree:   tree2,
		locale: "en-GB",
		path:   path(field, 2, 0, 0),
		result: "era",
	***REMOVED***, ***REMOVED***
		desc:      "en-GB/field month-narrow relativeTime future one",
		tree:      tree2,
		locale:    "en-GB",
		path:      path(field, 5, 2, 0, 1),
		isFeature: true,
		result:    "001NarrowFutMOne",
	***REMOVED***, ***REMOVED***
		// Don't fall back to the one of "en".
		desc:      "en-GB/field month-short relativeTime past one:other",
		tree:      tree2,
		locale:    "en-GB",
		path:      path(field, 4, 2, 1, 1),
		isFeature: true,
		result:    "001ShortPastMOther",
	***REMOVED***, ***REMOVED***
		desc:      "en-GB/field month relativeTime future two:other",
		tree:      tree2,
		locale:    "en-GB",
		path:      path(field, 3, 2, 0, 2),
		isFeature: true,
		result:    "enFutMOther",
	***REMOVED******REMOVED***

	for _, tc := range testCases ***REMOVED***
		t.Run(tc.desc, func(t *testing.T) ***REMOVED***
			tag, _ := language.CompactIndex(language.MustParse(tc.locale))
			s := tc.tree.lookup(tag, tc.isFeature, tc.path...)
			if s != tc.result ***REMOVED***
				t.Errorf("got %q; want %q", s, tc.result)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func path(e ...uint16) []uint16 ***REMOVED*** return e ***REMOVED***

func TestGen(t *testing.T) ***REMOVED***
	testCases := []string***REMOVED***"test1", "test2"***REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(tc, func(t *testing.T) ***REMOVED***
			_, got := loadTestdata(t, tc)

			// Remove sizes that may vary per architecture.
			re := regexp.MustCompile("// Size: [0-9]*")
			got = re.ReplaceAllLiteral(got, []byte("// Size: xxxx"))
			re = regexp.MustCompile("// Total table size [0-9]*")
			got = re.ReplaceAllLiteral(got, []byte("// Total table size: xxxx"))

			file := filepath.Join("testdata", tc, "output.go")
			if *genOutput ***REMOVED***
				ioutil.WriteFile(file, got, 0700)
				t.SkipNow()
			***REMOVED***

			b, err := ioutil.ReadFile(file)
			if err != nil ***REMOVED***
				t.Fatalf("failed to open file: %v", err)
			***REMOVED***
			if want := string(b); string(got) != want ***REMOVED***
				t.Log(string(got))
				t.Errorf("files differ")
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func loadTestdata(t *testing.T, test string) (tree *Tree, file []byte) ***REMOVED***
	b := New("test")

	var d cldr.Decoder

	data, err := d.DecodePath(filepath.Join("testdata", test))
	if err != nil ***REMOVED***
		t.Fatalf("error decoding testdata: %v", err)
	***REMOVED***

	context := Enum("context")
	widthMap := func(s string) string ***REMOVED***
		// Align era with width values.
		if r, ok := map[string]string***REMOVED***
			"eraAbbr":   "abbreviated",
			"eraNarrow": "narrow",
			"eraNames":  "wide",
		***REMOVED***[s]; ok ***REMOVED***
			s = r
		***REMOVED***
		return "w" + strings.Title(s)
	***REMOVED***
	width := EnumFunc("width", widthMap, "abbreviated", "narrow", "wide")
	month := Enum("month", "leap7")
	relative := EnumFunc("relative", func(s string) string ***REMOVED***
		x, err := strconv.ParseInt(s, 10, 8)
		if err != nil ***REMOVED***
			log.Fatal("Invalid number:", err)
		***REMOVED***
		return []string***REMOVED***
			"before1",
			"current",
			"after1",
		***REMOVED***[x+1]
	***REMOVED***)
	cycleType := EnumFunc("cycleType", func(s string) string ***REMOVED***
		return "cyc" + strings.Title(s)
	***REMOVED***)
	r := rand.New(rand.NewSource(0))

	for _, loc := range data.Locales() ***REMOVED***
		ldml := data.RawLDML(loc)
		x := b.Locale(language.Make(loc))

		if x := x.Index(ldml.Dates.Calendars); x != nil ***REMOVED***
			for _, cal := range ldml.Dates.Calendars.Calendar ***REMOVED***
				x := x.IndexFromType(cal)
				if x := x.Index(cal.Months); x != nil ***REMOVED***
					for _, mc := range cal.Months.MonthContext ***REMOVED***
						x := x.IndexFromType(mc, context)
						for _, mw := range mc.MonthWidth ***REMOVED***
							x := x.IndexFromType(mw, width)
							for _, m := range mw.Month ***REMOVED***
								x.SetValue(m.Yeartype+m.Type, m, month)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if x := x.Index(cal.CyclicNameSets); x != nil ***REMOVED***
					for _, cns := range cal.CyclicNameSets.CyclicNameSet ***REMOVED***
						x := x.IndexFromType(cns, cycleType)
						for _, cc := range cns.CyclicNameContext ***REMOVED***
							x := x.IndexFromType(cc, context)
							for _, cw := range cc.CyclicNameWidth ***REMOVED***
								x := x.IndexFromType(cw, width)
								for _, c := range cw.CyclicName ***REMOVED***
									x.SetValue(c.Type, c)
								***REMOVED***
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if x := x.Index(cal.Eras); x != nil ***REMOVED***
					opts := []Option***REMOVED***width, SharedType()***REMOVED***
					if x := x.Index(cal.Eras.EraNames, opts...); x != nil ***REMOVED***
						for _, e := range cal.Eras.EraNames.Era ***REMOVED***
							x.IndexFromAlt(e).SetValue(e.Type, e)
						***REMOVED***
					***REMOVED***
					if x := x.Index(cal.Eras.EraAbbr, opts...); x != nil ***REMOVED***
						for _, e := range cal.Eras.EraAbbr.Era ***REMOVED***
							x.IndexFromAlt(e).SetValue(e.Type, e)
						***REMOVED***
					***REMOVED***
					if x := x.Index(cal.Eras.EraNarrow, opts...); x != nil ***REMOVED***
						for _, e := range cal.Eras.EraNarrow.Era ***REMOVED***
							x.IndexFromAlt(e).SetValue(e.Type, e)
						***REMOVED***
					***REMOVED***
				***REMOVED***
				***REMOVED***
					// Ensure having more than 2 buckets.
					f := x.IndexWithName("filler")
					b := make([]byte, maxStrlen)
					opt := &options***REMOVED***parent: x***REMOVED***
					r.Read(b)
					f.setValue("0", string(b), opt)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if x := x.Index(ldml.Dates.Fields); x != nil ***REMOVED***
			for _, f := range ldml.Dates.Fields.Field ***REMOVED***
				x := x.IndexFromType(f)
				for _, d := range f.DisplayName ***REMOVED***
					x.Index(d).SetValue("", d)
				***REMOVED***
				for _, r := range f.Relative ***REMOVED***
					x.Index(r).SetValue(r.Type, r, relative)
				***REMOVED***
				for _, rt := range f.RelativeTime ***REMOVED***
					x := x.Index(rt).IndexFromType(rt)
					for _, p := range rt.RelativeTimePattern ***REMOVED***
						x.SetValue(p.Count, p)
					***REMOVED***
				***REMOVED***
				for _, rp := range f.RelativePeriod ***REMOVED***
					x.Index(rp).SetValue("", rp)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	tree, err = build(b)
	if err != nil ***REMOVED***
		t.Fatal("error building tree:", err)
	***REMOVED***
	w := gen.NewCodeWriter()
	generate(b, tree, w)
	generateTestData(b, w)
	buf := &bytes.Buffer***REMOVED******REMOVED***
	if _, err = w.WriteGo(buf, "test", ""); err != nil ***REMOVED***
		t.Log(buf.String())
		t.Fatal("error generating code:", err)
	***REMOVED***
	return tree, buf.Bytes()
***REMOVED***
