// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// Test the getNonceAccept function with values in
// http://tools.ietf.org/html/draft-ietf-hybi-thewebsocketprotocol-17
func TestSecWebSocketAccept(t *testing.T) ***REMOVED***
	nonce := []byte("dGhlIHNhbXBsZSBub25jZQ==")
	expected := []byte("s3pPLMBiTxaQ9kYGzzhZRbK+xOo=")
	accept, err := getNonceAccept(nonce)
	if err != nil ***REMOVED***
		t.Errorf("getNonceAccept: returned error %v", err)
		return
	***REMOVED***
	if !bytes.Equal(expected, accept) ***REMOVED***
		t.Errorf("getNonceAccept: expected %q got %q", expected, accept)
	***REMOVED***
***REMOVED***

func TestHybiClientHandshake(t *testing.T) ***REMOVED***
	type test struct ***REMOVED***
		url, host string
	***REMOVED***
	tests := []test***REMOVED***
		***REMOVED***"ws://server.example.com/chat", "server.example.com"***REMOVED***,
		***REMOVED***"ws://127.0.0.1/chat", "127.0.0.1"***REMOVED***,
	***REMOVED***
	if _, err := url.ParseRequestURI("http://[fe80::1%25lo0]"); err == nil ***REMOVED***
		tests = append(tests, test***REMOVED***"ws://[fe80::1%25lo0]/chat", "[fe80::1]"***REMOVED***)
	***REMOVED***

	for _, tt := range tests ***REMOVED***
		var b bytes.Buffer
		bw := bufio.NewWriter(&b)
		br := bufio.NewReader(strings.NewReader(`HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
Sec-WebSocket-Protocol: chat

`))
		var err error
		var config Config
		config.Location, err = url.ParseRequestURI(tt.url)
		if err != nil ***REMOVED***
			t.Fatal("location url", err)
		***REMOVED***
		config.Origin, err = url.ParseRequestURI("http://example.com")
		if err != nil ***REMOVED***
			t.Fatal("origin url", err)
		***REMOVED***
		config.Protocol = append(config.Protocol, "chat")
		config.Protocol = append(config.Protocol, "superchat")
		config.Version = ProtocolVersionHybi13
		config.handshakeData = map[string]string***REMOVED***
			"key": "dGhlIHNhbXBsZSBub25jZQ==",
		***REMOVED***
		if err := hybiClientHandshake(&config, br, bw); err != nil ***REMOVED***
			t.Fatal("handshake", err)
		***REMOVED***
		req, err := http.ReadRequest(bufio.NewReader(&b))
		if err != nil ***REMOVED***
			t.Fatal("read request", err)
		***REMOVED***
		if req.Method != "GET" ***REMOVED***
			t.Errorf("request method expected GET, but got %s", req.Method)
		***REMOVED***
		if req.URL.Path != "/chat" ***REMOVED***
			t.Errorf("request path expected /chat, but got %s", req.URL.Path)
		***REMOVED***
		if req.Proto != "HTTP/1.1" ***REMOVED***
			t.Errorf("request proto expected HTTP/1.1, but got %s", req.Proto)
		***REMOVED***
		if req.Host != tt.host ***REMOVED***
			t.Errorf("request host expected %s, but got %s", tt.host, req.Host)
		***REMOVED***
		var expectedHeader = map[string]string***REMOVED***
			"Connection":             "Upgrade",
			"Upgrade":                "websocket",
			"Sec-Websocket-Key":      config.handshakeData["key"],
			"Origin":                 config.Origin.String(),
			"Sec-Websocket-Protocol": "chat, superchat",
			"Sec-Websocket-Version":  fmt.Sprintf("%d", ProtocolVersionHybi13),
		***REMOVED***
		for k, v := range expectedHeader ***REMOVED***
			if req.Header.Get(k) != v ***REMOVED***
				t.Errorf("%s expected %s, but got %v", k, v, req.Header.Get(k))
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestHybiClientHandshakeWithHeader(t *testing.T) ***REMOVED***
	b := bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	bw := bufio.NewWriter(b)
	br := bufio.NewReader(strings.NewReader(`HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
Sec-WebSocket-Protocol: chat

`))
	var err error
	config := new(Config)
	config.Location, err = url.ParseRequestURI("ws://server.example.com/chat")
	if err != nil ***REMOVED***
		t.Fatal("location url", err)
	***REMOVED***
	config.Origin, err = url.ParseRequestURI("http://example.com")
	if err != nil ***REMOVED***
		t.Fatal("origin url", err)
	***REMOVED***
	config.Protocol = append(config.Protocol, "chat")
	config.Protocol = append(config.Protocol, "superchat")
	config.Version = ProtocolVersionHybi13
	config.Header = http.Header(make(map[string][]string))
	config.Header.Add("User-Agent", "test")

	config.handshakeData = map[string]string***REMOVED***
		"key": "dGhlIHNhbXBsZSBub25jZQ==",
	***REMOVED***
	err = hybiClientHandshake(config, br, bw)
	if err != nil ***REMOVED***
		t.Errorf("handshake failed: %v", err)
	***REMOVED***
	req, err := http.ReadRequest(bufio.NewReader(b))
	if err != nil ***REMOVED***
		t.Fatalf("read request: %v", err)
	***REMOVED***
	if req.Method != "GET" ***REMOVED***
		t.Errorf("request method expected GET, but got %q", req.Method)
	***REMOVED***
	if req.URL.Path != "/chat" ***REMOVED***
		t.Errorf("request path expected /chat, but got %q", req.URL.Path)
	***REMOVED***
	if req.Proto != "HTTP/1.1" ***REMOVED***
		t.Errorf("request proto expected HTTP/1.1, but got %q", req.Proto)
	***REMOVED***
	if req.Host != "server.example.com" ***REMOVED***
		t.Errorf("request Host expected server.example.com, but got %v", req.Host)
	***REMOVED***
	var expectedHeader = map[string]string***REMOVED***
		"Connection":             "Upgrade",
		"Upgrade":                "websocket",
		"Sec-Websocket-Key":      config.handshakeData["key"],
		"Origin":                 config.Origin.String(),
		"Sec-Websocket-Protocol": "chat, superchat",
		"Sec-Websocket-Version":  fmt.Sprintf("%d", ProtocolVersionHybi13),
		"User-Agent":             "test",
	***REMOVED***
	for k, v := range expectedHeader ***REMOVED***
		if req.Header.Get(k) != v ***REMOVED***
			t.Errorf(fmt.Sprintf("%s expected %q but got %q", k, v, req.Header.Get(k)))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestHybiServerHandshake(t *testing.T) ***REMOVED***
	config := new(Config)
	handshaker := &hybiServerHandshaker***REMOVED***Config: config***REMOVED***
	br := bufio.NewReader(strings.NewReader(`GET /chat HTTP/1.1
Host: server.example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Origin: http://example.com
Sec-WebSocket-Protocol: chat, superchat
Sec-WebSocket-Version: 13

`))
	req, err := http.ReadRequest(br)
	if err != nil ***REMOVED***
		t.Fatal("request", err)
	***REMOVED***
	code, err := handshaker.ReadHandshake(br, req)
	if err != nil ***REMOVED***
		t.Errorf("handshake failed: %v", err)
	***REMOVED***
	if code != http.StatusSwitchingProtocols ***REMOVED***
		t.Errorf("status expected %q but got %q", http.StatusSwitchingProtocols, code)
	***REMOVED***
	expectedProtocols := []string***REMOVED***"chat", "superchat"***REMOVED***
	if fmt.Sprintf("%v", config.Protocol) != fmt.Sprintf("%v", expectedProtocols) ***REMOVED***
		t.Errorf("protocol expected %q but got %q", expectedProtocols, config.Protocol)
	***REMOVED***
	b := bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	bw := bufio.NewWriter(b)

	config.Protocol = config.Protocol[:1]

	err = handshaker.AcceptHandshake(bw)
	if err != nil ***REMOVED***
		t.Errorf("handshake response failed: %v", err)
	***REMOVED***
	expectedResponse := strings.Join([]string***REMOVED***
		"HTTP/1.1 101 Switching Protocols",
		"Upgrade: websocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=",
		"Sec-WebSocket-Protocol: chat",
		"", ""***REMOVED***, "\r\n")

	if b.String() != expectedResponse ***REMOVED***
		t.Errorf("handshake expected %q but got %q", expectedResponse, b.String())
	***REMOVED***
