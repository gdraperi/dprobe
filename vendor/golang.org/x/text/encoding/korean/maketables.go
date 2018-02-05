// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

// This program generates tables.go:
//	go run maketables.go | gofmt > tables.go

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

func main() ***REMOVED***
	fmt.Printf("// generated by go run maketables.go; DO NOT EDIT\n\n")
	fmt.Printf("// Package korean provides Korean encodings such as EUC-KR.\n")
	fmt.Printf(`package korean // import "golang.org/x/text/encoding/korean"` + "\n\n")

	res, err := http.Get("http://encoding.spec.whatwg.org/index-euc-kr.txt")
	if err != nil ***REMOVED***
		log.Fatalf("Get: %v", err)
	***REMOVED***
	defer res.Body.Close()

	mapping := [65536]uint16***REMOVED******REMOVED***
	reverse := [65536]uint16***REMOVED******REMOVED***

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() ***REMOVED***
		s := strings.TrimSpace(scanner.Text())
		if s == "" || s[0] == '#' ***REMOVED***
			continue
		***REMOVED***
		x, y := uint16(0), uint16(0)
		if _, err := fmt.Sscanf(s, "%d 0x%x", &x, &y); err != nil ***REMOVED***
			log.Fatalf("could not parse %q", s)
		***REMOVED***
		if x < 0 || 178*(0xc7-0x81)+(0xfe-0xc7)*94+(0xff-0xa1) <= x ***REMOVED***
			log.Fatalf("EUC-KR code %d is out of range", x)
		***REMOVED***
		mapping[x] = y
		if reverse[y] == 0 ***REMOVED***
			c0, c1 := uint16(0), uint16(0)
			if x < 178*(0xc7-0x81) ***REMOVED***
				c0 = uint16(x/178) + 0x81
				c1 = uint16(x % 178)
				switch ***REMOVED***
				case c1 < 1*26:
					c1 += 0x41
				case c1 < 2*26:
					c1 += 0x47
				default:
					c1 += 0x4d
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				x -= 178 * (0xc7 - 0x81)
				c0 = uint16(x/94) + 0xc7
				c1 = uint16(x%94) + 0xa1
			***REMOVED***
			reverse[y] = c0<<8 | c1
		***REMOVED***
	***REMOVED***
	if err := scanner.Err(); err != nil ***REMOVED***
		log.Fatalf("scanner error: %v", err)
	***REMOVED***

	fmt.Printf("// decode is the decoding table from EUC-KR code to Unicode.\n")
	fmt.Printf("// It is defined at http://encoding.spec.whatwg.org/index-euc-kr.txt\n")
	fmt.Printf("var decode = [...]uint16***REMOVED***\n")
	for i, v := range mapping ***REMOVED***
		if v != 0 ***REMOVED***
			fmt.Printf("\t%d: 0x%04X,\n", i, v)
		***REMOVED***
	***REMOVED***
	fmt.Printf("***REMOVED***\n\n")

	// Any run of at least separation continuous zero entries in the reverse map will
	// be a separate encode table.
	const separation = 1024

	intervals := []interval(nil)
	low, high := -1, -1
	for i, v := range reverse ***REMOVED***
		if v == 0 ***REMOVED***
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

	fmt.Printf("const numEncodeTables = %d\n\n", len(intervals))
	fmt.Printf("// encodeX are the encoding tables from Unicode to EUC-KR code,\n")
	fmt.Printf("// sorted by decreasing length.\n")
	for i, v := range intervals ***REMOVED***
		fmt.Printf("// encode%d: %5d entries for runes in [%5d, %5d).\n", i, v.len(), v.low, v.high)
	***REMOVED***
	fmt.Printf("\n")

	for i, v := range intervals ***REMOVED***
		fmt.Printf("const encode%dLow, encode%dHigh = %d, %d\n\n", i, i, v.low, v.high)
		fmt.Printf("var encode%d = [...]uint16***REMOVED***\n", i)
		for j := v.low; j < v.high; j++ ***REMOVED***
			x := reverse[j]
			if x == 0 ***REMOVED***
				continue
			***REMOVED***
			fmt.Printf("\t%d-%d: 0x%04X,\n", j, v.low, x)
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
