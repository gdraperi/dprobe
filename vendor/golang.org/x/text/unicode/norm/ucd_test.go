// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/testtext"
)

var once sync.Once

func skipShort(t *testing.T) ***REMOVED***
	testtext.SkipIfNotLong(t)

	once.Do(func() ***REMOVED*** loadTestData(t) ***REMOVED***)
***REMOVED***

// This regression test runs the test set in NormalizationTest.txt
// (taken from http://www.unicode.org/Public/<unicode.Version>/ucd/).
//
// NormalizationTest.txt has form:
// @Part0 # Specific cases
// #
// 1E0A;1E0A;0044 0307;1E0A;0044 0307; # (Ḋ; Ḋ; D◌̇; Ḋ; D◌̇; ) LATIN CAPITAL LETTER D WITH DOT ABOVE
// 1E0C;1E0C;0044 0323;1E0C;0044 0323; # (Ḍ; Ḍ; D◌̣; Ḍ; D◌̣; ) LATIN CAPITAL LETTER D WITH DOT BELOW
//
// Each test has 5 columns (c1, c2, c3, c4, c5), where
// (c1, c2, c3, c4, c5) == (c1, NFC(c1), NFD(c1), NFKC(c1), NFKD(c1))
//
// CONFORMANCE:
// 1. The following invariants must be true for all conformant implementations
//
//    NFC
//      c2 ==  NFC(c1) ==  NFC(c2) ==  NFC(c3)
//      c4 ==  NFC(c4) ==  NFC(c5)
//
//    NFD
//      c3 ==  NFD(c1) ==  NFD(c2) ==  NFD(c3)
//      c5 ==  NFD(c4) ==  NFD(c5)
//
//    NFKC
//      c4 == NFKC(c1) == NFKC(c2) == NFKC(c3) == NFKC(c4) == NFKC(c5)
//
//    NFKD
//      c5 == NFKD(c1) == NFKD(c2) == NFKD(c3) == NFKD(c4) == NFKD(c5)
//
// 2. For every code point X assigned in this version of Unicode that is not
//    specifically listed in Part 1, the following invariants must be true
//    for all conformant implementations:
//
//      X == NFC(X) == NFD(X) == NFKC(X) == NFKD(X)
//

// Column types.
const (
	cRaw = iota
	cNFC
	cNFD
	cNFKC
	cNFKD
	cMaxColumns
)

// Holds data from NormalizationTest.txt
var part []Part

type Part struct ***REMOVED***
	name   string
	number int
	tests  []Test
***REMOVED***

type Test struct ***REMOVED***
	name   string
	partnr int
	number int
	r      rune                // used for character by character test
	cols   [cMaxColumns]string // Each has 5 entries, see below.
***REMOVED***

func (t Test) Name() string ***REMOVED***
	if t.number < 0 ***REMOVED***
		return part[t.partnr].name
	***REMOVED***
	return fmt.Sprintf("%s:%d", part[t.partnr].name, t.number)
***REMOVED***

var partRe = regexp.MustCompile(`@Part(\d) # (.*)$`)
var testRe = regexp.MustCompile(`^` + strings.Repeat(`([\dA-F ]+);`, 5) + ` # (.*)$`)

var counter int

// Load the data form NormalizationTest.txt
func loadTestData(t *testing.T) ***REMOVED***
	f := gen.OpenUCDFile("NormalizationTest.txt")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() ***REMOVED***
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' ***REMOVED***
			continue
		***REMOVED***
		m := partRe.FindStringSubmatch(line)
		if m != nil ***REMOVED***
			if len(m) < 3 ***REMOVED***
				t.Fatal("Failed to parse Part: ", line)
			***REMOVED***
			i, err := strconv.Atoi(m[1])
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			name := m[2]
			part = append(part, Part***REMOVED***name: name[:len(name)-1], number: i***REMOVED***)
			continue
		***REMOVED***
		m = testRe.FindStringSubmatch(line)
		if m == nil || len(m) < 7 ***REMOVED***
			t.Fatalf(`Failed to parse: "%s" result: %#v`, line, m)
		***REMOVED***
		test := Test***REMOVED***name: m[6], partnr: len(part) - 1, number: counter***REMOVED***
		counter++
		for j := 1; j < len(m)-1; j++ ***REMOVED***
			for _, split := range strings.Split(m[j], " ") ***REMOVED***
				r, err := strconv.ParseUint(split, 16, 64)
				if err != nil ***REMOVED***
					t.Fatal(err)
				***REMOVED***
				if test.r == 0 ***REMOVED***
					// save for CharacterByCharacterTests
					test.r = rune(r)
				***REMOVED***
				var buf [utf8.UTFMax]byte
				sz := utf8.EncodeRune(buf[:], rune(r))
				test.cols[j-1] += string(buf[:sz])
			***REMOVED***
		***REMOVED***
		part := &part[len(part)-1]
		part.tests = append(part.tests, test)
	***REMOVED***
	if scanner.Err() != nil ***REMOVED***
		t.Fatal(scanner.Err())
	***REMOVED***
***REMOVED***

