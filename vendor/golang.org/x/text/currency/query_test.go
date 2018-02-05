// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package currency

import (
	"testing"
	"time"

	"golang.org/x/text/language"
)

func TestQuery(t *testing.T) ***REMOVED***
	r := func(region string) language.Region ***REMOVED***
		return language.MustParseRegion(region)
	***REMOVED***
	t1800, _ := time.Parse("2006-01-02", "1800-01-01")
	type result struct ***REMOVED***
		region   language.Region
		unit     Unit
		isTender bool
		from, to string
	***REMOVED***
	testCases := []struct ***REMOVED***
		name    string
		opts    []QueryOption
		results []result
	***REMOVED******REMOVED******REMOVED***
		name:    "XA",
		opts:    []QueryOption***REMOVED***Region(r("XA"))***REMOVED***,
		results: []result***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		name: "AC",
		opts: []QueryOption***REMOVED***Region(r("AC"))***REMOVED***,
		results: []result***REMOVED***
			***REMOVED***r("AC"), MustParseISO("SHP"), true, "1976-01-01", ""***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		name: "US",
		opts: []QueryOption***REMOVED***Region(r("US"))***REMOVED***,
		results: []result***REMOVED***
			***REMOVED***r("US"), MustParseISO("USD"), true, "1792-01-01", ""***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		name: "US-hist",
		opts: []QueryOption***REMOVED***Region(r("US")), Historical***REMOVED***,
		results: []result***REMOVED***
			***REMOVED***r("US"), MustParseISO("USD"), true, "1792-01-01", ""***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		name: "US-non-tender",
		opts: []QueryOption***REMOVED***Region(r("US")), NonTender***REMOVED***,
		results: []result***REMOVED***
			***REMOVED***r("US"), MustParseISO("USD"), true, "1792-01-01", ""***REMOVED***,
			***REMOVED***r("US"), MustParseISO("USN"), false, "", ""***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		name: "US-historical+non-tender",
		opts: []QueryOption***REMOVED***Region(r("US")), Historical, NonTender***REMOVED***,
		results: []result***REMOVED***
			***REMOVED***r("US"), MustParseISO("USD"), true, "1792-01-01", ""***REMOVED***,
			***REMOVED***r("US"), MustParseISO("USN"), false, "", ""***REMOVED***,
			***REMOVED***r("US"), MustParseISO("USS"), false, "", "2014-03-01"***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		name: "1800",
		opts: []QueryOption***REMOVED***Date(t1800)***REMOVED***,
		results: []result***REMOVED***
			***REMOVED***r("CH"), MustParseISO("CHF"), true, "1799-03-17", ""***REMOVED***,
			***REMOVED***r("GB"), MustParseISO("GBP"), true, "1694-07-27", ""***REMOVED***,
			***REMOVED***r("GI"), MustParseISO("GIP"), true, "1713-01-01", ""***REMOVED***,
			// The date for IE and PR seem wrong, so these may be updated at
			// some point causing the tests to fail.
			***REMOVED***r("IE"), MustParseISO("GBP"), true, "1800-01-01", "1922-01-01"***REMOVED***,
			***REMOVED***r("PR"), MustParseISO("ESP"), true, "1800-01-01", "1898-12-10"***REMOVED***,
			***REMOVED***r("US"), MustParseISO("USD"), true, "1792-01-01", ""***REMOVED***,
		***REMOVED***,
	***REMOVED******REMOVED***
	for _, tc := range testCases ***REMOVED***
		n := 0
		for it := Query(tc.opts...); it.Next(); n++ ***REMOVED***
			if n < len(tc.results) ***REMOVED***
				got := result***REMOVED***
					it.Region(),
					it.Unit(),
					it.IsTender(),
					getTime(it.From()),
					getTime(it.To()),
				***REMOVED***
				if got != tc.results[n] ***REMOVED***
					t.Errorf("%s:%d: got %v; want %v", tc.name, n, got, tc.results[n])
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if n != len(tc.results) ***REMOVED***
			t.Errorf("%s: unexpected number of results: got %d; want %d", tc.name, n, len(tc.results))
		***REMOVED***
	***REMOVED***
***REMOVED***

func getTime(t time.Time, ok bool) string ***REMOVED***
	if !ok ***REMOVED***
		return ""
	***REMOVED***
	return t.Format("2006-01-02")
***REMOVED***
