package archive

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	minTime = time.Unix(0, 0)
	maxTime time.Time
)

func init() ***REMOVED***
	if unsafe.Sizeof(syscall.Timespec***REMOVED******REMOVED***.Nsec) == 8 ***REMOVED***
		// This is a 64 bit timespec
		// os.Chtimes limits time to the following
		maxTime = time.Unix(0, 1<<63-1)
	***REMOVED*** else ***REMOVED***
		// This is a 32 bit timespec
		maxTime = time.Unix(1<<31-1, 0)
	***REMOVED***
***REMOVED***

func boundTime(t time.Time) time.Time ***REMOVED***
	if t.Before(minTime) || t.After(maxTime) ***REMOVED***
		return minTime
	***REMOVED***

	return t
***REMOVED***

func latestTime(t1, t2 time.Time) time.Time ***REMOVED***
	if t1.Before(t2) ***REMOVED***
		return t2
	***REMOVED***
	return t1
***REMOVED***
