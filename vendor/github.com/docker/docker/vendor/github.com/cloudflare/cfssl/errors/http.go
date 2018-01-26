package errors

import (
	"errors"
	"net/http"
)

// HTTPError is an augmented error with a HTTP status code.
type HTTPError struct ***REMOVED***
	StatusCode int
	error
***REMOVED***

// Error implements the error interface.
func (e *HTTPError) Error() string ***REMOVED***
	return e.error.Error()
***REMOVED***

// NewMethodNotAllowed returns an appropriate error in the case that
// an HTTP client uses an invalid method (i.e. a GET in place of a POST)
// on an API endpoint.
func NewMethodNotAllowed(method string) *HTTPError ***REMOVED***
	return &HTTPError***REMOVED***http.StatusMethodNotAllowed, errors.New(`Method is not allowed:"` + method + `"`)***REMOVED***
***REMOVED***

// NewBadRequest creates a HttpError with the given error and error code 400.
func NewBadRequest(err error) *HTTPError ***REMOVED***
	return &HTTPError***REMOVED***http.StatusBadRequest, err***REMOVED***
***REMOVED***

// NewBadRequestString returns a HttpError with the supplied message
// and error code 400.
func NewBadRequestString(s string) *HTTPError ***REMOVED***
	return NewBadRequest(errors.New(s))
***REMOVED***

// NewBadRequestMissingParameter returns a 400 HttpError as a required
// parameter is missing in the HTTP request.
func NewBadRequestMissingParameter(s string) *HTTPError ***REMOVED***
	return NewBadRequestString(`Missing parameter "` + s + `"`)
***REMOVED***

// NewBadRequestUnwantedParameter returns a 400 HttpError as a unnecessary
// parameter is present in the HTTP request.
func NewBadRequestUnwantedParameter(s string) *HTTPError ***REMOVED***
	return NewBadRequestString(`Unwanted parameter "` + s + `"`)
***REMOVED***
