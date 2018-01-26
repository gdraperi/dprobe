// Package platform provides helper function to get the runtime architecture
// for different platforms.
package platform

import (
	"bytes"

	"golang.org/x/sys/unix"
)

// runtimeArchitecture gets the name of the current architecture (x86, x86_64, …)
func runtimeArchitecture() (string, error) ***REMOVED***
	utsname := &unix.Utsname***REMOVED******REMOVED***
	if err := unix.Uname(utsname); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return string(utsname.Machine[:bytes.IndexByte(utsname.Machine[:], 0)]), nil
***REMOVED***
