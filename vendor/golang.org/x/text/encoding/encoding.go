// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package encoding defines an interface for character encodings, such as Shift
// JIS and Windows 1252, that can convert to and from UTF-8.
//
// Encoding implementations are provided in other packages, such as
// golang.org/x/text/encoding/charmap and
// golang.org/x/text/encoding/japanese.
package encoding // import "golang.org/x/text/encoding"

import (
	"errors"
	"io"
	"strconv"
	"unicode/utf8"

	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/transform"
)

// TODO:
// - There seems to be some inconsistency in when decoders return errors
//   and when not. Also documentation seems to suggest they shouldn't return
//   errors at all (except for UTF-16).
// - Encoders seem to rely on or at least benefit from the input being in NFC
//   normal form. Perhaps add an example how users could prepare their output.

// Encoding is a character set encoding that can be transformed to and from
// UTF-8.
type Encoding interface ***REMOVED***
	// NewDecoder returns a Decoder.
	NewDecoder() *Decoder

	// NewEncoder returns an Encoder.
	NewEncoder() *Encoder
***REMOVED***

// A Decoder converts bytes to UTF-8. It implements transform.Transformer.
//
// Transforming source bytes that are not of that encoding will not result in an
// error per se. Each byte that cannot be transcoded will be represented in the
// output by the UTF-8 encoding of '\uFFFD', the replacement rune.
type Decoder struct ***REMOVED***
	transform.Transformer

	// This forces external creators of Decoders to use names in struct
	// initializers, allowing for future extendibility without having to break
	// code.
	_ struct***REMOVED******REMOVED***
***REMOVED***

