// +build !windows

package daemon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/volume"
	"github.com/docker/docker/volume/drivers"
	"github.com/docker/docker/volume/local"
	"github.com/pkg/errors"
)

// setupMounts iterates through each of the mount points for a container and
// calls Setup() on each. It also looks to see if is a network mount such as
// /etc/resolv.conf, and if it is not, appends it to the array of mounts.
func (daemon *Daemon) setupMounts(c *container.Container) ([]container.Mount, error) ***REMOVED***
	var mounts []container.Mount
	// TODO: tmpfs mounts should be part of Mountpoints
	tmpfsMounts := make(map[string]bool)
	tmpfsMountInfo, err := c.TmpfsMounts()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, m := range tmpfsMountInfo ***REMOVED***
		tmpfsMounts[m.Destination] = true
	***REMOVED***
	for _, m := range c.MountPoints ***REMOVED***
		if tmpfsMounts[m.Destination] ***REMOVED***
			continue
		***REMOVED***
		if err := daemon.lazyInitializeVolume(c.ID, m); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// If the daemon is being shutdown, we should not let a container start if it is trying to
		// mount the socket the daemon is listening on. During daemon shutdown, the socket
		// (/var/run/docker.sock by default) doesn't exist anymore causing the call to m.Setup to
		// create at directory instead. This in turn will prevent the daemon to restart.
		checkfunc := func(m *volume.MountPoint) error ***REMOVED***
			if _, exist := daemon.hosts[m.Source]; exist && daemon.IsShuttingDown() ***REMOVED***
				return fmt.Errorf("Could not mount %q to container while the daemon is shutting down", m.Source)
			***REMOVED***
			return nil
		***REMOVED***

		path, err := m.Setup(c.MountLabel, daemon.idMappings.RootPair(), checkfunc)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if !c.TrySetNetworkMount(m.Destination, path) ***REMOVED***
			mnt := container.Mount***REMOVED***
				Source:      path,
				Destination: m.Destination,
				Writable:    m.RW,
				Propagation: string(m.Propagation),
			***REMOVED***
			if m.Volume != nil ***REMOVED***
				attributes := map[string]string***REMOVED***
					"driver":      m.Volume.DriverName(),
					"container":   c.ID,
					"destination": m.Destination,
					"read/write":  strconv.FormatBool(m.RW),
					"propagation": string(m.Propagation),
				***REMOVED***
				daemon.LogVolumeEvent(m.Volume.Name(), "mount", attributes)
			***REMOVED***
			mounts = append(mounts, mnt)
		***REMOVED***
	***REMOVED***

	mounts = sortMounts(mounts)
	netMounts := c.NetworkMounts()
	// if we are going to mount any of the network files from container
	// metadata, the ownership must be set properly for potential container
	// remapped root (user namespaces)
	rootIDs := daemon.idMappings.RootPair()
	for _, mount := range netMounts ***REMOVED***
		// we should only modify ownership of network files within our own container
		// metadata repository. If the user specifies a mount path external, it is
		// up to the user to make sure the file has proper ownership for userns
		if strings.Index(mount.Source, daemon.repository) == 0 ***REMOVED***
			if err := os.Chown(mount.Source, rootIDs.UID, rootIDs.GID); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return append(mounts, netMounts...), nil
***REMOVED***

// sortMounts sorts an array of mounts in lexicographic order. This ensure that
// when mounting, the mounts don't shadow other mounts. For example, if mounting
// /etc and /etc/resolv.conf, /etc/resolv.conf must not be mounted first.
func sortMounts(m []container.Mount) []container.Mount ***REMOVED***
	sort.Sort(mounts(m))
	return m
***REMOVED***

// setBindModeIfNull is platform specific processing to ensure the
// shared mode is set to 'z' if it is null. This is called in the case
// of processing a named volume and not a typical bind.
func setBindModeIfNull(bind *volume.MountPoint) ***REMOVED***
	if bind.Mode == "" ***REMOVED***
		bind.Mode = "z"
	***REMOVED***
