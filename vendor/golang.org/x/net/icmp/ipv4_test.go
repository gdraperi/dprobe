// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp

import (
	"encoding/binary"
	"net"
	"reflect"
	"runtime"
	"testing"

	"golang.org/x/net/internal/socket"
	"golang.org/x/net/ipv4"
)

type ipv4HeaderTest struct ***REMOVED***
	wireHeaderFromKernel        [ipv4.HeaderLen]byte
	wireHeaderFromTradBSDKernel [ipv4.HeaderLen]byte
	Header                      *ipv4.Header
***REMOVED***

var ipv4HeaderLittleEndianTest = ipv4HeaderTest***REMOVED***
	// TODO(mikio): Add platform dependent wire header formats when
	// we support new platforms.
	wireHeaderFromKernel: [ipv4.HeaderLen]byte***REMOVED***
		0x45, 0x01, 0xbe, 0xef,
		0xca, 0xfe, 0x45, 0xdc,
		0xff, 0x01, 0xde, 0xad,
		172, 16, 254, 254,
		192, 168, 0, 1,
	***REMOVED***,
	wireHeaderFromTradBSDKernel: [ipv4.HeaderLen]byte***REMOVED***
		0x45, 0x01, 0xef, 0xbe,
		0xca, 0xfe, 0x45, 0xdc,
		0xff, 0x01, 0xde, 0xad,
		172, 16, 254, 254,
		192, 168, 0, 1,
	***REMOVED***,
	Header: &ipv4.Header***REMOVED***
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TOS:      1,
		TotalLen: 0xbeef,
		ID:       0xcafe,
		Flags:    ipv4.DontFragment,
		FragOff:  1500,
		TTL:      255,
		Protocol: 1,
		Checksum: 0xdead,
		Src:      net.IPv4(172, 16, 254, 254),
		Dst:      net.IPv4(192, 168, 0, 1),
	***REMOVED***,
***REMOVED***

func TestParseIPv4Header(t *testing.T) ***REMOVED***
	tt := &ipv4HeaderLittleEndianTest
	if socket.NativeEndian != binary.LittleEndian ***REMOVED***
		t.Skip("no test for non-little endian machine yet")
	***REMOVED***

	var wh []byte
	switch runtime.GOOS ***REMOVED***
	case "darwin":
		wh = tt.wireHeaderFromTradBSDKernel[:]
	case "freebsd":
		if freebsdVersion >= 1000000 ***REMOVED***
			wh = tt.wireHeaderFromKernel[:]
		***REMOVED*** else ***REMOVED***
			wh = tt.wireHeaderFromTradBSDKernel[:]
		***REMOVED***
	default:
		wh = tt.wireHeaderFromKernel[:]
	***REMOVED***
	h, err := ParseIPv4Header(wh)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(h, tt.Header) ***REMOVED***
		t.Fatalf("got %#v; want %#v", h, tt.Header)
	***REMOVED***
***REMOVED***
