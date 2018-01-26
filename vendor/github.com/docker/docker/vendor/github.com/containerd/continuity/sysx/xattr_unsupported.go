// +build freebsd solaris

package sysx

// Listxattr calls syscall listxattr and reads all content
// and returns a string array
func Listxattr(path string) ([]string, error) ***REMOVED***
	return []string***REMOVED******REMOVED***, nil
***REMOVED***

// Removexattr calls syscall removexattr
func Removexattr(path string, attr string) (err error) ***REMOVED***
	return unsupported
***REMOVED***

// Setxattr calls syscall setxattr
func Setxattr(path string, attr string, data []byte, flags int) (err error) ***REMOVED***
	return unsupported
***REMOVED***

// Getxattr calls syscall getxattr
func Getxattr(path, attr string) ([]byte, error) ***REMOVED***
	return []byte***REMOVED******REMOVED***, unsupported
***REMOVED***

// LListxattr lists xattrs, not following symlinks
func LListxattr(path string) ([]string, error) ***REMOVED***
	return []string***REMOVED******REMOVED***, nil
***REMOVED***

// LRemovexattr removes an xattr, not following symlinks
func LRemovexattr(path string, attr string) (err error) ***REMOVED***
	return unsupported
***REMOVED***

// LSetxattr sets an xattr, not following symlinks
func LSetxattr(path string, attr string, data []byte, flags int) (err error) ***REMOVED***
	return unsupported
***REMOVED***

// LGetxattr gets an xattr, not following symlinks
func LGetxattr(path, attr string) ([]byte, error) ***REMOVED***
	return []byte***REMOVED******REMOVED***, nil
***REMOVED***
