// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package user

import (
	"io"
	"os"

	"golang.org/x/sys/unix"
)

// Unix-specific path to the passwd and group formatted files.
const (
	unixPasswdPath = "/etc/passwd"
	unixGroupPath  = "/etc/group"
)

func GetPasswdPath() (string, error) ***REMOVED***
	return unixPasswdPath, nil
***REMOVED***

func GetPasswd() (io.ReadCloser, error) ***REMOVED***
	return os.Open(unixPasswdPath)
***REMOVED***

func GetGroupPath() (string, error) ***REMOVED***
	return unixGroupPath, nil
***REMOVED***

func GetGroup() (io.ReadCloser, error) ***REMOVED***
	return os.Open(unixGroupPath)
***REMOVED***

// CurrentUser looks up the current user by their user id in /etc/passwd. If the
// user cannot be found (or there is no /etc/passwd file on the filesystem),
// then CurrentUser returns an error.
func CurrentUser() (User, error) ***REMOVED***
	return LookupUid(unix.Getuid())
***REMOVED***

// CurrentGroup looks up the current user's group by their primary group id's
// entry in /etc/passwd. If the group cannot be found (or there is no
// /etc/group file on the filesystem), then CurrentGroup returns an error.
func CurrentGroup() (Group, error) ***REMOVED***
	return LookupGid(unix.Getgid())
***REMOVED***
