package ioutils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"

	"golang.org/x/net/context"
)

// ReadCloserWrapper wraps an io.Reader, and implements an io.ReadCloser
// It calls the given callback function when closed. It should be constructed
// with NewReadCloserWrapper
type ReadCloserWrapper struct ***REMOVED***
	io.Reader
	closer func() error
***REMOVED***

// Close calls back the passed closer function
func (r *ReadCloserWrapper) Close() error ***REMOVED***
	return r.closer()
***REMOVED***

// NewReadCloserWrapper returns a new io.ReadCloser.
func NewReadCloserWrapper(r io.Reader, closer func() error) io.ReadCloser ***REMOVED***
	return &ReadCloserWrapper***REMOVED***
		Reader: r,
		closer: closer,
	***REMOVED***
***REMOVED***

type readerErrWrapper struct ***REMOVED***
	reader io.Reader
	closer func()
***REMOVED***

func (r *readerErrWrapper) Read(p []byte) (int, error) ***REMOVED***
	n, err := r.reader.Read(p)
	if err != nil ***REMOVED***
		r.closer()
	***REMOVED***
	return n, err
***REMOVED***

// NewReaderErrWrapper returns a new io.Reader.
func NewReaderErrWrapper(r io.Reader, closer func()) io.Reader ***REMOVED***
	return &readerErrWrapper***REMOVED***
		reader: r,
		closer: closer,
	***REMOVED***
***REMOVED***

// HashData returns the sha256 sum of src.
func HashData(src io.Reader) (string, error) ***REMOVED***
	h := sha256.New()
	if _, err := io.Copy(h, src); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
***REMOVED***

// OnEOFReader wraps an io.ReadCloser and a function
// the function will run at the end of file or close the file.
type OnEOFReader struct ***REMOVED***
	Rc io.ReadCloser
	Fn func()
***REMOVED***

func (r *OnEOFReader) Read(p []byte) (n int, err error) ***REMOVED***
	n, err = r.Rc.Read(p)
	if err == io.EOF ***REMOVED***
		r.runFunc()
	***REMOVED***
	return
***REMOVED***

// Close closes the file and run the function.
func (r *OnEOFReader) Close() error ***REMOVED***
	err := r.Rc.Close()
	r.runFunc()
	return err
***REMOVED***

func (r *OnEOFReader) runFunc() ***REMOVED***
	if fn := r.Fn; fn != nil ***REMOVED***
		fn()
		r.Fn = nil
	***REMOVED***
***REMOVED***

// cancelReadCloser wraps an io.ReadCloser with a context for cancelling read
// operations.
type cancelReadCloser struct ***REMOVED***
	cancel func()
	pR     *io.PipeReader // Stream to read from
	pW     *io.PipeWriter
***REMOVED***

// NewCancelReadCloser creates a wrapper that closes the ReadCloser when the
// context is cancelled. The returned io.ReadCloser must be closed when it is
// no longer needed.
func NewCancelReadCloser(ctx context.Context, in io.ReadCloser) io.ReadCloser ***REMOVED***
	pR, pW := io.Pipe()

	// Create a context used to signal when the pipe is closed
	doneCtx, cancel := context.WithCancel(context.Background())

	p := &cancelReadCloser***REMOVED***
		cancel: cancel,
		pR:     pR,
		pW:     pW,
	***REMOVED***

	go func() ***REMOVED***
		_, err := io.Copy(pW, in)
		select ***REMOVED***
		case <-ctx.Done():
			// If the context was closed, p.closeWithError
			// was already called. Calling it again would
			// change the error that Read returns.
		default:
			p.closeWithError(err)
		***REMOVED***
		in.Close()
	***REMOVED***()
	go func() ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
				p.closeWithError(ctx.Err())
			case <-doneCtx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return p
***REMOVED***

// Read wraps the Read method of the pipe that provides data from the wrapped
// ReadCloser.
func (p *cancelReadCloser) Read(buf []byte) (n int, err error) ***REMOVED***
	return p.pR.Read(buf)
***REMOVED***

// closeWithError closes the wrapper and its underlying reader. It will
// cause future calls to Read to return err.
func (p *cancelReadCloser) closeWithError(err error) ***REMOVED***
	p.pW.CloseWithError(err)
	p.cancel()
***REMOVED***

// Close closes the wrapper its underlying reader. It will cause
// future calls to Read to return io.EOF.
func (p *cancelReadCloser) Close() error ***REMOVED***
	p.closeWithError(io.EOF)
	return nil
***REMOVED***
