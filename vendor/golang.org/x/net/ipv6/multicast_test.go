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

var packetConnReadWriteMulticastUDPTests = []struct ***REMOVED***
	addr     string
	grp, src *net.UDPAddr
***REMOVED******REMOVED***
	***REMOVED***"[ff02::]:0", &net.UDPAddr***REMOVED***IP: net.ParseIP("ff02::114")***REMOVED***, nil***REMOVED***, // see RFC 4727

	***REMOVED***"[ff30::8000:0]:0", &net.UDPAddr***REMOVED***IP: net.ParseIP("ff30::8000:1")***REMOVED***, &net.UDPAddr***REMOVED***IP: net.IPv6loopback***REMOVED******REMOVED***, // see RFC 5771
***REMOVED***

func TestPacketConnReadWriteMulticastUDP(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if !supportsIPv6 ***REMOVED***
		t.Skip("ipv6 is not supported")
	***REMOVED***
	if !nettest.SupportsIPv6MulticastDeliveryOnLoopback() ***REMOVED***
		t.Skipf("multicast delivery doesn't work correctly on %s", runtime.GOOS)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip6", net.FlagUp|net.FlagMulticast|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	for _, tt := range packetConnReadWriteMulticastUDPTests ***REMOVED***
		c, err := net.ListenPacket("udp6", tt.addr)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c.Close()

		grp := *tt.grp
		grp.Port = c.LocalAddr().(*net.UDPAddr).Port
		p := ipv6.NewPacketConn(c)
		defer p.Close()
		if tt.src == nil ***REMOVED***
			if err := p.JoinGroup(ifi, &grp); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			defer p.LeaveGroup(ifi, &grp)
		***REMOVED*** else ***REMOVED***
			if err := p.JoinSourceSpecificGroup(ifi, &grp, tt.src); err != nil ***REMOVED***
				switch runtime.GOOS ***REMOVED***
				case "freebsd", "linux":
				default: // platforms that don't support MLDv2 fail here
					t.Logf("not supported on %s", runtime.GOOS)
					continue
				***REMOVED***
				t.Fatal(err)
			***REMOVED***
			defer p.LeaveSourceSpecificGroup(ifi, &grp, tt.src)
		***REMOVED***
		if err := p.SetMulticastInterface(ifi); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if _, err := p.MulticastInterface(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if err := p.SetMulticastLoopback(true); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if _, err := p.MulticastLoopback(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		cm := ipv6.ControlMessage***REMOVED***
			TrafficClass: iana.DiffServAF11 | iana.CongestionExperienced,
			Src:          net.IPv6loopback,
			IfIndex:      ifi.Index,
		***REMOVED***
		cf := ipv6.FlagTrafficClass | ipv6.FlagHopLimit | ipv6.FlagSrc | ipv6.FlagDst | ipv6.FlagInterface | ipv6.FlagPathMTU
		wb := []byte("HELLO-R-U-THERE")

		for i, toggle := range []bool***REMOVED***true, false, true***REMOVED*** ***REMOVED***
			if err := p.SetControlMessage(cf, toggle); err != nil ***REMOVED***
				if nettest.ProtocolNotSupported(err) ***REMOVED***
					t.Logf("not supported on %s", runtime.GOOS)
					continue
				***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if err := p.SetDeadline(time.Now().Add(200 * time.Millisecond)); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			cm.HopLimit = i + 1
			if n, err := p.WriteTo(wb, &cm, &grp); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED*** else if n != len(wb) ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			rb := make([]byte, 128)
			if n, _, _, err := p.ReadFrom(rb); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED*** else if !bytes.Equal(rb[:n], wb) ***REMOVED***
				t.Fatalf("got %v; want %v", rb[:n], wb)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var packetConnReadWriteMulticastICMPTests = []struct ***REMOVED***
	grp, src *net.IPAddr
***REMOVED******REMOVED***
	***REMOVED***&net.IPAddr***REMOVED***IP: net.ParseIP("ff02::114")***REMOVED***, nil***REMOVED***, // see RFC 4727

	***REMOVED***&net.IPAddr***REMOVED***IP: net.ParseIP("ff30::8000:1")***REMOVED***, &net.IPAddr***REMOVED***IP: net.IPv6loopback***REMOVED******REMOVED***, // see RFC 5771
***REMOVED***

func TestPacketConnReadWriteMulticastICMP(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if !supportsIPv6 ***REMOVED***
		t.Skip("ipv6 is not supported")
	***REMOVED***
	if !nettest.SupportsIPv6MulticastDeliveryOnLoopback() ***REMOVED***
		t.Skipf("multicast delivery doesn't work correctly on %s", runtime.GOOS)
	***REMOVED***
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip6", net.FlagUp|net.FlagMulticast|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	for _, tt := range packetConnReadWriteMulticastICMPTests ***REMOVED***
		c, err := net.ListenPacket("ip6:ipv6-icmp", "::")
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c.Close()

		pshicmp := icmp.IPv6PseudoHeader(c.LocalAddr().(*net.IPAddr).IP, tt.grp.IP)
		p := ipv6.NewPacketConn(c)
		defer p.Close()
		if tt.src == nil ***REMOVED***
			if err := p.JoinGroup(ifi, tt.grp); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			defer p.LeaveGroup(ifi, tt.grp)
		***REMOVED*** else ***REMOVED***
			if err := p.JoinSourceSpecificGroup(ifi, tt.grp, tt.src); err != nil ***REMOVED***
				switch runtime.GOOS ***REMOVED***
				case "freebsd", "linux":
				default: // platforms that don't support MLDv2 fail here
					t.Logf("not supported on %s", runtime.GOOS)
					continue
				***REMOVED***
				t.Fatal(err)
			***REMOVED***
			defer p.LeaveSourceSpecificGroup(ifi, tt.grp, tt.src)
		***REMOVED***
		if err := p.SetMulticastInterface(ifi); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if _, err := p.MulticastInterface(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if err := p.SetMulticastLoopback(true); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if _, err := p.MulticastLoopback(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		cm := ipv6.ControlMessage***REMOVED***
			TrafficClass: iana.DiffServAF11 | iana.CongestionExperienced,
			Src:          net.IPv6loopback,
			IfIndex:      ifi.Index,
		***REMOVED***
		cf := ipv6.FlagTrafficClass | ipv6.FlagHopLimit | ipv6.FlagSrc | ipv6.FlagDst | ipv6.FlagInterface | ipv6.FlagPathMTU

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
					// Solaris never allows to
					// modify ICMP properties.
					if runtime.GOOS != "solaris" ***REMOVED***
						t.Fatal(err)
					***REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				psh = pshicmp
				// Some platforms never allow to
				// disable the kernel checksum
				// processing.
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
			if err := p.SetDeadline(time.Now().Add(200 * time.Millisecond)); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			cm.HopLimit = i + 1
			if n, err := p.WriteTo(wb, &cm, tt.grp); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED*** else if n != len(wb) ***REMOVED***
				t.Fatalf("got %v; want %v", n, len(wb))
			***REMOVED***
			rb := make([]byte, 128)
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
***REMOVED***
