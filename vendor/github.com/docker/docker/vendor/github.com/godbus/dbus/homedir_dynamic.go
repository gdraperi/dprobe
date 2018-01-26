// +build !static_build

package dbus

import (
	"os/user"
)

func lookupHomeDir() string ***REMOVED***
	u, err := user.Current()
	if err != nil ***REMOVED***
		return "/"
	***REMOVED***
	return u.HomeDir
***REMOVED***
