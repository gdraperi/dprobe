package ioutils

import "io"

// NopWriter represents a type which write operation is nop.
type NopWriter struct***REMOVED******REMOVED***

func (*NopWriter) Write(buf []byte) (int, error) ***REMOVED***
	return len(buf), nil
***REMOVED***

type nopWriteCloser struct ***REMOVED***
	io.Writer
***REMOVED***

func (w *nopWriteCloser) Close() error ***REMOVED*** return nil ***REMOVED***

// NopWriteCloser returns a nopWriteCloser.
func NopWriteCloser(w io.Writer) io.WriteCloser ***REMOVED***
	return &nopWriteCloser***REMOVED***w***REMOVED***
***REMOVED***

// NopFlusher represents a type which flush operation is nop.
type NopFlusher struct***REMOVED******REMOVED***

// Flush is a nop operation.
func (f *NopFlusher) Flush() ***REMOVED******REMOVED***

type writeCloserWrapper struct ***REMOVED***
	io.Writer
	closer func() error
***REMOVED***

func (r *writeCloserWrapper) Close() error ***REMOVED***
	return r.closer()
***REMOVED***

// NewWriteCloserWrapper returns a new io.WriteCloser.
func NewWriteCloserWrapper(r io.Writer, closer func() error) io.WriteCloser ***REMOVED***
	return &writeCloserWrapper***REMOVED***
		Writer: r,
		closer: closer,
	***REMOVED***
***REMOVED***

// WriteCounter wraps a concrete io.Writer and hold a count of the number
// of bytes written to the writer during a "session".
// This can be convenient when write return is masked
// (e.g., json.Encoder.Encode())
type WriteCounter struct ***REMOVED***
	Count  int64
	Writer io.Writer
***REMOVED***

// NewWriteCounter returns a new WriteCounter.
func NewWriteCounter(w io.Writer) *WriteCounter ***REMOVED***
	return &WriteCounter***REMOVED***
		Writer: w,
	***REMOVED***
***REMOVED***

func (wc *WriteCounter) Write(p []byte) (count int, err error) ***REMOVED***
	count, err = wc.Writer.Write(p)
	wc.Count += int64(count)
	return
***REMOVED***
