// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/http2/hpack"
)

var (
	extNet        = flag.Bool("extnet", false, "do external network tests")
	transportHost = flag.String("transporthost", "http2.golang.org", "hostname to use for TestTransport")
	insecure      = flag.Bool("insecure", false, "insecure TLS dials") // TODO: dead code. remove?
)

var tlsConfigInsecure = &tls.Config***REMOVED***InsecureSkipVerify: true***REMOVED***

type testContext struct***REMOVED******REMOVED***

func (testContext) Done() <-chan struct***REMOVED******REMOVED***                   ***REMOVED*** return make(chan struct***REMOVED******REMOVED***) ***REMOVED***
func (testContext) Err() error                              ***REMOVED*** panic("should not be called") ***REMOVED***
func (testContext) Deadline() (deadline time.Time, ok bool) ***REMOVED*** return time.Time***REMOVED******REMOVED***, false ***REMOVED***
func (testContext) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED***       ***REMOVED*** return nil ***REMOVED***

func TestTransportExternal(t *testing.T) ***REMOVED***
	if !*extNet ***REMOVED***
		t.Skip("skipping external network test")
	***REMOVED***
	req, _ := http.NewRequest("GET", "https://"+*transportHost+"/", nil)
	rt := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	res, err := rt.RoundTrip(req)
	if err != nil ***REMOVED***
		t.Fatalf("%v", err)
	***REMOVED***
	res.Write(os.Stdout)
***REMOVED***

type fakeTLSConn struct ***REMOVED***
	net.Conn
***REMOVED***

func (c *fakeTLSConn) ConnectionState() tls.ConnectionState ***REMOVED***
	return tls.ConnectionState***REMOVED***
		Version:     tls.VersionTLS12,
		CipherSuite: cipher_TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	***REMOVED***
***REMOVED***

func startH2cServer(t *testing.T) net.Listener ***REMOVED***
	h2Server := &Server***REMOVED******REMOVED***
	l := newLocalListener(t)
	go func() ***REMOVED***
		conn, err := l.Accept()
		if err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		h2Server.ServeConn(&fakeTLSConn***REMOVED***conn***REMOVED***, &ServeConnOpts***REMOVED***Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			fmt.Fprintf(w, "Hello, %v, http: %v", r.URL.Path, r.TLS == nil)
		***REMOVED***)***REMOVED***)
	***REMOVED***()
	return l
***REMOVED***

func TestTransportH2c(t *testing.T) ***REMOVED***
	l := startH2cServer(t)
	defer l.Close()
	req, err := http.NewRequest("GET", "http://"+l.Addr().String()+"/foobar", nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tr := &Transport***REMOVED***
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) ***REMOVED***
			return net.Dial(network, addr)
		***REMOVED***,
	***REMOVED***
	res, err := tr.RoundTrip(req)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if res.ProtoMajor != 2 ***REMOVED***
		t.Fatal("proto not h2c")
	***REMOVED***
	body, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if got, want := string(body), "Hello, /foobar, http: true"; got != want ***REMOVED***
		t.Fatalf("response got %v, want %v", got, want)
	***REMOVED***
***REMOVED***

func TestTransport(t *testing.T) ***REMOVED***
	const body = "sup"
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		io.WriteString(w, body)
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	req, err := http.NewRequest("GET", st.ts.URL, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	res, err := tr.RoundTrip(req)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer res.Body.Close()

	t.Logf("Got res: %+v", res)
	if g, w := res.StatusCode, 200; g != w ***REMOVED***
		t.Errorf("StatusCode = %v; want %v", g, w)
	***REMOVED***
	if g, w := res.Status, "200 OK"; g != w ***REMOVED***
		t.Errorf("Status = %q; want %q", g, w)
	***REMOVED***
	wantHeader := http.Header***REMOVED***
		"Content-Length": []string***REMOVED***"3"***REMOVED***,
		"Content-Type":   []string***REMOVED***"text/plain; charset=utf-8"***REMOVED***,
		"Date":           []string***REMOVED***"XXX"***REMOVED***, // see cleanDate
	***REMOVED***
	cleanDate(res)
	if !reflect.DeepEqual(res.Header, wantHeader) ***REMOVED***
		t.Errorf("res Header = %v; want %v", res.Header, wantHeader)
	***REMOVED***
	if res.Request != req ***REMOVED***
		t.Errorf("Response.Request = %p; want %p", res.Request, req)
	***REMOVED***
	if res.TLS == nil ***REMOVED***
		t.Error("Response.TLS = nil; want non-nil")
	***REMOVED***
	slurp, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		t.Errorf("Body read: %v", err)
	***REMOVED*** else if string(slurp) != body ***REMOVED***
		t.Errorf("Body = %q; want %q", slurp, body)
	***REMOVED***
***REMOVED***

func onSameConn(t *testing.T, modReq func(*http.Request)) bool ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		io.WriteString(w, r.RemoteAddr)
	***REMOVED***, optOnlyServer, func(c net.Conn, st http.ConnState) ***REMOVED***
		t.Logf("conn %v is now state %v", c.RemoteAddr(), st)
	***REMOVED***)
	defer st.Close()
	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()
	get := func() string ***REMOVED***
		req, err := http.NewRequest("GET", st.ts.URL, nil)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		modReq(req)
		res, err := tr.RoundTrip(req)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer res.Body.Close()
		slurp, err := ioutil.ReadAll(res.Body)
		if err != nil ***REMOVED***
			t.Fatalf("Body read: %v", err)
		***REMOVED***
		addr := strings.TrimSpace(string(slurp))
		if addr == "" ***REMOVED***
			t.Fatalf("didn't get an addr in response")
		***REMOVED***
		return addr
	***REMOVED***
	first := get()
	second := get()
	return first == second
***REMOVED***

func TestTransportReusesConns(t *testing.T) ***REMOVED***
	if !onSameConn(t, func(*http.Request) ***REMOVED******REMOVED***) ***REMOVED***
		t.Errorf("first and second responses were on different connections")
	***REMOVED***
***REMOVED***

func TestTransportReusesConn_RequestClose(t *testing.T) ***REMOVED***
	if onSameConn(t, func(r *http.Request) ***REMOVED*** r.Close = true ***REMOVED***) ***REMOVED***
		t.Errorf("first and second responses were not on different connections")
	***REMOVED***
***REMOVED***

func TestTransportReusesConn_ConnClose(t *testing.T) ***REMOVED***
	if onSameConn(t, func(r *http.Request) ***REMOVED*** r.Header.Set("Connection", "close") ***REMOVED***) ***REMOVED***
		t.Errorf("first and second responses were not on different connections")
	***REMOVED***
***REMOVED***

// Tests that the Transport only keeps one pending dial open per destination address.
// https://golang.org/issue/13397
func TestTransportGroupsPendingDials(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		io.WriteString(w, r.RemoteAddr)
	***REMOVED***, optOnlyServer)
	defer st.Close()
	tr := &Transport***REMOVED***
		TLSClientConfig: tlsConfigInsecure,
	***REMOVED***
	defer tr.CloseIdleConnections()
	var (
		mu    sync.Mutex
		dials = map[string]int***REMOVED******REMOVED***
	)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ ***REMOVED***
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			req, err := http.NewRequest("GET", st.ts.URL, nil)
			if err != nil ***REMOVED***
				t.Error(err)
				return
			***REMOVED***
			res, err := tr.RoundTrip(req)
			if err != nil ***REMOVED***
				t.Error(err)
				return
			***REMOVED***
			defer res.Body.Close()
			slurp, err := ioutil.ReadAll(res.Body)
			if err != nil ***REMOVED***
				t.Errorf("Body read: %v", err)
			***REMOVED***
			addr := strings.TrimSpace(string(slurp))
			if addr == "" ***REMOVED***
				t.Errorf("didn't get an addr in response")
			***REMOVED***
			mu.Lock()
			dials[addr]++
			mu.Unlock()
		***REMOVED***()
	***REMOVED***
	wg.Wait()
	if len(dials) != 1 ***REMOVED***
		t.Errorf("saw %d dials; want 1: %v", len(dials), dials)
	***REMOVED***
	tr.CloseIdleConnections()
	if err := retry(50, 10*time.Millisecond, func() error ***REMOVED***
		cp, ok := tr.connPool().(*clientConnPool)
		if !ok ***REMOVED***
			return fmt.Errorf("Conn pool is %T; want *clientConnPool", tr.connPool())
		***REMOVED***
		cp.mu.Lock()
		defer cp.mu.Unlock()
		if len(cp.dialing) != 0 ***REMOVED***
			return fmt.Errorf("dialing map = %v; want empty", cp.dialing)
		***REMOVED***
		if len(cp.conns) != 0 ***REMOVED***
			return fmt.Errorf("conns = %v; want empty", cp.conns)
		***REMOVED***
		if len(cp.keys) != 0 ***REMOVED***
			return fmt.Errorf("keys = %v; want empty", cp.keys)
		***REMOVED***
		return nil
	***REMOVED***); err != nil ***REMOVED***
		t.Errorf("State of pool after CloseIdleConnections: %v", err)
	***REMOVED***
***REMOVED***

func retry(tries int, delay time.Duration, fn func() error) error ***REMOVED***
	var err error
	for i := 0; i < tries; i++ ***REMOVED***
		err = fn()
		if err == nil ***REMOVED***
			return nil
		***REMOVED***
		time.Sleep(delay)
	***REMOVED***
	return err
***REMOVED***

func TestTransportAbortClosesPipes(t *testing.T) ***REMOVED***
	shutdown := make(chan struct***REMOVED******REMOVED***)
	st := newServerTester(t,
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			w.(http.Flusher).Flush()
			<-shutdown
		***REMOVED***,
		optOnlyServer,
	)
	defer st.Close()
	defer close(shutdown) // we must shutdown before st.Close() to avoid hanging

	done := make(chan struct***REMOVED******REMOVED***)
	requestMade := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		defer close(done)
		tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
		req, err := http.NewRequest("GET", st.ts.URL, nil)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		res, err := tr.RoundTrip(req)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer res.Body.Close()
		close(requestMade)
		_, err = ioutil.ReadAll(res.Body)
		if err == nil ***REMOVED***
			t.Error("expected error from res.Body.Read")
		***REMOVED***
	***REMOVED***()

	<-requestMade
	// Now force the serve loop to end, via closing the connection.
	st.closeConn()
	// deadlock? that's a bug.
	select ***REMOVED***
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	***REMOVED***
***REMOVED***

