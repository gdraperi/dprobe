package msgp

import (
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"sync"
	"time"
)

// Sizer is an interface implemented
// by types that can estimate their
// size when MessagePack encoded.
// This interface is optional, but
// encoding/marshaling implementations
// may use this as a way to pre-allocate
// memory for serialization.
type Sizer interface ***REMOVED***
	Msgsize() int
***REMOVED***

var (
	// Nowhere is an io.Writer to nowhere
	Nowhere io.Writer = nwhere***REMOVED******REMOVED***

	btsType    = reflect.TypeOf(([]byte)(nil))
	writerPool = sync.Pool***REMOVED***
		New: func() interface***REMOVED******REMOVED*** ***REMOVED***
			return &Writer***REMOVED***buf: make([]byte, 2048)***REMOVED***
		***REMOVED***,
	***REMOVED***
)

func popWriter(w io.Writer) *Writer ***REMOVED***
	wr := writerPool.Get().(*Writer)
	wr.Reset(w)
	return wr
***REMOVED***

func pushWriter(wr *Writer) ***REMOVED***
	wr.w = nil
	wr.wloc = 0
	writerPool.Put(wr)
***REMOVED***

// freeW frees a writer for use
// by other processes. It is not necessary
// to call freeW on a writer. However, maintaining
// a reference to a *Writer after calling freeW on
// it will cause undefined behavior.
func freeW(w *Writer) ***REMOVED*** pushWriter(w) ***REMOVED***

// Require ensures that cap(old)-len(old) >= extra.
func Require(old []byte, extra int) []byte ***REMOVED***
	l := len(old)
	c := cap(old)
	r := l + extra
	if c >= r ***REMOVED***
		return old
	***REMOVED*** else if l == 0 ***REMOVED***
		return make([]byte, 0, extra)
	***REMOVED***
	// the new size is the greater
	// of double the old capacity
	// and the sum of the old length
	// and the number of new bytes
	// necessary.
	c <<= 1
	if c < r ***REMOVED***
		c = r
	***REMOVED***
	n := make([]byte, l, c)
	copy(n, old)
	return n
***REMOVED***

// nowhere writer
type nwhere struct***REMOVED******REMOVED***

func (n nwhere) Write(p []byte) (int, error) ***REMOVED*** return len(p), nil ***REMOVED***

// Marshaler is the interface implemented
// by types that know how to marshal themselves
// as MessagePack. MarshalMsg appends the marshalled
// form of the object to the provided
// byte slice, returning the extended
// slice and any errors encountered.
type Marshaler interface ***REMOVED***
	MarshalMsg([]byte) ([]byte, error)
***REMOVED***

// Encodable is the interface implemented
// by types that know how to write themselves
// as MessagePack using a *msgp.Writer.
type Encodable interface ***REMOVED***
	EncodeMsg(*Writer) error
***REMOVED***

// Writer is a buffered writer
// that can be used to write
// MessagePack objects to an io.Writer.
// You must call *Writer.Flush() in order
// to flush all of the buffered data
// to the underlying writer.
type Writer struct ***REMOVED***
	w    io.Writer
	buf  []byte
	wloc int
***REMOVED***

// NewWriter returns a new *Writer.
func NewWriter(w io.Writer) *Writer ***REMOVED***
	if wr, ok := w.(*Writer); ok ***REMOVED***
		return wr
	***REMOVED***
	return popWriter(w)
***REMOVED***

// NewWriterSize returns a writer with a custom buffer size.
func NewWriterSize(w io.Writer, sz int) *Writer ***REMOVED***
	// we must be able to require() 18
	// contiguous bytes, so that is the
	// practical minimum buffer size
	if sz < 18 ***REMOVED***
		sz = 18
	***REMOVED***

	return &Writer***REMOVED***
		w:   w,
		buf: make([]byte, sz),
	***REMOVED***
***REMOVED***

// Encode encodes an Encodable to an io.Writer.
func Encode(w io.Writer, e Encodable) error ***REMOVED***
	wr := NewWriter(w)
	err := e.EncodeMsg(wr)
	if err == nil ***REMOVED***
		err = wr.Flush()
	***REMOVED***
	freeW(wr)
	return err
***REMOVED***

func (mw *Writer) flush() error ***REMOVED***
	if mw.wloc == 0 ***REMOVED***
		return nil
	***REMOVED***
	n, err := mw.w.Write(mw.buf[:mw.wloc])
	if err != nil ***REMOVED***
		if n > 0 ***REMOVED***
			mw.wloc = copy(mw.buf, mw.buf[n:mw.wloc])
		***REMOVED***
		return err
	***REMOVED***
	mw.wloc = 0
	return nil
