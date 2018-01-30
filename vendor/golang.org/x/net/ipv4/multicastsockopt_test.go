// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4_test

import (
	"net"
	"runtime"
	"testing"

	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/ipv4"
)

var packetConnMulticastSocketOptionTests = []struct ***REMOVED***
	net, proto, addr string
	grp, src         net.Addr
***REMOVED******REMOVED***
	***REMOVED***"udp4", "", "224.0.0.0:0", &net.UDPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 249)***REMOVED***, nil***REMOVED***, // see RFC 4727
	***REMOVED***"ip4", ":icmp", "0.0.0.0", &net.IPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 250)***REMOVED***, nil***REMOVED***,  // see RFC 4727

	***REMOVED***"udp4", "", "232.0.0.0:0", &net.UDPAddr***REMOVED***IP: net.IPv4(232, 0, 1, 249)***REMOVED***, &net.UDPAddr***REMOVED***IP: net.IPv4(127, 0, 0, 1)***REMOVED******REMOVED***, // see RFC 5771
	***REMOVED***"ip4", ":icmp", "0.0.0.0", &net.IPAddr***REMOVED***IP: net.IPv4(232, 0, 1, 250)***REMOVED***, &net.UDPAddr***REMOVED***IP: net.IPv4(127, 0, 0, 1)***REMOVED******REMOVED***,  // see RFC 5771
***REMOVED***

func TestPacketConnMulticastSocketOptions(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagMulticast|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	m, ok := nettest.SupportsRawIPSocket()
	for _, tt := range packetConnMulticastSocketOptionTests ***REMOVED***
		if tt.net == "ip4" && !ok ***REMOVED***
			t.Log(m)
			continue
		***REMOVED***
		c, err := net.ListenPacket(tt.net+tt.proto, tt.addr)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c.Close()
		p := ipv4.NewPacketConn(c)
		defer p.Close()

		if tt.src == nil ***REMOVED***
			testMulticastSocketOptions(t, p, ifi, tt.grp)
		***REMOVED*** else ***REMOVED***
			testSourceSpecificMulticastSocketOptions(t, p, ifi, tt.grp, tt.src)
		***REMOVED***
	***REMOVED***
***REMOVED***

var rawConnMulticastSocketOptionTests = []struct ***REMOVED***
	grp, src net.Addr
***REMOVED******REMOVED***
	***REMOVED***&net.IPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 250)***REMOVED***, nil***REMOVED***, // see RFC 4727

	***REMOVED***&net.IPAddr***REMOVED***IP: net.IPv4(232, 0, 1, 250)***REMOVED***, &net.IPAddr***REMOVED***IP: net.IPv4(127, 0, 0, 1)***REMOVED******REMOVED***, // see RFC 5771
***REMOVED***

func TestRawConnMulticastSocketOptions(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagMulticast|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	for _, tt := range rawConnMulticastSocketOptionTests ***REMOVED***
		c, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c.Close()
		r, err := ipv4.NewRawConn(c)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer r.Close()

		if tt.src == nil ***REMOVED***
			testMulticastSocketOptions(t, r, ifi, tt.grp)
		***REMOVED*** else ***REMOVED***
			testSourceSpecificMulticastSocketOptions(t, r, ifi, tt.grp, tt.src)
		***REMOVED***
	***REMOVED***
***REMOVED***

type testIPv4MulticastConn interface ***REMOVED***
	MulticastTTL() (int, error)
	SetMulticastTTL(ttl int) error
	MulticastLoopback() (bool, error)
	SetMulticastLoopback(bool) error
	JoinGroup(*net.Interface, net.Addr) error
	LeaveGroup(*net.Interface, net.Addr) error
	JoinSourceSpecificGroup(*net.Interface, net.Addr, net.Addr) error
	LeaveSourceSpecificGroup(*net.Interface, net.Addr, net.Addr) error
	ExcludeSourceSpecificGroup(*net.Interface, net.Addr, net.Addr) error
	IncludeSourceSpecificGroup(*net.Interface, net.Addr, net.Addr) error
***REMOVED***

func testMulticastSocketOptions(t *testing.T, c testIPv4MulticastConn, ifi *net.Interface, grp net.Addr) ***REMOVED***
	const ttl = 255
	if err := c.SetMulticastTTL(ttl); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	if v, err := c.MulticastTTL(); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED*** else if v != ttl ***REMOVED***
		t.Errorf("got %v; want %v", v, ttl)
		return
	***REMOVED***

	for _, toggle := range []bool***REMOVED***true, false***REMOVED*** ***REMOVED***
		if err := c.SetMulticastLoopback(toggle); err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		if v, err := c.MulticastLoopback(); err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED*** else if v != toggle ***REMOVED***
			t.Errorf("got %v; want %v", v, toggle)
			return
		***REMOVED***
	***REMOVED***

	if err := c.JoinGroup(ifi, grp); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	if err := c.LeaveGroup(ifi, grp); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
***REMOVED***

func testSourceSpecificMulticastSocketOptions(t *testing.T, c testIPv4MulticastConn, ifi *net.Interface, grp, src net.Addr) ***REMOVED***
	// MCAST_JOIN_GROUP -> MCAST_BLOCK_SOURCE -> MCAST_UNBLOCK_SOURCE -> MCAST_LEAVE_GROUP
	if err := c.JoinGroup(ifi, grp); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	if err := c.ExcludeSourceSpecificGroup(ifi, grp, src); err != nil ***REMOVED***
		switch runtime.GOOS ***REMOVED***
		case "freebsd", "linux":
		default: // platforms that don't support IGMPv2/3 fail here
			t.Logf("not supported on %s", runtime.GOOS)
			return
		***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	if err := c.IncludeSourceSpecificGroup(ifi, grp, src); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	if err := c.LeaveGroup(ifi, grp); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***

	// MCAST_JOIN_SOURCE_GROUP -> MCAST_LEAVE_SOURCE_GROUP
	if err := c.JoinSourceSpecificGroup(ifi, grp, src); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	if err := c.LeaveSourceSpecificGroup(ifi, grp, src); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***

	// MCAST_JOIN_SOURCE_GROUP -> MCAST_LEAVE_GROUP
	if err := c.JoinSourceSpecificGroup(ifi, grp, src); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	if err := c.LeaveGroup(ifi, grp); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
***REMOVED***
