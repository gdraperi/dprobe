// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.10

package bidirule

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/bidi"
)

var testCases = [][]ruleTest***REMOVED***
	// Go-specific rules.
	// Invalid UTF-8 is invalid.
	0: []ruleTest***REMOVED******REMOVED***
		in:  "",
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  "\x80",
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   0,
	***REMOVED***, ***REMOVED***
		in:  "\xcc",
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   0,
	***REMOVED***, ***REMOVED***
		in:  "abc\x80",
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   3,
	***REMOVED***, ***REMOVED***
		in:  "abc\xcc",
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   3,
	***REMOVED***, ***REMOVED***
		in:  "abc\xccdef",
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   3,
	***REMOVED***, ***REMOVED***
		in:  "\xccdef",
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   0,
	***REMOVED***, ***REMOVED***
		in:  strR + "\x80",
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   len(strR),
	***REMOVED***, ***REMOVED***
		in:  strR + "\xcc",
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   len(strR),
	***REMOVED***, ***REMOVED***
		in:  strAL + "\xcc" + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   len(strAL),
	***REMOVED***, ***REMOVED***
		in:  "\xcc" + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   0,
	***REMOVED******REMOVED***,

	// Rule 2.1: The first character must be a character with Bidi property L,
	// R, or AL.  If it has the R or AL property, it is an RTL label; if it has
	// the L property, it is an LTR label.
	1: []ruleTest***REMOVED******REMOVED***
		in:  strL,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAN,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strEN,
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   len(strEN),
	***REMOVED***, ***REMOVED***
		in:  strES,
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   len(strES),
	***REMOVED***, ***REMOVED***
		in:  strET,
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   len(strET),
	***REMOVED***, ***REMOVED***
		in:  strCS,
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   len(strCS),
	***REMOVED***, ***REMOVED***
		in:  strNSM,
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   len(strNSM),
	***REMOVED***, ***REMOVED***
		in:  strBN,
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   len(strBN),
	***REMOVED***, ***REMOVED***
		in:  strB,
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   len(strB),
	***REMOVED***, ***REMOVED***
		in:  strS,
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   len(strS),
	***REMOVED***, ***REMOVED***
		in:  strWS,
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   len(strWS),
	***REMOVED***, ***REMOVED***
		in:  strON,
		dir: bidi.LeftToRight,
		err: ErrInvalid,
		n:   len(strON),
	***REMOVED***, ***REMOVED***
		in:  strEN + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   3,
	***REMOVED***, ***REMOVED***
		in:  strES + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   2,
	***REMOVED***, ***REMOVED***
		in:  strET + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   1,
	***REMOVED***, ***REMOVED***
		in:  strCS + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   1,
	***REMOVED***, ***REMOVED***
		in:  strNSM + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   2,
	***REMOVED***, ***REMOVED***
		in:  strBN + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   3,
	***REMOVED***, ***REMOVED***
		in:  strB + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   3,
	***REMOVED***, ***REMOVED***
		in:  strS + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   1,
	***REMOVED***, ***REMOVED***
		in:  strWS + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   1,
	***REMOVED***, ***REMOVED***
		in:  strON + strR,
		dir: bidi.RightToLeft,
		err: ErrInvalid,
		n:   1,
	***REMOVED******REMOVED***,

	// Rule 2.2: In an RTL label, only characters with the Bidi properties R,
	// AL, AN, EN, ES, CS, ET, ON, BN, or NSM are allowed.
	2: []ruleTest***REMOVED******REMOVED***
		in:  strR + strR + strAL,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strAL + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strAN + strAL,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strEN + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strES + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strCS + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strET + strAL,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strON + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strBN + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strNSM + strAL,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strL + strR,
		dir: bidi.RightToLeft,
		n:   len(strR),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strB + strR,
		dir: bidi.RightToLeft,
		n:   len(strR),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strS + strAL,
		dir: bidi.RightToLeft,
		n:   len(strR),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strWS + strAL,
		dir: bidi.RightToLeft,
		n:   len(strR),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strR + strAL,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strAL + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strAN + strAL,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strEN + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strES + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strCS + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strET + strAL,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strON + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strBN + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strNSM + strAL,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strL + strR,
		dir: bidi.RightToLeft,
		n:   len(strAL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strB + strR,
		dir: bidi.RightToLeft,
		n:   len(strAL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strS + strAL,
		dir: bidi.RightToLeft,
		n:   len(strAL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strWS + strAL,
		dir: bidi.RightToLeft,
		n:   len(strAL),
		err: ErrInvalid,
	***REMOVED******REMOVED***,

	// Rule 2.3: In an RTL label, the end of the label must be a character with
	// Bidi property R, AL, EN, or AN, followed by zero or more characters with
	// Bidi property NSM.
	3: []ruleTest***REMOVED******REMOVED***
		in:  strR + strNSM,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strAL + strNSM,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strEN + strNSM + strNSM,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strAN,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strR + strES + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strR + strES + strNSM),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strCS + strNSM + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strR + strCS + strNSM + strNSM),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strET,
		dir: bidi.RightToLeft,
		n:   len(strR + strET),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strON + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strR + strON + strNSM),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strBN + strNSM + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strR + strBN + strNSM + strNSM),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strL + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strR),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strB + strNSM + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strR),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strS,
		dir: bidi.RightToLeft,
		n:   len(strR),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strWS,
		dir: bidi.RightToLeft,
		n:   len(strR),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strNSM,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strR,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strAL + strNSM,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strEN + strNSM + strNSM,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strAN,
		dir: bidi.RightToLeft,
	***REMOVED***, ***REMOVED***
		in:  strAL + strES + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strAL + strES + strNSM),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strCS + strNSM + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strAL + strCS + strNSM + strNSM),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strET,
		dir: bidi.RightToLeft,
		n:   len(strAL + strET),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strON + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strAL + strON + strNSM),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strBN + strNSM + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strAL + strBN + strNSM + strNSM),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strL + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strAL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strB + strNSM + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strAL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strS,
		dir: bidi.RightToLeft,
		n:   len(strAL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strWS,
		dir: bidi.RightToLeft,
		n:   len(strAL),
		err: ErrInvalid,
	***REMOVED******REMOVED***,

	// Rule 2.4: In an RTL label, if an EN is present, no AN may be present,
	// and vice versa.
	4: []ruleTest***REMOVED******REMOVED***
		in:  strR + strEN + strAN,
		dir: bidi.RightToLeft,
		n:   len(strR + strEN),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strR + strAN + strEN + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strR + strAN),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strEN + strAN,
		dir: bidi.RightToLeft,
		n:   len(strAL + strEN),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strAL + strAN + strEN + strNSM,
		dir: bidi.RightToLeft,
		n:   len(strAL + strAN),
		err: ErrInvalid,
	***REMOVED******REMOVED***,

	// Rule 2.5: In an LTR label, only characters with the Bidi properties L,
	// EN, ES, CS, ET, ON, BN, or NSM are allowed.
	5: []ruleTest***REMOVED******REMOVED***
		in:  strL + strL + strL,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strEN + strL,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strES + strL,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strCS + strL,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strET + strL,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strON + strL,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strBN + strL,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strNSM + strL,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strR + strL,
		dir: bidi.RightToLeft,
		n:   len(strL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strAL + strL,
		dir: bidi.RightToLeft,
		n:   len(strL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strAN + strL,
		dir: bidi.RightToLeft,
		n:   len(strL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strB + strL,
		dir: bidi.LeftToRight,
		n:   len(strL + strB + strL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strB + strL + strR,
		dir: bidi.RightToLeft,
		n:   len(strL + strB + strL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strS + strL,
		dir: bidi.LeftToRight,
		n:   len(strL + strS + strL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strS + strL + strR,
		dir: bidi.RightToLeft,
		n:   len(strL + strS + strL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strWS + strL,
		dir: bidi.LeftToRight,
		n:   len(strL + strWS + strL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strWS + strL + strR,
		dir: bidi.RightToLeft,
		n:   len(strL + strWS + strL),
		err: ErrInvalid,
	***REMOVED******REMOVED***,

	// Rule 2.6: In an LTR label, the end of the label must be a character with
	// Bidi property L or EN, followed by zero or more characters with Bidi
	// property NSM.
	6: []ruleTest***REMOVED******REMOVED***
		in:  strL,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strNSM,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strNSM + strNSM,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strEN,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strEN + strNSM,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strEN + strNSM + strNSM,
		dir: bidi.LeftToRight,
	***REMOVED***, ***REMOVED***
		in:  strL + strES,
		dir: bidi.LeftToRight,
		n:   len(strL + strES),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strES + strR,
		dir: bidi.RightToLeft,
		n:   len(strL + strES),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strCS,
		dir: bidi.LeftToRight,
		n:   len(strL + strCS),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strCS + strR,
		dir: bidi.RightToLeft,
		n:   len(strL + strCS),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strET,
		dir: bidi.LeftToRight,
		n:   len(strL + strET),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strET + strR,
		dir: bidi.RightToLeft,
		n:   len(strL + strET),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strON,
		dir: bidi.LeftToRight,
		n:   len(strL + strON),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strON + strR,
		dir: bidi.RightToLeft,
		n:   len(strL + strON),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strBN,
		dir: bidi.LeftToRight,
		n:   len(strL + strBN),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strBN + strR,
		dir: bidi.RightToLeft,
		n:   len(strL + strBN),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strR,
		dir: bidi.RightToLeft,
		n:   len(strL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strAL,
		dir: bidi.RightToLeft,
		n:   len(strL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strAN,
		dir: bidi.RightToLeft,
		n:   len(strL),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strB,
		dir: bidi.LeftToRight,
		n:   len(strL + strB),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strB + strR,
		dir: bidi.RightToLeft,
		n:   len(strL + strB),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strS,
		dir: bidi.LeftToRight,
		n:   len(strL + strS),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strS + strR,
		dir: bidi.RightToLeft,
		n:   len(strL + strS),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strWS,
		dir: bidi.LeftToRight,
		n:   len(strL + strWS),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  strL + strWS + strR,
		dir: bidi.RightToLeft,
		n:   len(strL + strWS),
		err: ErrInvalid,
	***REMOVED******REMOVED***,

	// Incremental processing.
	9: []ruleTest***REMOVED******REMOVED***
		in:  "e\u0301", // é
		dir: bidi.LeftToRight,

		pSrc: 2,
		nSrc: 1,
		err0: transform.ErrShortSrc,
	***REMOVED***, ***REMOVED***
		in:  "e\u1000f", // é
		dir: bidi.LeftToRight,

		pSrc: 3,
		nSrc: 1,
		err0: transform.ErrShortSrc,
	***REMOVED***, ***REMOVED***
		// Remain invalid once invalid.
		in:  strR + "ab",
		dir: bidi.RightToLeft,
		n:   len(strR),
		err: ErrInvalid,

		pSrc: len(strR) + 1,
		nSrc: len(strR),
		err0: ErrInvalid,
	***REMOVED***, ***REMOVED***
		// Short destination
		in:  "abcdefghij",
		dir: bidi.LeftToRight,

		pSrc:  10,
		szDst: 5,
		nSrc:  5,
		err0:  transform.ErrShortDst,
	***REMOVED***, ***REMOVED***
		in:  "\U000102f7",
		dir: bidi.LeftToRight,
		n:   len("\U000102f7"),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		// Short destination splitting input rune
		in:  "e\u0301",
		dir: bidi.LeftToRight,

		pSrc:  3,
		szDst: 2,
		nSrc:  1,
		err0:  transform.ErrShortDst,
	***REMOVED***, ***REMOVED***
		// Unicode 10.0.0 IDNA test string.
		in:  "FAX\u2a77\U0001d186",
		dir: bidi.LeftToRight,
		n:   len("FAX\u2a77\U0001d186"),
		err: ErrInvalid,
	***REMOVED***, ***REMOVED***
		in:  "\x80\u0660",
		dir: bidi.RightToLeft,
		n:   0,
		err: ErrInvalid,
	***REMOVED******REMOVED***,
***REMOVED***
