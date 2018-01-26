// +build linux

package selinux

import (
	"syscall"
	"unsafe"
)

var _zero uintptr

// Returns a []byte slice if the xattr is set and nil otherwise
// Requires path and its attribute as arguments
func lgetxattr(path string, attr string) ([]byte, error) ***REMOVED***
	var sz int
	pathBytes, err := syscall.BytePtrFromString(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	attrBytes, err := syscall.BytePtrFromString(attr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Start with a 128 length byte array
	sz = 128
	dest := make([]byte, sz)
	destBytes := unsafe.Pointer(&dest[0])
	_sz, _, errno := syscall.Syscall6(syscall.SYS_LGETXATTR, uintptr(unsafe.Pointer(pathBytes)), uintptr(unsafe.Pointer(attrBytes)), uintptr(destBytes), uintptr(len(dest)), 0, 0)

	switch ***REMOVED***
	case errno == syscall.ENODATA:
		return nil, errno
	case errno == syscall.ENOTSUP:
		return nil, errno
	case errno == syscall.ERANGE:
		// 128 byte array might just not be good enough,
		// A dummy buffer is used ``uintptr(0)`` to get real size
		// of the xattrs on disk
		_sz, _, errno = syscall.Syscall6(syscall.SYS_LGETXATTR, uintptr(unsafe.Pointer(pathBytes)), uintptr(unsafe.Pointer(attrBytes)), uintptr(unsafe.Pointer(nil)), uintptr(0), 0, 0)
		sz = int(_sz)
		if sz < 0 ***REMOVED***
			return nil, errno
		***REMOVED***
		dest = make([]byte, sz)
		destBytes := unsafe.Pointer(&dest[0])
		_sz, _, errno = syscall.Syscall6(syscall.SYS_LGETXATTR, uintptr(unsafe.Pointer(pathBytes)), uintptr(unsafe.Pointer(attrBytes)), uintptr(destBytes), uintptr(len(dest)), 0, 0)
		if errno != 0 ***REMOVED***
			return nil, errno
		***REMOVED***
	case errno != 0:
		return nil, errno
	***REMOVED***
	sz = int(_sz)
	return dest[:sz], nil
***REMOVED***

func lsetxattr(path string, attr string, data []byte, flags int) error ***REMOVED***
	pathBytes, err := syscall.BytePtrFromString(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	attrBytes, err := syscall.BytePtrFromString(attr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	var dataBytes unsafe.Pointer
	if len(data) > 0 ***REMOVED***
		dataBytes = unsafe.Pointer(&data[0])
	***REMOVED*** else ***REMOVED***
		dataBytes = unsafe.Pointer(&_zero)
	***REMOVED***
	_, _, errno := syscall.Syscall6(syscall.SYS_LSETXATTR, uintptr(unsafe.Pointer(pathBytes)), uintptr(unsafe.Pointer(attrBytes)), uintptr(dataBytes), uintptr(len(data)), uintptr(flags), 0)
	if errno != 0 ***REMOVED***
		return errno
	***REMOVED***
	return nil
***REMOVED***
