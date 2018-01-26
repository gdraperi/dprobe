// Package pools provides a collection of pools which provide various
// data types with buffers. These can be used to lower the number of
// memory allocations and reuse buffers.
//
// New pools should be added to this package to allow them to be
// shared across packages.
//
// Utility functions which operate on pools should be added to this
// package to allow them to be reused.
package pools

import (
	"bufio"
	"io"
	"sync"

	"github.com/docker/docker/pkg/ioutils"
)

const buffer32K = 32 * 1024

var (
	// BufioReader32KPool is a pool which returns bufio.Reader with a 32K buffer.
	BufioReader32KPool = newBufioReaderPoolWithSize(buffer32K)
	// BufioWriter32KPool is a pool which returns bufio.Writer with a 32K buffer.
	BufioWriter32KPool = newBufioWriterPoolWithSize(buffer32K)
	buffer32KPool      = newBufferPoolWithSize(buffer32K)
)

// BufioReaderPool is a bufio reader that uses sync.Pool.
type BufioReaderPool struct ***REMOVED***
	pool sync.Pool
***REMOVED***

// newBufioReaderPoolWithSize is unexported because new pools should be
// added here to be shared where required.
func newBufioReaderPoolWithSize(size int) *BufioReaderPool ***REMOVED***
	return &BufioReaderPool***REMOVED***
		pool: sync.Pool***REMOVED***
			New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return bufio.NewReaderSize(nil, size) ***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Get returns a bufio.Reader which reads from r. The buffer size is that of the pool.
func (bufPool *BufioReaderPool) Get(r io.Reader) *bufio.Reader ***REMOVED***
	buf := bufPool.pool.Get().(*bufio.Reader)
	buf.Reset(r)
	return buf
***REMOVED***

// Put puts the bufio.Reader back into the pool.
func (bufPool *BufioReaderPool) Put(b *bufio.Reader) ***REMOVED***
	b.Reset(nil)
	bufPool.pool.Put(b)
***REMOVED***

type bufferPool struct ***REMOVED***
	pool sync.Pool
***REMOVED***

func newBufferPoolWithSize(size int) *bufferPool ***REMOVED***
	return &bufferPool***REMOVED***
		pool: sync.Pool***REMOVED***
			New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return make([]byte, size) ***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (bp *bufferPool) Get() []byte ***REMOVED***
	return bp.pool.Get().([]byte)
***REMOVED***

func (bp *bufferPool) Put(b []byte) ***REMOVED***
	bp.pool.Put(b)
***REMOVED***

// Copy is a convenience wrapper which uses a buffer to avoid allocation in io.Copy.
func Copy(dst io.Writer, src io.Reader) (written int64, err error) ***REMOVED***
	buf := buffer32KPool.Get()
	written, err = io.CopyBuffer(dst, src, buf)
	buffer32KPool.Put(buf)
	return
***REMOVED***

// NewReadCloserWrapper returns a wrapper which puts the bufio.Reader back
// into the pool and closes the reader if it's an io.ReadCloser.
func (bufPool *BufioReaderPool) NewReadCloserWrapper(buf *bufio.Reader, r io.Reader) io.ReadCloser ***REMOVED***
	return ioutils.NewReadCloserWrapper(r, func() error ***REMOVED***
		if readCloser, ok := r.(io.ReadCloser); ok ***REMOVED***
			readCloser.Close()
		***REMOVED***
		bufPool.Put(buf)
		return nil
	***REMOVED***)
***REMOVED***

// BufioWriterPool is a bufio writer that uses sync.Pool.
type BufioWriterPool struct ***REMOVED***
	pool sync.Pool
***REMOVED***

// newBufioWriterPoolWithSize is unexported because new pools should be
// added here to be shared where required.
func newBufioWriterPoolWithSize(size int) *BufioWriterPool ***REMOVED***
	return &BufioWriterPool***REMOVED***
		pool: sync.Pool***REMOVED***
			New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return bufio.NewWriterSize(nil, size) ***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Get returns a bufio.Writer which writes to w. The buffer size is that of the pool.
func (bufPool *BufioWriterPool) Get(w io.Writer) *bufio.Writer ***REMOVED***
	buf := bufPool.pool.Get().(*bufio.Writer)
	buf.Reset(w)
	return buf
***REMOVED***

// Put puts the bufio.Writer back into the pool.
func (bufPool *BufioWriterPool) Put(b *bufio.Writer) ***REMOVED***
	b.Reset(nil)
	bufPool.pool.Put(b)
***REMOVED***

// NewWriteCloserWrapper returns a wrapper which puts the bufio.Writer back
// into the pool and closes the writer if it's an io.Writecloser.
func (bufPool *BufioWriterPool) NewWriteCloserWrapper(buf *bufio.Writer, w io.Writer) io.WriteCloser ***REMOVED***
	return ioutils.NewWriteCloserWrapper(w, func() error ***REMOVED***
		buf.Flush()
		if writeCloser, ok := w.(io.WriteCloser); ok ***REMOVED***
			writeCloser.Close()
		***REMOVED***
		bufPool.Put(buf)
		return nil
	***REMOVED***)
***REMOVED***
