// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

var cstUpgrader = Upgrader***REMOVED***
	Subprotocols:      []string***REMOVED***"p0", "p1"***REMOVED***,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) ***REMOVED***
		http.Error(w, reason.Error(), status)
	***REMOVED***,
***REMOVED***

var cstDialer = Dialer***REMOVED***
	Subprotocols:     []string***REMOVED***"p1", "p2"***REMOVED***,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 30 * time.Second,
***REMOVED***

type cstHandler struct***REMOVED*** *testing.T ***REMOVED***

type cstServer struct ***REMOVED***
	*httptest.Server
	URL string
***REMOVED***

const (
	cstPath       = "/a/b"
	cstRawQuery   = "x=y"
	cstRequestURI = cstPath + "?" + cstRawQuery
)

func newServer(t *testing.T) *cstServer ***REMOVED***
	var s cstServer
	s.Server = httptest.NewServer(cstHandler***REMOVED***t***REMOVED***)
	s.Server.URL += cstRequestURI
	s.URL = makeWsProto(s.Server.URL)
	return &s
***REMOVED***

func newTLSServer(t *testing.T) *cstServer ***REMOVED***
	var s cstServer
	s.Server = httptest.NewTLSServer(cstHandler***REMOVED***t***REMOVED***)
	s.Server.URL += cstRequestURI
	s.URL = makeWsProto(s.Server.URL)
	return &s
***REMOVED***

func (t cstHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	if r.URL.Path != cstPath ***REMOVED***
		t.Logf("path=%v, want %v", r.URL.Path, cstPath)
		http.Error(w, "bad path", 400)
		return
	***REMOVED***
	if r.URL.RawQuery != cstRawQuery ***REMOVED***
		t.Logf("query=%v, want %v", r.URL.RawQuery, cstRawQuery)
		http.Error(w, "bad path", 400)
		return
	***REMOVED***
	subprotos := Subprotocols(r)
	if !reflect.DeepEqual(subprotos, cstDialer.Subprotocols) ***REMOVED***
		t.Logf("subprotols=%v, want %v", subprotos, cstDialer.Subprotocols)
		http.Error(w, "bad protocol", 400)
		return
	***REMOVED***
	ws, err := cstUpgrader.Upgrade(w, r, http.Header***REMOVED***"Set-Cookie": ***REMOVED***"sessionID=1234"***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Logf("Upgrade: %v", err)
		return
	***REMOVED***
	defer ws.Close()

	if ws.Subprotocol() != "p1" ***REMOVED***
		t.Logf("Subprotocol() = %s, want p1", ws.Subprotocol())
		ws.Close()
		return
	***REMOVED***
	op, rd, err := ws.NextReader()
	if err != nil ***REMOVED***
		t.Logf("NextReader: %v", err)
		return
	***REMOVED***
	wr, err := ws.NextWriter(op)
	if err != nil ***REMOVED***
		t.Logf("NextWriter: %v", err)
		return
	***REMOVED***
	if _, err = io.Copy(wr, rd); err != nil ***REMOVED***
		t.Logf("NextWriter: %v", err)
		return
	***REMOVED***
	if err := wr.Close(); err != nil ***REMOVED***
		t.Logf("Close: %v", err)
		return
	***REMOVED***
***REMOVED***

func makeWsProto(s string) string ***REMOVED***
	return "ws" + strings.TrimPrefix(s, "http")
***REMOVED***

func sendRecv(t *testing.T, ws *Conn) ***REMOVED***
	const message = "Hello World!"
	if err := ws.SetWriteDeadline(time.Now().Add(time.Second)); err != nil ***REMOVED***
		t.Fatalf("SetWriteDeadline: %v", err)
	***REMOVED***
	if err := ws.WriteMessage(TextMessage, []byte(message)); err != nil ***REMOVED***
		t.Fatalf("WriteMessage: %v", err)
	***REMOVED***
	if err := ws.SetReadDeadline(time.Now().Add(time.Second)); err != nil ***REMOVED***
		t.Fatalf("SetReadDeadline: %v", err)
	***REMOVED***
	_, p, err := ws.ReadMessage()
	if err != nil ***REMOVED***
		t.Fatalf("ReadMessage: %v", err)
	***REMOVED***
	if string(p) != message ***REMOVED***
		t.Fatalf("message=%s, want %s", p, message)
	***REMOVED***
***REMOVED***

func TestProxyDial(t *testing.T) ***REMOVED***

	s := newServer(t)
	defer s.Close()

	surl, _ := url.Parse(s.Server.URL)

	cstDialer := cstDialer // make local copy for modification on next line.
	cstDialer.Proxy = http.ProxyURL(surl)

	connect := false
	origHandler := s.Server.Config.Handler

	// Capture the request Host header.
	s.Server.Config.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			if r.Method == "CONNECT" ***REMOVED***
				connect = true
				w.WriteHeader(200)
				return
			***REMOVED***

			if !connect ***REMOVED***
				t.Log("connect not received")
				http.Error(w, "connect not received", 405)
				return
			***REMOVED***
			origHandler.ServeHTTP(w, r)
		***REMOVED***)

	ws, _, err := cstDialer.Dial(s.URL, nil)
	if err != nil ***REMOVED***
		t.Fatalf("Dial: %v", err)
	***REMOVED***
	defer ws.Close()
	sendRecv(t, ws)
***REMOVED***

func TestProxyAuthorizationDial(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Close()

	surl, _ := url.Parse(s.Server.URL)
	surl.User = url.UserPassword("username", "password")

	cstDialer := cstDialer // make local copy for modification on next line.
	cstDialer.Proxy = http.ProxyURL(surl)

	connect := false
	origHandler := s.Server.Config.Handler

	// Capture the request Host header.
	s.Server.Config.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			proxyAuth := r.Header.Get("Proxy-Authorization")
			expectedProxyAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("username:password"))
			if r.Method == "CONNECT" && proxyAuth == expectedProxyAuth ***REMOVED***
				connect = true
				w.WriteHeader(200)
				return
			***REMOVED***

			if !connect ***REMOVED***
				t.Log("connect with proxy authorization not received")
				http.Error(w, "connect with proxy authorization not received", 405)
				return
			***REMOVED***
			origHandler.ServeHTTP(w, r)
		***REMOVED***)

	ws, _, err := cstDialer.Dial(s.URL, nil)
	if err != nil ***REMOVED***
		t.Fatalf("Dial: %v", err)
	***REMOVED***
	defer ws.Close()
	sendRecv(t, ws)
