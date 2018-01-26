// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"math"
	// "reflect"
	// "sync/atomic"
	"time"
	//"fmt"
)

const bincDoPrune = true // No longer needed. Needed before as C lib did not support pruning.

//var _ = fmt.Printf

// vd as low 4 bits (there are 16 slots)
const (
	bincVdSpecial byte = iota
	bincVdPosInt
	bincVdNegInt
	bincVdFloat

	bincVdString
	bincVdByteArray
	bincVdArray
	bincVdMap

	bincVdTimestamp
	bincVdSmallInt
	bincVdUnicodeOther
	bincVdSymbol

	bincVdDecimal
	_               // open slot
	_               // open slot
	bincVdCustomExt = 0x0f
)

const (
	bincSpNil byte = iota
	bincSpFalse
	bincSpTrue
	bincSpNan
	bincSpPosInf
	bincSpNegInf
	bincSpZeroFloat
	bincSpZero
	bincSpNegOne
)

const (
	bincFlBin16 byte = iota
	bincFlBin32
	_ // bincFlBin32e
	bincFlBin64
	_ // bincFlBin64e
	// others not currently supported
)

type bincEncDriver struct ***REMOVED***
	w encWriter
	m map[string]uint16 // symbols
	s uint32            // symbols sequencer
	b [8]byte
***REMOVED***

func (e *bincEncDriver) isBuiltinType(rt uintptr) bool ***REMOVED***
	return rt == timeTypId
***REMOVED***

func (e *bincEncDriver) encodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED***
	switch rt ***REMOVED***
	case timeTypId:
		bs := encodeTime(v.(time.Time))
		e.w.writen1(bincVdTimestamp<<4 | uint8(len(bs)))
		e.w.writeb(bs)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encodeNil() ***REMOVED***
	e.w.writen1(bincVdSpecial<<4 | bincSpNil)
***REMOVED***

func (e *bincEncDriver) encodeBool(b bool) ***REMOVED***
	if b ***REMOVED***
		e.w.writen1(bincVdSpecial<<4 | bincSpTrue)
	***REMOVED*** else ***REMOVED***
		e.w.writen1(bincVdSpecial<<4 | bincSpFalse)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encodeFloat32(f float32) ***REMOVED***
	if f == 0 ***REMOVED***
		e.w.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
		return
	***REMOVED***
	e.w.writen1(bincVdFloat<<4 | bincFlBin32)
	e.w.writeUint32(math.Float32bits(f))
***REMOVED***

func (e *bincEncDriver) encodeFloat64(f float64) ***REMOVED***
	if f == 0 ***REMOVED***
		e.w.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
		return
	***REMOVED***
	bigen.PutUint64(e.b[:], math.Float64bits(f))
	if bincDoPrune ***REMOVED***
		i := 7
		for ; i >= 0 && (e.b[i] == 0); i-- ***REMOVED***
		***REMOVED***
		i++
		if i <= 6 ***REMOVED***
			e.w.writen1(bincVdFloat<<4 | 0x8 | bincFlBin64)
			e.w.writen1(byte(i))
			e.w.writeb(e.b[:i])
			return
		***REMOVED***
	***REMOVED***
	e.w.writen1(bincVdFloat<<4 | bincFlBin64)
	e.w.writeb(e.b[:])
***REMOVED***

func (e *bincEncDriver) encIntegerPrune(bd byte, pos bool, v uint64, lim uint8) ***REMOVED***
	if lim == 4 ***REMOVED***
		bigen.PutUint32(e.b[:lim], uint32(v))
	***REMOVED*** else ***REMOVED***
		bigen.PutUint64(e.b[:lim], v)
	***REMOVED***
	if bincDoPrune ***REMOVED***
		i := pruneSignExt(e.b[:lim], pos)
		e.w.writen1(bd | lim - 1 - byte(i))
		e.w.writeb(e.b[i:lim])
	***REMOVED*** else ***REMOVED***
		e.w.writen1(bd | lim - 1)
		e.w.writeb(e.b[:lim])
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encodeInt(v int64) ***REMOVED***
	const nbd byte = bincVdNegInt << 4
	switch ***REMOVED***
	case v >= 0:
		e.encUint(bincVdPosInt<<4, true, uint64(v))
	case v == -1:
		e.w.writen1(bincVdSpecial<<4 | bincSpNegOne)
	default:
		e.encUint(bincVdNegInt<<4, false, uint64(-v))
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encodeUint(v uint64) ***REMOVED***
	e.encUint(bincVdPosInt<<4, true, v)
