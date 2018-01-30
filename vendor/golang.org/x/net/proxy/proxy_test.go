// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

type proxyFromEnvTest struct ***REMOVED***
	allProxyEnv string
	noProxyEnv  string
	wantTypeOf  Dialer
***REMOVED***

func (t proxyFromEnvTest) String() string ***REMOVED***
	var buf bytes.Buffer
	space := func() ***REMOVED***
		if buf.Len() > 0 ***REMOVED***
			buf.WriteByte(' ')
		***REMOVED***
	***REMOVED***
	if t.allProxyEnv != "" ***REMOVED***
		fmt.Fprintf(&buf, "all_proxy=%q", t.allProxyEnv)
	***REMOVED***
	if t.noProxyEnv != "" ***REMOVED***
		space()
		fmt.Fprintf(&buf, "no_proxy=%q", t.noProxyEnv)
	***REMOVED***
	return strings.TrimSpace(buf.String())
***REMOVED***

func TestFromEnvironment(t *testing.T) ***REMOVED***
	ResetProxyEnv()

	type dummyDialer struct ***REMOVED***
		direct
	***REMOVED***

	RegisterDialerType("irc", func(_ *url.URL, _ Dialer) (Dialer, error) ***REMOVED***
		return dummyDialer***REMOVED******REMOVED***, nil
	***REMOVED***)

	proxyFromEnvTests := []proxyFromEnvTest***REMOVED***
		***REMOVED***allProxyEnv: "127.0.0.1:8080", noProxyEnv: "localhost, 127.0.0.1", wantTypeOf: direct***REMOVED******REMOVED******REMOVED***,
		***REMOVED***allProxyEnv: "ftp://example.com:8000", noProxyEnv: "localhost, 127.0.0.1", wantTypeOf: direct***REMOVED******REMOVED******REMOVED***,
		***REMOVED***allProxyEnv: "socks5://example.com:8080", noProxyEnv: "localhost, 127.0.0.1", wantTypeOf: &PerHost***REMOVED******REMOVED******REMOVED***,
		***REMOVED***allProxyEnv: "irc://example.com:8000", wantTypeOf: dummyDialer***REMOVED******REMOVED******REMOVED***,
		***REMOVED***noProxyEnv: "localhost, 127.0.0.1", wantTypeOf: direct***REMOVED******REMOVED******REMOVED***,
		***REMOVED***wantTypeOf: direct***REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	for _, tt := range proxyFromEnvTests ***REMOVED***
		os.Setenv("ALL_PROXY", tt.allProxyEnv)
		os.Setenv("NO_PROXY", tt.noProxyEnv)
		ResetCachedEnvironment()

		d := FromEnvironment()
		if got, want := fmt.Sprintf("%T", d), fmt.Sprintf("%T", tt.wantTypeOf); got != want ***REMOVED***
			t.Errorf("%v: got type = %T, want %T", tt, d, tt.wantTypeOf)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFromURL(t *testing.T) ***REMOVED***
	endSystem, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil ***REMOVED***
		t.Fatalf("net.Listen failed: %v", err)
	***REMOVED***
	defer endSystem.Close()
	gateway, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil ***REMOVED***
		t.Fatalf("net.Listen failed: %v", err)
	***REMOVED***
	defer gateway.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go socks5Gateway(t, gateway, endSystem, socks5Domain, &wg)

	url, err := url.Parse("socks5://user:password@" + gateway.Addr().String())
	if err != nil ***REMOVED***
		t.Fatalf("url.Parse failed: %v", err)
	***REMOVED***
	proxy, err := FromURL(url, Direct)
	if err != nil ***REMOVED***
		t.Fatalf("FromURL failed: %v", err)
	***REMOVED***
	_, port, err := net.SplitHostPort(endSystem.Addr().String())
	if err != nil ***REMOVED***
		t.Fatalf("net.SplitHostPort failed: %v", err)
	***REMOVED***
	if c, err := proxy.Dial("tcp", "localhost:"+port); err != nil ***REMOVED***
		t.Fatalf("FromURL.Dial failed: %v", err)
	***REMOVED*** else ***REMOVED***
		c.Close()
	***REMOVED***

	wg.Wait()
***REMOVED***

