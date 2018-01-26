// +build !windows,!darwin

package pidfile

import (
	"os"
	"path/filepath"
	"strconv"
)

func processExists(pid int) bool ***REMOVED***
	if _, err := os.Stat(filepath.Join("/proc", strconv.Itoa(pid))); err == nil ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***
