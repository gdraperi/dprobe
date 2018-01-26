// +build linux darwin freebsd solaris

package devices

import (
	"fmt"
	"os"
	"syscall"
)

func DeviceInfo(fi os.FileInfo) (uint64, uint64, error) ***REMOVED***
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok ***REMOVED***
		return 0, 0, fmt.Errorf("cannot extract device from os.FileInfo")
	***REMOVED***

	return getmajor(sys.Rdev), getminor(sys.Rdev), nil
***REMOVED***

// mknod provides a shortcut for syscall.Mknod
func Mknod(p string, mode os.FileMode, maj, min int) error ***REMOVED***
	var (
		m   = syscallMode(mode.Perm())
		dev int
	)

	if mode&os.ModeDevice != 0 ***REMOVED***
		dev = makedev(maj, min)

		if mode&os.ModeCharDevice != 0 ***REMOVED***
			m |= syscall.S_IFCHR
		***REMOVED*** else ***REMOVED***
			m |= syscall.S_IFBLK
		***REMOVED***
	***REMOVED*** else if mode&os.ModeNamedPipe != 0 ***REMOVED***
		m |= syscall.S_IFIFO
	***REMOVED***

	return syscall.Mknod(p, m, dev)
***REMOVED***

// syscallMode returns the syscall-specific mode bits from Go's portable mode bits.
func syscallMode(i os.FileMode) (o uint32) ***REMOVED***
	o |= uint32(i.Perm())
	if i&os.ModeSetuid != 0 ***REMOVED***
		o |= syscall.S_ISUID
	***REMOVED***
	if i&os.ModeSetgid != 0 ***REMOVED***
		o |= syscall.S_ISGID
	***REMOVED***
	if i&os.ModeSticky != 0 ***REMOVED***
		o |= syscall.S_ISVTX
	***REMOVED***
	return
***REMOVED***
