// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package unicode provides Unicode encodings such as UTF-16.
package unicode // import "golang.org/x/text/encoding/unicode"

import (
	"errors"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/internal/utf8internal"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

// TODO: I think the Transformers really should return errors on unmatched
// surrogate pairs and odd numbers of bytes. This is not required by RFC 2781,
// which leaves it open, but is suggested by WhatWG. It will allow for all error
// modes as defined by WhatWG: fatal, HTML and Replacement. This would require
// the introduction of some kind of error type for conveying the erroneous code
// point.

// UTF8 is the UTF-8 encoding.
var UTF8 encoding.Encoding = utf8enc

var utf8enc = &internal.Encoding***REMOVED***
	&internal.SimpleEncoding***REMOVED***utf8Decoder***REMOVED******REMOVED***, runes.ReplaceIllFormed()***REMOVED***,
	"UTF-8",
	identifier.UTF8,
***REMOVED***

type utf8Decoder struct***REMOVED*** transform.NopResetter ***REMOVED***

func (utf8Decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	var pSrc int // point from which to start copy in src
	var accept utf8internal.AcceptRange

	// The decoder can only make the input larger, not smaller.
	n := len(src)
	if len(dst) < n ***REMOVED***
		err = transform.ErrShortDst
		n = len(dst)
		atEOF = false
	***REMOVED***
	for nSrc < n ***REMOVED***
		c := src[nSrc]
		if c < utf8.RuneSelf ***REMOVED***
			nSrc++
			continue
		***REMOVED***
		first := utf8internal.First[c]
		size := int(first & utf8internal.SizeMask)
		if first == utf8internal.FirstInvalid ***REMOVED***
			goto handleInvalid // invalid starter byte
		***REMOVED***
		accept = utf8internal.AcceptRanges[first>>utf8internal.AcceptShift]
		if nSrc+size > n ***REMOVED***
			if !atEOF ***REMOVED***
				// We may stop earlier than necessary here if the short sequence
				// has invalid bytes. Not checking for this simplifies the code
				// and may avoid duplicate computations in certain conditions.
				if err == nil ***REMOVED***
					err = transform.ErrShortSrc
				***REMOVED***
				break
			***REMOVED***
			// Determine the maximal subpart of an ill-formed subsequence.
			switch ***REMOVED***
			case nSrc+1 >= n || src[nSrc+1] < accept.Lo || accept.Hi < src[nSrc+1]:
				size = 1
			case nSrc+2 >= n || src[nSrc+2] < utf8internal.LoCB || utf8internal.HiCB < src[nSrc+2]:
				size = 2
			default:
				size = 3 // As we are short, the maximum is 3.
			***REMOVED***
			goto handleInvalid
		***REMOVED***
		if c = src[nSrc+1]; c < accept.Lo || accept.Hi < c ***REMOVED***
			size = 1
			goto handleInvalid // invalid continuation byte
		***REMOVED*** else if size == 2 ***REMOVED***
		***REMOVED*** else if c = src[nSrc+2]; c < utf8internal.LoCB || utf8internal.HiCB < c ***REMOVED***
			size = 2
			goto handleInvalid // invalid continuation byte
		***REMOVED*** else if size == 3 ***REMOVED***
		***REMOVED*** else if c = src[nSrc+3]; c < utf8internal.LoCB || utf8internal.HiCB < c ***REMOVED***
			size = 3
			goto handleInvalid // invalid continuation byte
		***REMOVED***
		nSrc += size
		continue

	handleInvalid:
		// Copy the scanned input so far.
		nDst += copy(dst[nDst:], src[pSrc:nSrc])

		// Append RuneError to the destination.
		const runeError = "\ufffd"
		if nDst+len(runeError) > len(dst) ***REMOVED***
			return nDst, nSrc, transform.ErrShortDst
		***REMOVED***
		nDst += copy(dst[nDst:], runeError)

		// Skip the maximal subpart of an ill-formed subsequence according to
		// the W3C standard way instead of the Go way. This Transform is
		// probably the only place in the text repo where it is warranted.
		nSrc += size
		pSrc = nSrc

		// Recompute the maximum source length.
		if sz := len(dst) - nDst; sz < len(src)-nSrc ***REMOVED***
			err = transform.ErrShortDst
			n = nSrc + sz
			atEOF = false
		***REMOVED***
	***REMOVED***
	return nDst + copy(dst[nDst:], src[pSrc:nSrc]), nSrc, err
***REMOVED***

// UTF16 returns a UTF-16 Encoding for the given default endianness and byte
// order mark (BOM) policy.
//
// When decoding from UTF-16 to UTF-8, if the BOMPolicy is IgnoreBOM then
// neither BOMs U+FEFF nor noncharacters U+FFFE in the input stream will affect
// the endianness used for decoding, and will instead be output as their
// standard UTF-8 encodings: "\xef\xbb\xbf" and "\xef\xbf\xbe". If the BOMPolicy
// is UseBOM or ExpectBOM a staring BOM is not written to the UTF-8 output.
// Instead, it overrides the default endianness e for the remainder of the
// transformation. Any subsequent BOMs U+FEFF or noncharacters U+FFFE will not
// affect the endianness used, and will instead be output as their standard
// UTF-8 encodings. For UseBOM, if there is no starting BOM, it will proceed
// with the default Endianness. For ExpectBOM, in that case, the transformation
// will return early with an ErrMissingBOM error.
//
// When encoding from UTF-8 to UTF-16, a BOM will be inserted at the start of
// the output if the BOMPolicy is UseBOM or ExpectBOM. Otherwise, a BOM will not
// be inserted. The UTF-8 input does not need to contain a BOM.
//
// There is no concept of a 'native' endianness. If the UTF-16 data is produced
// and consumed in a greater context that implies a certain endianness, use
// IgnoreBOM. Otherwise, use ExpectBOM and always produce and consume a BOM.
//
// In the language of http://www.unicode.org/faq/utf_bom.html#bom10, IgnoreBOM
// corresponds to "Where the precise type of the data stream is known... the
// BOM should not be used" and ExpectBOM corresponds to "A particular
// protocol... may require use of the BOM".
func UTF16(e Endianness, b BOMPolicy) encoding.Encoding ***REMOVED***
	return utf16Encoding***REMOVED***config***REMOVED***e, b***REMOVED***, mibValue[e][b&bomMask]***REMOVED***
***REMOVED***

// mibValue maps Endianness and BOMPolicy settings to MIB constants. Note that
// some configurations map to the same MIB identifier. RFC 2781 has requirements
// and recommendations. Some of the "configurations" are merely recommendations,
// so multiple configurations could match.
var mibValue = map[Endianness][numBOMValues]identifier.MIB***REMOVED***
	BigEndian: [numBOMValues]identifier.MIB***REMOVED***
		IgnoreBOM: identifier.UTF16BE,
		UseBOM:    identifier.UTF16, // BigEnding default is preferred by RFC 2781.
		// TODO: acceptBOM | strictBOM would map to UTF16BE as well.
	***REMOVED***,
	LittleEndian: [numBOMValues]identifier.MIB***REMOVED***
		IgnoreBOM: identifier.UTF16LE,
		UseBOM:    identifier.UTF16, // LittleEndian default is allowed and preferred on Windows.
		// TODO: acceptBOM | strictBOM would map to UTF16LE as well.
	***REMOVED***,
	// ExpectBOM is not widely used and has no valid MIB identifier.
***REMOVED***

// All lists a configuration for each IANA-defined UTF-16 variant.
var All = []encoding.Encoding***REMOVED***
	UTF8,
	UTF16(BigEndian, UseBOM),
	UTF16(BigEndian, IgnoreBOM),
	UTF16(LittleEndian, IgnoreBOM),
***REMOVED***

// BOMPolicy is a UTF-16 encoding's byte order mark policy.
type BOMPolicy uint8

const (
	writeBOM   BOMPolicy = 0x01
	acceptBOM  BOMPolicy = 0x02
	requireBOM BOMPolicy = 0x04
	bomMask    BOMPolicy = 0x07

	// HACK: numBOMValues == 8 triggers a bug in the 1.4 compiler (cannot have a
	// map of an array of length 8 of a type that is also used as a key or value
	// in another map). See golang.org/issue/11354.
	// TODO: consider changing this value back to 8 if the use of 1.4.* has
	// been minimized.
	numBOMValues = 8 + 1

	// IgnoreBOM means to ignore any byte order marks.
	IgnoreBOM BOMPolicy = 0
	// Common and RFC 2781-compliant interpretation for UTF-16BE/LE.

	// UseBOM means that the UTF-16 form may start with a byte order mark, which
	// will be used to override the default encoding.
	UseBOM BOMPolicy = writeBOM | acceptBOM
	// Common and RFC 2781-compliant interpretation for UTF-16.

	// ExpectBOM means that the UTF-16 form must start with a byte order mark,
	// which will be used to override the default encoding.
	ExpectBOM BOMPolicy = writeBOM | acceptBOM | requireBOM
	// Used in Java as Unicode (not to be confused with Java's UTF-16) and
	// ICU's UTF-16,version=1. Not compliant with RFC 2781.

	// TODO (maybe): strictBOM: BOM must match Endianness. This would allow:
	// - UTF-16(B|L)E,version=1: writeBOM | acceptBOM | requireBOM | strictBOM
	//    (UnicodeBig and UnicodeLittle in Java)
	// - RFC 2781-compliant, but less common interpretation for UTF-16(B|L)E:
	//    acceptBOM | strictBOM (e.g. assigned to CheckBOM).
	// This addition would be consistent with supporting ExpectBOM.
)

// Endianness is a UTF-16 encoding's default endianness.
type Endianness bool

const (
	// BigEndian is UTF-16BE.
	BigEndian Endianness = false
	// LittleEndian is UTF-16LE.
	LittleEndian Endianness = true
)

// ErrMissingBOM means that decoding UTF-16 input with ExpectBOM did not find a
// starting byte order mark.
var ErrMissingBOM = errors.New("encoding: missing byte order mark")

type utf16Encoding struct ***REMOVED***
	config
	mib identifier.MIB
***REMOVED***

type config struct ***REMOVED***
	endianness Endianness
	bomPolicy  BOMPolicy
***REMOVED***

func (u utf16Encoding) NewDecoder() *encoding.Decoder ***REMOVED***
	return &encoding.Decoder***REMOVED***Transformer: &utf16Decoder***REMOVED***
		initial: u.config,
		current: u.config,
	***REMOVED******REMOVED***
***REMOVED***

func (u utf16Encoding) NewEncoder() *encoding.Encoder ***REMOVED***
	return &encoding.Encoder***REMOVED***Transformer: &utf16Encoder***REMOVED***
		endianness:       u.endianness,
		initialBOMPolicy: u.bomPolicy,
		currentBOMPolicy: u.bomPolicy,
	***REMOVED******REMOVED***
***REMOVED***

func (u utf16Encoding) ID() (mib identifier.MIB, other string) ***REMOVED***
	return u.mib, ""
***REMOVED***

func (u utf16Encoding) String() string ***REMOVED***
	e, b := "B", ""
	if u.endianness == LittleEndian ***REMOVED***
		e = "L"
	***REMOVED***
	switch u.bomPolicy ***REMOVED***
	case ExpectBOM:
		b = "Expect"
	case UseBOM:
		b = "Use"
	case IgnoreBOM:
		b = "Ignore"
	***REMOVED***
	return "UTF-16" + e + "E (" + b + " BOM)"
***REMOVED***

type utf16Decoder struct ***REMOVED***
	initial config
	current config
***REMOVED***

func (u *utf16Decoder) Reset() ***REMOVED***
	u.current = u.initial
***REMOVED***

func (u *utf16Decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	if len(src) == 0 ***REMOVED***
		if atEOF && u.current.bomPolicy&requireBOM != 0 ***REMOVED***
			return 0, 0, ErrMissingBOM
		***REMOVED***
		return 0, 0, nil
	***REMOVED***
	if u.current.bomPolicy&acceptBOM != 0 ***REMOVED***
		if len(src) < 2 ***REMOVED***
			return 0, 0, transform.ErrShortSrc
		***REMOVED***
		switch ***REMOVED***
		case src[0] == 0xfe && src[1] == 0xff:
			u.current.endianness = BigEndian
			nSrc = 2
		case src[0] == 0xff && src[1] == 0xfe:
			u.current.endianness = LittleEndian
			nSrc = 2
		default:
			if u.current.bomPolicy&requireBOM != 0 ***REMOVED***
				return 0, 0, ErrMissingBOM
			***REMOVED***
		***REMOVED***
		u.current.bomPolicy = IgnoreBOM
	***REMOVED***

	var r rune
	var dSize, sSize int
	for nSrc < len(src) ***REMOVED***
		if nSrc+1 < len(src) ***REMOVED***
			x := uint16(src[nSrc+0])<<8 | uint16(src[nSrc+1])
			if u.current.endianness == LittleEndian ***REMOVED***
				x = x>>8 | x<<8
			***REMOVED***
			r, sSize = rune(x), 2
			if utf16.IsSurrogate(r) ***REMOVED***
				if nSrc+3 < len(src) ***REMOVED***
					x = uint16(src[nSrc+2])<<8 | uint16(src[nSrc+3])
					if u.current.endianness == LittleEndian ***REMOVED***
						x = x>>8 | x<<8
					***REMOVED***
					// Save for next iteration if it is not a high surrogate.
					if isHighSurrogate(rune(x)) ***REMOVED***
						r, sSize = utf16.DecodeRune(r, rune(x)), 4
					***REMOVED***
				***REMOVED*** else if !atEOF ***REMOVED***
					err = transform.ErrShortSrc
					break
				***REMOVED***
			***REMOVED***
			if dSize = utf8.RuneLen(r); dSize < 0 ***REMOVED***
				r, dSize = utf8.RuneError, 3
			***REMOVED***
		***REMOVED*** else if atEOF ***REMOVED***
			// Single trailing byte.
			r, dSize, sSize = utf8.RuneError, 3, 1
		***REMOVED*** else ***REMOVED***
			err = transform.ErrShortSrc
			break
		***REMOVED***
		if nDst+dSize > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break
		***REMOVED***
		nDst += utf8.EncodeRune(dst[nDst:], r)
		nSrc += sSize
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***

func isHighSurrogate(r rune) bool ***REMOVED***
	return 0xDC00 <= r && r <= 0xDFFF
***REMOVED***

type utf16Encoder struct ***REMOVED***
	endianness       Endianness
	initialBOMPolicy BOMPolicy
	currentBOMPolicy BOMPolicy
***REMOVED***

func (u *utf16Encoder) Reset() ***REMOVED***
	u.currentBOMPolicy = u.initialBOMPolicy
***REMOVED***

func (u *utf16Encoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	if u.currentBOMPolicy&writeBOM != 0 ***REMOVED***
		if len(dst) < 2 ***REMOVED***
			return 0, 0, transform.ErrShortDst
		***REMOVED***
		dst[0], dst[1] = 0xfe, 0xff
		u.currentBOMPolicy = IgnoreBOM
		nDst = 2
	***REMOVED***

	r, size := rune(0), 0
	for nSrc < len(src) ***REMOVED***
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
		***REMOVED***

		if r <= 0xffff ***REMOVED***
			if nDst+2 > len(dst) ***REMOVED***
				err = transform.ErrShortDst
				break
			***REMOVED***
			dst[nDst+0] = uint8(r >> 8)
			dst[nDst+1] = uint8(r)
			nDst += 2
		***REMOVED*** else ***REMOVED***
			if nDst+4 > len(dst) ***REMOVED***
				err = transform.ErrShortDst
				break
			***REMOVED***
			r1, r2 := utf16.EncodeRune(r)
			dst[nDst+0] = uint8(r1 >> 8)
			dst[nDst+1] = uint8(r1)
			dst[nDst+2] = uint8(r2 >> 8)
			dst[nDst+3] = uint8(r2)
			nDst += 4
		***REMOVED***
		nSrc += size
	***REMOVED***

	if u.endianness == LittleEndian ***REMOVED***
		for i := 0; i < nDst; i += 2 ***REMOVED***
			dst[i], dst[i+1] = dst[i+1], dst[i]
		***REMOVED***
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***
