// +build windows

package sys

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"unsafe"

	winio "github.com/Microsoft/go-winio"
)

// MkdirAllWithACL is a wrapper for MkdirAll that creates a directory
// ACL'd for Builtin Administrators and Local System.
func MkdirAllWithACL(path string, perm os.FileMode) error ***REMOVED***
	return mkdirall(path, true)
***REMOVED***

// MkdirAll implementation that is volume path aware for Windows.
func MkdirAll(path string, _ os.FileMode) error ***REMOVED***
	return mkdirall(path, false)
***REMOVED***

// mkdirall is a custom version of os.MkdirAll modified for use on Windows
// so that it is both volume path aware, and can create a directory with
// a DACL.
func mkdirall(path string, adminAndLocalSystem bool) error ***REMOVED***
	if re := regexp.MustCompile(`^\\\\\?\\Volume***REMOVED***[a-z0-9-]+***REMOVED***$`); re.MatchString(path) ***REMOVED***
		return nil
	***REMOVED***

	// The rest of this method is largely copied from os.MkdirAll and should be kept
	// as-is to ensure compatibility.

	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := os.Stat(path)
	if err == nil ***REMOVED***
		if dir.IsDir() ***REMOVED***
			return nil
		***REMOVED***
		return &os.PathError***REMOVED***
			Op:   "mkdir",
			Path: path,
			Err:  syscall.ENOTDIR,
		***REMOVED***
	***REMOVED***

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) ***REMOVED*** // Skip trailing path separator.
		i--
	***REMOVED***

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) ***REMOVED*** // Scan backward over element.
		j--
	***REMOVED***

	if j > 1 ***REMOVED***
		// Create parent
		err = mkdirall(path[0:j-1], false)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Parent now exists; invoke os.Mkdir or mkdirWithACL and use its result.
	if adminAndLocalSystem ***REMOVED***
		err = mkdirWithACL(path)
	***REMOVED*** else ***REMOVED***
		err = os.Mkdir(path, 0)
	***REMOVED***

	if err != nil ***REMOVED***
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// mkdirWithACL creates a new directory. If there is an error, it will be of
// type *PathError. .
//
// This is a modified and combined version of os.Mkdir and syscall.Mkdir
// in golang to cater for creating a directory am ACL permitting full
// access, with inheritance, to any subfolder/file for Built-in Administrators
// and Local System.
func mkdirWithACL(name string) error ***REMOVED***
	sa := syscall.SecurityAttributes***REMOVED***Length: 0***REMOVED***
	sddl := "D:P(A;OICI;GA;;;BA)(A;OICI;GA;;;SY)"
	sd, err := winio.SddlToSecurityDescriptor(sddl)
	if err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "mkdir", Path: name, Err: err***REMOVED***
	***REMOVED***
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	sa.SecurityDescriptor = uintptr(unsafe.Pointer(&sd[0]))

	namep, err := syscall.UTF16PtrFromString(name)
	if err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "mkdir", Path: name, Err: err***REMOVED***
	***REMOVED***

	e := syscall.CreateDirectory(namep, &sa)
	if e != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "mkdir", Path: name, Err: e***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// IsAbs is a platform-specific wrapper for filepath.IsAbs. On Windows,
// golang filepath.IsAbs does not consider a path \windows\system32 as absolute
// as it doesn't start with a drive-letter/colon combination. However, in
// docker we need to verify things such as WORKDIR /windows/system32 in
// a Dockerfile (which gets translated to \windows\system32 when being processed
// by the daemon. This SHOULD be treated as absolute from a docker processing
// perspective.
func IsAbs(path string) bool ***REMOVED***
	if !filepath.IsAbs(path) ***REMOVED***
		if !strings.HasPrefix(path, string(os.PathSeparator)) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// The origin of the functions below here are the golang OS and syscall packages,
// slightly modified to only cope with files, not directories due to the
// specific use case.
//
// The alteration is to allow a file on Windows to be opened with
// FILE_FLAG_SEQUENTIAL_SCAN (particular for docker load), to avoid eating
// the standby list, particularly when accessing large files such as layer.tar.

// CreateSequential creates the named file with mode 0666 (before umask), truncating
// it if it already exists. If successful, methods on the returned
// File can be used for I/O; the associated file descriptor has mode
// O_RDWR.
// If there is an error, it will be of type *PathError.
func CreateSequential(name string) (*os.File, error) ***REMOVED***
	return OpenFileSequential(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0)
