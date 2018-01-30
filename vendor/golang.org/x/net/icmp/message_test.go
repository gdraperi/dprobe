// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp_test

import (
	"net"
	"reflect"
	"testing"

	"golang.org/x/net/icmp"
	"golang.org/x/net/internal/iana"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var marshalAndParseMessageForIPv4Tests = []icmp.Message***REMOVED***
	***REMOVED***
		Type: ipv4.ICMPTypeDestinationUnreachable, Code: 15,
		Body: &icmp.DstUnreach***REMOVED***
			Data: []byte("ERROR-INVOKING-PACKET"),
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv4.ICMPTypeTimeExceeded, Code: 1,
		Body: &icmp.TimeExceeded***REMOVED***
			Data: []byte("ERROR-INVOKING-PACKET"),
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv4.ICMPTypeParameterProblem, Code: 2,
		Body: &icmp.ParamProb***REMOVED***
			Pointer: 8,
			Data:    []byte("ERROR-INVOKING-PACKET"),
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo***REMOVED***
			ID: 1, Seq: 2,
			Data: []byte("HELLO-R-U-THERE"),
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv4.ICMPTypePhoturis,
		Body: &icmp.DefaultMessageBody***REMOVED***
			Data: []byte***REMOVED***0x80, 0x40, 0x20, 0x10***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestMarshalAndParseMessageForIPv4(t *testing.T) ***REMOVED***
	for i, tt := range marshalAndParseMessageForIPv4Tests ***REMOVED***
		b, err := tt.Marshal(nil)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		m, err := icmp.ParseMessage(iana.ProtocolICMP, b)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if m.Type != tt.Type || m.Code != tt.Code ***REMOVED***
			t.Errorf("#%v: got %v; want %v", i, m, &tt)
		***REMOVED***
		if !reflect.DeepEqual(m.Body, tt.Body) ***REMOVED***
			t.Errorf("#%v: got %v; want %v", i, m.Body, tt.Body)
		***REMOVED***
	***REMOVED***
***REMOVED***

var marshalAndParseMessageForIPv6Tests = []icmp.Message***REMOVED***
	***REMOVED***
		Type: ipv6.ICMPTypeDestinationUnreachable, Code: 6,
		Body: &icmp.DstUnreach***REMOVED***
			Data: []byte("ERROR-INVOKING-PACKET"),
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv6.ICMPTypePacketTooBig, Code: 0,
		Body: &icmp.PacketTooBig***REMOVED***
			MTU:  1<<16 - 1,
			Data: []byte("ERROR-INVOKING-PACKET"),
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv6.ICMPTypeTimeExceeded, Code: 1,
		Body: &icmp.TimeExceeded***REMOVED***
			Data: []byte("ERROR-INVOKING-PACKET"),
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv6.ICMPTypeParameterProblem, Code: 2,
		Body: &icmp.ParamProb***REMOVED***
			Pointer: 8,
			Data:    []byte("ERROR-INVOKING-PACKET"),
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv6.ICMPTypeEchoRequest, Code: 0,
		Body: &icmp.Echo***REMOVED***
			ID: 1, Seq: 2,
			Data: []byte("HELLO-R-U-THERE"),
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv6.ICMPTypeDuplicateAddressConfirmation,
		Body: &icmp.DefaultMessageBody***REMOVED***
			Data: []byte***REMOVED***0x80, 0x40, 0x20, 0x10***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestMarshalAndParseMessageForIPv6(t *testing.T) ***REMOVED***
	pshicmp := icmp.IPv6PseudoHeader(net.ParseIP("fe80::1"), net.ParseIP("ff02::1"))
	for i, tt := range marshalAndParseMessageForIPv6Tests ***REMOVED***
		for _, psh := range [][]byte***REMOVED***pshicmp, nil***REMOVED*** ***REMOVED***
			b, err := tt.Marshal(psh)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			m, err := icmp.ParseMessage(iana.ProtocolIPv6ICMP, b)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if m.Type != tt.Type || m.Code != tt.Code ***REMOVED***
				t.Errorf("#%v: got %v; want %v", i, m, &tt)
			***REMOVED***
			if !reflect.DeepEqual(m.Body, tt.Body) ***REMOVED***
				t.Errorf("#%v: got %v; want %v", i, m.Body, tt.Body)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
