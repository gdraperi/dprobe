// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tar implements access to tar archives.
// It aims to cover most of the variations, including those produced
// by GNU and BSD tars.
//
// References:
//   http://www.freebsd.org/cgi/man.cgi?query=tar&sektion=5
//   http://www.gnu.org/software/tar/manual/html_node/Standard.html
//   http://pubs.opengroup.org/onlinepubs/9699919799/utilities/pax.html
package tar

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"
)

// BUG: Use of the Uid and Gid fields in Header could overflow on 32-bit
// architectures. If a large value is encountered when decoding, the result
// stored in Header will be the truncated version.

// Header type flags.
const (
	TypeReg           = '0'    // regular file
	TypeRegA          = '\x00' // regular file
	TypeLink          = '1'    // hard link
	TypeSymlink       = '2'    // symbolic link
	TypeChar          = '3'    // character device node
	TypeBlock         = '4'    // block device node
	TypeDir           = '5'    // directory
	TypeFifo          = '6'    // fifo node
	TypeCont          = '7'    // reserved
	TypeXHeader       = 'x'    // extended header
	TypeXGlobalHeader = 'g'    // global extended header
	TypeGNULongName   = 'L'    // Next file has a long name
	TypeGNULongLink   = 'K'    // Next file symlinks to a file w/ a long name
	TypeGNUSparse     = 'S'    // sparse file
)

// A Header represents a single header in a tar archive.
// Some fields may not be populated.
type Header struct ***REMOVED***
	Name       string    // name of header file entry
	Mode       int64     // permission and mode bits
	Uid        int       // user id of owner
	Gid        int       // group id of owner
	Size       int64     // length in bytes
	ModTime    time.Time // modified time
	Typeflag   byte      // type of header entry
	Linkname   string    // target name of link
	Uname      string    // user name of owner
	Gname      string    // group name of owner
	Devmajor   int64     // major number of character or block device
	Devminor   int64     // minor number of character or block device
	AccessTime time.Time // access time
	ChangeTime time.Time // status change time
	Xattrs     map[string]string
***REMOVED***

// FileInfo returns an os.FileInfo for the Header.
func (h *Header) FileInfo() os.FileInfo ***REMOVED***
	return headerFileInfo***REMOVED***h***REMOVED***
***REMOVED***

// headerFileInfo implements os.FileInfo.
type headerFileInfo struct ***REMOVED***
	h *Header
***REMOVED***

func (fi headerFileInfo) Size() int64        ***REMOVED*** return fi.h.Size ***REMOVED***
func (fi headerFileInfo) IsDir() bool        ***REMOVED*** return fi.Mode().IsDir() ***REMOVED***
func (fi headerFileInfo) ModTime() time.Time ***REMOVED*** return fi.h.ModTime ***REMOVED***
func (fi headerFileInfo) Sys() interface***REMOVED******REMOVED***   ***REMOVED*** return fi.h ***REMOVED***

// Name returns the base name of the file.
func (fi headerFileInfo) Name() string ***REMOVED***
	if fi.IsDir() ***REMOVED***
		return path.Base(path.Clean(fi.h.Name))
	***REMOVED***
	return path.Base(fi.h.Name)
***REMOVED***

// Mode returns the permission and mode bits for the headerFileInfo.
func (fi headerFileInfo) Mode() (mode os.FileMode) ***REMOVED***
	// Set file permission bits.
	mode = os.FileMode(fi.h.Mode).Perm()

	// Set setuid, setgid and sticky bits.
	if fi.h.Mode&c_ISUID != 0 ***REMOVED***
		// setuid
		mode |= os.ModeSetuid
	***REMOVED***
	if fi.h.Mode&c_ISGID != 0 ***REMOVED***
		// setgid
		mode |= os.ModeSetgid
	***REMOVED***
	if fi.h.Mode&c_ISVTX != 0 ***REMOVED***
		// sticky
		mode |= os.ModeSticky
	***REMOVED***

	// Set file mode bits.
	// clear perm, setuid, setgid and sticky bits.
	m := os.FileMode(fi.h.Mode) &^ 07777
	if m == c_ISDIR ***REMOVED***
		// directory
		mode |= os.ModeDir
	***REMOVED***
	if m == c_ISFIFO ***REMOVED***
		// named pipe (FIFO)
		mode |= os.ModeNamedPipe
	***REMOVED***
	if m == c_ISLNK ***REMOVED***
		// symbolic link
		mode |= os.ModeSymlink
	***REMOVED***
	if m == c_ISBLK ***REMOVED***
		// device file
		mode |= os.ModeDevice
	***REMOVED***
	if m == c_ISCHR ***REMOVED***
		// Unix character device
		mode |= os.ModeDevice
		mode |= os.ModeCharDevice
	***REMOVED***
	if m == c_ISSOCK ***REMOVED***
		// Unix domain socket
		mode |= os.ModeSocket
	***REMOVED***

	switch fi.h.Typeflag ***REMOVED***
	case TypeSymlink:
		// symbolic link
		mode |= os.ModeSymlink
	case TypeChar:
		// character device node
		mode |= os.ModeDevice
		mode |= os.ModeCharDevice
	case TypeBlock:
		// block device node
		mode |= os.ModeDevice
	case TypeDir:
		// directory
		mode |= os.ModeDir
	case TypeFifo:
		// fifo node
		mode |= os.ModeNamedPipe
	***REMOVED***

	return mode
