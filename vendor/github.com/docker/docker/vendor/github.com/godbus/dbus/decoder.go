package dbus

import (
	"encoding/binary"
	"io"
	"reflect"
)

type decoder struct ***REMOVED***
	in    io.Reader
	order binary.ByteOrder
	pos   int
***REMOVED***

// newDecoder returns a new decoder that reads values from in. The input is
// expected to be in the given byte order.
func newDecoder(in io.Reader, order binary.ByteOrder) *decoder ***REMOVED***
	dec := new(decoder)
	dec.in = in
	dec.order = order
	return dec
***REMOVED***

// align aligns the input to the given boundary and panics on error.
func (dec *decoder) align(n int) ***REMOVED***
	if dec.pos%n != 0 ***REMOVED***
		newpos := (dec.pos + n - 1) & ^(n - 1)
		empty := make([]byte, newpos-dec.pos)
		if _, err := io.ReadFull(dec.in, empty); err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		dec.pos = newpos
	***REMOVED***
***REMOVED***

// Calls binary.Read(dec.in, dec.order, v) and panics on read errors.
func (dec *decoder) binread(v interface***REMOVED******REMOVED***) ***REMOVED***
	if err := binary.Read(dec.in, dec.order, v); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (dec *decoder) Decode(sig Signature) (vs []interface***REMOVED******REMOVED***, err error) ***REMOVED***
	defer func() ***REMOVED***
		var ok bool
		v := recover()
		if err, ok = v.(error); ok ***REMOVED***
			if err == io.EOF || err == io.ErrUnexpectedEOF ***REMOVED***
				err = FormatError("unexpected EOF")
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	vs = make([]interface***REMOVED******REMOVED***, 0)
	s := sig.str
	for s != "" ***REMOVED***
		err, rem := validSingle(s, 0)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		v := dec.decode(s[:len(s)-len(rem)], 0)
		vs = append(vs, v)
		s = rem
	***REMOVED***
	return vs, nil
***REMOVED***

func (dec *decoder) decode(s string, depth int) interface***REMOVED******REMOVED*** ***REMOVED***
	dec.align(alignment(typeFor(s)))
	switch s[0] ***REMOVED***
	case 'y':
		var b [1]byte
		if _, err := dec.in.Read(b[:]); err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		dec.pos++
		return b[0]
	case 'b':
		i := dec.decode("u", depth).(uint32)
		switch ***REMOVED***
		case i == 0:
			return false
		case i == 1:
			return true
		default:
			panic(FormatError("invalid value for boolean"))
		***REMOVED***
	case 'n':
		var i int16
		dec.binread(&i)
		dec.pos += 2
		return i
	case 'i':
		var i int32
		dec.binread(&i)
		dec.pos += 4
		return i
	case 'x':
		var i int64
		dec.binread(&i)
		dec.pos += 8
		return i
	case 'q':
		var i uint16
		dec.binread(&i)
		dec.pos += 2
		return i
	case 'u':
		var i uint32
		dec.binread(&i)
		dec.pos += 4
		return i
	case 't':
		var i uint64
		dec.binread(&i)
		dec.pos += 8
		return i
	case 'd':
		var f float64
		dec.binread(&f)
		dec.pos += 8
		return f
	case 's':
		length := dec.decode("u", depth).(uint32)
		b := make([]byte, int(length)+1)
		if _, err := io.ReadFull(dec.in, b); err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		dec.pos += int(length) + 1
		return string(b[:len(b)-1])
	case 'o':
		return ObjectPath(dec.decode("s", depth).(string))
	case 'g':
		length := dec.decode("y", depth).(byte)
		b := make([]byte, int(length)+1)
		if _, err := io.ReadFull(dec.in, b); err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		dec.pos += int(length) + 1
		sig, err := ParseSignature(string(b[:len(b)-1]))
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		return sig
	case 'v':
		if depth >= 64 ***REMOVED***
			panic(FormatError("input exceeds container depth limit"))
		***REMOVED***
		var variant Variant
		sig := dec.decode("g", depth).(Signature)
		if len(sig.str) == 0 ***REMOVED***
			panic(FormatError("variant signature is empty"))
		***REMOVED***
		err, rem := validSingle(sig.str, 0)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		if rem != "" ***REMOVED***
			panic(FormatError("variant signature has multiple types"))
		***REMOVED***
		variant.sig = sig
		variant.value = dec.decode(sig.str, depth+1)
		return variant
	case 'h':
		return UnixFDIndex(dec.decode("u", depth).(uint32))
	case 'a':
		if len(s) > 1 && s[1] == '***REMOVED***' ***REMOVED***
			ksig := s[2:3]
			vsig := s[3 : len(s)-1]
			v := reflect.MakeMap(reflect.MapOf(typeFor(ksig), typeFor(vsig)))
			if depth >= 63 ***REMOVED***
				panic(FormatError("input exceeds container depth limit"))
			***REMOVED***
			length := dec.decode("u", depth).(uint32)
			// Even for empty maps, the correct padding must be included
			dec.align(8)
			spos := dec.pos
			for dec.pos < spos+int(length) ***REMOVED***
				dec.align(8)
				if !isKeyType(v.Type().Key()) ***REMOVED***
					panic(InvalidTypeError***REMOVED***v.Type()***REMOVED***)
				***REMOVED***
				kv := dec.decode(ksig, depth+2)
				vv := dec.decode(vsig, depth+2)
				v.SetMapIndex(reflect.ValueOf(kv), reflect.ValueOf(vv))
			***REMOVED***
			return v.Interface()
		***REMOVED***
		if depth >= 64 ***REMOVED***
			panic(FormatError("input exceeds container depth limit"))
		***REMOVED***
		length := dec.decode("u", depth).(uint32)
		v := reflect.MakeSlice(reflect.SliceOf(typeFor(s[1:])), 0, int(length))
		// Even for empty arrays, the correct padding must be included
		dec.align(alignment(typeFor(s[1:])))
		spos := dec.pos
		for dec.pos < spos+int(length) ***REMOVED***
			ev := dec.decode(s[1:], depth+1)
			v = reflect.Append(v, reflect.ValueOf(ev))
		***REMOVED***
		return v.Interface()
	case '(':
		if depth >= 64 ***REMOVED***
			panic(FormatError("input exceeds container depth limit"))
		***REMOVED***
		dec.align(8)
		v := make([]interface***REMOVED******REMOVED***, 0)
		s = s[1 : len(s)-1]
		for s != "" ***REMOVED***
			err, rem := validSingle(s, 0)
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***
			ev := dec.decode(s[:len(s)-len(rem)], depth+1)
			v = append(v, ev)
			s = rem
		***REMOVED***
		return v
	default:
		panic(SignatureError***REMOVED***Sig: s***REMOVED***)
	***REMOVED***
***REMOVED***

// A FormatError is an error in the wire format.
type FormatError string

func (e FormatError) Error() string ***REMOVED***
	return "dbus: wire format error: " + string(e)
***REMOVED***
