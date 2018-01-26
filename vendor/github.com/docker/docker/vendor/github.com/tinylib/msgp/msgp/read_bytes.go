package msgp

import (
	"bytes"
	"encoding/binary"
	"math"
	"time"
)

var big = binary.BigEndian

// NextType returns the type of the next
// object in the slice. If the length
// of the input is zero, it returns
// InvalidType.
func NextType(b []byte) Type ***REMOVED***
	if len(b) == 0 ***REMOVED***
		return InvalidType
	***REMOVED***
	spec := sizes[b[0]]
	t := spec.typ
	if t == ExtensionType && len(b) > int(spec.size) ***REMOVED***
		var tp int8
		if spec.extra == constsize ***REMOVED***
			tp = int8(b[1])
		***REMOVED*** else ***REMOVED***
			tp = int8(b[spec.size-1])
		***REMOVED***
		switch tp ***REMOVED***
		case TimeExtension:
			return TimeType
		case Complex128Extension:
			return Complex128Type
		case Complex64Extension:
			return Complex64Type
		default:
			return ExtensionType
		***REMOVED***
	***REMOVED***
	return t
***REMOVED***

// IsNil returns true if len(b)>0 and
// the leading byte is a 'nil' MessagePack
// byte; false otherwise
func IsNil(b []byte) bool ***REMOVED***
	if len(b) != 0 && b[0] == mnil ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// Raw is raw MessagePack.
// Raw allows you to read and write
// data without interpreting its contents.
type Raw []byte

// MarshalMsg implements msgp.Marshaler.
// It appends the raw contents of 'raw'
// to the provided byte slice. If 'raw'
// is 0 bytes, 'nil' will be appended instead.
func (r Raw) MarshalMsg(b []byte) ([]byte, error) ***REMOVED***
	i := len(r)
	if i == 0 ***REMOVED***
		return AppendNil(b), nil
	***REMOVED***
	o, l := ensure(b, i)
	copy(o[l:], []byte(r))
	return o, nil
***REMOVED***

// UnmarshalMsg implements msgp.Unmarshaler.
// It sets the contents of *Raw to be the next
// object in the provided byte slice.
func (r *Raw) UnmarshalMsg(b []byte) ([]byte, error) ***REMOVED***
	l := len(b)
	out, err := Skip(b)
	if err != nil ***REMOVED***
		return b, err
	***REMOVED***
	rlen := l - len(out)
	if cap(*r) < rlen ***REMOVED***
		*r = make(Raw, rlen)
	***REMOVED*** else ***REMOVED***
		*r = (*r)[0:rlen]
	***REMOVED***
	copy(*r, b[:rlen])
	return out, nil
***REMOVED***

// EncodeMsg implements msgp.Encodable.
// It writes the raw bytes to the writer.
// If r is empty, it writes 'nil' instead.
func (r Raw) EncodeMsg(w *Writer) error ***REMOVED***
	if len(r) == 0 ***REMOVED***
		return w.WriteNil()
	***REMOVED***
	_, err := w.Write([]byte(r))
	return err
***REMOVED***

// DecodeMsg implements msgp.Decodable.
// It sets the value of *Raw to be the
// next object on the wire.
func (r *Raw) DecodeMsg(f *Reader) error ***REMOVED***
	*r = (*r)[:0]
	return appendNext(f, (*[]byte)(r))
***REMOVED***

// Msgsize implements msgp.Sizer
func (r Raw) Msgsize() int ***REMOVED***
	l := len(r)
	if l == 0 ***REMOVED***
		return 1 // for 'nil'
	***REMOVED***
	return l
***REMOVED***

func appendNext(f *Reader, d *[]byte) error ***REMOVED***
	amt, o, err := getNextSize(f.R)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	var i int
	*d, i = ensure(*d, int(amt))
	_, err = f.R.ReadFull((*d)[i:])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for o > 0 ***REMOVED***
		err = appendNext(f, d)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		o--
	***REMOVED***
	return nil
***REMOVED***

// MarshalJSON implements json.Marshaler
func (r *Raw) MarshalJSON() ([]byte, error) ***REMOVED***
	var buf bytes.Buffer
	_, err := UnmarshalAsJSON(&buf, []byte(*r))
	return buf.Bytes(), err
***REMOVED***

