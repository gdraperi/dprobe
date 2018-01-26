package daemon

import (
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/chrootarchive"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/system"
	"github.com/pkg/errors"
)

// ErrExtractPointNotDirectory is used to convey that the operation to extract
// a tar archive to a directory in a container has failed because the specified
// path does not refer to a directory.
var ErrExtractPointNotDirectory = errors.New("extraction point is not a directory")

// The daemon will use the following interfaces if the container fs implements
// these for optimized copies to and from the container.
type extractor interface ***REMOVED***
	ExtractArchive(src io.Reader, dst string, opts *archive.TarOptions) error
***REMOVED***

type archiver interface ***REMOVED***
	ArchivePath(src string, opts *archive.TarOptions) (io.ReadCloser, error)
***REMOVED***

// helper functions to extract or archive
func extractArchive(i interface***REMOVED******REMOVED***, src io.Reader, dst string, opts *archive.TarOptions) error ***REMOVED***
	if ea, ok := i.(extractor); ok ***REMOVED***
		return ea.ExtractArchive(src, dst, opts)
	***REMOVED***
	return chrootarchive.Untar(src, dst, opts)
***REMOVED***

func archivePath(i interface***REMOVED******REMOVED***, src string, opts *archive.TarOptions) (io.ReadCloser, error) ***REMOVED***
	if ap, ok := i.(archiver); ok ***REMOVED***
		return ap.ArchivePath(src, opts)
	***REMOVED***
	return archive.TarWithOptions(src, opts)
***REMOVED***

// ContainerCopy performs a deprecated operation of archiving the resource at
// the specified path in the container identified by the given name.
func (daemon *Daemon) ContainerCopy(name string, res string) (io.ReadCloser, error) ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Make sure an online file-system operation is permitted.
	if err := daemon.isOnlineFSOperationPermitted(container); err != nil ***REMOVED***
		return nil, errdefs.System(err)
	***REMOVED***

	data, err := daemon.containerCopy(container, res)
	if err == nil ***REMOVED***
		return data, nil
	***REMOVED***

	if os.IsNotExist(err) ***REMOVED***
		return nil, containerFileNotFound***REMOVED***res, name***REMOVED***
	***REMOVED***
	return nil, errdefs.System(err)
***REMOVED***

// ContainerStatPath stats the filesystem resource at the specified path in the
// container identified by the given name.
func (daemon *Daemon) ContainerStatPath(name string, path string) (stat *types.ContainerPathStat, err error) ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Make sure an online file-system operation is permitted.
	if err := daemon.isOnlineFSOperationPermitted(container); err != nil ***REMOVED***
		return nil, errdefs.System(err)
	***REMOVED***

	stat, err = daemon.containerStatPath(container, path)
	if err == nil ***REMOVED***
		return stat, nil
	***REMOVED***

	if os.IsNotExist(err) ***REMOVED***
		return nil, containerFileNotFound***REMOVED***path, name***REMOVED***
	***REMOVED***
	return nil, errdefs.System(err)
***REMOVED***

// ContainerArchivePath creates an archive of the filesystem resource at the
// specified path in the container identified by the given name. Returns a
// tar archive of the resource and whether it was a directory or a single file.
func (daemon *Daemon) ContainerArchivePath(name string, path string) (content io.ReadCloser, stat *types.ContainerPathStat, err error) ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// Make sure an online file-system operation is permitted.
	if err := daemon.isOnlineFSOperationPermitted(container); err != nil ***REMOVED***
		return nil, nil, errdefs.System(err)
	***REMOVED***

	content, stat, err = daemon.containerArchivePath(container, path)
	if err == nil ***REMOVED***
		return content, stat, nil
	***REMOVED***

	if os.IsNotExist(err) ***REMOVED***
		return nil, nil, containerFileNotFound***REMOVED***path, name***REMOVED***
	***REMOVED***
	return nil, nil, errdefs.System(err)
***REMOVED***

// ContainerExtractToDir extracts the given archive to the specified location
// in the filesystem of the container identified by the given name. The given
// path must be of a directory in the container. If it is not, the error will
// be ErrExtractPointNotDirectory. If noOverwriteDirNonDir is true then it will
// be an error if unpacking the given content would cause an existing directory
// to be replaced with a non-directory and vice versa.
func (daemon *Daemon) ContainerExtractToDir(name, path string, copyUIDGID, noOverwriteDirNonDir bool, content io.Reader) error ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Make sure an online file-system operation is permitted.
	if err := daemon.isOnlineFSOperationPermitted(container); err != nil ***REMOVED***
		return errdefs.System(err)
	***REMOVED***

	err = daemon.containerExtractToDir(container, path, copyUIDGID, noOverwriteDirNonDir, content)
	if err == nil ***REMOVED***
		return nil
	***REMOVED***

	if os.IsNotExist(err) ***REMOVED***
		return containerFileNotFound***REMOVED***path, name***REMOVED***
	***REMOVED***
	return errdefs.System(err)