***REMOVED***

func TestHybiServerHandshakeNoSubProtocol(t *testing.T) ***REMOVED***
	config := new(Config)
	handshaker := &hybiServerHandshaker***REMOVED***Config: config***REMOVED***
	br := bufio.NewReader(strings.NewReader(`GET /chat HTTP/1.1
Host: server.example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Origin: http://example.com
Sec-WebSocket-Version: 13

`))
	req, err := http.ReadRequest(br)
	if err != nil ***REMOVED***
		t.Fatal("request", err)
	***REMOVED***
	code, err := handshaker.ReadHandshake(br, req)
	if err != nil ***REMOVED***
		t.Errorf("handshake failed: %v", err)
	***REMOVED***
	if code != http.StatusSwitchingProtocols ***REMOVED***
		t.Errorf("status expected %q but got %q", http.StatusSwitchingProtocols, code)
	***REMOVED***
	if len(config.Protocol) != 0 ***REMOVED***
		t.Errorf("len(config.Protocol) expected 0, but got %q", len(config.Protocol))
	***REMOVED***
	b := bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	bw := bufio.NewWriter(b)

	err = handshaker.AcceptHandshake(bw)
	if err != nil ***REMOVED***
		t.Errorf("handshake response failed: %v", err)
	***REMOVED***
	expectedResponse := strings.Join([]string***REMOVED***
		"HTTP/1.1 101 Switching Protocols",
		"Upgrade: websocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=",
		"", ""***REMOVED***, "\r\n")

	if b.String() != expectedResponse ***REMOVED***
		t.Errorf("handshake expected %q but got %q", expectedResponse, b.String())
	***REMOVED***
