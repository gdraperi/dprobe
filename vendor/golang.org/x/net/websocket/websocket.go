// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package websocket implements a client and server for the WebSocket protocol
// as specified in RFC 6455.
//
// This package currently lacks some features found in an alternative
// and more actively maintained WebSocket package:
//
//     https://godoc.org/github.com/gorilla/websocket
//
package websocket // import "golang.org/x/net/websocket"

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	ProtocolVersionHybi13    = 13
	ProtocolVersionHybi      = ProtocolVersionHybi13
	SupportedProtocolVersion = "13"

	ContinuationFrame = 0
	TextFrame         = 1
	BinaryFrame       = 2
	CloseFrame        = 8
	PingFrame         = 9
	PongFrame         = 10
	UnknownFrame      = 255

	DefaultMaxPayloadBytes = 32 << 20 // 32MB
)

// ProtocolError represents WebSocket protocol errors.
type ProtocolError struct ***REMOVED***
	ErrorString string
***REMOVED***

func (err *ProtocolError) Error() string ***REMOVED*** return err.ErrorString ***REMOVED***

var (
	ErrBadProtocolVersion   = &ProtocolError***REMOVED***"bad protocol version"***REMOVED***
	ErrBadScheme            = &ProtocolError***REMOVED***"bad scheme"***REMOVED***
	ErrBadStatus            = &ProtocolError***REMOVED***"bad status"***REMOVED***
	ErrBadUpgrade           = &ProtocolError***REMOVED***"missing or bad upgrade"***REMOVED***
	ErrBadWebSocketOrigin   = &ProtocolError***REMOVED***"missing or bad WebSocket-Origin"***REMOVED***
	ErrBadWebSocketLocation = &ProtocolError***REMOVED***"missing or bad WebSocket-Location"***REMOVED***
	ErrBadWebSocketProtocol = &ProtocolError***REMOVED***"missing or bad WebSocket-Protocol"***REMOVED***
	ErrBadWebSocketVersion  = &ProtocolError***REMOVED***"missing or bad WebSocket Version"***REMOVED***
	ErrChallengeResponse    = &ProtocolError***REMOVED***"mismatch challenge/response"***REMOVED***
	ErrBadFrame             = &ProtocolError***REMOVED***"bad frame"***REMOVED***
	ErrBadFrameBoundary     = &ProtocolError***REMOVED***"not on frame boundary"***REMOVED***
	ErrNotWebSocket         = &ProtocolError***REMOVED***"not websocket protocol"***REMOVED***
	ErrBadRequestMethod     = &ProtocolError***REMOVED***"bad method"***REMOVED***
	ErrNotSupported         = &ProtocolError***REMOVED***"not supported"***REMOVED***
)

// ErrFrameTooLarge is returned by Codec's Receive method if payload size
// exceeds limit set by Conn.MaxPayloadBytes
var ErrFrameTooLarge = errors.New("websocket: frame payload size exceeds limit")

// Addr is an implementation of net.Addr for WebSocket.
type Addr struct ***REMOVED***
	*url.URL
***REMOVED***

// Network returns the network type for a WebSocket, "websocket".
func (addr *Addr) Network() string ***REMOVED*** return "websocket" ***REMOVED***

// Config is a WebSocket configuration
type Config struct ***REMOVED***
	// A WebSocket server address.
	Location *url.URL

	// A Websocket client origin.
	Origin *url.URL

	// WebSocket subprotocols.
	Protocol []string

	// WebSocket protocol version.
	Version int

	// TLS config for secure WebSocket (wss).
	TlsConfig *tls.Config

	// Additional header fields to be sent in WebSocket opening handshake.
	Header http.Header

	// Dialer used when opening websocket connections.
	Dialer *net.Dialer

	handshakeData map[string]string
***REMOVED***

// serverHandshaker is an interface to handle WebSocket server side handshake.
type serverHandshaker interface ***REMOVED***
	// ReadHandshake reads handshake request message from client.
	// Returns http response code and error if any.
	ReadHandshake(buf *bufio.Reader, req *http.Request) (code int, err error)

	// AcceptHandshake accepts the client handshake request and sends
	// handshake response back to client.
	AcceptHandshake(buf *bufio.Writer) (err error)

	// NewServerConn creates a new WebSocket connection.
	NewServerConn(buf *bufio.ReadWriter, rwc io.ReadWriteCloser, request *http.Request) (conn *Conn)
***REMOVED***

// frameReader is an interface to read a WebSocket frame.
type frameReader interface ***REMOVED***
	// Reader is to read payload of the frame.
	io.Reader

	// PayloadType returns payload type.
	PayloadType() byte

	// HeaderReader returns a reader to read header of the frame.
	HeaderReader() io.Reader

	// TrailerReader returns a reader to read trailer of the frame.
	// If it returns nil, there is no trailer in the frame.
	TrailerReader() io.Reader

	// Len returns total length of the frame, including header and trailer.
	Len() int
***REMOVED***

