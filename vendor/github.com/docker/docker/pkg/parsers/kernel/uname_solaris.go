package kernel

import (
	"golang.org/x/sys/unix"
)

func uname() (*unix.Utsname, error) ***REMOVED***
	uts := &unix.Utsname***REMOVED******REMOVED***

	if err := unix.Uname(uts); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return uts, nil
***REMOVED***
