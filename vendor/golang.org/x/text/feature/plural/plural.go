// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go gen_common.go

// Package plural provides utilities for handling linguistic plurals in text.
//
// The definitions in this package are based on the plural rule handling defined
// in CLDR. See
// http://unicode.org/reports/tr35/tr35-numbers.html#Language_Plural_Rules for
// details.
package plural

import (
	"golang.org/x/text/internal/number"
	"golang.org/x/text/language"
)

// Rules defines the plural rules for all languages for a certain plural type.
//
//
// This package is UNDER CONSTRUCTION and its API may change.
type Rules struct ***REMOVED***
	rules          []pluralCheck
	index          []byte
	langToIndex    []byte
	inclusionMasks []uint64
***REMOVED***

var (
	// Cardinal defines the plural rules for numbers indicating quantities.
	Cardinal *Rules = cardinal

	// Ordinal defines the plural rules for numbers indicating position
	// (first, second, etc.).
	Ordinal *Rules = ordinal

	ordinal = &Rules***REMOVED***
		ordinalRules,
		ordinalIndex,
		ordinalLangToIndex,
		ordinalInclusionMasks[:],
	***REMOVED***

	cardinal = &Rules***REMOVED***
		cardinalRules,
		cardinalIndex,
		cardinalLangToIndex,
		cardinalInclusionMasks[:],
	***REMOVED***
)

// getIntApprox converts the digits in slice digits[start:end] to an integer
// according to the following rules:
//	- Let i be asInt(digits[start:end]), where out-of-range digits are assumed
//	  to be zero.
//	- Result n is big if i / 10^nMod > 1.
//	- Otherwise the result is i % 10^nMod.
//
// For example, if digits is ***REMOVED***1, 2, 3***REMOVED*** and start:end is 0:5, then the result
// for various values of nMod is:
//	- when nMod == 2, n == big
//	- when nMod == 3, n == big
//	- when nMod == 4, n == big
//	- when nMod == 5, n == 12300
//	- when nMod == 6, n == 12300
//	- when nMod == 7, n == 12300
func getIntApprox(digits []byte, start, end, nMod, big int) (n int) ***REMOVED***
	// Leading 0 digits just result in 0.
	p := start
	if p < 0 ***REMOVED***
		p = 0
	***REMOVED***
	// Range only over the part for which we have digits.
	mid := end
	if mid >= len(digits) ***REMOVED***
		mid = len(digits)
	***REMOVED***
	// Check digits more significant that nMod.
	if q := end - nMod; q > 0 ***REMOVED***
		if q > mid ***REMOVED***
			q = mid
		***REMOVED***
		for ; p < q; p++ ***REMOVED***
			if digits[p] != 0 ***REMOVED***
				return big
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for ; p < mid; p++ ***REMOVED***
		n = 10*n + int(digits[p])
	***REMOVED***
	// Multiply for trailing zeros.
	for ; p < end; p++ ***REMOVED***
		n *= 10
	***REMOVED***
	return n
***REMOVED***

// MatchDigits computes the plural form for the given language and the given
// decimal floating point digits. The digits are stored in big-endian order and
// are of value byte(0) - byte(9). The floating point position is indicated by
// exp and the number of visible decimals is scale. All leading and trailing
// zeros may be omitted from digits.
//
// The following table contains examples of possible arguments to represent
// the given numbers.
//      decimal    digits              exp    scale
//      123        []byte***REMOVED***1, 2, 3***REMOVED***     3      0
//      123.4      []byte***REMOVED***1, 2, 3, 4***REMOVED***  3      1
//      123.40     []byte***REMOVED***1, 2, 3, 4***REMOVED***  3      2
//      100000     []byte***REMOVED***1***REMOVED***           6      0
//      100000.00  []byte***REMOVED***1***REMOVED***           6      3
func (p *Rules) MatchDigits(t language.Tag, digits []byte, exp, scale int) Form ***REMOVED***
	index, _ := language.CompactIndex(t)

	// Differentiate up to including mod 1000000 for the integer part.
	n := getIntApprox(digits, 0, exp, 6, 1000000)

	// Differentiate up to including mod 100 for the fractional part.
	f := getIntApprox(digits, exp, exp+scale, 2, 100)

	return matchPlural(p, index, n, f, scale)
***REMOVED***

func (p *Rules) matchDisplayDigits(t language.Tag, d *number.Digits) (Form, int) ***REMOVED***
	n := getIntApprox(d.Digits, 0, int(d.Exp), 6, 1000000)
	return p.MatchDigits(t, d.Digits, int(d.Exp), d.NumFracDigits()), n
***REMOVED***

