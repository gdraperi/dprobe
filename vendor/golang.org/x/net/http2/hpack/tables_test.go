// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hpack

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestHeaderFieldTable(t *testing.T) ***REMOVED***
	table := &headerFieldTable***REMOVED******REMOVED***
	table.init()
	table.addEntry(pair("key1", "value1-1"))
	table.addEntry(pair("key2", "value2-1"))
	table.addEntry(pair("key1", "value1-2"))
	table.addEntry(pair("key3", "value3-1"))
	table.addEntry(pair("key4", "value4-1"))
	table.addEntry(pair("key2", "value2-2"))

	// Tests will be run twice: once before evicting anything, and
	// again after evicting the three oldest entries.
	tests := []struct ***REMOVED***
		f                 HeaderField
		beforeWantStaticI uint64
		beforeWantMatch   bool
		afterWantStaticI  uint64
		afterWantMatch    bool
	***REMOVED******REMOVED***
		***REMOVED***HeaderField***REMOVED***"key1", "value1-1", false***REMOVED***, 1, true, 0, false***REMOVED***,
		***REMOVED***HeaderField***REMOVED***"key1", "value1-2", false***REMOVED***, 3, true, 0, false***REMOVED***,
		***REMOVED***HeaderField***REMOVED***"key1", "value1-3", false***REMOVED***, 3, false, 0, false***REMOVED***,
		***REMOVED***HeaderField***REMOVED***"key2", "value2-1", false***REMOVED***, 2, true, 3, false***REMOVED***,
		***REMOVED***HeaderField***REMOVED***"key2", "value2-2", false***REMOVED***, 6, true, 3, true***REMOVED***,
		***REMOVED***HeaderField***REMOVED***"key2", "value2-3", false***REMOVED***, 6, false, 3, false***REMOVED***,
		***REMOVED***HeaderField***REMOVED***"key4", "value4-1", false***REMOVED***, 5, true, 2, true***REMOVED***,
		// Name match only, because sensitive.
		***REMOVED***HeaderField***REMOVED***"key4", "value4-1", true***REMOVED***, 5, false, 2, false***REMOVED***,
		// Key not found.
		***REMOVED***HeaderField***REMOVED***"key5", "value5-x", false***REMOVED***, 0, false, 0, false***REMOVED***,
	***REMOVED***

	staticToDynamic := func(i uint64) uint64 ***REMOVED***
		if i == 0 ***REMOVED***
			return 0
		***REMOVED***
		return uint64(table.len()) - i + 1 // dynamic is the reversed table
	***REMOVED***

	searchStatic := func(f HeaderField) (uint64, bool) ***REMOVED***
		old := staticTable
		staticTable = table
		defer func() ***REMOVED*** staticTable = old ***REMOVED***()
		return staticTable.search(f)
	***REMOVED***

	searchDynamic := func(f HeaderField) (uint64, bool) ***REMOVED***
		return table.search(f)
	***REMOVED***

	for _, test := range tests ***REMOVED***
		gotI, gotMatch := searchStatic(test.f)
		if wantI, wantMatch := test.beforeWantStaticI, test.beforeWantMatch; gotI != wantI || gotMatch != wantMatch ***REMOVED***
			t.Errorf("before evictions: searchStatic(%+v)=%v,%v want %v,%v", test.f, gotI, gotMatch, wantI, wantMatch)
		***REMOVED***
		gotI, gotMatch = searchDynamic(test.f)
		wantDynamicI := staticToDynamic(test.beforeWantStaticI)
		if wantI, wantMatch := wantDynamicI, test.beforeWantMatch; gotI != wantI || gotMatch != wantMatch ***REMOVED***
			t.Errorf("before evictions: searchDynamic(%+v)=%v,%v want %v,%v", test.f, gotI, gotMatch, wantI, wantMatch)
		***REMOVED***
	***REMOVED***

	table.evictOldest(3)

	for _, test := range tests ***REMOVED***
		gotI, gotMatch := searchStatic(test.f)
		if wantI, wantMatch := test.afterWantStaticI, test.afterWantMatch; gotI != wantI || gotMatch != wantMatch ***REMOVED***
			t.Errorf("after evictions: searchStatic(%+v)=%v,%v want %v,%v", test.f, gotI, gotMatch, wantI, wantMatch)
		***REMOVED***
		gotI, gotMatch = searchDynamic(test.f)
		wantDynamicI := staticToDynamic(test.afterWantStaticI)
		if wantI, wantMatch := wantDynamicI, test.afterWantMatch; gotI != wantI || gotMatch != wantMatch ***REMOVED***
			t.Errorf("after evictions: searchDynamic(%+v)=%v,%v want %v,%v", test.f, gotI, gotMatch, wantI, wantMatch)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestHeaderFieldTable_LookupMapEviction(t *testing.T) ***REMOVED***
	table := &headerFieldTable***REMOVED******REMOVED***
	table.init()
	table.addEntry(pair("key1", "value1-1"))
	table.addEntry(pair("key2", "value2-1"))
	table.addEntry(pair("key1", "value1-2"))
	table.addEntry(pair("key3", "value3-1"))
	table.addEntry(pair("key4", "value4-1"))
	table.addEntry(pair("key2", "value2-2"))

	// evict all pairs
	table.evictOldest(table.len())

	if l := table.len(); l > 0 ***REMOVED***
		t.Errorf("table.len() = %d, want 0", l)
	***REMOVED***

	if l := len(table.byName); l > 0 ***REMOVED***
		t.Errorf("len(table.byName) = %d, want 0", l)
	***REMOVED***

	if l := len(table.byNameValue); l > 0 ***REMOVED***
		t.Errorf("len(table.byNameValue) = %d, want 0", l)
	***REMOVED***
