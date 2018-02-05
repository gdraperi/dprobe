// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stringset

import "testing"

func TestStringSet(t *testing.T) ***REMOVED***
	testCases := [][]string***REMOVED***
		***REMOVED***""***REMOVED***,
		***REMOVED***"∫"***REMOVED***,
		***REMOVED***"a", "b", "c"***REMOVED***,
		***REMOVED***"", "a", "bb", "ccc"***REMOVED***,
		***REMOVED***"    ", "aaa", "bb", "c"***REMOVED***,
	***REMOVED***
	test := func(tc int, b *Builder) ***REMOVED***
		set := b.Set()
		if set.Len() != len(testCases[tc]) ***REMOVED***
			t.Errorf("%d:Len() = %d; want %d", tc, set.Len(), len(testCases[tc]))
		***REMOVED***
		for i, s := range testCases[tc] ***REMOVED***
			if x := b.Index(s); x != i ***REMOVED***
				t.Errorf("%d:Index(%q) = %d; want %d", tc, s, x, i)
			***REMOVED***
			if p := Search(&set, s); p != i ***REMOVED***
				t.Errorf("%d:Search(%q) = %d; want %d", tc, s, p, i)
			***REMOVED***
			if set.Elem(i) != s ***REMOVED***
				t.Errorf("%d:Elem(%d) = %s; want %s", tc, i, set.Elem(i), s)
			***REMOVED***
		***REMOVED***
		if p := Search(&set, "apple"); p != -1 ***REMOVED***
			t.Errorf(`%d:Search("apple") = %d; want -1`, tc, p)
		***REMOVED***
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		b := NewBuilder()
		for _, s := range tc ***REMOVED***
			b.Add(s)
		***REMOVED***
		b.Add(tc...)
		test(i, b)
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		b := NewBuilder()
		b.Add(tc...)
		for _, s := range tc ***REMOVED***
			b.Add(s)
		***REMOVED***
		test(i, b)
	***REMOVED***
***REMOVED***
