package container

import (
	"golang.org/x/sys/unix"
)

func detachMounted(path string) error ***REMOVED***
	return unix.Unmount(path, unix.MNT_DETACH)
***REMOVED***
