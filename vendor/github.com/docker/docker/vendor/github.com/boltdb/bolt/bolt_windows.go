package bolt

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

// LockFileEx code derived from golang build filemutex_windows.go @ v1.5.1
var (
	modkernel32      = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx   = modkernel32.NewProc("LockFileEx")
	procUnlockFileEx = modkernel32.NewProc("UnlockFileEx")
)

const (
	lockExt = ".lock"

	// see https://msdn.microsoft.com/en-us/library/windows/desktop/aa365203(v=vs.85).aspx
	flagLockExclusive       = 2
	flagLockFailImmediately = 1

	// see https://msdn.microsoft.com/en-us/library/windows/desktop/ms681382(v=vs.85).aspx
	errLockViolation syscall.Errno = 0x21
)

func lockFileEx(h syscall.Handle, flags, reserved, locklow, lockhigh uint32, ol *syscall.Overlapped) (err error) ***REMOVED***
	r, _, err := procLockFileEx.Call(uintptr(h), uintptr(flags), uintptr(reserved), uintptr(locklow), uintptr(lockhigh), uintptr(unsafe.Pointer(ol)))
	if r == 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func unlockFileEx(h syscall.Handle, reserved, locklow, lockhigh uint32, ol *syscall.Overlapped) (err error) ***REMOVED***
	r, _, err := procUnlockFileEx.Call(uintptr(h), uintptr(reserved), uintptr(locklow), uintptr(lockhigh), uintptr(unsafe.Pointer(ol)), 0)
	if r == 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// fdatasync flushes written data to a file descriptor.
func fdatasync(db *DB) error ***REMOVED***
	return db.file.Sync()
***REMOVED***

// flock acquires an advisory lock on a file descriptor.
func flock(db *DB, mode os.FileMode, exclusive bool, timeout time.Duration) error ***REMOVED***
	// Create a separate lock file on windows because a process
	// cannot share an exclusive lock on the same file. This is
	// needed during Tx.WriteTo().
	f, err := os.OpenFile(db.path+lockExt, os.O_CREATE, mode)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	db.lockfile = f

	var t time.Time
	for ***REMOVED***
		// If we're beyond our timeout then return an error.
		// This can only occur after we've attempted a flock once.
		if t.IsZero() ***REMOVED***
			t = time.Now()
		***REMOVED*** else if timeout > 0 && time.Since(t) > timeout ***REMOVED***
			return ErrTimeout
		***REMOVED***

		var flag uint32 = flagLockFailImmediately
		if exclusive ***REMOVED***
			flag |= flagLockExclusive
		***REMOVED***

		err := lockFileEx(syscall.Handle(db.lockfile.Fd()), flag, 0, 1, 0, &syscall.Overlapped***REMOVED******REMOVED***)
		if err == nil ***REMOVED***
			return nil
		***REMOVED*** else if err != errLockViolation ***REMOVED***
			return err
		***REMOVED***

		// Wait for a bit and try again.
		time.Sleep(50 * time.Millisecond)
	***REMOVED***
***REMOVED***

// funlock releases an advisory lock on a file descriptor.
func funlock(db *DB) error ***REMOVED***
	err := unlockFileEx(syscall.Handle(db.lockfile.Fd()), 0, 1, 0, &syscall.Overlapped***REMOVED******REMOVED***)
	db.lockfile.Close()
	os.Remove(db.path+lockExt)
	return err
***REMOVED***

// mmap memory maps a DB's data file.
// Based on: https://github.com/edsrzf/mmap-go
func mmap(db *DB, sz int) error ***REMOVED***
	if !db.readOnly ***REMOVED***
		// Truncate the database to the size of the mmap.
		if err := db.file.Truncate(int64(sz)); err != nil ***REMOVED***
			return fmt.Errorf("truncate: %s", err)
		***REMOVED***
	***REMOVED***

	// Open a file mapping handle.
	sizelo := uint32(sz >> 32)
	sizehi := uint32(sz) & 0xffffffff
	h, errno := syscall.CreateFileMapping(syscall.Handle(db.file.Fd()), nil, syscall.PAGE_READONLY, sizelo, sizehi, nil)
	if h == 0 ***REMOVED***
		return os.NewSyscallError("CreateFileMapping", errno)
	***REMOVED***

	// Create the memory map.
	addr, errno := syscall.MapViewOfFile(h, syscall.FILE_MAP_READ, 0, 0, uintptr(sz))
	if addr == 0 ***REMOVED***
		return os.NewSyscallError("MapViewOfFile", errno)
	***REMOVED***

	// Close mapping handle.
	if err := syscall.CloseHandle(syscall.Handle(h)); err != nil ***REMOVED***
		return os.NewSyscallError("CloseHandle", err)
	***REMOVED***

	// Convert to a byte array.
	db.data = ((*[maxMapSize]byte)(unsafe.Pointer(addr)))
	db.datasz = sz

	return nil
***REMOVED***

// munmap unmaps a pointer from a file.
// Based on: https://github.com/edsrzf/mmap-go
func munmap(db *DB) error ***REMOVED***
	if db.data == nil ***REMOVED***
		return nil
	***REMOVED***

	addr := (uintptr)(unsafe.Pointer(&db.data[0]))
	if err := syscall.UnmapViewOfFile(addr); err != nil ***REMOVED***
		return os.NewSyscallError("UnmapViewOfFile", err)
	***REMOVED***
	return nil
***REMOVED***
