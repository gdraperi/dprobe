// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

// Contains code shared by both encode and decode.

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	structTagName = "codec"

	// Support
	//    encoding.BinaryMarshaler: MarshalBinary() (data []byte, err error)
	//    encoding.BinaryUnmarshaler: UnmarshalBinary(data []byte) error
	// This constant flag will enable or disable it.
	supportBinaryMarshal = true

	// Each Encoder or Decoder uses a cache of functions based on conditionals,
	// so that the conditionals are not run every time.
	//
	// Either a map or a slice is used to keep track of the functions.
	// The map is more natural, but has a higher cost than a slice/array.
	// This flag (useMapForCodecCache) controls which is used.
	useMapForCodecCache = false

	// For some common container types, we can short-circuit an elaborate
	// reflection dance and call encode/decode directly.
	// The currently supported types are:
	//    - slices of strings, or id's (int64,uint64) or interfaces.
	//    - maps of str->str, str->intf, id(int64,uint64)->intf, intf->intf
	shortCircuitReflectToFastPath = true

	// for debugging, set this to false, to catch panic traces.
	// Note that this will always cause rpc tests to fail, since they need io.EOF sent via panic.
	recoverPanicToErr = true
)

type charEncoding uint8

const (
	c_RAW charEncoding = iota
	c_UTF8
	c_UTF16LE
	c_UTF16BE
	c_UTF32LE
	c_UTF32BE
)

// valueType is the stream type
type valueType uint8

const (
	valueTypeUnset valueType = iota
	valueTypeNil
	valueTypeInt
	valueTypeUint
	valueTypeFloat
	valueTypeBool
	valueTypeString
	valueTypeSymbol
	valueTypeBytes
	valueTypeMap
	valueTypeArray
	valueTypeTimestamp
	valueTypeExt

	valueTypeInvalid = 0xff
)

