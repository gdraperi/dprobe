package archive

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
)

// UnpackLayer unpack `layer` to a `dest`. The stream `layer` can be
// compressed or uncompressed.
// Returns the size in bytes of the contents of the layer.
func UnpackLayer(dest string, layer io.Reader, options *TarOptions) (size int64, err error) ***REMOVED***
	tr := tar.NewReader(layer)
	trBuf := pools.BufioReader32KPool.Get(tr)
	defer pools.BufioReader32KPool.Put(trBuf)

	var dirs []*tar.Header
	unpackedPaths := make(map[string]struct***REMOVED******REMOVED***)

	if options == nil ***REMOVED***
		options = &TarOptions***REMOVED******REMOVED***
	***REMOVED***
	if options.ExcludePatterns == nil ***REMOVED***
		options.ExcludePatterns = []string***REMOVED******REMOVED***
	***REMOVED***
	idMappings := idtools.NewIDMappingsFromMaps(options.UIDMaps, options.GIDMaps)

	aufsTempdir := ""
	aufsHardlinks := make(map[string]*tar.Header)

	// Iterate through the files in the archive.
	for ***REMOVED***
		hdr, err := tr.Next()
		if err == io.EOF ***REMOVED***
			// end of tar archive
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***

		size += hdr.Size

		// Normalize name, for safety and for a simple is-root check
		hdr.Name = filepath.Clean(hdr.Name)

		// Windows does not support filenames with colons in them. Ignore
		// these files. This is not a problem though (although it might
		// appear that it is). Let's suppose a client is running docker pull.
		// The daemon it points to is Windows. Would it make sense for the
		// client to be doing a docker pull Ubuntu for example (which has files
		// with colons in the name under /usr/share/man/man3)? No, absolutely
		// not as it would really only make sense that they were pulling a
		// Windows image. However, for development, it is necessary to be able
		// to pull Linux images which are in the repository.
		//
		// TODO Windows. Once the registry is aware of what images are Windows-
		// specific or Linux-specific, this warning should be changed to an error
		// to cater for the situation where someone does manage to upload a Linux
		// image but have it tagged as Windows inadvertently.
		if runtime.GOOS == "windows" ***REMOVED***
			if strings.Contains(hdr.Name, ":") ***REMOVED***
				logrus.Warnf("Windows: Ignoring %s (is this a Linux image?)", hdr.Name)
				continue
			***REMOVED***
		***REMOVED***

		// Note as these operations are platform specific, so must the slash be.
		if !strings.HasSuffix(hdr.Name, string(os.PathSeparator)) ***REMOVED***
			// Not the root directory, ensure that the parent directory exists.
			// This happened in some tests where an image had a tarfile without any
			// parent directories.
			parent := filepath.Dir(hdr.Name)
			parentPath := filepath.Join(dest, parent)

			if _, err := os.Lstat(parentPath); err != nil && os.IsNotExist(err) ***REMOVED***
				err = system.MkdirAll(parentPath, 0600, "")
				if err != nil ***REMOVED***
					return 0, err
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// Skip AUFS metadata dirs
		if strings.HasPrefix(hdr.Name, WhiteoutMetaPrefix) ***REMOVED***
			// Regular files inside /.wh..wh.plnk can be used as hardlink targets
			// We don't want this directory, but we need the files in them so that
			// such hardlinks can be resolved.
			if strings.HasPrefix(hdr.Name, WhiteoutLinkDir) && hdr.Typeflag == tar.TypeReg ***REMOVED***
				basename := filepath.Base(hdr.Name)
				aufsHardlinks[basename] = hdr
				if aufsTempdir == "" ***REMOVED***
					if aufsTempdir, err = ioutil.TempDir("", "dockerplnk"); err != nil ***REMOVED***
						return 0, err
					***REMOVED***
					defer os.RemoveAll(aufsTempdir)
				***REMOVED***
				if err := createTarFile(filepath.Join(aufsTempdir, basename), dest, hdr, tr, true, nil, options.InUserNS); err != nil ***REMOVED***
					return 0, err
				***REMOVED***
			***REMOVED***

			if hdr.Name != WhiteoutOpaqueDir ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		path := filepath.Join(dest, hdr.Name)
		rel, err := filepath.Rel(dest, path)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***

		// Note as these operations are platform specific, so must the slash be.
		if strings.HasPrefix(rel, ".."+string(os.PathSeparator)) ***REMOVED***
			return 0, breakoutError(fmt.Errorf("%q is outside of %q", hdr.Name, dest))
		***REMOVED***
		base := filepath.Base(path)

		if strings.HasPrefix(base, WhiteoutPrefix) ***REMOVED***
			dir := filepath.Dir(path)
			if base == WhiteoutOpaqueDir ***REMOVED***
				_, err := os.Lstat(dir)
				if err != nil ***REMOVED***
					return 0, err
				***REMOVED***
				err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error ***REMOVED***
					if err != nil ***REMOVED***
						if os.IsNotExist(err) ***REMOVED***
							err = nil // parent was deleted
						***REMOVED***
						return err
					***REMOVED***
					if path == dir ***REMOVED***
						return nil
					***REMOVED***
					if _, exists := unpackedPaths[path]; !exists ***REMOVED***
						err := os.RemoveAll(path)
						return err
					***REMOVED***
					return nil
				***REMOVED***)
				if err != nil ***REMOVED***
					return 0, err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				originalBase := base[len(WhiteoutPrefix):]
				originalPath := filepath.Join(dir, originalBase)
				if err := os.RemoveAll(originalPath); err != nil ***REMOVED***
					return 0, err
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// If path exits we almost always just want to remove and replace it.
			// The only exception is when it is a directory *and* the file from
			// the layer is also a directory. Then we want to merge them (i.e.
			// just apply the metadata from the layer).
			if fi, err := os.Lstat(path); err == nil ***REMOVED***
				if !(fi.IsDir() && hdr.Typeflag == tar.TypeDir) ***REMOVED***
					if err := os.RemoveAll(path); err != nil ***REMOVED***
						return 0, err
					***REMOVED***
				***REMOVED***
			***REMOVED***

			trBuf.Reset(tr)
			srcData := io.Reader(trBuf)
			srcHdr := hdr

			// Hard links into /.wh..wh.plnk don't work, as we don't extract that directory, so
			// we manually retarget these into the temporary files we extracted them into
			if hdr.Typeflag == tar.TypeLink && strings.HasPrefix(filepath.Clean(hdr.Linkname), WhiteoutLinkDir) ***REMOVED***
				linkBasename := filepath.Base(hdr.Linkname)
				srcHdr = aufsHardlinks[linkBasename]
				if srcHdr == nil ***REMOVED***
					return 0, fmt.Errorf("Invalid aufs hardlink")
				***REMOVED***
				tmpFile, err := os.Open(filepath.Join(aufsTempdir, linkBasename))
				if err != nil ***REMOVED***
					return 0, err
				***REMOVED***
				defer tmpFile.Close()
				srcData = tmpFile
			***REMOVED***

			if err := remapIDs(idMappings, srcHdr); err != nil ***REMOVED***
				return 0, err
			***REMOVED***

			if err := createTarFile(path, dest, srcHdr, srcData, true, nil, options.InUserNS); err != nil ***REMOVED***
				return 0, err
			***REMOVED***

			// Directory mtimes must be handled at the end to avoid further
			// file creation in them to modify the directory mtime
			if hdr.Typeflag == tar.TypeDir ***REMOVED***
				dirs = append(dirs, hdr)
			***REMOVED***
			unpackedPaths[path] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	for _, hdr := range dirs ***REMOVED***
		path := filepath.Join(dest, hdr.Name)
		if err := system.Chtimes(path, hdr.AccessTime, hdr.ModTime); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
	***REMOVED***

	return size, nil
