// +build linux freebsd

package container

import (
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	mounttypes "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/chrootarchive"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/volume"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

const (
	// DefaultStopTimeout is the timeout (in seconds) for the syscall signal used to stop a container.
	DefaultStopTimeout = 10

	containerSecretMountPath = "/run/secrets"
)

// TrySetNetworkMount attempts to set the network mounts given a provided destination and
// the path to use for it; return true if the given destination was a network mount file
func (container *Container) TrySetNetworkMount(destination string, path string) bool ***REMOVED***
	if destination == "/etc/resolv.conf" ***REMOVED***
		container.ResolvConfPath = path
		return true
	***REMOVED***
	if destination == "/etc/hostname" ***REMOVED***
		container.HostnamePath = path
		return true
	***REMOVED***
	if destination == "/etc/hosts" ***REMOVED***
		container.HostsPath = path
		return true
	***REMOVED***

	return false
***REMOVED***

// BuildHostnameFile writes the container's hostname file.
func (container *Container) BuildHostnameFile() error ***REMOVED***
	hostnamePath, err := container.GetRootResourcePath("hostname")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	container.HostnamePath = hostnamePath
	return ioutil.WriteFile(container.HostnamePath, []byte(container.Config.Hostname+"\n"), 0644)
***REMOVED***

// NetworkMounts returns the list of network mounts.
func (container *Container) NetworkMounts() []Mount ***REMOVED***
	var mounts []Mount
	shared := container.HostConfig.NetworkMode.IsContainer()
	parser := volume.NewParser(container.OS)
	if container.ResolvConfPath != "" ***REMOVED***
		if _, err := os.Stat(container.ResolvConfPath); err != nil ***REMOVED***
			logrus.Warnf("ResolvConfPath set to %q, but can't stat this filename (err = %v); skipping", container.ResolvConfPath, err)
		***REMOVED*** else ***REMOVED***
			writable := !container.HostConfig.ReadonlyRootfs
			if m, exists := container.MountPoints["/etc/resolv.conf"]; exists ***REMOVED***
				writable = m.RW
			***REMOVED*** else ***REMOVED***
				label.Relabel(container.ResolvConfPath, container.MountLabel, shared)
			***REMOVED***
			mounts = append(mounts, Mount***REMOVED***
				Source:      container.ResolvConfPath,
				Destination: "/etc/resolv.conf",
				Writable:    writable,
				Propagation: string(parser.DefaultPropagationMode()),
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	if container.HostnamePath != "" ***REMOVED***
		if _, err := os.Stat(container.HostnamePath); err != nil ***REMOVED***
			logrus.Warnf("HostnamePath set to %q, but can't stat this filename (err = %v); skipping", container.HostnamePath, err)
		***REMOVED*** else ***REMOVED***
			writable := !container.HostConfig.ReadonlyRootfs
			if m, exists := container.MountPoints["/etc/hostname"]; exists ***REMOVED***
				writable = m.RW
			***REMOVED*** else ***REMOVED***
				label.Relabel(container.HostnamePath, container.MountLabel, shared)
			***REMOVED***
			mounts = append(mounts, Mount***REMOVED***
				Source:      container.HostnamePath,
				Destination: "/etc/hostname",
				Writable:    writable,
				Propagation: string(parser.DefaultPropagationMode()),
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	if container.HostsPath != "" ***REMOVED***
		if _, err := os.Stat(container.HostsPath); err != nil ***REMOVED***
			logrus.Warnf("HostsPath set to %q, but can't stat this filename (err = %v); skipping", container.HostsPath, err)
		***REMOVED*** else ***REMOVED***
			writable := !container.HostConfig.ReadonlyRootfs
			if m, exists := container.MountPoints["/etc/hosts"]; exists ***REMOVED***
				writable = m.RW
			***REMOVED*** else ***REMOVED***
				label.Relabel(container.HostsPath, container.MountLabel, shared)
			***REMOVED***
			mounts = append(mounts, Mount***REMOVED***
				Source:      container.HostsPath,
				Destination: "/etc/hosts",
				Writable:    writable,
				Propagation: string(parser.DefaultPropagationMode()),
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return mounts
***REMOVED***

// CopyImagePathContent copies files in destination to the volume.
func (container *Container) CopyImagePathContent(v volume.Volume, destination string) error ***REMOVED***
	rootfs, err := container.GetResourcePath(destination)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if _, err = ioutil.ReadDir(rootfs); err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	id := stringid.GenerateNonCryptoID()
	path, err := v.Mount(id)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	defer func() ***REMOVED***
		if err := v.Unmount(id); err != nil ***REMOVED***
			logrus.Warnf("error while unmounting volume %s: %v", v.Name(), err)
		***REMOVED***
	***REMOVED***()
	if err := label.Relabel(path, container.MountLabel, true); err != nil && err != unix.ENOTSUP ***REMOVED***
		return err
	***REMOVED***
	return copyExistingContents(rootfs, path)
***REMOVED***

// ShmResourcePath returns path to shm
func (container *Container) ShmResourcePath() (string, error) ***REMOVED***
	return container.MountsResourcePath("shm")
***REMOVED***

// HasMountFor checks if path is a mountpoint
func (container *Container) HasMountFor(path string) bool ***REMOVED***
	_, exists := container.MountPoints[path]
	if exists ***REMOVED***
		return true
	***REMOVED***

	// Also search among the tmpfs mounts
	for dest := range container.HostConfig.Tmpfs ***REMOVED***
		if dest == path ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// UnmountIpcMount uses the provided unmount function to unmount shm if it was mounted
func (container *Container) UnmountIpcMount(unmount func(pth string) error) error ***REMOVED***
	if container.HasMountFor("/dev/shm") ***REMOVED***
		return nil
	***REMOVED***

	// container.ShmPath should not be used here as it may point
	// to the host's or other container's /dev/shm
	shmPath, err := container.ShmResourcePath()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if shmPath == "" ***REMOVED***
		return nil
	***REMOVED***
	if err = unmount(shmPath); err != nil && !os.IsNotExist(err) ***REMOVED***
		if mounted, mErr := mount.Mounted(shmPath); mounted || mErr != nil ***REMOVED***
			return errors.Wrapf(err, "umount %s", shmPath)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// IpcMounts returns the list of IPC mounts
func (container *Container) IpcMounts() []Mount ***REMOVED***
	var mounts []Mount
	parser := volume.NewParser(container.OS)

	if container.HasMountFor("/dev/shm") ***REMOVED***
		return mounts
	***REMOVED***
	if container.ShmPath == "" ***REMOVED***
		return mounts
	***REMOVED***

	label.SetFileLabel(container.ShmPath, container.MountLabel)
	mounts = append(mounts, Mount***REMOVED***
		Source:      container.ShmPath,
		Destination: "/dev/shm",
		Writable:    true,
		Propagation: string(parser.DefaultPropagationMode()),
	***REMOVED***)

	return mounts
***REMOVED***

// SecretMounts returns the mounts for the secret path.
func (container *Container) SecretMounts() ([]Mount, error) ***REMOVED***
	var mounts []Mount
	for _, r := range container.SecretReferences ***REMOVED***
		if r.File == nil ***REMOVED***
			continue
		***REMOVED***
		src, err := container.SecretFilePath(*r)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		mounts = append(mounts, Mount***REMOVED***
			Source:      src,
			Destination: getSecretTargetPath(r),
			Writable:    false,
		***REMOVED***)
	***REMOVED***

	return mounts, nil
***REMOVED***

// UnmountSecrets unmounts the local tmpfs for secrets
func (container *Container) UnmountSecrets() error ***REMOVED***
	p, err := container.SecretMountPath()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := os.Stat(p); err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	return mount.RecursiveUnmount(p)
***REMOVED***

// ConfigMounts returns the mounts for configs.
func (container *Container) ConfigMounts() ([]Mount, error) ***REMOVED***
	var mounts []Mount
	for _, configRef := range container.ConfigReferences ***REMOVED***
		if configRef.File == nil ***REMOVED***
			continue
		***REMOVED***
		src, err := container.ConfigFilePath(*configRef)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		mounts = append(mounts, Mount***REMOVED***
			Source:      src,
			Destination: configRef.File.Name,
			Writable:    false,
		***REMOVED***)
	***REMOVED***

	return mounts, nil
***REMOVED***

type conflictingUpdateOptions string

func (e conflictingUpdateOptions) Error() string ***REMOVED***
	return string(e)
***REMOVED***

func (e conflictingUpdateOptions) Conflict() ***REMOVED******REMOVED***

// UpdateContainer updates configuration of a container. Callers must hold a Lock on the Container.
func (container *Container) UpdateContainer(hostConfig *containertypes.HostConfig) error ***REMOVED***
	// update resources of container
	resources := hostConfig.Resources
	cResources := &container.HostConfig.Resources

	// validate NanoCPUs, CPUPeriod, and CPUQuota
	// Because NanoCPU effectively updates CPUPeriod/CPUQuota,
	// once NanoCPU is already set, updating CPUPeriod/CPUQuota will be blocked, and vice versa.
	// In the following we make sure the intended update (resources) does not conflict with the existing (cResource).
	if resources.NanoCPUs > 0 && cResources.CPUPeriod > 0 ***REMOVED***
		return conflictingUpdateOptions("Conflicting options: Nano CPUs cannot be updated as CPU Period has already been set")
	***REMOVED***
	if resources.NanoCPUs > 0 && cResources.CPUQuota > 0 ***REMOVED***
		return conflictingUpdateOptions("Conflicting options: Nano CPUs cannot be updated as CPU Quota has already been set")
	***REMOVED***
	if resources.CPUPeriod > 0 && cResources.NanoCPUs > 0 ***REMOVED***
		return conflictingUpdateOptions("Conflicting options: CPU Period cannot be updated as NanoCPUs has already been set")
	***REMOVED***
	if resources.CPUQuota > 0 && cResources.NanoCPUs > 0 ***REMOVED***
		return conflictingUpdateOptions("Conflicting options: CPU Quota cannot be updated as NanoCPUs has already been set")
	***REMOVED***

	if resources.BlkioWeight != 0 ***REMOVED***
		cResources.BlkioWeight = resources.BlkioWeight
	***REMOVED***
	if resources.CPUShares != 0 ***REMOVED***
		cResources.CPUShares = resources.CPUShares
	***REMOVED***
	if resources.NanoCPUs != 0 ***REMOVED***
		cResources.NanoCPUs = resources.NanoCPUs
	***REMOVED***
	if resources.CPUPeriod != 0 ***REMOVED***
		cResources.CPUPeriod = resources.CPUPeriod
	***REMOVED***
	if resources.CPUQuota != 0 ***REMOVED***
		cResources.CPUQuota = resources.CPUQuota
	***REMOVED***
	if resources.CpusetCpus != "" ***REMOVED***
		cResources.CpusetCpus = resources.CpusetCpus
	***REMOVED***
	if resources.CpusetMems != "" ***REMOVED***
		cResources.CpusetMems = resources.CpusetMems
	***REMOVED***
	if resources.Memory != 0 ***REMOVED***
		// if memory limit smaller than already set memoryswap limit and doesn't
		// update the memoryswap limit, then error out.
		if resources.Memory > cResources.MemorySwap && resources.MemorySwap == 0 ***REMOVED***
			return conflictingUpdateOptions("Memory limit should be smaller than already set memoryswap limit, update the memoryswap at the same time")
		***REMOVED***
		cResources.Memory = resources.Memory
	***REMOVED***
	if resources.MemorySwap != 0 ***REMOVED***
		cResources.MemorySwap = resources.MemorySwap
	***REMOVED***
	if resources.MemoryReservation != 0 ***REMOVED***
		cResources.MemoryReservation = resources.MemoryReservation
	***REMOVED***
	if resources.KernelMemory != 0 ***REMOVED***
		cResources.KernelMemory = resources.KernelMemory
	***REMOVED***
	if resources.CPURealtimePeriod != 0 ***REMOVED***
		cResources.CPURealtimePeriod = resources.CPURealtimePeriod
	***REMOVED***
	if resources.CPURealtimeRuntime != 0 ***REMOVED***
		cResources.CPURealtimeRuntime = resources.CPURealtimeRuntime
	***REMOVED***

	// update HostConfig of container
	if hostConfig.RestartPolicy.Name != "" ***REMOVED***
		if container.HostConfig.AutoRemove && !hostConfig.RestartPolicy.IsNone() ***REMOVED***
			return conflictingUpdateOptions("Restart policy cannot be updated because AutoRemove is enabled for the container")
		***REMOVED***
		container.HostConfig.RestartPolicy = hostConfig.RestartPolicy
	***REMOVED***

	return nil
***REMOVED***

// DetachAndUnmount uses a detached mount on all mount destinations, then
// unmounts each volume normally.
// This is used from daemon/archive for `docker cp`
func (container *Container) DetachAndUnmount(volumeEventLog func(name, action string, attributes map[string]string)) error ***REMOVED***
	networkMounts := container.NetworkMounts()
	mountPaths := make([]string, 0, len(container.MountPoints)+len(networkMounts))

	for _, mntPoint := range container.MountPoints ***REMOVED***
		dest, err := container.GetResourcePath(mntPoint.Destination)
		if err != nil ***REMOVED***
			logrus.Warnf("Failed to get volume destination path for container '%s' at '%s' while lazily unmounting: %v", container.ID, mntPoint.Destination, err)
			continue
		***REMOVED***
		mountPaths = append(mountPaths, dest)
	***REMOVED***

	for _, m := range networkMounts ***REMOVED***
		dest, err := container.GetResourcePath(m.Destination)
		if err != nil ***REMOVED***
			logrus.Warnf("Failed to get volume destination path for container '%s' at '%s' while lazily unmounting: %v", container.ID, m.Destination, err)
			continue
		***REMOVED***
		mountPaths = append(mountPaths, dest)
	***REMOVED***

	for _, mountPath := range mountPaths ***REMOVED***
		if err := detachMounted(mountPath); err != nil ***REMOVED***
			logrus.Warnf("%s unmountVolumes: Failed to do lazy umount fo volume '%s': %v", container.ID, mountPath, err)
		***REMOVED***
	***REMOVED***
	return container.UnmountVolumes(volumeEventLog)
***REMOVED***

// copyExistingContents copies from the source to the destination and
// ensures the ownership is appropriately set.
func copyExistingContents(source, destination string) error ***REMOVED***
	volList, err := ioutil.ReadDir(source)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(volList) > 0 ***REMOVED***
		srcList, err := ioutil.ReadDir(destination)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if len(srcList) == 0 ***REMOVED***
			// If the source volume is empty, copies files from the root into the volume
			if err := chrootarchive.NewArchiver(nil).CopyWithTar(source, destination); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return copyOwnership(source, destination)
***REMOVED***

// copyOwnership copies the permissions and uid:gid of the source file
// to the destination file
func copyOwnership(source, destination string) error ***REMOVED***
	stat, err := system.Stat(source)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	destStat, err := system.Stat(destination)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// In some cases, even though UID/GID match and it would effectively be a no-op,
	// this can return a permission denied error... for example if this is an NFS
	// mount.
	// Since it's not really an error that we can't chown to the same UID/GID, don't
	// even bother trying in such cases.
	if stat.UID() != destStat.UID() || stat.GID() != destStat.GID() ***REMOVED***
		if err := os.Chown(destination, int(stat.UID()), int(stat.GID())); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if stat.Mode() != destStat.Mode() ***REMOVED***
		return os.Chmod(destination, os.FileMode(stat.Mode()))
	***REMOVED***
	return nil
***REMOVED***

// TmpfsMounts returns the list of tmpfs mounts
func (container *Container) TmpfsMounts() ([]Mount, error) ***REMOVED***
	parser := volume.NewParser(container.OS)
	var mounts []Mount
	for dest, data := range container.HostConfig.Tmpfs ***REMOVED***
		mounts = append(mounts, Mount***REMOVED***
			Source:      "tmpfs",
			Destination: dest,
			Data:        data,
		***REMOVED***)
	***REMOVED***
	for dest, mnt := range container.MountPoints ***REMOVED***
		if mnt.Type == mounttypes.TypeTmpfs ***REMOVED***
			data, err := parser.ConvertTmpfsOptions(mnt.Spec.TmpfsOptions, mnt.Spec.ReadOnly)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			mounts = append(mounts, Mount***REMOVED***
				Source:      "tmpfs",
				Destination: dest,
				Data:        data,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return mounts, nil
***REMOVED***

// EnableServiceDiscoveryOnDefaultNetwork Enable service discovery on default network
func (container *Container) EnableServiceDiscoveryOnDefaultNetwork() bool ***REMOVED***
	return false
***REMOVED***

// GetMountPoints gives a platform specific transformation to types.MountPoint. Callers must hold a Container lock.
func (container *Container) GetMountPoints() []types.MountPoint ***REMOVED***
	mountPoints := make([]types.MountPoint, 0, len(container.MountPoints))
	for _, m := range container.MountPoints ***REMOVED***
		mountPoints = append(mountPoints, types.MountPoint***REMOVED***
			Type:        m.Type,
			Name:        m.Name,
			Source:      m.Path(),
			Destination: m.Destination,
			Driver:      m.Driver,
			Mode:        m.Mode,
			RW:          m.RW,
			Propagation: m.Propagation,
		***REMOVED***)
	***REMOVED***
	return mountPoints
***REMOVED***
