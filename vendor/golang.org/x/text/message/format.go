// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message

import (
	"bytes"
	"strconv"
	"unicode/utf8"

	"golang.org/x/text/internal/format"
)

const (
	ldigits = "0123456789abcdefx"
	udigits = "0123456789ABCDEFX"
)

const (
	signed   = true
	unsigned = false
)

// A formatInfo is the raw formatter used by Printf etc.
// It prints into a buffer that must be set up separately.
type formatInfo struct ***REMOVED***
	buf *bytes.Buffer

	format.Parser

	// intbuf is large enough to store %b of an int64 with a sign and
	// avoids padding at the end of the struct on 32 bit architectures.
	intbuf [68]byte
***REMOVED***

func (f *formatInfo) init(buf *bytes.Buffer) ***REMOVED***
	f.ClearFlags()
	f.buf = buf
***REMOVED***

// writePadding generates n bytes of padding.
func (f *formatInfo) writePadding(n int) ***REMOVED***
	if n <= 0 ***REMOVED*** // No padding bytes needed.
		return
	***REMOVED***
	f.buf.Grow(n)
	// Decide which byte the padding should be filled with.
	padByte := byte(' ')
	if f.Zero ***REMOVED***
		padByte = byte('0')
	***REMOVED***
	// Fill padding with padByte.
	for i := 0; i < n; i++ ***REMOVED***
		f.buf.WriteByte(padByte) // TODO: make more efficient.
	***REMOVED***
***REMOVED***

// pad appends b to f.buf, padded on left (!f.minus) or right (f.minus).
func (f *formatInfo) pad(b []byte) ***REMOVED***
	if !f.WidthPresent || f.Width == 0 ***REMOVED***
		f.buf.Write(b)
		return
	***REMOVED***
	width := f.Width - utf8.RuneCount(b)
	if !f.Minus ***REMOVED***
		// left padding
		f.writePadding(width)
		f.buf.Write(b)
	***REMOVED*** else ***REMOVED***
		// right padding
		f.buf.Write(b)
		f.writePadding(width)
	***REMOVED***
***REMOVED***

// padString appends s to f.buf, padded on left (!f.minus) or right (f.minus).
func (f *formatInfo) padString(s string) ***REMOVED***
	if !f.WidthPresent || f.Width == 0 ***REMOVED***
		f.buf.WriteString(s)
		return
	***REMOVED***
	width := f.Width - utf8.RuneCountInString(s)
	if !f.Minus ***REMOVED***
		// left padding
		f.writePadding(width)
		f.buf.WriteString(s)
	***REMOVED*** else ***REMOVED***
		// right padding
		f.buf.WriteString(s)
		f.writePadding(width)
	***REMOVED***
***REMOVED***

// fmt_boolean formats a boolean.
func (f *formatInfo) fmt_boolean(v bool) ***REMOVED***
	if v ***REMOVED***
		f.padString("true")
	***REMOVED*** else ***REMOVED***
		f.padString("false")
	***REMOVED***
***REMOVED***

// fmt_unicode formats a uint64 as "U+0078" or with f.sharp set as "U+0078 'x'".
func (f *formatInfo) fmt_unicode(u uint64) ***REMOVED***
	buf := f.intbuf[0:]

	// With default precision set the maximum needed buf length is 18
	// for formatting -1 with %#U ("U+FFFFFFFFFFFFFFFF") which fits
	// into the already allocated intbuf with a capacity of 68 bytes.
	prec := 4
	if f.PrecPresent && f.Prec > 4 ***REMOVED***
		prec = f.Prec
		// Compute space needed for "U+" , number, " '", character, "'".
		width := 2 + prec + 2 + utf8.UTFMax + 1
		if width > len(buf) ***REMOVED***
			buf = make([]byte, width)
		***REMOVED***
	***REMOVED***

	// Format into buf, ending at buf[i]. Formatting numbers is easier right-to-left.
	i := len(buf)

	// For %#U we want to add a space and a quoted character at the end of the buffer.
	if f.Sharp && u <= utf8.MaxRune && strconv.IsPrint(rune(u)) ***REMOVED***
		i--
		buf[i] = '\''
		i -= utf8.RuneLen(rune(u))
		utf8.EncodeRune(buf[i:], rune(u))
		i--
		buf[i] = '\''
		i--
		buf[i] = ' '
	***REMOVED***
	// Format the Unicode code point u as a hexadecimal number.
	for u >= 16 ***REMOVED***
		i--
		buf[i] = udigits[u&0xF]
		prec--
		u >>= 4
	***REMOVED***
	i--
	buf[i] = udigits[u]
	prec--
	// Add zeros in front of the number until requested precision is reached.
	for prec > 0 ***REMOVED***
		i--
		buf[i] = '0'
		prec--
	***REMOVED***
	// Add a leading "U+".
	i--
	buf[i] = '+'
	i--
	buf[i] = 'U'

	oldZero := f.Zero
	f.Zero = false
	f.pad(buf[i:])
	f.Zero = oldZero
