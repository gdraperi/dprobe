package bolt

import (
	"syscall"
)

// fdatasync flushes written data to a file descriptor.
func fdatasync(db *DB) error ***REMOVED***
	return syscall.Fdatasync(int(db.file.Fd()))
***REMOVED***
