// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd

package test

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"testing"
	"time"
)

type closeWriter interface ***REMOVED***
	CloseWrite() error
***REMOVED***

func testPortForward(t *testing.T, n, listenAddr string) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())
	defer conn.Close()

	sshListener, err := conn.Listen(n, listenAddr)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	go func() ***REMOVED***
		sshConn, err := sshListener.Accept()
		if err != nil ***REMOVED***
			t.Fatalf("listen.Accept failed: %v", err)
		***REMOVED***

		_, err = io.Copy(sshConn, sshConn)
		if err != nil && err != io.EOF ***REMOVED***
			t.Fatalf("ssh client copy: %v", err)
		***REMOVED***
		sshConn.Close()
	***REMOVED***()

	forwardedAddr := sshListener.Addr().String()
	netConn, err := net.Dial(n, forwardedAddr)
	if err != nil ***REMOVED***
		t.Fatalf("net dial failed: %v", err)
	***REMOVED***

	readChan := make(chan []byte)
	go func() ***REMOVED***
		data, _ := ioutil.ReadAll(netConn)
		readChan <- data
	***REMOVED***()

	// Invent some data.
	data := make([]byte, 100*1000)
	for i := range data ***REMOVED***
		data[i] = byte(i % 255)
	***REMOVED***

	var sent []byte
	for len(sent) < 1000*1000 ***REMOVED***
		// Send random sized chunks
		m := rand.Intn(len(data))
		n, err := netConn.Write(data[:m])
		if err != nil ***REMOVED***
			break
		***REMOVED***
		sent = append(sent, data[:n]...)
	***REMOVED***
	if err := netConn.(closeWriter).CloseWrite(); err != nil ***REMOVED***
		t.Errorf("netConn.CloseWrite: %v", err)
	***REMOVED***

	read := <-readChan

	if len(sent) != len(read) ***REMOVED***
		t.Fatalf("got %d bytes, want %d", len(read), len(sent))
	***REMOVED***
	if bytes.Compare(sent, read) != 0 ***REMOVED***
		t.Fatalf("read back data does not match")
	***REMOVED***

	if err := sshListener.Close(); err != nil ***REMOVED***
		t.Fatalf("sshListener.Close: %v", err)
	***REMOVED***

	// Check that the forward disappeared.
	netConn, err = net.Dial(n, forwardedAddr)
	if err == nil ***REMOVED***
		netConn.Close()
		t.Errorf("still listening to %s after closing", forwardedAddr)
	***REMOVED***
***REMOVED***

func TestPortForwardTCP(t *testing.T) ***REMOVED***
	testPortForward(t, "tcp", "localhost:0")
***REMOVED***

func TestPortForwardUnix(t *testing.T) ***REMOVED***
	addr, cleanup := newTempSocket(t)
	defer cleanup()
	testPortForward(t, "unix", addr)
***REMOVED***

func testAcceptClose(t *testing.T, n, listenAddr string) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())

	sshListener, err := conn.Listen(n, listenAddr)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	quit := make(chan error, 1)
	go func() ***REMOVED***
		for ***REMOVED***
			c, err := sshListener.Accept()
			if err != nil ***REMOVED***
				quit <- err
				break
			***REMOVED***
			c.Close()
		***REMOVED***
	***REMOVED***()
	sshListener.Close()

	select ***REMOVED***
	case <-time.After(1 * time.Second):
		t.Errorf("timeout: listener did not close.")
	case err := <-quit:
		t.Logf("quit as expected (error %v)", err)
	***REMOVED***
***REMOVED***

func TestAcceptCloseTCP(t *testing.T) ***REMOVED***
	testAcceptClose(t, "tcp", "localhost:0")
***REMOVED***

func TestAcceptCloseUnix(t *testing.T) ***REMOVED***
	addr, cleanup := newTempSocket(t)
	defer cleanup()
	testAcceptClose(t, "unix", addr)
***REMOVED***

// Check that listeners exit if the underlying client transport dies.
func testPortForwardConnectionClose(t *testing.T, n, listenAddr string) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())

	sshListener, err := conn.Listen(n, listenAddr)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	quit := make(chan error, 1)
	go func() ***REMOVED***
		for ***REMOVED***
			c, err := sshListener.Accept()
			if err != nil ***REMOVED***
				quit <- err
				break
			***REMOVED***
			c.Close()
		***REMOVED***
	***REMOVED***()

	// It would be even nicer if we closed the server side, but it
	// is more involved as the fd for that side is dup()ed.
	server.clientConn.Close()

	select ***REMOVED***
	case <-time.After(1 * time.Second):
		t.Errorf("timeout: listener did not close.")
	case err := <-quit:
		t.Logf("quit as expected (error %v)", err)
	***REMOVED***
***REMOVED***

func TestPortForwardConnectionCloseTCP(t *testing.T) ***REMOVED***
	testPortForwardConnectionClose(t, "tcp", "localhost:0")
***REMOVED***

func TestPortForwardConnectionCloseUnix(t *testing.T) ***REMOVED***
	addr, cleanup := newTempSocket(t)
	defer cleanup()
	testPortForwardConnectionClose(t, "unix", addr)
***REMOVED***
