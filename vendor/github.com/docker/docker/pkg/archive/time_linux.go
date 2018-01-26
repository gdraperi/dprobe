package archive

import (
	"syscall"
	"time"
)

func timeToTimespec(time time.Time) (ts syscall.Timespec) ***REMOVED***
	if time.IsZero() ***REMOVED***
		// Return UTIME_OMIT special value
		ts.Sec = 0
		ts.Nsec = ((1 << 30) - 2)
		return
	***REMOVED***
	return syscall.NsecToTimespec(time.UnixNano())
***REMOVED***
