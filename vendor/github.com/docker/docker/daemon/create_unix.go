// +build !windows

package daemon

import (
	"fmt"
	"os"
	"path/filepath"

	containertypes "github.com/docker/docker/api/types/container"
	mounttypes "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/stringid"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/sirupsen/logrus"
)

// createContainerOSSpecificSettings performs host-OS specific container create functionality
func (daemon *Daemon) createContainerOSSpecificSettings(container *container.Container, config *containertypes.Config, hostConfig *containertypes.HostConfig) error ***REMOVED***
	if err := daemon.Mount(container); err != nil ***REMOVED***
		return err
	***REMOVED***
	defer daemon.Unmount(container)

	rootIDs := daemon.idMappings.RootPair()
	if err := container.SetupWorkingDirectory(rootIDs); err != nil ***REMOVED***
		return err
	***REMOVED***

	for spec := range config.Volumes ***REMOVED***
		name := stringid.GenerateNonCryptoID()
		destination := filepath.Clean(spec)

		// Skip volumes for which we already have something mounted on that
		// destination because of a --volume-from.
		if container.IsDestinationMounted(destination) ***REMOVED***
			continue
		***REMOVED***
		path, err := container.GetResourcePath(destination)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		stat, err := os.Stat(path)
		if err == nil && !stat.IsDir() ***REMOVED***
			return fmt.Errorf("cannot mount volume over existing file, file exists %s", path)
		***REMOVED***

		v, err := daemon.volumes.CreateWithRef(name, hostConfig.VolumeDriver, container.ID, nil, nil)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := label.Relabel(v.Path(), container.MountLabel, true); err != nil ***REMOVED***
			return err
		***REMOVED***

		container.AddMountPointWithVolume(destination, v, true)
	***REMOVED***
	return daemon.populateVolumes(container)
***REMOVED***

// populateVolumes copies data from the container's rootfs into the volume for non-binds.
// this is only called when the container is created.
func (daemon *Daemon) populateVolumes(c *container.Container) error ***REMOVED***
	for _, mnt := range c.MountPoints ***REMOVED***
		if mnt.Volume == nil ***REMOVED***
			continue
		***REMOVED***

		if mnt.Type != mounttypes.TypeVolume || !mnt.CopyData ***REMOVED***
			continue
		***REMOVED***

		logrus.Debugf("copying image data from %s:%s, to %s", c.ID, mnt.Destination, mnt.Name)
		if err := c.CopyImagePathContent(mnt.Volume, mnt.Destination); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
