package homedir

import (
	"os"
)

// Key returns the env var name for the user's home dir based on
// the platform being run on
func Key() string ***REMOVED***
	return "USERPROFILE"
***REMOVED***

// Get returns the home directory of the current user with the help of
// environment variables depending on the target operating system.
// Returned path should be used with "path/filepath" to form new paths.
func Get() string ***REMOVED***
	return os.Getenv(Key())
***REMOVED***

// GetShortcutString returns the string that is shortcut to user's home directory
// in the native shell of the platform running on.
func GetShortcutString() string ***REMOVED***
	return "%USERPROFILE%" // be careful while using in format functions
***REMOVED***
