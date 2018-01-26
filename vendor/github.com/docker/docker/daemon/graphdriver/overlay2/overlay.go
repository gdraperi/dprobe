// +build linux

package overlay2

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/daemon/graphdriver/overlayutils"
	"github.com/docker/docker/daemon/graphdriver/quota"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/chrootarchive"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/directory"
	"github.com/docker/docker/pkg/fsutils"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/locker"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/parsers"
	"github.com/docker/docker/pkg/parsers/kernel"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/go-units"
	rsystem "github.com/opencontainers/runc/libcontainer/system"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

var (
	// untar defines the untar method
	untar = chrootarchive.UntarUncompressed
)

// This backend uses the overlay union filesystem for containers
// with diff directories for each layer.

// This version of the overlay driver requires at least kernel
// 4.0.0 in order to support mounting multiple diff directories.

// Each container/image has at least a "diff" directory and "link" file.
// If there is also a "lower" file when there are diff layers
// below as well as "merged" and "work" directories. The "diff" directory
// has the upper layer of the overlay and is used to capture any
// changes to the layer. The "lower" file contains all the lower layer
// mounts separated by ":" and ordered from uppermost to lowermost
// layers. The overlay itself is mounted in the "merged" directory,
// and the "work" dir is needed for overlay to work.

// The "link" file for each layer contains a unique string for the layer.
// Under the "l" directory at the root there will be a symbolic link
// with that unique string pointing the "diff" directory for the layer.
// The symbolic links are used to reference lower layers in the "lower"
// file and on mount. The links are used to shorten the total length
// of a layer reference without requiring changes to the layer identifier
// or root directory. Mounts are always done relative to root and
// referencing the symbolic links in order to ensure the number of
// lower directories can fit in a single page for making the mount
// syscall. A hard upper limit of 128 lower layers is enforced to ensure
// that mounts do not fail due to length.

const (
	driverName = "overlay2"
	linkDir    = "l"
	lowerFile  = "lower"
	maxDepth   = 128

	// idLength represents the number of random characters
	// which can be used to create the unique link identifier
	// for every layer. If this value is too long then the
	// page size limit for the mount command may be exceeded.
	// The idLength should be selected such that following equation
	// is true (512 is a buffer for label metadata).
	// ((idLength + len(linkDir) + 1) * maxDepth) <= (pageSize - 512)
	idLength = 26
)

type overlayOptions struct ***REMOVED***
	overrideKernelCheck bool
	quota               quota.Quota
***REMOVED***

// Driver contains information about the home directory and the list of active
// mounts that are created using this driver.
type Driver struct ***REMOVED***
	home          string
	uidMaps       []idtools.IDMap
	gidMaps       []idtools.IDMap
	ctr           *graphdriver.RefCounter
	quotaCtl      *quota.Control
	options       overlayOptions
	naiveDiff     graphdriver.DiffDriver
	supportsDType bool
	locker        *locker.Locker
***REMOVED***

var (
	backingFs             = "<unknown>"
	projectQuotaSupported = false

	useNaiveDiffLock sync.Once
	useNaiveDiffOnly bool
)

func init() ***REMOVED***
	graphdriver.Register(driverName, Init)
***REMOVED***

