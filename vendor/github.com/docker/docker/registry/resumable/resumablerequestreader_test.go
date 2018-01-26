package resumable

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResumableRequestHeaderSimpleErrors(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fmt.Fprintln(w, "Hello, world !")
	***REMOVED***))
	defer ts.Close()

	client := &http.Client***REMOVED******REMOVED***

	var req *http.Request
	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	resreq := &requestReader***REMOVED******REMOVED***
	_, err = resreq.Read([]byte***REMOVED******REMOVED***)
	assert.EqualError(t, err, "client and request can't be nil")

	resreq = &requestReader***REMOVED***
		client:    client,
		request:   req,
		totalSize: -1,
	***REMOVED***
	_, err = resreq.Read([]byte***REMOVED******REMOVED***)
	assert.EqualError(t, err, "failed to auto detect content length")
***REMOVED***

// Not too much failures, bails out after some wait
func TestResumableRequestHeaderNotTooMuchFailures(t *testing.T) ***REMOVED***
	client := &http.Client***REMOVED******REMOVED***

	var badReq *http.Request
	badReq, err := http.NewRequest("GET", "I'm not an url", nil)
	require.NoError(t, err)

	resreq := &requestReader***REMOVED***
		client:       client,
		request:      badReq,
		failures:     0,
		maxFailures:  2,
		waitDuration: 10 * time.Millisecond,
	***REMOVED***
	read, err := resreq.Read([]byte***REMOVED******REMOVED***)
	require.NoError(t, err)
	assert.Equal(t, 0, read)
***REMOVED***

// Too much failures, returns the error
func TestResumableRequestHeaderTooMuchFailures(t *testing.T) ***REMOVED***
	client := &http.Client***REMOVED******REMOVED***

	var badReq *http.Request
	badReq, err := http.NewRequest("GET", "I'm not an url", nil)
	require.NoError(t, err)

	resreq := &requestReader***REMOVED***
		client:      client,
		request:     badReq,
		failures:    0,
		maxFailures: 1,
	***REMOVED***
	defer resreq.Close()

	expectedError := `Get I%27m%20not%20an%20url: unsupported protocol scheme ""`
	read, err := resreq.Read([]byte***REMOVED******REMOVED***)
	assert.EqualError(t, err, expectedError)
	assert.Equal(t, 0, read)
***REMOVED***

type errorReaderCloser struct***REMOVED******REMOVED***

func (errorReaderCloser) Close() error ***REMOVED*** return nil ***REMOVED***

func (errorReaderCloser) Read(p []byte) (n int, err error) ***REMOVED***
	return 0, fmt.Errorf("An error occurred")
***REMOVED***

// If an unknown error is encountered, return 0, nil and log it
func TestResumableRequestReaderWithReadError(t *testing.T) ***REMOVED***
	var req *http.Request
	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	client := &http.Client***REMOVED******REMOVED***

	response := &http.Response***REMOVED***
		Status:        "500 Internal Server",
		StatusCode:    500,
		ContentLength: 0,
		Close:         true,
		Body:          errorReaderCloser***REMOVED******REMOVED***,
	***REMOVED***

	resreq := &requestReader***REMOVED***
		client:          client,
		request:         req,
		currentResponse: response,
		lastRange:       1,
		totalSize:       1,
	***REMOVED***
	defer resreq.Close()

	buf := make([]byte, 1)
	read, err := resreq.Read(buf)
	require.NoError(t, err)

	assert.Equal(t, 0, read)
***REMOVED***

func TestResumableRequestReaderWithEOFWith416Response(t *testing.T) ***REMOVED***
	var req *http.Request
	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	client := &http.Client***REMOVED******REMOVED***

	response := &http.Response***REMOVED***
		Status:        "416 Requested Range Not Satisfiable",
		StatusCode:    416,
		ContentLength: 0,
		Close:         true,
		Body:          ioutil.NopCloser(strings.NewReader("")),
	***REMOVED***

	resreq := &requestReader***REMOVED***
		client:          client,
		request:         req,
		currentResponse: response,
		lastRange:       1,
		totalSize:       1,
	***REMOVED***
	defer resreq.Close()

	buf := make([]byte, 1)
	_, err = resreq.Read(buf)
	assert.EqualError(t, err, io.EOF.Error())
***REMOVED***

func TestResumableRequestReaderWithServerDoesntSupportByteRanges(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Header.Get("Range") == "" ***REMOVED***
			t.Fatalf("Expected a Range HTTP header, got nothing")
		***REMOVED***
	***REMOVED***))
	defer ts.Close()

	var req *http.Request
	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	client := &http.Client***REMOVED******REMOVED***

	resreq := &requestReader***REMOVED***
		client:    client,
		request:   req,
		lastRange: 1,
	***REMOVED***
	defer resreq.Close()

	buf := make([]byte, 2)
	_, err = resreq.Read(buf)
	assert.EqualError(t, err, "the server doesn't support byte ranges")
***REMOVED***

func TestResumableRequestReaderWithZeroTotalSize(t *testing.T) ***REMOVED***
	srvtxt := "some response text data"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fmt.Fprintln(w, srvtxt)
	***REMOVED***))
	defer ts.Close()

	var req *http.Request
	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	client := &http.Client***REMOVED******REMOVED***
	retries := uint32(5)

	resreq := NewRequestReader(client, req, retries, 0)
	defer resreq.Close()

	data, err := ioutil.ReadAll(resreq)
	require.NoError(t, err)

	resstr := strings.TrimSuffix(string(data), "\n")
	assert.Equal(t, srvtxt, resstr)
***REMOVED***

func TestResumableRequestReader(t *testing.T) ***REMOVED***
	srvtxt := "some response text data"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fmt.Fprintln(w, srvtxt)
	***REMOVED***))
	defer ts.Close()

	var req *http.Request
	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	client := &http.Client***REMOVED******REMOVED***
	retries := uint32(5)
	imgSize := int64(len(srvtxt))

	resreq := NewRequestReader(client, req, retries, imgSize)
	defer resreq.Close()

	data, err := ioutil.ReadAll(resreq)
	require.NoError(t, err)

	resstr := strings.TrimSuffix(string(data), "\n")
	assert.Equal(t, srvtxt, resstr)
***REMOVED***

func TestResumableRequestReaderWithInitialResponse(t *testing.T) ***REMOVED***
	srvtxt := "some response text data"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fmt.Fprintln(w, srvtxt)
	***REMOVED***))
	defer ts.Close()

	var req *http.Request
	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	client := &http.Client***REMOVED******REMOVED***
	retries := uint32(5)
	imgSize := int64(len(srvtxt))

	res, err := client.Do(req)
	require.NoError(t, err)

	resreq := NewRequestReaderWithInitialResponse(client, req, retries, imgSize, res)
	defer resreq.Close()

	data, err := ioutil.ReadAll(resreq)
	require.NoError(t, err)

	resstr := strings.TrimSuffix(string(data), "\n")
	assert.Equal(t, srvtxt, resstr)
***REMOVED***
