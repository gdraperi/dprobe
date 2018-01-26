// +build linux

package devmapper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/pkg/devicemapper"
	"github.com/docker/docker/pkg/dmesg"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/loopback"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/parsers"
	"github.com/docker/docker/pkg/parsers/kernel"
	units "github.com/docker/go-units"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

var (
	defaultDataLoopbackSize      int64  = 100 * 1024 * 1024 * 1024
	defaultMetaDataLoopbackSize  int64  = 2 * 1024 * 1024 * 1024
	defaultBaseFsSize            uint64 = 10 * 1024 * 1024 * 1024
	defaultThinpBlockSize        uint32 = 128 // 64K = 128 512b sectors
	defaultUdevSyncOverride             = false
	maxDeviceID                         = 0xffffff // 24 bit, pool limit
	deviceIDMapSz                       = (maxDeviceID + 1) / 8
	driverDeferredRemovalSupport        = false
	enableDeferredRemoval               = false
	enableDeferredDeletion              = false
	userBaseSize                        = false
	defaultMinFreeSpacePercent   uint32 = 10
	lvmSetupConfigForce          bool
)

const deviceSetMetaFile string = "deviceset-metadata"
const transactionMetaFile string = "transaction-metadata"

type transaction struct ***REMOVED***
	OpenTransactionID uint64 `json:"open_transaction_id"`
	DeviceIDHash      string `json:"device_hash"`
	DeviceID          int    `json:"device_id"`
***REMOVED***

type devInfo struct ***REMOVED***
	Hash          string `json:"-"`
	DeviceID      int    `json:"device_id"`
	Size          uint64 `json:"size"`
	TransactionID uint64 `json:"transaction_id"`
	Initialized   bool   `json:"initialized"`
	Deleted       bool   `json:"deleted"`
	devices       *DeviceSet

	// The global DeviceSet lock guarantees that we serialize all
	// the calls to libdevmapper (which is not threadsafe), but we
	// sometimes release that lock while sleeping. In that case
	// this per-device lock is still held, protecting against
	// other accesses to the device that we're doing the wait on.
	//
	// WARNING: In order to avoid AB-BA deadlocks when releasing
	// the global lock while holding the per-device locks all
	// device locks must be acquired *before* the device lock, and
	// multiple device locks should be acquired parent before child.
	lock sync.Mutex
***REMOVED***

type metaData struct ***REMOVED***
	Devices map[string]*devInfo `json:"Devices"`
***REMOVED***

// DeviceSet holds information about list of devices
type DeviceSet struct ***REMOVED***
	metaData      `json:"-"`
	sync.Mutex    `json:"-"` // Protects all fields of DeviceSet and serializes calls into libdevmapper
	root          string
	devicePrefix  string
	TransactionID uint64 `json:"-"`
	NextDeviceID  int    `json:"next_device_id"`
	deviceIDMap   []byte

	// Options
	dataLoopbackSize      int64
	metaDataLoopbackSize  int64
	baseFsSize            uint64
	filesystem            string
	mountOptions          string
	mkfsArgs              []string
	dataDevice            string // block or loop dev
	dataLoopFile          string // loopback file, if used
	metadataDevice        string // block or loop dev
	metadataLoopFile      string // loopback file, if used
	doBlkDiscard          bool
	thinpBlockSize        uint32
	thinPoolDevice        string
	transaction           `json:"-"`
	overrideUdevSyncCheck bool
	deferredRemove        bool   // use deferred removal
	deferredDelete        bool   // use deferred deletion
	BaseDeviceUUID        string // save UUID of base device
	BaseDeviceFilesystem  string // save filesystem of base device
	nrDeletedDevices      uint   // number of deleted devices
	deletionWorkerTicker  *time.Ticker
	uidMaps               []idtools.IDMap
	gidMaps               []idtools.IDMap
	minFreeSpacePercent   uint32 //min free space percentage in thinpool
	xfsNospaceRetries     string // max retries when xfs receives ENOSPC
	lvmSetupConfig        directLVMConfig
***REMOVED***

// DiskUsage contains information about disk usage and is used when reporting Status of a device.
type DiskUsage struct ***REMOVED***
	// Used bytes on the disk.
	Used uint64
	// Total bytes on the disk.
	Total uint64
	// Available bytes on the disk.
	Available uint64
***REMOVED***

// Status returns the information about the device.
type Status struct ***REMOVED***
	// PoolName is the name of the data pool.
	PoolName string
	// DataFile is the actual block device for data.
	DataFile string
	// DataLoopback loopback file, if used.
	DataLoopback string
	// MetadataFile is the actual block device for metadata.
	MetadataFile string
	// MetadataLoopback is the loopback file, if used.
	MetadataLoopback string
	// Data is the disk used for data.
	Data DiskUsage
	// Metadata is the disk used for meta data.
	Metadata DiskUsage
	// BaseDeviceSize is base size of container and image
	BaseDeviceSize uint64
	// BaseDeviceFS is backing filesystem.
	BaseDeviceFS string
	// SectorSize size of the vector.
	SectorSize uint64
	// UdevSyncSupported is true if sync is supported.
	UdevSyncSupported bool
	// DeferredRemoveEnabled is true then the device is not unmounted.
	DeferredRemoveEnabled bool
	// True if deferred deletion is enabled. This is different from
	// deferred removal. "removal" means that device mapper device is
	// deactivated. Thin device is still in thin pool and can be activated
	// again. But "deletion" means that thin device will be deleted from
	// thin pool and it can't be activated again.
	DeferredDeleteEnabled      bool
	DeferredDeletedDeviceCount uint
	MinFreeSpace               uint64
***REMOVED***

// Structure used to export image/container metadata in docker inspect.
type deviceMetadata struct ***REMOVED***
	deviceID   int
	deviceSize uint64 // size in bytes
	deviceName string // Device name as used during activation
***REMOVED***

// DevStatus returns information about device mounted containing its id, size and sector information.
type DevStatus struct ***REMOVED***
	// DeviceID is the id of the device.
	DeviceID int
	// Size is the size of the filesystem.
	Size uint64
	// TransactionID is a unique integer per device set used to identify an operation on the file system, this number is incremental.
	TransactionID uint64
	// SizeInSectors indicates the size of the sectors allocated.
	SizeInSectors uint64
	// MappedSectors indicates number of mapped sectors.
	MappedSectors uint64
	// HighestMappedSector is the pointer to the highest mapped sector.
	HighestMappedSector uint64
***REMOVED***

func getDevName(name string) string ***REMOVED***
	return "/dev/mapper/" + name
***REMOVED***

func (info *devInfo) Name() string ***REMOVED***
	hash := info.Hash
	if hash == "" ***REMOVED***
		hash = "base"
	***REMOVED***
	return fmt.Sprintf("%s-%s", info.devices.devicePrefix, hash)
***REMOVED***

func (info *devInfo) DevName() string ***REMOVED***
	return getDevName(info.Name())
***REMOVED***

func (devices *DeviceSet) loopbackDir() string ***REMOVED***
	return path.Join(devices.root, "devicemapper")
***REMOVED***

func (devices *DeviceSet) metadataDir() string ***REMOVED***
	return path.Join(devices.root, "metadata")
***REMOVED***

func (devices *DeviceSet) metadataFile(info *devInfo) string ***REMOVED***
	file := info.Hash
	if file == "" ***REMOVED***
		file = "base"
	***REMOVED***
	return path.Join(devices.metadataDir(), file)
***REMOVED***

func (devices *DeviceSet) transactionMetaFile() string ***REMOVED***
	return path.Join(devices.metadataDir(), transactionMetaFile)
***REMOVED***

func (devices *DeviceSet) deviceSetMetaFile() string ***REMOVED***
	return path.Join(devices.metadataDir(), deviceSetMetaFile)
***REMOVED***

func (devices *DeviceSet) oldMetadataFile() string ***REMOVED***
	return path.Join(devices.loopbackDir(), "json")
***REMOVED***

func (devices *DeviceSet) getPoolName() string ***REMOVED***
	if devices.thinPoolDevice == "" ***REMOVED***
		return devices.devicePrefix + "-pool"
	***REMOVED***
	return devices.thinPoolDevice
***REMOVED***

func (devices *DeviceSet) getPoolDevName() string ***REMOVED***
	return getDevName(devices.getPoolName())
***REMOVED***

func (devices *DeviceSet) hasImage(name string) bool ***REMOVED***
	dirname := devices.loopbackDir()
	filename := path.Join(dirname, name)

	_, err := os.Stat(filename)
	return err == nil
***REMOVED***

// ensureImage creates a sparse file of <size> bytes at the path
// <root>/devicemapper/<name>.
// If the file already exists and new size is larger than its current size, it grows to the new size.
// Either way it returns the full path.
func (devices *DeviceSet) ensureImage(name string, size int64) (string, error) ***REMOVED***
	dirname := devices.loopbackDir()
	filename := path.Join(dirname, name)

	uid, gid, err := idtools.GetRootUIDGID(devices.uidMaps, devices.gidMaps)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if err := idtools.MkdirAllAndChown(dirname, 0700, idtools.IDPair***REMOVED***UID: uid, GID: gid***REMOVED***); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if fi, err := os.Stat(filename); err != nil ***REMOVED***
		if !os.IsNotExist(err) ***REMOVED***
			return "", err
		***REMOVED***
		logrus.Debugf("devmapper: Creating loopback file %s for device-manage use", filename)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		defer file.Close()

		if err := file.Truncate(size); err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if fi.Size() < size ***REMOVED***
			file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			defer file.Close()
			if err := file.Truncate(size); err != nil ***REMOVED***
				return "", fmt.Errorf("devmapper: Unable to grow loopback file %s: %v", filename, err)
			***REMOVED***
		***REMOVED*** else if fi.Size() > size ***REMOVED***
			logrus.Warnf("devmapper: Can't shrink loopback file %s", filename)
		***REMOVED***
	***REMOVED***
	return filename, nil
***REMOVED***

func (devices *DeviceSet) allocateTransactionID() uint64 ***REMOVED***
	devices.OpenTransactionID = devices.TransactionID + 1
	return devices.OpenTransactionID
***REMOVED***

func (devices *DeviceSet) updatePoolTransactionID() error ***REMOVED***
	if err := devicemapper.SetTransactionID(devices.getPoolDevName(), devices.TransactionID, devices.OpenTransactionID); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error setting devmapper transaction ID: %s", err)
	***REMOVED***
	devices.TransactionID = devices.OpenTransactionID
	return nil
***REMOVED***

func (devices *DeviceSet) removeMetadata(info *devInfo) error ***REMOVED***
	if err := os.RemoveAll(devices.metadataFile(info)); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error removing metadata file %s: %s", devices.metadataFile(info), err)
	***REMOVED***
	return nil
***REMOVED***

// Given json data and file path, write it to disk
func (devices *DeviceSet) writeMetaFile(jsonData []byte, filePath string) error ***REMOVED***
	tmpFile, err := ioutil.TempFile(devices.metadataDir(), ".tmp")
	if err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error creating metadata file: %s", err)
	***REMOVED***

	n, err := tmpFile.Write(jsonData)
	if err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error writing metadata to %s: %s", tmpFile.Name(), err)
	***REMOVED***
	if n < len(jsonData) ***REMOVED***
		return io.ErrShortWrite
	***REMOVED***
	if err := tmpFile.Sync(); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error syncing metadata file %s: %s", tmpFile.Name(), err)
	***REMOVED***
	if err := tmpFile.Close(); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error closing metadata file %s: %s", tmpFile.Name(), err)
	***REMOVED***
	if err := os.Rename(tmpFile.Name(), filePath); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error committing metadata file %s: %s", tmpFile.Name(), err)
	***REMOVED***

	return nil
***REMOVED***

func (devices *DeviceSet) saveMetadata(info *devInfo) error ***REMOVED***
	jsonData, err := json.Marshal(info)
	if err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error encoding metadata to json: %s", err)
	***REMOVED***
	return devices.writeMetaFile(jsonData, devices.metadataFile(info))
***REMOVED***

func (devices *DeviceSet) markDeviceIDUsed(deviceID int) ***REMOVED***
	var mask byte
	i := deviceID % 8
	mask = 1 << uint(i)
	devices.deviceIDMap[deviceID/8] = devices.deviceIDMap[deviceID/8] | mask
***REMOVED***

func (devices *DeviceSet) markDeviceIDFree(deviceID int) ***REMOVED***
	var mask byte
	i := deviceID % 8
	mask = ^(1 << uint(i))
	devices.deviceIDMap[deviceID/8] = devices.deviceIDMap[deviceID/8] & mask
***REMOVED***

