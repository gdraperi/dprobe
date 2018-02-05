// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colltab

import (
	"fmt"
	"testing"
	"unicode"
)

func (e Elem) String() string ***REMOVED***
	q := ""
	if v := e.Quaternary(); v == MaxQuaternary ***REMOVED***
		q = "max"
	***REMOVED*** else ***REMOVED***
		q = fmt.Sprint(v)
	***REMOVED***
	return fmt.Sprintf("[%d, %d, %d, %s]",
		e.Primary(),
		e.Secondary(),
		e.Tertiary(),
		q)
***REMOVED***

type ceTest struct ***REMOVED***
	f   func(inout []int) (Elem, ceType)
	arg []int
***REMOVED***

func makeCE(weights []int) Elem ***REMOVED***
	ce, _ := MakeElem(weights[0], weights[1], weights[2], uint8(weights[3]))
	return ce
***REMOVED***

var defaultValues = []int***REMOVED***0, defaultSecondary, defaultTertiary, 0***REMOVED***

func e(w ...int) Elem ***REMOVED***
	return makeCE(append(w, defaultValues[len(w):]...))
***REMOVED***

func makeContractIndex(index, n, offset int) Elem ***REMOVED***
	const (
		contractID            = 0xC0000000
		maxNBits              = 4
		maxTrieIndexBits      = 12
		maxContractOffsetBits = 13
	)
	ce := Elem(contractID)
	ce += Elem(offset << (maxNBits + maxTrieIndexBits))
	ce += Elem(index << maxNBits)
	ce += Elem(n)
	return ce
***REMOVED***

func makeExpandIndex(index int) Elem ***REMOVED***
	const expandID = 0xE0000000
	return expandID + Elem(index)
***REMOVED***

func makeDecompose(t1, t2 int) Elem ***REMOVED***
	const decompID = 0xF0000000
	return Elem(t2<<8+t1) + decompID
***REMOVED***

func normalCE(inout []int) (ce Elem, t ceType) ***REMOVED***
	ce = makeCE(inout)
	inout[0] = ce.Primary()
	inout[1] = ce.Secondary()
	inout[2] = int(ce.Tertiary())
	inout[3] = int(ce.CCC())
	return ce, ceNormal
***REMOVED***

func expandCE(inout []int) (ce Elem, t ceType) ***REMOVED***
	ce = makeExpandIndex(inout[0])
	inout[0] = splitExpandIndex(ce)
	return ce, ceExpansionIndex
***REMOVED***

func contractCE(inout []int) (ce Elem, t ceType) ***REMOVED***
	ce = makeContractIndex(inout[0], inout[1], inout[2])
	i, n, o := splitContractIndex(ce)
	inout[0], inout[1], inout[2] = i, n, o
	return ce, ceContractionIndex
***REMOVED***

func decompCE(inout []int) (ce Elem, t ceType) ***REMOVED***
	ce = makeDecompose(inout[0], inout[1])
	t1, t2 := splitDecompose(ce)
	inout[0], inout[1] = int(t1), int(t2)
	return ce, ceDecompose
***REMOVED***

