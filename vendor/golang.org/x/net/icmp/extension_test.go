// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp

import (
	"net"
	"reflect"
	"testing"

	"golang.org/x/net/internal/iana"
)

var marshalAndParseExtensionTests = []struct ***REMOVED***
	proto int
	hdr   []byte
	obj   []byte
	exts  []Extension
***REMOVED******REMOVED***
	// MPLS label stack with no label
	***REMOVED***
		proto: iana.ProtocolICMP,
		hdr: []byte***REMOVED***
			0x20, 0x00, 0x00, 0x00,
		***REMOVED***,
		obj: []byte***REMOVED***
			0x00, 0x04, 0x01, 0x01,
		***REMOVED***,
		exts: []Extension***REMOVED***
			&MPLSLabelStack***REMOVED***
				Class: classMPLSLabelStack,
				Type:  typeIncomingMPLSLabelStack,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	// MPLS label stack with a single label
	***REMOVED***
		proto: iana.ProtocolIPv6ICMP,
		hdr: []byte***REMOVED***
			0x20, 0x00, 0x00, 0x00,
		***REMOVED***,
		obj: []byte***REMOVED***
			0x00, 0x08, 0x01, 0x01,
			0x03, 0xe8, 0xe9, 0xff,
		***REMOVED***,
		exts: []Extension***REMOVED***
			&MPLSLabelStack***REMOVED***
				Class: classMPLSLabelStack,
				Type:  typeIncomingMPLSLabelStack,
				Labels: []MPLSLabel***REMOVED***
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
	// MPLS label stack with multiple labels
	***REMOVED***
		proto: iana.ProtocolICMP,
		hdr: []byte***REMOVED***
			0x20, 0x00, 0x00, 0x00,
		***REMOVED***,
		obj: []byte***REMOVED***
			0x00, 0x0c, 0x01, 0x01,
			0x03, 0xe8, 0xde, 0xfe,
			0x03, 0xe8, 0xe1, 0xff,
		***REMOVED***,
		exts: []Extension***REMOVED***
			&MPLSLabelStack***REMOVED***
				Class: classMPLSLabelStack,
				Type:  typeIncomingMPLSLabelStack,
				Labels: []MPLSLabel***REMOVED***
					***REMOVED***
						Label: 16013,
						TC:    0x7,
						S:     false,
						TTL:   254,
					***REMOVED***,
					***REMOVED***
						Label: 16014,
						TC:    0,
						S:     true,
						TTL:   255,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	// Interface information with no attribute
	***REMOVED***
		proto: iana.ProtocolICMP,
		hdr: []byte***REMOVED***
			0x20, 0x00, 0x00, 0x00,
		***REMOVED***,
		obj: []byte***REMOVED***
			0x00, 0x04, 0x02, 0x00,
		***REMOVED***,
		exts: []Extension***REMOVED***
			&InterfaceInfo***REMOVED***
				Class: classInterfaceInfo,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	// Interface information with ifIndex and name
	***REMOVED***
		proto: iana.ProtocolICMP,
		hdr: []byte***REMOVED***
			0x20, 0x00, 0x00, 0x00,
		***REMOVED***,
		obj: []byte***REMOVED***
			0x00, 0x10, 0x02, 0x0a,
			0x00, 0x00, 0x00, 0x10,
			0x08, byte('e'), byte('n'), byte('1'),
			byte('0'), byte('1'), 0x00, 0x00,
		***REMOVED***,
		exts: []Extension***REMOVED***
			&InterfaceInfo***REMOVED***
				Class: classInterfaceInfo,
				Type:  0x0a,
				Interface: &net.Interface***REMOVED***
					Index: 16,
					Name:  "en101",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	// Interface information with ifIndex, IPAddr, name and MTU
	***REMOVED***
		proto: iana.ProtocolIPv6ICMP,
		hdr: []byte***REMOVED***
			0x20, 0x00, 0x00, 0x00,
		***REMOVED***,
		obj: []byte***REMOVED***
			0x00, 0x28, 0x02, 0x0f,
			0x00, 0x00, 0x00, 0x0f,
			0x00, 0x02, 0x00, 0x00,
			0xfe, 0x80, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x01,
			0x08, byte('e'), byte('n'), byte('1'),
			byte('0'), byte('1'), 0x00, 0x00,
			0x00, 0x00, 0x20, 0x00,
		***REMOVED***,
		exts: []Extension***REMOVED***
			&InterfaceInfo***REMOVED***
				Class: classInterfaceInfo,
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
***REMOVED***