***REMOVED***

func TestHybiServerHandshakeHybiBadVersion(t *testing.T) ***REMOVED***
	config := new(Config)
	handshaker := &hybiServerHandshaker***REMOVED***Config: config***REMOVED***
	br := bufio.NewReader(strings.NewReader(`GET /chat HTTP/1.1
Host: server.example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Sec-WebSocket-Origin: http://example.com
Sec-WebSocket-Protocol: chat, superchat
Sec-WebSocket-Version: 9

`))
	req, err := http.ReadRequest(br)
	if err != nil ***REMOVED***
		t.Fatal("request", err)
	***REMOVED***
	code, err := handshaker.ReadHandshake(br, req)
	if err != ErrBadWebSocketVersion ***REMOVED***
		t.Errorf("handshake expected err %q but got %q", ErrBadWebSocketVersion, err)
	***REMOVED***
	if code != http.StatusBadRequest ***REMOVED***
		t.Errorf("status expected %q but got %q", http.StatusBadRequest, code)
	***REMOVED***
***REMOVED***

func testHybiFrame(t *testing.T, testHeader, testPayload, testMaskedPayload []byte, frameHeader *hybiFrameHeader) ***REMOVED***
	b := bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	frameWriterFactory := &hybiFrameWriterFactory***REMOVED***bufio.NewWriter(b), false***REMOVED***
	w, _ := frameWriterFactory.NewFrameWriter(TextFrame)
	w.(*hybiFrameWriter).header = frameHeader
	_, err := w.Write(testPayload)
	w.Close()
	if err != nil ***REMOVED***
		t.Errorf("Write error %q", err)
	***REMOVED***
	var expectedFrame []byte
	expectedFrame = append(expectedFrame, testHeader...)
	expectedFrame = append(expectedFrame, testMaskedPayload...)
	if !bytes.Equal(expectedFrame, b.Bytes()) ***REMOVED***
		t.Errorf("frame expected %q got %q", expectedFrame, b.Bytes())
	***REMOVED***
	frameReaderFactory := &hybiFrameReaderFactory***REMOVED***bufio.NewReader(b)***REMOVED***
	r, err := frameReaderFactory.NewFrameReader()
	if err != nil ***REMOVED***
		t.Errorf("Read error %q", err)
	***REMOVED***
	if header := r.HeaderReader(); header == nil ***REMOVED***
		t.Errorf("no header")
	***REMOVED*** else ***REMOVED***
		actualHeader := make([]byte, r.Len())
		n, err := header.Read(actualHeader)
		if err != nil ***REMOVED***
			t.Errorf("Read header error %q", err)
		***REMOVED*** else ***REMOVED***
			if n < len(testHeader) ***REMOVED***
				t.Errorf("header too short %q got %q", testHeader, actualHeader[:n])
			***REMOVED***
			if !bytes.Equal(testHeader, actualHeader[:n]) ***REMOVED***
				t.Errorf("header expected %q got %q", testHeader, actualHeader[:n])
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if trailer := r.TrailerReader(); trailer != nil ***REMOVED***
		t.Errorf("unexpected trailer %q", trailer)
	***REMOVED***
	frame := r.(*hybiFrameReader)
	if frameHeader.Fin != frame.header.Fin ||
		frameHeader.OpCode != frame.header.OpCode ||
		len(testPayload) != int(frame.header.Length) ***REMOVED***
		t.Errorf("mismatch %v (%d) vs %v", frameHeader, len(testPayload), frame)
	***REMOVED***
	payload := make([]byte, len(testPayload))
	_, err = r.Read(payload)
	if err != nil && err != io.EOF ***REMOVED***
		t.Errorf("read %v", err)
	***REMOVED***
	if !bytes.Equal(testPayload, payload) ***REMOVED***
		t.Errorf("payload %q vs %q", testPayload, payload)
	***REMOVED***
