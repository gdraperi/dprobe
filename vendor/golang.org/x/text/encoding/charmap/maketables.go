// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/internal/gen"
)

const ascii = "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f" +
	"\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f" +
	` !"#$%&'()*+,-./0123456789:;<=>?` +
	`@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_` +
	"`abcdefghijklmnopqrstuvwxyz***REMOVED***|***REMOVED***~\u007f"

var encodings = []struct ***REMOVED***
	name        string
	mib         string
	comment     string
	varName     string
	replacement byte
	mapping     string
***REMOVED******REMOVED***
	***REMOVED***
		"IBM Code Page 037",
		"IBM037",
		"",
		"CodePage037",
		0x3f,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/glibc-IBM037-2.1.2.ucm",
	***REMOVED***,
	***REMOVED***
		"IBM Code Page 437",
		"PC8CodePage437",
		"",
		"CodePage437",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/glibc-IBM437-2.1.2.ucm",
	***REMOVED***,
	***REMOVED***
		"IBM Code Page 850",
		"PC850Multilingual",
		"",
		"CodePage850",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/glibc-IBM850-2.1.2.ucm",
	***REMOVED***,
	***REMOVED***
		"IBM Code Page 852",
		"PCp852",
		"",
		"CodePage852",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/glibc-IBM852-2.1.2.ucm",
	***REMOVED***,
	***REMOVED***
		"IBM Code Page 855",
		"IBM855",
		"",
		"CodePage855",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/glibc-IBM855-2.1.2.ucm",
	***REMOVED***,
	***REMOVED***
		"Windows Code Page 858", // PC latin1 with Euro
		"IBM00858",
		"",
		"CodePage858",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/windows-858-2000.ucm",
	***REMOVED***,
	***REMOVED***
		"IBM Code Page 860",
		"IBM860",
		"",
		"CodePage860",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/glibc-IBM860-2.1.2.ucm",
	***REMOVED***,
	***REMOVED***
		"IBM Code Page 862",
		"PC862LatinHebrew",
		"",
		"CodePage862",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/glibc-IBM862-2.1.2.ucm",
	***REMOVED***,
	***REMOVED***
		"IBM Code Page 863",
		"IBM863",
		"",
		"CodePage863",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/glibc-IBM863-2.1.2.ucm",
	***REMOVED***,
	***REMOVED***
		"IBM Code Page 865",
		"IBM865",
		"",
		"CodePage865",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/glibc-IBM865-2.1.2.ucm",
	***REMOVED***,
	***REMOVED***
		"IBM Code Page 866",
		"IBM866",
		"",
		"CodePage866",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-ibm866.txt",
	***REMOVED***,
	***REMOVED***
		"IBM Code Page 1047",
		"IBM1047",
		"",
		"CodePage1047",
		0x3f,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/glibc-IBM1047-2.1.2.ucm",
	***REMOVED***,
	***REMOVED***
		"IBM Code Page 1140",
		"IBM01140",
		"",
		"CodePage1140",
		0x3f,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/ibm-1140_P100-1997.ucm",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-1",
		"ISOLatin1",
		"",
		"ISO8859_1",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/iso-8859_1-1998.ucm",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-2",
		"ISOLatin2",
		"",
		"ISO8859_2",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-2.txt",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-3",
		"ISOLatin3",
		"",
		"ISO8859_3",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-3.txt",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-4",
		"ISOLatin4",
		"",
		"ISO8859_4",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-4.txt",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-5",
		"ISOLatinCyrillic",
		"",
		"ISO8859_5",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-5.txt",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-6",
		"ISOLatinArabic",
		"",
		"ISO8859_6,ISO8859_6E,ISO8859_6I",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-6.txt",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-7",
		"ISOLatinGreek",
		"",
		"ISO8859_7",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-7.txt",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-8",
		"ISOLatinHebrew",
		"",
		"ISO8859_8,ISO8859_8E,ISO8859_8I",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-8.txt",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-9",
		"ISOLatin5",
		"",
		"ISO8859_9",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/iso-8859_9-1999.ucm",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-10",
		"ISOLatin6",
		"",
		"ISO8859_10",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-10.txt",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-13",
		"ISO885913",
		"",
		"ISO8859_13",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-13.txt",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-14",
		"ISO885914",
		"",
		"ISO8859_14",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-14.txt",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-15",
		"ISO885915",
		"",
		"ISO8859_15",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-15.txt",
	***REMOVED***,
	***REMOVED***
		"ISO 8859-16",
		"ISO885916",
		"",
		"ISO8859_16",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-iso-8859-16.txt",
	***REMOVED***,
	***REMOVED***
		"KOI8-R",
		"KOI8R",
		"",
		"KOI8R",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-koi8-r.txt",
	***REMOVED***,
	***REMOVED***
		"KOI8-U",
		"KOI8U",
		"",
		"KOI8U",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-koi8-u.txt",
	***REMOVED***,
	***REMOVED***
		"Macintosh",
		"Macintosh",
		"",
		"Macintosh",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-macintosh.txt",
	***REMOVED***,
	***REMOVED***
		"Macintosh Cyrillic",
		"MacintoshCyrillic",
		"",
		"MacintoshCyrillic",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-x-mac-cyrillic.txt",
	***REMOVED***,
	***REMOVED***
		"Windows 874",
		"Windows874",
		"",
		"Windows874",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-windows-874.txt",
	***REMOVED***,
	***REMOVED***
		"Windows 1250",
		"Windows1250",
		"",
		"Windows1250",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-windows-1250.txt",
	***REMOVED***,
	***REMOVED***
		"Windows 1251",
		"Windows1251",
		"",
		"Windows1251",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-windows-1251.txt",
	***REMOVED***,
	***REMOVED***
		"Windows 1252",
		"Windows1252",
		"",
		"Windows1252",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-windows-1252.txt",
	***REMOVED***,
	***REMOVED***
		"Windows 1253",
		"Windows1253",
		"",
		"Windows1253",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-windows-1253.txt",
	***REMOVED***,
	***REMOVED***
		"Windows 1254",
		"Windows1254",
		"",
		"Windows1254",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-windows-1254.txt",
	***REMOVED***,
	***REMOVED***
		"Windows 1255",
		"Windows1255",
		"",
		"Windows1255",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-windows-1255.txt",
	***REMOVED***,
	***REMOVED***
		"Windows 1256",
		"Windows1256",
		"",
		"Windows1256",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-windows-1256.txt",
	***REMOVED***,
	***REMOVED***
		"Windows 1257",
		"Windows1257",
		"",
		"Windows1257",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-windows-1257.txt",
	***REMOVED***,
	***REMOVED***
		"Windows 1258",
		"Windows1258",
		"",
		"Windows1258",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-windows-1258.txt",
	***REMOVED***,
	***REMOVED***
		"X-User-Defined",
		"XUserDefined",
		"It is defined at http://encoding.spec.whatwg.org/#x-user-defined",
		"XUserDefined",
		encoding.ASCIISub,
		ascii +
			"\uf780\uf781\uf782\uf783\uf784\uf785\uf786\uf787" +
			"\uf788\uf789\uf78a\uf78b\uf78c\uf78d\uf78e\uf78f" +
			"\uf790\uf791\uf792\uf793\uf794\uf795\uf796\uf797" +
			"\uf798\uf799\uf79a\uf79b\uf79c\uf79d\uf79e\uf79f" +
			"\uf7a0\uf7a1\uf7a2\uf7a3\uf7a4\uf7a5\uf7a6\uf7a7" +
			"\uf7a8\uf7a9\uf7aa\uf7ab\uf7ac\uf7ad\uf7ae\uf7af" +
			"\uf7b0\uf7b1\uf7b2\uf7b3\uf7b4\uf7b5\uf7b6\uf7b7" +
			"\uf7b8\uf7b9\uf7ba\uf7bb\uf7bc\uf7bd\uf7be\uf7bf" +
			"\uf7c0\uf7c1\uf7c2\uf7c3\uf7c4\uf7c5\uf7c6\uf7c7" +
			"\uf7c8\uf7c9\uf7ca\uf7cb\uf7cc\uf7cd\uf7ce\uf7cf" +
			"\uf7d0\uf7d1\uf7d2\uf7d3\uf7d4\uf7d5\uf7d6\uf7d7" +
			"\uf7d8\uf7d9\uf7da\uf7db\uf7dc\uf7dd\uf7de\uf7df" +
			"\uf7e0\uf7e1\uf7e2\uf7e3\uf7e4\uf7e5\uf7e6\uf7e7" +
			"\uf7e8\uf7e9\uf7ea\uf7eb\uf7ec\uf7ed\uf7ee\uf7ef" +
			"\uf7f0\uf7f1\uf7f2\uf7f3\uf7f4\uf7f5\uf7f6\uf7f7" +
			"\uf7f8\uf7f9\uf7fa\uf7fb\uf7fc\uf7fd\uf7fe\uf7ff",
	***REMOVED***,
