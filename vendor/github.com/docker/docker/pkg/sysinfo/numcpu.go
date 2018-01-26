// +build !linux,!windows

package sysinfo

import (
	"runtime"
)

// NumCPU returns the number of CPUs
func NumCPU() int ***REMOVED***
	return runtime.NumCPU()
***REMOVED***