***REMOVED***

func TestHybiShortTextFrame(t *testing.T) ***REMOVED***
	frameHeader := &hybiFrameHeader***REMOVED***Fin: true, OpCode: TextFrame***REMOVED***
	payload := []byte("hello")
	testHybiFrame(t, []byte***REMOVED***0x81, 0x05***REMOVED***, payload, payload, frameHeader)

	payload = make([]byte, 125)
	testHybiFrame(t, []byte***REMOVED***0x81, 125***REMOVED***, payload, payload, frameHeader)
***REMOVED***

func TestHybiShortMaskedTextFrame(t *testing.T) ***REMOVED***
	frameHeader := &hybiFrameHeader***REMOVED***Fin: true, OpCode: TextFrame,
		MaskingKey: []byte***REMOVED***0xcc, 0x55, 0x80, 0x20***REMOVED******REMOVED***
	payload := []byte("hello")
	maskedPayload := []byte***REMOVED***0xa4, 0x30, 0xec, 0x4c, 0xa3***REMOVED***
	header := []byte***REMOVED***0x81, 0x85***REMOVED***
	header = append(header, frameHeader.MaskingKey...)
	testHybiFrame(t, header, payload, maskedPayload, frameHeader)
***REMOVED***

func TestHybiShortBinaryFrame(t *testing.T) ***REMOVED***
	frameHeader := &hybiFrameHeader***REMOVED***Fin: true, OpCode: BinaryFrame***REMOVED***
	payload := []byte("hello")
	testHybiFrame(t, []byte***REMOVED***0x82, 0x05***REMOVED***, payload, payload, frameHeader)

	payload = make([]byte, 125)
	testHybiFrame(t, []byte***REMOVED***0x82, 125***REMOVED***, payload, payload, frameHeader)
***REMOVED***

