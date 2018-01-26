// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"io"
	"reflect"
)

const (
	// Some tagging information for error messages.
	msgTagEnc         = "codec.encoder"
	defEncByteBufSize = 1 << 6 // 4:16, 6:64, 8:256, 10:1024
	// maxTimeSecs32 = math.MaxInt32 / 60 / 24 / 366
)

// AsSymbolFlag defines what should be encoded as symbols.
type AsSymbolFlag uint8

const (
	// AsSymbolDefault is default.
	// Currently, this means only encode struct field names as symbols.
	// The default is subject to change.
	AsSymbolDefault AsSymbolFlag = iota

	// AsSymbolAll means encode anything which could be a symbol as a symbol.
	AsSymbolAll = 0xfe

	// AsSymbolNone means do not encode anything as a symbol.
	AsSymbolNone = 1 << iota

	// AsSymbolMapStringKeys means encode keys in map[string]XXX as symbols.
	AsSymbolMapStringKeysFlag

	// AsSymbolStructFieldName means encode struct field names as symbols.
	AsSymbolStructFieldNameFlag
)

// encWriter abstracting writing to a byte array or to an io.Writer.
type encWriter interface ***REMOVED***
	writeUint16(uint16)
	writeUint32(uint32)
	writeUint64(uint64)
	writeb([]byte)
	writestr(string)
	writen1(byte)
	writen2(byte, byte)
	atEndOfEncode()
***REMOVED***

// encDriver abstracts the actual codec (binc vs msgpack, etc)
type encDriver interface ***REMOVED***
	isBuiltinType(rt uintptr) bool
	encodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***)
	encodeNil()
	encodeInt(i int64)
	encodeUint(i uint64)
	encodeBool(b bool)
	encodeFloat32(f float32)
	encodeFloat64(f float64)
	encodeExtPreamble(xtag byte, length int)
	encodeArrayPreamble(length int)
	encodeMapPreamble(length int)
	encodeString(c charEncoding, v string)
	encodeSymbol(v string)
	encodeStringBytes(c charEncoding, v []byte)
	//TODO
	//encBignum(f *big.Int)
	//encStringRunes(c charEncoding, v []rune)
***REMOVED***

type ioEncWriterWriter interface ***REMOVED***
	WriteByte(c byte) error
	WriteString(s string) (n int, err error)
	Write(p []byte) (n int, err error)
***REMOVED***

type ioEncStringWriter interface ***REMOVED***
	WriteString(s string) (n int, err error)
***REMOVED***

type EncodeOptions struct ***REMOVED***
	// Encode a struct as an array, and not as a map.
	StructToArray bool

	// AsSymbols defines what should be encoded as symbols.
	//
	// Encoding as symbols can reduce the encoded size significantly.
	//
	// However, during decoding, each string to be encoded as a symbol must
	// be checked to see if it has been seen before. Consequently, encoding time
	// will increase if using symbols, because string comparisons has a clear cost.
	//
	// Sample values:
	//   AsSymbolNone
	//   AsSymbolAll
	//   AsSymbolMapStringKeys
	//   AsSymbolMapStringKeysFlag | AsSymbolStructFieldNameFlag
	AsSymbols AsSymbolFlag
***REMOVED***

// ---------------------------------------------

type simpleIoEncWriterWriter struct ***REMOVED***
	w  io.Writer
	bw io.ByteWriter
	sw ioEncStringWriter
***REMOVED***

func (o *simpleIoEncWriterWriter) WriteByte(c byte) (err error) ***REMOVED***
	if o.bw != nil ***REMOVED***
		return o.bw.WriteByte(c)
	***REMOVED***
	_, err = o.w.Write([]byte***REMOVED***c***REMOVED***)
	return
***REMOVED***

func (o *simpleIoEncWriterWriter) WriteString(s string) (n int, err error) ***REMOVED***
	if o.sw != nil ***REMOVED***
		return o.sw.WriteString(s)
	***REMOVED***
	return o.w.Write([]byte(s))
***REMOVED***

func (o *simpleIoEncWriterWriter) Write(p []byte) (n int, err error) ***REMOVED***
	return o.w.Write(p)
***REMOVED***

// ----------------------------------------

// ioEncWriter implements encWriter and can write to an io.Writer implementation
type ioEncWriter struct ***REMOVED***
	w ioEncWriterWriter
	x [8]byte // temp byte array re-used internally for efficiency
***REMOVED***

func (z *ioEncWriter) writeUint16(v uint16) ***REMOVED***
	bigen.PutUint16(z.x[:2], v)
	z.writeb(z.x[:2])
***REMOVED***

