// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	"golang.org/x/text/internal/cldrtree"
	"golang.org/x/text/internal/gen"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

var (
	draft = flag.String("draft",
		"contributed",
		`Minimal draft requirements (approved, contributed, provisional, unconfirmed).`)
)

// TODO:
// - Compile format patterns.
// - Compress the large amount of redundancy in metazones.
// - Split trees (with shared buckets) with data that is enough for default
//   formatting of Go Time values and and tables that are needed for larger
//   variants.
// - zone to metaZone mappings (in supplemental)
// - Add more enum values and also some key maps for some of the elements.

func main() ***REMOVED***
	gen.Init()

	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder***REMOVED******REMOVED***
	d.SetDirFilter("supplemental", "main")
	d.SetSectionFilter("dates")
	data, err := d.DecodeZip(r)
	if err != nil ***REMOVED***
		log.Fatalf("DecodeZip: %v", err)
	***REMOVED***

	dates := cldrtree.New("dates")
	buildCLDRTree(data, dates)

	w := gen.NewCodeWriter()
	if err := dates.Gen(w); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	gen.WriteCLDRVersion(w)
	w.WriteGoFile("tables.go", "date")

	w = gen.NewCodeWriter()
	if err := dates.GenTestData(w); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	w.WriteGoFile("data_test.go", "date")
***REMOVED***

