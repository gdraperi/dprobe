package mount

/*
#include <errno.h>
#include <stdlib.h>
#include <string.h>
#include <sys/_iovec.h>
#include <sys/mount.h>
#include <sys/param.h>
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

func allocateIOVecs(options []string) []C.struct_iovec ***REMOVED***
	out := make([]C.struct_iovec, len(options))
	for i, option := range options ***REMOVED***
		out[i].iov_base = unsafe.Pointer(C.CString(option))
		out[i].iov_len = C.size_t(len(option) + 1)
	***REMOVED***
	return out
***REMOVED***

func mount(device, target, mType string, flag uintptr, data string) error ***REMOVED***
	isNullFS := false

	xs := strings.Split(data, ",")
	for _, x := range xs ***REMOVED***
		if x == "bind" ***REMOVED***
			isNullFS = true
		***REMOVED***
	***REMOVED***

	options := []string***REMOVED***"fspath", target***REMOVED***
	if isNullFS ***REMOVED***
		options = append(options, "fstype", "nullfs", "target", device)
	***REMOVED*** else ***REMOVED***
		options = append(options, "fstype", mType, "from", device)
	***REMOVED***
	rawOptions := allocateIOVecs(options)
	for _, rawOption := range rawOptions ***REMOVED***
		defer C.free(rawOption.iov_base)
	***REMOVED***

	if errno := C.nmount(&rawOptions[0], C.uint(len(options)), C.int(flag)); errno != 0 ***REMOVED***
		reason := C.GoString(C.strerror(*C.__error()))
		return fmt.Errorf("Failed to call nmount: %s", reason)
	***REMOVED***
	return nil
***REMOVED***

func unmount(target string, flag int) error ***REMOVED***
	return unix.Unmount(target, flag)
***REMOVED***
