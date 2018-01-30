// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6_test

import (
	"net"
	"runtime"
	"testing"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/ipv6"
)

func TestConnUnicastSocketOptions(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if !supportsIPv6 ***REMOVED***
		t.Skip("ipv6 is not supported")
	***REMOVED***

	ln, err := net.Listen("tcp6", "[::1]:0")
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

	c, err := net.Dial("tcp6", ln.Addr().String())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()

	testUnicastSocketOptions(t, ipv6.NewConn(c))

	if err := <-errc; err != nil ***REMOVED***
		t.Errorf("server: %v", err)
	***REMOVED***
***REMOVED***

var packetConnUnicastSocketOptionTests = []struct ***REMOVED***
	net, proto, addr string
***REMOVED******REMOVED***
	***REMOVED***"udp6", "", "[::1]:0"***REMOVED***,
	***REMOVED***"ip6", ":ipv6-icmp", "::1"***REMOVED***,
***REMOVED***

func TestPacketConnUnicastSocketOptions(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if !supportsIPv6 ***REMOVED***
		t.Skip("ipv6 is not supported")
	***REMOVED***

	m, ok := nettest.SupportsRawIPSocket()
	for _, tt := range packetConnUnicastSocketOptionTests ***REMOVED***
		if tt.net == "ip6" && !ok ***REMOVED***
			t.Log(m)
			continue
		***REMOVED***
		c, err := net.ListenPacket(tt.net+tt.proto, tt.addr)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c.Close()

		testUnicastSocketOptions(t, ipv6.NewPacketConn(c))
	***REMOVED***
***REMOVED***

type testIPv6UnicastConn interface ***REMOVED***
	TrafficClass() (int, error)
	SetTrafficClass(int) error
	HopLimit() (int, error)
	SetHopLimit(int) error
***REMOVED***

func testUnicastSocketOptions(t *testing.T, c testIPv6UnicastConn) ***REMOVED***
	tclass := iana.DiffServCS0 | iana.NotECNTransport
	if err := c.SetTrafficClass(tclass); err != nil ***REMOVED***
		switch runtime.GOOS ***REMOVED***
		case "darwin": // older darwin kernels don't support IPV6_TCLASS option
			t.Logf("not supported on %s", runtime.GOOS)
			goto next
		***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if v, err := c.TrafficClass(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else if v != tclass ***REMOVED***
		t.Fatalf("got %v; want %v", v, tclass)
	***REMOVED***

next:
	hoplim := 255
	if err := c.SetHopLimit(hoplim); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if v, err := c.HopLimit(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else if v != hoplim ***REMOVED***
		t.Fatalf("got %v; want %v", v, hoplim)
	***REMOVED***
***REMOVED***
