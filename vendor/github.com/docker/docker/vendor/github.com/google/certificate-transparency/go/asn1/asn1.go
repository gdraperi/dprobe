// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package asn1 implements parsing of DER-encoded ASN.1 data structures,
// as defined in ITU-T Rec X.690.
//
// See also ``A Layman's Guide to a Subset of ASN.1, BER, and DER,''
// http://luca.ntop.org/Teaching/Appunti/asn1.html.
//
// START CT CHANGES
// This is a fork of the Go standard library ASN.1 implementation
// (encoding/asn1).  The main difference is that this version tries to correct
// for errors (e.g. use of tagPrintableString when the string data is really
// ISO8859-1 - a common error present in many x509 certificates in the wild.)
// END CT CHANGES
package asn1

// ASN.1 is a syntax for specifying abstract objects and BER, DER, PER, XER etc
// are different encoding formats for those objects. Here, we'll be dealing
// with DER, the Distinguished Encoding Rules. DER is used in X.509 because
// it's fast to parse and, unlike BER, has a unique encoding for every object.
// When calculating hashes over objects, it's important that the resulting
// bytes be the same at both ends and DER removes this margin of error.
//
// ASN.1 is very complex and this package doesn't attempt to implement
// everything by any means.

import (
	// START CT CHANGES
	"errors"
	"fmt"
	// END CT CHANGES
	"math/big"
	"reflect"
	// START CT CHANGES
	"strings"
	// END CT CHANGES
	"time"
)

// A StructuralError suggests that the ASN.1 data is valid, but the Go type
// which is receiving it doesn't match.
type StructuralError struct ***REMOVED***
	Msg string
***REMOVED***

func (e StructuralError) Error() string ***REMOVED*** return "asn1: structure error: " + e.Msg ***REMOVED***

// A SyntaxError suggests that the ASN.1 data is invalid.
type SyntaxError struct ***REMOVED***
	Msg string
***REMOVED***

func (e SyntaxError) Error() string ***REMOVED*** return "asn1: syntax error: " + e.Msg ***REMOVED***

// We start by dealing with each of the primitive types in turn.

// BOOLEAN

func parseBool(bytes []byte) (ret bool, err error) ***REMOVED***
	if len(bytes) != 1 ***REMOVED***
		err = SyntaxError***REMOVED***"invalid boolean"***REMOVED***
		return
	***REMOVED***

	// DER demands that "If the encoding represents the boolean value TRUE,
	// its single contents octet shall have all eight bits set to one."
	// Thus only 0 and 255 are valid encoded values.
	switch bytes[0] ***REMOVED***
	case 0:
		ret = false
	case 0xff:
		ret = true
	default:
		err = SyntaxError***REMOVED***"invalid boolean"***REMOVED***
	***REMOVED***

	return
***REMOVED***

// INTEGER

// parseInt64 treats the given bytes as a big-endian, signed integer and
// returns the result.
func parseInt64(bytes []byte) (ret int64, err error) ***REMOVED***
	if len(bytes) > 8 ***REMOVED***
		// We'll overflow an int64 in this case.
		err = StructuralError***REMOVED***"integer too large"***REMOVED***
		return
	***REMOVED***
	for bytesRead := 0; bytesRead < len(bytes); bytesRead++ ***REMOVED***
		ret <<= 8
		ret |= int64(bytes[bytesRead])
	***REMOVED***

	// Shift up and down in order to sign extend the result.
	ret <<= 64 - uint8(len(bytes))*8
	ret >>= 64 - uint8(len(bytes))*8
	return
***REMOVED***

// parseInt treats the given bytes as a big-endian, signed integer and returns
// the result.
func parseInt32(bytes []byte) (int32, error) ***REMOVED***
	ret64, err := parseInt64(bytes)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if ret64 != int64(int32(ret64)) ***REMOVED***
		return 0, StructuralError***REMOVED***"integer too large"***REMOVED***
	***REMOVED***
	return int32(ret64), nil
***REMOVED***

var bigOne = big.NewInt(1)