***REMOVED***

// Flush flushes all of the buffered
// data to the underlying writer.
func (mw *Writer) Flush() error ***REMOVED*** return mw.flush() ***REMOVED***

// Buffered returns the number bytes in the write buffer
func (mw *Writer) Buffered() int ***REMOVED*** return len(mw.buf) - mw.wloc ***REMOVED***

func (mw *Writer) avail() int ***REMOVED*** return len(mw.buf) - mw.wloc ***REMOVED***

func (mw *Writer) bufsize() int ***REMOVED*** return len(mw.buf) ***REMOVED***

// NOTE: this should only be called with
// a number that is guaranteed to be less than
// len(mw.buf). typically, it is called with a constant.
//
// NOTE: this is a hot code path
func (mw *Writer) require(n int) (int, error) ***REMOVED***
	c := len(mw.buf)
	wl := mw.wloc
	if c-wl < n ***REMOVED***
		if err := mw.flush(); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		wl = mw.wloc
	***REMOVED***
	mw.wloc += n
	return wl, nil
***REMOVED***

func (mw *Writer) Append(b ...byte) error ***REMOVED***
	if mw.avail() < len(b) ***REMOVED***
		err := mw.flush()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	mw.wloc += copy(mw.buf[mw.wloc:], b)
	return nil
***REMOVED***