var (
	bigen               = binary.BigEndian
	structInfoFieldName = "_struct"

	cachedTypeInfo      = make(map[uintptr]*typeInfo, 4)
	cachedTypeInfoMutex sync.RWMutex

	intfSliceTyp = reflect.TypeOf([]interface***REMOVED******REMOVED***(nil))
	intfTyp      = intfSliceTyp.Elem()

	strSliceTyp     = reflect.TypeOf([]string(nil))
	boolSliceTyp    = reflect.TypeOf([]bool(nil))
	uintSliceTyp    = reflect.TypeOf([]uint(nil))
	uint8SliceTyp   = reflect.TypeOf([]uint8(nil))
	uint16SliceTyp  = reflect.TypeOf([]uint16(nil))
	uint32SliceTyp  = reflect.TypeOf([]uint32(nil))
	uint64SliceTyp  = reflect.TypeOf([]uint64(nil))
	intSliceTyp     = reflect.TypeOf([]int(nil))
	int8SliceTyp    = reflect.TypeOf([]int8(nil))
	int16SliceTyp   = reflect.TypeOf([]int16(nil))
	int32SliceTyp   = reflect.TypeOf([]int32(nil))
	int64SliceTyp   = reflect.TypeOf([]int64(nil))
	float32SliceTyp = reflect.TypeOf([]float32(nil))
	float64SliceTyp = reflect.TypeOf([]float64(nil))

	mapIntfIntfTyp = reflect.TypeOf(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***(nil))
	mapStrIntfTyp  = reflect.TypeOf(map[string]interface***REMOVED******REMOVED***(nil))
	mapStrStrTyp   = reflect.TypeOf(map[string]string(nil))

	mapIntIntfTyp    = reflect.TypeOf(map[int]interface***REMOVED******REMOVED***(nil))
	mapInt64IntfTyp  = reflect.TypeOf(map[int64]interface***REMOVED******REMOVED***(nil))
	mapUintIntfTyp   = reflect.TypeOf(map[uint]interface***REMOVED******REMOVED***(nil))
	mapUint64IntfTyp = reflect.TypeOf(map[uint64]interface***REMOVED******REMOVED***(nil))

	stringTyp = reflect.TypeOf("")
	timeTyp   = reflect.TypeOf(time.Time***REMOVED******REMOVED***)
	rawExtTyp = reflect.TypeOf(RawExt***REMOVED******REMOVED***)

	mapBySliceTyp        = reflect.TypeOf((*MapBySlice)(nil)).Elem()
	binaryMarshalerTyp   = reflect.TypeOf((*binaryMarshaler)(nil)).Elem()
	binaryUnmarshalerTyp = reflect.TypeOf((*binaryUnmarshaler)(nil)).Elem()

	rawExtTypId = reflect.ValueOf(rawExtTyp).Pointer()
	intfTypId   = reflect.ValueOf(intfTyp).Pointer()
	timeTypId   = reflect.ValueOf(timeTyp).Pointer()

	intfSliceTypId = reflect.ValueOf(intfSliceTyp).Pointer()
	strSliceTypId  = reflect.ValueOf(strSliceTyp).Pointer()

	boolSliceTypId    = reflect.ValueOf(boolSliceTyp).Pointer()
	uintSliceTypId    = reflect.ValueOf(uintSliceTyp).Pointer()
	uint8SliceTypId   = reflect.ValueOf(uint8SliceTyp).Pointer()
	uint16SliceTypId  = reflect.ValueOf(uint16SliceTyp).Pointer()
	uint32SliceTypId  = reflect.ValueOf(uint32SliceTyp).Pointer()
	uint64SliceTypId  = reflect.ValueOf(uint64SliceTyp).Pointer()
	intSliceTypId     = reflect.ValueOf(intSliceTyp).Pointer()
	int8SliceTypId    = reflect.ValueOf(int8SliceTyp).Pointer()
	int16SliceTypId   = reflect.ValueOf(int16SliceTyp).Pointer()
	int32SliceTypId   = reflect.ValueOf(int32SliceTyp).Pointer()
	int64SliceTypId   = reflect.ValueOf(int64SliceTyp).Pointer()
	float32SliceTypId = reflect.ValueOf(float32SliceTyp).Pointer()
	float64SliceTypId = reflect.ValueOf(float64SliceTyp).Pointer()

	mapStrStrTypId     = reflect.ValueOf(mapStrStrTyp).Pointer()
	mapIntfIntfTypId   = reflect.ValueOf(mapIntfIntfTyp).Pointer()
	mapStrIntfTypId    = reflect.ValueOf(mapStrIntfTyp).Pointer()
	mapIntIntfTypId    = reflect.ValueOf(mapIntIntfTyp).Pointer()
	mapInt64IntfTypId  = reflect.ValueOf(mapInt64IntfTyp).Pointer()
	mapUintIntfTypId   = reflect.ValueOf(mapUintIntfTyp).Pointer()
	mapUint64IntfTypId = reflect.ValueOf(mapUint64IntfTyp).Pointer()
	// Id = reflect.ValueOf().Pointer()
	// mapBySliceTypId  = reflect.ValueOf(mapBySliceTyp).Pointer()

	binaryMarshalerTypId   = reflect.ValueOf(binaryMarshalerTyp).Pointer()
	binaryUnmarshalerTypId = reflect.ValueOf(binaryUnmarshalerTyp).Pointer()

	intBitsize  uint8 = uint8(reflect.TypeOf(int(0)).Bits())
	uintBitsize uint8 = uint8(reflect.TypeOf(uint(0)).Bits())

	bsAll0x00 = []byte***REMOVED***0, 0, 0, 0, 0, 0, 0, 0***REMOVED***
	bsAll0xff = []byte***REMOVED***0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff***REMOVED***
)

type binaryUnmarshaler interface ***REMOVED***
	UnmarshalBinary(data []byte) error
***REMOVED***

