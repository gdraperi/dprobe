// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp_test

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	"golang.org/x/net/icmp"
	"golang.org/x/net/internal/iana"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var marshalAndParseMultipartMessageForIPv4Tests = []icmp.Message***REMOVED***
	***REMOVED***
		Type: ipv4.ICMPTypeDestinationUnreachable, Code: 15,
		Body: &icmp.DstUnreach***REMOVED***
			Data: []byte("ERROR-INVOKING-PACKET"),
			Extensions: []icmp.Extension***REMOVED***
				&icmp.MPLSLabelStack***REMOVED***
					Class: 1,
					Type:  1,
					Labels: []icmp.MPLSLabel***REMOVED***
						***REMOVED***
							Label: 16014,
							TC:    0x4,
							S:     true,
							TTL:   255,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
				&icmp.InterfaceInfo***REMOVED***
					Class: 2,
					Type:  0x0f,
					Interface: &net.Interface***REMOVED***
						Index: 15,
						Name:  "en101",
						MTU:   8192,
					***REMOVED***,
					Addr: &net.IPAddr***REMOVED***
						IP: net.IPv4(192, 168, 0, 1).To4(),
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv4.ICMPTypeTimeExceeded, Code: 1,
		Body: &icmp.TimeExceeded***REMOVED***
			Data: []byte("ERROR-INVOKING-PACKET"),
			Extensions: []icmp.Extension***REMOVED***
				&icmp.InterfaceInfo***REMOVED***
					Class: 2,
					Type:  0x0f,
					Interface: &net.Interface***REMOVED***
						Index: 15,
						Name:  "en101",
						MTU:   8192,
					***REMOVED***,
					Addr: &net.IPAddr***REMOVED***
						IP: net.IPv4(192, 168, 0, 1).To4(),
					***REMOVED***,
				***REMOVED***,
				&icmp.MPLSLabelStack***REMOVED***
					Class: 1,
					Type:  1,
					Labels: []icmp.MPLSLabel***REMOVED***
						***REMOVED***
							Label: 16014,
							TC:    0x4,
							S:     true,
							TTL:   255,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv4.ICMPTypeParameterProblem, Code: 2,
		Body: &icmp.ParamProb***REMOVED***
			Pointer: 8,
			Data:    []byte("ERROR-INVOKING-PACKET"),
			Extensions: []icmp.Extension***REMOVED***
				&icmp.MPLSLabelStack***REMOVED***
					Class: 1,
					Type:  1,
					Labels: []icmp.MPLSLabel***REMOVED***
						***REMOVED***
							Label: 16014,
							TC:    0x4,
							S:     true,
							TTL:   255,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
				&icmp.InterfaceInfo***REMOVED***
					Class: 2,
					Type:  0x0f,
					Interface: &net.Interface***REMOVED***
						Index: 15,
						Name:  "en101",
						MTU:   8192,
					***REMOVED***,
					Addr: &net.IPAddr***REMOVED***
						IP: net.IPv4(192, 168, 0, 1).To4(),
					***REMOVED***,
				***REMOVED***,
				&icmp.InterfaceInfo***REMOVED***
					Class: 2,
					Type:  0x2f,
					Interface: &net.Interface***REMOVED***
						Index: 16,
						Name:  "en102",
						MTU:   8192,
					***REMOVED***,
					Addr: &net.IPAddr***REMOVED***
						IP: net.IPv4(192, 168, 0, 2).To4(),
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestMarshalAndParseMultipartMessageForIPv4(t *testing.T) ***REMOVED***
	for i, tt := range marshalAndParseMultipartMessageForIPv4Tests ***REMOVED***
		b, err := tt.Marshal(nil)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if b[5] != 32 ***REMOVED***
			t.Errorf("#%v: got %v; want 32", i, b[5])
		***REMOVED***
		m, err := icmp.ParseMessage(iana.ProtocolICMP, b)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if m.Type != tt.Type || m.Code != tt.Code ***REMOVED***
			t.Errorf("#%v: got %v; want %v", i, m, &tt)
		***REMOVED***
		switch m.Type ***REMOVED***
		case ipv4.ICMPTypeDestinationUnreachable:
			got, want := m.Body.(*icmp.DstUnreach), tt.Body.(*icmp.DstUnreach)
			if !reflect.DeepEqual(got.Extensions, want.Extensions) ***REMOVED***
				t.Error(dumpExtensions(i, got.Extensions, want.Extensions))
			***REMOVED***
			if len(got.Data) != 128 ***REMOVED***
				t.Errorf("#%v: got %v; want 128", i, len(got.Data))
			***REMOVED***
		case ipv4.ICMPTypeTimeExceeded:
			got, want := m.Body.(*icmp.TimeExceeded), tt.Body.(*icmp.TimeExceeded)
			if !reflect.DeepEqual(got.Extensions, want.Extensions) ***REMOVED***
				t.Error(dumpExtensions(i, got.Extensions, want.Extensions))
			***REMOVED***
			if len(got.Data) != 128 ***REMOVED***
				t.Errorf("#%v: got %v; want 128", i, len(got.Data))
			***REMOVED***
		case ipv4.ICMPTypeParameterProblem:
			got, want := m.Body.(*icmp.ParamProb), tt.Body.(*icmp.ParamProb)
			if !reflect.DeepEqual(got.Extensions, want.Extensions) ***REMOVED***
				t.Error(dumpExtensions(i, got.Extensions, want.Extensions))
			***REMOVED***
			if len(got.Data) != 128 ***REMOVED***
				t.Errorf("#%v: got %v; want 128", i, len(got.Data))
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var marshalAndParseMultipartMessageForIPv6Tests = []icmp.Message***REMOVED***
	***REMOVED***
		Type: ipv6.ICMPTypeDestinationUnreachable, Code: 6,
		Body: &icmp.DstUnreach***REMOVED***
			Data: []byte("ERROR-INVOKING-PACKET"),
			Extensions: []icmp.Extension***REMOVED***
				&icmp.MPLSLabelStack***REMOVED***
					Class: 1,
					Type:  1,
					Labels: []icmp.MPLSLabel***REMOVED***
						***REMOVED***
							Label: 16014,
							TC:    0x4,
							S:     true,
							TTL:   255,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
				&icmp.InterfaceInfo***REMOVED***
					Class: 2,
					Type:  0x0f,
					Interface: &net.Interface***REMOVED***
						Index: 15,
						Name:  "en101",
						MTU:   8192,
					***REMOVED***,
					Addr: &net.IPAddr***REMOVED***
						IP:   net.ParseIP("fe80::1"),
						Zone: "en101",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Type: ipv6.ICMPTypeTimeExceeded, Code: 1,
		Body: &icmp.TimeExceeded***REMOVED***
			Data: []byte("ERROR-INVOKING-PACKET"),
			Extensions: []icmp.Extension***REMOVED***
				&icmp.InterfaceInfo***REMOVED***
					Class: 2,
					Type:  0x0f,
					Interface: &net.Interface***REMOVED***
						Index: 15,
						Name:  "en101",
						MTU:   8192,
					***REMOVED***,
					Addr: &net.IPAddr***REMOVED***
						IP:   net.ParseIP("fe80::1"),
						Zone: "en101",
					***REMOVED***,
				***REMOVED***,
				&icmp.MPLSLabelStack***REMOVED***
					Class: 1,
					Type:  1,
					Labels: []icmp.MPLSLabel***REMOVED***
						***REMOVED***
							Label: 16014,
							TC:    0x4,
							S:     true,
							TTL:   255,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
				&icmp.InterfaceInfo***REMOVED***
					Class: 2,
					Type:  0x2f,
					Interface: &net.Interface***REMOVED***
						Index: 16,
						Name:  "en102",
						MTU:   8192,
					***REMOVED***,
					Addr: &net.IPAddr***REMOVED***
						IP:   net.ParseIP("fe80::1"),
						Zone: "en102",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestMarshalAndParseMultipartMessageForIPv6(t *testing.T) ***REMOVED***
	pshicmp := icmp.IPv6PseudoHeader(net.ParseIP("fe80::1"), net.ParseIP("ff02::1"))
	for i, tt := range marshalAndParseMultipartMessageForIPv6Tests ***REMOVED***
		for _, psh := range [][]byte***REMOVED***pshicmp, nil***REMOVED*** ***REMOVED***
			b, err := tt.Marshal(psh)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if b[4] != 16 ***REMOVED***
				t.Errorf("#%v: got %v; want 16", i, b[4])
			***REMOVED***
			m, err := icmp.ParseMessage(iana.ProtocolIPv6ICMP, b)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if m.Type != tt.Type || m.Code != tt.Code ***REMOVED***
				t.Errorf("#%v: got %v; want %v", i, m, &tt)
			***REMOVED***
			switch m.Type ***REMOVED***
			case ipv6.ICMPTypeDestinationUnreachable:
				got, want := m.Body.(*icmp.DstUnreach), tt.Body.(*icmp.DstUnreach)
				if !reflect.DeepEqual(got.Extensions, want.Extensions) ***REMOVED***
					t.Error(dumpExtensions(i, got.Extensions, want.Extensions))
				***REMOVED***
				if len(got.Data) != 128 ***REMOVED***
					t.Errorf("#%v: got %v; want 128", i, len(got.Data))
				***REMOVED***
			case ipv6.ICMPTypeTimeExceeded:
				got, want := m.Body.(*icmp.TimeExceeded), tt.Body.(*icmp.TimeExceeded)
				if !reflect.DeepEqual(got.Extensions, want.Extensions) ***REMOVED***
					t.Error(dumpExtensions(i, got.Extensions, want.Extensions))
				***REMOVED***
				if len(got.Data) != 128 ***REMOVED***
					t.Errorf("#%v: got %v; want 128", i, len(got.Data))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func dumpExtensions(i int, gotExts, wantExts []icmp.Extension) string ***REMOVED***
	var s string
	for j, got := range gotExts ***REMOVED***
		switch got := got.(type) ***REMOVED***
		case *icmp.MPLSLabelStack:
			want := wantExts[j].(*icmp.MPLSLabelStack)
			if !reflect.DeepEqual(got, want) ***REMOVED***
				s += fmt.Sprintf("#%v/%v: got %#v; want %#v\n", i, j, got, want)
			***REMOVED***
		case *icmp.InterfaceInfo:
			want := wantExts[j].(*icmp.InterfaceInfo)
			if !reflect.DeepEqual(got, want) ***REMOVED***
				s += fmt.Sprintf("#%v/%v: got %#v, %#v, %#v; want %#v, %#v, %#v\n", i, j, got, got.Interface, got.Addr, want, want.Interface, want.Addr)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return s[:len(s)-1]
