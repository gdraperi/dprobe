// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asn1

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"time"
	"unicode/utf8"
)

// A forkableWriter is an in-memory buffer that can be
// 'forked' to create new forkableWriters that bracket the
// original.  After
//    pre, post := w.fork();
// the overall sequence of bytes represented is logically w+pre+post.
type forkableWriter struct ***REMOVED***
	*bytes.Buffer
	pre, post *forkableWriter
***REMOVED***

func newForkableWriter() *forkableWriter ***REMOVED***
	return &forkableWriter***REMOVED***new(bytes.Buffer), nil, nil***REMOVED***
***REMOVED***

func (f *forkableWriter) fork() (pre, post *forkableWriter) ***REMOVED***
	if f.pre != nil || f.post != nil ***REMOVED***
		panic("have already forked")
	***REMOVED***
	f.pre = newForkableWriter()
	f.post = newForkableWriter()
	return f.pre, f.post
***REMOVED***

func (f *forkableWriter) Len() (l int) ***REMOVED***
	l += f.Buffer.Len()
	if f.pre != nil ***REMOVED***
		l += f.pre.Len()
	***REMOVED***
	if f.post != nil ***REMOVED***
		l += f.post.Len()
	***REMOVED***
	return
***REMOVED***

func (f *forkableWriter) writeTo(out io.Writer) (n int, err error) ***REMOVED***
	n, err = out.Write(f.Bytes())
	if err != nil ***REMOVED***
		return
	***REMOVED***

	var nn int

	if f.pre != nil ***REMOVED***
		nn, err = f.pre.writeTo(out)
		n += nn
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	if f.post != nil ***REMOVED***
		nn, err = f.post.writeTo(out)
		n += nn
	***REMOVED***
	return
***REMOVED***

func marshalBase128Int(out *forkableWriter, n int64) (err error) ***REMOVED***
	if n == 0 ***REMOVED***
		err = out.WriteByte(0)
		return
	***REMOVED***

	l := 0
	for i := n; i > 0; i >>= 7 ***REMOVED***
		l++
	***REMOVED***

	for i := l - 1; i >= 0; i-- ***REMOVED***
		o := byte(n >> uint(i*7))
		o &= 0x7f
		if i != 0 ***REMOVED***
			o |= 0x80
		***REMOVED***
		err = out.WriteByte(o)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func marshalInt64(out *forkableWriter, i int64) (err error) ***REMOVED***
	n := int64Length(i)

	for ; n > 0; n-- ***REMOVED***
		err = out.WriteByte(byte(i >> uint((n-1)*8)))
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func int64Length(i int64) (numBytes int) ***REMOVED***
	numBytes = 1

	for i > 127 ***REMOVED***
		numBytes++
		i >>= 8
	***REMOVED***

	for i < -128 ***REMOVED***
		numBytes++
		i >>= 8
	***REMOVED***

	return
***REMOVED***

func marshalBigInt(out *forkableWriter, n *big.Int) (err error) ***REMOVED***
	if n.Sign() < 0 ***REMOVED***
		// A negative number has to be converted to two's-complement
		// form. So we'll subtract 1 and invert. If the
		// most-significant-bit isn't set then we'll need to pad the
		// beginning with 0xff in order to keep the number negative.
		nMinus1 := new(big.Int).Neg(n)
		nMinus1.Sub(nMinus1, bigOne)
		bytes := nMinus1.Bytes()
		for i := range bytes ***REMOVED***
			bytes[i] ^= 0xff
		***REMOVED***
		if len(bytes) == 0 || bytes[0]&0x80 == 0 ***REMOVED***
			err = out.WriteByte(0xff)
			if err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		_, err = out.Write(bytes)
	***REMOVED*** else if n.Sign() == 0 ***REMOVED***
		// Zero is written as a single 0 zero rather than no bytes.
		err = out.WriteByte(0x00)
	***REMOVED*** else ***REMOVED***
		bytes := n.Bytes()
		if len(bytes) > 0 && bytes[0]&0x80 != 0 ***REMOVED***
			// We'll have to pad this with 0x00 in order to stop it
			// looking like a negative number.
			err = out.WriteByte(0)
			if err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		_, err = out.Write(bytes)
	***REMOVED***
	return
