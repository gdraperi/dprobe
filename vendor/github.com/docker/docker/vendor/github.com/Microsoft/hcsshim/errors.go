package hcsshim

import (
	"errors"
	"fmt"
	"syscall"
)

var (
	// ErrComputeSystemDoesNotExist is an error encountered when the container being operated on no longer exists
	ErrComputeSystemDoesNotExist = syscall.Errno(0xc037010e)

	// ErrElementNotFound is an error encountered when the object being referenced does not exist
	ErrElementNotFound = syscall.Errno(0x490)

	// ErrElementNotFound is an error encountered when the object being referenced does not exist
	ErrNotSupported = syscall.Errno(0x32)

	// ErrInvalidData is an error encountered when the request being sent to hcs is invalid/unsupported
	// decimal -2147024883 / hex 0x8007000d
	ErrInvalidData = syscall.Errno(0xd)

	// ErrHandleClose is an error encountered when the handle generating the notification being waited on has been closed
	ErrHandleClose = errors.New("hcsshim: the handle generating this notification has been closed")

	// ErrAlreadyClosed is an error encountered when using a handle that has been closed by the Close method
	ErrAlreadyClosed = errors.New("hcsshim: the handle has already been closed")

	// ErrInvalidNotificationType is an error encountered when an invalid notification type is used
	ErrInvalidNotificationType = errors.New("hcsshim: invalid notification type")

	// ErrInvalidProcessState is an error encountered when the process is not in a valid state for the requested operation
	ErrInvalidProcessState = errors.New("the process is in an invalid state for the attempted operation")

	// ErrTimeout is an error encountered when waiting on a notification times out
	ErrTimeout = errors.New("hcsshim: timeout waiting for notification")

	// ErrUnexpectedContainerExit is the error encountered when a container exits while waiting for
	// a different expected notification
	ErrUnexpectedContainerExit = errors.New("unexpected container exit")

	// ErrUnexpectedProcessAbort is the error encountered when communication with the compute service
	// is lost while waiting for a notification
	ErrUnexpectedProcessAbort = errors.New("lost communication with compute service")

	// ErrUnexpectedValue is an error encountered when hcs returns an invalid value
	ErrUnexpectedValue = errors.New("unexpected value returned from hcs")

	// ErrVmcomputeAlreadyStopped is an error encountered when a shutdown or terminate request is made on a stopped container
	ErrVmcomputeAlreadyStopped = syscall.Errno(0xc0370110)

	// ErrVmcomputeOperationPending is an error encountered when the operation is being completed asynchronously
	ErrVmcomputeOperationPending = syscall.Errno(0xC0370103)

	// ErrVmcomputeOperationInvalidState is an error encountered when the compute system is not in a valid state for the requested operation
	ErrVmcomputeOperationInvalidState = syscall.Errno(0xc0370105)

	// ErrProcNotFound is an error encountered when the the process cannot be found
	ErrProcNotFound = syscall.Errno(0x7f)

	// ErrVmcomputeOperationAccessIsDenied is an error which can be encountered when enumerating compute systems in RS1/RS2
	// builds when the underlying silo might be in the process of terminating. HCS was fixed in RS3.
	ErrVmcomputeOperationAccessIsDenied = syscall.Errno(0x5)

	// ErrVmcomputeInvalidJSON is an error encountered when the compute system does not support/understand the messages sent by management
	ErrVmcomputeInvalidJSON = syscall.Errno(0xc037010d)

	// ErrVmcomputeUnknownMessage is an error encountered guest compute system doesn't support the message
	ErrVmcomputeUnknownMessage = syscall.Errno(0xc037010b)

	// ErrNotSupported is an error encountered when hcs doesn't support the request
	ErrPlatformNotSupported = errors.New("unsupported platform request")
)

type EndpointNotFoundError struct ***REMOVED***
	EndpointName string
***REMOVED***

func (e EndpointNotFoundError) Error() string ***REMOVED***
	return fmt.Sprintf("Endpoint %s not found", e.EndpointName)
***REMOVED***

type NetworkNotFoundError struct ***REMOVED***
	NetworkName string
***REMOVED***

func (e NetworkNotFoundError) Error() string ***REMOVED***
	return fmt.Sprintf("Network %s not found", e.NetworkName)
***REMOVED***

// ProcessError is an error encountered in HCS during an operation on a Process object
type ProcessError struct ***REMOVED***
	Process   *process
	Operation string
	ExtraInfo string
	Err       error
***REMOVED***

// ContainerError is an error encountered in HCS during an operation on a Container object
type ContainerError struct ***REMOVED***
	Container *container
	Operation string
	ExtraInfo string
	Err       error
***REMOVED***

func (e *ContainerError) Error() string ***REMOVED***
	if e == nil ***REMOVED***
		return "<nil>"
	***REMOVED***

	if e.Container == nil ***REMOVED***
		return "unexpected nil container for error: " + e.Err.Error()
	***REMOVED***

	s := "container " + e.Container.id

	if e.Operation != "" ***REMOVED***
		s += " encountered an error during " + e.Operation
	***REMOVED***

	switch e.Err.(type) ***REMOVED***
	case nil:
		break
	case syscall.Errno:
		s += fmt.Sprintf(": failure in a Windows system call: %s (0x%x)", e.Err, win32FromError(e.Err))
	default:
		s += fmt.Sprintf(": %s", e.Err.Error())
	***REMOVED***

	if e.ExtraInfo != "" ***REMOVED***
		s += " extra info: " + e.ExtraInfo
	***REMOVED***

	return s
