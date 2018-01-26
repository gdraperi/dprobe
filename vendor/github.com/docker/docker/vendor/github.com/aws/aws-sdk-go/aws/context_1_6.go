// +build !go1.7

package aws

import "time"

// An emptyCtx is a copy of the Go 1.7 context.emptyCtx type. This is copied to
// provide a 1.6 and 1.5 safe version of context that is compatible with Go
// 1.7's Context.
//
// An emptyCtx is never canceled, has no values, and has no deadline. It is not
// struct***REMOVED******REMOVED***, since vars of this type must have distinct addresses.
type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) ***REMOVED***
	return
***REMOVED***

func (*emptyCtx) Done() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***

func (*emptyCtx) Err() error ***REMOVED***
	return nil
***REMOVED***

func (*emptyCtx) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***

func (e *emptyCtx) String() string ***REMOVED***
	switch e ***REMOVED***
	case backgroundCtx:
		return "aws.BackgroundContext"
	***REMOVED***
	return "unknown empty Context"
***REMOVED***

var (
	backgroundCtx = new(emptyCtx)
)
