// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

var serverAddr string
var once sync.Once

func echoServer(ws *Conn) ***REMOVED***
	defer ws.Close()
	io.Copy(ws, ws)
***REMOVED***

type Count struct ***REMOVED***
	S string
	N int
***REMOVED***

func countServer(ws *Conn) ***REMOVED***
	defer ws.Close()
	for ***REMOVED***
		var count Count
		err := JSON.Receive(ws, &count)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		count.N++
		count.S = strings.Repeat(count.S, count.N)
		err = JSON.Send(ws, count)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

type testCtrlAndDataHandler struct ***REMOVED***
	hybiFrameHandler
***REMOVED***

func (h *testCtrlAndDataHandler) WritePing(b []byte) (int, error) ***REMOVED***
	h.hybiFrameHandler.conn.wio.Lock()
	defer h.hybiFrameHandler.conn.wio.Unlock()
	w, err := h.hybiFrameHandler.conn.frameWriterFactory.NewFrameWriter(PingFrame)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	n, err := w.Write(b)
	w.Close()
	return n, err
***REMOVED***

func ctrlAndDataServer(ws *Conn) ***REMOVED***
	defer ws.Close()
	h := &testCtrlAndDataHandler***REMOVED***hybiFrameHandler: hybiFrameHandler***REMOVED***conn: ws***REMOVED******REMOVED***
	ws.frameHandler = h

	go func() ***REMOVED***
		for i := 0; ; i++ ***REMOVED***
			var b []byte
			if i%2 != 0 ***REMOVED*** // with or without payload
				b = []byte(fmt.Sprintf("#%d-CONTROL-FRAME-FROM-SERVER", i))
			***REMOVED***
			if _, err := h.WritePing(b); err != nil ***REMOVED***
				break
			***REMOVED***
			if _, err := h.WritePong(b); err != nil ***REMOVED*** // unsolicited pong
				break
			***REMOVED***
			time.Sleep(10 * time.Millisecond)
		***REMOVED***
	***REMOVED***()

	b := make([]byte, 128)
	for ***REMOVED***
		n, err := ws.Read(b)
		if err != nil ***REMOVED***
			break
		***REMOVED***
		if _, err := ws.Write(b[:n]); err != nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func subProtocolHandshake(config *Config, req *http.Request) error ***REMOVED***
	for _, proto := range config.Protocol ***REMOVED***
		if proto == "chat" ***REMOVED***
			config.Protocol = []string***REMOVED***proto***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return ErrBadWebSocketProtocol
***REMOVED***

func subProtoServer(ws *Conn) ***REMOVED***
	for _, proto := range ws.Config().Protocol ***REMOVED***
		io.WriteString(ws, proto)
	***REMOVED***
***REMOVED***

func startServer() ***REMOVED***
	http.Handle("/echo", Handler(echoServer))
	http.Handle("/count", Handler(countServer))
	http.Handle("/ctrldata", Handler(ctrlAndDataServer))
	subproto := Server***REMOVED***
		Handshake: subProtocolHandshake,
		Handler:   Handler(subProtoServer),
	***REMOVED***
	http.Handle("/subproto", subproto)
	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
	log.Print("Test WebSocket server listening on ", serverAddr)
***REMOVED***

func newConfig(t *testing.T, path string) *Config ***REMOVED***
	config, _ := NewConfig(fmt.Sprintf("ws://%s%s", serverAddr, path), "http://localhost")
	return config
***REMOVED***

func TestEcho(t *testing.T) ***REMOVED***
	once.Do(startServer)

	// websocket.Dial()
	client, err := net.Dial("tcp", serverAddr)
	if err != nil ***REMOVED***
		t.Fatal("dialing", err)
	***REMOVED***
	conn, err := NewClient(newConfig(t, "/echo"), client)
	if err != nil ***REMOVED***
		t.Errorf("WebSocket handshake error: %v", err)
		return
	***REMOVED***

	msg := []byte("hello, world\n")
	if _, err := conn.Write(msg); err != nil ***REMOVED***
		t.Errorf("Write: %v", err)
	***REMOVED***
	var actual_msg = make([]byte, 512)
	n, err := conn.Read(actual_msg)
	if err != nil ***REMOVED***
		t.Errorf("Read: %v", err)
	***REMOVED***
	actual_msg = actual_msg[0:n]
	if !bytes.Equal(msg, actual_msg) ***REMOVED***
		t.Errorf("Echo: expected %q got %q", msg, actual_msg)
	***REMOVED***
	conn.Close()
