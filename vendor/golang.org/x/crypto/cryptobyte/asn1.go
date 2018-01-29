// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cryptobyte

import (
	encoding_asn1 "encoding/asn1"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"golang.org/x/crypto/cryptobyte/asn1"
)

// This file contains ASN.1-related methods for String and Builder.

// Builder

// AddASN1Int64 appends a DER-encoded ASN.1 INTEGER.
func (b *Builder) AddASN1Int64(v int64) ***REMOVED***
	b.addASN1Signed(asn1.INTEGER, v)
***REMOVED***

// AddASN1Enum appends a DER-encoded ASN.1 ENUMERATION.
func (b *Builder) AddASN1Enum(v int64) ***REMOVED***
	b.addASN1Signed(asn1.ENUM, v)
***REMOVED***

func (b *Builder) addASN1Signed(tag asn1.Tag, v int64) ***REMOVED***
	b.AddASN1(tag, func(c *Builder) ***REMOVED***
		length := 1
		for i := v; i >= 0x80 || i < -0x80; i >>= 8 ***REMOVED***
			length++
		***REMOVED***

		for ; length > 0; length-- ***REMOVED***
			i := v >> uint((length-1)*8) & 0xff
			c.AddUint8(uint8(i))
		***REMOVED***
	***REMOVED***)
***REMOVED***

// AddASN1Uint64 appends a DER-encoded ASN.1 INTEGER.
func (b *Builder) AddASN1Uint64(v uint64) ***REMOVED***
	b.AddASN1(asn1.INTEGER, func(c *Builder) ***REMOVED***
		length := 1
		for i := v; i >= 0x80; i >>= 8 ***REMOVED***
			length++
		***REMOVED***

		for ; length > 0; length-- ***REMOVED***
			i := v >> uint((length-1)*8) & 0xff
			c.AddUint8(uint8(i))
		***REMOVED***
	***REMOVED***)
***REMOVED***

