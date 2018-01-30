// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6_test

import (
	"net"
	"reflect"
	"runtime"
	"testing"

	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/ipv6"
)

var icmpStringTests = []struct ***REMOVED***
	in  ipv6.ICMPType
	out string
***REMOVED******REMOVED***
	***REMOVED***ipv6.ICMPTypeDestinationUnreachable, "destination unreachable"***REMOVED***,

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
	case "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***

	var f ipv6.ICMPFilter
	for _, toggle := range []bool***REMOVED***false, true***REMOVED*** ***REMOVED***
		f.SetAll(toggle)
		for _, typ := range []ipv6.ICMPType***REMOVED***
			ipv6.ICMPTypeDestinationUnreachable,
			ipv6.ICMPTypeEchoReply,
			ipv6.ICMPTypeNeighborSolicitation,
			ipv6.ICMPTypeDuplicateAddressConfirmation,
		***REMOVED*** ***REMOVED***
			f.Accept(typ)
			if f.WillBlock(typ) ***REMOVED***
				t.Errorf("ipv6.ICMPFilter.Set(%v, false) failed", typ)
			***REMOVED***
			f.Block(typ)
			if !f.WillBlock(typ) ***REMOVED***
				t.Errorf("ipv6.ICMPFilter.Set(%v, true) failed", typ)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSetICMPFilter(t *testing.T) ***REMOVED***
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

	var f ipv6.ICMPFilter
	f.SetAll(true)
	f.Accept(ipv6.ICMPTypeEchoRequest)
	f.Accept(ipv6.ICMPTypeEchoReply)
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
