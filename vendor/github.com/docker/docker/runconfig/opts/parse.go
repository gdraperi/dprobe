package opts

import (
	"strings"
)

// ConvertKVStringsToMap converts ["key=value"] to ***REMOVED***"key":"value"***REMOVED***
func ConvertKVStringsToMap(values []string) map[string]string ***REMOVED***
	result := make(map[string]string, len(values))
	for _, value := range values ***REMOVED***
		kv := strings.SplitN(value, "=", 2)
		if len(kv) == 1 ***REMOVED***
			result[kv[0]] = ""
		***REMOVED*** else ***REMOVED***
			result[kv[0]] = kv[1]
		***REMOVED***
	***REMOVED***

	return result
***REMOVED***