// ReadMapHeaderBytes reads a map header size
// from 'b' and returns the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a map)
func ReadMapHeaderBytes(b []byte) (sz uint32, o []byte, err error) ***REMOVED***
	l := len(b)
	if l < 1 ***REMOVED***
		err = ErrShortBytes
		return
	***REMOVED***

	lead := b[0]
	if isfixmap(lead) ***REMOVED***
		sz = uint32(rfixmap(lead))
		o = b[1:]
		return
	***REMOVED***

	switch lead ***REMOVED***
	case mmap16:
		if l < 3 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		sz = uint32(big.Uint16(b[1:]))
		o = b[3:]
		return

	case mmap32:
		if l < 5 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		sz = big.Uint32(b[1:])
		o = b[5:]
		return

	default:
		err = badPrefix(MapType, lead)
		return
	***REMOVED***
***REMOVED***

// ReadMapKeyZC attempts to read a map key
// from 'b' and returns the key bytes and the remaining bytes
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a str or bin)
func ReadMapKeyZC(b []byte) ([]byte, []byte, error) ***REMOVED***
	o, b, err := ReadStringZC(b)
	if err != nil ***REMOVED***
		if tperr, ok := err.(TypeError); ok && tperr.Encoded == BinType ***REMOVED***
			return ReadBytesZC(b)
		***REMOVED***
		return nil, b, err
	***REMOVED***
	return o, b, nil
***REMOVED***

// ReadArrayHeaderBytes attempts to read
// the array header size off of 'b' and return
// the size and remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not an array)
func ReadArrayHeaderBytes(b []byte) (sz uint32, o []byte, err error) ***REMOVED***
	if len(b) < 1 ***REMOVED***
		return 0, nil, ErrShortBytes
	***REMOVED***
	lead := b[0]
	if isfixarray(lead) ***REMOVED***
		sz = uint32(rfixarray(lead))
		o = b[1:]
		return
	***REMOVED***

	switch lead ***REMOVED***
	case marray16:
		if len(b) < 3 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		sz = uint32(big.Uint16(b[1:]))
		o = b[3:]
		return

	case marray32:
		if len(b) < 5 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		sz = big.Uint32(b[1:])
		o = b[5:]
		return

	default:
		err = badPrefix(ArrayType, lead)
		return
	***REMOVED***
***REMOVED***

// ReadNilBytes tries to read a "nil" byte
// off of 'b' and return the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a 'nil')
// - InvalidPrefixError
func ReadNilBytes(b []byte) ([]byte, error) ***REMOVED***
	if len(b) < 1 ***REMOVED***
		return nil, ErrShortBytes
	***REMOVED***
	if b[0] != mnil ***REMOVED***
		return b, badPrefix(NilType, b[0])
	***REMOVED***
	return b[1:], nil
***REMOVED***

// ReadFloat64Bytes tries to read a float64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a float64)
func ReadFloat64Bytes(b []byte) (f float64, o []byte, err error) ***REMOVED***
	if len(b) < 9 ***REMOVED***
		if len(b) >= 5 && b[0] == mfloat32 ***REMOVED***
			var tf float32
			tf, o, err = ReadFloat32Bytes(b)
			f = float64(tf)
			return
		***REMOVED***
		err = ErrShortBytes
		return
	***REMOVED***

	if b[0] != mfloat64 ***REMOVED***
		if b[0] == mfloat32 ***REMOVED***
			var tf float32
			tf, o, err = ReadFloat32Bytes(b)
			f = float64(tf)
			return
		***REMOVED***
		err = badPrefix(Float64Type, b[0])
		return
	***REMOVED***

	f = math.Float64frombits(getMuint64(b))
	o = b[9:]
	return
***REMOVED***

// ReadFloat32Bytes tries to read a float64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a float32)
func ReadFloat32Bytes(b []byte) (f float32, o []byte, err error) ***REMOVED***
	if len(b) < 5 ***REMOVED***
		err = ErrShortBytes
		return
	***REMOVED***

	if b[0] != mfloat32 ***REMOVED***
		err = TypeError***REMOVED***Method: Float32Type, Encoded: getType(b[0])***REMOVED***
		return
	***REMOVED***

	f = math.Float32frombits(getMuint32(b))
	o = b[5:]
	return
***REMOVED***

