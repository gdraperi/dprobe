// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm_test

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// EqualSimple uses a norm.Iter to compare two non-normalized
// strings for equivalence.
func EqualSimple(a, b string) bool ***REMOVED***
	var ia, ib norm.Iter
	ia.InitString(norm.NFKD, a)
	ib.InitString(norm.NFKD, b)
	for !ia.Done() && !ib.Done() ***REMOVED***
		if !bytes.Equal(ia.Next(), ib.Next()) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return ia.Done() && ib.Done()
***REMOVED***

// FindPrefix finds the longest common prefix of ASCII characters
// of a and b.
func FindPrefix(a, b string) int ***REMOVED***
	i := 0
	for ; i < len(a) && i < len(b) && a[i] < utf8.RuneSelf && a[i] == b[i]; i++ ***REMOVED***
	***REMOVED***
	return i
***REMOVED***

// EqualOpt is like EqualSimple, but optimizes the special
// case for ASCII characters.
func EqualOpt(a, b string) bool ***REMOVED***
	n := FindPrefix(a, b)
	a, b = a[n:], b[n:]
	var ia, ib norm.Iter
	ia.InitString(norm.NFKD, a)
	ib.InitString(norm.NFKD, b)
	for !ia.Done() && !ib.Done() ***REMOVED***
		if !bytes.Equal(ia.Next(), ib.Next()) ***REMOVED***
			return false
		***REMOVED***
		if n := int64(FindPrefix(a[ia.Pos():], b[ib.Pos():])); n != 0 ***REMOVED***
			ia.Seek(n, 1)
			ib.Seek(n, 1)
		***REMOVED***
	***REMOVED***
	return ia.Done() && ib.Done()
***REMOVED***

var compareTests = []struct***REMOVED*** a, b string ***REMOVED******REMOVED***
	***REMOVED***"aaa", "aaa"***REMOVED***,
	***REMOVED***"aaa", "aab"***REMOVED***,
	***REMOVED***"a\u0300a", "\u00E0a"***REMOVED***,
	***REMOVED***"a\u0300\u0320b", "a\u0320\u0300b"***REMOVED***,
	***REMOVED***"\u1E0A\u0323", "\x44\u0323\u0307"***REMOVED***,
	// A character that decomposes into multiple segments
	// spans several iterations.
	***REMOVED***"\u3304", "\u30A4\u30CB\u30F3\u30AF\u3099"***REMOVED***,
***REMOVED***

func ExampleIter() ***REMOVED***
	for i, t := range compareTests ***REMOVED***
		r0 := EqualSimple(t.a, t.b)
		r1 := EqualOpt(t.a, t.b)
		fmt.Printf("%d: %v %v\n", i, r0, r1)
	***REMOVED***
	// Output:
	// 0: true true
	// 1: false false
	// 2: true true
	// 3: true true
	// 4: true true
	// 5: true true
***REMOVED***
