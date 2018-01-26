package errdefs

type causer interface ***REMOVED***
	Cause() error
***REMOVED***

func getImplementer(err error) error ***REMOVED***
	switch e := err.(type) ***REMOVED***
	case
		ErrNotFound,
		ErrInvalidParameter,
		ErrConflict,
		ErrUnauthorized,
		ErrUnavailable,
		ErrForbidden,
		ErrSystem,
		ErrNotModified,
		ErrAlreadyExists,
		ErrNotImplemented,
		ErrCancelled,
		ErrDeadline,
		ErrDataLoss,
		ErrUnknown:
		return e
	case causer:
		return getImplementer(e.Cause())
	default:
		return err
	***REMOVED***
***REMOVED***

// IsNotFound returns if the passed in error is an ErrNotFound
func IsNotFound(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrNotFound)
	return ok
***REMOVED***

// IsInvalidParameter returns if the passed in error is an ErrInvalidParameter
func IsInvalidParameter(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrInvalidParameter)
	return ok
***REMOVED***

// IsConflict returns if the passed in error is an ErrConflict
func IsConflict(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrConflict)
	return ok
***REMOVED***

// IsUnauthorized returns if the the passed in error is an ErrUnauthorized
func IsUnauthorized(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrUnauthorized)
	return ok
***REMOVED***

// IsUnavailable returns if the passed in error is an ErrUnavailable
func IsUnavailable(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrUnavailable)
	return ok
***REMOVED***

// IsForbidden returns if the passed in error is an ErrForbidden
func IsForbidden(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrForbidden)
	return ok
***REMOVED***

// IsSystem returns if the passed in error is an ErrSystem
func IsSystem(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrSystem)
	return ok
***REMOVED***

// IsNotModified returns if the passed in error is a NotModified error
func IsNotModified(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrNotModified)
	return ok
***REMOVED***

// IsAlreadyExists returns if the passed in error is a AlreadyExists error
func IsAlreadyExists(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrAlreadyExists)
	return ok
***REMOVED***

// IsNotImplemented returns if the passed in error is an ErrNotImplemented
func IsNotImplemented(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrNotImplemented)
	return ok
***REMOVED***

// IsUnknown returns if the passed in error is an ErrUnknown
func IsUnknown(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrUnknown)
	return ok
***REMOVED***

// IsCancelled returns if the passed in error is an ErrCancelled
func IsCancelled(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrCancelled)
	return ok
***REMOVED***

// IsDeadline returns if the passed in error is an ErrDeadline
func IsDeadline(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrDeadline)
	return ok
***REMOVED***

// IsDataLoss returns if the passed in error is an ErrDataLoss
func IsDataLoss(err error) bool ***REMOVED***
	_, ok := getImplementer(err).(ErrDataLoss)
	return ok
***REMOVED***