// ReadBoolBytes tries to read a float64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a bool)
func ReadBoolBytes(b []byte) (bool, []byte, error) ***REMOVED***
	if len(b) < 1 ***REMOVED***
		return false, b, ErrShortBytes
	***REMOVED***
	switch b[0] ***REMOVED***
	case mtrue:
		return true, b[1:], nil
	case mfalse:
		return false, b[1:], nil
	default:
		return false, b, badPrefix(BoolType, b[0])
	***REMOVED***
***REMOVED***

// ReadInt64Bytes tries to read an int64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError (not a int)
func ReadInt64Bytes(b []byte) (i int64, o []byte, err error) ***REMOVED***
	l := len(b)
	if l < 1 ***REMOVED***
		return 0, nil, ErrShortBytes
	***REMOVED***

	lead := b[0]
	if isfixint(lead) ***REMOVED***
		i = int64(rfixint(lead))
		o = b[1:]
		return
	***REMOVED***
	if isnfixint(lead) ***REMOVED***
		i = int64(rnfixint(lead))
		o = b[1:]
		return
	***REMOVED***

	switch lead ***REMOVED***
	case mint8:
		if l < 2 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		i = int64(getMint8(b))
		o = b[2:]
		return

	case mint16:
		if l < 3 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		i = int64(getMint16(b))
		o = b[3:]
		return

	case mint32:
		if l < 5 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		i = int64(getMint32(b))
		o = b[5:]
		return

	case mint64:
		if l < 9 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		i = getMint64(b)
		o = b[9:]
		return

	default:
		err = badPrefix(IntType, lead)
		return
	***REMOVED***
***REMOVED***

// ReadInt32Bytes tries to read an int32
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a int)
// - IntOverflow***REMOVED******REMOVED*** (value doesn't fit in int32)
func ReadInt32Bytes(b []byte) (int32, []byte, error) ***REMOVED***
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt32 || i < math.MinInt32 ***REMOVED***
		return 0, o, IntOverflow***REMOVED***Value: i, FailedBitsize: 32***REMOVED***
	***REMOVED***
	return int32(i), o, err
***REMOVED***

// ReadInt16Bytes tries to read an int16
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a int)
// - IntOverflow***REMOVED******REMOVED*** (value doesn't fit in int16)
func ReadInt16Bytes(b []byte) (int16, []byte, error) ***REMOVED***
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt16 || i < math.MinInt16 ***REMOVED***
		return 0, o, IntOverflow***REMOVED***Value: i, FailedBitsize: 16***REMOVED***
	***REMOVED***
	return int16(i), o, err
***REMOVED***

// ReadInt8Bytes tries to read an int16
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a int)
// - IntOverflow***REMOVED******REMOVED*** (value doesn't fit in int8)
func ReadInt8Bytes(b []byte) (int8, []byte, error) ***REMOVED***
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt8 || i < math.MinInt8 ***REMOVED***
		return 0, o, IntOverflow***REMOVED***Value: i, FailedBitsize: 8***REMOVED***
	***REMOVED***
	return int8(i), o, err
***REMOVED***

// ReadIntBytes tries to read an int
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a int)
// - IntOverflow***REMOVED******REMOVED*** (value doesn't fit in int; 32-bit platforms only)
func ReadIntBytes(b []byte) (int, []byte, error) ***REMOVED***
	if smallint ***REMOVED***
		i, b, err := ReadInt32Bytes(b)
		return int(i), b, err
	***REMOVED***
	i, b, err := ReadInt64Bytes(b)
	return int(i), b, err
***REMOVED***

// ReadUint64Bytes tries to read a uint64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a uint)
func ReadUint64Bytes(b []byte) (u uint64, o []byte, err error) ***REMOVED***
	l := len(b)
	if l < 1 ***REMOVED***
		return 0, nil, ErrShortBytes
	***REMOVED***

	lead := b[0]
	if isfixint(lead) ***REMOVED***
		u = uint64(rfixint(lead))
		o = b[1:]
		return
	***REMOVED***

	switch lead ***REMOVED***
	case muint8:
		if l < 2 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		u = uint64(getMuint8(b))
		o = b[2:]
		return

	case muint16:
		if l < 3 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		u = uint64(getMuint16(b))
		o = b[3:]
		return

	case muint32:
		if l < 5 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		u = uint64(getMuint32(b))
		o = b[5:]
		return

	case muint64:
		if l < 9 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		u = getMuint64(b)
		o = b[9:]
		return

	default:
		err = badPrefix(UintType, lead)
		return
	***REMOVED***