***REMOVED***

// ApplyLayer parses a diff in the standard layer format from `layer`,
// and applies it to the directory `dest`. The stream `layer` can be
// compressed or uncompressed.
// Returns the size in bytes of the contents of the layer.
func ApplyLayer(dest string, layer io.Reader) (int64, error) ***REMOVED***
	return applyLayerHandler(dest, layer, &TarOptions***REMOVED******REMOVED***, true)
***REMOVED***

// ApplyUncompressedLayer parses a diff in the standard layer format from
// `layer`, and applies it to the directory `dest`. The stream `layer`
// can only be uncompressed.
// Returns the size in bytes of the contents of the layer.
func ApplyUncompressedLayer(dest string, layer io.Reader, options *TarOptions) (int64, error) ***REMOVED***
	return applyLayerHandler(dest, layer, options, false)
***REMOVED***

// do the bulk load of ApplyLayer, but allow for not calling DecompressStream
func applyLayerHandler(dest string, layer io.Reader, options *TarOptions, decompress bool) (int64, error) ***REMOVED***
	dest = filepath.Clean(dest)

	// We need to be able to set any perms
	oldmask, err := system.Umask(0)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer system.Umask(oldmask) // ignore err, ErrNotSupportedPlatform

	if decompress ***REMOVED***
		layer, err = DecompressStream(layer)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
	***REMOVED***
	return UnpackLayer(dest, layer, options)
***REMOVED***
