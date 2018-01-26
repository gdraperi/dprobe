package system

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procGlobalMemoryStatusEx = modkernel32.NewProc("GlobalMemoryStatusEx")
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366589(v=vs.85).aspx
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366770(v=vs.85).aspx
type memorystatusex struct ***REMOVED***
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
***REMOVED***

// ReadMemInfo retrieves memory statistics of the host system and returns a
//  MemInfo type.
func ReadMemInfo() (*MemInfo, error) ***REMOVED***
	msi := &memorystatusex***REMOVED***
		dwLength: 64,
	***REMOVED***
	r1, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(msi)))
	if r1 == 0 ***REMOVED***
		return &MemInfo***REMOVED******REMOVED***, nil
	***REMOVED***
	return &MemInfo***REMOVED***
		MemTotal:  int64(msi.ullTotalPhys),
		MemFree:   int64(msi.ullAvailPhys),
		SwapTotal: int64(msi.ullTotalPageFile),
		SwapFree:  int64(msi.ullAvailPageFile),
	***REMOVED***, nil
***REMOVED***