type binaryMarshaler interface ***REMOVED***
	MarshalBinary() (data []byte, err error)
***REMOVED***

// MapBySlice represents a slice which should be encoded as a map in the stream.
// The slice contains a sequence of key-value pairs.
type MapBySlice interface ***REMOVED***
	MapBySlice()
***REMOVED***

// WARNING: DO NOT USE DIRECTLY. EXPORTED FOR GODOC BENEFIT. WILL BE REMOVED.
//
// BasicHandle encapsulates the common options and extension functions.
type BasicHandle struct ***REMOVED***
	extHandle
	EncodeOptions
	DecodeOptions
***REMOVED***

// Handle is the interface for a specific encoding format.
//
// Typically, a Handle is pre-configured before first time use,
// and not modified while in use. Such a pre-configured Handle
// is safe for concurrent access.
type Handle interface ***REMOVED***
	writeExt() bool
	getBasicHandle() *BasicHandle
	newEncDriver(w encWriter) encDriver
	newDecDriver(r decReader) decDriver
***REMOVED***

// RawExt represents raw unprocessed extension data.
type RawExt struct ***REMOVED***
	Tag  byte
	Data []byte
***REMOVED***

type extTypeTagFn struct ***REMOVED***
	rtid  uintptr
	rt    reflect.Type
	tag   byte
	encFn func(reflect.Value) ([]byte, error)
	decFn func(reflect.Value, []byte) error
***REMOVED***

type extHandle []*extTypeTagFn

// AddExt registers an encode and decode function for a reflect.Type.
// Note that the type must be a named type, and specifically not
// a pointer or Interface. An error is returned if that is not honored.
//
// To Deregister an ext, call AddExt with 0 tag, nil encfn and nil decfn.
func (o *extHandle) AddExt(
	rt reflect.Type,
	tag byte,
	encfn func(reflect.Value) ([]byte, error),
	decfn func(reflect.Value, []byte) error,
) (err error) ***REMOVED***
	// o is a pointer, because we may need to initialize it
	if rt.PkgPath() == "" || rt.Kind() == reflect.Interface ***REMOVED***
		err = fmt.Errorf("codec.Handle.AddExt: Takes named type, especially not a pointer or interface: %T",
			reflect.Zero(rt).Interface())
		return
	***REMOVED***

	// o cannot be nil, since it is always embedded in a Handle.
	// if nil, let it panic.
	// if o == nil ***REMOVED***
	// 	err = errors.New("codec.Handle.AddExt: extHandle cannot be a nil pointer.")
	// 	return
	// ***REMOVED***

	rtid := reflect.ValueOf(rt).Pointer()
	for _, v := range *o ***REMOVED***
		if v.rtid == rtid ***REMOVED***
			v.tag, v.encFn, v.decFn = tag, encfn, decfn
			return
		***REMOVED***
	***REMOVED***

	*o = append(*o, &extTypeTagFn***REMOVED***rtid, rt, tag, encfn, decfn***REMOVED***)
	return
***REMOVED***

func (o extHandle) getExt(rtid uintptr) *extTypeTagFn ***REMOVED***
	for _, v := range o ***REMOVED***
		if v.rtid == rtid ***REMOVED***
			return v
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (o extHandle) getExtForTag(tag byte) *extTypeTagFn ***REMOVED***
	for _, v := range o ***REMOVED***
		if v.tag == tag ***REMOVED***
			return v
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (o extHandle) getDecodeExtForTag(tag byte) (
	rv reflect.Value, fn func(reflect.Value, []byte) error) ***REMOVED***
	if x := o.getExtForTag(tag); x != nil ***REMOVED***
		// ext is only registered for base
		rv = reflect.New(x.rt).Elem()
		fn = x.decFn
	***REMOVED***
	return
***REMOVED***

func (o extHandle) getDecodeExt(rtid uintptr) (tag byte, fn func(reflect.Value, []byte) error) ***REMOVED***
	if x := o.getExt(rtid); x != nil ***REMOVED***
		tag = x.tag
		fn = x.decFn
	***REMOVED***
	return