***REMOVED***

func getWHATWG(url string) string ***REMOVED***
	res, err := http.Get(url)
	if err != nil ***REMOVED***
		log.Fatalf("%q: Get: %v", url, err)
	***REMOVED***
	defer res.Body.Close()

	mapping := make([]rune, 128)
	for i := range mapping ***REMOVED***
		mapping[i] = '\ufffd'
	***REMOVED***

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() ***REMOVED***
		s := strings.TrimSpace(scanner.Text())
		if s == "" || s[0] == '#' ***REMOVED***
			continue
		***REMOVED***
		x, y := 0, 0
		if _, err := fmt.Sscanf(s, "%d\t0x%x", &x, &y); err != nil ***REMOVED***
			log.Fatalf("could not parse %q", s)
		***REMOVED***
		if x < 0 || 128 <= x ***REMOVED***
			log.Fatalf("code %d is out of range", x)
		***REMOVED***
		if 0x80 <= y && y < 0xa0 ***REMOVED***
			// We diverge from the WHATWG spec by mapping control characters
			// in the range [0x80, 0xa0) to U+FFFD.
			continue
		***REMOVED***
		mapping[x] = rune(y)
	***REMOVED***
	return ascii + string(mapping)
***REMOVED***

func getUCM(url string) string ***REMOVED***
	res, err := http.Get(url)
	if err != nil ***REMOVED***
		log.Fatalf("%q: Get: %v", url, err)
	***REMOVED***
	defer res.Body.Close()

	mapping := make([]rune, 256)
	for i := range mapping ***REMOVED***
		mapping[i] = '\ufffd'
	***REMOVED***

	charsFound := 0
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() ***REMOVED***
		s := strings.TrimSpace(scanner.Text())
		if s == "" || s[0] == '#' ***REMOVED***
			continue
		***REMOVED***
		var c byte
		var r rune
		if _, err := fmt.Sscanf(s, `<U%x> \x%x |0`, &r, &c); err != nil ***REMOVED***
			continue
		***REMOVED***
		mapping[c] = r
		charsFound++
	***REMOVED***

	if charsFound < 200 ***REMOVED***
		log.Fatalf("%q: only %d characters found (wrong page format?)", url, charsFound)
	***REMOVED***

	return string(mapping)
