package containerfs

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
)

// TarFunc provides a function definition for a custom Tar function
type TarFunc func(string, *archive.TarOptions) (io.ReadCloser, error)

// UntarFunc provides a function definition for a custom Untar function
type UntarFunc func(io.Reader, string, *archive.TarOptions) error

// Archiver provides a similar implementation of the archive.Archiver package with the rootfs abstraction
type Archiver struct ***REMOVED***
	SrcDriver     Driver
	DstDriver     Driver
	Tar           TarFunc
	Untar         UntarFunc
	IDMappingsVar *idtools.IDMappings
***REMOVED***

// TarUntar is a convenience function which calls Tar and Untar, with the output of one piped into the other.
// If either Tar or Untar fails, TarUntar aborts and returns the error.
func (archiver *Archiver) TarUntar(src, dst string) error ***REMOVED***
	logrus.Debugf("TarUntar(%s %s)", src, dst)
	tarArchive, err := archiver.Tar(src, &archive.TarOptions***REMOVED***Compression: archive.Uncompressed***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer tarArchive.Close()
	options := &archive.TarOptions***REMOVED***
		UIDMaps: archiver.IDMappingsVar.UIDs(),
		GIDMaps: archiver.IDMappingsVar.GIDs(),
	***REMOVED***
	return archiver.Untar(tarArchive, dst, options)
***REMOVED***

// UntarPath untar a file from path to a destination, src is the source tar file path.
func (archiver *Archiver) UntarPath(src, dst string) error ***REMOVED***
	tarArchive, err := archiver.SrcDriver.Open(src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer tarArchive.Close()
	options := &archive.TarOptions***REMOVED***
		UIDMaps: archiver.IDMappingsVar.UIDs(),
		GIDMaps: archiver.IDMappingsVar.GIDs(),
	***REMOVED***
	return archiver.Untar(tarArchive, dst, options)
***REMOVED***

// CopyWithTar creates a tar archive of filesystem path `src`, and
// unpacks it at filesystem path `dst`.
// The archive is streamed directly with fixed buffering and no
// intermediary disk IO.
func (archiver *Archiver) CopyWithTar(src, dst string) error ***REMOVED***
	srcSt, err := archiver.SrcDriver.Stat(src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !srcSt.IsDir() ***REMOVED***
		return archiver.CopyFileWithTar(src, dst)
	***REMOVED***

	// if this archiver is set up with ID mapping we need to create
	// the new destination directory with the remapped root UID/GID pair
	// as owner
	rootIDs := archiver.IDMappingsVar.RootPair()
	// Create dst, copy src's content into it
	if err := idtools.MkdirAllAndChownNew(dst, 0755, rootIDs); err != nil ***REMOVED***
		return err
	***REMOVED***
	logrus.Debugf("Calling TarUntar(%s, %s)", src, dst)
	return archiver.TarUntar(src, dst)
***REMOVED***

// CopyFileWithTar emulates the behavior of the 'cp' command-line
// for a single file. It copies a regular file from path `src` to
// path `dst`, and preserves all its metadata.
func (archiver *Archiver) CopyFileWithTar(src, dst string) (err error) ***REMOVED***
	logrus.Debugf("CopyFileWithTar(%s, %s)", src, dst)
	srcDriver := archiver.SrcDriver
	dstDriver := archiver.DstDriver

	srcSt, err := srcDriver.Stat(src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if srcSt.IsDir() ***REMOVED***
		return fmt.Errorf("Can't copy a directory")
	***REMOVED***

	// Clean up the trailing slash. This must be done in an operating
	// system specific manner.
	if dst[len(dst)-1] == dstDriver.Separator() ***REMOVED***
		dst = dstDriver.Join(dst, srcDriver.Base(src))
	***REMOVED***

	// The original call was system.MkdirAll, which is just
	// os.MkdirAll on not-Windows and changed for Windows.
	if dstDriver.OS() == "windows" ***REMOVED***
		// Now we are WCOW
		if err := system.MkdirAll(filepath.Dir(dst), 0700, ""); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// We can just use the driver.MkdirAll function
		if err := dstDriver.MkdirAll(dstDriver.Dir(dst), 0700); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	r, w := io.Pipe()
	errC := make(chan error, 1)

	go func() ***REMOVED***
		defer close(errC)
		errC <- func() error ***REMOVED***
			defer w.Close()

			srcF, err := srcDriver.Open(src)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			defer srcF.Close()

			hdr, err := tar.FileInfoHeader(srcSt, "")
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hdr.Name = dstDriver.Base(dst)
			if dstDriver.OS() == "windows" ***REMOVED***
				hdr.Mode = int64(chmodTarEntry(os.FileMode(hdr.Mode)))
			***REMOVED*** else ***REMOVED***
				hdr.Mode = int64(os.FileMode(hdr.Mode))
			***REMOVED***

			if err := remapIDs(archiver.IDMappingsVar, hdr); err != nil ***REMOVED***
				return err
			***REMOVED***

			tw := tar.NewWriter(w)
			defer tw.Close()
			if err := tw.WriteHeader(hdr); err != nil ***REMOVED***
				return err
			***REMOVED***
			if _, err := io.Copy(tw, srcF); err != nil ***REMOVED***
				return err
			***REMOVED***
			return nil
		***REMOVED***()
	***REMOVED***()
	defer func() ***REMOVED***
		if er := <-errC; err == nil && er != nil ***REMOVED***
			err = er
		***REMOVED***
	***REMOVED***()

	err = archiver.Untar(r, dstDriver.Dir(dst), nil)
	if err != nil ***REMOVED***
		r.CloseWithError(err)
	***REMOVED***
	return err
***REMOVED***

// IDMappings returns the IDMappings of the archiver.
func (archiver *Archiver) IDMappings() *idtools.IDMappings ***REMOVED***
	return archiver.IDMappingsVar
***REMOVED***

func remapIDs(idMappings *idtools.IDMappings, hdr *tar.Header) error ***REMOVED***
	ids, err := idMappings.ToHost(idtools.IDPair***REMOVED***UID: hdr.Uid, GID: hdr.Gid***REMOVED***)
	hdr.Uid, hdr.Gid = ids.UID, ids.GID
	return err
***REMOVED***

// chmodTarEntry is used to adjust the file permissions used in tar header based
// on the platform the archival is done.
func chmodTarEntry(perm os.FileMode) os.FileMode ***REMOVED***
	//perm &= 0755 // this 0-ed out tar flags (like link, regular file, directory marker etc.)
	permPart := perm & os.ModePerm
	noPermPart := perm &^ os.ModePerm
	// Add the x bit: make everything +x from windows
	permPart |= 0111
	permPart &= 0755

	return noPermPart | permPart
***REMOVED***
