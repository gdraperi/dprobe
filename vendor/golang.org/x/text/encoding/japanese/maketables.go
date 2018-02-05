// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

// This program generates tables.go:
//	go run maketables.go | gofmt > tables.go

// TODO: Emoji extensions?
// http://www.unicode.org/faq/emoji_dingbats.html
// http://www.unicode.org/Public/UNIDATA/EmojiSources.txt

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

type entry struct ***REMOVED***
	jisCode, table int
***REMOVED***

func main() ***REMOVED***
	fmt.Printf("// generated by go run maketables.go; DO NOT EDIT\n\n")
	fmt.Printf("// Package japanese provides Japanese encodings such as EUC-JP and Shift JIS.\n")
	fmt.Printf(`package japanese // import "golang.org/x/text/encoding/japanese"` + "\n\n")

	reverse := [65536]entry***REMOVED******REMOVED***
	for i := range reverse ***REMOVED***
		reverse[i].table = -1
	***REMOVED***

	tables := []struct ***REMOVED***
		url  string
		name string
	***REMOVED******REMOVED***
		***REMOVED***"http://encoding.spec.whatwg.org/index-jis0208.txt", "0208"***REMOVED***,
		***REMOVED***"http://encoding.spec.whatwg.org/index-jis0212.txt", "0212"***REMOVED***,
	***REMOVED***
	for i, table := range tables ***REMOVED***
		res, err := http.Get(table.url)
		if err != nil ***REMOVED***
			log.Fatalf("%q: Get: %v", table.url, err)
		***REMOVED***
		defer res.Body.Close()

		mapping := [65536]uint16***REMOVED******REMOVED***

		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() ***REMOVED***
			s := strings.TrimSpace(scanner.Text())
			if s == "" || s[0] == '#' ***REMOVED***
				continue
			***REMOVED***
			x, y := 0, uint16(0)
			if _, err := fmt.Sscanf(s, "%d 0x%x", &x, &y); err != nil ***REMOVED***
				log.Fatalf("%q: could not parse %q", table.url, s)
			***REMOVED***
			if x < 0 || 120*94 <= x ***REMOVED***
				log.Fatalf("%q: JIS code %d is out of range", table.url, x)
			***REMOVED***
			mapping[x] = y
			if reverse[y].table == -1 ***REMOVED***
				reverse[y] = entry***REMOVED***jisCode: x, table: i***REMOVED***
			***REMOVED***
		***REMOVED***
		if err := scanner.Err(); err != nil ***REMOVED***
			log.Fatalf("%q: scanner error: %v", table.url, err)
		***REMOVED***

		fmt.Printf("// jis%sDecode is the decoding table from JIS %s code to Unicode.\n// It is defined at %s\n",
			table.name, table.name, table.url)
		fmt.Printf("var jis%sDecode = [...]uint16***REMOVED***\n", table.name)
		for i, m := range mapping ***REMOVED***
			if m != 0 ***REMOVED***
				fmt.Printf("\t%d: 0x%04X,\n", i, m)
			***REMOVED***
		***REMOVED***
		fmt.Printf("***REMOVED***\n\n")
	***REMOVED***

	// Any run of at least separation continuous zero entries in the reverse map will
	// be a separate encode table.
	const separation = 1024

	intervals := []interval(nil)
	low, high := -1, -1
	for i, v := range reverse ***REMOVED***
		if v.table == -1 ***REMOVED***
			continue
		***REMOVED***
		if low < 0 ***REMOVED***
			low = i
		***REMOVED*** else if i-high >= separation ***REMOVED***
			if high >= 0 ***REMOVED***
				intervals = append(intervals, interval***REMOVED***low, high***REMOVED***)
			***REMOVED***
			low = i
		***REMOVED***
		high = i + 1
	***REMOVED***
	if high >= 0 ***REMOVED***
		intervals = append(intervals, interval***REMOVED***low, high***REMOVED***)
	***REMOVED***
	sort.Sort(byDecreasingLength(intervals))

	fmt.Printf("const (\n")
	fmt.Printf("\tjis0208    = 1\n")
	fmt.Printf("\tjis0212    = 2\n")
	fmt.Printf("\tcodeMask   = 0x7f\n")
	fmt.Printf("\tcodeShift  = 7\n")
	fmt.Printf("\ttableShift = 14\n")
	fmt.Printf(")\n\n")

	fmt.Printf("const numEncodeTables = %d\n\n", len(intervals))
	fmt.Printf("// encodeX are the encoding tables from Unicode to JIS code,\n")
	fmt.Printf("// sorted by decreasing length.\n")
	for i, v := range intervals ***REMOVED***
		fmt.Printf("// encode%d: %5d entries for runes in [%5d, %5d).\n", i, v.len(), v.low, v.high)
	***REMOVED***
	fmt.Printf("//\n")
	fmt.Printf("// The high two bits of the value record whether the JIS code comes from the\n")
	fmt.Printf("// JIS0208 table (high bits == 1) or the JIS0212 table (high bits == 2).\n")
	fmt.Printf("// The low 14 bits are two 7-bit unsigned integers j1 and j2 that form the\n")
	fmt.Printf("// JIS code (94*j1 + j2) within that table.\n")
	fmt.Printf("\n")

	for i, v := range intervals ***REMOVED***
		fmt.Printf("const encode%dLow, encode%dHigh = %d, %d\n\n", i, i, v.low, v.high)
		fmt.Printf("var encode%d = [...]uint16***REMOVED***\n", i)
		for j := v.low; j < v.high; j++ ***REMOVED***
			x := reverse[j]
			if x.table == -1 ***REMOVED***
				continue
			***REMOVED***
			fmt.Printf("\t%d - %d: jis%s<<14 | 0x%02X<<7 | 0x%02X,\n",
				j, v.low, tables[x.table].name, x.jisCode/94, x.jisCode%94)
		***REMOVED***
		fmt.Printf("***REMOVED***\n\n")
	***REMOVED***
***REMOVED***

// interval is a half-open interval [low, high).
type interval struct ***REMOVED***
	low, high int
***REMOVED***

func (i interval) len() int ***REMOVED*** return i.high - i.low ***REMOVED***

// byDecreasingLength sorts intervals by decreasing length.
type byDecreasingLength []interval

func (b byDecreasingLength) Len() int           ***REMOVED*** return len(b) ***REMOVED***
func (b byDecreasingLength) Less(i, j int) bool ***REMOVED*** return b[i].len() > b[j].len() ***REMOVED***
func (b byDecreasingLength) Swap(i, j int)      ***REMOVED*** b[i], b[j] = b[j], b[i] ***REMOVED***