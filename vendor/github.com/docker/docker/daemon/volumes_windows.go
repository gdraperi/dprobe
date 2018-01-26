package daemon

import (
	"sort"

	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/volume"
)

// setupMounts configures the mount points for a container by appending each
// of the configured mounts on the container to the OCI mount structure
// which will ultimately be passed into the oci runtime during container creation.
// It also ensures each of the mounts are lexicographically sorted.

// BUGBUG TODO Windows containerd. This would be much better if it returned
// an array of runtime spec mounts, not container mounts. Then no need to
// do multiple transitions.

func (daemon *Daemon) setupMounts(c *container.Container) ([]container.Mount, error) ***REMOVED***
	var mnts []container.Mount
	for _, mount := range c.MountPoints ***REMOVED*** // type is volume.MountPoint
		if err := daemon.lazyInitializeVolume(c.ID, mount); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		s, err := mount.Setup(c.MountLabel, idtools.IDPair***REMOVED***0, 0***REMOVED***, nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		mnts = append(mnts, container.Mount***REMOVED***
			Source:      s,
			Destination: mount.Destination,
			Writable:    mount.RW,
		***REMOVED***)
	***REMOVED***

	sort.Sort(mounts(mnts))
	return mnts, nil
***REMOVED***

// setBindModeIfNull is platform specific processing which is a no-op on
// Windows.
func setBindModeIfNull(bind *volume.MountPoint) ***REMOVED***
	return
***REMOVED***
