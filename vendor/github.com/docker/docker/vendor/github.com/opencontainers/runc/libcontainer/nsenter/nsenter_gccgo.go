// +build linux,gccgo

package nsenter

/*
#cgo CFLAGS: -Wall
extern void nsexec();
void __attribute__((constructor)) init(void) ***REMOVED***
	nsexec();
***REMOVED***
*/
import "C"

// AlwaysFalse is here to stay false
// (and be exported so the compiler doesn't optimize out its reference)
var AlwaysFalse bool

func init() ***REMOVED***
	if AlwaysFalse ***REMOVED***
		// by referencing this C init() in a noop test, it will ensure the compiler
		// links in the C function.
		// https://gcc.gnu.org/bugzilla/show_bug.cgi?id=65134
		C.init()
	***REMOVED***
***REMOVED***
