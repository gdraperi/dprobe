// +build linux,cgo

package devicemapper

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func ioctlBlkGetSize64(fd uintptr) (int64, error) ***REMOVED***
	var size int64
	if _, _, err := unix.Syscall(unix.SYS_IOCTL, fd, BlkGetSize64, uintptr(unsafe.Pointer(&size))); err != 0 ***REMOVED***
		return 0, err
	***REMOVED***
	return size, nil
***REMOVED***

func ioctlBlkDiscard(fd uintptr, offset, length uint64) error ***REMOVED***
	var r [2]uint64
	r[0] = offset
	r[1] = length

	if _, _, err := unix.Syscall(unix.SYS_IOCTL, fd, BlkDiscard, uintptr(unsafe.Pointer(&r[0]))); err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
