package procfs

import (
	"fmt"
	"os"
	"path"
)

// FS represents the pseudo-filesystem proc, which provides an interface to
// kernel data structures.
type FS string

// DefaultMountPoint is the common mount point of the proc filesystem.
const DefaultMountPoint = "/proc"

// NewFS returns a new FS mounted under the given mountPoint. It will error
// if the mount point can't be read.
func NewFS(mountPoint string) (FS, error) ***REMOVED***
	info, err := os.Stat(mountPoint)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("could not read %s: %s", mountPoint, err)
	***REMOVED***
	if !info.IsDir() ***REMOVED***
		return "", fmt.Errorf("mount point %s is not a directory", mountPoint)
	***REMOVED***

	return FS(mountPoint), nil
***REMOVED***

// Path returns the path of the given subsystem relative to the procfs root.
func (fs FS) Path(p ...string) string ***REMOVED***
	return path.Join(append([]string***REMOVED***string(fs)***REMOVED***, p...)...)
***REMOVED***