***REMOVED***

// fmt_integer formats signed and unsigned integers.
func (f *formatInfo) fmt_integer(u uint64, base int, isSigned bool, digits string) ***REMOVED***
	negative := isSigned && int64(u) < 0
	if negative ***REMOVED***
		u = -u
	***REMOVED***

	buf := f.intbuf[0:]
	// The already allocated f.intbuf with a capacity of 68 bytes
	// is large enough for integer formatting when no precision or width is set.
	if f.WidthPresent || f.PrecPresent ***REMOVED***
		// Account 3 extra bytes for possible addition of a sign and "0x".
		width := 3 + f.Width + f.Prec // wid and prec are always positive.
		if width > len(buf) ***REMOVED***
			// We're going to need a bigger boat.
			buf = make([]byte, width)
		***REMOVED***
	***REMOVED***

	// Two ways to ask for extra leading zero digits: %.3d or %03d.
	// If both are specified the f.zero flag is ignored and
	// padding with spaces is used instead.
	prec := 0
	if f.PrecPresent ***REMOVED***
		prec = f.Prec
		// Precision of 0 and value of 0 means "print nothing" but padding.
		if prec == 0 && u == 0 ***REMOVED***
			oldZero := f.Zero
			f.Zero = false
			f.writePadding(f.Width)
			f.Zero = oldZero
			return
		***REMOVED***
	***REMOVED*** else if f.Zero && f.WidthPresent ***REMOVED***
		prec = f.Width
		if negative || f.Plus || f.Space ***REMOVED***
			prec-- // leave room for sign
		***REMOVED***
	***REMOVED***

	// Because printing is easier right-to-left: format u into buf, ending at buf[i].
	// We could make things marginally faster by splitting the 32-bit case out
	// into a separate block but it's not worth the duplication, so u has 64 bits.
	i := len(buf)
	// Use constants for the division and modulo for more efficient code.
	// Switch cases ordered by popularity.
	switch base ***REMOVED***
	case 10:
		for u >= 10 ***REMOVED***
			i--
			next := u / 10
			buf[i] = byte('0' + u - next*10)
			u = next
		***REMOVED***
	case 16:
		for u >= 16 ***REMOVED***
			i--
			buf[i] = digits[u&0xF]
			u >>= 4
		***REMOVED***
	case 8:
		for u >= 8 ***REMOVED***
			i--
			buf[i] = byte('0' + u&7)
			u >>= 3
		***REMOVED***
	case 2:
		for u >= 2 ***REMOVED***
			i--
			buf[i] = byte('0' + u&1)
			u >>= 1
		***REMOVED***
	default:
		panic("fmt: unknown base; can't happen")
	***REMOVED***
	i--
	buf[i] = digits[u]
	for i > 0 && prec > len(buf)-i ***REMOVED***
		i--
		buf[i] = '0'
	***REMOVED***

	// Various prefixes: 0x, -, etc.
	if f.Sharp ***REMOVED***
		switch base ***REMOVED***
		case 8:
			if buf[i] != '0' ***REMOVED***
				i--
				buf[i] = '0'
			***REMOVED***
		case 16:
			// Add a leading 0x or 0X.
			i--
			buf[i] = digits[16]
			i--
			buf[i] = '0'
		***REMOVED***
	***REMOVED***

	if negative ***REMOVED***
		i--
		buf[i] = '-'
	***REMOVED*** else if f.Plus ***REMOVED***
		i--
		buf[i] = '+'
	***REMOVED*** else if f.Space ***REMOVED***
		i--
		buf[i] = ' '
	***REMOVED***

	// Left padding with zeros has already been handled like precision earlier
	// or the f.zero flag is ignored due to an explicitly set precision.
	oldZero := f.Zero
	f.Zero = false
	f.pad(buf[i:])
	f.Zero = oldZero
***REMOVED***

