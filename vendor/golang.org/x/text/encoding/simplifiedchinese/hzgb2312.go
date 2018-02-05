// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simplifiedchinese

import (
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/transform"
)

// HZGB2312 is the HZ-GB2312 encoding.
var HZGB2312 encoding.Encoding = &hzGB2312

var hzGB2312 = internal.Encoding***REMOVED***
	internal.FuncEncoding***REMOVED***hzGB2312NewDecoder, hzGB2312NewEncoder***REMOVED***,
	"HZ-GB2312",
	identifier.HZGB2312,
***REMOVED***

func hzGB2312NewDecoder() transform.Transformer ***REMOVED***
	return new(hzGB2312Decoder)
***REMOVED***

func hzGB2312NewEncoder() transform.Transformer ***REMOVED***
	return new(hzGB2312Encoder)
***REMOVED***

const (
	asciiState = iota
	gbState
)

type hzGB2312Decoder int

func (d *hzGB2312Decoder) Reset() ***REMOVED***
	*d = asciiState
***REMOVED***

func (d *hzGB2312Decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	r, size := rune(0), 0
loop:
	for ; nSrc < len(src); nSrc += size ***REMOVED***
		c0 := src[nSrc]
		if c0 >= utf8.RuneSelf ***REMOVED***
			r, size = utf8.RuneError, 1
			goto write
		***REMOVED***

		if c0 == '~' ***REMOVED***
			if nSrc+1 >= len(src) ***REMOVED***
				if !atEOF ***REMOVED***
					err = transform.ErrShortSrc
					break loop
				***REMOVED***
				r = utf8.RuneError
				goto write
			***REMOVED***
			size = 2
			switch src[nSrc+1] ***REMOVED***
			case '***REMOVED***':
				*d = gbState
				continue
			case '***REMOVED***':
				*d = asciiState
				continue
			case '~':
				if nDst >= len(dst) ***REMOVED***
					err = transform.ErrShortDst
					break loop
				***REMOVED***
				dst[nDst] = '~'
				nDst++
				continue
			case '\n':
				continue
			default:
				r = utf8.RuneError
				goto write
			***REMOVED***
		***REMOVED***

		if *d == asciiState ***REMOVED***
			r, size = rune(c0), 1
		***REMOVED*** else ***REMOVED***
			if nSrc+1 >= len(src) ***REMOVED***
				if !atEOF ***REMOVED***
					err = transform.ErrShortSrc
					break loop
				***REMOVED***
				r, size = utf8.RuneError, 1
				goto write
			***REMOVED***
			size = 2
			c1 := src[nSrc+1]
			if c0 < 0x21 || 0x7e <= c0 || c1 < 0x21 || 0x7f <= c1 ***REMOVED***
				// error
			***REMOVED*** else if i := int(c0-0x01)*190 + int(c1+0x3f); i < len(decode) ***REMOVED***
				r = rune(decode[i])
				if r != 0 ***REMOVED***
					goto write
				***REMOVED***
			***REMOVED***
			if c1 > utf8.RuneSelf ***REMOVED***
				// Be consistent and always treat non-ASCII as a single error.
				size = 1
			***REMOVED***
			r = utf8.RuneError
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

type hzGB2312Encoder int

func (d *hzGB2312Encoder) Reset() ***REMOVED***
	*d = asciiState
***REMOVED***

func (e *hzGB2312Encoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	r, size := rune(0), 0
	for ; nSrc < len(src); nSrc += size ***REMOVED***
		r = rune(src[nSrc])

		// Decode a 1-byte rune.
		if r < utf8.RuneSelf ***REMOVED***
			size = 1
			if r == '~' ***REMOVED***
				if nDst+2 > len(dst) ***REMOVED***
					err = transform.ErrShortDst
					break
				***REMOVED***
				dst[nDst+0] = '~'
				dst[nDst+1] = '~'
				nDst += 2
				continue
			***REMOVED*** else if *e != asciiState ***REMOVED***
				if nDst+3 > len(dst) ***REMOVED***
					err = transform.ErrShortDst
					break
				***REMOVED***
				*e = asciiState
				dst[nDst+0] = '~'
				dst[nDst+1] = '***REMOVED***'
				nDst += 2
			***REMOVED*** else if nDst >= len(dst) ***REMOVED***
				err = transform.ErrShortDst
				break
			***REMOVED***
			dst[nDst] = uint8(r)
			nDst += 1
			continue

		***REMOVED***

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
				goto writeGB
			***REMOVED***
		case encode1Low <= r && r < encode1High:
			if r = rune(encode1[r-encode1Low]); r != 0 ***REMOVED***
				goto writeGB
			***REMOVED***
		case encode2Low <= r && r < encode2High:
			if r = rune(encode2[r-encode2Low]); r != 0 ***REMOVED***
				goto writeGB
			***REMOVED***
		case encode3Low <= r && r < encode3High:
			if r = rune(encode3[r-encode3Low]); r != 0 ***REMOVED***
				goto writeGB
			***REMOVED***
		case encode4Low <= r && r < encode4High:
			if r = rune(encode4[r-encode4Low]); r != 0 ***REMOVED***
				goto writeGB
			***REMOVED***
		***REMOVED***

	terminateInASCIIState:
		// Switch back to ASCII state in case of error so that an ASCII
		// replacement character can be written in the correct state.
		if *e != asciiState ***REMOVED***
			if nDst+2 > len(dst) ***REMOVED***
				err = transform.ErrShortDst
				break
			***REMOVED***
			dst[nDst+0] = '~'
			dst[nDst+1] = '***REMOVED***'
			nDst += 2
		***REMOVED***
		err = internal.ErrASCIIReplacement
		break

	writeGB:
		c0 := uint8(r>>8) - 0x80
		c1 := uint8(r) - 0x80
		if c0 < 0x21 || 0x7e <= c0 || c1 < 0x21 || 0x7f <= c1 ***REMOVED***
			goto terminateInASCIIState
		***REMOVED***
		if *e == asciiState ***REMOVED***
			if nDst+4 > len(dst) ***REMOVED***
				err = transform.ErrShortDst
				break
			***REMOVED***
			*e = gbState
			dst[nDst+0] = '~'
			dst[nDst+1] = '***REMOVED***'
			nDst += 2
		***REMOVED*** else if nDst+2 > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break
		***REMOVED***
		dst[nDst+0] = c0
		dst[nDst+1] = c1
		nDst += 2
		continue
	***REMOVED***
	// TODO: should one always terminate in ASCII state to make it safe to
	// concatenate two HZ-GB2312-encoded strings?
	return nDst, nSrc, err
***REMOVED***
