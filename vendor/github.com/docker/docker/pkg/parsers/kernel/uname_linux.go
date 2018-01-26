package kernel

import "golang.org/x/sys/unix"

// Utsname represents the system name structure.
// It is passthrough for unix.Utsname in order to make it portable with
// other platforms where it is not available.
type Utsname unix.Utsname

func uname() (*unix.Utsname, error) ***REMOVED***
	uts := &unix.Utsname***REMOVED******REMOVED***

	if err := unix.Uname(uts); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return uts, nil
***REMOVED***