func (devices *DeviceSet) isDeviceIDFree(deviceID int) bool ***REMOVED***
	var mask byte
	i := deviceID % 8
	mask = (1 << uint(i))
	return (devices.deviceIDMap[deviceID/8] & mask) == 0
***REMOVED***

// Should be called with devices.Lock() held.
func (devices *DeviceSet) lookupDevice(hash string) (*devInfo, error) ***REMOVED***
	info := devices.Devices[hash]
	if info == nil ***REMOVED***
		info = devices.loadMetadata(hash)
		if info == nil ***REMOVED***
			return nil, fmt.Errorf("devmapper: Unknown device %s", hash)
		***REMOVED***

		devices.Devices[hash] = info
	***REMOVED***
	return info, nil
***REMOVED***

func (devices *DeviceSet) lookupDeviceWithLock(hash string) (*devInfo, error) ***REMOVED***
	devices.Lock()
	defer devices.Unlock()
	info, err := devices.lookupDevice(hash)
	return info, err
***REMOVED***

// This function relies on that device hash map has been loaded in advance.
// Should be called with devices.Lock() held.
func (devices *DeviceSet) constructDeviceIDMap() ***REMOVED***
	logrus.Debug("devmapper: constructDeviceIDMap()")
	defer logrus.Debug("devmapper: constructDeviceIDMap() END")

	for _, info := range devices.Devices ***REMOVED***
		devices.markDeviceIDUsed(info.DeviceID)
		logrus.Debugf("devmapper: Added deviceId=%d to DeviceIdMap", info.DeviceID)
	***REMOVED***
***REMOVED***

func (devices *DeviceSet) deviceFileWalkFunction(path string, finfo os.FileInfo) error ***REMOVED***

	// Skip some of the meta files which are not device files.
	if strings.HasSuffix(finfo.Name(), ".migrated") ***REMOVED***
		logrus.Debugf("devmapper: Skipping file %s", path)
		return nil
	***REMOVED***

	if strings.HasPrefix(finfo.Name(), ".") ***REMOVED***
		logrus.Debugf("devmapper: Skipping file %s", path)
		return nil
	***REMOVED***

	if finfo.Name() == deviceSetMetaFile ***REMOVED***
		logrus.Debugf("devmapper: Skipping file %s", path)
		return nil
	***REMOVED***

	if finfo.Name() == transactionMetaFile ***REMOVED***
		logrus.Debugf("devmapper: Skipping file %s", path)
		return nil
	***REMOVED***

	logrus.Debugf("devmapper: Loading data for file %s", path)

	hash := finfo.Name()
	if hash == "base" ***REMOVED***
		hash = ""
	***REMOVED***

	// Include deleted devices also as cleanup delete device logic
	// will go through it and see if there are any deleted devices.
	if _, err := devices.lookupDevice(hash); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error looking up device %s:%v", hash, err)
	***REMOVED***

	return nil
***REMOVED***

func (devices *DeviceSet) loadDeviceFilesOnStart() error ***REMOVED***
	logrus.Debug("devmapper: loadDeviceFilesOnStart()")
	defer logrus.Debug("devmapper: loadDeviceFilesOnStart() END")

	var scan = func(path string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			logrus.Debugf("devmapper: Can't walk the file %s", path)
			return nil
		***REMOVED***

		// Skip any directories
		if info.IsDir() ***REMOVED***
			return nil
		***REMOVED***

		return devices.deviceFileWalkFunction(path, info)
	***REMOVED***

	return filepath.Walk(devices.metadataDir(), scan)
***REMOVED***

// Should be called with devices.Lock() held.
func (devices *DeviceSet) unregisterDevice(hash string) error ***REMOVED***
	logrus.Debugf("devmapper: unregisterDevice(%v)", hash)
	info := &devInfo***REMOVED***
		Hash: hash,
	***REMOVED***

	delete(devices.Devices, hash)

	if err := devices.removeMetadata(info); err != nil ***REMOVED***
		logrus.Debugf("devmapper: Error removing metadata: %s", err)
		return err
	***REMOVED***

	return nil
***REMOVED***

// Should be called with devices.Lock() held.
func (devices *DeviceSet) registerDevice(id int, hash string, size uint64, transactionID uint64) (*devInfo, error) ***REMOVED***
	logrus.Debugf("devmapper: registerDevice(%v, %v)", id, hash)
	info := &devInfo***REMOVED***
		Hash:          hash,
		DeviceID:      id,
		Size:          size,
		TransactionID: transactionID,
		Initialized:   false,
		devices:       devices,
	***REMOVED***

	devices.Devices[hash] = info

	if err := devices.saveMetadata(info); err != nil ***REMOVED***
		// Try to remove unused device
		delete(devices.Devices, hash)
		return nil, err
	***REMOVED***

	return info, nil
***REMOVED***

func (devices *DeviceSet) activateDeviceIfNeeded(info *devInfo, ignoreDeleted bool) error ***REMOVED***
	logrus.Debugf("devmapper: activateDeviceIfNeeded(%v)", info.Hash)

	if info.Deleted && !ignoreDeleted ***REMOVED***
		return fmt.Errorf("devmapper: Can't activate device %v as it is marked for deletion", info.Hash)
	***REMOVED***

	// Make sure deferred removal on device is canceled, if one was
	// scheduled.
	if err := devices.cancelDeferredRemovalIfNeeded(info); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Device Deferred Removal Cancellation Failed: %s", err)
	***REMOVED***

	if devinfo, _ := devicemapper.GetInfo(info.Name()); devinfo != nil && devinfo.Exists != 0 ***REMOVED***
		return nil
	***REMOVED***

	return devicemapper.ActivateDevice(devices.getPoolDevName(), info.Name(), info.DeviceID, info.Size)
***REMOVED***

// xfsSupported checks if xfs is supported, returns nil if it is, otherwise an error
func xfsSupported() error ***REMOVED***
	// Make sure mkfs.xfs is available
	if _, err := exec.LookPath("mkfs.xfs"); err != nil ***REMOVED***
		return err // error text is descriptive enough
	***REMOVED***

	// Check if kernel supports xfs filesystem or not.
	exec.Command("modprobe", "xfs").Run()

	f, err := os.Open("/proc/filesystems")
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "error checking for xfs support")
	***REMOVED***
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() ***REMOVED***
		if strings.HasSuffix(s.Text(), "\txfs") ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	if err := s.Err(); err != nil ***REMOVED***
		return errors.Wrapf(err, "error checking for xfs support")
	***REMOVED***

	return errors.New(`kernel does not support xfs, or "modprobe xfs" failed`)
***REMOVED***

func determineDefaultFS() string ***REMOVED***
	err := xfsSupported()
	if err == nil ***REMOVED***
		return "xfs"
	***REMOVED***

	logrus.Warnf("devmapper: XFS is not supported in your system (%v). Defaulting to ext4 filesystem", err)
	return "ext4"
***REMOVED***

// mkfsOptions tries to figure out whether some additional mkfs options are required
func mkfsOptions(fs string) []string ***REMOVED***
	if fs == "xfs" && !kernel.CheckKernelVersion(3, 16, 0) ***REMOVED***
		// For kernels earlier than 3.16 (and newer xfsutils),
		// some xfs features need to be explicitly disabled.
		return []string***REMOVED***"-m", "crc=0,finobt=0"***REMOVED***
	***REMOVED***

	return []string***REMOVED******REMOVED***
***REMOVED***

func (devices *DeviceSet) createFilesystem(info *devInfo) (err error) ***REMOVED***
	devname := info.DevName()

	if devices.filesystem == "" ***REMOVED***
		devices.filesystem = determineDefaultFS()
	***REMOVED***
	if err := devices.saveBaseDeviceFilesystem(devices.filesystem); err != nil ***REMOVED***
		return err
	***REMOVED***

	args := mkfsOptions(devices.filesystem)
	args = append(args, devices.mkfsArgs...)
	args = append(args, devname)

	logrus.Infof("devmapper: Creating filesystem %s on device %s, mkfs args: %v", devices.filesystem, info.Name(), args)
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			logrus.Infof("devmapper: Error while creating filesystem %s on device %s: %v", devices.filesystem, info.Name(), err)
		***REMOVED*** else ***REMOVED***
			logrus.Infof("devmapper: Successfully created filesystem %s on device %s", devices.filesystem, info.Name())
		***REMOVED***
	***REMOVED***()

	switch devices.filesystem ***REMOVED***
	case "xfs":
		err = exec.Command("mkfs.xfs", args...).Run()
	case "ext4":
		err = exec.Command("mkfs.ext4", append([]string***REMOVED***"-E", "nodiscard,lazy_itable_init=0,lazy_journal_init=0"***REMOVED***, args...)...).Run()
		if err != nil ***REMOVED***
			err = exec.Command("mkfs.ext4", append([]string***REMOVED***"-E", "nodiscard,lazy_itable_init=0"***REMOVED***, args...)...).Run()
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = exec.Command("tune2fs", append([]string***REMOVED***"-c", "-1", "-i", "0"***REMOVED***, devname)...).Run()
	default:
		err = fmt.Errorf("devmapper: Unsupported filesystem type %s", devices.filesystem)
	***REMOVED***
	return
***REMOVED***

func (devices *DeviceSet) migrateOldMetaData() error ***REMOVED***
	// Migrate old metadata file
	jsonData, err := ioutil.ReadFile(devices.oldMetadataFile())
	if err != nil && !os.IsNotExist(err) ***REMOVED***
		return err
	***REMOVED***

	if jsonData != nil ***REMOVED***
		m := metaData***REMOVED***Devices: make(map[string]*devInfo)***REMOVED***

		if err := json.Unmarshal(jsonData, &m); err != nil ***REMOVED***
			return err
		***REMOVED***

		for hash, info := range m.Devices ***REMOVED***
			info.Hash = hash
			devices.saveMetadata(info)
		***REMOVED***
		if err := os.Rename(devices.oldMetadataFile(), devices.oldMetadataFile()+".migrated"); err != nil ***REMOVED***
			return err
		***REMOVED***

	***REMOVED***

	return nil
***REMOVED***