func buildCLDRTree(data *cldr.CLDR, dates *cldrtree.Builder) ***REMOVED***
	context := cldrtree.Enum("context")
	widthMap := func(s string) string ***REMOVED***
		// Align era with width values.
		if r, ok := map[string]string***REMOVED***
			"eraAbbr":   "abbreviated",
			"eraNarrow": "narrow",
			"eraNames":  "wide",
		***REMOVED***[s]; ok ***REMOVED***
			s = r
		***REMOVED***
		// Prefix width to disambiguate with some overlapping length values.
		return "width" + strings.Title(s)
	***REMOVED***
	width := cldrtree.EnumFunc("width", widthMap, "abbreviated", "narrow", "wide")
	length := cldrtree.Enum("length", "short", "long")
	month := cldrtree.Enum("month", "leap7")
	relTime := cldrtree.EnumFunc("relTime", func(s string) string ***REMOVED***
		x, err := strconv.ParseInt(s, 10, 8)
		if err != nil ***REMOVED***
			log.Fatal("Invalid number:", err)
		***REMOVED***
		return []string***REMOVED***
			"before2",
			"before1",
			"current",
			"after1",
			"after2",
			"after3",
		***REMOVED***[x+2]
	***REMOVED***)
	// Disambiguate keys like 'months' and 'sun'.
	cycleType := cldrtree.EnumFunc("cycleType", func(s string) string ***REMOVED***
		return s + "CycleType"
	***REMOVED***)
	field := cldrtree.EnumFunc("field", func(s string) string ***REMOVED***
		return s + "Field"
	***REMOVED***)
	timeType := cldrtree.EnumFunc("timeType", func(s string) string ***REMOVED***
		if s == "" ***REMOVED***
			return "genericTime"
		***REMOVED***
		return s + "Time"
	***REMOVED***, "generic")

	zoneType := []cldrtree.Option***REMOVED***cldrtree.SharedType(), timeType***REMOVED***
	metaZoneType := []cldrtree.Option***REMOVED***cldrtree.SharedType(), timeType***REMOVED***

	for _, lang := range data.Locales() ***REMOVED***
		tag := language.Make(lang)
		ldml := data.RawLDML(lang)
		if ldml.Dates == nil ***REMOVED***
			continue
		***REMOVED***
		x := dates.Locale(tag)
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
				if x := x.Index(cal.MonthPatterns); x != nil ***REMOVED***
					for _, mc := range cal.MonthPatterns.MonthPatternContext ***REMOVED***
						x := x.IndexFromType(mc, context)
						for _, mw := range mc.MonthPatternWidth ***REMOVED***
							// Value is always leap, so no need to create a
							// subindex.
							for _, m := range mw.MonthPattern ***REMOVED***
								x.SetValue(mw.Type, m, width)
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
				if x := x.Index(cal.Days); x != nil ***REMOVED***
					for _, dc := range cal.Days.DayContext ***REMOVED***
						x := x.IndexFromType(dc, context)
						for _, dw := range dc.DayWidth ***REMOVED***
							x := x.IndexFromType(dw, width)
							for _, d := range dw.Day ***REMOVED***
								x.SetValue(d.Type, d)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if x := x.Index(cal.Quarters); x != nil ***REMOVED***
					for _, qc := range cal.Quarters.QuarterContext ***REMOVED***
						x := x.IndexFromType(qc, context)
						for _, qw := range qc.QuarterWidth ***REMOVED***
							x := x.IndexFromType(qw, width)
							for _, q := range qw.Quarter ***REMOVED***
								x.SetValue(q.Type, q)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if x := x.Index(cal.DayPeriods); x != nil ***REMOVED***
					for _, dc := range cal.DayPeriods.DayPeriodContext ***REMOVED***
						x := x.IndexFromType(dc, context)
						for _, dw := range dc.DayPeriodWidth ***REMOVED***
							x := x.IndexFromType(dw, width)
							for _, d := range dw.DayPeriod ***REMOVED***
								x.IndexFromType(d).SetValue(d.Alt, d)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if x := x.Index(cal.Eras); x != nil ***REMOVED***
					opts := []cldrtree.Option***REMOVED***width, cldrtree.SharedType()***REMOVED***
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
				if x := x.Index(cal.DateFormats); x != nil ***REMOVED***
					for _, dfl := range cal.DateFormats.DateFormatLength ***REMOVED***
						x := x.IndexFromType(dfl, length)
						for _, df := range dfl.DateFormat ***REMOVED***
							for _, p := range df.Pattern ***REMOVED***
								x.SetValue(p.Alt, p)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if x := x.Index(cal.TimeFormats); x != nil ***REMOVED***
					for _, tfl := range cal.TimeFormats.TimeFormatLength ***REMOVED***
						x := x.IndexFromType(tfl, length)
						for _, tf := range tfl.TimeFormat ***REMOVED***
							for _, p := range tf.Pattern ***REMOVED***
								x.SetValue(p.Alt, p)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if x := x.Index(cal.DateTimeFormats); x != nil ***REMOVED***
					for _, dtfl := range cal.DateTimeFormats.DateTimeFormatLength ***REMOVED***
						x := x.IndexFromType(dtfl, length)
						for _, dtf := range dtfl.DateTimeFormat ***REMOVED***
							for _, p := range dtf.Pattern ***REMOVED***
								x.SetValue(p.Alt, p)
							***REMOVED***
						***REMOVED***
					***REMOVED***
					// TODO:
					// - appendItems
					// - intervalFormats
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// TODO: this is a lot of data and is probably relatively little used.
		// Store this somewhere else.
		if x := x.Index(ldml.Dates.Fields); x != nil ***REMOVED***
			for _, f := range ldml.Dates.Fields.Field ***REMOVED***
				x := x.IndexFromType(f, field)
				for _, d := range f.DisplayName ***REMOVED***
					x.Index(d).SetValue(d.Alt, d)
				***REMOVED***
				for _, r := range f.Relative ***REMOVED***
					x.Index(r).SetValue(r.Type, r, relTime)
				***REMOVED***
				for _, rt := range f.RelativeTime ***REMOVED***
					x := x.Index(rt).IndexFromType(rt)
					for _, p := range rt.RelativeTimePattern ***REMOVED***
						x.SetValue(p.Count, p)
					***REMOVED***
				***REMOVED***
				for _, rp := range f.RelativePeriod ***REMOVED***
					x.Index(rp).SetValue(rp.Alt, rp)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if x := x.Index(ldml.Dates.TimeZoneNames); x != nil ***REMOVED***
			format := x.IndexWithName("zoneFormat")
			for _, h := range ldml.Dates.TimeZoneNames.HourFormat ***REMOVED***
				format.SetValue(h.Element(), h)
			***REMOVED***
			for _, g := range ldml.Dates.TimeZoneNames.GmtFormat ***REMOVED***
				format.SetValue(g.Element(), g)
			***REMOVED***
			for _, g := range ldml.Dates.TimeZoneNames.GmtZeroFormat ***REMOVED***
				format.SetValue(g.Element(), g)
			***REMOVED***
			for _, r := range ldml.Dates.TimeZoneNames.RegionFormat ***REMOVED***
				x.Index(r).SetValue(r.Type, r, timeType)
			***REMOVED***

			set := func(x *cldrtree.Index, e []*cldr.Common, zone string) ***REMOVED***
				for _, n := range e ***REMOVED***
					x.Index(n, zoneType...).SetValue(zone, n)
				***REMOVED***
			***REMOVED***
			zoneWidth := []cldrtree.Option***REMOVED***length, cldrtree.SharedType()***REMOVED***
			zs := x.IndexWithName("zone")
			for _, z := range ldml.Dates.TimeZoneNames.Zone ***REMOVED***
				for _, l := range z.Long ***REMOVED***
					x := zs.Index(l, zoneWidth...)
					set(x, l.Generic, z.Type)
					set(x, l.Standard, z.Type)
					set(x, l.Daylight, z.Type)
				***REMOVED***
				for _, s := range z.Short ***REMOVED***
					x := zs.Index(s, zoneWidth...)
					set(x, s.Generic, z.Type)
					set(x, s.Standard, z.Type)
					set(x, s.Daylight, z.Type)
				***REMOVED***
			***REMOVED***
			set = func(x *cldrtree.Index, e []*cldr.Common, zone string) ***REMOVED***
				for _, n := range e ***REMOVED***
					x.Index(n, metaZoneType...).SetValue(zone, n)
				***REMOVED***
			***REMOVED***
			zoneWidth = []cldrtree.Option***REMOVED***length, cldrtree.SharedType()***REMOVED***
			zs = x.IndexWithName("metaZone")
			for _, z := range ldml.Dates.TimeZoneNames.Metazone ***REMOVED***
				for _, l := range z.Long ***REMOVED***
					x := zs.Index(l, zoneWidth...)
					set(x, l.Generic, z.Type)
					set(x, l.Standard, z.Type)
					set(x, l.Daylight, z.Type)
				***REMOVED***
				for _, s := range z.Short ***REMOVED***
					x := zs.Index(s, zoneWidth...)
					set(x, s.Generic, z.Type)
					set(x, s.Standard, z.Type)
					set(x, s.Daylight, z.Type)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