***REMOVED***

func (o extHandle) getEncodeExt(rtid uintptr) (tag byte, fn func(reflect.Value) ([]byte, error)) ***REMOVED***
	if x := o.getExt(rtid); x != nil ***REMOVED***
		tag = x.tag
		fn = x.encFn
	***REMOVED***
	return
***REMOVED***

type structFieldInfo struct ***REMOVED***
	encName string // encode name

	// only one of 'i' or 'is' can be set. If 'i' is -1, then 'is' has been set.

	is        []int // (recursive/embedded) field index in struct
	i         int16 // field index in struct
	omitEmpty bool
	toArray   bool // if field is _struct, is the toArray set?

	// tag       string   // tag
	// name      string   // field name
	// encNameBs []byte   // encoded name as byte stream
	// ikind     int      // kind of the field as an int i.e. int(reflect.Kind)
***REMOVED***

func parseStructFieldInfo(fname string, stag string) *structFieldInfo ***REMOVED***
	if fname == "" ***REMOVED***
		panic("parseStructFieldInfo: No Field Name")
	***REMOVED***
	si := structFieldInfo***REMOVED***
		// name: fname,
		encName: fname,
		// tag: stag,
	***REMOVED***

	if stag != "" ***REMOVED***
		for i, s := range strings.Split(stag, ",") ***REMOVED***
			if i == 0 ***REMOVED***
				if s != "" ***REMOVED***
					si.encName = s
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				switch s ***REMOVED***
				case "omitempty":
					si.omitEmpty = true
				case "toarray":
					si.toArray = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// si.encNameBs = []byte(si.encName)
	return &si
***REMOVED***

type sfiSortedByEncName []*structFieldInfo

func (p sfiSortedByEncName) Len() int ***REMOVED***
	return len(p)
***REMOVED***

func (p sfiSortedByEncName) Less(i, j int) bool ***REMOVED***
	return p[i].encName < p[j].encName
***REMOVED***

func (p sfiSortedByEncName) Swap(i, j int) ***REMOVED***
	p[i], p[j] = p[j], p[i]
***REMOVED***

// typeInfo keeps information about each type referenced in the encode/decode sequence.
//
// During an encode/decode sequence, we work as below:
//   - If base is a built in type, en/decode base value
//   - If base is registered as an extension, en/decode base value
//   - If type is binary(M/Unm)arshaler, call Binary(M/Unm)arshal method
//   - Else decode appropriately based on the reflect.Kind
type typeInfo struct ***REMOVED***
	sfi  []*structFieldInfo // sorted. Used when enc/dec struct to map.
	sfip []*structFieldInfo // unsorted. Used when enc/dec struct to array.

	rt   reflect.Type
	rtid uintptr

	// baseId gives pointer to the base reflect.Type, after deferencing
	// the pointers. E.g. base type of ***time.Time is time.Time.
	base      reflect.Type
	baseId    uintptr
	baseIndir int8 // number of indirections to get to base

	mbs bool // base type (T or *T) is a MapBySlice

	m        bool // base type (T or *T) is a binaryMarshaler
	unm      bool // base type (T or *T) is a binaryUnmarshaler
	mIndir   int8 // number of indirections to get to binaryMarshaler type
	unmIndir int8 // number of indirections to get to binaryUnmarshaler type
	toArray  bool // whether this (struct) type should be encoded as an array
***REMOVED***