// Bytes converts the given encoded bytes to UTF-8. It returns the converted
// bytes or nil, err if any error occurred.
func (d *Decoder) Bytes(b []byte) ([]byte, error) ***REMOVED***
	b, _, err := transform.Bytes(d, b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b, nil
***REMOVED***

// String converts the given encoded string to UTF-8. It returns the converted
// string or "", err if any error occurred.
func (d *Decoder) String(s string) (string, error) ***REMOVED***
	s, _, err := transform.String(d, s)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return s, nil
***REMOVED***

// Reader wraps another Reader to decode its bytes.
//
// The Decoder may not be used for any other operation as long as the returned
// Reader is in use.
func (d *Decoder) Reader(r io.Reader) io.Reader ***REMOVED***
	return transform.NewReader(r, d)
***REMOVED***

// An Encoder converts bytes from UTF-8. It implements transform.Transformer.
//
// Each rune that cannot be transcoded will result in an error. In this case,
// the transform will consume all source byte up to, not including the offending
// rune. Transforming source bytes that are not valid UTF-8 will be replaced by
// `\uFFFD`. To return early with an error instead, use transform.Chain to
// preprocess the data with a UTF8Validator.
type Encoder struct ***REMOVED***
	transform.Transformer

	// This forces external creators of Encoders to use names in struct
	// initializers, allowing for future extendibility without having to break
	// code.
	_ struct***REMOVED******REMOVED***
***REMOVED***

// Bytes converts bytes from UTF-8. It returns the converted bytes or nil, err if
// any error occurred.
func (e *Encoder) Bytes(b []byte) ([]byte, error) ***REMOVED***
	b, _, err := transform.Bytes(e, b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b, nil
***REMOVED***

// String converts a string from UTF-8. It returns the converted string or
// "", err if any error occurred.
func (e *Encoder) String(s string) (string, error) ***REMOVED***
	s, _, err := transform.String(e, s)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return s, nil
***REMOVED***

// Writer wraps another Writer to encode its UTF-8 output.
//
// The Encoder may not be used for any other operation as long as the returned
// Writer is in use.
func (e *Encoder) Writer(w io.Writer) io.Writer ***REMOVED***
	return transform.NewWriter(w, e)
***REMOVED***

// ASCIISub is the ASCII substitute character, as recommended by
// http://unicode.org/reports/tr36/#Text_Comparison
const ASCIISub = '\x1a'

// Nop is the nop encoding. Its transformed bytes are the same as the source
// bytes; it does not replace invalid UTF-8 sequences.
var Nop Encoding = nop***REMOVED******REMOVED***

type nop struct***REMOVED******REMOVED***

func (nop) NewDecoder() *Decoder ***REMOVED***
	return &Decoder***REMOVED***Transformer: transform.Nop***REMOVED***
***REMOVED***
func (nop) NewEncoder() *Encoder ***REMOVED***
	return &Encoder***REMOVED***Transformer: transform.Nop***REMOVED***
***REMOVED***

// Replacement is the replacement encoding. Decoding from the replacement
// encoding yields a single '\uFFFD' replacement rune. Encoding from UTF-8 to
// the replacement encoding yields the same as the source bytes except that
// invalid UTF-8 is converted to '\uFFFD'.
//
// It is defined at http://encoding.spec.whatwg.org/#replacement
var Replacement Encoding = replacement***REMOVED******REMOVED***

type replacement struct***REMOVED******REMOVED***

func (replacement) NewDecoder() *Decoder ***REMOVED***
	return &Decoder***REMOVED***Transformer: replacementDecoder***REMOVED******REMOVED******REMOVED***
***REMOVED***

func (replacement) NewEncoder() *Encoder ***REMOVED***
	return &Encoder***REMOVED***Transformer: replacementEncoder***REMOVED******REMOVED******REMOVED***
***REMOVED***

func (replacement) ID() (mib identifier.MIB, other string) ***REMOVED***
	return identifier.Replacement, ""
***REMOVED***

type replacementDecoder struct***REMOVED*** transform.NopResetter ***REMOVED***

func (replacementDecoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	if len(dst) < 3 ***REMOVED***
		return 0, 0, transform.ErrShortDst
	***REMOVED***
	if atEOF ***REMOVED***
		const fffd = "\ufffd"
		dst[0] = fffd[0]
		dst[1] = fffd[1]
		dst[2] = fffd[2]
		nDst = 3
	***REMOVED***
	return nDst, len(src), nil
***REMOVED***

type replacementEncoder struct***REMOVED*** transform.NopResetter ***REMOVED***

func (replacementEncoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
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
				r = '\ufffd'
			***REMOVED***
		***REMOVED***

		if nDst+utf8.RuneLen(r) > len(dst) ***REMOVED***
			err = transform.ErrShortDst
			break
		***REMOVED***
		nDst += utf8.EncodeRune(dst[nDst:], r)
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***

// HTMLEscapeUnsupported wraps encoders to replace source runes outside the
// repertoire of the destination encoding with HTML escape sequences.
//
// This wrapper exists to comply to URL and HTML forms requiring a
// non-terminating legacy encoder. The produced sequences may lead to data
// loss as they are indistinguishable from legitimate input. To avoid this
// issue, use UTF-8 encodings whenever possible.
func HTMLEscapeUnsupported(e *Encoder) *Encoder ***REMOVED***
	return &Encoder***REMOVED***Transformer: &errorHandler***REMOVED***e, errorToHTML***REMOVED******REMOVED***
***REMOVED***

// ReplaceUnsupported wraps encoders to replace source runes outside the
// repertoire of the destination encoding with an encoding-specific
// replacement.
//
// This wrapper is only provided for backwards compatibility and legacy
// handling. Its use is strongly discouraged. Use UTF-8 whenever possible.
func ReplaceUnsupported(e *Encoder) *Encoder ***REMOVED***
	return &Encoder***REMOVED***Transformer: &errorHandler***REMOVED***e, errorToReplacement***REMOVED******REMOVED***
***REMOVED***

type errorHandler struct ***REMOVED***
	*Encoder
	handler func(dst []byte, r rune, err repertoireError) (n int, ok bool)
***REMOVED***

// TODO: consider making this error public in some form.
type repertoireError interface ***REMOVED***
	Replacement() byte
***REMOVED***

func (h errorHandler) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	nDst, nSrc, err = h.Transformer.Transform(dst, src, atEOF)
	for err != nil ***REMOVED***
		rerr, ok := err.(repertoireError)
		if !ok ***REMOVED***
			return nDst, nSrc, err
		***REMOVED***
		r, sz := utf8.DecodeRune(src[nSrc:])
		n, ok := h.handler(dst[nDst:], r, rerr)
		if !ok ***REMOVED***
			return nDst, nSrc, transform.ErrShortDst
		***REMOVED***
		err = nil
		nDst += n
		if nSrc += sz; nSrc < len(src) ***REMOVED***
			var dn, sn int
			dn, sn, err = h.Transformer.Transform(dst[nDst:], src[nSrc:], atEOF)
			nDst += dn
			nSrc += sn
		***REMOVED***
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***

func errorToHTML(dst []byte, r rune, err repertoireError) (n int, ok bool) ***REMOVED***
	buf := [8]byte***REMOVED******REMOVED***
	b := strconv.AppendUint(buf[:0], uint64(r), 10)
	if n = len(b) + len("&#;"); n >= len(dst) ***REMOVED***
		return 0, false
	***REMOVED***
	dst[0] = '&'
	dst[1] = '#'
	dst[copy(dst[2:], b)+2] = ';'
	return n, true
***REMOVED***

func errorToReplacement(dst []byte, r rune, err repertoireError) (n int, ok bool) ***REMOVED***
	if len(dst) == 0 ***REMOVED***
		return 0, false
	***REMOVED***
	dst[0] = err.Replacement()
	return 1, true
***REMOVED***

// ErrInvalidUTF8 means that a transformer encountered invalid UTF-8.
var ErrInvalidUTF8 = errors.New("encoding: invalid UTF-8")

// UTF8Validator is a transformer that returns ErrInvalidUTF8 on the first
// input byte that is not valid UTF-8.
var UTF8Validator transform.Transformer = utf8Validator***REMOVED******REMOVED***

type utf8Validator struct***REMOVED*** transform.NopResetter ***REMOVED***

func (utf8Validator) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	n := len(src)
	if n > len(dst) ***REMOVED***
		n = len(dst)
	***REMOVED***
	for i := 0; i < n; ***REMOVED***
		if c := src[i]; c < utf8.RuneSelf ***REMOVED***
			dst[i] = c
			i++
			continue
		***REMOVED***
		_, size := utf8.DecodeRune(src[i:])
		if size == 1 ***REMOVED***
			// All valid runes of size 1 (those below utf8.RuneSelf) were
			// handled above. We have invalid UTF-8 or we haven't seen the
			// full character yet.
			err = ErrInvalidUTF8
			if !atEOF && !utf8.FullRune(src[i:]) ***REMOVED***
				err = transform.ErrShortSrc
			***REMOVED***
			return i, i, err
		***REMOVED***
		if i+size > len(dst) ***REMOVED***
			return i, i, transform.ErrShortDst
		***REMOVED***
		for ; size > 0; size-- ***REMOVED***
			dst[i] = src[i]
			i++
		***REMOVED***
	***REMOVED***
	if len(src) > len(dst) ***REMOVED***
		err = transform.ErrShortDst
	***REMOVED***
	return n, n, err
***REMOVED***