// truncate truncates the string to the specified precision, if present.
func (f *formatInfo) truncate(s string) string ***REMOVED***
	if f.PrecPresent ***REMOVED***
		n := f.Prec
		for i := range s ***REMOVED***
			n--
			if n < 0 ***REMOVED***
				return s[:i]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***

// fmt_s formats a string.
func (f *formatInfo) fmt_s(s string) ***REMOVED***
	s = f.truncate(s)
	f.padString(s)
***REMOVED***

// fmt_sbx formats a string or byte slice as a hexadecimal encoding of its bytes.
func (f *formatInfo) fmt_sbx(s string, b []byte, digits string) ***REMOVED***
	length := len(b)
	if b == nil ***REMOVED***
		// No byte slice present. Assume string s should be encoded.
		length = len(s)
	***REMOVED***
	// Set length to not process more bytes than the precision demands.
	if f.PrecPresent && f.Prec < length ***REMOVED***
		length = f.Prec
	***REMOVED***
	// Compute width of the encoding taking into account the f.sharp and f.space flag.
	width := 2 * length
	if width > 0 ***REMOVED***
		if f.Space ***REMOVED***
			// Each element encoded by two hexadecimals will get a leading 0x or 0X.
			if f.Sharp ***REMOVED***
				width *= 2
			***REMOVED***
			// Elements will be separated by a space.
			width += length - 1
		***REMOVED*** else if f.Sharp ***REMOVED***
			// Only a leading 0x or 0X will be added for the whole string.
			width += 2
		***REMOVED***
	***REMOVED*** else ***REMOVED*** // The byte slice or string that should be encoded is empty.
		if f.WidthPresent ***REMOVED***
			f.writePadding(f.Width)
		***REMOVED***
		return
	***REMOVED***
	// Handle padding to the left.
	if f.WidthPresent && f.Width > width && !f.Minus ***REMOVED***
		f.writePadding(f.Width - width)
	***REMOVED***
	// Write the encoding directly into the output buffer.
	buf := f.buf
	if f.Sharp ***REMOVED***
		// Add leading 0x or 0X.
		buf.WriteByte('0')
		buf.WriteByte(digits[16])
	***REMOVED***
	var c byte
	for i := 0; i < length; i++ ***REMOVED***
		if f.Space && i > 0 ***REMOVED***
			// Separate elements with a space.
			buf.WriteByte(' ')
			if f.Sharp ***REMOVED***
				// Add leading 0x or 0X for each element.
				buf.WriteByte('0')
				buf.WriteByte(digits[16])
			***REMOVED***
		***REMOVED***
		if b != nil ***REMOVED***
			c = b[i] // Take a byte from the input byte slice.
		***REMOVED*** else ***REMOVED***
			c = s[i] // Take a byte from the input string.
		***REMOVED***
		// Encode each byte as two hexadecimal digits.
		buf.WriteByte(digits[c>>4])
		buf.WriteByte(digits[c&0xF])
	***REMOVED***
	// Handle padding to the right.
	if f.WidthPresent && f.Width > width && f.Minus ***REMOVED***
		f.writePadding(f.Width - width)
	***REMOVED***
***REMOVED***

// fmt_sx formats a string as a hexadecimal encoding of its bytes.
func (f *formatInfo) fmt_sx(s, digits string) ***REMOVED***
	f.fmt_sbx(s, nil, digits)
***REMOVED***

// fmt_bx formats a byte slice as a hexadecimal encoding of its bytes.
func (f *formatInfo) fmt_bx(b []byte, digits string) ***REMOVED***
	f.fmt_sbx("", b, digits)
***REMOVED***

// fmt_q formats a string as a double-quoted, escaped Go string constant.
// If f.sharp is set a raw (backquoted) string may be returned instead
// if the string does not contain any control characters other than tab.
func (f *formatInfo) fmt_q(s string) ***REMOVED***
	s = f.truncate(s)
	if f.Sharp && strconv.CanBackquote(s) ***REMOVED***
		f.padString("`" + s + "`")
		return
	***REMOVED***
	buf := f.intbuf[:0]
	if f.Plus ***REMOVED***
		f.pad(strconv.AppendQuoteToASCII(buf, s))
	***REMOVED*** else ***REMOVED***
		f.pad(strconv.AppendQuote(buf, s))
	***REMOVED***
***REMOVED***

