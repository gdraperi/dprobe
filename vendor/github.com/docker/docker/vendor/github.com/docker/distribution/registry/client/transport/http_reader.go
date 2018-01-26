package transport

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var (
	contentRangeRegexp = regexp.MustCompile(`bytes ([0-9]+)-([0-9]+)/([0-9]+|\\*)`)

	// ErrWrongCodeForByteRange is returned if the client sends a request
	// with a Range header but the server returns a 2xx or 3xx code other
	// than 206 Partial Content.
	ErrWrongCodeForByteRange = errors.New("expected HTTP 206 from byte range request")
)

// ReadSeekCloser combines io.ReadSeeker with io.Closer.
type ReadSeekCloser interface ***REMOVED***
	io.ReadSeeker
	io.Closer
***REMOVED***

// NewHTTPReadSeeker handles reading from an HTTP endpoint using a GET
// request. When seeking and starting a read from a non-zero offset
// the a "Range" header will be added which sets the offset.
// TODO(dmcgowan): Move this into a separate utility package
func NewHTTPReadSeeker(client *http.Client, url string, errorHandler func(*http.Response) error) ReadSeekCloser ***REMOVED***
	return &httpReadSeeker***REMOVED***
		client:       client,
		url:          url,
		errorHandler: errorHandler,
	***REMOVED***
***REMOVED***

type httpReadSeeker struct ***REMOVED***
	client *http.Client
	url    string

	// errorHandler creates an error from an unsuccessful HTTP response.
	// This allows the error to be created with the HTTP response body
	// without leaking the body through a returned error.
	errorHandler func(*http.Response) error

	size int64

	// rc is the remote read closer.
	rc io.ReadCloser
	// readerOffset tracks the offset as of the last read.
	readerOffset int64
	// seekOffset allows Seek to override the offset. Seek changes
	// seekOffset instead of changing readOffset directly so that
	// connection resets can be delayed and possibly avoided if the
	// seek is undone (i.e. seeking to the end and then back to the
	// beginning).
	seekOffset int64
	err        error
***REMOVED***

func (hrs *httpReadSeeker) Read(p []byte) (n int, err error) ***REMOVED***
	if hrs.err != nil ***REMOVED***
		return 0, hrs.err
	***REMOVED***

	// If we sought to a different position, we need to reset the
	// connection. This logic is here instead of Seek so that if
	// a seek is undone before the next read, the connection doesn't
	// need to be closed and reopened. A common example of this is
	// seeking to the end to determine the length, and then seeking
	// back to the original position.
	if hrs.readerOffset != hrs.seekOffset ***REMOVED***
		hrs.reset()
	***REMOVED***

	hrs.readerOffset = hrs.seekOffset

	rd, err := hrs.reader()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	n, err = rd.Read(p)
	hrs.seekOffset += int64(n)
	hrs.readerOffset += int64(n)

	return n, err
***REMOVED***

func (hrs *httpReadSeeker) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	if hrs.err != nil ***REMOVED***
		return 0, hrs.err
	***REMOVED***

	lastReaderOffset := hrs.readerOffset

	if whence == os.SEEK_SET && hrs.rc == nil ***REMOVED***
		// If no request has been made yet, and we are seeking to an
		// absolute position, set the read offset as well to avoid an
		// unnecessary request.
		hrs.readerOffset = offset
	***REMOVED***

	_, err := hrs.reader()
	if err != nil ***REMOVED***
		hrs.readerOffset = lastReaderOffset
		return 0, err
	***REMOVED***

	newOffset := hrs.seekOffset

	switch whence ***REMOVED***
	case os.SEEK_CUR:
		newOffset += offset
	case os.SEEK_END:
		if hrs.size < 0 ***REMOVED***
			return 0, errors.New("content length not known")
		***REMOVED***
		newOffset = hrs.size + offset
	case os.SEEK_SET:
		newOffset = offset
	***REMOVED***

	if newOffset < 0 ***REMOVED***
		err = errors.New("cannot seek to negative position")
	***REMOVED*** else ***REMOVED***
		hrs.seekOffset = newOffset
	***REMOVED***

	return hrs.seekOffset, err
