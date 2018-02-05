// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// This program generates the trie for width operations. The generated table
// includes width category information as well as the normalization mappings.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"unicode/utf8"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/triegen"
)

// See gen_common.go for flags.

func main() ***REMOVED***
	gen.Init()
	genTables()
	genTests()
	gen.Repackage("gen_trieval.go", "trieval.go", "width")
	gen.Repackage("gen_common.go", "common_test.go", "width")
***REMOVED***

func genTables() ***REMOVED***
	t := triegen.NewTrie("width")
	// fold and inverse mappings. See mapComment for a description of the format
	// of each entry. Add dummy value to make an index of 0 mean no mapping.
	inverse := [][4]byte***REMOVED******REMOVED******REMOVED******REMOVED***
	mapping := map[[4]byte]int***REMOVED***[4]byte***REMOVED******REMOVED***: 0***REMOVED***

	getWidthData(func(r rune, tag elem, alt rune) ***REMOVED***
		idx := 0
		if alt != 0 ***REMOVED***
			var buf [4]byte
			buf[0] = byte(utf8.EncodeRune(buf[1:], alt))
			s := string(r)
			buf[buf[0]] ^= s[len(s)-1]
			var ok bool
			if idx, ok = mapping[buf]; !ok ***REMOVED***
				idx = len(mapping)
				if idx > math.MaxUint8 ***REMOVED***
					log.Fatalf("Index %d does not fit in a byte.", idx)
				***REMOVED***
				mapping[buf] = idx
				inverse = append(inverse, buf)
			***REMOVED***
		***REMOVED***
		t.Insert(r, uint64(tag|elem(idx)))
	***REMOVED***)

	w := &bytes.Buffer***REMOVED******REMOVED***
	gen.WriteUnicodeVersion(w)

	sz, err := t.Gen(w)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	sz += writeMappings(w, inverse)

	fmt.Fprintf(w, "// Total table size %d bytes (%dKiB)\n", sz, sz/1024)

	gen.WriteVersionedGoFile(*outputFile, "width", w.Bytes())
***REMOVED***

const inverseDataComment = `
// inverseData contains 4-byte entries of the following format:
//   <length> <modified UTF-8-encoded rune> <0 padding>
// The last byte of the UTF-8-encoded rune is xor-ed with the last byte of the
// UTF-8 encoding of the original rune. Mappings often have the following
// pattern:
//   Ａ -> A  (U+FF21 -> U+0041)
//   Ｂ -> B  (U+FF22 -> U+0042)
//   ...
// By xor-ing the last byte the same entry can be shared by many mappings. This
// reduces the total number of distinct entries by about two thirds.
// The resulting entry for the aforementioned mappings is
//   ***REMOVED*** 0x01, 0xE0, 0x00, 0x00 ***REMOVED***
// Using this entry to map U+FF21 (UTF-8 [EF BC A1]), we get
//   E0 ^ A1 = 41.
// Similarly, for U+FF22 (UTF-8 [EF BC A2]), we get
//   E0 ^ A2 = 42.
// Note that because of the xor-ing, the byte sequence stored in the entry is
// not valid UTF-8.`

func writeMappings(w io.Writer, data [][4]byte) int ***REMOVED***
	fmt.Fprintln(w, inverseDataComment)
	fmt.Fprintf(w, "var inverseData = [%d][4]byte***REMOVED***\n", len(data))
	for _, x := range data ***REMOVED***
		fmt.Fprintf(w, "***REMOVED*** 0x%02x, 0x%02x, 0x%02x, 0x%02x ***REMOVED***,\n", x[0], x[1], x[2], x[3])
	***REMOVED***
	fmt.Fprintln(w, "***REMOVED***")
	return len(data) * 4
***REMOVED***

func genTests() ***REMOVED***
	w := &bytes.Buffer***REMOVED******REMOVED***
	fmt.Fprintf(w, "\nvar mapRunes = map[rune]struct***REMOVED***r rune; e elem***REMOVED******REMOVED***\n")
	getWidthData(func(r rune, tag elem, alt rune) ***REMOVED***
		if alt != 0 ***REMOVED***
			fmt.Fprintf(w, "\t0x%X: ***REMOVED***0x%X, 0x%X***REMOVED***,\n", r, alt, tag)
		***REMOVED***
	***REMOVED***)
	fmt.Fprintln(w, "***REMOVED***")
	gen.WriteGoFile("runes_test.go", "width", w.Bytes())
***REMOVED***