***REMOVED***

func main() ***REMOVED***
	mibs := map[string]bool***REMOVED******REMOVED***
	all := []string***REMOVED******REMOVED***

	w := gen.NewCodeWriter()
	defer w.WriteGoFile("tables.go", "charmap")

	printf := func(s string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** fmt.Fprintf(w, s, a...) ***REMOVED***

	printf("import (\n")
	printf("\t\"golang.org/x/text/encoding\"\n")
	printf("\t\"golang.org/x/text/encoding/internal/identifier\"\n")
	printf(")\n\n")
	for _, e := range encodings ***REMOVED***
		varNames := strings.Split(e.varName, ",")
		all = append(all, varNames...)
		varName := varNames[0]
		switch ***REMOVED***
		case strings.HasPrefix(e.mapping, "http://encoding.spec.whatwg.org/"):
			e.mapping = getWHATWG(e.mapping)
		case strings.HasPrefix(e.mapping, "http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/"):
			e.mapping = getUCM(e.mapping)
		***REMOVED***

		asciiSuperset, low := strings.HasPrefix(e.mapping, ascii), 0x00
		if asciiSuperset ***REMOVED***
			low = 0x80
		***REMOVED***
		lvn := 1
		if strings.HasPrefix(varName, "ISO") || strings.HasPrefix(varName, "KOI") ***REMOVED***
			lvn = 3
		***REMOVED***
		lowerVarName := strings.ToLower(varName[:lvn]) + varName[lvn:]
		printf("// %s is the %s encoding.\n", varName, e.name)
		if e.comment != "" ***REMOVED***
			printf("//\n// %s\n", e.comment)
		***REMOVED***
		printf("var %s *Charmap = &%s\n\nvar %s = Charmap***REMOVED***\nname: %q,\n",
			varName, lowerVarName, lowerVarName, e.name)
		if mibs[e.mib] ***REMOVED***
			log.Fatalf("MIB type %q declared multiple times.", e.mib)
		***REMOVED***
		printf("mib: identifier.%s,\n", e.mib)
		printf("asciiSuperset: %t,\n", asciiSuperset)
		printf("low: 0x%02x,\n", low)
		printf("replacement: 0x%02x,\n", e.replacement)

		printf("decode: [256]utf8Enc***REMOVED***\n")
		i, backMapping := 0, map[rune]byte***REMOVED******REMOVED***
		for _, c := range e.mapping ***REMOVED***
			if _, ok := backMapping[c]; !ok && c != utf8.RuneError ***REMOVED***
				backMapping[c] = byte(i)
			***REMOVED***
			var buf [8]byte
			n := utf8.EncodeRune(buf[:], c)
			if n > 3 ***REMOVED***
				panic(fmt.Sprintf("rune %q (%U) is too long", c, c))
			***REMOVED***
			printf("***REMOVED***%d,[3]byte***REMOVED***0x%02x,0x%02x,0x%02x***REMOVED******REMOVED***,", n, buf[0], buf[1], buf[2])
			if i%2 == 1 ***REMOVED***
				printf("\n")
			***REMOVED***
			i++
		***REMOVED***
		printf("***REMOVED***,\n")

		printf("encode: [256]uint32***REMOVED***\n")
		encode := make([]uint32, 0, 256)
		for c, i := range backMapping ***REMOVED***
			encode = append(encode, uint32(i)<<24|uint32(c))
		***REMOVED***
		sort.Sort(byRune(encode))
		for len(encode) < cap(encode) ***REMOVED***
			encode = append(encode, encode[len(encode)-1])
		***REMOVED***
		for i, enc := range encode ***REMOVED***
			printf("0x%08x,", enc)
			if i%8 == 7 ***REMOVED***
				printf("\n")
			***REMOVED***
		***REMOVED***
		printf("***REMOVED***,\n***REMOVED***\n")

		// Add an estimate of the size of a single Charmap***REMOVED******REMOVED*** struct value, which
		// includes two 256 elem arrays of 4 bytes and some extra fields, which
		// align to 3 uint64s on 64-bit architectures.
		w.Size += 2*4*256 + 3*8
	***REMOVED***
	// TODO: add proper line breaking.
	printf("var listAll = []encoding.Encoding***REMOVED***\n%s,\n***REMOVED***\n\n", strings.Join(all, ",\n"))
***REMOVED***

type byRune []uint32

func (b byRune) Len() int           ***REMOVED*** return len(b) ***REMOVED***
func (b byRune) Less(i, j int) bool ***REMOVED*** return b[i]&0xffffff < b[j]&0xffffff ***REMOVED***
func (b byRune) Swap(i, j int)      ***REMOVED*** b[i], b[j] = b[j], b[i] ***REMOVED***
