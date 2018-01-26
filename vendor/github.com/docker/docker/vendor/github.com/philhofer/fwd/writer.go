package fwd

import "io"

const (
	// DefaultWriterSize is the
	// default write buffer size.
	DefaultWriterSize = 2048

	minWriterSize = minReaderSize
)

// Writer is a buffered writer
type Writer struct ***REMOVED***
	w   io.Writer // writer
	buf []byte    // 0:len(buf) is bufered data
***REMOVED***

// NewWriter returns a new writer
// that writes to 'w' and has a buffer
// that is `DefaultWriterSize` bytes.
func NewWriter(w io.Writer) *Writer ***REMOVED***
	if wr, ok := w.(*Writer); ok ***REMOVED***
		return wr
	***REMOVED***
	return &Writer***REMOVED***
		w:   w,
		buf: make([]byte, 0, DefaultWriterSize),
	***REMOVED***
***REMOVED***

// NewWriterSize returns a new writer
// that writes to 'w' and has a buffer
// that is 'size' bytes.
func NewWriterSize(w io.Writer, size int) *Writer ***REMOVED***
	if wr, ok := w.(*Writer); ok && cap(wr.buf) >= size ***REMOVED***
		return wr
	***REMOVED***
	return &Writer***REMOVED***
		w:   w,
		buf: make([]byte, 0, max(size, minWriterSize)),
	***REMOVED***
***REMOVED***

// Buffered returns the number of buffered bytes
// in the reader.
func (w *Writer) Buffered() int ***REMOVED*** return len(w.buf) ***REMOVED***

// BufferSize returns the maximum size of the buffer.
func (w *Writer) BufferSize() int ***REMOVED*** return cap(w.buf) ***REMOVED***

// Flush flushes any buffered bytes
// to the underlying writer.
func (w *Writer) Flush() error ***REMOVED***
	l := len(w.buf)
	if l > 0 ***REMOVED***
		n, err := w.w.Write(w.buf)

		// if we didn't write the whole
		// thing, copy the unwritten
		// bytes to the beginnning of the
		// buffer.
		if n < l && n > 0 ***REMOVED***
			w.pushback(n)
			if err == nil ***REMOVED***
				err = io.ErrShortWrite
			***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		w.buf = w.buf[:0]
		return nil
	***REMOVED***
	return nil
***REMOVED***

// Write implements `io.Writer`
func (w *Writer) Write(p []byte) (int, error) ***REMOVED***
	c, l, ln := cap(w.buf), len(w.buf), len(p)
	avail := c - l

	// requires flush
	if avail < ln ***REMOVED***
		if err := w.Flush(); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		l = len(w.buf)
	***REMOVED***
	// too big to fit in buffer;
	// write directly to w.w
	if c < ln ***REMOVED***
		return w.w.Write(p)
	***REMOVED***

	// grow buf slice; copy; return
	w.buf = w.buf[:l+ln]
	return copy(w.buf[l:], p), nil
***REMOVED***

// WriteString is analogous to Write, but it takes a string.
func (w *Writer) WriteString(s string) (int, error) ***REMOVED***
	c, l, ln := cap(w.buf), len(w.buf), len(s)
	avail := c - l

	// requires flush
	if avail < ln ***REMOVED***
		if err := w.Flush(); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		l = len(w.buf)
	***REMOVED***
	// too big to fit in buffer;
	// write directly to w.w
	//
	// yes, this is unsafe. *but*
	// io.Writer is not allowed
	// to mutate its input or
	// maintain a reference to it,
	// per the spec in package io.
	//
	// plus, if the string is really
	// too big to fit in the buffer, then
	// creating a copy to write it is
	// expensive (and, strictly speaking,
	// unnecessary)
	if c < ln ***REMOVED***
		return w.w.Write(unsafestr(s))
	***REMOVED***

	// grow buf slice; copy; return
	w.buf = w.buf[:l+ln]
	return copy(w.buf[l:], s), nil
***REMOVED***

// WriteByte implements `io.ByteWriter`
func (w *Writer) WriteByte(b byte) error ***REMOVED***
	if len(w.buf) == cap(w.buf) ***REMOVED***
		if err := w.Flush(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	w.buf = append(w.buf, b)
	return nil
***REMOVED***

// Next returns the next 'n' free bytes
// in the write buffer, flushing the writer
// as necessary. Next will return `io.ErrShortBuffer`
// if 'n' is greater than the size of the write buffer.
// Calls to 'next' increment the write position by
// the size of the returned buffer.
func (w *Writer) Next(n int) ([]byte, error) ***REMOVED***
	c, l := cap(w.buf), len(w.buf)
	if n > c ***REMOVED***
		return nil, io.ErrShortBuffer
	***REMOVED***
	avail := c - l
	if avail < n ***REMOVED***
		if err := w.Flush(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		l = len(w.buf)
	***REMOVED***
	w.buf = w.buf[:l+n]
	return w.buf[l:], nil
***REMOVED***

// take the bytes from w.buf[n:len(w.buf)]
// and put them at the beginning of w.buf,
// and resize to the length of the copied segment.
func (w *Writer) pushback(n int) ***REMOVED***
	w.buf = w.buf[:copy(w.buf, w.buf[n:])]
***REMOVED***

// ReadFrom implements `io.ReaderFrom`
func (w *Writer) ReadFrom(r io.Reader) (int64, error) ***REMOVED***
	// anticipatory flush
	if err := w.Flush(); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	w.buf = w.buf[0:cap(w.buf)] // expand buffer

	var nn int64  // written
	var err error // error
	var x int     // read

	// 1:1 reads and writes
	for err == nil ***REMOVED***
		x, err = r.Read(w.buf)
		if x > 0 ***REMOVED***
			n, werr := w.w.Write(w.buf[:x])
			nn += int64(n)

			if err != nil ***REMOVED***
				if n < x && n > 0 ***REMOVED***
					w.pushback(n - x)
				***REMOVED***
				return nn, werr
			***REMOVED***
			if n < x ***REMOVED***
				w.pushback(n - x)
				return nn, io.ErrShortWrite
			***REMOVED***
		***REMOVED*** else if err == nil ***REMOVED***
			err = io.ErrNoProgress
			break
		***REMOVED***
	***REMOVED***
	if err != io.EOF ***REMOVED***
		return nn, err
	***REMOVED***

	// we only clear here
	// because we are sure
	// the writes have
	// succeeded. otherwise,
	// we retain the data in case
	// future writes succeed.
	w.buf = w.buf[0:0]

	return nn, nil
***REMOVED***
