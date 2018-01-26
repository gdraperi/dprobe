package archive

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
)

var unpigzPath string

func init() ***REMOVED***
	if path, err := exec.LookPath("unpigz"); err != nil ***REMOVED***
		logrus.Debug("unpigz binary not found in PATH, falling back to go gzip library")
	***REMOVED*** else ***REMOVED***
		logrus.Debugf("Using unpigz binary found at path %s", path)
		unpigzPath = path
	***REMOVED***
***REMOVED***

type (
	// Compression is the state represents if compressed or not.
	Compression int
	// WhiteoutFormat is the format of whiteouts unpacked
	WhiteoutFormat int

	// TarOptions wraps the tar options.
	TarOptions struct ***REMOVED***
		IncludeFiles     []string
		ExcludePatterns  []string
		Compression      Compression
		NoLchown         bool
		UIDMaps          []idtools.IDMap
		GIDMaps          []idtools.IDMap
		ChownOpts        *idtools.IDPair
		IncludeSourceDir bool
		// WhiteoutFormat is the expected on disk format for whiteout files.
		// This format will be converted to the standard format on pack
		// and from the standard format on unpack.
		WhiteoutFormat WhiteoutFormat
		// When unpacking, specifies whether overwriting a directory with a
		// non-directory is allowed and vice versa.
		NoOverwriteDirNonDir bool
		// For each include when creating an archive, the included name will be
		// replaced with the matching name from this map.
		RebaseNames map[string]string
		InUserNS    bool
	***REMOVED***
)

// Archiver implements the Archiver interface and allows the reuse of most utility functions of
// this package with a pluggable Untar function. Also, to facilitate the passing of specific id
// mappings for untar, an Archiver can be created with maps which will then be passed to Untar operations.
type Archiver struct ***REMOVED***
	Untar         func(io.Reader, string, *TarOptions) error
	IDMappingsVar *idtools.IDMappings
***REMOVED***

// NewDefaultArchiver returns a new Archiver without any IDMappings
func NewDefaultArchiver() *Archiver ***REMOVED***
	return &Archiver***REMOVED***Untar: Untar, IDMappingsVar: &idtools.IDMappings***REMOVED******REMOVED******REMOVED***
***REMOVED***

// breakoutError is used to differentiate errors related to breaking out
// When testing archive breakout in the unit tests, this error is expected
// in order for the test to pass.
type breakoutError error

const (
	// Uncompressed represents the uncompressed.
	Uncompressed Compression = iota
	// Bzip2 is bzip2 compression algorithm.
	Bzip2
	// Gzip is gzip compression algorithm.
	Gzip
	// Xz is xz compression algorithm.
	Xz
)

const (
	// AUFSWhiteoutFormat is the default format for whiteouts
	AUFSWhiteoutFormat WhiteoutFormat = iota
	// OverlayWhiteoutFormat formats whiteout according to the overlay
	// standard.
	OverlayWhiteoutFormat
)

const (
	modeISDIR  = 040000  // Directory
	modeISFIFO = 010000  // FIFO
	modeISREG  = 0100000 // Regular file
	modeISLNK  = 0120000 // Symbolic link
	modeISBLK  = 060000  // Block special file
	modeISCHR  = 020000  // Character special file
	modeISSOCK = 0140000 // Socket
)

// IsArchivePath checks if the (possibly compressed) file at the given path
// starts with a tar file header.
func IsArchivePath(path string) bool ***REMOVED***
	file, err := os.Open(path)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	defer file.Close()
	rdr, err := DecompressStream(file)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	r := tar.NewReader(rdr)
	_, err = r.Next()
	return err == nil
***REMOVED***

// DetectCompression detects the compression algorithm of the source.
func DetectCompression(source []byte) Compression ***REMOVED***
	for compression, m := range map[Compression][]byte***REMOVED***
		Bzip2: ***REMOVED***0x42, 0x5A, 0x68***REMOVED***,
		Gzip:  ***REMOVED***0x1F, 0x8B, 0x08***REMOVED***,
		Xz:    ***REMOVED***0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00***REMOVED***,
	***REMOVED*** ***REMOVED***
		if len(source) < len(m) ***REMOVED***
			logrus.Debug("Len too short")
			continue
		***REMOVED***
		if bytes.Equal(m, source[:len(m)]) ***REMOVED***
			return compression
		***REMOVED***
	***REMOVED***
	return Uncompressed
***REMOVED***

func xzDecompress(ctx context.Context, archive io.Reader) (io.ReadCloser, error) ***REMOVED***
	args := []string***REMOVED***"xz", "-d", "-c", "-q"***REMOVED***

	return cmdStream(exec.CommandContext(ctx, args[0], args[1:]...), archive)
***REMOVED***

func gzDecompress(ctx context.Context, buf io.Reader) (io.ReadCloser, error) ***REMOVED***
	if unpigzPath == "" ***REMOVED***
		return gzip.NewReader(buf)
	***REMOVED***

	disablePigzEnv := os.Getenv("MOBY_DISABLE_PIGZ")
	if disablePigzEnv != "" ***REMOVED***
		if disablePigz, err := strconv.ParseBool(disablePigzEnv); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else if disablePigz ***REMOVED***
			return gzip.NewReader(buf)
		***REMOVED***
	***REMOVED***

	return cmdStream(exec.CommandContext(ctx, unpigzPath, "-d", "-c"), buf)
