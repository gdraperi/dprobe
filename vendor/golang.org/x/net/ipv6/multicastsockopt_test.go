// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6_test

import (
	"net"
	"runtime"
	"testing"

	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/ipv6"
)

var packetConnMulticastSocketOptionTests = []struct ***REMOVED***
	net, proto, addr string
	grp, src         net.Addr
***REMOVED******REMOVED***
	***REMOVED***"udp6", "", "[ff02::]:0", &net.UDPAddr***REMOVED***IP: net.ParseIP("ff02::114")***REMOVED***, nil***REMOVED***, // see RFC 4727
	***REMOVED***"ip6", ":ipv6-icmp", "::", &net.IPAddr***REMOVED***IP: net.ParseIP("ff02::115")***REMOVED***, nil***REMOVED***, // see RFC 4727

	***REMOVED***"udp6", "", "[ff30::8000:0]:0", &net.UDPAddr***REMOVED***IP: net.ParseIP("ff30::8000:1")***REMOVED***, &net.UDPAddr***REMOVED***IP: net.IPv6loopback***REMOVED******REMOVED***, // see RFC 5771
	***REMOVED***"ip6", ":ipv6-icmp", "::", &net.IPAddr***REMOVED***IP: net.ParseIP("ff30::8000:2")***REMOVED***, &net.IPAddr***REMOVED***IP: net.IPv6loopback***REMOVED******REMOVED***,        // see RFC 5771
***REMOVED***

func TestPacketConnMulticastSocketOptions(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if !supportsIPv6 ***REMOVED***
		t.Skip("ipv6 is not supported")
	***REMOVED***
	ifi := nettest.RoutedInterface("ip6", net.FlagUp|net.FlagMulticast|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	m, ok := nettest.SupportsRawIPSocket()
	for _, tt := range packetConnMulticastSocketOptionTests ***REMOVED***
		if tt.net == "ip6" && !ok ***REMOVED***
			t.Log(m)
			continue
		***REMOVED***
		c, err := net.ListenPacket(tt.net+tt.proto, tt.addr)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c.Close()
		p := ipv6.NewPacketConn(c)
		defer p.Close()

		if tt.src == nil ***REMOVED***
			testMulticastSocketOptions(t, p, ifi, tt.grp)
		***REMOVED*** else ***REMOVED***
			testSourceSpecificMulticastSocketOptions(t, p, ifi, tt.grp, tt.src)
		***REMOVED***
	***REMOVED***
***REMOVED***

type testIPv6MulticastConn interface ***REMOVED***
	MulticastHopLimit() (int, error)
	SetMulticastHopLimit(ttl int) error
	MulticastLoopback() (bool, error)
	SetMulticastLoopback(bool) error
	JoinGroup(*net.Interface, net.Addr) error
	LeaveGroup(*net.Interface, net.Addr) error
	JoinSourceSpecificGroup(*net.Interface, net.Addr, net.Addr) error
	LeaveSourceSpecificGroup(*net.Interface, net.Addr, net.Addr) error
	ExcludeSourceSpecificGroup(*net.Interface, net.Addr, net.Addr) error
	IncludeSourceSpecificGroup(*net.Interface, net.Addr, net.Addr) error
***REMOVED***

func testMulticastSocketOptions(t *testing.T, c testIPv6MulticastConn, ifi *net.Interface, grp net.Addr) ***REMOVED***
	const hoplim = 255
	if err := c.SetMulticastHopLimit(hoplim); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	if v, err := c.MulticastHopLimit(); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED*** else if v != hoplim ***REMOVED***
		t.Errorf("got %v; want %v", v, hoplim)
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

func testSourceSpecificMulticastSocketOptions(t *testing.T, c testIPv6MulticastConn, ifi *net.Interface, grp, src net.Addr) ***REMOVED***
	// MCAST_JOIN_GROUP -> MCAST_BLOCK_SOURCE -> MCAST_UNBLOCK_SOURCE -> MCAST_LEAVE_GROUP
	if err := c.JoinGroup(ifi, grp); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	if err := c.ExcludeSourceSpecificGroup(ifi, grp, src); err != nil ***REMOVED***
		switch runtime.GOOS ***REMOVED***
		case "freebsd", "linux":
		default: // platforms that don't support MLDv2 fail here
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
