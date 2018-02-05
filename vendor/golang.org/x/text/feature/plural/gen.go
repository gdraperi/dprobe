// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

// This file generates data for the CLDR plural rules, as defined in
//    http://unicode.org/reports/tr35/tr35-numbers.html#Language_Plural_Rules
//
// We assume a slightly simplified grammar:
//
// 		condition     = and_condition ('or' and_condition)* samples
// 		and_condition = relation ('and' relation)*
// 		relation      = expr ('=' | '!=') range_list
// 		expr          = operand ('%' '10' '0'* )?
// 		operand       = 'n' | 'i' | 'f' | 't' | 'v' | 'w'
// 		range_list    = (range | value) (',' range_list)*
// 		range         = value'..'value
// 		value         = digit+
// 		digit         = 0|1|2|3|4|5|6|7|8|9
//
// 		samples       = ('@integer' sampleList)?
// 		                ('@decimal' sampleList)?
// 		sampleList    = sampleRange (',' sampleRange)* (',' ('…'|'...'))?
// 		sampleRange   = decimalValue ('~' decimalValue)?
// 		decimalValue  = value ('.' value)?
//
//		Symbol	Value
//		n	absolute value of the source number (integer and decimals).
//		i	integer digits of n.
//		v	number of visible fraction digits in n, with trailing zeros.
//		w	number of visible fraction digits in n, without trailing zeros.
//		f	visible fractional digits in n, with trailing zeros.
//		t	visible fractional digits in n, without trailing zeros.
//
// The algorithm for which the data is generated is based on the following
// observations
//
//    - the number of different sets of numbers which the plural rules use to
//      test inclusion is limited,
//    - most numbers that are tested on are < 100
//
// This allows us to define a bitmap for each number < 100 where a bit i
// indicates whether this number is included in some defined set i.
// The function matchPlural in plural.go defines how we can subsequently use
// this data to determine inclusion.
//
// There are a few languages for which this doesn't work. For one Italian and
// Azerbaijan, which both test against numbers > 100 for ordinals and Breton,
// which considers whether numbers are multiples of hundreds. The model here
// could be extended to handle Italian and Azerbaijan fairly easily (by
// considering the numbers 100, 200, 300, ..., 800, 900 in addition to the first
// 100), but for now it seems easier to just hard-code these cases.

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/text/internal"
	"golang.org/x/text/internal/gen"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

var (
	test = flag.Bool("test", false,
		"test existing tables; can be used to compare web data with package data.")
	outputFile     = flag.String("output", "tables.go", "output file")
	outputTestFile = flag.String("testoutput", "data_test.go", "output file")

	draft = flag.String("draft",
		"contributed",
		`Minimal draft requirements (approved, contributed, provisional, unconfirmed).`)
)

func main() ***REMOVED***
	gen.Init()

	const pkg = "plural"

	gen.Repackage("gen_common.go", "common.go", pkg)
	// Read the CLDR zip file.
	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder***REMOVED******REMOVED***
	d.SetDirFilter("supplemental", "main")
	d.SetSectionFilter("numbers", "plurals")
	data, err := d.DecodeZip(r)
	if err != nil ***REMOVED***
		log.Fatalf("DecodeZip: %v", err)
	***REMOVED***

	w := gen.NewCodeWriter()
	defer w.WriteGoFile(*outputFile, pkg)

	gen.WriteCLDRVersion(w)

	genPlurals(w, data)

	w = gen.NewCodeWriter()
	defer w.WriteGoFile(*outputTestFile, pkg)

	genPluralsTests(w, data)
***REMOVED***

type pluralTest struct ***REMOVED***
	locales string   // space-separated list of locales for this test
	form    int      // Use int instead of Form to simplify generation.
	integer []string // Entries of the form \d+ or \d+~\d+
	decimal []string // Entries of the form \f+ or \f+ +~\f+, where f is \d+\.\d+
***REMOVED***