***REMOVED***

// OpenSequential opens the named file for reading. If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.
func OpenSequential(name string) (*os.File, error) ***REMOVED***
	return OpenFileSequential(name, os.O_RDONLY, 0)
***REMOVED***

// OpenFileSequential is the generalized open call; most users will use Open
// or Create instead.
// If there is an error, it will be of type *PathError.
func OpenFileSequential(name string, flag int, _ os.FileMode) (*os.File, error) ***REMOVED***
	if name == "" ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "open", Path: name, Err: syscall.ENOENT***REMOVED***
	***REMOVED***
	r, errf := syscallOpenFileSequential(name, flag, 0)
	if errf == nil ***REMOVED***
		return r, nil
	***REMOVED***
	return nil, &os.PathError***REMOVED***Op: "open", Path: name, Err: errf***REMOVED***
***REMOVED***

func syscallOpenFileSequential(name string, flag int, _ os.FileMode) (file *os.File, err error) ***REMOVED***
	r, e := syscallOpenSequential(name, flag|syscall.O_CLOEXEC, 0)
	if e != nil ***REMOVED***
		return nil, e
	***REMOVED***
	return os.NewFile(uintptr(r), name), nil
***REMOVED***

func makeInheritSa() *syscall.SecurityAttributes ***REMOVED***
	var sa syscall.SecurityAttributes
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	return &sa
***REMOVED***

func syscallOpenSequential(path string, mode int, _ uint32) (fd syscall.Handle, err error) ***REMOVED***
	if len(path) == 0 ***REMOVED***
		return syscall.InvalidHandle, syscall.ERROR_FILE_NOT_FOUND
	***REMOVED***
	pathp, err := syscall.UTF16PtrFromString(path)
	if err != nil ***REMOVED***
		return syscall.InvalidHandle, err
	***REMOVED***
	var access uint32
	switch mode & (syscall.O_RDONLY | syscall.O_WRONLY | syscall.O_RDWR) ***REMOVED***
	case syscall.O_RDONLY:
		access = syscall.GENERIC_READ
	case syscall.O_WRONLY:
		access = syscall.GENERIC_WRITE
	case syscall.O_RDWR:
		access = syscall.GENERIC_READ | syscall.GENERIC_WRITE
	***REMOVED***
	if mode&syscall.O_CREAT != 0 ***REMOVED***
		access |= syscall.GENERIC_WRITE
	***REMOVED***
	if mode&syscall.O_APPEND != 0 ***REMOVED***
		access &^= syscall.GENERIC_WRITE
		access |= syscall.FILE_APPEND_DATA
	***REMOVED***
	sharemode := uint32(syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE)
	var sa *syscall.SecurityAttributes
	if mode&syscall.O_CLOEXEC == 0 ***REMOVED***
		sa = makeInheritSa()
	***REMOVED***
	var createmode uint32
	switch ***REMOVED***
	case mode&(syscall.O_CREAT|syscall.O_EXCL) == (syscall.O_CREAT | syscall.O_EXCL):
		createmode = syscall.CREATE_NEW
	case mode&(syscall.O_CREAT|syscall.O_TRUNC) == (syscall.O_CREAT | syscall.O_TRUNC):
		createmode = syscall.CREATE_ALWAYS
	case mode&syscall.O_CREAT == syscall.O_CREAT:
		createmode = syscall.OPEN_ALWAYS
	case mode&syscall.O_TRUNC == syscall.O_TRUNC:
		createmode = syscall.TRUNCATE_EXISTING
	default:
		createmode = syscall.OPEN_EXISTING
	***REMOVED***
	// Use FILE_FLAG_SEQUENTIAL_SCAN rather than FILE_ATTRIBUTE_NORMAL as implemented in golang.
	//https://msdn.microsoft.com/en-us/library/windows/desktop/aa363858(v=vs.85).aspx
	const fileFlagSequentialScan = 0x08000000 // FILE_FLAG_SEQUENTIAL_SCAN
	h, e := syscall.CreateFile(pathp, access, sharemode, sa, createmode, fileFlagSequentialScan, 0)
	return h, e
***REMOVED***