func cmpResult(t *testing.T, tc *Test, name string, f Form, gold, test, result string) ***REMOVED***
	if gold != result ***REMOVED***
		t.Errorf("%s:%s: %s(%+q)=%+q; want %+q: %s",
			tc.Name(), name, fstr[f], test, result, gold, tc.name)
	***REMOVED***
***REMOVED***

func cmpIsNormal(t *testing.T, tc *Test, name string, f Form, test string, result, want bool) ***REMOVED***
	if result != want ***REMOVED***
		t.Errorf("%s:%s: %s(%+q)=%v; want %v", tc.Name(), name, fstr[f], test, result, want)
	***REMOVED***
***REMOVED***

func doTest(t *testing.T, tc *Test, f Form, gold, test string) ***REMOVED***
	testb := []byte(test)
	result := f.Bytes(testb)
	cmpResult(t, tc, "Bytes", f, gold, test, string(result))

	sresult := f.String(test)
	cmpResult(t, tc, "String", f, gold, test, sresult)

	acc := []byte***REMOVED******REMOVED***
	i := Iter***REMOVED******REMOVED***
	i.InitString(f, test)
	for !i.Done() ***REMOVED***
		acc = append(acc, i.Next()...)
	***REMOVED***
	cmpResult(t, tc, "Iter.Next", f, gold, test, string(acc))

	buf := make([]byte, 128)
	acc = nil
	for p := 0; p < len(testb); ***REMOVED***
		nDst, nSrc, _ := f.Transform(buf, testb[p:], true)
		acc = append(acc, buf[:nDst]...)
		p += nSrc
	***REMOVED***
	cmpResult(t, tc, "Transform", f, gold, test, string(acc))

	for i := range test ***REMOVED***
		out := f.Append(f.Bytes([]byte(test[:i])), []byte(test[i:])...)
		cmpResult(t, tc, fmt.Sprintf(":Append:%d", i), f, gold, test, string(out))
	***REMOVED***
	cmpIsNormal(t, tc, "IsNormal", f, test, f.IsNormal([]byte(test)), test == gold)
	cmpIsNormal(t, tc, "IsNormalString", f, test, f.IsNormalString(test), test == gold)
***REMOVED***

func doConformanceTests(t *testing.T, tc *Test, partn int) ***REMOVED***
	for i := 0; i <= 2; i++ ***REMOVED***
		doTest(t, tc, NFC, tc.cols[1], tc.cols[i])
		doTest(t, tc, NFD, tc.cols[2], tc.cols[i])
		doTest(t, tc, NFKC, tc.cols[3], tc.cols[i])
		doTest(t, tc, NFKD, tc.cols[4], tc.cols[i])
	***REMOVED***
	for i := 3; i <= 4; i++ ***REMOVED***
		doTest(t, tc, NFC, tc.cols[3], tc.cols[i])
		doTest(t, tc, NFD, tc.cols[4], tc.cols[i])
		doTest(t, tc, NFKC, tc.cols[3], tc.cols[i])
		doTest(t, tc, NFKD, tc.cols[4], tc.cols[i])
	***REMOVED***
***REMOVED***

func TestCharacterByCharacter(t *testing.T) ***REMOVED***
	skipShort(t)
	tests := part[1].tests
	var last rune = 0
	for i := 0; i <= len(tests); i++ ***REMOVED*** // last one is special case
		var r rune
		if i == len(tests) ***REMOVED***
			r = 0x2FA1E // Don't have to go to 0x10FFFF
		***REMOVED*** else ***REMOVED***
			r = tests[i].r
		***REMOVED***
		for last++; last < r; last++ ***REMOVED***
			// Check all characters that were not explicitly listed in the test.
			tc := &Test***REMOVED***partnr: 1, number: -1***REMOVED***
			char := string(last)
			doTest(t, tc, NFC, char, char)
			doTest(t, tc, NFD, char, char)
			doTest(t, tc, NFKC, char, char)
			doTest(t, tc, NFKD, char, char)
		***REMOVED***
		if i < len(tests) ***REMOVED***
			doConformanceTests(t, &tests[i], 1)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestStandardTests(t *testing.T) ***REMOVED***
	skipShort(t)
	for _, j := range []int***REMOVED***0, 2, 3***REMOVED*** ***REMOVED***
		for _, test := range part[j].tests ***REMOVED***
			doConformanceTests(t, &test, j)
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestPerformance verifies that normalization is O(n). If any of the
// code does not properly check for maxCombiningChars, normalization
// may exhibit O(n**2) behavior.
func TestPerformance(t *testing.T) ***REMOVED***
	skipShort(t)
	runtime.GOMAXPROCS(2)
	success := make(chan bool, 1)
	go func() ***REMOVED***
		buf := bytes.Repeat([]byte("\u035D"), 1024*1024)
		buf = append(buf, "\u035B"...)
		NFC.Append(nil, buf...)
		success <- true
	***REMOVED***()
	timeout := time.After(1 * time.Second)
	select ***REMOVED***
	case <-success:
		// test completed before the timeout
	case <-timeout:
		t.Errorf(`unexpectedly long time to complete PerformanceTest`)
	***REMOVED***
***REMOVED***
