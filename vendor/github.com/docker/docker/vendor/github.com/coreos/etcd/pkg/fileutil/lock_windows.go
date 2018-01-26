// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build windows

package fileutil

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	modkernel32    = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx = modkernel32.NewProc("LockFileEx")

	errLocked = errors.New("The process cannot access the file because another process has locked a portion of the file.")
)

const (
	// https://msdn.microsoft.com/en-us/library/windows/desktop/aa365203(v=vs.85).aspx
	LOCKFILE_EXCLUSIVE_LOCK   = 2
	LOCKFILE_FAIL_IMMEDIATELY = 1

	// see https://msdn.microsoft.com/en-us/library/windows/desktop/ms681382(v=vs.85).aspx
	errLockViolation syscall.Errno = 0x21
)

func TryLockFile(path string, flag int, perm os.FileMode) (*LockedFile, error) ***REMOVED***
	f, err := open(path, flag, perm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := lockFile(syscall.Handle(f.Fd()), LOCKFILE_FAIL_IMMEDIATELY); err != nil ***REMOVED***
		f.Close()
		return nil, err
	***REMOVED***
	return &LockedFile***REMOVED***f***REMOVED***, nil
***REMOVED***

func LockFile(path string, flag int, perm os.FileMode) (*LockedFile, error) ***REMOVED***
	f, err := open(path, flag, perm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := lockFile(syscall.Handle(f.Fd()), 0); err != nil ***REMOVED***
		f.Close()
		return nil, err
	***REMOVED***
	return &LockedFile***REMOVED***f***REMOVED***, nil
***REMOVED***

func open(path string, flag int, perm os.FileMode) (*os.File, error) ***REMOVED***
	if path == "" ***REMOVED***
		return nil, fmt.Errorf("cannot open empty filename")
	***REMOVED***
	var access uint32
	switch flag ***REMOVED***
	case syscall.O_RDONLY:
		access = syscall.GENERIC_READ
	case syscall.O_WRONLY:
		access = syscall.GENERIC_WRITE
	case syscall.O_RDWR:
		access = syscall.GENERIC_READ | syscall.GENERIC_WRITE
	case syscall.O_WRONLY | syscall.O_CREAT:
		access = syscall.GENERIC_ALL
	default:
		panic(fmt.Errorf("flag %v is not supported", flag))
	***REMOVED***
	fd, err := syscall.CreateFile(&(syscall.StringToUTF16(path)[0]),
		access,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_ALWAYS,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return os.NewFile(uintptr(fd), path), nil
***REMOVED***

func lockFile(fd syscall.Handle, flags uint32) error ***REMOVED***
	var flag uint32 = LOCKFILE_EXCLUSIVE_LOCK
	flag |= flags
	if fd == syscall.InvalidHandle ***REMOVED***
		return nil
	***REMOVED***
	err := lockFileEx(fd, flag, 1, 0, &syscall.Overlapped***REMOVED******REMOVED***)
	if err == nil ***REMOVED***
		return nil
	***REMOVED*** else if err.Error() == errLocked.Error() ***REMOVED***
		return ErrLocked
	***REMOVED*** else if err != errLockViolation ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func lockFileEx(h syscall.Handle, flags, locklow, lockhigh uint32, ol *syscall.Overlapped) (err error) ***REMOVED***
	var reserved uint32 = 0
	r1, _, e1 := syscall.Syscall6(procLockFileEx.Addr(), 6, uintptr(h), uintptr(flags), uintptr(reserved), uintptr(locklow), uintptr(lockhigh), uintptr(unsafe.Pointer(ol)))
	if r1 == 0 ***REMOVED***
		if e1 != 0 ***REMOVED***
			err = error(e1)
		***REMOVED*** else ***REMOVED***
			err = syscall.EINVAL
		***REMOVED***
	***REMOVED***
	return
***REMOVED***
