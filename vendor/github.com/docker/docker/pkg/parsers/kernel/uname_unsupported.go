// +build !linux

package kernel

import (
	"errors"
)

// Utsname represents the system name structure.
// It is defined here to make it portable as it is available on linux but not
// on windows.
type Utsname struct ***REMOVED***
	Release [65]byte
***REMOVED***

func uname() (*Utsname, error) ***REMOVED***
	return nil, errors.New("Kernel version detection is available only on linux")
***REMOVED***