***REMOVED***

func (e *bincEncDriver) encUint(bd byte, pos bool, v uint64) ***REMOVED***
	switch ***REMOVED***
	case v == 0:
		e.w.writen1(bincVdSpecial<<4 | bincSpZero)
	case pos && v >= 1 && v <= 16:
		e.w.writen1(bincVdSmallInt<<4 | byte(v-1))
	case v <= math.MaxUint8:
		e.w.writen2(bd|0x0, byte(v))
	case v <= math.MaxUint16:
		e.w.writen1(bd | 0x01)
		e.w.writeUint16(uint16(v))
	case v <= math.MaxUint32:
		e.encIntegerPrune(bd, pos, v, 4)
	default:
		e.encIntegerPrune(bd, pos, v, 8)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encodeExtPreamble(xtag byte, length int) ***REMOVED***
	e.encLen(bincVdCustomExt<<4, uint64(length))
	e.w.writen1(xtag)
***REMOVED***

func (e *bincEncDriver) encodeArrayPreamble(length int) ***REMOVED***
	e.encLen(bincVdArray<<4, uint64(length))
***REMOVED***

func (e *bincEncDriver) encodeMapPreamble(length int) ***REMOVED***
	e.encLen(bincVdMap<<4, uint64(length))
***REMOVED***

func (e *bincEncDriver) encodeString(c charEncoding, v string) ***REMOVED***
	l := uint64(len(v))
	e.encBytesLen(c, l)
	if l > 0 ***REMOVED***
		e.w.writestr(v)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encodeSymbol(v string) ***REMOVED***
	// if WriteSymbolsNoRefs ***REMOVED***
	// 	e.encodeString(c_UTF8, v)
	// 	return
	// ***REMOVED***

	//symbols only offer benefit when string length > 1.
	//This is because strings with length 1 take only 2 bytes to store
	//(bd with embedded length, and single byte for string val).

	l := len(v)
	switch l ***REMOVED***
	case 0:
		e.encBytesLen(c_UTF8, 0)
		return
	case 1:
		e.encBytesLen(c_UTF8, 1)
		e.w.writen1(v[0])
		return
	***REMOVED***
	if e.m == nil ***REMOVED***
		e.m = make(map[string]uint16, 16)
	***REMOVED***
	ui, ok := e.m[v]
	if ok ***REMOVED***
		if ui <= math.MaxUint8 ***REMOVED***
			e.w.writen2(bincVdSymbol<<4, byte(ui))
		***REMOVED*** else ***REMOVED***
			e.w.writen1(bincVdSymbol<<4 | 0x8)
			e.w.writeUint16(ui)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		e.s++
		ui = uint16(e.s)
		//ui = uint16(atomic.AddUint32(&e.s, 1))
		e.m[v] = ui
		var lenprec uint8
		switch ***REMOVED***
		case l <= math.MaxUint8:
			// lenprec = 0
		case l <= math.MaxUint16:
			lenprec = 1
		case int64(l) <= math.MaxUint32:
			lenprec = 2
		default:
			lenprec = 3
		***REMOVED***
		if ui <= math.MaxUint8 ***REMOVED***
			e.w.writen2(bincVdSymbol<<4|0x0|0x4|lenprec, byte(ui))
		***REMOVED*** else ***REMOVED***
			e.w.writen1(bincVdSymbol<<4 | 0x8 | 0x4 | lenprec)
			e.w.writeUint16(ui)
		***REMOVED***
		switch lenprec ***REMOVED***
		case 0:
			e.w.writen1(byte(l))
		case 1:
			e.w.writeUint16(uint16(l))
		case 2:
			e.w.writeUint32(uint32(l))
		default:
			e.w.writeUint64(uint64(l))
		***REMOVED***
		e.w.writestr(v)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encodeStringBytes(c charEncoding, v []byte) ***REMOVED***
	l := uint64(len(v))
	e.encBytesLen(c, l)
	if l > 0 ***REMOVED***
		e.w.writeb(v)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encBytesLen(c charEncoding, length uint64) ***REMOVED***
	//TODO: support bincUnicodeOther (for now, just use string or bytearray)
	if c == c_RAW ***REMOVED***
		e.encLen(bincVdByteArray<<4, length)
	***REMOVED*** else ***REMOVED***
		e.encLen(bincVdString<<4, length)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encLen(bd byte, l uint64) ***REMOVED***
	if l < 12 ***REMOVED***
		e.w.writen1(bd | uint8(l+4))
	***REMOVED*** else ***REMOVED***
		e.encLenNumber(bd, l)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encLenNumber(bd byte, v uint64) ***REMOVED***
	switch ***REMOVED***
	case v <= math.MaxUint8:
		e.w.writen2(bd, byte(v))
	case v <= math.MaxUint16:
		e.w.writen1(bd | 0x01)
		e.w.writeUint16(uint16(v))
	case v <= math.MaxUint32:
		e.w.writen1(bd | 0x02)
		e.w.writeUint32(uint32(v))
	default:
		e.w.writen1(bd | 0x03)
		e.w.writeUint64(uint64(v))
	***REMOVED***
