// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import "math"

const (
	_               uint8 = iota
	simpleVdNil           = 1
	simpleVdFalse         = 2
	simpleVdTrue          = 3
	simpleVdFloat32       = 4
	simpleVdFloat64       = 5

	// each lasts for 4 (ie n, n+1, n+2, n+3)
	simpleVdPosInt = 8
	simpleVdNegInt = 12

	// containers: each lasts for 4 (ie n, n+1, n+2, ... n+7)
	simpleVdString    = 216
	simpleVdByteArray = 224
	simpleVdArray     = 232
	simpleVdMap       = 240
	simpleVdExt       = 248
)

type simpleEncDriver struct ***REMOVED***
	h *SimpleHandle
	w encWriter
	//b [8]byte
***REMOVED***

func (e *simpleEncDriver) isBuiltinType(rt uintptr) bool ***REMOVED***
	return false
***REMOVED***

func (e *simpleEncDriver) encodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED***
***REMOVED***

func (e *simpleEncDriver) encodeNil() ***REMOVED***
	e.w.writen1(simpleVdNil)
***REMOVED***

func (e *simpleEncDriver) encodeBool(b bool) ***REMOVED***
	if b ***REMOVED***
		e.w.writen1(simpleVdTrue)
	***REMOVED*** else ***REMOVED***
		e.w.writen1(simpleVdFalse)
	***REMOVED***
***REMOVED***

func (e *simpleEncDriver) encodeFloat32(f float32) ***REMOVED***
	e.w.writen1(simpleVdFloat32)
	e.w.writeUint32(math.Float32bits(f))
***REMOVED***

func (e *simpleEncDriver) encodeFloat64(f float64) ***REMOVED***
	e.w.writen1(simpleVdFloat64)
	e.w.writeUint64(math.Float64bits(f))
***REMOVED***

func (e *simpleEncDriver) encodeInt(v int64) ***REMOVED***
	if v < 0 ***REMOVED***
		e.encUint(uint64(-v), simpleVdNegInt)
	***REMOVED*** else ***REMOVED***
		e.encUint(uint64(v), simpleVdPosInt)
	***REMOVED***
***REMOVED***

func (e *simpleEncDriver) encodeUint(v uint64) ***REMOVED***
	e.encUint(v, simpleVdPosInt)
***REMOVED***

func (e *simpleEncDriver) encUint(v uint64, bd uint8) ***REMOVED***
	switch ***REMOVED***
	case v <= math.MaxUint8:
		e.w.writen2(bd, uint8(v))
	case v <= math.MaxUint16:
		e.w.writen1(bd + 1)
		e.w.writeUint16(uint16(v))
	case v <= math.MaxUint32:
		e.w.writen1(bd + 2)
		e.w.writeUint32(uint32(v))
	case v <= math.MaxUint64:
		e.w.writen1(bd + 3)
		e.w.writeUint64(v)
	***REMOVED***
***REMOVED***

func (e *simpleEncDriver) encLen(bd byte, length int) ***REMOVED***
	switch ***REMOVED***
	case length == 0:
		e.w.writen1(bd)
	case length <= math.MaxUint8:
		e.w.writen1(bd + 1)
		e.w.writen1(uint8(length))
	case length <= math.MaxUint16:
		e.w.writen1(bd + 2)
		e.w.writeUint16(uint16(length))
	case int64(length) <= math.MaxUint32:
		e.w.writen1(bd + 3)
		e.w.writeUint32(uint32(length))
	default:
		e.w.writen1(bd + 4)
		e.w.writeUint64(uint64(length))
	***REMOVED***
***REMOVED***

func (e *simpleEncDriver) encodeExtPreamble(xtag byte, length int) ***REMOVED***
	e.encLen(simpleVdExt, length)
	e.w.writen1(xtag)
***REMOVED***

func (e *simpleEncDriver) encodeArrayPreamble(length int) ***REMOVED***
	e.encLen(simpleVdArray, length)
***REMOVED***

func (e *simpleEncDriver) encodeMapPreamble(length int) ***REMOVED***
	e.encLen(simpleVdMap, length)
