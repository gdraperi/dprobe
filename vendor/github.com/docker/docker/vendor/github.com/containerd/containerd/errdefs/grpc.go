package errdefs

import (
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPC will attempt to map the backend containerd error into a grpc error,
// using the original error message as a description.
//
// Further information may be extracted from certain errors depending on their
// type.
//
// If the error is unmapped, the original error will be returned to be handled
// by the regular grpc error handling stack.
func ToGRPC(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***

	if isGRPCError(err) ***REMOVED***
		// error has already been mapped to grpc
		return err
	***REMOVED***

	switch ***REMOVED***
	case IsInvalidArgument(err):
		return status.Errorf(codes.InvalidArgument, err.Error())
	case IsNotFound(err):
		return status.Errorf(codes.NotFound, err.Error())
	case IsAlreadyExists(err):
		return status.Errorf(codes.AlreadyExists, err.Error())
	case IsFailedPrecondition(err):
		return status.Errorf(codes.FailedPrecondition, err.Error())
	case IsUnavailable(err):
		return status.Errorf(codes.Unavailable, err.Error())
	case IsNotImplemented(err):
		return status.Errorf(codes.Unimplemented, err.Error())
	***REMOVED***

	return err
***REMOVED***

// ToGRPCf maps the error to grpc error codes, assembling the formatting string
// and combining it with the target error string.
//
// This is equivalent to errors.ToGRPC(errors.Wrapf(err, format, args...))
func ToGRPCf(err error, format string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return ToGRPC(errors.Wrapf(err, format, args...))
***REMOVED***

// FromGRPC returns the underlying error from a grpc service based on the grpc error code
func FromGRPC(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***

	var cls error // divide these into error classes, becomes the cause

	switch code(err) ***REMOVED***
	case codes.InvalidArgument:
		cls = ErrInvalidArgument
	case codes.AlreadyExists:
		cls = ErrAlreadyExists
	case codes.NotFound:
		cls = ErrNotFound
	case codes.Unavailable:
		cls = ErrUnavailable
	case codes.FailedPrecondition:
		cls = ErrFailedPrecondition
	case codes.Unimplemented:
		cls = ErrNotImplemented
	default:
		cls = ErrUnknown
	***REMOVED***

	msg := rebaseMessage(cls, err)
	if msg != "" ***REMOVED***
		err = errors.Wrapf(cls, msg)
	***REMOVED*** else ***REMOVED***
		err = errors.WithStack(cls)
	***REMOVED***

	return err
***REMOVED***

// rebaseMessage removes the repeats for an error at the end of an error
// string. This will happen when taking an error over grpc then remapping it.
//
// Effectively, we just remove the string of cls from the end of err if it
// appears there.
func rebaseMessage(cls error, err error) string ***REMOVED***
	desc := errDesc(err)
	clss := cls.Error()
	if desc == clss ***REMOVED***
		return ""
	***REMOVED***

	return strings.TrimSuffix(desc, ": "+clss)
***REMOVED***

func isGRPCError(err error) bool ***REMOVED***
	_, ok := status.FromError(err)
	return ok
***REMOVED***

func code(err error) codes.Code ***REMOVED***
	if s, ok := status.FromError(err); ok ***REMOVED***
		return s.Code()
	***REMOVED***
	return codes.Unknown
***REMOVED***

func errDesc(err error) string ***REMOVED***
	if s, ok := status.FromError(err); ok ***REMOVED***
		return s.Message()
	***REMOVED***
	return err.Error()
***REMOVED***
