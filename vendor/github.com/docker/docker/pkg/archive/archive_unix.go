// +build !windows

package archive

import (
	"archive/tar"
	"errors"
	"os"
	"path/filepath"
	"syscall"

	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/system"
	rsystem "github.com/opencontainers/runc/libcontainer/system"
	"golang.org/x/sys/unix"
)

// fixVolumePathPrefix does platform specific processing to ensure that if
// the path being passed in is not in a volume path format, convert it to one.
func fixVolumePathPrefix(srcPath string) string ***REMOVED***
	return srcPath
***REMOVED***

// getWalkRoot calculates the root path when performing a TarWithOptions.
// We use a separate function as this is platform specific. On Linux, we
// can't use filepath.Join(srcPath,include) because this will clean away
// a trailing "." or "/" which may be important.
func getWalkRoot(srcPath string, include string) string ***REMOVED***
	return srcPath + string(filepath.Separator) + include
***REMOVED***

// CanonicalTarNameForPath returns platform-specific filepath
// to canonical posix-style path for tar archival. p is relative
// path.
func CanonicalTarNameForPath(p string) (string, error) ***REMOVED***
	return p, nil // already unix-style
***REMOVED***

// chmodTarEntry is used to adjust the file permissions used in tar header based
// on the platform the archival is done.

func chmodTarEntry(perm os.FileMode) os.FileMode ***REMOVED***
	return perm // noop for unix as golang APIs provide perm bits correctly
***REMOVED***

func setHeaderForSpecialDevice(hdr *tar.Header, name string, stat interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	s, ok := stat.(*syscall.Stat_t)

	if ok ***REMOVED***
		// Currently go does not fill in the major/minors
		if s.Mode&unix.S_IFBLK != 0 ||
			s.Mode&unix.S_IFCHR != 0 ***REMOVED***
			hdr.Devmajor = int64(unix.Major(uint64(s.Rdev))) // nolint: unconvert
			hdr.Devminor = int64(unix.Minor(uint64(s.Rdev))) // nolint: unconvert
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func getInodeFromStat(stat interface***REMOVED******REMOVED***) (inode uint64, err error) ***REMOVED***
	s, ok := stat.(*syscall.Stat_t)

	if ok ***REMOVED***
		inode = s.Ino
	***REMOVED***

	return
***REMOVED***

func getFileUIDGID(stat interface***REMOVED******REMOVED***) (idtools.IDPair, error) ***REMOVED***
	s, ok := stat.(*syscall.Stat_t)

	if !ok ***REMOVED***
		return idtools.IDPair***REMOVED******REMOVED***, errors.New("cannot convert stat value to syscall.Stat_t")
	***REMOVED***
	return idtools.IDPair***REMOVED***UID: int(s.Uid), GID: int(s.Gid)***REMOVED***, nil
***REMOVED***

// handleTarTypeBlockCharFifo is an OS-specific helper function used by
// createTarFile to handle the following types of header: Block; Char; Fifo
func handleTarTypeBlockCharFifo(hdr *tar.Header, path string) error ***REMOVED***
	if rsystem.RunningInUserNS() ***REMOVED***
		// cannot create a device if running in user namespace
		return nil
	***REMOVED***

	mode := uint32(hdr.Mode & 07777)
	switch hdr.Typeflag ***REMOVED***
	case tar.TypeBlock:
		mode |= unix.S_IFBLK
	case tar.TypeChar:
		mode |= unix.S_IFCHR
	case tar.TypeFifo:
		mode |= unix.S_IFIFO
	***REMOVED***

	return system.Mknod(path, mode, int(system.Mkdev(hdr.Devmajor, hdr.Devminor)))
***REMOVED***

func handleLChmod(hdr *tar.Header, path string, hdrInfo os.FileInfo) error ***REMOVED***
	if hdr.Typeflag == tar.TypeLink ***REMOVED***
		if fi, err := os.Lstat(hdr.Linkname); err == nil && (fi.Mode()&os.ModeSymlink == 0) ***REMOVED***
			if err := os.Chmod(path, hdrInfo.Mode()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if hdr.Typeflag != tar.TypeSymlink ***REMOVED***
		if err := os.Chmod(path, hdrInfo.Mode()); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
