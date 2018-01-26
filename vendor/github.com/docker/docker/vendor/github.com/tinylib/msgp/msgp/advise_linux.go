// +build linux,!appengine

package msgp

import (
	"os"
	"syscall"
)

func adviseRead(mem []byte) ***REMOVED***
	syscall.Madvise(mem, syscall.MADV_SEQUENTIAL|syscall.MADV_WILLNEED)
***REMOVED***

func adviseWrite(mem []byte) ***REMOVED***
	syscall.Madvise(mem, syscall.MADV_SEQUENTIAL)
***REMOVED***

func fallocate(f *os.File, sz int64) error ***REMOVED***
	err := syscall.Fallocate(int(f.Fd()), 0, 0, sz)
	if err == syscall.ENOTSUP ***REMOVED***
		return f.Truncate(sz)
	***REMOVED***
	return err
***REMOVED***
