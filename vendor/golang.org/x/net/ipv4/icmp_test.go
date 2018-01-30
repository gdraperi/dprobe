// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4_test

import (
	"net"
	"reflect"
	"runtime"
	"testing"

	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/ipv4"
)

var icmpStringTests = []struct ***REMOVED***
	in  ipv4.ICMPType
	out string
***REMOVED******REMOVED***
	***REMOVED***ipv4.ICMPTypeDestinationUnreachable, "destination unreachable"***REMOVED***,

	***REMOVED***256, "<nil>"***REMOVED***,
***REMOVED***

func TestICMPString(t *testing.T) ***REMOVED***
	for _, tt := range icmpStringTests ***REMOVED***
		s := tt.in.String()
		if s != tt.out ***REMOVED***
			t.Errorf("got %s; want %s", s, tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestICMPFilter(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "linux":
	default:
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***

	var f ipv4.ICMPFilter
	for _, toggle := range []bool***REMOVED***false, true***REMOVED*** ***REMOVED***
		f.SetAll(toggle)
		for _, typ := range []ipv4.ICMPType***REMOVED***
			ipv4.ICMPTypeDestinationUnreachable,
			ipv4.ICMPTypeEchoReply,
			ipv4.ICMPTypeTimeExceeded,
			ipv4.ICMPTypeParameterProblem,
		***REMOVED*** ***REMOVED***
			f.Accept(typ)
			if f.WillBlock(typ) ***REMOVED***
				t.Errorf("ipv4.ICMPFilter.Set(%v, false) failed", typ)
			***REMOVED***
			f.Block(typ)
			if !f.WillBlock(typ) ***REMOVED***
				t.Errorf("ipv4.ICMPFilter.Set(%v, true) failed", typ)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSetICMPFilter(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "linux":
	default:
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if m, ok := nettest.SupportsRawIPSocket(); !ok ***REMOVED***
		t.Skip(m)
	***REMOVED***

	c, err := net.ListenPacket("ip4:icmp", "127.0.0.1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer c.Close()

	p := ipv4.NewPacketConn(c)

	var f ipv4.ICMPFilter
	f.SetAll(true)
	f.Accept(ipv4.ICMPTypeEcho)
	f.Accept(ipv4.ICMPTypeEchoReply)
	if err := p.SetICMPFilter(&f); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	kf, err := p.ICMPFilter()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(kf, &f) ***REMOVED***
		t.Fatalf("got %#v; want %#v", kf, f)
	***REMOVED***
***REMOVED***
