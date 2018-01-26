// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

// This file implements a protocol of hybi draft.
// http://tools.ietf.org/html/draft-ietf-hybi-thewebsocketprotocol-17

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	websocketGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

	closeStatusNormal            = 1000
	closeStatusGoingAway         = 1001
	closeStatusProtocolError     = 1002
	closeStatusUnsupportedData   = 1003
	closeStatusFrameTooLarge     = 1004
	closeStatusNoStatusRcvd      = 1005
	closeStatusAbnormalClosure   = 1006
	closeStatusBadMessageData    = 1007
	closeStatusPolicyViolation   = 1008
	closeStatusTooBigData        = 1009
	closeStatusExtensionMismatch = 1010

	maxControlFramePayloadLength = 125
)

var (
	ErrBadMaskingKey         = &ProtocolError***REMOVED***"bad masking key"***REMOVED***
	ErrBadPongMessage        = &ProtocolError***REMOVED***"bad pong message"***REMOVED***
	ErrBadClosingStatus      = &ProtocolError***REMOVED***"bad closing status"***REMOVED***
	ErrUnsupportedExtensions = &ProtocolError***REMOVED***"unsupported extensions"***REMOVED***
	ErrNotImplemented        = &ProtocolError***REMOVED***"not implemented"***REMOVED***

	handshakeHeader = map[string]bool***REMOVED***
		"Host":                   true,
		"Upgrade":                true,
		"Connection":             true,
		"Sec-Websocket-Key":      true,
		"Sec-Websocket-Origin":   true,
		"Sec-Websocket-Version":  true,
		"Sec-Websocket-Protocol": true,
		"Sec-Websocket-Accept":   true,
	***REMOVED***
)

// A hybiFrameHeader is a frame header as defined in hybi draft.
type hybiFrameHeader struct ***REMOVED***
	Fin        bool
	Rsv        [3]bool
	OpCode     byte
	Length     int64
	MaskingKey []byte

	data *bytes.Buffer
***REMOVED***

// A hybiFrameReader is a reader for hybi frame.
type hybiFrameReader struct ***REMOVED***
	reader io.Reader

	header hybiFrameHeader
	pos    int64
	length int
***REMOVED***

func (frame *hybiFrameReader) Read(msg []byte) (n int, err error) ***REMOVED***
	n, err = frame.reader.Read(msg)
	if frame.header.MaskingKey != nil ***REMOVED***
		for i := 0; i < n; i++ ***REMOVED***
			msg[i] = msg[i] ^ frame.header.MaskingKey[frame.pos%4]
			frame.pos++
		***REMOVED***
	***REMOVED***
	return n, err
***REMOVED***

func (frame *hybiFrameReader) PayloadType() byte ***REMOVED*** return frame.header.OpCode ***REMOVED***

func (frame *hybiFrameReader) HeaderReader() io.Reader ***REMOVED***
	if frame.header.data == nil ***REMOVED***
		return nil
	***REMOVED***
	if frame.header.data.Len() == 0 ***REMOVED***
		return nil
	***REMOVED***
	return frame.header.data
***REMOVED***

func (frame *hybiFrameReader) TrailerReader() io.Reader ***REMOVED*** return nil ***REMOVED***

func (frame *hybiFrameReader) Len() (n int) ***REMOVED*** return frame.length ***REMOVED***

// A hybiFrameReaderFactory creates new frame reader based on its frame type.
type hybiFrameReaderFactory struct ***REMOVED***
	*bufio.Reader
***REMOVED***

