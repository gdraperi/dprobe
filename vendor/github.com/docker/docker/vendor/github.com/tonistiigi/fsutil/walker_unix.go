// +build !windows

package fsutil

import (
	"os"
	"syscall"

	"github.com/containerd/continuity/sysx"
	"github.com/pkg/errors"
)

func loadXattr(origpath string, stat *Stat) error ***REMOVED***
	xattrs, err := sysx.LListxattr(origpath)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to xattr %s", origpath)
	***REMOVED***
	if len(xattrs) > 0 ***REMOVED***
		m := make(map[string][]byte)
		for _, key := range xattrs ***REMOVED***
			v, err := sysx.LGetxattr(origpath, key)
			if err == nil ***REMOVED***
				m[key] = v
			***REMOVED***
		***REMOVED***
		stat.Xattrs = m
	***REMOVED***
	return nil
***REMOVED***

func setUnixOpt(fi os.FileInfo, stat *Stat, path string, seenFiles map[uint64]string) ***REMOVED***
	s := fi.Sys().(*syscall.Stat_t)

	stat.Uid = s.Uid
	stat.Gid = s.Gid

	if !fi.IsDir() ***REMOVED***
		if s.Mode&syscall.S_IFBLK != 0 ||
			s.Mode&syscall.S_IFCHR != 0 ***REMOVED***
			stat.Devmajor = int64(major(uint64(s.Rdev)))
			stat.Devminor = int64(minor(uint64(s.Rdev)))
		***REMOVED***

		ino := s.Ino
		if s.Nlink > 1 ***REMOVED***
			if oldpath, ok := seenFiles[ino]; ok ***REMOVED***
				stat.Linkname = oldpath
				stat.Size_ = 0
			***REMOVED***
		***REMOVED***
		seenFiles[ino] = path
	***REMOVED***
***REMOVED***

func major(device uint64) uint64 ***REMOVED***
	return (device >> 8) & 0xfff
***REMOVED***

func minor(device uint64) uint64 ***REMOVED***
	return (device & 0xff) | ((device >> 12) & 0xfff00)
***REMOVED***