// TODO: merge this with TestTransportBody to make TestTransportRequest? This
// could be a table-driven test with extra goodies.
func TestTransportPath(t *testing.T) ***REMOVED***
	gotc := make(chan *url.URL, 1)
	st := newServerTester(t,
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			gotc <- r.URL
		***REMOVED***,
		optOnlyServer,
	)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()
	const (
		path  = "/testpath"
		query = "q=1"
	)
	surl := st.ts.URL + path + "?" + query
	req, err := http.NewRequest("POST", surl, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	c := &http.Client***REMOVED***Transport: tr***REMOVED***
	res, err := c.Do(req)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer res.Body.Close()
	got := <-gotc
	if got.Path != path ***REMOVED***
		t.Errorf("Read Path = %q; want %q", got.Path, path)
	***REMOVED***
	if got.RawQuery != query ***REMOVED***
		t.Errorf("Read RawQuery = %q; want %q", got.RawQuery, query)
	***REMOVED***
***REMOVED***

func randString(n int) string ***REMOVED***
	rnd := rand.New(rand.NewSource(int64(n)))
	b := make([]byte, n)
	for i := range b ***REMOVED***
		b[i] = byte(rnd.Intn(256))
	***REMOVED***
	return string(b)
***REMOVED***

type panicReader struct***REMOVED******REMOVED***

func (panicReader) Read([]byte) (int, error) ***REMOVED*** panic("unexpected Read") ***REMOVED***
func (panicReader) Close() error             ***REMOVED*** panic("unexpected Close") ***REMOVED***

func TestActualContentLength(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		req  *http.Request
		want int64
	***REMOVED******REMOVED***
		// Verify we don't read from Body:
		0: ***REMOVED***
			req:  &http.Request***REMOVED***Body: panicReader***REMOVED******REMOVED******REMOVED***,
			want: -1,
		***REMOVED***,
		// nil Body means 0, regardless of ContentLength:
		1: ***REMOVED***
			req:  &http.Request***REMOVED***Body: nil, ContentLength: 5***REMOVED***,
			want: 0,
		***REMOVED***,
		// ContentLength is used if set.
		2: ***REMOVED***
			req:  &http.Request***REMOVED***Body: panicReader***REMOVED******REMOVED***, ContentLength: 5***REMOVED***,
			want: 5,
		***REMOVED***,
		// http.NoBody means 0, not -1.
		3: ***REMOVED***
			req:  &http.Request***REMOVED***Body: go18httpNoBody()***REMOVED***,
			want: 0,
		***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		got := actualContentLength(tt.req)
		if got != tt.want ***REMOVED***
			t.Errorf("test[%d]: got %d; want %d", i, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTransportBody(t *testing.T) ***REMOVED***
	bodyTests := []struct ***REMOVED***
		body         string
		noContentLen bool
	***REMOVED******REMOVED***
		***REMOVED***body: "some message"***REMOVED***,
		***REMOVED***body: "some message", noContentLen: true***REMOVED***,
		***REMOVED***body: strings.Repeat("a", 1<<20), noContentLen: true***REMOVED***,
		***REMOVED***body: strings.Repeat("a", 1<<20)***REMOVED***,
		***REMOVED***body: randString(16<<10 - 1)***REMOVED***,
		***REMOVED***body: randString(16 << 10)***REMOVED***,
		***REMOVED***body: randString(16<<10 + 1)***REMOVED***,
		***REMOVED***body: randString(512<<10 - 1)***REMOVED***,
		***REMOVED***body: randString(512 << 10)***REMOVED***,
		***REMOVED***body: randString(512<<10 + 1)***REMOVED***,
		***REMOVED***body: randString(1<<20 - 1)***REMOVED***,
		***REMOVED***body: randString(1 << 20)***REMOVED***,
		***REMOVED***body: randString(1<<20 + 2)***REMOVED***,
	***REMOVED***

	type reqInfo struct ***REMOVED***
		req   *http.Request
		slurp []byte
		err   error
	***REMOVED***
	gotc := make(chan reqInfo, 1)
	st := newServerTester(t,
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			slurp, err := ioutil.ReadAll(r.Body)
			if err != nil ***REMOVED***
				gotc <- reqInfo***REMOVED***err: err***REMOVED***
			***REMOVED*** else ***REMOVED***
				gotc <- reqInfo***REMOVED***req: r, slurp: slurp***REMOVED***
			***REMOVED***
		***REMOVED***,
		optOnlyServer,
	)
	defer st.Close()

	for i, tt := range bodyTests ***REMOVED***
		tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
		defer tr.CloseIdleConnections()

		var body io.Reader = strings.NewReader(tt.body)
		if tt.noContentLen ***REMOVED***
			body = struct***REMOVED*** io.Reader ***REMOVED******REMOVED***body***REMOVED*** // just a Reader, hiding concrete type and other methods
		***REMOVED***
		req, err := http.NewRequest("POST", st.ts.URL, body)
		if err != nil ***REMOVED***
			t.Fatalf("#%d: %v", i, err)
		***REMOVED***
		c := &http.Client***REMOVED***Transport: tr***REMOVED***
		res, err := c.Do(req)
		if err != nil ***REMOVED***
			t.Fatalf("#%d: %v", i, err)
		***REMOVED***
		defer res.Body.Close()
		ri := <-gotc
		if ri.err != nil ***REMOVED***
			t.Errorf("#%d: read error: %v", i, ri.err)
			continue
		***REMOVED***
		if got := string(ri.slurp); got != tt.body ***REMOVED***
			t.Errorf("#%d: Read body mismatch.\n got: %q (len %d)\nwant: %q (len %d)", i, shortString(got), len(got), shortString(tt.body), len(tt.body))
		***REMOVED***
		wantLen := int64(len(tt.body))
		if tt.noContentLen && tt.body != "" ***REMOVED***
			wantLen = -1
		***REMOVED***
		if ri.req.ContentLength != wantLen ***REMOVED***
			t.Errorf("#%d. handler got ContentLength = %v; want %v", i, ri.req.ContentLength, wantLen)
		***REMOVED***
	***REMOVED***
***REMOVED***

func shortString(v string) string ***REMOVED***
	const maxLen = 100
	if len(v) <= maxLen ***REMOVED***
		return v
	***REMOVED***
	return fmt.Sprintf("%v[...%d bytes omitted...]%v", v[:maxLen/2], len(v)-maxLen, v[len(v)-maxLen/2:])
***REMOVED***

func TestTransportDialTLS(t *testing.T) ***REMOVED***
	var mu sync.Mutex // guards following
	var gotReq, didDial bool

	ts := newServerTester(t,
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			mu.Lock()
			gotReq = true
			mu.Unlock()
		***REMOVED***,
		optOnlyServer,
	)
	defer ts.Close()
	tr := &Transport***REMOVED***
		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) ***REMOVED***
			mu.Lock()
			didDial = true
			mu.Unlock()
			cfg.InsecureSkipVerify = true
			c, err := tls.Dial(netw, addr, cfg)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return c, c.Handshake()
		***REMOVED***,
	***REMOVED***
	defer tr.CloseIdleConnections()
	client := &http.Client***REMOVED***Transport: tr***REMOVED***
	res, err := client.Get(ts.ts.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	res.Body.Close()
	mu.Lock()
	if !gotReq ***REMOVED***
		t.Error("didn't get request")
	***REMOVED***
	if !didDial ***REMOVED***
		t.Error("didn't use dial hook")
	***REMOVED***
***REMOVED***

func TestConfigureTransport(t *testing.T) ***REMOVED***
	t1 := &http.Transport***REMOVED******REMOVED***
	err := ConfigureTransport(t1)
	if err == errTransportVersion ***REMOVED***
		t.Skip(err)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if got := fmt.Sprintf("%#v", t1); !strings.Contains(got, `"h2"`) ***REMOVED***
		// Laziness, to avoid buildtags.
		t.Errorf("stringification of HTTP/1 transport didn't contain \"h2\": %v", got)
	***REMOVED***
	wantNextProtos := []string***REMOVED***"h2", "http/1.1"***REMOVED***
	if t1.TLSClientConfig == nil ***REMOVED***
		t.Errorf("nil t1.TLSClientConfig")
	***REMOVED*** else if !reflect.DeepEqual(t1.TLSClientConfig.NextProtos, wantNextProtos) ***REMOVED***
		t.Errorf("TLSClientConfig.NextProtos = %q; want %q", t1.TLSClientConfig.NextProtos, wantNextProtos)
	***REMOVED***
	if err := ConfigureTransport(t1); err == nil ***REMOVED***
		t.Error("unexpected success on second call to ConfigureTransport")
	***REMOVED***

	// And does it work?
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		io.WriteString(w, r.Proto)
	***REMOVED***, optOnlyServer)
	defer st.Close()

	t1.TLSClientConfig.InsecureSkipVerify = true
	c := &http.Client***REMOVED***Transport: t1***REMOVED***
	res, err := c.Get(st.ts.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	slurp, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if got, want := string(slurp), "HTTP/2.0"; got != want ***REMOVED***
		t.Errorf("body = %q; want %q", got, want)
	***REMOVED***
***REMOVED***

type capitalizeReader struct ***REMOVED***
	r io.Reader
***REMOVED***

func (cr capitalizeReader) Read(p []byte) (n int, err error) ***REMOVED***
	n, err = cr.r.Read(p)
	for i, b := range p[:n] ***REMOVED***
		if b >= 'a' && b <= 'z' ***REMOVED***
			p[i] = b - ('a' - 'A')
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

type flushWriter struct ***REMOVED***
	w io.Writer
***REMOVED***

func (fw flushWriter) Write(p []byte) (n int, err error) ***REMOVED***
	n, err = fw.w.Write(p)
	if f, ok := fw.w.(http.Flusher); ok ***REMOVED***
		f.Flush()
	***REMOVED***
	return
***REMOVED***

type clientTester struct ***REMOVED***
	t      *testing.T
	tr     *Transport
	sc, cc net.Conn // server and client conn
	fr     *Framer  // server's framer
	client func() error
	server func() error
***REMOVED***

func newClientTester(t *testing.T) *clientTester ***REMOVED***
	var dialOnce struct ***REMOVED***
		sync.Mutex
		dialed bool
	***REMOVED***
	ct := &clientTester***REMOVED***
		t: t,
	***REMOVED***
	ct.tr = &Transport***REMOVED***
		TLSClientConfig: tlsConfigInsecure,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) ***REMOVED***
			dialOnce.Lock()
			defer dialOnce.Unlock()
			if dialOnce.dialed ***REMOVED***
				return nil, errors.New("only one dial allowed in test mode")
			***REMOVED***
			dialOnce.dialed = true
			return ct.cc, nil
		***REMOVED***,
	***REMOVED***

	ln := newLocalListener(t)
	cc, err := net.Dial("tcp", ln.Addr().String())
	if err != nil ***REMOVED***
		t.Fatal(err)

	***REMOVED***
	sc, err := ln.Accept()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ln.Close()
	ct.cc = cc
	ct.sc = sc
	ct.fr = NewFramer(sc, sc)
	return ct
***REMOVED***

func newLocalListener(t *testing.T) net.Listener ***REMOVED***
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err == nil ***REMOVED***
		return ln
	***REMOVED***
	ln, err = net.Listen("tcp6", "[::1]:0")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	return ln
***REMOVED***

func (ct *clientTester) greet(settings ...Setting) ***REMOVED***
	buf := make([]byte, len(ClientPreface))
	_, err := io.ReadFull(ct.sc, buf)
	if err != nil ***REMOVED***
		ct.t.Fatalf("reading client preface: %v", err)
	***REMOVED***
	f, err := ct.fr.ReadFrame()
	if err != nil ***REMOVED***
		ct.t.Fatalf("Reading client settings frame: %v", err)
	***REMOVED***
	if sf, ok := f.(*SettingsFrame); !ok ***REMOVED***
		ct.t.Fatalf("Wanted client settings frame; got %v", f)
		_ = sf // stash it away?
	***REMOVED***
	if err := ct.fr.WriteSettings(settings...); err != nil ***REMOVED***
		ct.t.Fatal(err)
	***REMOVED***
	if err := ct.fr.WriteSettingsAck(); err != nil ***REMOVED***
		ct.t.Fatal(err)
	***REMOVED***
***REMOVED***

func (ct *clientTester) readNonSettingsFrame() (Frame, error) ***REMOVED***
	for ***REMOVED***
		f, err := ct.fr.ReadFrame()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if _, ok := f.(*SettingsFrame); ok ***REMOVED***
			continue
		***REMOVED***
		return f, nil
	***REMOVED***
***REMOVED***

func (ct *clientTester) cleanup() ***REMOVED***
	ct.tr.CloseIdleConnections()
***REMOVED***

func (ct *clientTester) run() ***REMOVED***
	errc := make(chan error, 2)
	ct.start("client", errc, ct.client)
	ct.start("server", errc, ct.server)
	defer ct.cleanup()
	for i := 0; i < 2; i++ ***REMOVED***
		if err := <-errc; err != nil ***REMOVED***
			ct.t.Error(err)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ct *clientTester) start(which string, errc chan<- error, fn func() error) ***REMOVED***
	go func() ***REMOVED***
		finished := false
		var err error
		defer func() ***REMOVED***
			if !finished ***REMOVED***
				err = fmt.Errorf("%s goroutine didn't finish.", which)
			***REMOVED*** else if err != nil ***REMOVED***
				err = fmt.Errorf("%s: %v", which, err)
			***REMOVED***
			errc <- err
		***REMOVED***()
		err = fn()
		finished = true
	***REMOVED***()
***REMOVED***

func (ct *clientTester) readFrame() (Frame, error) ***REMOVED***
	return readFrameTimeout(ct.fr, 2*time.Second)
***REMOVED***

func (ct *clientTester) firstHeaders() (*HeadersFrame, error) ***REMOVED***
	for ***REMOVED***
		f, err := ct.readFrame()
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("ReadFrame while waiting for Headers: %v", err)
		***REMOVED***
		switch f.(type) ***REMOVED***
		case *WindowUpdateFrame, *SettingsFrame:
			continue
		***REMOVED***
		hf, ok := f.(*HeadersFrame)
		if !ok ***REMOVED***
			return nil, fmt.Errorf("Got %T; want HeadersFrame", f)
		***REMOVED***
		return hf, nil
	***REMOVED***
***REMOVED***

type countingReader struct ***REMOVED***
	n *int64
***REMOVED***

func (r countingReader) Read(p []byte) (n int, err error) ***REMOVED***
	for i := range p ***REMOVED***
		p[i] = byte(i)
	***REMOVED***
	atomic.AddInt64(r.n, int64(len(p)))
	return len(p), err
***REMOVED***

func TestTransportReqBodyAfterResponse_200(t *testing.T) ***REMOVED*** testTransportReqBodyAfterResponse(t, 200) ***REMOVED***
func TestTransportReqBodyAfterResponse_403(t *testing.T) ***REMOVED*** testTransportReqBodyAfterResponse(t, 403) ***REMOVED***

func testTransportReqBodyAfterResponse(t *testing.T, status int) ***REMOVED***
	const bodySize = 10 << 20
	clientDone := make(chan struct***REMOVED******REMOVED***)
	ct := newClientTester(t)
	ct.client = func() error ***REMOVED***
		defer ct.cc.(*net.TCPConn).CloseWrite()
		defer close(clientDone)

		var n int64 // atomic
		req, err := http.NewRequest("PUT", "https://dummy.tld/", io.LimitReader(countingReader***REMOVED***&n***REMOVED***, bodySize))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		res, err := ct.tr.RoundTrip(req)
		if err != nil ***REMOVED***
			return fmt.Errorf("RoundTrip: %v", err)
		***REMOVED***
		defer res.Body.Close()
		if res.StatusCode != status ***REMOVED***
			return fmt.Errorf("status code = %v; want %v", res.StatusCode, status)
		***REMOVED***
		slurp, err := ioutil.ReadAll(res.Body)
		if err != nil ***REMOVED***
			return fmt.Errorf("Slurp: %v", err)
		***REMOVED***
		if len(slurp) > 0 ***REMOVED***
			return fmt.Errorf("unexpected body: %q", slurp)
		***REMOVED***
		if status == 200 ***REMOVED***
			if got := atomic.LoadInt64(&n); got != bodySize ***REMOVED***
				return fmt.Errorf("For 200 response, Transport wrote %d bytes; want %d", got, bodySize)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if got := atomic.LoadInt64(&n); got == 0 || got >= bodySize ***REMOVED***
				return fmt.Errorf("For %d response, Transport wrote %d bytes; want (0,%d) exclusive", status, got, bodySize)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)
		var dataRecv int64
		var closed bool
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				select ***REMOVED***
				case <-clientDone:
					// If the client's done, it
					// will have reported any
					// errors on its side.
					return nil
				default:
					return err
				***REMOVED***
			***REMOVED***
			//println(fmt.Sprintf("server got frame: %v", f))
			switch f := f.(type) ***REMOVED***
			case *WindowUpdateFrame, *SettingsFrame:
			case *HeadersFrame:
				if !f.HeadersEnded() ***REMOVED***
					return fmt.Errorf("headers should have END_HEADERS be ended: %v", f)
				***REMOVED***
				if f.StreamEnded() ***REMOVED***
					return fmt.Errorf("headers contains END_STREAM unexpectedly: %v", f)
				***REMOVED***
			case *DataFrame:
				dataLen := len(f.Data())
				if dataLen > 0 ***REMOVED***
					if dataRecv == 0 ***REMOVED***
						enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: strconv.Itoa(status)***REMOVED***)
						ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
							StreamID:      f.StreamID,
							EndHeaders:    true,
							EndStream:     false,
							BlockFragment: buf.Bytes(),
						***REMOVED***)
					***REMOVED***
					if err := ct.fr.WriteWindowUpdate(0, uint32(dataLen)); err != nil ***REMOVED***
						return err
					***REMOVED***
					if err := ct.fr.WriteWindowUpdate(f.StreamID, uint32(dataLen)); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				dataRecv += int64(dataLen)

				if !closed && ((status != 200 && dataRecv > 0) ||
					(status == 200 && dataRecv == bodySize)) ***REMOVED***
					closed = true
					if err := ct.fr.WriteData(f.StreamID, true, nil); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			default:
				return fmt.Errorf("Unexpected client frame %v", f)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

