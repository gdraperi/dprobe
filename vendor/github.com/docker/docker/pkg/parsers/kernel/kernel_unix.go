// +build linux freebsd openbsd

// Package kernel provides helper function to get, parse and compare kernel
// versions for different platforms.
package kernel

import (
	"bytes"

	"github.com/sirupsen/logrus"
)

// GetKernelVersion gets the current kernel version.
func GetKernelVersion() (*VersionInfo, error) ***REMOVED***
	uts, err := uname()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Remove the \x00 from the release for Atoi to parse correctly
	return ParseRelease(string(uts.Release[:bytes.IndexByte(uts.Release[:], 0)]))
***REMOVED***

// CheckKernelVersion checks if current kernel is newer than (or equal to)
// the given version.
func CheckKernelVersion(k, major, minor int) bool ***REMOVED***
	if v, err := GetKernelVersion(); err != nil ***REMOVED***
		logrus.Warnf("error getting kernel version: %s", err)
	***REMOVED*** else ***REMOVED***
		if CompareKernelVersion(*v, VersionInfo***REMOVED***Kernel: k, Major: major, Minor: minor***REMOVED***) < 0 ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
