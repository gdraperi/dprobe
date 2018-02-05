// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rangetable

import (
	"testing"
	"unicode"
)

var (
	maxRuneTable = &unicode.RangeTable***REMOVED***
		R32: []unicode.Range32***REMOVED***
			***REMOVED***unicode.MaxRune, unicode.MaxRune, 1***REMOVED***,
		***REMOVED***,
	***REMOVED***

	overlap1 = &unicode.RangeTable***REMOVED***
		R16: []unicode.Range16***REMOVED***
			***REMOVED***0x100, 0xfffc, 4***REMOVED***,
		***REMOVED***,
		R32: []unicode.Range32***REMOVED***
			***REMOVED***0x100000, 0x10fffc, 4***REMOVED***,
		***REMOVED***,
	***REMOVED***

	overlap2 = &unicode.RangeTable***REMOVED***
		R16: []unicode.Range16***REMOVED***
			***REMOVED***0x101, 0xfffd, 4***REMOVED***,
		***REMOVED***,
		R32: []unicode.Range32***REMOVED***
			***REMOVED***0x100001, 0x10fffd, 3***REMOVED***,
		***REMOVED***,
	***REMOVED***

	// The following table should be compacted into two entries for R16 and R32.
	optimize = &unicode.RangeTable***REMOVED***
		R16: []unicode.Range16***REMOVED***
			***REMOVED***0x1, 0x1, 1***REMOVED***,
			***REMOVED***0x2, 0x2, 1***REMOVED***,
			***REMOVED***0x3, 0x3, 1***REMOVED***,
			***REMOVED***0x5, 0x5, 1***REMOVED***,
			***REMOVED***0x7, 0x7, 1***REMOVED***,
			***REMOVED***0x9, 0x9, 1***REMOVED***,
			***REMOVED***0xb, 0xf, 2***REMOVED***,
		***REMOVED***,
		R32: []unicode.Range32***REMOVED***
			***REMOVED***0x10001, 0x10001, 1***REMOVED***,
			***REMOVED***0x10002, 0x10002, 1***REMOVED***,
			***REMOVED***0x10003, 0x10003, 1***REMOVED***,
			***REMOVED***0x10005, 0x10005, 1***REMOVED***,
			***REMOVED***0x10007, 0x10007, 1***REMOVED***,
			***REMOVED***0x10009, 0x10009, 1***REMOVED***,
			***REMOVED***0x1000b, 0x1000f, 2***REMOVED***,
		***REMOVED***,
	***REMOVED***
)

