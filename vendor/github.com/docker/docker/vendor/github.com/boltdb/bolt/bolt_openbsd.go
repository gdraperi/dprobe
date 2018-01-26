package bolt

import (
	"syscall"
	"unsafe"
)

const (
	msAsync      = 1 << iota // perform asynchronous writes
	msSync                   // perform synchronous writes
	msInvalidate             // invalidate cached data
)

func msync(db *DB) error ***REMOVED***
	_, _, errno := syscall.Syscall(syscall.SYS_MSYNC, uintptr(unsafe.Pointer(db.data)), uintptr(db.datasz), msInvalidate)
	if errno != 0 ***REMOVED***
		return errno
	***REMOVED***
	return nil
***REMOVED***

func fdatasync(db *DB) error ***REMOVED***
	if db.data != nil ***REMOVED***
		return msync(db)
	***REMOVED***
	return db.file.Sync()
***REMOVED***