// See golang.org/issue/13444
func TestTransportFullDuplex(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.WriteHeader(200) // redundant but for clarity
		w.(http.Flusher).Flush()
		io.Copy(flushWriter***REMOVED***w***REMOVED***, capitalizeReader***REMOVED***r.Body***REMOVED***)
		fmt.Fprintf(w, "bye.\n")
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()
	c := &http.Client***REMOVED***Transport: tr***REMOVED***

	pr, pw := io.Pipe()
	req, err := http.NewRequest("PUT", st.ts.URL, ioutil.NopCloser(pr))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	req.ContentLength = -1
	res, err := c.Do(req)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode != 200 ***REMOVED***
		t.Fatalf("StatusCode = %v; want %v", res.StatusCode, 200)
	***REMOVED***
	bs := bufio.NewScanner(res.Body)
	want := func(v string) ***REMOVED***
		if !bs.Scan() ***REMOVED***
			t.Fatalf("wanted to read %q but Scan() = false, err = %v", v, bs.Err())
		***REMOVED***
	***REMOVED***
	write := func(v string) ***REMOVED***
		_, err := io.WriteString(pw, v)
		if err != nil ***REMOVED***
			t.Fatalf("pipe write: %v", err)
		***REMOVED***
	***REMOVED***
	write("foo\n")
	want("FOO")
	write("bar\n")
	want("BAR")
	pw.Close()
	want("bye.")
	if err := bs.Err(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestTransportConnectRequest(t *testing.T) ***REMOVED***
	gotc := make(chan *http.Request, 1)
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		gotc <- r
	***REMOVED***, optOnlyServer)
	defer st.Close()

	u, err := url.Parse(st.ts.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()
	c := &http.Client***REMOVED***Transport: tr***REMOVED***

	tests := []struct ***REMOVED***
		req  *http.Request
		want string
	***REMOVED******REMOVED***
		***REMOVED***
			req: &http.Request***REMOVED***
				Method: "CONNECT",
				Header: http.Header***REMOVED******REMOVED***,
				URL:    u,
			***REMOVED***,
			want: u.Host,
		***REMOVED***,
		***REMOVED***
			req: &http.Request***REMOVED***
				Method: "CONNECT",
				Header: http.Header***REMOVED******REMOVED***,
				URL:    u,
				Host:   "example.com:123",
			***REMOVED***,
			want: "example.com:123",
		***REMOVED***,
	***REMOVED***

	for i, tt := range tests ***REMOVED***
		res, err := c.Do(tt.req)
		if err != nil ***REMOVED***
			t.Errorf("%d. RoundTrip = %v", i, err)
			continue
		***REMOVED***
		res.Body.Close()
		req := <-gotc
		if req.Method != "CONNECT" ***REMOVED***
			t.Errorf("method = %q; want CONNECT", req.Method)
		***REMOVED***
		if req.Host != tt.want ***REMOVED***
			t.Errorf("Host = %q; want %q", req.Host, tt.want)
		***REMOVED***
		if req.URL.Host != tt.want ***REMOVED***
			t.Errorf("URL.Host = %q; want %q", req.URL.Host, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

type headerType int

const (
	noHeader headerType = iota // omitted
	oneHeader
	splitHeader // broken into continuation on purpose
)

const (
	f0 = noHeader
	f1 = oneHeader
	f2 = splitHeader
	d0 = false
	d1 = true
)

// Test all 36 combinations of response frame orders:
//    (3 ways of 100-continue) * (2 ways of headers) * (2 ways of data) * (3 ways of trailers):func TestTransportResponsePattern_00f0(t *testing.T) ***REMOVED*** testTransportResponsePattern(h0, h1, false, h0) ***REMOVED***
// Generated by http://play.golang.org/p/SScqYKJYXd
func TestTransportResPattern_c0h1d0t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f1, d0, f0) ***REMOVED***
func TestTransportResPattern_c0h1d0t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f1, d0, f1) ***REMOVED***
func TestTransportResPattern_c0h1d0t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f1, d0, f2) ***REMOVED***
func TestTransportResPattern_c0h1d1t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f1, d1, f0) ***REMOVED***
func TestTransportResPattern_c0h1d1t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f1, d1, f1) ***REMOVED***
func TestTransportResPattern_c0h1d1t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f1, d1, f2) ***REMOVED***
func TestTransportResPattern_c0h2d0t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f2, d0, f0) ***REMOVED***
func TestTransportResPattern_c0h2d0t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f2, d0, f1) ***REMOVED***
func TestTransportResPattern_c0h2d0t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f2, d0, f2) ***REMOVED***
func TestTransportResPattern_c0h2d1t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f2, d1, f0) ***REMOVED***
func TestTransportResPattern_c0h2d1t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f2, d1, f1) ***REMOVED***
func TestTransportResPattern_c0h2d1t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f0, f2, d1, f2) ***REMOVED***
func TestTransportResPattern_c1h1d0t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f1, d0, f0) ***REMOVED***
func TestTransportResPattern_c1h1d0t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f1, d0, f1) ***REMOVED***
func TestTransportResPattern_c1h1d0t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f1, d0, f2) ***REMOVED***
func TestTransportResPattern_c1h1d1t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f1, d1, f0) ***REMOVED***
func TestTransportResPattern_c1h1d1t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f1, d1, f1) ***REMOVED***
func TestTransportResPattern_c1h1d1t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f1, d1, f2) ***REMOVED***
func TestTransportResPattern_c1h2d0t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f2, d0, f0) ***REMOVED***
func TestTransportResPattern_c1h2d0t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f2, d0, f1) ***REMOVED***
func TestTransportResPattern_c1h2d0t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f2, d0, f2) ***REMOVED***
func TestTransportResPattern_c1h2d1t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f2, d1, f0) ***REMOVED***
func TestTransportResPattern_c1h2d1t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f2, d1, f1) ***REMOVED***
func TestTransportResPattern_c1h2d1t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f1, f2, d1, f2) ***REMOVED***
func TestTransportResPattern_c2h1d0t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f1, d0, f0) ***REMOVED***
func TestTransportResPattern_c2h1d0t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f1, d0, f1) ***REMOVED***
func TestTransportResPattern_c2h1d0t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f1, d0, f2) ***REMOVED***
func TestTransportResPattern_c2h1d1t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f1, d1, f0) ***REMOVED***
func TestTransportResPattern_c2h1d1t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f1, d1, f1) ***REMOVED***
func TestTransportResPattern_c2h1d1t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f1, d1, f2) ***REMOVED***
func TestTransportResPattern_c2h2d0t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f2, d0, f0) ***REMOVED***
func TestTransportResPattern_c2h2d0t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f2, d0, f1) ***REMOVED***
func TestTransportResPattern_c2h2d0t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f2, d0, f2) ***REMOVED***
func TestTransportResPattern_c2h2d1t0(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f2, d1, f0) ***REMOVED***
func TestTransportResPattern_c2h2d1t1(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f2, d1, f1) ***REMOVED***
func TestTransportResPattern_c2h2d1t2(t *testing.T) ***REMOVED*** testTransportResPattern(t, f2, f2, d1, f2) ***REMOVED***

func testTransportResPattern(t *testing.T, expect100Continue, resHeader headerType, withData bool, trailers headerType) ***REMOVED***
	const reqBody = "some request body"
	const resBody = "some response body"

	if resHeader == noHeader ***REMOVED***
		// TODO: test 100-continue followed by immediate
		// server stream reset, without headers in the middle?
		panic("invalid combination")
	***REMOVED***

	ct := newClientTester(t)
	ct.client = func() error ***REMOVED***
		req, _ := http.NewRequest("POST", "https://dummy.tld/", strings.NewReader(reqBody))
		if expect100Continue != noHeader ***REMOVED***
			req.Header.Set("Expect", "100-continue")
		***REMOVED***
		res, err := ct.tr.RoundTrip(req)
		if err != nil ***REMOVED***
			return fmt.Errorf("RoundTrip: %v", err)
		***REMOVED***
		defer res.Body.Close()
		if res.StatusCode != 200 ***REMOVED***
			return fmt.Errorf("status code = %v; want 200", res.StatusCode)
		***REMOVED***
		slurp, err := ioutil.ReadAll(res.Body)
		if err != nil ***REMOVED***
			return fmt.Errorf("Slurp: %v", err)
		***REMOVED***
		wantBody := resBody
		if !withData ***REMOVED***
			wantBody = ""
		***REMOVED***
		if string(slurp) != wantBody ***REMOVED***
			return fmt.Errorf("body = %q; want %q", slurp, wantBody)
		***REMOVED***
		if trailers == noHeader ***REMOVED***
			if len(res.Trailer) > 0 ***REMOVED***
				t.Errorf("Trailer = %v; want none", res.Trailer)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			want := http.Header***REMOVED***"Some-Trailer": ***REMOVED***"some-value"***REMOVED******REMOVED***
			if !reflect.DeepEqual(res.Trailer, want) ***REMOVED***
				t.Errorf("Trailer = %v; want %v", res.Trailer, want)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)

		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			endStream := false
			send := func(mode headerType) ***REMOVED***
				hbf := buf.Bytes()
				switch mode ***REMOVED***
				case oneHeader:
					ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
						StreamID:      f.Header().StreamID,
						EndHeaders:    true,
						EndStream:     endStream,
						BlockFragment: hbf,
					***REMOVED***)
				case splitHeader:
					if len(hbf) < 2 ***REMOVED***
						panic("too small")
					***REMOVED***
					ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
						StreamID:      f.Header().StreamID,
						EndHeaders:    false,
						EndStream:     endStream,
						BlockFragment: hbf[:1],
					***REMOVED***)
					ct.fr.WriteContinuation(f.Header().StreamID, true, hbf[1:])
				default:
					panic("bogus mode")
				***REMOVED***
			***REMOVED***
			switch f := f.(type) ***REMOVED***
			case *WindowUpdateFrame, *SettingsFrame:
			case *DataFrame:
				if !f.StreamEnded() ***REMOVED***
					// No need to send flow control tokens. The test request body is tiny.
					continue
				***REMOVED***
				// Response headers (1+ frames; 1 or 2 in this test, but never 0)
				***REMOVED***
					buf.Reset()
					enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
					enc.WriteField(hpack.HeaderField***REMOVED***Name: "x-foo", Value: "blah"***REMOVED***)
					enc.WriteField(hpack.HeaderField***REMOVED***Name: "x-bar", Value: "more"***REMOVED***)
					if trailers != noHeader ***REMOVED***
						enc.WriteField(hpack.HeaderField***REMOVED***Name: "trailer", Value: "some-trailer"***REMOVED***)
					***REMOVED***
					endStream = withData == false && trailers == noHeader
					send(resHeader)
				***REMOVED***
				if withData ***REMOVED***
					endStream = trailers == noHeader
					ct.fr.WriteData(f.StreamID, endStream, []byte(resBody))
				***REMOVED***
				if trailers != noHeader ***REMOVED***
					endStream = true
					buf.Reset()
					enc.WriteField(hpack.HeaderField***REMOVED***Name: "some-trailer", Value: "some-value"***REMOVED***)
					send(trailers)
				***REMOVED***
				if endStream ***REMOVED***
					return nil
				***REMOVED***
			case *HeadersFrame:
				if expect100Continue != noHeader ***REMOVED***
					buf.Reset()
					enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "100"***REMOVED***)
					send(expect100Continue)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

func TestTransportReceiveUndeclaredTrailer(t *testing.T) ***REMOVED***
	ct := newClientTester(t)
	ct.client = func() error ***REMOVED***
		req, _ := http.NewRequest("GET", "https://dummy.tld/", nil)
		res, err := ct.tr.RoundTrip(req)
		if err != nil ***REMOVED***
			return fmt.Errorf("RoundTrip: %v", err)
		***REMOVED***
		defer res.Body.Close()
		if res.StatusCode != 200 ***REMOVED***
			return fmt.Errorf("status code = %v; want 200", res.StatusCode)
		***REMOVED***
		slurp, err := ioutil.ReadAll(res.Body)
		if err != nil ***REMOVED***
			return fmt.Errorf("res.Body ReadAll error = %q, %v; want %v", slurp, err, nil)
		***REMOVED***
		if len(slurp) > 0 ***REMOVED***
			return fmt.Errorf("body = %q; want nothing", slurp)
		***REMOVED***
		if _, ok := res.Trailer["Some-Trailer"]; !ok ***REMOVED***
			return fmt.Errorf("expected Some-Trailer")
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()

		var n int
		var hf *HeadersFrame
		for hf == nil && n < 10 ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hf, _ = f.(*HeadersFrame)
			n++
		***REMOVED***

		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)

		// send headers without Trailer header
		enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
		ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
			StreamID:      hf.StreamID,
			EndHeaders:    true,
			EndStream:     false,
			BlockFragment: buf.Bytes(),
		***REMOVED***)

		// send trailers
		buf.Reset()
		enc.WriteField(hpack.HeaderField***REMOVED***Name: "some-trailer", Value: "I'm an undeclared Trailer!"***REMOVED***)
		ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
			StreamID:      hf.StreamID,
			EndHeaders:    true,
			EndStream:     true,
			BlockFragment: buf.Bytes(),
		***REMOVED***)
		return nil
	***REMOVED***
	ct.run()
***REMOVED***

func TestTransportInvalidTrailer_Pseudo1(t *testing.T) ***REMOVED***
	testTransportInvalidTrailer_Pseudo(t, oneHeader)
***REMOVED***
func TestTransportInvalidTrailer_Pseudo2(t *testing.T) ***REMOVED***
	testTransportInvalidTrailer_Pseudo(t, splitHeader)
***REMOVED***
func testTransportInvalidTrailer_Pseudo(t *testing.T, trailers headerType) ***REMOVED***
	testInvalidTrailer(t, trailers, pseudoHeaderError(":colon"), func(enc *hpack.Encoder) ***REMOVED***
		enc.WriteField(hpack.HeaderField***REMOVED***Name: ":colon", Value: "foo"***REMOVED***)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: "foo", Value: "bar"***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestTransportInvalidTrailer_Capital1(t *testing.T) ***REMOVED***
	testTransportInvalidTrailer_Capital(t, oneHeader)
***REMOVED***
func TestTransportInvalidTrailer_Capital2(t *testing.T) ***REMOVED***
	testTransportInvalidTrailer_Capital(t, splitHeader)
***REMOVED***
func testTransportInvalidTrailer_Capital(t *testing.T, trailers headerType) ***REMOVED***
	testInvalidTrailer(t, trailers, headerFieldNameError("Capital"), func(enc *hpack.Encoder) ***REMOVED***
		enc.WriteField(hpack.HeaderField***REMOVED***Name: "foo", Value: "bar"***REMOVED***)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: "Capital", Value: "bad"***REMOVED***)
	***REMOVED***)
***REMOVED***
func TestTransportInvalidTrailer_EmptyFieldName(t *testing.T) ***REMOVED***
	testInvalidTrailer(t, oneHeader, headerFieldNameError(""), func(enc *hpack.Encoder) ***REMOVED***
		enc.WriteField(hpack.HeaderField***REMOVED***Name: "", Value: "bad"***REMOVED***)
	***REMOVED***)
***REMOVED***
func TestTransportInvalidTrailer_BinaryFieldValue(t *testing.T) ***REMOVED***
	testInvalidTrailer(t, oneHeader, headerFieldValueError("has\nnewline"), func(enc *hpack.Encoder) ***REMOVED***
		enc.WriteField(hpack.HeaderField***REMOVED***Name: "x", Value: "has\nnewline"***REMOVED***)
	***REMOVED***)
***REMOVED***

func testInvalidTrailer(t *testing.T, trailers headerType, wantErr error, writeTrailer func(*hpack.Encoder)) ***REMOVED***
	ct := newClientTester(t)
	ct.client = func() error ***REMOVED***
		req, _ := http.NewRequest("GET", "https://dummy.tld/", nil)
		res, err := ct.tr.RoundTrip(req)
		if err != nil ***REMOVED***
			return fmt.Errorf("RoundTrip: %v", err)
		***REMOVED***
		defer res.Body.Close()
		if res.StatusCode != 200 ***REMOVED***
			return fmt.Errorf("status code = %v; want 200", res.StatusCode)
		***REMOVED***
		slurp, err := ioutil.ReadAll(res.Body)
		se, ok := err.(StreamError)
		if !ok || se.Cause != wantErr ***REMOVED***
			return fmt.Errorf("res.Body ReadAll error = %q, %#v; want StreamError with cause %T, %#v", slurp, err, wantErr, wantErr)
		***REMOVED***
		if len(slurp) > 0 ***REMOVED***
			return fmt.Errorf("body = %q; want nothing", slurp)
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)

		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			switch f := f.(type) ***REMOVED***
			case *HeadersFrame:
				var endStream bool
				send := func(mode headerType) ***REMOVED***
					hbf := buf.Bytes()
					switch mode ***REMOVED***
					case oneHeader:
						ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
							StreamID:      f.StreamID,
							EndHeaders:    true,
							EndStream:     endStream,
							BlockFragment: hbf,
						***REMOVED***)
					case splitHeader:
						if len(hbf) < 2 ***REMOVED***
							panic("too small")
						***REMOVED***
						ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
							StreamID:      f.StreamID,
							EndHeaders:    false,
							EndStream:     endStream,
							BlockFragment: hbf[:1],
						***REMOVED***)
						ct.fr.WriteContinuation(f.StreamID, true, hbf[1:])
					default:
						panic("bogus mode")
					***REMOVED***
				***REMOVED***
				// Response headers (1+ frames; 1 or 2 in this test, but never 0)
				***REMOVED***
					buf.Reset()
					enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
					enc.WriteField(hpack.HeaderField***REMOVED***Name: "trailer", Value: "declared"***REMOVED***)
					endStream = false
					send(oneHeader)
				***REMOVED***
				// Trailers:
				***REMOVED***
					endStream = true
					buf.Reset()
					writeTrailer(enc)
					send(trailers)
				***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

// headerListSize returns the HTTP2 header list size of h.
//   http://httpwg.org/specs/rfc7540.html#SETTINGS_MAX_HEADER_LIST_SIZE
//   http://httpwg.org/specs/rfc7540.html#MaxHeaderBlock
func headerListSize(h http.Header) (size uint32) ***REMOVED***
	for k, vv := range h ***REMOVED***
		for _, v := range vv ***REMOVED***
			hf := hpack.HeaderField***REMOVED***Name: k, Value: v***REMOVED***
			size += hf.Size()
		***REMOVED***
	***REMOVED***
	return size
***REMOVED***

