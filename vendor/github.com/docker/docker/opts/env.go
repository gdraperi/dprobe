package opts

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// ValidateEnv validates an environment variable and returns it.
// If no value is specified, it returns the current value using os.Getenv.
//
// As on ParseEnvFile and related to #16585, environment variable names
// are not validate what so ever, it's up to application inside docker
// to validate them or not.
//
// The only validation here is to check if name is empty, per #25099
func ValidateEnv(val string) (string, error) ***REMOVED***
	arr := strings.Split(val, "=")
	if arr[0] == "" ***REMOVED***
		return "", errors.Errorf("invalid environment variable: %s", val)
	***REMOVED***
	if len(arr) > 1 ***REMOVED***
		return val, nil
	***REMOVED***
	if !doesEnvExist(val) ***REMOVED***
		return val, nil
	***REMOVED***
	return fmt.Sprintf("%s=%s", val, os.Getenv(val)), nil
***REMOVED***

func doesEnvExist(name string) bool ***REMOVED***
	for _, entry := range os.Environ() ***REMOVED***
		parts := strings.SplitN(entry, "=", 2)
		if runtime.GOOS == "windows" ***REMOVED***
			// Environment variable are case-insensitive on Windows. PaTh, path and PATH are equivalent.
			if strings.EqualFold(parts[0], name) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
		if parts[0] == name ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