// NewFrameReader reads a frame header from the connection, and creates new reader for the frame.
// See Section 5.2 Base Framing protocol for detail.
// http://tools.ietf.org/html/draft-ietf-hybi-thewebsocketprotocol-17#section-5.2
func (buf hybiFrameReaderFactory) NewFrameReader() (frame frameReader, err error) ***REMOVED***
	hybiFrame := new(hybiFrameReader)
	frame = hybiFrame
	var header []byte
	var b byte
	// First byte. FIN/RSV1/RSV2/RSV3/OpCode(4bits)
	b, err = buf.ReadByte()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	header = append(header, b)
	hybiFrame.header.Fin = ((header[0] >> 7) & 1) != 0
	for i := 0; i < 3; i++ ***REMOVED***
		j := uint(6 - i)
		hybiFrame.header.Rsv[i] = ((header[0] >> j) & 1) != 0
	***REMOVED***
	hybiFrame.header.OpCode = header[0] & 0x0f

	// Second byte. Mask/Payload len(7bits)
	b, err = buf.ReadByte()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	header = append(header, b)
	mask := (b & 0x80) != 0
	b &= 0x7f
	lengthFields := 0
	switch ***REMOVED***
	case b <= 125: // Payload length 7bits.
		hybiFrame.header.Length = int64(b)
	case b == 126: // Payload length 7+16bits
		lengthFields = 2
	case b == 127: // Payload length 7+64bits
		lengthFields = 8
	***REMOVED***
	for i := 0; i < lengthFields; i++ ***REMOVED***
		b, err = buf.ReadByte()
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if lengthFields == 8 && i == 0 ***REMOVED*** // MSB must be zero when 7+64 bits
			b &= 0x7f
		***REMOVED***
		header = append(header, b)
		hybiFrame.header.Length = hybiFrame.header.Length*256 + int64(b)
	***REMOVED***
	if mask ***REMOVED***
		// Masking key. 4 bytes.
		for i := 0; i < 4; i++ ***REMOVED***
			b, err = buf.ReadByte()
			if err != nil ***REMOVED***
				return
			***REMOVED***
			header = append(header, b)
			hybiFrame.header.MaskingKey = append(hybiFrame.header.MaskingKey, b)
		***REMOVED***
	***REMOVED***
	hybiFrame.reader = io.LimitReader(buf.Reader, hybiFrame.header.Length)
	hybiFrame.header.data = bytes.NewBuffer(header)
	hybiFrame.length = len(header) + int(hybiFrame.header.Length)
	return
***REMOVED***

// A HybiFrameWriter is a writer for hybi frame.
type hybiFrameWriter struct ***REMOVED***
	writer *bufio.Writer

	header *hybiFrameHeader
***REMOVED***

func (frame *hybiFrameWriter) Write(msg []byte) (n int, err error) ***REMOVED***
	var header []byte
	var b byte
	if frame.header.Fin ***REMOVED***
		b |= 0x80
	***REMOVED***
	for i := 0; i < 3; i++ ***REMOVED***
		if frame.header.Rsv[i] ***REMOVED***
			j := uint(6 - i)
			b |= 1 << j
		***REMOVED***
	***REMOVED***
	b |= frame.header.OpCode
	header = append(header, b)
	if frame.header.MaskingKey != nil ***REMOVED***
		b = 0x80
	***REMOVED*** else ***REMOVED***
		b = 0
	***REMOVED***
	lengthFields := 0
	length := len(msg)
	switch ***REMOVED***
	case length <= 125:
		b |= byte(length)
	case length < 65536:
		b |= 126
		lengthFields = 2
	default:
		b |= 127
		lengthFields = 8
	***REMOVED***
	header = append(header, b)
	for i := 0; i < lengthFields; i++ ***REMOVED***
		j := uint((lengthFields - i - 1) * 8)
		b = byte((length >> j) & 0xff)
		header = append(header, b)
	***REMOVED***
	if frame.header.MaskingKey != nil ***REMOVED***
		if len(frame.header.MaskingKey) != 4 ***REMOVED***
			return 0, ErrBadMaskingKey
		***REMOVED***
		header = append(header, frame.header.MaskingKey...)
		frame.writer.Write(header)
		data := make([]byte, length)
		for i := range data ***REMOVED***
			data[i] = msg[i] ^ frame.header.MaskingKey[i%4]
		***REMOVED***
		frame.writer.Write(data)
		err = frame.writer.Flush()
		return length, err
	***REMOVED***
	frame.writer.Write(header)
	frame.writer.Write(msg)
	err = frame.writer.Flush()
	return length, err
***REMOVED***

func (frame *hybiFrameWriter) Close() error ***REMOVED*** return nil ***REMOVED***

type hybiFrameWriterFactory struct ***REMOVED***
	*bufio.Writer
	needMaskingKey bool
***REMOVED***

func (buf hybiFrameWriterFactory) NewFrameWriter(payloadType byte) (frame frameWriter, err error) ***REMOVED***
	frameHeader := &hybiFrameHeader***REMOVED***Fin: true, OpCode: payloadType***REMOVED***
	if buf.needMaskingKey ***REMOVED***
		frameHeader.MaskingKey, err = generateMaskingKey()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return &hybiFrameWriter***REMOVED***writer: buf.Writer, header: frameHeader***REMOVED***, nil
