// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"io"
	"reflect"
	// "runtime/debug"
)

// Some tagging information for error messages.
const (
	msgTagDec             = "codec.decoder"
	msgBadDesc            = "Unrecognized descriptor byte"
	msgDecCannotExpandArr = "cannot expand go array from %v to stream length: %v"
)

// decReader abstracts the reading source, allowing implementations that can
// read from an io.Reader or directly off a byte slice with zero-copying.
type decReader interface ***REMOVED***
	readn(n int) []byte
	readb([]byte)
	readn1() uint8
	readUint16() uint16
	readUint32() uint32
	readUint64() uint64
***REMOVED***

type decDriver interface ***REMOVED***
	initReadNext()
	tryDecodeAsNil() bool
	currentEncodedType() valueType
	isBuiltinType(rt uintptr) bool
	decodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***)
	//decodeNaked: Numbers are decoded as int64, uint64, float64 only (no smaller sized number types).
	decodeNaked() (v interface***REMOVED******REMOVED***, vt valueType, decodeFurther bool)
	decodeInt(bitsize uint8) (i int64)
	decodeUint(bitsize uint8) (ui uint64)
	decodeFloat(chkOverflow32 bool) (f float64)
	decodeBool() (b bool)
	// decodeString can also decode symbols
	decodeString() (s string)
	decodeBytes(bs []byte) (bsOut []byte, changed bool)
	decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte)
	readMapLen() int
	readArrayLen() int
***REMOVED***

type DecodeOptions struct ***REMOVED***
	// An instance of MapType is used during schema-less decoding of a map in the stream.
	// If nil, we use map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
	MapType reflect.Type
	// An instance of SliceType is used during schema-less decoding of an array in the stream.
	// If nil, we use []interface***REMOVED******REMOVED***
	SliceType reflect.Type
	// ErrorIfNoField controls whether an error is returned when decoding a map
	// from a codec stream into a struct, and no matching struct field is found.
	ErrorIfNoField bool
***REMOVED***

// ------------------------------------

// ioDecReader is a decReader that reads off an io.Reader
type ioDecReader struct ***REMOVED***
	r  io.Reader
	br io.ByteReader
	x  [8]byte //temp byte array re-used internally for efficiency
***REMOVED***

func (z *ioDecReader) readn(n int) (bs []byte) ***REMOVED***
	if n <= 0 ***REMOVED***
		return
	***REMOVED***
	bs = make([]byte, n)
	if _, err := io.ReadAtLeast(z.r, bs, n); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) readb(bs []byte) ***REMOVED***
	if _, err := io.ReadAtLeast(z.r, bs, len(bs)); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (z *ioDecReader) readn1() uint8 ***REMOVED***
	if z.br != nil ***REMOVED***
		b, err := z.br.ReadByte()
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		return b
	***REMOVED***
	z.readb(z.x[:1])
	return z.x[0]
***REMOVED***

func (z *ioDecReader) readUint16() uint16 ***REMOVED***
	z.readb(z.x[:2])
	return bigen.Uint16(z.x[:2])
***REMOVED***

func (z *ioDecReader) readUint32() uint32 ***REMOVED***
	z.readb(z.x[:4])
	return bigen.Uint32(z.x[:4])
***REMOVED***

func (z *ioDecReader) readUint64() uint64 ***REMOVED***
	z.readb(z.x[:8])
	return bigen.Uint64(z.x[:8])
***REMOVED***

// ------------------------------------

// bytesDecReader is a decReader that reads off a byte slice with zero copying
type bytesDecReader struct ***REMOVED***
	b []byte // data
	c int    // cursor
	a int    // available
***REMOVED***

func (z *bytesDecReader) consume(n int) (oldcursor int) ***REMOVED***
	if z.a == 0 ***REMOVED***
		panic(io.EOF)
	***REMOVED***
	if n > z.a ***REMOVED***
		decErr("Trying to read %v bytes. Only %v available", n, z.a)
	***REMOVED***
	// z.checkAvailable(n)
	oldcursor = z.c
	z.c = oldcursor + n
	z.a = z.a - n
	return
***REMOVED***

func (z *bytesDecReader) readn(n int) (bs []byte) ***REMOVED***
	if n <= 0 ***REMOVED***
		return
	***REMOVED***
	c0 := z.consume(n)
	bs = z.b[c0:z.c]
	return
***REMOVED***

func (z *bytesDecReader) readb(bs []byte) ***REMOVED***
	copy(bs, z.readn(len(bs)))
***REMOVED***

func (z *bytesDecReader) readn1() uint8 ***REMOVED***
	c0 := z.consume(1)
	return z.b[c0]
***REMOVED***

// Use binaryEncoding helper for 4 and 8 bits, but inline it for 2 bits
// creating temp slice variable and copying it to helper function is expensive
// for just 2 bits.

func (z *bytesDecReader) readUint16() uint16 ***REMOVED***
	c0 := z.consume(2)
	return uint16(z.b[c0+1]) | uint16(z.b[c0])<<8
***REMOVED***