// padHeaders adds data to an http.Header until headerListSize(h) ==
// limit. Due to the way header list sizes are calculated, padHeaders
// cannot add fewer than len("Pad-Headers") + 32 bytes to h, and will
// call t.Fatal if asked to do so. PadHeaders first reserves enough
// space for an empty "Pad-Headers" key, then adds as many copies of
// filler as possible. Any remaining bytes necessary to push the
// header list size up to limit are added to h["Pad-Headers"].
func padHeaders(t *testing.T, h http.Header, limit uint64, filler string) ***REMOVED***
	if limit > 0xffffffff ***REMOVED***
		t.Fatalf("padHeaders: refusing to pad to more than 2^32-1 bytes. limit = %v", limit)
	***REMOVED***
	hf := hpack.HeaderField***REMOVED***Name: "Pad-Headers", Value: ""***REMOVED***
	minPadding := uint64(hf.Size())
	size := uint64(headerListSize(h))

	minlimit := size + minPadding
	if limit < minlimit ***REMOVED***
		t.Fatalf("padHeaders: limit %v < %v", limit, minlimit)
	***REMOVED***

	// Use a fixed-width format for name so that fieldSize
	// remains constant.
	nameFmt := "Pad-Headers-%06d"
	hf = hpack.HeaderField***REMOVED***Name: fmt.Sprintf(nameFmt, 1), Value: filler***REMOVED***
	fieldSize := uint64(hf.Size())

	// Add as many complete filler values as possible, leaving
	// room for at least one empty "Pad-Headers" key.
	limit = limit - minPadding
	for i := 0; size+fieldSize < limit; i++ ***REMOVED***
		name := fmt.Sprintf(nameFmt, i)
		h.Add(name, filler)
		size += fieldSize
	***REMOVED***

	// Add enough bytes to reach limit.
	remain := limit - size
	lastValue := strings.Repeat("*", int(remain))
	h.Add("Pad-Headers", lastValue)
***REMOVED***

func TestPadHeaders(t *testing.T) ***REMOVED***
	check := func(h http.Header, limit uint32, fillerLen int) ***REMOVED***
		if h == nil ***REMOVED***
			h = make(http.Header)
		***REMOVED***
		filler := strings.Repeat("f", fillerLen)
		padHeaders(t, h, uint64(limit), filler)
		gotSize := headerListSize(h)
		if gotSize != limit ***REMOVED***
			t.Errorf("Got size = %v; want %v", gotSize, limit)
		***REMOVED***
	***REMOVED***
	// Try all possible combinations for small fillerLen and limit.
	hf := hpack.HeaderField***REMOVED***Name: "Pad-Headers", Value: ""***REMOVED***
	minLimit := hf.Size()
	for limit := minLimit; limit <= 128; limit++ ***REMOVED***
		for fillerLen := 0; uint32(fillerLen) <= limit; fillerLen++ ***REMOVED***
			check(nil, limit, fillerLen)
		***REMOVED***
	***REMOVED***

	// Try a few tests with larger limits, plus cumulative
	// tests. Since these tests are cumulative, tests[i+1].limit
	// must be >= tests[i].limit + minLimit. See the comment on
	// padHeaders for more info on why the limit arg has this
	// restriction.
	tests := []struct ***REMOVED***
		fillerLen int
		limit     uint32
	***REMOVED******REMOVED***
		***REMOVED***
			fillerLen: 64,
			limit:     1024,
		***REMOVED***,
		***REMOVED***
			fillerLen: 1024,
			limit:     1286,
		***REMOVED***,
		***REMOVED***
			fillerLen: 256,
			limit:     2048,
		***REMOVED***,
		***REMOVED***
			fillerLen: 1024,
			limit:     10 * 1024,
		***REMOVED***,
		***REMOVED***
			fillerLen: 1023,
			limit:     11 * 1024,
		***REMOVED***,
	***REMOVED***
	h := make(http.Header)
	for _, tc := range tests ***REMOVED***
		check(nil, tc.limit, tc.fillerLen)
		check(h, tc.limit, tc.fillerLen)
	***REMOVED***
***REMOVED***

func TestTransportChecksRequestHeaderListSize(t *testing.T) ***REMOVED***
	st := newServerTester(t,
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			// Consume body & force client to send
			// trailers before writing response.
			// ioutil.ReadAll returns non-nil err for
			// requests that attempt to send greater than
			// maxHeaderListSize bytes of trailers, since
			// those requests generate a stream reset.
			ioutil.ReadAll(r.Body)
			r.Body.Close()
		***REMOVED***,
		func(ts *httptest.Server) ***REMOVED***
			ts.Config.MaxHeaderBytes = 16 << 10
		***REMOVED***,
		optOnlyServer,
		optQuiet,
	)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	checkRoundTrip := func(req *http.Request, wantErr error, desc string) ***REMOVED***
		res, err := tr.RoundTrip(req)
		if err != wantErr ***REMOVED***
			if res != nil ***REMOVED***
				res.Body.Close()
			***REMOVED***
			t.Errorf("%v: RoundTrip err = %v; want %v", desc, err, wantErr)
			return
		***REMOVED***
		if err == nil ***REMOVED***
			if res == nil ***REMOVED***
				t.Errorf("%v: response nil; want non-nil.", desc)
				return
			***REMOVED***
			defer res.Body.Close()
			if res.StatusCode != http.StatusOK ***REMOVED***
				t.Errorf("%v: response status = %v; want %v", desc, res.StatusCode, http.StatusOK)
			***REMOVED***
			return
		***REMOVED***
		if res != nil ***REMOVED***
			t.Errorf("%v: RoundTrip err = %v but response non-nil", desc, err)
		***REMOVED***
	***REMOVED***
	headerListSizeForRequest := func(req *http.Request) (size uint64) ***REMOVED***
		contentLen := actualContentLength(req)
		trailers, err := commaSeparatedTrailers(req)
		if err != nil ***REMOVED***
			t.Fatalf("headerListSizeForRequest: %v", err)
		***REMOVED***
		cc := &ClientConn***REMOVED***peerMaxHeaderListSize: 0xffffffffffffffff***REMOVED***
		cc.henc = hpack.NewEncoder(&cc.hbuf)
		cc.mu.Lock()
		hdrs, err := cc.encodeHeaders(req, true, trailers, contentLen)
		cc.mu.Unlock()
		if err != nil ***REMOVED***
			t.Fatalf("headerListSizeForRequest: %v", err)
		***REMOVED***
		hpackDec := hpack.NewDecoder(initialHeaderTableSize, func(hf hpack.HeaderField) ***REMOVED***
			size += uint64(hf.Size())
		***REMOVED***)
		if len(hdrs) > 0 ***REMOVED***
			if _, err := hpackDec.Write(hdrs); err != nil ***REMOVED***
				t.Fatalf("headerListSizeForRequest: %v", err)
			***REMOVED***
		***REMOVED***
		return size
	***REMOVED***
	// Create a new Request for each test, rather than reusing the
	// same Request, to avoid a race when modifying req.Headers.
	// See https://github.com/golang/go/issues/21316
	newRequest := func() *http.Request ***REMOVED***
		// Body must be non-nil to enable writing trailers.
		body := strings.NewReader("hello")
		req, err := http.NewRequest("POST", st.ts.URL, body)
		if err != nil ***REMOVED***
			t.Fatalf("newRequest: NewRequest: %v", err)
		***REMOVED***
		return req
	***REMOVED***

	// Make an arbitrary request to ensure we get the server's
	// settings frame and initialize peerMaxHeaderListSize.
	req := newRequest()
	checkRoundTrip(req, nil, "Initial request")

	// Get the ClientConn associated with the request and validate
	// peerMaxHeaderListSize.
	addr := authorityAddr(req.URL.Scheme, req.URL.Host)
	cc, err := tr.connPool().GetClientConn(req, addr)
	if err != nil ***REMOVED***
		t.Fatalf("GetClientConn: %v", err)
	***REMOVED***
	cc.mu.Lock()
	peerSize := cc.peerMaxHeaderListSize
	cc.mu.Unlock()
	st.scMu.Lock()
	wantSize := uint64(st.sc.maxHeaderListSize())
	st.scMu.Unlock()
	if peerSize != wantSize ***REMOVED***
		t.Errorf("peerMaxHeaderListSize = %v; want %v", peerSize, wantSize)
	***REMOVED***

	// Sanity check peerSize. (*serverConn) maxHeaderListSize adds
	// 320 bytes of padding.
	wantHeaderBytes := uint64(st.ts.Config.MaxHeaderBytes) + 320
	if peerSize != wantHeaderBytes ***REMOVED***
		t.Errorf("peerMaxHeaderListSize = %v; want %v.", peerSize, wantHeaderBytes)
	***REMOVED***

	// Pad headers & trailers, but stay under peerSize.
	req = newRequest()
	req.Header = make(http.Header)
	req.Trailer = make(http.Header)
	filler := strings.Repeat("*", 1024)
	padHeaders(t, req.Trailer, peerSize, filler)
	// cc.encodeHeaders adds some default headers to the request,
	// so we need to leave room for those.
	defaultBytes := headerListSizeForRequest(req)
	padHeaders(t, req.Header, peerSize-defaultBytes, filler)
	checkRoundTrip(req, nil, "Headers & Trailers under limit")

	// Add enough header bytes to push us over peerSize.
	req = newRequest()
	req.Header = make(http.Header)
	padHeaders(t, req.Header, peerSize, filler)
	checkRoundTrip(req, errRequestHeaderListSize, "Headers over limit")

	// Push trailers over the limit.
	req = newRequest()
	req.Trailer = make(http.Header)
	padHeaders(t, req.Trailer, peerSize+1, filler)
	checkRoundTrip(req, errRequestHeaderListSize, "Trailers over limit")

	// Send headers with a single large value.
	req = newRequest()
	filler = strings.Repeat("*", int(peerSize))
	req.Header = make(http.Header)
	req.Header.Set("Big", filler)
	checkRoundTrip(req, errRequestHeaderListSize, "Single large header")

	// Send trailers with a single large value.
	req = newRequest()
	req.Trailer = make(http.Header)
	req.Trailer.Set("Big", filler)
	checkRoundTrip(req, errRequestHeaderListSize, "Single large trailer")
***REMOVED***

func TestTransportChecksResponseHeaderListSize(t *testing.T) ***REMOVED***
	ct := newClientTester(t)
	ct.client = func() error ***REMOVED***
		req, _ := http.NewRequest("GET", "https://dummy.tld/", nil)
		res, err := ct.tr.RoundTrip(req)
		if err != errResponseHeaderListSize ***REMOVED***
			if res != nil ***REMOVED***
				res.Body.Close()
			***REMOVED***
			size := int64(0)
			for k, vv := range res.Header ***REMOVED***
				for _, v := range vv ***REMOVED***
					size += int64(len(k)) + int64(len(v)) + 32
				***REMOVED***
			***REMOVED***
			return fmt.Errorf("RoundTrip Error = %v (and %d bytes of response headers); want errResponseHeaderListSize", err, size)
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)

		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			switch f := f.(type) ***REMOVED***
			case *HeadersFrame:
				enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
				large := strings.Repeat("a", 1<<10)
				for i := 0; i < 5042; i++ ***REMOVED***
					enc.WriteField(hpack.HeaderField***REMOVED***Name: large, Value: large***REMOVED***)
				***REMOVED***
				if size, want := buf.Len(), 6329; size != want ***REMOVED***
					// Note: this number might change if
					// our hpack implementation
					// changes. That's fine. This is
					// just a sanity check that our
					// response can fit in a single
					// header block fragment frame.
					return fmt.Errorf("encoding over 10MB of duplicate keypairs took %d bytes; expected %d", size, want)
				***REMOVED***
				ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
					StreamID:      f.StreamID,
					EndHeaders:    true,
					EndStream:     true,
					BlockFragment: buf.Bytes(),
				***REMOVED***)
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

// Test that the the Transport returns a typed error from Response.Body.Read calls
// when the server sends an error. (here we use a panic, since that should generate
// a stream error, but others like cancel should be similar)
func TestTransportBodyReadErrorType(t *testing.T) ***REMOVED***
	doPanic := make(chan bool, 1)
	st := newServerTester(t,
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			w.(http.Flusher).Flush() // force headers out
			<-doPanic
			panic("boom")
		***REMOVED***,
		optOnlyServer,
		optQuiet,
	)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()
	c := &http.Client***REMOVED***Transport: tr***REMOVED***

	res, err := c.Get(st.ts.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer res.Body.Close()
	doPanic <- true
	buf := make([]byte, 100)
	n, err := res.Body.Read(buf)
	want := StreamError***REMOVED***StreamID: 0x1, Code: 0x2***REMOVED***
	if !reflect.DeepEqual(want, err) ***REMOVED***
		t.Errorf("Read = %v, %#v; want error %#v", n, err, want)
	***REMOVED***
***REMOVED***

// golang.org/issue/13924
// This used to fail after many iterations, especially with -race:
// go test -v -run=TestTransportDoubleCloseOnWriteError -count=500 -race
func TestTransportDoubleCloseOnWriteError(t *testing.T) ***REMOVED***
	var (
		mu   sync.Mutex
		conn net.Conn // to close if set
	)

	st := newServerTester(t,
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			mu.Lock()
			defer mu.Unlock()
			if conn != nil ***REMOVED***
				conn.Close()
			***REMOVED***
		***REMOVED***,
		optOnlyServer,
	)
	defer st.Close()

	tr := &Transport***REMOVED***
		TLSClientConfig: tlsConfigInsecure,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) ***REMOVED***
			tc, err := tls.Dial(network, addr, cfg)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			mu.Lock()
			defer mu.Unlock()
			conn = tc
			return tc, nil
		***REMOVED***,
	***REMOVED***
	defer tr.CloseIdleConnections()
	c := &http.Client***REMOVED***Transport: tr***REMOVED***
	c.Get(st.ts.URL)
***REMOVED***

// Test that the http1 Transport.DisableKeepAlives option is respected
// and connections are closed as soon as idle.
// See golang.org/issue/14008
func TestTransportDisableKeepAlives(t *testing.T) ***REMOVED***
	st := newServerTester(t,
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			io.WriteString(w, "hi")
		***REMOVED***,
		optOnlyServer,
	)
	defer st.Close()

	connClosed := make(chan struct***REMOVED******REMOVED***) // closed on tls.Conn.Close
	tr := &Transport***REMOVED***
		t1: &http.Transport***REMOVED***
			DisableKeepAlives: true,
		***REMOVED***,
		TLSClientConfig: tlsConfigInsecure,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) ***REMOVED***
			tc, err := tls.Dial(network, addr, cfg)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &noteCloseConn***REMOVED***Conn: tc, closefn: func() ***REMOVED*** close(connClosed) ***REMOVED******REMOVED***, nil
		***REMOVED***,
	***REMOVED***
	c := &http.Client***REMOVED***Transport: tr***REMOVED***
	res, err := c.Get(st.ts.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := ioutil.ReadAll(res.Body); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer res.Body.Close()

	select ***REMOVED***
	case <-connClosed:
	case <-time.After(1 * time.Second):
		t.Errorf("timeout")
	***REMOVED***

***REMOVED***

// Test concurrent requests with Transport.DisableKeepAlives. We can share connections,
// but when things are totally idle, it still needs to close.
func TestTransportDisableKeepAlives_Concurrency(t *testing.T) ***REMOVED***
	const D = 25 * time.Millisecond
	st := newServerTester(t,
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			time.Sleep(D)
			io.WriteString(w, "hi")
		***REMOVED***,
		optOnlyServer,
	)
	defer st.Close()

	var dials int32
	var conns sync.WaitGroup
	tr := &Transport***REMOVED***
		t1: &http.Transport***REMOVED***
			DisableKeepAlives: true,
		***REMOVED***,
		TLSClientConfig: tlsConfigInsecure,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) ***REMOVED***
			tc, err := tls.Dial(network, addr, cfg)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			atomic.AddInt32(&dials, 1)
			conns.Add(1)
			return &noteCloseConn***REMOVED***Conn: tc, closefn: func() ***REMOVED*** conns.Done() ***REMOVED******REMOVED***, nil
		***REMOVED***,
	***REMOVED***
	c := &http.Client***REMOVED***Transport: tr***REMOVED***
	var reqs sync.WaitGroup
	const N = 20
	for i := 0; i < N; i++ ***REMOVED***
		reqs.Add(1)
		if i == N-1 ***REMOVED***
			// For the final request, try to make all the
			// others close. This isn't verified in the
			// count, other than the Log statement, since
			// it's so timing dependent. This test is
			// really to make sure we don't interrupt a
			// valid request.
			time.Sleep(D * 2)
		***REMOVED***
		go func() ***REMOVED***
			defer reqs.Done()
			res, err := c.Get(st.ts.URL)
			if err != nil ***REMOVED***
				t.Error(err)
				return
			***REMOVED***
			if _, err := ioutil.ReadAll(res.Body); err != nil ***REMOVED***
				t.Error(err)
				return
			***REMOVED***
			res.Body.Close()
		***REMOVED***()
	***REMOVED***
	reqs.Wait()
	conns.Wait()
	t.Logf("did %d dials, %d requests", atomic.LoadInt32(&dials), N)
