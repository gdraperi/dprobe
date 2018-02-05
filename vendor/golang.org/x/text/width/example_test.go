// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package width_test

import (
	"fmt"

	"golang.org/x/text/width"
)

func ExampleTransformer_fold() ***REMOVED***
	s := "abｦ￦￮￥Ａ"
	f := width.Fold.String(s)
	fmt.Printf("%U: %s\n", []rune(s), s)
	fmt.Printf("%U: %s\n", []rune(f), f)

	// Output:
	// [U+0061 U+0062 U+FF66 U+FFE6 U+FFEE U+FFE5 U+FF21]: abｦ￦￮￥Ａ
	// [U+0061 U+0062 U+30F2 U+20A9 U+25CB U+00A5 U+0041]: abヲ₩○¥A
***REMOVED***

func ExampleTransformer_widen() ***REMOVED***
	s := "ab¥ｦ₩￮"
	w := width.Widen.String(s)
	fmt.Printf("%U: %s\n", []rune(s), s)
	fmt.Printf("%U: %s\n", []rune(w), w)

	// Output:
	// [U+0061 U+0062 U+00A5 U+FF66 U+20A9 U+FFEE]: ab¥ｦ₩￮
	// [U+FF41 U+FF42 U+FFE5 U+30F2 U+FFE6 U+25CB]: ａｂ￥ヲ￦○
***REMOVED***

func ExampleTransformer_narrow() ***REMOVED***
	s := "abヲ￦○￥Ａ"
	n := width.Narrow.String(s)
	fmt.Printf("%U: %s\n", []rune(s), s)
	fmt.Printf("%U: %s\n", []rune(n), n)

	// Ambiguous characters with a halfwidth equivalent get mapped as well.
	s = "←"
	n = width.Narrow.String(s)
	fmt.Printf("%U: %s\n", []rune(s), s)
	fmt.Printf("%U: %s\n", []rune(n), n)

	// Output:
	// [U+0061 U+0062 U+30F2 U+FFE6 U+25CB U+FFE5 U+FF21]: abヲ￦○￥Ａ
	// [U+0061 U+0062 U+FF66 U+20A9 U+FFEE U+00A5 U+0041]: abｦ₩￮¥A
	// [U+2190]: ←
	// [U+FFE9]: ￩
***REMOVED***
