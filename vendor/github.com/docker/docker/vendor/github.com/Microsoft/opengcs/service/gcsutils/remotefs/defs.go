package remotefs

import (
	"errors"
	"os"
	"time"
)

// RemotefsCmd is the name of the remotefs meta command
const RemotefsCmd = "remotefs"

// Name of the commands when called from the cli context (remotefs <CMD> ...)
const (
	StatCmd           = "stat"
	LstatCmd          = "lstat"
	ReadlinkCmd       = "readlink"
	MkdirCmd          = "mkdir"
	MkdirAllCmd       = "mkdirall"
	RemoveCmd         = "remove"
	RemoveAllCmd      = "removeall"
	LinkCmd           = "link"
	SymlinkCmd        = "symlink"
	LchmodCmd         = "lchmod"
	LchownCmd         = "lchown"
	MknodCmd          = "mknod"
	MkfifoCmd         = "mkfifo"
	OpenFileCmd       = "openfile"
	ReadFileCmd       = "readfile"
	WriteFileCmd      = "writefile"
	ReadDirCmd        = "readdir"
	ResolvePathCmd    = "resolvepath"
	ExtractArchiveCmd = "extractarchive"
	ArchivePathCmd    = "archivepath"
)

// ErrInvalid is returned if the parameters are invalid
var ErrInvalid = errors.New("invalid arguments")

// ErrUnknown is returned for an unknown remotefs command
var ErrUnknown = errors.New("unkown command")

// ExportedError is the serialized version of the a Go error.
// It also provides a trivial implementation of the error interface.
type ExportedError struct ***REMOVED***
	ErrString string
	ErrNum    int `json:",omitempty"`
***REMOVED***

// Error returns an error string
func (ee *ExportedError) Error() string ***REMOVED***
	return ee.ErrString
***REMOVED***

// FileInfo is the stat struct returned by the remotefs system. It
// fulfills the os.FileInfo interface.
type FileInfo struct ***REMOVED***
	NameVar    string
	SizeVar    int64
	ModeVar    os.FileMode
	ModTimeVar int64 // Serialization of time.Time breaks in travis, so use an int
	IsDirVar   bool
***REMOVED***

var _ os.FileInfo = &FileInfo***REMOVED******REMOVED***

// Name returns the filename from a FileInfo structure
func (f *FileInfo) Name() string ***REMOVED*** return f.NameVar ***REMOVED***

// Size returns the size from a FileInfo structure
func (f *FileInfo) Size() int64 ***REMOVED*** return f.SizeVar ***REMOVED***

// Mode returns the mode from a FileInfo structure
func (f *FileInfo) Mode() os.FileMode ***REMOVED*** return f.ModeVar ***REMOVED***

// ModTime returns the modification time from a FileInfo structure
func (f *FileInfo) ModTime() time.Time ***REMOVED*** return time.Unix(0, f.ModTimeVar) ***REMOVED***

// IsDir returns the is-directory indicator from a FileInfo structure
func (f *FileInfo) IsDir() bool ***REMOVED*** return f.IsDirVar ***REMOVED***

// Sys provides an interface to a FileInfo structure
func (f *FileInfo) Sys() interface***REMOVED******REMOVED*** ***REMOVED*** return nil ***REMOVED***

// FileHeader is a header for remote *os.File operations for remotefs.OpenFile
type FileHeader struct ***REMOVED***
	Cmd  uint32
	Size uint64
***REMOVED***

const (
	// Read request command.
	Read uint32 = iota
	// Write request command.
	Write
	// Seek request command.
	Seek
	// Close request command.
	Close
	// CmdOK is a response meaning request succeeded.
	CmdOK
	// CmdFailed is a response meaning request failed.
	CmdFailed
)

// SeekHeader is header for the Seek operation for remotefs.OpenFile
type SeekHeader struct ***REMOVED***
	Offset int64
	Whence int32
***REMOVED***
