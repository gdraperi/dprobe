// +build !windows

package system

import (
	"syscall"
)

// Lstat takes a path to a file and returns
// a system.StatT type pertaining to that file.
//
// Throws an error if the file does not exist
func Lstat(path string) (*StatT, error) ***REMOVED***
	s := &syscall.Stat_t***REMOVED******REMOVED***
	if err := syscall.Lstat(path, s); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return fromStatT(s)
***REMOVED***