// push one byte onto the buffer
//
// NOTE: this is a hot code path
func (mw *Writer) push(b byte) error ***REMOVED***
	if mw.wloc == len(mw.buf) ***REMOVED***
		if err := mw.flush(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	mw.buf[mw.wloc] = b
	mw.wloc++
	return nil
***REMOVED***

func (mw *Writer) prefix8(b byte, u uint8) error ***REMOVED***
	const need = 2
	if len(mw.buf)-mw.wloc < need ***REMOVED***
		if err := mw.flush(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	prefixu8(mw.buf[mw.wloc:], b, u)
	mw.wloc += need
	return nil
***REMOVED***

func (mw *Writer) prefix16(b byte, u uint16) error ***REMOVED***
	const need = 3
	if len(mw.buf)-mw.wloc < need ***REMOVED***
		if err := mw.flush(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	prefixu16(mw.buf[mw.wloc:], b, u)
	mw.wloc += need
	return nil
***REMOVED***

func (mw *Writer) prefix32(b byte, u uint32) error ***REMOVED***
	const need = 5
	if len(mw.buf)-mw.wloc < need ***REMOVED***
		if err := mw.flush(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	prefixu32(mw.buf[mw.wloc:], b, u)
	mw.wloc += need
	return nil
***REMOVED***

func (mw *Writer) prefix64(b byte, u uint64) error ***REMOVED***
	const need = 9
	if len(mw.buf)-mw.wloc < need ***REMOVED***
		if err := mw.flush(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	prefixu64(mw.buf[mw.wloc:], b, u)
	mw.wloc += need
	return nil
***REMOVED***

// Write implements io.Writer, and writes
// data directly to the buffer.
func (mw *Writer) Write(p []byte) (int, error) ***REMOVED***
	l := len(p)
	if mw.avail() < l ***REMOVED***
		if err := mw.flush(); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		if l > len(mw.buf) ***REMOVED***
			return mw.w.Write(p)
		***REMOVED***
	***REMOVED***
	mw.wloc += copy(mw.buf[mw.wloc:], p)
	return l, nil
***REMOVED***

// implements io.WriteString
func (mw *Writer) writeString(s string) error ***REMOVED***
	l := len(s)
	if mw.avail() < l ***REMOVED***
		if err := mw.flush(); err != nil ***REMOVED***
			return err
		***REMOVED***
		if l > len(mw.buf) ***REMOVED***
			_, err := io.WriteString(mw.w, s)
			return err
		***REMOVED***
	***REMOVED***
	mw.wloc += copy(mw.buf[mw.wloc:], s)
	return nil
***REMOVED***

// Reset changes the underlying writer used by the Writer
func (mw *Writer) Reset(w io.Writer) ***REMOVED***
	mw.buf = mw.buf[:cap(mw.buf)]
	mw.w = w
	mw.wloc = 0
***REMOVED***

// WriteMapHeader writes a map header of the given
// size to the writer
func (mw *Writer) WriteMapHeader(sz uint32) error ***REMOVED***
	switch ***REMOVED***
	case sz <= 15:
		return mw.push(wfixmap(uint8(sz)))
	case sz <= math.MaxUint16:
		return mw.prefix16(mmap16, uint16(sz))
	default:
		return mw.prefix32(mmap32, sz)
	***REMOVED***
***REMOVED***

// WriteArrayHeader writes an array header of the
// given size to the writer
func (mw *Writer) WriteArrayHeader(sz uint32) error ***REMOVED***
	switch ***REMOVED***
	case sz <= 15:
		return mw.push(wfixarray(uint8(sz)))
	case sz <= math.MaxUint16:
		return mw.prefix16(marray16, uint16(sz))
	default:
		return mw.prefix32(marray32, sz)
	***REMOVED***
***REMOVED***

// WriteNil writes a nil byte to the buffer
func (mw *Writer) WriteNil() error ***REMOVED***
	return mw.push(mnil)
***REMOVED***

// WriteFloat64 writes a float64 to the writer
func (mw *Writer) WriteFloat64(f float64) error ***REMOVED***
	return mw.prefix64(mfloat64, math.Float64bits(f))
***REMOVED***

// WriteFloat32 writes a float32 to the writer
func (mw *Writer) WriteFloat32(f float32) error ***REMOVED***
	return mw.prefix32(mfloat32, math.Float32bits(f))
***REMOVED***

// WriteInt64 writes an int64 to the writer
func (mw *Writer) WriteInt64(i int64) error ***REMOVED***
	if i >= 0 ***REMOVED***
		switch ***REMOVED***
		case i <= math.MaxInt8:
			return mw.push(wfixint(uint8(i)))
		case i <= math.MaxInt16:
			return mw.prefix16(mint16, uint16(i))
		case i <= math.MaxInt32:
			return mw.prefix32(mint32, uint32(i))
		default:
			return mw.prefix64(mint64, uint64(i))
		***REMOVED***
	***REMOVED***
	switch ***REMOVED***
	case i >= -32:
		return mw.push(wnfixint(int8(i)))
	case i >= math.MinInt8:
		return mw.prefix8(mint8, uint8(i))
	case i >= math.MinInt16:
		return mw.prefix16(mint16, uint16(i))
	case i >= math.MinInt32:
		return mw.prefix32(mint32, uint32(i))
	default:
		return mw.prefix64(mint64, uint64(i))
	***REMOVED***
***REMOVED***

// WriteInt8 writes an int8 to the writer
func (mw *Writer) WriteInt8(i int8) error ***REMOVED*** return mw.WriteInt64(int64(i)) ***REMOVED***

// WriteInt16 writes an int16 to the writer
func (mw *Writer) WriteInt16(i int16) error ***REMOVED*** return mw.WriteInt64(int64(i)) ***REMOVED***

// WriteInt32 writes an int32 to the writer
func (mw *Writer) WriteInt32(i int32) error ***REMOVED*** return mw.WriteInt64(int64(i)) ***REMOVED***

// WriteInt writes an int to the writer
func (mw *Writer) WriteInt(i int) error ***REMOVED*** return mw.WriteInt64(int64(i)) ***REMOVED***

// WriteUint64 writes a uint64 to the writer
func (mw *Writer) WriteUint64(u uint64) error ***REMOVED***
	switch ***REMOVED***
	case u <= (1<<7)-1:
		return mw.push(wfixint(uint8(u)))
	case u <= math.MaxUint8:
		return mw.prefix8(muint8, uint8(u))
	case u <= math.MaxUint16:
		return mw.prefix16(muint16, uint16(u))
	case u <= math.MaxUint32:
		return mw.prefix32(muint32, uint32(u))
	default:
		return mw.prefix64(muint64, u)
	***REMOVED***
***REMOVED***

// WriteByte is analogous to WriteUint8
func (mw *Writer) WriteByte(u byte) error ***REMOVED*** return mw.WriteUint8(uint8(u)) ***REMOVED***

// WriteUint8 writes a uint8 to the writer
func (mw *Writer) WriteUint8(u uint8) error ***REMOVED*** return mw.WriteUint64(uint64(u)) ***REMOVED***

// WriteUint16 writes a uint16 to the writer
func (mw *Writer) WriteUint16(u uint16) error ***REMOVED*** return mw.WriteUint64(uint64(u)) ***REMOVED***

// WriteUint32 writes a uint32 to the writer
func (mw *Writer) WriteUint32(u uint32) error ***REMOVED*** return mw.WriteUint64(uint64(u)) ***REMOVED***

// WriteUint writes a uint to the writer
func (mw *Writer) WriteUint(u uint) error ***REMOVED*** return mw.WriteUint64(uint64(u)) ***REMOVED***

// WriteBytes writes binary as 'bin' to the writer
func (mw *Writer) WriteBytes(b []byte) error ***REMOVED***
	sz := uint32(len(b))
	var err error
	switch ***REMOVED***
	case sz <= math.MaxUint8:
		err = mw.prefix8(mbin8, uint8(sz))
	case sz <= math.MaxUint16:
		err = mw.prefix16(mbin16, uint16(sz))
	default:
		err = mw.prefix32(mbin32, sz)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = mw.Write(b)
	return err
***REMOVED***

// WriteBytesHeader writes just the size header
// of a MessagePack 'bin' object. The user is responsible
// for then writing 'sz' more bytes into the stream.
func (mw *Writer) WriteBytesHeader(sz uint32) error ***REMOVED***
	switch ***REMOVED***
	case sz <= math.MaxUint8:
		return mw.prefix8(mbin8, uint8(sz))
	case sz <= math.MaxUint16:
		return mw.prefix16(mbin16, uint16(sz))
	default:
		return mw.prefix32(mbin32, sz)
	***REMOVED***
***REMOVED***

// WriteBool writes a bool to the writer
func (mw *Writer) WriteBool(b bool) error ***REMOVED***
	if b ***REMOVED***
		return mw.push(mtrue)
	***REMOVED***
	return mw.push(mfalse)
***REMOVED***

// WriteString writes a messagepack string to the writer.
// (This is NOT an implementation of io.StringWriter)
func (mw *Writer) WriteString(s string) error ***REMOVED***
	sz := uint32(len(s))
	var err error
	switch ***REMOVED***
	case sz <= 31:
		err = mw.push(wfixstr(uint8(sz)))
	case sz <= math.MaxUint8:
		err = mw.prefix8(mstr8, uint8(sz))
	case sz <= math.MaxUint16:
		err = mw.prefix16(mstr16, uint16(sz))
	default:
		err = mw.prefix32(mstr32, sz)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return mw.writeString(s)
***REMOVED***

// WriteStringHeader writes just the string size
// header of a MessagePack 'str' object. The user
// is responsible for writing 'sz' more valid UTF-8
// bytes to the stream.
func (mw *Writer) WriteStringHeader(sz uint32) error ***REMOVED***
	switch ***REMOVED***
	case sz <= 31:
		return mw.push(wfixstr(uint8(sz)))
	case sz <= math.MaxUint8:
		return mw.prefix8(mstr8, uint8(sz))
	case sz <= math.MaxUint16:
		return mw.prefix16(mstr16, uint16(sz))
	default:
		return mw.prefix32(mstr32, sz)
	***REMOVED***
***REMOVED***

// WriteStringFromBytes writes a 'str' object
// from a []byte.
func (mw *Writer) WriteStringFromBytes(str []byte) error ***REMOVED***
	sz := uint32(len(str))
	var err error
	switch ***REMOVED***
	case sz <= 31:
		err = mw.push(wfixstr(uint8(sz)))
	case sz <= math.MaxUint8:
		err = mw.prefix8(mstr8, uint8(sz))
	case sz <= math.MaxUint16:
		err = mw.prefix16(mstr16, uint16(sz))
	default:
		err = mw.prefix32(mstr32, sz)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = mw.Write(str)
	return err
***REMOVED***

// WriteComplex64 writes a complex64 to the writer
func (mw *Writer) WriteComplex64(f complex64) error ***REMOVED***
	o, err := mw.require(10)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	mw.buf[o] = mfixext8
	mw.buf[o+1] = Complex64Extension
	big.PutUint32(mw.buf[o+2:], math.Float32bits(real(f)))
	big.PutUint32(mw.buf[o+6:], math.Float32bits(imag(f)))
	return nil
***REMOVED***

// WriteComplex128 writes a complex128 to the writer
func (mw *Writer) WriteComplex128(f complex128) error ***REMOVED***
	o, err := mw.require(18)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	mw.buf[o] = mfixext16
	mw.buf[o+1] = Complex128Extension
	big.PutUint64(mw.buf[o+2:], math.Float64bits(real(f)))
	big.PutUint64(mw.buf[o+10:], math.Float64bits(imag(f)))
	return nil
***REMOVED***

// WriteMapStrStr writes a map[string]string to the writer
func (mw *Writer) WriteMapStrStr(mp map[string]string) (err error) ***REMOVED***
	err = mw.WriteMapHeader(uint32(len(mp)))
	if err != nil ***REMOVED***
		return
	***REMOVED***
	for key, val := range mp ***REMOVED***
		err = mw.WriteString(key)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		err = mw.WriteString(val)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// WriteMapStrIntf writes a map[string]interface to the writer
func (mw *Writer) WriteMapStrIntf(mp map[string]interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	err = mw.WriteMapHeader(uint32(len(mp)))
	if err != nil ***REMOVED***
		return
	***REMOVED***
	for key, val := range mp ***REMOVED***
		err = mw.WriteString(key)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		err = mw.WriteIntf(val)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// WriteTime writes a time.Time object to the wire.
//
// Time is encoded as Unix time, which means that
// location (time zone) data is removed from the object.
// The encoded object itself is 12 bytes: 8 bytes for
// a big-endian 64-bit integer denoting seconds
// elapsed since "zero" Unix time, followed by 4 bytes
// for a big-endian 32-bit signed integer denoting
// the nanosecond offset of the time. This encoding
// is intended to ease portability across languages.
// (Note that this is *not* the standard time.Time
// binary encoding, because its implementation relies
// heavily on the internal representation used by the
// time package.)
func (mw *Writer) WriteTime(t time.Time) error ***REMOVED***
	t = t.UTC()
	o, err := mw.require(15)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	mw.buf[o] = mext8
	mw.buf[o+1] = 12
	mw.buf[o+2] = TimeExtension
	putUnix(mw.buf[o+3:], t.Unix(), int32(t.Nanosecond()))
	return nil
***REMOVED***

// WriteIntf writes the concrete type of 'v'.
// WriteIntf will error if 'v' is not one of the following:
//  - A bool, float, string, []byte, int, uint, or complex
//  - A map of supported types (with string keys)
//  - An array or slice of supported types
//  - A pointer to a supported type
//  - A type that satisfies the msgp.Encodable interface
//  - A type that satisfies the msgp.Extension interface
func (mw *Writer) WriteIntf(v interface***REMOVED******REMOVED***) error ***REMOVED***
	if v == nil ***REMOVED***
		return mw.WriteNil()
	***REMOVED***
	switch v := v.(type) ***REMOVED***

	// preferred interfaces

	case Encodable:
		return v.EncodeMsg(mw)
	case Extension:
		return mw.WriteExtension(v)

	// concrete types

	case bool:
		return mw.WriteBool(v)
	case float32:
		return mw.WriteFloat32(v)
	case float64:
		return mw.WriteFloat64(v)
	case complex64:
		return mw.WriteComplex64(v)
	case complex128:
		return mw.WriteComplex128(v)
	case uint8:
		return mw.WriteUint8(v)
	case uint16:
		return mw.WriteUint16(v)
	case uint32:
		return mw.WriteUint32(v)
	case uint64:
		return mw.WriteUint64(v)
	case uint:
		return mw.WriteUint(v)
	case int8:
		return mw.WriteInt8(v)
	case int16:
		return mw.WriteInt16(v)
	case int32:
		return mw.WriteInt32(v)
	case int64:
		return mw.WriteInt64(v)
	case int:
		return mw.WriteInt(v)
	case string:
		return mw.WriteString(v)
	case []byte:
		return mw.WriteBytes(v)
	case map[string]string:
		return mw.WriteMapStrStr(v)
	case map[string]interface***REMOVED******REMOVED***:
		return mw.WriteMapStrIntf(v)
	case time.Time:
		return mw.WriteTime(v)
	***REMOVED***

	val := reflect.ValueOf(v)
	if !isSupported(val.Kind()) || !val.IsValid() ***REMOVED***
		return fmt.Errorf("msgp: type %s not supported", val)
	***REMOVED***

	switch val.Kind() ***REMOVED***
	case reflect.Ptr:
		if val.IsNil() ***REMOVED***
			return mw.WriteNil()
		***REMOVED***
		return mw.WriteIntf(val.Elem().Interface())
	case reflect.Slice:
		return mw.writeSlice(val)
	case reflect.Map:
		return mw.writeMap(val)
	***REMOVED***
	return &ErrUnsupportedType***REMOVED***val.Type()***REMOVED***
***REMOVED***

func (mw *Writer) writeMap(v reflect.Value) (err error) ***REMOVED***
	if v.Type().Key().Kind() != reflect.String ***REMOVED***
		return errors.New("msgp: map keys must be strings")
	***REMOVED***
	ks := v.MapKeys()
	err = mw.WriteMapHeader(uint32(len(ks)))
	if err != nil ***REMOVED***
		return
	***REMOVED***
	for _, key := range ks ***REMOVED***
		val := v.MapIndex(key)
		err = mw.WriteString(key.String())
		if err != nil ***REMOVED***
			return
		***REMOVED***
		err = mw.WriteIntf(val.Interface())
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (mw *Writer) writeSlice(v reflect.Value) (err error) ***REMOVED***
	// is []byte
	if v.Type().ConvertibleTo(btsType) ***REMOVED***
		return mw.WriteBytes(v.Bytes())
	***REMOVED***

	sz := uint32(v.Len())
	err = mw.WriteArrayHeader(sz)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	for i := uint32(0); i < sz; i++ ***REMOVED***
		err = mw.WriteIntf(v.Index(int(i)).Interface())
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (mw *Writer) writeStruct(v reflect.Value) error ***REMOVED***
	if enc, ok := v.Interface().(Encodable); ok ***REMOVED***
		return enc.EncodeMsg(mw)
	***REMOVED***
	return fmt.Errorf("msgp: unsupported type: %s", v.Type())
***REMOVED***

func (mw *Writer) writeVal(v reflect.Value) error ***REMOVED***
	if !isSupported(v.Kind()) ***REMOVED***
		return fmt.Errorf("msgp: msgp/enc: type %q not supported", v.Type())
	***REMOVED***

	// shortcut for nil values
	if v.IsNil() ***REMOVED***
		return mw.WriteNil()
	***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Bool:
		return mw.WriteBool(v.Bool())

	case reflect.Float32, reflect.Float64:
		return mw.WriteFloat64(v.Float())

	case reflect.Complex64, reflect.Complex128:
		return mw.WriteComplex128(v.Complex())

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return mw.WriteInt64(v.Int())

	case reflect.Interface, reflect.Ptr:
		if v.IsNil() ***REMOVED***
			mw.WriteNil()
		***REMOVED***
		return mw.writeVal(v.Elem())

	case reflect.Map:
		return mw.writeMap(v)

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return mw.WriteUint64(v.Uint())

	case reflect.String:
		return mw.WriteString(v.String())

	case reflect.Slice, reflect.Array:
		return mw.writeSlice(v)

	case reflect.Struct:
		return mw.writeStruct(v)

	***REMOVED***
	return fmt.Errorf("msgp: msgp/enc: type %q not supported", v.Type())
***REMOVED***

// is the reflect.Kind encodable?
func isSupported(k reflect.Kind) bool ***REMOVED***
	switch k ***REMOVED***
	case reflect.Func, reflect.Chan, reflect.Invalid, reflect.UnsafePointer:
		return false
	default:
		return true
	***REMOVED***
***REMOVED***

// GuessSize guesses the size of the underlying
// value of 'i'. If the underlying value is not
// a simple builtin (or []byte), GuessSize defaults
// to 512.
func GuessSize(i interface***REMOVED******REMOVED***) int ***REMOVED***
	if i == nil ***REMOVED***
		return NilSize
	***REMOVED***

	switch i := i.(type) ***REMOVED***
	case Sizer:
		return i.Msgsize()
	case Extension:
		return ExtensionPrefixSize + i.Len()
	case float64:
		return Float64Size
	case float32:
		return Float32Size
	case uint8, uint16, uint32, uint64, uint:
		return UintSize
	case int8, int16, int32, int64, int:
		return IntSize
	case []byte:
		return BytesPrefixSize + len(i)
	case string:
		return StringPrefixSize + len(i)
	case complex64:
		return Complex64Size
	case complex128:
		return Complex128Size
	case bool:
		return BoolSize
	case map[string]interface***REMOVED******REMOVED***:
		s := MapHeaderSize
		for key, val := range i ***REMOVED***
			s += StringPrefixSize + len(key) + GuessSize(val)
		***REMOVED***
		return s
	case map[string]string:
		s := MapHeaderSize
		for key, val := range i ***REMOVED***
			s += 2*StringPrefixSize + len(key) + len(val)
		***REMOVED***
		return s
	default:
		return 512
	***REMOVED***
***REMOVED***