***REMOVED***

func (e *simpleEncDriver) encodeString(c charEncoding, v string) ***REMOVED***
	e.encLen(simpleVdString, len(v))
	e.w.writestr(v)
***REMOVED***

func (e *simpleEncDriver) encodeSymbol(v string) ***REMOVED***
	e.encodeString(c_UTF8, v)
***REMOVED***

func (e *simpleEncDriver) encodeStringBytes(c charEncoding, v []byte) ***REMOVED***
	e.encLen(simpleVdByteArray, len(v))
	e.w.writeb(v)
***REMOVED***

//------------------------------------

type simpleDecDriver struct ***REMOVED***
	h      *SimpleHandle
	r      decReader
	bdRead bool
	bdType valueType
	bd     byte
	//b      [8]byte
***REMOVED***

func (d *simpleDecDriver) initReadNext() ***REMOVED***
	if d.bdRead ***REMOVED***
		return
	***REMOVED***
	d.bd = d.r.readn1()
	d.bdRead = true
	d.bdType = valueTypeUnset
***REMOVED***

func (d *simpleDecDriver) currentEncodedType() valueType ***REMOVED***
	if d.bdType == valueTypeUnset ***REMOVED***
		switch d.bd ***REMOVED***
		case simpleVdNil:
			d.bdType = valueTypeNil
		case simpleVdTrue, simpleVdFalse:
			d.bdType = valueTypeBool
		case simpleVdPosInt, simpleVdPosInt + 1, simpleVdPosInt + 2, simpleVdPosInt + 3:
			d.bdType = valueTypeUint
		case simpleVdNegInt, simpleVdNegInt + 1, simpleVdNegInt + 2, simpleVdNegInt + 3:
			d.bdType = valueTypeInt
		case simpleVdFloat32, simpleVdFloat64:
			d.bdType = valueTypeFloat
		case simpleVdString, simpleVdString + 1, simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
			d.bdType = valueTypeString
		case simpleVdByteArray, simpleVdByteArray + 1, simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
			d.bdType = valueTypeBytes
		case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
			d.bdType = valueTypeExt
		case simpleVdArray, simpleVdArray + 1, simpleVdArray + 2, simpleVdArray + 3, simpleVdArray + 4:
			d.bdType = valueTypeArray
		case simpleVdMap, simpleVdMap + 1, simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
			d.bdType = valueTypeMap
		default:
			decErr("currentEncodedType: Unrecognized d.vd: 0x%x", d.bd)
		***REMOVED***
	***REMOVED***
	return d.bdType
***REMOVED***

func (d *simpleDecDriver) tryDecodeAsNil() bool ***REMOVED***
	if d.bd == simpleVdNil ***REMOVED***
		d.bdRead = false
		return true
	***REMOVED***
	return false
***REMOVED***

func (d *simpleDecDriver) isBuiltinType(rt uintptr) bool ***REMOVED***
	return false
***REMOVED***

func (d *simpleDecDriver) decodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED***
***REMOVED***

func (d *simpleDecDriver) decIntAny() (ui uint64, i int64, neg bool) ***REMOVED***
	switch d.bd ***REMOVED***
	case simpleVdPosInt:
		ui = uint64(d.r.readn1())
		i = int64(ui)
	case simpleVdPosInt + 1:
		ui = uint64(d.r.readUint16())
		i = int64(ui)
	case simpleVdPosInt + 2:
		ui = uint64(d.r.readUint32())
		i = int64(ui)
	case simpleVdPosInt + 3:
		ui = uint64(d.r.readUint64())
		i = int64(ui)
	case simpleVdNegInt:
		ui = uint64(d.r.readn1())
		i = -(int64(ui))
		neg = true
	case simpleVdNegInt + 1:
		ui = uint64(d.r.readUint16())
		i = -(int64(ui))
		neg = true
	case simpleVdNegInt + 2:
		ui = uint64(d.r.readUint32())
		i = -(int64(ui))
		neg = true
	case simpleVdNegInt + 3:
		ui = uint64(d.r.readUint64())
		i = -(int64(ui))
		neg = true
	default:
		decErr("decIntAny: Integer only valid from pos/neg integer1..8. Invalid descriptor: %v", d.bd)
	***REMOVED***
	// don't do this check, because callers may only want the unsigned value.
	// if ui > math.MaxInt64 ***REMOVED***
	// 	decErr("decIntAny: Integer out of range for signed int64: %v", ui)
	// ***REMOVED***
	return
