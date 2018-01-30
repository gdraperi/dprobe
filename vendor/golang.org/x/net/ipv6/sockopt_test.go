// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6_test

import (
	"fmt"
	"net"
	"runtime"
	"testing"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/ipv6"
)

var supportsIPv6 bool = nettest.SupportsIPv6()

func TestConnInitiatorPathMTU(t *testing.T) ***REMOVED***
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

	done := make(chan bool)
	go acceptor(t, ln, done)

	c, err := net.Dial("tcp6", ln.Addr().String())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()

	if pmtu, err := ipv6.NewConn(c).PathMTU(); err != nil ***REMOVED***
		switch runtime.GOOS ***REMOVED***
		case "darwin": // older darwin kernels don't support IPV6_PATHMTU option
			t.Logf("not supported on %s", runtime.GOOS)
		default:
			t.Fatal(err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		t.Logf("path mtu for %v: %v", c.RemoteAddr(), pmtu)
	***REMOVED***

	<-done
***REMOVED***

func TestConnResponderPathMTU(t *testing.T) ***REMOVED***
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

	done := make(chan bool)
	go connector(t, "tcp6", ln.Addr().String(), done)

	c, err := ln.Accept()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()

	if pmtu, err := ipv6.NewConn(c).PathMTU(); err != nil ***REMOVED***
		switch runtime.GOOS ***REMOVED***
		case "darwin": // older darwin kernels don't support IPV6_PATHMTU option
			t.Logf("not supported on %s", runtime.GOOS)
		default:
			t.Fatal(err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		t.Logf("path mtu for %v: %v", c.RemoteAddr(), pmtu)
	***REMOVED***

	<-done
***REMOVED***

func TestPacketConnChecksum(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if !supportsIPv6 ***REMOVED***
		t.Skip("ipv6 is not supported")
	***REMOVED***
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***

	c, err := net.ListenPacket(fmt.Sprintf("ip6:%d", iana.ProtocolOSPFIGP), "::") // OSPF for IPv6
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()

	p := ipv6.NewPacketConn(c)
	offset := 12 // see RFC 5340

	for _, toggle := range []bool***REMOVED***false, true***REMOVED*** ***REMOVED***
		if err := p.SetChecksum(toggle, offset); err != nil ***REMOVED***
			if toggle ***REMOVED***
				t.Fatalf("ipv6.PacketConn.SetChecksum(%v, %v) failed: %v", toggle, offset, err)
			***REMOVED*** else ***REMOVED***
				// Some platforms never allow to disable the kernel
				// checksum processing.
				t.Logf("ipv6.PacketConn.SetChecksum(%v, %v) failed: %v", toggle, offset, err)
			***REMOVED***
		***REMOVED***
		if on, offset, err := p.Checksum(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED*** else ***REMOVED***
			t.Logf("kernel checksum processing enabled=%v, offset=%v", on, offset)
		***REMOVED***
	***REMOVED***
***REMOVED***