***REMOVED***

// sysStat, if non-nil, populates h from system-dependent fields of fi.
var sysStat func(fi os.FileInfo, h *Header) error

const (
	// Mode constants from the USTAR spec:
	// See http://pubs.opengroup.org/onlinepubs/9699919799/utilities/pax.html#tag_20_92_13_06
	c_ISUID = 04000 // Set uid
	c_ISGID = 02000 // Set gid
	c_ISVTX = 01000 // Save text (sticky bit)

	// Common Unix mode constants; these are not defined in any common tar standard.
	// Header.FileInfo understands these, but FileInfoHeader will never produce these.
	c_ISDIR  = 040000  // Directory
	c_ISFIFO = 010000  // FIFO
	c_ISREG  = 0100000 // Regular file
	c_ISLNK  = 0120000 // Symbolic link
	c_ISBLK  = 060000  // Block special file
	c_ISCHR  = 020000  // Character special file
	c_ISSOCK = 0140000 // Socket
)

// Keywords for the PAX Extended Header
const (
	paxAtime    = "atime"
	paxCharset  = "charset"
	paxComment  = "comment"
	paxCtime    = "ctime" // please note that ctime is not a valid pax header.
	paxGid      = "gid"
	paxGname    = "gname"
	paxLinkpath = "linkpath"
	paxMtime    = "mtime"
	paxPath     = "path"
	paxSize     = "size"
	paxUid      = "uid"
	paxUname    = "uname"
	paxXattr    = "SCHILY.xattr."
	paxNone     = ""
)

// FileInfoHeader creates a partially-populated Header from fi.
// If fi describes a symlink, FileInfoHeader records link as the link target.
// If fi describes a directory, a slash is appended to the name.
// Because os.FileInfo's Name method returns only the base name of
// the file it describes, it may be necessary to modify the Name field
// of the returned header to provide the full path name of the file.
func FileInfoHeader(fi os.FileInfo, link string) (*Header, error) ***REMOVED***
	if fi == nil ***REMOVED***
		return nil, errors.New("tar: FileInfo is nil")
	***REMOVED***
	fm := fi.Mode()
	h := &Header***REMOVED***
		Name:    fi.Name(),
		ModTime: fi.ModTime(),
		Mode:    int64(fm.Perm()), // or'd with c_IS* constants later
	***REMOVED***
	switch ***REMOVED***
	case fm.IsRegular():
		h.Typeflag = TypeReg
		h.Size = fi.Size()
	case fi.IsDir():
		h.Typeflag = TypeDir
		h.Name += "/"
	case fm&os.ModeSymlink != 0:
		h.Typeflag = TypeSymlink
		h.Linkname = link
	case fm&os.ModeDevice != 0:
		if fm&os.ModeCharDevice != 0 ***REMOVED***
			h.Typeflag = TypeChar
		***REMOVED*** else ***REMOVED***
			h.Typeflag = TypeBlock
		***REMOVED***
	case fm&os.ModeNamedPipe != 0:
		h.Typeflag = TypeFifo
	case fm&os.ModeSocket != 0:
		return nil, fmt.Errorf("archive/tar: sockets not supported")
	default:
		return nil, fmt.Errorf("archive/tar: unknown file mode %v", fm)
	***REMOVED***
	if fm&os.ModeSetuid != 0 ***REMOVED***
		h.Mode |= c_ISUID
	***REMOVED***
	if fm&os.ModeSetgid != 0 ***REMOVED***
		h.Mode |= c_ISGID
	***REMOVED***
	if fm&os.ModeSticky != 0 ***REMOVED***
		h.Mode |= c_ISVTX
	***REMOVED***
	// If possible, populate additional fields from OS-specific
	// FileInfo fields.
	if sys, ok := fi.Sys().(*Header); ok ***REMOVED***
		// This FileInfo came from a Header (not the OS). Use the
		// original Header to populate all remaining fields.
		h.Uid = sys.Uid
		h.Gid = sys.Gid
		h.Uname = sys.Uname
		h.Gname = sys.Gname
		h.AccessTime = sys.AccessTime
		h.ChangeTime = sys.ChangeTime
		if sys.Xattrs != nil ***REMOVED***
			h.Xattrs = make(map[string]string)
			for k, v := range sys.Xattrs ***REMOVED***
				h.Xattrs[k] = v
			***REMOVED***
		***REMOVED***
		if sys.Typeflag == TypeLink ***REMOVED***
			// hard link
			h.Typeflag = TypeLink
			h.Size = 0
			h.Linkname = sys.Linkname
		***REMOVED***
	***REMOVED***
	if sysStat != nil ***REMOVED***
		return h, sysStat(fi, h)
	***REMOVED***
	return h, nil
***REMOVED***

// isHeaderOnlyType checks if the given type flag is of the type that has no
// data section even if a size is specified.
func isHeaderOnlyType(flag byte) bool ***REMOVED***
	switch flag ***REMOVED***
	case TypeLink, TypeSymlink, TypeChar, TypeBlock, TypeDir, TypeFifo:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***