***REMOVED***

func marshalLength(out *forkableWriter, i int) (err error) ***REMOVED***
	n := lengthLength(i)

	for ; n > 0; n-- ***REMOVED***
		err = out.WriteByte(byte(i >> uint((n-1)*8)))
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func lengthLength(i int) (numBytes int) ***REMOVED***
	numBytes = 1
	for i > 255 ***REMOVED***
		numBytes++
		i >>= 8
	***REMOVED***
	return
***REMOVED***

func marshalTagAndLength(out *forkableWriter, t tagAndLength) (err error) ***REMOVED***
	b := uint8(t.class) << 6
	if t.isCompound ***REMOVED***
		b |= 0x20
	***REMOVED***
	if t.tag >= 31 ***REMOVED***
		b |= 0x1f
		err = out.WriteByte(b)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		err = marshalBase128Int(out, int64(t.tag))
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		b |= uint8(t.tag)
		err = out.WriteByte(b)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	if t.length >= 128 ***REMOVED***
		l := lengthLength(t.length)
		err = out.WriteByte(0x80 | byte(l))
		if err != nil ***REMOVED***
			return
		***REMOVED***
		err = marshalLength(out, t.length)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		err = out.WriteByte(byte(t.length))
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func marshalBitString(out *forkableWriter, b BitString) (err error) ***REMOVED***
	paddingBits := byte((8 - b.BitLength%8) % 8)
	err = out.WriteByte(paddingBits)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	_, err = out.Write(b.Bytes)
	return
***REMOVED***

func marshalObjectIdentifier(out *forkableWriter, oid []int) (err error) ***REMOVED***
	if len(oid) < 2 || oid[0] > 2 || (oid[0] < 2 && oid[1] >= 40) ***REMOVED***
		return StructuralError***REMOVED***"invalid object identifier"***REMOVED***
	***REMOVED***

	err = marshalBase128Int(out, int64(oid[0]*40+oid[1]))
	if err != nil ***REMOVED***
		return
	***REMOVED***
	for i := 2; i < len(oid); i++ ***REMOVED***
		err = marshalBase128Int(out, int64(oid[i]))
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func marshalPrintableString(out *forkableWriter, s string) (err error) ***REMOVED***
	b := []byte(s)
	for _, c := range b ***REMOVED***
		if !isPrintable(c) ***REMOVED***
			return StructuralError***REMOVED***"PrintableString contains invalid character"***REMOVED***
		***REMOVED***
	***REMOVED***

	_, err = out.Write(b)
	return
***REMOVED***

func marshalIA5String(out *forkableWriter, s string) (err error) ***REMOVED***
	b := []byte(s)
	for _, c := range b ***REMOVED***
		if c > 127 ***REMOVED***
			return StructuralError***REMOVED***"IA5String contains invalid character"***REMOVED***
		***REMOVED***
	***REMOVED***

	_, err = out.Write(b)
	return
***REMOVED***

func marshalUTF8String(out *forkableWriter, s string) (err error) ***REMOVED***
	_, err = out.Write([]byte(s))
	return
***REMOVED***

func marshalTwoDigits(out *forkableWriter, v int) (err error) ***REMOVED***
	err = out.WriteByte(byte('0' + (v/10)%10))
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return out.WriteByte(byte('0' + v%10))
***REMOVED***