func (z *bytesDecReader) readUint32() uint32 ***REMOVED***
	c0 := z.consume(4)
	return bigen.Uint32(z.b[c0:z.c])
***REMOVED***

func (z *bytesDecReader) readUint64() uint64 ***REMOVED***
	c0 := z.consume(8)
	return bigen.Uint64(z.b[c0:z.c])
***REMOVED***

// ------------------------------------

// decFnInfo has methods for registering handling decoding of a specific type
// based on some characteristics (builtin, extension, reflect Kind, etc)
type decFnInfo struct ***REMOVED***
	ti    *typeInfo
	d     *Decoder
	dd    decDriver
	xfFn  func(reflect.Value, []byte) error
	xfTag byte
	array bool
***REMOVED***

func (f *decFnInfo) builtin(rv reflect.Value) ***REMOVED***
	f.dd.decodeBuiltin(f.ti.rtid, rv.Addr().Interface())
***REMOVED***

func (f *decFnInfo) rawExt(rv reflect.Value) ***REMOVED***
	xtag, xbs := f.dd.decodeExt(false, 0)
	rv.Field(0).SetUint(uint64(xtag))
	rv.Field(1).SetBytes(xbs)
***REMOVED***

func (f *decFnInfo) ext(rv reflect.Value) ***REMOVED***
	_, xbs := f.dd.decodeExt(true, f.xfTag)
	if fnerr := f.xfFn(rv, xbs); fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

func (f *decFnInfo) binaryMarshal(rv reflect.Value) ***REMOVED***
	var bm binaryUnmarshaler
	if f.ti.unmIndir == -1 ***REMOVED***
		bm = rv.Addr().Interface().(binaryUnmarshaler)
	***REMOVED*** else if f.ti.unmIndir == 0 ***REMOVED***
		bm = rv.Interface().(binaryUnmarshaler)
	***REMOVED*** else ***REMOVED***
		for j, k := int8(0), f.ti.unmIndir; j < k; j++ ***REMOVED***
			if rv.IsNil() ***REMOVED***
				rv.Set(reflect.New(rv.Type().Elem()))
			***REMOVED***
			rv = rv.Elem()
		***REMOVED***
		bm = rv.Interface().(binaryUnmarshaler)
	***REMOVED***
	xbs, _ := f.dd.decodeBytes(nil)
	if fnerr := bm.UnmarshalBinary(xbs); fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

func (f *decFnInfo) kErr(rv reflect.Value) ***REMOVED***
	decErr("Unhandled value for kind: %v: %s", rv.Kind(), msgBadDesc)
***REMOVED***

func (f *decFnInfo) kString(rv reflect.Value) ***REMOVED***
	rv.SetString(f.dd.decodeString())
***REMOVED***

func (f *decFnInfo) kBool(rv reflect.Value) ***REMOVED***
	rv.SetBool(f.dd.decodeBool())
***REMOVED***

func (f *decFnInfo) kInt(rv reflect.Value) ***REMOVED***
	rv.SetInt(f.dd.decodeInt(intBitsize))
***REMOVED***

func (f *decFnInfo) kInt64(rv reflect.Value) ***REMOVED***
	rv.SetInt(f.dd.decodeInt(64))
***REMOVED***

func (f *decFnInfo) kInt32(rv reflect.Value) ***REMOVED***
	rv.SetInt(f.dd.decodeInt(32))
***REMOVED***

func (f *decFnInfo) kInt8(rv reflect.Value) ***REMOVED***
	rv.SetInt(f.dd.decodeInt(8))
***REMOVED***

func (f *decFnInfo) kInt16(rv reflect.Value) ***REMOVED***
	rv.SetInt(f.dd.decodeInt(16))
***REMOVED***

func (f *decFnInfo) kFloat32(rv reflect.Value) ***REMOVED***
	rv.SetFloat(f.dd.decodeFloat(true))
***REMOVED***

func (f *decFnInfo) kFloat64(rv reflect.Value) ***REMOVED***
	rv.SetFloat(f.dd.decodeFloat(false))
***REMOVED***

func (f *decFnInfo) kUint8(rv reflect.Value) ***REMOVED***
	rv.SetUint(f.dd.decodeUint(8))
***REMOVED***

func (f *decFnInfo) kUint64(rv reflect.Value) ***REMOVED***
	rv.SetUint(f.dd.decodeUint(64))
***REMOVED***

func (f *decFnInfo) kUint(rv reflect.Value) ***REMOVED***
	rv.SetUint(f.dd.decodeUint(uintBitsize))
***REMOVED***

func (f *decFnInfo) kUint32(rv reflect.Value) ***REMOVED***
	rv.SetUint(f.dd.decodeUint(32))
***REMOVED***

func (f *decFnInfo) kUint16(rv reflect.Value) ***REMOVED***
	rv.SetUint(f.dd.decodeUint(16))
***REMOVED***

// func (f *decFnInfo) kPtr(rv reflect.Value) ***REMOVED***
// 	debugf(">>>>>>> ??? decode kPtr called - shouldn't get called")
// 	if rv.IsNil() ***REMOVED***
// 		rv.Set(reflect.New(rv.Type().Elem()))
// 	***REMOVED***
// 	f.d.decodeValue(rv.Elem())
// ***REMOVED***

