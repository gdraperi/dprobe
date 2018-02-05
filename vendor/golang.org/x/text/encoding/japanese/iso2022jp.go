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

// ISO2022JP is the ISO-2022-JP encoding.
var ISO2022JP encoding.Encoding = &iso2022JP

var iso2022JP = internal.Encoding***REMOVED***
	internal.FuncEncoding***REMOVED***iso2022JPNewDecoder, iso2022JPNewEncoder***REMOVED***,
	"ISO-2022-JP",
	identifier.ISO2022JP,
***REMOVED***

func iso2022JPNewDecoder() transform.Transformer ***REMOVED***
	return new(iso2022JPDecoder)
***REMOVED***

func iso2022JPNewEncoder() transform.Transformer ***REMOVED***
	return new(iso2022JPEncoder)
***REMOVED***

const (
	asciiState = iota
	katakanaState
	jis0208State
	jis0212State
)

const asciiEsc = 0x1b

type iso2022JPDecoder int

func (d *iso2022JPDecoder) Reset() ***REMOVED***
	*d = asciiState
***REMOVED***

func (d *iso2022JPDecoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	r, size := rune(0), 0
	for ; nSrc < len(src); nSrc += size ***REMOVED***
		c0 := src[nSrc]
		if c0 >= utf8.RuneSelf ***REMOVED***
			r, size = '\ufffd', 1
			goto write
		***REMOVED***

		if c0 == asciiEsc ***REMOVED***
			if nSrc+2 >= len(src) ***REMOVED***
				if !atEOF ***REMOVED***
					return nDst, nSrc, transform.ErrShortSrc
				***REMOVED***
				// TODO: is it correct to only skip 1??
				r, size = '\ufffd', 1
				goto write
			***REMOVED***
			size = 3
			c1 := src[nSrc+1]
			c2 := src[nSrc+2]
			switch ***REMOVED***
			case c1 == '$' && (c2 == '@' || c2 == 'B'): // 0x24 ***REMOVED***0x40, 0x42***REMOVED***
				*d = jis0208State
				continue
			case c1 == '$' && c2 == '(': // 0x24 0x28
				if nSrc+3 >= len(src) ***REMOVED***
					if !atEOF ***REMOVED***
						return nDst, nSrc, transform.ErrShortSrc
					***REMOVED***
					r, size = '\ufffd', 1
					goto write
				***REMOVED***
				size = 4
				if src[nSrc+3] == 'D' ***REMOVED***
					*d = jis0212State
					continue
				***REMOVED***
			case c1 == '(' && (c2 == 'B' || c2 == 'J'): // 0x28 ***REMOVED***0x42, 0x4A***REMOVED***
				*d = asciiState
				continue
			case c1 == '(' && c2 == 'I': // 0x28 0x49
				*d = katakanaState
				continue
			***REMOVED***
			r, size = '\ufffd', 1
			goto write
		***REMOVED***

		switch *d ***REMOVED***
		case asciiState:
			r, size = rune(c0), 1

		case katakanaState:
			if c0 < 0x21 || 0x60 <= c0 ***REMOVED***
				r, size = '\ufffd', 1
				goto write
			***REMOVED***
			r, size = rune(c0)+(0xff61-0x21), 1

		default:
			if c0 == 0x0a ***REMOVED***
				*d = asciiState
				r, size = rune(c0), 1
				goto write
			***REMOVED***
			if nSrc+1 >= len(src) ***REMOVED***
				if !atEOF ***REMOVED***
					return nDst, nSrc, transform.ErrShortSrc
				***REMOVED***
				r, size = '\ufffd', 1
				goto write
			***REMOVED***
			size = 2
			c1 := src[nSrc+1]
			i := int(c0-0x21)*94 + int(c1-0x21)
			if *d == jis0208State && i < len(jis0208Decode) ***REMOVED***
				r = rune(jis0208Decode[i])
			***REMOVED*** else if *d == jis0212State && i < len(jis0212Decode) ***REMOVED***
				r = rune(jis0212Decode[i])
			***REMOVED*** else ***REMOVED***
				r = '\ufffd'
				goto write
			***REMOVED***
			if r == 0 ***REMOVED***
				r = '\ufffd'
			***REMOVED***
		***REMOVED***

	write:
		if nDst+utf8.RuneLen(r) > len(dst) ***REMOVED***
			return nDst, nSrc, transform.ErrShortDst
		***REMOVED***
		nDst += utf8.EncodeRune(dst[nDst:], r)
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***

type iso2022JPEncoder int

func (e *iso2022JPEncoder) Reset() ***REMOVED***
	*e = asciiState
***REMOVED***

func (e *iso2022JPEncoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
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
			//
			// http://encoding.spec.whatwg.org/#iso-2022-jp says that "the index jis0212
			// is not used by the iso-2022-jp encoder due to lack of widespread support".
			//
			// TODO: do we have to special-case U+00A5 and U+203E, as per
			// http://encoding.spec.whatwg.org/#iso-2022-jp
			// Doing so would mean that "\u00a5" would not be preserved
			// after an encode-decode round trip.
			switch ***REMOVED***
			case encode0Low <= r && r < encode0High:
				if r = rune(encode0[r-encode0Low]); r>>tableShift == jis0208 ***REMOVED***
					goto writeJIS
				***REMOVED***
			case encode1Low <= r && r < encode1High:
				if r = rune(encode1[r-encode1Low]); r>>tableShift == jis0208 ***REMOVED***
					goto writeJIS
				***REMOVED***
			case encode2Low <= r && r < encode2High:
				if r = rune(encode2[r-encode2Low]); r>>tableShift == jis0208 ***REMOVED***
					goto writeJIS
				***REMOVED***
			case encode3Low <= r && r < encode3High:
				if r = rune(encode3[r-encode3Low]); r>>tableShift == jis0208 ***REMOVED***
					goto writeJIS
				***REMOVED***
			case encode4Low <= r && r < encode4High:
				if r = rune(encode4[r-encode4Low]); r>>tableShift == jis0208 ***REMOVED***
					goto writeJIS
				***REMOVED***
			case encode5Low <= r && r < encode5High:
				if 0xff61 <= r && r < 0xffa0 ***REMOVED***
					goto writeKatakana
				***REMOVED***
				if r = rune(encode5[r-encode5Low]); r>>tableShift == jis0208 ***REMOVED***
					goto writeJIS
				***REMOVED***
			***REMOVED***

			// Switch back to ASCII state in case of error so that an ASCII
			// replacement character can be written in the correct state.
			if *e != asciiState ***REMOVED***
				if nDst+3 > len(dst) ***REMOVED***
					err = transform.ErrShortDst
					break
				***REMOVED***
				*e = asciiState
				dst[nDst+0] = asciiEsc
				dst[nDst+1] = '('
				dst[nDst+2] = 'B'
				nDst += 3
			***REMOVED***
			err = internal.ErrASCIIReplacement
			break
		***REMOVED***

		if *e != asciiState ***REMOVED***
			if nDst+4 > len(dst) ***REMOVED***
				err = transform.ErrShortDst
				break
			***REMOVED***
			*e = asciiState
			dst[nDst+0] = asciiEsc
			dst[nDst+1] = '('
			dst[nDst+2] = 'B'
			nDst += 3
		***REMOVED*** else if nDst >= len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break
		***REMOVED***
		dst[nDst] = uint8(r)
		nDst++
		continue

	writeJIS:
		if *e != jis0208State ***REMOVED***
			if nDst+5 > len(dst) ***REMOVED***
				err = transform.ErrShortDst
				break
			***REMOVED***
			*e = jis0208State
			dst[nDst+0] = asciiEsc
			dst[nDst+1] = '$'
			dst[nDst+2] = 'B'
			nDst += 3
		***REMOVED*** else if nDst+2 > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break
		***REMOVED***
		dst[nDst+0] = 0x21 + uint8(r>>codeShift)&codeMask
		dst[nDst+1] = 0x21 + uint8(r)&codeMask
		nDst += 2
		continue

	writeKatakana:
		if *e != katakanaState ***REMOVED***
			if nDst+4 > len(dst) ***REMOVED***
				err = transform.ErrShortDst
				break
			***REMOVED***
			*e = katakanaState
			dst[nDst+0] = asciiEsc
			dst[nDst+1] = '('
			dst[nDst+2] = 'I'
			nDst += 3
		***REMOVED*** else if nDst >= len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break
		***REMOVED***
		dst[nDst] = uint8(r - (0xff61 - 0x21))
		nDst++
		continue
	***REMOVED***
	if atEOF && err == nil && *e != asciiState ***REMOVED***
		if nDst+3 > len(dst) ***REMOVED***
			err = transform.ErrShortDst
		***REMOVED*** else ***REMOVED***
			*e = asciiState
			dst[nDst+0] = asciiEsc
			dst[nDst+1] = '('
			dst[nDst+2] = 'B'
			nDst += 3
		***REMOVED***
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***
