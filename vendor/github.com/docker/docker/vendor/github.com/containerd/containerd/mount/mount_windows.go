package mount

import "github.com/pkg/errors"

var (
	// ErrNotImplementOnWindows is returned when an action is not implemented for windows
	ErrNotImplementOnWindows = errors.New("not implemented under windows")
)

// Mount to the provided target
func (m *Mount) Mount(target string) error ***REMOVED***
	return ErrNotImplementOnWindows
***REMOVED***

// Unmount the mount at the provided path
func Unmount(mount string, flags int) error ***REMOVED***
	return ErrNotImplementOnWindows
***REMOVED***

// UnmountAll mounts at the provided path
func UnmountAll(mount string, flags int) error ***REMOVED***
	return ErrNotImplementOnWindows
***REMOVED***
