package errdefs

import "context"

type errNotFound struct***REMOVED*** error ***REMOVED***

func (errNotFound) NotFound() ***REMOVED******REMOVED***

func (e errNotFound) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// NotFound is a helper to create an error of the class with the same name from any error type
func NotFound(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errNotFound***REMOVED***err***REMOVED***
***REMOVED***

type errInvalidParameter struct***REMOVED*** error ***REMOVED***

func (errInvalidParameter) InvalidParameter() ***REMOVED******REMOVED***

func (e errInvalidParameter) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// InvalidParameter is a helper to create an error of the class with the same name from any error type
func InvalidParameter(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errInvalidParameter***REMOVED***err***REMOVED***
***REMOVED***

type errConflict struct***REMOVED*** error ***REMOVED***

func (errConflict) Conflict() ***REMOVED******REMOVED***

func (e errConflict) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// Conflict is a helper to create an error of the class with the same name from any error type
func Conflict(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errConflict***REMOVED***err***REMOVED***
***REMOVED***

type errUnauthorized struct***REMOVED*** error ***REMOVED***

func (errUnauthorized) Unauthorized() ***REMOVED******REMOVED***

func (e errUnauthorized) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// Unauthorized is a helper to create an error of the class with the same name from any error type
func Unauthorized(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errUnauthorized***REMOVED***err***REMOVED***
***REMOVED***

type errUnavailable struct***REMOVED*** error ***REMOVED***

func (errUnavailable) Unavailable() ***REMOVED******REMOVED***

func (e errUnavailable) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// Unavailable is a helper to create an error of the class with the same name from any error type
func Unavailable(err error) error ***REMOVED***
	return errUnavailable***REMOVED***err***REMOVED***
***REMOVED***

type errForbidden struct***REMOVED*** error ***REMOVED***

func (errForbidden) Forbidden() ***REMOVED******REMOVED***

func (e errForbidden) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// Forbidden is a helper to create an error of the class with the same name from any error type
func Forbidden(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errForbidden***REMOVED***err***REMOVED***
***REMOVED***

type errSystem struct***REMOVED*** error ***REMOVED***

func (errSystem) System() ***REMOVED******REMOVED***

func (e errSystem) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// System is a helper to create an error of the class with the same name from any error type
func System(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errSystem***REMOVED***err***REMOVED***
***REMOVED***

type errNotModified struct***REMOVED*** error ***REMOVED***

func (errNotModified) NotModified() ***REMOVED******REMOVED***

func (e errNotModified) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// NotModified is a helper to create an error of the class with the same name from any error type
func NotModified(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errNotModified***REMOVED***err***REMOVED***
***REMOVED***

type errAlreadyExists struct***REMOVED*** error ***REMOVED***

func (errAlreadyExists) AlreadyExists() ***REMOVED******REMOVED***

func (e errAlreadyExists) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// AlreadyExists is a helper to create an error of the class with the same name from any error type
func AlreadyExists(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errAlreadyExists***REMOVED***err***REMOVED***
***REMOVED***

type errNotImplemented struct***REMOVED*** error ***REMOVED***

func (errNotImplemented) NotImplemented() ***REMOVED******REMOVED***

func (e errNotImplemented) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// NotImplemented is a helper to create an error of the class with the same name from any error type
func NotImplemented(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errNotImplemented***REMOVED***err***REMOVED***
***REMOVED***

type errUnknown struct***REMOVED*** error ***REMOVED***

func (errUnknown) Unknown() ***REMOVED******REMOVED***

func (e errUnknown) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// Unknown is a helper to create an error of the class with the same name from any error type
func Unknown(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errUnknown***REMOVED***err***REMOVED***
***REMOVED***

type errCancelled struct***REMOVED*** error ***REMOVED***

func (errCancelled) Cancelled() ***REMOVED******REMOVED***

func (e errCancelled) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// Cancelled is a helper to create an error of the class with the same name from any error type
func Cancelled(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errCancelled***REMOVED***err***REMOVED***
***REMOVED***

type errDeadline struct***REMOVED*** error ***REMOVED***

func (errDeadline) DeadlineExceeded() ***REMOVED******REMOVED***

func (e errDeadline) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// Deadline is a helper to create an error of the class with the same name from any error type
func Deadline(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errDeadline***REMOVED***err***REMOVED***
***REMOVED***

type errDataLoss struct***REMOVED*** error ***REMOVED***

func (errDataLoss) DataLoss() ***REMOVED******REMOVED***

func (e errDataLoss) Cause() error ***REMOVED***
	return e.error
***REMOVED***

// DataLoss is a helper to create an error of the class with the same name from any error type
func DataLoss(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return errDataLoss***REMOVED***err***REMOVED***
***REMOVED***

// FromContext returns the error class from the passed in context
func FromContext(ctx context.Context) error ***REMOVED***
	e := ctx.Err()
	if e == nil ***REMOVED***
		return nil
	***REMOVED***

	if e == context.Canceled ***REMOVED***
		return Cancelled(e)
	***REMOVED***
	if e == context.DeadlineExceeded ***REMOVED***
		return Deadline(e)
	***REMOVED***
	return Unknown(e)
***REMOVED***
