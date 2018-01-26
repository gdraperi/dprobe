package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/docker/distribution/registry/api/errcode"
	"github.com/docker/distribution/registry/client/auth/challenge"
)

// ErrNoErrorsInBody is returned when an HTTP response body parses to an empty
// errcode.Errors slice.
var ErrNoErrorsInBody = errors.New("no error details found in HTTP response body")

// UnexpectedHTTPStatusError is returned when an unexpected HTTP status is
// returned when making a registry api call.
type UnexpectedHTTPStatusError struct ***REMOVED***
	Status string
***REMOVED***

func (e *UnexpectedHTTPStatusError) Error() string ***REMOVED***
	return fmt.Sprintf("received unexpected HTTP status: %s", e.Status)
***REMOVED***

// UnexpectedHTTPResponseError is returned when an expected HTTP status code
// is returned, but the content was unexpected and failed to be parsed.
type UnexpectedHTTPResponseError struct ***REMOVED***
	ParseErr   error
	StatusCode int
	Response   []byte
***REMOVED***

func (e *UnexpectedHTTPResponseError) Error() string ***REMOVED***
	return fmt.Sprintf("error parsing HTTP %d response body: %s: %q", e.StatusCode, e.ParseErr.Error(), string(e.Response))
***REMOVED***

func parseHTTPErrorResponse(statusCode int, r io.Reader) error ***REMOVED***
	var errors errcode.Errors
	body, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// For backward compatibility, handle irregularly formatted
	// messages that contain a "details" field.
	var detailsErr struct ***REMOVED***
		Details string `json:"details"`
	***REMOVED***
	err = json.Unmarshal(body, &detailsErr)
	if err == nil && detailsErr.Details != "" ***REMOVED***
		switch statusCode ***REMOVED***
		case http.StatusUnauthorized:
			return errcode.ErrorCodeUnauthorized.WithMessage(detailsErr.Details)
		case http.StatusTooManyRequests:
			return errcode.ErrorCodeTooManyRequests.WithMessage(detailsErr.Details)
		default:
			return errcode.ErrorCodeUnknown.WithMessage(detailsErr.Details)
		***REMOVED***
	***REMOVED***

	if err := json.Unmarshal(body, &errors); err != nil ***REMOVED***
		return &UnexpectedHTTPResponseError***REMOVED***
			ParseErr:   err,
			StatusCode: statusCode,
			Response:   body,
		***REMOVED***
	***REMOVED***

	if len(errors) == 0 ***REMOVED***
		// If there was no error specified in the body, return
		// UnexpectedHTTPResponseError.
		return &UnexpectedHTTPResponseError***REMOVED***
			ParseErr:   ErrNoErrorsInBody,
			StatusCode: statusCode,
			Response:   body,
		***REMOVED***
	***REMOVED***

	return errors
***REMOVED***

func makeErrorList(err error) []error ***REMOVED***
	if errL, ok := err.(errcode.Errors); ok ***REMOVED***
		return []error(errL)
	***REMOVED***
	return []error***REMOVED***err***REMOVED***
***REMOVED***

func mergeErrors(err1, err2 error) error ***REMOVED***
	return errcode.Errors(append(makeErrorList(err1), makeErrorList(err2)...))
***REMOVED***

// HandleErrorResponse returns error parsed from HTTP response for an
// unsuccessful HTTP response code (in the range 400 - 499 inclusive). An
// UnexpectedHTTPStatusError returned for response code outside of expected
// range.
func HandleErrorResponse(resp *http.Response) error ***REMOVED***
	if resp.StatusCode >= 400 && resp.StatusCode < 500 ***REMOVED***
		// Check for OAuth errors within the `WWW-Authenticate` header first
		// See https://tools.ietf.org/html/rfc6750#section-3
		for _, c := range challenge.ResponseChallenges(resp) ***REMOVED***
			if c.Scheme == "bearer" ***REMOVED***
				var err errcode.Error
				// codes defined at https://tools.ietf.org/html/rfc6750#section-3.1
				switch c.Parameters["error"] ***REMOVED***
				case "invalid_token":
					err.Code = errcode.ErrorCodeUnauthorized
				case "insufficient_scope":
					err.Code = errcode.ErrorCodeDenied
				default:
					continue
				***REMOVED***
				if description := c.Parameters["error_description"]; description != "" ***REMOVED***
					err.Message = description
				***REMOVED*** else ***REMOVED***
					err.Message = err.Code.Message()
				***REMOVED***

				return mergeErrors(err, parseHTTPErrorResponse(resp.StatusCode, resp.Body))
			***REMOVED***
		***REMOVED***
		err := parseHTTPErrorResponse(resp.StatusCode, resp.Body)
		if uErr, ok := err.(*UnexpectedHTTPResponseError); ok && resp.StatusCode == 401 ***REMOVED***
			return errcode.ErrorCodeUnauthorized.WithDetail(uErr.Response)
		***REMOVED***
		return err
	***REMOVED***
	return &UnexpectedHTTPStatusError***REMOVED***Status: resp.Status***REMOVED***
***REMOVED***

// SuccessStatus returns true if the argument is a successful HTTP response
// code (in the range 200 - 399 inclusive).
func SuccessStatus(status int) bool ***REMOVED***
	return status >= 200 && status <= 399
***REMOVED***