***REMOVED***

func TestStaticTable(t *testing.T) ***REMOVED***
	fromSpec := `
          +-------+-----------------------------+---------------+
          | 1     | :authority                  |               |
          | 2     | :method                     | GET           |
          | 3     | :method                     | POST          |
          | 4     | :path                       | /             |
          | 5     | :path                       | /index.html   |
          | 6     | :scheme                     | http          |
          | 7     | :scheme                     | https         |
          | 8     | :status                     | 200           |
          | 9     | :status                     | 204           |
          | 10    | :status                     | 206           |
          | 11    | :status                     | 304           |
          | 12    | :status                     | 400           |
          | 13    | :status                     | 404           |
          | 14    | :status                     | 500           |
          | 15    | accept-charset              |               |
          | 16    | accept-encoding             | gzip, deflate |
          | 17    | accept-language             |               |
          | 18    | accept-ranges               |               |
          | 19    | accept                      |               |
          | 20    | access-control-allow-origin |               |
          | 21    | age                         |               |
          | 22    | allow                       |               |
          | 23    | authorization               |               |
          | 24    | cache-control               |               |
          | 25    | content-disposition         |               |
          | 26    | content-encoding            |               |
          | 27    | content-language            |               |
          | 28    | content-length              |               |
          | 29    | content-location            |               |
          | 30    | content-range               |               |
          | 31    | content-type                |               |
          | 32    | cookie                      |               |
          | 33    | date                        |               |
          | 34    | etag                        |               |
          | 35    | expect                      |               |
          | 36    | expires                     |               |
          | 37    | from                        |               |
          | 38    | host                        |               |
          | 39    | if-match                    |               |
          | 40    | if-modified-since           |               |
          | 41    | if-none-match               |               |
          | 42    | if-range                    |               |
          | 43    | if-unmodified-since         |               |
          | 44    | last-modified               |               |
          | 45    | link                        |               |
          | 46    | location                    |               |
          | 47    | max-forwards                |               |
          | 48    | proxy-authenticate          |               |
          | 49    | proxy-authorization         |               |
          | 50    | range                       |               |
          | 51    | referer                     |               |
          | 52    | refresh                     |               |
          | 53    | retry-after                 |               |
          | 54    | server                      |               |
          | 55    | set-cookie                  |               |
          | 56    | strict-transport-security   |               |
          | 57    | transfer-encoding           |               |
          | 58    | user-agent                  |               |
          | 59    | vary                        |               |
          | 60    | via                         |               |
          | 61    | www-authenticate            |               |
          +-------+-----------------------------+---------------+
`
	bs := bufio.NewScanner(strings.NewReader(fromSpec))
	re := regexp.MustCompile(`\| (\d+)\s+\| (\S+)\s*\| (\S(.*\S)?)?\s+\|`)
	for bs.Scan() ***REMOVED***
		l := bs.Text()
		if !strings.Contains(l, "|") ***REMOVED***
			continue
		***REMOVED***
		m := re.FindStringSubmatch(l)
		if m == nil ***REMOVED***
			continue
		***REMOVED***
		i, err := strconv.Atoi(m[1])
		if err != nil ***REMOVED***
			t.Errorf("Bogus integer on line %q", l)
			continue
		***REMOVED***
		if i < 1 || i > staticTable.len() ***REMOVED***
			t.Errorf("Bogus index %d on line %q", i, l)
			continue
		***REMOVED***
		if got, want := staticTable.ents[i-1].Name, m[2]; got != want ***REMOVED***
			t.Errorf("header index %d name = %q; want %q", i, got, want)
		***REMOVED***
		if got, want := staticTable.ents[i-1].Value, m[3]; got != want ***REMOVED***
			t.Errorf("header index %d value = %q; want %q", i, got, want)
		***REMOVED***
	***REMOVED***
	if err := bs.Err(); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***
