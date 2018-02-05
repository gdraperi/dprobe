package testhelper

import (
	"runtime"
	"strings"
)

// Converts the line endings when on Windows
func Unix2dos(unix string) string ***REMOVED***
	if runtime.GOOS != "windows" ***REMOVED***
		return unix
	***REMOVED***

	return strings.Replace(unix, "\n", "\r\n", -1)
***REMOVED***
