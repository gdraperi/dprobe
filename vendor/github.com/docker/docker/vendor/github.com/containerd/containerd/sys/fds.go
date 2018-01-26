// +build !windows,!darwin

package sys

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
)

// GetOpenFds returns the number of open fds for the process provided by pid
func GetOpenFds(pid int) (int, error) ***REMOVED***
	dirs, err := ioutil.ReadDir(filepath.Join("/proc", strconv.Itoa(pid), "fd"))
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	return len(dirs), nil
***REMOVED***
