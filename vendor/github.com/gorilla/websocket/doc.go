// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package websocket implements the WebSocket protocol defined in RFC 6455.
//
// Overview
//
// The Conn type represents a WebSocket connection. A server application calls
// the Upgrader.Upgrade method from an HTTP request handler to get a *Conn:
//
//  var upgrader = websocket.Upgrader***REMOVED***
//      ReadBufferSize:  1024,
//      WriteBufferSize: 1024,
//  ***REMOVED***
//
//  func handler(w http.ResponseWriter, r *http.Request) ***REMOVED***
//      conn, err := upgrader.Upgrade(w, r, nil)
//      if err != nil ***REMOVED***
//          log.Println(err)
//          return
//  ***REMOVED***
//      ... Use conn to send and receive messages.
//  ***REMOVED***
//
// Call the connection's WriteMessage and ReadMessage methods to send and
// receive messages as a slice of bytes. This snippet of code shows how to echo
// messages using these methods:
//
//  for ***REMOVED***
//      messageType, p, err := conn.ReadMessage()
//      if err != nil ***REMOVED***
//          log.Println(err)
//          return
//  ***REMOVED***
//      if err := conn.WriteMessage(messageType, p); err != nil ***REMOVED***
//          log.Println(err)
//          return
//  ***REMOVED***
//  ***REMOVED***
//
// In above snippet of code, p is a []byte and messageType is an int with value
// websocket.BinaryMessage or websocket.TextMessage.
//
// An application can also send and receive messages using the io.WriteCloser
// and io.Reader interfaces. To send a message, call the connection NextWriter
// method to get an io.WriteCloser, write the message to the writer and close
// the writer when done. To receive a message, call the connection NextReader
// method to get an io.Reader and read until io.EOF is returned. This snippet
// shows how to echo messages using the NextWriter and NextReader methods:
//
//  for ***REMOVED***
//      messageType, r, err := conn.NextReader()
//      if err != nil ***REMOVED***
//          return
//  ***REMOVED***
//      w, err := conn.NextWriter(messageType)
//      if err != nil ***REMOVED***
//          return err
//  ***REMOVED***
//      if _, err := io.Copy(w, r); err != nil ***REMOVED***
//          return err
//  ***REMOVED***
//      if err := w.Close(); err != nil ***REMOVED***
//          return err
//  ***REMOVED***
//  ***REMOVED***
//
// Data Messages
//
// The WebSocket protocol distinguishes between text and binary data messages.
// Text messages are interpreted as UTF-8 encoded text. The interpretation of
// binary messages is left to the application.
//
// This package uses the TextMessage and BinaryMessage integer constants to
// identify the two data message types. The ReadMessage and NextReader methods
// return the type of the received message. The messageType argument to the
// WriteMessage and NextWriter methods specifies the type of a sent message.
//
// It is the application's responsibility to ensure that text messages are
// valid UTF-8 encoded text.
//
// Control Messages
//
// The WebSocket protocol defines three types of control messages: close, ping
// and pong. Call the connection WriteControl, WriteMessage or NextWriter
// methods to send a control message to the peer.
//
// Connections handle received close messages by calling the handler function
// set with the SetCloseHandler method and by returning a *CloseError from the
// NextReader, ReadMessage or the message Read method. The default close
// handler sends a close message to the peer.
//
// Connections handle received ping messages by calling the handler function
// set with the SetPingHandler method. The default ping handler sends a pong
// message to the peer.
//
// Connections handle received pong messages by calling the handler function
// set with the SetPongHandler method. The default pong handler does nothing.
// If an application sends ping messages, then the application should set a
// pong handler to receive the corresponding pong.
//
// The control message handler functions are called from the NextReader,
// ReadMessage and message reader Read methods. The default close and ping
// handlers can block these methods for a short time when the handler writes to
// the connection.
//
// The application must read the connection to process close, ping and pong
// messages sent from the peer. If the application is not otherwise interested
// in messages from the peer, then the application should start a goroutine to
// read and discard messages from the peer. A simple example is:
//
//  func readLoop(c *websocket.Conn) ***REMOVED***
//      for ***REMOVED***
//          if _, _, err := c.NextReader(); err != nil ***REMOVED***
//              c.Close()
//              break
//      ***REMOVED***
//  ***REMOVED***
//  ***REMOVED***
//
// Concurrency
//
// Connections support one concurrent reader and one concurrent writer.
//
// Applications are responsible for ensuring that no more than one goroutine
// calls the write methods (NextWriter, SetWriteDeadline, WriteMessage,
// WriteJSON, EnableWriteCompression, SetCompressionLevel) concurrently and
// that no more than one goroutine calls the read methods (NextReader,
// SetReadDeadline, ReadMessage, ReadJSON, SetPongHandler, SetPingHandler)
// concurrently.
//
// The Close and WriteControl methods can be called concurrently with all other
// methods.
//
// Origin Considerations
//
// Web browsers allow Javascript applications to open a WebSocket connection to
// any host. It's up to the server to enforce an origin policy using the Origin
// request header sent by the browser.
//
// The Upgrader calls the function specified in the CheckOrigin field to check
// the origin. If the CheckOrigin function returns false, then the Upgrade
// method fails the WebSocket handshake with HTTP status 403.
//
// If the CheckOrigin field is nil, then the Upgrader uses a safe default: fail
// the handshake if the Origin request header is present and the Origin host is
// not equal to the Host request header.
//
// The deprecated package-level Upgrade function does not perform origin
// checking. The application is responsible for checking the Origin header
// before calling the Upgrade function.
//
// Compression EXPERIMENTAL
//
// Per message compression extensions (RFC 7692) are experimentally supported
// by this package in a limited capacity. Setting the EnableCompression option
// to true in Dialer or Upgrader will attempt to negotiate per message deflate
// support.
//
//  var upgrader = websocket.Upgrader***REMOVED***
//      EnableCompression: true,
//  ***REMOVED***
//
// If compression was successfully negotiated with the connection's peer, any
// message received in compressed form will be automatically decompressed.
// All Read methods will return uncompressed bytes.
//
// Per message compression of messages written to a connection can be enabled
// or disabled by calling the corresponding Conn method:
//
//  conn.EnableWriteCompression(false)
//
// Currently this package does not support compression with "context takeover".
// This means that messages must be compressed and decompressed in isolation,
// without retaining sliding window or dictionary state across messages. For
// more details refer to RFC 7692.
//
// Use of compression is experimental and may result in decreased performance.
package websocket