func (ti *typeInfo) indexForEncName(name string) int ***REMOVED***
	//tisfi := ti.sfi
	const binarySearchThreshold = 16
	if sfilen := len(ti.sfi); sfilen < binarySearchThreshold ***REMOVED***
		// linear search. faster than binary search in my testing up to 16-field structs.
		for i, si := range ti.sfi ***REMOVED***
			if si.encName == name ***REMOVED***
				return i
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// binary search. adapted from sort/search.go.
		h, i, j := 0, 0, sfilen
		for i < j ***REMOVED***
			h = i + (j-i)/2
			if ti.sfi[h].encName < name ***REMOVED***
				i = h + 1
			***REMOVED*** else ***REMOVED***
				j = h
			***REMOVED***
		***REMOVED***
		if i < sfilen && ti.sfi[i].encName == name ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

func getTypeInfo(rtid uintptr, rt reflect.Type) (pti *typeInfo) ***REMOVED***
	var ok bool
	cachedTypeInfoMutex.RLock()
	pti, ok = cachedTypeInfo[rtid]
	cachedTypeInfoMutex.RUnlock()
	if ok ***REMOVED***
		return
	***REMOVED***

	cachedTypeInfoMutex.Lock()
	defer cachedTypeInfoMutex.Unlock()
	if pti, ok = cachedTypeInfo[rtid]; ok ***REMOVED***
		return
	***REMOVED***

	ti := typeInfo***REMOVED***rt: rt, rtid: rtid***REMOVED***
	pti = &ti

	var indir int8
	if ok, indir = implementsIntf(rt, binaryMarshalerTyp); ok ***REMOVED***
		ti.m, ti.mIndir = true, indir
	***REMOVED***
	if ok, indir = implementsIntf(rt, binaryUnmarshalerTyp); ok ***REMOVED***
		ti.unm, ti.unmIndir = true, indir
	***REMOVED***
	if ok, _ = implementsIntf(rt, mapBySliceTyp); ok ***REMOVED***
		ti.mbs = true
	***REMOVED***

	pt := rt
	var ptIndir int8
	// for ; pt.Kind() == reflect.Ptr; pt, ptIndir = pt.Elem(), ptIndir+1 ***REMOVED*** ***REMOVED***
	for pt.Kind() == reflect.Ptr ***REMOVED***
		pt = pt.Elem()
		ptIndir++
	***REMOVED***
	if ptIndir == 0 ***REMOVED***
		ti.base = rt
		ti.baseId = rtid
	***REMOVED*** else ***REMOVED***
		ti.base = pt
		ti.baseId = reflect.ValueOf(pt).Pointer()
		ti.baseIndir = ptIndir
	***REMOVED***

	if rt.Kind() == reflect.Struct ***REMOVED***
		var siInfo *structFieldInfo
		if f, ok := rt.FieldByName(structInfoFieldName); ok ***REMOVED***
			siInfo = parseStructFieldInfo(structInfoFieldName, f.Tag.Get(structTagName))
			ti.toArray = siInfo.toArray
		***REMOVED***
		sfip := make([]*structFieldInfo, 0, rt.NumField())
		rgetTypeInfo(rt, nil, make(map[string]bool), &sfip, siInfo)

		// // try to put all si close together
		// const tryToPutAllStructFieldInfoTogether = true
		// if tryToPutAllStructFieldInfoTogether ***REMOVED***
		// 	sfip2 := make([]structFieldInfo, len(sfip))
		// 	for i, si := range sfip ***REMOVED***
		// 		sfip2[i] = *si
		// 	***REMOVED***
		// 	for i := range sfip ***REMOVED***
		// 		sfip[i] = &sfip2[i]
		// 	***REMOVED***
		// ***REMOVED***

		ti.sfip = make([]*structFieldInfo, len(sfip))
		ti.sfi = make([]*structFieldInfo, len(sfip))
		copy(ti.sfip, sfip)
		sort.Sort(sfiSortedByEncName(sfip))
		copy(ti.sfi, sfip)
	***REMOVED***
	// sfi = sfip
	cachedTypeInfo[rtid] = pti
	return
***REMOVED***