// fmt_c formats an integer as a Unicode character.
// If the character is not valid Unicode, it will print '\ufffd'.
func (f *formatInfo) fmt_c(c uint64) ***REMOVED***
	r := rune(c)
	if c > utf8.MaxRune ***REMOVED***
		r = utf8.RuneError
	***REMOVED***
	buf := f.intbuf[:0]
	w := utf8.EncodeRune(buf[:utf8.UTFMax], r)
	f.pad(buf[:w])
***REMOVED***

// fmt_qc formats an integer as a single-quoted, escaped Go character constant.
// If the character is not valid Unicode, it will print '\ufffd'.
func (f *formatInfo) fmt_qc(c uint64) ***REMOVED***
	r := rune(c)
	if c > utf8.MaxRune ***REMOVED***
		r = utf8.RuneError
	***REMOVED***
	buf := f.intbuf[:0]
	if f.Plus ***REMOVED***
		f.pad(strconv.AppendQuoteRuneToASCII(buf, r))
	***REMOVED*** else ***REMOVED***
		f.pad(strconv.AppendQuoteRune(buf, r))
	***REMOVED***
***REMOVED***

// fmt_float formats a float64. It assumes that verb is a valid format specifier
// for strconv.AppendFloat and therefore fits into a byte.
func (f *formatInfo) fmt_float(v float64, size int, verb rune, prec int) ***REMOVED***
	// Explicit precision in format specifier overrules default precision.
	if f.PrecPresent ***REMOVED***
		prec = f.Prec
	***REMOVED***
	// Format number, reserving space for leading + sign if needed.
	num := strconv.AppendFloat(f.intbuf[:1], v, byte(verb), prec, size)
	if num[1] == '-' || num[1] == '+' ***REMOVED***
		num = num[1:]
	***REMOVED*** else ***REMOVED***
		num[0] = '+'
	***REMOVED***
	// f.space means to add a leading space instead of a "+" sign unless
	// the sign is explicitly asked for by f.plus.
	if f.Space && num[0] == '+' && !f.Plus ***REMOVED***
		num[0] = ' '
	***REMOVED***
	// Special handling for infinities and NaN,
	// which don't look like a number so shouldn't be padded with zeros.
	if num[1] == 'I' || num[1] == 'N' ***REMOVED***
		oldZero := f.Zero
		f.Zero = false
		// Remove sign before NaN if not asked for.
		if num[1] == 'N' && !f.Space && !f.Plus ***REMOVED***
			num = num[1:]
		***REMOVED***
		f.pad(num)
		f.Zero = oldZero
		return
	***REMOVED***
	// The sharp flag forces printing a decimal point for non-binary formats
	// and retains trailing zeros, which we may need to restore.
	if f.Sharp && verb != 'b' ***REMOVED***
		digits := 0
		switch verb ***REMOVED***
		case 'v', 'g', 'G':
			digits = prec
			// If no precision is set explicitly use a precision of 6.
			if digits == -1 ***REMOVED***
				digits = 6
			***REMOVED***
		***REMOVED***

		// Buffer pre-allocated with enough room for
		// exponent notations of the form "e+123".
		var tailBuf [5]byte
		tail := tailBuf[:0]

		hasDecimalPoint := false
		// Starting from i = 1 to skip sign at num[0].
		for i := 1; i < len(num); i++ ***REMOVED***
			switch num[i] ***REMOVED***
			case '.':
				hasDecimalPoint = true
			case 'e', 'E':
				tail = append(tail, num[i:]...)
				num = num[:i]
			default:
				digits--
			***REMOVED***
		***REMOVED***
		if !hasDecimalPoint ***REMOVED***
			num = append(num, '.')
		***REMOVED***
		for digits > 0 ***REMOVED***
			num = append(num, '0')
			digits--
		***REMOVED***
		num = append(num, tail...)
	***REMOVED***
	// We want a sign if asked for and if the sign is not positive.
	if f.Plus || num[0] != '+' ***REMOVED***
		// If we're zero padding to the left we want the sign before the leading zeros.
		// Achieve this by writing the sign out and then padding the unsigned number.
		if f.Zero && f.WidthPresent && f.Width > len(num) ***REMOVED***
			f.buf.WriteByte(num[0])
			f.writePadding(f.Width - len(num))
			f.buf.Write(num[1:])
			return
		***REMOVED***
		f.pad(num)
		return
	***REMOVED***
	// No sign to show and the number is positive; just print the unsigned number.
	f.pad(num[1:])
***REMOVED***
