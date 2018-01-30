// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package route

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestFetchAndParseRIB(t *testing.T) ***REMOVED***
	for _, typ := range []RIBType***REMOVED***sysNET_RT_DUMP, sysNET_RT_IFLIST***REMOVED*** ***REMOVED***
		var lastErr error
		var ms []Message
		for _, af := range []int***REMOVED***sysAF_UNSPEC, sysAF_INET, sysAF_INET6***REMOVED*** ***REMOVED***
			rs, err := fetchAndParseRIB(af, typ)
			if err != nil ***REMOVED***
				lastErr = err
				continue
			***REMOVED***
			ms = append(ms, rs...)
		***REMOVED***
		if len(ms) == 0 && lastErr != nil ***REMOVED***
			t.Error(typ, lastErr)
			continue
		***REMOVED***
		ss, err := msgs(ms).validate()
		if err != nil ***REMOVED***
			t.Error(typ, err)
			continue
		***REMOVED***
		for _, s := range ss ***REMOVED***
			t.Log(typ, s)
		***REMOVED***
	***REMOVED***
***REMOVED***

var (
	rtmonSock int
	rtmonErr  error
)

func init() ***REMOVED***
	// We need to keep rtmonSock alive to avoid treading on
	// recycled socket descriptors.
	rtmonSock, rtmonErr = syscall.Socket(sysAF_ROUTE, sysSOCK_RAW, sysAF_UNSPEC)
***REMOVED***

