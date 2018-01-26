// +build !windows

package archive

import (
	"os"
	"syscall"

	"github.com/docker/docker/pkg/system"
	"golang.org/x/sys/unix"
)

func statDifferent(oldStat *system.StatT, newStat *system.StatT) bool ***REMOVED***
	// Don't look at size for dirs, its not a good measure of change
	if oldStat.Mode() != newStat.Mode() ||
		oldStat.UID() != newStat.UID() ||
		oldStat.GID() != newStat.GID() ||
		oldStat.Rdev() != newStat.Rdev() ||
		// Don't look at size for dirs, its not a good measure of change
		(oldStat.Mode()&unix.S_IFDIR != unix.S_IFDIR &&
			(!sameFsTimeSpec(oldStat.Mtim(), newStat.Mtim()) || (oldStat.Size() != newStat.Size()))) ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (info *FileInfo) isDir() bool ***REMOVED***
	return info.parent == nil || info.stat.Mode()&unix.S_IFDIR != 0
***REMOVED***

func getIno(fi os.FileInfo) uint64 ***REMOVED***
	return fi.Sys().(*syscall.Stat_t).Ino
***REMOVED***

func hasHardlinks(fi os.FileInfo) bool ***REMOVED***
	return fi.Sys().(*syscall.Stat_t).Nlink > 1
***REMOVED***