***REMOVED***

type noteCloseConn struct ***REMOVED***
	net.Conn
	onceClose sync.Once
	closefn   func()
***REMOVED***

func (c *noteCloseConn) Close() error ***REMOVED***
	c.onceClose.Do(c.closefn)
	return c.Conn.Close()
***REMOVED***

func isTimeout(err error) bool ***REMOVED***
	switch err := err.(type) ***REMOVED***
	case nil:
		return false
	case *url.Error:
		return isTimeout(err.Err)
	case net.Error:
		return err.Timeout()
	***REMOVED***
	return false
***REMOVED***

// Test that the http1 Transport.ResponseHeaderTimeout option and cancel is sent.
func TestTransportResponseHeaderTimeout_NoBody(t *testing.T) ***REMOVED***
	testTransportResponseHeaderTimeout(t, false)
***REMOVED***
func TestTransportResponseHeaderTimeout_Body(t *testing.T) ***REMOVED***
	testTransportResponseHeaderTimeout(t, true)
***REMOVED***

func testTransportResponseHeaderTimeout(t *testing.T, body bool) ***REMOVED***
	ct := newClientTester(t)
	ct.tr.t1 = &http.Transport***REMOVED***
		ResponseHeaderTimeout: 5 * time.Millisecond,
	***REMOVED***
	ct.client = func() error ***REMOVED***
		c := &http.Client***REMOVED***Transport: ct.tr***REMOVED***
		var err error
		var n int64
		const bodySize = 4 << 20
		if body ***REMOVED***
			_, err = c.Post("https://dummy.tld/", "text/foo", io.LimitReader(countingReader***REMOVED***&n***REMOVED***, bodySize))
		***REMOVED*** else ***REMOVED***
			_, err = c.Get("https://dummy.tld/")
		***REMOVED***
		if !isTimeout(err) ***REMOVED***
			t.Errorf("client expected timeout error; got %#v", err)
		***REMOVED***
		if body && n != bodySize ***REMOVED***
			t.Errorf("only read %d bytes of body; want %d", n, bodySize)
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				t.Logf("ReadFrame: %v", err)
				return nil
			***REMOVED***
			switch f := f.(type) ***REMOVED***
			case *DataFrame:
				dataLen := len(f.Data())
				if dataLen > 0 ***REMOVED***
					if err := ct.fr.WriteWindowUpdate(0, uint32(dataLen)); err != nil ***REMOVED***
						return err
					***REMOVED***
					if err := ct.fr.WriteWindowUpdate(f.StreamID, uint32(dataLen)); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			case *RSTStreamFrame:
				if f.StreamID == 1 && f.ErrCode == ErrCodeCancel ***REMOVED***
					return nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

func TestTransportDisableCompression(t *testing.T) ***REMOVED***
	const body = "sup"
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		want := http.Header***REMOVED***
			"User-Agent": []string***REMOVED***"Go-http-client/2.0"***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(r.Header, want) ***REMOVED***
			t.Errorf("request headers = %v; want %v", r.Header, want)
		***REMOVED***
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***
		TLSClientConfig: tlsConfigInsecure,
		t1: &http.Transport***REMOVED***
			DisableCompression: true,
		***REMOVED***,
	***REMOVED***
	defer tr.CloseIdleConnections()

	req, err := http.NewRequest("GET", st.ts.URL, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	res, err := tr.RoundTrip(req)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer res.Body.Close()
***REMOVED***

// RFC 7540 section 8.1.2.2
func TestTransportRejectsConnHeaders(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		var got []string
		for k := range r.Header ***REMOVED***
			got = append(got, k)
		***REMOVED***
		sort.Strings(got)
		w.Header().Set("Got-Header", strings.Join(got, ","))
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	tests := []struct ***REMOVED***
		key   string
		value []string
		want  string
	***REMOVED******REMOVED***
		***REMOVED***
			key:   "Upgrade",
			value: []string***REMOVED***"anything"***REMOVED***,
			want:  "ERROR: http2: invalid Upgrade request header: [\"anything\"]",
		***REMOVED***,
		***REMOVED***
			key:   "Connection",
			value: []string***REMOVED***"foo"***REMOVED***,
			want:  "ERROR: http2: invalid Connection request header: [\"foo\"]",
		***REMOVED***,
		***REMOVED***
			key:   "Connection",
			value: []string***REMOVED***"close"***REMOVED***,
			want:  "Accept-Encoding,User-Agent",
		***REMOVED***,
		***REMOVED***
			key:   "Connection",
			value: []string***REMOVED***"close", "something-else"***REMOVED***,
			want:  "ERROR: http2: invalid Connection request header: [\"close\" \"something-else\"]",
		***REMOVED***,
		***REMOVED***
			key:   "Connection",
			value: []string***REMOVED***"keep-alive"***REMOVED***,
			want:  "Accept-Encoding,User-Agent",
		***REMOVED***,
		***REMOVED***
			key:   "Proxy-Connection", // just deleted and ignored
			value: []string***REMOVED***"keep-alive"***REMOVED***,
			want:  "Accept-Encoding,User-Agent",
		***REMOVED***,
		***REMOVED***
			key:   "Transfer-Encoding",
			value: []string***REMOVED***""***REMOVED***,
			want:  "Accept-Encoding,User-Agent",
		***REMOVED***,
		***REMOVED***
			key:   "Transfer-Encoding",
			value: []string***REMOVED***"foo"***REMOVED***,
			want:  "ERROR: http2: invalid Transfer-Encoding request header: [\"foo\"]",
		***REMOVED***,
		***REMOVED***
			key:   "Transfer-Encoding",
			value: []string***REMOVED***"chunked"***REMOVED***,
			want:  "Accept-Encoding,User-Agent",
		***REMOVED***,
		***REMOVED***
			key:   "Transfer-Encoding",
			value: []string***REMOVED***"chunked", "other"***REMOVED***,
			want:  "ERROR: http2: invalid Transfer-Encoding request header: [\"chunked\" \"other\"]",
		***REMOVED***,
		***REMOVED***
			key:   "Content-Length",
			value: []string***REMOVED***"123"***REMOVED***,
			want:  "Accept-Encoding,User-Agent",
		***REMOVED***,
		***REMOVED***
			key:   "Keep-Alive",
			value: []string***REMOVED***"doop"***REMOVED***,
			want:  "Accept-Encoding,User-Agent",
		***REMOVED***,
	***REMOVED***

	for _, tt := range tests ***REMOVED***
		req, _ := http.NewRequest("GET", st.ts.URL, nil)
		req.Header[tt.key] = tt.value
		res, err := tr.RoundTrip(req)
		var got string
		if err != nil ***REMOVED***
			got = fmt.Sprintf("ERROR: %v", err)
		***REMOVED*** else ***REMOVED***
			got = res.Header.Get("Got-Header")
			res.Body.Close()
		***REMOVED***
		if got != tt.want ***REMOVED***
			t.Errorf("For key %q, value %q, got = %q; want %q", tt.key, tt.value, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// golang.org/issue/14048
func TestTransportFailsOnInvalidHeaders(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		var got []string
		for k := range r.Header ***REMOVED***
			got = append(got, k)
		***REMOVED***
		sort.Strings(got)
		w.Header().Set("Got-Header", strings.Join(got, ","))
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tests := [...]struct ***REMOVED***
		h       http.Header
		wantErr string
	***REMOVED******REMOVED***
		0: ***REMOVED***
			h:       http.Header***REMOVED***"with space": ***REMOVED***"foo"***REMOVED******REMOVED***,
			wantErr: `invalid HTTP header name "with space"`,
		***REMOVED***,
		1: ***REMOVED***
			h:       http.Header***REMOVED***"name": ***REMOVED***"Брэд"***REMOVED******REMOVED***,
			wantErr: "", // okay
		***REMOVED***,
		2: ***REMOVED***
			h:       http.Header***REMOVED***"имя": ***REMOVED***"Brad"***REMOVED******REMOVED***,
			wantErr: `invalid HTTP header name "имя"`,
		***REMOVED***,
		3: ***REMOVED***
			h:       http.Header***REMOVED***"foo": ***REMOVED***"foo\x01bar"***REMOVED******REMOVED***,
			wantErr: `invalid HTTP header value "foo\x01bar" for header "foo"`,
		***REMOVED***,
	***REMOVED***

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	for i, tt := range tests ***REMOVED***
		req, _ := http.NewRequest("GET", st.ts.URL, nil)
		req.Header = tt.h
		res, err := tr.RoundTrip(req)
		var bad bool
		if tt.wantErr == "" ***REMOVED***
			if err != nil ***REMOVED***
				bad = true
				t.Errorf("case %d: error = %v; want no error", i, err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if !strings.Contains(fmt.Sprint(err), tt.wantErr) ***REMOVED***
				bad = true
				t.Errorf("case %d: error = %v; want error %q", i, err, tt.wantErr)
			***REMOVED***
		***REMOVED***
		if err == nil ***REMOVED***
			if bad ***REMOVED***
				t.Logf("case %d: server got headers %q", i, res.Header.Get("Got-Header"))
			***REMOVED***
			res.Body.Close()
		***REMOVED***
	***REMOVED***
***REMOVED***

// Tests that gzipReader doesn't crash on a second Read call following
// the first Read call's gzip.NewReader returning an error.
func TestGzipReader_DoubleReadCrash(t *testing.T) ***REMOVED***
	gz := &gzipReader***REMOVED***
		body: ioutil.NopCloser(strings.NewReader("0123456789")),
	***REMOVED***
	var buf [1]byte
	n, err1 := gz.Read(buf[:])
	if n != 0 || !strings.Contains(fmt.Sprint(err1), "invalid header") ***REMOVED***
		t.Fatalf("Read = %v, %v; want 0, invalid header", n, err1)
	***REMOVED***
	n, err2 := gz.Read(buf[:])
	if n != 0 || err2 != err1 ***REMOVED***
		t.Fatalf("second Read = %v, %v; want 0, %v", n, err2, err1)
	***REMOVED***
***REMOVED***

func TestTransportNewTLSConfig(t *testing.T) ***REMOVED***
	tests := [...]struct ***REMOVED***
		conf *tls.Config
		host string
		want *tls.Config
	***REMOVED******REMOVED***
		// Normal case.
		0: ***REMOVED***
			conf: nil,
			host: "foo.com",
			want: &tls.Config***REMOVED***
				ServerName: "foo.com",
				NextProtos: []string***REMOVED***NextProtoTLS***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		// User-provided name (bar.com) takes precedence:
		1: ***REMOVED***
			conf: &tls.Config***REMOVED***
				ServerName: "bar.com",
			***REMOVED***,
			host: "foo.com",
			want: &tls.Config***REMOVED***
				ServerName: "bar.com",
				NextProtos: []string***REMOVED***NextProtoTLS***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		// NextProto is prepended:
		2: ***REMOVED***
			conf: &tls.Config***REMOVED***
				NextProtos: []string***REMOVED***"foo", "bar"***REMOVED***,
			***REMOVED***,
			host: "example.com",
			want: &tls.Config***REMOVED***
				ServerName: "example.com",
				NextProtos: []string***REMOVED***NextProtoTLS, "foo", "bar"***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		// NextProto is not duplicated:
		3: ***REMOVED***
			conf: &tls.Config***REMOVED***
				NextProtos: []string***REMOVED***"foo", "bar", NextProtoTLS***REMOVED***,
			***REMOVED***,
			host: "example.com",
			want: &tls.Config***REMOVED***
				ServerName: "example.com",
				NextProtos: []string***REMOVED***"foo", "bar", NextProtoTLS***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		// Ignore the session ticket keys part, which ends up populating
		// unexported fields in the Config:
		if tt.conf != nil ***REMOVED***
			tt.conf.SessionTicketsDisabled = true
		***REMOVED***

		tr := &Transport***REMOVED***TLSClientConfig: tt.conf***REMOVED***
		got := tr.newTLSConfig(tt.host)

		got.SessionTicketsDisabled = false

		if !reflect.DeepEqual(got, tt.want) ***REMOVED***
			t.Errorf("%d. got %#v; want %#v", i, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// The Google GFE responds to HEAD requests with a HEADERS frame
// without END_STREAM, followed by a 0-length DATA frame with
// END_STREAM. Make sure we don't get confused by that. (We did.)
func TestTransportReadHeadResponse(t *testing.T) ***REMOVED***
	ct := newClientTester(t)
	clientDone := make(chan struct***REMOVED******REMOVED***)
	ct.client = func() error ***REMOVED***
		defer close(clientDone)
		req, _ := http.NewRequest("HEAD", "https://dummy.tld/", nil)
		res, err := ct.tr.RoundTrip(req)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if res.ContentLength != 123 ***REMOVED***
			return fmt.Errorf("Content-Length = %d; want 123", res.ContentLength)
		***REMOVED***
		slurp, err := ioutil.ReadAll(res.Body)
		if err != nil ***REMOVED***
			return fmt.Errorf("ReadAll: %v", err)
		***REMOVED***
		if len(slurp) > 0 ***REMOVED***
			return fmt.Errorf("Unexpected non-empty ReadAll body: %q", slurp)
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				t.Logf("ReadFrame: %v", err)
				return nil
			***REMOVED***
			hf, ok := f.(*HeadersFrame)
			if !ok ***REMOVED***
				continue
			***REMOVED***
			var buf bytes.Buffer
			enc := hpack.NewEncoder(&buf)
			enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
			enc.WriteField(hpack.HeaderField***REMOVED***Name: "content-length", Value: "123"***REMOVED***)
			ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
				StreamID:      hf.StreamID,
				EndHeaders:    true,
				EndStream:     false, // as the GFE does
				BlockFragment: buf.Bytes(),
			***REMOVED***)
			ct.fr.WriteData(hf.StreamID, true, nil)

			<-clientDone
			return nil
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

func TestTransportReadHeadResponseWithBody(t *testing.T) ***REMOVED***
	// This test use not valid response format.
	// Discarding logger output to not spam tests output.
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)

	response := "redirecting to /elsewhere"
	ct := newClientTester(t)
	clientDone := make(chan struct***REMOVED******REMOVED***)
	ct.client = func() error ***REMOVED***
		defer close(clientDone)
		req, _ := http.NewRequest("HEAD", "https://dummy.tld/", nil)
		res, err := ct.tr.RoundTrip(req)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if res.ContentLength != int64(len(response)) ***REMOVED***
			return fmt.Errorf("Content-Length = %d; want %d", res.ContentLength, len(response))
		***REMOVED***
		slurp, err := ioutil.ReadAll(res.Body)
		if err != nil ***REMOVED***
			return fmt.Errorf("ReadAll: %v", err)
		***REMOVED***
		if len(slurp) > 0 ***REMOVED***
			return fmt.Errorf("Unexpected non-empty ReadAll body: %q", slurp)
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				t.Logf("ReadFrame: %v", err)
				return nil
			***REMOVED***
			hf, ok := f.(*HeadersFrame)
			if !ok ***REMOVED***
				continue
			***REMOVED***
			var buf bytes.Buffer
			enc := hpack.NewEncoder(&buf)
			enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
			enc.WriteField(hpack.HeaderField***REMOVED***Name: "content-length", Value: strconv.Itoa(len(response))***REMOVED***)
			ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
				StreamID:      hf.StreamID,
				EndHeaders:    true,
				EndStream:     false,
				BlockFragment: buf.Bytes(),
			***REMOVED***)
			ct.fr.WriteData(hf.StreamID, true, []byte(response))

			<-clientDone
			return nil
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

type neverEnding byte

func (b neverEnding) Read(p []byte) (int, error) ***REMOVED***
	for i := range p ***REMOVED***
		p[i] = byte(b)
	***REMOVED***
	return len(p), nil
***REMOVED***

// golang.org/issue/15425: test that a handler closing the request
// body doesn't terminate the stream to the peer. (It just stops
// readability from the handler's side, and eventually the client
// runs out of flow control tokens)
func TestTransportHandlerBodyClose(t *testing.T) ***REMOVED***
	const bodySize = 10 << 20
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		r.Body.Close()
		io.Copy(w, io.LimitReader(neverEnding('A'), bodySize))
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	g0 := runtime.NumGoroutine()

	const numReq = 10
	for i := 0; i < numReq; i++ ***REMOVED***
		req, err := http.NewRequest("POST", st.ts.URL, struct***REMOVED*** io.Reader ***REMOVED******REMOVED***io.LimitReader(neverEnding('A'), bodySize)***REMOVED***)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		res, err := tr.RoundTrip(req)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		n, err := io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
		if n != bodySize || err != nil ***REMOVED***
			t.Fatalf("req#%d: Copy = %d, %v; want %d, nil", i, n, err, bodySize)
		***REMOVED***
	***REMOVED***
	tr.CloseIdleConnections()

	gd := runtime.NumGoroutine() - g0
	if gd > numReq/2 ***REMOVED***
		t.Errorf("appeared to leak goroutines")
	***REMOVED***

***REMOVED***

// https://golang.org/issue/15930
func TestTransportFlowControl(t *testing.T) ***REMOVED***
	const bufLen = 64 << 10
	var total int64 = 100 << 20 // 100MB
	if testing.Short() ***REMOVED***
		total = 10 << 20
	***REMOVED***

	var wrote int64 // updated atomically
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		b := make([]byte, bufLen)
		for wrote < total ***REMOVED***
			n, err := w.Write(b)
			atomic.AddInt64(&wrote, int64(n))
			if err != nil ***REMOVED***
				t.Errorf("ResponseWriter.Write error: %v", err)
				break
			***REMOVED***
			w.(http.Flusher).Flush()
		***REMOVED***
	***REMOVED***, optOnlyServer)

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()
	req, err := http.NewRequest("GET", st.ts.URL, nil)
	if err != nil ***REMOVED***
		t.Fatal("NewRequest error:", err)
	***REMOVED***
	resp, err := tr.RoundTrip(req)
	if err != nil ***REMOVED***
		t.Fatal("RoundTrip error:", err)
	***REMOVED***
	defer resp.Body.Close()

	var read int64
	b := make([]byte, bufLen)
	for ***REMOVED***
		n, err := resp.Body.Read(b)
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			t.Fatal("Read error:", err)
		***REMOVED***
		read += int64(n)

		const max = transportDefaultStreamFlow
		if w := atomic.LoadInt64(&wrote); -max > read-w || read-w > max ***REMOVED***
			t.Fatalf("Too much data inflight: server wrote %v bytes but client only received %v", w, read)
		***REMOVED***

		// Let the server get ahead of the client.
		time.Sleep(1 * time.Millisecond)
	***REMOVED***
***REMOVED***

// golang.org/issue/14627 -- if the server sends a GOAWAY frame, make
// the Transport remember it and return it back to users (via
// RoundTrip or request body reads) if needed (e.g. if the server
// proceeds to close the TCP connection before the client gets its
// response)
func TestTransportUsesGoAwayDebugError_RoundTrip(t *testing.T) ***REMOVED***
	testTransportUsesGoAwayDebugError(t, false)
***REMOVED***

func TestTransportUsesGoAwayDebugError_Body(t *testing.T) ***REMOVED***
	testTransportUsesGoAwayDebugError(t, true)
***REMOVED***

func testTransportUsesGoAwayDebugError(t *testing.T, failMidBody bool) ***REMOVED***
	ct := newClientTester(t)
	clientDone := make(chan struct***REMOVED******REMOVED***)

	const goAwayErrCode = ErrCodeHTTP11Required // arbitrary
	const goAwayDebugData = "some debug data"

	ct.client = func() error ***REMOVED***
		defer close(clientDone)
		req, _ := http.NewRequest("GET", "https://dummy.tld/", nil)
		res, err := ct.tr.RoundTrip(req)
		if failMidBody ***REMOVED***
			if err != nil ***REMOVED***
				return fmt.Errorf("unexpected client RoundTrip error: %v", err)
			***REMOVED***
			_, err = io.Copy(ioutil.Discard, res.Body)
			res.Body.Close()
		***REMOVED***
		want := GoAwayError***REMOVED***
			LastStreamID: 5,
			ErrCode:      goAwayErrCode,
			DebugData:    goAwayDebugData,
		***REMOVED***
		if !reflect.DeepEqual(err, want) ***REMOVED***
			t.Errorf("RoundTrip error = %T: %#v, want %T (%#v)", err, err, want, want)
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				t.Logf("ReadFrame: %v", err)
				return nil
			***REMOVED***
			hf, ok := f.(*HeadersFrame)
			if !ok ***REMOVED***
				continue
			***REMOVED***
			if failMidBody ***REMOVED***
				var buf bytes.Buffer
				enc := hpack.NewEncoder(&buf)
				enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
				enc.WriteField(hpack.HeaderField***REMOVED***Name: "content-length", Value: "123"***REMOVED***)
				ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
					StreamID:      hf.StreamID,
					EndHeaders:    true,
					EndStream:     false,
					BlockFragment: buf.Bytes(),
				***REMOVED***)
			***REMOVED***
			// Write two GOAWAY frames, to test that the Transport takes
			// the interesting parts of both.
			ct.fr.WriteGoAway(5, ErrCodeNo, []byte(goAwayDebugData))
			ct.fr.WriteGoAway(5, goAwayErrCode, nil)
			ct.sc.(*net.TCPConn).CloseWrite()
			<-clientDone
			return nil
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

func testTransportReturnsUnusedFlowControl(t *testing.T, oneDataFrame bool) ***REMOVED***
	ct := newClientTester(t)

	clientClosed := make(chan struct***REMOVED******REMOVED***)
	serverWroteFirstByte := make(chan struct***REMOVED******REMOVED***)

	ct.client = func() error ***REMOVED***
		req, _ := http.NewRequest("GET", "https://dummy.tld/", nil)
		res, err := ct.tr.RoundTrip(req)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		<-serverWroteFirstByte

		if n, err := res.Body.Read(make([]byte, 1)); err != nil || n != 1 ***REMOVED***
			return fmt.Errorf("body read = %v, %v; want 1, nil", n, err)
		***REMOVED***
		res.Body.Close() // leaving 4999 bytes unread
		close(clientClosed)

		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()

		var hf *HeadersFrame
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				return fmt.Errorf("ReadFrame while waiting for Headers: %v", err)
			***REMOVED***
			switch f.(type) ***REMOVED***
			case *WindowUpdateFrame, *SettingsFrame:
				continue
			***REMOVED***
			var ok bool
			hf, ok = f.(*HeadersFrame)
			if !ok ***REMOVED***
				return fmt.Errorf("Got %T; want HeadersFrame", f)
			***REMOVED***
			break
		***REMOVED***

		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: "content-length", Value: "5000"***REMOVED***)
		ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
			StreamID:      hf.StreamID,
			EndHeaders:    true,
			EndStream:     false,
			BlockFragment: buf.Bytes(),
		***REMOVED***)

		// Two cases:
		// - Send one DATA frame with 5000 bytes.
		// - Send two DATA frames with 1 and 4999 bytes each.
		//
		// In both cases, the client should consume one byte of data,
		// refund that byte, then refund the following 4999 bytes.
		//
		// In the second case, the server waits for the client connection to
		// close before seconding the second DATA frame. This tests the case
		// where the client receives a DATA frame after it has reset the stream.
		if oneDataFrame ***REMOVED***
			ct.fr.WriteData(hf.StreamID, false /* don't end stream */, make([]byte, 5000))
			close(serverWroteFirstByte)
			<-clientClosed
		***REMOVED*** else ***REMOVED***
			ct.fr.WriteData(hf.StreamID, false /* don't end stream */, make([]byte, 1))
			close(serverWroteFirstByte)
			<-clientClosed
			ct.fr.WriteData(hf.StreamID, false /* don't end stream */, make([]byte, 4999))
		***REMOVED***

		waitingFor := "RSTStreamFrame"
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				return fmt.Errorf("ReadFrame while waiting for %s: %v", waitingFor, err)
			***REMOVED***
			if _, ok := f.(*SettingsFrame); ok ***REMOVED***
				continue
			***REMOVED***
			switch waitingFor ***REMOVED***
			case "RSTStreamFrame":
				if rf, ok := f.(*RSTStreamFrame); !ok || rf.ErrCode != ErrCodeCancel ***REMOVED***
					return fmt.Errorf("Expected a RSTStreamFrame with code cancel; got %v", summarizeFrame(f))
				***REMOVED***
				waitingFor = "WindowUpdateFrame"
			case "WindowUpdateFrame":
				if wuf, ok := f.(*WindowUpdateFrame); !ok || wuf.Increment != 4999 ***REMOVED***
					return fmt.Errorf("Expected WindowUpdateFrame for 4999 bytes; got %v", summarizeFrame(f))
				***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

