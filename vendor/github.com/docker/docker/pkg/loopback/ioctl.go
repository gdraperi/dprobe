// +build linux,cgo

package loopback

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func ioctlLoopCtlGetFree(fd uintptr) (int, error) ***REMOVED***
	index, err := unix.IoctlGetInt(int(fd), LoopCtlGetFree)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return index, nil
***REMOVED***

func ioctlLoopSetFd(loopFd, sparseFd uintptr) error ***REMOVED***
	return unix.IoctlSetInt(int(loopFd), LoopSetFd, int(sparseFd))
***REMOVED***

func ioctlLoopSetStatus64(loopFd uintptr, loopInfo *loopInfo64) error ***REMOVED***
	if _, _, err := unix.Syscall(unix.SYS_IOCTL, loopFd, LoopSetStatus64, uintptr(unsafe.Pointer(loopInfo))); err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func ioctlLoopClrFd(loopFd uintptr) error ***REMOVED***
	if _, _, err := unix.Syscall(unix.SYS_IOCTL, loopFd, LoopClrFd, 0); err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func ioctlLoopGetStatus64(loopFd uintptr) (*loopInfo64, error) ***REMOVED***
	loopInfo := &loopInfo64***REMOVED******REMOVED***

	if _, _, err := unix.Syscall(unix.SYS_IOCTL, loopFd, LoopGetStatus64, uintptr(unsafe.Pointer(loopInfo))); err != 0 ***REMOVED***
		return nil, err
	***REMOVED***
	return loopInfo, nil
***REMOVED***

func ioctlLoopSetCapacity(loopFd uintptr, value int) error ***REMOVED***
	return unix.IoctlSetInt(int(loopFd), LoopSetCapacity, value)
***REMOVED***