***REMOVED***

func TestAddr(t *testing.T) ***REMOVED***
	once.Do(startServer)

	// websocket.Dial()
	client, err := net.Dial("tcp", serverAddr)
	if err != nil ***REMOVED***
		t.Fatal("dialing", err)
	***REMOVED***
	conn, err := NewClient(newConfig(t, "/echo"), client)
	if err != nil ***REMOVED***
		t.Errorf("WebSocket handshake error: %v", err)
		return
	***REMOVED***

	ra := conn.RemoteAddr().String()
	if !strings.HasPrefix(ra, "ws://") || !strings.HasSuffix(ra, "/echo") ***REMOVED***
		t.Errorf("Bad remote addr: %v", ra)
	***REMOVED***
	la := conn.LocalAddr().String()
	if !strings.HasPrefix(la, "http://") ***REMOVED***
		t.Errorf("Bad local addr: %v", la)
	***REMOVED***
	conn.Close()
***REMOVED***

func TestCount(t *testing.T) ***REMOVED***
	once.Do(startServer)

	// websocket.Dial()
	client, err := net.Dial("tcp", serverAddr)
	if err != nil ***REMOVED***
		t.Fatal("dialing", err)
	***REMOVED***
	conn, err := NewClient(newConfig(t, "/count"), client)
	if err != nil ***REMOVED***
		t.Errorf("WebSocket handshake error: %v", err)
		return
	***REMOVED***

	var count Count
	count.S = "hello"
	if err := JSON.Send(conn, count); err != nil ***REMOVED***
		t.Errorf("Write: %v", err)
	***REMOVED***
	if err := JSON.Receive(conn, &count); err != nil ***REMOVED***
		t.Errorf("Read: %v", err)
	***REMOVED***
	if count.N != 1 ***REMOVED***
		t.Errorf("count: expected %d got %d", 1, count.N)
	***REMOVED***
	if count.S != "hello" ***REMOVED***
		t.Errorf("count: expected %q got %q", "hello", count.S)
	***REMOVED***
	if err := JSON.Send(conn, count); err != nil ***REMOVED***
		t.Errorf("Write: %v", err)
	***REMOVED***
	if err := JSON.Receive(conn, &count); err != nil ***REMOVED***
		t.Errorf("Read: %v", err)
	***REMOVED***
	if count.N != 2 ***REMOVED***
		t.Errorf("count: expected %d got %d", 2, count.N)
	***REMOVED***
	if count.S != "hellohello" ***REMOVED***
		t.Errorf("count: expected %q got %q", "hellohello", count.S)
	***REMOVED***
	conn.Close()
***REMOVED***

func TestWithQuery(t *testing.T) ***REMOVED***
	once.Do(startServer)

	client, err := net.Dial("tcp", serverAddr)
	if err != nil ***REMOVED***
		t.Fatal("dialing", err)
	***REMOVED***

	config := newConfig(t, "/echo")
	config.Location, err = url.ParseRequestURI(fmt.Sprintf("ws://%s/echo?q=v", serverAddr))
	if err != nil ***REMOVED***
		t.Fatal("location url", err)
	***REMOVED***

	ws, err := NewClient(config, client)
	if err != nil ***REMOVED***
		t.Errorf("WebSocket handshake: %v", err)
		return
	***REMOVED***
	ws.Close()
***REMOVED***

func testWithProtocol(t *testing.T, subproto []string) (string, error) ***REMOVED***
	once.Do(startServer)

	client, err := net.Dial("tcp", serverAddr)
	if err != nil ***REMOVED***
		t.Fatal("dialing", err)
	***REMOVED***

	config := newConfig(t, "/subproto")
	config.Protocol = subproto

	ws, err := NewClient(config, client)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	msg := make([]byte, 16)
	n, err := ws.Read(msg)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	ws.Close()
	return string(msg[:n]), nil
***REMOVED***

func TestWithProtocol(t *testing.T) ***REMOVED***
	proto, err := testWithProtocol(t, []string***REMOVED***"chat"***REMOVED***)
	if err != nil ***REMOVED***
		t.Errorf("SubProto: unexpected error: %v", err)
	***REMOVED***
	if proto != "chat" ***REMOVED***
		t.Errorf("SubProto: expected %q, got %q", "chat", proto)
	***REMOVED***
***REMOVED***