***REMOVED***

func wrapReadCloser(readBuf io.ReadCloser, cancel context.CancelFunc) io.ReadCloser ***REMOVED***
	return ioutils.NewReadCloserWrapper(readBuf, func() error ***REMOVED***
		cancel()
		return readBuf.Close()
	***REMOVED***)
***REMOVED***

// DecompressStream decompresses the archive and returns a ReaderCloser with the decompressed archive.
func DecompressStream(archive io.Reader) (io.ReadCloser, error) ***REMOVED***
	p := pools.BufioReader32KPool
	buf := p.Get(archive)
	bs, err := buf.Peek(10)
	if err != nil && err != io.EOF ***REMOVED***
		// Note: we'll ignore any io.EOF error because there are some odd
		// cases where the layer.tar file will be empty (zero bytes) and
		// that results in an io.EOF from the Peek() call. So, in those
		// cases we'll just treat it as a non-compressed stream and
		// that means just create an empty layer.
		// See Issue 18170
		return nil, err
	***REMOVED***

	compression := DetectCompression(bs)
	switch compression ***REMOVED***
	case Uncompressed:
		readBufWrapper := p.NewReadCloserWrapper(buf, buf)
		return readBufWrapper, nil
	case Gzip:
		ctx, cancel := context.WithCancel(context.Background())

		gzReader, err := gzDecompress(ctx, buf)
		if err != nil ***REMOVED***
			cancel()
			return nil, err
		***REMOVED***
		readBufWrapper := p.NewReadCloserWrapper(buf, gzReader)
		return wrapReadCloser(readBufWrapper, cancel), nil
	case Bzip2:
		bz2Reader := bzip2.NewReader(buf)
		readBufWrapper := p.NewReadCloserWrapper(buf, bz2Reader)
		return readBufWrapper, nil
	case Xz:
		ctx, cancel := context.WithCancel(context.Background())

		xzReader, err := xzDecompress(ctx, buf)
		if err != nil ***REMOVED***
			cancel()
			return nil, err
		***REMOVED***
		readBufWrapper := p.NewReadCloserWrapper(buf, xzReader)
		return wrapReadCloser(readBufWrapper, cancel), nil
	default:
		return nil, fmt.Errorf("Unsupported compression format %s", (&compression).Extension())
	***REMOVED***
***REMOVED***

// CompressStream compresses the dest with specified compression algorithm.
func CompressStream(dest io.Writer, compression Compression) (io.WriteCloser, error) ***REMOVED***
	p := pools.BufioWriter32KPool
	buf := p.Get(dest)
	switch compression ***REMOVED***
	case Uncompressed:
		writeBufWrapper := p.NewWriteCloserWrapper(buf, buf)
		return writeBufWrapper, nil
	case Gzip:
		gzWriter := gzip.NewWriter(dest)
		writeBufWrapper := p.NewWriteCloserWrapper(buf, gzWriter)
		return writeBufWrapper, nil
	case Bzip2, Xz:
		// archive/bzip2 does not support writing, and there is no xz support at all
		// However, this is not a problem as docker only currently generates gzipped tars
		return nil, fmt.Errorf("Unsupported compression format %s", (&compression).Extension())
	default:
		return nil, fmt.Errorf("Unsupported compression format %s", (&compression).Extension())
	***REMOVED***
***REMOVED***

// TarModifierFunc is a function that can be passed to ReplaceFileTarWrapper to
// modify the contents or header of an entry in the archive. If the file already
// exists in the archive the TarModifierFunc will be called with the Header and
// a reader which will return the files content. If the file does not exist both
// header and content will be nil.
type TarModifierFunc func(path string, header *tar.Header, content io.Reader) (*tar.Header, []byte, error)

// ReplaceFileTarWrapper converts inputTarStream to a new tar stream. Files in the
// tar stream are modified if they match any of the keys in mods.
func ReplaceFileTarWrapper(inputTarStream io.ReadCloser, mods map[string]TarModifierFunc) io.ReadCloser ***REMOVED***
	pipeReader, pipeWriter := io.Pipe()

	go func() ***REMOVED***
		tarReader := tar.NewReader(inputTarStream)
		tarWriter := tar.NewWriter(pipeWriter)
		defer inputTarStream.Close()
		defer tarWriter.Close()

		modify := func(name string, original *tar.Header, modifier TarModifierFunc, tarReader io.Reader) error ***REMOVED***
			header, data, err := modifier(name, original, tarReader)
			switch ***REMOVED***
			case err != nil:
				return err
			case header == nil:
				return nil
			***REMOVED***

			header.Name = name
			header.Size = int64(len(data))
			if err := tarWriter.WriteHeader(header); err != nil ***REMOVED***
				return err
			***REMOVED***
			if len(data) != 0 ***REMOVED***
				if _, err := tarWriter.Write(data); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			return nil
		***REMOVED***

		var err error
		var originalHeader *tar.Header
		for ***REMOVED***
			originalHeader, err = tarReader.Next()
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			if err != nil ***REMOVED***
				pipeWriter.CloseWithError(err)
				return
			***REMOVED***

			modifier, ok := mods[originalHeader.Name]
			if !ok ***REMOVED***
				// No modifiers for this file, copy the header and data
				if err := tarWriter.WriteHeader(originalHeader); err != nil ***REMOVED***
					pipeWriter.CloseWithError(err)
					return
				***REMOVED***
				if _, err := pools.Copy(tarWriter, tarReader); err != nil ***REMOVED***
					pipeWriter.CloseWithError(err)
					return
				***REMOVED***
				continue
			***REMOVED***
			delete(mods, originalHeader.Name)

			if err := modify(originalHeader.Name, originalHeader, modifier, tarReader); err != nil ***REMOVED***
				pipeWriter.CloseWithError(err)
				return
			***REMOVED***
		***REMOVED***

		// Apply the modifiers that haven't matched any files in the archive
		for name, modifier := range mods ***REMOVED***
			if err := modify(name, nil, modifier, nil); err != nil ***REMOVED***
				pipeWriter.CloseWithError(err)
				return
			***REMOVED***
		***REMOVED***

		pipeWriter.Close()

	***REMOVED***()
	return pipeReader
