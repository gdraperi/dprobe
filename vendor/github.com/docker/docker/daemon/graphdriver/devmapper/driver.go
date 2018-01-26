// +build linux

package devmapper

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/devicemapper"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/locker"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/system"
	units "github.com/docker/go-units"
)

func init() ***REMOVED***
	graphdriver.Register("devicemapper", Init)
***REMOVED***

// Driver contains the device set mounted and the home directory
type Driver struct ***REMOVED***
	*DeviceSet
	home    string
	uidMaps []idtools.IDMap
	gidMaps []idtools.IDMap
	ctr     *graphdriver.RefCounter
	locker  *locker.Locker
***REMOVED***

// Init creates a driver with the given home and the set of options.
func Init(home string, options []string, uidMaps, gidMaps []idtools.IDMap) (graphdriver.Driver, error) ***REMOVED***
	deviceSet, err := NewDeviceSet(home, true, options, uidMaps, gidMaps)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	d := &Driver***REMOVED***
		DeviceSet: deviceSet,
		home:      home,
		uidMaps:   uidMaps,
		gidMaps:   gidMaps,
		ctr:       graphdriver.NewRefCounter(graphdriver.NewDefaultChecker()),
		locker:    locker.New(),
	***REMOVED***

	return graphdriver.NewNaiveDiffDriver(d, uidMaps, gidMaps), nil
***REMOVED***

func (d *Driver) String() string ***REMOVED***
	return "devicemapper"
***REMOVED***

// Status returns the status about the driver in a printable format.
// Information returned contains Pool Name, Data File, Metadata file, disk usage by
// the data and metadata, etc.
func (d *Driver) Status() [][2]string ***REMOVED***
	s := d.DeviceSet.Status()

	status := [][2]string***REMOVED***
		***REMOVED***"Pool Name", s.PoolName***REMOVED***,
		***REMOVED***"Pool Blocksize", units.HumanSize(float64(s.SectorSize))***REMOVED***,
		***REMOVED***"Base Device Size", units.HumanSize(float64(s.BaseDeviceSize))***REMOVED***,
		***REMOVED***"Backing Filesystem", s.BaseDeviceFS***REMOVED***,
		***REMOVED***"Udev Sync Supported", fmt.Sprintf("%v", s.UdevSyncSupported)***REMOVED***,
	***REMOVED***

	if len(s.DataFile) > 0 ***REMOVED***
		status = append(status, [2]string***REMOVED***"Data file", s.DataFile***REMOVED***)
	***REMOVED***
	if len(s.MetadataFile) > 0 ***REMOVED***
		status = append(status, [2]string***REMOVED***"Metadata file", s.MetadataFile***REMOVED***)
	***REMOVED***
	if len(s.DataLoopback) > 0 ***REMOVED***
		status = append(status, [2]string***REMOVED***"Data loop file", s.DataLoopback***REMOVED***)
	***REMOVED***
	if len(s.MetadataLoopback) > 0 ***REMOVED***
		status = append(status, [2]string***REMOVED***"Metadata loop file", s.MetadataLoopback***REMOVED***)
	***REMOVED***

	status = append(status, [][2]string***REMOVED***
		***REMOVED***"Data Space Used", units.HumanSize(float64(s.Data.Used))***REMOVED***,
		***REMOVED***"Data Space Total", units.HumanSize(float64(s.Data.Total))***REMOVED***,
		***REMOVED***"Data Space Available", units.HumanSize(float64(s.Data.Available))***REMOVED***,
		***REMOVED***"Metadata Space Used", units.HumanSize(float64(s.Metadata.Used))***REMOVED***,
		***REMOVED***"Metadata Space Total", units.HumanSize(float64(s.Metadata.Total))***REMOVED***,
		***REMOVED***"Metadata Space Available", units.HumanSize(float64(s.Metadata.Available))***REMOVED***,
		***REMOVED***"Thin Pool Minimum Free Space", units.HumanSize(float64(s.MinFreeSpace))***REMOVED***,
		***REMOVED***"Deferred Removal Enabled", fmt.Sprintf("%v", s.DeferredRemoveEnabled)***REMOVED***,
		***REMOVED***"Deferred Deletion Enabled", fmt.Sprintf("%v", s.DeferredDeleteEnabled)***REMOVED***,
		***REMOVED***"Deferred Deleted Device Count", fmt.Sprintf("%v", s.DeferredDeletedDeviceCount)***REMOVED***,
	***REMOVED***...)

	if vStr, err := devicemapper.GetLibraryVersion(); err == nil ***REMOVED***
		status = append(status, [2]string***REMOVED***"Library Version", vStr***REMOVED***)
	***REMOVED***
	return status
***REMOVED***

// GetMetadata returns a map of information about the device.
func (d *Driver) GetMetadata(id string) (map[string]string, error) ***REMOVED***
	m, err := d.DeviceSet.exportDeviceMetadata(id)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	metadata := make(map[string]string)
	metadata["DeviceId"] = strconv.Itoa(m.deviceID)
	metadata["DeviceSize"] = strconv.FormatUint(m.deviceSize, 10)
	metadata["DeviceName"] = m.deviceName
	return metadata, nil