***REMOVED***

type hybiFrameHandler struct ***REMOVED***
	conn        *Conn
	payloadType byte
***REMOVED***

func (handler *hybiFrameHandler) HandleFrame(frame frameReader) (frameReader, error) ***REMOVED***
	if handler.conn.IsServerConn() ***REMOVED***
		// The client MUST mask all frames sent to the server.
		if frame.(*hybiFrameReader).header.MaskingKey == nil ***REMOVED***
			handler.WriteClose(closeStatusProtocolError)
			return nil, io.EOF
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// The server MUST NOT mask all frames.
		if frame.(*hybiFrameReader).header.MaskingKey != nil ***REMOVED***
			handler.WriteClose(closeStatusProtocolError)
			return nil, io.EOF
		***REMOVED***
	***REMOVED***
	if header := frame.HeaderReader(); header != nil ***REMOVED***
		io.Copy(ioutil.Discard, header)
	***REMOVED***
	switch frame.PayloadType() ***REMOVED***
	case ContinuationFrame:
		frame.(*hybiFrameReader).header.OpCode = handler.payloadType
	case TextFrame, BinaryFrame:
		handler.payloadType = frame.PayloadType()
	case CloseFrame:
		return nil, io.EOF
	case PingFrame, PongFrame:
		b := make([]byte, maxControlFramePayloadLength)
		n, err := io.ReadFull(frame, b)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF ***REMOVED***
			return nil, err
		***REMOVED***
		io.Copy(ioutil.Discard, frame)
		if frame.PayloadType() == PingFrame ***REMOVED***
			if _, err := handler.WritePong(b[:n]); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		return nil, nil
	***REMOVED***
	return frame, nil
***REMOVED***

func (handler *hybiFrameHandler) WriteClose(status int) (err error) ***REMOVED***
	handler.conn.wio.Lock()
	defer handler.conn.wio.Unlock()
	w, err := handler.conn.frameWriterFactory.NewFrameWriter(CloseFrame)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	msg := make([]byte, 2)
	binary.BigEndian.PutUint16(msg, uint16(status))
	_, err = w.Write(msg)
	w.Close()
	return err
***REMOVED***

func (handler *hybiFrameHandler) WritePong(msg []byte) (n int, err error) ***REMOVED***
	handler.conn.wio.Lock()
	defer handler.conn.wio.Unlock()
	w, err := handler.conn.frameWriterFactory.NewFrameWriter(PongFrame)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	n, err = w.Write(msg)
	w.Close()
	return n, err
***REMOVED***

// newHybiConn creates a new WebSocket connection speaking hybi draft protocol.
func newHybiConn(config *Config, buf *bufio.ReadWriter, rwc io.ReadWriteCloser, request *http.Request) *Conn ***REMOVED***
	if buf == nil ***REMOVED***
		br := bufio.NewReader(rwc)
		bw := bufio.NewWriter(rwc)
		buf = bufio.NewReadWriter(br, bw)
	***REMOVED***
	ws := &Conn***REMOVED***config: config, request: request, buf: buf, rwc: rwc,
		frameReaderFactory: hybiFrameReaderFactory***REMOVED***buf.Reader***REMOVED***,
		frameWriterFactory: hybiFrameWriterFactory***REMOVED***
			buf.Writer, request == nil***REMOVED***,
		PayloadType:        TextFrame,
		defaultCloseStatus: closeStatusNormal***REMOVED***
	ws.frameHandler = &hybiFrameHandler***REMOVED***conn: ws***REMOVED***
	return ws
***REMOVED***

// generateMaskingKey generates a masking key for a frame.
func generateMaskingKey() (maskingKey []byte, err error) ***REMOVED***
	maskingKey = make([]byte, 4)
	if _, err = io.ReadFull(rand.Reader, maskingKey); err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// generateNonce generates a nonce consisting of a randomly selected 16-byte
// value that has been base64-encoded.
func generateNonce() (nonce []byte) ***REMOVED***
	key := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, key); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	nonce = make([]byte, 24)
	base64.StdEncoding.Encode(nonce, key)
	return
***REMOVED***

