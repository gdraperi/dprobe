package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	mounttypes "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/volume"
	"github.com/docker/docker/volume/drivers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	// ErrVolumeReadonly is used to signal an error when trying to copy data into
	// a volume mount that is not writable.
	ErrVolumeReadonly = errors.New("mounted volume is marked read-only")
)

type mounts []container.Mount

// volumeToAPIType converts a volume.Volume to the type used by the Engine API
func volumeToAPIType(v volume.Volume) *types.Volume ***REMOVED***
	createdAt, _ := v.CreatedAt()
	tv := &types.Volume***REMOVED***
		Name:      v.Name(),
		Driver:    v.DriverName(),
		CreatedAt: createdAt.Format(time.RFC3339),
	***REMOVED***
	if v, ok := v.(volume.DetailedVolume); ok ***REMOVED***
		tv.Labels = v.Labels()
		tv.Options = v.Options()
		tv.Scope = v.Scope()
	***REMOVED***

	return tv
***REMOVED***

// Len returns the number of mounts. Used in sorting.
func (m mounts) Len() int ***REMOVED***
	return len(m)
***REMOVED***

// Less returns true if the number of parts (a/b/c would be 3 parts) in the
// mount indexed by parameter 1 is less than that of the mount indexed by
// parameter 2. Used in sorting.
func (m mounts) Less(i, j int) bool ***REMOVED***
	return m.parts(i) < m.parts(j)
***REMOVED***

// Swap swaps two items in an array of mounts. Used in sorting
func (m mounts) Swap(i, j int) ***REMOVED***
	m[i], m[j] = m[j], m[i]
***REMOVED***

// parts returns the number of parts in the destination of a mount. Used in sorting.
func (m mounts) parts(i int) int ***REMOVED***
	return strings.Count(filepath.Clean(m[i].Destination), string(os.PathSeparator))
***REMOVED***