func (f *decFnInfo) kInterface(rv reflect.Value) ***REMOVED***
	// debugf("\t===> kInterface")
	if !rv.IsNil() ***REMOVED***
		f.d.decodeValue(rv.Elem())
		return
	***REMOVED***
	// nil interface:
	// use some hieristics to set the nil interface to an
	// appropriate value based on the first byte read (byte descriptor bd)
	v, vt, decodeFurther := f.dd.decodeNaked()
	if vt == valueTypeNil ***REMOVED***
		return
	***REMOVED***
	// Cannot decode into nil interface with methods (e.g. error, io.Reader, etc)
	// if non-nil value in stream.
	if num := f.ti.rt.NumMethod(); num > 0 ***REMOVED***
		decErr("decodeValue: Cannot decode non-nil codec value into nil %v (%v methods)",
			f.ti.rt, num)
	***REMOVED***
	var rvn reflect.Value
	var useRvn bool
	switch vt ***REMOVED***
	case valueTypeMap:
		if f.d.h.MapType == nil ***REMOVED***
			var m2 map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
			v = &m2
		***REMOVED*** else ***REMOVED***
			rvn = reflect.New(f.d.h.MapType).Elem()
			useRvn = true
		***REMOVED***
	case valueTypeArray:
		if f.d.h.SliceType == nil ***REMOVED***
			var m2 []interface***REMOVED******REMOVED***
			v = &m2
		***REMOVED*** else ***REMOVED***
			rvn = reflect.New(f.d.h.SliceType).Elem()
			useRvn = true
		***REMOVED***
	case valueTypeExt:
		re := v.(*RawExt)
		var bfn func(reflect.Value, []byte) error
		rvn, bfn = f.d.h.getDecodeExtForTag(re.Tag)
		if bfn == nil ***REMOVED***
			rvn = reflect.ValueOf(*re)
		***REMOVED*** else if fnerr := bfn(rvn, re.Data); fnerr != nil ***REMOVED***
			panic(fnerr)
		***REMOVED***
		rv.Set(rvn)
		return
	***REMOVED***
	if decodeFurther ***REMOVED***
		if useRvn ***REMOVED***
			f.d.decodeValue(rvn)
		***REMOVED*** else if v != nil ***REMOVED***
			// this v is a pointer, so we need to dereference it when done
			f.d.decode(v)
			rvn = reflect.ValueOf(v).Elem()
			useRvn = true
		***REMOVED***
	***REMOVED***
	if useRvn ***REMOVED***
		rv.Set(rvn)
	***REMOVED*** else if v != nil ***REMOVED***
		rv.Set(reflect.ValueOf(v))
	***REMOVED***
***REMOVED***

func (f *decFnInfo) kStruct(rv reflect.Value) ***REMOVED***
	fti := f.ti
	if currEncodedType := f.dd.currentEncodedType(); currEncodedType == valueTypeMap ***REMOVED***
		containerLen := f.dd.readMapLen()
		if containerLen == 0 ***REMOVED***
			return
		***REMOVED***
		tisfi := fti.sfi
		for j := 0; j < containerLen; j++ ***REMOVED***
			// var rvkencname string
			// ddecode(&rvkencname)
			f.dd.initReadNext()
			rvkencname := f.dd.decodeString()
			// rvksi := ti.getForEncName(rvkencname)
			if k := fti.indexForEncName(rvkencname); k > -1 ***REMOVED***
				sfik := tisfi[k]
				if sfik.i != -1 ***REMOVED***
					f.d.decodeValue(rv.Field(int(sfik.i)))
				***REMOVED*** else ***REMOVED***
					f.d.decEmbeddedField(rv, sfik.is)
				***REMOVED***
				// f.d.decodeValue(ti.field(k, rv))
			***REMOVED*** else ***REMOVED***
				if f.d.h.ErrorIfNoField ***REMOVED***
					decErr("No matching struct field found when decoding stream map with key: %v",
						rvkencname)
				***REMOVED*** else ***REMOVED***
					var nilintf0 interface***REMOVED******REMOVED***
					f.d.decodeValue(reflect.ValueOf(&nilintf0).Elem())
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if currEncodedType == valueTypeArray ***REMOVED***
		containerLen := f.dd.readArrayLen()
		if containerLen == 0 ***REMOVED***
			return
		***REMOVED***
		for j, si := range fti.sfip ***REMOVED***
			if j == containerLen ***REMOVED***
				break
			***REMOVED***
			if si.i != -1 ***REMOVED***
				f.d.decodeValue(rv.Field(int(si.i)))
			***REMOVED*** else ***REMOVED***
				f.d.decEmbeddedField(rv, si.is)
			***REMOVED***
		***REMOVED***
		if containerLen > len(fti.sfip) ***REMOVED***
			// read remaining values and throw away
			for j := len(fti.sfip); j < containerLen; j++ ***REMOVED***
				var nilintf0 interface***REMOVED******REMOVED***
				f.d.decodeValue(reflect.ValueOf(&nilintf0).Elem())
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		decErr("Only encoded map or array can be decoded into a struct. (valueType: %x)",
			currEncodedType)
	***REMOVED***
