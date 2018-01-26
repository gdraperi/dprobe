package system

import "golang.org/x/sys/unix"

// Lgetxattr retrieves the value of the extended attribute identified by attr
// and associated with the given path in the file system.
// It will returns a nil slice and nil error if the xattr is not set.
func Lgetxattr(path string, attr string) ([]byte, error) ***REMOVED***
	dest := make([]byte, 128)
	sz, errno := unix.Lgetxattr(path, attr, dest)
	if errno == unix.ENODATA ***REMOVED***
		return nil, nil
	***REMOVED***
	if errno == unix.ERANGE ***REMOVED***
		dest = make([]byte, sz)
		sz, errno = unix.Lgetxattr(path, attr, dest)
	***REMOVED***
	if errno != nil ***REMOVED***
		return nil, errno
	***REMOVED***

	return dest[:sz], nil
***REMOVED***

// Lsetxattr sets the value of the extended attribute identified by attr
// and associated with the given path in the file system.
func Lsetxattr(path string, attr string, data []byte, flags int) error ***REMOVED***
	return unix.Lsetxattr(path, attr, data, flags)
***REMOVED***
