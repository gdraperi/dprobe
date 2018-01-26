package msgp

import (
	"io"
	"math"
	"sync"
	"time"

	"github.com/philhofer/fwd"
)

// where we keep old *Readers
var readerPool = sync.Pool***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return &Reader***REMOVED******REMOVED*** ***REMOVED******REMOVED***

// Type is a MessagePack wire type,
// including this package's built-in
// extension types.
type Type byte

// MessagePack Types
//
// The zero value of Type
// is InvalidType.
const (
	InvalidType Type = iota

	// MessagePack built-in types

	StrType
	BinType
	MapType
	ArrayType
	Float64Type
	Float32Type
	BoolType
	IntType
	UintType
	NilType
	ExtensionType

	// pseudo-types provided
	// by extensions

	Complex64Type
	Complex128Type
	TimeType

	_maxtype
)

// String implements fmt.Stringer
func (t Type) String() string ***REMOVED***
	switch t ***REMOVED***
	case StrType:
		return "str"
	case BinType:
		return "bin"
	case MapType:
		return "map"
	case ArrayType:
		return "array"
	case Float64Type:
		return "float64"
	case Float32Type:
		return "float32"
	case BoolType:
		return "bool"
	case UintType:
		return "uint"
	case IntType:
		return "int"
	case ExtensionType:
		return "ext"
	case NilType:
		return "nil"
	default:
		return "<invalid>"
	***REMOVED***
***REMOVED***

func freeR(m *Reader) ***REMOVED***
	readerPool.Put(m)
***REMOVED***

// Unmarshaler is the interface fulfilled
// by objects that know how to unmarshal
// themselves from MessagePack.
// UnmarshalMsg unmarshals the object
// from binary, returing any leftover
// bytes and any errors encountered.
type Unmarshaler interface ***REMOVED***
	UnmarshalMsg([]byte) ([]byte, error)
***REMOVED***

// Decodable is the interface fulfilled
// by objects that know how to read
// themselves from a *Reader.
type Decodable interface ***REMOVED***
	DecodeMsg(*Reader) error
***REMOVED***

// Decode decodes 'd' from 'r'.
func Decode(r io.Reader, d Decodable) error ***REMOVED***
	rd := NewReader(r)
	err := d.DecodeMsg(rd)
	freeR(rd)
	return err
***REMOVED***

// NewReader returns a *Reader that
// reads from the provided reader. The
// reader will be buffered.
func NewReader(r io.Reader) *Reader ***REMOVED***
	p := readerPool.Get().(*Reader)
	if p.R == nil ***REMOVED***
		p.R = fwd.NewReader(r)
	***REMOVED*** else ***REMOVED***
		p.R.Reset(r)
	***REMOVED***
	return p
***REMOVED***

// NewReaderSize returns a *Reader with a buffer of the given size.
// (This is vastly preferable to passing the decoder a reader that is already buffered.)
func NewReaderSize(r io.Reader, sz int) *Reader ***REMOVED***
	return &Reader***REMOVED***R: fwd.NewReaderSize(r, sz)***REMOVED***
***REMOVED***

// Reader wraps an io.Reader and provides
// methods to read MessagePack-encoded values
// from it. Readers are buffered.
type Reader struct ***REMOVED***
	// R is the buffered reader
	// that the Reader uses
	// to decode MessagePack.
	// The Reader itself
	// is stateless; all the
	// buffering is done
	// within R.
	R       *fwd.Reader
	scratch []byte
***REMOVED***

// Read implements `io.Reader`
func (m *Reader) Read(p []byte) (int, error) ***REMOVED***
	return m.R.Read(p)
***REMOVED***

