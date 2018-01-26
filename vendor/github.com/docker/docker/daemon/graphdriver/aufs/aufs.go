// +build linux

/*

aufs driver directory structure

  .
  ├── layers // Metadata of layers
  │   ├── 1
  │   ├── 2
  │   └── 3
  ├── diff  // Content of the layer
  │   ├── 1  // Contains layers that need to be mounted for the id
  │   ├── 2
  │   └── 3
  └── mnt    // Mount points for the rw layers to be mounted
      ├── 1
      ├── 2
      └── 3

*/

package aufs

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/chrootarchive"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/directory"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/locker"
	mountpk "github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/system"
	rsystem "github.com/opencontainers/runc/libcontainer/system"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vbatts/tar-split/tar/storage"
	"golang.org/x/sys/unix"
)

var (
	// ErrAufsNotSupported is returned if aufs is not supported by the host.
	ErrAufsNotSupported = fmt.Errorf("AUFS was not found in /proc/filesystems")
	// ErrAufsNested means aufs cannot be used bc we are in a user namespace
	ErrAufsNested = fmt.Errorf("AUFS cannot be used in non-init user namespace")
	backingFs     = "<unknown>"

	enableDirpermLock sync.Once
	enableDirperm     bool
)

func init() ***REMOVED***
	graphdriver.Register("aufs", Init)
***REMOVED***

// Driver contains information about the filesystem mounted.
type Driver struct ***REMOVED***
	sync.Mutex
	root          string
	uidMaps       []idtools.IDMap
	gidMaps       []idtools.IDMap
	ctr           *graphdriver.RefCounter
	pathCacheLock sync.Mutex
	pathCache     map[string]string
	naiveDiff     graphdriver.DiffDriver
	locker        *locker.Locker
***REMOVED***

