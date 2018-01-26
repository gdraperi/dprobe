// The `fwd` package provides a buffered reader
// and writer. Each has methods that help improve
// the encoding/decoding performance of some binary
// protocols.
//
// The `fwd.Writer` and `fwd.Reader` type provide similar
// functionality to their counterparts in `bufio`, plus
// a few extra utility methods that simplify read-ahead
// and write-ahead. I wrote this package to improve serialization
// performance for http://github.com/tinylib/msgp,
// where it provided about a 2x speedup over `bufio` for certain
// workloads. However, care must be taken to understand the semantics of the
// extra methods provided by this package, as they allow
// the user to access and manipulate the buffer memory
// directly.
//
// The extra methods for `fwd.Reader` are `Peek`, `Skip`
// and `Next`. `(*fwd.Reader).Peek`, unlike `(*bufio.Reader).Peek`,
// will re-allocate the read buffer in order to accommodate arbitrarily
// large read-ahead. `(*fwd.Reader).Skip` skips the next `n` bytes
// in the stream, and uses the `io.Seeker` interface if the underlying
// stream implements it. `(*fwd.Reader).Next` returns a slice pointing
// to the next `n` bytes in the read buffer (like `Peek`), but also
// increments the read position. This allows users to process streams
// in arbitrary block sizes without having to manage appropriately-sized
// slices. Additionally, obviating the need to copy the data from the
// buffer to another location in memory can improve performance dramatically
// in CPU-bound applications.
//
// `fwd.Writer` only has one extra method, which is `(*fwd.Writer).Next`, which
// returns a slice pointing to the next `n` bytes of the writer, and increments
// the write position by the length of the returned slice. This allows users
// to write directly to the end of the buffer.
//
package fwd

import "io"

const (
	// DefaultReaderSize is the default size of the read buffer
	DefaultReaderSize = 2048

	// minimum read buffer; straight from bufio
	minReaderSize = 16
)

// NewReader returns a new *Reader that reads from 'r'
func NewReader(r io.Reader) *Reader ***REMOVED***
	return NewReaderSize(r, DefaultReaderSize)
***REMOVED***

// NewReaderSize returns a new *Reader that
// reads from 'r' and has a buffer size 'n'
func NewReaderSize(r io.Reader, n int) *Reader ***REMOVED***
	rd := &Reader***REMOVED***
		r:    r,
		data: make([]byte, 0, max(minReaderSize, n)),
	***REMOVED***
	if s, ok := r.(io.Seeker); ok ***REMOVED***
		rd.rs = s
	***REMOVED***
	return rd
***REMOVED***

// Reader is a buffered look-ahead reader
type Reader struct ***REMOVED***
	r io.Reader // underlying reader

	// data[n:len(data)] is buffered data; data[len(data):cap(data)] is free buffer space
	data  []byte // data
	n     int    // read offset
	state error  // last read error

	// if the reader past to NewReader was
	// also an io.Seeker, this is non-nil
	rs io.Seeker
***REMOVED***

// Reset resets the underlying reader
// and the read buffer.
func (r *Reader) Reset(rd io.Reader) ***REMOVED***
	r.r = rd
	r.data = r.data[0:0]
	r.n = 0
	r.state = nil
	if s, ok := rd.(io.Seeker); ok ***REMOVED***
		r.rs = s
	***REMOVED*** else ***REMOVED***
		r.rs = nil
	***REMOVED***
***REMOVED***

// more() does one read on the underlying reader
func (r *Reader) more() ***REMOVED***
	// move data backwards so that
	// the read offset is 0; this way
	// we can supply the maximum number of
	// bytes to the reader
	if r.n != 0 ***REMOVED***
		if r.n < len(r.data) ***REMOVED***
			r.data = r.data[:copy(r.data[0:], r.data[r.n:])]
		***REMOVED*** else ***REMOVED***
			r.data = r.data[:0]
		***REMOVED***
		r.n = 0
	***REMOVED***
	var a int
	a, r.state = r.r.Read(r.data[len(r.data):cap(r.data)])
	if a == 0 && r.state == nil ***REMOVED***
		r.state = io.ErrNoProgress
		return
	***REMOVED***
	r.data = r.data[:len(r.data)+a]
