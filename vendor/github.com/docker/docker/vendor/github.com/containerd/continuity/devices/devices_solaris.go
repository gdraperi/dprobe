// +build cgo

package devices

//#include <sys/mkdev.h>
import "C"

func getmajor(dev uint64) uint64 ***REMOVED***
	return uint64(C.major(C.dev_t(dev)))
***REMOVED***

func getminor(dev uint64) uint64 ***REMOVED***
	return uint64(C.minor(C.dev_t(dev)))
***REMOVED***

func makedev(major int, minor int) int ***REMOVED***
	return int(C.makedev(C.major_t(major), C.minor_t(minor)))
***REMOVED***