// Cleanup deleted devices. It assumes that all the devices have been
// loaded in the hash table.
func (devices *DeviceSet) cleanupDeletedDevices() error ***REMOVED***
	devices.Lock()

	// If there are no deleted devices, there is nothing to do.
	if devices.nrDeletedDevices == 0 ***REMOVED***
		devices.Unlock()
		return nil
	***REMOVED***

	var deletedDevices []*devInfo

	for _, info := range devices.Devices ***REMOVED***
		if !info.Deleted ***REMOVED***
			continue
		***REMOVED***
		logrus.Debugf("devmapper: Found deleted device %s.", info.Hash)
		deletedDevices = append(deletedDevices, info)
	***REMOVED***

	// Delete the deleted devices. DeleteDevice() first takes the info lock
	// and then devices.Lock(). So drop it to avoid deadlock.
	devices.Unlock()

	for _, info := range deletedDevices ***REMOVED***
		// This will again try deferred deletion.
		if err := devices.DeleteDevice(info.Hash, false); err != nil ***REMOVED***
			logrus.Warnf("devmapper: Deletion of device %s, device_id=%v failed:%v", info.Hash, info.DeviceID, err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (devices *DeviceSet) countDeletedDevices() ***REMOVED***
	for _, info := range devices.Devices ***REMOVED***
		if !info.Deleted ***REMOVED***
			continue
		***REMOVED***
		devices.nrDeletedDevices++
	***REMOVED***
***REMOVED***

func (devices *DeviceSet) startDeviceDeletionWorker() ***REMOVED***
	// Deferred deletion is not enabled. Don't do anything.
	if !devices.deferredDelete ***REMOVED***
		return
	***REMOVED***

	logrus.Debug("devmapper: Worker to cleanup deleted devices started")
	for range devices.deletionWorkerTicker.C ***REMOVED***
		devices.cleanupDeletedDevices()
	***REMOVED***
***REMOVED***

func (devices *DeviceSet) initMetaData() error ***REMOVED***
	devices.Lock()
	defer devices.Unlock()

	if err := devices.migrateOldMetaData(); err != nil ***REMOVED***
		return err
	***REMOVED***

	_, transactionID, _, _, _, _, err := devices.poolStatus()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	devices.TransactionID = transactionID

	if err := devices.loadDeviceFilesOnStart(); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Failed to load device files:%v", err)
	***REMOVED***

	devices.constructDeviceIDMap()
	devices.countDeletedDevices()

	if err := devices.processPendingTransaction(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Start a goroutine to cleanup Deleted Devices
	go devices.startDeviceDeletionWorker()
	return nil
***REMOVED***

func (devices *DeviceSet) incNextDeviceID() ***REMOVED***
	// IDs are 24bit, so wrap around
	devices.NextDeviceID = (devices.NextDeviceID + 1) & maxDeviceID
***REMOVED***

func (devices *DeviceSet) getNextFreeDeviceID() (int, error) ***REMOVED***
	devices.incNextDeviceID()
	for i := 0; i <= maxDeviceID; i++ ***REMOVED***
		if devices.isDeviceIDFree(devices.NextDeviceID) ***REMOVED***
			devices.markDeviceIDUsed(devices.NextDeviceID)
			return devices.NextDeviceID, nil
		***REMOVED***
		devices.incNextDeviceID()
	***REMOVED***

	return 0, fmt.Errorf("devmapper: Unable to find a free device ID")
***REMOVED***

func (devices *DeviceSet) poolHasFreeSpace() error ***REMOVED***
	if devices.minFreeSpacePercent == 0 ***REMOVED***
		return nil
	***REMOVED***

	_, _, dataUsed, dataTotal, metadataUsed, metadataTotal, err := devices.poolStatus()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	minFreeData := (dataTotal * uint64(devices.minFreeSpacePercent)) / 100
	if minFreeData < 1 ***REMOVED***
		minFreeData = 1
	***REMOVED***
	dataFree := dataTotal - dataUsed
	if dataFree < minFreeData ***REMOVED***
		return fmt.Errorf("devmapper: Thin Pool has %v free data blocks which is less than minimum required %v free data blocks. Create more free space in thin pool or use dm.min_free_space option to change behavior", (dataTotal - dataUsed), minFreeData)
	***REMOVED***

	minFreeMetadata := (metadataTotal * uint64(devices.minFreeSpacePercent)) / 100
	if minFreeMetadata < 1 ***REMOVED***
		minFreeMetadata = 1
	***REMOVED***

	metadataFree := metadataTotal - metadataUsed
	if metadataFree < minFreeMetadata ***REMOVED***
		return fmt.Errorf("devmapper: Thin Pool has %v free metadata blocks which is less than minimum required %v free metadata blocks. Create more free metadata space in thin pool or use dm.min_free_space option to change behavior", (metadataTotal - metadataUsed), minFreeMetadata)
	***REMOVED***

	return nil
***REMOVED***

func (devices *DeviceSet) createRegisterDevice(hash string) (*devInfo, error) ***REMOVED***
	devices.Lock()
	defer devices.Unlock()

	deviceID, err := devices.getNextFreeDeviceID()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := devices.openTransaction(hash, deviceID); err != nil ***REMOVED***
		logrus.Debugf("devmapper: Error opening transaction hash = %s deviceID = %d", hash, deviceID)
		devices.markDeviceIDFree(deviceID)
		return nil, err
	***REMOVED***

	for ***REMOVED***
		if err := devicemapper.CreateDevice(devices.getPoolDevName(), deviceID); err != nil ***REMOVED***
			if devicemapper.DeviceIDExists(err) ***REMOVED***
				// Device ID already exists. This should not
				// happen. Now we have a mechanism to find
				// a free device ID. So something is not right.
				// Give a warning and continue.
				logrus.Errorf("devmapper: Device ID %d exists in pool but it is supposed to be unused", deviceID)
				deviceID, err = devices.getNextFreeDeviceID()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				// Save new device id into transaction
				devices.refreshTransaction(deviceID)
				continue
			***REMOVED***
			logrus.Debugf("devmapper: Error creating device: %s", err)
			devices.markDeviceIDFree(deviceID)
			return nil, err
		***REMOVED***
		break
	***REMOVED***

	logrus.Debugf("devmapper: Registering device (id %v) with FS size %v", deviceID, devices.baseFsSize)
	info, err := devices.registerDevice(deviceID, hash, devices.baseFsSize, devices.OpenTransactionID)
	if err != nil ***REMOVED***
		_ = devicemapper.DeleteDevice(devices.getPoolDevName(), deviceID)
		devices.markDeviceIDFree(deviceID)
		return nil, err
	***REMOVED***

	if err := devices.closeTransaction(); err != nil ***REMOVED***
		devices.unregisterDevice(hash)
		devicemapper.DeleteDevice(devices.getPoolDevName(), deviceID)
		devices.markDeviceIDFree(deviceID)
		return nil, err
	***REMOVED***
	return info, nil
***REMOVED***

func (devices *DeviceSet) takeSnapshot(hash string, baseInfo *devInfo, size uint64) error ***REMOVED***
	var (
		devinfo *devicemapper.Info
		err     error
	)

	if err = devices.poolHasFreeSpace(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if devices.deferredRemove ***REMOVED***
		devinfo, err = devicemapper.GetInfoWithDeferred(baseInfo.Name())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if devinfo != nil && devinfo.DeferredRemove != 0 ***REMOVED***
			err = devices.cancelDeferredRemoval(baseInfo)
			if err != nil ***REMOVED***
				// If Error is ErrEnxio. Device is probably already gone. Continue.
				if err != devicemapper.ErrEnxio ***REMOVED***
					return err
				***REMOVED***
				devinfo = nil
			***REMOVED*** else ***REMOVED***
				defer devices.deactivateDevice(baseInfo)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		devinfo, err = devicemapper.GetInfo(baseInfo.Name())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	doSuspend := devinfo != nil && devinfo.Exists != 0

	if doSuspend ***REMOVED***
		if err = devicemapper.SuspendDevice(baseInfo.Name()); err != nil ***REMOVED***
			return err
		***REMOVED***
		defer devicemapper.ResumeDevice(baseInfo.Name())
	***REMOVED***

	return devices.createRegisterSnapDevice(hash, baseInfo, size)
***REMOVED***

func (devices *DeviceSet) createRegisterSnapDevice(hash string, baseInfo *devInfo, size uint64) error ***REMOVED***
	deviceID, err := devices.getNextFreeDeviceID()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := devices.openTransaction(hash, deviceID); err != nil ***REMOVED***
		logrus.Debugf("devmapper: Error opening transaction hash = %s deviceID = %d", hash, deviceID)
		devices.markDeviceIDFree(deviceID)
		return err
	***REMOVED***

	for ***REMOVED***
		if err := devicemapper.CreateSnapDeviceRaw(devices.getPoolDevName(), deviceID, baseInfo.DeviceID); err != nil ***REMOVED***
			if devicemapper.DeviceIDExists(err) ***REMOVED***
				// Device ID already exists. This should not
				// happen. Now we have a mechanism to find
				// a free device ID. So something is not right.
				// Give a warning and continue.
				logrus.Errorf("devmapper: Device ID %d exists in pool but it is supposed to be unused", deviceID)
				deviceID, err = devices.getNextFreeDeviceID()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				// Save new device id into transaction
				devices.refreshTransaction(deviceID)
				continue
			***REMOVED***
			logrus.Debugf("devmapper: Error creating snap device: %s", err)
			devices.markDeviceIDFree(deviceID)
			return err
		***REMOVED***
		break
	***REMOVED***

	if _, err := devices.registerDevice(deviceID, hash, size, devices.OpenTransactionID); err != nil ***REMOVED***
		devicemapper.DeleteDevice(devices.getPoolDevName(), deviceID)
		devices.markDeviceIDFree(deviceID)
		logrus.Debugf("devmapper: Error registering device: %s", err)
		return err
	***REMOVED***

	if err := devices.closeTransaction(); err != nil ***REMOVED***
		devices.unregisterDevice(hash)
		devicemapper.DeleteDevice(devices.getPoolDevName(), deviceID)
		devices.markDeviceIDFree(deviceID)
		return err
	***REMOVED***
	return nil
***REMOVED***

func (devices *DeviceSet) loadMetadata(hash string) *devInfo ***REMOVED***
	info := &devInfo***REMOVED***Hash: hash, devices: devices***REMOVED***

	jsonData, err := ioutil.ReadFile(devices.metadataFile(info))
	if err != nil ***REMOVED***
		logrus.Debugf("devmapper: Failed to read %s with err: %v", devices.metadataFile(info), err)
		return nil
	***REMOVED***

	if err := json.Unmarshal(jsonData, &info); err != nil ***REMOVED***
		logrus.Debugf("devmapper: Failed to unmarshal devInfo from %s with err: %v", devices.metadataFile(info), err)
		return nil
	***REMOVED***

	if info.DeviceID > maxDeviceID ***REMOVED***
		logrus.Errorf("devmapper: Ignoring Invalid DeviceId=%d", info.DeviceID)
		return nil
	***REMOVED***

	return info
***REMOVED***

func getDeviceUUID(device string) (string, error) ***REMOVED***
	out, err := exec.Command("blkid", "-s", "UUID", "-o", "value", device).Output()
	if err != nil ***REMOVED***
		return "", fmt.Errorf("devmapper: Failed to find uuid for device %s:%v", device, err)
	***REMOVED***

	uuid := strings.TrimSuffix(string(out), "\n")
	uuid = strings.TrimSpace(uuid)
	logrus.Debugf("devmapper: UUID for device: %s is:%s", device, uuid)
	return uuid, nil
***REMOVED***

func (devices *DeviceSet) getBaseDeviceSize() uint64 ***REMOVED***
	info, _ := devices.lookupDevice("")
	if info == nil ***REMOVED***
		return 0
	***REMOVED***
	return info.Size
***REMOVED***

func (devices *DeviceSet) getBaseDeviceFS() string ***REMOVED***
	return devices.BaseDeviceFilesystem
***REMOVED***

func (devices *DeviceSet) verifyBaseDeviceUUIDFS(baseInfo *devInfo) error ***REMOVED***
	devices.Lock()
	defer devices.Unlock()

	if err := devices.activateDeviceIfNeeded(baseInfo, false); err != nil ***REMOVED***
		return err
	***REMOVED***
	defer devices.deactivateDevice(baseInfo)

	uuid, err := getDeviceUUID(baseInfo.DevName())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if devices.BaseDeviceUUID != uuid ***REMOVED***
		return fmt.Errorf("devmapper: Current Base Device UUID:%s does not match with stored UUID:%s. Possibly using a different thin pool than last invocation", uuid, devices.BaseDeviceUUID)
	***REMOVED***

	if devices.BaseDeviceFilesystem == "" ***REMOVED***
		fsType, err := ProbeFsType(baseInfo.DevName())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := devices.saveBaseDeviceFilesystem(fsType); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// If user specified a filesystem using dm.fs option and current
	// file system of base image is not same, warn user that dm.fs
	// will be ignored.
	if devices.BaseDeviceFilesystem != devices.filesystem ***REMOVED***
		logrus.Warnf("devmapper: Base device already exists and has filesystem %s on it. User specified filesystem %s will be ignored.", devices.BaseDeviceFilesystem, devices.filesystem)
		devices.filesystem = devices.BaseDeviceFilesystem
	***REMOVED***
	return nil
***REMOVED***

func (devices *DeviceSet) saveBaseDeviceFilesystem(fs string) error ***REMOVED***
	devices.BaseDeviceFilesystem = fs
	return devices.saveDeviceSetMetaData()
***REMOVED***

func (devices *DeviceSet) saveBaseDeviceUUID(baseInfo *devInfo) error ***REMOVED***
	devices.Lock()
	defer devices.Unlock()

	if err := devices.activateDeviceIfNeeded(baseInfo, false); err != nil ***REMOVED***
		return err
	***REMOVED***
	defer devices.deactivateDevice(baseInfo)

	uuid, err := getDeviceUUID(baseInfo.DevName())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	devices.BaseDeviceUUID = uuid
	return devices.saveDeviceSetMetaData()
***REMOVED***

func (devices *DeviceSet) createBaseImage() error ***REMOVED***
	logrus.Debug("devmapper: Initializing base device-mapper thin volume")

	// Create initial device
	info, err := devices.createRegisterDevice("")
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	logrus.Debug("devmapper: Creating filesystem on base device-mapper thin volume")

	if err := devices.activateDeviceIfNeeded(info, false); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := devices.createFilesystem(info); err != nil ***REMOVED***
		return err
	***REMOVED***

	info.Initialized = true
	if err := devices.saveMetadata(info); err != nil ***REMOVED***
		info.Initialized = false
		return err
	***REMOVED***

	if err := devices.saveBaseDeviceUUID(info); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Could not query and save base device UUID:%v", err)
	***REMOVED***

	return nil
***REMOVED***

// Returns if thin pool device exists or not. If device exists, also makes
// sure it is a thin pool device and not some other type of device.
func (devices *DeviceSet) thinPoolExists(thinPoolDevice string) (bool, error) ***REMOVED***
	logrus.Debugf("devmapper: Checking for existence of the pool %s", thinPoolDevice)

	info, err := devicemapper.GetInfo(thinPoolDevice)
	if err != nil ***REMOVED***
		return false, fmt.Errorf("devmapper: GetInfo() on device %s failed: %v", thinPoolDevice, err)
	***REMOVED***

	// Device does not exist.
	if info.Exists == 0 ***REMOVED***
		return false, nil
	***REMOVED***

	_, _, deviceType, _, err := devicemapper.GetStatus(thinPoolDevice)
	if err != nil ***REMOVED***
		return false, fmt.Errorf("devmapper: GetStatus() on device %s failed: %v", thinPoolDevice, err)
	***REMOVED***

	if deviceType != "thin-pool" ***REMOVED***
		return false, fmt.Errorf("devmapper: Device %s is not a thin pool", thinPoolDevice)
	***REMOVED***

	return true, nil
***REMOVED***

func (devices *DeviceSet) checkThinPool() error ***REMOVED***
	_, transactionID, dataUsed, _, _, _, err := devices.poolStatus()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if dataUsed != 0 ***REMOVED***
		return fmt.Errorf("devmapper: Unable to take ownership of thin-pool (%s) that already has used data blocks",
			devices.thinPoolDevice)
	***REMOVED***
	if transactionID != 0 ***REMOVED***
		return fmt.Errorf("devmapper: Unable to take ownership of thin-pool (%s) with non-zero transaction ID",
			devices.thinPoolDevice)
	***REMOVED***
	return nil
***REMOVED***

// Base image is initialized properly. Either save UUID for first time (for
// upgrade case or verify UUID.
func (devices *DeviceSet) setupVerifyBaseImageUUIDFS(baseInfo *devInfo) error ***REMOVED***
	// If BaseDeviceUUID is nil (upgrade case), save it and return success.
	if devices.BaseDeviceUUID == "" ***REMOVED***
		if err := devices.saveBaseDeviceUUID(baseInfo); err != nil ***REMOVED***
			return fmt.Errorf("devmapper: Could not query and save base device UUID:%v", err)
		***REMOVED***
		return nil
	***REMOVED***

	if err := devices.verifyBaseDeviceUUIDFS(baseInfo); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Base Device UUID and Filesystem verification failed: %v", err)
	***REMOVED***

	return nil
***REMOVED***

func (devices *DeviceSet) checkGrowBaseDeviceFS(info *devInfo) error ***REMOVED***

	if !userBaseSize ***REMOVED***
		return nil
	***REMOVED***

	if devices.baseFsSize < devices.getBaseDeviceSize() ***REMOVED***
		return fmt.Errorf("devmapper: Base device size cannot be smaller than %s", units.HumanSize(float64(devices.getBaseDeviceSize())))
	***REMOVED***

	if devices.baseFsSize == devices.getBaseDeviceSize() ***REMOVED***
		return nil
	***REMOVED***

	info.lock.Lock()
	defer info.lock.Unlock()

	devices.Lock()
	defer devices.Unlock()

	info.Size = devices.baseFsSize

	if err := devices.saveMetadata(info); err != nil ***REMOVED***
		// Try to remove unused device
		delete(devices.Devices, info.Hash)
		return err
	***REMOVED***

	return devices.growFS(info)
***REMOVED***

func (devices *DeviceSet) growFS(info *devInfo) error ***REMOVED***
	if err := devices.activateDeviceIfNeeded(info, false); err != nil ***REMOVED***
		return fmt.Errorf("Error activating devmapper device: %s", err)
	***REMOVED***

	defer devices.deactivateDevice(info)

	fsMountPoint := "/run/docker/mnt"
	if _, err := os.Stat(fsMountPoint); os.IsNotExist(err) ***REMOVED***
		if err := os.MkdirAll(fsMountPoint, 0700); err != nil ***REMOVED***
			return err
		***REMOVED***
		defer os.RemoveAll(fsMountPoint)
	***REMOVED***

	options := ""
	if devices.BaseDeviceFilesystem == "xfs" ***REMOVED***
		// XFS needs nouuid or it can't mount filesystems with the same fs
		options = joinMountOptions(options, "nouuid")
	***REMOVED***
	options = joinMountOptions(options, devices.mountOptions)

	if err := mount.Mount(info.DevName(), fsMountPoint, devices.BaseDeviceFilesystem, options); err != nil ***REMOVED***
		return fmt.Errorf("Error mounting '%s' on '%s' (fstype='%s' options='%s'): %s\n%v", info.DevName(), fsMountPoint, devices.BaseDeviceFilesystem, options, err, string(dmesg.Dmesg(256)))
	***REMOVED***

	defer unix.Unmount(fsMountPoint, unix.MNT_DETACH)

	switch devices.BaseDeviceFilesystem ***REMOVED***
	case "ext4":
		if out, err := exec.Command("resize2fs", info.DevName()).CombinedOutput(); err != nil ***REMOVED***
			return fmt.Errorf("Failed to grow rootfs:%v:%s", err, string(out))
		***REMOVED***
	case "xfs":
		if out, err := exec.Command("xfs_growfs", info.DevName()).CombinedOutput(); err != nil ***REMOVED***
			return fmt.Errorf("Failed to grow rootfs:%v:%s", err, string(out))
		***REMOVED***
	default:
		return fmt.Errorf("Unsupported filesystem type %s", devices.BaseDeviceFilesystem)
	***REMOVED***
	return nil
***REMOVED***

func (devices *DeviceSet) setupBaseImage() error ***REMOVED***
	oldInfo, _ := devices.lookupDeviceWithLock("")

	// base image already exists. If it is initialized properly, do UUID
	// verification and return. Otherwise remove image and set it up
	// fresh.

	if oldInfo != nil ***REMOVED***
		if oldInfo.Initialized && !oldInfo.Deleted ***REMOVED***
			if err := devices.setupVerifyBaseImageUUIDFS(oldInfo); err != nil ***REMOVED***
				return err
			***REMOVED***
			return devices.checkGrowBaseDeviceFS(oldInfo)
		***REMOVED***

		logrus.Debug("devmapper: Removing uninitialized base image")
		// If previous base device is in deferred delete state,
		// that needs to be cleaned up first. So don't try
		// deferred deletion.
		if err := devices.DeleteDevice("", true); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// If we are setting up base image for the first time, make sure
	// thin pool is empty.
	if devices.thinPoolDevice != "" && oldInfo == nil ***REMOVED***
		if err := devices.checkThinPool(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Create new base image device
	return devices.createBaseImage()
***REMOVED***

func setCloseOnExec(name string) ***REMOVED***
	fileInfos, _ := ioutil.ReadDir("/proc/self/fd")
	for _, i := range fileInfos ***REMOVED***
		link, _ := os.Readlink(filepath.Join("/proc/self/fd", i.Name()))
		if link == name ***REMOVED***
			fd, err := strconv.Atoi(i.Name())
			if err == nil ***REMOVED***
				unix.CloseOnExec(fd)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func major(device uint64) uint64 ***REMOVED***
	return (device >> 8) & 0xfff
***REMOVED***

func minor(device uint64) uint64 ***REMOVED***
	return (device & 0xff) | ((device >> 12) & 0xfff00)
***REMOVED***

// ResizePool increases the size of the pool.
func (devices *DeviceSet) ResizePool(size int64) error ***REMOVED***
	dirname := devices.loopbackDir()
	datafilename := path.Join(dirname, "data")
	if len(devices.dataDevice) > 0 ***REMOVED***
		datafilename = devices.dataDevice
	***REMOVED***
	metadatafilename := path.Join(dirname, "metadata")
	if len(devices.metadataDevice) > 0 ***REMOVED***
		metadatafilename = devices.metadataDevice
	***REMOVED***

	datafile, err := os.OpenFile(datafilename, os.O_RDWR, 0)
	if datafile == nil ***REMOVED***
		return err
	***REMOVED***
	defer datafile.Close()

	fi, err := datafile.Stat()
	if fi == nil ***REMOVED***
		return err
	***REMOVED***

	if fi.Size() > size ***REMOVED***
		return fmt.Errorf("devmapper: Can't shrink file")
	***REMOVED***

	dataloopback := loopback.FindLoopDeviceFor(datafile)
	if dataloopback == nil ***REMOVED***
		return fmt.Errorf("devmapper: Unable to find loopback mount for: %s", datafilename)
	***REMOVED***
	defer dataloopback.Close()

	metadatafile, err := os.OpenFile(metadatafilename, os.O_RDWR, 0)
	if metadatafile == nil ***REMOVED***
		return err
	***REMOVED***
	defer metadatafile.Close()

	metadataloopback := loopback.FindLoopDeviceFor(metadatafile)
	if metadataloopback == nil ***REMOVED***
		return fmt.Errorf("devmapper: Unable to find loopback mount for: %s", metadatafilename)
	***REMOVED***
	defer metadataloopback.Close()

	// Grow loopback file
	if err := datafile.Truncate(size); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Unable to grow loopback file: %s", err)
	***REMOVED***

	// Reload size for loopback device
	if err := loopback.SetCapacity(dataloopback); err != nil ***REMOVED***
		return fmt.Errorf("Unable to update loopback capacity: %s", err)
	***REMOVED***

	// Suspend the pool
	if err := devicemapper.SuspendDevice(devices.getPoolName()); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Unable to suspend pool: %s", err)
	***REMOVED***

	// Reload with the new block sizes
	if err := devicemapper.ReloadPool(devices.getPoolName(), dataloopback, metadataloopback, devices.thinpBlockSize); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Unable to reload pool: %s", err)
	***REMOVED***

	// Resume the pool
	if err := devicemapper.ResumeDevice(devices.getPoolName()); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Unable to resume pool: %s", err)
	***REMOVED***

	return nil
***REMOVED***

func (devices *DeviceSet) loadTransactionMetaData() error ***REMOVED***
	jsonData, err := ioutil.ReadFile(devices.transactionMetaFile())
	if err != nil ***REMOVED***
		// There is no active transaction. This will be the case
		// during upgrade.
		if os.IsNotExist(err) ***REMOVED***
			devices.OpenTransactionID = devices.TransactionID
			return nil
		***REMOVED***
		return err
	***REMOVED***

	json.Unmarshal(jsonData, &devices.transaction)
	return nil
***REMOVED***

func (devices *DeviceSet) saveTransactionMetaData() error ***REMOVED***
	jsonData, err := json.Marshal(&devices.transaction)
	if err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error encoding metadata to json: %s", err)
	***REMOVED***

	return devices.writeMetaFile(jsonData, devices.transactionMetaFile())
***REMOVED***

func (devices *DeviceSet) removeTransactionMetaData() error ***REMOVED***
	return os.RemoveAll(devices.transactionMetaFile())
***REMOVED***

func (devices *DeviceSet) rollbackTransaction() error ***REMOVED***
	logrus.Debugf("devmapper: Rolling back open transaction: TransactionID=%d hash=%s device_id=%d", devices.OpenTransactionID, devices.DeviceIDHash, devices.DeviceID)

	// A device id might have already been deleted before transaction
	// closed. In that case this call will fail. Just leave a message
	// in case of failure.
	if err := devicemapper.DeleteDevice(devices.getPoolDevName(), devices.DeviceID); err != nil ***REMOVED***
		logrus.Errorf("devmapper: Unable to delete device: %s", err)
	***REMOVED***

	dinfo := &devInfo***REMOVED***Hash: devices.DeviceIDHash***REMOVED***
	if err := devices.removeMetadata(dinfo); err != nil ***REMOVED***
		logrus.Errorf("devmapper: Unable to remove metadata: %s", err)
	***REMOVED*** else ***REMOVED***
		devices.markDeviceIDFree(devices.DeviceID)
	***REMOVED***

	if err := devices.removeTransactionMetaData(); err != nil ***REMOVED***
		logrus.Errorf("devmapper: Unable to remove transaction meta file %s: %s", devices.transactionMetaFile(), err)
	***REMOVED***

	return nil
***REMOVED***

func (devices *DeviceSet) processPendingTransaction() error ***REMOVED***
	if err := devices.loadTransactionMetaData(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// If there was open transaction but pool transaction ID is same
	// as open transaction ID, nothing to roll back.
	if devices.TransactionID == devices.OpenTransactionID ***REMOVED***
		return nil
	***REMOVED***

	// If open transaction ID is less than pool transaction ID, something
	// is wrong. Bail out.
	if devices.OpenTransactionID < devices.TransactionID ***REMOVED***
		logrus.Errorf("devmapper: Open Transaction id %d is less than pool transaction id %d", devices.OpenTransactionID, devices.TransactionID)
		return nil
	***REMOVED***

	// Pool transaction ID is not same as open transaction. There is
	// a transaction which was not completed.
	if err := devices.rollbackTransaction(); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Rolling back open transaction failed: %s", err)
	***REMOVED***

	devices.OpenTransactionID = devices.TransactionID
	return nil
***REMOVED***

func (devices *DeviceSet) loadDeviceSetMetaData() error ***REMOVED***
	jsonData, err := ioutil.ReadFile(devices.deviceSetMetaFile())
	if err != nil ***REMOVED***
		// For backward compatibility return success if file does
		// not exist.
		if os.IsNotExist(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	return json.Unmarshal(jsonData, devices)
***REMOVED***

func (devices *DeviceSet) saveDeviceSetMetaData() error ***REMOVED***
	jsonData, err := json.Marshal(devices)
	if err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error encoding metadata to json: %s", err)
	***REMOVED***

	return devices.writeMetaFile(jsonData, devices.deviceSetMetaFile())
***REMOVED***

func (devices *DeviceSet) openTransaction(hash string, DeviceID int) error ***REMOVED***
	devices.allocateTransactionID()
	devices.DeviceIDHash = hash
	devices.DeviceID = DeviceID
	if err := devices.saveTransactionMetaData(); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error saving transaction metadata: %s", err)
	***REMOVED***
	return nil
***REMOVED***

func (devices *DeviceSet) refreshTransaction(DeviceID int) error ***REMOVED***
	devices.DeviceID = DeviceID
	if err := devices.saveTransactionMetaData(); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error saving transaction metadata: %s", err)
	***REMOVED***
	return nil
***REMOVED***

func (devices *DeviceSet) closeTransaction() error ***REMOVED***
	if err := devices.updatePoolTransactionID(); err != nil ***REMOVED***
		logrus.Debug("devmapper: Failed to close Transaction")
		return err
	***REMOVED***
	return nil
***REMOVED***

func determineDriverCapabilities(version string) error ***REMOVED***
	// Kernel driver version >= 4.27.0 support deferred removal

	logrus.Debugf("devicemapper: kernel dm driver version is %s", version)

	versionSplit := strings.Split(version, ".")
	major, err := strconv.Atoi(versionSplit[0])
	if err != nil ***REMOVED***
		return graphdriver.ErrNotSupported
	***REMOVED***

	if major > 4 ***REMOVED***
		driverDeferredRemovalSupport = true
		return nil
	***REMOVED***

	if major < 4 ***REMOVED***
		return nil
	***REMOVED***

	minor, err := strconv.Atoi(versionSplit[1])
	if err != nil ***REMOVED***
		return graphdriver.ErrNotSupported
	***REMOVED***

	/*
	 * If major is 4 and minor is 27, then there is no need to
	 * check for patch level as it can not be less than 0.
	 */
	if minor >= 27 ***REMOVED***
		driverDeferredRemovalSupport = true
		return nil
	***REMOVED***

	return nil
***REMOVED***

// Determine the major and minor number of loopback device
func getDeviceMajorMinor(file *os.File) (uint64, uint64, error) ***REMOVED***
	var stat unix.Stat_t
	err := unix.Stat(file.Name(), &stat)
	if err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***

	dev := stat.Rdev
	majorNum := major(dev)
	minorNum := minor(dev)

	logrus.Debugf("devmapper: Major:Minor for device: %s is:%v:%v", file.Name(), majorNum, minorNum)
	return majorNum, minorNum, nil
***REMOVED***

// Given a file which is backing file of a loop back device, find the
// loopback device name and its major/minor number.
func getLoopFileDeviceMajMin(filename string) (string, uint64, uint64, error) ***REMOVED***
	file, err := os.Open(filename)
	if err != nil ***REMOVED***
		logrus.Debugf("devmapper: Failed to open file %s", filename)
		return "", 0, 0, err
	***REMOVED***

	defer file.Close()
	loopbackDevice := loopback.FindLoopDeviceFor(file)
	if loopbackDevice == nil ***REMOVED***
		return "", 0, 0, fmt.Errorf("devmapper: Unable to find loopback mount for: %s", filename)
	***REMOVED***
	defer loopbackDevice.Close()

	Major, Minor, err := getDeviceMajorMinor(loopbackDevice)
	if err != nil ***REMOVED***
		return "", 0, 0, err
	***REMOVED***
	return loopbackDevice.Name(), Major, Minor, nil
***REMOVED***

// Get the major/minor numbers of thin pool data and metadata devices
func (devices *DeviceSet) getThinPoolDataMetaMajMin() (uint64, uint64, uint64, uint64, error) ***REMOVED***
	var params, poolDataMajMin, poolMetadataMajMin string

	_, _, _, params, err := devicemapper.GetTable(devices.getPoolName())
	if err != nil ***REMOVED***
		return 0, 0, 0, 0, err
	***REMOVED***

	if _, err = fmt.Sscanf(params, "%s %s", &poolMetadataMajMin, &poolDataMajMin); err != nil ***REMOVED***
		return 0, 0, 0, 0, err
	***REMOVED***

	logrus.Debugf("devmapper: poolDataMajMin=%s poolMetaMajMin=%s\n", poolDataMajMin, poolMetadataMajMin)

	poolDataMajMinorSplit := strings.Split(poolDataMajMin, ":")
	poolDataMajor, err := strconv.ParseUint(poolDataMajMinorSplit[0], 10, 32)
	if err != nil ***REMOVED***
		return 0, 0, 0, 0, err
	***REMOVED***

	poolDataMinor, err := strconv.ParseUint(poolDataMajMinorSplit[1], 10, 32)
	if err != nil ***REMOVED***
		return 0, 0, 0, 0, err
	***REMOVED***

	poolMetadataMajMinorSplit := strings.Split(poolMetadataMajMin, ":")
	poolMetadataMajor, err := strconv.ParseUint(poolMetadataMajMinorSplit[0], 10, 32)
	if err != nil ***REMOVED***
		return 0, 0, 0, 0, err
	***REMOVED***

	poolMetadataMinor, err := strconv.ParseUint(poolMetadataMajMinorSplit[1], 10, 32)
	if err != nil ***REMOVED***
		return 0, 0, 0, 0, err
	***REMOVED***

	return poolDataMajor, poolDataMinor, poolMetadataMajor, poolMetadataMinor, nil
***REMOVED***

func (devices *DeviceSet) loadThinPoolLoopBackInfo() error ***REMOVED***
	poolDataMajor, poolDataMinor, poolMetadataMajor, poolMetadataMinor, err := devices.getThinPoolDataMetaMajMin()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	dirname := devices.loopbackDir()

	// data device has not been passed in. So there should be a data file
	// which is being mounted as loop device.
	if devices.dataDevice == "" ***REMOVED***
		datafilename := path.Join(dirname, "data")
		dataLoopDevice, dataMajor, dataMinor, err := getLoopFileDeviceMajMin(datafilename)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Compare the two
		if poolDataMajor == dataMajor && poolDataMinor == dataMinor ***REMOVED***
			devices.dataDevice = dataLoopDevice
			devices.dataLoopFile = datafilename
		***REMOVED***

	***REMOVED***

	// metadata device has not been passed in. So there should be a
	// metadata file which is being mounted as loop device.
	if devices.metadataDevice == "" ***REMOVED***
		metadatafilename := path.Join(dirname, "metadata")
		metadataLoopDevice, metadataMajor, metadataMinor, err := getLoopFileDeviceMajMin(metadatafilename)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if poolMetadataMajor == metadataMajor && poolMetadataMinor == metadataMinor ***REMOVED***
			devices.metadataDevice = metadataLoopDevice
			devices.metadataLoopFile = metadatafilename
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (devices *DeviceSet) enableDeferredRemovalDeletion() error ***REMOVED***

	// If user asked for deferred removal then check both libdm library
	// and kernel driver support deferred removal otherwise error out.
	if enableDeferredRemoval ***REMOVED***
		if !driverDeferredRemovalSupport ***REMOVED***
			return fmt.Errorf("devmapper: Deferred removal can not be enabled as kernel does not support it")
		***REMOVED***
		if !devicemapper.LibraryDeferredRemovalSupport ***REMOVED***
			return fmt.Errorf("devmapper: Deferred removal can not be enabled as libdm does not support it")
		***REMOVED***
		logrus.Debug("devmapper: Deferred removal support enabled.")
		devices.deferredRemove = true
	***REMOVED***

	if enableDeferredDeletion ***REMOVED***
		if !devices.deferredRemove ***REMOVED***
			return fmt.Errorf("devmapper: Deferred deletion can not be enabled as deferred removal is not enabled. Enable deferred removal using --storage-opt dm.use_deferred_removal=true parameter")
		***REMOVED***
		logrus.Debug("devmapper: Deferred deletion support enabled.")
		devices.deferredDelete = true
	***REMOVED***
	return nil
***REMOVED***

func (devices *DeviceSet) initDevmapper(doInit bool) (retErr error) ***REMOVED***
	if err := devices.enableDeferredRemovalDeletion(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// https://github.com/docker/docker/issues/4036
	if supported := devicemapper.UdevSetSyncSupport(true); !supported ***REMOVED***
		if dockerversion.IAmStatic == "true" ***REMOVED***
			logrus.Error("devmapper: Udev sync is not supported. This will lead to data loss and unexpected behavior. Install a dynamic binary to use devicemapper or select a different storage driver. For more information, see https://docs.docker.com/engine/reference/commandline/dockerd/#storage-driver-options")
		***REMOVED*** else ***REMOVED***
			logrus.Error("devmapper: Udev sync is not supported. This will lead to data loss and unexpected behavior. Install a more recent version of libdevmapper or select a different storage driver. For more information, see https://docs.docker.com/engine/reference/commandline/dockerd/#storage-driver-options")
		***REMOVED***

		if !devices.overrideUdevSyncCheck ***REMOVED***
			return graphdriver.ErrNotSupported
		***REMOVED***
	***REMOVED***

	//create the root dir of the devmapper driver ownership to match this
	//daemon's remapped root uid/gid so containers can start properly
	uid, gid, err := idtools.GetRootUIDGID(devices.uidMaps, devices.gidMaps)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := idtools.MkdirAndChown(devices.root, 0700, idtools.IDPair***REMOVED***UID: uid, GID: gid***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := os.MkdirAll(devices.metadataDir(), 0700); err != nil ***REMOVED***
		return err
	***REMOVED***

	prevSetupConfig, err := readLVMConfig(devices.root)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !reflect.DeepEqual(devices.lvmSetupConfig, directLVMConfig***REMOVED******REMOVED***) ***REMOVED***
		if devices.thinPoolDevice != "" ***REMOVED***
			return errors.New("cannot setup direct-lvm when `dm.thinpooldev` is also specified")
		***REMOVED***

		if !reflect.DeepEqual(prevSetupConfig, devices.lvmSetupConfig) ***REMOVED***
			if !reflect.DeepEqual(prevSetupConfig, directLVMConfig***REMOVED******REMOVED***) ***REMOVED***
				return errors.New("changing direct-lvm config is not supported")
			***REMOVED***
			logrus.WithField("storage-driver", "devicemapper").WithField("direct-lvm-config", devices.lvmSetupConfig).Debugf("Setting up direct lvm mode")
			if err := verifyBlockDevice(devices.lvmSetupConfig.Device, lvmSetupConfigForce); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := setupDirectLVM(devices.lvmSetupConfig); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := writeLVMConfig(devices.root, devices.lvmSetupConfig); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		devices.thinPoolDevice = "docker-thinpool"
		logrus.WithField("storage-driver", "devicemapper").Debugf("Setting dm.thinpooldev to %q", devices.thinPoolDevice)
	***REMOVED***

	// Set the device prefix from the device id and inode of the docker root dir
	var st unix.Stat_t
	if err := unix.Stat(devices.root, &st); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error looking up dir %s: %s", devices.root, err)
	***REMOVED***
	// "reg-" stands for "regular file".
	// In the future we might use "dev-" for "device file", etc.
	// docker-maj,min[-inode] stands for:
	//	- Managed by docker
	//	- The target of this device is at major <maj> and minor <min>
	//	- If <inode> is defined, use that file inside the device as a loopback image. Otherwise use the device itself.
	devices.devicePrefix = fmt.Sprintf("docker-%d:%d-%d", major(st.Dev), minor(st.Dev), st.Ino)
	logrus.Debugf("devmapper: Generated prefix: %s", devices.devicePrefix)

	// Check for the existence of the thin-pool device
	poolExists, err := devices.thinPoolExists(devices.getPoolName())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// It seems libdevmapper opens this without O_CLOEXEC, and go exec will not close files
	// that are not Close-on-exec,
	// so we add this badhack to make sure it closes itself
	setCloseOnExec("/dev/mapper/control")

	// Make sure the sparse images exist in <root>/devicemapper/data and
	// <root>/devicemapper/metadata

	createdLoopback := false

	// If the pool doesn't exist, create it
	if !poolExists && devices.thinPoolDevice == "" ***REMOVED***
		logrus.Debug("devmapper: Pool doesn't exist. Creating it.")

		var (
			dataFile     *os.File
			metadataFile *os.File
		)

		if devices.dataDevice == "" ***REMOVED***
			// Make sure the sparse images exist in <root>/devicemapper/data

			hasData := devices.hasImage("data")

			if !doInit && !hasData ***REMOVED***
				return errors.New("loopback data file not found")
			***REMOVED***

			if !hasData ***REMOVED***
				createdLoopback = true
			***REMOVED***

			data, err := devices.ensureImage("data", devices.dataLoopbackSize)
			if err != nil ***REMOVED***
				logrus.Debugf("devmapper: Error device ensureImage (data): %s", err)
				return err
			***REMOVED***

			dataFile, err = loopback.AttachLoopDevice(data)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			devices.dataLoopFile = data
			devices.dataDevice = dataFile.Name()
		***REMOVED*** else ***REMOVED***
			dataFile, err = os.OpenFile(devices.dataDevice, os.O_RDWR, 0600)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		defer dataFile.Close()

		if devices.metadataDevice == "" ***REMOVED***
			// Make sure the sparse images exist in <root>/devicemapper/metadata

			hasMetadata := devices.hasImage("metadata")

			if !doInit && !hasMetadata ***REMOVED***
				return errors.New("loopback metadata file not found")
			***REMOVED***

			if !hasMetadata ***REMOVED***
				createdLoopback = true
			***REMOVED***

			metadata, err := devices.ensureImage("metadata", devices.metaDataLoopbackSize)
			if err != nil ***REMOVED***
				logrus.Debugf("devmapper: Error device ensureImage (metadata): %s", err)
				return err
			***REMOVED***

			metadataFile, err = loopback.AttachLoopDevice(metadata)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			devices.metadataLoopFile = metadata
			devices.metadataDevice = metadataFile.Name()
		***REMOVED*** else ***REMOVED***
			metadataFile, err = os.OpenFile(devices.metadataDevice, os.O_RDWR, 0600)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		defer metadataFile.Close()

		if err := devicemapper.CreatePool(devices.getPoolName(), dataFile, metadataFile, devices.thinpBlockSize); err != nil ***REMOVED***
			return err
		***REMOVED***
		defer func() ***REMOVED***
			if retErr != nil ***REMOVED***
				err = devices.deactivatePool()
				if err != nil ***REMOVED***
					logrus.Warnf("devmapper: Failed to deactivatePool: %v", err)
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// Pool already exists and caller did not pass us a pool. That means
	// we probably created pool earlier and could not remove it as some
	// containers were still using it. Detect some of the properties of
	// pool, like is it using loop devices.
	if poolExists && devices.thinPoolDevice == "" ***REMOVED***
		if err := devices.loadThinPoolLoopBackInfo(); err != nil ***REMOVED***
			logrus.Debugf("devmapper: Failed to load thin pool loopback device information:%v", err)
			return err
		***REMOVED***
	***REMOVED***

	// If we didn't just create the data or metadata image, we need to
	// load the transaction id and migrate old metadata
	if !createdLoopback ***REMOVED***
		if err := devices.initMetaData(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if devices.thinPoolDevice == "" ***REMOVED***
		if devices.metadataLoopFile != "" || devices.dataLoopFile != "" ***REMOVED***
			logrus.Warn("devmapper: Usage of loopback devices is strongly discouraged for production use. Please use `--storage-opt dm.thinpooldev` or use `man dockerd` to refer to dm.thinpooldev section.")
		***REMOVED***
	***REMOVED***

	// Right now this loads only NextDeviceID. If there is more metadata
	// down the line, we might have to move it earlier.
	if err := devices.loadDeviceSetMetaData(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Setup the base image
	if doInit ***REMOVED***
		if err := devices.setupBaseImage(); err != nil ***REMOVED***
			logrus.Debugf("devmapper: Error device setupBaseImage: %s", err)
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// AddDevice adds a device and registers in the hash.
func (devices *DeviceSet) AddDevice(hash, baseHash string, storageOpt map[string]string) error ***REMOVED***
	logrus.Debugf("devmapper: AddDevice START(hash=%s basehash=%s)", hash, baseHash)
	defer logrus.Debugf("devmapper: AddDevice END(hash=%s basehash=%s)", hash, baseHash)

	// If a deleted device exists, return error.
	baseInfo, err := devices.lookupDeviceWithLock(baseHash)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if baseInfo.Deleted ***REMOVED***
		return fmt.Errorf("devmapper: Base device %v has been marked for deferred deletion", baseInfo.Hash)
	***REMOVED***

	baseInfo.lock.Lock()
	defer baseInfo.lock.Unlock()

	devices.Lock()
	defer devices.Unlock()

	// Also include deleted devices in case hash of new device is
	// same as one of the deleted devices.
	if info, _ := devices.lookupDevice(hash); info != nil ***REMOVED***
		return fmt.Errorf("devmapper: device %s already exists. Deleted=%v", hash, info.Deleted)
	***REMOVED***

	size, err := devices.parseStorageOpt(storageOpt)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if size == 0 ***REMOVED***
		size = baseInfo.Size
	***REMOVED***

	if size < baseInfo.Size ***REMOVED***
		return fmt.Errorf("devmapper: Container size cannot be smaller than %s", units.HumanSize(float64(baseInfo.Size)))
	***REMOVED***

	if err := devices.takeSnapshot(hash, baseInfo, size); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Grow the container rootfs.
	if size > baseInfo.Size ***REMOVED***
		info, err := devices.lookupDevice(hash)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := devices.growFS(info); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (devices *DeviceSet) parseStorageOpt(storageOpt map[string]string) (uint64, error) ***REMOVED***

	// Read size to change the block device size per container.
	for key, val := range storageOpt ***REMOVED***
		key := strings.ToLower(key)
		switch key ***REMOVED***
		case "size":
			size, err := units.RAMInBytes(val)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			return uint64(size), nil
		default:
			return 0, fmt.Errorf("Unknown option %s", key)
		***REMOVED***
	***REMOVED***

	return 0, nil
***REMOVED***

func (devices *DeviceSet) markForDeferredDeletion(info *devInfo) error ***REMOVED***
	// If device is already in deleted state, there is nothing to be done.
	if info.Deleted ***REMOVED***
		return nil
	***REMOVED***

	logrus.Debugf("devmapper: Marking device %s for deferred deletion.", info.Hash)

	info.Deleted = true

	// save device metadata to reflect deleted state.
	if err := devices.saveMetadata(info); err != nil ***REMOVED***
		info.Deleted = false
		return err
	***REMOVED***

	devices.nrDeletedDevices++
	return nil
***REMOVED***

// Should be called with devices.Lock() held.
func (devices *DeviceSet) deleteTransaction(info *devInfo, syncDelete bool) error ***REMOVED***
	if err := devices.openTransaction(info.Hash, info.DeviceID); err != nil ***REMOVED***
		logrus.Debugf("devmapper: Error opening transaction hash = %s deviceId = %d", "", info.DeviceID)
		return err
	***REMOVED***

	defer devices.closeTransaction()

	err := devicemapper.DeleteDevice(devices.getPoolDevName(), info.DeviceID)
	if err != nil ***REMOVED***
		// If syncDelete is true, we want to return error. If deferred
		// deletion is not enabled, we return an error. If error is
		// something other then EBUSY, return an error.
		if syncDelete || !devices.deferredDelete || err != devicemapper.ErrBusy ***REMOVED***
			logrus.Debugf("devmapper: Error deleting device: %s", err)
			return err
		***REMOVED***
	***REMOVED***

	if err == nil ***REMOVED***
		if err := devices.unregisterDevice(info.Hash); err != nil ***REMOVED***
			return err
		***REMOVED***
		// If device was already in deferred delete state that means
		// deletion was being tried again later. Reduce the deleted
		// device count.
		if info.Deleted ***REMOVED***
			devices.nrDeletedDevices--
		***REMOVED***
		devices.markDeviceIDFree(info.DeviceID)
	***REMOVED*** else ***REMOVED***
		if err := devices.markForDeferredDeletion(info); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Issue discard only if device open count is zero.
func (devices *DeviceSet) issueDiscard(info *devInfo) error ***REMOVED***
	logrus.Debugf("devmapper: issueDiscard START(device: %s).", info.Hash)
	defer logrus.Debugf("devmapper: issueDiscard END(device: %s).", info.Hash)
	// This is a workaround for the kernel not discarding block so
	// on the thin pool when we remove a thinp device, so we do it
	// manually.
	// Even if device is deferred deleted, activate it and issue
	// discards.
	if err := devices.activateDeviceIfNeeded(info, true); err != nil ***REMOVED***
		return err
	***REMOVED***

	devinfo, err := devicemapper.GetInfo(info.Name())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if devinfo.OpenCount != 0 ***REMOVED***
		logrus.Debugf("devmapper: Device: %s is in use. OpenCount=%d. Not issuing discards.", info.Hash, devinfo.OpenCount)
		return nil
	***REMOVED***

	if err := devicemapper.BlockDeviceDiscard(info.DevName()); err != nil ***REMOVED***
		logrus.Debugf("devmapper: Error discarding block on device: %s (ignoring)", err)
	***REMOVED***
	return nil
***REMOVED***

// Should be called with devices.Lock() held.
func (devices *DeviceSet) deleteDevice(info *devInfo, syncDelete bool) error ***REMOVED***
	if devices.doBlkDiscard ***REMOVED***
		devices.issueDiscard(info)
	***REMOVED***

	// Try to deactivate device in case it is active.
	// If deferred removal is enabled and deferred deletion is disabled
	// then make sure device is removed synchronously. There have been
	// some cases of device being busy for short duration and we would
	// rather busy wait for device removal to take care of these cases.
	deferredRemove := devices.deferredRemove
	if !devices.deferredDelete ***REMOVED***
		deferredRemove = false
	***REMOVED***

	if err := devices.deactivateDeviceMode(info, deferredRemove); err != nil ***REMOVED***
		logrus.Debugf("devmapper: Error deactivating device: %s", err)
		return err
	***REMOVED***

	return devices.deleteTransaction(info, syncDelete)
***REMOVED***

// DeleteDevice will return success if device has been marked for deferred
// removal. If one wants to override that and want DeleteDevice() to fail if
// device was busy and could not be deleted, set syncDelete=true.
func (devices *DeviceSet) DeleteDevice(hash string, syncDelete bool) error ***REMOVED***
	logrus.Debugf("devmapper: DeleteDevice START(hash=%v syncDelete=%v)", hash, syncDelete)
	defer logrus.Debugf("devmapper: DeleteDevice END(hash=%v syncDelete=%v)", hash, syncDelete)
	info, err := devices.lookupDeviceWithLock(hash)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	info.lock.Lock()
	defer info.lock.Unlock()

	devices.Lock()
	defer devices.Unlock()

	return devices.deleteDevice(info, syncDelete)
***REMOVED***

func (devices *DeviceSet) deactivatePool() error ***REMOVED***
	logrus.Debug("devmapper: deactivatePool() START")
	defer logrus.Debug("devmapper: deactivatePool() END")
	devname := devices.getPoolDevName()

	devinfo, err := devicemapper.GetInfo(devname)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if devinfo.Exists == 0 ***REMOVED***
		return nil
	***REMOVED***
	if err := devicemapper.RemoveDevice(devname); err != nil ***REMOVED***
		return err
	***REMOVED***

	if d, err := devicemapper.GetDeps(devname); err == nil ***REMOVED***
		logrus.Warnf("devmapper: device %s still has %d active dependents", devname, d.Count)
	***REMOVED***

	return nil
***REMOVED***

func (devices *DeviceSet) deactivateDevice(info *devInfo) error ***REMOVED***
	return devices.deactivateDeviceMode(info, devices.deferredRemove)
***REMOVED***

func (devices *DeviceSet) deactivateDeviceMode(info *devInfo, deferredRemove bool) error ***REMOVED***
	var err error
	logrus.Debugf("devmapper: deactivateDevice START(%s)", info.Hash)
	defer logrus.Debugf("devmapper: deactivateDevice END(%s)", info.Hash)

	devinfo, err := devicemapper.GetInfo(info.Name())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if devinfo.Exists == 0 ***REMOVED***
		return nil
	***REMOVED***

	if deferredRemove ***REMOVED***
		err = devicemapper.RemoveDeviceDeferred(info.Name())
	***REMOVED*** else ***REMOVED***
		err = devices.removeDevice(info.Name())
	***REMOVED***

	// This function's semantics is such that it does not return an
	// error if device does not exist. So if device went away by
	// the time we actually tried to remove it, do not return error.
	if err != devicemapper.ErrEnxio ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Issues the underlying dm remove operation.
func (devices *DeviceSet) removeDevice(devname string) error ***REMOVED***
	var err error

	logrus.Debugf("devmapper: removeDevice START(%s)", devname)
	defer logrus.Debugf("devmapper: removeDevice END(%s)", devname)

	for i := 0; i < 200; i++ ***REMOVED***
		err = devicemapper.RemoveDevice(devname)
		if err == nil ***REMOVED***
			break
		***REMOVED***
		if err != devicemapper.ErrBusy ***REMOVED***
			return err
		***REMOVED***

		// If we see EBUSY it may be a transient error,
		// sleep a bit a retry a few times.
		devices.Unlock()
		time.Sleep(100 * time.Millisecond)
		devices.Lock()
	***REMOVED***

	return err
***REMOVED***

func (devices *DeviceSet) cancelDeferredRemovalIfNeeded(info *devInfo) error ***REMOVED***
	if !devices.deferredRemove ***REMOVED***
		return nil
	***REMOVED***

	logrus.Debugf("devmapper: cancelDeferredRemovalIfNeeded START(%s)", info.Name())
	defer logrus.Debugf("devmapper: cancelDeferredRemovalIfNeeded END(%s)", info.Name())

	devinfo, err := devicemapper.GetInfoWithDeferred(info.Name())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if devinfo != nil && devinfo.DeferredRemove == 0 ***REMOVED***
		return nil
	***REMOVED***

	// Cancel deferred remove
	if err := devices.cancelDeferredRemoval(info); err != nil ***REMOVED***
		// If Error is ErrEnxio. Device is probably already gone. Continue.
		if err != devicemapper.ErrEnxio ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (devices *DeviceSet) cancelDeferredRemoval(info *devInfo) error ***REMOVED***
	logrus.Debugf("devmapper: cancelDeferredRemoval START(%s)", info.Name())
	defer logrus.Debugf("devmapper: cancelDeferredRemoval END(%s)", info.Name())

	var err error

	// Cancel deferred remove
	for i := 0; i < 100; i++ ***REMOVED***
		err = devicemapper.CancelDeferredRemove(info.Name())
		if err != nil ***REMOVED***
			if err == devicemapper.ErrBusy ***REMOVED***
				// If we see EBUSY it may be a transient error,
				// sleep a bit a retry a few times.
				devices.Unlock()
				time.Sleep(100 * time.Millisecond)
				devices.Lock()
				continue
			***REMOVED***
		***REMOVED***
		break
	***REMOVED***
	return err
***REMOVED***

// Shutdown shuts down the device by unmounting the root.
func (devices *DeviceSet) Shutdown(home string) error ***REMOVED***
	logrus.Debugf("devmapper: [deviceset %s] Shutdown()", devices.devicePrefix)
	logrus.Debugf("devmapper: Shutting down DeviceSet: %s", devices.root)
	defer logrus.Debugf("devmapper: [deviceset %s] Shutdown() END", devices.devicePrefix)

	// Stop deletion worker. This should start delivering new events to
	// ticker channel. That means no new instance of cleanupDeletedDevice()
	// will run after this call. If one instance is already running at
	// the time of the call, it must be holding devices.Lock() and
	// we will block on this lock till cleanup function exits.
	devices.deletionWorkerTicker.Stop()

	devices.Lock()
	// Save DeviceSet Metadata first. Docker kills all threads if they
	// don't finish in certain time. It is possible that Shutdown()
	// routine does not finish in time as we loop trying to deactivate
	// some devices while these are busy. In that case shutdown() routine
	// will be killed and we will not get a chance to save deviceset
	// metadata. Hence save this early before trying to deactivate devices.
	devices.saveDeviceSetMetaData()

	// ignore the error since it's just a best effort to not try to unmount something that's mounted
	mounts, _ := mount.GetMounts()
	mounted := make(map[string]bool, len(mounts))
	for _, mnt := range mounts ***REMOVED***
		mounted[mnt.Mountpoint] = true
	***REMOVED***

	if err := filepath.Walk(path.Join(home, "mnt"), func(p string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if !info.IsDir() ***REMOVED***
			return nil
		***REMOVED***

		if mounted[p] ***REMOVED***
			// We use MNT_DETACH here in case it is still busy in some running
			// container. This means it'll go away from the global scope directly,
			// and the device will be released when that container dies.
			if err := unix.Unmount(p, unix.MNT_DETACH); err != nil ***REMOVED***
				logrus.Debugf("devmapper: Shutdown unmounting %s, error: %s", p, err)
			***REMOVED***
		***REMOVED***

		if devInfo, err := devices.lookupDevice(path.Base(p)); err != nil ***REMOVED***
			logrus.Debugf("devmapper: Shutdown lookup device %s, error: %s", path.Base(p), err)
		***REMOVED*** else ***REMOVED***
			if err := devices.deactivateDevice(devInfo); err != nil ***REMOVED***
				logrus.Debugf("devmapper: Shutdown deactivate %s , error: %s", devInfo.Hash, err)
			***REMOVED***
		***REMOVED***

		return nil
	***REMOVED***); err != nil && !os.IsNotExist(err) ***REMOVED***
		devices.Unlock()
		return err
	***REMOVED***

	devices.Unlock()

	info, _ := devices.lookupDeviceWithLock("")
	if info != nil ***REMOVED***
		info.lock.Lock()
		devices.Lock()
		if err := devices.deactivateDevice(info); err != nil ***REMOVED***
			logrus.Debugf("devmapper: Shutdown deactivate base , error: %s", err)
		***REMOVED***
		devices.Unlock()
		info.lock.Unlock()
	***REMOVED***

	devices.Lock()
	if devices.thinPoolDevice == "" ***REMOVED***
		if err := devices.deactivatePool(); err != nil ***REMOVED***
			logrus.Debugf("devmapper: Shutdown deactivate pool , error: %s", err)
		***REMOVED***
	***REMOVED***
	devices.Unlock()

	return nil
***REMOVED***

// Recent XFS changes allow changing behavior of filesystem in case of errors.
// When thin pool gets full and XFS gets ENOSPC error, currently it tries
// IO infinitely and sometimes it can block the container process
// and process can't be killWith 0 value, XFS will not retry upon error
// and instead will shutdown filesystem.

func (devices *DeviceSet) xfsSetNospaceRetries(info *devInfo) error ***REMOVED***
	dmDevicePath, err := os.Readlink(info.DevName())
	if err != nil ***REMOVED***
		return fmt.Errorf("devmapper: readlink failed for device %v:%v", info.DevName(), err)
	***REMOVED***

	dmDeviceName := path.Base(dmDevicePath)
	filePath := "/sys/fs/xfs/" + dmDeviceName + "/error/metadata/ENOSPC/max_retries"
	maxRetriesFile, err := os.OpenFile(filePath, os.O_WRONLY, 0)
	if err != nil ***REMOVED***
		return fmt.Errorf("devmapper: user specified daemon option dm.xfs_nospace_max_retries but it does not seem to be supported on this system :%v", err)
	***REMOVED***
	defer maxRetriesFile.Close()

	// Set max retries to 0
	_, err = maxRetriesFile.WriteString(devices.xfsNospaceRetries)
	if err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Failed to write string %v to file %v:%v", devices.xfsNospaceRetries, filePath, err)
	***REMOVED***
	return nil
***REMOVED***

// MountDevice mounts the device if not already mounted.
func (devices *DeviceSet) MountDevice(hash, path, mountLabel string) error ***REMOVED***
	info, err := devices.lookupDeviceWithLock(hash)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if info.Deleted ***REMOVED***
		return fmt.Errorf("devmapper: Can't mount device %v as it has been marked for deferred deletion", info.Hash)
	***REMOVED***

	info.lock.Lock()
	defer info.lock.Unlock()

	devices.Lock()
	defer devices.Unlock()

	if err := devices.activateDeviceIfNeeded(info, false); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error activating devmapper device for '%s': %s", hash, err)
	***REMOVED***

	fstype, err := ProbeFsType(info.DevName())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	options := ""

	if fstype == "xfs" ***REMOVED***
		// XFS needs nouuid or it can't mount filesystems with the same fs
		options = joinMountOptions(options, "nouuid")
	***REMOVED***

	options = joinMountOptions(options, devices.mountOptions)
	options = joinMountOptions(options, label.FormatMountLabel("", mountLabel))

	if err := mount.Mount(info.DevName(), path, fstype, options); err != nil ***REMOVED***
		return fmt.Errorf("devmapper: Error mounting '%s' on '%s' (fstype='%s' options='%s'): %s\n%v", info.DevName(), path, fstype, options, err, string(dmesg.Dmesg(256)))
	***REMOVED***

	if fstype == "xfs" && devices.xfsNospaceRetries != "" ***REMOVED***
		if err := devices.xfsSetNospaceRetries(info); err != nil ***REMOVED***
			unix.Unmount(path, unix.MNT_DETACH)
			devices.deactivateDevice(info)
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// UnmountDevice unmounts the device and removes it from hash.
func (devices *DeviceSet) UnmountDevice(hash, mountPath string) error ***REMOVED***
	logrus.Debugf("devmapper: UnmountDevice START(hash=%s)", hash)
	defer logrus.Debugf("devmapper: UnmountDevice END(hash=%s)", hash)

	info, err := devices.lookupDeviceWithLock(hash)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	info.lock.Lock()
	defer info.lock.Unlock()

	devices.Lock()
	defer devices.Unlock()

	logrus.Debugf("devmapper: Unmount(%s)", mountPath)
	if err := unix.Unmount(mountPath, unix.MNT_DETACH); err != nil ***REMOVED***
		return err
	***REMOVED***
	logrus.Debug("devmapper: Unmount done")

	// Remove the mountpoint here. Removing the mountpoint (in newer kernels)
	// will cause all other instances of this mount in other mount namespaces
	// to be killed (this is an anti-DoS measure that is necessary for things
	// like devicemapper). This is necessary to avoid cases where a libdm mount
	// that is present in another namespace will cause subsequent RemoveDevice
	// operations to fail. We ignore any errors here because this may fail on
	// older kernels which don't have
	// torvalds/linux@8ed936b5671bfb33d89bc60bdcc7cf0470ba52fe applied.
	if err := os.Remove(mountPath); err != nil ***REMOVED***
		logrus.Debugf("devmapper: error doing a remove on unmounted device %s: %v", mountPath, err)
	***REMOVED***

	return devices.deactivateDevice(info)
***REMOVED***

// HasDevice returns true if the device metadata exists.
func (devices *DeviceSet) HasDevice(hash string) bool ***REMOVED***
	info, _ := devices.lookupDeviceWithLock(hash)
	return info != nil
***REMOVED***

// List returns a list of device ids.
func (devices *DeviceSet) List() []string ***REMOVED***
	devices.Lock()
	defer devices.Unlock()

	ids := make([]string, len(devices.Devices))
	i := 0
	for k := range devices.Devices ***REMOVED***
		ids[i] = k
		i++
	***REMOVED***
	return ids
***REMOVED***

func (devices *DeviceSet) deviceStatus(devName string) (sizeInSectors, mappedSectors, highestMappedSector uint64, err error) ***REMOVED***
	var params string
	_, sizeInSectors, _, params, err = devicemapper.GetStatus(devName)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if _, err = fmt.Sscanf(params, "%d %d", &mappedSectors, &highestMappedSector); err == nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// GetDeviceStatus provides size, mapped sectors
func (devices *DeviceSet) GetDeviceStatus(hash string) (*DevStatus, error) ***REMOVED***
	info, err := devices.lookupDeviceWithLock(hash)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	info.lock.Lock()
	defer info.lock.Unlock()

	devices.Lock()
	defer devices.Unlock()

	status := &DevStatus***REMOVED***
		DeviceID:      info.DeviceID,
		Size:          info.Size,
		TransactionID: info.TransactionID,
	***REMOVED***

	if err := devices.activateDeviceIfNeeded(info, false); err != nil ***REMOVED***
		return nil, fmt.Errorf("devmapper: Error activating devmapper device for '%s': %s", hash, err)
	***REMOVED***

	sizeInSectors, mappedSectors, highestMappedSector, err := devices.deviceStatus(info.DevName())

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	status.SizeInSectors = sizeInSectors
	status.MappedSectors = mappedSectors
	status.HighestMappedSector = highestMappedSector

	return status, nil
***REMOVED***

func (devices *DeviceSet) poolStatus() (totalSizeInSectors, transactionID, dataUsed, dataTotal, metadataUsed, metadataTotal uint64, err error) ***REMOVED***
	var params string
	if _, totalSizeInSectors, _, params, err = devicemapper.GetStatus(devices.getPoolName()); err == nil ***REMOVED***
		_, err = fmt.Sscanf(params, "%d %d/%d %d/%d", &transactionID, &metadataUsed, &metadataTotal, &dataUsed, &dataTotal)
	***REMOVED***
	return
***REMOVED***

// DataDevicePath returns the path to the data storage for this deviceset,
// regardless of loopback or block device
func (devices *DeviceSet) DataDevicePath() string ***REMOVED***
	return devices.dataDevice
***REMOVED***

// MetadataDevicePath returns the path to the metadata storage for this deviceset,
// regardless of loopback or block device
func (devices *DeviceSet) MetadataDevicePath() string ***REMOVED***
	return devices.metadataDevice
***REMOVED***

func (devices *DeviceSet) getUnderlyingAvailableSpace(loopFile string) (uint64, error) ***REMOVED***
	buf := new(unix.Statfs_t)
	if err := unix.Statfs(loopFile, buf); err != nil ***REMOVED***
		logrus.Warnf("devmapper: Couldn't stat loopfile filesystem %v: %v", loopFile, err)
		return 0, err
	***REMOVED***
	return buf.Bfree * uint64(buf.Bsize), nil
***REMOVED***

func (devices *DeviceSet) isRealFile(loopFile string) (bool, error) ***REMOVED***
	if loopFile != "" ***REMOVED***
		fi, err := os.Stat(loopFile)
		if err != nil ***REMOVED***
			logrus.Warnf("devmapper: Couldn't stat loopfile %v: %v", loopFile, err)
			return false, err
		***REMOVED***
		return fi.Mode().IsRegular(), nil
	***REMOVED***
	return false, nil
***REMOVED***

// Status returns the current status of this deviceset
func (devices *DeviceSet) Status() *Status ***REMOVED***
	devices.Lock()
	defer devices.Unlock()

	status := &Status***REMOVED******REMOVED***

	status.PoolName = devices.getPoolName()
	status.DataFile = devices.DataDevicePath()
	status.DataLoopback = devices.dataLoopFile
	status.MetadataFile = devices.MetadataDevicePath()
	status.MetadataLoopback = devices.metadataLoopFile
	status.UdevSyncSupported = devicemapper.UdevSyncSupported()
	status.DeferredRemoveEnabled = devices.deferredRemove
	status.DeferredDeleteEnabled = devices.deferredDelete
	status.DeferredDeletedDeviceCount = devices.nrDeletedDevices
	status.BaseDeviceSize = devices.getBaseDeviceSize()
	status.BaseDeviceFS = devices.getBaseDeviceFS()

	totalSizeInSectors, _, dataUsed, dataTotal, metadataUsed, metadataTotal, err := devices.poolStatus()
	if err == nil ***REMOVED***
		// Convert from blocks to bytes
		blockSizeInSectors := totalSizeInSectors / dataTotal

		status.Data.Used = dataUsed * blockSizeInSectors * 512
		status.Data.Total = dataTotal * blockSizeInSectors * 512
		status.Data.Available = status.Data.Total - status.Data.Used

		// metadata blocks are always 4k
		status.Metadata.Used = metadataUsed * 4096
		status.Metadata.Total = metadataTotal * 4096
		status.Metadata.Available = status.Metadata.Total - status.Metadata.Used

		status.SectorSize = blockSizeInSectors * 512

		if check, _ := devices.isRealFile(devices.dataLoopFile); check ***REMOVED***
			actualSpace, err := devices.getUnderlyingAvailableSpace(devices.dataLoopFile)
			if err == nil && actualSpace < status.Data.Available ***REMOVED***
				status.Data.Available = actualSpace
			***REMOVED***
		***REMOVED***

		if check, _ := devices.isRealFile(devices.metadataLoopFile); check ***REMOVED***
			actualSpace, err := devices.getUnderlyingAvailableSpace(devices.metadataLoopFile)
			if err == nil && actualSpace < status.Metadata.Available ***REMOVED***
				status.Metadata.Available = actualSpace
			***REMOVED***
		***REMOVED***

		minFreeData := (dataTotal * uint64(devices.minFreeSpacePercent)) / 100
		status.MinFreeSpace = minFreeData * blockSizeInSectors * 512
	***REMOVED***

	return status
***REMOVED***

// Status returns the current status of this deviceset
func (devices *DeviceSet) exportDeviceMetadata(hash string) (*deviceMetadata, error) ***REMOVED***
	info, err := devices.lookupDeviceWithLock(hash)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	info.lock.Lock()
	defer info.lock.Unlock()

	metadata := &deviceMetadata***REMOVED***info.DeviceID, info.Size, info.Name()***REMOVED***
	return metadata, nil
***REMOVED***

// NewDeviceSet creates the device set based on the options provided.
func NewDeviceSet(root string, doInit bool, options []string, uidMaps, gidMaps []idtools.IDMap) (*DeviceSet, error) ***REMOVED***
	devicemapper.SetDevDir("/dev")

	devices := &DeviceSet***REMOVED***
		root:                  root,
		metaData:              metaData***REMOVED***Devices: make(map[string]*devInfo)***REMOVED***,
		dataLoopbackSize:      defaultDataLoopbackSize,
		metaDataLoopbackSize:  defaultMetaDataLoopbackSize,
		baseFsSize:            defaultBaseFsSize,
		overrideUdevSyncCheck: defaultUdevSyncOverride,
		doBlkDiscard:          true,
		thinpBlockSize:        defaultThinpBlockSize,
		deviceIDMap:           make([]byte, deviceIDMapSz),
		deletionWorkerTicker:  time.NewTicker(time.Second * 30),
		uidMaps:               uidMaps,
		gidMaps:               gidMaps,
		minFreeSpacePercent:   defaultMinFreeSpacePercent,
	***REMOVED***

	version, err := devicemapper.GetDriverVersion()
	if err != nil ***REMOVED***
		// Can't even get driver version, assume not supported
		return nil, graphdriver.ErrNotSupported
	***REMOVED***

	if err := determineDriverCapabilities(version); err != nil ***REMOVED***
		return nil, graphdriver.ErrNotSupported
	***REMOVED***

	if driverDeferredRemovalSupport && devicemapper.LibraryDeferredRemovalSupport ***REMOVED***
		// enable deferred stuff by default
		enableDeferredDeletion = true
		enableDeferredRemoval = true
	***REMOVED***

	foundBlkDiscard := false
	var lvmSetupConfig directLVMConfig
	for _, option := range options ***REMOVED***
		key, val, err := parsers.ParseKeyValueOpt(option)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		key = strings.ToLower(key)
		switch key ***REMOVED***
		case "dm.basesize":
			size, err := units.RAMInBytes(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			userBaseSize = true
			devices.baseFsSize = uint64(size)
		case "dm.loopdatasize":
			size, err := units.RAMInBytes(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			devices.dataLoopbackSize = size
		case "dm.loopmetadatasize":
			size, err := units.RAMInBytes(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			devices.metaDataLoopbackSize = size
		case "dm.fs":
			if val != "ext4" && val != "xfs" ***REMOVED***
				return nil, fmt.Errorf("devmapper: Unsupported filesystem %s", val)
			***REMOVED***
			devices.filesystem = val
		case "dm.mkfsarg":
			devices.mkfsArgs = append(devices.mkfsArgs, val)
		case "dm.mountopt":
			devices.mountOptions = joinMountOptions(devices.mountOptions, val)
		case "dm.metadatadev":
			devices.metadataDevice = val
		case "dm.datadev":
			devices.dataDevice = val
		case "dm.thinpooldev":
			devices.thinPoolDevice = strings.TrimPrefix(val, "/dev/mapper/")
		case "dm.blkdiscard":
			foundBlkDiscard = true
			devices.doBlkDiscard, err = strconv.ParseBool(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		case "dm.blocksize":
			size, err := units.RAMInBytes(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			// convert to 512b sectors
			devices.thinpBlockSize = uint32(size) >> 9
		case "dm.override_udev_sync_check":
			devices.overrideUdevSyncCheck, err = strconv.ParseBool(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

		case "dm.use_deferred_removal":
			enableDeferredRemoval, err = strconv.ParseBool(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

		case "dm.use_deferred_deletion":
			enableDeferredDeletion, err = strconv.ParseBool(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

		case "dm.min_free_space":
			if !strings.HasSuffix(val, "%") ***REMOVED***
				return nil, fmt.Errorf("devmapper: Option dm.min_free_space requires %% suffix")
			***REMOVED***

			valstring := strings.TrimSuffix(val, "%")
			minFreeSpacePercent, err := strconv.ParseUint(valstring, 10, 32)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if minFreeSpacePercent >= 100 ***REMOVED***
				return nil, fmt.Errorf("devmapper: Invalid value %v for option dm.min_free_space", val)
			***REMOVED***

			devices.minFreeSpacePercent = uint32(minFreeSpacePercent)
		case "dm.xfs_nospace_max_retries":
			_, err := strconv.ParseUint(val, 10, 64)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			devices.xfsNospaceRetries = val
		case "dm.directlvm_device":
			lvmSetupConfig.Device = val
		case "dm.directlvm_device_force":
			lvmSetupConfigForce, err = strconv.ParseBool(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		case "dm.thinp_percent":
			per, err := strconv.ParseUint(strings.TrimSuffix(val, "%"), 10, 32)
			if err != nil ***REMOVED***
				return nil, errors.Wrapf(err, "could not parse `dm.thinp_percent=%s`", val)
			***REMOVED***
			if per >= 100 ***REMOVED***
				return nil, errors.New("dm.thinp_percent must be greater than 0 and less than 100")
			***REMOVED***
			lvmSetupConfig.ThinpPercent = per
		case "dm.thinp_metapercent":
			per, err := strconv.ParseUint(strings.TrimSuffix(val, "%"), 10, 32)
			if err != nil ***REMOVED***
				return nil, errors.Wrapf(err, "could not parse `dm.thinp_metapercent=%s`", val)
			***REMOVED***
			if per >= 100 ***REMOVED***
				return nil, errors.New("dm.thinp_metapercent must be greater than 0 and less than 100")
			***REMOVED***
			lvmSetupConfig.ThinpMetaPercent = per
		case "dm.thinp_autoextend_percent":
			per, err := strconv.ParseUint(strings.TrimSuffix(val, "%"), 10, 32)
			if err != nil ***REMOVED***
				return nil, errors.Wrapf(err, "could not parse `dm.thinp_autoextend_percent=%s`", val)
			***REMOVED***
			if per > 100 ***REMOVED***
				return nil, errors.New("dm.thinp_autoextend_percent must be greater than 0 and less than 100")
			***REMOVED***
			lvmSetupConfig.AutoExtendPercent = per
		case "dm.thinp_autoextend_threshold":
			per, err := strconv.ParseUint(strings.TrimSuffix(val, "%"), 10, 32)
			if err != nil ***REMOVED***
				return nil, errors.Wrapf(err, "could not parse `dm.thinp_autoextend_threshold=%s`", val)
			***REMOVED***
			if per > 100 ***REMOVED***
				return nil, errors.New("dm.thinp_autoextend_threshold must be greater than 0 and less than 100")
			***REMOVED***
			lvmSetupConfig.AutoExtendThreshold = per
		case "dm.libdm_log_level":
			level, err := strconv.ParseInt(val, 10, 32)
			if err != nil ***REMOVED***
				return nil, errors.Wrapf(err, "could not parse `dm.libdm_log_level=%s`", val)
			***REMOVED***
			if level < devicemapper.LogLevelFatal || level > devicemapper.LogLevelDebug ***REMOVED***
				return nil, errors.Errorf("dm.libdm_log_level must be in range [%d,%d]", devicemapper.LogLevelFatal, devicemapper.LogLevelDebug)
			***REMOVED***
			// Register a new logging callback with the specified level.
			devicemapper.LogInit(devicemapper.DefaultLogger***REMOVED***
				Level: int(level),
			***REMOVED***)
		default:
			return nil, fmt.Errorf("devmapper: Unknown option %s", key)
		***REMOVED***
	***REMOVED***

	if err := validateLVMConfig(lvmSetupConfig); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	devices.lvmSetupConfig = lvmSetupConfig

	// By default, don't do blk discard hack on raw devices, its rarely useful and is expensive
	if !foundBlkDiscard && (devices.dataDevice != "" || devices.thinPoolDevice != "") ***REMOVED***
		devices.doBlkDiscard = false
	***REMOVED***

	if err := devices.initDevmapper(doInit); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return devices, nil
***REMOVED***