// parseBigInt treats the given bytes as a big-endian, signed integer and returns
// the result.
func parseBigInt(bytes []byte) *big.Int ***REMOVED***
	ret := new(big.Int)
	if len(bytes) > 0 && bytes[0]&0x80 == 0x80 ***REMOVED***
		// This is a negative number.
		notBytes := make([]byte, len(bytes))
		for i := range notBytes ***REMOVED***
			notBytes[i] = ^bytes[i]
		***REMOVED***
		ret.SetBytes(notBytes)
		ret.Add(ret, bigOne)
		ret.Neg(ret)
		return ret
	***REMOVED***
	ret.SetBytes(bytes)
	return ret
***REMOVED***

// BIT STRING

// BitString is the structure to use when you want an ASN.1 BIT STRING type. A
// bit string is padded up to the nearest byte in memory and the number of
// valid bits is recorded. Padding bits will be zero.
type BitString struct ***REMOVED***
	Bytes     []byte // bits packed into bytes.
	BitLength int    // length in bits.
***REMOVED***

// At returns the bit at the given index. If the index is out of range it
// returns false.
func (b BitString) At(i int) int ***REMOVED***
	if i < 0 || i >= b.BitLength ***REMOVED***
		return 0
	***REMOVED***
	x := i / 8
	y := 7 - uint(i%8)
	return int(b.Bytes[x]>>y) & 1
***REMOVED***

// RightAlign returns a slice where the padding bits are at the beginning. The
// slice may share memory with the BitString.
func (b BitString) RightAlign() []byte ***REMOVED***
	shift := uint(8 - (b.BitLength % 8))
	if shift == 8 || len(b.Bytes) == 0 ***REMOVED***
		return b.Bytes
	***REMOVED***

	a := make([]byte, len(b.Bytes))
	a[0] = b.Bytes[0] >> shift
	for i := 1; i < len(b.Bytes); i++ ***REMOVED***
		a[i] = b.Bytes[i-1] << (8 - shift)
		a[i] |= b.Bytes[i] >> shift
	***REMOVED***

	return a
***REMOVED***

// parseBitString parses an ASN.1 bit string from the given byte slice and returns it.
func parseBitString(bytes []byte) (ret BitString, err error) ***REMOVED***
	if len(bytes) == 0 ***REMOVED***
		err = SyntaxError***REMOVED***"zero length BIT STRING"***REMOVED***
		return
	***REMOVED***
	paddingBits := int(bytes[0])
	if paddingBits > 7 ||
		len(bytes) == 1 && paddingBits > 0 ||
		bytes[len(bytes)-1]&((1<<bytes[0])-1) != 0 ***REMOVED***
		err = SyntaxError***REMOVED***"invalid padding bits in BIT STRING"***REMOVED***
		return
	***REMOVED***
	ret.BitLength = (len(bytes)-1)*8 - paddingBits
	ret.Bytes = bytes[1:]
	return
***REMOVED***

// OBJECT IDENTIFIER

// An ObjectIdentifier represents an ASN.1 OBJECT IDENTIFIER.
type ObjectIdentifier []int