func TestWithTwoProtocol(t *testing.T) ***REMOVED***
	proto, err := testWithProtocol(t, []string***REMOVED***"test", "chat"***REMOVED***)
	if err != nil ***REMOVED***
		t.Errorf("SubProto: unexpected error: %v", err)
	***REMOVED***
	if proto != "chat" ***REMOVED***
		t.Errorf("SubProto: expected %q, got %q", "chat", proto)
	***REMOVED***
***REMOVED***

func TestWithBadProtocol(t *testing.T) ***REMOVED***
	_, err := testWithProtocol(t, []string***REMOVED***"test"***REMOVED***)
	if err != ErrBadStatus ***REMOVED***
		t.Errorf("SubProto: expected %v, got %v", ErrBadStatus, err)
	***REMOVED***
***REMOVED***

func TestHTTP(t *testing.T) ***REMOVED***
	once.Do(startServer)

	// If the client did not send a handshake that matches the protocol
	// specification, the server MUST return an HTTP response with an
	// appropriate error code (such as 400 Bad Request)
	resp, err := http.Get(fmt.Sprintf("http://%s/echo", serverAddr))
	if err != nil ***REMOVED***
		t.Errorf("Get: error %#v", err)
		return
	***REMOVED***
	if resp == nil ***REMOVED***
		t.Error("Get: resp is null")
		return
	***REMOVED***
	if resp.StatusCode != http.StatusBadRequest ***REMOVED***
		t.Errorf("Get: expected %q got %q", http.StatusBadRequest, resp.StatusCode)
	***REMOVED***
***REMOVED***

func TestTrailingSpaces(t *testing.T) ***REMOVED***
	// http://code.google.com/p/go/issues/detail?id=955
	// The last runs of this create keys with trailing spaces that should not be
	// generated by the client.
	once.Do(startServer)
	config := newConfig(t, "/echo")
	for i := 0; i < 30; i++ ***REMOVED***
		// body
		ws, err := DialConfig(config)
		if err != nil ***REMOVED***
			t.Errorf("Dial #%d failed: %v", i, err)
			break
		***REMOVED***
		ws.Close()
	***REMOVED***
***REMOVED***