***REMOVED***

// ReadUint32Bytes tries to read a uint32
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a uint)
// - UintOverflow***REMOVED******REMOVED*** (value too large for uint32)
func ReadUint32Bytes(b []byte) (uint32, []byte, error) ***REMOVED***
	v, o, err := ReadUint64Bytes(b)
	if v > math.MaxUint32 ***REMOVED***
		return 0, nil, UintOverflow***REMOVED***Value: v, FailedBitsize: 32***REMOVED***
	***REMOVED***
	return uint32(v), o, err
***REMOVED***

// ReadUint16Bytes tries to read a uint16
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a uint)
// - UintOverflow***REMOVED******REMOVED*** (value too large for uint16)
func ReadUint16Bytes(b []byte) (uint16, []byte, error) ***REMOVED***
	v, o, err := ReadUint64Bytes(b)
	if v > math.MaxUint16 ***REMOVED***
		return 0, nil, UintOverflow***REMOVED***Value: v, FailedBitsize: 16***REMOVED***
	***REMOVED***
	return uint16(v), o, err
***REMOVED***

// ReadUint8Bytes tries to read a uint8
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a uint)
// - UintOverflow***REMOVED******REMOVED*** (value too large for uint8)
func ReadUint8Bytes(b []byte) (uint8, []byte, error) ***REMOVED***
	v, o, err := ReadUint64Bytes(b)
	if v > math.MaxUint8 ***REMOVED***
		return 0, nil, UintOverflow***REMOVED***Value: v, FailedBitsize: 8***REMOVED***
	***REMOVED***
	return uint8(v), o, err
***REMOVED***

// ReadUintBytes tries to read a uint
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a uint)
// - UintOverflow***REMOVED******REMOVED*** (value too large for uint; 32-bit platforms only)
func ReadUintBytes(b []byte) (uint, []byte, error) ***REMOVED***
	if smallint ***REMOVED***
		u, b, err := ReadUint32Bytes(b)
		return uint(u), b, err
	***REMOVED***
	u, b, err := ReadUint64Bytes(b)
	return uint(u), b, err
***REMOVED***

// ReadByteBytes is analogous to ReadUint8Bytes
func ReadByteBytes(b []byte) (byte, []byte, error) ***REMOVED***
	return ReadUint8Bytes(b)
***REMOVED***

// ReadBytesBytes reads a 'bin' object
// from 'b' and returns its vaue and
// the remaining bytes in 'b'.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError***REMOVED******REMOVED*** (not a 'bin' object)
func ReadBytesBytes(b []byte, scratch []byte) (v []byte, o []byte, err error) ***REMOVED***
	return readBytesBytes(b, scratch, false)
***REMOVED***

func readBytesBytes(b []byte, scratch []byte, zc bool) (v []byte, o []byte, err error) ***REMOVED***
	l := len(b)
	if l < 1 ***REMOVED***
		return nil, nil, ErrShortBytes
	***REMOVED***

	lead := b[0]
	var read int
	switch lead ***REMOVED***
	case mbin8:
		if l < 2 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***

		read = int(b[1])
		b = b[2:]

	case mbin16:
		if l < 3 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		read = int(big.Uint16(b[1:]))
		b = b[3:]

	case mbin32:
		if l < 5 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		read = int(big.Uint32(b[1:]))
		b = b[5:]

	default:
		err = badPrefix(BinType, lead)
		return
	***REMOVED***

	if len(b) < read ***REMOVED***
		err = ErrShortBytes
		return
	***REMOVED***

	// zero-copy
	if zc ***REMOVED***
		v = b[0:read]
		o = b[read:]
		return
	***REMOVED***

	if cap(scratch) >= read ***REMOVED***
		v = scratch[0:read]
	***REMOVED*** else ***REMOVED***
		v = make([]byte, read)
	***REMOVED***

	o = b[copy(v, b):]
	return
***REMOVED***