// See golang.org/issue/16481
func TestTransportReturnsUnusedFlowControlSingleWrite(t *testing.T) ***REMOVED***
	testTransportReturnsUnusedFlowControl(t, true)
***REMOVED***

// See golang.org/issue/20469
func TestTransportReturnsUnusedFlowControlMultipleWrites(t *testing.T) ***REMOVED***
	testTransportReturnsUnusedFlowControl(t, false)
***REMOVED***

// Issue 16612: adjust flow control on open streams when transport
// receives SETTINGS with INITIAL_WINDOW_SIZE from server.
func TestTransportAdjustsFlowControl(t *testing.T) ***REMOVED***
	ct := newClientTester(t)
	clientDone := make(chan struct***REMOVED******REMOVED***)

	const bodySize = 1 << 20

	ct.client = func() error ***REMOVED***
		defer ct.cc.(*net.TCPConn).CloseWrite()
		defer close(clientDone)

		req, _ := http.NewRequest("POST", "https://dummy.tld/", struct***REMOVED*** io.Reader ***REMOVED******REMOVED***io.LimitReader(neverEnding('A'), bodySize)***REMOVED***)
		res, err := ct.tr.RoundTrip(req)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		res.Body.Close()
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		_, err := io.ReadFull(ct.sc, make([]byte, len(ClientPreface)))
		if err != nil ***REMOVED***
			return fmt.Errorf("reading client preface: %v", err)
		***REMOVED***

		var gotBytes int64
		var sentSettings bool
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				select ***REMOVED***
				case <-clientDone:
					return nil
				default:
					return fmt.Errorf("ReadFrame while waiting for Headers: %v", err)
				***REMOVED***
			***REMOVED***
			switch f := f.(type) ***REMOVED***
			case *DataFrame:
				gotBytes += int64(len(f.Data()))
				// After we've got half the client's
				// initial flow control window's worth
				// of request body data, give it just
				// enough flow control to finish.
				if gotBytes >= initialWindowSize/2 && !sentSettings ***REMOVED***
					sentSettings = true

					ct.fr.WriteSettings(Setting***REMOVED***ID: SettingInitialWindowSize, Val: bodySize***REMOVED***)
					ct.fr.WriteWindowUpdate(0, bodySize)
					ct.fr.WriteSettingsAck()
				***REMOVED***

				if f.StreamEnded() ***REMOVED***
					var buf bytes.Buffer
					enc := hpack.NewEncoder(&buf)
					enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
					ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
						StreamID:      f.StreamID,
						EndHeaders:    true,
						EndStream:     true,
						BlockFragment: buf.Bytes(),
					***REMOVED***)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

// See golang.org/issue/16556
func TestTransportReturnsDataPaddingFlowControl(t *testing.T) ***REMOVED***
	ct := newClientTester(t)

	unblockClient := make(chan bool, 1)

	ct.client = func() error ***REMOVED***
		req, _ := http.NewRequest("GET", "https://dummy.tld/", nil)
		res, err := ct.tr.RoundTrip(req)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer res.Body.Close()
		<-unblockClient
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()

		var hf *HeadersFrame
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				return fmt.Errorf("ReadFrame while waiting for Headers: %v", err)
			***REMOVED***
			switch f.(type) ***REMOVED***
			case *WindowUpdateFrame, *SettingsFrame:
				continue
			***REMOVED***
			var ok bool
			hf, ok = f.(*HeadersFrame)
			if !ok ***REMOVED***
				return fmt.Errorf("Got %T; want HeadersFrame", f)
			***REMOVED***
			break
		***REMOVED***

		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: "content-length", Value: "5000"***REMOVED***)
		ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
			StreamID:      hf.StreamID,
			EndHeaders:    true,
			EndStream:     false,
			BlockFragment: buf.Bytes(),
		***REMOVED***)
		pad := make([]byte, 5)
		ct.fr.WriteDataPadded(hf.StreamID, false, make([]byte, 5000), pad) // without ending stream

		f, err := ct.readNonSettingsFrame()
		if err != nil ***REMOVED***
			return fmt.Errorf("ReadFrame while waiting for first WindowUpdateFrame: %v", err)
		***REMOVED***
		wantBack := uint32(len(pad)) + 1 // one byte for the length of the padding
		if wuf, ok := f.(*WindowUpdateFrame); !ok || wuf.Increment != wantBack || wuf.StreamID != 0 ***REMOVED***
			return fmt.Errorf("Expected conn WindowUpdateFrame for %d bytes; got %v", wantBack, summarizeFrame(f))
		***REMOVED***

		f, err = ct.readNonSettingsFrame()
		if err != nil ***REMOVED***
			return fmt.Errorf("ReadFrame while waiting for second WindowUpdateFrame: %v", err)
		***REMOVED***
		if wuf, ok := f.(*WindowUpdateFrame); !ok || wuf.Increment != wantBack || wuf.StreamID == 0 ***REMOVED***
			return fmt.Errorf("Expected stream WindowUpdateFrame for %d bytes; got %v", wantBack, summarizeFrame(f))
		***REMOVED***
		unblockClient <- true
		return nil
	***REMOVED***
	ct.run()
***REMOVED***

// golang.org/issue/16572 -- RoundTrip shouldn't hang when it gets a
// StreamError as a result of the response HEADERS
func TestTransportReturnsErrorOnBadResponseHeaders(t *testing.T) ***REMOVED***
	ct := newClientTester(t)

	ct.client = func() error ***REMOVED***
		req, _ := http.NewRequest("GET", "https://dummy.tld/", nil)
		res, err := ct.tr.RoundTrip(req)
		if err == nil ***REMOVED***
			res.Body.Close()
			return errors.New("unexpected successful GET")
		***REMOVED***
		want := StreamError***REMOVED***1, ErrCodeProtocol, headerFieldNameError("  content-type")***REMOVED***
		if !reflect.DeepEqual(want, err) ***REMOVED***
			t.Errorf("RoundTrip error = %#v; want %#v", err, want)
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()

		hf, err := ct.firstHeaders()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: "  content-type", Value: "bogus"***REMOVED***) // bogus spaces
		ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
			StreamID:      hf.StreamID,
			EndHeaders:    true,
			EndStream:     false,
			BlockFragment: buf.Bytes(),
		***REMOVED***)

		for ***REMOVED***
			fr, err := ct.readFrame()
			if err != nil ***REMOVED***
				return fmt.Errorf("error waiting for RST_STREAM from client: %v", err)
			***REMOVED***
			if _, ok := fr.(*SettingsFrame); ok ***REMOVED***
				continue
			***REMOVED***
			if rst, ok := fr.(*RSTStreamFrame); !ok || rst.StreamID != 1 || rst.ErrCode != ErrCodeProtocol ***REMOVED***
				t.Errorf("Frame = %v; want RST_STREAM for stream 1 with ErrCodeProtocol", summarizeFrame(fr))
			***REMOVED***
			break
		***REMOVED***

		return nil
	***REMOVED***
	ct.run()
