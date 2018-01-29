// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package test

// direct-tcpip and direct-streamlocal functional tests

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"
	"testing"
)

type dialTester interface ***REMOVED***
	TestServerConn(t *testing.T, c net.Conn)
	TestClientConn(t *testing.T, c net.Conn)
***REMOVED***

func testDial(t *testing.T, n, listenAddr string, x dialTester) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	sshConn := server.Dial(clientConfig())
	defer sshConn.Close()

	l, err := net.Listen(n, listenAddr)
	if err != nil ***REMOVED***
		t.Fatalf("Listen: %v", err)
	***REMOVED***
	defer l.Close()

	testData := fmt.Sprintf("hello from %s, %s", n, listenAddr)
	go func() ***REMOVED***
		for ***REMOVED***
			c, err := l.Accept()
			if err != nil ***REMOVED***
				break
			***REMOVED***
			x.TestServerConn(t, c)

			io.WriteString(c, testData)
			c.Close()
		***REMOVED***
	***REMOVED***()

	conn, err := sshConn.Dial(n, l.Addr().String())
	if err != nil ***REMOVED***
		t.Fatalf("Dial: %v", err)
	***REMOVED***
	x.TestClientConn(t, conn)
	defer conn.Close()
	b, err := ioutil.ReadAll(conn)
	if err != nil ***REMOVED***
		t.Fatalf("ReadAll: %v", err)
	***REMOVED***
	t.Logf("got %q", string(b))
	if string(b) != testData ***REMOVED***
		t.Fatalf("expected %q, got %q", testData, string(b))
	***REMOVED***
***REMOVED***

type tcpDialTester struct ***REMOVED***
	listenAddr string
***REMOVED***

func (x *tcpDialTester) TestServerConn(t *testing.T, c net.Conn) ***REMOVED***
	host := strings.Split(x.listenAddr, ":")[0]
	prefix := host + ":"
	if !strings.HasPrefix(c.LocalAddr().String(), prefix) ***REMOVED***
		t.Fatalf("expected to start with %q, got %q", prefix, c.LocalAddr().String())
	***REMOVED***
	if !strings.HasPrefix(c.RemoteAddr().String(), prefix) ***REMOVED***
		t.Fatalf("expected to start with %q, got %q", prefix, c.RemoteAddr().String())
	***REMOVED***
***REMOVED***

func (x *tcpDialTester) TestClientConn(t *testing.T, c net.Conn) ***REMOVED***
	// we use zero addresses. see *Client.Dial.
	if c.LocalAddr().String() != "0.0.0.0:0" ***REMOVED***
		t.Fatalf("expected \"0.0.0.0:0\", got %q", c.LocalAddr().String())
	***REMOVED***
	if c.RemoteAddr().String() != "0.0.0.0:0" ***REMOVED***
		t.Fatalf("expected \"0.0.0.0:0\", got %q", c.RemoteAddr().String())
	***REMOVED***
***REMOVED***

func TestDialTCP(t *testing.T) ***REMOVED***
	x := &tcpDialTester***REMOVED***
		listenAddr: "127.0.0.1:0",
	***REMOVED***
	testDial(t, "tcp", x.listenAddr, x)
***REMOVED***

type unixDialTester struct ***REMOVED***
	listenAddr string
***REMOVED***

func (x *unixDialTester) TestServerConn(t *testing.T, c net.Conn) ***REMOVED***
	if c.LocalAddr().String() != x.listenAddr ***REMOVED***
		t.Fatalf("expected %q, got %q", x.listenAddr, c.LocalAddr().String())
	***REMOVED***
	if c.RemoteAddr().String() != "@" ***REMOVED***
		t.Fatalf("expected \"@\", got %q", c.RemoteAddr().String())
	***REMOVED***
***REMOVED***

func (x *unixDialTester) TestClientConn(t *testing.T, c net.Conn) ***REMOVED***
	if c.RemoteAddr().String() != x.listenAddr ***REMOVED***
		t.Fatalf("expected %q, got %q", x.listenAddr, c.RemoteAddr().String())
	***REMOVED***
	if c.LocalAddr().String() != "@" ***REMOVED***
		t.Fatalf("expected \"@\", got %q", c.LocalAddr().String())
	***REMOVED***
***REMOVED***

func TestDialUnix(t *testing.T) ***REMOVED***
	addr, cleanup := newTempSocket(t)
	defer cleanup()
	x := &unixDialTester***REMOVED***
		listenAddr: addr,
	***REMOVED***
	testDial(t, "unix", x.listenAddr, x)
***REMOVED***