***REMOVED***

// Extension returns the extension of a file that uses the specified compression algorithm.
func (compression *Compression) Extension() string ***REMOVED***
	switch *compression ***REMOVED***
	case Uncompressed:
		return "tar"
	case Bzip2:
		return "tar.bz2"
	case Gzip:
		return "tar.gz"
	case Xz:
		return "tar.xz"
	***REMOVED***
	return ""
***REMOVED***

// FileInfoHeader creates a populated Header from fi.
// Compared to archive pkg this function fills in more information.
// Also, regardless of Go version, this function fills file type bits (e.g. hdr.Mode |= modeISDIR),
// which have been deleted since Go 1.9 archive/tar.
func FileInfoHeader(name string, fi os.FileInfo, link string) (*tar.Header, error) ***REMOVED***
	hdr, err := tar.FileInfoHeader(fi, link)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	hdr.Mode = fillGo18FileTypeBits(int64(chmodTarEntry(os.FileMode(hdr.Mode))), fi)
	name, err = canonicalTarName(name, fi.IsDir())
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("tar: cannot canonicalize path: %v", err)
	***REMOVED***
	hdr.Name = name
	if err := setHeaderForSpecialDevice(hdr, name, fi.Sys()); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return hdr, nil
***REMOVED***

// fillGo18FileTypeBits fills type bits which have been removed on Go 1.9 archive/tar
// https://github.com/golang/go/commit/66b5a2f
func fillGo18FileTypeBits(mode int64, fi os.FileInfo) int64 ***REMOVED***
	fm := fi.Mode()
	switch ***REMOVED***
	case fm.IsRegular():
		mode |= modeISREG
	case fi.IsDir():
		mode |= modeISDIR
	case fm&os.ModeSymlink != 0:
		mode |= modeISLNK
	case fm&os.ModeDevice != 0:
		if fm&os.ModeCharDevice != 0 ***REMOVED***
			mode |= modeISCHR
		***REMOVED*** else ***REMOVED***
			mode |= modeISBLK
		***REMOVED***
	case fm&os.ModeNamedPipe != 0:
		mode |= modeISFIFO
	case fm&os.ModeSocket != 0:
		mode |= modeISSOCK
	***REMOVED***
	return mode
***REMOVED***

// ReadSecurityXattrToTarHeader reads security.capability xattr from filesystem
// to a tar header
func ReadSecurityXattrToTarHeader(path string, hdr *tar.Header) error ***REMOVED***
	capability, _ := system.Lgetxattr(path, "security.capability")
	if capability != nil ***REMOVED***
		hdr.Xattrs = make(map[string]string)
		hdr.Xattrs["security.capability"] = string(capability)
	***REMOVED***
	return nil
***REMOVED***

type tarWhiteoutConverter interface ***REMOVED***
	ConvertWrite(*tar.Header, string, os.FileInfo) (*tar.Header, error)
	ConvertRead(*tar.Header, string) (bool, error)
***REMOVED***

type tarAppender struct ***REMOVED***
	TarWriter *tar.Writer
	Buffer    *bufio.Writer

	// for hardlink mapping
	SeenFiles  map[uint64]string
	IDMappings *idtools.IDMappings
	ChownOpts  *idtools.IDPair

	// For packing and unpacking whiteout files in the
	// non standard format. The whiteout files defined
	// by the AUFS standard are used as the tar whiteout
	// standard.
	WhiteoutConverter tarWhiteoutConverter
***REMOVED***

func newTarAppender(idMapping *idtools.IDMappings, writer io.Writer, chownOpts *idtools.IDPair) *tarAppender ***REMOVED***
	return &tarAppender***REMOVED***
		SeenFiles:  make(map[uint64]string),
		TarWriter:  tar.NewWriter(writer),
		Buffer:     pools.BufioWriter32KPool.Get(nil),
		IDMappings: idMapping,
		ChownOpts:  chownOpts,
	***REMOVED***
***REMOVED***

// canonicalTarName provides a platform-independent and consistent posix-style
//path for files and directories to be archived regardless of the platform.
func canonicalTarName(name string, isDir bool) (string, error) ***REMOVED***
	name, err := CanonicalTarNameForPath(name)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	// suffix with '/' for directories
	if isDir && !strings.HasSuffix(name, "/") ***REMOVED***
		name += "/"
	***REMOVED***
	return name, nil
***REMOVED***

