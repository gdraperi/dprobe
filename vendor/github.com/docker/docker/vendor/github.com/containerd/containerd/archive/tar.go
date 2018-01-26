package archive

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/containerd/containerd/fs"
	"github.com/containerd/containerd/log"
	"github.com/dmcgowan/go-tar"
	"github.com/pkg/errors"
)

var bufferPool = &sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		buffer := make([]byte, 32*1024)
		return &buffer
	***REMOVED***,
***REMOVED***

// Diff returns a tar stream of the computed filesystem
// difference between the provided directories.
//
// Produces a tar using OCI style file markers for deletions. Deleted
// files will be prepended with the prefix ".wh.". This style is
// based off AUFS whiteouts.
// See https://github.com/opencontainers/image-spec/blob/master/layer.md
func Diff(ctx context.Context, a, b string) io.ReadCloser ***REMOVED***
	r, w := io.Pipe()

	go func() ***REMOVED***
		err := WriteDiff(ctx, w, a, b)
		if err = w.CloseWithError(err); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Debugf("closing tar pipe failed")
		***REMOVED***
	***REMOVED***()

	return r
***REMOVED***

// WriteDiff writes a tar stream of the computed difference between the
// provided directories.
//
// Produces a tar using OCI style file markers for deletions. Deleted
// files will be prepended with the prefix ".wh.". This style is
// based off AUFS whiteouts.
// See https://github.com/opencontainers/image-spec/blob/master/layer.md
func WriteDiff(ctx context.Context, w io.Writer, a, b string) error ***REMOVED***
	cw := newChangeWriter(w, b)
	err := fs.Changes(ctx, a, b, cw.HandleChange)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to create diff tar stream")
	***REMOVED***
	return cw.Close()
***REMOVED***

const (
	// whiteoutPrefix prefix means file is a whiteout. If this is followed by a
	// filename this means that file has been removed from the base layer.
	// See https://github.com/opencontainers/image-spec/blob/master/layer.md#whiteouts
	whiteoutPrefix = ".wh."

	// whiteoutMetaPrefix prefix means whiteout has a special meaning and is not
	// for removing an actual file. Normally these files are excluded from exported
	// archives.
	whiteoutMetaPrefix = whiteoutPrefix + whiteoutPrefix

	// whiteoutLinkDir is a directory AUFS uses for storing hardlink links to other
	// layers. Normally these should not go into exported archives and all changed
	// hardlinks should be copied to the top layer.
	whiteoutLinkDir = whiteoutMetaPrefix + "plnk"

	// whiteoutOpaqueDir file means directory has been made opaque - meaning
	// readdir calls to this directory do not follow to lower layers.
	whiteoutOpaqueDir = whiteoutMetaPrefix + ".opq"

	paxSchilyXattr = "SCHILY.xattrs."
)

