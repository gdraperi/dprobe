// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bidi

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/internal/ucd"
	"golang.org/x/text/unicode/norm"
)

var testLevels = flag.Bool("levels", false, "enable testing of levels")

// TestBidiCore performs the tests in BidiTest.txt.
// See http://www.unicode.org/Public/UCD/latest/ucd/BidiTest.txt.
func TestBidiCore(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	r := gen.OpenUCDFile("BidiTest.txt")
	defer r.Close()

	var wantLevels, wantOrder []string
	p := ucd.New(r, ucd.Part(func(p *ucd.Parser) ***REMOVED***
		s := strings.Split(p.String(0), ":")
		switch s[0] ***REMOVED***
		case "Levels":
			wantLevels = strings.Fields(s[1])
		case "Reorder":
			wantOrder = strings.Fields(s[1])
		default:
			log.Fatalf("Unknown part %q.", s[0])
		***REMOVED***
	***REMOVED***))

	for p.Next() ***REMOVED***
		types := []Class***REMOVED******REMOVED***
		for _, s := range p.Strings(0) ***REMOVED***
			types = append(types, bidiClass[s])
		***REMOVED***
		// We ignore the bracketing part of the algorithm.
		pairTypes := make([]bracketType, len(types))
		pairValues := make([]rune, len(types))

		for i := uint(0); i < 3; i++ ***REMOVED***
			if p.Uint(1)&(1<<i) == 0 ***REMOVED***
				continue
			***REMOVED***
			lev := level(int(i) - 1)
			par := newParagraph(types, pairTypes, pairValues, lev)

			if *testLevels ***REMOVED***
				levels := par.getLevels([]int***REMOVED***len(types)***REMOVED***)
				for i, s := range wantLevels ***REMOVED***
					if s == "x" ***REMOVED***
						continue
					***REMOVED***
					l, _ := strconv.ParseUint(s, 10, 8)
					if level(l)&1 != levels[i]&1 ***REMOVED***
						t.Errorf("%s:%d:levels: got %v; want %v", p.String(0), lev, levels, wantLevels)
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***

			order := par.getReordering([]int***REMOVED***len(types)***REMOVED***)
			gotOrder := filterOrder(types, order)
			if got, want := fmt.Sprint(gotOrder), fmt.Sprint(wantOrder); got != want ***REMOVED***
				t.Errorf("%s:%d:order: got %v; want %v\noriginal %v", p.String(0), lev, got, want, order)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if err := p.Err(); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
***REMOVED***

var removeClasses = map[Class]bool***REMOVED***
	LRO: true,
	RLO: true,
	RLE: true,
	LRE: true,
	PDF: true,
	BN:  true,
***REMOVED***

// TestBidiCharacters performs the tests in BidiCharacterTest.txt.
// See http://www.unicode.org/Public/UCD/latest/ucd/BidiCharacterTest.txt
func TestBidiCharacters(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	ucd.Parse(gen.OpenUCDFile("BidiCharacterTest.txt"), func(p *ucd.Parser) ***REMOVED***
		var (
			types      []Class
			pairTypes  []bracketType
			pairValues []rune
			parLevel   level

			wantLevel       = level(p.Int(2))
			wantLevels      = p.Strings(3)
			wantVisualOrder = p.Strings(4)
		)

		switch l := p.Int(1); l ***REMOVED***
		case 0, 1:
			parLevel = level(l)
		case 2:
			parLevel = implicitLevel
		default:
			// Spec says to ignore unknown parts.
		***REMOVED***

		runes := p.Runes(0)

		for _, r := range runes ***REMOVED***
			// Assign the bracket type.
			if d := norm.NFKD.PropertiesString(string(r)).Decomposition(); d != nil ***REMOVED***
				r = []rune(string(d))[0]
			***REMOVED***
			p, _ := LookupRune(r)

			// Assign the class for this rune.
			types = append(types, p.Class())

			switch ***REMOVED***
			case !p.IsBracket():
				pairTypes = append(pairTypes, bpNone)
				pairValues = append(pairValues, 0)
			case p.IsOpeningBracket():
				pairTypes = append(pairTypes, bpOpen)
				pairValues = append(pairValues, r)
			default:
				pairTypes = append(pairTypes, bpClose)
				pairValues = append(pairValues, p.reverseBracket(r))
			***REMOVED***
		***REMOVED***
		par := newParagraph(types, pairTypes, pairValues, parLevel)

		// Test results:
		if got := par.embeddingLevel; got != wantLevel ***REMOVED***
			t.Errorf("%v:level: got %d; want %d", string(runes), got, wantLevel)
		***REMOVED***

		if *testLevels ***REMOVED***
			gotLevels := getLevelStrings(types, par.getLevels([]int***REMOVED***len(types)***REMOVED***))
			if got, want := fmt.Sprint(gotLevels), fmt.Sprint(wantLevels); got != want ***REMOVED***
				t.Errorf("%04X %q:%d: got %v; want %v\nval: %x\npair: %v", runes, string(runes), parLevel, got, want, pairValues, pairTypes)
			***REMOVED***
		***REMOVED***

		order := par.getReordering([]int***REMOVED***len(types)***REMOVED***)
		order = filterOrder(types, order)
		if got, want := fmt.Sprint(order), fmt.Sprint(wantVisualOrder); got != want ***REMOVED***
			t.Errorf("%04X %q:%d: got %v; want %v\ngot order: %s", runes, string(runes), parLevel, got, want, reorder(runes, order))
		***REMOVED***
	***REMOVED***)
***REMOVED***

func getLevelStrings(cl []Class, levels []level) []string ***REMOVED***
	var results []string
	for i, l := range levels ***REMOVED***
		if !removeClasses[cl[i]] ***REMOVED***
			results = append(results, fmt.Sprint(l))
		***REMOVED*** else ***REMOVED***
			results = append(results, "x")
		***REMOVED***
	***REMOVED***
	return results
***REMOVED***

func filterOrder(cl []Class, order []int) []int ***REMOVED***
	no := []int***REMOVED******REMOVED***
	for _, o := range order ***REMOVED***
		if !removeClasses[cl[o]] ***REMOVED***
			no = append(no, o)
		***REMOVED***
	***REMOVED***
	return no
***REMOVED***

func reorder(r []rune, order []int) string ***REMOVED***
	nr := make([]rune, len(order))
	for i, o := range order ***REMOVED***
		nr[i] = r[o]
	***REMOVED***
	return string(nr)
***REMOVED***

// bidiClass names and codes taken from class "bc" in
// http://www.unicode.org/Public/8.0.0/ucd/PropertyValueAliases.txt
var bidiClass = map[string]Class***REMOVED***
	"AL":  AL,  // classArabicLetter,
	"AN":  AN,  // classArabicNumber,
	"B":   B,   // classParagraphSeparator,
	"BN":  BN,  // classBoundaryNeutral,
	"CS":  CS,  // classCommonSeparator,
	"EN":  EN,  // classEuropeanNumber,
	"ES":  ES,  // classEuropeanSeparator,
	"ET":  ET,  // classEuropeanTerminator,
	"L":   L,   // classLeftToRight,
	"NSM": NSM, // classNonspacingMark,
	"ON":  ON,  // classOtherNeutral,
	"R":   R,   // classRightToLeft,
	"S":   S,   // classSegmentSeparator,
	"WS":  WS,  // classWhiteSpace,

	"LRO": LRO, // classLeftToRightOverride,
	"RLO": RLO, // classRightToLeftOverride,
	"LRE": LRE, // classLeftToRightEmbedding,
	"RLE": RLE, // classRightToLeftEmbedding,
	"PDF": PDF, // classPopDirectionalFormat,
	"LRI": LRI, // classLeftToRightIsolate,
	"RLI": RLI, // classRightToLeftIsolate,
	"FSI": FSI, // classFirstStrongIsolate,
	"PDI": PDI, // classPopDirectionalIsolate,
***REMOVED***