***REMOVED***

//------------------------------------

type bincDecDriver struct ***REMOVED***
	r      decReader
	bdRead bool
	bdType valueType
	bd     byte
	vd     byte
	vs     byte
	b      [8]byte
	m      map[uint32]string // symbols (use uint32 as key, as map optimizes for it)
***REMOVED***

func (d *bincDecDriver) initReadNext() ***REMOVED***
	if d.bdRead ***REMOVED***
		return
	***REMOVED***
	d.bd = d.r.readn1()
	d.vd = d.bd >> 4
	d.vs = d.bd & 0x0f
	d.bdRead = true
	d.bdType = valueTypeUnset
***REMOVED***

func (d *bincDecDriver) currentEncodedType() valueType ***REMOVED***
	if d.bdType == valueTypeUnset ***REMOVED***
		switch d.vd ***REMOVED***
		case bincVdSpecial:
			switch d.vs ***REMOVED***
			case bincSpNil:
				d.bdType = valueTypeNil
			case bincSpFalse, bincSpTrue:
				d.bdType = valueTypeBool
			case bincSpNan, bincSpNegInf, bincSpPosInf, bincSpZeroFloat:
				d.bdType = valueTypeFloat
			case bincSpZero:
				d.bdType = valueTypeUint
			case bincSpNegOne:
				d.bdType = valueTypeInt
			default:
				decErr("currentEncodedType: Unrecognized special value 0x%x", d.vs)
			***REMOVED***
		case bincVdSmallInt:
			d.bdType = valueTypeUint
		case bincVdPosInt:
			d.bdType = valueTypeUint
		case bincVdNegInt:
			d.bdType = valueTypeInt
		case bincVdFloat:
			d.bdType = valueTypeFloat
		case bincVdString:
			d.bdType = valueTypeString
		case bincVdSymbol:
			d.bdType = valueTypeSymbol
		case bincVdByteArray:
			d.bdType = valueTypeBytes
		case bincVdTimestamp:
			d.bdType = valueTypeTimestamp
		case bincVdCustomExt:
			d.bdType = valueTypeExt
		case bincVdArray:
			d.bdType = valueTypeArray
		case bincVdMap:
			d.bdType = valueTypeMap
		default:
			decErr("currentEncodedType: Unrecognized d.vd: 0x%x", d.vd)
		***REMOVED***
	***REMOVED***
	return d.bdType
***REMOVED***

func (d *bincDecDriver) tryDecodeAsNil() bool ***REMOVED***
	if d.bd == bincVdSpecial<<4|bincSpNil ***REMOVED***
		d.bdRead = false
		return true
	***REMOVED***
	return false
***REMOVED***

func (d *bincDecDriver) isBuiltinType(rt uintptr) bool ***REMOVED***
	return rt == timeTypId
***REMOVED***

func (d *bincDecDriver) decodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED***
	switch rt ***REMOVED***
	case timeTypId:
		if d.vd != bincVdTimestamp ***REMOVED***
			decErr("Invalid d.vd. Expecting 0x%x. Received: 0x%x", bincVdTimestamp, d.vd)
		***REMOVED***
		tt, err := decodeTime(d.r.readn(int(d.vs)))
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		var vt *time.Time = v.(*time.Time)
		*vt = tt
		d.bdRead = false
	***REMOVED***