func TestHybiControlFrame(t *testing.T) ***REMOVED***
	payload := []byte("hello")

	frameHeader := &hybiFrameHeader***REMOVED***Fin: true, OpCode: PingFrame***REMOVED***
	testHybiFrame(t, []byte***REMOVED***0x89, 0x05***REMOVED***, payload, payload, frameHeader)

	frameHeader = &hybiFrameHeader***REMOVED***Fin: true, OpCode: PingFrame***REMOVED***
	testHybiFrame(t, []byte***REMOVED***0x89, 0x00***REMOVED***, nil, nil, frameHeader)

	frameHeader = &hybiFrameHeader***REMOVED***Fin: true, OpCode: PongFrame***REMOVED***
	testHybiFrame(t, []byte***REMOVED***0x8A, 0x05***REMOVED***, payload, payload, frameHeader)

	frameHeader = &hybiFrameHeader***REMOVED***Fin: true, OpCode: PongFrame***REMOVED***
	testHybiFrame(t, []byte***REMOVED***0x8A, 0x00***REMOVED***, nil, nil, frameHeader)

	frameHeader = &hybiFrameHeader***REMOVED***Fin: true, OpCode: CloseFrame***REMOVED***
	payload = []byte***REMOVED***0x03, 0xe8***REMOVED*** // 1000
	testHybiFrame(t, []byte***REMOVED***0x88, 0x02***REMOVED***, payload, payload, frameHeader)
***REMOVED***

func TestHybiLongFrame(t *testing.T) ***REMOVED***
	frameHeader := &hybiFrameHeader***REMOVED***Fin: true, OpCode: TextFrame***REMOVED***
	payload := make([]byte, 126)
	testHybiFrame(t, []byte***REMOVED***0x81, 126, 0x00, 126***REMOVED***, payload, payload, frameHeader)

	payload = make([]byte, 65535)
	testHybiFrame(t, []byte***REMOVED***0x81, 126, 0xff, 0xff***REMOVED***, payload, payload, frameHeader)

	payload = make([]byte, 65536)
	testHybiFrame(t, []byte***REMOVED***0x81, 127, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00***REMOVED***, payload, payload, frameHeader)
***REMOVED***

func TestHybiClientRead(t *testing.T) ***REMOVED***
	wireData := []byte***REMOVED***0x81, 0x05, 'h', 'e', 'l', 'l', 'o',
		0x89, 0x05, 'h', 'e', 'l', 'l', 'o', // ping
		0x81, 0x05, 'w', 'o', 'r', 'l', 'd'***REMOVED***
	br := bufio.NewReader(bytes.NewBuffer(wireData))
	bw := bufio.NewWriter(bytes.NewBuffer([]byte***REMOVED******REMOVED***))
	conn := newHybiConn(newConfig(t, "/"), bufio.NewReadWriter(br, bw), nil, nil)

	msg := make([]byte, 512)
	n, err := conn.Read(msg)
	if err != nil ***REMOVED***
		t.Errorf("read 1st frame, error %q", err)
	***REMOVED***
	if n != 5 ***REMOVED***
		t.Errorf("read 1st frame, expect 5, got %d", n)
	***REMOVED***
	if !bytes.Equal(wireData[2:7], msg[:n]) ***REMOVED***
		t.Errorf("read 1st frame %v, got %v", wireData[2:7], msg[:n])
	***REMOVED***
	n, err = conn.Read(msg)
	if err != nil ***REMOVED***
		t.Errorf("read 2nd frame, error %q", err)
	***REMOVED***
	if n != 5 ***REMOVED***
		t.Errorf("read 2nd frame, expect 5, got %d", n)
	***REMOVED***
	if !bytes.Equal(wireData[16:21], msg[:n]) ***REMOVED***
		t.Errorf("read 2nd frame %v, got %v", wireData[16:21], msg[:n])
	***REMOVED***
	n, err = conn.Read(msg)
	if err == nil ***REMOVED***
		t.Errorf("read not EOF")
	***REMOVED***
	if n != 0 ***REMOVED***
		t.Errorf("expect read 0, got %d", n)
	***REMOVED***
