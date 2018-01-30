// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp_test

import (
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func googleAddr(c *icmp.PacketConn, protocol int) (net.Addr, error) ***REMOVED***
	const host = "www.google.com"
	ips, err := net.LookupIP(host)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	netaddr := func(ip net.IP) (net.Addr, error) ***REMOVED***
		switch c.LocalAddr().(type) ***REMOVED***
		case *net.UDPAddr:
			return &net.UDPAddr***REMOVED***IP: ip***REMOVED***, nil
		case *net.IPAddr:
			return &net.IPAddr***REMOVED***IP: ip***REMOVED***, nil
		default:
			return nil, errors.New("neither UDPAddr nor IPAddr")
		***REMOVED***
	***REMOVED***
	for _, ip := range ips ***REMOVED***
		switch protocol ***REMOVED***
		case iana.ProtocolICMP:
			if ip.To4() != nil ***REMOVED***
				return netaddr(ip)
			***REMOVED***
		case iana.ProtocolIPv6ICMP:
			if ip.To16() != nil && ip.To4() == nil ***REMOVED***
				return netaddr(ip)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil, errors.New("no A or AAAA record")
***REMOVED***

type pingTest struct ***REMOVED***
	network, address string
	protocol         int
	mtype            icmp.Type
***REMOVED***

var nonPrivilegedPingTests = []pingTest***REMOVED***
	***REMOVED***"udp4", "0.0.0.0", iana.ProtocolICMP, ipv4.ICMPTypeEcho***REMOVED***,

	***REMOVED***"udp6", "::", iana.ProtocolIPv6ICMP, ipv6.ICMPTypeEchoRequest***REMOVED***,
***REMOVED***

func TestNonPrivilegedPing(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("avoid external network")
	***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "darwin":
	case "linux":
		t.Log("you may need to adjust the net.ipv4.ping_group_range kernel state")
	default:
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***

	for i, tt := range nonPrivilegedPingTests ***REMOVED***
		if err := doPing(tt, i); err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

var privilegedPingTests = []pingTest***REMOVED***
	***REMOVED***"ip4:icmp", "0.0.0.0", iana.ProtocolICMP, ipv4.ICMPTypeEcho***REMOVED***,

	***REMOVED***"ip6:ipv6-icmp", "::", iana.ProtocolIPv6ICMP, ipv6.ICMPTypeEchoRequest***REMOVED***,
***REMOVED***

func TestPrivilegedPing(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("avoid external network")
	***REMOVED***
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***

	for i, tt := range privilegedPingTests ***REMOVED***
		if err := doPing(tt, i); err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func doPing(tt pingTest, seq int) error ***REMOVED***
	c, err := icmp.ListenPacket(tt.network, tt.address)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer c.Close()

	dst, err := googleAddr(c, tt.protocol)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if tt.network != "udp6" && tt.protocol == iana.ProtocolIPv6ICMP ***REMOVED***
		var f ipv6.ICMPFilter
		f.SetAll(true)
		f.Accept(ipv6.ICMPTypeDestinationUnreachable)
		f.Accept(ipv6.ICMPTypePacketTooBig)
		f.Accept(ipv6.ICMPTypeTimeExceeded)
		f.Accept(ipv6.ICMPTypeParameterProblem)
		f.Accept(ipv6.ICMPTypeEchoReply)
		if err := c.IPv6PacketConn().SetICMPFilter(&f); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	wm := icmp.Message***REMOVED***
		Type: tt.mtype, Code: 0,
		Body: &icmp.Echo***REMOVED***
			ID: os.Getpid() & 0xffff, Seq: 1 << uint(seq),
			Data: []byte("HELLO-R-U-THERE"),
		***REMOVED***,
	***REMOVED***
	wb, err := wm.Marshal(nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n, err := c.WriteTo(wb, dst); err != nil ***REMOVED***
		return err
	***REMOVED*** else if n != len(wb) ***REMOVED***
		return fmt.Errorf("got %v; want %v", n, len(wb))
	***REMOVED***

	rb := make([]byte, 1500)
	if err := c.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil ***REMOVED***
		return err
	***REMOVED***
	n, peer, err := c.ReadFrom(rb)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	rm, err := icmp.ParseMessage(tt.protocol, rb[:n])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	switch rm.Type ***REMOVED***
	case ipv4.ICMPTypeEchoReply, ipv6.ICMPTypeEchoReply:
		return nil
	default:
		return fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	***REMOVED***
***REMOVED***

func TestConcurrentNonPrivilegedListenPacket(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("avoid external network")
	***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "darwin":
	case "linux":
		t.Log("you may need to adjust the net.ipv4.ping_group_range kernel state")
	default:
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***

	network, address := "udp4", "127.0.0.1"
	if !nettest.SupportsIPv4() ***REMOVED***
		network, address = "udp6", "::1"
	***REMOVED***
	const N = 1000
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ ***REMOVED***
		go func() ***REMOVED***
			defer wg.Done()
			c, err := icmp.ListenPacket(network, address)
			if err != nil ***REMOVED***
				t.Error(err)
				return
			***REMOVED***
			c.Close()
		***REMOVED***()
	***REMOVED***
	wg.Wait()
***REMOVED***
