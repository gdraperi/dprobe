// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.9

package ipv6_test

import (
	"bytes"
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"
	"testing"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/ipv6"
)

func BenchmarkPacketConnReadWriteUnicast(b *testing.B) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		b.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***

	payload := []byte("HELLO-R-U-THERE")
	iph := []byte***REMOVED***
		0x69, 0x8b, 0xee, 0xf1, 0xca, 0xfe, 0xff, 0x01,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x02, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	***REMOVED***
	greh := []byte***REMOVED***0x00, 0x00, 0x86, 0xdd, 0x00, 0x00, 0x00, 0x00***REMOVED***
	datagram := append(greh, append(iph, payload...)...)
	bb := make([]byte, 128)
	cm := ipv6.ControlMessage***REMOVED***
		TrafficClass: iana.DiffServAF11 | iana.CongestionExperienced,
		HopLimit:     1,
		Src:          net.IPv6loopback,
	***REMOVED***
	if ifi := nettest.RoutedInterface("ip6", net.FlagUp|net.FlagLoopback); ifi != nil ***REMOVED***
		cm.IfIndex = ifi.Index
	***REMOVED***

	b.Run("UDP", func(b *testing.B) ***REMOVED***
		c, err := nettest.NewLocalPacketListener("udp6")
		if err != nil ***REMOVED***
			b.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
		***REMOVED***
		defer c.Close()
		p := ipv6.NewPacketConn(c)
		dst := c.LocalAddr()
		cf := ipv6.FlagHopLimit | ipv6.FlagInterface
		if err := p.SetControlMessage(cf, true); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		b.Run("Net", func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				if _, err := c.WriteTo(payload, dst); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
				if _, _, err := c.ReadFrom(bb); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
		b.Run("ToFrom", func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				if _, err := p.WriteTo(payload, &cm, dst); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
				if _, _, _, err := p.ReadFrom(bb); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	b.Run("IP", func(b *testing.B) ***REMOVED***
		switch runtime.GOOS ***REMOVED***
		case "netbsd":
			b.Skip("need to configure gre on netbsd")
		case "openbsd":
			b.Skip("net.inet.gre.allow=0 by default on openbsd")
		***REMOVED***

		c, err := net.ListenPacket(fmt.Sprintf("ip6:%d", iana.ProtocolGRE), "::1")
		if err != nil ***REMOVED***
			b.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
		***REMOVED***
		defer c.Close()
		p := ipv6.NewPacketConn(c)
		dst := c.LocalAddr()
		cf := ipv6.FlagTrafficClass | ipv6.FlagHopLimit | ipv6.FlagSrc | ipv6.FlagDst | ipv6.FlagInterface | ipv6.FlagPathMTU
		if err := p.SetControlMessage(cf, true); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		b.Run("Net", func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				if _, err := c.WriteTo(datagram, dst); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
				if _, _, err := c.ReadFrom(bb); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
		b.Run("ToFrom", func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				if _, err := p.WriteTo(datagram, &cm, dst); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
				if _, _, _, err := p.ReadFrom(bb); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestPacketConnConcurrentReadWriteUnicast(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***

	payload := []byte("HELLO-R-U-THERE")
	iph := []byte***REMOVED***
		0x69, 0x8b, 0xee, 0xf1, 0xca, 0xfe, 0xff, 0x01,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x02, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	***REMOVED***
	greh := []byte***REMOVED***0x00, 0x00, 0x86, 0xdd, 0x00, 0x00, 0x00, 0x00***REMOVED***
	datagram := append(greh, append(iph, payload...)...)

	t.Run("UDP", func(t *testing.T) ***REMOVED***
		c, err := nettest.NewLocalPacketListener("udp6")
		if err != nil ***REMOVED***
			t.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
		***REMOVED***
		defer c.Close()
		p := ipv6.NewPacketConn(c)
		t.Run("ToFrom", func(t *testing.T) ***REMOVED***
			testPacketConnConcurrentReadWriteUnicast(t, p, payload, c.LocalAddr())
		***REMOVED***)
	***REMOVED***)
	t.Run("IP", func(t *testing.T) ***REMOVED***
		switch runtime.GOOS ***REMOVED***
		case "netbsd":
			t.Skip("need to configure gre on netbsd")
		case "openbsd":
			t.Skip("net.inet.gre.allow=0 by default on openbsd")
		***REMOVED***

		c, err := net.ListenPacket(fmt.Sprintf("ip6:%d", iana.ProtocolGRE), "::1")
		if err != nil ***REMOVED***
			t.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
		***REMOVED***
		defer c.Close()
		p := ipv6.NewPacketConn(c)
		t.Run("ToFrom", func(t *testing.T) ***REMOVED***
			testPacketConnConcurrentReadWriteUnicast(t, p, datagram, c.LocalAddr())
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func testPacketConnConcurrentReadWriteUnicast(t *testing.T, p *ipv6.PacketConn, data []byte, dst net.Addr) ***REMOVED***
	ifi := nettest.RoutedInterface("ip6", net.FlagUp|net.FlagLoopback)
	cf := ipv6.FlagTrafficClass | ipv6.FlagHopLimit | ipv6.FlagSrc | ipv6.FlagDst | ipv6.FlagInterface | ipv6.FlagPathMTU

	if err := p.SetControlMessage(cf, true); err != nil ***REMOVED*** // probe before test
		if nettest.ProtocolNotSupported(err) ***REMOVED***
			t.Skipf("not supported on %s", runtime.GOOS)
		***REMOVED***
		t.Fatal(err)
	***REMOVED***

	var wg sync.WaitGroup
	reader := func() ***REMOVED***
		defer wg.Done()
		b := make([]byte, 128)
		n, cm, _, err := p.ReadFrom(b)
		if err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		if !bytes.Equal(b[:n], data) ***REMOVED***
			t.Errorf("got %#v; want %#v", b[:n], data)
			return
		***REMOVED***
		s := cm.String()
		if strings.Contains(s, ",") ***REMOVED***
			t.Errorf("should be space-separated values: %s", s)
			return
		***REMOVED***
	***REMOVED***
	writer := func(toggle bool) ***REMOVED***
		defer wg.Done()
		cm := ipv6.ControlMessage***REMOVED***
			TrafficClass: iana.DiffServAF11 | iana.CongestionExperienced,
			HopLimit:     1,
			Src:          net.IPv6loopback,
		***REMOVED***
		if ifi != nil ***REMOVED***
			cm.IfIndex = ifi.Index
		***REMOVED***
		if err := p.SetControlMessage(cf, toggle); err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		n, err := p.WriteTo(data, &cm, dst)
		if err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		if n != len(data) ***REMOVED***
			t.Errorf("got %d; want %d", n, len(data))
			return
		***REMOVED***
	***REMOVED***

	const N = 10
	wg.Add(N)
	for i := 0; i < N; i++ ***REMOVED***
		go reader()
	***REMOVED***
	wg.Add(2 * N)
	for i := 0; i < 2*N; i++ ***REMOVED***
		go writer(i%2 != 0)

	***REMOVED***
	wg.Add(N)
	for i := 0; i < N; i++ ***REMOVED***
		go reader()
	***REMOVED***
	wg.Wait()
***REMOVED***
