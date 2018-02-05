// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traditionalchinese

import (
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/transform"
)

// All is a list of all defined encodings in this package.
var All = []encoding.Encoding***REMOVED***Big5***REMOVED***

// Big5 is the Big5 encoding, also known as Code Page 950.
var Big5 encoding.Encoding = &big5

var big5 = internal.Encoding***REMOVED***
	&internal.SimpleEncoding***REMOVED***big5Decoder***REMOVED******REMOVED***, big5Encoder***REMOVED******REMOVED******REMOVED***,
	"Big5",
	identifier.Big5,
***REMOVED***

type big5Decoder struct***REMOVED*** transform.NopResetter ***REMOVED***

func (big5Decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	r, size, s := rune(0), 0, ""
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
				goto write
			***REMOVED***
			c1 := src[nSrc+1]
			switch ***REMOVED***
			case 0x40 <= c1 && c1 < 0x7f:
				c1 -= 0x40
			case 0xa1 <= c1 && c1 < 0xff:
				c1 -= 0x62
			case c1 < 0x40:
				r, size = utf8.RuneError, 1
				goto write
			default:
				r, size = utf8.RuneError, 2
				goto write
			***REMOVED***
			r, size = '\ufffd', 2
			if i := int(c0-0x81)*157 + int(c1); i < len(decode) ***REMOVED***
				if 1133 <= i && i < 1167 ***REMOVED***
					// The two-rune special cases for LATIN CAPITAL / SMALL E WITH CIRCUMFLEX
					// AND MACRON / CARON are from http://encoding.spec.whatwg.org/#big5
					switch i ***REMOVED***
					case 1133:
						s = "\u00CA\u0304"
						goto writeStr
					case 1135:
						s = "\u00CA\u030C"
						goto writeStr
					case 1164:
						s = "\u00EA\u0304"
						goto writeStr
					case 1166:
						s = "\u00EA\u030C"
						goto writeStr
					***REMOVED***
				***REMOVED***
				r = rune(decode[i])
				if r == 0 ***REMOVED***
					r = '\ufffd'
				***REMOVED***
			***REMOVED***

		default:
			r, size = utf8.RuneError, 1
		***REMOVED***

	write:
		if nDst+utf8.RuneLen(r) > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break loop
		***REMOVED***
		nDst += utf8.EncodeRune(dst[nDst:], r)
		continue loop

	writeStr:
		if nDst+len(s) > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break loop
		***REMOVED***
		nDst += copy(dst[nDst:], s)
		continue loop
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***

type big5Encoder struct***REMOVED*** transform.NopResetter ***REMOVED***

func (big5Encoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
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
		***REMOVED***

		if r >= utf8.RuneSelf ***REMOVED***
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
			case encode7Low <= r && r < encode7High:
				if r = rune(encode7[r-encode7Low]); r != 0 ***REMOVED***
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
	if numEncodeTables != 8 ***REMOVED***
		panic("bad numEncodeTables")
	***REMOVED***
***REMOVED***