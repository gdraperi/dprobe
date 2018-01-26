// +build !windows

package daemon

import (
	"github.com/docker/docker/container"
	"github.com/docker/docker/volume"
)

// checkIfPathIsInAVolume checks if the path is in a volume. If it is, it
// cannot be in a read-only volume. If it  is not in a volume, the container
// cannot be configured with a read-only rootfs.
func checkIfPathIsInAVolume(container *container.Container, absPath string) (bool, error) ***REMOVED***
	var toVolume bool
	parser := volume.NewParser(container.OS)
	for _, mnt := range container.MountPoints ***REMOVED***
		if toVolume = parser.HasResource(mnt, absPath); toVolume ***REMOVED***
			if mnt.RW ***REMOVED***
				break
			***REMOVED***
			return false, ErrVolumeReadonly
		***REMOVED***
	***REMOVED***
	return toVolume, nil
***REMOVED***

// isOnlineFSOperationPermitted returns an error if an online filesystem operation
// is not permitted.
func (daemon *Daemon) isOnlineFSOperationPermitted(container *container.Container) error ***REMOVED***
	return nil
***REMOVED***
