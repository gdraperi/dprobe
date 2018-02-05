// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runes_test

import (
	"fmt"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
)

func ExampleRemove() ***REMOVED***
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ := transform.String(t, "résumé")
	fmt.Println(s)

	// Output:
	// resume
***REMOVED***

func ExampleMap() ***REMOVED***
	replaceHyphens := runes.Map(func(r rune) rune ***REMOVED***
		if unicode.Is(unicode.Hyphen, r) ***REMOVED***
			return '|'
		***REMOVED***
		return r
	***REMOVED***)
	s, _, _ := transform.String(replaceHyphens, "a-b‐c⸗d﹣e")
	fmt.Println(s)

	// Output:
	// a|b|c|d|e
***REMOVED***

func ExampleIn() ***REMOVED***
	// Convert Latin characters to their canonical form, while keeping other
	// width distinctions.
	t := runes.If(runes.In(unicode.Latin), width.Fold, nil)
	s, _, _ := transform.String(t, "ｱﾙｱﾉﾘｳ tech / アルアノリウ ｔｅｃｈ")
	fmt.Println(s)

	// Output:
	// ｱﾙｱﾉﾘｳ tech / アルアノリウ tech
***REMOVED***

func ExampleIf() ***REMOVED***
	// Widen everything but ASCII.
	isASCII := func(r rune) bool ***REMOVED*** return r <= unicode.MaxASCII ***REMOVED***
	t := runes.If(runes.Predicate(isASCII), nil, width.Widen)
	s, _, _ := transform.String(t, "ｱﾙｱﾉﾘｳ tech / 中國 / 5₩")
	fmt.Println(s)

	// Output:
	// アルアノリウ tech / 中國 / 5￦
***REMOVED***