// Apply applies a tar stream of an OCI style diff tar.
// See https://github.com/opencontainers/image-spec/blob/master/layer.md#applying-changesets
func Apply(ctx context.Context, root string, r io.Reader) (int64, error) ***REMOVED***
	root = filepath.Clean(root)

	var (
		tr   = tar.NewReader(r)
		size int64
		dirs []*tar.Header

		// Used for handling opaque directory markers which
		// may occur out of order
		unpackedPaths = make(map[string]struct***REMOVED******REMOVED***)

		// Used for aufs plink directory
		aufsTempdir   = ""
		aufsHardlinks = make(map[string]*tar.Header)
	)

	// Iterate through the files in the archive.
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		***REMOVED***

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

		if skipFile(hdr) ***REMOVED***
			log.G(ctx).Warnf("file %q ignored: archive may not be supported on system", hdr.Name)
			continue
		***REMOVED***

		// Split name and resolve symlinks for root directory.
		ppath, base := filepath.Split(hdr.Name)
		ppath, err = fs.RootPath(root, ppath)
		if err != nil ***REMOVED***
			return 0, errors.Wrap(err, "failed to get root path")
		***REMOVED***

		// Join to root before joining to parent path to ensure relative links are
		// already resolved based on the root before adding to parent.
		path := filepath.Join(ppath, filepath.Join("/", base))
		if path == root ***REMOVED***
			log.G(ctx).Debugf("file %q ignored: resolved to root", hdr.Name)
			continue
		***REMOVED***

		// If file is not directly under root, ensure parent directory
		// exists or is created.
		if ppath != root ***REMOVED***
			parentPath := ppath
			if base == "" ***REMOVED***
				parentPath = filepath.Dir(path)
			***REMOVED***
			if _, err := os.Lstat(parentPath); err != nil && os.IsNotExist(err) ***REMOVED***
				err = mkdirAll(parentPath, 0700)
				if err != nil ***REMOVED***
					return 0, err
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// Skip AUFS metadata dirs
		if strings.HasPrefix(hdr.Name, whiteoutMetaPrefix) ***REMOVED***
			// Regular files inside /.wh..wh.plnk can be used as hardlink targets
			// We don't want this directory, but we need the files in them so that
			// such hardlinks can be resolved.
			if strings.HasPrefix(hdr.Name, whiteoutLinkDir) && hdr.Typeflag == tar.TypeReg ***REMOVED***
				basename := filepath.Base(hdr.Name)
				aufsHardlinks[basename] = hdr
				if aufsTempdir == "" ***REMOVED***
					if aufsTempdir, err = ioutil.TempDir("", "dockerplnk"); err != nil ***REMOVED***
						return 0, err
					***REMOVED***
					defer os.RemoveAll(aufsTempdir)
				***REMOVED***
				p, err := fs.RootPath(aufsTempdir, basename)
				if err != nil ***REMOVED***
					return 0, err
				***REMOVED***
				if err := createTarFile(ctx, p, root, hdr, tr); err != nil ***REMOVED***
					return 0, err
				***REMOVED***
			***REMOVED***

			if hdr.Name != whiteoutOpaqueDir ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		if strings.HasPrefix(base, whiteoutPrefix) ***REMOVED***
			dir := filepath.Dir(path)
			if base == whiteoutOpaqueDir ***REMOVED***
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
				continue
			***REMOVED***

			originalBase := base[len(whiteoutPrefix):]
			originalPath := filepath.Join(dir, originalBase)
			if err := os.RemoveAll(originalPath); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			continue
		***REMOVED***
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

		srcData := io.Reader(tr)
		srcHdr := hdr

		// Hard links into /.wh..wh.plnk don't work, as we don't extract that directory, so
		// we manually retarget these into the temporary files we extracted them into
		if hdr.Typeflag == tar.TypeLink && strings.HasPrefix(filepath.Clean(hdr.Linkname), whiteoutLinkDir) ***REMOVED***
			linkBasename := filepath.Base(hdr.Linkname)
			srcHdr = aufsHardlinks[linkBasename]
			if srcHdr == nil ***REMOVED***
				return 0, fmt.Errorf("Invalid aufs hardlink")
			***REMOVED***
			p, err := fs.RootPath(aufsTempdir, linkBasename)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			tmpFile, err := os.Open(p)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			defer tmpFile.Close()
			srcData = tmpFile
		***REMOVED***

		if err := createTarFile(ctx, path, root, srcHdr, srcData); err != nil ***REMOVED***
			return 0, err
		***REMOVED***

		// Directory mtimes must be handled at the end to avoid further
		// file creation in them to modify the directory mtime
		if hdr.Typeflag == tar.TypeDir ***REMOVED***
			dirs = append(dirs, hdr)
		***REMOVED***
		unpackedPaths[path] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	for _, hdr := range dirs ***REMOVED***
		path, err := fs.RootPath(root, hdr.Name)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		if err := chtimes(path, boundTime(latestTime(hdr.AccessTime, hdr.ModTime)), boundTime(hdr.ModTime)); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
	***REMOVED***

	return size, nil
***REMOVED***

type changeWriter struct ***REMOVED***
	tw        *tar.Writer
	source    string
	whiteoutT time.Time
	inodeSrc  map[uint64]string
	inodeRefs map[uint64][]string
***REMOVED***

func newChangeWriter(w io.Writer, source string) *changeWriter ***REMOVED***
	return &changeWriter***REMOVED***
		tw:        tar.NewWriter(w),
		source:    source,
		whiteoutT: time.Now(),
		inodeSrc:  map[uint64]string***REMOVED******REMOVED***,
		inodeRefs: map[uint64][]string***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

