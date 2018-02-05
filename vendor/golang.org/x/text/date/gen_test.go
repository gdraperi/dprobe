// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"strconv"
	"strings"
	"testing"

	"golang.org/x/text/internal/cldrtree"
	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

func TestTables(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder***REMOVED******REMOVED***
	d.SetDirFilter("supplemental", "main")
	d.SetSectionFilter("dates")
	data, err := d.DecodeZip(r)
	if err != nil ***REMOVED***
		t.Fatalf("DecodeZip: %v", err)
	***REMOVED***

	count := 0
	for _, lang := range data.Locales() ***REMOVED***
		ldml := data.RawLDML(lang)
		if ldml.Dates == nil ***REMOVED***
			continue
		***REMOVED***
		tag, _ := language.CompactIndex(language.MustParse(lang))

		test := func(want cldrtree.Element, path ...string) ***REMOVED***
			if count > 30 ***REMOVED***
				return
			***REMOVED***
			t.Run(lang+"/"+strings.Join(path, "/"), func(t *testing.T) ***REMOVED***
				p := make([]uint16, len(path))
				for i, s := range path ***REMOVED***
					if v, err := strconv.Atoi(s); err == nil ***REMOVED***
						p[i] = uint16(v)
					***REMOVED*** else if v, ok := enumMap[s]; ok ***REMOVED***
						p[i] = v
					***REMOVED*** else ***REMOVED***
						count++
						t.Fatalf("Unknown key %q", s)
					***REMOVED***
				***REMOVED***
				wantStr := want.GetCommon().Data()
				if got := tree.Lookup(tag, p...); got != wantStr ***REMOVED***
					count++
					t.Errorf("got %q; want %q", got, wantStr)
				***REMOVED***
			***REMOVED***)
		***REMOVED***

		width := func(s string) string ***REMOVED*** return "width" + strings.Title(s) ***REMOVED***

		if ldml.Dates.Calendars != nil ***REMOVED***
			for _, cal := range ldml.Dates.Calendars.Calendar ***REMOVED***
				if cal.Months != nil ***REMOVED***
					for _, mc := range cal.Months.MonthContext ***REMOVED***
						for _, mw := range mc.MonthWidth ***REMOVED***
							for _, m := range mw.Month ***REMOVED***
								test(m, "calendars", cal.Type, "months", mc.Type, width(mw.Type), m.Yeartype+m.Type)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if cal.MonthPatterns != nil ***REMOVED***
					for _, mc := range cal.MonthPatterns.MonthPatternContext ***REMOVED***
						for _, mw := range mc.MonthPatternWidth ***REMOVED***
							for _, m := range mw.MonthPattern ***REMOVED***
								test(m, "calendars", cal.Type, "monthPatterns", mc.Type, width(mw.Type))
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if cal.CyclicNameSets != nil ***REMOVED***
					for _, cns := range cal.CyclicNameSets.CyclicNameSet ***REMOVED***
						for _, cc := range cns.CyclicNameContext ***REMOVED***
							for _, cw := range cc.CyclicNameWidth ***REMOVED***
								for _, c := range cw.CyclicName ***REMOVED***
									test(c, "calendars", cal.Type, "cyclicNameSets", cns.Type+"CycleType", cc.Type, width(cw.Type), c.Type)

								***REMOVED***
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if cal.Days != nil ***REMOVED***
					for _, dc := range cal.Days.DayContext ***REMOVED***
						for _, dw := range dc.DayWidth ***REMOVED***
							for _, d := range dw.Day ***REMOVED***
								test(d, "calendars", cal.Type, "days", dc.Type, width(dw.Type), d.Type)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if cal.Quarters != nil ***REMOVED***
					for _, qc := range cal.Quarters.QuarterContext ***REMOVED***
						for _, qw := range qc.QuarterWidth ***REMOVED***
							for _, q := range qw.Quarter ***REMOVED***
								test(q, "calendars", cal.Type, "quarters", qc.Type, width(qw.Type), q.Type)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if cal.DayPeriods != nil ***REMOVED***
					for _, dc := range cal.DayPeriods.DayPeriodContext ***REMOVED***
						for _, dw := range dc.DayPeriodWidth ***REMOVED***
							for _, d := range dw.DayPeriod ***REMOVED***
								test(d, "calendars", cal.Type, "dayPeriods", dc.Type, width(dw.Type), d.Type, d.Alt)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if cal.Eras != nil ***REMOVED***
					if cal.Eras.EraNames != nil ***REMOVED***
						for _, e := range cal.Eras.EraNames.Era ***REMOVED***
							test(e, "calendars", cal.Type, "eras", "widthWide", e.Alt, e.Type)
						***REMOVED***
					***REMOVED***
					if cal.Eras.EraAbbr != nil ***REMOVED***
						for _, e := range cal.Eras.EraAbbr.Era ***REMOVED***
							test(e, "calendars", cal.Type, "eras", "widthAbbreviated", e.Alt, e.Type)
						***REMOVED***
					***REMOVED***
					if cal.Eras.EraNarrow != nil ***REMOVED***
						for _, e := range cal.Eras.EraNarrow.Era ***REMOVED***
							test(e, "calendars", cal.Type, "eras", "widthNarrow", e.Alt, e.Type)
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if cal.DateFormats != nil ***REMOVED***
					for _, dfl := range cal.DateFormats.DateFormatLength ***REMOVED***
						for _, df := range dfl.DateFormat ***REMOVED***
							for _, p := range df.Pattern ***REMOVED***
								test(p, "calendars", cal.Type, "dateFormats", dfl.Type, p.Alt)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if cal.TimeFormats != nil ***REMOVED***
					for _, tfl := range cal.TimeFormats.TimeFormatLength ***REMOVED***
						for _, tf := range tfl.TimeFormat ***REMOVED***
							for _, p := range tf.Pattern ***REMOVED***
								test(p, "calendars", cal.Type, "timeFormats", tfl.Type, p.Alt)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if cal.DateTimeFormats != nil ***REMOVED***
					for _, dtfl := range cal.DateTimeFormats.DateTimeFormatLength ***REMOVED***
						for _, dtf := range dtfl.DateTimeFormat ***REMOVED***
							for _, p := range dtf.Pattern ***REMOVED***
								test(p, "calendars", cal.Type, "dateTimeFormats", dtfl.Type, p.Alt)
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
		if ldml.Dates.Fields != nil ***REMOVED***
			for _, f := range ldml.Dates.Fields.Field ***REMOVED***
				field := f.Type + "Field"
				for _, d := range f.DisplayName ***REMOVED***
					test(d, "fields", field, "displayName", d.Alt)
				***REMOVED***
				for _, r := range f.Relative ***REMOVED***
					i, _ := strconv.Atoi(r.Type)
					v := []string***REMOVED***"before2", "before1", "current", "after1", "after2", "after3"***REMOVED***[i+2]
					test(r, "fields", field, "relative", v)
				***REMOVED***
				for _, rt := range f.RelativeTime ***REMOVED***
					for _, p := range rt.RelativeTimePattern ***REMOVED***
						test(p, "fields", field, "relativeTime", rt.Type, p.Count)
					***REMOVED***
				***REMOVED***
				for _, rp := range f.RelativePeriod ***REMOVED***
					test(rp, "fields", field, "relativePeriod", rp.Alt)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if ldml.Dates.TimeZoneNames != nil ***REMOVED***
			for _, h := range ldml.Dates.TimeZoneNames.HourFormat ***REMOVED***
				test(h, "timeZoneNames", "zoneFormat", h.Element())
			***REMOVED***
			for _, g := range ldml.Dates.TimeZoneNames.GmtFormat ***REMOVED***
				test(g, "timeZoneNames", "zoneFormat", g.Element())
			***REMOVED***
			for _, g := range ldml.Dates.TimeZoneNames.GmtZeroFormat ***REMOVED***
				test(g, "timeZoneNames", "zoneFormat", g.Element())
			***REMOVED***
			for _, r := range ldml.Dates.TimeZoneNames.RegionFormat ***REMOVED***
				s := r.Type
				if s == "" ***REMOVED***
					s = "generic"
				***REMOVED***
				test(r, "timeZoneNames", "regionFormat", s+"Time")
			***REMOVED***

			testZone := func(zoneType, zoneWidth, zone string, a ...[]*cldr.Common) ***REMOVED***
				for _, e := range a ***REMOVED***
					for _, n := range e ***REMOVED***
						test(n, "timeZoneNames", zoneType, zoneWidth, n.Element()+"Time", zone)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			for _, z := range ldml.Dates.TimeZoneNames.Zone ***REMOVED***
				for _, l := range z.Long ***REMOVED***
					testZone("zone", l.Element(), z.Type, l.Generic, l.Standard, l.Daylight)
				***REMOVED***
				for _, l := range z.Short ***REMOVED***
					testZone("zone", l.Element(), z.Type, l.Generic, l.Standard, l.Daylight)
				***REMOVED***
			***REMOVED***
			for _, z := range ldml.Dates.TimeZoneNames.Metazone ***REMOVED***
				for _, l := range z.Long ***REMOVED***
					testZone("metaZone", l.Element(), z.Type, l.Generic, l.Standard, l.Daylight)
				***REMOVED***
				for _, l := range z.Short ***REMOVED***
					testZone("metaZone", l.Element(), z.Type, l.Generic, l.Standard, l.Daylight)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
