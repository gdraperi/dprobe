package docker

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/log"
	"github.com/pkg/errors"
)

type httpReadSeeker struct ***REMOVED***
	size   int64
	offset int64
	rc     io.ReadCloser
	open   func(offset int64) (io.ReadCloser, error)
	closed bool
***REMOVED***

func newHTTPReadSeeker(size int64, open func(offset int64) (io.ReadCloser, error)) (io.ReadCloser, error) ***REMOVED***
	return &httpReadSeeker***REMOVED***
		size: size,
		open: open,
	***REMOVED***, nil
***REMOVED***

func (hrs *httpReadSeeker) Read(p []byte) (n int, err error) ***REMOVED***
	if hrs.closed ***REMOVED***
		return 0, io.EOF
	***REMOVED***

	rd, err := hrs.reader()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	n, err = rd.Read(p)
	hrs.offset += int64(n)
	return
***REMOVED***

func (hrs *httpReadSeeker) Close() error ***REMOVED***
	if hrs.closed ***REMOVED***
		return nil
	***REMOVED***
	hrs.closed = true
	if hrs.rc != nil ***REMOVED***
		return hrs.rc.Close()
	***REMOVED***

	return nil
***REMOVED***

func (hrs *httpReadSeeker) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	if hrs.closed ***REMOVED***
		return 0, errors.Wrap(errdefs.ErrUnavailable, "Fetcher.Seek: closed")
	***REMOVED***

	abs := hrs.offset
	switch whence ***REMOVED***
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs += offset
	case io.SeekEnd:
		if hrs.size == -1 ***REMOVED***
			return 0, errors.Wrap(errdefs.ErrUnavailable, "Fetcher.Seek: unknown size, cannot seek from end")
		***REMOVED***
		abs = hrs.size + offset
	default:
		return 0, errors.Wrap(errdefs.ErrInvalidArgument, "Fetcher.Seek: invalid whence")
	***REMOVED***

	if abs < 0 ***REMOVED***
		return 0, errors.Wrapf(errdefs.ErrInvalidArgument, "Fetcher.Seek: negative offset")
	***REMOVED***

	if abs != hrs.offset ***REMOVED***
		if hrs.rc != nil ***REMOVED***
			if err := hrs.rc.Close(); err != nil ***REMOVED***
				log.L.WithError(err).Errorf("Fetcher.Seek: failed to close ReadCloser")
			***REMOVED***

			hrs.rc = nil
		***REMOVED***

		hrs.offset = abs
	***REMOVED***

	return hrs.offset, nil
***REMOVED***

func (hrs *httpReadSeeker) reader() (io.Reader, error) ***REMOVED***
	if hrs.rc != nil ***REMOVED***
		return hrs.rc, nil
	***REMOVED***

	if hrs.size == -1 || hrs.offset < hrs.size ***REMOVED***
		// only try to reopen the body request if we are seeking to a value
		// less than the actual size.
		if hrs.open == nil ***REMOVED***
			return nil, errors.Wrapf(errdefs.ErrNotImplemented, "cannot open")
		***REMOVED***

		rc, err := hrs.open(hrs.offset)
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "httpReaderSeeker: failed open")
		***REMOVED***

		if hrs.rc != nil ***REMOVED***
			if err := hrs.rc.Close(); err != nil ***REMOVED***
				log.L.WithError(err).Errorf("httpReadSeeker: failed to close ReadCloser")
			***REMOVED***
		***REMOVED***
		hrs.rc = rc
	***REMOVED*** else ***REMOVED***
		// There is an edge case here where offset == size of the content. If
		// we seek, we will probably get an error for content that cannot be
		// sought (?). In that case, we should err on committing the content,
		// as the length is already satisified but we just return the empty
		// reader instead.

		hrs.rc = ioutil.NopCloser(bytes.NewReader([]byte***REMOVED******REMOVED***))
	***REMOVED***

	return hrs.rc, nil
***REMOVED***
