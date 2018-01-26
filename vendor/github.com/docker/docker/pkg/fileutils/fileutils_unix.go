// +build linux freebsd

package fileutils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// GetTotalUsedFds Returns the number of used File Descriptors by
// reading it via /proc filesystem.
func GetTotalUsedFds() int ***REMOVED***
	if fds, err := ioutil.ReadDir(fmt.Sprintf("/proc/%d/fd", os.Getpid())); err != nil ***REMOVED***
		logrus.Errorf("Error opening /proc/%d/fd: %s", os.Getpid(), err)
	***REMOVED*** else ***REMOVED***
		return len(fds)
	***REMOVED***
	return -1
***REMOVED***
