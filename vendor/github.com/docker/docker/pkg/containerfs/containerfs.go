package containerfs

import (
	"path/filepath"
	"runtime"

	"github.com/containerd/continuity/driver"
	"github.com/containerd/continuity/pathdriver"
	"github.com/docker/docker/pkg/symlink"
)

// ContainerFS is that represents a root file system
type ContainerFS interface ***REMOVED***
	// Path returns the path to the root. Note that this may not exist
	// on the local system, so the continuity operations must be used
	Path() string

	// ResolveScopedPath evaluates the given path scoped to the root.
	// For example, if root=/a, and path=/b/c, then this function would return /a/b/c.
	// If rawPath is true, then the function will not preform any modifications
	// before path resolution. Otherwise, the function will clean the given path
	// by making it an absolute path.
	ResolveScopedPath(path string, rawPath bool) (string, error)

	Driver
***REMOVED***

// Driver combines both continuity's Driver and PathDriver interfaces with a Platform
// field to determine the OS.
type Driver interface ***REMOVED***
	// OS returns the OS where the rootfs is located. Essentially,
	// runtime.GOOS for everything aside from LCOW, which is "linux"
	OS() string

	// Architecture returns the hardware architecture where the
	// container is located.
	Architecture() string

	// Driver & PathDriver provide methods to manipulate files & paths
	driver.Driver
	pathdriver.PathDriver
***REMOVED***

// NewLocalContainerFS is a helper function to implement daemon's Mount interface
// when the graphdriver mount point is a local path on the machine.
func NewLocalContainerFS(path string) ContainerFS ***REMOVED***
	return &local***REMOVED***
		path:       path,
		Driver:     driver.LocalDriver,
		PathDriver: pathdriver.LocalPathDriver,
	***REMOVED***
***REMOVED***

// NewLocalDriver provides file and path drivers for a local file system. They are
// essentially a wrapper around the `os` and `filepath` functions.
func NewLocalDriver() Driver ***REMOVED***
	return &local***REMOVED***
		Driver:     driver.LocalDriver,
		PathDriver: pathdriver.LocalPathDriver,
	***REMOVED***
***REMOVED***

type local struct ***REMOVED***
	path string
	driver.Driver
	pathdriver.PathDriver
***REMOVED***

func (l *local) Path() string ***REMOVED***
	return l.path
***REMOVED***

func (l *local) ResolveScopedPath(path string, rawPath bool) (string, error) ***REMOVED***
	cleanedPath := path
	if !rawPath ***REMOVED***
		cleanedPath = cleanScopedPath(path)
	***REMOVED***
	return symlink.FollowSymlinkInScope(filepath.Join(l.path, cleanedPath), l.path)
***REMOVED***

func (l *local) OS() string ***REMOVED***
	return runtime.GOOS
***REMOVED***

func (l *local) Architecture() string ***REMOVED***
	return runtime.GOARCH
***REMOVED***