// Equal reports whether oi and other represent the same identifier.
func (oi ObjectIdentifier) Equal(other ObjectIdentifier) bool ***REMOVED***
	if len(oi) != len(other) ***REMOVED***
		return false
	***REMOVED***
	for i := 0; i < len(oi); i++ ***REMOVED***
		if oi[i] != other[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// parseObjectIdentifier parses an OBJECT IDENTIFIER from the given bytes and
// returns it. An object identifier is a sequence of variable length integers
// that are assigned in a hierarchy.
func parseObjectIdentifier(bytes []byte) (s []int, err error) ***REMOVED***
	if len(bytes) == 0 ***REMOVED***
		err = SyntaxError***REMOVED***"zero length OBJECT IDENTIFIER"***REMOVED***
		return
	***REMOVED***

	// In the worst case, we get two elements from the first byte (which is
	// encoded differently) and then every varint is a single byte long.
	s = make([]int, len(bytes)+1)

	// The first varint is 40*value1 + value2:
	// According to this packing, value1 can take the values 0, 1 and 2 only.
	// When value1 = 0 or value1 = 1, then value2 is <= 39. When value1 = 2,
	// then there are no restrictions on value2.
	v, offset, err := parseBase128Int(bytes, 0)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if v < 80 ***REMOVED***
		s[0] = v / 40
		s[1] = v % 40
	***REMOVED*** else ***REMOVED***
		s[0] = 2
		s[1] = v - 80
	***REMOVED***

	i := 2
	for ; offset < len(bytes); i++ ***REMOVED***
		v, offset, err = parseBase128Int(bytes, offset)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		s[i] = v
	***REMOVED***
	s = s[0:i]
	return
***REMOVED***

// ENUMERATED

// An Enumerated is represented as a plain int.
type Enumerated int

// FLAG

// A Flag accepts any data and is set to true if present.
type Flag bool

// parseBase128Int parses a base-128 encoded int from the given offset in the
// given byte slice. It returns the value and the new offset.
func parseBase128Int(bytes []byte, initOffset int) (ret, offset int, err error) ***REMOVED***
	offset = initOffset
	for shifted := 0; offset < len(bytes); shifted++ ***REMOVED***
		if shifted > 4 ***REMOVED***
			err = StructuralError***REMOVED***"base 128 integer too large"***REMOVED***
			return
		***REMOVED***
		ret <<= 7
		b := bytes[offset]
		ret |= int(b & 0x7f)
		offset++
		if b&0x80 == 0 ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	err = SyntaxError***REMOVED***"truncated base 128 integer"***REMOVED***
	return
***REMOVED***

// UTCTime

func parseUTCTime(bytes []byte) (ret time.Time, err error) ***REMOVED***
	s := string(bytes)
	ret, err = time.Parse("0601021504Z0700", s)
	if err != nil ***REMOVED***
		ret, err = time.Parse("060102150405Z0700", s)
	***REMOVED***
	if err == nil && ret.Year() >= 2050 ***REMOVED***
		// UTCTime only encodes times prior to 2050. See https://tools.ietf.org/html/rfc5280#section-4.1.2.5.1
		ret = ret.AddDate(-100, 0, 0)
	***REMOVED***

	return
***REMOVED***

// parseGeneralizedTime parses the GeneralizedTime from the given byte slice
// and returns the resulting time.
func parseGeneralizedTime(bytes []byte) (ret time.Time, err error) ***REMOVED***
	return time.Parse("20060102150405Z0700", string(bytes))
***REMOVED***

// PrintableString

// parsePrintableString parses a ASN.1 PrintableString from the given byte
// array and returns it.
func parsePrintableString(bytes []byte) (ret string, err error) ***REMOVED***
	for _, b := range bytes ***REMOVED***
		if !isPrintable(b) ***REMOVED***
			err = SyntaxError***REMOVED***"PrintableString contains invalid character"***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	ret = string(bytes)
	return
***REMOVED***

// isPrintable returns true iff the given b is in the ASN.1 PrintableString set.
func isPrintable(b byte) bool ***REMOVED***
	return 'a' <= b && b <= 'z' ||
		'A' <= b && b <= 'Z' ||
		'0' <= b && b <= '9' ||
		'\'' <= b && b <= ')' ||
		'+' <= b && b <= '/' ||
		b == ' ' ||
		b == ':' ||
		b == '=' ||
		b == '?' ||
		// This is technically not allowed in a PrintableString.
		// However, x509 certificates with wildcard strings don't
		// always use the correct string type so we permit it.
		b == '*'
***REMOVED***

// IA5String

// parseIA5String parses a ASN.1 IA5String (ASCII string) from the given
// byte slice and returns it.
func parseIA5String(bytes []byte) (ret string, err error) ***REMOVED***
	for _, b := range bytes ***REMOVED***
		if b >= 0x80 ***REMOVED***
			err = SyntaxError***REMOVED***"IA5String contains invalid character"***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	ret = string(bytes)
	return
***REMOVED***

// T61String

// parseT61String parses a ASN.1 T61String (8-bit clean string) from the given
// byte slice and returns it.
func parseT61String(bytes []byte) (ret string, err error) ***REMOVED***
	return string(bytes), nil
***REMOVED***

// UTF8String

// parseUTF8String parses a ASN.1 UTF8String (raw UTF-8) from the given byte
// array and returns it.
func parseUTF8String(bytes []byte) (ret string, err error) ***REMOVED***
	return string(bytes), nil
***REMOVED***

// A RawValue represents an undecoded ASN.1 object.
type RawValue struct ***REMOVED***
	Class, Tag int
	IsCompound bool
	Bytes      []byte
	FullBytes  []byte // includes the tag and length
***REMOVED***

// RawContent is used to signal that the undecoded, DER data needs to be
// preserved for a struct. To use it, the first field of the struct must have
// this type. It's an error for any of the other fields to have this type.
type RawContent []byte

// Tagging

// parseTagAndLength parses an ASN.1 tag and length pair from the given offset
// into a byte slice. It returns the parsed data and the new offset. SET and
// SET OF (tag 17) are mapped to SEQUENCE and SEQUENCE OF (tag 16) since we
// don't distinguish between ordered and unordered objects in this code.
func parseTagAndLength(bytes []byte, initOffset int) (ret tagAndLength, offset int, err error) ***REMOVED***
	offset = initOffset
	b := bytes[offset]
	offset++
	ret.class = int(b >> 6)
	ret.isCompound = b&0x20 == 0x20
	ret.tag = int(b & 0x1f)

	// If the bottom five bits are set, then the tag number is actually base 128
	// encoded afterwards
	if ret.tag == 0x1f ***REMOVED***
		ret.tag, offset, err = parseBase128Int(bytes, offset)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	if offset >= len(bytes) ***REMOVED***
		err = SyntaxError***REMOVED***"truncated tag or length"***REMOVED***
		return
	***REMOVED***
	b = bytes[offset]
	offset++
	if b&0x80 == 0 ***REMOVED***
		// The length is encoded in the bottom 7 bits.
		ret.length = int(b & 0x7f)
	***REMOVED*** else ***REMOVED***
		// Bottom 7 bits give the number of length bytes to follow.
		numBytes := int(b & 0x7f)
		if numBytes == 0 ***REMOVED***
			err = SyntaxError***REMOVED***"indefinite length found (not DER)"***REMOVED***
			return
		***REMOVED***
		ret.length = 0
		for i := 0; i < numBytes; i++ ***REMOVED***
			if offset >= len(bytes) ***REMOVED***
				err = SyntaxError***REMOVED***"truncated tag or length"***REMOVED***
				return
			***REMOVED***
			b = bytes[offset]
			offset++
			if ret.length >= 1<<23 ***REMOVED***
				// We can't shift ret.length up without
				// overflowing.
				err = StructuralError***REMOVED***"length too large"***REMOVED***
				return
			***REMOVED***
			ret.length <<= 8
			ret.length |= int(b)
			if ret.length == 0 ***REMOVED***
				// DER requires that lengths be minimal.
				err = StructuralError***REMOVED***"superfluous leading zeros in length"***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

// parseSequenceOf is used for SEQUENCE OF and SET OF values. It tries to parse
// a number of ASN.1 values from the given byte slice and returns them as a
// slice of Go values of the given type.
func parseSequenceOf(bytes []byte, sliceType reflect.Type, elemType reflect.Type) (ret reflect.Value, err error) ***REMOVED***
	expectedTag, compoundType, ok := getUniversalType(elemType)
	if !ok ***REMOVED***
		err = StructuralError***REMOVED***"unknown Go type for slice"***REMOVED***
		return
	***REMOVED***

	// First we iterate over the input and count the number of elements,
	// checking that the types are correct in each case.
	numElements := 0
	for offset := 0; offset < len(bytes); ***REMOVED***
		var t tagAndLength
		t, offset, err = parseTagAndLength(bytes, offset)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		// We pretend that GENERAL STRINGs are PRINTABLE STRINGs so
		// that a sequence of them can be parsed into a []string.
		if t.tag == tagGeneralString ***REMOVED***
			t.tag = tagPrintableString
		***REMOVED***
		if t.class != classUniversal || t.isCompound != compoundType || t.tag != expectedTag ***REMOVED***
			err = StructuralError***REMOVED***"sequence tag mismatch"***REMOVED***
			return
		***REMOVED***
		if invalidLength(offset, t.length, len(bytes)) ***REMOVED***
			err = SyntaxError***REMOVED***"truncated sequence"***REMOVED***
			return
		***REMOVED***
		offset += t.length
		numElements++
	***REMOVED***
	ret = reflect.MakeSlice(sliceType, numElements, numElements)
	params := fieldParameters***REMOVED******REMOVED***
	offset := 0
	for i := 0; i < numElements; i++ ***REMOVED***
		offset, err = parseField(ret.Index(i), bytes, offset, params)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

var (
	bitStringType        = reflect.TypeOf(BitString***REMOVED******REMOVED***)
	objectIdentifierType = reflect.TypeOf(ObjectIdentifier***REMOVED******REMOVED***)
	enumeratedType       = reflect.TypeOf(Enumerated(0))
	flagType             = reflect.TypeOf(Flag(false))
	timeType             = reflect.TypeOf(time.Time***REMOVED******REMOVED***)
	rawValueType         = reflect.TypeOf(RawValue***REMOVED******REMOVED***)
	rawContentsType      = reflect.TypeOf(RawContent(nil))
	bigIntType           = reflect.TypeOf(new(big.Int))
)

// invalidLength returns true iff offset + length > sliceLength, or if the
// addition would overflow.
func invalidLength(offset, length, sliceLength int) bool ***REMOVED***
	return offset+length < offset || offset+length > sliceLength
***REMOVED***

// START CT CHANGES

// Tests whether the data in |bytes| would be a valid ISO8859-1 string.
// Clearly, a sequence of bytes comprised solely of valid ISO8859-1
// codepoints does not imply that the encoding MUST be ISO8859-1, rather that
// you would not encounter an error trying to interpret the data as such.
func couldBeISO8859_1(bytes []byte) bool ***REMOVED***
	for _, b := range bytes ***REMOVED***
		if b < 0x20 || (b >= 0x7F && b < 0xA0) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Checks whether the data in |bytes| would be a valid T.61 string.
// Clearly, a sequence of bytes comprised solely of valid T.61
// codepoints does not imply that the encoding MUST be T.61, rather that
// you would not encounter an error trying to interpret the data as such.
func couldBeT61(bytes []byte) bool ***REMOVED***
	for _, b := range bytes ***REMOVED***
		switch b ***REMOVED***
		case 0x00:
			// Since we're guessing at (incorrect) encodings for a
			// PrintableString, we'll err on the side of caution and disallow
			// strings with a NUL in them, don't want to re-create a PayPal NUL
			// situation in monitors.
			fallthrough
		case 0x23, 0x24, 0x5C, 0x5E, 0x60, 0x7B, 0x7D, 0x7E, 0xA5, 0xA6, 0xAC, 0xAD, 0xAE, 0xAF,
			0xB9, 0xBA, 0xC0, 0xC9, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9,
			0xDA, 0xDB, 0xDC, 0xDE, 0xDF, 0xE5, 0xFF:
			// These are all invalid code points in T.61, so it can't be a T.61 string.
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Converts the data in |bytes| to the equivalent UTF-8 string.
func iso8859_1ToUTF8(bytes []byte) string ***REMOVED***
	buf := make([]rune, len(bytes))
	for i, b := range bytes ***REMOVED***
		buf[i] = rune(b)
	***REMOVED***
	return string(buf)
***REMOVED***

// END CT CHANGES

// parseField is the main parsing function. Given a byte slice and an offset
// into the array, it will try to parse a suitable ASN.1 value out and store it
// in the given Value.
func parseField(v reflect.Value, bytes []byte, initOffset int, params fieldParameters) (offset int, err error) ***REMOVED***
	offset = initOffset
	fieldType := v.Type()

	// If we have run out of data, it may be that there are optional elements at the end.
	if offset == len(bytes) ***REMOVED***
		if !setDefaultValue(v, params) ***REMOVED***
			err = SyntaxError***REMOVED***"sequence truncated"***REMOVED***
		***REMOVED***
		return
	***REMOVED***

	// Deal with raw values.
	if fieldType == rawValueType ***REMOVED***
		var t tagAndLength
		t, offset, err = parseTagAndLength(bytes, offset)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if invalidLength(offset, t.length, len(bytes)) ***REMOVED***
			err = SyntaxError***REMOVED***"data truncated"***REMOVED***
			return
		***REMOVED***
		result := RawValue***REMOVED***t.class, t.tag, t.isCompound, bytes[offset : offset+t.length], bytes[initOffset : offset+t.length]***REMOVED***
		offset += t.length
		v.Set(reflect.ValueOf(result))
		return
	***REMOVED***

	// Deal with the ANY type.
	if ifaceType := fieldType; ifaceType.Kind() == reflect.Interface && ifaceType.NumMethod() == 0 ***REMOVED***
		var t tagAndLength
		t, offset, err = parseTagAndLength(bytes, offset)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if invalidLength(offset, t.length, len(bytes)) ***REMOVED***
			err = SyntaxError***REMOVED***"data truncated"***REMOVED***
			return
		***REMOVED***
		var result interface***REMOVED******REMOVED***
		if !t.isCompound && t.class == classUniversal ***REMOVED***
			innerBytes := bytes[offset : offset+t.length]
			switch t.tag ***REMOVED***
			case tagPrintableString:
				result, err = parsePrintableString(innerBytes)
				// START CT CHANGES
				if err != nil && strings.Contains(err.Error(), "PrintableString contains invalid character") ***REMOVED***
					// Probably an ISO8859-1 string stuffed in, check if it
					// would be valid and assume that's what's happened if so,
					// otherwise try T.61, failing that give up and just assign
					// the bytes
					switch ***REMOVED***
					case couldBeISO8859_1(innerBytes):
						result, err = iso8859_1ToUTF8(innerBytes), nil
					case couldBeT61(innerBytes):
						result, err = parseT61String(innerBytes)
					default:
						result = nil
						err = errors.New("PrintableString contains invalid character, but couldn't determine correct String type.")
					***REMOVED***
				***REMOVED***
				// END CT CHANGES
			case tagIA5String:
				result, err = parseIA5String(innerBytes)
			case tagT61String:
				result, err = parseT61String(innerBytes)
			case tagUTF8String:
				result, err = parseUTF8String(innerBytes)
			case tagInteger:
				result, err = parseInt64(innerBytes)
			case tagBitString:
				result, err = parseBitString(innerBytes)
			case tagOID:
				result, err = parseObjectIdentifier(innerBytes)
			case tagUTCTime:
				result, err = parseUTCTime(innerBytes)
			case tagOctetString:
				result = innerBytes
			default:
				// If we don't know how to handle the type, we just leave Value as nil.
			***REMOVED***
		***REMOVED***
		offset += t.length
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if result != nil ***REMOVED***
			v.Set(reflect.ValueOf(result))
		***REMOVED***
		return
	***REMOVED***
	universalTag, compoundType, ok1 := getUniversalType(fieldType)
	if !ok1 ***REMOVED***
		err = StructuralError***REMOVED***fmt.Sprintf("unknown Go type: %v", fieldType)***REMOVED***
		return
	***REMOVED***

	t, offset, err := parseTagAndLength(bytes, offset)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if params.explicit ***REMOVED***
		expectedClass := classContextSpecific
		if params.application ***REMOVED***
			expectedClass = classApplication
		***REMOVED***
		if t.class == expectedClass && t.tag == *params.tag && (t.length == 0 || t.isCompound) ***REMOVED***
			if t.length > 0 ***REMOVED***
				t, offset, err = parseTagAndLength(bytes, offset)
				if err != nil ***REMOVED***
					return
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if fieldType != flagType ***REMOVED***
					err = StructuralError***REMOVED***"zero length explicit tag was not an asn1.Flag"***REMOVED***
					return
				***REMOVED***
				v.SetBool(true)
				return
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// The tags didn't match, it might be an optional element.
			ok := setDefaultValue(v, params)
			if ok ***REMOVED***
				offset = initOffset
			***REMOVED*** else ***REMOVED***
				err = StructuralError***REMOVED***"explicitly tagged member didn't match"***REMOVED***
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	// Special case for strings: all the ASN.1 string types map to the Go
	// type string. getUniversalType returns the tag for PrintableString
	// when it sees a string, so if we see a different string type on the
	// wire, we change the universal type to match.
	if universalTag == tagPrintableString ***REMOVED***
		switch t.tag ***REMOVED***
		case tagIA5String, tagGeneralString, tagT61String, tagUTF8String:
			universalTag = t.tag
		***REMOVED***
	***REMOVED***

	// Special case for time: UTCTime and GeneralizedTime both map to the
	// Go type time.Time.
	if universalTag == tagUTCTime && t.tag == tagGeneralizedTime ***REMOVED***
		universalTag = tagGeneralizedTime
	***REMOVED***

	expectedClass := classUniversal
	expectedTag := universalTag

	if !params.explicit && params.tag != nil ***REMOVED***
		expectedClass = classContextSpecific
		expectedTag = *params.tag
	***REMOVED***

	if !params.explicit && params.application && params.tag != nil ***REMOVED***
		expectedClass = classApplication
		expectedTag = *params.tag
	***REMOVED***

	// We have unwrapped any explicit tagging at this point.
	if t.class != expectedClass || t.tag != expectedTag || t.isCompound != compoundType ***REMOVED***
		// Tags don't match. Again, it could be an optional element.
		ok := setDefaultValue(v, params)
		if ok ***REMOVED***
			offset = initOffset
		***REMOVED*** else ***REMOVED***
			err = StructuralError***REMOVED***fmt.Sprintf("tags don't match (%d vs %+v) %+v %s @%d", expectedTag, t, params, fieldType.Name(), offset)***REMOVED***
		***REMOVED***
		return
	***REMOVED***
	if invalidLength(offset, t.length, len(bytes)) ***REMOVED***
		err = SyntaxError***REMOVED***"data truncated"***REMOVED***
		return
	***REMOVED***
	innerBytes := bytes[offset : offset+t.length]
	offset += t.length

	// We deal with the structures defined in this package first.
	switch fieldType ***REMOVED***
	case objectIdentifierType:
		newSlice, err1 := parseObjectIdentifier(innerBytes)
		v.Set(reflect.MakeSlice(v.Type(), len(newSlice), len(newSlice)))
		if err1 == nil ***REMOVED***
			reflect.Copy(v, reflect.ValueOf(newSlice))
		***REMOVED***
		err = err1
		return
	case bitStringType:
		bs, err1 := parseBitString(innerBytes)
		if err1 == nil ***REMOVED***
			v.Set(reflect.ValueOf(bs))
		***REMOVED***
		err = err1
		return
	case timeType:
		var time time.Time
		var err1 error
		if universalTag == tagUTCTime ***REMOVED***
			time, err1 = parseUTCTime(innerBytes)
		***REMOVED*** else ***REMOVED***
			time, err1 = parseGeneralizedTime(innerBytes)
		***REMOVED***
		if err1 == nil ***REMOVED***
			v.Set(reflect.ValueOf(time))
		***REMOVED***
		err = err1
		return
	case enumeratedType:
		parsedInt, err1 := parseInt32(innerBytes)
		if err1 == nil ***REMOVED***
			v.SetInt(int64(parsedInt))
		***REMOVED***
		err = err1
		return
	case flagType:
		v.SetBool(true)
		return
	case bigIntType:
		parsedInt := parseBigInt(innerBytes)
		v.Set(reflect.ValueOf(parsedInt))
		return
	***REMOVED***
	switch val := v; val.Kind() ***REMOVED***
	case reflect.Bool:
		parsedBool, err1 := parseBool(innerBytes)
		if err1 == nil ***REMOVED***
			val.SetBool(parsedBool)
		***REMOVED***
		err = err1
		return
	case reflect.Int, reflect.Int32, reflect.Int64:
		if val.Type().Size() == 4 ***REMOVED***
			parsedInt, err1 := parseInt32(innerBytes)
			if err1 == nil ***REMOVED***
				val.SetInt(int64(parsedInt))
			***REMOVED***
			err = err1
		***REMOVED*** else ***REMOVED***
			parsedInt, err1 := parseInt64(innerBytes)
			if err1 == nil ***REMOVED***
				val.SetInt(parsedInt)
			***REMOVED***
			err = err1
		***REMOVED***
		return
	// TODO(dfc) Add support for the remaining integer types
	case reflect.Struct:
		structType := fieldType

		if structType.NumField() > 0 &&
			structType.Field(0).Type == rawContentsType ***REMOVED***
			bytes := bytes[initOffset:offset]
			val.Field(0).Set(reflect.ValueOf(RawContent(bytes)))
		***REMOVED***

		innerOffset := 0
		for i := 0; i < structType.NumField(); i++ ***REMOVED***
			field := structType.Field(i)
			if i == 0 && field.Type == rawContentsType ***REMOVED***
				continue
			***REMOVED***
			innerOffset, err = parseField(val.Field(i), innerBytes, innerOffset, parseFieldParameters(field.Tag.Get("asn1")))
			if err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		// We allow extra bytes at the end of the SEQUENCE because
		// adding elements to the end has been used in X.509 as the
		// version numbers have increased.
		return
	case reflect.Slice:
		sliceType := fieldType
		if sliceType.Elem().Kind() == reflect.Uint8 ***REMOVED***
			val.Set(reflect.MakeSlice(sliceType, len(innerBytes), len(innerBytes)))
			reflect.Copy(val, reflect.ValueOf(innerBytes))
			return
		***REMOVED***
		newSlice, err1 := parseSequenceOf(innerBytes, sliceType, sliceType.Elem())
		if err1 == nil ***REMOVED***
			val.Set(newSlice)
		***REMOVED***
		err = err1
		return
	case reflect.String:
		var v string
		switch universalTag ***REMOVED***
		case tagPrintableString:
			v, err = parsePrintableString(innerBytes)
		case tagIA5String:
			v, err = parseIA5String(innerBytes)
		case tagT61String:
			v, err = parseT61String(innerBytes)
		case tagUTF8String:
			v, err = parseUTF8String(innerBytes)
		case tagGeneralString:
			// GeneralString is specified in ISO-2022/ECMA-35,
			// A brief review suggests that it includes structures
			// that allow the encoding to change midstring and
			// such. We give up and pass it as an 8-bit string.
			v, err = parseT61String(innerBytes)
		default:
			err = SyntaxError***REMOVED***fmt.Sprintf("internal error: unknown string type %d", universalTag)***REMOVED***
		***REMOVED***
		if err == nil ***REMOVED***
			val.SetString(v)
		***REMOVED***
		return
	***REMOVED***
	err = StructuralError***REMOVED***"unsupported: " + v.Type().String()***REMOVED***
	return
***REMOVED***

// setDefaultValue is used to install a default value, from a tag string, into
// a Value. It is successful is the field was optional, even if a default value
// wasn't provided or it failed to install it into the Value.
func setDefaultValue(v reflect.Value, params fieldParameters) (ok bool) ***REMOVED***
	if !params.optional ***REMOVED***
		return
	***REMOVED***
	ok = true
	if params.defaultValue == nil ***REMOVED***
		return
	***REMOVED***
	switch val := v; val.Kind() ***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val.SetInt(*params.defaultValue)
	***REMOVED***
	return
***REMOVED***

// Unmarshal parses the DER-encoded ASN.1 data structure b
// and uses the reflect package to fill in an arbitrary value pointed at by val.
// Because Unmarshal uses the reflect package, the structs
// being written to must use upper case field names.
//
// An ASN.1 INTEGER can be written to an int, int32, int64,
// or *big.Int (from the math/big package).
// If the encoded value does not fit in the Go type,
// Unmarshal returns a parse error.
//
// An ASN.1 BIT STRING can be written to a BitString.
//
// An ASN.1 OCTET STRING can be written to a []byte.
//
// An ASN.1 OBJECT IDENTIFIER can be written to an
// ObjectIdentifier.
//
// An ASN.1 ENUMERATED can be written to an Enumerated.
//
// An ASN.1 UTCTIME or GENERALIZEDTIME can be written to a time.Time.
//
// An ASN.1 PrintableString or IA5String can be written to a string.
//
// Any of the above ASN.1 values can be written to an interface***REMOVED******REMOVED***.
// The value stored in the interface has the corresponding Go type.
// For integers, that type is int64.
//
// An ASN.1 SEQUENCE OF x or SET OF x can be written
// to a slice if an x can be written to the slice's element type.
//
// An ASN.1 SEQUENCE or SET can be written to a struct
// if each of the elements in the sequence can be
// written to the corresponding element in the struct.
//
// The following tags on struct fields have special meaning to Unmarshal:
//
//	optional		marks the field as ASN.1 OPTIONAL
//	[explicit] tag:x	specifies the ASN.1 tag number; implies ASN.1 CONTEXT SPECIFIC
//	default:x		sets the default value for optional integer fields
//
// If the type of the first field of a structure is RawContent then the raw
// ASN1 contents of the struct will be stored in it.
//
// Other ASN.1 types are not supported; if it encounters them,
// Unmarshal returns a parse error.
func Unmarshal(b []byte, val interface***REMOVED******REMOVED***) (rest []byte, err error) ***REMOVED***
	return UnmarshalWithParams(b, val, "")
***REMOVED***

// UnmarshalWithParams allows field parameters to be specified for the
// top-level element. The form of the params is the same as the field tags.
func UnmarshalWithParams(b []byte, val interface***REMOVED******REMOVED***, params string) (rest []byte, err error) ***REMOVED***
	v := reflect.ValueOf(val).Elem()
	offset, err := parseField(v, b, 0, parseFieldParameters(params))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b[offset:], nil
***REMOVED***