// ReadBytesZC extracts the messagepack-encoded
// binary field without copying. The returned []byte
// points to the same memory as the input slice.
// Possible errors:
// - ErrShortBytes (b not long enough)
// - TypeError***REMOVED******REMOVED*** (object not 'bin')
func ReadBytesZC(b []byte) (v []byte, o []byte, err error) ***REMOVED***
	return readBytesBytes(b, nil, true)
***REMOVED***

func ReadExactBytes(b []byte, into []byte) (o []byte, err error) ***REMOVED***
	l := len(b)
	if l < 1 ***REMOVED***
		err = ErrShortBytes
		return
	***REMOVED***

	lead := b[0]
	var read uint32
	var skip int
	switch lead ***REMOVED***
	case mbin8:
		if l < 2 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***

		read = uint32(b[1])
		skip = 2

	case mbin16:
		if l < 3 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		read = uint32(big.Uint16(b[1:]))
		skip = 3

	case mbin32:
		if l < 5 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		read = uint32(big.Uint32(b[1:]))
		skip = 5

	default:
		err = badPrefix(BinType, lead)
		return
	***REMOVED***

	if read != uint32(len(into)) ***REMOVED***
		err = ArrayError***REMOVED***Wanted: uint32(len(into)), Got: read***REMOVED***
		return
	***REMOVED***

	o = b[skip+copy(into, b[skip:]):]
	return
***REMOVED***

// ReadStringZC reads a messagepack string field
// without copying. The returned []byte points
// to the same memory as the input slice.
// Possible errors:
// - ErrShortBytes (b not long enough)
// - TypeError***REMOVED******REMOVED*** (object not 'str')
func ReadStringZC(b []byte) (v []byte, o []byte, err error) ***REMOVED***
	l := len(b)
	if l < 1 ***REMOVED***
		return nil, nil, ErrShortBytes
	***REMOVED***

	lead := b[0]
	var read int

	if isfixstr(lead) ***REMOVED***
		read = int(rfixstr(lead))
		b = b[1:]
	***REMOVED*** else ***REMOVED***
		switch lead ***REMOVED***
		case mstr8:
			if l < 2 ***REMOVED***
				err = ErrShortBytes
				return
			***REMOVED***
			read = int(b[1])
			b = b[2:]

		case mstr16:
			if l < 3 ***REMOVED***
				err = ErrShortBytes
				return
			***REMOVED***
			read = int(big.Uint16(b[1:]))
			b = b[3:]

		case mstr32:
			if l < 5 ***REMOVED***
				err = ErrShortBytes
				return
			***REMOVED***
			read = int(big.Uint32(b[1:]))
			b = b[5:]

		default:
			err = TypeError***REMOVED***Method: StrType, Encoded: getType(lead)***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	if len(b) < read ***REMOVED***
		err = ErrShortBytes
		return
	***REMOVED***

	v = b[0:read]
	o = b[read:]
	return
***REMOVED***

// ReadStringBytes reads a 'str' object
// from 'b' and returns its value and the
// remaining bytes in 'b'.
// Possible errors:
// - ErrShortBytes (b not long enough)
// - TypeError***REMOVED******REMOVED*** (not 'str' type)
// - InvalidPrefixError
func ReadStringBytes(b []byte) (string, []byte, error) ***REMOVED***
	v, o, err := ReadStringZC(b)
	return string(v), o, err
***REMOVED***

// ReadStringAsBytes reads a 'str' object
// into a slice of bytes. 'v' is the value of
// the 'str' object, which may reside in memory
// pointed to by 'scratch.' 'o' is the remaining bytes
// in 'b.''
// Possible errors:
// - ErrShortBytes (b not long enough)
// - TypeError***REMOVED******REMOVED*** (not 'str' type)
// - InvalidPrefixError (unknown type marker)
func ReadStringAsBytes(b []byte, scratch []byte) (v []byte, o []byte, err error) ***REMOVED***
	var tmp []byte
	tmp, o, err = ReadStringZC(b)
	v = append(scratch[:0], tmp...)
	return
***REMOVED***