func (cw *changeWriter) HandleChange(k fs.ChangeKind, p string, f os.FileInfo, err error) error ***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if k == fs.ChangeKindDelete ***REMOVED***
		whiteOutDir := filepath.Dir(p)
		whiteOutBase := filepath.Base(p)
		whiteOut := filepath.Join(whiteOutDir, whiteoutPrefix+whiteOutBase)
		hdr := &tar.Header***REMOVED***
			Name:       whiteOut[1:],
			Size:       0,
			ModTime:    cw.whiteoutT,
			AccessTime: cw.whiteoutT,
			ChangeTime: cw.whiteoutT,
		***REMOVED***
		if err := cw.tw.WriteHeader(hdr); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to write whiteout header")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var (
			link   string
			err    error
			source = filepath.Join(cw.source, p)
		)

		if f.Mode()&os.ModeSymlink != 0 ***REMOVED***
			if link, err = os.Readlink(source); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		hdr, err := tar.FileInfoHeader(f, link)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		hdr.Mode = int64(chmodTarEntry(os.FileMode(hdr.Mode)))

		name := p
		if strings.HasPrefix(name, string(filepath.Separator)) ***REMOVED***
			name, err = filepath.Rel(string(filepath.Separator), name)
			if err != nil ***REMOVED***
				return errors.Wrap(err, "failed to make path relative")
			***REMOVED***
		***REMOVED***
		name, err = tarName(name)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "cannot canonicalize path")
		***REMOVED***
		// suffix with '/' for directories
		if f.IsDir() && !strings.HasSuffix(name, "/") ***REMOVED***
			name += "/"
		***REMOVED***
		hdr.Name = name

		if err := setHeaderForSpecialDevice(hdr, name, f); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to set device headers")
		***REMOVED***

		// additionalLinks stores file names which must be linked to
		// this file when this file is added
		var additionalLinks []string
		inode, isHardlink := fs.GetLinkInfo(f)
		if isHardlink ***REMOVED***
			// If the inode has a source, always link to it
			if source, ok := cw.inodeSrc[inode]; ok ***REMOVED***
				hdr.Typeflag = tar.TypeLink
				hdr.Linkname = source
				hdr.Size = 0
			***REMOVED*** else ***REMOVED***
				if k == fs.ChangeKindUnmodified ***REMOVED***
					cw.inodeRefs[inode] = append(cw.inodeRefs[inode], name)
					return nil
				***REMOVED***
				cw.inodeSrc[inode] = name
				additionalLinks = cw.inodeRefs[inode]
				delete(cw.inodeRefs, inode)
			***REMOVED***
		***REMOVED*** else if k == fs.ChangeKindUnmodified && !f.IsDir() ***REMOVED***
			// Nothing to write to diff
			// Unmodified directories should still be written to keep
			// directory permissions correct on direct unpack
			return nil
		***REMOVED***

		if capability, err := getxattr(source, "security.capability"); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to get capabilities xattr")
		***REMOVED*** else if capability != nil ***REMOVED***
			if hdr.PAXRecords == nil ***REMOVED***
				hdr.PAXRecords = map[string]string***REMOVED******REMOVED***
			***REMOVED***
			hdr.PAXRecords[paxSchilyXattr+"security.capability"] = string(capability)
		***REMOVED***

		if err := cw.tw.WriteHeader(hdr); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to write file header")
		***REMOVED***

		if hdr.Typeflag == tar.TypeReg && hdr.Size > 0 ***REMOVED***
			file, err := open(source)
			if err != nil ***REMOVED***
				return errors.Wrapf(err, "failed to open path: %v", source)
			***REMOVED***
			defer file.Close()

			buf := bufferPool.Get().(*[]byte)
			n, err := io.CopyBuffer(cw.tw, file, *buf)
			bufferPool.Put(buf)
			if err != nil ***REMOVED***
				return errors.Wrap(err, "failed to copy")
			***REMOVED***
			if n != hdr.Size ***REMOVED***
				return errors.New("short write copying file")
			***REMOVED***
		***REMOVED***

		if additionalLinks != nil ***REMOVED***
			source = hdr.Name
			for _, extra := range additionalLinks ***REMOVED***
				hdr.Name = extra
				hdr.Typeflag = tar.TypeLink
				hdr.Linkname = source
				hdr.Size = 0
				if err := cw.tw.WriteHeader(hdr); err != nil ***REMOVED***
					return errors.Wrap(err, "failed to write file header")
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (cw *changeWriter) Close() error ***REMOVED***
	if err := cw.tw.Close(); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to close tar writer")
	***REMOVED***
	return nil
