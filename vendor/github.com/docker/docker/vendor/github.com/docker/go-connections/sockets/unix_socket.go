// +build !windows

package sockets

import (
	"net"
	"os"
	"syscall"
)

// NewUnixSocket creates a unix socket with the specified path and group.
func NewUnixSocket(path string, gid int) (net.Listener, error) ***REMOVED***
	if err := syscall.Unlink(path); err != nil && !os.IsNotExist(err) ***REMOVED***
		return nil, err
	***REMOVED***
	mask := syscall.Umask(0777)
	defer syscall.Umask(mask)

	l, err := net.Listen("unix", path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := os.Chown(path, 0, gid); err != nil ***REMOVED***
		l.Close()
		return nil, err
	***REMOVED***
	if err := os.Chmod(path, 0660); err != nil ***REMOVED***
		l.Close()
		return nil, err
	***REMOVED***
	return l, nil
***REMOVED***
