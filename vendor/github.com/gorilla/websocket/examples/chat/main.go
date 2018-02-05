// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) ***REMOVED***
	log.Println(r.URL)
	if r.URL.Path != "/" ***REMOVED***
		http.Error(w, "Not found", 404)
		return
	***REMOVED***
	if r.Method != "GET" ***REMOVED***
		http.Error(w, "Method not allowed", 405)
		return
	***REMOVED***
	http.ServeFile(w, r, "home.html")
***REMOVED***

func main() ***REMOVED***
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		serveWs(hub, w, r)
	***REMOVED***)
	err := http.ListenAndServe(*addr, nil)
	if err != nil ***REMOVED***
		log.Fatal("ListenAndServe: ", err)
	***REMOVED***
***REMOVED***
