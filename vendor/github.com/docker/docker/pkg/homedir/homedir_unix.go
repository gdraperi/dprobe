// +build !windows

package homedir

import (
	"os"

	"github.com/opencontainers/runc/libcontainer/user"
)

// Key returns the env var name for the user's home dir based on
// the platform being run on
func Key() string ***REMOVED***
	return "HOME"
***REMOVED***

// Get returns the home directory of the current user with the help of
// environment variables depending on the target operating system.
// Returned path should be used with "path/filepath" to form new paths.
func Get() string ***REMOVED***
	home := os.Getenv(Key())
	if home == "" ***REMOVED***
		if u, err := user.CurrentUser(); err == nil ***REMOVED***
			return u.Home
		***REMOVED***
	***REMOVED***
	return home
***REMOVED***

// GetShortcutString returns the string that is shortcut to user's home directory
// in the native shell of the platform running on.
func GetShortcutString() string ***REMOVED***
	return "~"
***REMOVED***
