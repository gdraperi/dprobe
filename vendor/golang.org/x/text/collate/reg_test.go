// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collate

import (
	"archive/zip"
	"bufio"
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"golang.org/x/text/collate/build"
	"golang.org/x/text/internal/gen"
	"golang.org/x/text/language"
)

var long = flag.Bool("long", false,
	"run time-consuming tests, such as tests that fetch data online")

// This regression test runs tests for the test files in CollationTest.zip
// (taken from http://www.unicode.org/Public/UCA/<gen.UnicodeVersion()>/).
//
// The test files have the following form:
// # header
// 0009 0021;	# ('\u0009') <CHARACTER TABULATION>	[| | | 0201 025E]
// 0009 003F;	# ('\u0009') <CHARACTER TABULATION>	[| | | 0201 0263]
// 000A 0021;	# ('\u000A') <LINE FEED (LF)>	[| | | 0202 025E]
// 000A 003F;	# ('\u000A') <LINE FEED (LF)>	[| | | 0202 0263]
//
// The part before the semicolon is the hex representation of a sequence
// of runes. After the hash mark is a comment. The strings
// represented by rune sequence are in the file in sorted order, as
// defined by the DUCET.

type Test struct ***REMOVED***
	name    string
	str     [][]byte
	comment []string
***REMOVED***

var versionRe = regexp.MustCompile(`# UCA Version: (.*)\n?$`)
var testRe = regexp.MustCompile(`^([\dA-F ]+);.*# (.*)\n?$`)

func TestCollation(t *testing.T) ***REMOVED***
	if !gen.IsLocal() && !*long ***REMOVED***
		t.Skip("skipping test to prevent downloading; to run use -long or use -local to specify a local source")
	***REMOVED***
	t.Skip("must first update to new file format to support test")
	for _, test := range loadTestData() ***REMOVED***
		doTest(t, test)
	***REMOVED***
***REMOVED***

func Error(e error) ***REMOVED***
	if e != nil ***REMOVED***
		log.Fatal(e)
	***REMOVED***
***REMOVED***

// parseUCA parses a Default Unicode Collation Element Table of the format
// specified in http://www.unicode.org/reports/tr10/#File_Format.
// It returns the variable top.
func parseUCA(builder *build.Builder) ***REMOVED***
	r := gen.OpenUnicodeFile("UCA", "", "allkeys.txt")
	defer r.Close()
	input := bufio.NewReader(r)
	colelem := regexp.MustCompile(`\[([.*])([0-9A-F.]+)\]`)
	for i := 1; true; i++ ***REMOVED***
		l, prefix, err := input.ReadLine()
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		Error(err)
		line := string(l)
		if prefix ***REMOVED***
			log.Fatalf("%d: buffer overflow", i)
		***REMOVED***
		if len(line) == 0 || line[0] == '#' ***REMOVED***
			continue
		***REMOVED***
		if line[0] == '@' ***REMOVED***
			if strings.HasPrefix(line[1:], "version ") ***REMOVED***
				if v := strings.Split(line[1:], " ")[1]; v != gen.UnicodeVersion() ***REMOVED***
					log.Fatalf("incompatible version %s; want %s", v, gen.UnicodeVersion())
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// parse entries
			part := strings.Split(line, " ; ")
			if len(part) != 2 ***REMOVED***
				log.Fatalf("%d: production rule without ';': %v", i, line)
			***REMOVED***
			lhs := []rune***REMOVED******REMOVED***
			for _, v := range strings.Split(part[0], " ") ***REMOVED***
				if v != "" ***REMOVED***
					lhs = append(lhs, rune(convHex(i, v)))
				***REMOVED***
			***REMOVED***
			vars := []int***REMOVED******REMOVED***
			rhs := [][]int***REMOVED******REMOVED***
			for i, m := range colelem.FindAllStringSubmatch(part[1], -1) ***REMOVED***
				if m[1] == "*" ***REMOVED***
					vars = append(vars, i)
				***REMOVED***
				elem := []int***REMOVED******REMOVED***
				for _, h := range strings.Split(m[2], ".") ***REMOVED***
					elem = append(elem, convHex(i, h))
				***REMOVED***
				rhs = append(rhs, elem)
			***REMOVED***
			builder.Add(lhs, rhs, vars)
		***REMOVED***
	***REMOVED***
***REMOVED***

