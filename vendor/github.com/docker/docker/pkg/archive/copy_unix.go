// +build !windows

package archive

import (
	"path/filepath"
)

func normalizePath(path string) string ***REMOVED***
	return filepath.ToSlash(path)
***REMOVED***