***REMOVED***

func (d *simpleDecDriver) decodeInt(bitsize uint8) (i int64) ***REMOVED***
	_, i, _ = d.decIntAny()
	checkOverflow(0, i, bitsize)
	d.bdRead = false
	return
***REMOVED***

func (d *simpleDecDriver) decodeUint(bitsize uint8) (ui uint64) ***REMOVED***
	ui, i, neg := d.decIntAny()
	if neg ***REMOVED***
		decErr("Assigning negative signed value: %v, to unsigned type", i)
	***REMOVED***
	checkOverflow(ui, 0, bitsize)
	d.bdRead = false
	return
***REMOVED***

func (d *simpleDecDriver) decodeFloat(chkOverflow32 bool) (f float64) ***REMOVED***
	switch d.bd ***REMOVED***
	case simpleVdFloat32:
		f = float64(math.Float32frombits(d.r.readUint32()))
	case simpleVdFloat64:
		f = math.Float64frombits(d.r.readUint64())
	default:
		if d.bd >= simpleVdPosInt && d.bd <= simpleVdNegInt+3 ***REMOVED***
			_, i, _ := d.decIntAny()
			f = float64(i)
		***REMOVED*** else ***REMOVED***
			decErr("Float only valid from float32/64: Invalid descriptor: %v", d.bd)
		***REMOVED***
	***REMOVED***
	checkOverflowFloat32(f, chkOverflow32)
	d.bdRead = false
	return
***REMOVED***

// bool can be decoded from bool only (single byte).
func (d *simpleDecDriver) decodeBool() (b bool) ***REMOVED***
	switch d.bd ***REMOVED***
	case simpleVdTrue:
		b = true
	case simpleVdFalse:
	default:
		decErr("Invalid single-byte value for bool: %s: %x", msgBadDesc, d.bd)
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *simpleDecDriver) readMapLen() (length int) ***REMOVED***
	d.bdRead = false
	return d.decLen()
***REMOVED***

func (d *simpleDecDriver) readArrayLen() (length int) ***REMOVED***
	d.bdRead = false
	return d.decLen()
***REMOVED***

func (d *simpleDecDriver) decLen() int ***REMOVED***
	switch d.bd % 8 ***REMOVED***
	case 0:
		return 0
	case 1:
		return int(d.r.readn1())
	case 2:
		return int(d.r.readUint16())
	case 3:
		ui := uint64(d.r.readUint32())
		checkOverflow(ui, 0, intBitsize)
		return int(ui)
	case 4:
		ui := d.r.readUint64()
		checkOverflow(ui, 0, intBitsize)
		return int(ui)
	***REMOVED***
	decErr("decLen: Cannot read length: bd%8 must be in range 0..4. Got: %d", d.bd%8)
	return -1
***REMOVED***

func (d *simpleDecDriver) decodeString() (s string) ***REMOVED***
	s = string(d.r.readn(d.decLen()))
	d.bdRead = false
	return
***REMOVED***

func (d *simpleDecDriver) decodeBytes(bs []byte) (bsOut []byte, changed bool) ***REMOVED***
	if clen := d.decLen(); clen > 0 ***REMOVED***
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

