// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"golang.org/x/text/internal/gen"
)

type group struct ***REMOVED***
	Encodings []struct ***REMOVED***
		Labels []string
		Name   string
	***REMOVED***
***REMOVED***

func main() ***REMOVED***
	gen.Init()

	r := gen.Open("https://encoding.spec.whatwg.org", "whatwg", "encodings.json")
	var groups []group
	if err := json.NewDecoder(r).Decode(&groups); err != nil ***REMOVED***
		log.Fatalf("Error reading encodings.json: %v", err)
	***REMOVED***

	w := &bytes.Buffer***REMOVED******REMOVED***
	fmt.Fprintln(w, "type htmlEncoding byte")
	fmt.Fprintln(w, "const (")
	for i, g := range groups ***REMOVED***
		for _, e := range g.Encodings ***REMOVED***
			key := strings.ToLower(e.Name)
			name := consts[key]
			if name == "" ***REMOVED***
				log.Fatalf("No const defined for %s.", key)
			***REMOVED***
			if i == 0 ***REMOVED***
				fmt.Fprintf(w, "%s htmlEncoding = iota\n", name)
			***REMOVED*** else ***REMOVED***
				fmt.Fprintf(w, "%s\n", name)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	fmt.Fprintln(w, "numEncodings")
	fmt.Fprint(w, ")\n\n")

	fmt.Fprintln(w, "var canonical = [numEncodings]string***REMOVED***")
	for _, g := range groups ***REMOVED***
		for _, e := range g.Encodings ***REMOVED***
			fmt.Fprintf(w, "%q,\n", strings.ToLower(e.Name))
		***REMOVED***
	***REMOVED***
	fmt.Fprint(w, "***REMOVED***\n\n")

	fmt.Fprintln(w, "var nameMap = map[string]htmlEncoding***REMOVED***")
	for _, g := range groups ***REMOVED***
		for _, e := range g.Encodings ***REMOVED***
			for _, l := range e.Labels ***REMOVED***
				key := strings.ToLower(e.Name)
				name := consts[key]
				fmt.Fprintf(w, "%q: %s,\n", l, name)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	fmt.Fprint(w, "***REMOVED***\n\n")

	var tags []string
	fmt.Fprintln(w, "var localeMap = []htmlEncoding***REMOVED***")
	for _, loc := range locales ***REMOVED***
		tags = append(tags, loc.tag)
		fmt.Fprintf(w, "%s, // %s \n", consts[loc.name], loc.tag)
	***REMOVED***
	fmt.Fprint(w, "***REMOVED***\n\n")

	fmt.Fprintf(w, "const locales = %q\n", strings.Join(tags, " "))

	gen.WriteGoFile("tables.go", "htmlindex", w.Bytes())
***REMOVED***

// consts maps canonical encoding name to internal constant.
var consts = map[string]string***REMOVED***
	"utf-8":          "utf8",
	"ibm866":         "ibm866",
	"iso-8859-2":     "iso8859_2",
	"iso-8859-3":     "iso8859_3",
	"iso-8859-4":     "iso8859_4",
	"iso-8859-5":     "iso8859_5",
	"iso-8859-6":     "iso8859_6",
	"iso-8859-7":     "iso8859_7",
	"iso-8859-8":     "iso8859_8",
	"iso-8859-8-i":   "iso8859_8I",
	"iso-8859-10":    "iso8859_10",
	"iso-8859-13":    "iso8859_13",
	"iso-8859-14":    "iso8859_14",
	"iso-8859-15":    "iso8859_15",
	"iso-8859-16":    "iso8859_16",
	"koi8-r":         "koi8r",
	"koi8-u":         "koi8u",
	"macintosh":      "macintosh",
	"windows-874":    "windows874",
	"windows-1250":   "windows1250",
	"windows-1251":   "windows1251",
	"windows-1252":   "windows1252",
	"windows-1253":   "windows1253",
	"windows-1254":   "windows1254",
	"windows-1255":   "windows1255",
	"windows-1256":   "windows1256",
	"windows-1257":   "windows1257",
	"windows-1258":   "windows1258",
	"x-mac-cyrillic": "macintoshCyrillic",
	"gbk":            "gbk",
	"gb18030":        "gb18030",
	// "hz-gb-2312":     "hzgb2312", // Was removed from WhatWG
	"big5":           "big5",
	"euc-jp":         "eucjp",
	"iso-2022-jp":    "iso2022jp",
	"shift_jis":      "shiftJIS",
	"euc-kr":         "euckr",
	"replacement":    "replacement",
	"utf-16be":       "utf16be",
	"utf-16le":       "utf16le",
	"x-user-defined": "xUserDefined",
***REMOVED***

// locales is taken from
// https://html.spec.whatwg.org/multipage/syntax.html#encoding-sniffing-algorithm.
var locales = []struct***REMOVED*** tag, name string ***REMOVED******REMOVED***
	// The default value. Explicitly state latin to benefit from the exact
	// script option, while still making 1252 the default encoding for languages
	// written in Latin script.
	***REMOVED***"und_Latn", "windows-1252"***REMOVED***,
	***REMOVED***"ar", "windows-1256"***REMOVED***,
	***REMOVED***"ba", "windows-1251"***REMOVED***,
	***REMOVED***"be", "windows-1251"***REMOVED***,
	***REMOVED***"bg", "windows-1251"***REMOVED***,
	***REMOVED***"cs", "windows-1250"***REMOVED***,
	***REMOVED***"el", "iso-8859-7"***REMOVED***,
	***REMOVED***"et", "windows-1257"***REMOVED***,
	***REMOVED***"fa", "windows-1256"***REMOVED***,
	***REMOVED***"he", "windows-1255"***REMOVED***,
	***REMOVED***"hr", "windows-1250"***REMOVED***,
	***REMOVED***"hu", "iso-8859-2"***REMOVED***,
	***REMOVED***"ja", "shift_jis"***REMOVED***,
	***REMOVED***"kk", "windows-1251"***REMOVED***,
	***REMOVED***"ko", "euc-kr"***REMOVED***,
	***REMOVED***"ku", "windows-1254"***REMOVED***,
	***REMOVED***"ky", "windows-1251"***REMOVED***,
	***REMOVED***"lt", "windows-1257"***REMOVED***,
	***REMOVED***"lv", "windows-1257"***REMOVED***,
	***REMOVED***"mk", "windows-1251"***REMOVED***,
	***REMOVED***"pl", "iso-8859-2"***REMOVED***,
	***REMOVED***"ru", "windows-1251"***REMOVED***,
	***REMOVED***"sah", "windows-1251"***REMOVED***,
	***REMOVED***"sk", "windows-1250"***REMOVED***,
	***REMOVED***"sl", "iso-8859-2"***REMOVED***,
	***REMOVED***"sr", "windows-1251"***REMOVED***,
	***REMOVED***"tg", "windows-1251"***REMOVED***,
	***REMOVED***"th", "windows-874"***REMOVED***,
	***REMOVED***"tr", "windows-1254"***REMOVED***,
	***REMOVED***"tt", "windows-1251"***REMOVED***,
	***REMOVED***"uk", "windows-1251"***REMOVED***,
	***REMOVED***"vi", "windows-1258"***REMOVED***,
	***REMOVED***"zh-hans", "gb18030"***REMOVED***,
	***REMOVED***"zh-hant", "big5"***REMOVED***,
***REMOVED***
