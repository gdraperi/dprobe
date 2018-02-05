// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import (
	"strings"
	"testing"
)

var strdata = []string***REMOVED***
	"aa  ",
	"aaa ",
	"aaaa",
	"aaab",
	"aab ",
	"ab  ",
	"ba  ",
	"xxxx",
	"\xff\xff\xff\xff",
***REMOVED***

var testCases = map[string]int***REMOVED***
	"a":    0,
	"aa":   0,
	"aaa":  1,
	"aa ":  0,
	"aaaa": 2,
	"aaab": 3,
	"b":    6,
	"ba":   6,
	"    ": -1,
	"aaax": -1,
	"bbbb": -1,
	"zzzz": -1,
***REMOVED***

func TestIndex(t *testing.T) ***REMOVED***
	index := Index(strings.Join(strdata, ""))
	for k, v := range testCases ***REMOVED***
		if i := index.Index([]byte(k)); i != v ***REMOVED***
			t.Errorf("%s: got %d; want %d", k, i, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFixCase(t *testing.T) ***REMOVED***
	tests := []string***REMOVED***
		"aaaa", "AbCD", "abcd",
		"Zzzz", "AbCD", "Abcd",
		"Zzzz", "AbC", "",
		"XXX", "ab ", "",
		"XXX", "usd", "USD",
		"cmn", "AB ", "",
		"gsw", "CMN", "cmn",
	***REMOVED***
	for tc := tests; len(tc) > 0; tc = tc[3:] ***REMOVED***
		b := []byte(tc[1])
		if !FixCase(tc[0], b) ***REMOVED***
			b = nil
		***REMOVED***
		if string(b) != tc[2] ***REMOVED***
			t.Errorf("FixCase(%q, %q) = %q; want %q", tc[0], tc[1], b, tc[2])
		***REMOVED***
	***REMOVED***
***REMOVED***