// CopyNext reads the next object from m without decoding it and writes it to w.
// It avoids unnecessary copies internally.
func (m *Reader) CopyNext(w io.Writer) (int64, error) ***REMOVED***
	sz, o, err := getNextSize(m.R)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	var n int64
	// Opportunistic optimization: if we can fit the whole thing in the m.R
	// buffer, then just get a pointer to that, and pass it to w.Write,
	// avoiding an allocation.
	if int(sz) <= m.R.BufferSize() ***REMOVED***
		var nn int
		var buf []byte
		buf, err = m.R.Next(int(sz))
		if err != nil ***REMOVED***
			if err == io.ErrUnexpectedEOF ***REMOVED***
				err = ErrShortBytes
			***REMOVED***
			return 0, err
		***REMOVED***
		nn, err = w.Write(buf)
		n += int64(nn)
	***REMOVED*** else ***REMOVED***
		// Fall back to io.CopyN.
		// May avoid allocating if w is a ReaderFrom (e.g. bytes.Buffer)
		n, err = io.CopyN(w, m.R, int64(sz))
		if err == io.ErrUnexpectedEOF ***REMOVED***
			err = ErrShortBytes
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return n, err
	***REMOVED*** else if n < int64(sz) ***REMOVED***
		return n, io.ErrShortWrite
	***REMOVED***

	// for maps and slices, read elements
	for x := uintptr(0); x < o; x++ ***REMOVED***
		var n2 int64
		n2, err = m.CopyNext(w)
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		n += n2
	***REMOVED***
	return n, nil
***REMOVED***

// ReadFull implements `io.ReadFull`
func (m *Reader) ReadFull(p []byte) (int, error) ***REMOVED***
	return m.R.ReadFull(p)
***REMOVED***

// Reset resets the underlying reader.
func (m *Reader) Reset(r io.Reader) ***REMOVED*** m.R.Reset(r) ***REMOVED***

// Buffered returns the number of bytes currently in the read buffer.
func (m *Reader) Buffered() int ***REMOVED*** return m.R.Buffered() ***REMOVED***

// BufferSize returns the capacity of the read buffer.
func (m *Reader) BufferSize() int ***REMOVED*** return m.R.BufferSize() ***REMOVED***

// NextType returns the next object type to be decoded.
func (m *Reader) NextType() (Type, error) ***REMOVED***
	p, err := m.R.Peek(1)
	if err != nil ***REMOVED***
		return InvalidType, err
	***REMOVED***
	t := getType(p[0])
	if t == InvalidType ***REMOVED***
		return t, InvalidPrefixError(p[0])
	***REMOVED***
	if t == ExtensionType ***REMOVED***
		v, err := m.peekExtensionType()
		if err != nil ***REMOVED***
			return InvalidType, err
		***REMOVED***
		switch v ***REMOVED***
		case Complex64Extension:
			return Complex64Type, nil
		case Complex128Extension:
			return Complex128Type, nil
		case TimeExtension:
			return TimeType, nil
		***REMOVED***
	***REMOVED***
	return t, nil
***REMOVED***

// IsNil returns whether or not
// the next byte is a null messagepack byte
func (m *Reader) IsNil() bool ***REMOVED***
	p, err := m.R.Peek(1)
	return err == nil && p[0] == mnil
***REMOVED***

// getNextSize returns the size of the next object on the wire.
// returns (obj size, obj elements, error)
// only maps and arrays have non-zero obj elements
// for maps and arrays, obj size does not include elements
//
// use uintptr b/c it's guaranteed to be large enough
// to hold whatever we can fit in memory.
func getNextSize(r *fwd.Reader) (uintptr, uintptr, error) ***REMOVED***
	b, err := r.Peek(1)
	if err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***
	lead := b[0]
	spec := &sizes[lead]
	size, mode := spec.size, spec.extra
	if size == 0 ***REMOVED***
		return 0, 0, InvalidPrefixError(lead)
	***REMOVED***
	if mode >= 0 ***REMOVED***
		return uintptr(size), uintptr(mode), nil
	***REMOVED***
	b, err = r.Peek(int(size))
	if err != nil ***REMOVED***
		return 0, 0, err
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

