// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6_test

import (
	"bytes"
	"net"
	"os"
	"runtime"
	"testing"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/ipv6"
)

func TestPacketConnReadWriteUnicastUDP(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if !supportsIPv6 ***REMOVED***
		t.Skip("ipv6 is not supported")
	***REMOVED***

	c, err := nettest.NewLocalPacketListener("udp6")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()
	p := ipv6.NewPacketConn(c)
	defer p.Close()

	dst := c.LocalAddr()
	cm := ipv6.ControlMessage***REMOVED***
		TrafficClass: iana.DiffServAF11 | iana.CongestionExperienced,
		Src:          net.IPv6loopback,
	***REMOVED***
	cf := ipv6.FlagTrafficClass | ipv6.FlagHopLimit | ipv6.FlagSrc | ipv6.FlagDst | ipv6.FlagInterface | ipv6.FlagPathMTU
	ifi := nettest.RoutedInterface("ip6", net.FlagUp|net.FlagLoopback)
	if ifi != nil ***REMOVED***
		cm.IfIndex = ifi.Index
	***REMOVED***
	wb := []byte("HELLO-R-U-THERE")

	for i, toggle := range []bool***REMOVED***true, false, true***REMOVED*** ***REMOVED***
		if err := p.SetControlMessage(cf, toggle); err != nil ***REMOVED***
			if nettest.ProtocolNotSupported(err) ***REMOVED***
				t.Logf("not supported on %s", runtime.GOOS)
				continue
			***REMOVED***
			t.Fatal(err)
		***REMOVED***
		cm.HopLimit = i + 1
		if err := p.SetWriteDeadline(time.Now().Add(100 * time.Millisecond)); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if n, err := p.WriteTo(wb, &cm, dst); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED*** else if n != len(wb) ***REMOVED***
			t.Fatalf("got %v; want %v", n, len(wb))
		***REMOVED***
		rb := make([]byte, 128)
		if err := p.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if n, _, _, err := p.ReadFrom(rb); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED*** else if !bytes.Equal(rb[:n], wb) ***REMOVED***
			t.Fatalf("got %v; want %v", rb[:n], wb)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestPacketConnReadWriteUnicastICMP(t *testing.T) ***REMOVED***
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

	c, err := net.ListenPacket("ip6:ipv6-icmp", "::1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()
	p := ipv6.NewPacketConn(c)
	defer p.Close()

	dst, err := net.ResolveIPAddr("ip6", "::1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	pshicmp := icmp.IPv6PseudoHeader(c.LocalAddr().(*net.IPAddr).IP, dst.IP)
	cm := ipv6.ControlMessage***REMOVED***
		TrafficClass: iana.DiffServAF11 | iana.CongestionExperienced,
		Src:          net.IPv6loopback,
	***REMOVED***
	cf := ipv6.FlagTrafficClass | ipv6.FlagHopLimit | ipv6.FlagSrc | ipv6.FlagDst | ipv6.FlagInterface | ipv6.FlagPathMTU
	ifi := nettest.RoutedInterface("ip6", net.FlagUp|net.FlagLoopback)
	if ifi != nil ***REMOVED***
		cm.IfIndex = ifi.Index
	***REMOVED***

	var f ipv6.ICMPFilter
	f.SetAll(true)
	f.Accept(ipv6.ICMPTypeEchoReply)
	if err := p.SetICMPFilter(&f); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	var psh []byte
	for i, toggle := range []bool***REMOVED***true, false, true***REMOVED*** ***REMOVED***
		if toggle ***REMOVED***
			psh = nil
			if err := p.SetChecksum(true, 2); err != nil ***REMOVED***
				// Solaris never allows to modify
				// ICMP properties.
				if runtime.GOOS != "solaris" ***REMOVED***
					t.Fatal(err)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			psh = pshicmp
			// Some platforms never allow to disable the
			// kernel checksum processing.
			p.SetChecksum(false, -1)
		***REMOVED***
		wb, err := (&icmp.Message***REMOVED***
			Type: ipv6.ICMPTypeEchoRequest, Code: 0,
			Body: &icmp.Echo***REMOVED***
				ID: os.Getpid() & 0xffff, Seq: i + 1,
				Data: []byte("HELLO-R-U-THERE"),
			***REMOVED***,
		***REMOVED***).Marshal(psh)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if err := p.SetControlMessage(cf, toggle); err != nil ***REMOVED***
			if nettest.ProtocolNotSupported(err) ***REMOVED***
				t.Logf("not supported on %s", runtime.GOOS)
				continue
			***REMOVED***
			t.Fatal(err)
		***REMOVED***
		cm.HopLimit = i + 1
		if err := p.SetWriteDeadline(time.Now().Add(100 * time.Millisecond)); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if n, err := p.WriteTo(wb, &cm, dst); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED*** else if n != len(wb) ***REMOVED***
			t.Fatalf("got %v; want %v", n, len(wb))
		***REMOVED***
		rb := make([]byte, 128)
		if err := p.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if n, _, _, err := p.ReadFrom(rb); err != nil ***REMOVED***
			switch runtime.GOOS ***REMOVED***
			case "darwin": // older darwin kernels have some limitation on receiving icmp packet through raw socket
				t.Logf("not supported on %s", runtime.GOOS)
				continue
			***REMOVED***
			t.Fatal(err)
		***REMOVED*** else ***REMOVED***
			if m, err := icmp.ParseMessage(iana.ProtocolIPv6ICMP, rb[:n]); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED*** else if m.Type != ipv6.ICMPTypeEchoReply || m.Code != 0 ***REMOVED***
				t.Fatalf("got type=%v, code=%v; want type=%v, code=%v", m.Type, m.Code, ipv6.ICMPTypeEchoReply, 0)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