// TestMonitorAndParseRIB leaks a worker goroutine and a socket
// descriptor but that's intentional.
func TestMonitorAndParseRIB(t *testing.T) ***REMOVED***
	if testing.Short() || os.Getuid() != 0 ***REMOVED***
		t.Skip("must be root")
	***REMOVED***

	if rtmonErr != nil ***REMOVED***
		t.Fatal(rtmonErr)
	***REMOVED***

	// We suppose that using an IPv4 link-local address and the
	// dot1Q ID for Token Ring and FDDI doesn't harm anyone.
	pv := &propVirtual***REMOVED***addr: "169.254.0.1", mask: "255.255.255.0"***REMOVED***
	if err := pv.configure(1002); err != nil ***REMOVED***
		t.Skip(err)
	***REMOVED***
	if err := pv.setup(); err != nil ***REMOVED***
		t.Skip(err)
	***REMOVED***
	pv.teardown()

	go func() ***REMOVED***
		b := make([]byte, os.Getpagesize())
		for ***REMOVED***
			// There's no easy way to unblock this read
			// call because the routing message exchange
			// over routing socket is a connectionless
			// message-oriented protocol, no control plane
			// for signaling connectivity, and we cannot
			// use the net package of standard library due
			// to the lack of support for routing socket
			// and circular dependency.
			n, err := syscall.Read(rtmonSock, b)
			if err != nil ***REMOVED***
				return
			***REMOVED***
			ms, err := ParseRIB(0, b[:n])
			if err != nil ***REMOVED***
				t.Error(err)
				return
			***REMOVED***
			ss, err := msgs(ms).validate()
			if err != nil ***REMOVED***
				t.Error(err)
				return
			***REMOVED***
			for _, s := range ss ***REMOVED***
				t.Log(s)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	for _, vid := range []int***REMOVED***1002, 1003, 1004, 1005***REMOVED*** ***REMOVED***
		pv := &propVirtual***REMOVED***addr: "169.254.0.1", mask: "255.255.255.0"***REMOVED***
		if err := pv.configure(vid); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if err := pv.setup(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		time.Sleep(200 * time.Millisecond)
		if err := pv.teardown(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		time.Sleep(200 * time.Millisecond)
	***REMOVED***
***REMOVED***

func TestParseRIBWithFuzz(t *testing.T) ***REMOVED***
	for _, fuzz := range []string***REMOVED***
		"0\x00\x05\x050000000000000000" +
			"00000000000000000000" +
			"00000000000000000000" +
			"00000000000000000000" +
			"0000000000000\x02000000" +
			"00000000",
		"\x02\x00\x05\f0000000000000000" +
			"0\x0200000000000000",
		"\x02\x00\x05\x100000000000000\x1200" +
			"0\x00\xff\x00",
		"\x02\x00\x05\f0000000000000000" +
			"0\x12000\x00\x02\x0000",
		"\x00\x00\x00\x01\x00",
		"00000",
	***REMOVED*** ***REMOVED***
		for typ := RIBType(0); typ < 256; typ++ ***REMOVED***
			ParseRIB(typ, []byte(fuzz))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRouteMessage(t *testing.T) ***REMOVED***
	s, err := syscall.Socket(sysAF_ROUTE, sysSOCK_RAW, sysAF_UNSPEC)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer syscall.Close(s)

	var ms []RouteMessage
	for _, af := range []int***REMOVED***sysAF_INET, sysAF_INET6***REMOVED*** ***REMOVED***
		if _, err := fetchAndParseRIB(af, sysNET_RT_DUMP); err != nil ***REMOVED***
			t.Log(err)
			continue
		***REMOVED***
		switch af ***REMOVED***
		case sysAF_INET:
			ms = append(ms, []RouteMessage***REMOVED***
				***REMOVED***
					Type: sysRTM_GET,
					Addrs: []Addr***REMOVED***
						&Inet4Addr***REMOVED***IP: [4]byte***REMOVED***127, 0, 0, 1***REMOVED******REMOVED***,
						nil,
						nil,
						nil,
						&LinkAddr***REMOVED******REMOVED***,
						&Inet4Addr***REMOVED******REMOVED***,
						nil,
						&Inet4Addr***REMOVED******REMOVED***,
					***REMOVED***,
				***REMOVED***,
				***REMOVED***
					Type: sysRTM_GET,
					Addrs: []Addr***REMOVED***
						&Inet4Addr***REMOVED***IP: [4]byte***REMOVED***127, 0, 0, 1***REMOVED******REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***...)
		case sysAF_INET6:
			ms = append(ms, []RouteMessage***REMOVED***
				***REMOVED***
					Type: sysRTM_GET,
					Addrs: []Addr***REMOVED***
						&Inet6Addr***REMOVED***IP: [16]byte***REMOVED***0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1***REMOVED******REMOVED***,
						nil,
						nil,
						nil,
						&LinkAddr***REMOVED******REMOVED***,
						&Inet6Addr***REMOVED******REMOVED***,
						nil,
						&Inet6Addr***REMOVED******REMOVED***,
					***REMOVED***,
				***REMOVED***,
				***REMOVED***
					Type: sysRTM_GET,
					Addrs: []Addr***REMOVED***
						&Inet6Addr***REMOVED***IP: [16]byte***REMOVED***0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1***REMOVED******REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***...)
		***REMOVED***
	***REMOVED***
	for i, m := range ms ***REMOVED***
		m.ID = uintptr(os.Getpid())
		m.Seq = i + 1
		wb, err := m.Marshal()
		if err != nil ***REMOVED***
			t.Fatalf("%v: %v", m, err)
		***REMOVED***
		if _, err := syscall.Write(s, wb); err != nil ***REMOVED***
			t.Fatalf("%v: %v", m, err)
		***REMOVED***
		rb := make([]byte, os.Getpagesize())
		n, err := syscall.Read(s, rb)
		if err != nil ***REMOVED***
			t.Fatalf("%v: %v", m, err)
		***REMOVED***
		rms, err := ParseRIB(0, rb[:n])
		if err != nil ***REMOVED***
			t.Fatalf("%v: %v", m, err)
		***REMOVED***
		for _, rm := range rms ***REMOVED***
			err := rm.(*RouteMessage).Err
			if err != nil ***REMOVED***
				t.Errorf("%v: %v", m, err)
			***REMOVED***
		***REMOVED***
		ss, err := msgs(rms).validate()
		if err != nil ***REMOVED***
			t.Fatalf("%v: %v", m, err)
		***REMOVED***
		for _, s := range ss ***REMOVED***
			t.Log(s)
		***REMOVED***
	***REMOVED***
***REMOVED***