func (z *ioEncWriter) writeUint32(v uint32) ***REMOVED***
	bigen.PutUint32(z.x[:4], v)
	z.writeb(z.x[:4])
***REMOVED***

func (z *ioEncWriter) writeUint64(v uint64) ***REMOVED***
	bigen.PutUint64(z.x[:8], v)
	z.writeb(z.x[:8])
***REMOVED***

func (z *ioEncWriter) writeb(bs []byte) ***REMOVED***
	if len(bs) == 0 ***REMOVED***
		return
	***REMOVED***
	n, err := z.w.Write(bs)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	if n != len(bs) ***REMOVED***
		encErr("write: Incorrect num bytes written. Expecting: %v, Wrote: %v", len(bs), n)
	***REMOVED***
***REMOVED***

func (z *ioEncWriter) writestr(s string) ***REMOVED***
	n, err := z.w.WriteString(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	if n != len(s) ***REMOVED***
		encErr("write: Incorrect num bytes written. Expecting: %v, Wrote: %v", len(s), n)
	***REMOVED***
***REMOVED***

func (z *ioEncWriter) writen1(b byte) ***REMOVED***
	if err := z.w.WriteByte(b); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (z *ioEncWriter) writen2(b1 byte, b2 byte) ***REMOVED***
	z.writen1(b1)
	z.writen1(b2)
***REMOVED***

func (z *ioEncWriter) atEndOfEncode() ***REMOVED******REMOVED***

// ----------------------------------------

// bytesEncWriter implements encWriter and can write to an byte slice.
// It is used by Marshal function.
type bytesEncWriter struct ***REMOVED***
	b   []byte
	c   int     // cursor
	out *[]byte // write out on atEndOfEncode
***REMOVED***

func (z *bytesEncWriter) writeUint16(v uint16) ***REMOVED***
	c := z.grow(2)
	z.b[c] = byte(v >> 8)
	z.b[c+1] = byte(v)
***REMOVED***

func (z *bytesEncWriter) writeUint32(v uint32) ***REMOVED***
	c := z.grow(4)
	z.b[c] = byte(v >> 24)
	z.b[c+1] = byte(v >> 16)
	z.b[c+2] = byte(v >> 8)
	z.b[c+3] = byte(v)
***REMOVED***

func (z *bytesEncWriter) writeUint64(v uint64) ***REMOVED***
	c := z.grow(8)
	z.b[c] = byte(v >> 56)
	z.b[c+1] = byte(v >> 48)
	z.b[c+2] = byte(v >> 40)
	z.b[c+3] = byte(v >> 32)
	z.b[c+4] = byte(v >> 24)
	z.b[c+5] = byte(v >> 16)
	z.b[c+6] = byte(v >> 8)
	z.b[c+7] = byte(v)
***REMOVED***

func (z *bytesEncWriter) writeb(s []byte) ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return
	***REMOVED***
	c := z.grow(len(s))
	copy(z.b[c:], s)
***REMOVED***

func (z *bytesEncWriter) writestr(s string) ***REMOVED***
	c := z.grow(len(s))
	copy(z.b[c:], s)
***REMOVED***

func (z *bytesEncWriter) writen1(b1 byte) ***REMOVED***
	c := z.grow(1)
	z.b[c] = b1
***REMOVED***

func (z *bytesEncWriter) writen2(b1 byte, b2 byte) ***REMOVED***
	c := z.grow(2)
	z.b[c] = b1
	z.b[c+1] = b2
***REMOVED***

func (z *bytesEncWriter) atEndOfEncode() ***REMOVED***
	*(z.out) = z.b[:z.c]
***REMOVED***

func (z *bytesEncWriter) grow(n int) (oldcursor int) ***REMOVED***
	oldcursor = z.c
	z.c = oldcursor + n
	if z.c > cap(z.b) ***REMOVED***
		// Tried using appendslice logic: (if cap < 1024, *2, else *1.25).
		// However, it was too expensive, causing too many iterations of copy.
		// Using bytes.Buffer model was much better (2*cap + n)
		bs := make([]byte, 2*cap(z.b)+n)
		copy(bs, z.b[:oldcursor])
		z.b = bs
	***REMOVED*** else if z.c > len(z.b) ***REMOVED***
		z.b = z.b[:cap(z.b)]
	***REMOVED***
	return
***REMOVED***

// ---------------------------------------------

type encFnInfo struct ***REMOVED***
	ti    *typeInfo
	e     *Encoder
	ee    encDriver
	xfFn  func(reflect.Value) ([]byte, error)
	xfTag byte
***REMOVED***

func (f *encFnInfo) builtin(rv reflect.Value) ***REMOVED***
	f.ee.encodeBuiltin(f.ti.rtid, rv.Interface())
***REMOVED***

func (f *encFnInfo) rawExt(rv reflect.Value) ***REMOVED***
	f.e.encRawExt(rv.Interface().(RawExt))
***REMOVED***

func (f *encFnInfo) ext(rv reflect.Value) ***REMOVED***
	bs, fnerr := f.xfFn(rv)
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
	if bs == nil ***REMOVED***
		f.ee.encodeNil()
		return
	***REMOVED***
	if f.e.hh.writeExt() ***REMOVED***
		f.ee.encodeExtPreamble(f.xfTag, len(bs))
		f.e.w.writeb(bs)
	***REMOVED*** else ***REMOVED***
		f.ee.encodeStringBytes(c_RAW, bs)
	***REMOVED***

***REMOVED***

func (f *encFnInfo) binaryMarshal(rv reflect.Value) ***REMOVED***
	var bm binaryMarshaler
	if f.ti.mIndir == 0 ***REMOVED***
		bm = rv.Interface().(binaryMarshaler)
	***REMOVED*** else if f.ti.mIndir == -1 ***REMOVED***
		bm = rv.Addr().Interface().(binaryMarshaler)
	***REMOVED*** else ***REMOVED***
		for j, k := int8(0), f.ti.mIndir; j < k; j++ ***REMOVED***
			if rv.IsNil() ***REMOVED***
				f.ee.encodeNil()
				return
			***REMOVED***
			rv = rv.Elem()
		***REMOVED***
		bm = rv.Interface().(binaryMarshaler)
	***REMOVED***
	// debugf(">>>> binaryMarshaler: %T", rv.Interface())
	bs, fnerr := bm.MarshalBinary()
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
	if bs == nil ***REMOVED***
		f.ee.encodeNil()
	***REMOVED*** else ***REMOVED***
		f.ee.encodeStringBytes(c_RAW, bs)
	***REMOVED***
***REMOVED***

func (f *encFnInfo) kBool(rv reflect.Value) ***REMOVED***
	f.ee.encodeBool(rv.Bool())
***REMOVED***

func (f *encFnInfo) kString(rv reflect.Value) ***REMOVED***
	f.ee.encodeString(c_UTF8, rv.String())
***REMOVED***

func (f *encFnInfo) kFloat64(rv reflect.Value) ***REMOVED***
	f.ee.encodeFloat64(rv.Float())
***REMOVED***

func (f *encFnInfo) kFloat32(rv reflect.Value) ***REMOVED***
	f.ee.encodeFloat32(float32(rv.Float()))
***REMOVED***

func (f *encFnInfo) kInt(rv reflect.Value) ***REMOVED***
	f.ee.encodeInt(rv.Int())
***REMOVED***

func (f *encFnInfo) kUint(rv reflect.Value) ***REMOVED***
	f.ee.encodeUint(rv.Uint())
***REMOVED***

func (f *encFnInfo) kInvalid(rv reflect.Value) ***REMOVED***
	f.ee.encodeNil()
***REMOVED***

func (f *encFnInfo) kErr(rv reflect.Value) ***REMOVED***
	encErr("Unsupported kind: %s, for: %#v", rv.Kind(), rv)
***REMOVED***

func (f *encFnInfo) kSlice(rv reflect.Value) ***REMOVED***
	if rv.IsNil() ***REMOVED***
		f.ee.encodeNil()
		return
	***REMOVED***

	if shortCircuitReflectToFastPath ***REMOVED***
		switch f.ti.rtid ***REMOVED***
		case intfSliceTypId:
			f.e.encSliceIntf(rv.Interface().([]interface***REMOVED******REMOVED***))
			return
		case strSliceTypId:
			f.e.encSliceStr(rv.Interface().([]string))
			return
		case uint64SliceTypId:
			f.e.encSliceUint64(rv.Interface().([]uint64))
			return
		case int64SliceTypId:
			f.e.encSliceInt64(rv.Interface().([]int64))
			return
		***REMOVED***
	***REMOVED***

	// If in this method, then there was no extension function defined.
	// So it's okay to treat as []byte.
	if f.ti.rtid == uint8SliceTypId || f.ti.rt.Elem().Kind() == reflect.Uint8 ***REMOVED***
		f.ee.encodeStringBytes(c_RAW, rv.Bytes())
		return
	***REMOVED***

	l := rv.Len()
	if f.ti.mbs ***REMOVED***
		if l%2 == 1 ***REMOVED***
			encErr("mapBySlice: invalid length (must be divisible by 2): %v", l)
		***REMOVED***
		f.ee.encodeMapPreamble(l / 2)
	***REMOVED*** else ***REMOVED***
		f.ee.encodeArrayPreamble(l)
	***REMOVED***
	if l == 0 ***REMOVED***
		return
	***REMOVED***
	for j := 0; j < l; j++ ***REMOVED***
		// TODO: Consider perf implication of encoding odd index values as symbols if type is string
		f.e.encodeValue(rv.Index(j))
	***REMOVED***
***REMOVED***

func (f *encFnInfo) kArray(rv reflect.Value) ***REMOVED***
	// We cannot share kSlice method, because the array may be non-addressable.
	// E.g. type struct S***REMOVED***B [2]byte***REMOVED***; Encode(S***REMOVED******REMOVED***) will bomb on "panic: slice of unaddressable array".
	// So we have to duplicate the functionality here.
	// f.e.encodeValue(rv.Slice(0, rv.Len()))
	// f.kSlice(rv.Slice(0, rv.Len()))

	l := rv.Len()
	// Handle an array of bytes specially (in line with what is done for slices)
	if f.ti.rt.Elem().Kind() == reflect.Uint8 ***REMOVED***
		if l == 0 ***REMOVED***
			f.ee.encodeStringBytes(c_RAW, nil)
			return
		***REMOVED***
		var bs []byte
		if rv.CanAddr() ***REMOVED***
			bs = rv.Slice(0, l).Bytes()
		***REMOVED*** else ***REMOVED***
			bs = make([]byte, l)
			for i := 0; i < l; i++ ***REMOVED***
				bs[i] = byte(rv.Index(i).Uint())
			***REMOVED***
		***REMOVED***
		f.ee.encodeStringBytes(c_RAW, bs)
		return
	***REMOVED***

	if f.ti.mbs ***REMOVED***
		if l%2 == 1 ***REMOVED***
			encErr("mapBySlice: invalid length (must be divisible by 2): %v", l)
		***REMOVED***
		f.ee.encodeMapPreamble(l / 2)
	***REMOVED*** else ***REMOVED***
		f.ee.encodeArrayPreamble(l)
	***REMOVED***
	if l == 0 ***REMOVED***
		return
	***REMOVED***
	for j := 0; j < l; j++ ***REMOVED***
		// TODO: Consider perf implication of encoding odd index values as symbols if type is string
		f.e.encodeValue(rv.Index(j))
	***REMOVED***
***REMOVED***

func (f *encFnInfo) kStruct(rv reflect.Value) ***REMOVED***
	fti := f.ti
	newlen := len(fti.sfi)
	rvals := make([]reflect.Value, newlen)
	var encnames []string
	e := f.e
	tisfi := fti.sfip
	toMap := !(fti.toArray || e.h.StructToArray)
	// if toMap, use the sorted array. If toArray, use unsorted array (to match sequence in struct)
	if toMap ***REMOVED***
		tisfi = fti.sfi
		encnames = make([]string, newlen)
	***REMOVED***
	newlen = 0
	for _, si := range tisfi ***REMOVED***
		if si.i != -1 ***REMOVED***
			rvals[newlen] = rv.Field(int(si.i))
		***REMOVED*** else ***REMOVED***
			rvals[newlen] = rv.FieldByIndex(si.is)
		***REMOVED***
		if toMap ***REMOVED***
			if si.omitEmpty && isEmptyValue(rvals[newlen]) ***REMOVED***
				continue
			***REMOVED***
			encnames[newlen] = si.encName
		***REMOVED*** else ***REMOVED***
			if si.omitEmpty && isEmptyValue(rvals[newlen]) ***REMOVED***
				rvals[newlen] = reflect.Value***REMOVED******REMOVED*** //encode as nil
			***REMOVED***
		***REMOVED***
		newlen++
	***REMOVED***

	// debugf(">>>> kStruct: newlen: %v", newlen)
	if toMap ***REMOVED***
		ee := f.ee //don't dereference everytime
		ee.encodeMapPreamble(newlen)
		// asSymbols := e.h.AsSymbols&AsSymbolStructFieldNameFlag != 0
		asSymbols := e.h.AsSymbols == AsSymbolDefault || e.h.AsSymbols&AsSymbolStructFieldNameFlag != 0
		for j := 0; j < newlen; j++ ***REMOVED***
			if asSymbols ***REMOVED***
				ee.encodeSymbol(encnames[j])
			***REMOVED*** else ***REMOVED***
				ee.encodeString(c_UTF8, encnames[j])
			***REMOVED***
			e.encodeValue(rvals[j])
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		f.ee.encodeArrayPreamble(newlen)
		for j := 0; j < newlen; j++ ***REMOVED***
			e.encodeValue(rvals[j])
		***REMOVED***
	***REMOVED***
***REMOVED***

// func (f *encFnInfo) kPtr(rv reflect.Value) ***REMOVED***
// 	debugf(">>>>>>> ??? encode kPtr called - shouldn't get called")
// 	if rv.IsNil() ***REMOVED***
// 		f.ee.encodeNil()
// 		return
// 	***REMOVED***
// 	f.e.encodeValue(rv.Elem())
// ***REMOVED***

func (f *encFnInfo) kInterface(rv reflect.Value) ***REMOVED***
	if rv.IsNil() ***REMOVED***
		f.ee.encodeNil()
		return
	***REMOVED***
	f.e.encodeValue(rv.Elem())
***REMOVED***

func (f *encFnInfo) kMap(rv reflect.Value) ***REMOVED***
	if rv.IsNil() ***REMOVED***
		f.ee.encodeNil()
		return
	***REMOVED***

	if shortCircuitReflectToFastPath ***REMOVED***
		switch f.ti.rtid ***REMOVED***
		case mapIntfIntfTypId:
			f.e.encMapIntfIntf(rv.Interface().(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***))
			return
		case mapStrIntfTypId:
			f.e.encMapStrIntf(rv.Interface().(map[string]interface***REMOVED******REMOVED***))
			return
		case mapStrStrTypId:
			f.e.encMapStrStr(rv.Interface().(map[string]string))
			return
		case mapInt64IntfTypId:
			f.e.encMapInt64Intf(rv.Interface().(map[int64]interface***REMOVED******REMOVED***))
			return
		case mapUint64IntfTypId:
			f.e.encMapUint64Intf(rv.Interface().(map[uint64]interface***REMOVED******REMOVED***))
			return
		***REMOVED***
	***REMOVED***

	l := rv.Len()
	f.ee.encodeMapPreamble(l)
	if l == 0 ***REMOVED***
		return
	***REMOVED***
	// keyTypeIsString := f.ti.rt.Key().Kind() == reflect.String
	keyTypeIsString := f.ti.rt.Key() == stringTyp
	var asSymbols bool
	if keyTypeIsString ***REMOVED***
		asSymbols = f.e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	***REMOVED***
	mks := rv.MapKeys()
	// for j, lmks := 0, len(mks); j < lmks; j++ ***REMOVED***
	for j := range mks ***REMOVED***
		if keyTypeIsString ***REMOVED***
			if asSymbols ***REMOVED***
				f.ee.encodeSymbol(mks[j].String())
			***REMOVED*** else ***REMOVED***
				f.ee.encodeString(c_UTF8, mks[j].String())
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			f.e.encodeValue(mks[j])
		***REMOVED***
		f.e.encodeValue(rv.MapIndex(mks[j]))
	***REMOVED***

***REMOVED***

// --------------------------------------------------

// encFn encapsulates the captured variables and the encode function.
// This way, we only do some calculations one times, and pass to the
// code block that should be called (encapsulated in a function)
// instead of executing the checks every time.
type encFn struct ***REMOVED***
	i *encFnInfo
	f func(*encFnInfo, reflect.Value)
***REMOVED***

// --------------------------------------------------

// An Encoder writes an object to an output stream in the codec format.
type Encoder struct ***REMOVED***
	w  encWriter
	e  encDriver
	h  *BasicHandle
	hh Handle
	f  map[uintptr]encFn
	x  []uintptr
	s  []encFn
***REMOVED***

// NewEncoder returns an Encoder for encoding into an io.Writer.
//
// For efficiency, Users are encouraged to pass in a memory buffered writer
// (eg bufio.Writer, bytes.Buffer).
func NewEncoder(w io.Writer, h Handle) *Encoder ***REMOVED***
	ww, ok := w.(ioEncWriterWriter)
	if !ok ***REMOVED***
		sww := simpleIoEncWriterWriter***REMOVED***w: w***REMOVED***
		sww.bw, _ = w.(io.ByteWriter)
		sww.sw, _ = w.(ioEncStringWriter)
		ww = &sww
		//ww = bufio.NewWriterSize(w, defEncByteBufSize)
	***REMOVED***
	z := ioEncWriter***REMOVED***
		w: ww,
	***REMOVED***
	return &Encoder***REMOVED***w: &z, hh: h, h: h.getBasicHandle(), e: h.newEncDriver(&z)***REMOVED***
***REMOVED***

// NewEncoderBytes returns an encoder for encoding directly and efficiently
// into a byte slice, using zero-copying to temporary slices.
//
// It will potentially replace the output byte slice pointed to.
// After encoding, the out parameter contains the encoded contents.
func NewEncoderBytes(out *[]byte, h Handle) *Encoder ***REMOVED***
	in := *out
	if in == nil ***REMOVED***
		in = make([]byte, defEncByteBufSize)
	***REMOVED***
	z := bytesEncWriter***REMOVED***
		b:   in,
		out: out,
	***REMOVED***
	return &Encoder***REMOVED***w: &z, hh: h, h: h.getBasicHandle(), e: h.newEncDriver(&z)***REMOVED***
***REMOVED***

// Encode writes an object into a stream in the codec format.
//
// Encoding can be configured via the "codec" struct tag for the fields.
//
// The "codec" key in struct field's tag value is the key name,
// followed by an optional comma and options.
//
// To set an option on all fields (e.g. omitempty on all fields), you
// can create a field called _struct, and set flags on it.
//
// Struct values "usually" encode as maps. Each exported struct field is encoded unless:
//    - the field's codec tag is "-", OR
//    - the field is empty and its codec tag specifies the "omitempty" option.
//
// When encoding as a map, the first string in the tag (before the comma)
// is the map key string to use when encoding.
//
// However, struct values may encode as arrays. This happens when:
//    - StructToArray Encode option is set, OR
//    - the codec tag on the _struct field sets the "toarray" option
//
// Values with types that implement MapBySlice are encoded as stream maps.
//
// The empty values (for omitempty option) are false, 0, any nil pointer
// or interface value, and any array, slice, map, or string of length zero.
//
// Anonymous fields are encoded inline if no struct tag is present.
// Else they are encoded as regular fields.
//
// Examples:
//
//      type MyStruct struct ***REMOVED***
//          _struct bool    `codec:",omitempty"`   //set omitempty for every field
//          Field1 string   `codec:"-"`            //skip this field
//          Field2 int      `codec:"myName"`       //Use key "myName" in encode stream
//          Field3 int32    `codec:",omitempty"`   //use key "Field3". Omit if empty.
//          Field4 bool     `codec:"f4,omitempty"` //use key "f4". Omit if empty.
//          ...
//  ***REMOVED***
//
//      type MyStruct struct ***REMOVED***
//          _struct bool    `codec:",omitempty,toarray"`   //set omitempty for every field
//                                                         //and encode struct as an array
//  ***REMOVED***
//
// The mode of encoding is based on the type of the value. When a value is seen:
//   - If an extension is registered for it, call that extension function
//   - If it implements BinaryMarshaler, call its MarshalBinary() (data []byte, err error)
//   - Else encode it based on its reflect.Kind
//
// Note that struct field names and keys in map[string]XXX will be treated as symbols.
// Some formats support symbols (e.g. binc) and will properly encode the string
// only once in the stream, and use a tag to refer to it thereafter.
func (e *Encoder) Encode(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	defer panicToErr(&err)
	e.encode(v)
	e.w.atEndOfEncode()
	return
***REMOVED***

func (e *Encoder) encode(iv interface***REMOVED******REMOVED***) ***REMOVED***
	switch v := iv.(type) ***REMOVED***
	case nil:
		e.e.encodeNil()

	case reflect.Value:
		e.encodeValue(v)

	case string:
		e.e.encodeString(c_UTF8, v)
	case bool:
		e.e.encodeBool(v)
	case int:
		e.e.encodeInt(int64(v))
	case int8:
		e.e.encodeInt(int64(v))
	case int16:
		e.e.encodeInt(int64(v))
	case int32:
		e.e.encodeInt(int64(v))
	case int64:
		e.e.encodeInt(v)
	case uint:
		e.e.encodeUint(uint64(v))
	case uint8:
		e.e.encodeUint(uint64(v))
	case uint16:
		e.e.encodeUint(uint64(v))
	case uint32:
		e.e.encodeUint(uint64(v))
	case uint64:
		e.e.encodeUint(v)
	case float32:
		e.e.encodeFloat32(v)
	case float64:
		e.e.encodeFloat64(v)

	case []interface***REMOVED******REMOVED***:
		e.encSliceIntf(v)
	case []string:
		e.encSliceStr(v)
	case []int64:
		e.encSliceInt64(v)
	case []uint64:
		e.encSliceUint64(v)
	case []uint8:
		e.e.encodeStringBytes(c_RAW, v)

	case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
		e.encMapIntfIntf(v)
	case map[string]interface***REMOVED******REMOVED***:
		e.encMapStrIntf(v)
	case map[string]string:
		e.encMapStrStr(v)
	case map[int64]interface***REMOVED******REMOVED***:
		e.encMapInt64Intf(v)
	case map[uint64]interface***REMOVED******REMOVED***:
		e.encMapUint64Intf(v)

	case *string:
		e.e.encodeString(c_UTF8, *v)
	case *bool:
		e.e.encodeBool(*v)
	case *int:
		e.e.encodeInt(int64(*v))
	case *int8:
		e.e.encodeInt(int64(*v))
	case *int16:
		e.e.encodeInt(int64(*v))
	case *int32:
		e.e.encodeInt(int64(*v))
	case *int64:
		e.e.encodeInt(*v)
	case *uint:
		e.e.encodeUint(uint64(*v))
	case *uint8:
		e.e.encodeUint(uint64(*v))
	case *uint16:
		e.e.encodeUint(uint64(*v))
	case *uint32:
		e.e.encodeUint(uint64(*v))
	case *uint64:
		e.e.encodeUint(*v)
	case *float32:
		e.e.encodeFloat32(*v)
	case *float64:
		e.e.encodeFloat64(*v)

	case *[]interface***REMOVED******REMOVED***:
		e.encSliceIntf(*v)
	case *[]string:
		e.encSliceStr(*v)
	case *[]int64:
		e.encSliceInt64(*v)
	case *[]uint64:
		e.encSliceUint64(*v)
	case *[]uint8:
		e.e.encodeStringBytes(c_RAW, *v)

	case *map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
		e.encMapIntfIntf(*v)
	case *map[string]interface***REMOVED******REMOVED***:
		e.encMapStrIntf(*v)
	case *map[string]string:
		e.encMapStrStr(*v)
	case *map[int64]interface***REMOVED******REMOVED***:
		e.encMapInt64Intf(*v)
	case *map[uint64]interface***REMOVED******REMOVED***:
		e.encMapUint64Intf(*v)

	default:
		e.encodeValue(reflect.ValueOf(iv))
	***REMOVED***
***REMOVED***

func (e *Encoder) encodeValue(rv reflect.Value) ***REMOVED***
	for rv.Kind() == reflect.Ptr ***REMOVED***
		if rv.IsNil() ***REMOVED***
			e.e.encodeNil()
			return
		***REMOVED***
		rv = rv.Elem()
	***REMOVED***

	rt := rv.Type()
	rtid := reflect.ValueOf(rt).Pointer()

	// if e.f == nil && e.s == nil ***REMOVED*** debugf("---->Creating new enc f map for type: %v\n", rt) ***REMOVED***
	var fn encFn
	var ok bool
	if useMapForCodecCache ***REMOVED***
		fn, ok = e.f[rtid]
	***REMOVED*** else ***REMOVED***
		for i, v := range e.x ***REMOVED***
			if v == rtid ***REMOVED***
				fn, ok = e.s[i], true
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if !ok ***REMOVED***
		// debugf("\tCreating new enc fn for type: %v\n", rt)
		fi := encFnInfo***REMOVED***ti: getTypeInfo(rtid, rt), e: e, ee: e.e***REMOVED***
		fn.i = &fi
		if rtid == rawExtTypId ***REMOVED***
			fn.f = (*encFnInfo).rawExt
		***REMOVED*** else if e.e.isBuiltinType(rtid) ***REMOVED***
			fn.f = (*encFnInfo).builtin
		***REMOVED*** else if xfTag, xfFn := e.h.getEncodeExt(rtid); xfFn != nil ***REMOVED***
			fi.xfTag, fi.xfFn = xfTag, xfFn
			fn.f = (*encFnInfo).ext
		***REMOVED*** else if supportBinaryMarshal && fi.ti.m ***REMOVED***
			fn.f = (*encFnInfo).binaryMarshal
		***REMOVED*** else ***REMOVED***
			switch rk := rt.Kind(); rk ***REMOVED***
			case reflect.Bool:
				fn.f = (*encFnInfo).kBool
			case reflect.String:
				fn.f = (*encFnInfo).kString
			case reflect.Float64:
				fn.f = (*encFnInfo).kFloat64
			case reflect.Float32:
				fn.f = (*encFnInfo).kFloat32
			case reflect.Int, reflect.Int8, reflect.Int64, reflect.Int32, reflect.Int16:
				fn.f = (*encFnInfo).kInt
			case reflect.Uint8, reflect.Uint64, reflect.Uint, reflect.Uint32, reflect.Uint16:
				fn.f = (*encFnInfo).kUint
			case reflect.Invalid:
				fn.f = (*encFnInfo).kInvalid
			case reflect.Slice:
				fn.f = (*encFnInfo).kSlice
			case reflect.Array:
				fn.f = (*encFnInfo).kArray
			case reflect.Struct:
				fn.f = (*encFnInfo).kStruct
			// case reflect.Ptr:
			// 	fn.f = (*encFnInfo).kPtr
			case reflect.Interface:
				fn.f = (*encFnInfo).kInterface
			case reflect.Map:
				fn.f = (*encFnInfo).kMap
			default:
				fn.f = (*encFnInfo).kErr
			***REMOVED***
		***REMOVED***
		if useMapForCodecCache ***REMOVED***
			if e.f == nil ***REMOVED***
				e.f = make(map[uintptr]encFn, 16)
			***REMOVED***
			e.f[rtid] = fn
		***REMOVED*** else ***REMOVED***
			e.s = append(e.s, fn)
			e.x = append(e.x, rtid)
		***REMOVED***
	***REMOVED***

	fn.f(fn.i, rv)

***REMOVED***

func (e *Encoder) encRawExt(re RawExt) ***REMOVED***
	if re.Data == nil ***REMOVED***
		e.e.encodeNil()
		return
	***REMOVED***
	if e.hh.writeExt() ***REMOVED***
		e.e.encodeExtPreamble(re.Tag, len(re.Data))
		e.w.writeb(re.Data)
	***REMOVED*** else ***REMOVED***
		e.e.encodeStringBytes(c_RAW, re.Data)
	***REMOVED***
***REMOVED***

// ---------------------------------------------
// short circuit functions for common maps and slices

func (e *Encoder) encSliceIntf(v []interface***REMOVED******REMOVED***) ***REMOVED***
	e.e.encodeArrayPreamble(len(v))
	for _, v2 := range v ***REMOVED***
		e.encode(v2)
	***REMOVED***
***REMOVED***

func (e *Encoder) encSliceStr(v []string) ***REMOVED***
	e.e.encodeArrayPreamble(len(v))
	for _, v2 := range v ***REMOVED***
		e.e.encodeString(c_UTF8, v2)
	***REMOVED***
***REMOVED***

func (e *Encoder) encSliceInt64(v []int64) ***REMOVED***
	e.e.encodeArrayPreamble(len(v))
	for _, v2 := range v ***REMOVED***
		e.e.encodeInt(v2)
	***REMOVED***
***REMOVED***

func (e *Encoder) encSliceUint64(v []uint64) ***REMOVED***
	e.e.encodeArrayPreamble(len(v))
	for _, v2 := range v ***REMOVED***
		e.e.encodeUint(v2)
	***REMOVED***
***REMOVED***

func (e *Encoder) encMapStrStr(v map[string]string) ***REMOVED***
	e.e.encodeMapPreamble(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	for k2, v2 := range v ***REMOVED***
		if asSymbols ***REMOVED***
			e.e.encodeSymbol(k2)
		***REMOVED*** else ***REMOVED***
			e.e.encodeString(c_UTF8, k2)
		***REMOVED***
		e.e.encodeString(c_UTF8, v2)
	***REMOVED***
***REMOVED***

func (e *Encoder) encMapStrIntf(v map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	e.e.encodeMapPreamble(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	for k2, v2 := range v ***REMOVED***
		if asSymbols ***REMOVED***
			e.e.encodeSymbol(k2)
		***REMOVED*** else ***REMOVED***
			e.e.encodeString(c_UTF8, k2)
		***REMOVED***
		e.encode(v2)
	***REMOVED***
***REMOVED***

func (e *Encoder) encMapInt64Intf(v map[int64]interface***REMOVED******REMOVED***) ***REMOVED***
	e.e.encodeMapPreamble(len(v))
	for k2, v2 := range v ***REMOVED***
		e.e.encodeInt(k2)
		e.encode(v2)
	***REMOVED***
***REMOVED***

func (e *Encoder) encMapUint64Intf(v map[uint64]interface***REMOVED******REMOVED***) ***REMOVED***
	e.e.encodeMapPreamble(len(v))
	for k2, v2 := range v ***REMOVED***
		e.e.encodeUint(uint64(k2))
		e.encode(v2)
	***REMOVED***
***REMOVED***

func (e *Encoder) encMapIntfIntf(v map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***) ***REMOVED***
	e.e.encodeMapPreamble(len(v))
	for k2, v2 := range v ***REMOVED***
		e.encode(k2)
		e.encode(v2)
	***REMOVED***
***REMOVED***

// ----------------------------------------

func encErr(format string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
	doPanic(msgTagEnc, format, params...)
***REMOVED***
