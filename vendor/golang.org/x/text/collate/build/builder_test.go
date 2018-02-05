// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import "testing"

// cjk returns an implicit collation element for a CJK rune.
func cjk(r rune) []rawCE ***REMOVED***
	// A CJK character C is represented in the DUCET as
	//   [.AAAA.0020.0002.C][.BBBB.0000.0000.C]
	// Where AAAA is the most significant 15 bits plus a base value.
	// Any base value will work for the test, so we pick the common value of FB40.
	const base = 0xFB40
	return []rawCE***REMOVED***
		***REMOVED***w: []int***REMOVED***base + int(r>>15), defaultSecondary, defaultTertiary, int(r)***REMOVED******REMOVED***,
		***REMOVED***w: []int***REMOVED***int(r&0x7FFF) | 0x8000, 0, 0, int(r)***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

func pCE(p int) []rawCE ***REMOVED***
	return mkCE([]int***REMOVED***p, defaultSecondary, defaultTertiary, 0***REMOVED***, 0)
***REMOVED***

func pqCE(p, q int) []rawCE ***REMOVED***
	return mkCE([]int***REMOVED***p, defaultSecondary, defaultTertiary, q***REMOVED***, 0)
***REMOVED***

func ptCE(p, t int) []rawCE ***REMOVED***
	return mkCE([]int***REMOVED***p, defaultSecondary, t, 0***REMOVED***, 0)
***REMOVED***

func ptcCE(p, t int, ccc uint8) []rawCE ***REMOVED***
	return mkCE([]int***REMOVED***p, defaultSecondary, t, 0***REMOVED***, ccc)
***REMOVED***

func sCE(s int) []rawCE ***REMOVED***
	return mkCE([]int***REMOVED***0, s, defaultTertiary, 0***REMOVED***, 0)
***REMOVED***

func stCE(s, t int) []rawCE ***REMOVED***
	return mkCE([]int***REMOVED***0, s, t, 0***REMOVED***, 0)
***REMOVED***

func scCE(s int, ccc uint8) []rawCE ***REMOVED***
	return mkCE([]int***REMOVED***0, s, defaultTertiary, 0***REMOVED***, ccc)
***REMOVED***

func mkCE(w []int, ccc uint8) []rawCE ***REMOVED***
	return []rawCE***REMOVED***rawCE***REMOVED***w, ccc***REMOVED******REMOVED***
***REMOVED***

// ducetElem is used to define test data that is used to generate a table.
type ducetElem struct ***REMOVED***
	str string
	ces []rawCE
***REMOVED***

func newBuilder(t *testing.T, ducet []ducetElem) *Builder ***REMOVED***
	b := NewBuilder()
	for _, e := range ducet ***REMOVED***
		ces := [][]int***REMOVED******REMOVED***
		for _, ce := range e.ces ***REMOVED***
			ces = append(ces, ce.w)
		***REMOVED***
		if err := b.Add([]rune(e.str), ces, nil); err != nil ***REMOVED***
			t.Errorf(err.Error())
		***REMOVED***
	***REMOVED***
	b.t = &table***REMOVED******REMOVED***
	b.root.sort()
	return b
***REMOVED***

type convertTest struct ***REMOVED***
	in, out []rawCE
	err     bool
***REMOVED***

var convLargeTests = []convertTest***REMOVED***
	***REMOVED***pCE(0xFB39), pCE(0xFB39), false***REMOVED***,
	***REMOVED***cjk(0x2F9B2), pqCE(0x3F9B2, 0x2F9B2), false***REMOVED***,
	***REMOVED***pCE(0xFB40), pCE(0), true***REMOVED***,
	***REMOVED***append(pCE(0xFB40), pCE(0)[0]), pCE(0), true***REMOVED***,
	***REMOVED***pCE(0xFFFE), pCE(illegalOffset), false***REMOVED***,
	***REMOVED***pCE(0xFFFF), pCE(illegalOffset + 1), false***REMOVED***,
***REMOVED***

