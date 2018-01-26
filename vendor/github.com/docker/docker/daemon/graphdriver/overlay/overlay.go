// +build linux

package overlay

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/daemon/graphdriver/copy"
	"github.com/docker/docker/daemon/graphdriver/overlayutils"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/fsutils"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/locker"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/system"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// This is a small wrapper over the NaiveDiffWriter that lets us have a custom
// implementation of ApplyDiff()

var (
	// ErrApplyDiffFallback is returned to indicate that a normal ApplyDiff is applied as a fallback from Naive diff writer.
	ErrApplyDiffFallback = fmt.Errorf("Fall back to normal ApplyDiff")
	backingFs            = "<unknown>"
)

// ApplyDiffProtoDriver wraps the ProtoDriver by extending the interface with ApplyDiff method.
type ApplyDiffProtoDriver interface ***REMOVED***
	graphdriver.ProtoDriver
	// ApplyDiff writes the diff to the archive for the given id and parent id.
	// It returns the size in bytes written if successful, an error ErrApplyDiffFallback is returned otherwise.
	ApplyDiff(id, parent string, diff io.Reader) (size int64, err error)
***REMOVED***

type naiveDiffDriverWithApply struct ***REMOVED***
	graphdriver.Driver
	applyDiff ApplyDiffProtoDriver
***REMOVED***

// NaiveDiffDriverWithApply returns a NaiveDiff driver with custom ApplyDiff.
func NaiveDiffDriverWithApply(driver ApplyDiffProtoDriver, uidMaps, gidMaps []idtools.IDMap) graphdriver.Driver ***REMOVED***
	return &naiveDiffDriverWithApply***REMOVED***
		Driver:    graphdriver.NewNaiveDiffDriver(driver, uidMaps, gidMaps),
		applyDiff: driver,
	***REMOVED***
***REMOVED***

// ApplyDiff creates a diff layer with either the NaiveDiffDriver or with a fallback.
func (d *naiveDiffDriverWithApply) ApplyDiff(id, parent string, diff io.Reader) (int64, error) ***REMOVED***
	b, err := d.applyDiff.ApplyDiff(id, parent, diff)
	if err == ErrApplyDiffFallback ***REMOVED***
		return d.Driver.ApplyDiff(id, parent, diff)
	***REMOVED***
	return b, err
***REMOVED***

// This backend uses the overlay union filesystem for containers
// plus hard link file sharing for images.

// Each container/image can have a "root" subdirectory which is a plain
// filesystem hierarchy, or they can use overlay.

// If they use overlay there is a "upper" directory and a "lower-id"
// file, as well as "merged" and "work" directories. The "upper"
// directory has the upper layer of the overlay, and "lower-id" contains
// the id of the parent whose "root" directory shall be used as the lower
// layer in the overlay. The overlay itself is mounted in the "merged"
// directory, and the "work" dir is needed for overlay to work.

// When an overlay layer is created there are two cases, either the
// parent has a "root" dir, then we start out with an empty "upper"
// directory overlaid on the parents root. This is typically the
// case with the init layer of a container which is based on an image.
// If there is no "root" in the parent, we inherit the lower-id from
// the parent and start by making a copy in the parent's "upper" dir.
// This is typically the case for a container layer which copies
// its parent -init upper layer.

// Additionally we also have a custom implementation of ApplyLayer
// which makes a recursive copy of the parent "root" layer using
// hardlinks to share file data, and then applies the layer on top
// of that. This means all child images share file (but not directory)
// data with the parent.

// Driver contains information about the home directory and the list of active mounts that are created using this driver.
type Driver struct ***REMOVED***
	home          string
	uidMaps       []idtools.IDMap
	gidMaps       []idtools.IDMap
	ctr           *graphdriver.RefCounter
	supportsDType bool
	locker        *locker.Locker
***REMOVED***

func init() ***REMOVED***
	graphdriver.Register("overlay", Init)
***REMOVED***

