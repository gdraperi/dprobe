package graphdriver

const (
	// ErrNotSupported returned when driver is not supported.
	ErrNotSupported NotSupportedError = "driver not supported"
	// ErrPrerequisites returned when driver does not meet prerequisites.
	ErrPrerequisites NotSupportedError = "prerequisites for driver not satisfied (wrong filesystem?)"
	// ErrIncompatibleFS returned when file system is not supported.
	ErrIncompatibleFS NotSupportedError = "backing file system is unsupported for this graph driver"
)

// ErrUnSupported signals that the graph-driver is not supported on the current configuration
type ErrUnSupported interface ***REMOVED***
	NotSupported()
***REMOVED***

// NotSupportedError signals that the graph-driver is not supported on the current configuration
type NotSupportedError string

func (e NotSupportedError) Error() string ***REMOVED***
	return string(e)
***REMOVED***

// NotSupported signals that a graph-driver is not supported.
func (e NotSupportedError) NotSupported() ***REMOVED******REMOVED***

// IsDriverNotSupported returns true if the error initializing
// the graph driver is a non-supported error.
func IsDriverNotSupported(err error) bool ***REMOVED***
	switch err.(type) ***REMOVED***
	case ErrUnSupported:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***
