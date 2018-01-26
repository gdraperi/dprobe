package driver

import (
	"fmt"
	"io"
	"os"
)

var ErrNotSupported = fmt.Errorf("not supported")

// Driver provides all of the system-level functions in a common interface.
// The context should call these with full paths and should never use the `os`
// package or any other package to access resources on the filesystem. This
// mechanism let's us carefully control access to the context and maintain
// path and resource integrity. It also gives us an interface to reason about
// direct resource access.
//
// Implementations don't need to do much other than meet the interface. For
// example, it is not required to wrap os.FileInfo to return correct paths for
// the call to Name().
type Driver interface ***REMOVED***
	// Note that Open() returns a File interface instead of *os.File. This
	// is because os.File is a struct, so if Open was to return *os.File,
	// the only way to fulfill the interface would be to call os.Open()
	Open(path string) (File, error)
	OpenFile(path string, flag int, perm os.FileMode) (File, error)

	Stat(path string) (os.FileInfo, error)
	Lstat(path string) (os.FileInfo, error)
	Readlink(p string) (string, error)
	Mkdir(path string, mode os.FileMode) error
	Remove(path string) error

	Link(oldname, newname string) error
	Lchmod(path string, mode os.FileMode) error
	Lchown(path string, uid, gid int64) error
	Symlink(oldname, newname string) error

	MkdirAll(path string, perm os.FileMode) error
	RemoveAll(path string) error

	// TODO(aaronl): These methods might move outside the main Driver
	// interface in the future as more platforms are added.
	Mknod(path string, mode os.FileMode, major int, minor int) error
	Mkfifo(path string, mode os.FileMode) error
***REMOVED***

// File is the interface for interacting with files returned by continuity's Open
// This is needed since os.File is a struct, instead of an interface, so it can't
// be used.
type File interface ***REMOVED***
	io.ReadWriteCloser
	io.Seeker
	Readdir(n int) ([]os.FileInfo, error)
***REMOVED***

func NewSystemDriver() (Driver, error) ***REMOVED***
	// TODO(stevvooe): Consider having this take a "hint" path argument, which
	// would be the context root. The hint could be used to resolve required
	// filesystem support when assembling the driver to use.
	return &driver***REMOVED******REMOVED***, nil
***REMOVED***

// XAttrDriver should be implemented on operation systems and filesystems that
// have xattr support for regular files and directories.
type XAttrDriver interface ***REMOVED***
	// Getxattr returns all of the extended attributes for the file at path.
	// Typically, this takes a syscall call to Listxattr and Getxattr.
	Getxattr(path string) (map[string][]byte, error)

	// Setxattr sets all of the extended attributes on file at path, following
	// any symbolic links, if necessary. All attributes on the target are
	// replaced by the values from attr. If the operation fails to set any
	// attribute, those already applied will not be rolled back.
	Setxattr(path string, attr map[string][]byte) error
***REMOVED***

// LXAttrDriver should be implemented by drivers on operating systems and
// filesystems that support setting and getting extended attributes on
// symbolic links. If this is not implemented, extended attributes will be
// ignored on symbolic links.
type LXAttrDriver interface ***REMOVED***
	// LGetxattr returns all of the extended attributes for the file at path
	// and does not follow symlinks. Typically, this takes a syscall call to
	// Llistxattr and Lgetxattr.
	LGetxattr(path string) (map[string][]byte, error)

	// LSetxattr sets all of the extended attributes on file at path, without
	// following symbolic links. All attributes on the target are replaced by
	// the values from attr. If the operation fails to set any attribute,
	// those already applied will not be rolled back.
	LSetxattr(path string, attr map[string][]byte) error
***REMOVED***

type DeviceInfoDriver interface ***REMOVED***
	DeviceInfo(fi os.FileInfo) (maj uint64, min uint64, err error)
***REMOVED***

// driver is a simple default implementation that sends calls out to the "os"
// package. Extend the "driver" type in system-specific files to add support,
// such as xattrs, which can add support at compile time.
type driver struct***REMOVED******REMOVED***

var _ File = &os.File***REMOVED******REMOVED***

// LocalDriver is the exported Driver struct for convenience.
var LocalDriver Driver = &driver***REMOVED******REMOVED***

func (d *driver) Open(p string) (File, error) ***REMOVED***
	return os.Open(p)
***REMOVED***

func (d *driver) OpenFile(path string, flag int, perm os.FileMode) (File, error) ***REMOVED***
	return os.OpenFile(path, flag, perm)
***REMOVED***

func (d *driver) Stat(p string) (os.FileInfo, error) ***REMOVED***
	return os.Stat(p)
***REMOVED***

func (d *driver) Lstat(p string) (os.FileInfo, error) ***REMOVED***
	return os.Lstat(p)
***REMOVED***

func (d *driver) Readlink(p string) (string, error) ***REMOVED***
	return os.Readlink(p)
***REMOVED***

func (d *driver) Mkdir(p string, mode os.FileMode) error ***REMOVED***
	return os.Mkdir(p, mode)
***REMOVED***

// Remove is used to unlink files and remove directories.
// This is following the golang os package api which
// combines the operations into a higher level Remove
// function. If explicit unlinking or directory removal
// to mirror system call is required, they should be
// split up at that time.
func (d *driver) Remove(path string) error ***REMOVED***
	return os.Remove(path)
***REMOVED***

func (d *driver) Link(oldname, newname string) error ***REMOVED***
	return os.Link(oldname, newname)
***REMOVED***

func (d *driver) Lchown(name string, uid, gid int64) error ***REMOVED***
	// TODO: error out if uid excesses int bit width?
	return os.Lchown(name, int(uid), int(gid))
***REMOVED***

func (d *driver) Symlink(oldname, newname string) error ***REMOVED***
	return os.Symlink(oldname, newname)
***REMOVED***

func (d *driver) MkdirAll(path string, perm os.FileMode) error ***REMOVED***
	return os.MkdirAll(path, perm)
***REMOVED***

func (d *driver) RemoveAll(path string) error ***REMOVED***
	return os.RemoveAll(path)
***REMOVED***