func genPluralsTests(w *gen.CodeWriter, data *cldr.CLDR) ***REMOVED***
	w.WriteType(pluralTest***REMOVED******REMOVED***)

	for _, plurals := range data.Supplemental().Plurals ***REMOVED***
		if plurals.Type == "" ***REMOVED***
			// The empty type is reserved for plural ranges.
			continue
		***REMOVED***
		tests := []pluralTest***REMOVED******REMOVED***

		for _, pRules := range plurals.PluralRules ***REMOVED***
			for _, rule := range pRules.PluralRule ***REMOVED***
				test := pluralTest***REMOVED***
					locales: pRules.Locales,
					form:    int(countMap[rule.Count]),
				***REMOVED***
				scan := bufio.NewScanner(strings.NewReader(rule.Data()))
				scan.Split(splitTokens)
				var p *[]string
				for scan.Scan() ***REMOVED***
					switch t := scan.Text(); t ***REMOVED***
					case "@integer":
						p = &test.integer
					case "@decimal":
						p = &test.decimal
					case ",", "…":
					default:
						if p != nil ***REMOVED***
							*p = append(*p, t)
						***REMOVED***
					***REMOVED***
				***REMOVED***
				tests = append(tests, test)
			***REMOVED***
		***REMOVED***
		w.WriteVar(plurals.Type+"Tests", tests)
	***REMOVED***
***REMOVED***

func genPlurals(w *gen.CodeWriter, data *cldr.CLDR) ***REMOVED***
	for _, plurals := range data.Supplemental().Plurals ***REMOVED***
		if plurals.Type == "" ***REMOVED***
			continue
		***REMOVED***
		// Initialize setMap and inclusionMasks. They are already populated with
		// a few entries to serve as an example and to assign nice numbers to
		// common cases.

		// setMap contains sets of numbers represented by boolean arrays where
		// a true value for element i means that the number i is included.
		setMap := map[[numN]bool]int***REMOVED***
			// The above init func adds an entry for including all numbers.
			[numN]bool***REMOVED***1: true***REMOVED***: 1, // fix ***REMOVED***1***REMOVED*** to a nice value
			[numN]bool***REMOVED***2: true***REMOVED***: 2, // fix ***REMOVED***2***REMOVED*** to a nice value
			[numN]bool***REMOVED***0: true***REMOVED***: 3, // fix ***REMOVED***0***REMOVED*** to a nice value
		***REMOVED***

		// inclusionMasks contains bit masks for every number under numN to
		// indicate in which set the number is included. Bit 1 << x will be set
		// if it is included in set x.
		inclusionMasks := [numN]uint64***REMOVED***
			// Note: these entries are not complete: more bits will be set along the way.
			0: 1 << 3,
			1: 1 << 1,
			2: 1 << 2,
		***REMOVED***

		// Create set ***REMOVED***0..99***REMOVED***. We will assign this set the identifier 0.
		var all [numN]bool
		for i := range all ***REMOVED***
			// Mark number i as being included in the set (which has identifier 0).
			inclusionMasks[i] |= 1 << 0
			// Mark number i as included in the set.
			all[i] = true
		***REMOVED***
		// Register the identifier for the set.
		setMap[all] = 0

		rules := []pluralCheck***REMOVED******REMOVED***
		index := []byte***REMOVED***0***REMOVED***
		langMap := map[int]byte***REMOVED***0: 0***REMOVED*** // From compact language index to index

		for _, pRules := range plurals.PluralRules ***REMOVED***
			// Parse the rules.
			var conds []orCondition
			for _, rule := range pRules.PluralRule ***REMOVED***
				form := countMap[rule.Count]
				conds = parsePluralCondition(conds, rule.Data(), form)
			***REMOVED***
			// Encode the rules.
			for _, c := range conds ***REMOVED***
				// If an or condition only has filters, we create an entry for
				// this filter and the set that contains all values.
				empty := true
				for _, b := range c.used ***REMOVED***
					empty = empty && !b
				***REMOVED***
				if empty ***REMOVED***
					rules = append(rules, pluralCheck***REMOVED***
						cat:   byte(opMod<<opShift) | byte(c.form),
						setID: 0, // all values
					***REMOVED***)
					continue
				***REMOVED***
				// We have some entries with values.
				for i, set := range c.set ***REMOVED***
					if !c.used[i] ***REMOVED***
						continue
					***REMOVED***
					index, ok := setMap[set]
					if !ok ***REMOVED***
						index = len(setMap)
						setMap[set] = index
						for i := range inclusionMasks ***REMOVED***
							if set[i] ***REMOVED***
								inclusionMasks[i] |= 1 << uint64(index)
							***REMOVED***
						***REMOVED***
					***REMOVED***
					rules = append(rules, pluralCheck***REMOVED***
						cat:   byte(i<<opShift | andNext),
						setID: byte(index),
					***REMOVED***)
				***REMOVED***
				// Now set the last entry to the plural form the rule matches.
				rules[len(rules)-1].cat &^= formMask
				rules[len(rules)-1].cat |= byte(c.form)
			***REMOVED***
			// Point the relevant locales to the created entries.
			for _, loc := range strings.Split(pRules.Locales, " ") ***REMOVED***
				if strings.TrimSpace(loc) == "" ***REMOVED***
					continue
				***REMOVED***
				lang, ok := language.CompactIndex(language.MustParse(loc))
				if !ok ***REMOVED***
					log.Printf("No compact index for locale %q", loc)
				***REMOVED***
				langMap[lang] = byte(len(index) - 1)
			***REMOVED***
			index = append(index, byte(len(rules)))
		***REMOVED***
		w.WriteVar(plurals.Type+"Rules", rules)
		w.WriteVar(plurals.Type+"Index", index)
		// Expand the values.
		langToIndex := make([]byte, language.NumCompactTags)
		for i := range langToIndex ***REMOVED***
			for p := i; ; p = int(internal.Parent[p]) ***REMOVED***
				if x, ok := langMap[p]; ok ***REMOVED***
					langToIndex[i] = x
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
		w.WriteVar(plurals.Type+"LangToIndex", langToIndex)
		// Need to convert array to slice because of golang.org/issue/7651.
		// This will allow tables to be dropped when unused. This is especially
		// relevant for the ordinal data, which I suspect won't be used as much.
		w.WriteVar(plurals.Type+"InclusionMasks", inclusionMasks[:])

		if len(rules) > 0xFF ***REMOVED***
			log.Fatalf("Too many entries for rules: %#x", len(rules))
		***REMOVED***
		if len(index) > 0xFF ***REMOVED***
			log.Fatalf("Too many entries for index: %#x", len(index))
		***REMOVED***
		if len(setMap) > 64 ***REMOVED*** // maximum number of bits.
			log.Fatalf("Too many entries for setMap: %d", len(setMap))
		***REMOVED***
		w.WriteComment(
			"Slots used for %s: %X of 0xFF rules; %X of 0xFF indexes; %d of 64 sets",
			plurals.Type, len(rules), len(index), len(setMap))
		// Prevent comment from attaching to the next entry.
		fmt.Fprint(w, "\n\n")
	***REMOVED***
