// +build !windows,!plan9,!solaris

package bolt

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

// flock acquires an advisory lock on a file descriptor.
func flock(db *DB, mode os.FileMode, exclusive bool, timeout time.Duration) error ***REMOVED***
	var t time.Time
	for ***REMOVED***
		// If we're beyond our timeout then return an error.
		// This can only occur after we've attempted a flock once.
		if t.IsZero() ***REMOVED***
			t = time.Now()
		***REMOVED*** else if timeout > 0 && time.Since(t) > timeout ***REMOVED***
			return ErrTimeout
		***REMOVED***
		flag := syscall.LOCK_SH
		if exclusive ***REMOVED***
			flag = syscall.LOCK_EX
		***REMOVED***

		// Otherwise attempt to obtain an exclusive lock.
		err := syscall.Flock(int(db.file.Fd()), flag|syscall.LOCK_NB)
		if err == nil ***REMOVED***
			return nil
		***REMOVED*** else if err != syscall.EWOULDBLOCK ***REMOVED***
			return err
		***REMOVED***

		// Wait for a bit and try again.
		time.Sleep(50 * time.Millisecond)
	***REMOVED***
***REMOVED***

// funlock releases an advisory lock on a file descriptor.
func funlock(db *DB) error ***REMOVED***
	return syscall.Flock(int(db.file.Fd()), syscall.LOCK_UN)
***REMOVED***

// mmap memory maps a DB's data file.
func mmap(db *DB, sz int) error ***REMOVED***
	// Map the data file to memory.
	b, err := syscall.Mmap(int(db.file.Fd()), 0, sz, syscall.PROT_READ, syscall.MAP_SHARED|db.MmapFlags)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Advise the kernel that the mmap is accessed randomly.
	if err := madvise(b, syscall.MADV_RANDOM); err != nil ***REMOVED***
		return fmt.Errorf("madvise: %s", err)
	***REMOVED***

	// Save the original byte slice and convert to a byte array pointer.
	db.dataref = b
	db.data = (*[maxMapSize]byte)(unsafe.Pointer(&b[0]))
	db.datasz = sz
	return nil
***REMOVED***

// munmap unmaps a DB's data file from memory.
func munmap(db *DB) error ***REMOVED***
	// Ignore the unmap if we have no mapped data.
	if db.dataref == nil ***REMOVED***
		return nil
	***REMOVED***

	// Unmap using the original byte slice.
	err := syscall.Munmap(db.dataref)
	db.dataref = nil
	db.data = nil
	db.datasz = 0
	return err
***REMOVED***

// NOTE: This function is copied from stdlib because it is not available on darwin.
func madvise(b []byte, advice int) (err error) ***REMOVED***
	_, _, e1 := syscall.Syscall(syscall.SYS_MADVISE, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), uintptr(advice))
	if e1 != 0 ***REMOVED***
		err = e1
	***REMOVED***
	return
***REMOVED***
