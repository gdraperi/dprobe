// +build windows

package winio

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

//sys getFileInformationByHandleEx(h syscall.Handle, class uint32, buffer *byte, size uint32) (err error) = GetFileInformationByHandleEx
//sys setFileInformationByHandle(h syscall.Handle, class uint32, buffer *byte, size uint32) (err error) = SetFileInformationByHandle

const (
	fileBasicInfo = 0
	fileIDInfo    = 0x12
)

// FileBasicInfo contains file access time and file attributes information.
type FileBasicInfo struct ***REMOVED***
	CreationTime, LastAccessTime, LastWriteTime, ChangeTime syscall.Filetime
	FileAttributes                                          uintptr // includes padding
***REMOVED***

// GetFileBasicInfo retrieves times and attributes for a file.
func GetFileBasicInfo(f *os.File) (*FileBasicInfo, error) ***REMOVED***
	bi := &FileBasicInfo***REMOVED******REMOVED***
	if err := getFileInformationByHandleEx(syscall.Handle(f.Fd()), fileBasicInfo, (*byte)(unsafe.Pointer(bi)), uint32(unsafe.Sizeof(*bi))); err != nil ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "GetFileInformationByHandleEx", Path: f.Name(), Err: err***REMOVED***
	***REMOVED***
	runtime.KeepAlive(f)
	return bi, nil
***REMOVED***

// SetFileBasicInfo sets times and attributes for a file.
func SetFileBasicInfo(f *os.File, bi *FileBasicInfo) error ***REMOVED***
	if err := setFileInformationByHandle(syscall.Handle(f.Fd()), fileBasicInfo, (*byte)(unsafe.Pointer(bi)), uint32(unsafe.Sizeof(*bi))); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "SetFileInformationByHandle", Path: f.Name(), Err: err***REMOVED***
	***REMOVED***
	runtime.KeepAlive(f)
	return nil
***REMOVED***

// FileIDInfo contains the volume serial number and file ID for a file. This pair should be
// unique on a system.
type FileIDInfo struct ***REMOVED***
	VolumeSerialNumber uint64
	FileID             [16]byte
***REMOVED***

// GetFileID retrieves the unique (volume, file ID) pair for a file.
func GetFileID(f *os.File) (*FileIDInfo, error) ***REMOVED***
	fileID := &FileIDInfo***REMOVED******REMOVED***
	if err := getFileInformationByHandleEx(syscall.Handle(f.Fd()), fileIDInfo, (*byte)(unsafe.Pointer(fileID)), uint32(unsafe.Sizeof(*fileID))); err != nil ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "GetFileInformationByHandleEx", Path: f.Name(), Err: err***REMOVED***
	***REMOVED***
	runtime.KeepAlive(f)
	return fileID, nil
***REMOVED***