***REMOVED***

// containerStatPath stats the filesystem resource at the specified path in this
// container. Returns stat info about the resource.
func (daemon *Daemon) containerStatPath(container *container.Container, path string) (stat *types.ContainerPathStat, err error) ***REMOVED***
	container.Lock()
	defer container.Unlock()

	if err = daemon.Mount(container); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer daemon.Unmount(container)

	err = daemon.mountVolumes(container)
	defer container.DetachAndUnmount(daemon.LogVolumeEvent)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Normalize path before sending to rootfs
	path = container.BaseFS.FromSlash(path)

	resolvedPath, absPath, err := container.ResolvePath(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return container.StatPath(resolvedPath, absPath)
***REMOVED***

// containerArchivePath creates an archive of the filesystem resource at the specified
// path in this container. Returns a tar archive of the resource and stat info
// about the resource.
func (daemon *Daemon) containerArchivePath(container *container.Container, path string) (content io.ReadCloser, stat *types.ContainerPathStat, err error) ***REMOVED***
	container.Lock()

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			// Wait to unlock the container until the archive is fully read
			// (see the ReadCloseWrapper func below) or if there is an error
			// before that occurs.
			container.Unlock()
		***REMOVED***
	***REMOVED***()

	if err = daemon.Mount(container); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			// unmount any volumes
			container.DetachAndUnmount(daemon.LogVolumeEvent)
			// unmount the container's rootfs
			daemon.Unmount(container)
		***REMOVED***
	***REMOVED***()

	if err = daemon.mountVolumes(container); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// Normalize path before sending to rootfs
	path = container.BaseFS.FromSlash(path)

	resolvedPath, absPath, err := container.ResolvePath(path)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	stat, err = container.StatPath(resolvedPath, absPath)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// We need to rebase the archive entries if the last element of the
	// resolved path was a symlink that was evaluated and is now different
	// than the requested path. For example, if the given path was "/foo/bar/",
	// but it resolved to "/var/lib/docker/containers/***REMOVED***id***REMOVED***/foo/baz/", we want
	// to ensure that the archive entries start with "bar" and not "baz". This
	// also catches the case when the root directory of the container is
	// requested: we want the archive entries to start with "/" and not the
	// container ID.
	driver := container.BaseFS

	// Get the source and the base paths of the container resolved path in order
	// to get the proper tar options for the rebase tar.
	resolvedPath = driver.Clean(resolvedPath)
	if driver.Base(resolvedPath) == "." ***REMOVED***
		resolvedPath += string(driver.Separator()) + "."
	***REMOVED***
	sourceDir, sourceBase := driver.Dir(resolvedPath), driver.Base(resolvedPath)
	opts := archive.TarResourceRebaseOpts(sourceBase, driver.Base(absPath))

	data, err := archivePath(driver, sourceDir, opts)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	content = ioutils.NewReadCloserWrapper(data, func() error ***REMOVED***
		err := data.Close()
		container.DetachAndUnmount(daemon.LogVolumeEvent)
		daemon.Unmount(container)
		container.Unlock()
		return err
	***REMOVED***)

	daemon.LogContainerEvent(container, "archive-path")

	return content, stat, nil
***REMOVED***

// containerExtractToDir extracts the given tar archive to the specified location in the
// filesystem of this container. The given path must be of a directory in the
// container. If it is not, the error will be ErrExtractPointNotDirectory. If
// noOverwriteDirNonDir is true then it will be an error if unpacking the
// given content would cause an existing directory to be replaced with a non-
// directory and vice versa.
func (daemon *Daemon) containerExtractToDir(container *container.Container, path string, copyUIDGID, noOverwriteDirNonDir bool, content io.Reader) (err error) ***REMOVED***
	container.Lock()
	defer container.Unlock()

	if err = daemon.Mount(container); err != nil ***REMOVED***
		return err
	***REMOVED***
	defer daemon.Unmount(container)

	err = daemon.mountVolumes(container)
	defer container.DetachAndUnmount(daemon.LogVolumeEvent)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Normalize path before sending to rootfs'
	path = container.BaseFS.FromSlash(path)
	driver := container.BaseFS

	// Check if a drive letter supplied, it must be the system drive. No-op except on Windows
	path, err = system.CheckSystemDriveAndRemoveDriveLetter(path, driver)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// The destination path needs to be resolved to a host path, with all
	// symbolic links followed in the scope of the container's rootfs. Note
	// that we do not use `container.ResolvePath(path)` here because we need
	// to also evaluate the last path element if it is a symlink. This is so
	// that you can extract an archive to a symlink that points to a directory.

	// Consider the given path as an absolute path in the container.
	absPath := archive.PreserveTrailingDotOrSeparator(
		driver.Join(string(driver.Separator()), path),
		path,
		driver.Separator())

	// This will evaluate the last path element if it is a symlink.
	resolvedPath, err := container.GetResourcePath(absPath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	stat, err := driver.Lstat(resolvedPath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !stat.IsDir() ***REMOVED***
		return ErrExtractPointNotDirectory
	***REMOVED***

	// Need to check if the path is in a volume. If it is, it cannot be in a
	// read-only volume. If it is not in a volume, the container cannot be
	// configured with a read-only rootfs.

	// Use the resolved path relative to the container rootfs as the new
	// absPath. This way we fully follow any symlinks in a volume that may
	// lead back outside the volume.
	//
	// The Windows implementation of filepath.Rel in golang 1.4 does not
	// support volume style file path semantics. On Windows when using the
	// filter driver, we are guaranteed that the path will always be
	// a volume file path.
	var baseRel string
	if strings.HasPrefix(resolvedPath, `\\?\Volume***REMOVED***`) ***REMOVED***
		if strings.HasPrefix(resolvedPath, driver.Path()) ***REMOVED***
			baseRel = resolvedPath[len(driver.Path()):]
			if baseRel[:1] == `\` ***REMOVED***
				baseRel = baseRel[1:]
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		baseRel, err = driver.Rel(driver.Path(), resolvedPath)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// Make it an absolute path.
	absPath = driver.Join(string(driver.Separator()), baseRel)

	// @ TODO: gupta-ak: Technically, this works since it no-ops
	// on Windows and the file system is local anyway on linux.
	// But eventually, it should be made driver aware.
	toVolume, err := checkIfPathIsInAVolume(container, absPath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !toVolume && container.HostConfig.ReadonlyRootfs ***REMOVED***
		return ErrRootFSReadOnly
	***REMOVED***

	options := daemon.defaultTarCopyOptions(noOverwriteDirNonDir)

	if copyUIDGID ***REMOVED***
		var err error
		// tarCopyOptions will appropriately pull in the right uid/gid for the
		// user/group and will set the options.
		options, err = daemon.tarCopyOptions(container, noOverwriteDirNonDir)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if err := extractArchive(driver, content, resolvedPath, options); err != nil ***REMOVED***
		return err
	***REMOVED***

	daemon.LogContainerEvent(container, "extract-to-dir")

	return nil
***REMOVED***

func (daemon *Daemon) containerCopy(container *container.Container, resource string) (rc io.ReadCloser, err error) ***REMOVED***
	if resource[0] == '/' || resource[0] == '\\' ***REMOVED***
		resource = resource[1:]
	***REMOVED***
	container.Lock()

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			// Wait to unlock the container until the archive is fully read
			// (see the ReadCloseWrapper func below) or if there is an error
			// before that occurs.
			container.Unlock()
		***REMOVED***
	***REMOVED***()

	if err := daemon.Mount(container); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			// unmount any volumes
			container.DetachAndUnmount(daemon.LogVolumeEvent)
			// unmount the container's rootfs
			daemon.Unmount(container)
		***REMOVED***
	***REMOVED***()

	if err := daemon.mountVolumes(container); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Normalize path before sending to rootfs
	resource = container.BaseFS.FromSlash(resource)
	driver := container.BaseFS

	basePath, err := container.GetResourcePath(resource)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	stat, err := driver.Stat(basePath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var filter []string
	if !stat.IsDir() ***REMOVED***
		d, f := driver.Split(basePath)
		basePath = d
		filter = []string***REMOVED***f***REMOVED***
	***REMOVED*** else ***REMOVED***
		filter = []string***REMOVED***driver.Base(basePath)***REMOVED***
		basePath = driver.Dir(basePath)
	***REMOVED***
	archive, err := archivePath(driver, basePath, &archive.TarOptions***REMOVED***
		Compression:  archive.Uncompressed,
		IncludeFiles: filter,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	reader := ioutils.NewReadCloserWrapper(archive, func() error ***REMOVED***
		err := archive.Close()
		container.DetachAndUnmount(daemon.LogVolumeEvent)
		daemon.Unmount(container)
		container.Unlock()
		return err
	***REMOVED***)
	daemon.LogContainerEvent(container, "copy")
	return reader, nil
***REMOVED***
