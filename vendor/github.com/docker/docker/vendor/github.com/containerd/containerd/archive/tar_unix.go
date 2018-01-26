// +build !windows

package archive

import (
	"os"
	"sync"
	"syscall"

	"github.com/containerd/continuity/sysx"
	"github.com/dmcgowan/go-tar"
	"github.com/opencontainers/runc/libcontainer/system"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

func tarName(p string) (string, error) ***REMOVED***
	return p, nil
***REMOVED***

func chmodTarEntry(perm os.FileMode) os.FileMode ***REMOVED***
	return perm
***REMOVED***

func setHeaderForSpecialDevice(hdr *tar.Header, name string, fi os.FileInfo) error ***REMOVED***
	s, ok := fi.Sys().(*syscall.Stat_t)
	if !ok ***REMOVED***
		return errors.New("unsupported stat type")
	***REMOVED***

	// Currently go does not fill in the major/minors
	if s.Mode&syscall.S_IFBLK != 0 ||
		s.Mode&syscall.S_IFCHR != 0 ***REMOVED***
		hdr.Devmajor = int64(unix.Major(uint64(s.Rdev)))
		hdr.Devminor = int64(unix.Minor(uint64(s.Rdev)))
	***REMOVED***

	return nil
***REMOVED***

func open(p string) (*os.File, error) ***REMOVED***
	return os.Open(p)
***REMOVED***

func openFile(name string, flag int, perm os.FileMode) (*os.File, error) ***REMOVED***
	f, err := os.OpenFile(name, flag, perm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Call chmod to avoid permission mask
	if err := os.Chmod(name, perm); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return f, err
***REMOVED***

func mkdirAll(path string, perm os.FileMode) error ***REMOVED***
	return os.MkdirAll(path, perm)
***REMOVED***

func mkdir(path string, perm os.FileMode) error ***REMOVED***
	if err := os.Mkdir(path, perm); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Only final created directory gets explicit permission
	// call to avoid permission mask
	return os.Chmod(path, perm)
***REMOVED***

func skipFile(*tar.Header) bool ***REMOVED***
	return false
***REMOVED***

var (
	inUserNS bool
	nsOnce   sync.Once
)

func setInUserNS() ***REMOVED***
	inUserNS = system.RunningInUserNS()
***REMOVED***

// handleTarTypeBlockCharFifo is an OS-specific helper function used by
// createTarFile to handle the following types of header: Block; Char; Fifo
func handleTarTypeBlockCharFifo(hdr *tar.Header, path string) error ***REMOVED***
	nsOnce.Do(setInUserNS)
	if inUserNS ***REMOVED***
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

	return unix.Mknod(path, mode, int(unix.Mkdev(uint32(hdr.Devmajor), uint32(hdr.Devminor))))
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

func getxattr(path, attr string) ([]byte, error) ***REMOVED***
	b, err := sysx.LGetxattr(path, attr)
	if err == unix.ENOTSUP || err == sysx.ENODATA ***REMOVED***
		return nil, nil
	***REMOVED***
	return b, err
***REMOVED***

func setxattr(path, key, value string) error ***REMOVED***
	return sysx.LSetxattr(path, key, []byte(value), 0)
***REMOVED***
