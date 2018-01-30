// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package charset

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/text/transform"
)

func transformString(t transform.Transformer, s string) (string, error) ***REMOVED***
	r := transform.NewReader(strings.NewReader(s), t)
	b, err := ioutil.ReadAll(r)
	return string(b), err
***REMOVED***

type testCase struct ***REMOVED***
	utf8, other, otherEncoding string
***REMOVED***

// testCases for encoding and decoding.
var testCases = []testCase***REMOVED***
	***REMOVED***"Résumé", "Résumé", "utf8"***REMOVED***,
	***REMOVED***"Résumé", "R\xe9sum\xe9", "latin1"***REMOVED***,
	***REMOVED***"これは漢字です。", "S0\x8c0o0\"oW[g0Y0\x020", "UTF-16LE"***REMOVED***,
	***REMOVED***"これは漢字です。", "0S0\x8c0oo\"[W0g0Y0\x02", "UTF-16BE"***REMOVED***,
	***REMOVED***"Hello, world", "Hello, world", "ASCII"***REMOVED***,
	***REMOVED***"Gdańsk", "Gda\xf1sk", "ISO-8859-2"***REMOVED***,
	***REMOVED***"Ââ Čč Đđ Ŋŋ Õõ Šš Žž Åå Ää", "\xc2\xe2 \xc8\xe8 \xa9\xb9 \xaf\xbf \xd5\xf5 \xaa\xba \xac\xbc \xc5\xe5 \xc4\xe4", "ISO-8859-10"***REMOVED***,
	***REMOVED***"สำหรับ", "\xca\xd3\xcb\xc3\u047a", "ISO-8859-11"***REMOVED***,
	***REMOVED***"latviešu", "latvie\xf0u", "ISO-8859-13"***REMOVED***,
	***REMOVED***"Seònaid", "Se\xf2naid", "ISO-8859-14"***REMOVED***,
	***REMOVED***"€1 is cheap", "\xa41 is cheap", "ISO-8859-15"***REMOVED***,
	***REMOVED***"românește", "rom\xe2ne\xbate", "ISO-8859-16"***REMOVED***,
	***REMOVED***"nutraĵo", "nutra\xbco", "ISO-8859-3"***REMOVED***,
	***REMOVED***"Kalâdlit", "Kal\xe2dlit", "ISO-8859-4"***REMOVED***,
	***REMOVED***"русский", "\xe0\xe3\xe1\xe1\xda\xd8\xd9", "ISO-8859-5"***REMOVED***,
	***REMOVED***"ελληνικά", "\xe5\xeb\xeb\xe7\xed\xe9\xea\xdc", "ISO-8859-7"***REMOVED***,
	***REMOVED***"Kağan", "Ka\xf0an", "ISO-8859-9"***REMOVED***,
	***REMOVED***"Résumé", "R\x8esum\x8e", "macintosh"***REMOVED***,
	***REMOVED***"Gdańsk", "Gda\xf1sk", "windows-1250"***REMOVED***,
	***REMOVED***"русский", "\xf0\xf3\xf1\xf1\xea\xe8\xe9", "windows-1251"***REMOVED***,
	***REMOVED***"Résumé", "R\xe9sum\xe9", "windows-1252"***REMOVED***,
	***REMOVED***"ελληνικά", "\xe5\xeb\xeb\xe7\xed\xe9\xea\xdc", "windows-1253"***REMOVED***,
	***REMOVED***"Kağan", "Ka\xf0an", "windows-1254"***REMOVED***,
	***REMOVED***"עִבְרִית", "\xf2\xc4\xe1\xc0\xf8\xc4\xe9\xfa", "windows-1255"***REMOVED***,
	***REMOVED***"العربية", "\xc7\xe1\xda\xd1\xc8\xed\xc9", "windows-1256"***REMOVED***,
	***REMOVED***"latviešu", "latvie\xf0u", "windows-1257"***REMOVED***,
	***REMOVED***"Việt", "Vi\xea\xf2t", "windows-1258"***REMOVED***,
	***REMOVED***"สำหรับ", "\xca\xd3\xcb\xc3\u047a", "windows-874"***REMOVED***,
	***REMOVED***"русский", "\xd2\xd5\xd3\xd3\xcb\xc9\xca", "KOI8-R"***REMOVED***,
	***REMOVED***"українська", "\xd5\xcb\xd2\xc1\xa7\xce\xd3\xd8\xcb\xc1", "KOI8-U"***REMOVED***,
	***REMOVED***"Hello 常用國字標準字體表", "Hello \xb1`\xa5\u03b0\xea\xa6r\xbc\u0437\u01e6r\xc5\xe9\xaa\xed", "big5"***REMOVED***,
	***REMOVED***"Hello 常用國字標準字體表", "Hello \xb3\xa3\xd3\xc3\x87\xf8\xd7\xd6\x98\xcb\x9c\xca\xd7\xd6\xf3\x77\xb1\xed", "gbk"***REMOVED***,
	***REMOVED***"Hello 常用國字標準字體表", "Hello \xb3\xa3\xd3\xc3\x87\xf8\xd7\xd6\x98\xcb\x9c\xca\xd7\xd6\xf3\x77\xb1\xed", "gb18030"***REMOVED***,
	***REMOVED***"עִבְרִית", "\x81\x30\xfb\x30\x81\x30\xf6\x34\x81\x30\xf9\x33\x81\x30\xf6\x30\x81\x30\xfb\x36\x81\x30\xf6\x34\x81\x30\xfa\x31\x81\x30\xfb\x38", "gb18030"***REMOVED***,
	***REMOVED***"㧯", "\x82\x31\x89\x38", "gb18030"***REMOVED***,
	***REMOVED***"これは漢字です。", "\x82\xb1\x82\xea\x82\xcd\x8a\xbf\x8e\x9a\x82\xc5\x82\xb7\x81B", "SJIS"***REMOVED***,
	***REMOVED***"Hello, 世界!", "Hello, \x90\xa2\x8aE!", "SJIS"***REMOVED***,
	***REMOVED***"ｲｳｴｵｶ", "\xb2\xb3\xb4\xb5\xb6", "SJIS"***REMOVED***,
	***REMOVED***"これは漢字です。", "\xa4\xb3\xa4\xec\xa4\u03f4\xc1\xbb\xfa\xa4\u01e4\xb9\xa1\xa3", "EUC-JP"***REMOVED***,
	***REMOVED***"Hello, 世界!", "Hello, \x1b$B@$3&\x1b(B!", "ISO-2022-JP"***REMOVED***,
	***REMOVED***"다음과 같은 조건을 따라야 합니다: 저작자표시", "\xb4\xd9\xc0\xbd\xb0\xfa \xb0\xb0\xc0\xba \xc1\xb6\xb0\xc7\xc0\xbb \xb5\xfb\xb6\xf3\xbe\xdf \xc7մϴ\xd9: \xc0\xfa\xc0\xdb\xc0\xdaǥ\xbd\xc3", "EUC-KR"***REMOVED***,