***REMOVED***

func (f *decFnInfo) kSlice(rv reflect.Value) ***REMOVED***
	// A slice can be set from a map or array in stream.
	currEncodedType := f.dd.currentEncodedType()

	switch currEncodedType ***REMOVED***
	case valueTypeBytes, valueTypeString:
		if f.ti.rtid == uint8SliceTypId || f.ti.rt.Elem().Kind() == reflect.Uint8 ***REMOVED***
			if bs2, changed2 := f.dd.decodeBytes(rv.Bytes()); changed2 ***REMOVED***
				rv.SetBytes(bs2)
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	if shortCircuitReflectToFastPath && rv.CanAddr() ***REMOVED***
		switch f.ti.rtid ***REMOVED***
		case intfSliceTypId:
			f.d.decSliceIntf(rv.Addr().Interface().(*[]interface***REMOVED******REMOVED***), currEncodedType, f.array)
			return
		case uint64SliceTypId:
			f.d.decSliceUint64(rv.Addr().Interface().(*[]uint64), currEncodedType, f.array)
			return
		case int64SliceTypId:
			f.d.decSliceInt64(rv.Addr().Interface().(*[]int64), currEncodedType, f.array)
			return
		case strSliceTypId:
			f.d.decSliceStr(rv.Addr().Interface().(*[]string), currEncodedType, f.array)
			return
		***REMOVED***
	***REMOVED***

	containerLen, containerLenS := decContLens(f.dd, currEncodedType)

	// an array can never return a nil slice. so no need to check f.array here.

	if rv.IsNil() ***REMOVED***
		rv.Set(reflect.MakeSlice(f.ti.rt, containerLenS, containerLenS))
	***REMOVED***

	if containerLen == 0 ***REMOVED***
		return
	***REMOVED***

	if rvcap, rvlen := rv.Len(), rv.Cap(); containerLenS > rvcap ***REMOVED***
		if f.array ***REMOVED*** // !rv.CanSet()
			decErr(msgDecCannotExpandArr, rvcap, containerLenS)
		***REMOVED***
		rvn := reflect.MakeSlice(f.ti.rt, containerLenS, containerLenS)
		if rvlen > 0 ***REMOVED***
			reflect.Copy(rvn, rv)
		***REMOVED***
		rv.Set(rvn)
	***REMOVED*** else if containerLenS > rvlen ***REMOVED***
		rv.SetLen(containerLenS)
	***REMOVED***

	for j := 0; j < containerLenS; j++ ***REMOVED***
		f.d.decodeValue(rv.Index(j))
	***REMOVED***
***REMOVED***

func (f *decFnInfo) kArray(rv reflect.Value) ***REMOVED***
	// f.d.decodeValue(rv.Slice(0, rv.Len()))
	f.kSlice(rv.Slice(0, rv.Len()))
***REMOVED***

func (f *decFnInfo) kMap(rv reflect.Value) ***REMOVED***
	if shortCircuitReflectToFastPath && rv.CanAddr() ***REMOVED***
		switch f.ti.rtid ***REMOVED***
		case mapStrIntfTypId:
			f.d.decMapStrIntf(rv.Addr().Interface().(*map[string]interface***REMOVED******REMOVED***))
			return
		case mapIntfIntfTypId:
			f.d.decMapIntfIntf(rv.Addr().Interface().(*map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***))
			return
		case mapInt64IntfTypId:
			f.d.decMapInt64Intf(rv.Addr().Interface().(*map[int64]interface***REMOVED******REMOVED***))
			return
		case mapUint64IntfTypId:
			f.d.decMapUint64Intf(rv.Addr().Interface().(*map[uint64]interface***REMOVED******REMOVED***))
			return
		***REMOVED***
	***REMOVED***

	containerLen := f.dd.readMapLen()

	if rv.IsNil() ***REMOVED***
		rv.Set(reflect.MakeMap(f.ti.rt))
	***REMOVED***

	if containerLen == 0 ***REMOVED***
		return
	***REMOVED***

	ktype, vtype := f.ti.rt.Key(), f.ti.rt.Elem()
	ktypeId := reflect.ValueOf(ktype).Pointer()
	for j := 0; j < containerLen; j++ ***REMOVED***
		rvk := reflect.New(ktype).Elem()
		f.d.decodeValue(rvk)

		// special case if a byte array.
		// if ktype == intfTyp ***REMOVED***
		if ktypeId == intfTypId ***REMOVED***
			rvk = rvk.Elem()
			if rvk.Type() == uint8SliceTyp ***REMOVED***
				rvk = reflect.ValueOf(string(rvk.Bytes()))
			***REMOVED***
		***REMOVED***
		rvv := rv.MapIndex(rvk)
		if !rvv.IsValid() ***REMOVED***
			rvv = reflect.New(vtype).Elem()
		***REMOVED***

		f.d.decodeValue(rvv)
		rv.SetMapIndex(rvk, rvv)
	***REMOVED***