func TestMerge(t *testing.T) ***REMOVED***
	for i, tt := range [][]*unicode.RangeTable***REMOVED***
		***REMOVED***unicode.Cc, unicode.Cf***REMOVED***,
		***REMOVED***unicode.L, unicode.Ll***REMOVED***,
		***REMOVED***unicode.L, unicode.Ll, unicode.Lu***REMOVED***,
		***REMOVED***unicode.Ll, unicode.Lu***REMOVED***,
		***REMOVED***unicode.M***REMOVED***,
		unicode.GraphicRanges,
		cased,

		// Merge R16 only and R32 only and vice versa.
		***REMOVED***unicode.Khmer, unicode.Khudawadi***REMOVED***,
		***REMOVED***unicode.Imperial_Aramaic, unicode.Radical***REMOVED***,

		// Merge with empty.
		***REMOVED***&unicode.RangeTable***REMOVED******REMOVED******REMOVED***,
		***REMOVED***&unicode.RangeTable***REMOVED******REMOVED***, &unicode.RangeTable***REMOVED******REMOVED******REMOVED***,
		***REMOVED***&unicode.RangeTable***REMOVED******REMOVED***, &unicode.RangeTable***REMOVED******REMOVED***, &unicode.RangeTable***REMOVED******REMOVED******REMOVED***,
		***REMOVED***&unicode.RangeTable***REMOVED******REMOVED***, unicode.Hiragana***REMOVED***,
		***REMOVED***unicode.Inherited, &unicode.RangeTable***REMOVED******REMOVED******REMOVED***,
		***REMOVED***&unicode.RangeTable***REMOVED******REMOVED***, unicode.Hanunoo, &unicode.RangeTable***REMOVED******REMOVED******REMOVED***,

		// Hypothetical tables.
		***REMOVED***maxRuneTable***REMOVED***,
		***REMOVED***overlap1, overlap2***REMOVED***,

		// Optimization
		***REMOVED***optimize***REMOVED***,
	***REMOVED*** ***REMOVED***
		rt := Merge(tt...)
		for r := rune(0); r <= unicode.MaxRune; r++ ***REMOVED***
			if got, want := unicode.Is(rt, r), unicode.In(r, tt...); got != want ***REMOVED***
				t.Fatalf("%d:%U: got %v; want %v", i, r, got, want)
			***REMOVED***
		***REMOVED***
		// Test optimization and correctness for R16.
		for k := 0; k < len(rt.R16)-1; k++ ***REMOVED***
			if lo, hi := rt.R16[k].Lo, rt.R16[k].Hi; lo > hi ***REMOVED***
				t.Errorf("%d: Lo (%x) > Hi (%x)", i, lo, hi)
			***REMOVED***
			if hi, lo := rt.R16[k].Hi, rt.R16[k+1].Lo; hi >= lo ***REMOVED***
				t.Errorf("%d: Hi (%x) >= next Lo (%x)", i, hi, lo)
			***REMOVED***
			if rt.R16[k].Hi+rt.R16[k].Stride == rt.R16[k+1].Lo ***REMOVED***
				t.Errorf("%d: missed optimization for R16 at %d between %X and %x",
					i, k, rt.R16[k], rt.R16[k+1])
			***REMOVED***
		***REMOVED***
		// Test optimization and correctness for R32.
		for k := 0; k < len(rt.R32)-1; k++ ***REMOVED***
			if lo, hi := rt.R32[k].Lo, rt.R32[k].Hi; lo > hi ***REMOVED***
				t.Errorf("%d: Lo (%x) > Hi (%x)", i, lo, hi)
			***REMOVED***
			if hi, lo := rt.R32[k].Hi, rt.R32[k+1].Lo; hi >= lo ***REMOVED***
				t.Errorf("%d: Hi (%x) >= next Lo (%x)", i, hi, lo)
			***REMOVED***
			if rt.R32[k].Hi+rt.R32[k].Stride == rt.R32[k+1].Lo ***REMOVED***
				t.Errorf("%d: missed optimization for R32 at %d between %X and %X",
					i, k, rt.R32[k], rt.R32[k+1])
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

const runes = "Hello World in 2015!,\U0010fffd"

func BenchmarkNotMerged(t *testing.B) ***REMOVED***
	for i := 0; i < t.N; i++ ***REMOVED***
		for _, r := range runes ***REMOVED***
			unicode.In(r, unicode.GraphicRanges...)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkMerged(t *testing.B) ***REMOVED***
	rt := Merge(unicode.GraphicRanges...)

	for i := 0; i < t.N; i++ ***REMOVED***
		for _, r := range runes ***REMOVED***
			unicode.Is(rt, r)
		***REMOVED***
	***REMOVED***
***REMOVED***

var cased = []*unicode.RangeTable***REMOVED***
	unicode.Lower,
	unicode.Upper,
	unicode.Title,
	unicode.Other_Lowercase,
	unicode.Other_Uppercase,
***REMOVED***

func BenchmarkNotMergedCased(t *testing.B) ***REMOVED***
	for i := 0; i < t.N; i++ ***REMOVED***
		for _, r := range runes ***REMOVED***
			unicode.In(r, cased...)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkMergedCased(t *testing.B) ***REMOVED***
	// This reduces len(R16) from 243 to 82 and len(R32) from 65 to 35 for
	// Unicode 7.0.0.
	rt := Merge(cased...)

	for i := 0; i < t.N; i++ ***REMOVED***
		for _, r := range runes ***REMOVED***
			unicode.Is(rt, r)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkInit(t *testing.B) ***REMOVED***
	for i := 0; i < t.N; i++ ***REMOVED***
		Merge(cased...)
		Merge(unicode.GraphicRanges...)
	***REMOVED***
***REMOVED***

func BenchmarkInit2(t *testing.B) ***REMOVED***
	// Hypothetical near-worst-case performance.
	for i := 0; i < t.N; i++ ***REMOVED***
		Merge(overlap1, overlap2)
	***REMOVED***
***REMOVED***
