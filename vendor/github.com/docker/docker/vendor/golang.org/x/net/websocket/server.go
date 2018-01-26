// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
)

func newServerConn(rwc io.ReadWriteCloser, buf *bufio.ReadWriter, req *http.Request, config *Config, handshake func(*Config, *http.Request) error) (conn *Conn, err error) ***REMOVED***
	var hs serverHandshaker = &hybiServerHandshaker***REMOVED***Config: config***REMOVED***
	code, err := hs.ReadHandshake(buf.Reader, req)
	if err == ErrBadWebSocketVersion ***REMOVED***
		fmt.Fprintf(buf, "HTTP/1.1 %03d %s\r\n", code, http.StatusText(code))
		fmt.Fprintf(buf, "Sec-WebSocket-Version: %s\r\n", SupportedProtocolVersion)
		buf.WriteString("\r\n")
		buf.WriteString(err.Error())
		buf.Flush()
		return
	***REMOVED***
	if err != nil ***REMOVED***
		fmt.Fprintf(buf, "HTTP/1.1 %03d %s\r\n", code, http.StatusText(code))
		buf.WriteString("\r\n")
		buf.WriteString(err.Error())
		buf.Flush()
		return
	***REMOVED***
	if handshake != nil ***REMOVED***
		err = handshake(config, req)
		if err != nil ***REMOVED***
			code = http.StatusForbidden
			fmt.Fprintf(buf, "HTTP/1.1 %03d %s\r\n", code, http.StatusText(code))
			buf.WriteString("\r\n")
			buf.Flush()
			return
		***REMOVED***
	***REMOVED***
	err = hs.AcceptHandshake(buf.Writer)
	if err != nil ***REMOVED***
		code = http.StatusBadRequest
		fmt.Fprintf(buf, "HTTP/1.1 %03d %s\r\n", code, http.StatusText(code))
		buf.WriteString("\r\n")
		buf.Flush()
		return
	***REMOVED***
	conn = hs.NewServerConn(buf, rwc, req)
	return
***REMOVED***

// Server represents a server of a WebSocket.
type Server struct ***REMOVED***
	// Config is a WebSocket configuration for new WebSocket connection.
	Config

	// Handshake is an optional function in WebSocket handshake.
	// For example, you can check, or don't check Origin header.
	// Another example, you can select config.Protocol.
	Handshake func(*Config, *http.Request) error

	// Handler handles a WebSocket connection.
	Handler
***REMOVED***

// ServeHTTP implements the http.Handler interface for a WebSocket
func (s Server) ServeHTTP(w http.ResponseWriter, req *http.Request) ***REMOVED***
	s.serveWebSocket(w, req)
***REMOVED***

func (s Server) serveWebSocket(w http.ResponseWriter, req *http.Request) ***REMOVED***
	rwc, buf, err := w.(http.Hijacker).Hijack()
	if err != nil ***REMOVED***
		panic("Hijack failed: " + err.Error())
	***REMOVED***
	// The server should abort the WebSocket connection if it finds
	// the client did not send a handshake that matches with protocol
	// specification.
	defer rwc.Close()
	conn, err := newServerConn(rwc, buf, req, &s.Config, s.Handshake)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if conn == nil ***REMOVED***
		panic("unexpected nil conn")
	***REMOVED***
	s.Handler(conn)
***REMOVED***

// Handler is a simple interface to a WebSocket browser client.
// It checks if Origin header is valid URL by default.
// You might want to verify websocket.Conn.Config().Origin in the func.
// If you use Server instead of Handler, you could call websocket.Origin and
// check the origin in your Handshake func. So, if you want to accept
// non-browser clients, which do not send an Origin header, set a
// Server.Handshake that does not check the origin.
type Handler func(*Conn)

func checkOrigin(config *Config, req *http.Request) (err error) ***REMOVED***
	config.Origin, err = Origin(config, req)
	if err == nil && config.Origin == nil ***REMOVED***
		return fmt.Errorf("null origin")
	***REMOVED***
	return err
***REMOVED***

// ServeHTTP implements the http.Handler interface for a WebSocket
func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) ***REMOVED***
	s := Server***REMOVED***Handler: h, Handshake: checkOrigin***REMOVED***
	s.serveWebSocket(w, req)
***REMOVED***