***REMOVED***

func TestHybiShortRead(t *testing.T) ***REMOVED***
	wireData := []byte***REMOVED***0x81, 0x05, 'h', 'e', 'l', 'l', 'o',
		0x89, 0x05, 'h', 'e', 'l', 'l', 'o', // ping
		0x81, 0x05, 'w', 'o', 'r', 'l', 'd'***REMOVED***
	br := bufio.NewReader(bytes.NewBuffer(wireData))
	bw := bufio.NewWriter(bytes.NewBuffer([]byte***REMOVED******REMOVED***))
	conn := newHybiConn(newConfig(t, "/"), bufio.NewReadWriter(br, bw), nil, nil)

	step := 0
	pos := 0
	expectedPos := []int***REMOVED***2, 5, 16, 19***REMOVED***
	expectedLen := []int***REMOVED***3, 2, 3, 2***REMOVED***
	for ***REMOVED***
		msg := make([]byte, 3)
		n, err := conn.Read(msg)
		if step >= len(expectedPos) ***REMOVED***
			if err == nil ***REMOVED***
				t.Errorf("read not EOF")
			***REMOVED***
			if n != 0 ***REMOVED***
				t.Errorf("expect read 0, got %d", n)
			***REMOVED***
			return
		***REMOVED***
		pos = expectedPos[step]
		endPos := pos + expectedLen[step]
		if err != nil ***REMOVED***
			t.Errorf("read from %d, got error %q", pos, err)
			return
		***REMOVED***
		if n != endPos-pos ***REMOVED***
			t.Errorf("read from %d, expect %d, got %d", pos, endPos-pos, n)
		***REMOVED***
		if !bytes.Equal(wireData[pos:endPos], msg[:n]) ***REMOVED***
			t.Errorf("read from %d, frame %v, got %v", pos, wireData[pos:endPos], msg[:n])
		***REMOVED***
		step++
	***REMOVED***
***REMOVED***

func TestHybiServerRead(t *testing.T) ***REMOVED***
	wireData := []byte***REMOVED***0x81, 0x85, 0xcc, 0x55, 0x80, 0x20,
		0xa4, 0x30, 0xec, 0x4c, 0xa3, // hello
		0x89, 0x85, 0xcc, 0x55, 0x80, 0x20,
		0xa4, 0x30, 0xec, 0x4c, 0xa3, // ping: hello
		0x81, 0x85, 0xed, 0x83, 0xb4, 0x24,
		0x9a, 0xec, 0xc6, 0x48, 0x89, // world
	***REMOVED***
	br := bufio.NewReader(bytes.NewBuffer(wireData))
	bw := bufio.NewWriter(bytes.NewBuffer([]byte***REMOVED******REMOVED***))
	conn := newHybiConn(newConfig(t, "/"), bufio.NewReadWriter(br, bw), nil, new(http.Request))

	expected := [][]byte***REMOVED***[]byte("hello"), []byte("world")***REMOVED***

	msg := make([]byte, 512)
	n, err := conn.Read(msg)
	if err != nil ***REMOVED***
		t.Errorf("read 1st frame, error %q", err)
	***REMOVED***
	if n != 5 ***REMOVED***
		t.Errorf("read 1st frame, expect 5, got %d", n)
	***REMOVED***
	if !bytes.Equal(expected[0], msg[:n]) ***REMOVED***
		t.Errorf("read 1st frame %q, got %q", expected[0], msg[:n])
	***REMOVED***

	n, err = conn.Read(msg)
	if err != nil ***REMOVED***
		t.Errorf("read 2nd frame, error %q", err)
	***REMOVED***
	if n != 5 ***REMOVED***
		t.Errorf("read 2nd frame, expect 5, got %d", n)
	***REMOVED***
	if !bytes.Equal(expected[1], msg[:n]) ***REMOVED***
		t.Errorf("read 2nd frame %q, got %q", expected[1], msg[:n])
	***REMOVED***

	n, err = conn.Read(msg)
	if err == nil ***REMOVED***
		t.Errorf("read not EOF")
	***REMOVED***
	if n != 0 ***REMOVED***
		t.Errorf("expect read 0, got %d", n)
	***REMOVED***