***REMOVED***

var multipartMessageBodyLenTests = []struct ***REMOVED***
	proto int
	in    icmp.MessageBody
	out   int
***REMOVED******REMOVED***
	***REMOVED***
		iana.ProtocolICMP,
		&icmp.DstUnreach***REMOVED***
			Data: make([]byte, ipv4.HeaderLen),
		***REMOVED***,
		4 + ipv4.HeaderLen, // unused and original datagram
	***REMOVED***,
	***REMOVED***
		iana.ProtocolICMP,
		&icmp.TimeExceeded***REMOVED***
			Data: make([]byte, ipv4.HeaderLen),
		***REMOVED***,
		4 + ipv4.HeaderLen, // unused and original datagram
	***REMOVED***,
	***REMOVED***
		iana.ProtocolICMP,
		&icmp.ParamProb***REMOVED***
			Data: make([]byte, ipv4.HeaderLen),
		***REMOVED***,
		4 + ipv4.HeaderLen, // [pointer, unused] and original datagram
	***REMOVED***,

	***REMOVED***
		iana.ProtocolICMP,
		&icmp.ParamProb***REMOVED***
			Data: make([]byte, ipv4.HeaderLen),
			Extensions: []icmp.Extension***REMOVED***
				&icmp.MPLSLabelStack***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		4 + 4 + 4 + 0 + 128, // [pointer, length, unused], extension header, object header, object payload, original datagram
	***REMOVED***,
	***REMOVED***
		iana.ProtocolICMP,
		&icmp.ParamProb***REMOVED***
			Data: make([]byte, 128),
			Extensions: []icmp.Extension***REMOVED***
				&icmp.MPLSLabelStack***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		4 + 4 + 4 + 0 + 128, // [pointer, length, unused], extension header, object header, object payload and original datagram
	***REMOVED***,
	***REMOVED***
		iana.ProtocolICMP,
		&icmp.ParamProb***REMOVED***
			Data: make([]byte, 129),
			Extensions: []icmp.Extension***REMOVED***
				&icmp.MPLSLabelStack***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		4 + 4 + 4 + 0 + 132, // [pointer, length, unused], extension header, object header, object payload and original datagram
	***REMOVED***,

	***REMOVED***
		iana.ProtocolIPv6ICMP,
		&icmp.DstUnreach***REMOVED***
			Data: make([]byte, ipv6.HeaderLen),
		***REMOVED***,
		4 + ipv6.HeaderLen, // unused and original datagram
	***REMOVED***,
	***REMOVED***
		iana.ProtocolIPv6ICMP,
		&icmp.PacketTooBig***REMOVED***
			Data: make([]byte, ipv6.HeaderLen),
		***REMOVED***,
		4 + ipv6.HeaderLen, // mtu and original datagram
	***REMOVED***,
	***REMOVED***
		iana.ProtocolIPv6ICMP,
		&icmp.TimeExceeded***REMOVED***
			Data: make([]byte, ipv6.HeaderLen),
		***REMOVED***,
		4 + ipv6.HeaderLen, // unused and original datagram
	***REMOVED***,
	***REMOVED***
		iana.ProtocolIPv6ICMP,
		&icmp.ParamProb***REMOVED***
			Data: make([]byte, ipv6.HeaderLen),
		***REMOVED***,
		4 + ipv6.HeaderLen, // pointer and original datagram
	***REMOVED***,

	***REMOVED***
		iana.ProtocolIPv6ICMP,
		&icmp.DstUnreach***REMOVED***
			Data: make([]byte, 127),
			Extensions: []icmp.Extension***REMOVED***
				&icmp.MPLSLabelStack***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		4 + 4 + 4 + 0 + 128, // [length, unused], extension header, object header, object payload and original datagram
	***REMOVED***,
	***REMOVED***
		iana.ProtocolIPv6ICMP,
		&icmp.DstUnreach***REMOVED***
			Data: make([]byte, 128),
			Extensions: []icmp.Extension***REMOVED***
				&icmp.MPLSLabelStack***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		4 + 4 + 4 + 0 + 128, // [length, unused], extension header, object header, object payload and original datagram
	***REMOVED***,
	***REMOVED***
		iana.ProtocolIPv6ICMP,
		&icmp.DstUnreach***REMOVED***
			Data: make([]byte, 129),
			Extensions: []icmp.Extension***REMOVED***
				&icmp.MPLSLabelStack***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		4 + 4 + 4 + 0 + 136, // [length, unused], extension header, object header, object payload and original datagram
	***REMOVED***,
***REMOVED***

func TestMultipartMessageBodyLen(t *testing.T) ***REMOVED***
	for i, tt := range multipartMessageBodyLenTests ***REMOVED***
		if out := tt.in.Len(tt.proto); out != tt.out ***REMOVED***
			t.Errorf("#%d: got %d; want %d", i, out, tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***