***REMOVED***

// Cleanup unmounts a device.
func (d *Driver) Cleanup() error ***REMOVED***
	err := d.DeviceSet.Shutdown(d.home)

	if err2 := mount.RecursiveUnmount(d.home); err == nil ***REMOVED***
		err = err2
	***REMOVED***

	return err
***REMOVED***

// CreateReadWrite creates a layer that is writable for use as a container
// file system.
func (d *Driver) CreateReadWrite(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	return d.Create(id, parent, opts)
***REMOVED***

// Create adds a device with a given id and the parent.
func (d *Driver) Create(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	var storageOpt map[string]string
	if opts != nil ***REMOVED***
		storageOpt = opts.StorageOpt
	***REMOVED***
	return d.DeviceSet.AddDevice(id, parent, storageOpt)
***REMOVED***

// Remove removes a device with a given id, unmounts the filesystem.
func (d *Driver) Remove(id string) error ***REMOVED***
	d.locker.Lock(id)
	defer d.locker.Unlock(id)
	if !d.DeviceSet.HasDevice(id) ***REMOVED***
		// Consider removing a non-existing device a no-op
		// This is useful to be able to progress on container removal
		// if the underlying device has gone away due to earlier errors
		return nil
	***REMOVED***

	// This assumes the device has been properly Get/Put:ed and thus is unmounted
	if err := d.DeviceSet.DeleteDevice(id, false); err != nil ***REMOVED***
		return fmt.Errorf("failed to remove device %s: %v", id, err)
	***REMOVED***
	return system.EnsureRemoveAll(path.Join(d.home, "mnt", id))
***REMOVED***

// Get mounts a device with given id into the root filesystem
func (d *Driver) Get(id, mountLabel string) (containerfs.ContainerFS, error) ***REMOVED***
	d.locker.Lock(id)
	defer d.locker.Unlock(id)
	mp := path.Join(d.home, "mnt", id)
	rootFs := path.Join(mp, "rootfs")
	if count := d.ctr.Increment(mp); count > 1 ***REMOVED***
		return containerfs.NewLocalContainerFS(rootFs), nil
	***REMOVED***

	uid, gid, err := idtools.GetRootUIDGID(d.uidMaps, d.gidMaps)
	if err != nil ***REMOVED***
		d.ctr.Decrement(mp)
		return nil, err
	***REMOVED***

	// Create the target directories if they don't exist
	if err := idtools.MkdirAllAndChown(path.Join(d.home, "mnt"), 0755, idtools.IDPair***REMOVED***UID: uid, GID: gid***REMOVED***); err != nil ***REMOVED***
		d.ctr.Decrement(mp)
		return nil, err
	***REMOVED***
	if err := idtools.MkdirAndChown(mp, 0755, idtools.IDPair***REMOVED***UID: uid, GID: gid***REMOVED***); err != nil && !os.IsExist(err) ***REMOVED***
		d.ctr.Decrement(mp)
		return nil, err
	***REMOVED***

	// Mount the device
	if err := d.DeviceSet.MountDevice(id, mp, mountLabel); err != nil ***REMOVED***
		d.ctr.Decrement(mp)
		return nil, err
	***REMOVED***

	if err := idtools.MkdirAllAndChown(rootFs, 0755, idtools.IDPair***REMOVED***UID: uid, GID: gid***REMOVED***); err != nil ***REMOVED***
		d.ctr.Decrement(mp)
		d.DeviceSet.UnmountDevice(id, mp)
		return nil, err
	***REMOVED***

	idFile := path.Join(mp, "id")
	if _, err := os.Stat(idFile); err != nil && os.IsNotExist(err) ***REMOVED***
		// Create an "id" file with the container/image id in it to help reconstruct this in case
		// of later problems
		if err := ioutil.WriteFile(idFile, []byte(id), 0600); err != nil ***REMOVED***
			d.ctr.Decrement(mp)
			d.DeviceSet.UnmountDevice(id, mp)
			return nil, err
		***REMOVED***
	***REMOVED***

	return containerfs.NewLocalContainerFS(rootFs), nil
***REMOVED***

// Put unmounts a device and removes it.
func (d *Driver) Put(id string) error ***REMOVED***
	d.locker.Lock(id)
	defer d.locker.Unlock(id)
	mp := path.Join(d.home, "mnt", id)
	if count := d.ctr.Decrement(mp); count > 0 ***REMOVED***
		return nil
	***REMOVED***

	err := d.DeviceSet.UnmountDevice(id, mp)
	if err != nil ***REMOVED***
		logrus.Errorf("devmapper: Error unmounting device %s: %v", id, err)
	***REMOVED***

	return err
***REMOVED***

// Exists checks to see if the device exists.
func (d *Driver) Exists(id string) bool ***REMOVED***
	return d.DeviceSet.HasDevice(id)
***REMOVED***