func rgetTypeInfo(rt reflect.Type, indexstack []int, fnameToHastag map[string]bool,
	sfi *[]*structFieldInfo, siInfo *structFieldInfo,
) ***REMOVED***
	// for rt.Kind() == reflect.Ptr ***REMOVED***
	// 	// indexstack = append(indexstack, 0)
	// 	rt = rt.Elem()
	// ***REMOVED***
	for j := 0; j < rt.NumField(); j++ ***REMOVED***
		f := rt.Field(j)
		stag := f.Tag.Get(structTagName)
		if stag == "-" ***REMOVED***
			continue
		***REMOVED***
		if r1, _ := utf8.DecodeRuneInString(f.Name); r1 == utf8.RuneError || !unicode.IsUpper(r1) ***REMOVED***
			continue
		***REMOVED***
		// if anonymous and there is no struct tag and its a struct (or pointer to struct), inline it.
		if f.Anonymous && stag == "" ***REMOVED***
			ft := f.Type
			for ft.Kind() == reflect.Ptr ***REMOVED***
				ft = ft.Elem()
			***REMOVED***
			if ft.Kind() == reflect.Struct ***REMOVED***
				indexstack2 := append(append(make([]int, 0, len(indexstack)+4), indexstack...), j)
				rgetTypeInfo(ft, indexstack2, fnameToHastag, sfi, siInfo)
				continue
			***REMOVED***
		***REMOVED***
		// do not let fields with same name in embedded structs override field at higher level.
		// this must be done after anonymous check, to allow anonymous field
		// still include their child fields
		if _, ok := fnameToHastag[f.Name]; ok ***REMOVED***
			continue
		***REMOVED***
		si := parseStructFieldInfo(f.Name, stag)
		// si.ikind = int(f.Type.Kind())
		if len(indexstack) == 0 ***REMOVED***
			si.i = int16(j)
		***REMOVED*** else ***REMOVED***
			si.i = -1
			si.is = append(append(make([]int, 0, len(indexstack)+4), indexstack...), j)
		***REMOVED***

		if siInfo != nil ***REMOVED***
			if siInfo.omitEmpty ***REMOVED***
				si.omitEmpty = true
			***REMOVED***
		***REMOVED***
		*sfi = append(*sfi, si)
		fnameToHastag[f.Name] = stag != ""
	***REMOVED***
***REMOVED***

func panicToErr(err *error) ***REMOVED***
	if recoverPanicToErr ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			//debug.PrintStack()
			panicValToErr(x, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func doPanic(tag string, format string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
	params2 := make([]interface***REMOVED******REMOVED***, len(params)+1)
	params2[0] = tag
	copy(params2[1:], params)
	panic(fmt.Errorf("%s: "+format, params2...))
***REMOVED***

func checkOverflowFloat32(f float64, doCheck bool) ***REMOVED***
	if !doCheck ***REMOVED***
		return
	***REMOVED***
	// check overflow (logic adapted from std pkg reflect/value.go OverflowFloat()
	f2 := f
	if f2 < 0 ***REMOVED***
		f2 = -f
	***REMOVED***
	if math.MaxFloat32 < f2 && f2 <= math.MaxFloat64 ***REMOVED***
		decErr("Overflow float32 value: %v", f2)
	***REMOVED***
***REMOVED***

func checkOverflow(ui uint64, i int64, bitsize uint8) ***REMOVED***
	// check overflow (logic adapted from std pkg reflect/value.go OverflowUint()
	if bitsize == 0 ***REMOVED***
		return
	***REMOVED***
	if i != 0 ***REMOVED***
		if trunc := (i << (64 - bitsize)) >> (64 - bitsize); i != trunc ***REMOVED***
			decErr("Overflow int value: %v", i)
		***REMOVED***
	***REMOVED***
	if ui != 0 ***REMOVED***
		if trunc := (ui << (64 - bitsize)) >> (64 - bitsize); ui != trunc ***REMOVED***
			decErr("Overflow uint value: %v", ui)
		***REMOVED***
	***REMOVED***
***REMOVED***