***REMOVED***

// pop error
func (r *Reader) err() (e error) ***REMOVED***
	e, r.state = r.state, nil
	return
***REMOVED***

// pop error; EOF -> io.ErrUnexpectedEOF
func (r *Reader) noEOF() (e error) ***REMOVED***
	e, r.state = r.state, nil
	if e == io.EOF ***REMOVED***
		e = io.ErrUnexpectedEOF
	***REMOVED***
	return
***REMOVED***

// buffered bytes
func (r *Reader) buffered() int ***REMOVED*** return len(r.data) - r.n ***REMOVED***

// Buffered returns the number of bytes currently in the buffer
func (r *Reader) Buffered() int ***REMOVED*** return len(r.data) - r.n ***REMOVED***

// BufferSize returns the total size of the buffer
func (r *Reader) BufferSize() int ***REMOVED*** return cap(r.data) ***REMOVED***

// Peek returns the next 'n' buffered bytes,
// reading from the underlying reader if necessary.
// It will only return a slice shorter than 'n' bytes
// if it also returns an error. Peek does not advance
// the reader. EOF errors are *not* returned as
// io.ErrUnexpectedEOF.
func (r *Reader) Peek(n int) ([]byte, error) ***REMOVED***
	// in the degenerate case,
	// we may need to realloc
	// (the caller asked for more
	// bytes than the size of the buffer)
	if cap(r.data) < n ***REMOVED***
		old := r.data[r.n:]
		r.data = make([]byte, n+r.buffered())
		r.data = r.data[:copy(r.data, old)]
		r.n = 0
	***REMOVED***

	// keep filling until
	// we hit an error or
	// read enough bytes
	for r.buffered() < n && r.state == nil ***REMOVED***
		r.more()
	***REMOVED***

	// we must have hit an error
	if r.buffered() < n ***REMOVED***
		return r.data[r.n:], r.err()
	***REMOVED***

	return r.data[r.n : r.n+n], nil
***REMOVED***

// Skip moves the reader forward 'n' bytes.
// Returns the number of bytes skipped and any
// errors encountered. It is analogous to Seek(n, 1).
// If the underlying reader implements io.Seeker, then
// that method will be used to skip forward.
//
// If the reader encounters
// an EOF before skipping 'n' bytes, it
// returns io.ErrUnexpectedEOF. If the
// underlying reader implements io.Seeker, then
// those rules apply instead. (Many implementations
// will not return `io.EOF` until the next call
// to Read.)
func (r *Reader) Skip(n int) (int, error) ***REMOVED***

	// fast path
	if r.buffered() >= n ***REMOVED***
		r.n += n
		return n, nil
	***REMOVED***

	// use seeker implementation
	// if we can
	if r.rs != nil ***REMOVED***
		return r.skipSeek(n)
	***REMOVED***

	// loop on filling
	// and then erasing
	o := n
	for r.buffered() < n && r.state == nil ***REMOVED***
		r.more()
		// we can skip forward
		// up to r.buffered() bytes
		step := min(r.buffered(), n)
		r.n += step
		n -= step
	***REMOVED***
	// at this point, n should be
	// 0 if everything went smoothly
	return o - n, r.noEOF()
***REMOVED***

// Next returns the next 'n' bytes in the stream.
// Unlike Peek, Next advances the reader position.
// The returned bytes point to the same
// data as the buffer, so the slice is
// only valid until the next reader method call.
// An EOF is considered an unexpected error.
// If an the returned slice is less than the
// length asked for, an error will be returned,
// and the reader position will not be incremented.
func (r *Reader) Next(n int) ([]byte, error) ***REMOVED***

	// in case the buffer is too small
	if cap(r.data) < n ***REMOVED***
		old := r.data[r.n:]
		r.data = make([]byte, n+r.buffered())
		r.data = r.data[:copy(r.data, old)]
		r.n = 0
	***REMOVED***

	// fill at least 'n' bytes
	for r.buffered() < n && r.state == nil ***REMOVED***
		r.more()
	***REMOVED***

	if r.buffered() < n ***REMOVED***
		return r.data[r.n:], r.noEOF()
	***REMOVED***
	out := r.data[r.n : r.n+n]
	r.n += n
	return out, nil
