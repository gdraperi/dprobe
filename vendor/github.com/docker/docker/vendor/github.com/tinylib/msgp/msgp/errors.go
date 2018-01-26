package msgp

import (
	"fmt"
	"reflect"
)

var (
	// ErrShortBytes is returned when the
	// slice being decoded is too short to
	// contain the contents of the message
	ErrShortBytes error = errShort***REMOVED******REMOVED***

	// this error is only returned
	// if we reach code that should
	// be unreachable
	fatal error = errFatal***REMOVED******REMOVED***
)

// Error is the interface satisfied
// by all of the errors that originate
// from this package.
type Error interface ***REMOVED***
	error

	// Resumable returns whether
	// or not the error means that
	// the stream of data is malformed
	// and  the information is unrecoverable.
	Resumable() bool
***REMOVED***

type errShort struct***REMOVED******REMOVED***

func (e errShort) Error() string   ***REMOVED*** return "msgp: too few bytes left to read object" ***REMOVED***
func (e errShort) Resumable() bool ***REMOVED*** return false ***REMOVED***

type errFatal struct***REMOVED******REMOVED***

func (f errFatal) Error() string   ***REMOVED*** return "msgp: fatal decoding error (unreachable code)" ***REMOVED***
func (f errFatal) Resumable() bool ***REMOVED*** return false ***REMOVED***

// ArrayError is an error returned
// when decoding a fix-sized array
// of the wrong size
type ArrayError struct ***REMOVED***
	Wanted uint32
	Got    uint32
***REMOVED***

// Error implements the error interface
func (a ArrayError) Error() string ***REMOVED***
	return fmt.Sprintf("msgp: wanted array of size %d; got %d", a.Wanted, a.Got)
***REMOVED***

// Resumable is always 'true' for ArrayErrors
func (a ArrayError) Resumable() bool ***REMOVED*** return true ***REMOVED***

// IntOverflow is returned when a call
// would downcast an integer to a type
// with too few bits to hold its value.
type IntOverflow struct ***REMOVED***
	Value         int64 // the value of the integer
	FailedBitsize int   // the bit size that the int64 could not fit into
***REMOVED***

// Error implements the error interface
func (i IntOverflow) Error() string ***REMOVED***
	return fmt.Sprintf("msgp: %d overflows int%d", i.Value, i.FailedBitsize)
***REMOVED***

// Resumable is always 'true' for overflows
func (i IntOverflow) Resumable() bool ***REMOVED*** return true ***REMOVED***

// UintOverflow is returned when a call
// would downcast an unsigned integer to a type
// with too few bits to hold its value
type UintOverflow struct ***REMOVED***
	Value         uint64 // value of the uint
	FailedBitsize int    // the bit size that couldn't fit the value
***REMOVED***

// Error implements the error interface
func (u UintOverflow) Error() string ***REMOVED***
	return fmt.Sprintf("msgp: %d overflows uint%d", u.Value, u.FailedBitsize)
***REMOVED***

// Resumable is always 'true' for overflows
func (u UintOverflow) Resumable() bool ***REMOVED*** return true ***REMOVED***

// A TypeError is returned when a particular
// decoding method is unsuitable for decoding
// a particular MessagePack value.
type TypeError struct ***REMOVED***
	Method  Type // Type expected by method
	Encoded Type // Type actually encoded
***REMOVED***

// Error implements the error interface
func (t TypeError) Error() string ***REMOVED***
	return fmt.Sprintf("msgp: attempted to decode type %q with method for %q", t.Encoded, t.Method)
***REMOVED***

// Resumable returns 'true' for TypeErrors
func (t TypeError) Resumable() bool ***REMOVED*** return true ***REMOVED***

// returns either InvalidPrefixError or
// TypeError depending on whether or not
// the prefix is recognized
func badPrefix(want Type, lead byte) error ***REMOVED***
	t := sizes[lead].typ
	if t == InvalidType ***REMOVED***
		return InvalidPrefixError(lead)
	***REMOVED***
	return TypeError***REMOVED***Method: want, Encoded: t***REMOVED***
***REMOVED***

// InvalidPrefixError is returned when a bad encoding
// uses a prefix that is not recognized in the MessagePack standard.
// This kind of error is unrecoverable.
type InvalidPrefixError byte

// Error implements the error interface
func (i InvalidPrefixError) Error() string ***REMOVED***
	return fmt.Sprintf("msgp: unrecognized type prefix 0x%x", byte(i))
***REMOVED***

// Resumable returns 'false' for InvalidPrefixErrors
func (i InvalidPrefixError) Resumable() bool ***REMOVED*** return false ***REMOVED***

// ErrUnsupportedType is returned
// when a bad argument is supplied
// to a function that takes `interface***REMOVED******REMOVED***`.
type ErrUnsupportedType struct ***REMOVED***
	T reflect.Type
***REMOVED***

// Error implements error
func (e *ErrUnsupportedType) Error() string ***REMOVED*** return fmt.Sprintf("msgp: type %q not supported", e.T) ***REMOVED***

// Resumable returns 'true' for ErrUnsupportedType
func (e *ErrUnsupportedType) Resumable() bool ***REMOVED*** return true ***REMOVED***
