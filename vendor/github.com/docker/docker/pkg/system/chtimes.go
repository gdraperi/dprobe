package system

import (
	"os"
	"time"
)

// Chtimes changes the access time and modified time of a file at the given path
func Chtimes(name string, atime time.Time, mtime time.Time) error ***REMOVED***
	unixMinTime := time.Unix(0, 0)
	unixMaxTime := maxTime

	// If the modified time is prior to the Unix Epoch, or after the
	// end of Unix Time, os.Chtimes has undefined behavior
	// default to Unix Epoch in this case, just in case

	if atime.Before(unixMinTime) || atime.After(unixMaxTime) ***REMOVED***
		atime = unixMinTime
	***REMOVED***

	if mtime.Before(unixMinTime) || mtime.After(unixMaxTime) ***REMOVED***
		mtime = unixMinTime
	***REMOVED***

	if err := os.Chtimes(name, atime, mtime); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Take platform specific action for setting create time.
	return setCTime(name, mtime)
***REMOVED***