// Init returns a new AUFS driver.
// An error is returned if AUFS is not supported.
func Init(root string, options []string, uidMaps, gidMaps []idtools.IDMap) (graphdriver.Driver, error) ***REMOVED***

	// Try to load the aufs kernel module
	if err := supportsAufs(); err != nil ***REMOVED***
		return nil, graphdriver.ErrNotSupported
	***REMOVED***

	// Perform feature detection on /var/lib/docker/aufs if it's an existing directory.
	// This covers situations where /var/lib/docker/aufs is a mount, and on a different
	// filesystem than /var/lib/docker.
	// If the path does not exist, fall back to using /var/lib/docker for feature detection.
	testdir := root
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
	case graphdriver.FsMagicAufs, graphdriver.FsMagicBtrfs, graphdriver.FsMagicEcryptfs:
		logrus.Errorf("AUFS is not supported over %s", backingFs)
		return nil, graphdriver.ErrIncompatibleFS
	***REMOVED***

	paths := []string***REMOVED***
		"mnt",
		"diff",
		"layers",
	***REMOVED***

	a := &Driver***REMOVED***
		root:      root,
		uidMaps:   uidMaps,
		gidMaps:   gidMaps,
		pathCache: make(map[string]string),
		ctr:       graphdriver.NewRefCounter(graphdriver.NewFsChecker(graphdriver.FsMagicAufs)),
		locker:    locker.New(),
	***REMOVED***

	rootUID, rootGID, err := idtools.GetRootUIDGID(uidMaps, gidMaps)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Create the root aufs driver dir
	if err := idtools.MkdirAllAndChown(root, 0700, idtools.IDPair***REMOVED***UID: rootUID, GID: rootGID***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Populate the dir structure
	for _, p := range paths ***REMOVED***
		if err := idtools.MkdirAllAndChown(path.Join(root, p), 0700, idtools.IDPair***REMOVED***UID: rootUID, GID: rootGID***REMOVED***); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	logger := logrus.WithFields(logrus.Fields***REMOVED***
		"module": "graphdriver",
		"driver": "aufs",
	***REMOVED***)

	for _, path := range []string***REMOVED***"mnt", "diff"***REMOVED*** ***REMOVED***
		p := filepath.Join(root, path)
		entries, err := ioutil.ReadDir(p)
		if err != nil ***REMOVED***
			logger.WithError(err).WithField("dir", p).Error("error reading dir entries")
			continue
		***REMOVED***
		for _, entry := range entries ***REMOVED***
			if !entry.IsDir() ***REMOVED***
				continue
			***REMOVED***
			if strings.HasSuffix(entry.Name(), "-removing") ***REMOVED***
				logger.WithField("dir", entry.Name()).Debug("Cleaning up stale layer dir")
				if err := system.EnsureRemoveAll(filepath.Join(p, entry.Name())); err != nil ***REMOVED***
					logger.WithField("dir", entry.Name()).WithError(err).Error("Error removing stale layer dir")
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	a.naiveDiff = graphdriver.NewNaiveDiffDriver(a, uidMaps, gidMaps)
	return a, nil
***REMOVED***

// Return a nil error if the kernel supports aufs
// We cannot modprobe because inside dind modprobe fails
// to run
func supportsAufs() error ***REMOVED***
	// We can try to modprobe aufs first before looking at
	// proc/filesystems for when aufs is supported
	exec.Command("modprobe", "aufs").Run()

	if rsystem.RunningInUserNS() ***REMOVED***
		return ErrAufsNested
	***REMOVED***

	f, err := os.Open("/proc/filesystems")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() ***REMOVED***
		if strings.Contains(s.Text(), "aufs") ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return ErrAufsNotSupported
***REMOVED***

func (a *Driver) rootPath() string ***REMOVED***
	return a.root
***REMOVED***

func (*Driver) String() string ***REMOVED***
	return "aufs"
***REMOVED***

// Status returns current information about the filesystem such as root directory, number of directories mounted, etc.
func (a *Driver) Status() [][2]string ***REMOVED***
	ids, _ := loadIds(path.Join(a.rootPath(), "layers"))
	return [][2]string***REMOVED***
		***REMOVED***"Root Dir", a.rootPath()***REMOVED***,
		***REMOVED***"Backing Filesystem", backingFs***REMOVED***,
		***REMOVED***"Dirs", fmt.Sprintf("%d", len(ids))***REMOVED***,
		***REMOVED***"Dirperm1 Supported", fmt.Sprintf("%v", useDirperm())***REMOVED***,
	***REMOVED***
***REMOVED***

// GetMetadata not implemented
func (a *Driver) GetMetadata(id string) (map[string]string, error) ***REMOVED***
	return nil, nil
***REMOVED***

// Exists returns true if the given id is registered with
// this driver
func (a *Driver) Exists(id string) bool ***REMOVED***
	if _, err := os.Lstat(path.Join(a.rootPath(), "layers", id)); err != nil ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// CreateReadWrite creates a layer that is writable for use as a container
// file system.
func (a *Driver) CreateReadWrite(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	return a.Create(id, parent, opts)
***REMOVED***

// Create three folders for each id
// mnt, layers, and diff
func (a *Driver) Create(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***

	if opts != nil && len(opts.StorageOpt) != 0 ***REMOVED***
		return fmt.Errorf("--storage-opt is not supported for aufs")
	***REMOVED***

	if err := a.createDirsFor(id); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Write the layers metadata
	f, err := os.Create(path.Join(a.rootPath(), "layers", id))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	if parent != "" ***REMOVED***
		ids, err := getParentIDs(a.rootPath(), parent)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if _, err := fmt.Fprintln(f, parent); err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, i := range ids ***REMOVED***
			if _, err := fmt.Fprintln(f, i); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// createDirsFor creates two directories for the given id.
// mnt and diff
func (a *Driver) createDirsFor(id string) error ***REMOVED***
	paths := []string***REMOVED***
		"mnt",
		"diff",
	***REMOVED***

	rootUID, rootGID, err := idtools.GetRootUIDGID(a.uidMaps, a.gidMaps)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// Directory permission is 0755.
	// The path of directories are <aufs_root_path>/mnt/<image_id>
	// and <aufs_root_path>/diff/<image_id>
	for _, p := range paths ***REMOVED***
		if err := idtools.MkdirAllAndChown(path.Join(a.rootPath(), p, id), 0755, idtools.IDPair***REMOVED***UID: rootUID, GID: rootGID***REMOVED***); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Remove will unmount and remove the given id.
func (a *Driver) Remove(id string) error ***REMOVED***
	a.locker.Lock(id)
	defer a.locker.Unlock(id)
	a.pathCacheLock.Lock()
	mountpoint, exists := a.pathCache[id]
	a.pathCacheLock.Unlock()
	if !exists ***REMOVED***
		mountpoint = a.getMountpoint(id)
	***REMOVED***

	logger := logrus.WithFields(logrus.Fields***REMOVED***
		"module": "graphdriver",
		"driver": "aufs",
		"layer":  id,
	***REMOVED***)

	var retries int
	for ***REMOVED***
		mounted, err := a.mounted(mountpoint)
		if err != nil ***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				break
			***REMOVED***
			return err
		***REMOVED***
		if !mounted ***REMOVED***
			break
		***REMOVED***

		err = a.unmount(mountpoint)
		if err == nil ***REMOVED***
			break
		***REMOVED***

		if err != unix.EBUSY ***REMOVED***
			return errors.Wrapf(err, "aufs: unmount error: %s", mountpoint)
		***REMOVED***
		if retries >= 5 ***REMOVED***
			return errors.Wrapf(err, "aufs: unmount error after retries: %s", mountpoint)
		***REMOVED***
		// If unmount returns EBUSY, it could be a transient error. Sleep and retry.
		retries++
		logger.Warnf("unmount failed due to EBUSY: retry count: %d", retries)
		time.Sleep(100 * time.Millisecond)
	***REMOVED***

	// Remove the layers file for the id
	if err := os.Remove(path.Join(a.rootPath(), "layers", id)); err != nil && !os.IsNotExist(err) ***REMOVED***
		return errors.Wrapf(err, "error removing layers dir for %s", id)
	***REMOVED***

	if err := atomicRemove(a.getDiffPath(id)); err != nil ***REMOVED***
		return errors.Wrapf(err, "could not remove diff path for id %s", id)
	***REMOVED***

	// Atomically remove each directory in turn by first moving it out of the
	// way (so that docker doesn't find it anymore) before doing removal of
	// the whole tree.
	if err := atomicRemove(mountpoint); err != nil ***REMOVED***
		if errors.Cause(err) == unix.EBUSY ***REMOVED***
			logger.WithField("dir", mountpoint).WithError(err).Warn("error performing atomic remove due to EBUSY")
		***REMOVED***
		return errors.Wrapf(err, "could not remove mountpoint for id %s", id)
	***REMOVED***

	a.pathCacheLock.Lock()
	delete(a.pathCache, id)
	a.pathCacheLock.Unlock()
	return nil
***REMOVED***

func atomicRemove(source string) error ***REMOVED***
	target := source + "-removing"

	err := os.Rename(source, target)
	switch ***REMOVED***
	case err == nil, os.IsNotExist(err):
	case os.IsExist(err):
		// Got error saying the target dir already exists, maybe the source doesn't exist due to a previous (failed) remove
		if _, e := os.Stat(source); !os.IsNotExist(e) ***REMOVED***
			return errors.Wrapf(err, "target rename dir '%s' exists but should not, this needs to be manually cleaned up")
		***REMOVED***
	default:
		return errors.Wrapf(err, "error preparing atomic delete")
	***REMOVED***

	return system.EnsureRemoveAll(target)
***REMOVED***

// Get returns the rootfs path for the id.
// This will mount the dir at its given path
func (a *Driver) Get(id, mountLabel string) (containerfs.ContainerFS, error) ***REMOVED***
	a.locker.Lock(id)
	defer a.locker.Unlock(id)
	parents, err := a.getParentLayerPaths(id)
	if err != nil && !os.IsNotExist(err) ***REMOVED***
		return nil, err
	***REMOVED***

	a.pathCacheLock.Lock()
	m, exists := a.pathCache[id]
	a.pathCacheLock.Unlock()

	if !exists ***REMOVED***
		m = a.getDiffPath(id)
		if len(parents) > 0 ***REMOVED***
			m = a.getMountpoint(id)
		***REMOVED***
	***REMOVED***
	if count := a.ctr.Increment(m); count > 1 ***REMOVED***
		return containerfs.NewLocalContainerFS(m), nil
	***REMOVED***

	// If a dir does not have a parent ( no layers )do not try to mount
	// just return the diff path to the data
	if len(parents) > 0 ***REMOVED***
		if err := a.mount(id, m, mountLabel, parents); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	a.pathCacheLock.Lock()
	a.pathCache[id] = m
	a.pathCacheLock.Unlock()
	return containerfs.NewLocalContainerFS(m), nil
***REMOVED***

// Put unmounts and updates list of active mounts.
func (a *Driver) Put(id string) error ***REMOVED***
	a.locker.Lock(id)
	defer a.locker.Unlock(id)
	a.pathCacheLock.Lock()
	m, exists := a.pathCache[id]
	if !exists ***REMOVED***
		m = a.getMountpoint(id)
		a.pathCache[id] = m
	***REMOVED***
	a.pathCacheLock.Unlock()
	if count := a.ctr.Decrement(m); count > 0 ***REMOVED***
		return nil
	***REMOVED***

	err := a.unmount(m)
	if err != nil ***REMOVED***
		logrus.Debugf("Failed to unmount %s aufs: %v", id, err)
	***REMOVED***
	return err
***REMOVED***

// isParent returns if the passed in parent is the direct parent of the passed in layer
func (a *Driver) isParent(id, parent string) bool ***REMOVED***
	parents, _ := getParentIDs(a.rootPath(), id)
	if parent == "" && len(parents) > 0 ***REMOVED***
		return false
	***REMOVED***
	return !(len(parents) > 0 && parent != parents[0])
***REMOVED***

// Diff produces an archive of the changes between the specified
// layer and its parent layer which may be "".
func (a *Driver) Diff(id, parent string) (io.ReadCloser, error) ***REMOVED***
	if !a.isParent(id, parent) ***REMOVED***
		return a.naiveDiff.Diff(id, parent)
	***REMOVED***

	// AUFS doesn't need the parent layer to produce a diff.
	return archive.TarWithOptions(path.Join(a.rootPath(), "diff", id), &archive.TarOptions***REMOVED***
		Compression:     archive.Uncompressed,
		ExcludePatterns: []string***REMOVED***archive.WhiteoutMetaPrefix + "*", "!" + archive.WhiteoutOpaqueDir***REMOVED***,
		UIDMaps:         a.uidMaps,
		GIDMaps:         a.gidMaps,
	***REMOVED***)
***REMOVED***

type fileGetNilCloser struct ***REMOVED***
	storage.FileGetter
***REMOVED***

func (f fileGetNilCloser) Close() error ***REMOVED***
	return nil
***REMOVED***

// DiffGetter returns a FileGetCloser that can read files from the directory that
// contains files for the layer differences. Used for direct access for tar-split.
func (a *Driver) DiffGetter(id string) (graphdriver.FileGetCloser, error) ***REMOVED***
	p := path.Join(a.rootPath(), "diff", id)
	return fileGetNilCloser***REMOVED***storage.NewPathFileGetter(p)***REMOVED***, nil
***REMOVED***

func (a *Driver) applyDiff(id string, diff io.Reader) error ***REMOVED***
	return chrootarchive.UntarUncompressed(diff, path.Join(a.rootPath(), "diff", id), &archive.TarOptions***REMOVED***
		UIDMaps: a.uidMaps,
		GIDMaps: a.gidMaps,
	***REMOVED***)
***REMOVED***

// DiffSize calculates the changes between the specified id
// and its parent and returns the size in bytes of the changes
// relative to its base filesystem directory.
func (a *Driver) DiffSize(id, parent string) (size int64, err error) ***REMOVED***
	if !a.isParent(id, parent) ***REMOVED***
		return a.naiveDiff.DiffSize(id, parent)
	***REMOVED***
	// AUFS doesn't need the parent layer to calculate the diff size.
	return directory.Size(path.Join(a.rootPath(), "diff", id))
***REMOVED***

// ApplyDiff extracts the changeset from the given diff into the
// layer with the specified id and parent, returning the size of the
// new layer in bytes.
func (a *Driver) ApplyDiff(id, parent string, diff io.Reader) (size int64, err error) ***REMOVED***
	if !a.isParent(id, parent) ***REMOVED***
		return a.naiveDiff.ApplyDiff(id, parent, diff)
	***REMOVED***

	// AUFS doesn't need the parent id to apply the diff if it is the direct parent.
	if err = a.applyDiff(id, diff); err != nil ***REMOVED***
		return
	***REMOVED***

	return a.DiffSize(id, parent)
***REMOVED***

// Changes produces a list of changes between the specified layer
// and its parent layer. If parent is "", then all changes will be ADD changes.
func (a *Driver) Changes(id, parent string) ([]archive.Change, error) ***REMOVED***
	if !a.isParent(id, parent) ***REMOVED***
		return a.naiveDiff.Changes(id, parent)
	***REMOVED***

	// AUFS doesn't have snapshots, so we need to get changes from all parent
	// layers.
	layers, err := a.getParentLayerPaths(id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return archive.Changes(layers, path.Join(a.rootPath(), "diff", id))
***REMOVED***

func (a *Driver) getParentLayerPaths(id string) ([]string, error) ***REMOVED***
	parentIds, err := getParentIDs(a.rootPath(), id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	layers := make([]string, len(parentIds))

	// Get the diff paths for all the parent ids
	for i, p := range parentIds ***REMOVED***
		layers[i] = path.Join(a.rootPath(), "diff", p)
	***REMOVED***
	return layers, nil
***REMOVED***

func (a *Driver) mount(id string, target string, mountLabel string, layers []string) error ***REMOVED***
	a.Lock()
	defer a.Unlock()

	// If the id is mounted or we get an error return
	if mounted, err := a.mounted(target); err != nil || mounted ***REMOVED***
		return err
	***REMOVED***

	rw := a.getDiffPath(id)

	if err := a.aufsMount(layers, rw, target, mountLabel); err != nil ***REMOVED***
		return fmt.Errorf("error creating aufs mount to %s: %v", target, err)
	***REMOVED***
	return nil
***REMOVED***

func (a *Driver) unmount(mountPath string) error ***REMOVED***
	a.Lock()
	defer a.Unlock()

	if mounted, err := a.mounted(mountPath); err != nil || !mounted ***REMOVED***
		return err
	***REMOVED***
	return Unmount(mountPath)
***REMOVED***

func (a *Driver) mounted(mountpoint string) (bool, error) ***REMOVED***
	return graphdriver.Mounted(graphdriver.FsMagicAufs, mountpoint)
***REMOVED***

// Cleanup aufs and unmount all mountpoints
func (a *Driver) Cleanup() error ***REMOVED***
	var dirs []string
	if err := filepath.Walk(a.mntPath(), func(path string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if !info.IsDir() ***REMOVED***
			return nil
		***REMOVED***
		dirs = append(dirs, path)
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, m := range dirs ***REMOVED***
		if err := a.unmount(m); err != nil ***REMOVED***
			logrus.Debugf("aufs error unmounting %s: %s", m, err)
		***REMOVED***
	***REMOVED***
	return mountpk.RecursiveUnmount(a.root)
***REMOVED***

func (a *Driver) aufsMount(ro []string, rw, target, mountLabel string) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			Unmount(target)
		***REMOVED***
	***REMOVED***()

	// Mount options are clipped to page size(4096 bytes). If there are more
	// layers then these are remounted individually using append.

	offset := 54
	if useDirperm() ***REMOVED***
		offset += len(",dirperm1")
	***REMOVED***
	b := make([]byte, unix.Getpagesize()-len(mountLabel)-offset) // room for xino & mountLabel
	bp := copy(b, fmt.Sprintf("br:%s=rw", rw))

	index := 0
	for ; index < len(ro); index++ ***REMOVED***
		layer := fmt.Sprintf(":%s=ro+wh", ro[index])
		if bp+len(layer) > len(b) ***REMOVED***
			break
		***REMOVED***
		bp += copy(b[bp:], layer)
	***REMOVED***

	opts := "dio,xino=/dev/shm/aufs.xino"
	if useDirperm() ***REMOVED***
		opts += ",dirperm1"
	***REMOVED***
	data := label.FormatMountLabel(fmt.Sprintf("%s,%s", string(b[:bp]), opts), mountLabel)
	if err = mount("none", target, "aufs", 0, data); err != nil ***REMOVED***
		return
	***REMOVED***

	for ; index < len(ro); index++ ***REMOVED***
		layer := fmt.Sprintf(":%s=ro+wh", ro[index])
		data := label.FormatMountLabel(fmt.Sprintf("append%s", layer), mountLabel)
		if err = mount("none", target, "aufs", unix.MS_REMOUNT, data); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

// useDirperm checks dirperm1 mount option can be used with the current
// version of aufs.
func useDirperm() bool ***REMOVED***
	enableDirpermLock.Do(func() ***REMOVED***
		base, err := ioutil.TempDir("", "docker-aufs-base")
		if err != nil ***REMOVED***
			logrus.Errorf("error checking dirperm1: %v", err)
			return
		***REMOVED***
		defer os.RemoveAll(base)

		union, err := ioutil.TempDir("", "docker-aufs-union")
		if err != nil ***REMOVED***
			logrus.Errorf("error checking dirperm1: %v", err)
			return
		***REMOVED***
		defer os.RemoveAll(union)

		opts := fmt.Sprintf("br:%s,dirperm1,xino=/dev/shm/aufs.xino", base)
		if err := mount("none", union, "aufs", 0, opts); err != nil ***REMOVED***
			return
		***REMOVED***
		enableDirperm = true
		if err := Unmount(union); err != nil ***REMOVED***
			logrus.Errorf("error checking dirperm1: failed to unmount %v", err)
		***REMOVED***
	***REMOVED***)
	return enableDirperm
***REMOVED***
