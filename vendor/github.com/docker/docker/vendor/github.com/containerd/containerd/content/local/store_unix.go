// +build linux solaris darwin freebsd

package local

import (
	"os"
	"syscall"
	"time"

	"github.com/containerd/containerd/sys"
)

func getATime(fi os.FileInfo) time.Time ***REMOVED***
	if st, ok := fi.Sys().(*syscall.Stat_t); ok ***REMOVED***
		return time.Unix(int64(sys.StatAtime(st).Sec),
			int64(sys.StatAtime(st).Nsec))
	***REMOVED***

	return fi.ModTime()
***REMOVED***