***REMOVED***

func (d *bincDecDriver) decFloatPre(vs, defaultLen byte) ***REMOVED***
	if vs&0x8 == 0 ***REMOVED***
		d.r.readb(d.b[0:defaultLen])
	***REMOVED*** else ***REMOVED***
		l := d.r.readn1()
		if l > 8 ***REMOVED***
			decErr("At most 8 bytes used to represent float. Received: %v bytes", l)
		***REMOVED***
		for i := l; i < 8; i++ ***REMOVED***
			d.b[i] = 0
		***REMOVED***
		d.r.readb(d.b[0:l])
	***REMOVED***
***REMOVED***

func (d *bincDecDriver) decFloat() (f float64) ***REMOVED***
	//if true ***REMOVED*** f = math.Float64frombits(d.r.readUint64()); break; ***REMOVED***
	switch vs := d.vs; vs & 0x7 ***REMOVED***
	case bincFlBin32:
		d.decFloatPre(vs, 4)
		f = float64(math.Float32frombits(bigen.Uint32(d.b[0:4])))
	case bincFlBin64:
		d.decFloatPre(vs, 8)
		f = math.Float64frombits(bigen.Uint64(d.b[0:8]))
	default:
		decErr("only float32 and float64 are supported. d.vd: 0x%x, d.vs: 0x%x", d.vd, d.vs)
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) decUint() (v uint64) ***REMOVED***
	// need to inline the code (interface conversion and type assertion expensive)
	switch d.vs ***REMOVED***
	case 0:
		v = uint64(d.r.readn1())
	case 1:
		d.r.readb(d.b[6:])
		v = uint64(bigen.Uint16(d.b[6:]))
	case 2:
		d.b[4] = 0
		d.r.readb(d.b[5:])
		v = uint64(bigen.Uint32(d.b[4:]))
	case 3:
		d.r.readb(d.b[4:])
		v = uint64(bigen.Uint32(d.b[4:]))
	case 4, 5, 6:
		lim := int(7 - d.vs)
		d.r.readb(d.b[lim:])
		for i := 0; i < lim; i++ ***REMOVED***
			d.b[i] = 0
		***REMOVED***
		v = uint64(bigen.Uint64(d.b[:]))
	case 7:
		d.r.readb(d.b[:])
		v = uint64(bigen.Uint64(d.b[:]))
	default:
		decErr("unsigned integers with greater than 64 bits of precision not supported")
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) decIntAny() (ui uint64, i int64, neg bool) ***REMOVED***
	switch d.vd ***REMOVED***
	case bincVdPosInt:
		ui = d.decUint()
		i = int64(ui)
	case bincVdNegInt:
		ui = d.decUint()
		i = -(int64(ui))
		neg = true
	case bincVdSmallInt:
		i = int64(d.vs) + 1
		ui = uint64(d.vs) + 1
	case bincVdSpecial:
		switch d.vs ***REMOVED***
		case bincSpZero:
			//i = 0
		case bincSpNegOne:
			neg = true
			ui = 1
			i = -1
		default:
			decErr("numeric decode fails for special value: d.vs: 0x%x", d.vs)
		***REMOVED***
	default:
		decErr("number can only be decoded from uint or int values. d.bd: 0x%x, d.vd: 0x%x", d.bd, d.vd)
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) decodeInt(bitsize uint8) (i int64) ***REMOVED***
	_, i, _ = d.decIntAny()
	checkOverflow(0, i, bitsize)
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) decodeUint(bitsize uint8) (ui uint64) ***REMOVED***
	ui, i, neg := d.decIntAny()
	if neg ***REMOVED***
		decErr("Assigning negative signed value: %v, to unsigned type", i)
	***REMOVED***
	checkOverflow(ui, 0, bitsize)
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) decodeFloat(chkOverflow32 bool) (f float64) ***REMOVED***
	switch d.vd ***REMOVED***
	case bincVdSpecial:
		d.bdRead = false
		switch d.vs ***REMOVED***
		case bincSpNan:
			return math.NaN()
		case bincSpPosInf:
			return math.Inf(1)
		case bincSpZeroFloat, bincSpZero:
			return
		case bincSpNegInf:
			return math.Inf(-1)
		default:
			decErr("Invalid d.vs decoding float where d.vd=bincVdSpecial: %v", d.vs)
		***REMOVED***
	case bincVdFloat:
		f = d.decFloat()
	default:
		_, i, _ := d.decIntAny()
		f = float64(i)
	***REMOVED***
	checkOverflowFloat32(f, chkOverflow32)
	d.bdRead = false
	return
