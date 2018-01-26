// Copyright (c) 2013, Suryandaru Triandana <syndtr@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package capability

import (
	"syscall"
	"unsafe"
)

type capHeader struct ***REMOVED***
	version uint32
	pid     int
***REMOVED***

type capData struct ***REMOVED***
	effective   uint32
	permitted   uint32
	inheritable uint32
***REMOVED***

func capget(hdr *capHeader, data *capData) (err error) ***REMOVED***
	_, _, e1 := syscall.Syscall(syscall.SYS_CAPGET, uintptr(unsafe.Pointer(hdr)), uintptr(unsafe.Pointer(data)), 0)
	if e1 != 0 ***REMOVED***
		err = e1
	***REMOVED***
	return
***REMOVED***

func capset(hdr *capHeader, data *capData) (err error) ***REMOVED***
	_, _, e1 := syscall.Syscall(syscall.SYS_CAPSET, uintptr(unsafe.Pointer(hdr)), uintptr(unsafe.Pointer(data)), 0)
	if e1 != 0 ***REMOVED***
		err = e1
	***REMOVED***
	return
***REMOVED***

func prctl(option int, arg2, arg3, arg4, arg5 uintptr) (err error) ***REMOVED***
	_, _, e1 := syscall.Syscall6(syscall.SYS_PRCTL, uintptr(option), arg2, arg3, arg4, arg5, 0)
	if e1 != 0 ***REMOVED***
		err = e1
	***REMOVED***
	return
***REMOVED***

const (
	vfsXattrName = "security.capability"

	vfsCapVerMask = 0xff000000
	vfsCapVer1    = 0x01000000
	vfsCapVer2    = 0x02000000

	vfsCapFlagMask      = ^vfsCapVerMask
	vfsCapFlageffective = 0x000001

	vfscapDataSizeV1 = 4 * (1 + 2*1)
	vfscapDataSizeV2 = 4 * (1 + 2*2)
)

type vfscapData struct ***REMOVED***
	magic uint32
	data  [2]struct ***REMOVED***
		permitted   uint32
		inheritable uint32
	***REMOVED***
	effective [2]uint32
	version   int8
***REMOVED***

var (
	_vfsXattrName *byte
)

func init() ***REMOVED***
	_vfsXattrName, _ = syscall.BytePtrFromString(vfsXattrName)
***REMOVED***

func getVfsCap(path string, dest *vfscapData) (err error) ***REMOVED***
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(path)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	r0, _, e1 := syscall.Syscall6(syscall.SYS_GETXATTR, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_vfsXattrName)), uintptr(unsafe.Pointer(dest)), vfscapDataSizeV2, 0, 0)
	if e1 != 0 ***REMOVED***
		if e1 == syscall.ENODATA ***REMOVED***
			dest.version = 2
			return
		***REMOVED***
		err = e1
	***REMOVED***
	switch dest.magic & vfsCapVerMask ***REMOVED***
	case vfsCapVer1:
		dest.version = 1
		if r0 != vfscapDataSizeV1 ***REMOVED***
			return syscall.EINVAL
		***REMOVED***
		dest.data[1].permitted = 0
		dest.data[1].inheritable = 0
	case vfsCapVer2:
		dest.version = 2
		if r0 != vfscapDataSizeV2 ***REMOVED***
			return syscall.EINVAL
		***REMOVED***
	default:
		return syscall.EINVAL
	***REMOVED***
	if dest.magic&vfsCapFlageffective != 0 ***REMOVED***
		dest.effective[0] = dest.data[0].permitted | dest.data[0].inheritable
		dest.effective[1] = dest.data[1].permitted | dest.data[1].inheritable
	***REMOVED*** else ***REMOVED***
		dest.effective[0] = 0
		dest.effective[1] = 0
	***REMOVED***
	return
***REMOVED***

func setVfsCap(path string, data *vfscapData) (err error) ***REMOVED***
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(path)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	var size uintptr
	if data.version == 1 ***REMOVED***
		data.magic = vfsCapVer1
		size = vfscapDataSizeV1
	***REMOVED*** else if data.version == 2 ***REMOVED***
		data.magic = vfsCapVer2
		if data.effective[0] != 0 || data.effective[1] != 0 ***REMOVED***
			data.magic |= vfsCapFlageffective
		***REMOVED***
		size = vfscapDataSizeV2
	***REMOVED*** else ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	_, _, e1 := syscall.Syscall6(syscall.SYS_SETXATTR, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_vfsXattrName)), uintptr(unsafe.Pointer(data)), size, 0, 0)
	if e1 != 0 ***REMOVED***
		err = e1
	***REMOVED***
	return
***REMOVED***
