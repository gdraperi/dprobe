// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package japanese

import (
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/transform"
)

// EUCJP is the EUC-JP encoding.
var EUCJP encoding.Encoding = &eucJP

var eucJP = internal.Encoding***REMOVED***
	&internal.SimpleEncoding***REMOVED***eucJPDecoder***REMOVED******REMOVED***, eucJPEncoder***REMOVED******REMOVED******REMOVED***,
	"EUC-JP",
	identifier.EUCPkdFmtJapanese,
***REMOVED***

type eucJPDecoder struct***REMOVED*** transform.NopResetter ***REMOVED***

// See https://encoding.spec.whatwg.org/#euc-jp-decoder.
func (eucJPDecoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	r, size := rune(0), 0
loop:
	for ; nSrc < len(src); nSrc += size ***REMOVED***
		switch c0 := src[nSrc]; ***REMOVED***
		case c0 < utf8.RuneSelf:
			r, size = rune(c0), 1

		case c0 == 0x8e:
			if nSrc+1 >= len(src) ***REMOVED***
				if !atEOF ***REMOVED***
					err = transform.ErrShortSrc
					break loop
				***REMOVED***
				r, size = utf8.RuneError, 1
				break
			***REMOVED***
			c1 := src[nSrc+1]
			switch ***REMOVED***
			case c1 < 0xa1:
				r, size = utf8.RuneError, 1
			case c1 > 0xdf:
				r, size = utf8.RuneError, 2
				if c1 == 0xff ***REMOVED***
					size = 1
				***REMOVED***
			default:
				r, size = rune(c1)+(0xff61-0xa1), 2
			***REMOVED***
		case c0 == 0x8f:
			if nSrc+2 >= len(src) ***REMOVED***
				if !atEOF ***REMOVED***
					err = transform.ErrShortSrc
					break loop
				***REMOVED***
				r, size = utf8.RuneError, 1
				if p := nSrc + 1; p < len(src) && 0xa1 <= src[p] && src[p] < 0xfe ***REMOVED***
					size = 2
				***REMOVED***
				break
			***REMOVED***
			c1 := src[nSrc+1]
			if c1 < 0xa1 || 0xfe < c1 ***REMOVED***
				r, size = utf8.RuneError, 1
				break
			***REMOVED***
			c2 := src[nSrc+2]
			if c2 < 0xa1 || 0xfe < c2 ***REMOVED***
				r, size = utf8.RuneError, 2
				break
			***REMOVED***
			r, size = utf8.RuneError, 3
			if i := int(c1-0xa1)*94 + int(c2-0xa1); i < len(jis0212Decode) ***REMOVED***
				r = rune(jis0212Decode[i])
				if r == 0 ***REMOVED***
					r = utf8.RuneError
				***REMOVED***
			***REMOVED***

		case 0xa1 <= c0 && c0 <= 0xfe:
			if nSrc+1 >= len(src) ***REMOVED***
				if !atEOF ***REMOVED***
					err = transform.ErrShortSrc
					break loop
				***REMOVED***
				r, size = utf8.RuneError, 1
				break
			***REMOVED***
			c1 := src[nSrc+1]
			if c1 < 0xa1 || 0xfe < c1 ***REMOVED***
				r, size = utf8.RuneError, 1
				break
			***REMOVED***
			r, size = utf8.RuneError, 2
			if i := int(c0-0xa1)*94 + int(c1-0xa1); i < len(jis0208Decode) ***REMOVED***
				r = rune(jis0208Decode[i])
				if r == 0 ***REMOVED***
					r = utf8.RuneError
				***REMOVED***
			***REMOVED***

		default:
			r, size = utf8.RuneError, 1
		***REMOVED***

		if nDst+utf8.RuneLen(r) > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break loop
		***REMOVED***
		nDst += utf8.EncodeRune(dst[nDst:], r)
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***

type eucJPEncoder struct***REMOVED*** transform.NopResetter ***REMOVED***

func (eucJPEncoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	r, size := rune(0), 0
	for ; nSrc < len(src); nSrc += size ***REMOVED***
		r = rune(src[nSrc])

		// Decode a 1-byte rune.
		if r < utf8.RuneSelf ***REMOVED***
			size = 1

		***REMOVED*** else ***REMOVED***
			// Decode a multi-byte rune.
			r, size = utf8.DecodeRune(src[nSrc:])
			if size == 1 ***REMOVED***
				// All valid runes of size 1 (those below utf8.RuneSelf) were
				// handled above. We have invalid UTF-8 or we haven't seen the
				// full character yet.
				if !atEOF && !utf8.FullRune(src[nSrc:]) ***REMOVED***
					err = transform.ErrShortSrc
					break
				***REMOVED***
			***REMOVED***

			// func init checks that the switch covers all tables.
			switch ***REMOVED***
			case encode0Low <= r && r < encode0High:
				if r = rune(encode0[r-encode0Low]); r != 0 ***REMOVED***
					goto write2or3
				***REMOVED***
			case encode1Low <= r && r < encode1High:
				if r = rune(encode1[r-encode1Low]); r != 0 ***REMOVED***
					goto write2or3
				***REMOVED***
			case encode2Low <= r && r < encode2High:
				if r = rune(encode2[r-encode2Low]); r != 0 ***REMOVED***
					goto write2or3
				***REMOVED***
			case encode3Low <= r && r < encode3High:
				if r = rune(encode3[r-encode3Low]); r != 0 ***REMOVED***
					goto write2or3
				***REMOVED***
			case encode4Low <= r && r < encode4High:
				if r = rune(encode4[r-encode4Low]); r != 0 ***REMOVED***
					goto write2or3
				***REMOVED***
			case encode5Low <= r && r < encode5High:
				if 0xff61 <= r && r < 0xffa0 ***REMOVED***
					goto write2
				***REMOVED***
				if r = rune(encode5[r-encode5Low]); r != 0 ***REMOVED***
					goto write2or3
				***REMOVED***
			***REMOVED***
			err = internal.ErrASCIIReplacement
			break
		***REMOVED***

		if nDst >= len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break
		***REMOVED***
		dst[nDst] = uint8(r)
		nDst++
		continue

	write2or3:
		if r>>tableShift == jis0208 ***REMOVED***
			if nDst+2 > len(dst) ***REMOVED***
				err = transform.ErrShortDst
				break
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if nDst+3 > len(dst) ***REMOVED***
				err = transform.ErrShortDst
				break
			***REMOVED***
			dst[nDst] = 0x8f
			nDst++
		***REMOVED***
		dst[nDst+0] = 0xa1 + uint8(r>>codeShift)&codeMask
		dst[nDst+1] = 0xa1 + uint8(r)&codeMask
		nDst += 2
		continue

	write2:
		if nDst+2 > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break
		***REMOVED***
		dst[nDst+0] = 0x8e
		dst[nDst+1] = uint8(r - (0xff61 - 0xa1))
		nDst += 2
		continue
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***

func init() ***REMOVED***
	// Check that the hard-coded encode switch covers all tables.
	if numEncodeTables != 6 ***REMOVED***
		panic("bad numEncodeTables")
	***REMOVED***
***REMOVED***