// Skip skips over the next object, regardless of
// its type. If it is an array or map, the whole array
// or map will be skipped.
func (m *Reader) Skip() error ***REMOVED***
	var (
		v   uintptr // bytes
		o   uintptr // objects
		err error
		p   []byte
	)

	// we can use the faster
	// method if we have enough
	// buffered data
	if m.R.Buffered() >= 5 ***REMOVED***
		p, err = m.R.Peek(5)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v, o, err = getSize(p)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		v, o, err = getNextSize(m.R)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// 'v' is always non-zero
	// if err == nil
	_, err = m.R.Skip(int(v))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// for maps and slices, skip elements
	for x := uintptr(0); x < o; x++ ***REMOVED***
		err = m.Skip()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// ReadMapHeader reads the next object
// as a map header and returns the size
// of the map and the number of bytes written.
// It will return a TypeError***REMOVED******REMOVED*** if the next
// object is not a map.
func (m *Reader) ReadMapHeader() (sz uint32, err error) ***REMOVED***
	var p []byte
	var lead byte
	p, err = m.R.Peek(1)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	lead = p[0]
	if isfixmap(lead) ***REMOVED***
		sz = uint32(rfixmap(lead))
		_, err = m.R.Skip(1)
		return
	***REMOVED***
	switch lead ***REMOVED***
	case mmap16:
		p, err = m.R.Next(3)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		sz = uint32(big.Uint16(p[1:]))
		return
	case mmap32:
		p, err = m.R.Next(5)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		sz = big.Uint32(p[1:])
		return
	default:
		err = badPrefix(MapType, lead)
		return
	***REMOVED***
***REMOVED***

// ReadMapKey reads either a 'str' or 'bin' field from
// the reader and returns the value as a []byte. It uses
// scratch for storage if it is large enough.
func (m *Reader) ReadMapKey(scratch []byte) ([]byte, error) ***REMOVED***
	out, err := m.ReadStringAsBytes(scratch)
	if err != nil ***REMOVED***
		if tperr, ok := err.(TypeError); ok && tperr.Encoded == BinType ***REMOVED***
			return m.ReadBytes(scratch)
		***REMOVED***
		return nil, err
	***REMOVED***
	return out, nil
***REMOVED***

// MapKeyPtr returns a []byte pointing to the contents
// of a valid map key. The key cannot be empty, and it
// must be shorter than the total buffer size of the
// *Reader. Additionally, the returned slice is only
// valid until the next *Reader method call. Users
// should exercise extreme care when using this
// method; writing into the returned slice may
// corrupt future reads.
func (m *Reader) ReadMapKeyPtr() ([]byte, error) ***REMOVED***
	p, err := m.R.Peek(1)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	lead := p[0]
	var read int
	if isfixstr(lead) ***REMOVED***
		read = int(rfixstr(lead))
		m.R.Skip(1)
		goto fill
	***REMOVED***
	switch lead ***REMOVED***
	case mstr8, mbin8:
		p, err = m.R.Next(2)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		read = int(p[1])
	case mstr16, mbin16:
		p, err = m.R.Next(3)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		read = int(big.Uint16(p[1:]))
	case mstr32, mbin32:
		p, err = m.R.Next(5)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		read = int(big.Uint32(p[1:]))
	default:
		return nil, badPrefix(StrType, lead)
	***REMOVED***
fill:
	if read == 0 ***REMOVED***
		return nil, ErrShortBytes
	***REMOVED***
	return m.R.Next(read)
***REMOVED***

// ReadArrayHeader reads the next object as an
// array header and returns the size of the array
// and the number of bytes read.
func (m *Reader) ReadArrayHeader() (sz uint32, err error) ***REMOVED***
	var lead byte
	var p []byte
	p, err = m.R.Peek(1)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	lead = p[0]
	if isfixarray(lead) ***REMOVED***
		sz = uint32(rfixarray(lead))
		_, err = m.R.Skip(1)
		return
	***REMOVED***
	switch lead ***REMOVED***
	case marray16:
		p, err = m.R.Next(3)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		sz = uint32(big.Uint16(p[1:]))
		return

	case marray32:
		p, err = m.R.Next(5)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		sz = big.Uint32(p[1:])
		return

	default:
		err = badPrefix(ArrayType, lead)
		return
	***REMOVED***
