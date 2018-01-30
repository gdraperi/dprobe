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

func TestPacketConnReadWriteUnicastUDP(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	c, err := nettest.NewLocalPacketListener("udp4")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()
	p := ipv4.NewPacketConn(c)
	defer p.Close()

	dst := c.LocalAddr()
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
		p.SetTTL(i + 1)
		if err := p.SetWriteDeadline(time.Now().Add(100 * time.Millisecond)); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if n, err := p.WriteTo(wb, nil, dst); err != nil ***REMOVED***
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
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	c, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()

	dst, err := net.ResolveIPAddr("ip4", "127.0.0.1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	p := ipv4.NewPacketConn(c)
	defer p.Close()
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
		p.SetTTL(i + 1)
		if err := p.SetWriteDeadline(time.Now().Add(100 * time.Millisecond)); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if n, err := p.WriteTo(wb, nil, dst); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED*** else if n != len(wb) ***REMOVED***
			t.Fatalf("got %v; want %v", n, len(wb))
		***REMOVED***
		rb := make([]byte, 128)
	loop:
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
			m, err := icmp.ParseMessage(iana.ProtocolICMP, rb[:n])
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if runtime.GOOS == "linux" && m.Type == ipv4.ICMPTypeEcho ***REMOVED***
				// On Linux we must handle own sent packets.
				goto loop
			***REMOVED***
			if m.Type != ipv4.ICMPTypeEchoReply || m.Code != 0 ***REMOVED***
				t.Fatalf("got type=%v, code=%v; want type=%v, code=%v", m.Type, m.Code, ipv4.ICMPTypeEchoReply, 0)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRawConnReadWriteUnicastICMP(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***
	ifi := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagLoopback)
	if ifi == nil ***REMOVED***
		t.Skipf("not available on %s", runtime.GOOS)
	***REMOVED***

	c, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()

	dst, err := net.ResolveIPAddr("ip4", "127.0.0.1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	r, err := ipv4.NewRawConn(c)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer r.Close()
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
			TTL:      i + 1,
			Protocol: 1,
			Dst:      dst.IP,
		***REMOVED***
		if err := r.SetControlMessage(cf, toggle); err != nil ***REMOVED***
			if nettest.ProtocolNotSupported(err) ***REMOVED***
				t.Logf("not supported on %s", runtime.GOOS)
				continue
			***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if err := r.SetWriteDeadline(time.Now().Add(100 * time.Millisecond)); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if err := r.WriteTo(wh, wb, nil); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		rb := make([]byte, ipv4.HeaderLen+128)
	loop:
		if err := r.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if _, b, _, err := r.ReadFrom(rb); err != nil ***REMOVED***
			switch runtime.GOOS ***REMOVED***
			case "darwin": // older darwin kernels have some limitation on receiving icmp packet through raw socket
				t.Logf("not supported on %s", runtime.GOOS)
				continue
			***REMOVED***
			t.Fatal(err)
		***REMOVED*** else ***REMOVED***
			m, err := icmp.ParseMessage(iana.ProtocolICMP, b)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if runtime.GOOS == "linux" && m.Type == ipv4.ICMPTypeEcho ***REMOVED***
				// On Linux we must handle own sent packets.
				goto loop
			***REMOVED***
			if m.Type != ipv4.ICMPTypeEchoReply || m.Code != 0 ***REMOVED***
				t.Fatalf("got type=%v, code=%v; want type=%v, code=%v", m.Type, m.Code, ipv4.ICMPTypeEchoReply, 0)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
