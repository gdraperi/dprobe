// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket_test

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

// This example demonstrates a trivial client.
func ExampleDial() ***REMOVED***
	origin := "http://localhost/"
	url := "ws://localhost:12345/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	if _, err := ws.Write([]byte("hello, world!\n")); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	fmt.Printf("Received: %s.\n", msg[:n])
***REMOVED***
