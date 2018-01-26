package exec

import "github.com/pkg/errors"

var (
	// ErrRuntimeUnsupported encountered when a task requires a runtime
	// unsupported by the executor.
	ErrRuntimeUnsupported = errors.New("exec: unsupported runtime")

	// ErrTaskPrepared is called if the task is already prepared.
	ErrTaskPrepared = errors.New("exec: task already prepared")

	// ErrTaskStarted can be returned from any operation that cannot be
	// performed because the task has already been started. This does not imply
	// that the task is running but rather that it is no longer valid to call
	// Start.
	ErrTaskStarted = errors.New("exec: task already started")

	// ErrTaskUpdateRejected is returned if a task update is rejected by a controller.
	ErrTaskUpdateRejected = errors.New("exec: task update rejected")

	// ErrControllerClosed returned when a task controller has been closed.
	ErrControllerClosed = errors.New("exec: controller closed")

	// ErrTaskRetry is returned by Do when an operation failed by should be
	// retried. The status should still be reported in this case.
	ErrTaskRetry = errors.New("exec: task retry")

	// ErrTaskNoop returns when the a subsequent call to Do will not result in
	// advancing the task. Callers should avoid calling Do until the task has been updated.
	ErrTaskNoop = errors.New("exec: task noop")
)

// ExitCoder is implemented by errors that have an exit code.
type ExitCoder interface ***REMOVED***
	// ExitCode returns the exit code.
	ExitCode() int
***REMOVED***

// Temporary indicates whether or not the error condition is temporary.
//
// If this is encountered in the controller, the failing operation will be
// retried when this returns true. Otherwise, the operation is considered
// fatal.
type Temporary interface ***REMOVED***
	Temporary() bool
***REMOVED***

// MakeTemporary makes the error temporary.
func MakeTemporary(err error) error ***REMOVED***
	if IsTemporary(err) ***REMOVED***
		return err
	***REMOVED***

	return temporary***REMOVED***err***REMOVED***
***REMOVED***

type temporary struct ***REMOVED***
	error
***REMOVED***

func (t temporary) Cause() error    ***REMOVED*** return t.error ***REMOVED***
func (t temporary) Temporary() bool ***REMOVED*** return true ***REMOVED***

// IsTemporary returns true if the error or a recursive cause returns true for
// temporary.
func IsTemporary(err error) bool ***REMOVED***
	for err != nil ***REMOVED***
		if tmp, ok := err.(Temporary); ok && tmp.Temporary() ***REMOVED***
			return true
		***REMOVED***

		cause := errors.Cause(err)
		if cause == err ***REMOVED***
			break
		***REMOVED***

		err = cause
	***REMOVED***

	return false
***REMOVED***
