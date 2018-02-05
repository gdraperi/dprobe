// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package korean

import (
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/transform"
)

// All is a list of all defined encodings in this package.
var All = []encoding.Encoding***REMOVED***EUCKR***REMOVED***

// EUCKR is the EUC-KR encoding, also known as Code Page 949.
var EUCKR encoding.Encoding = &eucKR

var eucKR = internal.Encoding***REMOVED***
	&internal.SimpleEncoding***REMOVED***eucKRDecoder***REMOVED******REMOVED***, eucKREncoder***REMOVED******REMOVED******REMOVED***,
	"EUC-KR",
	identifier.EUCKR,
***REMOVED***

type eucKRDecoder struct***REMOVED*** transform.NopResetter ***REMOVED***

func (eucKRDecoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	r, size := rune(0), 0
loop:
	for ; nSrc < len(src); nSrc += size ***REMOVED***
		switch c0 := src[nSrc]; ***REMOVED***
		case c0 < utf8.RuneSelf:
			r, size = rune(c0), 1

		case 0x81 <= c0 && c0 < 0xff:
			if nSrc+1 >= len(src) ***REMOVED***
				if !atEOF ***REMOVED***
					err = transform.ErrShortSrc
					break loop
				***REMOVED***
				r, size = utf8.RuneError, 1
				break
			***REMOVED***
			c1 := src[nSrc+1]
			size = 2
			if c0 < 0xc7 ***REMOVED***
				r = 178 * rune(c0-0x81)
				switch ***REMOVED***
				case 0x41 <= c1 && c1 < 0x5b:
					r += rune(c1) - (0x41 - 0*26)
				case 0x61 <= c1 && c1 < 0x7b:
					r += rune(c1) - (0x61 - 1*26)
				case 0x81 <= c1 && c1 < 0xff:
					r += rune(c1) - (0x81 - 2*26)
				default:
					goto decError
				***REMOVED***
			***REMOVED*** else if 0xa1 <= c1 && c1 < 0xff ***REMOVED***
				r = 178*(0xc7-0x81) + rune(c0-0xc7)*94 + rune(c1-0xa1)
			***REMOVED*** else ***REMOVED***
				goto decError
			***REMOVED***
			if int(r) < len(decode) ***REMOVED***
				r = rune(decode[r])
				if r != 0 ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
		decError:
			r = utf8.RuneError
			if c1 < utf8.RuneSelf ***REMOVED***
				size = 1
			***REMOVED***

		default:
			r, size = utf8.RuneError, 1
			break
		***REMOVED***

		if nDst+utf8.RuneLen(r) > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break
		***REMOVED***
		nDst += utf8.EncodeRune(dst[nDst:], r)
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***

type eucKREncoder struct***REMOVED*** transform.NopResetter ***REMOVED***

func (eucKREncoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	r, size := rune(0), 0
	for ; nSrc < len(src); nSrc += size ***REMOVED***
		r = rune(src[nSrc])

		// Decode a 1-byte rune.
		if r < utf8.RuneSelf ***REMOVED***
			size = 1

			if nDst >= len(dst) ***REMOVED***
				err = transform.ErrShortDst
				break
			***REMOVED***
			dst[nDst] = uint8(r)
			nDst++
			continue

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
					goto write2
				***REMOVED***
			case encode1Low <= r && r < encode1High:
				if r = rune(encode1[r-encode1Low]); r != 0 ***REMOVED***
					goto write2
				***REMOVED***
			case encode2Low <= r && r < encode2High:
				if r = rune(encode2[r-encode2Low]); r != 0 ***REMOVED***
					goto write2
				***REMOVED***
			case encode3Low <= r && r < encode3High:
				if r = rune(encode3[r-encode3Low]); r != 0 ***REMOVED***
					goto write2
				***REMOVED***
			case encode4Low <= r && r < encode4High:
				if r = rune(encode4[r-encode4Low]); r != 0 ***REMOVED***
					goto write2
				***REMOVED***
			case encode5Low <= r && r < encode5High:
				if r = rune(encode5[r-encode5Low]); r != 0 ***REMOVED***
					goto write2
				***REMOVED***
			case encode6Low <= r && r < encode6High:
				if r = rune(encode6[r-encode6Low]); r != 0 ***REMOVED***
					goto write2
				***REMOVED***
			***REMOVED***
			err = internal.ErrASCIIReplacement
			break
		***REMOVED***

	write2:
		if nDst+2 > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break
		***REMOVED***
		dst[nDst+0] = uint8(r >> 8)
		dst[nDst+1] = uint8(r)
		nDst += 2
		continue
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***

func init() ***REMOVED***
	// Check that the hard-coded encode switch covers all tables.
	if numEncodeTables != 7 ***REMOVED***
		panic("bad numEncodeTables")
	***REMOVED***
***REMOVED***