// AddASN1BigInt appends a DER-encoded ASN.1 INTEGER.
func (b *Builder) AddASN1BigInt(n *big.Int) ***REMOVED***
	if b.err != nil ***REMOVED***
		return
	***REMOVED***

	b.AddASN1(asn1.INTEGER, func(c *Builder) ***REMOVED***
		if n.Sign() < 0 ***REMOVED***
			// A negative number has to be converted to two's-complement form. So we
			// invert and subtract 1. If the most-significant-bit isn't set then
			// we'll need to pad the beginning with 0xff in order to keep the number
			// negative.
			nMinus1 := new(big.Int).Neg(n)
			nMinus1.Sub(nMinus1, bigOne)
			bytes := nMinus1.Bytes()
			for i := range bytes ***REMOVED***
				bytes[i] ^= 0xff
			***REMOVED***
			if bytes[0]&0x80 == 0 ***REMOVED***
				c.add(0xff)
			***REMOVED***
			c.add(bytes...)
		***REMOVED*** else if n.Sign() == 0 ***REMOVED***
			c.add(0)
		***REMOVED*** else ***REMOVED***
			bytes := n.Bytes()
			if bytes[0]&0x80 != 0 ***REMOVED***
				c.add(0)
			***REMOVED***
			c.add(bytes...)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// AddASN1OctetString appends a DER-encoded ASN.1 OCTET STRING.
func (b *Builder) AddASN1OctetString(bytes []byte) ***REMOVED***
	b.AddASN1(asn1.OCTET_STRING, func(c *Builder) ***REMOVED***
		c.AddBytes(bytes)
	***REMOVED***)
***REMOVED***

const generalizedTimeFormatStr = "20060102150405Z0700"

// AddASN1GeneralizedTime appends a DER-encoded ASN.1 GENERALIZEDTIME.
func (b *Builder) AddASN1GeneralizedTime(t time.Time) ***REMOVED***
	if t.Year() < 0 || t.Year() > 9999 ***REMOVED***
		b.err = fmt.Errorf("cryptobyte: cannot represent %v as a GeneralizedTime", t)
		return
	***REMOVED***
	b.AddASN1(asn1.GeneralizedTime, func(c *Builder) ***REMOVED***
		c.AddBytes([]byte(t.Format(generalizedTimeFormatStr)))
	***REMOVED***)
***REMOVED***

// AddASN1BitString appends a DER-encoded ASN.1 BIT STRING. This does not
// support BIT STRINGs that are not a whole number of bytes.
func (b *Builder) AddASN1BitString(data []byte) ***REMOVED***
	b.AddASN1(asn1.BIT_STRING, func(b *Builder) ***REMOVED***
		b.AddUint8(0)
		b.AddBytes(data)
	***REMOVED***)
***REMOVED***

func (b *Builder) addBase128Int(n int64) ***REMOVED***
	var length int
	if n == 0 ***REMOVED***
		length = 1
	***REMOVED*** else ***REMOVED***
		for i := n; i > 0; i >>= 7 ***REMOVED***
			length++
		***REMOVED***
	***REMOVED***

	for i := length - 1; i >= 0; i-- ***REMOVED***
		o := byte(n >> uint(i*7))
		o &= 0x7f
		if i != 0 ***REMOVED***
			o |= 0x80
		***REMOVED***

		b.add(o)
	***REMOVED***
***REMOVED***

func isValidOID(oid encoding_asn1.ObjectIdentifier) bool ***REMOVED***
	if len(oid) < 2 ***REMOVED***
		return false
	***REMOVED***

	if oid[0] > 2 || (oid[0] <= 1 && oid[1] >= 40) ***REMOVED***
		return false
	***REMOVED***

	for _, v := range oid ***REMOVED***
		if v < 0 ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

func (b *Builder) AddASN1ObjectIdentifier(oid encoding_asn1.ObjectIdentifier) ***REMOVED***
	b.AddASN1(asn1.OBJECT_IDENTIFIER, func(b *Builder) ***REMOVED***
		if !isValidOID(oid) ***REMOVED***
			b.err = fmt.Errorf("cryptobyte: invalid OID: %v", oid)
			return
		***REMOVED***

		b.addBase128Int(int64(oid[0])*40 + int64(oid[1]))
		for _, v := range oid[2:] ***REMOVED***
			b.addBase128Int(int64(v))
		***REMOVED***
	***REMOVED***)
***REMOVED***

func (b *Builder) AddASN1Boolean(v bool) ***REMOVED***
	b.AddASN1(asn1.BOOLEAN, func(b *Builder) ***REMOVED***
		if v ***REMOVED***
			b.AddUint8(0xff)
		***REMOVED*** else ***REMOVED***
			b.AddUint8(0)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func (b *Builder) AddASN1NULL() ***REMOVED***
	b.add(uint8(asn1.NULL), 0)
***REMOVED***

// MarshalASN1 calls encoding_asn1.Marshal on its input and appends the result if
// successful or records an error if one occurred.
func (b *Builder) MarshalASN1(v interface***REMOVED******REMOVED***) ***REMOVED***
	// NOTE(martinkr): This is somewhat of a hack to allow propagation of
	// encoding_asn1.Marshal errors into Builder.err. N.B. if you call MarshalASN1 with a
	// value embedded into a struct, its tag information is lost.
	if b.err != nil ***REMOVED***
		return
	***REMOVED***
	bytes, err := encoding_asn1.Marshal(v)
	if err != nil ***REMOVED***
		b.err = err
		return
	***REMOVED***
	b.AddBytes(bytes)
***REMOVED***

// AddASN1 appends an ASN.1 object. The object is prefixed with the given tag.
// Tags greater than 30 are not supported and result in an error (i.e.
// low-tag-number form only). The child builder passed to the
// BuilderContinuation can be used to build the content of the ASN.1 object.
func (b *Builder) AddASN1(tag asn1.Tag, f BuilderContinuation) ***REMOVED***
	if b.err != nil ***REMOVED***
		return
	***REMOVED***
	// Identifiers with the low five bits set indicate high-tag-number format
	// (two or more octets), which we don't support.
	if tag&0x1f == 0x1f ***REMOVED***
		b.err = fmt.Errorf("cryptobyte: high-tag number identifier octects not supported: 0x%x", tag)
		return
	***REMOVED***
	b.AddUint8(uint8(tag))
	b.addLengthPrefixed(1, true, f)
***REMOVED***

// String

func (s *String) ReadASN1Boolean(out *bool) bool ***REMOVED***
	var bytes String
	if !s.ReadASN1(&bytes, asn1.INTEGER) || len(bytes) != 1 ***REMOVED***
		return false
	***REMOVED***

	switch bytes[0] ***REMOVED***
	case 0:
		*out = false
	case 0xff:
		*out = true
	default:
		return false
	***REMOVED***

	return true
***REMOVED***

var bigIntType = reflect.TypeOf((*big.Int)(nil)).Elem()

// ReadASN1Integer decodes an ASN.1 INTEGER into out and advances. If out does
// not point to an integer or to a big.Int, it panics. It returns true on
// success and false on error.
func (s *String) ReadASN1Integer(out interface***REMOVED******REMOVED***) bool ***REMOVED***
	if reflect.TypeOf(out).Kind() != reflect.Ptr ***REMOVED***
		panic("out is not a pointer")
	***REMOVED***
	switch reflect.ValueOf(out).Elem().Kind() ***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		if !s.readASN1Int64(&i) || reflect.ValueOf(out).Elem().OverflowInt(i) ***REMOVED***
			return false
		***REMOVED***
		reflect.ValueOf(out).Elem().SetInt(i)
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var u uint64
		if !s.readASN1Uint64(&u) || reflect.ValueOf(out).Elem().OverflowUint(u) ***REMOVED***
			return false
		***REMOVED***
		reflect.ValueOf(out).Elem().SetUint(u)
		return true
	case reflect.Struct:
		if reflect.TypeOf(out).Elem() == bigIntType ***REMOVED***
			return s.readASN1BigInt(out.(*big.Int))
		***REMOVED***
	***REMOVED***
	panic("out does not point to an integer type")
***REMOVED***

func checkASN1Integer(bytes []byte) bool ***REMOVED***
	if len(bytes) == 0 ***REMOVED***
		// An INTEGER is encoded with at least one octet.
		return false
	***REMOVED***
	if len(bytes) == 1 ***REMOVED***
		return true
	***REMOVED***
	if bytes[0] == 0 && bytes[1]&0x80 == 0 || bytes[0] == 0xff && bytes[1]&0x80 == 0x80 ***REMOVED***
		// Value is not minimally encoded.
		return false
	***REMOVED***
	return true
***REMOVED***

var bigOne = big.NewInt(1)

func (s *String) readASN1BigInt(out *big.Int) bool ***REMOVED***
	var bytes String
	if !s.ReadASN1(&bytes, asn1.INTEGER) || !checkASN1Integer(bytes) ***REMOVED***
		return false
	***REMOVED***
	if bytes[0]&0x80 == 0x80 ***REMOVED***
		// Negative number.
		neg := make([]byte, len(bytes))
		for i, b := range bytes ***REMOVED***
			neg[i] = ^b
		***REMOVED***
		out.SetBytes(neg)
		out.Add(out, bigOne)
		out.Neg(out)
	***REMOVED*** else ***REMOVED***
		out.SetBytes(bytes)
	***REMOVED***
	return true
***REMOVED***

func (s *String) readASN1Int64(out *int64) bool ***REMOVED***
	var bytes String
	if !s.ReadASN1(&bytes, asn1.INTEGER) || !checkASN1Integer(bytes) || !asn1Signed(out, bytes) ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func asn1Signed(out *int64, n []byte) bool ***REMOVED***
	length := len(n)
	if length > 8 ***REMOVED***
		return false
	***REMOVED***
	for i := 0; i < length; i++ ***REMOVED***
		*out <<= 8
		*out |= int64(n[i])
	***REMOVED***
	// Shift up and down in order to sign extend the result.
	*out <<= 64 - uint8(length)*8
	*out >>= 64 - uint8(length)*8
	return true
***REMOVED***

func (s *String) readASN1Uint64(out *uint64) bool ***REMOVED***
	var bytes String
	if !s.ReadASN1(&bytes, asn1.INTEGER) || !checkASN1Integer(bytes) || !asn1Unsigned(out, bytes) ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func asn1Unsigned(out *uint64, n []byte) bool ***REMOVED***
	length := len(n)
	if length > 9 || length == 9 && n[0] != 0 ***REMOVED***
		// Too large for uint64.
		return false
	***REMOVED***
	if n[0]&0x80 != 0 ***REMOVED***
		// Negative number.
		return false
	***REMOVED***
	for i := 0; i < length; i++ ***REMOVED***
		*out <<= 8
		*out |= uint64(n[i])
	***REMOVED***
	return true
***REMOVED***

// ReadASN1Enum decodes an ASN.1 ENUMERATION into out and advances. It returns
// true on success and false on error.
func (s *String) ReadASN1Enum(out *int) bool ***REMOVED***
	var bytes String
	var i int64
	if !s.ReadASN1(&bytes, asn1.ENUM) || !checkASN1Integer(bytes) || !asn1Signed(&i, bytes) ***REMOVED***
		return false
	***REMOVED***
	if int64(int(i)) != i ***REMOVED***
		return false
	***REMOVED***
	*out = int(i)
	return true
***REMOVED***

func (s *String) readBase128Int(out *int) bool ***REMOVED***
	ret := 0
	for i := 0; len(*s) > 0; i++ ***REMOVED***
		if i == 4 ***REMOVED***
			return false
		***REMOVED***
		ret <<= 7
		b := s.read(1)[0]
		ret |= int(b & 0x7f)
		if b&0x80 == 0 ***REMOVED***
			*out = ret
			return true
		***REMOVED***
	***REMOVED***
	return false // truncated
***REMOVED***

// ReadASN1ObjectIdentifier decodes an ASN.1 OBJECT IDENTIFIER into out and
// advances. It returns true on success and false on error.
func (s *String) ReadASN1ObjectIdentifier(out *encoding_asn1.ObjectIdentifier) bool ***REMOVED***
	var bytes String
	if !s.ReadASN1(&bytes, asn1.OBJECT_IDENTIFIER) || len(bytes) == 0 ***REMOVED***
		return false
	***REMOVED***

	// In the worst case, we get two elements from the first byte (which is
	// encoded differently) and then every varint is a single byte long.
	components := make([]int, len(bytes)+1)

	// The first varint is 40*value1 + value2:
	// According to this packing, value1 can take the values 0, 1 and 2 only.
	// When value1 = 0 or value1 = 1, then value2 is <= 39. When value1 = 2,
	// then there are no restrictions on value2.
	var v int
	if !bytes.readBase128Int(&v) ***REMOVED***
		return false
	***REMOVED***
	if v < 80 ***REMOVED***
		components[0] = v / 40
		components[1] = v % 40
	***REMOVED*** else ***REMOVED***
		components[0] = 2
		components[1] = v - 80
	***REMOVED***

	i := 2
	for ; len(bytes) > 0; i++ ***REMOVED***
		if !bytes.readBase128Int(&v) ***REMOVED***
			return false
		***REMOVED***
		components[i] = v
	***REMOVED***
	*out = components[:i]
	return true
***REMOVED***

// ReadASN1GeneralizedTime decodes an ASN.1 GENERALIZEDTIME into out and
// advances. It returns true on success and false on error.
func (s *String) ReadASN1GeneralizedTime(out *time.Time) bool ***REMOVED***
	var bytes String
	if !s.ReadASN1(&bytes, asn1.GeneralizedTime) ***REMOVED***
		return false
	***REMOVED***
	t := string(bytes)
	res, err := time.Parse(generalizedTimeFormatStr, t)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	if serialized := res.Format(generalizedTimeFormatStr); serialized != t ***REMOVED***
		return false
	***REMOVED***
	*out = res
	return true
***REMOVED***

// ReadASN1BitString decodes an ASN.1 BIT STRING into out and advances. It
// returns true on success and false on error.
func (s *String) ReadASN1BitString(out *encoding_asn1.BitString) bool ***REMOVED***
	var bytes String
	if !s.ReadASN1(&bytes, asn1.BIT_STRING) || len(bytes) == 0 ***REMOVED***
		return false
	***REMOVED***

	paddingBits := uint8(bytes[0])
	bytes = bytes[1:]
	if paddingBits > 7 ||
		len(bytes) == 0 && paddingBits != 0 ||
		len(bytes) > 0 && bytes[len(bytes)-1]&(1<<paddingBits-1) != 0 ***REMOVED***
		return false
	***REMOVED***

	out.BitLength = len(bytes)*8 - int(paddingBits)
	out.Bytes = bytes
	return true
***REMOVED***

// ReadASN1BitString decodes an ASN.1 BIT STRING into out and advances. It is
// an error if the BIT STRING is not a whole number of bytes. This function
// returns true on success and false on error.
func (s *String) ReadASN1BitStringAsBytes(out *[]byte) bool ***REMOVED***
	var bytes String
	if !s.ReadASN1(&bytes, asn1.BIT_STRING) || len(bytes) == 0 ***REMOVED***
		return false
	***REMOVED***

	paddingBits := uint8(bytes[0])
	if paddingBits != 0 ***REMOVED***
		return false
	***REMOVED***
	*out = bytes[1:]
	return true
***REMOVED***

// ReadASN1Bytes reads the contents of a DER-encoded ASN.1 element (not including
// tag and length bytes) into out, and advances. The element must match the
// given tag. It returns true on success and false on error.
func (s *String) ReadASN1Bytes(out *[]byte, tag asn1.Tag) bool ***REMOVED***
	return s.ReadASN1((*String)(out), tag)
***REMOVED***

// ReadASN1 reads the contents of a DER-encoded ASN.1 element (not including
// tag and length bytes) into out, and advances. The element must match the
// given tag. It returns true on success and false on error.
//
// Tags greater than 30 are not supported (i.e. low-tag-number format only).
func (s *String) ReadASN1(out *String, tag asn1.Tag) bool ***REMOVED***
	var t asn1.Tag
	if !s.ReadAnyASN1(out, &t) || t != tag ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// ReadASN1Element reads the contents of a DER-encoded ASN.1 element (including
// tag and length bytes) into out, and advances. The element must match the
// given tag. It returns true on success and false on error.
//
// Tags greater than 30 are not supported (i.e. low-tag-number format only).
func (s *String) ReadASN1Element(out *String, tag asn1.Tag) bool ***REMOVED***
	var t asn1.Tag
	if !s.ReadAnyASN1Element(out, &t) || t != tag ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// ReadAnyASN1 reads the contents of a DER-encoded ASN.1 element (not including
// tag and length bytes) into out, sets outTag to its tag, and advances. It
// returns true on success and false on error.
//
// Tags greater than 30 are not supported (i.e. low-tag-number format only).
func (s *String) ReadAnyASN1(out *String, outTag *asn1.Tag) bool ***REMOVED***
	return s.readASN1(out, outTag, true /* skip header */)
***REMOVED***

// ReadAnyASN1Element reads the contents of a DER-encoded ASN.1 element
// (including tag and length bytes) into out, sets outTag to is tag, and
// advances. It returns true on success and false on error.
//
// Tags greater than 30 are not supported (i.e. low-tag-number format only).
func (s *String) ReadAnyASN1Element(out *String, outTag *asn1.Tag) bool ***REMOVED***
	return s.readASN1(out, outTag, false /* include header */)
***REMOVED***

// PeekASN1Tag returns true if the next ASN.1 value on the string starts with
// the given tag.
func (s String) PeekASN1Tag(tag asn1.Tag) bool ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return false
	***REMOVED***
	return asn1.Tag(s[0]) == tag
***REMOVED***

// SkipASN1 reads and discards an ASN.1 element with the given tag.
func (s *String) SkipASN1(tag asn1.Tag) bool ***REMOVED***
	var unused String
	return s.ReadASN1(&unused, tag)
***REMOVED***

// ReadOptionalASN1 attempts to read the contents of a DER-encoded ASN.1
// element (not including tag and length bytes) tagged with the given tag into
// out. It stores whether an element with the tag was found in outPresent,
// unless outPresent is nil. It returns true on success and false on error.
func (s *String) ReadOptionalASN1(out *String, outPresent *bool, tag asn1.Tag) bool ***REMOVED***
	present := s.PeekASN1Tag(tag)
	if outPresent != nil ***REMOVED***
		*outPresent = present
	***REMOVED***
	if present && !s.ReadASN1(out, tag) ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// SkipOptionalASN1 advances s over an ASN.1 element with the given tag, or
// else leaves s unchanged.
func (s *String) SkipOptionalASN1(tag asn1.Tag) bool ***REMOVED***
	if !s.PeekASN1Tag(tag) ***REMOVED***
		return true
	***REMOVED***
	var unused String
	return s.ReadASN1(&unused, tag)
***REMOVED***

// ReadOptionalASN1Integer attempts to read an optional ASN.1 INTEGER
// explicitly tagged with tag into out and advances. If no element with a
// matching tag is present, it writes defaultValue into out instead. If out
// does not point to an integer or to a big.Int, it panics. It returns true on
// success and false on error.
func (s *String) ReadOptionalASN1Integer(out interface***REMOVED******REMOVED***, tag asn1.Tag, defaultValue interface***REMOVED******REMOVED***) bool ***REMOVED***
	if reflect.TypeOf(out).Kind() != reflect.Ptr ***REMOVED***
		panic("out is not a pointer")
	***REMOVED***
	var present bool
	var i String
	if !s.ReadOptionalASN1(&i, &present, tag) ***REMOVED***
		return false
	***REMOVED***
	if !present ***REMOVED***
		switch reflect.ValueOf(out).Elem().Kind() ***REMOVED***
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(defaultValue))
		case reflect.Struct:
			if reflect.TypeOf(out).Elem() != bigIntType ***REMOVED***
				panic("invalid integer type")
			***REMOVED***
			if reflect.TypeOf(defaultValue).Kind() != reflect.Ptr ||
				reflect.TypeOf(defaultValue).Elem() != bigIntType ***REMOVED***
				panic("out points to big.Int, but defaultValue does not")
			***REMOVED***
			out.(*big.Int).Set(defaultValue.(*big.Int))
		default:
			panic("invalid integer type")
		***REMOVED***
		return true
	***REMOVED***
	if !i.ReadASN1Integer(out) || !i.Empty() ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// ReadOptionalASN1OctetString attempts to read an optional ASN.1 OCTET STRING
// explicitly tagged with tag into out and advances. If no element with a
// matching tag is present, it writes defaultValue into out instead. It returns
// true on success and false on error.
func (s *String) ReadOptionalASN1OctetString(out *[]byte, outPresent *bool, tag asn1.Tag) bool ***REMOVED***
	var present bool
	var child String
	if !s.ReadOptionalASN1(&child, &present, tag) ***REMOVED***
		return false
	***REMOVED***
	if outPresent != nil ***REMOVED***
		*outPresent = present
	***REMOVED***
	if present ***REMOVED***
		var oct String
		if !child.ReadASN1(&oct, asn1.OCTET_STRING) || !child.Empty() ***REMOVED***
			return false
		***REMOVED***
		*out = oct
	***REMOVED*** else ***REMOVED***
		*out = nil
	***REMOVED***
	return true
***REMOVED***

// ReadOptionalASN1Boolean sets *out to the value of the next ASN.1 BOOLEAN or,
// if the next bytes are not an ASN.1 BOOLEAN, to the value of defaultValue.
func (s *String) ReadOptionalASN1Boolean(out *bool, defaultValue bool) bool ***REMOVED***
	var present bool
	var child String
	if !s.ReadOptionalASN1(&child, &present, asn1.BOOLEAN) ***REMOVED***
		return false
	***REMOVED***

	if !present ***REMOVED***
		*out = defaultValue
		return true
	***REMOVED***

	return s.ReadASN1Boolean(out)
***REMOVED***

func (s *String) readASN1(out *String, outTag *asn1.Tag, skipHeader bool) bool ***REMOVED***
	if len(*s) < 2 ***REMOVED***
		return false
	***REMOVED***
	tag, lenByte := (*s)[0], (*s)[1]

	if tag&0x1f == 0x1f ***REMOVED***
		// ITU-T X.690 section 8.1.2
		//
		// An identifier octet with a tag part of 0x1f indicates a high-tag-number
		// form identifier with two or more octets. We only support tags less than
		// 31 (i.e. low-tag-number form, single octet identifier).
		return false
	***REMOVED***

	if outTag != nil ***REMOVED***
		*outTag = asn1.Tag(tag)
	***REMOVED***

	// ITU-T X.690 section 8.1.3
	//
	// Bit 8 of the first length byte indicates whether the length is short- or
	// long-form.
	var length, headerLen uint32 // length includes headerLen
	if lenByte&0x80 == 0 ***REMOVED***
		// Short-form length (section 8.1.3.4), encoded in bits 1-7.
		length = uint32(lenByte) + 2
		headerLen = 2
	***REMOVED*** else ***REMOVED***
		// Long-form length (section 8.1.3.5). Bits 1-7 encode the number of octets
		// used to encode the length.
		lenLen := lenByte & 0x7f
		var len32 uint32

		if lenLen == 0 || lenLen > 4 || len(*s) < int(2+lenLen) ***REMOVED***
			return false
		***REMOVED***

		lenBytes := String((*s)[2 : 2+lenLen])
		if !lenBytes.readUnsigned(&len32, int(lenLen)) ***REMOVED***
			return false
		***REMOVED***

		// ITU-T X.690 section 10.1 (DER length forms) requires encoding the length
		// with the minimum number of octets.
		if len32 < 128 ***REMOVED***
			// Length should have used short-form encoding.
			return false
		***REMOVED***
		if len32>>((lenLen-1)*8) == 0 ***REMOVED***
			// Leading octet is 0. Length should have been at least one byte shorter.
			return false
		***REMOVED***

		headerLen = 2 + uint32(lenLen)
		if headerLen+len32 < len32 ***REMOVED***
			// Overflow.
			return false
		***REMOVED***
		length = headerLen + len32
	***REMOVED***

	if uint32(int(length)) != length || !s.ReadBytes((*[]byte)(out), int(length)) ***REMOVED***
		return false
	***REMOVED***
	if skipHeader && !out.Skip(int(headerLen)) ***REMOVED***
		panic("cryptobyte: internal error")
	***REMOVED***

	return true
***REMOVED***
