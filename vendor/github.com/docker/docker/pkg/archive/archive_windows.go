package archive

import (
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/longpath"
)

// fixVolumePathPrefix does platform specific processing to ensure that if
// the path being passed in is not in a volume path format, convert it to one.
func fixVolumePathPrefix(srcPath string) string ***REMOVED***
	return longpath.AddPrefix(srcPath)
***REMOVED***

// getWalkRoot calculates the root path when performing a TarWithOptions.
// We use a separate function as this is platform specific.
func getWalkRoot(srcPath string, include string) string ***REMOVED***
	return filepath.Join(srcPath, include)
***REMOVED***

// CanonicalTarNameForPath returns platform-specific filepath
// to canonical posix-style path for tar archival. p is relative
// path.
func CanonicalTarNameForPath(p string) (string, error) ***REMOVED***
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
	//perm &= 0755 // this 0-ed out tar flags (like link, regular file, directory marker etc.)
	permPart := perm & os.ModePerm
	noPermPart := perm &^ os.ModePerm
	// Add the x bit: make everything +x from windows
	permPart |= 0111
	permPart &= 0755

	return noPermPart | permPart
***REMOVED***

func setHeaderForSpecialDevice(hdr *tar.Header, name string, stat interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	// do nothing. no notion of Rdev, Nlink in stat on Windows
	return
***REMOVED***

func getInodeFromStat(stat interface***REMOVED******REMOVED***) (inode uint64, err error) ***REMOVED***
	// do nothing. no notion of Inode in stat on Windows
	return
***REMOVED***

// handleTarTypeBlockCharFifo is an OS-specific helper function used by
// createTarFile to handle the following types of header: Block; Char; Fifo
func handleTarTypeBlockCharFifo(hdr *tar.Header, path string) error ***REMOVED***
	return nil
***REMOVED***

func handleLChmod(hdr *tar.Header, path string, hdrInfo os.FileInfo) error ***REMOVED***
	return nil
***REMOVED***

func getFileUIDGID(stat interface***REMOVED******REMOVED***) (idtools.IDPair, error) ***REMOVED***
	// no notion of file ownership mapping yet on Windows
	return idtools.IDPair***REMOVED***0, 0***REMOVED***, nil
***REMOVED***