func marshalUTCTime(out *forkableWriter, t time.Time) (err error) ***REMOVED***
	year, month, day := t.Date()

	switch ***REMOVED***
	case 1950 <= year && year < 2000:
		err = marshalTwoDigits(out, int(year-1900))
	case 2000 <= year && year < 2050:
		err = marshalTwoDigits(out, int(year-2000))
	default:
		return StructuralError***REMOVED***"cannot represent time as UTCTime"***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return
	***REMOVED***

	err = marshalTwoDigits(out, int(month))
	if err != nil ***REMOVED***
		return
	***REMOVED***

	err = marshalTwoDigits(out, day)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	hour, min, sec := t.Clock()

	err = marshalTwoDigits(out, hour)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	err = marshalTwoDigits(out, min)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	err = marshalTwoDigits(out, sec)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	_, offset := t.Zone()

	switch ***REMOVED***
	case offset/60 == 0:
		err = out.WriteByte('Z')
		return
	case offset > 0:
		err = out.WriteByte('+')
	case offset < 0:
		err = out.WriteByte('-')
	***REMOVED***

	if err != nil ***REMOVED***
		return
	***REMOVED***

	offsetMinutes := offset / 60
	if offsetMinutes < 0 ***REMOVED***
		offsetMinutes = -offsetMinutes
	***REMOVED***

	err = marshalTwoDigits(out, offsetMinutes/60)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	err = marshalTwoDigits(out, offsetMinutes%60)
	return
***REMOVED***

func stripTagAndLength(in []byte) []byte ***REMOVED***
	_, offset, err := parseTagAndLength(in, 0)
	if err != nil ***REMOVED***
		return in
	***REMOVED***
	return in[offset:]
***REMOVED***

func marshalBody(out *forkableWriter, value reflect.Value, params fieldParameters) (err error) ***REMOVED***
	switch value.Type() ***REMOVED***
	case timeType:
		return marshalUTCTime(out, value.Interface().(time.Time))
	case bitStringType:
		return marshalBitString(out, value.Interface().(BitString))
	case objectIdentifierType:
		return marshalObjectIdentifier(out, value.Interface().(ObjectIdentifier))
	case bigIntType:
		return marshalBigInt(out, value.Interface().(*big.Int))
	***REMOVED***

	switch v := value; v.Kind() ***REMOVED***
	case reflect.Bool:
		if v.Bool() ***REMOVED***
			return out.WriteByte(255)
		***REMOVED*** else ***REMOVED***
			return out.WriteByte(0)
		***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return marshalInt64(out, int64(v.Int()))
	case reflect.Struct:
		t := v.Type()

		startingField := 0

		// If the first element of the structure is a non-empty
		// RawContents, then we don't bother serializing the rest.
		if t.NumField() > 0 && t.Field(0).Type == rawContentsType ***REMOVED***
			s := v.Field(0)
			if s.Len() > 0 ***REMOVED***
				bytes := make([]byte, s.Len())
				for i := 0; i < s.Len(); i++ ***REMOVED***
					bytes[i] = uint8(s.Index(i).Uint())
				***REMOVED***
				/* The RawContents will contain the tag and
				 * length fields but we'll also be writing
				 * those ourselves, so we strip them out of
				 * bytes */
				_, err = out.Write(stripTagAndLength(bytes))
				return
			***REMOVED*** else ***REMOVED***
				startingField = 1
			***REMOVED***
		***REMOVED***

		for i := startingField; i < t.NumField(); i++ ***REMOVED***
			var pre *forkableWriter
			pre, out = out.fork()
			err = marshalField(pre, v.Field(i), parseFieldParameters(t.Field(i).Tag.Get("asn1")))
			if err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		return
	case reflect.Slice:
		sliceType := v.Type()
		if sliceType.Elem().Kind() == reflect.Uint8 ***REMOVED***
			bytes := make([]byte, v.Len())
			for i := 0; i < v.Len(); i++ ***REMOVED***
				bytes[i] = uint8(v.Index(i).Uint())
			***REMOVED***
			_, err = out.Write(bytes)
			return
		***REMOVED***

		var fp fieldParameters
		for i := 0; i < v.Len(); i++ ***REMOVED***
			var pre *forkableWriter
			pre, out = out.fork()
			err = marshalField(pre, v.Index(i), fp)
			if err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		return
	case reflect.String:
		switch params.stringType ***REMOVED***
		case tagIA5String:
			return marshalIA5String(out, v.String())
		case tagPrintableString:
			return marshalPrintableString(out, v.String())
		default:
			return marshalUTF8String(out, v.String())
		***REMOVED***
	***REMOVED***

	return StructuralError***REMOVED***"unknown Go type"***REMOVED***