func (d *simpleDecDriver) decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte) ***REMOVED***
	switch d.bd ***REMOVED***
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag ***REMOVED***
			decErr("Wrong extension tag. Got %b. Expecting: %v", xtag, tag)
		***REMOVED***
		xbs = d.r.readn(l)
	case simpleVdByteArray, simpleVdByteArray + 1, simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		xbs, _ = d.decodeBytes(nil)
	default:
		decErr("Invalid d.vd for extensions (Expecting extensions or byte array). Got: 0x%x", d.bd)
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *simpleDecDriver) decodeNaked() (v interface***REMOVED******REMOVED***, vt valueType, decodeFurther bool) ***REMOVED***
	d.initReadNext()

	switch d.bd ***REMOVED***
	case simpleVdNil:
		vt = valueTypeNil
	case simpleVdFalse:
		vt = valueTypeBool
		v = false
	case simpleVdTrue:
		vt = valueTypeBool
		v = true
	case simpleVdPosInt, simpleVdPosInt + 1, simpleVdPosInt + 2, simpleVdPosInt + 3:
		vt = valueTypeUint
		ui, _, _ := d.decIntAny()
		v = ui
	case simpleVdNegInt, simpleVdNegInt + 1, simpleVdNegInt + 2, simpleVdNegInt + 3:
		vt = valueTypeInt
		_, i, _ := d.decIntAny()
		v = i
	case simpleVdFloat32:
		vt = valueTypeFloat
		v = d.decodeFloat(true)
	case simpleVdFloat64:
		vt = valueTypeFloat
		v = d.decodeFloat(false)
	case simpleVdString, simpleVdString + 1, simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		vt = valueTypeString
		v = d.decodeString()
	case simpleVdByteArray, simpleVdByteArray + 1, simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		vt = valueTypeBytes
		v, _ = d.decodeBytes(nil)
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		vt = valueTypeExt
		l := d.decLen()
		var re RawExt
		re.Tag = d.r.readn1()
		re.Data = d.r.readn(l)
		v = &re
		vt = valueTypeExt
	case simpleVdArray, simpleVdArray + 1, simpleVdArray + 2, simpleVdArray + 3, simpleVdArray + 4:
		vt = valueTypeArray
		decodeFurther = true
	case simpleVdMap, simpleVdMap + 1, simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		vt = valueTypeMap
		decodeFurther = true
	default:
		decErr("decodeNaked: Unrecognized d.vd: 0x%x", d.bd)
	***REMOVED***

	if !decodeFurther ***REMOVED***
		d.bdRead = false
	***REMOVED***
	return
***REMOVED***

//------------------------------------

// SimpleHandle is a Handle for a very simple encoding format.
//
// simple is a simplistic codec similar to binc, but not as compact.
//   - Encoding of a value is always preceeded by the descriptor byte (bd)
//   - True, false, nil are encoded fully in 1 byte (the descriptor)
//   - Integers (intXXX, uintXXX) are encoded in 1, 2, 4 or 8 bytes (plus a descriptor byte).
//     There are positive (uintXXX and intXXX >= 0) and negative (intXXX < 0) integers.
//   - Floats are encoded in 4 or 8 bytes (plus a descriptor byte)
//   - Lenght of containers (strings, bytes, array, map, extensions)
//     are encoded in 0, 1, 2, 4 or 8 bytes.
//     Zero-length containers have no length encoded.
//     For others, the number of bytes is given by pow(2, bd%3)
//   - maps are encoded as [bd] [length] [[key][value]]...
//   - arrays are encoded as [bd] [length] [value]...
//   - extensions are encoded as [bd] [length] [tag] [byte]...
//   - strings/bytearrays are encoded as [bd] [length] [byte]...
//
// The full spec will be published soon.
type SimpleHandle struct ***REMOVED***
	BasicHandle
***REMOVED***

func (h *SimpleHandle) newEncDriver(w encWriter) encDriver ***REMOVED***
	return &simpleEncDriver***REMOVED***w: w, h: h***REMOVED***
***REMOVED***

func (h *SimpleHandle) newDecDriver(r decReader) decDriver ***REMOVED***
	return &simpleDecDriver***REMOVED***r: r, h: h***REMOVED***
***REMOVED***

func (_ *SimpleHandle) writeExt() bool ***REMOVED***
	return true
***REMOVED***

func (h *SimpleHandle) getBasicHandle() *BasicHandle ***REMOVED***
	return &h.BasicHandle
***REMOVED***

var _ decDriver = (*simpleDecDriver)(nil)
var _ encDriver = (*simpleEncDriver)(nil)
