package deepcopy

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/types"
)

// CopierFrom can be implemented if an object knows how to copy another into itself.
type CopierFrom interface ***REMOVED***
	// Copy takes the fields from src and copies them into the target object.
	//
	// Calling this method with a nil receiver or a nil src may panic.
	CopyFrom(src interface***REMOVED******REMOVED***)
***REMOVED***

// Copy copies src into dst. dst and src must have the same type.
//
// If the type has a copy function defined, it will be used.
//
// Default implementations for builtin types and well known protobuf types may
// be provided.
//
// If the copy cannot be performed, this function will panic. Make sure to test
// types that use this function.
func Copy(dst, src interface***REMOVED******REMOVED***) ***REMOVED***
	switch dst := dst.(type) ***REMOVED***
	case *types.Any:
		src := src.(*types.Any)
		dst.TypeUrl = src.TypeUrl
		if src.Value != nil ***REMOVED***
			dst.Value = make([]byte, len(src.Value))
			copy(dst.Value, src.Value)
		***REMOVED*** else ***REMOVED***
			dst.Value = nil
		***REMOVED***
	case *types.Duration:
		src := src.(*types.Duration)
		*dst = *src
	case *time.Duration:
		src := src.(*time.Duration)
		*dst = *src
	case *types.Timestamp:
		src := src.(*types.Timestamp)
		*dst = *src
	case CopierFrom:
		dst.CopyFrom(src)
	default:
		panic(fmt.Sprintf("Copy for %T not implemented", dst))
	***REMOVED***

***REMOVED***
