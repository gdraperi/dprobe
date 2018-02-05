// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() ***REMOVED***
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL***REMOVED***Scheme: "ws", Host: *addr, Path: "/echo"***REMOVED***
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil ***REMOVED***
		log.Fatal("dial:", err)
	***REMOVED***
	defer c.Close()

	done := make(chan struct***REMOVED******REMOVED***)

	go func() ***REMOVED***
		defer c.Close()
		defer close(done)
		for ***REMOVED***
			_, message, err := c.ReadMessage()
			if err != nil ***REMOVED***
				log.Println("read:", err)
				return
			***REMOVED***
			log.Printf("recv: %s", message)
		***REMOVED***
	***REMOVED***()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for ***REMOVED***
		select ***REMOVED***
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil ***REMOVED***
				log.Println("write:", err)
				return
			***REMOVED***
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil ***REMOVED***
				log.Println("write close:", err)
				return
			***REMOVED***
			select ***REMOVED***
			case <-done:
			case <-time.After(time.Second):
			***REMOVED***
			c.Close()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***