***REMOVED***

// ReadNil reads a 'nil' MessagePack byte from the reader
func (m *Reader) ReadNil() error ***REMOVED***
	p, err := m.R.Peek(1)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if p[0] != mnil ***REMOVED***
		return badPrefix(NilType, p[0])
	***REMOVED***
	_, err = m.R.Skip(1)
	return err
***REMOVED***

// ReadFloat64 reads a float64 from the reader.
// (If the value on the wire is encoded as a float32,
// it will be up-cast to a float64.)
func (m *Reader) ReadFloat64() (f float64, err error) ***REMOVED***
	var p []byte
	p, err = m.R.Peek(9)
	if err != nil ***REMOVED***
		// we'll allow a coversion from float32 to float64,
		// since we don't lose any precision
		if err == io.EOF && len(p) > 0 && p[0] == mfloat32 ***REMOVED***
			ef, err := m.ReadFloat32()
			return float64(ef), err
		***REMOVED***
		return
	***REMOVED***
	if p[0] != mfloat64 ***REMOVED***
		// see above
		if p[0] == mfloat32 ***REMOVED***
			ef, err := m.ReadFloat32()
			return float64(ef), err
		***REMOVED***
		err = badPrefix(Float64Type, p[0])
		return
	***REMOVED***
	f = math.Float64frombits(getMuint64(p))
	_, err = m.R.Skip(9)
	return
***REMOVED***

// ReadFloat32 reads a float32 from the reader
func (m *Reader) ReadFloat32() (f float32, err error) ***REMOVED***
	var p []byte
	p, err = m.R.Peek(5)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if p[0] != mfloat32 ***REMOVED***
		err = badPrefix(Float32Type, p[0])
		return
	***REMOVED***
	f = math.Float32frombits(getMuint32(p))
	_, err = m.R.Skip(5)
	return
***REMOVED***

// ReadBool reads a bool from the reader
func (m *Reader) ReadBool() (b bool, err error) ***REMOVED***
	var p []byte
	p, err = m.R.Peek(1)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	switch p[0] ***REMOVED***
	case mtrue:
		b = true
	case mfalse:
	default:
		err = badPrefix(BoolType, p[0])
		return
	***REMOVED***
	_, err = m.R.Skip(1)
	return
***REMOVED***

// ReadInt64 reads an int64 from the reader
func (m *Reader) ReadInt64() (i int64, err error) ***REMOVED***
	var p []byte
	var lead byte
	p, err = m.R.Peek(1)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	lead = p[0]

	if isfixint(lead) ***REMOVED***
		i = int64(rfixint(lead))
		_, err = m.R.Skip(1)
		return
	***REMOVED*** else if isnfixint(lead) ***REMOVED***
		i = int64(rnfixint(lead))
		_, err = m.R.Skip(1)
		return
	***REMOVED***

	switch lead ***REMOVED***
	case mint8:
		p, err = m.R.Next(2)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		i = int64(getMint8(p))
		return

	case mint16:
		p, err = m.R.Next(3)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		i = int64(getMint16(p))
		return

	case mint32:
		p, err = m.R.Next(5)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		i = int64(getMint32(p))
		return

	case mint64:
		p, err = m.R.Next(9)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		i = getMint64(p)
		return

	default:
		err = badPrefix(IntType, lead)
		return
	***REMOVED***
***REMOVED***

// ReadInt32 reads an int32 from the reader
func (m *Reader) ReadInt32() (i int32, err error) ***REMOVED***
	var in int64
	in, err = m.ReadInt64()
	if in > math.MaxInt32 || in < math.MinInt32 ***REMOVED***
		err = IntOverflow***REMOVED***Value: in, FailedBitsize: 32***REMOVED***
		return
	***REMOVED***
	i = int32(in)
	return
***REMOVED***

