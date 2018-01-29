// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package unix_test

import (
	"bytes"
	"go/build"
	"net"
	"os"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

// TestSCMCredentials tests the sending and receiving of credentials
// (PID, UID, GID) in an ancillary message between two UNIX
// sockets. The SO_PASSCRED socket option is enabled on the sending
// socket for this to work.
func TestSCMCredentials(t *testing.T) ***REMOVED***
	socketTypeTests := []struct ***REMOVED***
		socketType int
		dataLen    int
	***REMOVED******REMOVED***
		***REMOVED***
			unix.SOCK_STREAM,
			1,
		***REMOVED***, ***REMOVED***
			unix.SOCK_DGRAM,
			0,
		***REMOVED***,
	***REMOVED***

	for _, tt := range socketTypeTests ***REMOVED***
		if tt.socketType == unix.SOCK_DGRAM && !atLeast1p10() ***REMOVED***
			t.Log("skipping DGRAM test on pre-1.10")
			continue
		***REMOVED***

		fds, err := unix.Socketpair(unix.AF_LOCAL, tt.socketType, 0)
		if err != nil ***REMOVED***
			t.Fatalf("Socketpair: %v", err)
		***REMOVED***
		defer unix.Close(fds[0])
		defer unix.Close(fds[1])

		err = unix.SetsockoptInt(fds[0], unix.SOL_SOCKET, unix.SO_PASSCRED, 1)
		if err != nil ***REMOVED***
			t.Fatalf("SetsockoptInt: %v", err)
		***REMOVED***

		srvFile := os.NewFile(uintptr(fds[0]), "server")
		defer srvFile.Close()
		srv, err := net.FileConn(srvFile)
		if err != nil ***REMOVED***
			t.Errorf("FileConn: %v", err)
			return
		***REMOVED***
		defer srv.Close()

		cliFile := os.NewFile(uintptr(fds[1]), "client")
		defer cliFile.Close()
		cli, err := net.FileConn(cliFile)
		if err != nil ***REMOVED***
			t.Errorf("FileConn: %v", err)
			return
		***REMOVED***
		defer cli.Close()

		var ucred unix.Ucred
		if os.Getuid() != 0 ***REMOVED***
			ucred.Pid = int32(os.Getpid())
			ucred.Uid = 0
			ucred.Gid = 0
			oob := unix.UnixCredentials(&ucred)
			_, _, err := cli.(*net.UnixConn).WriteMsgUnix(nil, oob, nil)
			if op, ok := err.(*net.OpError); ok ***REMOVED***
				err = op.Err
			***REMOVED***
			if sys, ok := err.(*os.SyscallError); ok ***REMOVED***
				err = sys.Err
			***REMOVED***
			if err != syscall.EPERM ***REMOVED***
				t.Fatalf("WriteMsgUnix failed with %v, want EPERM", err)
			***REMOVED***
		***REMOVED***

		ucred.Pid = int32(os.Getpid())
		ucred.Uid = uint32(os.Getuid())
		ucred.Gid = uint32(os.Getgid())
		oob := unix.UnixCredentials(&ucred)

		// On SOCK_STREAM, this is internally going to send a dummy byte
		n, oobn, err := cli.(*net.UnixConn).WriteMsgUnix(nil, oob, nil)
		if err != nil ***REMOVED***
			t.Fatalf("WriteMsgUnix: %v", err)
		***REMOVED***
		if n != 0 ***REMOVED***
			t.Fatalf("WriteMsgUnix n = %d, want 0", n)
		***REMOVED***
		if oobn != len(oob) ***REMOVED***
			t.Fatalf("WriteMsgUnix oobn = %d, want %d", oobn, len(oob))
		***REMOVED***

		oob2 := make([]byte, 10*len(oob))
		n, oobn2, flags, _, err := srv.(*net.UnixConn).ReadMsgUnix(nil, oob2)
		if err != nil ***REMOVED***
			t.Fatalf("ReadMsgUnix: %v", err)
		***REMOVED***
		if flags != 0 ***REMOVED***
			t.Fatalf("ReadMsgUnix flags = 0x%x, want 0", flags)
		***REMOVED***
		if n != tt.dataLen ***REMOVED***
			t.Fatalf("ReadMsgUnix n = %d, want %d", n, tt.dataLen)
		***REMOVED***
		if oobn2 != oobn ***REMOVED***
			// without SO_PASSCRED set on the socket, ReadMsgUnix will
			// return zero oob bytes
			t.Fatalf("ReadMsgUnix oobn = %d, want %d", oobn2, oobn)
		***REMOVED***
		oob2 = oob2[:oobn2]
		if !bytes.Equal(oob, oob2) ***REMOVED***
			t.Fatal("ReadMsgUnix oob bytes don't match")
		***REMOVED***

		scm, err := unix.ParseSocketControlMessage(oob2)
		if err != nil ***REMOVED***
			t.Fatalf("ParseSocketControlMessage: %v", err)
		***REMOVED***
		newUcred, err := unix.ParseUnixCredentials(&scm[0])
		if err != nil ***REMOVED***
			t.Fatalf("ParseUnixCredentials: %v", err)
		***REMOVED***
		if *newUcred != ucred ***REMOVED***
			t.Fatalf("ParseUnixCredentials = %+v, want %+v", newUcred, ucred)
		***REMOVED***
	***REMOVED***
***REMOVED***

// atLeast1p10 reports whether we are running on Go 1.10 or later.
func atLeast1p10() bool ***REMOVED***
	for _, ver := range build.Default.ReleaseTags ***REMOVED***
		if ver == "go1.10" ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
