package aws

import (
	"io"
	"sync"
)

// ReadSeekCloser wraps a io.Reader returning a ReaderSeekerCloser. Should
// only be used with an io.Reader that is also an io.Seeker. Doing so may
// cause request signature errors, or request body's not sent for GET, HEAD
// and DELETE HTTP methods.
//
// Deprecated: Should only be used with io.ReadSeeker. If using for
// S3 PutObject to stream content use s3manager.Uploader instead.
func ReadSeekCloser(r io.Reader) ReaderSeekerCloser ***REMOVED***
	return ReaderSeekerCloser***REMOVED***r***REMOVED***
***REMOVED***

// ReaderSeekerCloser represents a reader that can also delegate io.Seeker and
// io.Closer interfaces to the underlying object if they are available.
type ReaderSeekerCloser struct ***REMOVED***
	r io.Reader
***REMOVED***

// Read reads from the reader up to size of p. The number of bytes read, and
// error if it occurred will be returned.
//
// If the reader is not an io.Reader zero bytes read, and nil error will be returned.
//
// Performs the same functionality as io.Reader Read
func (r ReaderSeekerCloser) Read(p []byte) (int, error) ***REMOVED***
	switch t := r.r.(type) ***REMOVED***
	case io.Reader:
		return t.Read(p)
	***REMOVED***
	return 0, nil
***REMOVED***

// Seek sets the offset for the next Read to offset, interpreted according to
// whence: 0 means relative to the origin of the file, 1 means relative to the
// current offset, and 2 means relative to the end. Seek returns the new offset
// and an error, if any.
//
// If the ReaderSeekerCloser is not an io.Seeker nothing will be done.
func (r ReaderSeekerCloser) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	switch t := r.r.(type) ***REMOVED***
	case io.Seeker:
		return t.Seek(offset, whence)
	***REMOVED***
	return int64(0), nil
***REMOVED***

// IsSeeker returns if the underlying reader is also a seeker.
func (r ReaderSeekerCloser) IsSeeker() bool ***REMOVED***
	_, ok := r.r.(io.Seeker)
	return ok
***REMOVED***

// Close closes the ReaderSeekerCloser.
//
// If the ReaderSeekerCloser is not an io.Closer nothing will be done.
func (r ReaderSeekerCloser) Close() error ***REMOVED***
	switch t := r.r.(type) ***REMOVED***
	case io.Closer:
		return t.Close()
	***REMOVED***
	return nil
***REMOVED***

// A WriteAtBuffer provides a in memory buffer supporting the io.WriterAt interface
// Can be used with the s3manager.Downloader to download content to a buffer
// in memory. Safe to use concurrently.
type WriteAtBuffer struct ***REMOVED***
	buf []byte
	m   sync.Mutex

	// GrowthCoeff defines the growth rate of the internal buffer. By
	// default, the growth rate is 1, where expanding the internal
	// buffer will allocate only enough capacity to fit the new expected
	// length.
	GrowthCoeff float64
***REMOVED***

// NewWriteAtBuffer creates a WriteAtBuffer with an internal buffer
// provided by buf.
func NewWriteAtBuffer(buf []byte) *WriteAtBuffer ***REMOVED***
	return &WriteAtBuffer***REMOVED***buf: buf***REMOVED***
***REMOVED***

// WriteAt writes a slice of bytes to a buffer starting at the position provided
// The number of bytes written will be returned, or error. Can overwrite previous
// written slices if the write ats overlap.
func (b *WriteAtBuffer) WriteAt(p []byte, pos int64) (n int, err error) ***REMOVED***
	pLen := len(p)
	expLen := pos + int64(pLen)
	b.m.Lock()
	defer b.m.Unlock()
	if int64(len(b.buf)) < expLen ***REMOVED***
		if int64(cap(b.buf)) < expLen ***REMOVED***
			if b.GrowthCoeff < 1 ***REMOVED***
				b.GrowthCoeff = 1
			***REMOVED***
			newBuf := make([]byte, expLen, int64(b.GrowthCoeff*float64(expLen)))
			copy(newBuf, b.buf)
			b.buf = newBuf
		***REMOVED***
		b.buf = b.buf[:expLen]
	***REMOVED***
	copy(b.buf[pos:], p)
	return pLen, nil
***REMOVED***

// Bytes returns a slice of bytes written to the buffer.
func (b *WriteAtBuffer) Bytes() []byte ***REMOVED***
	b.m.Lock()
	defer b.m.Unlock()
	return b.buf
***REMOVED***