// frameReaderFactory is an interface to creates new frame reader.
type frameReaderFactory interface ***REMOVED***
	NewFrameReader() (r frameReader, err error)
***REMOVED***

// frameWriter is an interface to write a WebSocket frame.
type frameWriter interface ***REMOVED***
	// Writer is to write payload of the frame.
	io.WriteCloser
***REMOVED***

// frameWriterFactory is an interface to create new frame writer.
type frameWriterFactory interface ***REMOVED***
	NewFrameWriter(payloadType byte) (w frameWriter, err error)
***REMOVED***

type frameHandler interface ***REMOVED***
	HandleFrame(frame frameReader) (r frameReader, err error)
	WriteClose(status int) (err error)
***REMOVED***

// Conn represents a WebSocket connection.
//
// Multiple goroutines may invoke methods on a Conn simultaneously.
type Conn struct ***REMOVED***
	config  *Config
	request *http.Request

	buf *bufio.ReadWriter
	rwc io.ReadWriteCloser

	rio sync.Mutex
	frameReaderFactory
	frameReader

	wio sync.Mutex
	frameWriterFactory

	frameHandler
	PayloadType        byte
	defaultCloseStatus int

	// MaxPayloadBytes limits the size of frame payload received over Conn
	// by Codec's Receive method. If zero, DefaultMaxPayloadBytes is used.
	MaxPayloadBytes int
***REMOVED***

// Read implements the io.Reader interface:
// it reads data of a frame from the WebSocket connection.
// if msg is not large enough for the frame data, it fills the msg and next Read
// will read the rest of the frame data.
// it reads Text frame or Binary frame.
func (ws *Conn) Read(msg []byte) (n int, err error) ***REMOVED***
	ws.rio.Lock()
	defer ws.rio.Unlock()
again:
	if ws.frameReader == nil ***REMOVED***
		frame, err := ws.frameReaderFactory.NewFrameReader()
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		ws.frameReader, err = ws.frameHandler.HandleFrame(frame)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		if ws.frameReader == nil ***REMOVED***
			goto again
		***REMOVED***
	***REMOVED***
	n, err = ws.frameReader.Read(msg)
	if err == io.EOF ***REMOVED***
		if trailer := ws.frameReader.TrailerReader(); trailer != nil ***REMOVED***
			io.Copy(ioutil.Discard, trailer)
		***REMOVED***
		ws.frameReader = nil
		goto again
	***REMOVED***
	return n, err
***REMOVED***

// Write implements the io.Writer interface:
// it writes data as a frame to the WebSocket connection.
func (ws *Conn) Write(msg []byte) (n int, err error) ***REMOVED***
	ws.wio.Lock()
	defer ws.wio.Unlock()
	w, err := ws.frameWriterFactory.NewFrameWriter(ws.PayloadType)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	n, err = w.Write(msg)
	w.Close()
	return n, err
***REMOVED***

// Close implements the io.Closer interface.
func (ws *Conn) Close() error ***REMOVED***
	err := ws.frameHandler.WriteClose(ws.defaultCloseStatus)
	err1 := ws.rwc.Close()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return err1
***REMOVED***

func (ws *Conn) IsClientConn() bool ***REMOVED*** return ws.request == nil ***REMOVED***
func (ws *Conn) IsServerConn() bool ***REMOVED*** return ws.request != nil ***REMOVED***

// LocalAddr returns the WebSocket Origin for the connection for client, or
// the WebSocket location for server.
func (ws *Conn) LocalAddr() net.Addr ***REMOVED***
	if ws.IsClientConn() ***REMOVED***
		return &Addr***REMOVED***ws.config.Origin***REMOVED***
	***REMOVED***
	return &Addr***REMOVED***ws.config.Location***REMOVED***
***REMOVED***

// RemoteAddr returns the WebSocket location for the connection for client, or
// the Websocket Origin for server.
func (ws *Conn) RemoteAddr() net.Addr ***REMOVED***
	if ws.IsClientConn() ***REMOVED***
		return &Addr***REMOVED***ws.config.Location***REMOVED***
	***REMOVED***
	return &Addr***REMOVED***ws.config.Origin***REMOVED***
***REMOVED***

var errSetDeadline = errors.New("websocket: cannot set deadline: not using a net.Conn")

// SetDeadline sets the connection's network read & write deadlines.
func (ws *Conn) SetDeadline(t time.Time) error ***REMOVED***
	if conn, ok := ws.rwc.(net.Conn); ok ***REMOVED***
		return conn.SetDeadline(t)
	***REMOVED***
	return errSetDeadline
***REMOVED***

// SetReadDeadline sets the connection's network read deadline.
func (ws *Conn) SetReadDeadline(t time.Time) error ***REMOVED***
	if conn, ok := ws.rwc.(net.Conn); ok ***REMOVED***
		return conn.SetReadDeadline(t)
	***REMOVED***
	return errSetDeadline
***REMOVED***

// SetWriteDeadline sets the connection's network write deadline.
func (ws *Conn) SetWriteDeadline(t time.Time) error ***REMOVED***
	if conn, ok := ws.rwc.(net.Conn); ok ***REMOVED***
		return conn.SetWriteDeadline(t)
	***REMOVED***
	return errSetDeadline