***REMOVED***

func makeContainerError(container *container, operation string, extraInfo string, err error) error ***REMOVED***
	// Don't double wrap errors
	if _, ok := err.(*ContainerError); ok ***REMOVED***
		return err
	***REMOVED***
	containerError := &ContainerError***REMOVED***Container: container, Operation: operation, ExtraInfo: extraInfo, Err: err***REMOVED***
	return containerError
***REMOVED***

func (e *ProcessError) Error() string ***REMOVED***
	if e == nil ***REMOVED***
		return "<nil>"
	***REMOVED***

	if e.Process == nil ***REMOVED***
		return "Unexpected nil process for error: " + e.Err.Error()
	***REMOVED***

	s := fmt.Sprintf("process %d", e.Process.processID)

	if e.Process.container != nil ***REMOVED***
		s += " in container " + e.Process.container.id
	***REMOVED***

	if e.Operation != "" ***REMOVED***
		s += " encountered an error during " + e.Operation
	***REMOVED***

	switch e.Err.(type) ***REMOVED***
	case nil:
		break
	case syscall.Errno:
		s += fmt.Sprintf(": failure in a Windows system call: %s (0x%x)", e.Err, win32FromError(e.Err))
	default:
		s += fmt.Sprintf(": %s", e.Err.Error())
	***REMOVED***

	return s
***REMOVED***

func makeProcessError(process *process, operation string, extraInfo string, err error) error ***REMOVED***
	// Don't double wrap errors
	if _, ok := err.(*ProcessError); ok ***REMOVED***
		return err
	***REMOVED***
	processError := &ProcessError***REMOVED***Process: process, Operation: operation, ExtraInfo: extraInfo, Err: err***REMOVED***
	return processError
***REMOVED***

// IsNotExist checks if an error is caused by the Container or Process not existing.
// Note: Currently, ErrElementNotFound can mean that a Process has either
// already exited, or does not exist. Both IsAlreadyStopped and IsNotExist
// will currently return true when the error is ErrElementNotFound or ErrProcNotFound.
func IsNotExist(err error) bool ***REMOVED***
	err = getInnerError(err)
	if _, ok := err.(EndpointNotFoundError); ok ***REMOVED***
		return true
	***REMOVED***
	if _, ok := err.(NetworkNotFoundError); ok ***REMOVED***
		return true
	***REMOVED***
	return err == ErrComputeSystemDoesNotExist ||
		err == ErrElementNotFound ||
		err == ErrProcNotFound
***REMOVED***

// IsAlreadyClosed checks if an error is caused by the Container or Process having been
// already closed by a call to the Close() method.
func IsAlreadyClosed(err error) bool ***REMOVED***
	err = getInnerError(err)
	return err == ErrAlreadyClosed
***REMOVED***

// IsPending returns a boolean indicating whether the error is that
// the requested operation is being completed in the background.
func IsPending(err error) bool ***REMOVED***
	err = getInnerError(err)
	return err == ErrVmcomputeOperationPending
***REMOVED***

// IsTimeout returns a boolean indicating whether the error is caused by
// a timeout waiting for the operation to complete.
func IsTimeout(err error) bool ***REMOVED***
	err = getInnerError(err)
	return err == ErrTimeout
***REMOVED***

// IsAlreadyStopped returns a boolean indicating whether the error is caused by
// a Container or Process being already stopped.
// Note: Currently, ErrElementNotFound can mean that a Process has either
// already exited, or does not exist. Both IsAlreadyStopped and IsNotExist
// will currently return true when the error is ErrElementNotFound or ErrProcNotFound.
func IsAlreadyStopped(err error) bool ***REMOVED***
	err = getInnerError(err)
	return err == ErrVmcomputeAlreadyStopped ||
		err == ErrElementNotFound ||
		err == ErrProcNotFound
***REMOVED***

// IsNotSupported returns a boolean indicating whether the error is caused by
// unsupported platform requests
// Note: Currently Unsupported platform requests can be mean either
// ErrVmcomputeInvalidJSON, ErrInvalidData, ErrNotSupported or ErrVmcomputeUnknownMessage
// is thrown from the Platform
func IsNotSupported(err error) bool ***REMOVED***
	err = getInnerError(err)
	// If Platform doesn't recognize or support the request sent, below errors are seen
	return err == ErrVmcomputeInvalidJSON ||
		err == ErrInvalidData ||
		err == ErrNotSupported ||
		err == ErrVmcomputeUnknownMessage
***REMOVED***

func getInnerError(err error) error ***REMOVED***
	switch pe := err.(type) ***REMOVED***
	case nil:
		return nil
	case *ContainerError:
		err = pe.Err
	case *ProcessError:
		err = pe.Err
	***REMOVED***
	return err
***REMOVED***