***REMOVED***

func createTarFile(ctx context.Context, path, extractDir string, hdr *tar.Header, reader io.Reader) error ***REMOVED***
	// hdr.Mode is in linux format, which we can use for syscalls,
	// but for os.Foo() calls we need the mode converted to os.FileMode,
	// so use hdrInfo.Mode() (they differ for e.g. setuid bits)
	hdrInfo := hdr.FileInfo()

	switch hdr.Typeflag ***REMOVED***
	case tar.TypeDir:
		// Create directory unless it exists as a directory already.
		// In that case we just want to merge the two
		if fi, err := os.Lstat(path); !(err == nil && fi.IsDir()) ***REMOVED***
			if err := mkdir(path, hdrInfo.Mode()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

	case tar.TypeReg, tar.TypeRegA:
		file, err := openFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, hdrInfo.Mode())
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		_, err = copyBuffered(ctx, file, reader)
		if err1 := file.Close(); err == nil ***REMOVED***
			err = err1
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

	case tar.TypeBlock, tar.TypeChar:
		// Handle this is an OS-specific way
		if err := handleTarTypeBlockCharFifo(hdr, path); err != nil ***REMOVED***
			return err
		***REMOVED***

	case tar.TypeFifo:
		// Handle this is an OS-specific way
		if err := handleTarTypeBlockCharFifo(hdr, path); err != nil ***REMOVED***
			return err
		***REMOVED***

	case tar.TypeLink:
		targetPath, err := fs.RootPath(extractDir, hdr.Linkname)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := os.Link(targetPath, path); err != nil ***REMOVED***
			return err
		***REMOVED***

	case tar.TypeSymlink:
		if err := os.Symlink(hdr.Linkname, path); err != nil ***REMOVED***
			return err
		***REMOVED***

	case tar.TypeXGlobalHeader:
		log.G(ctx).Debug("PAX Global Extended Headers found and ignored")
		return nil

	default:
		return errors.Errorf("unhandled tar header type %d\n", hdr.Typeflag)
	***REMOVED***

	// Lchown is not supported on Windows.
	if runtime.GOOS != "windows" ***REMOVED***
		if err := os.Lchown(path, hdr.Uid, hdr.Gid); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	for key, value := range hdr.PAXRecords ***REMOVED***
		if strings.HasPrefix(key, paxSchilyXattr) ***REMOVED***
			key = key[len(paxSchilyXattr):]
			if err := setxattr(path, key, value); err != nil ***REMOVED***
				if errors.Cause(err) == syscall.ENOTSUP ***REMOVED***
					log.G(ctx).WithError(err).Warnf("ignored xattr %s in archive", key)
					continue
				***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// There is no LChmod, so ignore mode for symlink. Also, this
	// must happen after chown, as that can modify the file mode
	if err := handleLChmod(hdr, path, hdrInfo); err != nil ***REMOVED***
		return err
	***REMOVED***

	return chtimes(path, boundTime(latestTime(hdr.AccessTime, hdr.ModTime)), boundTime(hdr.ModTime))
***REMOVED***

func copyBuffered(ctx context.Context, dst io.Writer, src io.Reader) (written int64, err error) ***REMOVED***
	buf := bufferPool.Get().(*[]byte)
	defer bufferPool.Put(buf)

	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			err = ctx.Err()
			return
		default:
		***REMOVED***

		nr, er := src.Read(*buf)
		if nr > 0 ***REMOVED***
			nw, ew := dst.Write((*buf)[0:nr])
			if nw > 0 ***REMOVED***
				written += int64(nw)
			***REMOVED***
			if ew != nil ***REMOVED***
				err = ew
				break
			***REMOVED***
			if nr != nw ***REMOVED***
				err = io.ErrShortWrite
				break
			***REMOVED***
		***REMOVED***
		if er != nil ***REMOVED***
			if er != io.EOF ***REMOVED***
				err = er
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return written, err

***REMOVED***
