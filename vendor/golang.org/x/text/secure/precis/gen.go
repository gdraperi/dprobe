// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Unicode table generator.
// Data read from the web.

// +build ignore

package main

import (
	"flag"
	"log"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/triegen"
	"golang.org/x/text/internal/ucd"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/rangetable"
)

var outputFile = flag.String("output", "tables.go", "output file for generated tables; default tables.go")

var assigned, disallowedRunes *unicode.RangeTable

var runeCategory = map[rune]category***REMOVED******REMOVED***

var overrides = map[category]category***REMOVED***
	viramaModifier: viramaJoinT,
	greek:          greekJoinT,
	hebrew:         hebrewJoinT,
***REMOVED***

func setCategory(r rune, cat category) ***REMOVED***
	if c, ok := runeCategory[r]; ok ***REMOVED***
		if override, ok := overrides[c]; cat == joiningT && ok ***REMOVED***
			cat = override
		***REMOVED*** else ***REMOVED***
			log.Fatalf("%U: multiple categories for rune (%v and %v)", r, c, cat)
		***REMOVED***
	***REMOVED***
	runeCategory[r] = cat
***REMOVED***

func init() ***REMOVED***
	if numCategories > 1<<propShift ***REMOVED***
		log.Fatalf("Number of categories is %d; may at most be %d", numCategories, 1<<propShift)
	***REMOVED***
***REMOVED***

func main() ***REMOVED***
	gen.Init()

	// Load data
	runes := []rune***REMOVED******REMOVED***
	// PrecisIgnorableProperties: https://tools.ietf.org/html/rfc7564#section-9.13
	ucd.Parse(gen.OpenUCDFile("DerivedCoreProperties.txt"), func(p *ucd.Parser) ***REMOVED***
		if p.String(1) == "Default_Ignorable_Code_Point" ***REMOVED***
			runes = append(runes, p.Rune(0))
		***REMOVED***
	***REMOVED***)
	ucd.Parse(gen.OpenUCDFile("PropList.txt"), func(p *ucd.Parser) ***REMOVED***
		switch p.String(1) ***REMOVED***
		case "Noncharacter_Code_Point":
			runes = append(runes, p.Rune(0))
		***REMOVED***
	***REMOVED***)
	// OldHangulJamo: https://tools.ietf.org/html/rfc5892#section-2.9
	ucd.Parse(gen.OpenUCDFile("HangulSyllableType.txt"), func(p *ucd.Parser) ***REMOVED***
		switch p.String(1) ***REMOVED***
		case "L", "V", "T":
			runes = append(runes, p.Rune(0))
		***REMOVED***
	***REMOVED***)

	disallowedRunes = rangetable.New(runes...)
	assigned = rangetable.Assigned(unicode.Version)

	// Load category data.
	runeCategory['l'] = latinSmallL
	ucd.Parse(gen.OpenUCDFile("UnicodeData.txt"), func(p *ucd.Parser) ***REMOVED***
		const cccVirama = 9
		if p.Int(ucd.CanonicalCombiningClass) == cccVirama ***REMOVED***
			setCategory(p.Rune(0), viramaModifier)
		***REMOVED***
	***REMOVED***)
	ucd.Parse(gen.OpenUCDFile("Scripts.txt"), func(p *ucd.Parser) ***REMOVED***
		switch p.String(1) ***REMOVED***
		case "Greek":
			setCategory(p.Rune(0), greek)
		case "Hebrew":
			setCategory(p.Rune(0), hebrew)
		case "Hiragana", "Katakana", "Han":
			setCategory(p.Rune(0), japanese)
		***REMOVED***
	***REMOVED***)

	// Set the rule categories associated with exceptions. This overrides any
	// previously set categories. The original categories are manually
	// reintroduced in the categoryTransitions table.
	for r, e := range exceptions ***REMOVED***
		if e.cat != 0 ***REMOVED***
			runeCategory[r] = e.cat
		***REMOVED***
	***REMOVED***
	cat := map[string]category***REMOVED***
		"L": joiningL,
		"D": joiningD,
		"T": joiningT,

		"R": joiningR,
	***REMOVED***
	ucd.Parse(gen.OpenUCDFile("extracted/DerivedJoiningType.txt"), func(p *ucd.Parser) ***REMOVED***
		switch v := p.String(1); v ***REMOVED***
		case "L", "D", "T", "R":
			setCategory(p.Rune(0), cat[v])
		***REMOVED***
	***REMOVED***)

	writeTables()
	gen.Repackage("gen_trieval.go", "trieval.go", "precis")
***REMOVED***

type exception struct ***REMOVED***
	prop property
	cat  category
***REMOVED***