***REMOVED***

func TestHybiServerReadWithoutMasking(t *testing.T) ***REMOVED***
	wireData := []byte***REMOVED***0x81, 0x05, 'h', 'e', 'l', 'l', 'o'***REMOVED***
	br := bufio.NewReader(bytes.NewBuffer(wireData))
	bw := bufio.NewWriter(bytes.NewBuffer([]byte***REMOVED******REMOVED***))
	conn := newHybiConn(newConfig(t, "/"), bufio.NewReadWriter(br, bw), nil, new(http.Request))
	// server MUST close the connection upon receiving a non-masked frame.
	msg := make([]byte, 512)
	_, err := conn.Read(msg)
	if err != io.EOF ***REMOVED***
		t.Errorf("read 1st frame, expect %q, but got %q", io.EOF, err)
	***REMOVED***
***REMOVED***

func TestHybiClientReadWithMasking(t *testing.T) ***REMOVED***
	wireData := []byte***REMOVED***0x81, 0x85, 0xcc, 0x55, 0x80, 0x20,
		0xa4, 0x30, 0xec, 0x4c, 0xa3, // hello
	***REMOVED***
	br := bufio.NewReader(bytes.NewBuffer(wireData))
	bw := bufio.NewWriter(bytes.NewBuffer([]byte***REMOVED******REMOVED***))
	conn := newHybiConn(newConfig(t, "/"), bufio.NewReadWriter(br, bw), nil, nil)

	// client MUST close the connection upon receiving a masked frame.
	msg := make([]byte, 512)
	_, err := conn.Read(msg)
	if err != io.EOF ***REMOVED***
		t.Errorf("read 1st frame, expect %q, but got %q", io.EOF, err)
	***REMOVED***
***REMOVED***

// Test the hybiServerHandshaker supports firefox implementation and
// checks Connection request header include (but it's not necessary
// equal to) "upgrade"
func TestHybiServerFirefoxHandshake(t *testing.T) ***REMOVED***
	config := new(Config)
	handshaker := &hybiServerHandshaker***REMOVED***Config: config***REMOVED***
	br := bufio.NewReader(strings.NewReader(`GET /chat HTTP/1.1
Host: server.example.com
Upgrade: websocket
Connection: keep-alive, upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Origin: http://example.com
Sec-WebSocket-Protocol: chat, superchat
Sec-WebSocket-Version: 13

`))
	req, err := http.ReadRequest(br)
	if err != nil ***REMOVED***
		t.Fatal("request", err)
	***REMOVED***
	code, err := handshaker.ReadHandshake(br, req)
	if err != nil ***REMOVED***
		t.Errorf("handshake failed: %v", err)
	***REMOVED***
	if code != http.StatusSwitchingProtocols ***REMOVED***
		t.Errorf("status expected %q but got %q", http.StatusSwitchingProtocols, code)
	***REMOVED***
	b := bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	bw := bufio.NewWriter(b)

	config.Protocol = []string***REMOVED***"chat"***REMOVED***

	err = handshaker.AcceptHandshake(bw)
	if err != nil ***REMOVED***
		t.Errorf("handshake response failed: %v", err)
	***REMOVED***
	expectedResponse := strings.Join([]string***REMOVED***
		"HTTP/1.1 101 Switching Protocols",
		"Upgrade: websocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=",
		"Sec-WebSocket-Protocol: chat",
		"", ""***REMOVED***, "\r\n")

	if b.String() != expectedResponse ***REMOVED***
		t.Errorf("handshake expected %q but got %q", expectedResponse, b.String())
	***REMOVED***
***REMOVED***