***REMOVED***

func marshalField(out *forkableWriter, v reflect.Value, params fieldParameters) (err error) ***REMOVED***
	// If the field is an interface***REMOVED******REMOVED*** then recurse into it.
	if v.Kind() == reflect.Interface && v.Type().NumMethod() == 0 ***REMOVED***
		return marshalField(out, v.Elem(), params)
	***REMOVED***

	if v.Kind() == reflect.Slice && v.Len() == 0 && params.omitEmpty ***REMOVED***
		return
	***REMOVED***

	if params.optional && reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface()) ***REMOVED***
		return
	***REMOVED***

	if v.Type() == rawValueType ***REMOVED***
		rv := v.Interface().(RawValue)
		if len(rv.FullBytes) != 0 ***REMOVED***
			_, err = out.Write(rv.FullBytes)
		***REMOVED*** else ***REMOVED***
			err = marshalTagAndLength(out, tagAndLength***REMOVED***rv.Class, rv.Tag, len(rv.Bytes), rv.IsCompound***REMOVED***)
			if err != nil ***REMOVED***
				return
			***REMOVED***
			_, err = out.Write(rv.Bytes)
		***REMOVED***
		return
	***REMOVED***

	tag, isCompound, ok := getUniversalType(v.Type())
	if !ok ***REMOVED***
		err = StructuralError***REMOVED***fmt.Sprintf("unknown Go type: %v", v.Type())***REMOVED***
		return
	***REMOVED***
	class := classUniversal

	if params.stringType != 0 && tag != tagPrintableString ***REMOVED***
		return StructuralError***REMOVED***"explicit string type given to non-string member"***REMOVED***
	***REMOVED***

	if tag == tagPrintableString ***REMOVED***
		if params.stringType == 0 ***REMOVED***
			// This is a string without an explicit string type. We'll use
			// a PrintableString if the character set in the string is
			// sufficiently limited, otherwise we'll use a UTF8String.
			for _, r := range v.String() ***REMOVED***
				if r >= utf8.RuneSelf || !isPrintable(byte(r)) ***REMOVED***
					if !utf8.ValidString(v.String()) ***REMOVED***
						return errors.New("asn1: string not valid UTF-8")
					***REMOVED***
					tag = tagUTF8String
					break
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			tag = params.stringType
		***REMOVED***
	***REMOVED***

	if params.set ***REMOVED***
		if tag != tagSequence ***REMOVED***
			return StructuralError***REMOVED***"non sequence tagged as set"***REMOVED***
		***REMOVED***
		tag = tagSet
	***REMOVED***

	tags, body := out.fork()

	err = marshalBody(body, v, params)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	bodyLen := body.Len()

	var explicitTag *forkableWriter
	if params.explicit ***REMOVED***
		explicitTag, tags = tags.fork()
	***REMOVED***

	if !params.explicit && params.tag != nil ***REMOVED***
		// implicit tag.
		tag = *params.tag
		class = classContextSpecific
	***REMOVED***

	err = marshalTagAndLength(tags, tagAndLength***REMOVED***class, tag, bodyLen, isCompound***REMOVED***)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if params.explicit ***REMOVED***
		err = marshalTagAndLength(explicitTag, tagAndLength***REMOVED***
			class:      classContextSpecific,
			tag:        *params.tag,
			length:     bodyLen + tags.Len(),
			isCompound: true,
		***REMOVED***)
	***REMOVED***

	return nil
***REMOVED***

// Marshal returns the ASN.1 encoding of val.
func Marshal(val interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	var out bytes.Buffer
	v := reflect.ValueOf(val)
	f := newForkableWriter()
	err := marshalField(f, v, fieldParameters***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	_, err = f.writeTo(&out)
	return out.Bytes(), nil
***REMOVED***