func init() ***REMOVED***
	// Programmatically add the Arabic and Indic digits to the exceptions map.
	// See comment in the exceptions map below why these are marked disallowed.
	for i := rune(0); i <= 9; i++ ***REMOVED***
		exceptions[0x0660+i] = exception***REMOVED***
			prop: disallowed,
			cat:  arabicIndicDigit,
		***REMOVED***
		exceptions[0x06F0+i] = exception***REMOVED***
			prop: disallowed,
			cat:  extendedArabicIndicDigit,
		***REMOVED***
	***REMOVED***
***REMOVED***

// The Exceptions class as defined in RFC 5892
// https://tools.ietf.org/html/rfc5892#section-2.6
var exceptions = map[rune]exception***REMOVED***
	0x00DF: ***REMOVED***prop: pValid***REMOVED***,
	0x03C2: ***REMOVED***prop: pValid***REMOVED***,
	0x06FD: ***REMOVED***prop: pValid***REMOVED***,
	0x06FE: ***REMOVED***prop: pValid***REMOVED***,
	0x0F0B: ***REMOVED***prop: pValid***REMOVED***,
	0x3007: ***REMOVED***prop: pValid***REMOVED***,

	// ContextO|J rules are marked as disallowed, taking a "guilty until proven
	// innocent" approach. The main reason for this is that the check for
	// whether a context rule should be applied can be moved to the logic for
	// handing disallowed runes, taken it off the common path. The exception to
	// this rule is for katakanaMiddleDot, as the rule logic is handled without
	// using a rule function.

	// ContextJ (Join control)
	0x200C: ***REMOVED***prop: disallowed, cat: zeroWidthNonJoiner***REMOVED***,
	0x200D: ***REMOVED***prop: disallowed, cat: zeroWidthJoiner***REMOVED***,

	// ContextO
	0x00B7: ***REMOVED***prop: disallowed, cat: middleDot***REMOVED***,
	0x0375: ***REMOVED***prop: disallowed, cat: greekLowerNumeralSign***REMOVED***,
	0x05F3: ***REMOVED***prop: disallowed, cat: hebrewPreceding***REMOVED***, // punctuation Geresh
	0x05F4: ***REMOVED***prop: disallowed, cat: hebrewPreceding***REMOVED***, // punctuation Gershayim
	0x30FB: ***REMOVED***prop: pValid, cat: katakanaMiddleDot***REMOVED***,

	// These are officially ContextO, but the implementation does not require
	// special treatment of these, so we simply mark them as valid.
	0x0660: ***REMOVED***prop: pValid***REMOVED***,
	0x0661: ***REMOVED***prop: pValid***REMOVED***,
	0x0662: ***REMOVED***prop: pValid***REMOVED***,
	0x0663: ***REMOVED***prop: pValid***REMOVED***,
	0x0664: ***REMOVED***prop: pValid***REMOVED***,
	0x0665: ***REMOVED***prop: pValid***REMOVED***,
	0x0666: ***REMOVED***prop: pValid***REMOVED***,
	0x0667: ***REMOVED***prop: pValid***REMOVED***,
	0x0668: ***REMOVED***prop: pValid***REMOVED***,
	0x0669: ***REMOVED***prop: pValid***REMOVED***,
	0x06F0: ***REMOVED***prop: pValid***REMOVED***,
	0x06F1: ***REMOVED***prop: pValid***REMOVED***,
	0x06F2: ***REMOVED***prop: pValid***REMOVED***,
	0x06F3: ***REMOVED***prop: pValid***REMOVED***,
	0x06F4: ***REMOVED***prop: pValid***REMOVED***,
	0x06F5: ***REMOVED***prop: pValid***REMOVED***,
	0x06F6: ***REMOVED***prop: pValid***REMOVED***,
	0x06F7: ***REMOVED***prop: pValid***REMOVED***,
	0x06F8: ***REMOVED***prop: pValid***REMOVED***,
	0x06F9: ***REMOVED***prop: pValid***REMOVED***,

	0x0640: ***REMOVED***prop: disallowed***REMOVED***,
	0x07FA: ***REMOVED***prop: disallowed***REMOVED***,
	0x302E: ***REMOVED***prop: disallowed***REMOVED***,
	0x302F: ***REMOVED***prop: disallowed***REMOVED***,
	0x3031: ***REMOVED***prop: disallowed***REMOVED***,
	0x3032: ***REMOVED***prop: disallowed***REMOVED***,
	0x3033: ***REMOVED***prop: disallowed***REMOVED***,
	0x3034: ***REMOVED***prop: disallowed***REMOVED***,
	0x3035: ***REMOVED***prop: disallowed***REMOVED***,
	0x303B: ***REMOVED***prop: disallowed***REMOVED***,