***REMOVED***

// migrateVolume links the contents of a volume created pre Docker 1.7
// into the location expected by the local driver.
// It creates a symlink from DOCKER_ROOT/vfs/dir/VOLUME_ID to DOCKER_ROOT/volumes/VOLUME_ID/_container_data.
// It preserves the volume json configuration generated pre Docker 1.7 to be able to
// downgrade from Docker 1.7 to Docker 1.6 without losing volume compatibility.
func migrateVolume(id, vfs string) error ***REMOVED***
	l, err := volumedrivers.GetDriver(volume.DefaultDriverName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	newDataPath := l.(*local.Root).DataPath(id)
	fi, err := os.Stat(newDataPath)
	if err != nil && !os.IsNotExist(err) ***REMOVED***
		return err
	***REMOVED***

	if fi != nil && fi.IsDir() ***REMOVED***
		return nil
	***REMOVED***

	return os.Symlink(vfs, newDataPath)
***REMOVED***

// verifyVolumesInfo ports volumes configured for the containers pre docker 1.7.
// It reads the container configuration and creates valid mount points for the old volumes.
func (daemon *Daemon) verifyVolumesInfo(container *container.Container) error ***REMOVED***
	container.Lock()
	defer container.Unlock()

	// Inspect old structures only when we're upgrading from old versions
	// to versions >= 1.7 and the MountPoints has not been populated with volumes data.
	type volumes struct ***REMOVED***
		Volumes   map[string]string
		VolumesRW map[string]bool
	***REMOVED***
	cfgPath, err := container.ConfigPath()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	f, err := os.Open(cfgPath)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "could not open container config")
	***REMOVED***
	defer f.Close()
	var cv volumes
	if err := json.NewDecoder(f).Decode(&cv); err != nil ***REMOVED***
		return errors.Wrap(err, "could not decode container config")
	***REMOVED***

	if len(container.MountPoints) == 0 && len(cv.Volumes) > 0 ***REMOVED***
		for destination, hostPath := range cv.Volumes ***REMOVED***
			vfsPath := filepath.Join(daemon.root, "vfs", "dir")
			rw := cv.VolumesRW != nil && cv.VolumesRW[destination]

			if strings.HasPrefix(hostPath, vfsPath) ***REMOVED***
				id := filepath.Base(hostPath)
				v, err := daemon.volumes.CreateWithRef(id, volume.DefaultDriverName, container.ID, nil, nil)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := migrateVolume(id, hostPath); err != nil ***REMOVED***
					return err
				***REMOVED***
				container.AddMountPointWithVolume(destination, v, true)
			***REMOVED*** else ***REMOVED*** // Bind mount
				m := volume.MountPoint***REMOVED***Source: hostPath, Destination: destination, RW: rw***REMOVED***
				container.MountPoints[destination] = &m
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) mountVolumes(container *container.Container) error ***REMOVED***
	mounts, err := daemon.setupMounts(container)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, m := range mounts ***REMOVED***
		dest, err := container.GetResourcePath(m.Destination)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		var stat os.FileInfo
		stat, err = os.Stat(m.Source)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err = fileutils.CreateIfNotExists(dest, stat.IsDir()); err != nil ***REMOVED***
			return err
		***REMOVED***

		opts := "rbind,ro"
		if m.Writable ***REMOVED***
			opts = "rbind,rw"
		***REMOVED***

		if err := mount.Mount(m.Source, dest, bindMountType, opts); err != nil ***REMOVED***
			return err
		***REMOVED***

		// mountVolumes() seems to be called for temporary mounts
		// outside the container. Soon these will be unmounted with
		// lazy unmount option and given we have mounted the rbind,
		// all the submounts will propagate if these are shared. If
		// daemon is running in host namespace and has / as shared
		// then these unmounts will propagate and unmount original
		// mount as well. So make all these mounts rprivate.
		// Do not use propagation property of volume as that should
		// apply only when mounting happen inside the container.
		if err := mount.MakeRPrivate(dest); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
