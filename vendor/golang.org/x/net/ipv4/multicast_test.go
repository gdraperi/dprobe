// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4_test

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
	"golang.org/x/net/ipv4"
)

var packetConnReadWriteMulticastUDPTests = []struct ***REMOVED***
	addr     string
	grp, src *net.UDPAddr
***REMOVED******REMOVED***
	***REMOVED***"224.0.0.0:0", &net.UDPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 254)***REMOVED***, nil***REMOVED***, // see RFC 4727

	***REMOVED***"232.0.1.0:0", &net.UDPAddr***REMOVED***IP: net.IPv4(232, 0, 1, 254)***REMOVED***, &net.UDPAddr***REMOVED***IP: net.IPv4(127, 0, 0, 1)***REMOVED******REMOVED***, // see RFC 5771
***REMOVED***

func TestPacketConnReadWriteMulticastUDP(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "solaris", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagMulticast|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	for _, tt := range packetConnReadWriteMulticastUDPTests ***REMOVED***
		c, err := net.ListenPacket("udp4", tt.addr)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c.Close()

		grp := *tt.grp
		grp.Port = c.LocalAddr().(*net.UDPAddr).Port
		p := ipv4.NewPacketConn(c)
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
				default: // platforms that don't support IGMPv2/3 fail here
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
		cf := ipv4.FlagTTL | ipv4.FlagDst | ipv4.FlagInterface
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
			p.SetMulticastTTL(i + 1)
			if n, err := p.WriteTo(wb, nil, &grp); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED*** else if n != len(wb) ***REMOVED***
				t.Fatalf("got %v; want %v", n, len(wb))
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
	***REMOVED***&net.IPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 254)***REMOVED***, nil***REMOVED***, // see RFC 4727

	***REMOVED***&net.IPAddr***REMOVED***IP: net.IPv4(232, 0, 1, 254)***REMOVED***, &net.IPAddr***REMOVED***IP: net.IPv4(127, 0, 0, 1)***REMOVED******REMOVED***, // see RFC 5771
***REMOVED***

func TestPacketConnReadWriteMulticastICMP(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "solaris", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagMulticast|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	for _, tt := range packetConnReadWriteMulticastICMPTests ***REMOVED***
		c, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c.Close()

		p := ipv4.NewPacketConn(c)
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
				default: // platforms that don't support IGMPv2/3 fail here
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
		cf := ipv4.FlagDst | ipv4.FlagInterface
		if runtime.GOOS != "solaris" ***REMOVED***
			// Solaris never allows to modify ICMP properties.
			cf |= ipv4.FlagTTL
		***REMOVED***

		for i, toggle := range []bool***REMOVED***true, false, true***REMOVED*** ***REMOVED***
			wb, err := (&icmp.Message***REMOVED***
				Type: ipv4.ICMPTypeEcho, Code: 0,
				Body: &icmp.Echo***REMOVED***
					ID: os.Getpid() & 0xffff, Seq: i + 1,
					Data: []byte("HELLO-R-U-THERE"),
				***REMOVED***,
			***REMOVED***).Marshal(nil)
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
			p.SetMulticastTTL(i + 1)
			if n, err := p.WriteTo(wb, nil, tt.grp); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED*** else if n != len(wb) ***REMOVED***
				t.Fatalf("got %v; want %v", n, len(wb))
			***REMOVED***
			rb := make([]byte, 128)
			if n, _, _, err := p.ReadFrom(rb); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED*** else ***REMOVED***
				m, err := icmp.ParseMessage(iana.ProtocolICMP, rb[:n])
				if err != nil ***REMOVED***
					t.Fatal(err)
				***REMOVED***
				switch ***REMOVED***
				case m.Type == ipv4.ICMPTypeEchoReply && m.Code == 0: // net.inet.icmp.bmcastecho=1
				case m.Type == ipv4.ICMPTypeEcho && m.Code == 0: // net.inet.icmp.bmcastecho=0
				default:
					t.Fatalf("got type=%v, code=%v; want type=%v, code=%v", m.Type, m.Code, ipv4.ICMPTypeEchoReply, 0)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var rawConnReadWriteMulticastICMPTests = []struct ***REMOVED***
	grp, src *net.IPAddr
***REMOVED******REMOVED***
	***REMOVED***&net.IPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 254)***REMOVED***, nil***REMOVED***, // see RFC 4727

	***REMOVED***&net.IPAddr***REMOVED***IP: net.IPv4(232, 0, 1, 254)***REMOVED***, &net.IPAddr***REMOVED***IP: net.IPv4(127, 0, 0, 1)***REMOVED******REMOVED***, // see RFC 5771