***REMOVED***

// LetterDigits: https://tools.ietf.org/html/rfc5892#section-2.1
// r in ***REMOVED***Ll, Lu, Lo, Nd, Lm, Mn, Mc***REMOVED***.
func isLetterDigits(r rune) bool ***REMOVED***
	return unicode.In(r,
		unicode.Ll, unicode.Lu, unicode.Lm, unicode.Lo, // Letters
		unicode.Mn, unicode.Mc, // Modifiers
		unicode.Nd, // Digits
	)
***REMOVED***

func isIdDisAndFreePVal(r rune) bool ***REMOVED***
	return unicode.In(r,
		// OtherLetterDigits: https://tools.ietf.org/html/rfc7564#section-9.18
		// r in in ***REMOVED***Lt, Nl, No, Me***REMOVED***
		unicode.Lt, unicode.Nl, unicode.No, // Other letters / numbers
		unicode.Me, // Modifiers

		// Spaces: https://tools.ietf.org/html/rfc7564#section-9.14
		// r in in ***REMOVED***Zs***REMOVED***
		unicode.Zs,

		// Symbols: https://tools.ietf.org/html/rfc7564#section-9.15
		// r in ***REMOVED***Sm, Sc, Sk, So***REMOVED***
		unicode.Sm, unicode.Sc, unicode.Sk, unicode.So,

		// Punctuation: https://tools.ietf.org/html/rfc7564#section-9.16
		// r in ***REMOVED***Pc, Pd, Ps, Pe, Pi, Pf, Po***REMOVED***
		unicode.Pc, unicode.Pd, unicode.Ps, unicode.Pe,
		unicode.Pi, unicode.Pf, unicode.Po,
	)
***REMOVED***

// HasCompat: https://tools.ietf.org/html/rfc7564#section-9.17
func hasCompat(r rune) bool ***REMOVED***
	return !norm.NFKC.IsNormalString(string(r))
***REMOVED***

// From https://tools.ietf.org/html/rfc5892:
//
// If .cp. .in. Exceptions Then Exceptions(cp);
//   Else If .cp. .in. BackwardCompatible Then BackwardCompatible(cp);
//   Else If .cp. .in. Unassigned Then UNASSIGNED;
//   Else If .cp. .in. ASCII7 Then PVALID;
//   Else If .cp. .in. JoinControl Then CONTEXTJ;
//   Else If .cp. .in. OldHangulJamo Then DISALLOWED;
//   Else If .cp. .in. PrecisIgnorableProperties Then DISALLOWED;
//   Else If .cp. .in. Controls Then DISALLOWED;
//   Else If .cp. .in. HasCompat Then ID_DIS or FREE_PVAL;
//   Else If .cp. .in. LetterDigits Then PVALID;
//   Else If .cp. .in. OtherLetterDigits Then ID_DIS or FREE_PVAL;
//   Else If .cp. .in. Spaces Then ID_DIS or FREE_PVAL;
//   Else If .cp. .in. Symbols Then ID_DIS or FREE_PVAL;
//   Else If .cp. .in. Punctuation Then ID_DIS or FREE_PVAL;
//   Else DISALLOWED;

func writeTables() ***REMOVED***
	propTrie := triegen.NewTrie("derivedProperties")
	w := gen.NewCodeWriter()
	defer w.WriteVersionedGoFile(*outputFile, "precis")
	gen.WriteUnicodeVersion(w)

	// Iterate over all the runes...
	for i := rune(0); i < unicode.MaxRune; i++ ***REMOVED***
		r := rune(i)

		if !utf8.ValidRune(r) ***REMOVED***
			continue
		***REMOVED***

		e, ok := exceptions[i]
		p := e.prop
		switch ***REMOVED***
		case ok:
		case !unicode.In(r, assigned):
			p = unassigned
		case r >= 0x0021 && r <= 0x007e: // Is ASCII 7
			p = pValid
		case unicode.In(r, disallowedRunes, unicode.Cc):
			p = disallowed
		case hasCompat(r):
			p = idDisOrFreePVal
		case isLetterDigits(r):
			p = pValid
		case isIdDisAndFreePVal(r):
			p = idDisOrFreePVal
		default:
			p = disallowed
		***REMOVED***
		cat := runeCategory[r]
		// Don't set category for runes that are disallowed.
		if p == disallowed ***REMOVED***
			cat = exceptions[r].cat
		***REMOVED***
		propTrie.Insert(r, uint64(p)|uint64(cat))
	***REMOVED***
	sz, err := propTrie.Gen(w)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	w.Size += sz
***REMOVED***
