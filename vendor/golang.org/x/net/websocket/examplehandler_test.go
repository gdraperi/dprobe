// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket_test

import (
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) ***REMOVED***
	io.Copy(ws, ws)
***REMOVED***

// This example demonstrates a trivial echo server.
func ExampleHandler() ***REMOVED***
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil ***REMOVED***
		panic("ListenAndServe: " + err.Error())
	***REMOVED***
***REMOVED***