***REMOVED***

// ----------------------------------------

type decFn struct ***REMOVED***
	i *decFnInfo
	f func(*decFnInfo, reflect.Value)
***REMOVED***

// A Decoder reads and decodes an object from an input stream in the codec format.
type Decoder struct ***REMOVED***
	r decReader
	d decDriver
	h *BasicHandle
	f map[uintptr]decFn
	x []uintptr
	s []decFn
***REMOVED***

// NewDecoder returns a Decoder for decoding a stream of bytes from an io.Reader.
//
// For efficiency, Users are encouraged to pass in a memory buffered writer
// (eg bufio.Reader, bytes.Buffer).
func NewDecoder(r io.Reader, h Handle) *Decoder ***REMOVED***
	z := ioDecReader***REMOVED***
		r: r,
	***REMOVED***
	z.br, _ = r.(io.ByteReader)
	return &Decoder***REMOVED***r: &z, d: h.newDecDriver(&z), h: h.getBasicHandle()***REMOVED***
***REMOVED***

// NewDecoderBytes returns a Decoder which efficiently decodes directly
// from a byte slice with zero copying.
func NewDecoderBytes(in []byte, h Handle) *Decoder ***REMOVED***
	z := bytesDecReader***REMOVED***
		b: in,
		a: len(in),
	***REMOVED***
	return &Decoder***REMOVED***r: &z, d: h.newDecDriver(&z), h: h.getBasicHandle()***REMOVED***
***REMOVED***