// ReadInt16 reads an int16 from the reader
func (m *Reader) ReadInt16() (i int16, err error) ***REMOVED***
	var in int64
	in, err = m.ReadInt64()
	if in > math.MaxInt16 || in < math.MinInt16 ***REMOVED***
		err = IntOverflow***REMOVED***Value: in, FailedBitsize: 16***REMOVED***
		return
	***REMOVED***
	i = int16(in)
	return
***REMOVED***

// ReadInt8 reads an int8 from the reader
func (m *Reader) ReadInt8() (i int8, err error) ***REMOVED***
	var in int64
	in, err = m.ReadInt64()
	if in > math.MaxInt8 || in < math.MinInt8 ***REMOVED***
		err = IntOverflow***REMOVED***Value: in, FailedBitsize: 8***REMOVED***
		return
	***REMOVED***
	i = int8(in)
	return
***REMOVED***

// ReadInt reads an int from the reader
func (m *Reader) ReadInt() (i int, err error) ***REMOVED***
	if smallint ***REMOVED***
		var in int32
		in, err = m.ReadInt32()
		i = int(in)
		return
	***REMOVED***
	var in int64
	in, err = m.ReadInt64()
	i = int(in)
	return
***REMOVED***

// ReadUint64 reads a uint64 from the reader
func (m *Reader) ReadUint64() (u uint64, err error) ***REMOVED***
	var p []byte
	var lead byte
	p, err = m.R.Peek(1)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	lead = p[0]
	if isfixint(lead) ***REMOVED***
		u = uint64(rfixint(lead))
		_, err = m.R.Skip(1)
		return
	***REMOVED***
	switch lead ***REMOVED***
	case muint8:
		p, err = m.R.Next(2)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		u = uint64(getMuint8(p))
		return

	case muint16:
		p, err = m.R.Next(3)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		u = uint64(getMuint16(p))
		return

	case muint32:
		p, err = m.R.Next(5)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		u = uint64(getMuint32(p))
		return

	case muint64:
		p, err = m.R.Next(9)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		u = getMuint64(p)
		return

	default:
		err = badPrefix(UintType, lead)
		return

	***REMOVED***
***REMOVED***

// ReadUint32 reads a uint32 from the reader
func (m *Reader) ReadUint32() (u uint32, err error) ***REMOVED***
	var in uint64
	in, err = m.ReadUint64()
	if in > math.MaxUint32 ***REMOVED***
		err = UintOverflow***REMOVED***Value: in, FailedBitsize: 32***REMOVED***
		return
	***REMOVED***
	u = uint32(in)
	return
***REMOVED***

// ReadUint16 reads a uint16 from the reader
func (m *Reader) ReadUint16() (u uint16, err error) ***REMOVED***
	var in uint64
	in, err = m.ReadUint64()
	if in > math.MaxUint16 ***REMOVED***
		err = UintOverflow***REMOVED***Value: in, FailedBitsize: 16***REMOVED***
		return
	***REMOVED***
	u = uint16(in)
	return
***REMOVED***

// ReadUint8 reads a uint8 from the reader
func (m *Reader) ReadUint8() (u uint8, err error) ***REMOVED***
	var in uint64
	in, err = m.ReadUint64()
	if in > math.MaxUint8 ***REMOVED***
		err = UintOverflow***REMOVED***Value: in, FailedBitsize: 8***REMOVED***
		return
	***REMOVED***
	u = uint8(in)
	return
***REMOVED***

// ReadUint reads a uint from the reader
func (m *Reader) ReadUint() (u uint, err error) ***REMOVED***
	if smallint ***REMOVED***
		var un uint32
		un, err = m.ReadUint32()
		u = uint(un)
		return
	***REMOVED***
	var un uint64
	un, err = m.ReadUint64()
	u = uint(un)
	return
***REMOVED***

// ReadByte is analogous to ReadUint8.
//
// NOTE: this is *not* an implementation
// of io.ByteReader.
func (m *Reader) ReadByte() (b byte, err error) ***REMOVED***
	var in uint64
	in, err = m.ReadUint64()
	if in > math.MaxUint8 ***REMOVED***
		err = UintOverflow***REMOVED***Value: in, FailedBitsize: 8***REMOVED***
		return
	***REMOVED***
	b = byte(in)
	return