// addTarFile adds to the tar archive a file from `path` as `name`
func (ta *tarAppender) addTarFile(path, name string) error ***REMOVED***
	fi, err := os.Lstat(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var link string
	if fi.Mode()&os.ModeSymlink != 0 ***REMOVED***
		var err error
		link, err = os.Readlink(path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	hdr, err := FileInfoHeader(name, fi, link)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := ReadSecurityXattrToTarHeader(path, hdr); err != nil ***REMOVED***
		return err
	***REMOVED***

	// if it's not a directory and has more than 1 link,
	// it's hard linked, so set the type flag accordingly
	if !fi.IsDir() && hasHardlinks(fi) ***REMOVED***
		inode, err := getInodeFromStat(fi.Sys())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// a link should have a name that it links too
		// and that linked name should be first in the tar archive
		if oldpath, ok := ta.SeenFiles[inode]; ok ***REMOVED***
			hdr.Typeflag = tar.TypeLink
			hdr.Linkname = oldpath
			hdr.Size = 0 // This Must be here for the writer math to add up!
		***REMOVED*** else ***REMOVED***
			ta.SeenFiles[inode] = name
		***REMOVED***
	***REMOVED***

	//check whether the file is overlayfs whiteout
	//if yes, skip re-mapping container ID mappings.
	isOverlayWhiteout := fi.Mode()&os.ModeCharDevice != 0 && hdr.Devmajor == 0 && hdr.Devminor == 0

	//handle re-mapping container ID mappings back to host ID mappings before
	//writing tar headers/files. We skip whiteout files because they were written
	//by the kernel and already have proper ownership relative to the host
	if !isOverlayWhiteout &&
		!strings.HasPrefix(filepath.Base(hdr.Name), WhiteoutPrefix) &&
		!ta.IDMappings.Empty() ***REMOVED***
		fileIDPair, err := getFileUIDGID(fi.Sys())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		hdr.Uid, hdr.Gid, err = ta.IDMappings.ToContainer(fileIDPair)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// explicitly override with ChownOpts
	if ta.ChownOpts != nil ***REMOVED***
		hdr.Uid = ta.ChownOpts.UID
		hdr.Gid = ta.ChownOpts.GID
	***REMOVED***

	if ta.WhiteoutConverter != nil ***REMOVED***
		wo, err := ta.WhiteoutConverter.ConvertWrite(hdr, path, fi)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// If a new whiteout file exists, write original hdr, then
		// replace hdr with wo to be written after. Whiteouts should
		// always be written after the original. Note the original
		// hdr may have been updated to be a whiteout with returning
		// a whiteout header
		if wo != nil ***REMOVED***
			if err := ta.TarWriter.WriteHeader(hdr); err != nil ***REMOVED***
				return err
			***REMOVED***
			if hdr.Typeflag == tar.TypeReg && hdr.Size > 0 ***REMOVED***
				return fmt.Errorf("tar: cannot use whiteout for non-empty file")
			***REMOVED***
			hdr = wo
		***REMOVED***
	***REMOVED***

	if err := ta.TarWriter.WriteHeader(hdr); err != nil ***REMOVED***
		return err
	***REMOVED***

	if hdr.Typeflag == tar.TypeReg && hdr.Size > 0 ***REMOVED***
		// We use system.OpenSequential to ensure we use sequential file
		// access on Windows to avoid depleting the standby list.
		// On Linux, this equates to a regular os.Open.
		file, err := system.OpenSequential(path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		ta.Buffer.Reset(ta.TarWriter)
		defer ta.Buffer.Reset(nil)
		_, err = io.Copy(ta.Buffer, file)
		file.Close()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = ta.Buffer.Flush()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func createTarFile(path, extractDir string, hdr *tar.Header, reader io.Reader, Lchown bool, chownOpts *idtools.IDPair, inUserns bool) error ***REMOVED***
	// hdr.Mode is in linux format, which we can use for sycalls,
	// but for os.Foo() calls we need the mode converted to os.FileMode,
	// so use hdrInfo.Mode() (they differ for e.g. setuid bits)
	hdrInfo := hdr.FileInfo()

	switch hdr.Typeflag ***REMOVED***
	case tar.TypeDir:
		// Create directory unless it exists as a directory already.
		// In that case we just want to merge the two
		if fi, err := os.Lstat(path); !(err == nil && fi.IsDir()) ***REMOVED***
			if err := os.Mkdir(path, hdrInfo.Mode()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

	case tar.TypeReg, tar.TypeRegA:
		// Source is regular file. We use system.OpenFileSequential to use sequential
		// file access to avoid depleting the standby list on Windows.
		// On Linux, this equates to a regular os.OpenFile
		file, err := system.OpenFileSequential(path, os.O_CREATE|os.O_WRONLY, hdrInfo.Mode())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, err := io.Copy(file, reader); err != nil ***REMOVED***
			file.Close()
			return err
		***REMOVED***
		file.Close()

	case tar.TypeBlock, tar.TypeChar:
		if inUserns ***REMOVED*** // cannot create devices in a userns
			return nil
		***REMOVED***
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
		targetPath := filepath.Join(extractDir, hdr.Linkname)
		// check for hardlink breakout
		if !strings.HasPrefix(targetPath, extractDir) ***REMOVED***
			return breakoutError(fmt.Errorf("invalid hardlink %q -> %q", targetPath, hdr.Linkname))
		***REMOVED***
		if err := os.Link(targetPath, path); err != nil ***REMOVED***
			return err
		***REMOVED***

	case tar.TypeSymlink:
		// 	path 				-> hdr.Linkname = targetPath
		// e.g. /extractDir/path/to/symlink 	-> ../2/file	= /extractDir/path/2/file
		targetPath := filepath.Join(filepath.Dir(path), hdr.Linkname)

		// the reason we don't need to check symlinks in the path (with FollowSymlinkInScope) is because
		// that symlink would first have to be created, which would be caught earlier, at this very check:
		if !strings.HasPrefix(targetPath, extractDir) ***REMOVED***
			return breakoutError(fmt.Errorf("invalid symlink %q -> %q", path, hdr.Linkname))
		***REMOVED***
		if err := os.Symlink(hdr.Linkname, path); err != nil ***REMOVED***
			return err
		***REMOVED***

	case tar.TypeXGlobalHeader:
		logrus.Debug("PAX Global Extended Headers found and ignored")
		return nil

	default:
		return fmt.Errorf("unhandled tar header type %d", hdr.Typeflag)
	***REMOVED***

	// Lchown is not supported on Windows.
	if Lchown && runtime.GOOS != "windows" ***REMOVED***
		if chownOpts == nil ***REMOVED***
			chownOpts = &idtools.IDPair***REMOVED***UID: hdr.Uid, GID: hdr.Gid***REMOVED***
		***REMOVED***
		if err := os.Lchown(path, chownOpts.UID, chownOpts.GID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	var errors []string
	for key, value := range hdr.Xattrs ***REMOVED***
		if err := system.Lsetxattr(path, key, []byte(value), 0); err != nil ***REMOVED***
			if err == syscall.ENOTSUP ***REMOVED***
				// We ignore errors here because not all graphdrivers support
				// xattrs *cough* old versions of AUFS *cough*. However only
				// ENOTSUP should be emitted in that case, otherwise we still
				// bail.
				errors = append(errors, err.Error())
				continue
			***REMOVED***
			return err
		***REMOVED***

	***REMOVED***

	if len(errors) > 0 ***REMOVED***
		logrus.WithFields(logrus.Fields***REMOVED***
			"errors": errors,
		***REMOVED***).Warn("ignored xattrs in archive: underlying filesystem doesn't support them")
	***REMOVED***

	// There is no LChmod, so ignore mode for symlink. Also, this
	// must happen after chown, as that can modify the file mode
	if err := handleLChmod(hdr, path, hdrInfo); err != nil ***REMOVED***
		return err
	***REMOVED***

	aTime := hdr.AccessTime
	if aTime.Before(hdr.ModTime) ***REMOVED***
		// Last access time should never be before last modified time.
		aTime = hdr.ModTime
	***REMOVED***

	// system.Chtimes doesn't support a NOFOLLOW flag atm
	if hdr.Typeflag == tar.TypeLink ***REMOVED***
		if fi, err := os.Lstat(hdr.Linkname); err == nil && (fi.Mode()&os.ModeSymlink == 0) ***REMOVED***
			if err := system.Chtimes(path, aTime, hdr.ModTime); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if hdr.Typeflag != tar.TypeSymlink ***REMOVED***
		if err := system.Chtimes(path, aTime, hdr.ModTime); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		ts := []syscall.Timespec***REMOVED***timeToTimespec(aTime), timeToTimespec(hdr.ModTime)***REMOVED***
		if err := system.LUtimesNano(path, ts); err != nil && err != system.ErrNotSupportedPlatform ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Tar creates an archive from the directory at `path`, and returns it as a
// stream of bytes.
func Tar(path string, compression Compression) (io.ReadCloser, error) ***REMOVED***
	return TarWithOptions(path, &TarOptions***REMOVED***Compression: compression***REMOVED***)
***REMOVED***

// TarWithOptions creates an archive from the directory at `path`, only including files whose relative
// paths are included in `options.IncludeFiles` (if non-nil) or not in `options.ExcludePatterns`.
func TarWithOptions(srcPath string, options *TarOptions) (io.ReadCloser, error) ***REMOVED***

	// Fix the source path to work with long path names. This is a no-op
	// on platforms other than Windows.
	srcPath = fixVolumePathPrefix(srcPath)

	pm, err := fileutils.NewPatternMatcher(options.ExcludePatterns)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pipeReader, pipeWriter := io.Pipe()

	compressWriter, err := CompressStream(pipeWriter, options.Compression)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	go func() ***REMOVED***
		ta := newTarAppender(
			idtools.NewIDMappingsFromMaps(options.UIDMaps, options.GIDMaps),
			compressWriter,
			options.ChownOpts,
		)
		ta.WhiteoutConverter = getWhiteoutConverter(options.WhiteoutFormat)

		defer func() ***REMOVED***
			// Make sure to check the error on Close.
			if err := ta.TarWriter.Close(); err != nil ***REMOVED***
				logrus.Errorf("Can't close tar writer: %s", err)
			***REMOVED***
			if err := compressWriter.Close(); err != nil ***REMOVED***
				logrus.Errorf("Can't close compress writer: %s", err)
			***REMOVED***
			if err := pipeWriter.Close(); err != nil ***REMOVED***
				logrus.Errorf("Can't close pipe writer: %s", err)
			***REMOVED***
		***REMOVED***()

		// this buffer is needed for the duration of this piped stream
		defer pools.BufioWriter32KPool.Put(ta.Buffer)

		// In general we log errors here but ignore them because
		// during e.g. a diff operation the container can continue
		// mutating the filesystem and we can see transient errors
		// from this

		stat, err := os.Lstat(srcPath)
		if err != nil ***REMOVED***
			return
		***REMOVED***

		if !stat.IsDir() ***REMOVED***
			// We can't later join a non-dir with any includes because the
			// 'walk' will error if "file/." is stat-ed and "file" is not a
			// directory. So, we must split the source path and use the
			// basename as the include.
			if len(options.IncludeFiles) > 0 ***REMOVED***
				logrus.Warn("Tar: Can't archive a file with includes")
			***REMOVED***

			dir, base := SplitPathDirEntry(srcPath)
			srcPath = dir
			options.IncludeFiles = []string***REMOVED***base***REMOVED***
		***REMOVED***

		if len(options.IncludeFiles) == 0 ***REMOVED***
			options.IncludeFiles = []string***REMOVED***"."***REMOVED***
		***REMOVED***

		seen := make(map[string]bool)

		for _, include := range options.IncludeFiles ***REMOVED***
			rebaseName := options.RebaseNames[include]

			walkRoot := getWalkRoot(srcPath, include)
			filepath.Walk(walkRoot, func(filePath string, f os.FileInfo, err error) error ***REMOVED***
				if err != nil ***REMOVED***
					logrus.Errorf("Tar: Can't stat file %s to tar: %s", srcPath, err)
					return nil
				***REMOVED***

				relFilePath, err := filepath.Rel(srcPath, filePath)
				if err != nil || (!options.IncludeSourceDir && relFilePath == "." && f.IsDir()) ***REMOVED***
					// Error getting relative path OR we are looking
					// at the source directory path. Skip in both situations.
					return nil
				***REMOVED***

				if options.IncludeSourceDir && include == "." && relFilePath != "." ***REMOVED***
					relFilePath = strings.Join([]string***REMOVED***".", relFilePath***REMOVED***, string(filepath.Separator))
				***REMOVED***

				skip := false

				// If "include" is an exact match for the current file
				// then even if there's an "excludePatterns" pattern that
				// matches it, don't skip it. IOW, assume an explicit 'include'
				// is asking for that file no matter what - which is true
				// for some files, like .dockerignore and Dockerfile (sometimes)
				if include != relFilePath ***REMOVED***
					skip, err = pm.Matches(relFilePath)
					if err != nil ***REMOVED***
						logrus.Errorf("Error matching %s: %v", relFilePath, err)
						return err
					***REMOVED***
				***REMOVED***

				if skip ***REMOVED***
					// If we want to skip this file and its a directory
					// then we should first check to see if there's an
					// excludes pattern (e.g. !dir/file) that starts with this
					// dir. If so then we can't skip this dir.

					// Its not a dir then so we can just return/skip.
					if !f.IsDir() ***REMOVED***
						return nil
					***REMOVED***

					// No exceptions (!...) in patterns so just skip dir
					if !pm.Exclusions() ***REMOVED***
						return filepath.SkipDir
					***REMOVED***

					dirSlash := relFilePath + string(filepath.Separator)

					for _, pat := range pm.Patterns() ***REMOVED***
						if !pat.Exclusion() ***REMOVED***
							continue
						***REMOVED***
						if strings.HasPrefix(pat.String()+string(filepath.Separator), dirSlash) ***REMOVED***
							// found a match - so can't skip this dir
							return nil
						***REMOVED***
					***REMOVED***

					// No matching exclusion dir so just skip dir
					return filepath.SkipDir
				***REMOVED***

				if seen[relFilePath] ***REMOVED***
					return nil
				***REMOVED***
				seen[relFilePath] = true

				// Rename the base resource.
				if rebaseName != "" ***REMOVED***
					var replacement string
					if rebaseName != string(filepath.Separator) ***REMOVED***
						// Special case the root directory to replace with an
						// empty string instead so that we don't end up with
						// double slashes in the paths.
						replacement = rebaseName
					***REMOVED***

					relFilePath = strings.Replace(relFilePath, include, replacement, 1)
				***REMOVED***

				if err := ta.addTarFile(filePath, relFilePath); err != nil ***REMOVED***
					logrus.Errorf("Can't add file %s to tar: %s", filePath, err)
					// if pipe is broken, stop writing tar stream to it
					if err == io.ErrClosedPipe ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				return nil
			***REMOVED***)
		***REMOVED***
	***REMOVED***()

	return pipeReader, nil
***REMOVED***

// Unpack unpacks the decompressedArchive to dest with options.
func Unpack(decompressedArchive io.Reader, dest string, options *TarOptions) error ***REMOVED***
	tr := tar.NewReader(decompressedArchive)
	trBuf := pools.BufioReader32KPool.Get(nil)
	defer pools.BufioReader32KPool.Put(trBuf)

	var dirs []*tar.Header
	idMappings := idtools.NewIDMappingsFromMaps(options.UIDMaps, options.GIDMaps)
	rootIDs := idMappings.RootPair()
	whiteoutConverter := getWhiteoutConverter(options.WhiteoutFormat)

	// Iterate through the files in the archive.
loop:
	for ***REMOVED***
		hdr, err := tr.Next()
		if err == io.EOF ***REMOVED***
			// end of tar archive
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Normalize name, for safety and for a simple is-root check
		// This keeps "../" as-is, but normalizes "/../" to "/". Or Windows:
		// This keeps "..\" as-is, but normalizes "\..\" to "\".
		hdr.Name = filepath.Clean(hdr.Name)

		for _, exclude := range options.ExcludePatterns ***REMOVED***
			if strings.HasPrefix(hdr.Name, exclude) ***REMOVED***
				continue loop
			***REMOVED***
		***REMOVED***

		// After calling filepath.Clean(hdr.Name) above, hdr.Name will now be in
		// the filepath format for the OS on which the daemon is running. Hence
		// the check for a slash-suffix MUST be done in an OS-agnostic way.
		if !strings.HasSuffix(hdr.Name, string(os.PathSeparator)) ***REMOVED***
			// Not the root directory, ensure that the parent directory exists
			parent := filepath.Dir(hdr.Name)
			parentPath := filepath.Join(dest, parent)
			if _, err := os.Lstat(parentPath); err != nil && os.IsNotExist(err) ***REMOVED***
				err = idtools.MkdirAllAndChownNew(parentPath, 0777, rootIDs)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***

		path := filepath.Join(dest, hdr.Name)
		rel, err := filepath.Rel(dest, path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if strings.HasPrefix(rel, ".."+string(os.PathSeparator)) ***REMOVED***
			return breakoutError(fmt.Errorf("%q is outside of %q", hdr.Name, dest))
		***REMOVED***

		// If path exits we almost always just want to remove and replace it
		// The only exception is when it is a directory *and* the file from
		// the layer is also a directory. Then we want to merge them (i.e.
		// just apply the metadata from the layer).
		if fi, err := os.Lstat(path); err == nil ***REMOVED***
			if options.NoOverwriteDirNonDir && fi.IsDir() && hdr.Typeflag != tar.TypeDir ***REMOVED***
				// If NoOverwriteDirNonDir is true then we cannot replace
				// an existing directory with a non-directory from the archive.
				return fmt.Errorf("cannot overwrite directory %q with non-directory %q", path, dest)
			***REMOVED***

			if options.NoOverwriteDirNonDir && !fi.IsDir() && hdr.Typeflag == tar.TypeDir ***REMOVED***
				// If NoOverwriteDirNonDir is true then we cannot replace
				// an existing non-directory with a directory from the archive.
				return fmt.Errorf("cannot overwrite non-directory %q with directory %q", path, dest)
			***REMOVED***

			if fi.IsDir() && hdr.Name == "." ***REMOVED***
				continue
			***REMOVED***

			if !(fi.IsDir() && hdr.Typeflag == tar.TypeDir) ***REMOVED***
				if err := os.RemoveAll(path); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
		trBuf.Reset(tr)

		if err := remapIDs(idMappings, hdr); err != nil ***REMOVED***
			return err
		***REMOVED***

		if whiteoutConverter != nil ***REMOVED***
			writeFile, err := whiteoutConverter.ConvertRead(hdr, path)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if !writeFile ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		if err := createTarFile(path, dest, hdr, trBuf, !options.NoLchown, options.ChownOpts, options.InUserNS); err != nil ***REMOVED***
			return err
		***REMOVED***

		// Directory mtimes must be handled at the end to avoid further
		// file creation in them to modify the directory mtime
		if hdr.Typeflag == tar.TypeDir ***REMOVED***
			dirs = append(dirs, hdr)
		***REMOVED***
	***REMOVED***

	for _, hdr := range dirs ***REMOVED***
		path := filepath.Join(dest, hdr.Name)

		if err := system.Chtimes(path, hdr.AccessTime, hdr.ModTime); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Untar reads a stream of bytes from `archive`, parses it as a tar archive,
// and unpacks it into the directory at `dest`.
// The archive may be compressed with one of the following algorithms:
//  identity (uncompressed), gzip, bzip2, xz.
// FIXME: specify behavior when target path exists vs. doesn't exist.
func Untar(tarArchive io.Reader, dest string, options *TarOptions) error ***REMOVED***
	return untarHandler(tarArchive, dest, options, true)
***REMOVED***

// UntarUncompressed reads a stream of bytes from `archive`, parses it as a tar archive,
// and unpacks it into the directory at `dest`.
// The archive must be an uncompressed stream.
func UntarUncompressed(tarArchive io.Reader, dest string, options *TarOptions) error ***REMOVED***
	return untarHandler(tarArchive, dest, options, false)
***REMOVED***

// Handler for teasing out the automatic decompression
func untarHandler(tarArchive io.Reader, dest string, options *TarOptions, decompress bool) error ***REMOVED***
	if tarArchive == nil ***REMOVED***
		return fmt.Errorf("Empty archive")
	***REMOVED***
	dest = filepath.Clean(dest)
	if options == nil ***REMOVED***
		options = &TarOptions***REMOVED******REMOVED***
	***REMOVED***
	if options.ExcludePatterns == nil ***REMOVED***
		options.ExcludePatterns = []string***REMOVED******REMOVED***
	***REMOVED***

	r := tarArchive
	if decompress ***REMOVED***
		decompressedArchive, err := DecompressStream(tarArchive)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer decompressedArchive.Close()
		r = decompressedArchive
	***REMOVED***

	return Unpack(r, dest, options)
***REMOVED***

// TarUntar is a convenience function which calls Tar and Untar, with the output of one piped into the other.
// If either Tar or Untar fails, TarUntar aborts and returns the error.
func (archiver *Archiver) TarUntar(src, dst string) error ***REMOVED***
	logrus.Debugf("TarUntar(%s %s)", src, dst)
	archive, err := TarWithOptions(src, &TarOptions***REMOVED***Compression: Uncompressed***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer archive.Close()
	options := &TarOptions***REMOVED***
		UIDMaps: archiver.IDMappingsVar.UIDs(),
		GIDMaps: archiver.IDMappingsVar.GIDs(),
	***REMOVED***
	return archiver.Untar(archive, dst, options)
***REMOVED***

// UntarPath untar a file from path to a destination, src is the source tar file path.
func (archiver *Archiver) UntarPath(src, dst string) error ***REMOVED***
	archive, err := os.Open(src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer archive.Close()
	options := &TarOptions***REMOVED***
		UIDMaps: archiver.IDMappingsVar.UIDs(),
		GIDMaps: archiver.IDMappingsVar.GIDs(),
	***REMOVED***
	return archiver.Untar(archive, dst, options)
***REMOVED***

// CopyWithTar creates a tar archive of filesystem path `src`, and
// unpacks it at filesystem path `dst`.
// The archive is streamed directly with fixed buffering and no
// intermediary disk IO.
func (archiver *Archiver) CopyWithTar(src, dst string) error ***REMOVED***
	srcSt, err := os.Stat(src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !srcSt.IsDir() ***REMOVED***
		return archiver.CopyFileWithTar(src, dst)
	***REMOVED***

	// if this Archiver is set up with ID mapping we need to create
	// the new destination directory with the remapped root UID/GID pair
	// as owner
	rootIDs := archiver.IDMappingsVar.RootPair()
	// Create dst, copy src's content into it
	logrus.Debugf("Creating dest directory: %s", dst)
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
	srcSt, err := os.Stat(src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if srcSt.IsDir() ***REMOVED***
		return fmt.Errorf("Can't copy a directory")
	***REMOVED***

	// Clean up the trailing slash. This must be done in an operating
	// system specific manner.
	if dst[len(dst)-1] == os.PathSeparator ***REMOVED***
		dst = filepath.Join(dst, filepath.Base(src))
	***REMOVED***
	// Create the holding directory if necessary
	if err := system.MkdirAll(filepath.Dir(dst), 0700, ""); err != nil ***REMOVED***
		return err
	***REMOVED***

	r, w := io.Pipe()
	errC := make(chan error, 1)

	go func() ***REMOVED***
		defer close(errC)

		errC <- func() error ***REMOVED***
			defer w.Close()

			srcF, err := os.Open(src)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			defer srcF.Close()

			hdr, err := tar.FileInfoHeader(srcSt, "")
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hdr.Name = filepath.Base(dst)
			hdr.Mode = int64(chmodTarEntry(os.FileMode(hdr.Mode)))

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

	err = archiver.Untar(r, filepath.Dir(dst), nil)
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

// cmdStream executes a command, and returns its stdout as a stream.
// If the command fails to run or doesn't complete successfully, an error
// will be returned, including anything written on stderr.
func cmdStream(cmd *exec.Cmd, input io.Reader) (io.ReadCloser, error) ***REMOVED***
	cmd.Stdin = input
	pipeR, pipeW := io.Pipe()
	cmd.Stdout = pipeW
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	// Run the command and return the pipe
	if err := cmd.Start(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Copy stdout to the returned pipe
	go func() ***REMOVED***
		if err := cmd.Wait(); err != nil ***REMOVED***
			pipeW.CloseWithError(fmt.Errorf("%s: %s", err, errBuf.String()))
		***REMOVED*** else ***REMOVED***
			pipeW.Close()
		***REMOVED***
	***REMOVED***()

	return pipeR, nil
***REMOVED***

// NewTempArchive reads the content of src into a temporary file, and returns the contents
// of that file as an archive. The archive can only be read once - as soon as reading completes,
// the file will be deleted.
func NewTempArchive(src io.Reader, dir string) (*TempArchive, error) ***REMOVED***
	f, err := ioutil.TempFile(dir, "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, err := io.Copy(f, src); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, err := f.Seek(0, 0); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	st, err := f.Stat()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	size := st.Size()
	return &TempArchive***REMOVED***File: f, Size: size***REMOVED***, nil
***REMOVED***

// TempArchive is a temporary archive. The archive can only be read once - as soon as reading completes,
// the file will be deleted.
type TempArchive struct ***REMOVED***
	*os.File
	Size   int64 // Pre-computed from Stat().Size() as a convenience
	read   int64
	closed bool
***REMOVED***

// Close closes the underlying file if it's still open, or does a no-op
// to allow callers to try to close the TempArchive multiple times safely.
func (archive *TempArchive) Close() error ***REMOVED***
	if archive.closed ***REMOVED***
		return nil
	***REMOVED***

	archive.closed = true

	return archive.File.Close()
***REMOVED***

func (archive *TempArchive) Read(data []byte) (int, error) ***REMOVED***
	n, err := archive.File.Read(data)
	archive.read += int64(n)
	if err != nil || archive.read == archive.Size ***REMOVED***
		archive.Close()
		os.Remove(archive.File.Name())
	***REMOVED***
	return n, err
***REMOVED***
