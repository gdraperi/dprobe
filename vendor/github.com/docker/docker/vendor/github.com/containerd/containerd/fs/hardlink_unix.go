// +build !windows

package fs

import (
	"os"
	"syscall"
)

func getLinkInfo(fi os.FileInfo) (uint64, bool) ***REMOVED***
	s, ok := fi.Sys().(*syscall.Stat_t)
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***

	return uint64(s.Ino), !fi.IsDir() && s.Nlink > 1
***REMOVED***