***REMOVED***

// skipSeek uses the io.Seeker to seek forward.
// only call this function when n > r.buffered()
func (r *Reader) skipSeek(n int) (int, error) ***REMOVED***
	o := r.buffered()
	// first, clear buffer
	n -= o
	r.n = 0
	r.data = r.data[:0]

	// then seek forward remaning bytes
	i, err := r.rs.Seek(int64(n), 1)
	return int(i) + o, err
***REMOVED***

// Read implements `io.Reader`
func (r *Reader) Read(b []byte) (int, error) ***REMOVED***
	// if we have data in the buffer, just
	// return that.
	if r.buffered() != 0 ***REMOVED***
		x := copy(b, r.data[r.n:])
		r.n += x
		return x, nil
	***REMOVED***
	var n int
	// we have no buffered data; determine
	// whether or not to buffer or call
	// the underlying reader directly
	if len(b) >= cap(r.data) ***REMOVED***
		n, r.state = r.r.Read(b)
	***REMOVED*** else ***REMOVED***
		r.more()
		n = copy(b, r.data)
		r.n = n
	***REMOVED***
	if n == 0 ***REMOVED***
		return 0, r.err()
	***REMOVED***
	return n, nil
***REMOVED***

// ReadFull attempts to read len(b) bytes into
// 'b'. It returns the number of bytes read into
// 'b', and an error if it does not return len(b).
// EOF is considered an unexpected error.
func (r *Reader) ReadFull(b []byte) (int, error) ***REMOVED***
	var n int  // read into b
	var nn int // scratch
	l := len(b)
	// either read buffered data,
	// or read directly for the underlying
	// buffer, or fetch more buffered data.
	for n < l && r.state == nil ***REMOVED***
		if r.buffered() != 0 ***REMOVED***
			nn = copy(b[n:], r.data[r.n:])
			n += nn
			r.n += nn
		***REMOVED*** else if l-n > cap(r.data) ***REMOVED***
			nn, r.state = r.r.Read(b[n:])
			n += nn
		***REMOVED*** else ***REMOVED***
			r.more()
		***REMOVED***
	***REMOVED***
	if n < l ***REMOVED***
		return n, r.noEOF()
	***REMOVED***
	return n, nil
***REMOVED***

// ReadByte implements `io.ByteReader`
func (r *Reader) ReadByte() (byte, error) ***REMOVED***
	for r.buffered() < 1 && r.state == nil ***REMOVED***
		r.more()
	***REMOVED***
	if r.buffered() < 1 ***REMOVED***
		return 0, r.err()
	***REMOVED***
	b := r.data[r.n]
	r.n++
	return b, nil
***REMOVED***

// WriteTo implements `io.WriterTo`
func (r *Reader) WriteTo(w io.Writer) (int64, error) ***REMOVED***
	var (
		i   int64
		ii  int
		err error
	)
	// first, clear buffer
	if r.buffered() > 0 ***REMOVED***
		ii, err = w.Write(r.data[r.n:])
		i += int64(ii)
		if err != nil ***REMOVED***
			return i, err
		***REMOVED***
		r.data = r.data[0:0]
		r.n = 0
	***REMOVED***
	for r.state == nil ***REMOVED***
		// here we just do
		// 1:1 reads and writes
		r.more()
		if r.buffered() > 0 ***REMOVED***
			ii, err = w.Write(r.data)
			i += int64(ii)
			if err != nil ***REMOVED***
				return i, err
			***REMOVED***
			r.data = r.data[0:0]
			r.n = 0
		***REMOVED***
	***REMOVED***
	if r.state != io.EOF ***REMOVED***
		return i, r.err()
	***REMOVED***
	return i, nil
***REMOVED***

func min(a int, b int) int ***REMOVED***
	if a < b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

func max(a int, b int) int ***REMOVED***
	if a < b ***REMOVED***
		return b
	***REMOVED***
	return a
***REMOVED***