***REMOVED***

// ReadBytes reads a MessagePack 'bin' object
// from the reader and returns its value. It may
// use 'scratch' for storage if it is non-nil.
func (m *Reader) ReadBytes(scratch []byte) (b []byte, err error) ***REMOVED***
	var p []byte
	var lead byte
	p, err = m.R.Peek(2)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	lead = p[0]
	var read int64
	switch lead ***REMOVED***
	case mbin8:
		read = int64(p[1])
		m.R.Skip(2)
	case mbin16:
		p, err = m.R.Next(3)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		read = int64(big.Uint16(p[1:]))
	case mbin32:
		p, err = m.R.Next(5)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		read = int64(big.Uint32(p[1:]))
	default:
		err = badPrefix(BinType, lead)
		return
	***REMOVED***
	if int64(cap(scratch)) < read ***REMOVED***
		b = make([]byte, read)
	***REMOVED*** else ***REMOVED***
		b = scratch[0:read]
	***REMOVED***
	_, err = m.R.ReadFull(b)
	return
***REMOVED***

// ReadBytesHeader reads the size header
// of a MessagePack 'bin' object. The user
// is responsible for dealing with the next
// 'sz' bytes from the reader in an application-specific
// way.
func (m *Reader) ReadBytesHeader() (sz uint32, err error) ***REMOVED***
	var p []byte
	p, err = m.R.Peek(1)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	switch p[0] ***REMOVED***
	case mbin8:
		p, err = m.R.Next(2)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		sz = uint32(p[1])
		return
	case mbin16:
		p, err = m.R.Next(3)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		sz = uint32(big.Uint16(p[1:]))
		return
	case mbin32:
		p, err = m.R.Next(5)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		sz = uint32(big.Uint32(p[1:]))
		return
	default:
		err = badPrefix(BinType, p[0])
		return
	***REMOVED***
***REMOVED***

// ReadExactBytes reads a MessagePack 'bin'-encoded
// object off of the wire into the provided slice. An
// ArrayError will be returned if the object is not
// exactly the length of the input slice.
func (m *Reader) ReadExactBytes(into []byte) error ***REMOVED***
	p, err := m.R.Peek(2)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	lead := p[0]
	var read int64 // bytes to read
	var skip int   // prefix size to skip
	switch lead ***REMOVED***
	case mbin8:
		read = int64(p[1])
		skip = 2
	case mbin16:
		p, err = m.R.Peek(3)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		read = int64(big.Uint16(p[1:]))
		skip = 3
	case mbin32:
		p, err = m.R.Peek(5)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		read = int64(big.Uint32(p[1:]))
		skip = 5
	default:
		return badPrefix(BinType, lead)
	***REMOVED***
	if read != int64(len(into)) ***REMOVED***
		return ArrayError***REMOVED***Wanted: uint32(len(into)), Got: uint32(read)***REMOVED***
	***REMOVED***
	m.R.Skip(skip)
	_, err = m.R.ReadFull(into)
	return err
***REMOVED***

// ReadStringAsBytes reads a MessagePack 'str' (utf-8) string
// and returns its value as bytes. It may use 'scratch' for storage
// if it is non-nil.
func (m *Reader) ReadStringAsBytes(scratch []byte) (b []byte, err error) ***REMOVED***
	var p []byte
	var lead byte
	p, err = m.R.Peek(1)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	lead = p[0]
	var read int64

	if isfixstr(lead) ***REMOVED***
		read = int64(rfixstr(lead))
		m.R.Skip(1)
		goto fill
	***REMOVED***

	switch lead ***REMOVED***
	case mstr8:
		p, err = m.R.Next(2)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		read = int64(uint8(p[1]))
	case mstr16:
		p, err = m.R.Next(3)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		read = int64(big.Uint16(p[1:]))
	case mstr32:
		p, err = m.R.Next(5)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		read = int64(big.Uint32(p[1:]))
	default:
		err = badPrefix(StrType, lead)
		return
	***REMOVED***
