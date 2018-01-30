// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4_test

import (
	"net"
	"runtime"
	"testing"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/ipv4"
)

func TestConnUnicastSocketOptions(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer ln.Close()

	errc := make(chan error, 1)
	go func() ***REMOVED***
		c, err := ln.Accept()
		if err != nil ***REMOVED***
			errc <- err
			return
		***REMOVED***
		errc <- c.Close()
	***REMOVED***()

	c, err := net.Dial("tcp4", ln.Addr().String())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()

	testUnicastSocketOptions(t, ipv4.NewConn(c))

	if err := <-errc; err != nil ***REMOVED***
		t.Errorf("server: %v", err)
	***REMOVED***
***REMOVED***

var packetConnUnicastSocketOptionTests = []struct ***REMOVED***
	net, proto, addr string
***REMOVED******REMOVED***
	***REMOVED***"udp4", "", "127.0.0.1:0"***REMOVED***,
	***REMOVED***"ip4", ":icmp", "127.0.0.1"***REMOVED***,
***REMOVED***

func TestPacketConnUnicastSocketOptions(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	m, ok := nettest.SupportsRawIPSocket()
	for _, tt := range packetConnUnicastSocketOptionTests ***REMOVED***
		if tt.net == "ip4" && !ok ***REMOVED***
			t.Log(m)
			continue
		***REMOVED***
		c, err := net.ListenPacket(tt.net+tt.proto, tt.addr)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c.Close()

		testUnicastSocketOptions(t, ipv4.NewPacketConn(c))
	***REMOVED***
***REMOVED***

func TestRawConnUnicastSocketOptions(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	c, err := net.ListenPacket("ip4:icmp", "127.0.0.1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()

	r, err := ipv4.NewRawConn(c)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	testUnicastSocketOptions(t, r)
***REMOVED***

type testIPv4UnicastConn interface ***REMOVED***
	TOS() (int, error)
	SetTOS(int) error
	TTL() (int, error)
	SetTTL(int) error
***REMOVED***

func testUnicastSocketOptions(t *testing.T, c testIPv4UnicastConn) ***REMOVED***
	tos := iana.DiffServCS0 | iana.NotECNTransport
	switch runtime.GOOS ***REMOVED***
	case "windows":
		// IP_TOS option is supported on Windows 8 and beyond.
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***

	if err := c.SetTOS(tos); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if v, err := c.TOS(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else if v != tos ***REMOVED***
		t.Fatalf("got %v; want %v", v, tos)
	***REMOVED***
	const ttl = 255
	if err := c.SetTTL(ttl); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if v, err := c.TTL(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else if v != ttl ***REMOVED***
		t.Fatalf("got %v; want %v", v, ttl)
	***REMOVED***
***REMOVED***
