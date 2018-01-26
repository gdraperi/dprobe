// The UnixCredentials system call is currently only implemented on Linux
// http://golang.org/src/pkg/syscall/sockcmsg_linux.go
// https://golang.org/s/go1.4-syscall
// http://code.google.com/p/go/source/browse/unix/sockcmsg_linux.go?repo=sys

package dbus

import (
	"io"
	"os"
	"syscall"
)

func (t *unixTransport) SendNullByte() error ***REMOVED***
	ucred := &syscall.Ucred***REMOVED***Pid: int32(os.Getpid()), Uid: uint32(os.Getuid()), Gid: uint32(os.Getgid())***REMOVED***
	b := syscall.UnixCredentials(ucred)
	_, oobn, err := t.UnixConn.WriteMsgUnix([]byte***REMOVED***0***REMOVED***, b, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if oobn != len(b) ***REMOVED***
		return io.ErrShortWrite
	***REMOVED***
	return nil
***REMOVED***
