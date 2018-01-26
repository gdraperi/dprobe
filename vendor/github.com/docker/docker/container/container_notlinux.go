// +build freebsd

package container

import (
	"golang.org/x/sys/unix"
)

func detachMounted(path string) error ***REMOVED***
	// FreeBSD do not support the lazy unmount or MNT_DETACH feature.
	// Therefore there are separate definitions for this.
	return unix.Unmount(path, 0)
***REMOVED***

// SecretMounts returns the mounts for the secret path
func (container *Container) SecretMounts() []Mount ***REMOVED***
	return nil
***REMOVED***

// UnmountSecrets unmounts the fs for secrets
func (container *Container) UnmountSecrets() error ***REMOVED***
	return nil
***REMOVED***
