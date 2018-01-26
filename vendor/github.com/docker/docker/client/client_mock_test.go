package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types"
)

// transportFunc allows us to inject a mock transport for testing. We define it
// here so we can detect the tlsconfig and return nil for only this type.
type transportFunc func(*http.Request) (*http.Response, error)

func (tf transportFunc) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	return tf(req)
***REMOVED***

func newMockClient(doer func(*http.Request) (*http.Response, error)) *http.Client ***REMOVED***
	return &http.Client***REMOVED***
		Transport: transportFunc(doer),
	***REMOVED***
***REMOVED***

func errorMock(statusCode int, message string) func(req *http.Request) (*http.Response, error) ***REMOVED***
	return func(req *http.Request) (*http.Response, error) ***REMOVED***
		header := http.Header***REMOVED******REMOVED***
		header.Set("Content-Type", "application/json")

		body, err := json.Marshal(&types.ErrorResponse***REMOVED***
			Message: message,
		***REMOVED***)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return &http.Response***REMOVED***
			StatusCode: statusCode,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
			Header:     header,
		***REMOVED***, nil
	***REMOVED***
***REMOVED***

func plainTextErrorMock(statusCode int, message string) func(req *http.Request) (*http.Response, error) ***REMOVED***
	return func(req *http.Request) (*http.Response, error) ***REMOVED***
		return &http.Response***REMOVED***
			StatusCode: statusCode,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(message))),
		***REMOVED***, nil
	***REMOVED***
***REMOVED***