***REMOVED***

// byteAndEOFReader returns is in an io.Reader which reads one byte
// (the underlying byte) and io.EOF at once in its Read call.
type byteAndEOFReader byte

func (b byteAndEOFReader) Read(p []byte) (n int, err error) ***REMOVED***
	if len(p) == 0 ***REMOVED***
		panic("unexpected useless call")
	***REMOVED***
	p[0] = byte(b)
	return 1, io.EOF
***REMOVED***

// Issue 16788: the Transport had a regression where it started
// sending a spurious DATA frame with a duplicate END_STREAM bit after
// the request body writer goroutine had already read an EOF from the
// Request.Body and included the END_STREAM on a data-carrying DATA
// frame.
//
// Notably, to trigger this, the requests need to use a Request.Body
// which returns (non-0, io.EOF) and also needs to set the ContentLength
// explicitly.
func TestTransportBodyDoubleEndStream(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		// Nothing.
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	for i := 0; i < 2; i++ ***REMOVED***
		req, _ := http.NewRequest("POST", st.ts.URL, byteAndEOFReader('a'))
		req.ContentLength = 1
		res, err := tr.RoundTrip(req)
		if err != nil ***REMOVED***
			t.Fatalf("failure on req %d: %v", i+1, err)
		***REMOVED***
		defer res.Body.Close()
	***REMOVED***
***REMOVED***

// golang.org/issue/16847, golang.org/issue/19103
func TestTransportRequestPathPseudo(t *testing.T) ***REMOVED***
	type result struct ***REMOVED***
		path string
		err  string
	***REMOVED***
	tests := []struct ***REMOVED***
		req  *http.Request
		want result
	***REMOVED******REMOVED***
		0: ***REMOVED***
			req: &http.Request***REMOVED***
				Method: "GET",
				URL: &url.URL***REMOVED***
					Host: "foo.com",
					Path: "/foo",
				***REMOVED***,
			***REMOVED***,
			want: result***REMOVED***path: "/foo"***REMOVED***,
		***REMOVED***,
		// In Go 1.7, we accepted paths of "//foo".
		// In Go 1.8, we rejected it (issue 16847).
		// In Go 1.9, we accepted it again (issue 19103).
		1: ***REMOVED***
			req: &http.Request***REMOVED***
				Method: "GET",
				URL: &url.URL***REMOVED***
					Host: "foo.com",
					Path: "//foo",
				***REMOVED***,
			***REMOVED***,
			want: result***REMOVED***path: "//foo"***REMOVED***,
		***REMOVED***,

		// Opaque with //$Matching_Hostname/path
		2: ***REMOVED***
			req: &http.Request***REMOVED***
				Method: "GET",
				URL: &url.URL***REMOVED***
					Scheme: "https",
					Opaque: "//foo.com/path",
					Host:   "foo.com",
					Path:   "/ignored",
				***REMOVED***,
			***REMOVED***,
			want: result***REMOVED***path: "/path"***REMOVED***,
		***REMOVED***,

		// Opaque with some other Request.Host instead:
		3: ***REMOVED***
			req: &http.Request***REMOVED***
				Method: "GET",
				Host:   "bar.com",
				URL: &url.URL***REMOVED***
					Scheme: "https",
					Opaque: "//bar.com/path",
					Host:   "foo.com",
					Path:   "/ignored",
				***REMOVED***,
			***REMOVED***,
			want: result***REMOVED***path: "/path"***REMOVED***,
		***REMOVED***,

		// Opaque without the leading "//":
		4: ***REMOVED***
			req: &http.Request***REMOVED***
				Method: "GET",
				URL: &url.URL***REMOVED***
					Opaque: "/path",
					Host:   "foo.com",
					Path:   "/ignored",
				***REMOVED***,
			***REMOVED***,
			want: result***REMOVED***path: "/path"***REMOVED***,
		***REMOVED***,

		// Opaque we can't handle:
		5: ***REMOVED***
			req: &http.Request***REMOVED***
				Method: "GET",
				URL: &url.URL***REMOVED***
					Scheme: "https",
					Opaque: "//unknown_host/path",
					Host:   "foo.com",
					Path:   "/ignored",
				***REMOVED***,
			***REMOVED***,
			want: result***REMOVED***err: `invalid request :path "https://unknown_host/path" from URL.Opaque = "//unknown_host/path"`***REMOVED***,
		***REMOVED***,

		// A CONNECT request:
		6: ***REMOVED***
			req: &http.Request***REMOVED***
				Method: "CONNECT",
				URL: &url.URL***REMOVED***
					Host: "foo.com",
				***REMOVED***,
			***REMOVED***,
			want: result***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		cc := &ClientConn***REMOVED***peerMaxHeaderListSize: 0xffffffffffffffff***REMOVED***
		cc.henc = hpack.NewEncoder(&cc.hbuf)
		cc.mu.Lock()
		hdrs, err := cc.encodeHeaders(tt.req, false, "", -1)
		cc.mu.Unlock()
		var got result
		hpackDec := hpack.NewDecoder(initialHeaderTableSize, func(f hpack.HeaderField) ***REMOVED***
			if f.Name == ":path" ***REMOVED***
				got.path = f.Value
			***REMOVED***
		***REMOVED***)
		if err != nil ***REMOVED***
			got.err = err.Error()
		***REMOVED*** else if len(hdrs) > 0 ***REMOVED***
			if _, err := hpackDec.Write(hdrs); err != nil ***REMOVED***
				t.Errorf("%d. bogus hpack: %v", i, err)
				continue
			***REMOVED***
		***REMOVED***
		if got != tt.want ***REMOVED***
			t.Errorf("%d. got %+v; want %+v", i, got, tt.want)
		***REMOVED***

	***REMOVED***

***REMOVED***

// golang.org/issue/17071 -- don't sniff the first byte of the request body
// before we've determined that the ClientConn is usable.
func TestRoundTripDoesntConsumeRequestBodyEarly(t *testing.T) ***REMOVED***
	const body = "foo"
	req, _ := http.NewRequest("POST", "http://foo.com/", ioutil.NopCloser(strings.NewReader(body)))
	cc := &ClientConn***REMOVED***
		closed: true,
	***REMOVED***
	_, err := cc.RoundTrip(req)
	if err != errClientConnUnusable ***REMOVED***
		t.Fatalf("RoundTrip = %v; want errClientConnUnusable", err)
	***REMOVED***
	slurp, err := ioutil.ReadAll(req.Body)
	if err != nil ***REMOVED***
		t.Errorf("ReadAll = %v", err)
	***REMOVED***
	if string(slurp) != body ***REMOVED***
		t.Errorf("Body = %q; want %q", slurp, body)
	***REMOVED***
***REMOVED***

func TestClientConnPing(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED******REMOVED***, optOnlyServer)
	defer st.Close()
	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()
	cc, err := tr.dialClientConn(st.ts.Listener.Addr().String(), false)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err = cc.Ping(testContext***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Issue 16974: if the server sent a DATA frame after the user
// canceled the Transport's Request, the Transport previously wrote to a
// closed pipe, got an error, and ended up closing the whole TCP
// connection.
func TestTransportCancelDataResponseRace(t *testing.T) ***REMOVED***
	cancel := make(chan struct***REMOVED******REMOVED***)
	clientGotError := make(chan bool, 1)

	const msg = "Hello."
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if strings.Contains(r.URL.Path, "/hello") ***REMOVED***
			time.Sleep(50 * time.Millisecond)
			io.WriteString(w, msg)
			return
		***REMOVED***
		for i := 0; i < 50; i++ ***REMOVED***
			io.WriteString(w, "Some data.")
			w.(http.Flusher).Flush()
			if i == 2 ***REMOVED***
				close(cancel)
				<-clientGotError
			***REMOVED***
			time.Sleep(10 * time.Millisecond)
		***REMOVED***
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	c := &http.Client***REMOVED***Transport: tr***REMOVED***
	req, _ := http.NewRequest("GET", st.ts.URL, nil)
	req.Cancel = cancel
	res, err := c.Do(req)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err = io.Copy(ioutil.Discard, res.Body); err == nil ***REMOVED***
		t.Fatal("unexpected success")
	***REMOVED***
	clientGotError <- true

	res, err = c.Get(st.ts.URL + "/hello")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	slurp, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(slurp) != msg ***REMOVED***
		t.Errorf("Got = %q; want %q", slurp, msg)
	***REMOVED***
***REMOVED***

// Issue 21316: It should be safe to reuse an http.Request after the
// request has completed.
func TestTransportNoRaceOnRequestObjectAfterRequestComplete(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.WriteHeader(200)
		io.WriteString(w, "body")
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	req, _ := http.NewRequest("GET", st.ts.URL, nil)
	resp, err := tr.RoundTrip(req)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err = io.Copy(ioutil.Discard, resp.Body); err != nil ***REMOVED***
		t.Fatalf("error reading response body: %v", err)
	***REMOVED***
	if err := resp.Body.Close(); err != nil ***REMOVED***
		t.Fatalf("error closing response body: %v", err)
	***REMOVED***

	// This access of req.Header should not race with code in the transport.
	req.Header = http.Header***REMOVED******REMOVED***
***REMOVED***

func TestTransportRetryAfterGOAWAY(t *testing.T) ***REMOVED***
	var dialer struct ***REMOVED***
		sync.Mutex
		count int
	***REMOVED***
	ct1 := make(chan *clientTester)
	ct2 := make(chan *clientTester)

	ln := newLocalListener(t)
	defer ln.Close()

	tr := &Transport***REMOVED***
		TLSClientConfig: tlsConfigInsecure,
	***REMOVED***
	tr.DialTLS = func(network, addr string, cfg *tls.Config) (net.Conn, error) ***REMOVED***
		dialer.Lock()
		defer dialer.Unlock()
		dialer.count++
		if dialer.count == 3 ***REMOVED***
			return nil, errors.New("unexpected number of dials")
		***REMOVED***
		cc, err := net.Dial("tcp", ln.Addr().String())
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("dial error: %v", err)
		***REMOVED***
		sc, err := ln.Accept()
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("accept error: %v", err)
		***REMOVED***
		ct := &clientTester***REMOVED***
			t:  t,
			tr: tr,
			cc: cc,
			sc: sc,
			fr: NewFramer(sc, sc),
		***REMOVED***
		switch dialer.count ***REMOVED***
		case 1:
			ct1 <- ct
		case 2:
			ct2 <- ct
		***REMOVED***
		return cc, nil
	***REMOVED***

	errs := make(chan error, 3)
	done := make(chan struct***REMOVED******REMOVED***)
	defer close(done)

	// Client.
	go func() ***REMOVED***
		req, _ := http.NewRequest("GET", "https://dummy.tld/", nil)
		res, err := tr.RoundTrip(req)
		if res != nil ***REMOVED***
			res.Body.Close()
			if got := res.Header.Get("Foo"); got != "bar" ***REMOVED***
				err = fmt.Errorf("foo header = %q; want bar", got)
			***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			err = fmt.Errorf("RoundTrip: %v", err)
		***REMOVED***
		errs <- err
	***REMOVED***()

	connToClose := make(chan io.Closer, 2)

	// Server for the first request.
	go func() ***REMOVED***
		var ct *clientTester
		select ***REMOVED***
		case ct = <-ct1:
		case <-done:
			return
		***REMOVED***

		connToClose <- ct.cc
		ct.greet()
		hf, err := ct.firstHeaders()
		if err != nil ***REMOVED***
			errs <- fmt.Errorf("server1 failed reading HEADERS: %v", err)
			return
		***REMOVED***
		t.Logf("server1 got %v", hf)
		if err := ct.fr.WriteGoAway(0 /*max id*/, ErrCodeNo, nil); err != nil ***REMOVED***
			errs <- fmt.Errorf("server1 failed writing GOAWAY: %v", err)
			return
		***REMOVED***
		errs <- nil
	***REMOVED***()

	// Server for the second request.
	go func() ***REMOVED***
		var ct *clientTester
		select ***REMOVED***
		case ct = <-ct2:
		case <-done:
			return
		***REMOVED***

		connToClose <- ct.cc
		ct.greet()
		hf, err := ct.firstHeaders()
		if err != nil ***REMOVED***
			errs <- fmt.Errorf("server2 failed reading HEADERS: %v", err)
			return
		***REMOVED***
		t.Logf("server2 got %v", hf)

		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: "foo", Value: "bar"***REMOVED***)
		err = ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
			StreamID:      hf.StreamID,
			EndHeaders:    true,
			EndStream:     false,
			BlockFragment: buf.Bytes(),
		***REMOVED***)
		if err != nil ***REMOVED***
			errs <- fmt.Errorf("server2 failed writing response HEADERS: %v", err)
		***REMOVED*** else ***REMOVED***
			errs <- nil
		***REMOVED***
	***REMOVED***()

	for k := 0; k < 3; k++ ***REMOVED***
		select ***REMOVED***
		case err := <-errs:
			if err != nil ***REMOVED***
				t.Error(err)
			***REMOVED***
		case <-time.After(1 * time.Second):
			t.Errorf("timed out")
		***REMOVED***
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case c := <-connToClose:
			c.Close()
		default:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTransportRetryAfterRefusedStream(t *testing.T) ***REMOVED***
	clientDone := make(chan struct***REMOVED******REMOVED***)
	ct := newClientTester(t)
	ct.client = func() error ***REMOVED***
		defer ct.cc.(*net.TCPConn).CloseWrite()
		defer close(clientDone)
		req, _ := http.NewRequest("GET", "https://dummy.tld/", nil)
		resp, err := ct.tr.RoundTrip(req)
		if err != nil ***REMOVED***
			return fmt.Errorf("RoundTrip: %v", err)
		***REMOVED***
		resp.Body.Close()
		if resp.StatusCode != 204 ***REMOVED***
			return fmt.Errorf("Status = %v; want 204", resp.StatusCode)
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)
		nreq := 0

		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				select ***REMOVED***
				case <-clientDone:
					// If the client's done, it
					// will have reported any
					// errors on its side.
					return nil
				default:
					return err
				***REMOVED***
			***REMOVED***
			switch f := f.(type) ***REMOVED***
			case *WindowUpdateFrame, *SettingsFrame:
			case *HeadersFrame:
				if !f.HeadersEnded() ***REMOVED***
					return fmt.Errorf("headers should have END_HEADERS be ended: %v", f)
				***REMOVED***
				nreq++
				if nreq == 1 ***REMOVED***
					ct.fr.WriteRSTStream(f.StreamID, ErrCodeRefusedStream)
				***REMOVED*** else ***REMOVED***
					enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "204"***REMOVED***)
					ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
						StreamID:      f.StreamID,
						EndHeaders:    true,
						EndStream:     true,
						BlockFragment: buf.Bytes(),
					***REMOVED***)
				***REMOVED***
			default:
				return fmt.Errorf("Unexpected client frame %v", f)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