func TestSOCKS5(t *testing.T) ***REMOVED***
	endSystem, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil ***REMOVED***
		t.Fatalf("net.Listen failed: %v", err)
	***REMOVED***
	defer endSystem.Close()
	gateway, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil ***REMOVED***
		t.Fatalf("net.Listen failed: %v", err)
	***REMOVED***
	defer gateway.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go socks5Gateway(t, gateway, endSystem, socks5IP4, &wg)

	proxy, err := SOCKS5("tcp", gateway.Addr().String(), nil, Direct)
	if err != nil ***REMOVED***
		t.Fatalf("SOCKS5 failed: %v", err)
	***REMOVED***
	if c, err := proxy.Dial("tcp", endSystem.Addr().String()); err != nil ***REMOVED***
		t.Fatalf("SOCKS5.Dial failed: %v", err)
	***REMOVED*** else ***REMOVED***
		c.Close()
	***REMOVED***

	wg.Wait()
***REMOVED***

func socks5Gateway(t *testing.T, gateway, endSystem net.Listener, typ byte, wg *sync.WaitGroup) ***REMOVED***
	defer wg.Done()

	c, err := gateway.Accept()
	if err != nil ***REMOVED***
		t.Errorf("net.Listener.Accept failed: %v", err)
		return
	***REMOVED***
	defer c.Close()

	b := make([]byte, 32)
	var n int
	if typ == socks5Domain ***REMOVED***
		n = 4
	***REMOVED*** else ***REMOVED***
		n = 3
	***REMOVED***
	if _, err := io.ReadFull(c, b[:n]); err != nil ***REMOVED***
		t.Errorf("io.ReadFull failed: %v", err)
		return
	***REMOVED***
	if _, err := c.Write([]byte***REMOVED***socks5Version, socks5AuthNone***REMOVED***); err != nil ***REMOVED***
		t.Errorf("net.Conn.Write failed: %v", err)
		return
	***REMOVED***
	if typ == socks5Domain ***REMOVED***
		n = 16
	***REMOVED*** else ***REMOVED***
		n = 10
	***REMOVED***
	if _, err := io.ReadFull(c, b[:n]); err != nil ***REMOVED***
		t.Errorf("io.ReadFull failed: %v", err)
		return
	***REMOVED***
	if b[0] != socks5Version || b[1] != socks5Connect || b[2] != 0x00 || b[3] != typ ***REMOVED***
		t.Errorf("got an unexpected packet: %#02x %#02x %#02x %#02x", b[0], b[1], b[2], b[3])
		return
	***REMOVED***
	if typ == socks5Domain ***REMOVED***
		copy(b[:5], []byte***REMOVED***socks5Version, 0x00, 0x00, socks5Domain, 9***REMOVED***)
		b = append(b, []byte("localhost")...)
	***REMOVED*** else ***REMOVED***
		copy(b[:4], []byte***REMOVED***socks5Version, 0x00, 0x00, socks5IP4***REMOVED***)
	***REMOVED***
	host, port, err := net.SplitHostPort(endSystem.Addr().String())
	if err != nil ***REMOVED***
		t.Errorf("net.SplitHostPort failed: %v", err)
		return
	***REMOVED***
	b = append(b, []byte(net.ParseIP(host).To4())...)
	p, err := strconv.Atoi(port)
	if err != nil ***REMOVED***
		t.Errorf("strconv.Atoi failed: %v", err)
		return
	***REMOVED***
	b = append(b, []byte***REMOVED***byte(p >> 8), byte(p)***REMOVED***...)
	if _, err := c.Write(b); err != nil ***REMOVED***
		t.Errorf("net.Conn.Write failed: %v", err)
		return
	***REMOVED***
***REMOVED***

func ResetProxyEnv() ***REMOVED***
	for _, env := range []*envOnce***REMOVED***allProxyEnv, noProxyEnv***REMOVED*** ***REMOVED***
		for _, v := range env.names ***REMOVED***
			os.Setenv(v, "")
		***REMOVED***
	***REMOVED***
	ResetCachedEnvironment()
***REMOVED***

func ResetCachedEnvironment() ***REMOVED***
	allProxyEnv.reset()
	noProxyEnv.reset()
***REMOVED***
