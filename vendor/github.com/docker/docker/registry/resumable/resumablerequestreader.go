package resumable

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type requestReader struct ***REMOVED***
	client          *http.Client
	request         *http.Request
	lastRange       int64
	totalSize       int64
	currentResponse *http.Response
	failures        uint32
	maxFailures     uint32
	waitDuration    time.Duration
***REMOVED***

// NewRequestReader makes it possible to resume reading a request's body transparently
// maxfail is the number of times we retry to make requests again (not resumes)
// totalsize is the total length of the body; auto detect if not provided
func NewRequestReader(c *http.Client, r *http.Request, maxfail uint32, totalsize int64) io.ReadCloser ***REMOVED***
	return &requestReader***REMOVED***client: c, request: r, maxFailures: maxfail, totalSize: totalsize, waitDuration: 5 * time.Second***REMOVED***
***REMOVED***

// NewRequestReaderWithInitialResponse makes it possible to resume
// reading the body of an already initiated request.
func NewRequestReaderWithInitialResponse(c *http.Client, r *http.Request, maxfail uint32, totalsize int64, initialResponse *http.Response) io.ReadCloser ***REMOVED***
	return &requestReader***REMOVED***client: c, request: r, maxFailures: maxfail, totalSize: totalsize, currentResponse: initialResponse, waitDuration: 5 * time.Second***REMOVED***
***REMOVED***

func (r *requestReader) Read(p []byte) (n int, err error) ***REMOVED***
	if r.client == nil || r.request == nil ***REMOVED***
		return 0, fmt.Errorf("client and request can't be nil")
	***REMOVED***
	isFreshRequest := false
	if r.lastRange != 0 && r.currentResponse == nil ***REMOVED***
		readRange := fmt.Sprintf("bytes=%d-%d", r.lastRange, r.totalSize)
		r.request.Header.Set("Range", readRange)
		time.Sleep(r.waitDuration)
	***REMOVED***
	if r.currentResponse == nil ***REMOVED***
		r.currentResponse, err = r.client.Do(r.request)
		isFreshRequest = true
	***REMOVED***
	if err != nil && r.failures+1 != r.maxFailures ***REMOVED***
		r.cleanUpResponse()
		r.failures++
		time.Sleep(time.Duration(r.failures) * r.waitDuration)
		return 0, nil
	***REMOVED*** else if err != nil ***REMOVED***
		r.cleanUpResponse()
		return 0, err
	***REMOVED***
	if r.currentResponse.StatusCode == 416 && r.lastRange == r.totalSize && r.currentResponse.ContentLength == 0 ***REMOVED***
		r.cleanUpResponse()
		return 0, io.EOF
	***REMOVED*** else if r.currentResponse.StatusCode != 206 && r.lastRange != 0 && isFreshRequest ***REMOVED***
		r.cleanUpResponse()
		return 0, fmt.Errorf("the server doesn't support byte ranges")
	***REMOVED***
	if r.totalSize == 0 ***REMOVED***
		r.totalSize = r.currentResponse.ContentLength
	***REMOVED*** else if r.totalSize <= 0 ***REMOVED***
		r.cleanUpResponse()
		return 0, fmt.Errorf("failed to auto detect content length")
	***REMOVED***
	n, err = r.currentResponse.Body.Read(p)
	r.lastRange += int64(n)
	if err != nil ***REMOVED***
		r.cleanUpResponse()
	***REMOVED***
	if err != nil && err != io.EOF ***REMOVED***
		logrus.Infof("encountered error during pull and clearing it before resume: %s", err)
		err = nil
	***REMOVED***
	return n, err
***REMOVED***

func (r *requestReader) Close() error ***REMOVED***
	r.cleanUpResponse()
	r.client = nil
	r.request = nil
	return nil
***REMOVED***

func (r *requestReader) cleanUpResponse() ***REMOVED***
	if r.currentResponse != nil ***REMOVED***
		r.currentResponse.Body.Close()
		r.currentResponse = nil
	***REMOVED***
***REMOVED***
