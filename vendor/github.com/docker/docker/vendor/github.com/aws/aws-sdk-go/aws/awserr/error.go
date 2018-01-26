// Package awserr represents API error interface accessors for the SDK.
package awserr

// An Error wraps lower level errors with code, message and an original error.
// The underlying concrete error type may also satisfy other interfaces which
// can be to used to obtain more specific information about the error.
//
// Calling Error() or String() will always include the full information about
// an error based on its underlying type.
//
// Example:
//
//     output, err := s3manage.Upload(svc, input, opts)
//     if err != nil ***REMOVED***
//         if awsErr, ok := err.(awserr.Error); ok ***REMOVED***
//             // Get error details
//             log.Println("Error:", awsErr.Code(), awsErr.Message())
//
//             // Prints out full error message, including original error if there was one.
//             log.Println("Error:", awsErr.Error())
//
//             // Get original error
//             if origErr := awsErr.OrigErr(); origErr != nil ***REMOVED***
//                 // operate on original error.
//         ***REMOVED***
//     ***REMOVED*** else ***REMOVED***
//             fmt.Println(err.Error())
//     ***REMOVED***
// ***REMOVED***
//
type Error interface ***REMOVED***
	// Satisfy the generic error interface.
	error

	// Returns the short phrase depicting the classification of the error.
	Code() string

	// Returns the error details message.
	Message() string

	// Returns the original error if one was set.  Nil is returned if not set.
	OrigErr() error
***REMOVED***

// BatchError is a batch of errors which also wraps lower level errors with
// code, message, and original errors. Calling Error() will include all errors
// that occurred in the batch.
//
// Deprecated: Replaced with BatchedErrors. Only defined for backwards
// compatibility.
type BatchError interface ***REMOVED***
	// Satisfy the generic error interface.
	error

	// Returns the short phrase depicting the classification of the error.
	Code() string

	// Returns the error details message.
	Message() string

	// Returns the original error if one was set.  Nil is returned if not set.
	OrigErrs() []error
***REMOVED***

// BatchedErrors is a batch of errors which also wraps lower level errors with
// code, message, and original errors. Calling Error() will include all errors
// that occurred in the batch.
//
// Replaces BatchError
type BatchedErrors interface ***REMOVED***
	// Satisfy the base Error interface.
	Error

	// Returns the original error if one was set.  Nil is returned if not set.
	OrigErrs() []error
***REMOVED***

// New returns an Error object described by the code, message, and origErr.
//
// If origErr satisfies the Error interface it will not be wrapped within a new
// Error object and will instead be returned.
func New(code, message string, origErr error) Error ***REMOVED***
	var errs []error
	if origErr != nil ***REMOVED***
		errs = append(errs, origErr)
	***REMOVED***
	return newBaseError(code, message, errs)
***REMOVED***

// NewBatchError returns an BatchedErrors with a collection of errors as an
// array of errors.
func NewBatchError(code, message string, errs []error) BatchedErrors ***REMOVED***
	return newBaseError(code, message, errs)
***REMOVED***

// A RequestFailure is an interface to extract request failure information from
// an Error such as the request ID of the failed request returned by a service.
// RequestFailures may not always have a requestID value if the request failed
// prior to reaching the service such as a connection error.
//
// Example:
//
//     output, err := s3manage.Upload(svc, input, opts)
//     if err != nil ***REMOVED***
//         if reqerr, ok := err.(RequestFailure); ok ***REMOVED***
//             log.Println("Request failed", reqerr.Code(), reqerr.Message(), reqerr.RequestID())
//     ***REMOVED*** else ***REMOVED***
//             log.Println("Error:", err.Error())
//     ***REMOVED***
// ***REMOVED***
//
// Combined with awserr.Error:
//
//    output, err := s3manage.Upload(svc, input, opts)
//    if err != nil ***REMOVED***
//        if awsErr, ok := err.(awserr.Error); ok ***REMOVED***
//            // Generic AWS Error with Code, Message, and original error (if any)
//            fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
//
//            if reqErr, ok := err.(awserr.RequestFailure); ok ***REMOVED***
//                // A service error occurred
//                fmt.Println(reqErr.StatusCode(), reqErr.RequestID())
//        ***REMOVED***
//    ***REMOVED*** else ***REMOVED***
//            fmt.Println(err.Error())
//    ***REMOVED***
//***REMOVED***
//
type RequestFailure interface ***REMOVED***
	Error

	// The status code of the HTTP response.
	StatusCode() int

	// The request ID returned by the service for a request failure. This will
	// be empty if no request ID is available such as the request failed due
	// to a connection error.
	RequestID() string
***REMOVED***

// NewRequestFailure returns a new request error wrapper for the given Error
// provided.
func NewRequestFailure(err Error, statusCode int, reqID string) RequestFailure ***REMOVED***
	return newRequestError(err, statusCode, reqID)
***REMOVED***