***REMOVED***

func (hrs *httpReadSeeker) Close() error ***REMOVED***
	if hrs.err != nil ***REMOVED***
		return hrs.err
	***REMOVED***

	// close and release reader chain
	if hrs.rc != nil ***REMOVED***
		hrs.rc.Close()
	***REMOVED***

	hrs.rc = nil

	hrs.err = errors.New("httpLayer: closed")

	return nil
***REMOVED***

func (hrs *httpReadSeeker) reset() ***REMOVED***
	if hrs.err != nil ***REMOVED***
		return
	***REMOVED***
	if hrs.rc != nil ***REMOVED***
		hrs.rc.Close()
		hrs.rc = nil
	***REMOVED***
***REMOVED***

func (hrs *httpReadSeeker) reader() (io.Reader, error) ***REMOVED***
	if hrs.err != nil ***REMOVED***
		return nil, hrs.err
	***REMOVED***

	if hrs.rc != nil ***REMOVED***
		return hrs.rc, nil
	***REMOVED***

	req, err := http.NewRequest("GET", hrs.url, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if hrs.readerOffset > 0 ***REMOVED***
		// If we are at different offset, issue a range request from there.
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-", hrs.readerOffset))
		// TODO: get context in here
		// context.GetLogger(hrs.context).Infof("Range: %s", req.Header.Get("Range"))
	***REMOVED***

	req.Header.Add("Accept-Encoding", "identity")
	resp, err := hrs.client.Do(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Normally would use client.SuccessStatus, but that would be a cyclic
	// import
	if resp.StatusCode >= 200 && resp.StatusCode <= 399 ***REMOVED***
		if hrs.readerOffset > 0 ***REMOVED***
			if resp.StatusCode != http.StatusPartialContent ***REMOVED***
				return nil, ErrWrongCodeForByteRange
			***REMOVED***

			contentRange := resp.Header.Get("Content-Range")
			if contentRange == "" ***REMOVED***
				return nil, errors.New("no Content-Range header found in HTTP 206 response")
			***REMOVED***

			submatches := contentRangeRegexp.FindStringSubmatch(contentRange)
			if len(submatches) < 4 ***REMOVED***
				return nil, fmt.Errorf("could not parse Content-Range header: %s", contentRange)
			***REMOVED***

			startByte, err := strconv.ParseUint(submatches[1], 10, 64)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("could not parse start of range in Content-Range header: %s", contentRange)
			***REMOVED***

			if startByte != uint64(hrs.readerOffset) ***REMOVED***
				return nil, fmt.Errorf("received Content-Range starting at offset %d instead of requested %d", startByte, hrs.readerOffset)
			***REMOVED***

			endByte, err := strconv.ParseUint(submatches[2], 10, 64)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("could not parse end of range in Content-Range header: %s", contentRange)
			***REMOVED***

			if submatches[3] == "*" ***REMOVED***
				hrs.size = -1
			***REMOVED*** else ***REMOVED***
				size, err := strconv.ParseUint(submatches[3], 10, 64)
				if err != nil ***REMOVED***
					return nil, fmt.Errorf("could not parse total size in Content-Range header: %s", contentRange)
				***REMOVED***

				if endByte+1 != size ***REMOVED***
					return nil, fmt.Errorf("range in Content-Range stops before the end of the content: %s", contentRange)
				***REMOVED***

				hrs.size = int64(size)
			***REMOVED***
		***REMOVED*** else if resp.StatusCode == http.StatusOK ***REMOVED***
			hrs.size = resp.ContentLength
		***REMOVED*** else ***REMOVED***
			hrs.size = -1
		***REMOVED***
		hrs.rc = resp.Body
	***REMOVED*** else ***REMOVED***
		defer resp.Body.Close()
		if hrs.errorHandler != nil ***REMOVED***
			return nil, hrs.errorHandler(resp)
		***REMOVED***
		return nil, fmt.Errorf("unexpected status resolving reader: %v", resp.Status)
	***REMOVED***

	return hrs.rc, nil
***REMOVED***
