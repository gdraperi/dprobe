package errcode

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
)

var (
	errorCodeToDescriptors = map[ErrorCode]ErrorDescriptor***REMOVED******REMOVED***
	idToDescriptors        = map[string]ErrorDescriptor***REMOVED******REMOVED***
	groupToDescriptors     = map[string][]ErrorDescriptor***REMOVED******REMOVED***
)

var (
	// ErrorCodeUnknown is a generic error that can be used as a last
	// resort if there is no situation-specific error message that can be used
	ErrorCodeUnknown = Register("errcode", ErrorDescriptor***REMOVED***
		Value:   "UNKNOWN",
		Message: "unknown error",
		Description: `Generic error returned when the error does not have an
			                                            API classification.`,
		HTTPStatusCode: http.StatusInternalServerError,
	***REMOVED***)

	// ErrorCodeUnsupported is returned when an operation is not supported.
	ErrorCodeUnsupported = Register("errcode", ErrorDescriptor***REMOVED***
		Value:   "UNSUPPORTED",
		Message: "The operation is unsupported.",
		Description: `The operation was unsupported due to a missing
		implementation or invalid set of parameters.`,
		HTTPStatusCode: http.StatusMethodNotAllowed,
	***REMOVED***)

	// ErrorCodeUnauthorized is returned if a request requires
	// authentication.
	ErrorCodeUnauthorized = Register("errcode", ErrorDescriptor***REMOVED***
		Value:   "UNAUTHORIZED",
		Message: "authentication required",
		Description: `The access controller was unable to authenticate
		the client. Often this will be accompanied by a
		Www-Authenticate HTTP response header indicating how to
		authenticate.`,
		HTTPStatusCode: http.StatusUnauthorized,
	***REMOVED***)

	// ErrorCodeDenied is returned if a client does not have sufficient
	// permission to perform an action.
	ErrorCodeDenied = Register("errcode", ErrorDescriptor***REMOVED***
		Value:   "DENIED",
		Message: "requested access to the resource is denied",
		Description: `The access controller denied access for the
		operation on a resource.`,
		HTTPStatusCode: http.StatusForbidden,
	***REMOVED***)

	// ErrorCodeUnavailable provides a common error to report unavailability
	// of a service or endpoint.
	ErrorCodeUnavailable = Register("errcode", ErrorDescriptor***REMOVED***
		Value:          "UNAVAILABLE",
		Message:        "service unavailable",
		Description:    "Returned when a service is not available",
		HTTPStatusCode: http.StatusServiceUnavailable,
	***REMOVED***)

	// ErrorCodeTooManyRequests is returned if a client attempts too many
	// times to contact a service endpoint.
	ErrorCodeTooManyRequests = Register("errcode", ErrorDescriptor***REMOVED***
		Value:   "TOOMANYREQUESTS",
		Message: "too many requests",
		Description: `Returned when a client attempts to contact a
		service too many times`,
		HTTPStatusCode: http.StatusTooManyRequests,
	***REMOVED***)
)

var nextCode = 1000
var registerLock sync.Mutex

// Register will make the passed-in error known to the environment and
// return a new ErrorCode
func Register(group string, descriptor ErrorDescriptor) ErrorCode ***REMOVED***
	registerLock.Lock()
	defer registerLock.Unlock()

	descriptor.Code = ErrorCode(nextCode)

	if _, ok := idToDescriptors[descriptor.Value]; ok ***REMOVED***
		panic(fmt.Sprintf("ErrorValue %q is already registered", descriptor.Value))
	***REMOVED***
	if _, ok := errorCodeToDescriptors[descriptor.Code]; ok ***REMOVED***
		panic(fmt.Sprintf("ErrorCode %v is already registered", descriptor.Code))
	***REMOVED***

	groupToDescriptors[group] = append(groupToDescriptors[group], descriptor)
	errorCodeToDescriptors[descriptor.Code] = descriptor
	idToDescriptors[descriptor.Value] = descriptor

	nextCode++
	return descriptor.Code
***REMOVED***

type byValue []ErrorDescriptor

func (a byValue) Len() int           ***REMOVED*** return len(a) ***REMOVED***
func (a byValue) Swap(i, j int)      ***REMOVED*** a[i], a[j] = a[j], a[i] ***REMOVED***
func (a byValue) Less(i, j int) bool ***REMOVED*** return a[i].Value < a[j].Value ***REMOVED***

// GetGroupNames returns the list of Error group names that are registered
func GetGroupNames() []string ***REMOVED***
	keys := []string***REMOVED******REMOVED***

	for k := range groupToDescriptors ***REMOVED***
		keys = append(keys, k)
	***REMOVED***
	sort.Strings(keys)
	return keys
***REMOVED***

// GetErrorCodeGroup returns the named group of error descriptors
func GetErrorCodeGroup(name string) []ErrorDescriptor ***REMOVED***
	desc := groupToDescriptors[name]
	sort.Sort(byValue(desc))
	return desc
***REMOVED***

// GetErrorAllDescriptors returns a slice of all ErrorDescriptors that are
// registered, irrespective of what group they're in
func GetErrorAllDescriptors() []ErrorDescriptor ***REMOVED***
	result := []ErrorDescriptor***REMOVED******REMOVED***

	for _, group := range GetGroupNames() ***REMOVED***
		result = append(result, GetErrorCodeGroup(group)...)
	***REMOVED***
	sort.Sort(byValue(result))
	return result
***REMOVED***