***REMOVED***

func TestDial(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Close()

	ws, _, err := cstDialer.Dial(s.URL, nil)
	if err != nil ***REMOVED***
		t.Fatalf("Dial: %v", err)
	***REMOVED***
	defer ws.Close()
	sendRecv(t, ws)
***REMOVED***

func TestDialCookieJar(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Close()

	jar, _ := cookiejar.New(nil)
	d := cstDialer
	d.Jar = jar

	u, _ := url.Parse(s.URL)

	switch u.Scheme ***REMOVED***
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	***REMOVED***

	cookies := []*http.Cookie***REMOVED******REMOVED***Name: "gorilla", Value: "ws", Path: "/"***REMOVED******REMOVED***
	d.Jar.SetCookies(u, cookies)

	ws, _, err := d.Dial(s.URL, nil)
	if err != nil ***REMOVED***
		t.Fatalf("Dial: %v", err)
	***REMOVED***
	defer ws.Close()

	var gorilla string
	var sessionID string
	for _, c := range d.Jar.Cookies(u) ***REMOVED***
		if c.Name == "gorilla" ***REMOVED***
			gorilla = c.Value
		***REMOVED***

		if c.Name == "sessionID" ***REMOVED***
			sessionID = c.Value
		***REMOVED***
	***REMOVED***
	if gorilla != "ws" ***REMOVED***
		t.Error("Cookie not present in jar.")
	***REMOVED***

	if sessionID != "1234" ***REMOVED***
		t.Error("Set-Cookie not received from the server.")
	***REMOVED***

	sendRecv(t, ws)
***REMOVED***

func TestDialTLS(t *testing.T) ***REMOVED***
	s := newTLSServer(t)
	defer s.Close()

	certs := x509.NewCertPool()
	for _, c := range s.TLS.Certificates ***REMOVED***
		roots, err := x509.ParseCertificates(c.Certificate[len(c.Certificate)-1])
		if err != nil ***REMOVED***
			t.Fatalf("error parsing server's root cert: %v", err)
		***REMOVED***
		for _, root := range roots ***REMOVED***
			certs.AddCert(root)
		***REMOVED***
	***REMOVED***

	d := cstDialer
	d.TLSClientConfig = &tls.Config***REMOVED***RootCAs: certs***REMOVED***
	ws, _, err := d.Dial(s.URL, nil)
	if err != nil ***REMOVED***
		t.Fatalf("Dial: %v", err)
	***REMOVED***
	defer ws.Close()
	sendRecv(t, ws)
***REMOVED***

func xTestDialTLSBadCert(t *testing.T) ***REMOVED***
	// This test is deactivated because of noisy logging from the net/http package.
	s := newTLSServer(t)
	defer s.Close()

	ws, _, err := cstDialer.Dial(s.URL, nil)
	if err == nil ***REMOVED***
		ws.Close()
		t.Fatalf("Dial: nil")
	***REMOVED***
***REMOVED***

func TestDialTLSNoVerify(t *testing.T) ***REMOVED***
	s := newTLSServer(t)
	defer s.Close()

	d := cstDialer
	d.TLSClientConfig = &tls.Config***REMOVED***InsecureSkipVerify: true***REMOVED***
	ws, _, err := d.Dial(s.URL, nil)
	if err != nil ***REMOVED***
		t.Fatalf("Dial: %v", err)
	***REMOVED***
	defer ws.Close()
	sendRecv(t, ws)
***REMOVED***

func TestDialTimeout(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Close()

	d := cstDialer
	d.HandshakeTimeout = -1
	ws, _, err := d.Dial(s.URL, nil)
	if err == nil ***REMOVED***
		ws.Close()
		t.Fatalf("Dial: nil")
	***REMOVED***
***REMOVED***

func TestDialBadScheme(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Close()

	ws, _, err := cstDialer.Dial(s.Server.URL, nil)
	if err == nil ***REMOVED***
		ws.Close()
		t.Fatalf("Dial: nil")
	***REMOVED***
***REMOVED***

func TestDialBadOrigin(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Close()

	ws, resp, err := cstDialer.Dial(s.URL, http.Header***REMOVED***"Origin": ***REMOVED***"bad"***REMOVED******REMOVED***)
	if err == nil ***REMOVED***
		ws.Close()
		t.Fatalf("Dial: nil")
	***REMOVED***
	if resp == nil ***REMOVED***
		t.Fatalf("resp=nil, err=%v", err)
	***REMOVED***
	if resp.StatusCode != http.StatusForbidden ***REMOVED***
		t.Fatalf("status=%d, want %d", resp.StatusCode, http.StatusForbidden)
	***REMOVED***
***REMOVED***

func TestDialBadHeader(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Close()

	for _, k := range []string***REMOVED***"Upgrade",
		"Connection",
		"Sec-Websocket-Key",
		"Sec-Websocket-Version",
		"Sec-Websocket-Protocol"***REMOVED*** ***REMOVED***
		h := http.Header***REMOVED******REMOVED***
		h.Set(k, "bad")
		ws, _, err := cstDialer.Dial(s.URL, http.Header***REMOVED***"Origin": ***REMOVED***"bad"***REMOVED******REMOVED***)
		if err == nil ***REMOVED***
			ws.Close()
			t.Errorf("Dial with header %s returned nil", k)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBadMethod(t *testing.T) ***REMOVED***
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		ws, err := cstUpgrader.Upgrade(w, r, nil)
		if err == nil ***REMOVED***
			t.Errorf("handshake succeeded, expect fail")
			ws.Close()
		***REMOVED***
	***REMOVED***))
	defer s.Close()

	req, err := http.NewRequest("POST", s.URL, strings.NewReader(""))
	if err != nil ***REMOVED***
		t.Fatalf("NewRequest returned error %v", err)
	***REMOVED***
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-Websocket-Version", "13")

	resp, err := http.DefaultClient.Do(req)
	if err != nil ***REMOVED***
		t.Fatalf("Do returned error %v", err)
	***REMOVED***
	resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed ***REMOVED***
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	***REMOVED***
***REMOVED***

func TestHandshake(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Close()

	ws, resp, err := cstDialer.Dial(s.URL, http.Header***REMOVED***"Origin": ***REMOVED***s.URL***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("Dial: %v", err)
	***REMOVED***
	defer ws.Close()

	var sessionID string
	for _, c := range resp.Cookies() ***REMOVED***
		if c.Name == "sessionID" ***REMOVED***
			sessionID = c.Value
		***REMOVED***
	***REMOVED***
	if sessionID != "1234" ***REMOVED***
		t.Error("Set-Cookie not received from the server.")
	***REMOVED***

	if ws.Subprotocol() != "p1" ***REMOVED***
		t.Errorf("ws.Subprotocol() = %s, want p1", ws.Subprotocol())
	***REMOVED***
	sendRecv(t, ws)
***REMOVED***

func TestRespOnBadHandshake(t *testing.T) ***REMOVED***
	const expectedStatus = http.StatusGone
	const expectedBody = "This is the response body."

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.WriteHeader(expectedStatus)
		io.WriteString(w, expectedBody)
	***REMOVED***))
	defer s.Close()

	ws, resp, err := cstDialer.Dial(makeWsProto(s.URL), nil)
	if err == nil ***REMOVED***
		ws.Close()
		t.Fatalf("Dial: nil")
	***REMOVED***

	if resp == nil ***REMOVED***
		t.Fatalf("resp=nil, err=%v", err)
	***REMOVED***

	if resp.StatusCode != expectedStatus ***REMOVED***
		t.Errorf("resp.StatusCode=%d, want %d", resp.StatusCode, expectedStatus)
	***REMOVED***

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil ***REMOVED***
		t.Fatalf("ReadFull(resp.Body) returned error %v", err)
	***REMOVED***

	if string(p) != expectedBody ***REMOVED***
		t.Errorf("resp.Body=%s, want %s", p, expectedBody)
	***REMOVED***
***REMOVED***

// TestHostHeader confirms that the host header provided in the call to Dial is
// sent to the server.
func TestHostHeader(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Close()

	specifiedHost := make(chan string, 1)
	origHandler := s.Server.Config.Handler

	// Capture the request Host header.
	s.Server.Config.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			specifiedHost <- r.Host
			origHandler.ServeHTTP(w, r)
		***REMOVED***)

	ws, _, err := cstDialer.Dial(s.URL, http.Header***REMOVED***"Host": ***REMOVED***"testhost"***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("Dial: %v", err)
	***REMOVED***
	defer ws.Close()

	if gotHost := <-specifiedHost; gotHost != "testhost" ***REMOVED***
		t.Fatalf("gotHost = %q, want \"testhost\"", gotHost)
	***REMOVED***

	sendRecv(t, ws)
***REMOVED***

func TestDialCompression(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Close()

	dialer := cstDialer
	dialer.EnableCompression = true
	ws, _, err := dialer.Dial(s.URL, nil)
	if err != nil ***REMOVED***
		t.Fatalf("Dial: %v", err)
	***REMOVED***
	defer ws.Close()
	sendRecv(t, ws)
***REMOVED***

func TestSocksProxyDial(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Close()

	proxyListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil ***REMOVED***
		t.Fatalf("listen failed: %v", err)
	***REMOVED***
	defer proxyListener.Close()
	go func() ***REMOVED***
		c1, err := proxyListener.Accept()
		if err != nil ***REMOVED***
			t.Errorf("proxy accept failed: %v", err)
			return
		***REMOVED***
		defer c1.Close()

		c1.SetDeadline(time.Now().Add(30 * time.Second))

		buf := make([]byte, 32)
		if _, err := io.ReadFull(c1, buf[:3]); err != nil ***REMOVED***
			t.Errorf("read failed: %v", err)
			return
		***REMOVED***
		if want := []byte***REMOVED***5, 1, 0***REMOVED***; !bytes.Equal(want, buf[:len(want)]) ***REMOVED***
			t.Errorf("read %x, want %x", buf[:len(want)], want)
		***REMOVED***
		if _, err := c1.Write([]byte***REMOVED***5, 0***REMOVED***); err != nil ***REMOVED***
			t.Errorf("write failed: %v", err)
			return
		***REMOVED***
		if _, err := io.ReadFull(c1, buf[:10]); err != nil ***REMOVED***
			t.Errorf("read failed: %v", err)
			return
		***REMOVED***
		if want := []byte***REMOVED***5, 1, 0, 1***REMOVED***; !bytes.Equal(want, buf[:len(want)]) ***REMOVED***
			t.Errorf("read %x, want %x", buf[:len(want)], want)
			return
		***REMOVED***
		buf[1] = 0
		if _, err := c1.Write(buf[:10]); err != nil ***REMOVED***
			t.Errorf("write failed: %v", err)
			return
		***REMOVED***

		ip := net.IP(buf[4:8])
		port := binary.BigEndian.Uint16(buf[8:10])

		c2, err := net.DialTCP("tcp", nil, &net.TCPAddr***REMOVED***IP: ip, Port: int(port)***REMOVED***)
		if err != nil ***REMOVED***
			t.Errorf("dial failed; %v", err)
			return
		***REMOVED***
		defer c2.Close()
		done := make(chan struct***REMOVED******REMOVED***)
		go func() ***REMOVED***
			io.Copy(c1, c2)
			close(done)
		***REMOVED***()
		io.Copy(c2, c1)
		<-done
	***REMOVED***()

	purl, err := url.Parse("socks5://" + proxyListener.Addr().String())
	if err != nil ***REMOVED***
		t.Fatalf("parse failed: %v", err)
	***REMOVED***

	cstDialer := cstDialer // make local copy for modification on next line.
	cstDialer.Proxy = http.ProxyURL(purl)

	ws, _, err := cstDialer.Dial(s.URL, nil)
	if err != nil ***REMOVED***
		t.Fatalf("Dial: %v", err)
	***REMOVED***
	defer ws.Close()
	sendRecv(t, ws)
***REMOVED***