// Init returns the native diff driver for overlay filesystem.
// If overlay filesystem is not supported on the host, the error
// graphdriver.ErrNotSupported is returned.
// If an overlay filesystem is not supported over an existing filesystem then
// the error graphdriver.ErrIncompatibleFS is returned.
func Init(home string, options []string, uidMaps, gidMaps []idtools.IDMap) (graphdriver.Driver, error) ***REMOVED***
	opts, err := parseOptions(options)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := supportsOverlay(); err != nil ***REMOVED***
		return nil, graphdriver.ErrNotSupported
	***REMOVED***

	// require kernel 4.0.0 to ensure multiple lower dirs are supported
	v, err := kernel.GetKernelVersion()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Perform feature detection on /var/lib/docker/overlay2 if it's an existing directory.
	// This covers situations where /var/lib/docker/overlay2 is a mount, and on a different
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
	if fsName, ok := graphdriver.FsNames[fsMagic]; ok ***REMOVED***
		backingFs = fsName
	***REMOVED***

	switch fsMagic ***REMOVED***
	case graphdriver.FsMagicAufs, graphdriver.FsMagicEcryptfs, graphdriver.FsMagicNfsFs, graphdriver.FsMagicOverlay, graphdriver.FsMagicZfs:
		logrus.Errorf("'overlay2' is not supported over %s", backingFs)
		return nil, graphdriver.ErrIncompatibleFS
	case graphdriver.FsMagicBtrfs:
		// Support for OverlayFS on BTRFS was added in kernel 4.7
		// See https://btrfs.wiki.kernel.org/index.php/Changelog
		if kernel.CompareKernelVersion(*v, kernel.VersionInfo***REMOVED***Kernel: 4, Major: 7, Minor: 0***REMOVED***) < 0 ***REMOVED***
			if !opts.overrideKernelCheck ***REMOVED***
				logrus.Errorf("'overlay2' requires kernel 4.7 to use on %s", backingFs)
				return nil, graphdriver.ErrIncompatibleFS
			***REMOVED***
			logrus.Warn("Using pre-4.7.0 kernel for overlay2 on btrfs, may require kernel update")
		***REMOVED***
	***REMOVED***

	if kernel.CompareKernelVersion(*v, kernel.VersionInfo***REMOVED***Kernel: 4, Major: 0, Minor: 0***REMOVED***) < 0 ***REMOVED***
		if opts.overrideKernelCheck ***REMOVED***
			logrus.Warn("Using pre-4.0.0 kernel for overlay2, mount failures may require kernel update")
		***REMOVED*** else ***REMOVED***
			if err := supportsMultipleLowerDir(testdir); err != nil ***REMOVED***
				logrus.Debugf("Multiple lower dirs not supported: %v", err)
				return nil, graphdriver.ErrNotSupported
			***REMOVED***
		***REMOVED***
	***REMOVED***
	supportsDType, err := fsutils.SupportsDType(testdir)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !supportsDType ***REMOVED***
		if !graphdriver.IsInitialized(home) ***REMOVED***
			return nil, overlayutils.ErrDTypeNotSupported("overlay2", backingFs)
		***REMOVED***
		// allow running without d_type only for existing setups (#27443)
		logrus.Warn(overlayutils.ErrDTypeNotSupported("overlay2", backingFs))
	***REMOVED***

	rootUID, rootGID, err := idtools.GetRootUIDGID(uidMaps, gidMaps)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Create the driver home dir
	if err := idtools.MkdirAllAndChown(path.Join(home, linkDir), 0700, idtools.IDPair***REMOVED***rootUID, rootGID***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	d := &Driver***REMOVED***
		home:          home,
		uidMaps:       uidMaps,
		gidMaps:       gidMaps,
		ctr:           graphdriver.NewRefCounter(graphdriver.NewFsChecker(graphdriver.FsMagicOverlay)),
		supportsDType: supportsDType,
		locker:        locker.New(),
		options:       *opts,
	***REMOVED***

	d.naiveDiff = graphdriver.NewNaiveDiffDriver(d, uidMaps, gidMaps)

	if backingFs == "xfs" ***REMOVED***
		// Try to enable project quota support over xfs.
		if d.quotaCtl, err = quota.NewControl(home); err == nil ***REMOVED***
			projectQuotaSupported = true
		***REMOVED*** else if opts.quota.Size > 0 ***REMOVED***
			return nil, fmt.Errorf("Storage option overlay2.size not supported. Filesystem does not support Project Quota: %v", err)
		***REMOVED***
	***REMOVED*** else if opts.quota.Size > 0 ***REMOVED***
		// if xfs is not the backing fs then error out if the storage-opt overlay2.size is used.
		return nil, fmt.Errorf("Storage Option overlay2.size only supported for backingFS XFS. Found %v", backingFs)
	***REMOVED***

	logrus.Debugf("backingFs=%s,  projectQuotaSupported=%v", backingFs, projectQuotaSupported)

	return d, nil
***REMOVED***

func parseOptions(options []string) (*overlayOptions, error) ***REMOVED***
	o := &overlayOptions***REMOVED******REMOVED***
	for _, option := range options ***REMOVED***
		key, val, err := parsers.ParseKeyValueOpt(option)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		key = strings.ToLower(key)
		switch key ***REMOVED***
		case "overlay2.override_kernel_check":
			o.overrideKernelCheck, err = strconv.ParseBool(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		case "overlay2.size":
			size, err := units.RAMInBytes(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			o.quota.Size = uint64(size)
		default:
			return nil, fmt.Errorf("overlay2: unknown option %s", key)
		***REMOVED***
	***REMOVED***
	return o, nil
***REMOVED***

func supportsOverlay() error ***REMOVED***
	// We can try to modprobe overlay first before looking at
	// proc/filesystems for when overlay is supported
	exec.Command("modprobe", "overlay").Run()

	f, err := os.Open("/proc/filesystems")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() ***REMOVED***
		if s.Text() == "nodev\toverlay" ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	logrus.Error("'overlay' not found as a supported filesystem on this host. Please ensure kernel is new enough and has overlay support loaded.")
	return graphdriver.ErrNotSupported
***REMOVED***

func useNaiveDiff(home string) bool ***REMOVED***
	useNaiveDiffLock.Do(func() ***REMOVED***
		if err := doesSupportNativeDiff(home); err != nil ***REMOVED***
			logrus.Warnf("Not using native diff for overlay2, this may cause degraded performance for building images: %v", err)
			useNaiveDiffOnly = true
		***REMOVED***
	***REMOVED***)
	return useNaiveDiffOnly
***REMOVED***

func (d *Driver) String() string ***REMOVED***
	return driverName
***REMOVED***

// Status returns current driver information in a two dimensional string array.
// Output contains "Backing Filesystem" used in this implementation.
func (d *Driver) Status() [][2]string ***REMOVED***
	return [][2]string***REMOVED***
		***REMOVED***"Backing Filesystem", backingFs***REMOVED***,
		***REMOVED***"Supports d_type", strconv.FormatBool(d.supportsDType)***REMOVED***,
		***REMOVED***"Native Overlay Diff", strconv.FormatBool(!useNaiveDiff(d.home))***REMOVED***,
	***REMOVED***
***REMOVED***

// GetMetadata returns metadata about the overlay driver such as the LowerDir,
// UpperDir, WorkDir, and MergeDir used to store data.
func (d *Driver) GetMetadata(id string) (map[string]string, error) ***REMOVED***
	dir := d.dir(id)
	if _, err := os.Stat(dir); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	metadata := map[string]string***REMOVED***
		"WorkDir":   path.Join(dir, "work"),
		"MergedDir": path.Join(dir, "merged"),
		"UpperDir":  path.Join(dir, "diff"),
	***REMOVED***

	lowerDirs, err := d.getLowerDirs(id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(lowerDirs) > 0 ***REMOVED***
		metadata["LowerDir"] = strings.Join(lowerDirs, ":")
	***REMOVED***

	return metadata, nil
***REMOVED***

// Cleanup any state created by overlay which should be cleaned when daemon
// is being shutdown. For now, we just have to unmount the bind mounted
// we had created.
func (d *Driver) Cleanup() error ***REMOVED***
	return mount.RecursiveUnmount(d.home)
***REMOVED***

// CreateReadWrite creates a layer that is writable for use as a container
// file system.
func (d *Driver) CreateReadWrite(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	if opts != nil && len(opts.StorageOpt) != 0 && !projectQuotaSupported ***REMOVED***
		return fmt.Errorf("--storage-opt is supported only for overlay over xfs with 'pquota' mount option")
	***REMOVED***

	if opts == nil ***REMOVED***
		opts = &graphdriver.CreateOpts***REMOVED***
			StorageOpt: map[string]string***REMOVED******REMOVED***,
		***REMOVED***
	***REMOVED***

	if _, ok := opts.StorageOpt["size"]; !ok ***REMOVED***
		if opts.StorageOpt == nil ***REMOVED***
			opts.StorageOpt = map[string]string***REMOVED******REMOVED***
		***REMOVED***
		opts.StorageOpt["size"] = strconv.FormatUint(d.options.quota.Size, 10)
	***REMOVED***

	return d.create(id, parent, opts)
***REMOVED***

// Create is used to create the upper, lower, and merge directories required for overlay fs for a given id.
// The parent filesystem is used to configure these directories for the overlay.
func (d *Driver) Create(id, parent string, opts *graphdriver.CreateOpts) (retErr error) ***REMOVED***
	if opts != nil && len(opts.StorageOpt) != 0 ***REMOVED***
		if _, ok := opts.StorageOpt["size"]; ok ***REMOVED***
			return fmt.Errorf("--storage-opt size is only supported for ReadWrite Layers")
		***REMOVED***
	***REMOVED***
	return d.create(id, parent, opts)
***REMOVED***

func (d *Driver) create(id, parent string, opts *graphdriver.CreateOpts) (retErr error) ***REMOVED***
	dir := d.dir(id)

	rootUID, rootGID, err := idtools.GetRootUIDGID(d.uidMaps, d.gidMaps)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	root := idtools.IDPair***REMOVED***UID: rootUID, GID: rootGID***REMOVED***

	if err := idtools.MkdirAllAndChown(path.Dir(dir), 0700, root); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := idtools.MkdirAndChown(dir, 0700, root); err != nil ***REMOVED***
		return err
	***REMOVED***

	defer func() ***REMOVED***
		// Clean up on failure
		if retErr != nil ***REMOVED***
			os.RemoveAll(dir)
		***REMOVED***
	***REMOVED***()

	if opts != nil && len(opts.StorageOpt) > 0 ***REMOVED***
		driver := &Driver***REMOVED******REMOVED***
		if err := d.parseStorageOpt(opts.StorageOpt, driver); err != nil ***REMOVED***
			return err
		***REMOVED***

		if driver.options.quota.Size > 0 ***REMOVED***
			// Set container disk quota limit
			if err := d.quotaCtl.SetQuota(dir, driver.options.quota); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if err := idtools.MkdirAndChown(path.Join(dir, "diff"), 0755, root); err != nil ***REMOVED***
		return err
	***REMOVED***

	lid := generateID(idLength)
	if err := os.Symlink(path.Join("..", id, "diff"), path.Join(d.home, linkDir, lid)); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Write link id to link file
	if err := ioutil.WriteFile(path.Join(dir, "link"), []byte(lid), 0644); err != nil ***REMOVED***
		return err
	***REMOVED***

	// if no parent directory, done
	if parent == "" ***REMOVED***
		return nil
	***REMOVED***

	if err := idtools.MkdirAndChown(path.Join(dir, "work"), 0700, root); err != nil ***REMOVED***
		return err
	***REMOVED***

	lower, err := d.getLower(parent)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if lower != "" ***REMOVED***
		if err := ioutil.WriteFile(path.Join(dir, lowerFile), []byte(lower), 0666); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Parse overlay storage options
func (d *Driver) parseStorageOpt(storageOpt map[string]string, driver *Driver) error ***REMOVED***
	// Read size to set the disk project quota per container
	for key, val := range storageOpt ***REMOVED***
		key := strings.ToLower(key)
		switch key ***REMOVED***
		case "size":
			size, err := units.RAMInBytes(val)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			driver.options.quota.Size = uint64(size)
		default:
			return fmt.Errorf("Unknown option %s", key)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *Driver) getLower(parent string) (string, error) ***REMOVED***
	parentDir := d.dir(parent)

	// Ensure parent exists
	if _, err := os.Lstat(parentDir); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	// Read Parent link fileA
	parentLink, err := ioutil.ReadFile(path.Join(parentDir, "link"))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	lowers := []string***REMOVED***path.Join(linkDir, string(parentLink))***REMOVED***

	parentLower, err := ioutil.ReadFile(path.Join(parentDir, lowerFile))
	if err == nil ***REMOVED***
		parentLowers := strings.Split(string(parentLower), ":")
		lowers = append(lowers, parentLowers...)
	***REMOVED***
	if len(lowers) > maxDepth ***REMOVED***
		return "", errors.New("max depth exceeded")
	***REMOVED***
	return strings.Join(lowers, ":"), nil
***REMOVED***

func (d *Driver) dir(id string) string ***REMOVED***
	return path.Join(d.home, id)
***REMOVED***

func (d *Driver) getLowerDirs(id string) ([]string, error) ***REMOVED***
	var lowersArray []string
	lowers, err := ioutil.ReadFile(path.Join(d.dir(id), lowerFile))
	if err == nil ***REMOVED***
		for _, s := range strings.Split(string(lowers), ":") ***REMOVED***
			lp, err := os.Readlink(path.Join(d.home, s))
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			lowersArray = append(lowersArray, path.Clean(path.Join(d.home, linkDir, lp)))
		***REMOVED***
	***REMOVED*** else if !os.IsNotExist(err) ***REMOVED***
		return nil, err
	***REMOVED***
	return lowersArray, nil
***REMOVED***

// Remove cleans the directories that are created for this id.
func (d *Driver) Remove(id string) error ***REMOVED***
	d.locker.Lock(id)
	defer d.locker.Unlock(id)
	dir := d.dir(id)
	lid, err := ioutil.ReadFile(path.Join(dir, "link"))
	if err == nil ***REMOVED***
		if err := os.RemoveAll(path.Join(d.home, linkDir, string(lid))); err != nil ***REMOVED***
			logrus.Debugf("Failed to remove link: %v", err)
		***REMOVED***
	***REMOVED***

	if err := system.EnsureRemoveAll(dir); err != nil && !os.IsNotExist(err) ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Get creates and mounts the required file system for the given id and returns the mount path.
func (d *Driver) Get(id, mountLabel string) (_ containerfs.ContainerFS, retErr error) ***REMOVED***
	d.locker.Lock(id)
	defer d.locker.Unlock(id)
	dir := d.dir(id)
	if _, err := os.Stat(dir); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	diffDir := path.Join(dir, "diff")
	lowers, err := ioutil.ReadFile(path.Join(dir, lowerFile))
	if err != nil ***REMOVED***
		// If no lower, just return diff directory
		if os.IsNotExist(err) ***REMOVED***
			return containerfs.NewLocalContainerFS(diffDir), nil
		***REMOVED***
		return nil, err
	***REMOVED***

	mergedDir := path.Join(dir, "merged")
	if count := d.ctr.Increment(mergedDir); count > 1 ***REMOVED***
		return containerfs.NewLocalContainerFS(mergedDir), nil
	***REMOVED***
	defer func() ***REMOVED***
		if retErr != nil ***REMOVED***
			if c := d.ctr.Decrement(mergedDir); c <= 0 ***REMOVED***
				if mntErr := unix.Unmount(mergedDir, 0); mntErr != nil ***REMOVED***
					logrus.Errorf("error unmounting %v: %v", mergedDir, mntErr)
				***REMOVED***
				// Cleanup the created merged directory; see the comment in Put's rmdir
				if rmErr := unix.Rmdir(mergedDir); rmErr != nil && !os.IsNotExist(rmErr) ***REMOVED***
					logrus.Debugf("Failed to remove %s: %v: %v", id, rmErr, err)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	workDir := path.Join(dir, "work")
	splitLowers := strings.Split(string(lowers), ":")
	absLowers := make([]string, len(splitLowers))
	for i, s := range splitLowers ***REMOVED***
		absLowers[i] = path.Join(d.home, s)
	***REMOVED***
	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", strings.Join(absLowers, ":"), path.Join(dir, "diff"), path.Join(dir, "work"))
	mountData := label.FormatMountLabel(opts, mountLabel)
	mount := unix.Mount
	mountTarget := mergedDir

	rootUID, rootGID, err := idtools.GetRootUIDGID(d.uidMaps, d.gidMaps)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := idtools.MkdirAndChown(mergedDir, 0700, idtools.IDPair***REMOVED***rootUID, rootGID***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pageSize := unix.Getpagesize()

	// Go can return a larger page size than supported by the system
	// as of go 1.7. This will be fixed in 1.8 and this block can be
	// removed when building with 1.8.
	// See https://github.com/golang/go/commit/1b9499b06989d2831e5b156161d6c07642926ee1
	// See https://github.com/docker/docker/issues/27384
	if pageSize > 4096 ***REMOVED***
		pageSize = 4096
	***REMOVED***

	// Use relative paths and mountFrom when the mount data has exceeded
	// the page size. The mount syscall fails if the mount data cannot
	// fit within a page and relative links make the mount data much
	// smaller at the expense of requiring a fork exec to chroot.
	if len(mountData) > pageSize ***REMOVED***
		opts = fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", string(lowers), path.Join(id, "diff"), path.Join(id, "work"))
		mountData = label.FormatMountLabel(opts, mountLabel)
		if len(mountData) > pageSize ***REMOVED***
			return nil, fmt.Errorf("cannot mount layer, mount label too large %d", len(mountData))
		***REMOVED***

		mount = func(source string, target string, mType string, flags uintptr, label string) error ***REMOVED***
			return mountFrom(d.home, source, target, mType, flags, label)
		***REMOVED***
		mountTarget = path.Join(id, "merged")
	***REMOVED***

	if err := mount("overlay", mountTarget, "overlay", 0, mountData); err != nil ***REMOVED***
		return nil, fmt.Errorf("error creating overlay mount to %s: %v", mergedDir, err)
	***REMOVED***

	// chown "workdir/work" to the remapped root UID/GID. Overlay fs inside a
	// user namespace requires this to move a directory from lower to upper.
	if err := os.Chown(path.Join(workDir, "work"), rootUID, rootGID); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return containerfs.NewLocalContainerFS(mergedDir), nil
***REMOVED***

// Put unmounts the mount path created for the give id.
// It also removes the 'merged' directory to force the kernel to unmount the
// overlay mount in other namespaces.
func (d *Driver) Put(id string) error ***REMOVED***
	d.locker.Lock(id)
	defer d.locker.Unlock(id)
	dir := d.dir(id)
	_, err := ioutil.ReadFile(path.Join(dir, lowerFile))
	if err != nil ***REMOVED***
		// If no lower, no mount happened and just return directly
		if os.IsNotExist(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	mountpoint := path.Join(dir, "merged")
	if count := d.ctr.Decrement(mountpoint); count > 0 ***REMOVED***
		return nil
	***REMOVED***
	if err := unix.Unmount(mountpoint, unix.MNT_DETACH); err != nil ***REMOVED***
		logrus.Debugf("Failed to unmount %s overlay: %s - %v", id, mountpoint, err)
	***REMOVED***
	// Remove the mountpoint here. Removing the mountpoint (in newer kernels)
	// will cause all other instances of this mount in other mount namespaces
	// to be unmounted. This is necessary to avoid cases where an overlay mount
	// that is present in another namespace will cause subsequent mounts
	// operations to fail with ebusy.  We ignore any errors here because this may
	// fail on older kernels which don't have
	// torvalds/linux@8ed936b5671bfb33d89bc60bdcc7cf0470ba52fe applied.
	if err := unix.Rmdir(mountpoint); err != nil && !os.IsNotExist(err) ***REMOVED***
		logrus.Debugf("Failed to remove %s overlay: %v", id, err)
	***REMOVED***
	return nil
***REMOVED***

// Exists checks to see if the id is already mounted.
func (d *Driver) Exists(id string) bool ***REMOVED***
	_, err := os.Stat(d.dir(id))
	return err == nil
***REMOVED***

// isParent determines whether the given parent is the direct parent of the
// given layer id
func (d *Driver) isParent(id, parent string) bool ***REMOVED***
	lowers, err := d.getLowerDirs(id)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	if parent == "" && len(lowers) > 0 ***REMOVED***
		return false
	***REMOVED***

	parentDir := d.dir(parent)
	var ld string
	if len(lowers) > 0 ***REMOVED***
		ld = filepath.Dir(lowers[0])
	***REMOVED***
	if ld == "" && parent == "" ***REMOVED***
		return true
	***REMOVED***
	return ld == parentDir
***REMOVED***

// ApplyDiff applies the new layer into a root
func (d *Driver) ApplyDiff(id string, parent string, diff io.Reader) (size int64, err error) ***REMOVED***
	if !d.isParent(id, parent) ***REMOVED***
		return d.naiveDiff.ApplyDiff(id, parent, diff)
	***REMOVED***

	applyDir := d.getDiffPath(id)

	logrus.Debugf("Applying tar in %s", applyDir)
	// Overlay doesn't need the parent id to apply the diff
	if err := untar(diff, applyDir, &archive.TarOptions***REMOVED***
		UIDMaps:        d.uidMaps,
		GIDMaps:        d.gidMaps,
		WhiteoutFormat: archive.OverlayWhiteoutFormat,
		InUserNS:       rsystem.RunningInUserNS(),
	***REMOVED***); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return directory.Size(applyDir)
***REMOVED***

func (d *Driver) getDiffPath(id string) string ***REMOVED***
	dir := d.dir(id)

	return path.Join(dir, "diff")
***REMOVED***

// DiffSize calculates the changes between the specified id
// and its parent and returns the size in bytes of the changes
// relative to its base filesystem directory.
func (d *Driver) DiffSize(id, parent string) (size int64, err error) ***REMOVED***
	if useNaiveDiff(d.home) || !d.isParent(id, parent) ***REMOVED***
		return d.naiveDiff.DiffSize(id, parent)
	***REMOVED***
	return directory.Size(d.getDiffPath(id))
***REMOVED***

// Diff produces an archive of the changes between the specified
// layer and its parent layer which may be "".
func (d *Driver) Diff(id, parent string) (io.ReadCloser, error) ***REMOVED***
	if useNaiveDiff(d.home) || !d.isParent(id, parent) ***REMOVED***
		return d.naiveDiff.Diff(id, parent)
	***REMOVED***

	diffPath := d.getDiffPath(id)
	logrus.Debugf("Tar with options on %s", diffPath)
	return archive.TarWithOptions(diffPath, &archive.TarOptions***REMOVED***
		Compression:    archive.Uncompressed,
		UIDMaps:        d.uidMaps,
		GIDMaps:        d.gidMaps,
		WhiteoutFormat: archive.OverlayWhiteoutFormat,
	***REMOVED***)
***REMOVED***

// Changes produces a list of changes between the specified layer and its
// parent layer. If parent is "", then all changes will be ADD changes.
func (d *Driver) Changes(id, parent string) ([]archive.Change, error) ***REMOVED***
	if useNaiveDiff(d.home) || !d.isParent(id, parent) ***REMOVED***
		return d.naiveDiff.Changes(id, parent)
	***REMOVED***
	// Overlay doesn't have snapshots, so we need to get changes from all parent
	// layers.
	diffPath := d.getDiffPath(id)
	layers, err := d.getLowerDirs(id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return archive.OverlayChanges(layers, diffPath)
***REMOVED***
