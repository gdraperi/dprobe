// +build !windows

package containerfs

import "path/filepath"

// cleanScopedPath preappends a to combine with a mnt path.
func cleanScopedPath(path string) string ***REMOVED***
	return filepath.Join(string(filepath.Separator), path)
***REMOVED***
