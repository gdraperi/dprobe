// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6_test

import (
	"bytes"
	"net"
	"runtime"
	"strings"
	"sync"
	"testing"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/ipv6"
)

func BenchmarkReadWriteUnicast(b *testing.B) ***REMOVED***
	c, err := nettest.NewLocalPacketListener("udp6")
	if err != nil ***REMOVED***
		b.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
	***REMOVED***
	defer c.Close()

	dst := c.LocalAddr()
	wb, rb := []byte("HELLO-R-U-THERE"), make([]byte, 128)

	b.Run("NetUDP", func(b *testing.B) ***REMOVED***
		for i := 0; i < b.N; i++ ***REMOVED***
			if _, err := c.WriteTo(wb, dst); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
			if _, _, err := c.ReadFrom(rb); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	b.Run("IPv6UDP", func(b *testing.B) ***REMOVED***
		p := ipv6.NewPacketConn(c)
		cf := ipv6.FlagTrafficClass | ipv6.FlagHopLimit | ipv6.FlagSrc | ipv6.FlagDst | ipv6.FlagInterface | ipv6.FlagPathMTU
		if err := p.SetControlMessage(cf, true); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		cm := ipv6.ControlMessage***REMOVED***
			TrafficClass: iana.DiffServAF11 | iana.CongestionExperienced,
			HopLimit:     1,
		***REMOVED***
		ifi := nettest.RoutedInterface("ip6", net.FlagUp|net.FlagLoopback)
		if ifi != nil ***REMOVED***
			cm.IfIndex = ifi.Index
		***REMOVED***

		for i := 0; i < b.N; i++ ***REMOVED***
			if _, err := p.WriteTo(wb, &cm, dst); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
			if _, _, _, err := p.ReadFrom(rb); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestPacketConnConcurrentReadWriteUnicastUDP(t *testing.T) ***REMOVED***
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
	ifi := nettest.RoutedInterface("ip6", net.FlagUp|net.FlagLoopback)
	cf := ipv6.FlagTrafficClass | ipv6.FlagHopLimit | ipv6.FlagSrc | ipv6.FlagDst | ipv6.FlagInterface | ipv6.FlagPathMTU
	wb := []byte("HELLO-R-U-THERE")

	if err := p.SetControlMessage(cf, true); err != nil ***REMOVED*** // probe before test
		if nettest.ProtocolNotSupported(err) ***REMOVED***
			t.Skipf("not supported on %s", runtime.GOOS)
		***REMOVED***
		t.Fatal(err)
	***REMOVED***

	var wg sync.WaitGroup
	reader := func() ***REMOVED***
		defer wg.Done()
		rb := make([]byte, 128)
		if n, cm, _, err := p.ReadFrom(rb); err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED*** else if !bytes.Equal(rb[:n], wb) ***REMOVED***
			t.Errorf("got %v; want %v", rb[:n], wb)
			return
		***REMOVED*** else ***REMOVED***
			s := cm.String()
			if strings.Contains(s, ",") ***REMOVED***
				t.Errorf("should be space-separated values: %s", s)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	writer := func(toggle bool) ***REMOVED***
		defer wg.Done()
		cm := ipv6.ControlMessage***REMOVED***
			TrafficClass: iana.DiffServAF11 | iana.CongestionExperienced,
			Src:          net.IPv6loopback,
		***REMOVED***
		if ifi != nil ***REMOVED***
			cm.IfIndex = ifi.Index
		***REMOVED***
		if err := p.SetControlMessage(cf, toggle); err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		if n, err := p.WriteTo(wb, &cm, dst); err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED*** else if n != len(wb) ***REMOVED***
			t.Errorf("got %d; want %d", n, len(wb))
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
