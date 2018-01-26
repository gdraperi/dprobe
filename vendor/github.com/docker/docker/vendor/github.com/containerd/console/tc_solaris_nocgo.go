// +build solaris,!cgo

//
// Implementing the functions below requires cgo support.  Non-cgo stubs
// versions are defined below to enable cross-compilation of source code
// that depends on these functions, but the resultant cross-compiled
// binaries cannot actually be used.  If the stub function(s) below are
// actually invoked they will display an error message and cause the
// calling process to exit.
//

package console

import (
	"os"

	"golang.org/x/sys/unix"
)

const (
	cmdTcGet = unix.TCGETS
	cmdTcSet = unix.TCSETS
)

func ptsname(f *os.File) (string, error) ***REMOVED***
	panic("ptsname() support requires cgo.")
***REMOVED***

func unlockpt(f *os.File) error ***REMOVED***
	panic("unlockpt() support requires cgo.")
***REMOVED***
