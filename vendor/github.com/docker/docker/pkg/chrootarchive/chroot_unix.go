// +build !windows,!linux

package chrootarchive

import "golang.org/x/sys/unix"

func chroot(path string) error ***REMOVED***
	if err := unix.Chroot(path); err != nil ***REMOVED***
		return err
	***REMOVED***
	return unix.Chdir("/")
***REMOVED***