***REMOVED***

// bool can be decoded from bool only (single byte).
func (d *bincDecDriver) decodeBool() (b bool) ***REMOVED***
	switch d.bd ***REMOVED***
	case (bincVdSpecial | bincSpFalse):
		// b = false
	case (bincVdSpecial | bincSpTrue):
		b = true
	default:
		decErr("Invalid single-byte value for bool: %s: %x", msgBadDesc, d.bd)
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) readMapLen() (length int) ***REMOVED***
	if d.vd != bincVdMap ***REMOVED***
		decErr("Invalid d.vd for map. Expecting 0x%x. Got: 0x%x", bincVdMap, d.vd)
	***REMOVED***
	length = d.decLen()
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) readArrayLen() (length int) ***REMOVED***
	if d.vd != bincVdArray ***REMOVED***
		decErr("Invalid d.vd for array. Expecting 0x%x. Got: 0x%x", bincVdArray, d.vd)
	***REMOVED***
	length = d.decLen()
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) decLen() int ***REMOVED***
	if d.vs <= 3 ***REMOVED***
		return int(d.decUint())
	***REMOVED***
	return int(d.vs - 4)
***REMOVED***

func (d *bincDecDriver) decodeString() (s string) ***REMOVED***
	switch d.vd ***REMOVED***
	case bincVdString, bincVdByteArray:
		if length := d.decLen(); length > 0 ***REMOVED***
			s = string(d.r.readn(length))
		***REMOVED***
	case bincVdSymbol:
		//from vs: extract numSymbolBytes, containsStringVal, strLenPrecision,
		//extract symbol
		//if containsStringVal, read it and put in map
		//else look in map for string value
		var symbol uint32
		vs := d.vs
		//fmt.Printf(">>>> d.vs: 0b%b, & 0x8: %v, & 0x4: %v\n", d.vs, vs & 0x8, vs & 0x4)
		if vs&0x8 == 0 ***REMOVED***
			symbol = uint32(d.r.readn1())
		***REMOVED*** else ***REMOVED***
			symbol = uint32(d.r.readUint16())
		***REMOVED***
		if d.m == nil ***REMOVED***
			d.m = make(map[uint32]string, 16)
		***REMOVED***

		if vs&0x4 == 0 ***REMOVED***
			s = d.m[symbol]
		***REMOVED*** else ***REMOVED***
			var slen int
			switch vs & 0x3 ***REMOVED***
			case 0:
				slen = int(d.r.readn1())
			case 1:
				slen = int(d.r.readUint16())
			case 2:
				slen = int(d.r.readUint32())
			case 3:
				slen = int(d.r.readUint64())
			***REMOVED***
			s = string(d.r.readn(slen))
			d.m[symbol] = s
		***REMOVED***
	default:
		decErr("Invalid d.vd for string. Expecting string:0x%x, bytearray:0x%x or symbol: 0x%x. Got: 0x%x",
			bincVdString, bincVdByteArray, bincVdSymbol, d.vd)
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) decodeBytes(bs []byte) (bsOut []byte, changed bool) ***REMOVED***
	var clen int
	switch d.vd ***REMOVED***
	case bincVdString, bincVdByteArray:
		clen = d.decLen()
	default:
		decErr("Invalid d.vd for bytes. Expecting string:0x%x or bytearray:0x%x. Got: 0x%x",
			bincVdString, bincVdByteArray, d.vd)
	***REMOVED***
	if clen > 0 ***REMOVED***
		// if no contents in stream, don't update the passed byteslice
		if len(bs) != clen ***REMOVED***
			if len(bs) > clen ***REMOVED***
				bs = bs[:clen]
			***REMOVED*** else ***REMOVED***
				bs = make([]byte, clen)
			***REMOVED***
			bsOut = bs
			changed = true
		***REMOVED***
		d.r.readb(bs)
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte) ***REMOVED***
	switch d.vd ***REMOVED***
	case bincVdCustomExt:
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag ***REMOVED***
			decErr("Wrong extension tag. Got %b. Expecting: %v", xtag, tag)
		***REMOVED***
		xbs = d.r.readn(l)
	case bincVdByteArray:
		xbs, _ = d.decodeBytes(nil)
	default:
		decErr("Invalid d.vd for extensions (Expecting extensions or byte array). Got: 0x%x", d.vd)
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) decodeNaked() (v interface***REMOVED******REMOVED***, vt valueType, decodeFurther bool) ***REMOVED***
	d.initReadNext()

	switch d.vd ***REMOVED***
	case bincVdSpecial:
		switch d.vs ***REMOVED***
		case bincSpNil:
			vt = valueTypeNil
		case bincSpFalse:
			vt = valueTypeBool
			v = false
		case bincSpTrue:
			vt = valueTypeBool
			v = true
		case bincSpNan:
			vt = valueTypeFloat
			v = math.NaN()
		case bincSpPosInf:
			vt = valueTypeFloat
			v = math.Inf(1)
		case bincSpNegInf:
			vt = valueTypeFloat
			v = math.Inf(-1)
		case bincSpZeroFloat:
			vt = valueTypeFloat
			v = float64(0)
		case bincSpZero:
			vt = valueTypeUint
			v = int64(0) // int8(0)
		case bincSpNegOne:
			vt = valueTypeInt
			v = int64(-1) // int8(-1)
		default:
			decErr("decodeNaked: Unrecognized special value 0x%x", d.vs)
		***REMOVED***
	case bincVdSmallInt:
		vt = valueTypeUint
		v = uint64(int8(d.vs)) + 1 // int8(d.vs) + 1
	case bincVdPosInt:
		vt = valueTypeUint
		v = d.decUint()
	case bincVdNegInt:
		vt = valueTypeInt
		v = -(int64(d.decUint()))
	case bincVdFloat:
		vt = valueTypeFloat
		v = d.decFloat()
	case bincVdSymbol:
		vt = valueTypeSymbol
		v = d.decodeString()
	case bincVdString:
		vt = valueTypeString
		v = d.decodeString()
	case bincVdByteArray:
		vt = valueTypeBytes
		v, _ = d.decodeBytes(nil)
	case bincVdTimestamp:
		vt = valueTypeTimestamp
		tt, err := decodeTime(d.r.readn(int(d.vs)))
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		v = tt
	case bincVdCustomExt:
		vt = valueTypeExt
		l := d.decLen()
		var re RawExt
		re.Tag = d.r.readn1()
		re.Data = d.r.readn(l)
		v = &re
		vt = valueTypeExt
	case bincVdArray:
		vt = valueTypeArray
		decodeFurther = true
	case bincVdMap:
		vt = valueTypeMap
		decodeFurther = true
	default:
		decErr("decodeNaked: Unrecognized d.vd: 0x%x", d.vd)
	***REMOVED***

	if !decodeFurther ***REMOVED***
		d.bdRead = false
	***REMOVED***
	return
