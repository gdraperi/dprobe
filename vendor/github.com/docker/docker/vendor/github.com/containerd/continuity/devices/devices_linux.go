package devices

// from /usr/include/linux/kdev_t.h

func getmajor(dev uint64) uint64 ***REMOVED***
	return dev >> 8
***REMOVED***

func getminor(dev uint64) uint64 ***REMOVED***
	return dev & 0xff
***REMOVED***

func makedev(major int, minor int) int ***REMOVED***
	return ((major << 8) | minor)
***REMOVED***