// Decode decodes the stream from reader and stores the result in the
// value pointed to by v. v cannot be a nil pointer. v can also be
// a reflect.Value of a pointer.
//
// Note that a pointer to a nil interface is not a nil pointer.
// If you do not know what type of stream it is, pass in a pointer to a nil interface.
// We will decode and store a value in that nil interface.
//
// Sample usages:
//   // Decoding into a non-nil typed value
//   var f float32
//   err = codec.NewDecoder(r, handle).Decode(&f)
//
//   // Decoding into nil interface
//   var v interface***REMOVED******REMOVED***
//   dec := codec.NewDecoder(r, handle)
//   err = dec.Decode(&v)
//
// When decoding into a nil interface***REMOVED******REMOVED***, we will decode into an appropriate value based
// on the contents of the stream:
//   - Numbers are decoded as float64, int64 or uint64.
//   - Other values are decoded appropriately depending on the type:
//     bool, string, []byte, time.Time, etc
//   - Extensions are decoded as RawExt (if no ext function registered for the tag)
// Configurations exist on the Handle to override defaults
// (e.g. for MapType, SliceType and how to decode raw bytes).
//
// When decoding into a non-nil interface***REMOVED******REMOVED*** value, the mode of encoding is based on the
// type of the value. When a value is seen:
//   - If an extension is registered for it, call that extension function
//   - If it implements BinaryUnmarshaler, call its UnmarshalBinary(data []byte) error
//   - Else decode it based on its reflect.Kind
//
// There are some special rules when decoding into containers (slice/array/map/struct).
// Decode will typically use the stream contents to UPDATE the container.
//   - A map can be decoded from a stream map, by updating matching keys.
//   - A slice can be decoded from a stream array,
//     by updating the first n elements, where n is length of the stream.
//   - A slice can be decoded from a stream map, by decoding as if
//     it contains a sequence of key-value pairs.
//   - A struct can be decoded from a stream map, by updating matching fields.
//   - A struct can be decoded from a stream array,
//     by updating fields as they occur in the struct (by index).
//
// When decoding a stream map or array with length of 0 into a nil map or slice,
// we reset the destination map or slice to a zero-length value.
//
// However, when decoding a stream nil, we reset the destination container
// to its "zero" value (e.g. nil for slice/map, etc).
//
func (d *Decoder) Decode(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	defer panicToErr(&err)
	d.decode(v)
	return
***REMOVED***

func (d *Decoder) decode(iv interface***REMOVED******REMOVED***) ***REMOVED***
	d.d.initReadNext()

	switch v := iv.(type) ***REMOVED***
	case nil:
		decErr("Cannot decode into nil.")

	case reflect.Value:
		d.chkPtrValue(v)
		d.decodeValue(v.Elem())

	case *string:
		*v = d.d.decodeString()
	case *bool:
		*v = d.d.decodeBool()
	case *int:
		*v = int(d.d.decodeInt(intBitsize))
	case *int8:
		*v = int8(d.d.decodeInt(8))
	case *int16:
		*v = int16(d.d.decodeInt(16))
	case *int32:
		*v = int32(d.d.decodeInt(32))
	case *int64:
		*v = d.d.decodeInt(64)
	case *uint:
		*v = uint(d.d.decodeUint(uintBitsize))
	case *uint8:
		*v = uint8(d.d.decodeUint(8))
	case *uint16:
		*v = uint16(d.d.decodeUint(16))
	case *uint32:
		*v = uint32(d.d.decodeUint(32))
	case *uint64:
		*v = d.d.decodeUint(64)
	case *float32:
		*v = float32(d.d.decodeFloat(true))
	case *float64:
		*v = d.d.decodeFloat(false)
	case *[]byte:
		*v, _ = d.d.decodeBytes(*v)

	case *[]interface***REMOVED******REMOVED***:
		d.decSliceIntf(v, valueTypeInvalid, false)
	case *[]uint64:
		d.decSliceUint64(v, valueTypeInvalid, false)
	case *[]int64:
		d.decSliceInt64(v, valueTypeInvalid, false)
	case *[]string:
		d.decSliceStr(v, valueTypeInvalid, false)
	case *map[string]interface***REMOVED******REMOVED***:
		d.decMapStrIntf(v)
	case *map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
		d.decMapIntfIntf(v)
	case *map[uint64]interface***REMOVED******REMOVED***:
		d.decMapUint64Intf(v)
	case *map[int64]interface***REMOVED******REMOVED***:
		d.decMapInt64Intf(v)

	case *interface***REMOVED******REMOVED***:
		d.decodeValue(reflect.ValueOf(iv).Elem())

	default:
		rv := reflect.ValueOf(iv)
		d.chkPtrValue(rv)
		d.decodeValue(rv.Elem())
	***REMOVED***
***REMOVED***

func (d *Decoder) decodeValue(rv reflect.Value) ***REMOVED***
	d.d.initReadNext()

	if d.d.tryDecodeAsNil() ***REMOVED***
		// If value in stream is nil, set the dereferenced value to its "zero" value (if settable)
		if rv.Kind() == reflect.Ptr ***REMOVED***
			if !rv.IsNil() ***REMOVED***
				rv.Set(reflect.Zero(rv.Type()))
			***REMOVED***
			return
		***REMOVED***
		// for rv.Kind() == reflect.Ptr ***REMOVED***
		// 	rv = rv.Elem()
		// ***REMOVED***
		if rv.IsValid() ***REMOVED*** // rv.CanSet() // always settable, except it's invalid
			rv.Set(reflect.Zero(rv.Type()))
		***REMOVED***
		return
	***REMOVED***

	// If stream is not containing a nil value, then we can deref to the base
	// non-pointer value, and decode into that.
	for rv.Kind() == reflect.Ptr ***REMOVED***
		if rv.IsNil() ***REMOVED***
			rv.Set(reflect.New(rv.Type().Elem()))
		***REMOVED***
		rv = rv.Elem()
	***REMOVED***

	rt := rv.Type()
	rtid := reflect.ValueOf(rt).Pointer()

	// retrieve or register a focus'ed function for this type
	// to eliminate need to do the retrieval multiple times

	// if d.f == nil && d.s == nil ***REMOVED*** debugf("---->Creating new dec f map for type: %v\n", rt) ***REMOVED***
	var fn decFn
	var ok bool
	if useMapForCodecCache ***REMOVED***
		fn, ok = d.f[rtid]
	***REMOVED*** else ***REMOVED***
		for i, v := range d.x ***REMOVED***
			if v == rtid ***REMOVED***
				fn, ok = d.s[i], true
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if !ok ***REMOVED***
		// debugf("\tCreating new dec fn for type: %v\n", rt)
		fi := decFnInfo***REMOVED***ti: getTypeInfo(rtid, rt), d: d, dd: d.d***REMOVED***
		fn.i = &fi
		// An extension can be registered for any type, regardless of the Kind
		// (e.g. type BitSet int64, type MyStruct ***REMOVED*** / * unexported fields * / ***REMOVED***, type X []int, etc.
		//
		// We can't check if it's an extension byte here first, because the user may have
		// registered a pointer or non-pointer type, meaning we may have to recurse first
		// before matching a mapped type, even though the extension byte is already detected.
		//
		// NOTE: if decoding into a nil interface***REMOVED******REMOVED***, we return a non-nil
		// value except even if the container registers a length of 0.
		if rtid == rawExtTypId ***REMOVED***
			fn.f = (*decFnInfo).rawExt
		***REMOVED*** else if d.d.isBuiltinType(rtid) ***REMOVED***
			fn.f = (*decFnInfo).builtin
		***REMOVED*** else if xfTag, xfFn := d.h.getDecodeExt(rtid); xfFn != nil ***REMOVED***
			fi.xfTag, fi.xfFn = xfTag, xfFn
			fn.f = (*decFnInfo).ext
		***REMOVED*** else if supportBinaryMarshal && fi.ti.unm ***REMOVED***
			fn.f = (*decFnInfo).binaryMarshal
		***REMOVED*** else ***REMOVED***
			switch rk := rt.Kind(); rk ***REMOVED***
			case reflect.String:
				fn.f = (*decFnInfo).kString
			case reflect.Bool:
				fn.f = (*decFnInfo).kBool
			case reflect.Int:
				fn.f = (*decFnInfo).kInt
			case reflect.Int64:
				fn.f = (*decFnInfo).kInt64
			case reflect.Int32:
				fn.f = (*decFnInfo).kInt32
			case reflect.Int8:
				fn.f = (*decFnInfo).kInt8
			case reflect.Int16:
				fn.f = (*decFnInfo).kInt16
			case reflect.Float32:
				fn.f = (*decFnInfo).kFloat32
			case reflect.Float64:
				fn.f = (*decFnInfo).kFloat64
			case reflect.Uint8:
				fn.f = (*decFnInfo).kUint8
			case reflect.Uint64:
				fn.f = (*decFnInfo).kUint64
			case reflect.Uint:
				fn.f = (*decFnInfo).kUint
			case reflect.Uint32:
				fn.f = (*decFnInfo).kUint32
			case reflect.Uint16:
				fn.f = (*decFnInfo).kUint16
			// case reflect.Ptr:
			// 	fn.f = (*decFnInfo).kPtr
			case reflect.Interface:
				fn.f = (*decFnInfo).kInterface
			case reflect.Struct:
				fn.f = (*decFnInfo).kStruct
			case reflect.Slice:
				fn.f = (*decFnInfo).kSlice
			case reflect.Array:
				fi.array = true
				fn.f = (*decFnInfo).kArray
			case reflect.Map:
				fn.f = (*decFnInfo).kMap
			default:
				fn.f = (*decFnInfo).kErr
			***REMOVED***
		***REMOVED***
		if useMapForCodecCache ***REMOVED***
			if d.f == nil ***REMOVED***
				d.f = make(map[uintptr]decFn, 16)
			***REMOVED***
			d.f[rtid] = fn
		***REMOVED*** else ***REMOVED***
			d.s = append(d.s, fn)
			d.x = append(d.x, rtid)
		***REMOVED***
	***REMOVED***

	fn.f(fn.i, rv)

	return
***REMOVED***

func (d *Decoder) chkPtrValue(rv reflect.Value) ***REMOVED***
	// We can only decode into a non-nil pointer
	if rv.Kind() == reflect.Ptr && !rv.IsNil() ***REMOVED***
		return
	***REMOVED***
	if !rv.IsValid() ***REMOVED***
		decErr("Cannot decode into a zero (ie invalid) reflect.Value")
	***REMOVED***
	if !rv.CanInterface() ***REMOVED***
		decErr("Cannot decode into a value without an interface: %v", rv)
	***REMOVED***
	rvi := rv.Interface()
	decErr("Cannot decode into non-pointer or nil pointer. Got: %v, %T, %v",
		rv.Kind(), rvi, rvi)
***REMOVED***

func (d *Decoder) decEmbeddedField(rv reflect.Value, index []int) ***REMOVED***
	// d.decodeValue(rv.FieldByIndex(index))
	// nil pointers may be here; so reproduce FieldByIndex logic + enhancements
	for _, j := range index ***REMOVED***
		if rv.Kind() == reflect.Ptr ***REMOVED***
			if rv.IsNil() ***REMOVED***
				rv.Set(reflect.New(rv.Type().Elem()))
			***REMOVED***
			// If a pointer, it must be a pointer to struct (based on typeInfo contract)
			rv = rv.Elem()
		***REMOVED***
		rv = rv.Field(j)
	***REMOVED***
	d.decodeValue(rv)
***REMOVED***

// --------------------------------------------------

// short circuit functions for common maps and slices

func (d *Decoder) decSliceIntf(v *[]interface***REMOVED******REMOVED***, currEncodedType valueType, doNotReset bool) ***REMOVED***
	_, containerLenS := decContLens(d.d, currEncodedType)
	s := *v
	if s == nil ***REMOVED***
		s = make([]interface***REMOVED******REMOVED***, containerLenS, containerLenS)
	***REMOVED*** else if containerLenS > cap(s) ***REMOVED***
		if doNotReset ***REMOVED***
			decErr(msgDecCannotExpandArr, cap(s), containerLenS)
		***REMOVED***
		s = make([]interface***REMOVED******REMOVED***, containerLenS, containerLenS)
		copy(s, *v)
	***REMOVED*** else if containerLenS > len(s) ***REMOVED***
		s = s[:containerLenS]
	***REMOVED***
	for j := 0; j < containerLenS; j++ ***REMOVED***
		d.decode(&s[j])
	***REMOVED***
	*v = s
***REMOVED***

func (d *Decoder) decSliceInt64(v *[]int64, currEncodedType valueType, doNotReset bool) ***REMOVED***
	_, containerLenS := decContLens(d.d, currEncodedType)
	s := *v
	if s == nil ***REMOVED***
		s = make([]int64, containerLenS, containerLenS)
	***REMOVED*** else if containerLenS > cap(s) ***REMOVED***
		if doNotReset ***REMOVED***
			decErr(msgDecCannotExpandArr, cap(s), containerLenS)
		***REMOVED***
		s = make([]int64, containerLenS, containerLenS)
		copy(s, *v)
	***REMOVED*** else if containerLenS > len(s) ***REMOVED***
		s = s[:containerLenS]
	***REMOVED***
	for j := 0; j < containerLenS; j++ ***REMOVED***
		// d.decode(&s[j])
		d.d.initReadNext()
		s[j] = d.d.decodeInt(intBitsize)
	***REMOVED***
	*v = s
***REMOVED***

func (d *Decoder) decSliceUint64(v *[]uint64, currEncodedType valueType, doNotReset bool) ***REMOVED***
	_, containerLenS := decContLens(d.d, currEncodedType)
	s := *v
	if s == nil ***REMOVED***
		s = make([]uint64, containerLenS, containerLenS)
	***REMOVED*** else if containerLenS > cap(s) ***REMOVED***
		if doNotReset ***REMOVED***
			decErr(msgDecCannotExpandArr, cap(s), containerLenS)
		***REMOVED***
		s = make([]uint64, containerLenS, containerLenS)
		copy(s, *v)
	***REMOVED*** else if containerLenS > len(s) ***REMOVED***
		s = s[:containerLenS]
	***REMOVED***
	for j := 0; j < containerLenS; j++ ***REMOVED***
		// d.decode(&s[j])
		d.d.initReadNext()
		s[j] = d.d.decodeUint(intBitsize)
	***REMOVED***
	*v = s
***REMOVED***

func (d *Decoder) decSliceStr(v *[]string, currEncodedType valueType, doNotReset bool) ***REMOVED***
	_, containerLenS := decContLens(d.d, currEncodedType)
	s := *v
	if s == nil ***REMOVED***
		s = make([]string, containerLenS, containerLenS)
	***REMOVED*** else if containerLenS > cap(s) ***REMOVED***
		if doNotReset ***REMOVED***
			decErr(msgDecCannotExpandArr, cap(s), containerLenS)
		***REMOVED***
		s = make([]string, containerLenS, containerLenS)
		copy(s, *v)
	***REMOVED*** else if containerLenS > len(s) ***REMOVED***
		s = s[:containerLenS]
	***REMOVED***
	for j := 0; j < containerLenS; j++ ***REMOVED***
		// d.decode(&s[j])
		d.d.initReadNext()
		s[j] = d.d.decodeString()
	***REMOVED***
	*v = s
***REMOVED***

func (d *Decoder) decMapIntfIntf(v *map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***) ***REMOVED***
	containerLen := d.d.readMapLen()
	m := *v
	if m == nil ***REMOVED***
		m = make(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***, containerLen)
		*v = m
	***REMOVED***
	for j := 0; j < containerLen; j++ ***REMOVED***
		var mk interface***REMOVED******REMOVED***
		d.decode(&mk)
		// special case if a byte array.
		if bv, bok := mk.([]byte); bok ***REMOVED***
			mk = string(bv)
		***REMOVED***
		mv := m[mk]
		d.decode(&mv)
		m[mk] = mv
	***REMOVED***
***REMOVED***

func (d *Decoder) decMapInt64Intf(v *map[int64]interface***REMOVED******REMOVED***) ***REMOVED***
	containerLen := d.d.readMapLen()
	m := *v
	if m == nil ***REMOVED***
		m = make(map[int64]interface***REMOVED******REMOVED***, containerLen)
		*v = m
	***REMOVED***
	for j := 0; j < containerLen; j++ ***REMOVED***
		d.d.initReadNext()
		mk := d.d.decodeInt(intBitsize)
		mv := m[mk]
		d.decode(&mv)
		m[mk] = mv
	***REMOVED***
***REMOVED***

func (d *Decoder) decMapUint64Intf(v *map[uint64]interface***REMOVED******REMOVED***) ***REMOVED***
	containerLen := d.d.readMapLen()
	m := *v
	if m == nil ***REMOVED***
		m = make(map[uint64]interface***REMOVED******REMOVED***, containerLen)
		*v = m
	***REMOVED***
	for j := 0; j < containerLen; j++ ***REMOVED***
		d.d.initReadNext()
		mk := d.d.decodeUint(intBitsize)
		mv := m[mk]
		d.decode(&mv)
		m[mk] = mv
	***REMOVED***
***REMOVED***

func (d *Decoder) decMapStrIntf(v *map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	containerLen := d.d.readMapLen()
	m := *v
	if m == nil ***REMOVED***
		m = make(map[string]interface***REMOVED******REMOVED***, containerLen)
		*v = m
	***REMOVED***
	for j := 0; j < containerLen; j++ ***REMOVED***
		d.d.initReadNext()
		mk := d.d.decodeString()
		mv := m[mk]
		d.decode(&mv)
		m[mk] = mv
	***REMOVED***
***REMOVED***

// ----------------------------------------

func decContLens(dd decDriver, currEncodedType valueType) (containerLen, containerLenS int) ***REMOVED***
	if currEncodedType == valueTypeInvalid ***REMOVED***
		currEncodedType = dd.currentEncodedType()
	***REMOVED***
	switch currEncodedType ***REMOVED***
	case valueTypeArray:
		containerLen = dd.readArrayLen()
		containerLenS = containerLen
	case valueTypeMap:
		containerLen = dd.readMapLen()
		containerLenS = containerLen * 2
	default:
		decErr("Only encoded map or array can be decoded into a slice. (valueType: %0x)",
			currEncodedType)
	***REMOVED***
	return
***REMOVED***

func decErr(format string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
	doPanic(msgTagDec, format, params...)
***REMOVED***
