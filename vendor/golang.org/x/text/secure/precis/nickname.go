// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package precis

import (
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/transform"
)

type nickAdditionalMapping struct ***REMOVED***
	// TODO: This transformer needs to be stateless somehow…
	notStart  bool
	prevSpace bool
***REMOVED***

func (t *nickAdditionalMapping) Reset() ***REMOVED***
	t.prevSpace = false
	t.notStart = false
***REMOVED***

func (t *nickAdditionalMapping) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	// RFC 8266 §2.1.  Rules
	//
	// 2.  Additional Mapping Rule: The additional mapping rule consists of
	//     the following sub-rules.
	//
	//     a.  Map any instances of non-ASCII space to SPACE (U+0020); a
	//         non-ASCII space is any Unicode code point having a general
	//         category of "Zs", naturally with the exception of SPACE
	//         (U+0020).  (The inclusion of only ASCII space prevents
	//         confusion with various non-ASCII space code points, many of
	//         which are difficult to reproduce across different input
	//         methods.)
	//
	//     b.  Remove any instances of the ASCII space character at the
	//         beginning or end of a nickname (e.g., "stpeter " is mapped to
	//         "stpeter").
	//
	//     c.  Map interior sequences of more than one ASCII space character
	//         to a single ASCII space character (e.g., "St  Peter" is
	//         mapped to "St Peter").
	for nSrc < len(src) ***REMOVED***
		r, size := utf8.DecodeRune(src[nSrc:])
		if size == 0 ***REMOVED*** // Incomplete UTF-8 encoding
			if !atEOF ***REMOVED***
				return nDst, nSrc, transform.ErrShortSrc
			***REMOVED***
			size = 1
		***REMOVED***
		if unicode.Is(unicode.Zs, r) ***REMOVED***
			t.prevSpace = true
		***REMOVED*** else ***REMOVED***
			if t.prevSpace && t.notStart ***REMOVED***
				dst[nDst] = ' '
				nDst += 1
			***REMOVED***
			if size != copy(dst[nDst:], src[nSrc:nSrc+size]) ***REMOVED***
				nDst += size
				return nDst, nSrc, transform.ErrShortDst
			***REMOVED***
			nDst += size
			t.prevSpace = false
			t.notStart = true
		***REMOVED***
		nSrc += size
	***REMOVED***
	return nDst, nSrc, nil
***REMOVED***