// ReadComplex128Bytes reads a complex128
// extension object from 'b' and returns the
// remaining bytes.
// Possible errors:
// - ErrShortBytes (not enough bytes in 'b')
// - TypeError***REMOVED******REMOVED*** (object not a complex128)
// - InvalidPrefixError
// - ExtensionTypeError***REMOVED******REMOVED*** (object an extension of the correct size, but not a complex128)
func ReadComplex128Bytes(b []byte) (c complex128, o []byte, err error) ***REMOVED***
	if len(b) < 18 ***REMOVED***
		err = ErrShortBytes
		return
	***REMOVED***
	if b[0] != mfixext16 ***REMOVED***
		err = badPrefix(Complex128Type, b[0])
		return
	***REMOVED***
	if int8(b[1]) != Complex128Extension ***REMOVED***
		err = errExt(int8(b[1]), Complex128Extension)
		return
	***REMOVED***
	c = complex(math.Float64frombits(big.Uint64(b[2:])),
		math.Float64frombits(big.Uint64(b[10:])))
	o = b[18:]
	return
***REMOVED***

// ReadComplex64Bytes reads a complex64
// extension object from 'b' and returns the
// remaining bytes.
// Possible errors:
// - ErrShortBytes (not enough bytes in 'b')
// - TypeError***REMOVED******REMOVED*** (object not a complex64)
// - ExtensionTypeError***REMOVED******REMOVED*** (object an extension of the correct size, but not a complex64)
func ReadComplex64Bytes(b []byte) (c complex64, o []byte, err error) ***REMOVED***
	if len(b) < 10 ***REMOVED***
		err = ErrShortBytes
		return
	***REMOVED***
	if b[0] != mfixext8 ***REMOVED***
		err = badPrefix(Complex64Type, b[0])
		return
	***REMOVED***
	if b[1] != Complex64Extension ***REMOVED***
		err = errExt(int8(b[1]), Complex64Extension)
		return
	***REMOVED***
	c = complex(math.Float32frombits(big.Uint32(b[2:])),
		math.Float32frombits(big.Uint32(b[6:])))
	o = b[10:]
	return
***REMOVED***

// ReadTimeBytes reads a time.Time
// extension object from 'b' and returns the
// remaining bytes.
// Possible errors:
// - ErrShortBytes (not enough bytes in 'b')
// - TypeError***REMOVED******REMOVED*** (object not a complex64)
// - ExtensionTypeError***REMOVED******REMOVED*** (object an extension of the correct size, but not a time.Time)
func ReadTimeBytes(b []byte) (t time.Time, o []byte, err error) ***REMOVED***
	if len(b) < 15 ***REMOVED***
		err = ErrShortBytes
		return
	***REMOVED***
	if b[0] != mext8 || b[1] != 12 ***REMOVED***
		err = badPrefix(TimeType, b[0])
		return
	***REMOVED***
	if int8(b[2]) != TimeExtension ***REMOVED***
		err = errExt(int8(b[2]), TimeExtension)
		return
	***REMOVED***
	sec, nsec := getUnix(b[3:])
	t = time.Unix(sec, int64(nsec)).Local()
	o = b[15:]
	return
***REMOVED***

// ReadMapStrIntfBytes reads a map[string]interface***REMOVED******REMOVED***
// out of 'b' and returns the map and remaining bytes.
// If 'old' is non-nil, the values will be read into that map.
func ReadMapStrIntfBytes(b []byte, old map[string]interface***REMOVED******REMOVED***) (v map[string]interface***REMOVED******REMOVED***, o []byte, err error) ***REMOVED***
	var sz uint32
	o = b
	sz, o, err = ReadMapHeaderBytes(o)

	if err != nil ***REMOVED***
		return
	***REMOVED***

	if old != nil ***REMOVED***
		for key := range old ***REMOVED***
			delete(old, key)
		***REMOVED***
		v = old
	***REMOVED*** else ***REMOVED***
		v = make(map[string]interface***REMOVED******REMOVED***, int(sz))
	***REMOVED***

	for z := uint32(0); z < sz; z++ ***REMOVED***
		if len(o) < 1 ***REMOVED***
			err = ErrShortBytes
			return
		***REMOVED***
		var key []byte
		key, o, err = ReadMapKeyZC(o)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		var val interface***REMOVED******REMOVED***
		val, o, err = ReadIntfBytes(o)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		v[string(key)] = val
	***REMOVED***
	return
***REMOVED***

