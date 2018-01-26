// +build !windows

package symlink

import (
	"path/filepath"
)

func evalSymlinks(path string) (string, error) ***REMOVED***
	return filepath.EvalSymlinks(path)
***REMOVED***

func isDriveOrRoot(p string) bool ***REMOVED***
	return p == string(filepath.Separator)
***REMOVED***