***REMOVED***

func TestDecode(t *testing.T) ***REMOVED***
	testCases := append(testCases, []testCase***REMOVED***
		// Replace multi-byte maximum subpart of ill-formed subsequence with
		// single replacement character (WhatWG requirement).
		***REMOVED***"Rés\ufffdumé", "Rés\xe1\x80umé", "utf8"***REMOVED***,
	***REMOVED***...)
	for _, tc := range testCases ***REMOVED***
		e, _ := Lookup(tc.otherEncoding)
		if e == nil ***REMOVED***
			t.Errorf("%s: not found", tc.otherEncoding)
			continue
		***REMOVED***
		s, err := transformString(e.NewDecoder(), tc.other)
		if err != nil ***REMOVED***
			t.Errorf("%s: decode %q: %v", tc.otherEncoding, tc.other, err)
			continue
		***REMOVED***
		if s != tc.utf8 ***REMOVED***
			t.Errorf("%s: got %q, want %q", tc.otherEncoding, s, tc.utf8)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEncode(t *testing.T) ***REMOVED***
	testCases := append(testCases, []testCase***REMOVED***
		// Use Go-style replacement.
		***REMOVED***"Rés\xe1\x80umé", "Rés\ufffd\ufffdumé", "utf8"***REMOVED***,
		// U+0144 LATIN SMALL LETTER N WITH ACUTE not supported by encoding.
		***REMOVED***"Gdańsk", "Gda&#324;sk", "ISO-8859-11"***REMOVED***,
		***REMOVED***"\ufffd", "&#65533;", "ISO-8859-11"***REMOVED***,
		***REMOVED***"a\xe1\x80b", "a&#65533;&#65533;b", "ISO-8859-11"***REMOVED***,
	***REMOVED***...)
	for _, tc := range testCases ***REMOVED***
		e, _ := Lookup(tc.otherEncoding)
		if e == nil ***REMOVED***
			t.Errorf("%s: not found", tc.otherEncoding)
			continue
		***REMOVED***
		s, err := transformString(e.NewEncoder(), tc.utf8)
		if err != nil ***REMOVED***
			t.Errorf("%s: encode %q: %s", tc.otherEncoding, tc.utf8, err)
			continue
		***REMOVED***
		if s != tc.other ***REMOVED***
			t.Errorf("%s: got %q, want %q", tc.otherEncoding, s, tc.other)
		***REMOVED***
	***REMOVED***
***REMOVED***

var sniffTestCases = []struct ***REMOVED***
	filename, declared, want string
