package request

import (
	"io"
	"sync"
)

// offsetReader is a thread-safe io.ReadCloser to prevent racing
// with retrying requests
type offsetReader struct ***REMOVED***
	buf    io.ReadSeeker
	lock   sync.Mutex
	closed bool
***REMOVED***

func newOffsetReader(buf io.ReadSeeker, offset int64) *offsetReader ***REMOVED***
	reader := &offsetReader***REMOVED******REMOVED***
	buf.Seek(offset, 0)

	reader.buf = buf
	return reader
***REMOVED***

// Close will close the instance of the offset reader's access to
// the underlying io.ReadSeeker.
func (o *offsetReader) Close() error ***REMOVED***
	o.lock.Lock()
	defer o.lock.Unlock()
	o.closed = true
	return nil
***REMOVED***

// Read is a thread-safe read of the underlying io.ReadSeeker
func (o *offsetReader) Read(p []byte) (int, error) ***REMOVED***
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.closed ***REMOVED***
		return 0, io.EOF
	***REMOVED***

	return o.buf.Read(p)
***REMOVED***

// Seek is a thread-safe seeking operation.
func (o *offsetReader) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	o.lock.Lock()
	defer o.lock.Unlock()

	return o.buf.Seek(offset, whence)
***REMOVED***

// CloseAndCopy will return a new offsetReader with a copy of the old buffer
// and close the old buffer.
func (o *offsetReader) CloseAndCopy(offset int64) *offsetReader ***REMOVED***
	o.Close()
	return newOffsetReader(o.buf, offset)
***REMOVED***