func TestTransportRetryHasLimit(t *testing.T) ***REMOVED***
	// Skip in short mode because the total expected delay is 1s+2s+4s+8s+16s=29s.
	if testing.Short() ***REMOVED***
		t.Skip("skipping long test in short mode")
	***REMOVED***
	clientDone := make(chan struct***REMOVED******REMOVED***)
	ct := newClientTester(t)
	ct.client = func() error ***REMOVED***
		defer ct.cc.(*net.TCPConn).CloseWrite()
		defer close(clientDone)
		req, _ := http.NewRequest("GET", "https://dummy.tld/", nil)
		resp, err := ct.tr.RoundTrip(req)
		if err == nil ***REMOVED***
			return fmt.Errorf("RoundTrip expected error, got response: %+v", resp)
		***REMOVED***
		t.Logf("expected error, got: %v", err)
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				select ***REMOVED***
				case <-clientDone:
					// If the client's done, it
					// will have reported any
					// errors on its side.
					return nil
				default:
					return err
				***REMOVED***
			***REMOVED***
			switch f := f.(type) ***REMOVED***
			case *WindowUpdateFrame, *SettingsFrame:
			case *HeadersFrame:
				if !f.HeadersEnded() ***REMOVED***
					return fmt.Errorf("headers should have END_HEADERS be ended: %v", f)
				***REMOVED***
				ct.fr.WriteRSTStream(f.StreamID, ErrCodeRefusedStream)
			default:
				return fmt.Errorf("Unexpected client frame %v", f)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

func TestTransportResponseDataBeforeHeaders(t *testing.T) ***REMOVED***
	// This test use not valid response format.
	// Discarding logger output to not spam tests output.
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)

	ct := newClientTester(t)
	ct.client = func() error ***REMOVED***
		defer ct.cc.(*net.TCPConn).CloseWrite()
		req := httptest.NewRequest("GET", "https://dummy.tld/", nil)
		// First request is normal to ensure the check is per stream and not per connection.
		_, err := ct.tr.RoundTrip(req)
		if err != nil ***REMOVED***
			return fmt.Errorf("RoundTrip expected no error, got: %v", err)
		***REMOVED***
		// Second request returns a DATA frame with no HEADERS.
		resp, err := ct.tr.RoundTrip(req)
		if err == nil ***REMOVED***
			return fmt.Errorf("RoundTrip expected error, got response: %+v", resp)
		***REMOVED***
		if err, ok := err.(StreamError); !ok || err.Code != ErrCodeProtocol ***REMOVED***
			return fmt.Errorf("expected stream PROTOCOL_ERROR, got: %v", err)
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err == io.EOF ***REMOVED***
				return nil
			***REMOVED*** else if err != nil ***REMOVED***
				return err
			***REMOVED***
			switch f := f.(type) ***REMOVED***
			case *WindowUpdateFrame, *SettingsFrame:
			case *HeadersFrame:
				switch f.StreamID ***REMOVED***
				case 1:
					// Send a valid response to first request.
					var buf bytes.Buffer
					enc := hpack.NewEncoder(&buf)
					enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
					ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
						StreamID:      f.StreamID,
						EndHeaders:    true,
						EndStream:     true,
						BlockFragment: buf.Bytes(),
					***REMOVED***)
				case 3:
					ct.fr.WriteData(f.StreamID, true, []byte("payload"))
				***REMOVED***
			default:
				return fmt.Errorf("Unexpected client frame %v", f)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***
func TestTransportRequestsStallAtServerLimit(t *testing.T) ***REMOVED***
	const maxConcurrent = 2

	greet := make(chan struct***REMOVED******REMOVED***)      // server sends initial SETTINGS frame
	gotRequest := make(chan struct***REMOVED******REMOVED***) // server received a request
	clientDone := make(chan struct***REMOVED******REMOVED***)

	// Collect errors from goroutines.
	var wg sync.WaitGroup
	errs := make(chan error, 100)
	defer func() ***REMOVED***
		wg.Wait()
		close(errs)
		for err := range errs ***REMOVED***
			t.Error(err)
		***REMOVED***
	***REMOVED***()

	// We will send maxConcurrent+2 requests. This checker goroutine waits for the
	// following stages:
	//   1. The first maxConcurrent requests are received by the server.
	//   2. The client will cancel the next request
	//   3. The server is unblocked so it can service the first maxConcurrent requests
	//   4. The client will send the final request
	wg.Add(1)
	unblockClient := make(chan struct***REMOVED******REMOVED***)
	clientRequestCancelled := make(chan struct***REMOVED******REMOVED***)
	unblockServer := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		defer wg.Done()
		// Stage 1.
		for k := 0; k < maxConcurrent; k++ ***REMOVED***
			<-gotRequest
		***REMOVED***
		// Stage 2.
		close(unblockClient)
		<-clientRequestCancelled
		// Stage 3: give some time for the final RoundTrip call to be scheduled and
		// verify that the final request is not sent.
		time.Sleep(50 * time.Millisecond)
		select ***REMOVED***
		case <-gotRequest:
			errs <- errors.New("last request did not stall")
			close(unblockServer)
			return
		default:
		***REMOVED***
		close(unblockServer)
		// Stage 4.
		<-gotRequest
	***REMOVED***()

	ct := newClientTester(t)
	ct.client = func() error ***REMOVED***
		var wg sync.WaitGroup
		defer func() ***REMOVED***
			wg.Wait()
			close(clientDone)
			ct.cc.(*net.TCPConn).CloseWrite()
		***REMOVED***()
		for k := 0; k < maxConcurrent+2; k++ ***REMOVED***
			wg.Add(1)
			go func(k int) ***REMOVED***
				defer wg.Done()
				// Don't send the second request until after receiving SETTINGS from the server
				// to avoid a race where we use the default SettingMaxConcurrentStreams, which
				// is much larger than maxConcurrent. We have to send the first request before
				// waiting because the first request triggers the dial and greet.
				if k > 0 ***REMOVED***
					<-greet
				***REMOVED***
				// Block until maxConcurrent requests are sent before sending any more.
				if k >= maxConcurrent ***REMOVED***
					<-unblockClient
				***REMOVED***
				req, _ := http.NewRequest("GET", fmt.Sprintf("https://dummy.tld/%d", k), nil)
				if k == maxConcurrent ***REMOVED***
					// This request will be canceled.
					cancel := make(chan struct***REMOVED******REMOVED***)
					req.Cancel = cancel
					close(cancel)
					_, err := ct.tr.RoundTrip(req)
					close(clientRequestCancelled)
					if err == nil ***REMOVED***
						errs <- fmt.Errorf("RoundTrip(%d) should have failed due to cancel", k)
						return
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					resp, err := ct.tr.RoundTrip(req)
					if err != nil ***REMOVED***
						errs <- fmt.Errorf("RoundTrip(%d): %v", k, err)
						return
					***REMOVED***
					ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					if resp.StatusCode != 204 ***REMOVED***
						errs <- fmt.Errorf("Status = %v; want 204", resp.StatusCode)
						return
					***REMOVED***
				***REMOVED***
			***REMOVED***(k)
		***REMOVED***
		return nil
	***REMOVED***

	ct.server = func() error ***REMOVED***
		var wg sync.WaitGroup
		defer wg.Wait()

		ct.greet(Setting***REMOVED***SettingMaxConcurrentStreams, maxConcurrent***REMOVED***)

		// Server write loop.
		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)
		writeResp := make(chan uint32, maxConcurrent+1)

		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			<-unblockServer
			for id := range writeResp ***REMOVED***
				buf.Reset()
				enc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "204"***REMOVED***)
				ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
					StreamID:      id,
					EndHeaders:    true,
					EndStream:     true,
					BlockFragment: buf.Bytes(),
				***REMOVED***)
			***REMOVED***
		***REMOVED***()

		// Server read loop.
		var nreq int
		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				select ***REMOVED***
				case <-clientDone:
					// If the client's done, it will have reported any errors on its side.
					return nil
				default:
					return err
				***REMOVED***
			***REMOVED***
			switch f := f.(type) ***REMOVED***
			case *WindowUpdateFrame:
			case *SettingsFrame:
				// Wait for the client SETTINGS ack until ending the greet.
				close(greet)
			case *HeadersFrame:
				if !f.HeadersEnded() ***REMOVED***
					return fmt.Errorf("headers should have END_HEADERS be ended: %v", f)
				***REMOVED***
				gotRequest <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
				nreq++
				writeResp <- f.StreamID
				if nreq == maxConcurrent+1 ***REMOVED***
					close(writeResp)
				***REMOVED***
			default:
				return fmt.Errorf("Unexpected client frame %v", f)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	ct.run()
***REMOVED***

func TestAuthorityAddr(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		scheme, authority string
		want              string
	***REMOVED******REMOVED***
		***REMOVED***"http", "foo.com", "foo.com:80"***REMOVED***,
		***REMOVED***"https", "foo.com", "foo.com:443"***REMOVED***,
		***REMOVED***"https", "foo.com:1234", "foo.com:1234"***REMOVED***,
		***REMOVED***"https", "1.2.3.4:1234", "1.2.3.4:1234"***REMOVED***,
		***REMOVED***"https", "1.2.3.4", "1.2.3.4:443"***REMOVED***,
		***REMOVED***"https", "[::1]:1234", "[::1]:1234"***REMOVED***,
		***REMOVED***"https", "[::1]", "[::1]:443"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		got := authorityAddr(tt.scheme, tt.authority)
		if got != tt.want ***REMOVED***
			t.Errorf("authorityAddr(%q, %q) = %q; want %q", tt.scheme, tt.authority, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Issue 20448: stop allocating for DATA frames' payload after
// Response.Body.Close is called.
func TestTransportAllocationsAfterResponseBodyClose(t *testing.T) ***REMOVED***
	megabyteZero := make([]byte, 1<<20)

	writeErr := make(chan error, 1)

	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.(http.Flusher).Flush()
		var sum int64
		for i := 0; i < 100; i++ ***REMOVED***
			n, err := w.Write(megabyteZero)
			sum += int64(n)
			if err != nil ***REMOVED***
				writeErr <- err
				return
			***REMOVED***
		***REMOVED***
		t.Logf("wrote all %d bytes", sum)
		writeErr <- nil
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()
	c := &http.Client***REMOVED***Transport: tr***REMOVED***
	res, err := c.Get(st.ts.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	var buf [1]byte
	if _, err := res.Body.Read(buf[:]); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	if err := res.Body.Close(); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	trb, ok := res.Body.(transportResponseBody)
	if !ok ***REMOVED***
		t.Fatalf("res.Body = %T; want transportResponseBody", res.Body)
	***REMOVED***
	if trb.cs.bufPipe.b != nil ***REMOVED***
		t.Errorf("response body pipe is still open")
	***REMOVED***

	gotErr := <-writeErr
	if gotErr == nil ***REMOVED***
		t.Errorf("Handler unexpectedly managed to write its entire response without getting an error")
	***REMOVED*** else if gotErr != errStreamClosed ***REMOVED***
		t.Errorf("Handler Write err = %v; want errStreamClosed", gotErr)
	***REMOVED***
***REMOVED***

// Issue 18891: make sure Request.Body == NoBody means no DATA frame
// is ever sent, even if empty.
func TestTransportNoBodyMeansNoDATA(t *testing.T) ***REMOVED***
	ct := newClientTester(t)

	unblockClient := make(chan bool)

	ct.client = func() error ***REMOVED***
		req, _ := http.NewRequest("GET", "https://dummy.tld/", go18httpNoBody())
		ct.tr.RoundTrip(req)
		<-unblockClient
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		defer close(unblockClient)
		defer ct.cc.(*net.TCPConn).Close()
		ct.greet()

		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				return fmt.Errorf("ReadFrame while waiting for Headers: %v", err)
			***REMOVED***
			switch f := f.(type) ***REMOVED***
			default:
				return fmt.Errorf("Got %T; want HeadersFrame", f)
			case *WindowUpdateFrame, *SettingsFrame:
				continue
			case *HeadersFrame:
				if !f.StreamEnded() ***REMOVED***
					return fmt.Errorf("got headers frame without END_STREAM")
				***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

func benchSimpleRoundTrip(b *testing.B, nHeaders int) ***REMOVED***
	defer disableGoroutineTracking()()
	b.ReportAllocs()
	st := newServerTester(b,
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		***REMOVED***,
		optOnlyServer,
		optQuiet,
	)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	req, err := http.NewRequest("GET", st.ts.URL, nil)
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	for i := 0; i < nHeaders; i++ ***REMOVED***
		name := fmt.Sprint("A-", i)
		req.Header.Set(name, "*")
	***REMOVED***

	b.ResetTimer()

	for i := 0; i < b.N; i++ ***REMOVED***
		res, err := tr.RoundTrip(req)
		if err != nil ***REMOVED***
			if res != nil ***REMOVED***
				res.Body.Close()
			***REMOVED***
			b.Fatalf("RoundTrip err = %v; want nil", err)
		***REMOVED***
		res.Body.Close()
		if res.StatusCode != http.StatusOK ***REMOVED***
			b.Fatalf("Response code = %v; want %v", res.StatusCode, http.StatusOK)
		***REMOVED***
	***REMOVED***
***REMOVED***

type infiniteReader struct***REMOVED******REMOVED***

func (r infiniteReader) Read(b []byte) (int, error) ***REMOVED***
	return len(b), nil
***REMOVED***

// Issue 20521: it is not an error to receive a response and end stream
// from the server without the body being consumed.
func TestTransportResponseAndResetWithoutConsumingBodyRace(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.WriteHeader(http.StatusOK)
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	// The request body needs to be big enough to trigger flow control.
	req, _ := http.NewRequest("PUT", st.ts.URL, infiniteReader***REMOVED******REMOVED***)
	res, err := tr.RoundTrip(req)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if res.StatusCode != http.StatusOK ***REMOVED***
		t.Fatalf("Response code = %v; want %v", res.StatusCode, http.StatusOK)
	***REMOVED***
***REMOVED***

// Verify transport doesn't crash when receiving bogus response lacking a :status header.
// Issue 22880.
func TestTransportHandlesInvalidStatuslessResponse(t *testing.T) ***REMOVED***
	ct := newClientTester(t)
	ct.client = func() error ***REMOVED***
		req, _ := http.NewRequest("GET", "https://dummy.tld/", nil)
		_, err := ct.tr.RoundTrip(req)
		const substr = "malformed response from server: missing status pseudo header"
		if !strings.Contains(fmt.Sprint(err), substr) ***REMOVED***
			return fmt.Errorf("RoundTrip error = %v; want substring %q", err, substr)
		***REMOVED***
		return nil
	***REMOVED***
	ct.server = func() error ***REMOVED***
		ct.greet()
		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)

		for ***REMOVED***
			f, err := ct.fr.ReadFrame()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			switch f := f.(type) ***REMOVED***
			case *HeadersFrame:
				enc.WriteField(hpack.HeaderField***REMOVED***Name: "content-type", Value: "text/html"***REMOVED***) // no :status header
				ct.fr.WriteHeaders(HeadersFrameParam***REMOVED***
					StreamID:      f.StreamID,
					EndHeaders:    true,
					EndStream:     false, // we'll send some DATA to try to crash the transport
					BlockFragment: buf.Bytes(),
				***REMOVED***)
				ct.fr.WriteData(f.StreamID, true, []byte("payload"))
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ct.run()
***REMOVED***

func BenchmarkClientRequestHeaders(b *testing.B) ***REMOVED***
	b.Run("   0 Headers", func(b *testing.B) ***REMOVED*** benchSimpleRoundTrip(b, 0) ***REMOVED***)
	b.Run("  10 Headers", func(b *testing.B) ***REMOVED*** benchSimpleRoundTrip(b, 10) ***REMOVED***)
	b.Run(" 100 Headers", func(b *testing.B) ***REMOVED*** benchSimpleRoundTrip(b, 100) ***REMOVED***)
	b.Run("1000 Headers", func(b *testing.B) ***REMOVED*** benchSimpleRoundTrip(b, 1000) ***REMOVED***)
***REMOVED***