func validForms(p *Rules, t language.Tag) (forms []Form) ***REMOVED***
	index, _ := language.CompactIndex(t)
	offset := p.langToIndex[index]
	rules := p.rules[p.index[offset]:p.index[offset+1]]

	forms = append(forms, Other)
	last := Other
	for _, r := range rules ***REMOVED***
		if cat := Form(r.cat & formMask); cat != andNext && last != cat ***REMOVED***
			forms = append(forms, cat)
			last = cat
		***REMOVED***
	***REMOVED***
	return forms
***REMOVED***

func (p *Rules) matchComponents(t language.Tag, n, f, scale int) Form ***REMOVED***
	index, _ := language.CompactIndex(t)
	return matchPlural(p, index, n, f, scale)
***REMOVED***

// MatchPlural returns the plural form for the given language and plural
// operands (as defined in
// http://unicode.org/reports/tr35/tr35-numbers.html#Language_Plural_Rules):
//  where
//  	n  absolute value of the source number (integer and decimals)
//  input
//  	i  integer digits of n.
//  	v  number of visible fraction digits in n, with trailing zeros.
//  	w  number of visible fraction digits in n, without trailing zeros.
//  	f  visible fractional digits in n, with trailing zeros (f = t * 10^(v-w))
//  	t  visible fractional digits in n, without trailing zeros.
//
// If any of the operand values is too large to fit in an int, it is okay to
// pass the value modulo 10,000,000.
func (p *Rules) MatchPlural(lang language.Tag, i, v, w, f, t int) Form ***REMOVED***
	index, _ := language.CompactIndex(lang)
	return matchPlural(p, index, i, f, v)
***REMOVED***

func matchPlural(p *Rules, index int, n, f, v int) Form ***REMOVED***
	nMask := p.inclusionMasks[n%maxMod]
	// Compute the fMask inline in the rules below, as it is relatively rare.
	// fMask := p.inclusionMasks[f%maxMod]
	vMask := p.inclusionMasks[v%maxMod]

	// Do the matching
	offset := p.langToIndex[index]
	rules := p.rules[p.index[offset]:p.index[offset+1]]
	for i := 0; i < len(rules); i++ ***REMOVED***
		rule := rules[i]
		setBit := uint64(1 << rule.setID)
		var skip bool
		switch op := opID(rule.cat >> opShift); op ***REMOVED***
		case opI: // i = x
			skip = n >= numN || nMask&setBit == 0

		case opI | opNotEqual: // i != x
			skip = n < numN && nMask&setBit != 0

		case opI | opMod: // i % m = x
			skip = nMask&setBit == 0

		case opI | opMod | opNotEqual: // i % m != x
			skip = nMask&setBit != 0

		case opN: // n = x
			skip = f != 0 || n >= numN || nMask&setBit == 0

		case opN | opNotEqual: // n != x
			skip = f == 0 && n < numN && nMask&setBit != 0

		case opN | opMod: // n % m = x
			skip = f != 0 || nMask&setBit == 0

		case opN | opMod | opNotEqual: // n % m != x
			skip = f == 0 && nMask&setBit != 0

		case opF: // f = x
			skip = f >= numN || p.inclusionMasks[f%maxMod]&setBit == 0

		case opF | opNotEqual: // f != x
			skip = f < numN && p.inclusionMasks[f%maxMod]&setBit != 0

		case opF | opMod: // f % m = x
			skip = p.inclusionMasks[f%maxMod]&setBit == 0

		case opF | opMod | opNotEqual: // f % m != x
			skip = p.inclusionMasks[f%maxMod]&setBit != 0

		case opV: // v = x
			skip = v < numN && vMask&setBit == 0

		case opV | opNotEqual: // v != x
			skip = v < numN && vMask&setBit != 0

		case opW: // w == 0
			skip = f != 0

		case opW | opNotEqual: // w != 0
			skip = f == 0

		// Hard-wired rules that cannot be handled by our algorithm.

		case opBretonM:
			skip = f != 0 || n == 0 || n%1000000 != 0

		case opAzerbaijan00s:
			// 100,200,300,400,500,600,700,800,900
			skip = n == 0 || n >= 1000 || n%100 != 0

		case opItalian800:
			skip = (f != 0 || n >= numN || nMask&setBit == 0) && n != 800
		***REMOVED***
		if skip ***REMOVED***
			// advance over AND entries.
			for ; i < len(rules) && rules[i].cat&formMask == andNext; i++ ***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
		// return if we have a final entry.
		if cat := rule.cat & formMask; cat != andNext ***REMOVED***
			return Form(cat)
		***REMOVED***
	***REMOVED***
	return Other
***REMOVED***
