// +build !windows

// Package kernel provides helper function to get, parse and compare kernel
// versions for different platforms.
package kernel

import (
	"errors"
	"fmt"
)

// VersionInfo holds information about the kernel.
type VersionInfo struct ***REMOVED***
	Kernel int    // Version of the kernel (e.g. 4.1.2-generic -> 4)
	Major  int    // Major part of the kernel version (e.g. 4.1.2-generic -> 1)
	Minor  int    // Minor part of the kernel version (e.g. 4.1.2-generic -> 2)
	Flavor string // Flavor of the kernel version (e.g. 4.1.2-generic -> generic)
***REMOVED***

func (k *VersionInfo) String() string ***REMOVED***
	return fmt.Sprintf("%d.%d.%d%s", k.Kernel, k.Major, k.Minor, k.Flavor)
***REMOVED***

// CompareKernelVersion compares two kernel.VersionInfo structs.
// Returns -1 if a < b, 0 if a == b, 1 it a > b
func CompareKernelVersion(a, b VersionInfo) int ***REMOVED***
	if a.Kernel < b.Kernel ***REMOVED***
		return -1
	***REMOVED*** else if a.Kernel > b.Kernel ***REMOVED***
		return 1
	***REMOVED***

	if a.Major < b.Major ***REMOVED***
		return -1
	***REMOVED*** else if a.Major > b.Major ***REMOVED***
		return 1
	***REMOVED***

	if a.Minor < b.Minor ***REMOVED***
		return -1
	***REMOVED*** else if a.Minor > b.Minor ***REMOVED***
		return 1
	***REMOVED***

	return 0
***REMOVED***

// ParseRelease parses a string and creates a VersionInfo based on it.
func ParseRelease(release string) (*VersionInfo, error) ***REMOVED***
	var (
		kernel, major, minor, parsed int
		flavor, partial              string
	)

	// Ignore error from Sscanf to allow an empty flavor.  Instead, just
	// make sure we got all the version numbers.
	parsed, _ = fmt.Sscanf(release, "%d.%d%s", &kernel, &major, &partial)
	if parsed < 2 ***REMOVED***
		return nil, errors.New("Can't parse kernel version " + release)
	***REMOVED***

	// sometimes we have 3.12.25-gentoo, but sometimes we just have 3.12-1-amd64
	parsed, _ = fmt.Sscanf(partial, ".%d%s", &minor, &flavor)
	if parsed < 1 ***REMOVED***
		flavor = partial
	***REMOVED***

	return &VersionInfo***REMOVED***
		Kernel: kernel,
		Major:  major,
		Minor:  minor,
		Flavor: flavor,
	***REMOVED***, nil
***REMOVED***
