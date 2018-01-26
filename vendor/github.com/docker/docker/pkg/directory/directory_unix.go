// +build linux freebsd

package directory

import (
	"os"
	"path/filepath"
	"syscall"
)

// Size walks a directory tree and returns its total size in bytes.
func Size(dir string) (size int64, err error) ***REMOVED***
	data := make(map[uint64]struct***REMOVED******REMOVED***)
	err = filepath.Walk(dir, func(d string, fileInfo os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			// if dir does not exist, Size() returns the error.
			// if dir/x disappeared while walking, Size() ignores dir/x.
			if os.IsNotExist(err) && d != dir ***REMOVED***
				return nil
			***REMOVED***
			return err
		***REMOVED***

		// Ignore directory sizes
		if fileInfo == nil ***REMOVED***
			return nil
		***REMOVED***

		s := fileInfo.Size()
		if fileInfo.IsDir() || s == 0 ***REMOVED***
			return nil
		***REMOVED***

		// Check inode to handle hard links correctly
		inode := fileInfo.Sys().(*syscall.Stat_t).Ino
		// inode is not a uint64 on all platforms. Cast it to avoid issues.
		if _, exists := data[inode]; exists ***REMOVED***
			return nil
		***REMOVED***
		// inode is not a uint64 on all platforms. Cast it to avoid issues.
		data[inode] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

		size += s

		return nil
	***REMOVED***)
	return
***REMOVED***
