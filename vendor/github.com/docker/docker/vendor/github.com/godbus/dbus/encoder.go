package dbus

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
)

// An encoder encodes values to the D-Bus wire format.
type encoder struct ***REMOVED***
	out   io.Writer
	order binary.ByteOrder
	pos   int
***REMOVED***

// NewEncoder returns a new encoder that writes to out in the given byte order.
func newEncoder(out io.Writer, order binary.ByteOrder) *encoder ***REMOVED***
	return newEncoderAtOffset(out, 0, order)
***REMOVED***

// newEncoderAtOffset returns a new encoder that writes to out in the given
// byte order. Specify the offset to initialize pos for proper alignment
// computation.
func newEncoderAtOffset(out io.Writer, offset int, order binary.ByteOrder) *encoder ***REMOVED***
	enc := new(encoder)
	enc.out = out
	enc.order = order
	enc.pos = offset
	return enc
***REMOVED***

// Aligns the next output to be on a multiple of n. Panics on write errors.
func (enc *encoder) align(n int) ***REMOVED***
	pad := enc.padding(0, n)
	if pad > 0 ***REMOVED***
		empty := make([]byte, pad)
		if _, err := enc.out.Write(empty); err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		enc.pos += pad
	***REMOVED***
***REMOVED***

// pad returns the number of bytes of padding, based on current position and additional offset.
// and alignment.
func (enc *encoder) padding(offset, algn int) int ***REMOVED***
	abs := enc.pos + offset
	if abs%algn != 0 ***REMOVED***
		newabs := (abs + algn - 1) & ^(algn - 1)
		return newabs - abs
	***REMOVED***
	return 0
***REMOVED***

// Calls binary.Write(enc.out, enc.order, v) and panics on write errors.
func (enc *encoder) binwrite(v interface***REMOVED******REMOVED***) ***REMOVED***
	if err := binary.Write(enc.out, enc.order, v); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

// Encode encodes the given values to the underyling reader. All written values
// are aligned properly as required by the D-Bus spec.
func (enc *encoder) Encode(vs ...interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	defer func() ***REMOVED***
		err, _ = recover().(error)
	***REMOVED***()
	for _, v := range vs ***REMOVED***
		enc.encode(reflect.ValueOf(v), 0)
	***REMOVED***
	return nil
***REMOVED***

// encode encodes the given value to the writer and panics on error. depth holds
// the depth of the container nesting.
func (enc *encoder) encode(v reflect.Value, depth int) ***REMOVED***
	enc.align(alignment(v.Type()))
	switch v.Kind() ***REMOVED***
	case reflect.Uint8:
		var b [1]byte
		b[0] = byte(v.Uint())
		if _, err := enc.out.Write(b[:]); err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		enc.pos++
	case reflect.Bool:
		if v.Bool() ***REMOVED***
			enc.encode(reflect.ValueOf(uint32(1)), depth)
		***REMOVED*** else ***REMOVED***
			enc.encode(reflect.ValueOf(uint32(0)), depth)
		***REMOVED***
	case reflect.Int16:
		enc.binwrite(int16(v.Int()))
		enc.pos += 2
	case reflect.Uint16:
		enc.binwrite(uint16(v.Uint()))
		enc.pos += 2
	case reflect.Int32:
		enc.binwrite(int32(v.Int()))
		enc.pos += 4
	case reflect.Uint32:
		enc.binwrite(uint32(v.Uint()))
		enc.pos += 4
	case reflect.Int64:
		enc.binwrite(v.Int())
		enc.pos += 8
	case reflect.Uint64:
		enc.binwrite(v.Uint())
		enc.pos += 8
	case reflect.Float64:
		enc.binwrite(v.Float())
		enc.pos += 8
	case reflect.String:
		enc.encode(reflect.ValueOf(uint32(len(v.String()))), depth)
		b := make([]byte, v.Len()+1)
		copy(b, v.String())
		b[len(b)-1] = 0
		n, err := enc.out.Write(b)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		enc.pos += n
	case reflect.Ptr:
		enc.encode(v.Elem(), depth)
	case reflect.Slice, reflect.Array:
		if depth >= 64 ***REMOVED***
			panic(FormatError("input exceeds container depth limit"))
		***REMOVED***
		// Lookahead offset: 4 bytes for uint32 length (with alignment),
		// plus alignment for elements.
		n := enc.padding(0, 4) + 4
		offset := enc.pos + n + enc.padding(n, alignment(v.Type().Elem()))

		var buf bytes.Buffer
		bufenc := newEncoderAtOffset(&buf, offset, enc.order)

		for i := 0; i < v.Len(); i++ ***REMOVED***
			bufenc.encode(v.Index(i), depth+1)
		***REMOVED***
		enc.encode(reflect.ValueOf(uint32(buf.Len())), depth)
		length := buf.Len()
		enc.align(alignment(v.Type().Elem()))
		if _, err := buf.WriteTo(enc.out); err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		enc.pos += length
	case reflect.Struct:
		if depth >= 64 && v.Type() != signatureType ***REMOVED***
			panic(FormatError("input exceeds container depth limit"))
		***REMOVED***
		switch t := v.Type(); t ***REMOVED***
		case signatureType:
			str := v.Field(0)
			enc.encode(reflect.ValueOf(byte(str.Len())), depth+1)
			b := make([]byte, str.Len()+1)
			copy(b, str.String())
			b[len(b)-1] = 0
			n, err := enc.out.Write(b)
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***
			enc.pos += n
		case variantType:
			variant := v.Interface().(Variant)
			enc.encode(reflect.ValueOf(variant.sig), depth+1)
			enc.encode(reflect.ValueOf(variant.value), depth+1)
		default:
			for i := 0; i < v.Type().NumField(); i++ ***REMOVED***
				field := t.Field(i)
				if field.PkgPath == "" && field.Tag.Get("dbus") != "-" ***REMOVED***
					enc.encode(v.Field(i), depth+1)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	case reflect.Map:
		// Maps are arrays of structures, so they actually increase the depth by
		// 2.
		if depth >= 63 ***REMOVED***
			panic(FormatError("input exceeds container depth limit"))
		***REMOVED***
		if !isKeyType(v.Type().Key()) ***REMOVED***
			panic(InvalidTypeError***REMOVED***v.Type()***REMOVED***)
		***REMOVED***
		keys := v.MapKeys()
		// Lookahead offset: 4 bytes for uint32 length (with alignment),
		// plus 8-byte alignment
		n := enc.padding(0, 4) + 4
		offset := enc.pos + n + enc.padding(n, 8)

		var buf bytes.Buffer
		bufenc := newEncoderAtOffset(&buf, offset, enc.order)
		for _, k := range keys ***REMOVED***
			bufenc.align(8)
			bufenc.encode(k, depth+2)
			bufenc.encode(v.MapIndex(k), depth+2)
		***REMOVED***
		enc.encode(reflect.ValueOf(uint32(buf.Len())), depth)
		length := buf.Len()
		enc.align(8)
		if _, err := buf.WriteTo(enc.out); err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		enc.pos += length
	default:
		panic(InvalidTypeError***REMOVED***v.Type()***REMOVED***)
	***REMOVED***
***REMOVED***
