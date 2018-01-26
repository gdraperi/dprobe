package devices

// from /usr/include/sys/types.h

func getmajor(dev uint32) uint64 ***REMOVED***
	return (uint64(dev) >> 24) & 0xff
***REMOVED***

func getminor(dev uint32) uint64 ***REMOVED***
	return uint64(dev) & 0xffffff
***REMOVED***

func makedev(major int, minor int) int ***REMOVED***
	return ((major << 24) | minor)
***REMOVED***
