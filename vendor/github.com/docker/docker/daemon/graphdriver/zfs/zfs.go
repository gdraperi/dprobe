// +build linux freebsd

package zfs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/parsers"
	zfs "github.com/mistifyio/go-zfs"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type zfsOptions struct ***REMOVED***
	fsName    string
	mountPath string
***REMOVED***

func init() ***REMOVED***
	graphdriver.Register("zfs", Init)
***REMOVED***

// Logger returns a zfs logger implementation.
type Logger struct***REMOVED******REMOVED***

// Log wraps log message from ZFS driver with a prefix '[zfs]'.
func (*Logger) Log(cmd []string) ***REMOVED***
	logrus.Debugf("[zfs] %s", strings.Join(cmd, " "))
***REMOVED***

// Init returns a new ZFS driver.
// It takes base mount path and an array of options which are represented as key value pairs.
// Each option is in the for key=value. 'zfs.fsname' is expected to be a valid key in the options.
func Init(base string, opt []string, uidMaps, gidMaps []idtools.IDMap) (graphdriver.Driver, error) ***REMOVED***
	var err error

	if _, err := exec.LookPath("zfs"); err != nil ***REMOVED***
		logrus.Debugf("[zfs] zfs command is not available: %v", err)
		return nil, graphdriver.ErrPrerequisites
	***REMOVED***

	file, err := os.OpenFile("/dev/zfs", os.O_RDWR, 600)
	if err != nil ***REMOVED***
		logrus.Debugf("[zfs] cannot open /dev/zfs: %v", err)
		return nil, graphdriver.ErrPrerequisites
	***REMOVED***
	defer file.Close()

	options, err := parseOptions(opt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	options.mountPath = base

	rootdir := path.Dir(base)

	if options.fsName == "" ***REMOVED***
		err = checkRootdirFs(rootdir)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if options.fsName == "" ***REMOVED***
		options.fsName, err = lookupZfsDataset(rootdir)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	zfs.SetLogger(new(Logger))

	filesystems, err := zfs.Filesystems(options.fsName)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Cannot find root filesystem %s: %v", options.fsName, err)
	***REMOVED***

	filesystemsCache := make(map[string]bool, len(filesystems))
	var rootDataset *zfs.Dataset
	for _, fs := range filesystems ***REMOVED***
		if fs.Name == options.fsName ***REMOVED***
			rootDataset = fs
		***REMOVED***
		filesystemsCache[fs.Name] = true
	***REMOVED***

	if rootDataset == nil ***REMOVED***
		return nil, fmt.Errorf("BUG: zfs get all -t filesystem -rHp '%s' should contain '%s'", options.fsName, options.fsName)
	***REMOVED***

	rootUID, rootGID, err := idtools.GetRootUIDGID(uidMaps, gidMaps)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Failed to get root uid/guid: %v", err)
	***REMOVED***
	if err := idtools.MkdirAllAndChown(base, 0700, idtools.IDPair***REMOVED***rootUID, rootGID***REMOVED***); err != nil ***REMOVED***
		return nil, fmt.Errorf("Failed to create '%s': %v", base, err)
	***REMOVED***

	d := &Driver***REMOVED***
		dataset:          rootDataset,
		options:          options,
		filesystemsCache: filesystemsCache,
		uidMaps:          uidMaps,
		gidMaps:          gidMaps,
		ctr:              graphdriver.NewRefCounter(graphdriver.NewDefaultChecker()),
	***REMOVED***
	return graphdriver.NewNaiveDiffDriver(d, uidMaps, gidMaps), nil
***REMOVED***

func parseOptions(opt []string) (zfsOptions, error) ***REMOVED***
	var options zfsOptions
	options.fsName = ""
	for _, option := range opt ***REMOVED***
		key, val, err := parsers.ParseKeyValueOpt(option)
		if err != nil ***REMOVED***
			return options, err
		***REMOVED***
		key = strings.ToLower(key)
		switch key ***REMOVED***
		case "zfs.fsname":
			options.fsName = val
		default:
			return options, fmt.Errorf("Unknown option %s", key)
		***REMOVED***
	***REMOVED***
	return options, nil
***REMOVED***

func lookupZfsDataset(rootdir string) (string, error) ***REMOVED***
	var stat unix.Stat_t
	if err := unix.Stat(rootdir, &stat); err != nil ***REMOVED***
		return "", fmt.Errorf("Failed to access '%s': %s", rootdir, err)
	***REMOVED***
	wantedDev := stat.Dev

	mounts, err := mount.GetMounts()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	for _, m := range mounts ***REMOVED***
		if err := unix.Stat(m.Mountpoint, &stat); err != nil ***REMOVED***
			logrus.Debugf("[zfs] failed to stat '%s' while scanning for zfs mount: %v", m.Mountpoint, err)
			continue // may fail on fuse file systems
		***REMOVED***

		if stat.Dev == wantedDev && m.Fstype == "zfs" ***REMOVED***
			return m.Source, nil
		***REMOVED***
	***REMOVED***

	return "", fmt.Errorf("Failed to find zfs dataset mounted on '%s' in /proc/mounts", rootdir)
***REMOVED***

// Driver holds information about the driver, such as zfs dataset, options and cache.
type Driver struct ***REMOVED***
	dataset          *zfs.Dataset
	options          zfsOptions
	sync.Mutex       // protects filesystem cache against concurrent access
	filesystemsCache map[string]bool
	uidMaps          []idtools.IDMap
	gidMaps          []idtools.IDMap
	ctr              *graphdriver.RefCounter
***REMOVED***

func (d *Driver) String() string ***REMOVED***
	return "zfs"
***REMOVED***

// Cleanup is called on daemon shutdown, it is used to clean up any remaining mounts
func (d *Driver) Cleanup() error ***REMOVED***
	return mount.RecursiveUnmount(d.options.mountPath)
***REMOVED***

// Status returns information about the ZFS filesystem. It returns a two dimensional array of information
// such as pool name, dataset name, disk usage, parent quota and compression used.
// Currently it return 'Zpool', 'Zpool Health', 'Parent Dataset', 'Space Used By Parent',
// 'Space Available', 'Parent Quota' and 'Compression'.
func (d *Driver) Status() [][2]string ***REMOVED***
	parts := strings.Split(d.dataset.Name, "/")
	pool, err := zfs.GetZpool(parts[0])

	var poolName, poolHealth string
	if err == nil ***REMOVED***
		poolName = pool.Name
		poolHealth = pool.Health
	***REMOVED*** else ***REMOVED***
		poolName = fmt.Sprintf("error while getting pool information %v", err)
		poolHealth = "not available"
	***REMOVED***

	quota := "no"
	if d.dataset.Quota != 0 ***REMOVED***
		quota = strconv.FormatUint(d.dataset.Quota, 10)
	***REMOVED***

	return [][2]string***REMOVED***
		***REMOVED***"Zpool", poolName***REMOVED***,
		***REMOVED***"Zpool Health", poolHealth***REMOVED***,
		***REMOVED***"Parent Dataset", d.dataset.Name***REMOVED***,
		***REMOVED***"Space Used By Parent", strconv.FormatUint(d.dataset.Used, 10)***REMOVED***,
		***REMOVED***"Space Available", strconv.FormatUint(d.dataset.Avail, 10)***REMOVED***,
		***REMOVED***"Parent Quota", quota***REMOVED***,
		***REMOVED***"Compression", d.dataset.Compression***REMOVED***,
	***REMOVED***
***REMOVED***

// GetMetadata returns image/container metadata related to graph driver
func (d *Driver) GetMetadata(id string) (map[string]string, error) ***REMOVED***
	return map[string]string***REMOVED***
		"Mountpoint": d.mountPath(id),
		"Dataset":    d.zfsPath(id),
	***REMOVED***, nil
***REMOVED***

func (d *Driver) cloneFilesystem(name, parentName string) error ***REMOVED***
	snapshotName := fmt.Sprintf("%d", time.Now().Nanosecond())
	parentDataset := zfs.Dataset***REMOVED***Name: parentName***REMOVED***
	snapshot, err := parentDataset.Snapshot(snapshotName /*recursive */, false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	_, err = snapshot.Clone(name, map[string]string***REMOVED***"mountpoint": "legacy"***REMOVED***)
	if err == nil ***REMOVED***
		d.Lock()
		d.filesystemsCache[name] = true
		d.Unlock()
	***REMOVED***

	if err != nil ***REMOVED***
		snapshot.Destroy(zfs.DestroyDeferDeletion)
		return err
	***REMOVED***
	return snapshot.Destroy(zfs.DestroyDeferDeletion)
***REMOVED***

func (d *Driver) zfsPath(id string) string ***REMOVED***
	return d.options.fsName + "/" + id
***REMOVED***

func (d *Driver) mountPath(id string) string ***REMOVED***
	return path.Join(d.options.mountPath, "graph", getMountpoint(id))
***REMOVED***

// CreateReadWrite creates a layer that is writable for use as a container
// file system.
func (d *Driver) CreateReadWrite(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	return d.Create(id, parent, opts)
***REMOVED***

// Create prepares the dataset and filesystem for the ZFS driver for the given id under the parent.
func (d *Driver) Create(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	var storageOpt map[string]string
	if opts != nil ***REMOVED***
		storageOpt = opts.StorageOpt
	***REMOVED***

	err := d.create(id, parent, storageOpt)
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	if zfsError, ok := err.(*zfs.Error); ok ***REMOVED***
		if !strings.HasSuffix(zfsError.Stderr, "dataset already exists\n") ***REMOVED***
			return err
		***REMOVED***
		// aborted build -> cleanup
	***REMOVED*** else ***REMOVED***
		return err
	***REMOVED***

	dataset := zfs.Dataset***REMOVED***Name: d.zfsPath(id)***REMOVED***
	if err := dataset.Destroy(zfs.DestroyRecursiveClones); err != nil ***REMOVED***
		return err
	***REMOVED***

	// retry
	return d.create(id, parent, storageOpt)
***REMOVED***

func (d *Driver) create(id, parent string, storageOpt map[string]string) error ***REMOVED***
	name := d.zfsPath(id)
	quota, err := parseStorageOpt(storageOpt)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if parent == "" ***REMOVED***
		mountoptions := map[string]string***REMOVED***"mountpoint": "legacy"***REMOVED***
		fs, err := zfs.CreateFilesystem(name, mountoptions)
		if err == nil ***REMOVED***
			err = setQuota(name, quota)
			if err == nil ***REMOVED***
				d.Lock()
				d.filesystemsCache[fs.Name] = true
				d.Unlock()
			***REMOVED***
		***REMOVED***
		return err
	***REMOVED***
	err = d.cloneFilesystem(name, d.zfsPath(parent))
	if err == nil ***REMOVED***
		err = setQuota(name, quota)
	***REMOVED***
	return err
***REMOVED***

func parseStorageOpt(storageOpt map[string]string) (string, error) ***REMOVED***
	// Read size to change the disk quota per container
	for k, v := range storageOpt ***REMOVED***
		key := strings.ToLower(k)
		switch key ***REMOVED***
		case "size":
			return v, nil
		default:
			return "0", fmt.Errorf("Unknown option %s", key)
		***REMOVED***
	***REMOVED***
	return "0", nil
***REMOVED***

func setQuota(name string, quota string) error ***REMOVED***
	if quota == "0" ***REMOVED***
		return nil
	***REMOVED***
	fs, err := zfs.GetDataset(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return fs.SetProperty("quota", quota)
***REMOVED***

// Remove deletes the dataset, filesystem and the cache for the given id.
func (d *Driver) Remove(id string) error ***REMOVED***
	name := d.zfsPath(id)
	dataset := zfs.Dataset***REMOVED***Name: name***REMOVED***
	err := dataset.Destroy(zfs.DestroyRecursive)
	if err == nil ***REMOVED***
		d.Lock()
		delete(d.filesystemsCache, name)
		d.Unlock()
	***REMOVED***
	return err
***REMOVED***

// Get returns the mountpoint for the given id after creating the target directories if necessary.
func (d *Driver) Get(id, mountLabel string) (_ containerfs.ContainerFS, retErr error) ***REMOVED***
	mountpoint := d.mountPath(id)
	if count := d.ctr.Increment(mountpoint); count > 1 ***REMOVED***
		return containerfs.NewLocalContainerFS(mountpoint), nil
	***REMOVED***
	defer func() ***REMOVED***
		if retErr != nil ***REMOVED***
			if c := d.ctr.Decrement(mountpoint); c <= 0 ***REMOVED***
				if mntErr := unix.Unmount(mountpoint, 0); mntErr != nil ***REMOVED***
					logrus.Errorf("Error unmounting %v: %v", mountpoint, mntErr)
				***REMOVED***
				if rmErr := unix.Rmdir(mountpoint); rmErr != nil && !os.IsNotExist(rmErr) ***REMOVED***
					logrus.Debugf("Failed to remove %s: %v", id, rmErr)
				***REMOVED***

			***REMOVED***
		***REMOVED***
	***REMOVED***()

	filesystem := d.zfsPath(id)
	options := label.FormatMountLabel("", mountLabel)
	logrus.Debugf(`[zfs] mount("%s", "%s", "%s")`, filesystem, mountpoint, options)

	rootUID, rootGID, err := idtools.GetRootUIDGID(d.uidMaps, d.gidMaps)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Create the target directories if they don't exist
	if err := idtools.MkdirAllAndChown(mountpoint, 0755, idtools.IDPair***REMOVED***rootUID, rootGID***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := mount.Mount(filesystem, mountpoint, "zfs", options); err != nil ***REMOVED***
		return nil, fmt.Errorf("error creating zfs mount of %s to %s: %v", filesystem, mountpoint, err)
	***REMOVED***

	// this could be our first mount after creation of the filesystem, and the root dir may still have root
	// permissions instead of the remapped root uid:gid (if user namespaces are enabled):
	if err := os.Chown(mountpoint, rootUID, rootGID); err != nil ***REMOVED***
		return nil, fmt.Errorf("error modifying zfs mountpoint (%s) directory ownership: %v", mountpoint, err)
	***REMOVED***

	return containerfs.NewLocalContainerFS(mountpoint), nil
***REMOVED***

// Put removes the existing mountpoint for the given id if it exists.
func (d *Driver) Put(id string) error ***REMOVED***
	mountpoint := d.mountPath(id)
	if count := d.ctr.Decrement(mountpoint); count > 0 ***REMOVED***
		return nil
	***REMOVED***

	logrus.Debugf(`[zfs] unmount("%s")`, mountpoint)

	if err := unix.Unmount(mountpoint, unix.MNT_DETACH); err != nil ***REMOVED***
		logrus.Warnf("Failed to unmount %s mount %s: %v", id, mountpoint, err)
	***REMOVED***
	if err := unix.Rmdir(mountpoint); err != nil && !os.IsNotExist(err) ***REMOVED***
		logrus.Debugf("Failed to remove %s mount point %s: %v", id, mountpoint, err)
	***REMOVED***

	return nil
***REMOVED***

// Exists checks to see if the cache entry exists for the given id.
func (d *Driver) Exists(id string) bool ***REMOVED***
	d.Lock()
	defer d.Unlock()
	return d.filesystemsCache[d.zfsPath(id)]
***REMOVED***