func TestConvertLarge(t *testing.T) ***REMOVED***
	for i, tt := range convLargeTests ***REMOVED***
		e := new(entry)
		for _, ce := range tt.in ***REMOVED***
			e.elems = append(e.elems, makeRawCE(ce.w, ce.ccc))
		***REMOVED***
		elems, err := convertLargeWeights(e.elems)
		if tt.err ***REMOVED***
			if err == nil ***REMOVED***
				t.Errorf("%d: expected error; none found", i)
			***REMOVED***
			continue
		***REMOVED*** else if err != nil ***REMOVED***
			t.Errorf("%d: unexpected error: %v", i, err)
		***REMOVED***
		if !equalCEArrays(elems, tt.out) ***REMOVED***
			t.Errorf("%d: conversion was %x; want %x", i, elems, tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Collation element table for simplify tests.
var simplifyTest = []ducetElem***REMOVED***
	***REMOVED***"\u0300", sCE(30)***REMOVED***, // grave
	***REMOVED***"\u030C", sCE(40)***REMOVED***, // caron
	***REMOVED***"A", ptCE(100, 8)***REMOVED***,
	***REMOVED***"D", ptCE(104, 8)***REMOVED***,
	***REMOVED***"E", ptCE(105, 8)***REMOVED***,
	***REMOVED***"I", ptCE(110, 8)***REMOVED***,
	***REMOVED***"z", ptCE(130, 8)***REMOVED***,
	***REMOVED***"\u05F2", append(ptCE(200, 4), ptCE(200, 4)[0])***REMOVED***,
	***REMOVED***"\u05B7", sCE(80)***REMOVED***,
	***REMOVED***"\u00C0", append(ptCE(100, 8), sCE(30)...)***REMOVED***,                                // A with grave, can be removed
	***REMOVED***"\u00C8", append(ptCE(105, 8), sCE(30)...)***REMOVED***,                                // E with grave
	***REMOVED***"\uFB1F", append(ptCE(200, 4), ptCE(200, 4)[0], sCE(80)[0])***REMOVED***,               // eliminated by NFD
	***REMOVED***"\u00C8\u0302", ptCE(106, 8)***REMOVED***,                                              // block previous from simplifying
	***REMOVED***"\u01C5", append(ptCE(104, 9), ptCE(130, 4)[0], stCE(40, maxTertiary)[0])***REMOVED***, // eliminated by NFKD
	// no removal: tertiary value of third element is not maxTertiary
	***REMOVED***"\u2162", append(ptCE(110, 9), ptCE(110, 4)[0], ptCE(110, 8)[0])***REMOVED***,
***REMOVED***

var genColTests = []ducetElem***REMOVED***
	***REMOVED***"\uFA70", pqCE(0x1FA70, 0xFA70)***REMOVED***,
	***REMOVED***"A\u0300", append(ptCE(100, 8), sCE(30)...)***REMOVED***,
	***REMOVED***"A\u0300\uFA70", append(ptCE(100, 8), sCE(30)[0], pqCE(0x1FA70, 0xFA70)[0])***REMOVED***,
	***REMOVED***"A\u0300A\u0300", append(ptCE(100, 8), sCE(30)[0], ptCE(100, 8)[0], sCE(30)[0])***REMOVED***,
***REMOVED***

func TestGenColElems(t *testing.T) ***REMOVED***
	b := newBuilder(t, simplifyTest[:5])

	for i, tt := range genColTests ***REMOVED***
		res := b.root.genColElems(tt.str)
		if !equalCEArrays(tt.ces, res) ***REMOVED***
			t.Errorf("%d: result %X; want %X", i, res, tt.ces)
		***REMOVED***
	***REMOVED***
***REMOVED***

type strArray []string

func (sa strArray) contains(s string) bool ***REMOVED***
	for _, e := range sa ***REMOVED***
		if e == s ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

var simplifyRemoved = strArray***REMOVED***"\u00C0", "\uFB1F"***REMOVED***
var simplifyMarked = strArray***REMOVED***"\u01C5"***REMOVED***

func TestSimplify(t *testing.T) ***REMOVED***
	b := newBuilder(t, simplifyTest)
	o := &b.root
	simplify(o)

	for i, tt := range simplifyTest ***REMOVED***
		if simplifyRemoved.contains(tt.str) ***REMOVED***
			continue
		***REMOVED***
		e := o.find(tt.str)
		if e.str != tt.str || !equalCEArrays(e.elems, tt.ces) ***REMOVED***
			t.Errorf("%d: found element %s -> %X; want %s -> %X", i, e.str, e.elems, tt.str, tt.ces)
			break
		***REMOVED***
	***REMOVED***
	var i, k int
	for e := o.front(); e != nil; e, _ = e.nextIndexed() ***REMOVED***
		gold := simplifyMarked.contains(e.str)
		if gold ***REMOVED***
			k++
		***REMOVED***
		if gold != e.decompose ***REMOVED***
			t.Errorf("%d: %s has decompose %v; want %v", i, e.str, e.decompose, gold)
		***REMOVED***
		i++
	***REMOVED***
	if k != len(simplifyMarked) ***REMOVED***
		t.Errorf(" an entry that should be marked as decompose was deleted")
	***REMOVED***
***REMOVED***

var expandTest = []ducetElem***REMOVED***
	***REMOVED***"\u0300", append(scCE(29, 230), scCE(30, 230)...)***REMOVED***,
	***REMOVED***"\u00C0", append(ptCE(100, 8), scCE(30, 230)...)***REMOVED***,
	***REMOVED***"\u00C8", append(ptCE(105, 8), scCE(30, 230)...)***REMOVED***,
	***REMOVED***"\u00C9", append(ptCE(105, 8), scCE(30, 230)...)***REMOVED***, // identical expansion
	***REMOVED***"\u05F2", append(ptCE(200, 4), ptCE(200, 4)[0], ptCE(200, 4)[0])***REMOVED***,
	***REMOVED***"\u01FF", append(ptCE(200, 4), ptcCE(201, 4, 0)[0], scCE(30, 230)[0])***REMOVED***,
***REMOVED***

func TestExpand(t *testing.T) ***REMOVED***
	const (
		totalExpansions = 5
		totalElements   = 2 + 2 + 2 + 3 + 3 + totalExpansions
	)
	b := newBuilder(t, expandTest)
	o := &b.root
	b.processExpansions(o)

	e := o.front()
	for _, tt := range expandTest ***REMOVED***
		exp := b.t.ExpandElem[e.expansionIndex:]
		if int(exp[0]) != len(tt.ces) ***REMOVED***
			t.Errorf("%U: len(expansion)==%d; want %d", []rune(tt.str)[0], exp[0], len(tt.ces))
		***REMOVED***
		exp = exp[1:]
		for j, w := range tt.ces ***REMOVED***
			if ce, _ := makeCE(w); exp[j] != ce ***REMOVED***
				t.Errorf("%U: element %d is %X; want %X", []rune(tt.str)[0], j, exp[j], ce)
			***REMOVED***
		***REMOVED***
		e, _ = e.nextIndexed()
	***REMOVED***
	// Verify uniquing.
	if len(b.t.ExpandElem) != totalElements ***REMOVED***
		t.Errorf("len(expandElem)==%d; want %d", len(b.t.ExpandElem), totalElements)
	***REMOVED***
***REMOVED***

var contractTest = []ducetElem***REMOVED***
	***REMOVED***"abc", pCE(102)***REMOVED***,
	***REMOVED***"abd", pCE(103)***REMOVED***,
	***REMOVED***"a", pCE(100)***REMOVED***,
	***REMOVED***"ab", pCE(101)***REMOVED***,
	***REMOVED***"ac", pCE(104)***REMOVED***,
	***REMOVED***"bcd", pCE(202)***REMOVED***,
	***REMOVED***"b", pCE(200)***REMOVED***,
	***REMOVED***"bc", pCE(201)***REMOVED***,
	***REMOVED***"bd", pCE(203)***REMOVED***,
	// shares suffixes with a*
	***REMOVED***"Ab", pCE(301)***REMOVED***,
	***REMOVED***"A", pCE(300)***REMOVED***,
	***REMOVED***"Ac", pCE(304)***REMOVED***,
	***REMOVED***"Abc", pCE(302)***REMOVED***,
	***REMOVED***"Abd", pCE(303)***REMOVED***,
	// starter to be ignored
	***REMOVED***"z", pCE(1000)***REMOVED***,
***REMOVED***

func TestContract(t *testing.T) ***REMOVED***
	const (
		totalElements = 5 + 5 + 4
	)
	b := newBuilder(t, contractTest)
	o := &b.root
	b.processContractions(o)

	indexMap := make(map[int]bool)
	handleMap := make(map[rune]*entry)
	for e := o.front(); e != nil; e, _ = e.nextIndexed() ***REMOVED***
		if e.contractionHandle.n > 0 ***REMOVED***
			handleMap[e.runes[0]] = e
			indexMap[e.contractionHandle.index] = true
		***REMOVED***
	***REMOVED***
	// Verify uniquing.
	if len(indexMap) != 2 ***REMOVED***
		t.Errorf("number of tries is %d; want %d", len(indexMap), 2)
	***REMOVED***
	for _, tt := range contractTest ***REMOVED***
		e, ok := handleMap[[]rune(tt.str)[0]]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		str := tt.str[1:]
		offset, n := lookup(&b.t.ContractTries, e.contractionHandle, []byte(str))
		if len(str) != n ***REMOVED***
			t.Errorf("%s: bytes consumed==%d; want %d", tt.str, n, len(str))
		***REMOVED***
		ce := b.t.ContractElem[offset+e.contractionIndex]
		if want, _ := makeCE(tt.ces[0]); want != ce ***REMOVED***
			t.Errorf("%s: element %X; want %X", tt.str, ce, want)
		***REMOVED***
	***REMOVED***
	if len(b.t.ContractElem) != totalElements ***REMOVED***
		t.Errorf("len(expandElem)==%d; want %d", len(b.t.ContractElem), totalElements)
	***REMOVED***
***REMOVED***