func TestDialConfigBadVersion(t *testing.T) ***REMOVED***
	once.Do(startServer)
	config := newConfig(t, "/echo")
	config.Version = 1234

	_, err := DialConfig(config)

	if dialerr, ok := err.(*DialError); ok ***REMOVED***
		if dialerr.Err != ErrBadProtocolVersion ***REMOVED***
			t.Errorf("dial expected err %q but got %q", ErrBadProtocolVersion, dialerr.Err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDialConfigWithDialer(t *testing.T) ***REMOVED***
	once.Do(startServer)
	config := newConfig(t, "/echo")
	config.Dialer = &net.Dialer***REMOVED***
		Deadline: time.Now().Add(-time.Minute),
	***REMOVED***
	_, err := DialConfig(config)
	dialerr, ok := err.(*DialError)
	if !ok ***REMOVED***
		t.Fatalf("DialError expected, got %#v", err)
	***REMOVED***
	neterr, ok := dialerr.Err.(*net.OpError)
	if !ok ***REMOVED***
		t.Fatalf("net.OpError error expected, got %#v", dialerr.Err)
	***REMOVED***
	if !neterr.Timeout() ***REMOVED***
		t.Fatalf("expected timeout error, got %#v", neterr)
	***REMOVED***
***REMOVED***

func TestSmallBuffer(t *testing.T) ***REMOVED***
	// http://code.google.com/p/go/issues/detail?id=1145
	// Read should be able to handle reading a fragment of a frame.
	once.Do(startServer)

	// websocket.Dial()
	client, err := net.Dial("tcp", serverAddr)
	if err != nil ***REMOVED***
		t.Fatal("dialing", err)
	***REMOVED***
	conn, err := NewClient(newConfig(t, "/echo"), client)
	if err != nil ***REMOVED***
		t.Errorf("WebSocket handshake error: %v", err)
		return
	***REMOVED***

	msg := []byte("hello, world\n")
	if _, err := conn.Write(msg); err != nil ***REMOVED***
		t.Errorf("Write: %v", err)
	***REMOVED***
	var small_msg = make([]byte, 8)
	n, err := conn.Read(small_msg)
	if err != nil ***REMOVED***
		t.Errorf("Read: %v", err)
	***REMOVED***
	if !bytes.Equal(msg[:len(small_msg)], small_msg) ***REMOVED***
		t.Errorf("Echo: expected %q got %q", msg[:len(small_msg)], small_msg)
	***REMOVED***
	var second_msg = make([]byte, len(msg))
	n, err = conn.Read(second_msg)
	if err != nil ***REMOVED***
		t.Errorf("Read: %v", err)
	***REMOVED***
	second_msg = second_msg[0:n]
	if !bytes.Equal(msg[len(small_msg):], second_msg) ***REMOVED***
		t.Errorf("Echo: expected %q got %q", msg[len(small_msg):], second_msg)
	***REMOVED***
	conn.Close()
***REMOVED***

var parseAuthorityTests = []struct ***REMOVED***
	in  *url.URL
	out string
***REMOVED******REMOVED***
	***REMOVED***
		&url.URL***REMOVED***
			Scheme: "ws",
			Host:   "www.google.com",
		***REMOVED***,
		"www.google.com:80",
	***REMOVED***,
	***REMOVED***
		&url.URL***REMOVED***
			Scheme: "wss",
			Host:   "www.google.com",
		***REMOVED***,
		"www.google.com:443",
	***REMOVED***,
	***REMOVED***
		&url.URL***REMOVED***
			Scheme: "ws",
			Host:   "www.google.com:80",
		***REMOVED***,
		"www.google.com:80",
	***REMOVED***,
	***REMOVED***
		&url.URL***REMOVED***
			Scheme: "wss",
			Host:   "www.google.com:443",
		***REMOVED***,
		"www.google.com:443",
	***REMOVED***,
	// some invalid ones for parseAuthority. parseAuthority doesn't
	// concern itself with the scheme unless it actually knows about it
	***REMOVED***
		&url.URL***REMOVED***
			Scheme: "http",
			Host:   "www.google.com",
		***REMOVED***,
		"www.google.com",
	***REMOVED***,
	***REMOVED***
		&url.URL***REMOVED***
			Scheme: "http",
			Host:   "www.google.com:80",
		***REMOVED***,
		"www.google.com:80",
	***REMOVED***,
	***REMOVED***
		&url.URL***REMOVED***
			Scheme: "asdf",
			Host:   "127.0.0.1",
		***REMOVED***,
		"127.0.0.1",
	***REMOVED***,
	***REMOVED***
		&url.URL***REMOVED***
			Scheme: "asdf",
			Host:   "www.google.com",
		***REMOVED***,
		"www.google.com",
	***REMOVED***,
***REMOVED***

func TestParseAuthority(t *testing.T) ***REMOVED***
	for _, tt := range parseAuthorityTests ***REMOVED***
		out := parseAuthority(tt.in)
		if out != tt.out ***REMOVED***
			t.Errorf("got %v; want %v", out, tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

type closerConn struct ***REMOVED***
	net.Conn
	closed int // count of the number of times Close was called
***REMOVED***

func (c *closerConn) Close() error ***REMOVED***
	c.closed++
	return c.Conn.Close()
***REMOVED***

func TestClose(t *testing.T) ***REMOVED***
	if runtime.GOOS == "plan9" ***REMOVED***
		t.Skip("see golang.org/issue/11454")
	***REMOVED***

	once.Do(startServer)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil ***REMOVED***
		t.Fatal("dialing", err)
	***REMOVED***

	cc := closerConn***REMOVED***Conn: conn***REMOVED***

	client, err := NewClient(newConfig(t, "/echo"), &cc)
	if err != nil ***REMOVED***
		t.Fatalf("WebSocket handshake: %v", err)
	***REMOVED***

	// set the deadline to ten minutes ago, which will have expired by the time
	// client.Close sends the close status frame.
	conn.SetDeadline(time.Now().Add(-10 * time.Minute))

	if err := client.Close(); err == nil ***REMOVED***
		t.Errorf("ws.Close(): expected error, got %v", err)
	***REMOVED***
	if cc.closed < 1 ***REMOVED***
		t.Fatalf("ws.Close(): expected underlying ws.rwc.Close to be called > 0 times, got: %v", cc.closed)
	***REMOVED***
***REMOVED***

var originTests = []struct ***REMOVED***
	req    *http.Request
	origin *url.URL
***REMOVED******REMOVED***
	***REMOVED***
		req: &http.Request***REMOVED***
			Header: http.Header***REMOVED***
				"Origin": []string***REMOVED***"http://www.example.com"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		origin: &url.URL***REMOVED***
			Scheme: "http",
			Host:   "www.example.com",
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		req: &http.Request***REMOVED******REMOVED***,
	***REMOVED***,
***REMOVED***

func TestOrigin(t *testing.T) ***REMOVED***
	conf := newConfig(t, "/echo")
	conf.Version = ProtocolVersionHybi13
	for i, tt := range originTests ***REMOVED***
		origin, err := Origin(conf, tt.req)
		if err != nil ***REMOVED***
			t.Error(err)
			continue
		***REMOVED***
		if !reflect.DeepEqual(origin, tt.origin) ***REMOVED***
			t.Errorf("#%d: got origin %v; want %v", i, origin, tt.origin)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCtrlAndData(t *testing.T) ***REMOVED***
	once.Do(startServer)

	c, err := net.Dial("tcp", serverAddr)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ws, err := NewClient(newConfig(t, "/ctrldata"), c)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer ws.Close()

	h := &testCtrlAndDataHandler***REMOVED***hybiFrameHandler: hybiFrameHandler***REMOVED***conn: ws***REMOVED******REMOVED***
	ws.frameHandler = h

	b := make([]byte, 128)
	for i := 0; i < 2; i++ ***REMOVED***
		data := []byte(fmt.Sprintf("#%d-DATA-FRAME-FROM-CLIENT", i))
		if _, err := ws.Write(data); err != nil ***REMOVED***
			t.Fatalf("#%d: %v", i, err)
		***REMOVED***
		var ctrl []byte
		if i%2 != 0 ***REMOVED*** // with or without payload
			ctrl = []byte(fmt.Sprintf("#%d-CONTROL-FRAME-FROM-CLIENT", i))
		***REMOVED***
		if _, err := h.WritePing(ctrl); err != nil ***REMOVED***
			t.Fatalf("#%d: %v", i, err)
		***REMOVED***
		n, err := ws.Read(b)
		if err != nil ***REMOVED***
			t.Fatalf("#%d: %v", i, err)
		***REMOVED***
		if !bytes.Equal(b[:n], data) ***REMOVED***
			t.Fatalf("#%d: got %v; want %v", i, b[:n], data)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCodec_ReceiveLimited(t *testing.T) ***REMOVED***
	const limit = 2048
	var payloads [][]byte
	for _, size := range []int***REMOVED***
		1024,
		2048,
		4096, // receive of this message would be interrupted due to limit
		2048, // this one is to make sure next receive recovers discarding leftovers
	***REMOVED*** ***REMOVED***
		b := make([]byte, size)
		rand.Read(b)
		payloads = append(payloads, b)
	***REMOVED***
	handlerDone := make(chan struct***REMOVED******REMOVED***)
	limitedHandler := func(ws *Conn) ***REMOVED***
		defer close(handlerDone)
		ws.MaxPayloadBytes = limit
		defer ws.Close()
		for i, p := range payloads ***REMOVED***
			t.Logf("payload #%d (size %d, exceeds limit: %v)", i, len(p), len(p) > limit)
			var recv []byte
			err := Message.Receive(ws, &recv)
			switch err ***REMOVED***
			case nil:
			case ErrFrameTooLarge:
				if len(p) <= limit ***REMOVED***
					t.Fatalf("unexpected frame size limit: expected %d bytes of payload having limit at %d", len(p), limit)
				***REMOVED***
				continue
			default:
				t.Fatalf("unexpected error: %v (want either nil or ErrFrameTooLarge)", err)
			***REMOVED***
			if len(recv) > limit ***REMOVED***
				t.Fatalf("received %d bytes of payload having limit at %d", len(recv), limit)
			***REMOVED***
			if !bytes.Equal(p, recv) ***REMOVED***
				t.Fatalf("received payload differs:\ngot:\t%v\nwant:\t%v", recv, p)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	server := httptest.NewServer(Handler(limitedHandler))
	defer server.CloseClientConnections()
	defer server.Close()
	addr := server.Listener.Addr().String()
	ws, err := Dial("ws://"+addr+"/", "", "http://localhost/")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer ws.Close()
	for i, p := range payloads ***REMOVED***
		if err := Message.Send(ws, p); err != nil ***REMOVED***
			t.Fatalf("payload #%d (size %d): %v", i, len(p), err)
		***REMOVED***
	***REMOVED***
	<-handlerDone
***REMOVED***
