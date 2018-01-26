// +build !linux

package archive

import (
	"syscall"
	"time"
)

func timeToTimespec(time time.Time) (ts syscall.Timespec) ***REMOVED***
	nsec := int64(0)
	if !time.IsZero() ***REMOVED***
		nsec = time.UnixNano()
	***REMOVED***
	return syscall.NsecToTimespec(nsec)
***REMOVED***
