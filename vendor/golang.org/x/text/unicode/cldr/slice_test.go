// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldr

import (
	"reflect"
	"testing"
)

type testSlice []*Common

func mkElem(alt, typ, ref string) *Common ***REMOVED***
	return &Common***REMOVED***
		Type:      typ,
		Reference: ref,
		Alt:       alt,
	***REMOVED***
***REMOVED***

var (
	testSlice1 = testSlice***REMOVED***
		mkElem("1", "a", "i.a"),
		mkElem("1", "b", "i.b"),
		mkElem("1", "c", "i.c"),
		mkElem("2", "b", "ii"),
		mkElem("3", "c", "iii"),
		mkElem("4", "a", "iv.a"),
		mkElem("4", "d", "iv.d"),
	***REMOVED***
	testSliceE = testSlice***REMOVED******REMOVED***
)

func panics(f func()) (panics bool) ***REMOVED***
	defer func() ***REMOVED***
		if err := recover(); err != nil ***REMOVED***
			panics = true
		***REMOVED***
	***REMOVED***()
	f()
	return panics
***REMOVED***

func TestMakeSlice(t *testing.T) ***REMOVED***
	foo := 1
	bar := []int***REMOVED******REMOVED***
	tests := []struct ***REMOVED***
		i      interface***REMOVED******REMOVED***
		panics bool
		err    string
	***REMOVED******REMOVED***
		***REMOVED***&foo, true, "should panic when passed a pointer to the wrong type"***REMOVED***,
		***REMOVED***&bar, true, "should panic when slice element of the wrong type"***REMOVED***,
		***REMOVED***testSlice1, true, "should panic when passed a slice"***REMOVED***,
		***REMOVED***&testSlice1, false, "should not panic"***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		if panics(func() ***REMOVED*** MakeSlice(tt.i) ***REMOVED***) != tt.panics ***REMOVED***
			t.Errorf("%d: %s", i, tt.err)
		***REMOVED***
	***REMOVED***
***REMOVED***

var anyOfTests = []struct ***REMOVED***
	sl     testSlice
	values []string
	n      int
***REMOVED******REMOVED***
	***REMOVED***testSliceE, []string***REMOVED******REMOVED***, 0***REMOVED***,
	***REMOVED***testSliceE, []string***REMOVED***"1", "2", "3"***REMOVED***, 0***REMOVED***,
	***REMOVED***testSlice1, []string***REMOVED******REMOVED***, 0***REMOVED***,
	***REMOVED***testSlice1, []string***REMOVED***"1"***REMOVED***, 3***REMOVED***,
	***REMOVED***testSlice1, []string***REMOVED***"2"***REMOVED***, 1***REMOVED***,
	***REMOVED***testSlice1, []string***REMOVED***"5"***REMOVED***, 0***REMOVED***,
	***REMOVED***testSlice1, []string***REMOVED***"1", "2", "3"***REMOVED***, 5***REMOVED***,
***REMOVED***

func TestSelectAnyOf(t *testing.T) ***REMOVED***
	for i, tt := range anyOfTests ***REMOVED***
		sl := tt.sl
		s := MakeSlice(&sl)
		s.SelectAnyOf("alt", tt.values...)
		if len(sl) != tt.n ***REMOVED***
			t.Errorf("%d: found len == %d; want %d", i, len(sl), tt.n)
		***REMOVED***
	***REMOVED***
	sl := testSlice1
	s := MakeSlice(&sl)
	if !panics(func() ***REMOVED*** s.SelectAnyOf("foo") ***REMOVED***) ***REMOVED***
		t.Errorf("should panic on non-existing attribute")
	***REMOVED***
***REMOVED***

func TestFilter(t *testing.T) ***REMOVED***
	for i, tt := range anyOfTests ***REMOVED***
		sl := tt.sl
		s := MakeSlice(&sl)
		s.Filter(func(e Elem) bool ***REMOVED***
			v, _ := findField(reflect.ValueOf(e), "alt")
			return in(tt.values, v.String())
		***REMOVED***)
		if len(sl) != tt.n ***REMOVED***
			t.Errorf("%d: found len == %d; want %d", i, len(sl), tt.n)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestGroup(t *testing.T) ***REMOVED***
	f := func(excl ...string) func(Elem) string ***REMOVED***
		return func(e Elem) string ***REMOVED***
			return Key(e, excl...)
		***REMOVED***
	***REMOVED***
	tests := []struct ***REMOVED***
		sl   testSlice
		f    func(Elem) string
		lens []int
	***REMOVED******REMOVED***
		***REMOVED***testSliceE, f(), []int***REMOVED******REMOVED******REMOVED***,
		***REMOVED***testSlice1, f(), []int***REMOVED***1, 1, 1, 1, 1, 1, 1***REMOVED******REMOVED***,
		***REMOVED***testSlice1, f("type"), []int***REMOVED***3, 1, 1, 2***REMOVED******REMOVED***,
		***REMOVED***testSlice1, f("alt"), []int***REMOVED***2, 2, 2, 1***REMOVED******REMOVED***,
		***REMOVED***testSlice1, f("alt", "type"), []int***REMOVED***7***REMOVED******REMOVED***,
		***REMOVED***testSlice1, f("alt", "type"), []int***REMOVED***7***REMOVED******REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		sl := tt.sl
		s := MakeSlice(&sl)
		g := s.Group(tt.f)
		if len(tt.lens) != len(g) ***REMOVED***
			t.Errorf("%d: found %d; want %d", i, len(g), len(tt.lens))
			continue
		***REMOVED***
		for j, v := range tt.lens ***REMOVED***
			if n := g[j].Value().Len(); n != v ***REMOVED***
				t.Errorf("%d: found %d for length of group %d; want %d", i, n, j, v)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSelectOnePerGroup(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		sl     testSlice
		attr   string
		values []string
		refs   []string
	***REMOVED******REMOVED***
		***REMOVED***testSliceE, "alt", []string***REMOVED***"1"***REMOVED***, []string***REMOVED******REMOVED******REMOVED***,
		***REMOVED***testSliceE, "type", []string***REMOVED***"a"***REMOVED***, []string***REMOVED******REMOVED******REMOVED***,
		***REMOVED***testSlice1, "alt", []string***REMOVED***"2", "3", "1"***REMOVED***, []string***REMOVED***"i.a", "ii", "iii"***REMOVED******REMOVED***,
		***REMOVED***testSlice1, "alt", []string***REMOVED***"1", "4"***REMOVED***, []string***REMOVED***"i.a", "i.b", "i.c", "iv.d"***REMOVED******REMOVED***,
		***REMOVED***testSlice1, "type", []string***REMOVED***"c", "d"***REMOVED***, []string***REMOVED***"i.c", "iii", "iv.d"***REMOVED******REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		sl := tt.sl
		s := MakeSlice(&sl)
		s.SelectOnePerGroup(tt.attr, tt.values)
		if len(sl) != len(tt.refs) ***REMOVED***
			t.Errorf("%d: found result length %d; want %d", i, len(sl), len(tt.refs))
			continue
		***REMOVED***
		for j, e := range sl ***REMOVED***
			if tt.refs[j] != e.Reference ***REMOVED***
				t.Errorf("%d:%d found %s; want %s", i, j, e.Reference, tt.refs[i])
			***REMOVED***
		***REMOVED***
	***REMOVED***
	sl := testSlice1
	s := MakeSlice(&sl)
	if !panics(func() ***REMOVED*** s.SelectOnePerGroup("foo", nil) ***REMOVED***) ***REMOVED***
		t.Errorf("should panic on non-existing attribute")
	***REMOVED***
***REMOVED***
