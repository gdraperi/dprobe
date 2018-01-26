// +build linux

package btrfs

/*
#include <stdlib.h>
#include <dirent.h>
#include <btrfs/ioctl.h>
#include <btrfs/ctree.h>

static void set_name_btrfs_ioctl_vol_args_v2(struct btrfs_ioctl_vol_args_v2* btrfs_struct, const char* value) ***REMOVED***
    snprintf(btrfs_struct->name, BTRFS_SUBVOL_NAME_MAX, "%s", value);
***REMOVED***
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/parsers"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/go-units"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func init() ***REMOVED***
	graphdriver.Register("btrfs", Init)
***REMOVED***

type btrfsOptions struct ***REMOVED***
	minSpace uint64
	size     uint64
***REMOVED***

// Init returns a new BTRFS driver.
// An error is returned if BTRFS is not supported.
func Init(home string, options []string, uidMaps, gidMaps []idtools.IDMap) (graphdriver.Driver, error) ***REMOVED***

	// Perform feature detection on /var/lib/docker/btrfs if it's an existing directory.
	// This covers situations where /var/lib/docker/btrfs is a mount, and on a different
	// filesystem than /var/lib/docker.
	// If the path does not exist, fall back to using /var/lib/docker for feature detection.
	testdir := home
	if _, err := os.Stat(testdir); os.IsNotExist(err) ***REMOVED***
		testdir = filepath.Dir(testdir)
	***REMOVED***

	fsMagic, err := graphdriver.GetFSMagic(testdir)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if fsMagic != graphdriver.FsMagicBtrfs ***REMOVED***
		return nil, graphdriver.ErrPrerequisites
	***REMOVED***

	rootUID, rootGID, err := idtools.GetRootUIDGID(uidMaps, gidMaps)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := idtools.MkdirAllAndChown(home, 0700, idtools.IDPair***REMOVED***UID: rootUID, GID: rootGID***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	opt, userDiskQuota, err := parseOptions(options)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	driver := &Driver***REMOVED***
		home:    home,
		uidMaps: uidMaps,
		gidMaps: gidMaps,
		options: opt,
	***REMOVED***

	if userDiskQuota ***REMOVED***
		if err := driver.subvolEnableQuota(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return graphdriver.NewNaiveDiffDriver(driver, uidMaps, gidMaps), nil
***REMOVED***

func parseOptions(opt []string) (btrfsOptions, bool, error) ***REMOVED***
	var options btrfsOptions
	userDiskQuota := false
	for _, option := range opt ***REMOVED***
		key, val, err := parsers.ParseKeyValueOpt(option)
		if err != nil ***REMOVED***
			return options, userDiskQuota, err
		***REMOVED***
		key = strings.ToLower(key)
		switch key ***REMOVED***
		case "btrfs.min_space":
			minSpace, err := units.RAMInBytes(val)
			if err != nil ***REMOVED***
				return options, userDiskQuota, err
			***REMOVED***
			userDiskQuota = true
			options.minSpace = uint64(minSpace)
		default:
			return options, userDiskQuota, fmt.Errorf("Unknown option %s", key)
		***REMOVED***
	***REMOVED***
	return options, userDiskQuota, nil
***REMOVED***

// Driver contains information about the filesystem mounted.
type Driver struct ***REMOVED***
	//root of the file system
	home         string
	uidMaps      []idtools.IDMap
	gidMaps      []idtools.IDMap
	options      btrfsOptions
	quotaEnabled bool
	once         sync.Once
***REMOVED***

// String prints the name of the driver (btrfs).
func (d *Driver) String() string ***REMOVED***
	return "btrfs"
***REMOVED***

// Status returns current driver information in a two dimensional string array.
// Output contains "Build Version" and "Library Version" of the btrfs libraries used.
// Version information can be used to check compatibility with your kernel.
func (d *Driver) Status() [][2]string ***REMOVED***
	status := [][2]string***REMOVED******REMOVED***
	if bv := btrfsBuildVersion(); bv != "-" ***REMOVED***
		status = append(status, [2]string***REMOVED***"Build Version", bv***REMOVED***)
	***REMOVED***
	if lv := btrfsLibVersion(); lv != -1 ***REMOVED***
		status = append(status, [2]string***REMOVED***"Library Version", fmt.Sprintf("%d", lv)***REMOVED***)
	***REMOVED***
	return status
***REMOVED***

// GetMetadata returns empty metadata for this driver.
func (d *Driver) GetMetadata(id string) (map[string]string, error) ***REMOVED***
	return nil, nil
***REMOVED***

// Cleanup unmounts the home directory.
func (d *Driver) Cleanup() error ***REMOVED***
	if err := d.subvolDisableQuota(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return mount.RecursiveUnmount(d.home)
***REMOVED***

func free(p *C.char) ***REMOVED***
	C.free(unsafe.Pointer(p))
***REMOVED***

func openDir(path string) (*C.DIR, error) ***REMOVED***
	Cpath := C.CString(path)
	defer free(Cpath)

	dir := C.opendir(Cpath)
	if dir == nil ***REMOVED***
		return nil, fmt.Errorf("Can't open dir")
	***REMOVED***
	return dir, nil
***REMOVED***

func closeDir(dir *C.DIR) ***REMOVED***
	if dir != nil ***REMOVED***
		C.closedir(dir)
	***REMOVED***
***REMOVED***

func getDirFd(dir *C.DIR) uintptr ***REMOVED***
	return uintptr(C.dirfd(dir))
***REMOVED***

func subvolCreate(path, name string) error ***REMOVED***
	dir, err := openDir(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer closeDir(dir)

	var args C.struct_btrfs_ioctl_vol_args
	for i, c := range []byte(name) ***REMOVED***
		args.name[i] = C.char(c)
	***REMOVED***

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, getDirFd(dir), C.BTRFS_IOC_SUBVOL_CREATE,
		uintptr(unsafe.Pointer(&args)))
	if errno != 0 ***REMOVED***
		return fmt.Errorf("Failed to create btrfs subvolume: %v", errno.Error())
	***REMOVED***
	return nil
***REMOVED***

func subvolSnapshot(src, dest, name string) error ***REMOVED***
	srcDir, err := openDir(src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer closeDir(srcDir)

	destDir, err := openDir(dest)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer closeDir(destDir)

	var args C.struct_btrfs_ioctl_vol_args_v2
	args.fd = C.__s64(getDirFd(srcDir))

	var cs = C.CString(name)
	C.set_name_btrfs_ioctl_vol_args_v2(&args, cs)
	C.free(unsafe.Pointer(cs))

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, getDirFd(destDir), C.BTRFS_IOC_SNAP_CREATE_V2,
		uintptr(unsafe.Pointer(&args)))
	if errno != 0 ***REMOVED***
		return fmt.Errorf("Failed to create btrfs snapshot: %v", errno.Error())
	***REMOVED***
	return nil
***REMOVED***

func isSubvolume(p string) (bool, error) ***REMOVED***
	var bufStat unix.Stat_t
	if err := unix.Lstat(p, &bufStat); err != nil ***REMOVED***
		return false, err
	***REMOVED***

	// return true if it is a btrfs subvolume
	return bufStat.Ino == C.BTRFS_FIRST_FREE_OBJECTID, nil
***REMOVED***

func subvolDelete(dirpath, name string, quotaEnabled bool) error ***REMOVED***
	dir, err := openDir(dirpath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer closeDir(dir)
	fullPath := path.Join(dirpath, name)

	var args C.struct_btrfs_ioctl_vol_args

	// walk the btrfs subvolumes
	walkSubvolumes := func(p string, f os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			if os.IsNotExist(err) && p != fullPath ***REMOVED***
				// missing most likely because the path was a subvolume that got removed in the previous iteration
				// since it's gone anyway, we don't care
				return nil
			***REMOVED***
			return fmt.Errorf("error walking subvolumes: %v", err)
		***REMOVED***
		// we want to check children only so skip itself
		// it will be removed after the filepath walk anyways
		if f.IsDir() && p != fullPath ***REMOVED***
			sv, err := isSubvolume(p)
			if err != nil ***REMOVED***
				return fmt.Errorf("Failed to test if %s is a btrfs subvolume: %v", p, err)
			***REMOVED***
			if sv ***REMOVED***
				if err := subvolDelete(path.Dir(p), f.Name(), quotaEnabled); err != nil ***REMOVED***
					return fmt.Errorf("Failed to destroy btrfs child subvolume (%s) of parent (%s): %v", p, dirpath, err)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	if err := filepath.Walk(path.Join(dirpath, name), walkSubvolumes); err != nil ***REMOVED***
		return fmt.Errorf("Recursively walking subvolumes for %s failed: %v", dirpath, err)
	***REMOVED***

	if quotaEnabled ***REMOVED***
		if qgroupid, err := subvolLookupQgroup(fullPath); err == nil ***REMOVED***
			var args C.struct_btrfs_ioctl_qgroup_create_args
			args.qgroupid = C.__u64(qgroupid)

			_, _, errno := unix.Syscall(unix.SYS_IOCTL, getDirFd(dir), C.BTRFS_IOC_QGROUP_CREATE,
				uintptr(unsafe.Pointer(&args)))
			if errno != 0 ***REMOVED***
				logrus.Errorf("Failed to delete btrfs qgroup %v for %s: %v", qgroupid, fullPath, errno.Error())
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			logrus.Errorf("Failed to lookup btrfs qgroup for %s: %v", fullPath, err.Error())
		***REMOVED***
	***REMOVED***

	// all subvolumes have been removed
	// now remove the one originally passed in
	for i, c := range []byte(name) ***REMOVED***
		args.name[i] = C.char(c)
	***REMOVED***
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, getDirFd(dir), C.BTRFS_IOC_SNAP_DESTROY,
		uintptr(unsafe.Pointer(&args)))
	if errno != 0 ***REMOVED***
		return fmt.Errorf("Failed to destroy btrfs snapshot %s for %s: %v", dirpath, name, errno.Error())
	***REMOVED***
	return nil
***REMOVED***

func (d *Driver) updateQuotaStatus() ***REMOVED***
	d.once.Do(func() ***REMOVED***
		if !d.quotaEnabled ***REMOVED***
			// In case quotaEnabled is not set, check qgroup and update quotaEnabled as needed
			if err := subvolQgroupStatus(d.home); err != nil ***REMOVED***
				// quota is still not enabled
				return
			***REMOVED***
			d.quotaEnabled = true
		***REMOVED***
	***REMOVED***)
***REMOVED***

func (d *Driver) subvolEnableQuota() error ***REMOVED***
	d.updateQuotaStatus()

	if d.quotaEnabled ***REMOVED***
		return nil
	***REMOVED***

	dir, err := openDir(d.home)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer closeDir(dir)

	var args C.struct_btrfs_ioctl_quota_ctl_args
	args.cmd = C.BTRFS_QUOTA_CTL_ENABLE
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, getDirFd(dir), C.BTRFS_IOC_QUOTA_CTL,
		uintptr(unsafe.Pointer(&args)))
	if errno != 0 ***REMOVED***
		return fmt.Errorf("Failed to enable btrfs quota for %s: %v", dir, errno.Error())
	***REMOVED***

	d.quotaEnabled = true

	return nil
***REMOVED***

func (d *Driver) subvolDisableQuota() error ***REMOVED***
	d.updateQuotaStatus()

	if !d.quotaEnabled ***REMOVED***
		return nil
	***REMOVED***

	dir, err := openDir(d.home)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer closeDir(dir)

	var args C.struct_btrfs_ioctl_quota_ctl_args
	args.cmd = C.BTRFS_QUOTA_CTL_DISABLE
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, getDirFd(dir), C.BTRFS_IOC_QUOTA_CTL,
		uintptr(unsafe.Pointer(&args)))
	if errno != 0 ***REMOVED***
		return fmt.Errorf("Failed to disable btrfs quota for %s: %v", dir, errno.Error())
	***REMOVED***

	d.quotaEnabled = false

	return nil
***REMOVED***

func (d *Driver) subvolRescanQuota() error ***REMOVED***
	d.updateQuotaStatus()

	if !d.quotaEnabled ***REMOVED***
		return nil
	***REMOVED***

	dir, err := openDir(d.home)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer closeDir(dir)

	var args C.struct_btrfs_ioctl_quota_rescan_args
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, getDirFd(dir), C.BTRFS_IOC_QUOTA_RESCAN_WAIT,
		uintptr(unsafe.Pointer(&args)))
	if errno != 0 ***REMOVED***
		return fmt.Errorf("Failed to rescan btrfs quota for %s: %v", dir, errno.Error())
	***REMOVED***

	return nil
***REMOVED***

func subvolLimitQgroup(path string, size uint64) error ***REMOVED***
	dir, err := openDir(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer closeDir(dir)

	var args C.struct_btrfs_ioctl_qgroup_limit_args
	args.lim.max_referenced = C.__u64(size)
	args.lim.flags = C.BTRFS_QGROUP_LIMIT_MAX_RFER
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, getDirFd(dir), C.BTRFS_IOC_QGROUP_LIMIT,
		uintptr(unsafe.Pointer(&args)))
	if errno != 0 ***REMOVED***
		return fmt.Errorf("Failed to limit qgroup for %s: %v", dir, errno.Error())
	***REMOVED***

	return nil
***REMOVED***

// subvolQgroupStatus performs a BTRFS_IOC_TREE_SEARCH on the root path
// with search key of BTRFS_QGROUP_STATUS_KEY.
// In case qgroup is enabled, the retuned key type will match BTRFS_QGROUP_STATUS_KEY.
// For more details please see https://github.com/kdave/btrfs-progs/blob/v4.9/qgroup.c#L1035
func subvolQgroupStatus(path string) error ***REMOVED***
	dir, err := openDir(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer closeDir(dir)

	var args C.struct_btrfs_ioctl_search_args
	args.key.tree_id = C.BTRFS_QUOTA_TREE_OBJECTID
	args.key.min_type = C.BTRFS_QGROUP_STATUS_KEY
	args.key.max_type = C.BTRFS_QGROUP_STATUS_KEY
	args.key.max_objectid = C.__u64(math.MaxUint64)
	args.key.max_offset = C.__u64(math.MaxUint64)
	args.key.max_transid = C.__u64(math.MaxUint64)
	args.key.nr_items = 4096

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, getDirFd(dir), C.BTRFS_IOC_TREE_SEARCH,
		uintptr(unsafe.Pointer(&args)))
	if errno != 0 ***REMOVED***
		return fmt.Errorf("Failed to search qgroup for %s: %v", path, errno.Error())
	***REMOVED***
	sh := (*C.struct_btrfs_ioctl_search_header)(unsafe.Pointer(&args.buf))
	if sh._type != C.BTRFS_QGROUP_STATUS_KEY ***REMOVED***
		return fmt.Errorf("Invalid qgroup search header type for %s: %v", path, sh._type)
	***REMOVED***
	return nil
***REMOVED***

func subvolLookupQgroup(path string) (uint64, error) ***REMOVED***
	dir, err := openDir(path)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer closeDir(dir)

	var args C.struct_btrfs_ioctl_ino_lookup_args
	args.objectid = C.BTRFS_FIRST_FREE_OBJECTID

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, getDirFd(dir), C.BTRFS_IOC_INO_LOOKUP,
		uintptr(unsafe.Pointer(&args)))
	if errno != 0 ***REMOVED***
		return 0, fmt.Errorf("Failed to lookup qgroup for %s: %v", dir, errno.Error())
	***REMOVED***
	if args.treeid == 0 ***REMOVED***
		return 0, fmt.Errorf("Invalid qgroup id for %s: 0", dir)
	***REMOVED***

	return uint64(args.treeid), nil
***REMOVED***

func (d *Driver) subvolumesDir() string ***REMOVED***
	return path.Join(d.home, "subvolumes")
***REMOVED***

func (d *Driver) subvolumesDirID(id string) string ***REMOVED***
	return path.Join(d.subvolumesDir(), id)
***REMOVED***

func (d *Driver) quotasDir() string ***REMOVED***
	return path.Join(d.home, "quotas")
***REMOVED***

func (d *Driver) quotasDirID(id string) string ***REMOVED***
	return path.Join(d.quotasDir(), id)
***REMOVED***

// CreateReadWrite creates a layer that is writable for use as a container
// file system.
func (d *Driver) CreateReadWrite(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	return d.Create(id, parent, opts)
***REMOVED***

// Create the filesystem with given id.
func (d *Driver) Create(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	quotas := path.Join(d.home, "quotas")
	subvolumes := path.Join(d.home, "subvolumes")
	rootUID, rootGID, err := idtools.GetRootUIDGID(d.uidMaps, d.gidMaps)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := idtools.MkdirAllAndChown(subvolumes, 0700, idtools.IDPair***REMOVED***UID: rootUID, GID: rootGID***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	if parent == "" ***REMOVED***
		if err := subvolCreate(subvolumes, id); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		parentDir := d.subvolumesDirID(parent)
		st, err := os.Stat(parentDir)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if !st.IsDir() ***REMOVED***
			return fmt.Errorf("%s: not a directory", parentDir)
		***REMOVED***
		if err := subvolSnapshot(parentDir, subvolumes, id); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	var storageOpt map[string]string
	if opts != nil ***REMOVED***
		storageOpt = opts.StorageOpt
	***REMOVED***

	if _, ok := storageOpt["size"]; ok ***REMOVED***
		driver := &Driver***REMOVED******REMOVED***
		if err := d.parseStorageOpt(storageOpt, driver); err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := d.setStorageSize(path.Join(subvolumes, id), driver); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := idtools.MkdirAllAndChown(quotas, 0700, idtools.IDPair***REMOVED***UID: rootUID, GID: rootGID***REMOVED***); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := ioutil.WriteFile(path.Join(quotas, id), []byte(fmt.Sprint(driver.options.size)), 0644); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// if we have a remapped root (user namespaces enabled), change the created snapshot
	// dir ownership to match
	if rootUID != 0 || rootGID != 0 ***REMOVED***
		if err := os.Chown(path.Join(subvolumes, id), rootUID, rootGID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	mountLabel := ""
	if opts != nil ***REMOVED***
		mountLabel = opts.MountLabel
	***REMOVED***

	return label.Relabel(path.Join(subvolumes, id), mountLabel, false)
***REMOVED***

// Parse btrfs storage options
func (d *Driver) parseStorageOpt(storageOpt map[string]string, driver *Driver) error ***REMOVED***
	// Read size to change the subvolume disk quota per container
	for key, val := range storageOpt ***REMOVED***
		key := strings.ToLower(key)
		switch key ***REMOVED***
		case "size":
			size, err := units.RAMInBytes(val)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			driver.options.size = uint64(size)
		default:
			return fmt.Errorf("Unknown option %s", key)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Set btrfs storage size
func (d *Driver) setStorageSize(dir string, driver *Driver) error ***REMOVED***
	if driver.options.size <= 0 ***REMOVED***
		return fmt.Errorf("btrfs: invalid storage size: %s", units.HumanSize(float64(driver.options.size)))
	***REMOVED***
	if d.options.minSpace > 0 && driver.options.size < d.options.minSpace ***REMOVED***
		return fmt.Errorf("btrfs: storage size cannot be less than %s", units.HumanSize(float64(d.options.minSpace)))
	***REMOVED***
	if err := d.subvolEnableQuota(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return subvolLimitQgroup(dir, driver.options.size)
***REMOVED***

// Remove the filesystem with given id.
func (d *Driver) Remove(id string) error ***REMOVED***
	dir := d.subvolumesDirID(id)
	if _, err := os.Stat(dir); err != nil ***REMOVED***
		return err
	***REMOVED***
	quotasDir := d.quotasDirID(id)
	if _, err := os.Stat(quotasDir); err == nil ***REMOVED***
		if err := os.Remove(quotasDir); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else if !os.IsNotExist(err) ***REMOVED***
		return err
	***REMOVED***

	// Call updateQuotaStatus() to invoke status update
	d.updateQuotaStatus()

	if err := subvolDelete(d.subvolumesDir(), id, d.quotaEnabled); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := system.EnsureRemoveAll(dir); err != nil ***REMOVED***
		return err
	***REMOVED***
	return d.subvolRescanQuota()
***REMOVED***

// Get the requested filesystem id.
func (d *Driver) Get(id, mountLabel string) (containerfs.ContainerFS, error) ***REMOVED***
	dir := d.subvolumesDirID(id)
	st, err := os.Stat(dir)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if !st.IsDir() ***REMOVED***
		return nil, fmt.Errorf("%s: not a directory", dir)
	***REMOVED***

	if quota, err := ioutil.ReadFile(d.quotasDirID(id)); err == nil ***REMOVED***
		if size, err := strconv.ParseUint(string(quota), 10, 64); err == nil && size >= d.options.minSpace ***REMOVED***
			if err := d.subvolEnableQuota(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if err := subvolLimitQgroup(dir, size); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return containerfs.NewLocalContainerFS(dir), nil
***REMOVED***

// Put is not implemented for BTRFS as there is no cleanup required for the id.
func (d *Driver) Put(id string) error ***REMOVED***
	// Get() creates no runtime resources (like e.g. mounts)
	// so this doesn't need to do anything.
	return nil
***REMOVED***

// Exists checks if the id exists in the filesystem.
func (d *Driver) Exists(id string) bool ***REMOVED***
	dir := d.subvolumesDirID(id)
	_, err := os.Stat(dir)
	return err == nil
***REMOVED***
