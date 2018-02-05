// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"testing"

	"golang.org/x/text/internal/colltab"
)

type ceTest struct ***REMOVED***
	f   func(in []int) (uint32, error)
	arg []int
	val uint32
***REMOVED***

func normalCE(in []int) (ce uint32, err error) ***REMOVED***
	return makeCE(rawCE***REMOVED***w: in[:3], ccc: uint8(in[3])***REMOVED***)
***REMOVED***

func expandCE(in []int) (ce uint32, err error) ***REMOVED***
	return makeExpandIndex(in[0])
***REMOVED***

func contractCE(in []int) (ce uint32, err error) ***REMOVED***
	return makeContractIndex(ctHandle***REMOVED***in[0], in[1]***REMOVED***, in[2])
***REMOVED***

func decompCE(in []int) (ce uint32, err error) ***REMOVED***
	return makeDecompose(in[0], in[1])
***REMOVED***

var ceTests = []ceTest***REMOVED***
	***REMOVED***normalCE, []int***REMOVED***0, 0, 0, 0***REMOVED***, 0xA0000000***REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***0, 0x28, 3, 0***REMOVED***, 0xA0002803***REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***0, 0x28, 3, 0xFF***REMOVED***, 0xAFF02803***REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***100, defaultSecondary, 3, 0***REMOVED***, 0x0000C883***REMOVED***,
	// non-ignorable primary with non-default secondary
	***REMOVED***normalCE, []int***REMOVED***100, 0x28, defaultTertiary, 0***REMOVED***, 0x4000C828***REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***100, defaultSecondary + 8, 3, 0***REMOVED***, 0x0000C983***REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***100, 0, 3, 0***REMOVED***, 0xFFFF***REMOVED***, // non-ignorable primary with non-supported secondary
	***REMOVED***normalCE, []int***REMOVED***100, 1, 3, 0***REMOVED***, 0xFFFF***REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***1 << maxPrimaryBits, defaultSecondary, 0, 0***REMOVED***, 0xFFFF***REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***0, 1 << maxSecondaryBits, 0, 0***REMOVED***, 0xFFFF***REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***100, defaultSecondary, 1 << maxTertiaryBits, 0***REMOVED***, 0xFFFF***REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***0x123, defaultSecondary, 8, 0xFF***REMOVED***, 0x88FF0123***REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***0x123, defaultSecondary + 1, 8, 0xFF***REMOVED***, 0xFFFF***REMOVED***,

	***REMOVED***contractCE, []int***REMOVED***0, 0, 0***REMOVED***, 0xC0000000***REMOVED***,
	***REMOVED***contractCE, []int***REMOVED***1, 1, 1***REMOVED***, 0xC0010011***REMOVED***,
	***REMOVED***contractCE, []int***REMOVED***1, (1 << maxNBits) - 1, 1***REMOVED***, 0xC001001F***REMOVED***,
	***REMOVED***contractCE, []int***REMOVED***(1 << maxTrieIndexBits) - 1, 1, 1***REMOVED***, 0xC001FFF1***REMOVED***,
	***REMOVED***contractCE, []int***REMOVED***1, 1, (1 << maxContractOffsetBits) - 1***REMOVED***, 0xDFFF0011***REMOVED***,
	***REMOVED***contractCE, []int***REMOVED***1, (1 << maxNBits), 1***REMOVED***, 0xFFFF***REMOVED***,
	***REMOVED***contractCE, []int***REMOVED***(1 << maxTrieIndexBits), 1, 1***REMOVED***, 0xFFFF***REMOVED***,
	***REMOVED***contractCE, []int***REMOVED***1, (1 << maxContractOffsetBits), 1***REMOVED***, 0xFFFF***REMOVED***,

	***REMOVED***expandCE, []int***REMOVED***0***REMOVED***, 0xE0000000***REMOVED***,
	***REMOVED***expandCE, []int***REMOVED***5***REMOVED***, 0xE0000005***REMOVED***,
	***REMOVED***expandCE, []int***REMOVED***(1 << maxExpandIndexBits) - 1***REMOVED***, 0xE000FFFF***REMOVED***,
	***REMOVED***expandCE, []int***REMOVED***1 << maxExpandIndexBits***REMOVED***, 0xFFFF***REMOVED***,

	***REMOVED***decompCE, []int***REMOVED***0, 0***REMOVED***, 0xF0000000***REMOVED***,
	***REMOVED***decompCE, []int***REMOVED***1, 1***REMOVED***, 0xF0000101***REMOVED***,
	***REMOVED***decompCE, []int***REMOVED***0x1F, 0x1F***REMOVED***, 0xF0001F1F***REMOVED***,
	***REMOVED***decompCE, []int***REMOVED***256, 0x1F***REMOVED***, 0xFFFF***REMOVED***,
	***REMOVED***decompCE, []int***REMOVED***0x1F, 256***REMOVED***, 0xFFFF***REMOVED***,
***REMOVED***

func TestColElem(t *testing.T) ***REMOVED***
	for i, tt := range ceTests ***REMOVED***
		in := make([]int, len(tt.arg))
		copy(in, tt.arg)
		ce, err := tt.f(in)
		if tt.val == 0xFFFF ***REMOVED***
			if err == nil ***REMOVED***
				t.Errorf("%d: expected error for args %x", i, tt.arg)
			***REMOVED***
			continue
		***REMOVED***
		if err != nil ***REMOVED***
			t.Errorf("%d: unexpected error: %v", i, err.Error())
		***REMOVED***
		if ce != tt.val ***REMOVED***
			t.Errorf("%d: colElem=%X; want %X", i, ce, tt.val)
		***REMOVED***
	***REMOVED***
