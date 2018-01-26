package format

import "fmt"

// Message accepts a msgAndArgs varargs and formats it using fmt.Sprintf
func Message(msgAndArgs ...interface***REMOVED******REMOVED***) string ***REMOVED***
	switch len(msgAndArgs) ***REMOVED***
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%v", msgAndArgs[0])
	default:
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	***REMOVED***
***REMOVED***

// WithCustomMessage accepts one or two messages and formats them appropriately
func WithCustomMessage(source string, msgAndArgs ...interface***REMOVED******REMOVED***) string ***REMOVED***
	custom := Message(msgAndArgs...)
	switch ***REMOVED***
	case custom == "":
		return source
	case source == "":
		return custom
	***REMOVED***
	return fmt.Sprintf("%s: %s", source, custom)
***REMOVED***
