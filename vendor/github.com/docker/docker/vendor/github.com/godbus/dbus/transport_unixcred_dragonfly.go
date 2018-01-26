// The UnixCredentials system call is currently only implemented on Linux
// http://golang.org/src/pkg/syscall/sockcmsg_linux.go
// https://golang.org/s/go1.4-syscall
// http://code.google.com/p/go/source/browse/unix/sockcmsg_linux.go?repo=sys

// Local implementation of the UnixCredentials system call for DragonFly BSD

package dbus

/*
#include <sys/ucred.h>
*/
import "C"

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

// http://golang.org/src/pkg/syscall/ztypes_linux_amd64.go
// http://golang.org/src/pkg/syscall/ztypes_dragonfly_amd64.go
type Ucred struct ***REMOVED***
	Pid int32
	Uid uint32
	Gid uint32
***REMOVED***

// http://golang.org/src/pkg/syscall/types_linux.go
// http://golang.org/src/pkg/syscall/types_dragonfly.go
// https://github.com/DragonFlyBSD/DragonFlyBSD/blob/master/sys/sys/ucred.h
const (
	SizeofUcred = C.sizeof_struct_ucred
)

// http://golang.org/src/pkg/syscall/sockcmsg_unix.go
func cmsgAlignOf(salen int) int ***REMOVED***
	// From http://golang.org/src/pkg/syscall/sockcmsg_unix.go
	//salign := sizeofPtr
	// NOTE: It seems like 64-bit Darwin and DragonFly BSD kernels
	// still require 32-bit aligned access to network subsystem.
	//if darwin64Bit || dragonfly64Bit ***REMOVED***
	//	salign = 4
	//***REMOVED***
	salign := 4
	return (salen + salign - 1) & ^(salign - 1)
***REMOVED***

// http://golang.org/src/pkg/syscall/sockcmsg_unix.go
func cmsgData(h *syscall.Cmsghdr) unsafe.Pointer ***REMOVED***
	return unsafe.Pointer(uintptr(unsafe.Pointer(h)) + uintptr(cmsgAlignOf(syscall.SizeofCmsghdr)))
***REMOVED***

// http://golang.org/src/pkg/syscall/sockcmsg_linux.go
// UnixCredentials encodes credentials into a socket control message
// for sending to another process. This can be used for
// authentication.
func UnixCredentials(ucred *Ucred) []byte ***REMOVED***
	b := make([]byte, syscall.CmsgSpace(SizeofUcred))
	h := (*syscall.Cmsghdr)(unsafe.Pointer(&b[0]))
	h.Level = syscall.SOL_SOCKET
	h.Type = syscall.SCM_CREDS
	h.SetLen(syscall.CmsgLen(SizeofUcred))
	*((*Ucred)(cmsgData(h))) = *ucred
	return b
***REMOVED***

// http://golang.org/src/pkg/syscall/sockcmsg_linux.go
// ParseUnixCredentials decodes a socket control message that contains
// credentials in a Ucred structure. To receive such a message, the
// SO_PASSCRED option must be enabled on the socket.
func ParseUnixCredentials(m *syscall.SocketControlMessage) (*Ucred, error) ***REMOVED***
	if m.Header.Level != syscall.SOL_SOCKET ***REMOVED***
		return nil, syscall.EINVAL
	***REMOVED***
	if m.Header.Type != syscall.SCM_CREDS ***REMOVED***
		return nil, syscall.EINVAL
	***REMOVED***
	ucred := *(*Ucred)(unsafe.Pointer(&m.Data[0]))
	return &ucred, nil
***REMOVED***

func (t *unixTransport) SendNullByte() error ***REMOVED***
	ucred := &Ucred***REMOVED***Pid: int32(os.Getpid()), Uid: uint32(os.Getuid()), Gid: uint32(os.Getgid())***REMOVED***
	b := UnixCredentials(ucred)
	_, oobn, err := t.UnixConn.WriteMsgUnix([]byte***REMOVED***0***REMOVED***, b, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if oobn != len(b) ***REMOVED***
		return io.ErrShortWrite
	***REMOVED***
	return nil
***REMOVED***
