package operatingsystem

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// See https://code.google.com/p/go/source/browse/src/pkg/mime/type_windows.go?r=d14520ac25bf6940785aabb71f5be453a286f58c
// for a similar sample

// GetOperatingSystem gets the name of the current operating system.
func GetOperatingSystem() (string, error) ***REMOVED***

	var h windows.Handle

	// Default return value
	ret := "Unknown Operating System"

	if err := windows.RegOpenKeyEx(windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\`),
		0,
		windows.KEY_READ,
		&h); err != nil ***REMOVED***
		return ret, err
	***REMOVED***
	defer windows.RegCloseKey(h)

	var buf [1 << 10]uint16
	var typ uint32
	n := uint32(len(buf) * 2) // api expects array of bytes, not uint16

	if err := windows.RegQueryValueEx(h,
		windows.StringToUTF16Ptr("ProductName"),
		nil,
		&typ,
		(*byte)(unsafe.Pointer(&buf[0])),
		&n); err != nil ***REMOVED***
		return ret, err
	***REMOVED***
	ret = windows.UTF16ToString(buf[:])

	return ret, nil
***REMOVED***

// IsContainerized returns true if we are running inside a container.
// No-op on Windows, always returns false.
func IsContainerized() (bool, error) ***REMOVED***
	return false, nil
***REMOVED***