***REMOVED******REMOVED***
	***REMOVED***"HTTP-charset.html", "text/html; charset=iso-8859-15", "iso-8859-15"***REMOVED***,
	***REMOVED***"UTF-16LE-BOM.html", "", "utf-16le"***REMOVED***,
	***REMOVED***"UTF-16BE-BOM.html", "", "utf-16be"***REMOVED***,
	***REMOVED***"meta-content-attribute.html", "text/html", "iso-8859-15"***REMOVED***,
	***REMOVED***"meta-charset-attribute.html", "text/html", "iso-8859-15"***REMOVED***,
	***REMOVED***"No-encoding-declaration.html", "text/html", "utf-8"***REMOVED***,
	***REMOVED***"HTTP-vs-UTF-8-BOM.html", "text/html; charset=iso-8859-15", "utf-8"***REMOVED***,
	***REMOVED***"HTTP-vs-meta-content.html", "text/html; charset=iso-8859-15", "iso-8859-15"***REMOVED***,
	***REMOVED***"HTTP-vs-meta-charset.html", "text/html; charset=iso-8859-15", "iso-8859-15"***REMOVED***,
	***REMOVED***"UTF-8-BOM-vs-meta-content.html", "text/html", "utf-8"***REMOVED***,
	***REMOVED***"UTF-8-BOM-vs-meta-charset.html", "text/html", "utf-8"***REMOVED***,
***REMOVED***

func TestSniff(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl": // platforms that don't permit direct file system access
		t.Skipf("not supported on %q", runtime.GOOS)
	***REMOVED***

	for _, tc := range sniffTestCases ***REMOVED***
		content, err := ioutil.ReadFile("testdata/" + tc.filename)
		if err != nil ***REMOVED***
			t.Errorf("%s: error reading file: %v", tc.filename, err)
			continue
		***REMOVED***

		_, name, _ := DetermineEncoding(content, tc.declared)
		if name != tc.want ***REMOVED***
			t.Errorf("%s: got %q, want %q", tc.filename, name, tc.want)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReader(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl": // platforms that don't permit direct file system access
		t.Skipf("not supported on %q", runtime.GOOS)
	***REMOVED***

	for _, tc := range sniffTestCases ***REMOVED***
		content, err := ioutil.ReadFile("testdata/" + tc.filename)
		if err != nil ***REMOVED***
			t.Errorf("%s: error reading file: %v", tc.filename, err)
			continue
		***REMOVED***

		r, err := NewReader(bytes.NewReader(content), tc.declared)
		if err != nil ***REMOVED***
			t.Errorf("%s: error creating reader: %v", tc.filename, err)
			continue
		***REMOVED***

		got, err := ioutil.ReadAll(r)
		if err != nil ***REMOVED***
			t.Errorf("%s: error reading from charset.NewReader: %v", tc.filename, err)
			continue
		***REMOVED***

		e, _ := Lookup(tc.want)
		want, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(content), e.NewDecoder()))
		if err != nil ***REMOVED***
			t.Errorf("%s: error decoding with hard-coded charset name: %v", tc.filename, err)
			continue
		***REMOVED***

		if !bytes.Equal(got, want) ***REMOVED***
			t.Errorf("%s: got %q, want %q", tc.filename, got, want)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

var metaTestCases = []struct ***REMOVED***
	meta, want string
***REMOVED******REMOVED***
	***REMOVED***"", ""***REMOVED***,
	***REMOVED***"text/html", ""***REMOVED***,
	***REMOVED***"text/html; charset utf-8", ""***REMOVED***,
	***REMOVED***"text/html; charset=latin-2", "latin-2"***REMOVED***,
	***REMOVED***"text/html; charset; charset = utf-8", "utf-8"***REMOVED***,
	***REMOVED***`charset="big5"`, "big5"***REMOVED***,
	***REMOVED***"charset='shift_jis'", "shift_jis"***REMOVED***,
***REMOVED***

func TestFromMeta(t *testing.T) ***REMOVED***
	for _, tc := range metaTestCases ***REMOVED***
		got := fromMetaElement(tc.meta)
		if got != tc.want ***REMOVED***
			t.Errorf("%q: got %q, want %q", tc.meta, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestXML(t *testing.T) ***REMOVED***
	const s = "<?xml version=\"1.0\" encoding=\"windows-1252\"?><a><Word>r\xe9sum\xe9</Word></a>"

	d := xml.NewDecoder(strings.NewReader(s))
	d.CharsetReader = NewReaderLabel

	var a struct ***REMOVED***
		Word string
	***REMOVED***
	err := d.Decode(&a)
	if err != nil ***REMOVED***
		t.Fatalf("Decode: %v", err)
	***REMOVED***

	want := "résumé"
	if a.Word != want ***REMOVED***
		t.Errorf("got %q, want %q", a.Word, want)
	***REMOVED***
***REMOVED***