// registerMountPoints initializes the container mount points with the configured volumes and bind mounts.
// It follows the next sequence to decide what to mount in each final destination:
//
// 1. Select the previously configured mount points for the containers, if any.
// 2. Select the volumes mounted from another containers. Overrides previously configured mount point destination.
// 3. Select the bind mounts set by the client. Overrides previously configured mount point destinations.
// 4. Cleanup old volumes that are about to be reassigned.
func (daemon *Daemon) registerMountPoints(container *container.Container, hostConfig *containertypes.HostConfig) (retErr error) ***REMOVED***
	binds := map[string]bool***REMOVED******REMOVED***
	mountPoints := map[string]*volume.MountPoint***REMOVED******REMOVED***
	parser := volume.NewParser(container.OS)
	defer func() ***REMOVED***
		// clean up the container mountpoints once return with error
		if retErr != nil ***REMOVED***
			for _, m := range mountPoints ***REMOVED***
				if m.Volume == nil ***REMOVED***
					continue
				***REMOVED***
				daemon.volumes.Dereference(m.Volume, container.ID)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	dereferenceIfExists := func(destination string) ***REMOVED***
		if v, ok := mountPoints[destination]; ok ***REMOVED***
			logrus.Debugf("Duplicate mount point '%s'", destination)
			if v.Volume != nil ***REMOVED***
				daemon.volumes.Dereference(v.Volume, container.ID)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// 1. Read already configured mount points.
	for destination, point := range container.MountPoints ***REMOVED***
		mountPoints[destination] = point
	***REMOVED***

	// 2. Read volumes from other containers.
	for _, v := range hostConfig.VolumesFrom ***REMOVED***
		containerID, mode, err := parser.ParseVolumesFrom(v)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		c, err := daemon.GetContainer(containerID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		for _, m := range c.MountPoints ***REMOVED***
			cp := &volume.MountPoint***REMOVED***
				Type:        m.Type,
				Name:        m.Name,
				Source:      m.Source,
				RW:          m.RW && parser.ReadWrite(mode),
				Driver:      m.Driver,
				Destination: m.Destination,
				Propagation: m.Propagation,
				Spec:        m.Spec,
				CopyData:    false,
			***REMOVED***

			if len(cp.Source) == 0 ***REMOVED***
				v, err := daemon.volumes.GetWithRef(cp.Name, cp.Driver, container.ID)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				cp.Volume = v
			***REMOVED***
			dereferenceIfExists(cp.Destination)
			mountPoints[cp.Destination] = cp
		***REMOVED***
	***REMOVED***

	// 3. Read bind mounts
	for _, b := range hostConfig.Binds ***REMOVED***
		bind, err := parser.ParseMountRaw(b, hostConfig.VolumeDriver)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// #10618
		_, tmpfsExists := hostConfig.Tmpfs[bind.Destination]
		if binds[bind.Destination] || tmpfsExists ***REMOVED***
			return duplicateMountPointError(bind.Destination)
		***REMOVED***

		if bind.Type == mounttypes.TypeVolume ***REMOVED***
			// create the volume
			v, err := daemon.volumes.CreateWithRef(bind.Name, bind.Driver, container.ID, nil, nil)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			bind.Volume = v
			bind.Source = v.Path()
			// bind.Name is an already existing volume, we need to use that here
			bind.Driver = v.DriverName()
			if bind.Driver == volume.DefaultDriverName ***REMOVED***
				setBindModeIfNull(bind)
			***REMOVED***
		***REMOVED***

		binds[bind.Destination] = true
		dereferenceIfExists(bind.Destination)
		mountPoints[bind.Destination] = bind
	***REMOVED***

	for _, cfg := range hostConfig.Mounts ***REMOVED***
		mp, err := parser.ParseMountSpec(cfg)
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***

		if binds[mp.Destination] ***REMOVED***
			return duplicateMountPointError(cfg.Target)
		***REMOVED***

		if mp.Type == mounttypes.TypeVolume ***REMOVED***
			var v volume.Volume
			if cfg.VolumeOptions != nil ***REMOVED***
				var driverOpts map[string]string
				if cfg.VolumeOptions.DriverConfig != nil ***REMOVED***
					driverOpts = cfg.VolumeOptions.DriverConfig.Options
				***REMOVED***
				v, err = daemon.volumes.CreateWithRef(mp.Name, mp.Driver, container.ID, driverOpts, cfg.VolumeOptions.Labels)
			***REMOVED*** else ***REMOVED***
				v, err = daemon.volumes.CreateWithRef(mp.Name, mp.Driver, container.ID, nil, nil)
			***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			mp.Volume = v
			mp.Name = v.Name()
			mp.Driver = v.DriverName()

			// only use the cached path here since getting the path is not necessary right now and calling `Path()` may be slow
			if cv, ok := v.(interface ***REMOVED***
				CachedPath() string
			***REMOVED***); ok ***REMOVED***
				mp.Source = cv.CachedPath()
			***REMOVED***
			if mp.Driver == volume.DefaultDriverName ***REMOVED***
				setBindModeIfNull(mp)
			***REMOVED***
		***REMOVED***

		binds[mp.Destination] = true
		dereferenceIfExists(mp.Destination)
		mountPoints[mp.Destination] = mp
	***REMOVED***

	container.Lock()

	// 4. Cleanup old volumes that are about to be reassigned.
	for _, m := range mountPoints ***REMOVED***
		if parser.IsBackwardCompatible(m) ***REMOVED***
			if mp, exists := container.MountPoints[m.Destination]; exists && mp.Volume != nil ***REMOVED***
				daemon.volumes.Dereference(mp.Volume, container.ID)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	container.MountPoints = mountPoints

	container.Unlock()

	return nil
***REMOVED***

// lazyInitializeVolume initializes a mountpoint's volume if needed.
// This happens after a daemon restart.
func (daemon *Daemon) lazyInitializeVolume(containerID string, m *volume.MountPoint) error ***REMOVED***
	if len(m.Driver) > 0 && m.Volume == nil ***REMOVED***
		v, err := daemon.volumes.GetWithRef(m.Name, m.Driver, containerID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		m.Volume = v
	***REMOVED***
	return nil
***REMOVED***

// backportMountSpec resolves mount specs (introduced in 1.13) from pre-1.13
// mount configurations
// The container lock should not be held when calling this function.
// Changes are only made in-memory and may make changes to containers referenced
// by `container.HostConfig.VolumesFrom`
func (daemon *Daemon) backportMountSpec(container *container.Container) ***REMOVED***
	container.Lock()
	defer container.Unlock()

	parser := volume.NewParser(container.OS)

	maybeUpdate := make(map[string]bool)
	for _, mp := range container.MountPoints ***REMOVED***
		if mp.Spec.Source != "" && mp.Type != "" ***REMOVED***
			continue
		***REMOVED***
		maybeUpdate[mp.Destination] = true
	***REMOVED***
	if len(maybeUpdate) == 0 ***REMOVED***
		return
	***REMOVED***

	mountSpecs := make(map[string]bool, len(container.HostConfig.Mounts))
	for _, m := range container.HostConfig.Mounts ***REMOVED***
		mountSpecs[m.Target] = true
	***REMOVED***

	binds := make(map[string]*volume.MountPoint, len(container.HostConfig.Binds))
	for _, rawSpec := range container.HostConfig.Binds ***REMOVED***
		mp, err := parser.ParseMountRaw(rawSpec, container.HostConfig.VolumeDriver)
		if err != nil ***REMOVED***
			logrus.WithError(err).Error("Got unexpected error while re-parsing raw volume spec during spec backport")
			continue
		***REMOVED***
		binds[mp.Destination] = mp
	***REMOVED***

	volumesFrom := make(map[string]volume.MountPoint)
	for _, fromSpec := range container.HostConfig.VolumesFrom ***REMOVED***
		from, _, err := parser.ParseVolumesFrom(fromSpec)
		if err != nil ***REMOVED***
			logrus.WithError(err).WithField("id", container.ID).Error("Error reading volumes-from spec during mount spec backport")
			continue
		***REMOVED***
		fromC, err := daemon.GetContainer(from)
		if err != nil ***REMOVED***
			logrus.WithError(err).WithField("from-container", from).Error("Error looking up volumes-from container")
			continue
		***REMOVED***

		// make sure from container's specs have been backported
		daemon.backportMountSpec(fromC)

		fromC.Lock()
		for t, mp := range fromC.MountPoints ***REMOVED***
			volumesFrom[t] = *mp
		***REMOVED***
		fromC.Unlock()
	***REMOVED***

	needsUpdate := func(containerMount, other *volume.MountPoint) bool ***REMOVED***
		if containerMount.Type != other.Type || !reflect.DeepEqual(containerMount.Spec, other.Spec) ***REMOVED***
			return true
		***REMOVED***
		return false
	***REMOVED***

	// main
	for _, cm := range container.MountPoints ***REMOVED***
		if !maybeUpdate[cm.Destination] ***REMOVED***
			continue
		***REMOVED***
		// nothing to backport if from hostconfig.Mounts
		if mountSpecs[cm.Destination] ***REMOVED***
			continue
		***REMOVED***

		if mp, exists := binds[cm.Destination]; exists ***REMOVED***
			if needsUpdate(cm, mp) ***REMOVED***
				cm.Spec = mp.Spec
				cm.Type = mp.Type
			***REMOVED***
			continue
		***REMOVED***

		if cm.Name != "" ***REMOVED***
			if mp, exists := volumesFrom[cm.Destination]; exists ***REMOVED***
				if needsUpdate(cm, &mp) ***REMOVED***
					cm.Spec = mp.Spec
					cm.Type = mp.Type
				***REMOVED***
				continue
			***REMOVED***

			if cm.Type != "" ***REMOVED***
				// probably specified via the hostconfig.Mounts
				continue
			***REMOVED***

			// anon volume
			cm.Type = mounttypes.TypeVolume
			cm.Spec.Type = mounttypes.TypeVolume
		***REMOVED*** else ***REMOVED***
			if cm.Type != "" ***REMOVED***
				// already updated
				continue
			***REMOVED***

			cm.Type = mounttypes.TypeBind
			cm.Spec.Type = mounttypes.TypeBind
			cm.Spec.Source = cm.Source
			if cm.Propagation != "" ***REMOVED***
				cm.Spec.BindOptions = &mounttypes.BindOptions***REMOVED***
					Propagation: cm.Propagation,
				***REMOVED***
			***REMOVED***
		***REMOVED***

		cm.Spec.Target = cm.Destination
		cm.Spec.ReadOnly = !cm.RW
	***REMOVED***
***REMOVED***

func (daemon *Daemon) traverseLocalVolumes(fn func(volume.Volume) error) error ***REMOVED***
	localVolumeDriver, err := volumedrivers.GetDriver(volume.DefaultDriverName)
	if err != nil ***REMOVED***
		return fmt.Errorf("can't retrieve local volume driver: %v", err)
	***REMOVED***
	vols, err := localVolumeDriver.List()
	if err != nil ***REMOVED***
		return fmt.Errorf("can't retrieve local volumes: %v", err)
	***REMOVED***

	for _, v := range vols ***REMOVED***
		name := v.Name()
		vol, err := daemon.volumes.Get(name)
		if err != nil ***REMOVED***
			logrus.Warnf("failed to retrieve volume %s from store: %v", name, err)
		***REMOVED*** else ***REMOVED***
			// daemon.volumes.Get will return DetailedVolume
			v = vol
		***REMOVED***

		err = fn(v)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