func convHex(line int, s string) int ***REMOVED***
	r, e := strconv.ParseInt(s, 16, 32)
	if e != nil ***REMOVED***
		log.Fatalf("%d: %v", line, e)
	***REMOVED***
	return int(r)
***REMOVED***

func loadTestData() []Test ***REMOVED***
	f := gen.OpenUnicodeFile("UCA", "", "CollationTest.zip")
	buffer, err := ioutil.ReadAll(f)
	f.Close()
	Error(err)
	archive, err := zip.NewReader(bytes.NewReader(buffer), int64(len(buffer)))
	Error(err)
	tests := []Test***REMOVED******REMOVED***
	for _, f := range archive.File ***REMOVED***
		// Skip the short versions, which are simply duplicates of the long versions.
		if strings.Contains(f.Name, "SHORT") || f.FileInfo().IsDir() ***REMOVED***
			continue
		***REMOVED***
		ff, err := f.Open()
		Error(err)
		defer ff.Close()
		scanner := bufio.NewScanner(ff)
		test := Test***REMOVED***name: path.Base(f.Name)***REMOVED***
		for scanner.Scan() ***REMOVED***
			line := scanner.Text()
			if len(line) <= 1 || line[0] == '#' ***REMOVED***
				if m := versionRe.FindStringSubmatch(line); m != nil ***REMOVED***
					if m[1] != gen.UnicodeVersion() ***REMOVED***
						log.Printf("warning:%s: version is %s; want %s", f.Name, m[1], gen.UnicodeVersion())
					***REMOVED***
				***REMOVED***
				continue
			***REMOVED***
			m := testRe.FindStringSubmatch(line)
			if m == nil || len(m) < 3 ***REMOVED***
				log.Fatalf(`Failed to parse: "%s" result: %#v`, line, m)
			***REMOVED***
			str := []byte***REMOVED******REMOVED***
			// In the regression test data (unpaired) surrogates are assigned a weight
			// corresponding to their code point value.  However, utf8.DecodeRune,
			// which is used to compute the implicit weight, assigns FFFD to surrogates.
			// We therefore skip tests with surrogates.  This skips about 35 entries
			// per test.
			valid := true
			for _, split := range strings.Split(m[1], " ") ***REMOVED***
				r, err := strconv.ParseUint(split, 16, 64)
				Error(err)
				valid = valid && utf8.ValidRune(rune(r))
				str = append(str, string(rune(r))...)
			***REMOVED***
			if valid ***REMOVED***
				test.str = append(test.str, str)
				test.comment = append(test.comment, m[2])
			***REMOVED***
		***REMOVED***
		if scanner.Err() != nil ***REMOVED***
			log.Fatal(scanner.Err())
		***REMOVED***
		tests = append(tests, test)
	***REMOVED***
	return tests
***REMOVED***

var errorCount int

func runes(b []byte) []rune ***REMOVED***
	return []rune(string(b))
***REMOVED***

var shifted = language.MustParse("und-u-ka-shifted-ks-level4")

func doTest(t *testing.T, tc Test) ***REMOVED***
	bld := build.NewBuilder()
	parseUCA(bld)
	w, err := bld.Build()
	Error(err)
	var tag language.Tag
	if !strings.Contains(tc.name, "NON_IGNOR") ***REMOVED***
		tag = shifted
	***REMOVED***
	c := NewFromTable(w, OptionsFromTag(tag))
	b := &Buffer***REMOVED******REMOVED***
	prev := tc.str[0]
	for i := 1; i < len(tc.str); i++ ***REMOVED***
		b.Reset()
		s := tc.str[i]
		ka := c.Key(b, prev)
		kb := c.Key(b, s)
		if r := bytes.Compare(ka, kb); r == 1 ***REMOVED***
			t.Errorf("%s:%d: Key(%.4X) < Key(%.4X) (%X < %X) == %d; want -1 or 0", tc.name, i, []rune(string(prev)), []rune(string(s)), ka, kb, r)
			prev = s
			continue
		***REMOVED***
		if r := c.Compare(prev, s); r == 1 ***REMOVED***
			t.Errorf("%s:%d: Compare(%.4X, %.4X) == %d; want -1 or 0", tc.name, i, runes(prev), runes(s), r)
		***REMOVED***
		if r := c.Compare(s, prev); r == -1 ***REMOVED***
			t.Errorf("%s:%d: Compare(%.4X, %.4X) == %d; want 1 or 0", tc.name, i, runes(s), runes(prev), r)
		***REMOVED***
		prev = s
	***REMOVED***
***REMOVED***