***REMOVED***

type orCondition struct ***REMOVED***
	original string // for debugging

	form Form
	used [32]bool
	set  [32][numN]bool
***REMOVED***

func (o *orCondition) add(op opID, mod int, v []int) (ok bool) ***REMOVED***
	ok = true
	for _, x := range v ***REMOVED***
		if x >= maxMod ***REMOVED***
			ok = false
			break
		***REMOVED***
	***REMOVED***
	for i := 0; i < numN; i++ ***REMOVED***
		m := i
		if mod != 0 ***REMOVED***
			m = i % mod
		***REMOVED***
		if !intIn(m, v) ***REMOVED***
			o.set[op][i] = false
		***REMOVED***
	***REMOVED***
	if ok ***REMOVED***
		o.used[op] = true
	***REMOVED***
	return ok
***REMOVED***

func intIn(x int, a []int) bool ***REMOVED***
	for _, y := range a ***REMOVED***
		if x == y ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

var operandIndex = map[string]opID***REMOVED***
	"i": opI,
	"n": opN,
	"f": opF,
	"v": opV,
	"w": opW,
***REMOVED***

// parsePluralCondition parses the condition of a single pluralRule and appends
// the resulting or conditions to conds.
//
// Example rules:
//   // Category "one" in English: only allow 1 with no visible fraction
//   i = 1 and v = 0 @integer 1
//
//   // Category "few" in Czech: all numbers with visible fractions
//   v != 0   @decimal ...
//
//   // Category "zero" in Latvian: all multiples of 10 or the numbers 11-19 or
//   // numbers with a fraction 11..19 and no trailing zeros.
//   n % 10 = 0 or n % 100 = 11..19 or v = 2 and f % 100 = 11..19 @integer ...
//
// @integer and @decimal are followed by examples and are not relevant for the
// rule itself. The are used here to signal the termination of the rule.
func parsePluralCondition(conds []orCondition, s string, f Form) []orCondition ***REMOVED***
	scan := bufio.NewScanner(strings.NewReader(s))
	scan.Split(splitTokens)
	for ***REMOVED***
		cond := orCondition***REMOVED***original: s, form: f***REMOVED***
		// Set all numbers to be allowed for all number classes and restrict
		// from here on.
		for i := range cond.set ***REMOVED***
			for j := range cond.set[i] ***REMOVED***
				cond.set[i][j] = true
			***REMOVED***
		***REMOVED***
	andLoop:
		for ***REMOVED***
			var token string
			scan.Scan() // Must exist.
			switch class := scan.Text(); class ***REMOVED***
			case "t":
				class = "w" // equal to w for t == 0
				fallthrough
			case "n", "i", "f", "v", "w":
				op := scanToken(scan)
				opCode := operandIndex[class]
				mod := 0
				if op == "%" ***REMOVED***
					opCode |= opMod

					switch v := scanUint(scan); v ***REMOVED***
					case 10, 100:
						mod = v
					case 1000:
						// A more general solution would be to allow checking
						// against multiples of 100 and include entries for the
						// numbers 100..900 in the inclusion masks. At the
						// moment this would only help Azerbaijan and Italian.

						// Italian doesn't use '%', so this must be Azerbaijan.
						cond.used[opAzerbaijan00s] = true
						return append(conds, cond)

					case 1000000:
						cond.used[opBretonM] = true
						return append(conds, cond)

					default:
						log.Fatalf("Modulo value not supported %d", v)
					***REMOVED***
					op = scanToken(scan)
				***REMOVED***
				if op != "=" && op != "!=" ***REMOVED***
					log.Fatalf("Unexpected op %q", op)
				***REMOVED***
				if op == "!=" ***REMOVED***
					opCode |= opNotEqual
				***REMOVED***
				a := []int***REMOVED******REMOVED***
				v := scanUint(scan)
				if class == "w" && v != 0 ***REMOVED***
					log.Fatalf("Must compare against zero for operand type %q", class)
				***REMOVED***
				token = scanToken(scan)
				for ***REMOVED***
					switch token ***REMOVED***
					case "..":
						end := scanUint(scan)
						for ; v <= end; v++ ***REMOVED***
							a = append(a, v)
						***REMOVED***
						token = scanToken(scan)
					default: // ",", "or", "and", "@..."
						a = append(a, v)
					***REMOVED***
					if token != "," ***REMOVED***
						break
					***REMOVED***
					v = scanUint(scan)
					token = scanToken(scan)
				***REMOVED***
				if !cond.add(opCode, mod, a) ***REMOVED***
					// Detected large numbers. As we ruled out Azerbaijan, this
					// must be the many rule for Italian ordinals.
					cond.set[opItalian800] = cond.set[opN]
					cond.used[opItalian800] = true
				***REMOVED***

			case "@integer", "@decimal": // "other" entry: tests only.
				return conds
			default:
				log.Fatalf("Unexpected operand class %q (%s)", class, s)
			***REMOVED***
			switch token ***REMOVED***
			case "or":
				conds = append(conds, cond)
				break andLoop
			case "@integer", "@decimal": // examples
				// There is always an example in practice, so we always terminate here.
				if err := scan.Err(); err != nil ***REMOVED***
					log.Fatal(err)
				***REMOVED***
				return append(conds, cond)
			case "and":
				// keep accumulating
			default:
				log.Fatalf("Unexpected token %q", token)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func scanToken(scan *bufio.Scanner) string ***REMOVED***
	scan.Scan()
	return scan.Text()