// removeZone removes IPv6 zone identifer from host.
// E.g., "[fe80::1%en0]:8080" to "[fe80::1]:8080"
func removeZone(host string) string ***REMOVED***
	if !strings.HasPrefix(host, "[") ***REMOVED***
		return host
	***REMOVED***
	i := strings.LastIndex(host, "]")
	if i < 0 ***REMOVED***
		return host
	***REMOVED***
	j := strings.LastIndex(host[:i], "%")
	if j < 0 ***REMOVED***
		return host
	***REMOVED***
	return host[:j] + host[i:]
***REMOVED***

// getNonceAccept computes the base64-encoded SHA-1 of the concatenation of
// the nonce ("Sec-WebSocket-Key" value) with the websocket GUID string.
func getNonceAccept(nonce []byte) (expected []byte, err error) ***REMOVED***
	h := sha1.New()
	if _, err = h.Write(nonce); err != nil ***REMOVED***
		return
	***REMOVED***
	if _, err = h.Write([]byte(websocketGUID)); err != nil ***REMOVED***
		return
	***REMOVED***
	expected = make([]byte, 28)
	base64.StdEncoding.Encode(expected, h.Sum(nil))
	return
***REMOVED***

// Client handshake described in draft-ietf-hybi-thewebsocket-protocol-17
func hybiClientHandshake(config *Config, br *bufio.Reader, bw *bufio.Writer) (err error) ***REMOVED***
	bw.WriteString("GET " + config.Location.RequestURI() + " HTTP/1.1\r\n")

	// According to RFC 6874, an HTTP client, proxy, or other
	// intermediary must remove any IPv6 zone identifier attached
	// to an outgoing URI.
	bw.WriteString("Host: " + removeZone(config.Location.Host) + "\r\n")
	bw.WriteString("Upgrade: websocket\r\n")
	bw.WriteString("Connection: Upgrade\r\n")
	nonce := generateNonce()
	if config.handshakeData != nil ***REMOVED***
		nonce = []byte(config.handshakeData["key"])
	***REMOVED***
	bw.WriteString("Sec-WebSocket-Key: " + string(nonce) + "\r\n")
	bw.WriteString("Origin: " + strings.ToLower(config.Origin.String()) + "\r\n")

	if config.Version != ProtocolVersionHybi13 ***REMOVED***
		return ErrBadProtocolVersion
	***REMOVED***

	bw.WriteString("Sec-WebSocket-Version: " + fmt.Sprintf("%d", config.Version) + "\r\n")
	if len(config.Protocol) > 0 ***REMOVED***
		bw.WriteString("Sec-WebSocket-Protocol: " + strings.Join(config.Protocol, ", ") + "\r\n")
	***REMOVED***
	// TODO(ukai): send Sec-WebSocket-Extensions.
	err = config.Header.WriteSubset(bw, handshakeHeader)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	bw.WriteString("\r\n")
	if err = bw.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***

	resp, err := http.ReadResponse(br, &http.Request***REMOVED***Method: "GET"***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if resp.StatusCode != 101 ***REMOVED***
		return ErrBadStatus
	***REMOVED***
	if strings.ToLower(resp.Header.Get("Upgrade")) != "websocket" ||
		strings.ToLower(resp.Header.Get("Connection")) != "upgrade" ***REMOVED***
		return ErrBadUpgrade
	***REMOVED***
	expectedAccept, err := getNonceAccept(nonce)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if resp.Header.Get("Sec-WebSocket-Accept") != string(expectedAccept) ***REMOVED***
		return ErrChallengeResponse
	***REMOVED***
	if resp.Header.Get("Sec-WebSocket-Extensions") != "" ***REMOVED***
		return ErrUnsupportedExtensions
	***REMOVED***
	offeredProtocol := resp.Header.Get("Sec-WebSocket-Protocol")
	if offeredProtocol != "" ***REMOVED***
		protocolMatched := false
		for i := 0; i < len(config.Protocol); i++ ***REMOVED***
			if config.Protocol[i] == offeredProtocol ***REMOVED***
				protocolMatched = true
				break
			***REMOVED***
		***REMOVED***
		if !protocolMatched ***REMOVED***
			return ErrBadWebSocketProtocol
		***REMOVED***
		config.Protocol = []string***REMOVED***offeredProtocol***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// newHybiClientConn creates a client WebSocket connection after handshake.
func newHybiClientConn(config *Config, buf *bufio.ReadWriter, rwc io.ReadWriteCloser) *Conn ***REMOVED***
	return newHybiConn(config, buf, rwc, nil)
***REMOVED***

// A HybiServerHandshaker performs a server handshake using hybi draft protocol.
type hybiServerHandshaker struct ***REMOVED***
	*Config
	accept []byte
***REMOVED***

func (c *hybiServerHandshaker) ReadHandshake(buf *bufio.Reader, req *http.Request) (code int, err error) ***REMOVED***
	c.Version = ProtocolVersionHybi13
	if req.Method != "GET" ***REMOVED***
		return http.StatusMethodNotAllowed, ErrBadRequestMethod
	***REMOVED***
	// HTTP version can be safely ignored.

	if strings.ToLower(req.Header.Get("Upgrade")) != "websocket" ||
		!strings.Contains(strings.ToLower(req.Header.Get("Connection")), "upgrade") ***REMOVED***
		return http.StatusBadRequest, ErrNotWebSocket
	***REMOVED***

	key := req.Header.Get("Sec-Websocket-Key")
	if key == "" ***REMOVED***
		return http.StatusBadRequest, ErrChallengeResponse
	***REMOVED***
	version := req.Header.Get("Sec-Websocket-Version")
	switch version ***REMOVED***
	case "13":
		c.Version = ProtocolVersionHybi13
	default:
		return http.StatusBadRequest, ErrBadWebSocketVersion
	***REMOVED***
	var scheme string
	if req.TLS != nil ***REMOVED***
		scheme = "wss"
	***REMOVED*** else ***REMOVED***
		scheme = "ws"
	***REMOVED***
	c.Location, err = url.ParseRequestURI(scheme + "://" + req.Host + req.URL.RequestURI())
	if err != nil ***REMOVED***
		return http.StatusBadRequest, err
	***REMOVED***
	protocol := strings.TrimSpace(req.Header.Get("Sec-Websocket-Protocol"))
	if protocol != "" ***REMOVED***
		protocols := strings.Split(protocol, ",")
		for i := 0; i < len(protocols); i++ ***REMOVED***
			c.Protocol = append(c.Protocol, strings.TrimSpace(protocols[i]))
		***REMOVED***
	***REMOVED***
	c.accept, err = getNonceAccept([]byte(key))
	if err != nil ***REMOVED***
		return http.StatusInternalServerError, err
	***REMOVED***
	return http.StatusSwitchingProtocols, nil
***REMOVED***

// Origin parses the Origin header in req.
// If the Origin header is not set, it returns nil and nil.
func Origin(config *Config, req *http.Request) (*url.URL, error) ***REMOVED***
	var origin string
	switch config.Version ***REMOVED***
	case ProtocolVersionHybi13:
		origin = req.Header.Get("Origin")
	***REMOVED***
	if origin == "" ***REMOVED***
		return nil, nil
	***REMOVED***
	return url.ParseRequestURI(origin)
***REMOVED***

func (c *hybiServerHandshaker) AcceptHandshake(buf *bufio.Writer) (err error) ***REMOVED***
	if len(c.Protocol) > 0 ***REMOVED***
		if len(c.Protocol) != 1 ***REMOVED***
			// You need choose a Protocol in Handshake func in Server.
			return ErrBadWebSocketProtocol
		***REMOVED***
	***REMOVED***
	buf.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	buf.WriteString("Upgrade: websocket\r\n")
	buf.WriteString("Connection: Upgrade\r\n")
	buf.WriteString("Sec-WebSocket-Accept: " + string(c.accept) + "\r\n")
	if len(c.Protocol) > 0 ***REMOVED***
		buf.WriteString("Sec-WebSocket-Protocol: " + c.Protocol[0] + "\r\n")
	***REMOVED***
	// TODO(ukai): send Sec-WebSocket-Extensions.
	if c.Header != nil ***REMOVED***
		err := c.Header.WriteSubset(buf, handshakeHeader)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	buf.WriteString("\r\n")
	return buf.Flush()
***REMOVED***

func (c *hybiServerHandshaker) NewServerConn(buf *bufio.ReadWriter, rwc io.ReadWriteCloser, request *http.Request) *Conn ***REMOVED***
	return newHybiServerConn(c.Config, buf, rwc, request)
***REMOVED***

// newHybiServerConn returns a new WebSocket connection speaking hybi draft protocol.
func newHybiServerConn(config *Config, buf *bufio.ReadWriter, rwc io.ReadWriteCloser, request *http.Request) *Conn ***REMOVED***
	return newHybiConn(config, buf, rwc, request)
***REMOVED***