func TestMarshalAndParseExtension(t *testing.T) ***REMOVED***
	for i, tt := range marshalAndParseExtensionTests ***REMOVED***
		for j, ext := range tt.exts ***REMOVED***
			var err error
			var b []byte
			switch ext := ext.(type) ***REMOVED***
			case *MPLSLabelStack:
				b, err = ext.Marshal(tt.proto)
				if err != nil ***REMOVED***
					t.Errorf("#%v/%v: %v", i, j, err)
					continue
				***REMOVED***
			case *InterfaceInfo:
				b, err = ext.Marshal(tt.proto)
				if err != nil ***REMOVED***
					t.Errorf("#%v/%v: %v", i, j, err)
					continue
				***REMOVED***
			***REMOVED***
			if !reflect.DeepEqual(b, tt.obj) ***REMOVED***
				t.Errorf("#%v/%v: got %#v; want %#v", i, j, b, tt.obj)
				continue
			***REMOVED***
		***REMOVED***

		for j, wire := range []struct ***REMOVED***
			data     []byte // original datagram
			inlattr  int    // length of padded original datagram, a hint
			outlattr int    // length of padded original datagram, a want
			err      error
		***REMOVED******REMOVED***
			***REMOVED***nil, 0, -1, errNoExtension***REMOVED***,
			***REMOVED***make([]byte, 127), 128, -1, errNoExtension***REMOVED***,

			***REMOVED***make([]byte, 128), 127, -1, errNoExtension***REMOVED***,
			***REMOVED***make([]byte, 128), 128, -1, errNoExtension***REMOVED***,
			***REMOVED***make([]byte, 128), 129, -1, errNoExtension***REMOVED***,

			***REMOVED***append(make([]byte, 128), append(tt.hdr, tt.obj...)...), 127, 128, nil***REMOVED***,
			***REMOVED***append(make([]byte, 128), append(tt.hdr, tt.obj...)...), 128, 128, nil***REMOVED***,
			***REMOVED***append(make([]byte, 128), append(tt.hdr, tt.obj...)...), 129, 128, nil***REMOVED***,

			***REMOVED***append(make([]byte, 512), append(tt.hdr, tt.obj...)...), 511, -1, errNoExtension***REMOVED***,
			***REMOVED***append(make([]byte, 512), append(tt.hdr, tt.obj...)...), 512, 512, nil***REMOVED***,
			***REMOVED***append(make([]byte, 512), append(tt.hdr, tt.obj...)...), 513, -1, errNoExtension***REMOVED***,
		***REMOVED*** ***REMOVED***
			exts, l, err := parseExtensions(wire.data, wire.inlattr)
			if err != wire.err ***REMOVED***
				t.Errorf("#%v/%v: got %v; want %v", i, j, err, wire.err)
				continue
			***REMOVED***
			if wire.err != nil ***REMOVED***
				continue
			***REMOVED***
			if l != wire.outlattr ***REMOVED***
				t.Errorf("#%v/%v: got %v; want %v", i, j, l, wire.outlattr)
			***REMOVED***
			if !reflect.DeepEqual(exts, tt.exts) ***REMOVED***
				for j, ext := range exts ***REMOVED***
					switch ext := ext.(type) ***REMOVED***
					case *MPLSLabelStack:
						want := tt.exts[j].(*MPLSLabelStack)
						t.Errorf("#%v/%v: got %#v; want %#v", i, j, ext, want)
					case *InterfaceInfo:
						want := tt.exts[j].(*InterfaceInfo)
						t.Errorf("#%v/%v: got %#v; want %#v", i, j, ext, want)
					***REMOVED***
				***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var parseInterfaceNameTests = []struct ***REMOVED***
	b []byte
	error
***REMOVED******REMOVED***
	***REMOVED***[]byte***REMOVED***0, 'e', 'n', '0'***REMOVED***, errInvalidExtension***REMOVED***,
	***REMOVED***[]byte***REMOVED***4, 'e', 'n', '0'***REMOVED***, nil***REMOVED***,
	***REMOVED***[]byte***REMOVED***7, 'e', 'n', '0', 0xff, 0xff, 0xff, 0xff***REMOVED***, errInvalidExtension***REMOVED***,
	***REMOVED***[]byte***REMOVED***8, 'e', 'n', '0', 0xff, 0xff, 0xff***REMOVED***, errMessageTooShort***REMOVED***,
***REMOVED***

func TestParseInterfaceName(t *testing.T) ***REMOVED***
	ifi := InterfaceInfo***REMOVED***Interface: &net.Interface***REMOVED******REMOVED******REMOVED***
	for i, tt := range parseInterfaceNameTests ***REMOVED***
		if _, err := ifi.parseName(tt.b); err != tt.error ***REMOVED***
			t.Errorf("#%d: got %v; want %v", i, err, tt.error)
		***REMOVED***
	***REMOVED***
***REMOVED***