***REMOVED***

func TestRawConnReadWriteMulticastICMP(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("to avoid external network")
	***REMOVED***
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagMulticast|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	for _, tt := range rawConnReadWriteMulticastICMPTests ***REMOVED***
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
			if err := r.JoinGroup(ifi, tt.grp); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			defer r.LeaveGroup(ifi, tt.grp)
		***REMOVED*** else ***REMOVED***
			if err := r.JoinSourceSpecificGroup(ifi, tt.grp, tt.src); err != nil ***REMOVED***
				switch runtime.GOOS ***REMOVED***
				case "freebsd", "linux":
				default: // platforms that don't support IGMPv2/3 fail here
					t.Logf("not supported on %s", runtime.GOOS)
					continue
				***REMOVED***
				t.Fatal(err)
			***REMOVED***
			defer r.LeaveSourceSpecificGroup(ifi, tt.grp, tt.src)
		***REMOVED***
		if err := r.SetMulticastInterface(ifi); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if _, err := r.MulticastInterface(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if err := r.SetMulticastLoopback(true); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if _, err := r.MulticastLoopback(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		cf := ipv4.FlagTTL | ipv4.FlagDst | ipv4.FlagInterface

		for i, toggle := range []bool***REMOVED***true, false, true***REMOVED*** ***REMOVED***
			wb, err := (&icmp.Message***REMOVED***
				Type: ipv4.ICMPTypeEcho, Code: 0,
				Body: &icmp.Echo***REMOVED***
					ID: os.Getpid() & 0xffff, Seq: i + 1,
					Data: []byte("HELLO-R-U-THERE"),
				***REMOVED***,
			***REMOVED***).Marshal(nil)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			wh := &ipv4.Header***REMOVED***
				Version:  ipv4.Version,
				Len:      ipv4.HeaderLen,
				TOS:      i + 1,
				TotalLen: ipv4.HeaderLen + len(wb),
				Protocol: 1,
				Dst:      tt.grp.IP,
			***REMOVED***
			if err := r.SetControlMessage(cf, toggle); err != nil ***REMOVED***
				if nettest.ProtocolNotSupported(err) ***REMOVED***
					t.Logf("not supported on %s", runtime.GOOS)
					continue
				***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if err := r.SetDeadline(time.Now().Add(200 * time.Millisecond)); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			r.SetMulticastTTL(i + 1)
			if err := r.WriteTo(wh, wb, nil); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			rb := make([]byte, ipv4.HeaderLen+128)
			if rh, b, _, err := r.ReadFrom(rb); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED*** else ***REMOVED***
				m, err := icmp.ParseMessage(iana.ProtocolICMP, b)
				if err != nil ***REMOVED***
					t.Fatal(err)
				***REMOVED***
				switch ***REMOVED***
				case (rh.Dst.IsLoopback() || rh.Dst.IsLinkLocalUnicast() || rh.Dst.IsGlobalUnicast()) && m.Type == ipv4.ICMPTypeEchoReply && m.Code == 0: // net.inet.icmp.bmcastecho=1
				case rh.Dst.IsMulticast() && m.Type == ipv4.ICMPTypeEcho && m.Code == 0: // net.inet.icmp.bmcastecho=0
				default:
					t.Fatalf("got type=%v, code=%v; want type=%v, code=%v", m.Type, m.Code, ipv4.ICMPTypeEchoReply, 0)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