***REMOVED***

func mkRawCES(in [][]int) []rawCE ***REMOVED***
	out := []rawCE***REMOVED******REMOVED***
	for _, w := range in ***REMOVED***
		out = append(out, rawCE***REMOVED***w: w***REMOVED***)
	***REMOVED***
	return out
***REMOVED***

type weightsTest struct ***REMOVED***
	a, b   [][]int
	level  colltab.Level
	result int
***REMOVED***

var nextWeightTests = []weightsTest***REMOVED***
	***REMOVED***
		a:     [][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED******REMOVED***,
		b:     [][]int***REMOVED******REMOVED***101, defaultSecondary, defaultTertiary, 0***REMOVED******REMOVED***,
		level: colltab.Primary,
	***REMOVED***,
	***REMOVED***
		a:     [][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED******REMOVED***,
		b:     [][]int***REMOVED******REMOVED***100, 21, defaultTertiary, 0***REMOVED******REMOVED***,
		level: colltab.Secondary,
	***REMOVED***,
	***REMOVED***
		a:     [][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED******REMOVED***,
		b:     [][]int***REMOVED******REMOVED***100, 20, 6, 0***REMOVED******REMOVED***,
		level: colltab.Tertiary,
	***REMOVED***,
	***REMOVED***
		a:     [][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED******REMOVED***,
		b:     [][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED******REMOVED***,
		level: colltab.Identity,
	***REMOVED***,
***REMOVED***

var extra = [][]int***REMOVED******REMOVED***200, 32, 8, 0***REMOVED***, ***REMOVED***0, 32, 8, 0***REMOVED***, ***REMOVED***0, 0, 8, 0***REMOVED***, ***REMOVED***0, 0, 0, 0***REMOVED******REMOVED***

func TestNextWeight(t *testing.T) ***REMOVED***
	for i, tt := range nextWeightTests ***REMOVED***
		test := func(l colltab.Level, tt weightsTest, a, gold [][]int) ***REMOVED***
			res := nextWeight(tt.level, mkRawCES(a))
			if !equalCEArrays(mkRawCES(gold), res) ***REMOVED***
				t.Errorf("%d:%d: expected weights %d; found %d", i, l, gold, res)
			***REMOVED***
		***REMOVED***
		test(-1, tt, tt.a, tt.b)
		for l := colltab.Primary; l <= colltab.Tertiary; l++ ***REMOVED***
			if tt.level <= l ***REMOVED***
				test(l, tt, append(tt.a, extra[l]), tt.b)
			***REMOVED*** else ***REMOVED***
				test(l, tt, append(tt.a, extra[l]), append(tt.b, extra[l]))
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var compareTests = []weightsTest***REMOVED***
	***REMOVED***
		[][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED******REMOVED***,
		[][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED******REMOVED***,
		colltab.Identity,
		0,
	***REMOVED***,
	***REMOVED***
		[][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED***, extra[0]***REMOVED***,
		[][]int***REMOVED******REMOVED***100, 20, 5, 1***REMOVED******REMOVED***,
		colltab.Primary,
		1,
	***REMOVED***,
	***REMOVED***
		[][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED******REMOVED***,
		[][]int***REMOVED******REMOVED***101, 20, 5, 0***REMOVED******REMOVED***,
		colltab.Primary,
		-1,
	***REMOVED***,
	***REMOVED***
		[][]int***REMOVED******REMOVED***101, 20, 5, 0***REMOVED******REMOVED***,
		[][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED******REMOVED***,
		colltab.Primary,
		1,
	***REMOVED***,
	***REMOVED***
		[][]int***REMOVED******REMOVED***100, 0, 0, 0***REMOVED***, ***REMOVED***0, 20, 5, 0***REMOVED******REMOVED***,
		[][]int***REMOVED******REMOVED***0, 20, 5, 0***REMOVED***, ***REMOVED***100, 0, 0, 0***REMOVED******REMOVED***,
		colltab.Identity,
		0,
	***REMOVED***,
	***REMOVED***
		[][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED******REMOVED***,
		[][]int***REMOVED******REMOVED***100, 21, 5, 0***REMOVED******REMOVED***,
		colltab.Secondary,
		-1,
	***REMOVED***,
	***REMOVED***
		[][]int***REMOVED******REMOVED***100, 20, 5, 0***REMOVED******REMOVED***,
		[][]int***REMOVED******REMOVED***100, 20, 2, 0***REMOVED******REMOVED***,
		colltab.Tertiary,
		1,
	***REMOVED***,
	***REMOVED***
		[][]int***REMOVED******REMOVED***100, 20, 5, 1***REMOVED******REMOVED***,
		[][]int***REMOVED******REMOVED***100, 20, 5, 2***REMOVED******REMOVED***,
		colltab.Quaternary,
		-1,
	***REMOVED***,
***REMOVED***

func TestCompareWeights(t *testing.T) ***REMOVED***
	for i, tt := range compareTests ***REMOVED***
		test := func(tt weightsTest, a, b [][]int) ***REMOVED***
			res, level := compareWeights(mkRawCES(a), mkRawCES(b))
			if res != tt.result ***REMOVED***
				t.Errorf("%d: expected comparison result %d; found %d", i, tt.result, res)
			***REMOVED***
			if level != tt.level ***REMOVED***
				t.Errorf("%d: expected level %d; found %d", i, tt.level, level)
			***REMOVED***
		***REMOVED***
		test(tt, tt.a, tt.b)
		test(tt, append(tt.a, extra[0]), append(tt.b, extra[0]))
	***REMOVED***
***REMOVED***