***REMOVED***

// Config returns the WebSocket config.
func (ws *Conn) Config() *Config ***REMOVED*** return ws.config ***REMOVED***

// Request returns the http request upgraded to the WebSocket.
// It is nil for client side.
func (ws *Conn) Request() *http.Request ***REMOVED*** return ws.request ***REMOVED***

// Codec represents a symmetric pair of functions that implement a codec.
type Codec struct ***REMOVED***
	Marshal   func(v interface***REMOVED******REMOVED***) (data []byte, payloadType byte, err error)
	Unmarshal func(data []byte, payloadType byte, v interface***REMOVED******REMOVED***) (err error)
***REMOVED***

// Send sends v marshaled by cd.Marshal as single frame to ws.
func (cd Codec) Send(ws *Conn, v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	data, payloadType, err := cd.Marshal(v)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	ws.wio.Lock()
	defer ws.wio.Unlock()
	w, err := ws.frameWriterFactory.NewFrameWriter(payloadType)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = w.Write(data)
	w.Close()
	return err
***REMOVED***

// Receive receives single frame from ws, unmarshaled by cd.Unmarshal and stores
// in v. The whole frame payload is read to an in-memory buffer; max size of
// payload is defined by ws.MaxPayloadBytes. If frame payload size exceeds
// limit, ErrFrameTooLarge is returned; in this case frame is not read off wire
// completely. The next call to Receive would read and discard leftover data of
// previous oversized frame before processing next frame.
func (cd Codec) Receive(ws *Conn, v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	ws.rio.Lock()
	defer ws.rio.Unlock()
	if ws.frameReader != nil ***REMOVED***
		_, err = io.Copy(ioutil.Discard, ws.frameReader)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		ws.frameReader = nil
	***REMOVED***
again:
	frame, err := ws.frameReaderFactory.NewFrameReader()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	frame, err = ws.frameHandler.HandleFrame(frame)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if frame == nil ***REMOVED***
		goto again
	***REMOVED***
	maxPayloadBytes := ws.MaxPayloadBytes
	if maxPayloadBytes == 0 ***REMOVED***
		maxPayloadBytes = DefaultMaxPayloadBytes
	***REMOVED***
	if hf, ok := frame.(*hybiFrameReader); ok && hf.header.Length > int64(maxPayloadBytes) ***REMOVED***
		// payload size exceeds limit, no need to call Unmarshal
		//
		// set frameReader to current oversized frame so that
		// the next call to this function can drain leftover
		// data before processing the next frame
		ws.frameReader = frame
		return ErrFrameTooLarge
	***REMOVED***
	payloadType := frame.PayloadType()
	data, err := ioutil.ReadAll(frame)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return cd.Unmarshal(data, payloadType, v)
***REMOVED***

func marshal(v interface***REMOVED******REMOVED***) (msg []byte, payloadType byte, err error) ***REMOVED***
	switch data := v.(type) ***REMOVED***
	case string:
		return []byte(data), TextFrame, nil
	case []byte:
		return data, BinaryFrame, nil
	***REMOVED***
	return nil, UnknownFrame, ErrNotSupported
***REMOVED***

func unmarshal(msg []byte, payloadType byte, v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	switch data := v.(type) ***REMOVED***
	case *string:
		*data = string(msg)
		return nil
	case *[]byte:
		*data = msg
		return nil
	***REMOVED***
	return ErrNotSupported
***REMOVED***

/*
Message is a codec to send/receive text/binary data in a frame on WebSocket connection.
To send/receive text frame, use string type.
To send/receive binary frame, use []byte type.

Trivial usage:

	import "websocket"

	// receive text frame
	var message string
	websocket.Message.Receive(ws, &message)

	// send text frame
	message = "hello"
	websocket.Message.Send(ws, message)

	// receive binary frame
	var data []byte
	websocket.Message.Receive(ws, &data)

	// send binary frame
	data = []byte***REMOVED***0, 1, 2***REMOVED***
	websocket.Message.Send(ws, data)

*/
var Message = Codec***REMOVED***marshal, unmarshal***REMOVED***

func jsonMarshal(v interface***REMOVED******REMOVED***) (msg []byte, payloadType byte, err error) ***REMOVED***
	msg, err = json.Marshal(v)
	return msg, TextFrame, err
***REMOVED***

func jsonUnmarshal(msg []byte, payloadType byte, v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	return json.Unmarshal(msg, v)
***REMOVED***

/*
JSON is a codec to send/receive JSON data in a frame from a WebSocket connection.

Trivial usage:

	import "websocket"

	type T struct ***REMOVED***
		Msg string
		Count int
	***REMOVED***

	// receive JSON type T
	var data T
	websocket.JSON.Receive(ws, &data)

	// send JSON type T
	websocket.JSON.Send(ws, data)
*/
var JSON = Codec***REMOVED***jsonMarshal, jsonUnmarshal***REMOVED***
