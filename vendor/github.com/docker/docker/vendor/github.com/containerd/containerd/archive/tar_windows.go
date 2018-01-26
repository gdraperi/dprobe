package archive

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/containerd/containerd/sys"
	"github.com/dmcgowan/go-tar"
)

// tarName returns platform-specific filepath
// to canonical posix-style path for tar archival. p is relative
// path.
func tarName(p string) (string, error) ***REMOVED***
	// windows: convert windows style relative path with backslashes
	// into forward slashes. Since windows does not allow '/' or '\'
	// in file names, it is mostly safe to replace however we must
	// check just in case
	if strings.Contains(p, "/") ***REMOVED***
		return "", fmt.Errorf("Windows path contains forward slash: %s", p)
	***REMOVED***

	return strings.Replace(p, string(os.PathSeparator), "/", -1), nil
***REMOVED***

// chmodTarEntry is used to adjust the file permissions used in tar header based
// on the platform the archival is done.
func chmodTarEntry(perm os.FileMode) os.FileMode ***REMOVED***
	perm &= 0755
	// Add the x bit: make everything +x from windows
	perm |= 0111

	return perm
***REMOVED***

func setHeaderForSpecialDevice(*tar.Header, string, os.FileInfo) error ***REMOVED***
	// do nothing. no notion of Rdev, Inode, Nlink in stat on Windows
	return nil
***REMOVED***

func open(p string) (*os.File, error) ***REMOVED***
	// We use sys.OpenSequential to ensure we use sequential file
	// access on Windows to avoid depleting the standby list.
	return sys.OpenSequential(p)
***REMOVED***

func openFile(name string, flag int, perm os.FileMode) (*os.File, error) ***REMOVED***
	// Source is regular file. We use sys.OpenFileSequential to use sequential
	// file access to avoid depleting the standby list on Windows.
	return sys.OpenFileSequential(name, flag, perm)
***REMOVED***

func mkdirAll(path string, perm os.FileMode) error ***REMOVED***
	return sys.MkdirAll(path, perm)
***REMOVED***

func mkdir(path string, perm os.FileMode) error ***REMOVED***
	return os.Mkdir(path, perm)
***REMOVED***

func skipFile(hdr *tar.Header) bool ***REMOVED***
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
	if strings.Contains(hdr.Name, ":") ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

// handleTarTypeBlockCharFifo is an OS-specific helper function used by
// createTarFile to handle the following types of header: Block; Char; Fifo
func handleTarTypeBlockCharFifo(hdr *tar.Header, path string) error ***REMOVED***
	return nil
***REMOVED***

func handleLChmod(hdr *tar.Header, path string, hdrInfo os.FileInfo) error ***REMOVED***
	return nil
***REMOVED***

func getxattr(path, attr string) ([]byte, error) ***REMOVED***
	return nil, nil
***REMOVED***

func setxattr(path, key, value string) error ***REMOVED***
	// Return not support error, do not wrap underlying not supported
	// since xattrs should not exist in windows diff archives
	return errors.New("xattrs not supported on Windows")
***REMOVED***
