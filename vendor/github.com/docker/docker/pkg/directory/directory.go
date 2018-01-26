package directory

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// MoveToSubdir moves all contents of a directory to a subdirectory underneath the original path
func MoveToSubdir(oldpath, subdir string) error ***REMOVED***

	infos, err := ioutil.ReadDir(oldpath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, info := range infos ***REMOVED***
		if info.Name() != subdir ***REMOVED***
			oldName := filepath.Join(oldpath, info.Name())
			newName := filepath.Join(oldpath, subdir, info.Name())
			if err := os.Rename(oldName, newName); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
