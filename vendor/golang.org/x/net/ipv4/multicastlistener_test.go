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

var udpMultipleGroupListenerTests = []net.Addr***REMOVED***
	&net.UDPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 249)***REMOVED***, // see RFC 4727
	&net.UDPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 250)***REMOVED***,
	&net.UDPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 254)***REMOVED***,
***REMOVED***

func TestUDPSinglePacketConnWithMultipleGroupListeners(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("to avoid external network")
	***REMOVED***

	for _, gaddr := range udpMultipleGroupListenerTests ***REMOVED***
		c, err := net.ListenPacket("udp4", "0.0.0.0:0") // wildcard address with no reusable port
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c.Close()

		p := ipv4.NewPacketConn(c)
		var mift []*net.Interface

		ift, err := net.Interfaces()
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		for i, ifi := range ift ***REMOVED***
			if _, ok := nettest.IsMulticastCapable("ip4", &ifi); !ok ***REMOVED***
				continue
			***REMOVED***
			if err := p.JoinGroup(&ifi, gaddr); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			mift = append(mift, &ift[i])
		***REMOVED***
		for _, ifi := range mift ***REMOVED***
			if err := p.LeaveGroup(ifi, gaddr); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUDPMultiplePacketConnWithMultipleGroupListeners(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("to avoid external network")
	***REMOVED***

	for _, gaddr := range udpMultipleGroupListenerTests ***REMOVED***
		c1, err := net.ListenPacket("udp4", "224.0.0.0:0") // wildcard address with reusable port
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c1.Close()
		_, port, err := net.SplitHostPort(c1.LocalAddr().String())
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		c2, err := net.ListenPacket("udp4", net.JoinHostPort("224.0.0.0", port)) // wildcard address with reusable port
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c2.Close()

		var ps [2]*ipv4.PacketConn
		ps[0] = ipv4.NewPacketConn(c1)
		ps[1] = ipv4.NewPacketConn(c2)
		var mift []*net.Interface

		ift, err := net.Interfaces()
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		for i, ifi := range ift ***REMOVED***
			if _, ok := nettest.IsMulticastCapable("ip4", &ifi); !ok ***REMOVED***
				continue
			***REMOVED***
			for _, p := range ps ***REMOVED***
				if err := p.JoinGroup(&ifi, gaddr); err != nil ***REMOVED***
					t.Fatal(err)
				***REMOVED***
			***REMOVED***
			mift = append(mift, &ift[i])
		***REMOVED***
		for _, ifi := range mift ***REMOVED***
			for _, p := range ps ***REMOVED***
				if err := p.LeaveGroup(ifi, gaddr); err != nil ***REMOVED***
					t.Fatal(err)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUDPPerInterfaceSinglePacketConnWithSingleGroupListener(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("to avoid external network")
	***REMOVED***

	gaddr := net.IPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 254)***REMOVED*** // see RFC 4727
	type ml struct ***REMOVED***
		c   *ipv4.PacketConn
		ifi *net.Interface
	***REMOVED***
	var mlt []*ml

	ift, err := net.Interfaces()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	port := "0"
	for i, ifi := range ift ***REMOVED***
		ip, ok := nettest.IsMulticastCapable("ip4", &ifi)
		if !ok ***REMOVED***
			continue
		***REMOVED***
		c, err := net.ListenPacket("udp4", net.JoinHostPort(ip.String(), port)) // unicast address with non-reusable port
		if err != nil ***REMOVED***
			// The listen may fail when the serivce is
			// already in use, but it's fine because the
			// purpose of this is not to test the
			// bookkeeping of IP control block inside the
			// kernel.
			t.Log(err)
			continue
		***REMOVED***
		defer c.Close()
		if port == "0" ***REMOVED***
			_, port, err = net.SplitHostPort(c.LocalAddr().String())
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
		p := ipv4.NewPacketConn(c)
		if err := p.JoinGroup(&ifi, &gaddr); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		mlt = append(mlt, &ml***REMOVED***p, &ift[i]***REMOVED***)
	***REMOVED***
	for _, m := range mlt ***REMOVED***
		if err := m.c.LeaveGroup(m.ifi, &gaddr); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIPSingleRawConnWithSingleGroupListener(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("to avoid external network")
	***REMOVED***
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***

	c, err := net.ListenPacket("ip4:icmp", "0.0.0.0") // wildcard address
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()

	r, err := ipv4.NewRawConn(c)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	gaddr := net.IPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 254)***REMOVED*** // see RFC 4727
	var mift []*net.Interface

	ift, err := net.Interfaces()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	for i, ifi := range ift ***REMOVED***
		if _, ok := nettest.IsMulticastCapable("ip4", &ifi); !ok ***REMOVED***
			continue
		***REMOVED***
		if err := r.JoinGroup(&ifi, &gaddr); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		mift = append(mift, &ift[i])
	***REMOVED***
	for _, ifi := range mift ***REMOVED***
		if err := r.LeaveGroup(ifi, &gaddr); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIPPerInterfaceSingleRawConnWithSingleGroupListener(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("to avoid external network")
	***REMOVED***
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***

	gaddr := net.IPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 254)***REMOVED*** // see RFC 4727
	type ml struct ***REMOVED***
		c   *ipv4.RawConn
		ifi *net.Interface
	***REMOVED***
	var mlt []*ml

	ift, err := net.Interfaces()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	for i, ifi := range ift ***REMOVED***
		ip, ok := nettest.IsMulticastCapable("ip4", &ifi)
		if !ok ***REMOVED***
			continue
		***REMOVED***
		c, err := net.ListenPacket("ip4:253", ip.String()) // unicast address
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer c.Close()
		r, err := ipv4.NewRawConn(c)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if err := r.JoinGroup(&ifi, &gaddr); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		mlt = append(mlt, &ml***REMOVED***r, &ift[i]***REMOVED***)
	***REMOVED***
	for _, m := range mlt ***REMOVED***
		if err := m.c.LeaveGroup(m.ifi, &gaddr); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