fill:
	if int64(cap(scratch)) < read ***REMOVED***
		b = make([]byte, read)
	***REMOVED*** else ***REMOVED***
		b = scratch[0:read]
	***REMOVED***
	_, err = m.R.ReadFull(b)
	return
***REMOVED***

// ReadStringHeader reads a string header
// off of the wire. The user is then responsible
// for dealing with the next 'sz' bytes from
// the reader in an application-specific manner.
func (m *Reader) ReadStringHeader() (sz uint32, err error) ***REMOVED***
	var p []byte
	p, err = m.R.Peek(1)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	lead := p[0]
	if isfixstr(lead) ***REMOVED***
		sz = uint32(rfixstr(lead))
		m.R.Skip(1)
		return
	***REMOVED***
	switch lead ***REMOVED***
	case mstr8:
		p, err = m.R.Next(2)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		sz = uint32(p[1])
		return
	case mstr16:
		p, err = m.R.Next(3)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		sz = uint32(big.Uint16(p[1:]))
		return
	case mstr32:
		p, err = m.R.Next(5)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		sz = big.Uint32(p[1:])
		return
	default:
		err = badPrefix(StrType, lead)
		return
	***REMOVED***
***REMOVED***

// ReadString reads a utf-8 string from the reader
func (m *Reader) ReadString() (s string, err error) ***REMOVED***
	var p []byte
	var lead byte
	var read int64
	p, err = m.R.Peek(1)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	lead = p[0]

	if isfixstr(lead) ***REMOVED***
		read = int64(rfixstr(lead))
		m.R.Skip(1)
		goto fill
	***REMOVED***

	switch lead ***REMOVED***
	case mstr8:
		p, err = m.R.Next(2)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		read = int64(uint8(p[1]))
	case mstr16:
		p, err = m.R.Next(3)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		read = int64(big.Uint16(p[1:]))
	case mstr32:
		p, err = m.R.Next(5)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		read = int64(big.Uint32(p[1:]))
	default:
		err = badPrefix(StrType, lead)
		return
	***REMOVED***
fill:
	if read == 0 ***REMOVED***
		s, err = "", nil
		return
	***REMOVED***
	// reading into the memory
	// that will become the string
	// itself has vastly superior
	// worst-case performance, because
	// the reader buffer doesn't have
	// to be large enough to hold the string.
	// the idea here is to make it more
	// difficult for someone malicious
	// to cause the system to run out of
	// memory by sending very large strings.
	//
	// NOTE: this works because the argument
	// passed to (*fwd.Reader).ReadFull escapes
	// to the heap; its argument may, in turn,
	// be passed to the underlying reader, and
	// thus escape analysis *must* conclude that
	// 'out' escapes.
	out := make([]byte, read)
	_, err = m.R.ReadFull(out)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	s = UnsafeString(out)
	return
***REMOVED***

// ReadComplex64 reads a complex64 from the reader
func (m *Reader) ReadComplex64() (f complex64, err error) ***REMOVED***
	var p []byte
	p, err = m.R.Peek(10)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if p[0] != mfixext8 ***REMOVED***
		err = badPrefix(Complex64Type, p[0])
		return
	***REMOVED***
	if int8(p[1]) != Complex64Extension ***REMOVED***
		err = errExt(int8(p[1]), Complex64Extension)
		return
	***REMOVED***
	f = complex(math.Float32frombits(big.Uint32(p[2:])),
		math.Float32frombits(big.Uint32(p[6:])))
	_, err = m.R.Skip(10)
	return
***REMOVED***

// ReadComplex128 reads a complex128 from the reader
func (m *Reader) ReadComplex128() (f complex128, err error) ***REMOVED***
	var p []byte
	p, err = m.R.Peek(18)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if p[0] != mfixext16 ***REMOVED***
		err = badPrefix(Complex128Type, p[0])
		return
	***REMOVED***
	if int8(p[1]) != Complex128Extension ***REMOVED***
		err = errExt(int8(p[1]), Complex128Extension)
		return
	***REMOVED***
	f = complex(math.Float64frombits(big.Uint64(p[2:])),
		math.Float64frombits(big.Uint64(p[10:])))
	_, err = m.R.Skip(18)
	return
