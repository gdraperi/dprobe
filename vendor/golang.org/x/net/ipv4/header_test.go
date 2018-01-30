// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

import (
	"bytes"
	"encoding/binary"
	"net"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/net/internal/socket"
)

type headerTest struct ***REMOVED***
	wireHeaderFromKernel          []byte
	wireHeaderToKernel            []byte
	wireHeaderFromTradBSDKernel   []byte
	wireHeaderToTradBSDKernel     []byte
	wireHeaderFromFreeBSD10Kernel []byte
	wireHeaderToFreeBSD10Kernel   []byte
	*Header
***REMOVED***

var headerLittleEndianTests = []headerTest***REMOVED***
	// TODO(mikio): Add platform dependent wire header formats when
	// we support new platforms.
	***REMOVED***
		wireHeaderFromKernel: []byte***REMOVED***
			0x45, 0x01, 0xbe, 0xef,
			0xca, 0xfe, 0x45, 0xdc,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
		***REMOVED***,
		wireHeaderToKernel: []byte***REMOVED***
			0x45, 0x01, 0xbe, 0xef,
			0xca, 0xfe, 0x45, 0xdc,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
		***REMOVED***,
		wireHeaderFromTradBSDKernel: []byte***REMOVED***
			0x45, 0x01, 0xdb, 0xbe,
			0xca, 0xfe, 0xdc, 0x45,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
		***REMOVED***,
		wireHeaderToTradBSDKernel: []byte***REMOVED***
			0x45, 0x01, 0xef, 0xbe,
			0xca, 0xfe, 0xdc, 0x45,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
		***REMOVED***,
		wireHeaderFromFreeBSD10Kernel: []byte***REMOVED***
			0x45, 0x01, 0xef, 0xbe,
			0xca, 0xfe, 0xdc, 0x45,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
		***REMOVED***,
		wireHeaderToFreeBSD10Kernel: []byte***REMOVED***
			0x45, 0x01, 0xef, 0xbe,
			0xca, 0xfe, 0xdc, 0x45,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
		***REMOVED***,
		Header: &Header***REMOVED***
			Version:  Version,
			Len:      HeaderLen,
			TOS:      1,
			TotalLen: 0xbeef,
			ID:       0xcafe,
			Flags:    DontFragment,
			FragOff:  1500,
			TTL:      255,
			Protocol: 1,
			Checksum: 0xdead,
			Src:      net.IPv4(172, 16, 254, 254),
			Dst:      net.IPv4(192, 168, 0, 1),
		***REMOVED***,
	***REMOVED***,

	// with option headers
	***REMOVED***
		wireHeaderFromKernel: []byte***REMOVED***
			0x46, 0x01, 0xbe, 0xf3,
			0xca, 0xfe, 0x45, 0xdc,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
			0xff, 0xfe, 0xfe, 0xff,
		***REMOVED***,
		wireHeaderToKernel: []byte***REMOVED***
			0x46, 0x01, 0xbe, 0xf3,
			0xca, 0xfe, 0x45, 0xdc,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
			0xff, 0xfe, 0xfe, 0xff,
		***REMOVED***,
		wireHeaderFromTradBSDKernel: []byte***REMOVED***
			0x46, 0x01, 0xdb, 0xbe,
			0xca, 0xfe, 0xdc, 0x45,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
			0xff, 0xfe, 0xfe, 0xff,
		***REMOVED***,
		wireHeaderToTradBSDKernel: []byte***REMOVED***
			0x46, 0x01, 0xf3, 0xbe,
			0xca, 0xfe, 0xdc, 0x45,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
			0xff, 0xfe, 0xfe, 0xff,
		***REMOVED***,
		wireHeaderFromFreeBSD10Kernel: []byte***REMOVED***
			0x46, 0x01, 0xf3, 0xbe,
			0xca, 0xfe, 0xdc, 0x45,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
			0xff, 0xfe, 0xfe, 0xff,
		***REMOVED***,
		wireHeaderToFreeBSD10Kernel: []byte***REMOVED***
			0x46, 0x01, 0xf3, 0xbe,
			0xca, 0xfe, 0xdc, 0x45,
			0xff, 0x01, 0xde, 0xad,
			172, 16, 254, 254,
			192, 168, 0, 1,
			0xff, 0xfe, 0xfe, 0xff,
		***REMOVED***,
		Header: &Header***REMOVED***
			Version:  Version,
			Len:      HeaderLen + 4,
			TOS:      1,
			TotalLen: 0xbef3,
			ID:       0xcafe,
			Flags:    DontFragment,
			FragOff:  1500,
			TTL:      255,
			Protocol: 1,
			Checksum: 0xdead,
			Src:      net.IPv4(172, 16, 254, 254),
			Dst:      net.IPv4(192, 168, 0, 1),
			Options:  []byte***REMOVED***0xff, 0xfe, 0xfe, 0xff***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestMarshalHeader(t *testing.T) ***REMOVED***
	if socket.NativeEndian != binary.LittleEndian ***REMOVED***
		t.Skip("no test for non-little endian machine yet")
	***REMOVED***

	for _, tt := range headerLittleEndianTests ***REMOVED***
		b, err := tt.Header.Marshal()
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		var wh []byte
		switch runtime.GOOS ***REMOVED***
		case "darwin", "dragonfly", "netbsd":
			wh = tt.wireHeaderToTradBSDKernel
		case "freebsd":
			switch ***REMOVED***
			case freebsdVersion < 1000000:
				wh = tt.wireHeaderToTradBSDKernel
			case 1000000 <= freebsdVersion && freebsdVersion < 1100000:
				wh = tt.wireHeaderToFreeBSD10Kernel
			default:
				wh = tt.wireHeaderToKernel
			***REMOVED***
		default:
			wh = tt.wireHeaderToKernel
		***REMOVED***
		if !bytes.Equal(b, wh) ***REMOVED***
			t.Fatalf("got %#v; want %#v", b, wh)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseHeader(t *testing.T) ***REMOVED***
	if socket.NativeEndian != binary.LittleEndian ***REMOVED***
		t.Skip("no test for big endian machine yet")
	***REMOVED***

	for _, tt := range headerLittleEndianTests ***REMOVED***
		var wh []byte
		switch runtime.GOOS ***REMOVED***
		case "darwin", "dragonfly", "netbsd":
			wh = tt.wireHeaderFromTradBSDKernel
		case "freebsd":
			switch ***REMOVED***
			case freebsdVersion < 1000000:
				wh = tt.wireHeaderFromTradBSDKernel
			case 1000000 <= freebsdVersion && freebsdVersion < 1100000:
				wh = tt.wireHeaderFromFreeBSD10Kernel
			default:
				wh = tt.wireHeaderFromKernel
			***REMOVED***
		default:
			wh = tt.wireHeaderFromKernel
		***REMOVED***
		h, err := ParseHeader(wh)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if err := h.Parse(wh); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if !reflect.DeepEqual(h, tt.Header) ***REMOVED***
			t.Fatalf("got %#v; want %#v", h, tt.Header)
		***REMOVED***
		s := h.String()
		if strings.Contains(s, ",") ***REMOVED***
			t.Fatalf("should be space-separated values: %s", s)
		***REMOVED***
	***REMOVED***
***REMOVED***
