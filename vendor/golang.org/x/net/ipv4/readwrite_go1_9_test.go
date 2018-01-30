// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9

package ipv4_test

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
	"golang.org/x/net/ipv4"
)

func BenchmarkPacketConnReadWriteUnicast(b *testing.B) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		b.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***

	payload := []byte("HELLO-R-U-THERE")
	iph, err := (&ipv4.Header***REMOVED***
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TotalLen: ipv4.HeaderLen + len(payload),
		TTL:      1,
		Protocol: iana.ProtocolReserved,
		Src:      net.IPv4(192, 0, 2, 1),
		Dst:      net.IPv4(192, 0, 2, 254),
	***REMOVED***).Marshal()
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	greh := []byte***REMOVED***0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00***REMOVED***
	datagram := append(greh, append(iph, payload...)...)
	bb := make([]byte, 128)
	cm := ipv4.ControlMessage***REMOVED***
		Src: net.IPv4(127, 0, 0, 1),
	***REMOVED***
	if ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagLoopback); ifi != nil ***REMOVED***
		cm.IfIndex = ifi.Index
	***REMOVED***

	b.Run("UDP", func(b *testing.B) ***REMOVED***
		c, err := nettest.NewLocalPacketListener("udp4")
		if err != nil ***REMOVED***
			b.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
		***REMOVED***
		defer c.Close()
		p := ipv4.NewPacketConn(c)
		dst := c.LocalAddr()
		cf := ipv4.FlagTTL | ipv4.FlagInterface
		if err := p.SetControlMessage(cf, true); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		wms := []ipv4.Message***REMOVED***
			***REMOVED***
				Buffers: [][]byte***REMOVED***payload***REMOVED***,
				Addr:    dst,
				OOB:     cm.Marshal(),
			***REMOVED***,
		***REMOVED***
		rms := []ipv4.Message***REMOVED***
			***REMOVED***
				Buffers: [][]byte***REMOVED***bb***REMOVED***,
				OOB:     ipv4.NewControlMessage(cf),
			***REMOVED***,
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
		b.Run("Batch", func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				if _, err := p.WriteBatch(wms, 0); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
				if _, err := p.ReadBatch(rms, 0); err != nil ***REMOVED***
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

		c, err := net.ListenPacket(fmt.Sprintf("ip4:%d", iana.ProtocolGRE), "127.0.0.1")
		if err != nil ***REMOVED***
			b.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
		***REMOVED***
		defer c.Close()
		p := ipv4.NewPacketConn(c)
		dst := c.LocalAddr()
		cf := ipv4.FlagTTL | ipv4.FlagInterface
		if err := p.SetControlMessage(cf, true); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		wms := []ipv4.Message***REMOVED***
			***REMOVED***
				Buffers: [][]byte***REMOVED***datagram***REMOVED***,
				Addr:    dst,
				OOB:     cm.Marshal(),
			***REMOVED***,
		***REMOVED***
		rms := []ipv4.Message***REMOVED***
			***REMOVED***
				Buffers: [][]byte***REMOVED***bb***REMOVED***,
				OOB:     ipv4.NewControlMessage(cf),
			***REMOVED***,
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
		b.Run("Batch", func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				if _, err := p.WriteBatch(wms, 0); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
				if _, err := p.ReadBatch(rms, 0); err != nil ***REMOVED***
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
	iph, err := (&ipv4.Header***REMOVED***
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TotalLen: ipv4.HeaderLen + len(payload),
		TTL:      1,
		Protocol: iana.ProtocolReserved,
		Src:      net.IPv4(192, 0, 2, 1),
		Dst:      net.IPv4(192, 0, 2, 254),
	***REMOVED***).Marshal()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	greh := []byte***REMOVED***0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00***REMOVED***
	datagram := append(greh, append(iph, payload...)...)

	t.Run("UDP", func(t *testing.T) ***REMOVED***
		c, err := nettest.NewLocalPacketListener("udp4")
		if err != nil ***REMOVED***
			t.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
		***REMOVED***
		defer c.Close()
		p := ipv4.NewPacketConn(c)
		t.Run("ToFrom", func(t *testing.T) ***REMOVED***
			testPacketConnConcurrentReadWriteUnicast(t, p, payload, c.LocalAddr(), false)
		***REMOVED***)
		t.Run("Batch", func(t *testing.T) ***REMOVED***
			testPacketConnConcurrentReadWriteUnicast(t, p, payload, c.LocalAddr(), true)
		***REMOVED***)
	***REMOVED***)
	t.Run("IP", func(t *testing.T) ***REMOVED***
		switch runtime.GOOS ***REMOVED***
		case "netbsd":
			t.Skip("need to configure gre on netbsd")
		case "openbsd":
			t.Skip("net.inet.gre.allow=0 by default on openbsd")
		***REMOVED***

		c, err := net.ListenPacket(fmt.Sprintf("ip4:%d", iana.ProtocolGRE), "127.0.0.1")
		if err != nil ***REMOVED***
			t.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
		***REMOVED***
		defer c.Close()
		p := ipv4.NewPacketConn(c)
		t.Run("ToFrom", func(t *testing.T) ***REMOVED***
			testPacketConnConcurrentReadWriteUnicast(t, p, datagram, c.LocalAddr(), false)
		***REMOVED***)
		t.Run("Batch", func(t *testing.T) ***REMOVED***
			testPacketConnConcurrentReadWriteUnicast(t, p, datagram, c.LocalAddr(), true)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func testPacketConnConcurrentReadWriteUnicast(t *testing.T, p *ipv4.PacketConn, data []byte, dst net.Addr, batch bool) ***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagLoopback)
	cf := ipv4.FlagTTL | ipv4.FlagSrc | ipv4.FlagDst | ipv4.FlagInterface

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
	batchReader := func() ***REMOVED***
		defer wg.Done()
		ms := []ipv4.Message***REMOVED***
			***REMOVED***
				Buffers: [][]byte***REMOVED***make([]byte, 128)***REMOVED***,
				OOB:     ipv4.NewControlMessage(cf),
			***REMOVED***,
		***REMOVED***
		n, err := p.ReadBatch(ms, 0)
		if err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		if n != len(ms) ***REMOVED***
			t.Errorf("got %d; want %d", n, len(ms))
			return
		***REMOVED***
		var cm ipv4.ControlMessage
		if err := cm.Parse(ms[0].OOB[:ms[0].NN]); err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		var b []byte
		if _, ok := dst.(*net.IPAddr); ok ***REMOVED***
			var h ipv4.Header
			if err := h.Parse(ms[0].Buffers[0][:ms[0].N]); err != nil ***REMOVED***
				t.Error(err)
				return
			***REMOVED***
			b = ms[0].Buffers[0][h.Len:ms[0].N]
		***REMOVED*** else ***REMOVED***
			b = ms[0].Buffers[0][:ms[0].N]
		***REMOVED***
		if !bytes.Equal(b, data) ***REMOVED***
			t.Errorf("got %#v; want %#v", b, data)
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
		cm := ipv4.ControlMessage***REMOVED***
			Src: net.IPv4(127, 0, 0, 1),
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
	batchWriter := func(toggle bool) ***REMOVED***
		defer wg.Done()
		cm := ipv4.ControlMessage***REMOVED***
			Src: net.IPv4(127, 0, 0, 1),
		***REMOVED***
		if ifi != nil ***REMOVED***
			cm.IfIndex = ifi.Index
		***REMOVED***
		if err := p.SetControlMessage(cf, toggle); err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		ms := []ipv4.Message***REMOVED***
			***REMOVED***
				Buffers: [][]byte***REMOVED***data***REMOVED***,
				OOB:     cm.Marshal(),
				Addr:    dst,
			***REMOVED***,
		***REMOVED***
		n, err := p.WriteBatch(ms, 0)
		if err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		if n != len(ms) ***REMOVED***
			t.Errorf("got %d; want %d", n, len(ms))
			return
		***REMOVED***
		if ms[0].N != len(data) ***REMOVED***
			t.Errorf("got %d; want %d", ms[0].N, len(data))
			return
		***REMOVED***
	***REMOVED***

	const N = 10
	wg.Add(N)
	for i := 0; i < N; i++ ***REMOVED***
		if batch ***REMOVED***
			go batchReader()
		***REMOVED*** else ***REMOVED***
			go reader()
		***REMOVED***
	***REMOVED***
	wg.Add(2 * N)
	for i := 0; i < 2*N; i++ ***REMOVED***
		if batch ***REMOVED***
			go batchWriter(i%2 != 0)
		***REMOVED*** else ***REMOVED***
			go writer(i%2 != 0)
		***REMOVED***

	***REMOVED***
	wg.Add(N)
	for i := 0; i < N; i++ ***REMOVED***
		if batch ***REMOVED***
			go batchReader()
		***REMOVED*** else ***REMOVED***
			go reader()
		***REMOVED***
	***REMOVED***
	wg.Wait()
***REMOVED***
