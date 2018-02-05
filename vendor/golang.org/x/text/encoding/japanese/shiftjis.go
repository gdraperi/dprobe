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

// ShiftJIS is the Shift JIS encoding, also known as Code Page 932 and
// Windows-31J.
var ShiftJIS encoding.Encoding = &shiftJIS

var shiftJIS = internal.Encoding***REMOVED***
	&internal.SimpleEncoding***REMOVED***shiftJISDecoder***REMOVED******REMOVED***, shiftJISEncoder***REMOVED******REMOVED******REMOVED***,
	"Shift JIS",
	identifier.ShiftJIS,
***REMOVED***

type shiftJISDecoder struct***REMOVED*** transform.NopResetter ***REMOVED***

func (shiftJISDecoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	r, size := rune(0), 0
loop:
	for ; nSrc < len(src); nSrc += size ***REMOVED***
		switch c0 := src[nSrc]; ***REMOVED***
		case c0 < utf8.RuneSelf:
			r, size = rune(c0), 1

		case 0xa1 <= c0 && c0 < 0xe0:
			r, size = rune(c0)+(0xff61-0xa1), 1

		case (0x81 <= c0 && c0 < 0xa0) || (0xe0 <= c0 && c0 < 0xfd):
			if c0 <= 0x9f ***REMOVED***
				c0 -= 0x70
			***REMOVED*** else ***REMOVED***
				c0 -= 0xb0
			***REMOVED***
			c0 = 2*c0 - 0x21

			if nSrc+1 >= len(src) ***REMOVED***
				if !atEOF ***REMOVED***
					err = transform.ErrShortSrc
					break loop
				***REMOVED***
				r, size = '\ufffd', 1
				goto write
			***REMOVED***
			c1 := src[nSrc+1]
			switch ***REMOVED***
			case c1 < 0x40:
				r, size = '\ufffd', 1 // c1 is ASCII so output on next round
				goto write
			case c1 < 0x7f:
				c0--
				c1 -= 0x40
			case c1 == 0x7f:
				r, size = '\ufffd', 1 // c1 is ASCII so output on next round
				goto write
			case c1 < 0x9f:
				c0--
				c1 -= 0x41
			case c1 < 0xfd:
				c1 -= 0x9f
			default:
				r, size = '\ufffd', 2
				goto write
			***REMOVED***
			r, size = '\ufffd', 2
			if i := int(c0)*94 + int(c1); i < len(jis0208Decode) ***REMOVED***
				r = rune(jis0208Decode[i])
				if r == 0 ***REMOVED***
					r = '\ufffd'
				***REMOVED***
			***REMOVED***

		case c0 == 0x80:
			r, size = 0x80, 1

		default:
			r, size = '\ufffd', 1
		***REMOVED***
	write:
		if nDst+utf8.RuneLen(r) > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break loop
		***REMOVED***
		nDst += utf8.EncodeRune(dst[nDst:], r)
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***

type shiftJISEncoder struct***REMOVED*** transform.NopResetter ***REMOVED***

func (shiftJISEncoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	r, size := rune(0), 0
loop:
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
					break loop
				***REMOVED***
			***REMOVED***

			// func init checks that the switch covers all tables.
			switch ***REMOVED***
			case encode0Low <= r && r < encode0High:
				if r = rune(encode0[r-encode0Low]); r>>tableShift == jis0208 ***REMOVED***
					goto write2
				***REMOVED***
			case encode1Low <= r && r < encode1High:
				if r = rune(encode1[r-encode1Low]); r>>tableShift == jis0208 ***REMOVED***
					goto write2
				***REMOVED***
			case encode2Low <= r && r < encode2High:
				if r = rune(encode2[r-encode2Low]); r>>tableShift == jis0208 ***REMOVED***
					goto write2
				***REMOVED***
			case encode3Low <= r && r < encode3High:
				if r = rune(encode3[r-encode3Low]); r>>tableShift == jis0208 ***REMOVED***
					goto write2
				***REMOVED***
			case encode4Low <= r && r < encode4High:
				if r = rune(encode4[r-encode4Low]); r>>tableShift == jis0208 ***REMOVED***
					goto write2
				***REMOVED***
			case encode5Low <= r && r < encode5High:
				if 0xff61 <= r && r < 0xffa0 ***REMOVED***
					r -= 0xff61 - 0xa1
					goto write1
				***REMOVED***
				if r = rune(encode5[r-encode5Low]); r>>tableShift == jis0208 ***REMOVED***
					goto write2
				***REMOVED***
			***REMOVED***
			err = internal.ErrASCIIReplacement
			break
		***REMOVED***

	write1:
		if nDst >= len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break
		***REMOVED***
		dst[nDst] = uint8(r)
		nDst++
		continue

	write2:
		j1 := uint8(r>>codeShift) & codeMask
		j2 := uint8(r) & codeMask
		if nDst+2 > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break loop
		***REMOVED***
		if j1 <= 61 ***REMOVED***
			dst[nDst+0] = 129 + j1/2
		***REMOVED*** else ***REMOVED***
			dst[nDst+0] = 193 + j1/2
		***REMOVED***
		if j1&1 == 0 ***REMOVED***
			dst[nDst+1] = j2 + j2/63 + 64
		***REMOVED*** else ***REMOVED***
			dst[nDst+1] = j2 + 159
		***REMOVED***
		nDst += 2
		continue
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***
