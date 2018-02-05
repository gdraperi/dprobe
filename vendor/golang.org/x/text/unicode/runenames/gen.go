// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"log"
	"strings"
	"unicode"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/ucd"
)

// snippet is a slice of data; data is the concatenation of all of the names.
type snippet struct ***REMOVED***
	offset int
	length int
	s      string
***REMOVED***

func makeTable0EntryDirect(rOffset, rLength, dOffset, dLength int) uint64 ***REMOVED***
	if rOffset >= 1<<bitsRuneOffset ***REMOVED***
		log.Fatalf("makeTable0EntryDirect: rOffset %d is too large", rOffset)
	***REMOVED***
	if rLength >= 1<<bitsRuneLength ***REMOVED***
		log.Fatalf("makeTable0EntryDirect: rLength %d is too large", rLength)
	***REMOVED***
	if dOffset >= 1<<bitsDataOffset ***REMOVED***
		log.Fatalf("makeTable0EntryDirect: dOffset %d is too large", dOffset)
	***REMOVED***
	if dLength >= 1<<bitsRuneLength ***REMOVED***
		log.Fatalf("makeTable0EntryDirect: dLength %d is too large", dLength)
	***REMOVED***
	return uint64(rOffset)<<shiftRuneOffset |
		uint64(rLength)<<shiftRuneLength |
		uint64(dOffset)<<shiftDataOffset |
		uint64(dLength)<<shiftDataLength |
		1 // Direct bit.
***REMOVED***

func makeTable0EntryIndirect(rOffset, rLength, dBase, t1Offset int) uint64 ***REMOVED***
	if rOffset >= 1<<bitsRuneOffset ***REMOVED***
		log.Fatalf("makeTable0EntryIndirect: rOffset %d is too large", rOffset)
	***REMOVED***
	if rLength >= 1<<bitsRuneLength ***REMOVED***
		log.Fatalf("makeTable0EntryIndirect: rLength %d is too large", rLength)
	***REMOVED***
	if dBase >= 1<<bitsDataBase ***REMOVED***
		log.Fatalf("makeTable0EntryIndirect: dBase %d is too large", dBase)
	***REMOVED***
	if t1Offset >= 1<<bitsTable1Offset ***REMOVED***
		log.Fatalf("makeTable0EntryIndirect: t1Offset %d is too large", t1Offset)
	***REMOVED***
	return uint64(rOffset)<<shiftRuneOffset |
		uint64(rLength)<<shiftRuneLength |
		uint64(dBase)<<shiftDataBase |
		uint64(t1Offset)<<shiftTable1Offset |
		0 // Direct bit.
***REMOVED***

func makeTable1Entry(x int) uint16 ***REMOVED***
	if x < 0 || 0xffff < x ***REMOVED***
		log.Fatalf("makeTable1Entry: entry %d is out of range", x)
	***REMOVED***
	return uint16(x)
***REMOVED***

var (
	data     []byte
	snippets = make([]snippet, 1+unicode.MaxRune)
)

func main() ***REMOVED***
	gen.Init()

	names, counts := parse()
	appendRepeatNames(names, counts)
	appendUniqueNames(names, counts)

	table0, table1 := makeTables()

	gen.Repackage("gen_bits.go", "bits.go", "runenames")

	w := gen.NewCodeWriter()
	w.WriteVar("table0", table0)
	w.WriteVar("table1", table1)
	w.WriteConst("data", string(data))
	w.WriteGoFile("tables.go", "runenames")
***REMOVED***

func parse() (names []string, counts map[string]int) ***REMOVED***
	names = make([]string, 1+unicode.MaxRune)
	counts = map[string]int***REMOVED******REMOVED***
	ucd.Parse(gen.OpenUCDFile("UnicodeData.txt"), func(p *ucd.Parser) ***REMOVED***
		r, s := p.Rune(0), p.String(ucd.Name)
		if s == "" ***REMOVED***
			return
		***REMOVED***
		if s[0] == '<' ***REMOVED***
			const first = ", First>"
			if i := strings.Index(s, first); i >= 0 ***REMOVED***
				s = s[:i] + ">"
			***REMOVED***
		***REMOVED***
		names[r] = s
		counts[s]++
	***REMOVED***)
	return names, counts
***REMOVED***

func appendRepeatNames(names []string, counts map[string]int) ***REMOVED***
	alreadySeen := map[string]snippet***REMOVED******REMOVED***
	for r, s := range names ***REMOVED***
		if s == "" || counts[s] == 1 ***REMOVED***
			continue
		***REMOVED***
		if s[0] != '<' ***REMOVED***
			log.Fatalf("Repeated name %q does not start with a '<'", s)
		***REMOVED***

		if z, ok := alreadySeen[s]; ok ***REMOVED***
			snippets[r] = z
			continue
		***REMOVED***

		z := snippet***REMOVED***
			offset: len(data),
			length: len(s),
			s:      s,
		***REMOVED***
		data = append(data, s...)
		snippets[r] = z
		alreadySeen[s] = z
	***REMOVED***
***REMOVED***

func appendUniqueNames(names []string, counts map[string]int) ***REMOVED***
	for r, s := range names ***REMOVED***
		if s == "" || counts[s] != 1 ***REMOVED***
			continue
		***REMOVED***
		if s[0] == '<' ***REMOVED***
			log.Fatalf("Unique name %q starts with a '<'", s)
		***REMOVED***

		z := snippet***REMOVED***
			offset: len(data),
			length: len(s),
			s:      s,
		***REMOVED***
		data = append(data, s...)
		snippets[r] = z
	***REMOVED***
***REMOVED***

func makeTables() (table0 []uint64, table1 []uint16) ***REMOVED***
	for i := 0; i < len(snippets); ***REMOVED***
		zi := snippets[i]
		if zi == (snippet***REMOVED******REMOVED***) ***REMOVED***
			i++
			continue
		***REMOVED***

		// Look for repeat names. If we have one, we only need a table0 entry.
		j := i + 1
		for ; j < len(snippets) && zi == snippets[j]; j++ ***REMOVED***
		***REMOVED***
		if j > i+1 ***REMOVED***
			table0 = append(table0, makeTable0EntryDirect(i, j-i, zi.offset, zi.length))
			i = j
			continue
		***REMOVED***

		// Otherwise, we have a run of unique names. We need one table0 entry
		// and two or more table1 entries.
		base := zi.offset &^ (1<<dataBaseUnit - 1)
		t1Offset := len(table1) + 1
		table1 = append(table1, makeTable1Entry(zi.offset-base))
		table1 = append(table1, makeTable1Entry(zi.offset+zi.length-base))
		for ; j < len(snippets) && snippets[j] != (snippet***REMOVED******REMOVED***); j++ ***REMOVED***
			zj := snippets[j]
			if data[zj.offset] == '<' ***REMOVED***
				break
			***REMOVED***
			table1 = append(table1, makeTable1Entry(zj.offset+zj.length-base))
		***REMOVED***
		table0 = append(table0, makeTable0EntryIndirect(i, j-i, base>>dataBaseUnit, t1Offset))
		i = j
	***REMOVED***
	return table0, table1
***REMOVED***
