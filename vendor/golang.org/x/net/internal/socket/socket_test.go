// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package socket_test

import (
	"net"
	"runtime"
	"syscall"
	"testing"

	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/internal/socket"
)

func TestSocket(t *testing.T) ***REMOVED***
	t.Run("Option", func(t *testing.T) ***REMOVED***
		testSocketOption(t, &socket.Option***REMOVED***Level: syscall.SOL_SOCKET, Name: syscall.SO_RCVBUF, Len: 4***REMOVED***)
	***REMOVED***)
***REMOVED***

func testSocketOption(t *testing.T, so *socket.Option) ***REMOVED***
	c, err := nettest.NewLocalPacketListener("udp")
	if err != nil ***REMOVED***
		t.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
	***REMOVED***
	defer c.Close()
	cc, err := socket.NewConn(c.(net.Conn))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	const N = 2048
	if err := so.SetInt(cc, N); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	n, err := so.GetInt(cc)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if n < N ***REMOVED***
		t.Fatalf("got %d; want greater than or equal to %d", n, N)
	***REMOVED***
***REMOVED***
