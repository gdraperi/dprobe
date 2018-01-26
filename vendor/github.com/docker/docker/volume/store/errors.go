package store

import (
	"strings"
)

const (
	// errVolumeInUse is a typed error returned when trying to remove a volume that is currently in use by a container
	errVolumeInUse conflictError = "volume is in use"
	// errNoSuchVolume is a typed error returned if the requested volume doesn't exist in the volume store
	errNoSuchVolume notFoundError = "no such volume"
	// errNameConflict is a typed error returned on create when a volume exists with the given name, but for a different driver
	errNameConflict conflictError = "volume name must be unique"
)

type conflictError string

func (e conflictError) Error() string ***REMOVED***
	return string(e)
***REMOVED***
func (conflictError) Conflict() ***REMOVED******REMOVED***

type notFoundError string

func (e notFoundError) Error() string ***REMOVED***
	return string(e)
***REMOVED***

func (notFoundError) NotFound() ***REMOVED******REMOVED***

// OpErr is the error type returned by functions in the store package. It describes
// the operation, volume name, and error.
type OpErr struct ***REMOVED***
	// Err is the error that occurred during the operation.
	Err error
	// Op is the operation which caused the error, such as "create", or "list".
	Op string
	// Name is the name of the resource being requested for this op, typically the volume name or the driver name.
	Name string
	// Refs is the list of references associated with the resource.
	Refs []string
***REMOVED***

// Error satisfies the built-in error interface type.
func (e *OpErr) Error() string ***REMOVED***
	if e == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	s := e.Op
	if e.Name != "" ***REMOVED***
		s = s + " " + e.Name
	***REMOVED***

	s = s + ": " + e.Err.Error()
	if len(e.Refs) > 0 ***REMOVED***
		s = s + " - " + "[" + strings.Join(e.Refs, ", ") + "]"
	***REMOVED***
	return s
***REMOVED***

// Cause returns the error the caused this error
func (e *OpErr) Cause() error ***REMOVED***
	return e.Err
***REMOVED***

// IsInUse returns a boolean indicating whether the error indicates that a
// volume is in use
func IsInUse(err error) bool ***REMOVED***
	return isErr(err, errVolumeInUse)
***REMOVED***

// IsNotExist returns a boolean indicating whether the error indicates that the volume does not exist
func IsNotExist(err error) bool ***REMOVED***
	return isErr(err, errNoSuchVolume)
***REMOVED***

// IsNameConflict returns a boolean indicating whether the error indicates that a
// volume name is already taken
func IsNameConflict(err error) bool ***REMOVED***
	return isErr(err, errNameConflict)
***REMOVED***

type causal interface ***REMOVED***
	Cause() error
***REMOVED***

func isErr(err error, expected error) bool ***REMOVED***
	switch pe := err.(type) ***REMOVED***
	case nil:
		return false
	case causal:
		return isErr(pe.Cause(), expected)
	***REMOVED***
	return err == expected
***REMOVED***
