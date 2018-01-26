// +build windows

package windowsconsole

import (
	"os"

	"github.com/Azure/go-ansiterm/winterm"
)

// GetHandleInfo returns file descriptor and bool indicating whether the file is a console.
func GetHandleInfo(in interface***REMOVED******REMOVED***) (uintptr, bool) ***REMOVED***
	switch t := in.(type) ***REMOVED***
	case *ansiReader:
		return t.Fd(), true
	case *ansiWriter:
		return t.Fd(), true
	***REMOVED***

	var inFd uintptr
	var isTerminal bool

	if file, ok := in.(*os.File); ok ***REMOVED***
		inFd = file.Fd()
		isTerminal = IsConsole(inFd)
	***REMOVED***
	return inFd, isTerminal
***REMOVED***

// IsConsole returns true if the given file descriptor is a Windows Console.
// The code assumes that GetConsoleMode will return an error for file descriptors that are not a console.
func IsConsole(fd uintptr) bool ***REMOVED***
	_, e := winterm.GetConsoleMode(fd)
	return e == nil
***REMOVED***
