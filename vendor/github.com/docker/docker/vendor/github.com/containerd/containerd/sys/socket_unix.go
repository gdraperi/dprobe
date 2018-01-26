// +build !windows

package sys

import (
	"net"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// CreateUnixSocket creates a unix socket and returns the listener
func CreateUnixSocket(path string) (net.Listener, error) ***REMOVED***
	if err := os.MkdirAll(filepath.Dir(path), 0660); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := unix.Unlink(path); err != nil && !os.IsNotExist(err) ***REMOVED***
		return nil, err
	***REMOVED***
	return net.Listen("unix", path)
***REMOVED***

// GetLocalListener returns a listerner out of a unix socket.
func GetLocalListener(path string, uid, gid int) (net.Listener, error) ***REMOVED***
	// Ensure parent directory is created
	if err := mkdirAs(filepath.Dir(path), uid, gid); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	l, err := CreateUnixSocket(path)
	if err != nil ***REMOVED***
		return l, err
	***REMOVED***

	if err := os.Chmod(path, 0660); err != nil ***REMOVED***
		l.Close()
		return nil, err
	***REMOVED***

	if err := os.Chown(path, uid, gid); err != nil ***REMOVED***
		l.Close()
		return nil, err
	***REMOVED***

	return l, nil
***REMOVED***

func mkdirAs(path string, uid, gid int) error ***REMOVED***
	if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) ***REMOVED***
		return err
	***REMOVED***

	if err := os.Mkdir(path, 0770); err != nil ***REMOVED***
		return err
	***REMOVED***

	return os.Chown(path, uid, gid)
***REMOVED***
