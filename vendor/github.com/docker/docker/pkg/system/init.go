package system

import (
	"syscall"
	"time"
	"unsafe"
)

// Used by chtimes
var maxTime time.Time

func init() ***REMOVED***
	// chtimes initialization
	if unsafe.Sizeof(syscall.Timespec***REMOVED******REMOVED***.Nsec) == 8 ***REMOVED***
		// This is a 64 bit timespec
		// os.Chtimes limits time to the following
		maxTime = time.Unix(0, 1<<63-1)
	***REMOVED*** else ***REMOVED***
		// This is a 32 bit timespec
		maxTime = time.Unix(1<<31-1, 0)
	***REMOVED***
***REMOVED***