// ReadIntfBytes attempts to read
// the next object out of 'b' as a raw interface***REMOVED******REMOVED*** and
// return the remaining bytes.
func ReadIntfBytes(b []byte) (i interface***REMOVED******REMOVED***, o []byte, err error) ***REMOVED***
	if len(b) < 1 ***REMOVED***
		err = ErrShortBytes
		return
	***REMOVED***

	k := NextType(b)

	switch k ***REMOVED***
	case MapType:
		i, o, err = ReadMapStrIntfBytes(b, nil)
		return

	case ArrayType:
		var sz uint32
		sz, o, err = ReadArrayHeaderBytes(b)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		j := make([]interface***REMOVED******REMOVED***, int(sz))
		i = j
		for d := range j ***REMOVED***
			j[d], o, err = ReadIntfBytes(o)
			if err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		return

	case Float32Type:
		i, o, err = ReadFloat32Bytes(b)
		return

	case Float64Type:
		i, o, err = ReadFloat64Bytes(b)
		return

	case IntType:
		i, o, err = ReadInt64Bytes(b)
		return

	case UintType:
		i, o, err = ReadUint64Bytes(b)
		return

	case BoolType:
		i, o, err = ReadBoolBytes(b)
		return

	case TimeType:
		i, o, err = ReadTimeBytes(b)
		return

	case Complex64Type:
		i, o, err = ReadComplex64Bytes(b)
		return

	case Complex128Type:
		i, o, err = ReadComplex128Bytes(b)
		return

	case ExtensionType:
		var t int8
		t, err = peekExtension(b)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		// use a user-defined extension,
		// if it's been registered
		f, ok := extensionReg[t]
		if ok ***REMOVED***
			e := f()
			o, err = ReadExtensionBytes(b, e)
			i = e
			return
		***REMOVED***
		// last resort is a raw extension
		e := RawExtension***REMOVED******REMOVED***
		e.Type = int8(t)
		o, err = ReadExtensionBytes(b, &e)
		i = &e
		return

	case NilType:
		o, err = ReadNilBytes(b)
		return

	case BinType:
		i, o, err = ReadBytesBytes(b, nil)
		return

	case StrType:
		i, o, err = ReadStringBytes(b)
		return

	default:
		err = InvalidPrefixError(b[0])
		return
	***REMOVED***
***REMOVED***

// Skip skips the next object in 'b' and
// returns the remaining bytes. If the object
// is a map or array, all of its elements
// will be skipped.
// Possible Errors:
// - ErrShortBytes (not enough bytes in b)
// - InvalidPrefixError (bad encoding)
func Skip(b []byte) ([]byte, error) ***REMOVED***
	sz, asz, err := getSize(b)
	if err != nil ***REMOVED***
		return b, err
	***REMOVED***
	if uintptr(len(b)) < sz ***REMOVED***
		return b, ErrShortBytes
	***REMOVED***
	b = b[sz:]
	for asz > 0 ***REMOVED***
		b, err = Skip(b)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
		asz--
	***REMOVED***
	return b, nil
***REMOVED***

// returns (skip N bytes, skip M objects, error)
func getSize(b []byte) (uintptr, uintptr, error) ***REMOVED***
	l := len(b)
	if l == 0 ***REMOVED***
		return 0, 0, ErrShortBytes
	***REMOVED***
	lead := b[0]
	spec := &sizes[lead] // get type information
	size, mode := spec.size, spec.extra
	if size == 0 ***REMOVED***
		return 0, 0, InvalidPrefixError(lead)
	***REMOVED***
	if mode >= 0 ***REMOVED*** // fixed composites
		return uintptr(size), uintptr(mode), nil
	***REMOVED***
	if l < int(size) ***REMOVED***
		return 0, 0, ErrShortBytes
	***REMOVED***
	switch mode ***REMOVED***
	case extra8:
		return uintptr(size) + uintptr(b[1]), 0, nil
	case extra16:
		return uintptr(size) + uintptr(big.Uint16(b[1:])), 0, nil
	case extra32:
		return uintptr(size) + uintptr(big.Uint32(b[1:])), 0, nil
	case map16v:
		return uintptr(size), 2 * uintptr(big.Uint16(b[1:])), nil
	case map32v:
		return uintptr(size), 2 * uintptr(big.Uint32(b[1:])), nil
	case array16v:
		return uintptr(size), uintptr(big.Uint16(b[1:])), nil
	case array32v:
		return uintptr(size), uintptr(big.Uint32(b[1:])), nil
	default:
		return 0, 0, fatal
	***REMOVED***
***REMOVED***