***REMOVED***

//------------------------------------

//BincHandle is a Handle for the Binc Schema-Free Encoding Format
//defined at https://github.com/ugorji/binc .
//
//BincHandle currently supports all Binc features with the following EXCEPTIONS:
//  - only integers up to 64 bits of precision are supported.
//    big integers are unsupported.
//  - Only IEEE 754 binary32 and binary64 floats are supported (ie Go float32 and float64 types).
//    extended precision and decimal IEEE 754 floats are unsupported.
//  - Only UTF-8 strings supported.
//    Unicode_Other Binc types (UTF16, UTF32) are currently unsupported.
//Note that these EXCEPTIONS are temporary and full support is possible and may happen soon.
type BincHandle struct ***REMOVED***
	BasicHandle
***REMOVED***

func (h *BincHandle) newEncDriver(w encWriter) encDriver ***REMOVED***
	return &bincEncDriver***REMOVED***w: w***REMOVED***
***REMOVED***

func (h *BincHandle) newDecDriver(r decReader) decDriver ***REMOVED***
	return &bincDecDriver***REMOVED***r: r***REMOVED***
***REMOVED***

func (_ *BincHandle) writeExt() bool ***REMOVED***
	return true
***REMOVED***

func (h *BincHandle) getBasicHandle() *BasicHandle ***REMOVED***
	return &h.BasicHandle
***REMOVED***
