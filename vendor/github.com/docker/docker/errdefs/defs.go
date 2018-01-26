package errdefs

// ErrNotFound signals that the requested object doesn't exist
type ErrNotFound interface ***REMOVED***
	NotFound()
***REMOVED***

// ErrInvalidParameter signals that the user input is invalid
type ErrInvalidParameter interface ***REMOVED***
	InvalidParameter()
***REMOVED***

// ErrConflict signals that some internal state conflicts with the requested action and can't be performed.
// A change in state should be able to clear this error.
type ErrConflict interface ***REMOVED***
	Conflict()
***REMOVED***

// ErrUnauthorized is used to signify that the user is not authorized to perform a specific action
type ErrUnauthorized interface ***REMOVED***
	Unauthorized()
***REMOVED***

// ErrUnavailable signals that the requested action/subsystem is not available.
type ErrUnavailable interface ***REMOVED***
	Unavailable()
***REMOVED***

// ErrForbidden signals that the requested action cannot be performed under any circumstances.
// When a ErrForbidden is returned, the caller should never retry the action.
type ErrForbidden interface ***REMOVED***
	Forbidden()
***REMOVED***

// ErrSystem signals that some internal error occurred.
// An example of this would be a failed mount request.
type ErrSystem interface ***REMOVED***
	ErrSystem()
***REMOVED***

// ErrNotModified signals that an action can't be performed because it's already in the desired state
type ErrNotModified interface ***REMOVED***
	NotModified()
***REMOVED***

// ErrAlreadyExists is a special case of ErrConflict which signals that the desired object already exists
type ErrAlreadyExists interface ***REMOVED***
	AlreadyExists()
***REMOVED***

// ErrNotImplemented signals that the requested action/feature is not implemented on the system as configured.
type ErrNotImplemented interface ***REMOVED***
	NotImplemented()
***REMOVED***

// ErrUnknown signals that the kind of error that occurred is not known.
type ErrUnknown interface ***REMOVED***
	Unknown()
***REMOVED***

// ErrCancelled signals that the action was cancelled.
type ErrCancelled interface ***REMOVED***
	Cancelled()
***REMOVED***

// ErrDeadline signals that the deadline was reached before the action completed.
type ErrDeadline interface ***REMOVED***
	DeadlineExceeded()
***REMOVED***

// ErrDataLoss indicates that data was lost or there is data corruption.
type ErrDataLoss interface ***REMOVED***
	DataLoss()
***REMOVED***