var ceTests = []ceTest***REMOVED***
	***REMOVED***normalCE, []int***REMOVED***0, 0, 0, 0***REMOVED******REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***0, 30, 3, 0***REMOVED******REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***0, 30, 3, 0xFF***REMOVED******REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***100, defaultSecondary, defaultTertiary, 0***REMOVED******REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***100, defaultSecondary, defaultTertiary, 0xFF***REMOVED******REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***100, defaultSecondary, 3, 0***REMOVED******REMOVED***,
	***REMOVED***normalCE, []int***REMOVED***0x123, defaultSecondary, 8, 0xFF***REMOVED******REMOVED***,

	***REMOVED***contractCE, []int***REMOVED***0, 0, 0***REMOVED******REMOVED***,
	***REMOVED***contractCE, []int***REMOVED***1, 1, 1***REMOVED******REMOVED***,
	***REMOVED***contractCE, []int***REMOVED***1, (1 << maxNBits) - 1, 1***REMOVED******REMOVED***,
	***REMOVED***contractCE, []int***REMOVED***(1 << maxTrieIndexBits) - 1, 1, 1***REMOVED******REMOVED***,
	***REMOVED***contractCE, []int***REMOVED***1, 1, (1 << maxContractOffsetBits) - 1***REMOVED******REMOVED***,

	***REMOVED***expandCE, []int***REMOVED***0***REMOVED******REMOVED***,
	***REMOVED***expandCE, []int***REMOVED***5***REMOVED******REMOVED***,
	***REMOVED***expandCE, []int***REMOVED***(1 << maxExpandIndexBits) - 1***REMOVED******REMOVED***,

	***REMOVED***decompCE, []int***REMOVED***0, 0***REMOVED******REMOVED***,
	***REMOVED***decompCE, []int***REMOVED***1, 1***REMOVED******REMOVED***,
	***REMOVED***decompCE, []int***REMOVED***0x1F, 0x1F***REMOVED******REMOVED***,
***REMOVED***

func TestColElem(t *testing.T) ***REMOVED***
	for i, tt := range ceTests ***REMOVED***
		inout := make([]int, len(tt.arg))
		copy(inout, tt.arg)
		ce, typ := tt.f(inout)
		if ce.ctype() != typ ***REMOVED***
			t.Errorf("%d: type is %d; want %d (ColElem: %X)", i, ce.ctype(), typ, ce)
		***REMOVED***
		for j, a := range tt.arg ***REMOVED***
			if inout[j] != a ***REMOVED***
				t.Errorf("%d: argument %d is %X; want %X (ColElem: %X)", i, j, inout[j], a, ce)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type implicitTest struct ***REMOVED***
	r rune
	p int
***REMOVED***

var implicitTests = []implicitTest***REMOVED***
	***REMOVED***0x33FF, 0x533FF***REMOVED***,
	***REMOVED***0x3400, 0x23400***REMOVED***,
	***REMOVED***0x4DC0, 0x54DC0***REMOVED***,
	***REMOVED***0x4DFF, 0x54DFF***REMOVED***,
	***REMOVED***0x4E00, 0x14E00***REMOVED***,
	***REMOVED***0x9FCB, 0x19FCB***REMOVED***,
	***REMOVED***0xA000, 0x5A000***REMOVED***,
	***REMOVED***0xF8FF, 0x5F8FF***REMOVED***,
	***REMOVED***0xF900, 0x1F900***REMOVED***,
	***REMOVED***0xFA23, 0x1FA23***REMOVED***,
	***REMOVED***0xFAD9, 0x1FAD9***REMOVED***,
	***REMOVED***0xFB00, 0x5FB00***REMOVED***,
	***REMOVED***0x20000, 0x40000***REMOVED***,
	***REMOVED***0x2B81C, 0x4B81C***REMOVED***,
	***REMOVED***unicode.MaxRune, 0x15FFFF***REMOVED***, // maximum primary value
***REMOVED***

func TestImplicit(t *testing.T) ***REMOVED***
	for _, tt := range implicitTests ***REMOVED***
		if p := implicitPrimary(tt.r); p != tt.p ***REMOVED***
			t.Errorf("%U: was %X; want %X", tt.r, p, tt.p)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUpdateTertiary(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		in, out Elem
		t       uint8
	***REMOVED******REMOVED***
		***REMOVED***0x4000FE20, 0x0000FE8A, 0x0A***REMOVED***,
		***REMOVED***0x4000FE21, 0x0000FEAA, 0x0A***REMOVED***,
		***REMOVED***0x0000FE8B, 0x0000FE83, 0x03***REMOVED***,
		***REMOVED***0x82FF0188, 0x9BFF0188, 0x1B***REMOVED***,
		***REMOVED***0xAFF0CC02, 0xAFF0CC1B, 0x1B***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		if out := tt.in.updateTertiary(tt.t); out != tt.out ***REMOVED***
			t.Errorf("%d: was %X; want %X", i, out, tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***
