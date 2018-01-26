package directory

import (
	"os"
	"path/filepath"
)

// Size walks a directory tree and returns its total size in bytes.
func Size(dir string) (size int64, err error) ***REMOVED***
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

		size += s

		return nil
	***REMOVED***)
	return
***REMOVED***
