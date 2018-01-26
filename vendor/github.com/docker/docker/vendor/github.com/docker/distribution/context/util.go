package context

import (
	"time"
)

// Since looks up key, which should be a time.Time, and returns the duration
// since that time. If the key is not found, the value returned will be zero.
// This is helpful when inferring metrics related to context execution times.
func Since(ctx Context, key interface***REMOVED******REMOVED***) time.Duration ***REMOVED***
	if startedAt, ok := ctx.Value(key).(time.Time); ok ***REMOVED***
		return time.Since(startedAt)
	***REMOVED***
	return 0
***REMOVED***

// GetStringValue returns a string value from the context. The empty string
// will be returned if not found.
func GetStringValue(ctx Context, key interface***REMOVED******REMOVED***) (value string) ***REMOVED***
	if valuev, ok := ctx.Value(key).(string); ok ***REMOVED***
		value = valuev
	***REMOVED***
	return value
***REMOVED***