***REMOVED***

func scanUint(scan *bufio.Scanner) int ***REMOVED***
	scan.Scan()
	val, err := strconv.ParseUint(scan.Text(), 10, 32)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	return int(val)
***REMOVED***

// splitTokens can be used with bufio.Scanner to tokenize CLDR plural rules.
func splitTokens(data []byte, atEOF bool) (advance int, token []byte, err error) ***REMOVED***
	condTokens := [][]byte***REMOVED***
		[]byte(".."),
		[]byte(","),
		[]byte("!="),
		[]byte("="),
	***REMOVED***
	advance, token, err = bufio.ScanWords(data, atEOF)
	for _, t := range condTokens ***REMOVED***
		if len(t) >= len(token) ***REMOVED***
			continue
		***REMOVED***
		switch p := bytes.Index(token, t); ***REMOVED***
		case p == -1:
		case p == 0:
			advance = len(t)
			token = token[:len(t)]
			return advance - len(token) + len(t), token[:len(t)], err
		case p < advance:
			// Don't split when "=" overlaps "!=".
			if t[0] == '=' && token[p-1] == '!' ***REMOVED***
				continue
			***REMOVED***
			advance = p
			token = token[:p]
		***REMOVED***
	***REMOVED***
	return advance, token, err
***REMOVED***
