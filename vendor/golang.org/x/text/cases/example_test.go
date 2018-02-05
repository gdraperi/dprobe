// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cases_test

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Example() ***REMOVED***
	src := []string***REMOVED***
		"hello world!",
		"i with dot",
		"'n ijsberg",
		"here comes O'Brian",
	***REMOVED***
	for _, c := range []cases.Caser***REMOVED***
		cases.Lower(language.Und),
		cases.Upper(language.Turkish),
		cases.Title(language.Dutch),
		cases.Title(language.Und, cases.NoLower),
	***REMOVED*** ***REMOVED***
		fmt.Println()
		for _, s := range src ***REMOVED***
			fmt.Println(c.String(s))
		***REMOVED***
	***REMOVED***

	// Output:
	// hello world!
	// i with dot
	// 'n ijsberg
	// here comes o'brian
	//
	// HELLO WORLD!
	// İ WİTH DOT
	// 'N İJSBERG
	// HERE COMES O'BRİAN
	//
	// Hello World!
	// I With Dot
	// 'n IJsberg
	// Here Comes O'brian
	//
	// Hello World!
	// I With Dot
	// 'N Ijsberg
	// Here Comes O'Brian
***REMOVED***