// Init returns the NaiveDiffDriver, a native diff driver for overlay filesystem.
// If overlay filesystem is not supported on the host, the error
// graphdriver.ErrNotSupported is returned.
// If an overlay filesystem is not supported over an existing filesystem then
// error graphdriver.ErrIncompatibleFS is returned.
func Init(home string, options []string, uidMaps, gidMaps []idtools.IDMap) (graphdriver.Driver, error) ***REMOVED***

	if err := supportsOverlay(); err != nil ***REMOVED***
		return nil, graphdriver.ErrNotSupported
	***REMOVED***

	// Perform feature detection on /var/lib/docker/overlay if it's an existing directory.
	// This covers situations where /var/lib/docker/overlay is a mount, and on a different
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
	case graphdriver.FsMagicAufs, graphdriver.FsMagicBtrfs, graphdriver.FsMagicEcryptfs, graphdriver.FsMagicNfsFs, graphdriver.FsMagicOverlay, graphdriver.FsMagicZfs:
		logrus.Errorf("'overlay' is not supported over %s", backingFs)
		return nil, graphdriver.ErrIncompatibleFS
	***REMOVED***

	supportsDType, err := fsutils.SupportsDType(testdir)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !supportsDType ***REMOVED***
		if !graphdriver.IsInitialized(home) ***REMOVED***
			return nil, overlayutils.ErrDTypeNotSupported("overlay", backingFs)
		***REMOVED***
		// allow running without d_type only for existing setups (#27443)
		logrus.Warn(overlayutils.ErrDTypeNotSupported("overlay", backingFs))
	***REMOVED***

	rootUID, rootGID, err := idtools.GetRootUIDGID(uidMaps, gidMaps)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Create the driver home dir
	if err := idtools.MkdirAllAndChown(home, 0700, idtools.IDPair***REMOVED***rootUID, rootGID***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	d := &Driver***REMOVED***
		home:          home,
		uidMaps:       uidMaps,
		gidMaps:       gidMaps,
		ctr:           graphdriver.NewRefCounter(graphdriver.NewFsChecker(graphdriver.FsMagicOverlay)),
		supportsDType: supportsDType,
		locker:        locker.New(),
	***REMOVED***

	return NaiveDiffDriverWithApply(d, uidMaps, gidMaps), nil
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

func (d *Driver) String() string ***REMOVED***
	return "overlay"
***REMOVED***

// Status returns current driver information in a two dimensional string array.
// Output contains "Backing Filesystem" used in this implementation.
func (d *Driver) Status() [][2]string ***REMOVED***
	return [][2]string***REMOVED***
		***REMOVED***"Backing Filesystem", backingFs***REMOVED***,
		***REMOVED***"Supports d_type", strconv.FormatBool(d.supportsDType)***REMOVED***,
	***REMOVED***
***REMOVED***

// GetMetadata returns metadata about the overlay driver such as root,
// LowerDir, UpperDir, WorkDir and MergeDir used to store data.
func (d *Driver) GetMetadata(id string) (map[string]string, error) ***REMOVED***
	dir := d.dir(id)
	if _, err := os.Stat(dir); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	metadata := make(map[string]string)

	// If id has a root, it is an image
	rootDir := path.Join(dir, "root")
	if _, err := os.Stat(rootDir); err == nil ***REMOVED***
		metadata["RootDir"] = rootDir
		return metadata, nil
	***REMOVED***

	lowerID, err := ioutil.ReadFile(path.Join(dir, "lower-id"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	metadata["LowerDir"] = path.Join(d.dir(string(lowerID)), "root")
	metadata["UpperDir"] = path.Join(dir, "upper")
	metadata["WorkDir"] = path.Join(dir, "work")
	metadata["MergedDir"] = path.Join(dir, "merged")

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
	return d.Create(id, parent, opts)
***REMOVED***

// Create is used to create the upper, lower, and merge directories required for overlay fs for a given id.
// The parent filesystem is used to configure these directories for the overlay.
func (d *Driver) Create(id, parent string, opts *graphdriver.CreateOpts) (retErr error) ***REMOVED***

	if opts != nil && len(opts.StorageOpt) != 0 ***REMOVED***
		return fmt.Errorf("--storage-opt is not supported for overlay")
	***REMOVED***

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

	// Toplevel images are just a "root" dir
	if parent == "" ***REMOVED***
		return idtools.MkdirAndChown(path.Join(dir, "root"), 0755, root)
	***REMOVED***

	parentDir := d.dir(parent)

	// Ensure parent exists
	if _, err := os.Lstat(parentDir); err != nil ***REMOVED***
		return err
	***REMOVED***

	// If parent has a root, just do an overlay to it
	parentRoot := path.Join(parentDir, "root")

	if s, err := os.Lstat(parentRoot); err == nil ***REMOVED***
		if err := idtools.MkdirAndChown(path.Join(dir, "upper"), s.Mode(), root); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := idtools.MkdirAndChown(path.Join(dir, "work"), 0700, root); err != nil ***REMOVED***
			return err
		***REMOVED***
		return ioutil.WriteFile(path.Join(dir, "lower-id"), []byte(parent), 0666)
	***REMOVED***

	// Otherwise, copy the upper and the lower-id from the parent

	lowerID, err := ioutil.ReadFile(path.Join(parentDir, "lower-id"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := ioutil.WriteFile(path.Join(dir, "lower-id"), lowerID, 0666); err != nil ***REMOVED***
		return err
	***REMOVED***

	parentUpperDir := path.Join(parentDir, "upper")
	s, err := os.Lstat(parentUpperDir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	upperDir := path.Join(dir, "upper")
	if err := idtools.MkdirAndChown(upperDir, s.Mode(), root); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := idtools.MkdirAndChown(path.Join(dir, "work"), 0700, root); err != nil ***REMOVED***
		return err
	***REMOVED***

	return copy.DirCopy(parentUpperDir, upperDir, copy.Content, true)
***REMOVED***

func (d *Driver) dir(id string) string ***REMOVED***
	return path.Join(d.home, id)
***REMOVED***

// Remove cleans the directories that are created for this id.
func (d *Driver) Remove(id string) error ***REMOVED***
	d.locker.Lock(id)
	defer d.locker.Unlock(id)
	return system.EnsureRemoveAll(d.dir(id))
***REMOVED***

// Get creates and mounts the required file system for the given id and returns the mount path.
func (d *Driver) Get(id, mountLabel string) (_ containerfs.ContainerFS, err error) ***REMOVED***
	d.locker.Lock(id)
	defer d.locker.Unlock(id)
	dir := d.dir(id)
	if _, err := os.Stat(dir); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// If id has a root, just return it
	rootDir := path.Join(dir, "root")
	if _, err := os.Stat(rootDir); err == nil ***REMOVED***
		return containerfs.NewLocalContainerFS(rootDir), nil
	***REMOVED***

	mergedDir := path.Join(dir, "merged")
	if count := d.ctr.Increment(mergedDir); count > 1 ***REMOVED***
		return containerfs.NewLocalContainerFS(mergedDir), nil
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if c := d.ctr.Decrement(mergedDir); c <= 0 ***REMOVED***
				if mntErr := unix.Unmount(mergedDir, 0); mntErr != nil ***REMOVED***
					logrus.Debugf("Failed to unmount %s: %v: %v", id, mntErr, err)
				***REMOVED***
				// Cleanup the created merged directory; see the comment in Put's rmdir
				if rmErr := unix.Rmdir(mergedDir); rmErr != nil && !os.IsNotExist(rmErr) ***REMOVED***
					logrus.Warnf("Failed to remove %s: %v: %v", id, rmErr, err)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	lowerID, err := ioutil.ReadFile(path.Join(dir, "lower-id"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	rootUID, rootGID, err := idtools.GetRootUIDGID(d.uidMaps, d.gidMaps)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := idtools.MkdirAndChown(mergedDir, 0700, idtools.IDPair***REMOVED***rootUID, rootGID***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var (
		lowerDir = path.Join(d.dir(string(lowerID)), "root")
		upperDir = path.Join(dir, "upper")
		workDir  = path.Join(dir, "work")
		opts     = fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir)
	)
	if err := unix.Mount("overlay", mergedDir, "overlay", 0, label.FormatMountLabel(opts, mountLabel)); err != nil ***REMOVED***
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
	// If id has a root, just return
	if _, err := os.Stat(path.Join(d.dir(id), "root")); err == nil ***REMOVED***
		return nil
	***REMOVED***
	mountpoint := path.Join(d.dir(id), "merged")
	if count := d.ctr.Decrement(mountpoint); count > 0 ***REMOVED***
		return nil
	***REMOVED***
	if err := unix.Unmount(mountpoint, unix.MNT_DETACH); err != nil ***REMOVED***
		logrus.Debugf("Failed to unmount %s overlay: %v", id, err)
	***REMOVED***

	// Remove the mountpoint here. Removing the mountpoint (in newer kernels)
	// will cause all other instances of this mount in other mount namespaces
	// to be unmounted. This is necessary to avoid cases where an overlay mount
	// that is present in another namespace will cause subsequent mounts
	// operations to fail with ebusy.  We ignore any errors here because this may
	// fail on older kernels which don't have
	// torvalds/linux@8ed936b5671bfb33d89bc60bdcc7cf0470ba52fe applied.
	if err := unix.Rmdir(mountpoint); err != nil ***REMOVED***
		logrus.Debugf("Failed to remove %s overlay: %v", id, err)
	***REMOVED***
	return nil
***REMOVED***

// ApplyDiff applies the new layer on top of the root, if parent does not exist with will return an ErrApplyDiffFallback error.
func (d *Driver) ApplyDiff(id string, parent string, diff io.Reader) (size int64, err error) ***REMOVED***
	dir := d.dir(id)

	if parent == "" ***REMOVED***
		return 0, ErrApplyDiffFallback
	***REMOVED***

	parentRootDir := path.Join(d.dir(parent), "root")
	if _, err := os.Stat(parentRootDir); err != nil ***REMOVED***
		return 0, ErrApplyDiffFallback
	***REMOVED***

	// We now know there is a parent, and it has a "root" directory containing
	// the full root filesystem. We can just hardlink it and apply the
	// layer. This relies on two things:
	// 1) ApplyDiff is only run once on a clean (no writes to upper layer) container
	// 2) ApplyDiff doesn't do any in-place writes to files (would break hardlinks)
	// These are all currently true and are not expected to break

	tmpRootDir, err := ioutil.TempDir(dir, "tmproot")
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			os.RemoveAll(tmpRootDir)
		***REMOVED*** else ***REMOVED***
			os.RemoveAll(path.Join(dir, "upper"))
			os.RemoveAll(path.Join(dir, "work"))
			os.RemoveAll(path.Join(dir, "merged"))
			os.RemoveAll(path.Join(dir, "lower-id"))
		***REMOVED***
	***REMOVED***()

	if err = copy.DirCopy(parentRootDir, tmpRootDir, copy.Hardlink, true); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	options := &archive.TarOptions***REMOVED***UIDMaps: d.uidMaps, GIDMaps: d.gidMaps***REMOVED***
	if size, err = graphdriver.ApplyUncompressedLayer(tmpRootDir, diff, options); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	rootDir := path.Join(dir, "root")
	if err := os.Rename(tmpRootDir, rootDir); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return
***REMOVED***

// Exists checks to see if the id is already mounted.
func (d *Driver) Exists(id string) bool ***REMOVED***
	_, err := os.Stat(d.dir(id))
	return err == nil
***REMOVED***
