// +build linux,!gccgo

package nsenter

/*
#cgo CFLAGS: -Wall
extern void nsexec();
void __attribute__((constructor)) init(void) ***REMOVED***
	nsexec();
***REMOVED***
*/
import "C"
