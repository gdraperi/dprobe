package multierror

import (
	"fmt"
)

// Error is an error type to track multiple errors. This is used to
// accumulate errors in cases and return them as a single "error".
type Error struct ***REMOVED***
	Errors      []error
	ErrorFormat ErrorFormatFunc
***REMOVED***

func (e *Error) Error() string ***REMOVED***
	fn := e.ErrorFormat
	if fn == nil ***REMOVED***
		fn = ListFormatFunc
	***REMOVED***

	return fn(e.Errors)
***REMOVED***

// ErrorOrNil returns an error interface if this Error represents
// a list of errors, or returns nil if the list of errors is empty. This
// function is useful at the end of accumulation to make sure that the value
// returned represents the existence of errors.
func (e *Error) ErrorOrNil() error ***REMOVED***
	if e == nil ***REMOVED***
		return nil
	***REMOVED***
	if len(e.Errors) == 0 ***REMOVED***
		return nil
	***REMOVED***

	return e
***REMOVED***

func (e *Error) GoString() string ***REMOVED***
	return fmt.Sprintf("*%#v", *e)
***REMOVED***

// WrappedErrors returns the list of errors that this Error is wrapping.
// It is an implementatin of the errwrap.Wrapper interface so that
// multierror.Error can be used with that library.
//
// This method is not safe to be called concurrently and is no different
// than accessing the Errors field directly. It is implementd only to
// satisfy the errwrap.Wrapper interface.
func (e *Error) WrappedErrors() []error ***REMOVED***
	return e.Errors
***REMOVED***
