// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSupported(t *testing.T) ***REMOVED***
	// To prove the results are correct for a type, we test that the number of
	// results is identical to the number of results on record, that all results
	// are distinct and that all results are valid.
	tests := map[string]int***REMOVED***
		"BaseLanguages": numLanguages,
		"Scripts":       numScripts,
		"Regions":       numRegions,
		"Tags":          0,
	***REMOVED***
	sup := reflect.ValueOf(Supported)
	for name, num := range tests ***REMOVED***
		v := sup.MethodByName(name).Call(nil)[0]
		if n := v.Len(); n != num ***REMOVED***
			t.Errorf("len(%s()) was %d; want %d", name, n, num)
		***REMOVED***
		dup := make(map[string]bool)
		for i := 0; i < v.Len(); i++ ***REMOVED***
			x := v.Index(i).Interface()
			// An invalid value will either cause a crash or result in a
			// duplicate when passed to Sprint.
			s := fmt.Sprint(x)
			if dup[s] ***REMOVED***
				t.Errorf("%s: duplicate entry %q", name, s)
			***REMOVED***
			dup[s] = true
		***REMOVED***
		if len(dup) != v.Len() ***REMOVED***
			t.Errorf("%s: # unique entries was %d; want %d", name, len(dup), v.Len())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNewCoverage(t *testing.T) ***REMOVED***
	bases := []Base***REMOVED***Base***REMOVED***0***REMOVED***, Base***REMOVED***3***REMOVED***, Base***REMOVED***7***REMOVED******REMOVED***
	scripts := []Script***REMOVED***Script***REMOVED***11***REMOVED***, Script***REMOVED***17***REMOVED***, Script***REMOVED***23***REMOVED******REMOVED***
	regions := []Region***REMOVED***Region***REMOVED***101***REMOVED***, Region***REMOVED***103***REMOVED***, Region***REMOVED***107***REMOVED******REMOVED***
	tags := []Tag***REMOVED***Make("pt"), Make("en"), Make("en-GB"), Make("en-US"), Make("pt-PT")***REMOVED***
	fbases := func() []Base ***REMOVED*** return bases ***REMOVED***
	fscripts := func() []Script ***REMOVED*** return scripts ***REMOVED***
	fregions := func() []Region ***REMOVED*** return regions ***REMOVED***
	ftags := func() []Tag ***REMOVED*** return tags ***REMOVED***

	tests := []struct ***REMOVED***
		desc    string
		list    []interface***REMOVED******REMOVED***
		bases   []Base
		scripts []Script
		regions []Region
		tags    []Tag
	***REMOVED******REMOVED***
		***REMOVED***
			desc: "empty",
		***REMOVED***,
		***REMOVED***
			desc:  "bases",
			list:  []interface***REMOVED******REMOVED******REMOVED***bases***REMOVED***,
			bases: bases,
		***REMOVED***,
		***REMOVED***
			desc:    "scripts",
			list:    []interface***REMOVED******REMOVED******REMOVED***scripts***REMOVED***,
			scripts: scripts,
		***REMOVED***,
		***REMOVED***
			desc:    "regions",
			list:    []interface***REMOVED******REMOVED******REMOVED***regions***REMOVED***,
			regions: regions,
		***REMOVED***,
		***REMOVED***
			desc:  "bases derives from tags",
			list:  []interface***REMOVED******REMOVED******REMOVED***tags***REMOVED***,
			bases: []Base***REMOVED***Base***REMOVED***_en***REMOVED***, Base***REMOVED***_pt***REMOVED******REMOVED***,
			tags:  tags,
		***REMOVED***,
		***REMOVED***
			desc:  "tags and bases",
			list:  []interface***REMOVED******REMOVED******REMOVED***tags, bases***REMOVED***,
			bases: bases,
			tags:  tags,
		***REMOVED***,
		***REMOVED***
			desc:    "fully specified",
			list:    []interface***REMOVED******REMOVED******REMOVED***tags, bases, scripts, regions***REMOVED***,
			bases:   bases,
			scripts: scripts,
			regions: regions,
			tags:    tags,
		***REMOVED***,
		***REMOVED***
			desc:  "bases func",
			list:  []interface***REMOVED******REMOVED******REMOVED***fbases***REMOVED***,
			bases: bases,
		***REMOVED***,
		***REMOVED***
			desc:    "scripts func",
			list:    []interface***REMOVED******REMOVED******REMOVED***fscripts***REMOVED***,
			scripts: scripts,
		***REMOVED***,
		***REMOVED***
			desc:    "regions func",
			list:    []interface***REMOVED******REMOVED******REMOVED***fregions***REMOVED***,
			regions: regions,
		***REMOVED***,
		***REMOVED***
			desc:  "tags func",
			list:  []interface***REMOVED******REMOVED******REMOVED***ftags***REMOVED***,
			bases: []Base***REMOVED***Base***REMOVED***_en***REMOVED***, Base***REMOVED***_pt***REMOVED******REMOVED***,
			tags:  tags,
		***REMOVED***,
		***REMOVED***
			desc:  "tags and bases",
			list:  []interface***REMOVED******REMOVED******REMOVED***ftags, fbases***REMOVED***,
			bases: bases,
			tags:  tags,
		***REMOVED***,
		***REMOVED***
			desc:    "fully specified",
			list:    []interface***REMOVED******REMOVED******REMOVED***ftags, fbases, fscripts, fregions***REMOVED***,
			bases:   bases,
			scripts: scripts,
			regions: regions,
			tags:    tags,
		***REMOVED***,
	***REMOVED***

	for i, tt := range tests ***REMOVED***
		l := NewCoverage(tt.list...)
		if a := l.BaseLanguages(); !reflect.DeepEqual(a, tt.bases) ***REMOVED***
			t.Errorf("%d:%s: BaseLanguages was %v; want %v", i, tt.desc, a, tt.bases)
		***REMOVED***
		if a := l.Scripts(); !reflect.DeepEqual(a, tt.scripts) ***REMOVED***
			t.Errorf("%d:%s: Scripts was %v; want %v", i, tt.desc, a, tt.scripts)
		***REMOVED***
		if a := l.Regions(); !reflect.DeepEqual(a, tt.regions) ***REMOVED***
			t.Errorf("%d:%s: Regions was %v; want %v", i, tt.desc, a, tt.regions)
		***REMOVED***
		if a := l.Tags(); !reflect.DeepEqual(a, tt.tags) ***REMOVED***
			t.Errorf("%d:%s: Tags was %v; want %v", i, tt.desc, a, tt.tags)
		***REMOVED***
	***REMOVED***
***REMOVED***
