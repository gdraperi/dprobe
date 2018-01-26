package remotecontext

import (
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/builder"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/chrootarchive"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/tarsum"
	"github.com/pkg/errors"
)

type archiveContext struct ***REMOVED***
	root containerfs.ContainerFS
	sums tarsum.FileInfoSums
***REMOVED***

func (c *archiveContext) Close() error ***REMOVED***
	return c.root.RemoveAll(c.root.Path())
***REMOVED***

func convertPathError(err error, cleanpath string) error ***REMOVED***
	if err, ok := err.(*os.PathError); ok ***REMOVED***
		err.Path = cleanpath
		return err
	***REMOVED***
	return err
***REMOVED***

type modifiableContext interface ***REMOVED***
	builder.Source
	// Remove deletes the entry specified by `path`.
	// It is usual for directory entries to delete all its subentries.
	Remove(path string) error
***REMOVED***

// FromArchive returns a build source from a tar stream.
//
// It extracts the tar stream to a temporary folder that is deleted as soon as
// the Context is closed.
// As the extraction happens, a tarsum is calculated for every file, and the set of
// all those sums then becomes the source of truth for all operations on this Context.
//
// Closing tarStream has to be done by the caller.
func FromArchive(tarStream io.Reader) (builder.Source, error) ***REMOVED***
	root, err := ioutils.TempDir("", "docker-builder")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Assume local file system. Since it's coming from a tar file.
	tsc := &archiveContext***REMOVED***root: containerfs.NewLocalContainerFS(root)***REMOVED***

	// Make sure we clean-up upon error.  In the happy case the caller
	// is expected to manage the clean-up
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			tsc.Close()
		***REMOVED***
	***REMOVED***()

	decompressedStream, err := archive.DecompressStream(tarStream)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sum, err := tarsum.NewTarSum(decompressedStream, true, tarsum.Version1)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = chrootarchive.Untar(sum, root, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	tsc.sums = sum.GetSums()
	return tsc, nil
***REMOVED***

func (c *archiveContext) Root() containerfs.ContainerFS ***REMOVED***
	return c.root
***REMOVED***

func (c *archiveContext) Remove(path string) error ***REMOVED***
	_, fullpath, err := normalize(path, c.root)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return c.root.RemoveAll(fullpath)
***REMOVED***

func (c *archiveContext) Hash(path string) (string, error) ***REMOVED***
	cleanpath, fullpath, err := normalize(path, c.root)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	rel, err := c.root.Rel(c.root.Path(), fullpath)
	if err != nil ***REMOVED***
		return "", convertPathError(err, cleanpath)
	***REMOVED***

	// Use the checksum of the followed path(not the possible symlink) because
	// this is the file that is actually copied.
	if tsInfo := c.sums.GetFile(filepath.ToSlash(rel)); tsInfo != nil ***REMOVED***
		return tsInfo.Sum(), nil
	***REMOVED***
	// We set sum to path by default for the case where GetFile returns nil.
	// The usual case is if relative path is empty.
	return path, nil // backwards compat TODO: see if really needed
***REMOVED***

func normalize(path string, root containerfs.ContainerFS) (cleanPath, fullPath string, err error) ***REMOVED***
	cleanPath = root.Clean(string(root.Separator()) + path)[1:]
	fullPath, err = root.ResolveScopedPath(path, true)
	if err != nil ***REMOVED***
		return "", "", errors.Wrapf(err, "forbidden path outside the build context: %s (%s)", path, cleanPath)
	***REMOVED***
	return
***REMOVED***
