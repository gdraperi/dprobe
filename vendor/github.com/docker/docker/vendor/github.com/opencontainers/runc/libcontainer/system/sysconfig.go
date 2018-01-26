// +build cgo,linux

package system

/*
#include <unistd.h>
*/
import "C"

func GetClockTicks() int ***REMOVED***
	return int(C.sysconf(C._SC_CLK_TCK))
***REMOVED***