***REMOVED***

// ReadMapStrIntf reads a MessagePack map into a map[string]interface***REMOVED******REMOVED***.
// (You must pass a non-nil map into the function.)
func (m *Reader) ReadMapStrIntf(mp map[string]interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	var sz uint32
	sz, err = m.ReadMapHeader()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	for key := range mp ***REMOVED***
		delete(mp, key)
	***REMOVED***
	for i := uint32(0); i < sz; i++ ***REMOVED***
		var key string
		var val interface***REMOVED******REMOVED***
		key, err = m.ReadString()
		if err != nil ***REMOVED***
			return
		***REMOVED***
		val, err = m.ReadIntf()
		if err != nil ***REMOVED***
			return
		***REMOVED***
		mp[key] = val
	***REMOVED***
	return
***REMOVED***

// ReadTime reads a time.Time object from the reader.
// The returned time's location will be set to time.Local.
func (m *Reader) ReadTime() (t time.Time, err error) ***REMOVED***
	var p []byte
	p, err = m.R.Peek(15)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if p[0] != mext8 || p[1] != 12 ***REMOVED***
		err = badPrefix(TimeType, p[0])
		return
	***REMOVED***
	if int8(p[2]) != TimeExtension ***REMOVED***
		err = errExt(int8(p[2]), TimeExtension)
		return
	***REMOVED***
	sec, nsec := getUnix(p[3:])
	t = time.Unix(sec, int64(nsec)).Local()
	_, err = m.R.Skip(15)
	return
***REMOVED***

// ReadIntf reads out the next object as a raw interface***REMOVED******REMOVED***.
// Arrays are decoded as []interface***REMOVED******REMOVED***, and maps are decoded
// as map[string]interface***REMOVED******REMOVED***. Integers are decoded as int64
// and unsigned integers are decoded as uint64.
func (m *Reader) ReadIntf() (i interface***REMOVED******REMOVED***, err error) ***REMOVED***
	var t Type
	t, err = m.NextType()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	switch t ***REMOVED***
	case BoolType:
		i, err = m.ReadBool()
		return

	case IntType:
		i, err = m.ReadInt64()
		return

	case UintType:
		i, err = m.ReadUint64()
		return

	case BinType:
		i, err = m.ReadBytes(nil)
		return

	case StrType:
		i, err = m.ReadString()
		return

	case Complex64Type:
		i, err = m.ReadComplex64()
		return

	case Complex128Type:
		i, err = m.ReadComplex128()
		return

	case TimeType:
		i, err = m.ReadTime()
		return

	case ExtensionType:
		var t int8
		t, err = m.peekExtensionType()
		if err != nil ***REMOVED***
			return
		***REMOVED***
		f, ok := extensionReg[t]
		if ok ***REMOVED***
			e := f()
			err = m.ReadExtension(e)
			i = e
			return
		***REMOVED***
		var e RawExtension
		e.Type = t
		err = m.ReadExtension(&e)
		i = &e
		return

	case MapType:
		mp := make(map[string]interface***REMOVED******REMOVED***)
		err = m.ReadMapStrIntf(mp)
		i = mp
		return

	case NilType:
		err = m.ReadNil()
		i = nil
		return

	case Float32Type:
		i, err = m.ReadFloat32()
		return

	case Float64Type:
		i, err = m.ReadFloat64()
		return

	case ArrayType:
		var sz uint32
		sz, err = m.ReadArrayHeader()

		if err != nil ***REMOVED***
			return
		***REMOVED***
		out := make([]interface***REMOVED******REMOVED***, int(sz))
		for j := range out ***REMOVED***
			out[j], err = m.ReadIntf()
			if err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		i = out
		return

	default:
		return nil, fatal // unreachable
	***REMOVED***
***REMOVED***
