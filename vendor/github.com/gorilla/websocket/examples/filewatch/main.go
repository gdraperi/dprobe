// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

var (
	addr      = flag.String("addr", ":8080", "http service address")
	homeTempl = template.Must(template.New("").Parse(homeHTML))
	filename  string
	upgrader  = websocket.Upgrader***REMOVED***
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	***REMOVED***
)

func readFileIfModified(lastMod time.Time) ([]byte, time.Time, error) ***REMOVED***
	fi, err := os.Stat(filename)
	if err != nil ***REMOVED***
		return nil, lastMod, err
	***REMOVED***
	if !fi.ModTime().After(lastMod) ***REMOVED***
		return nil, lastMod, nil
	***REMOVED***
	p, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		return nil, fi.ModTime(), err
	***REMOVED***
	return p, fi.ModTime(), nil
***REMOVED***

func reader(ws *websocket.Conn) ***REMOVED***
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error ***REMOVED*** ws.SetReadDeadline(time.Now().Add(pongWait)); return nil ***REMOVED***)
	for ***REMOVED***
		_, _, err := ws.ReadMessage()
		if err != nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func writer(ws *websocket.Conn, lastMod time.Time) ***REMOVED***
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(filePeriod)
	defer func() ***REMOVED***
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	***REMOVED***()
	for ***REMOVED***
		select ***REMOVED***
		case <-fileTicker.C:
			var p []byte
			var err error

			p, lastMod, err = readFileIfModified(lastMod)

			if err != nil ***REMOVED***
				if s := err.Error(); s != lastError ***REMOVED***
					lastError = s
					p = []byte(lastError)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				lastError = ""
			***REMOVED***

			if p != nil ***REMOVED***
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, p); err != nil ***REMOVED***
					return
				***REMOVED***
			***REMOVED***
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte***REMOVED******REMOVED***); err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func serveWs(w http.ResponseWriter, r *http.Request) ***REMOVED***
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil ***REMOVED***
		if _, ok := err.(websocket.HandshakeError); !ok ***REMOVED***
			log.Println(err)
		***REMOVED***
		return
	***REMOVED***

	var lastMod time.Time
	if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil ***REMOVED***
		lastMod = time.Unix(0, n)
	***REMOVED***

	go writer(ws, lastMod)
	reader(ws)
***REMOVED***

func serveHome(w http.ResponseWriter, r *http.Request) ***REMOVED***
	if r.URL.Path != "/" ***REMOVED***
		http.Error(w, "Not found", 404)
		return
	***REMOVED***
	if r.Method != "GET" ***REMOVED***
		http.Error(w, "Method not allowed", 405)
		return
	***REMOVED***
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	p, lastMod, err := readFileIfModified(time.Time***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		p = []byte(err.Error())
		lastMod = time.Unix(0, 0)
	***REMOVED***
	var v = struct ***REMOVED***
		Host    string
		Data    string
		LastMod string
	***REMOVED******REMOVED***
		r.Host,
		string(p),
		strconv.FormatInt(lastMod.UnixNano(), 16),
	***REMOVED***
	homeTempl.Execute(w, &v)
***REMOVED***

func main() ***REMOVED***
	flag.Parse()
	if flag.NArg() != 1 ***REMOVED***
		log.Fatal("filename not specified")
	***REMOVED***
	filename = flag.Args()[0]
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	if err := http.ListenAndServe(*addr, nil); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
***REMOVED***

const homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
        <pre id="fileData">***REMOVED******REMOVED***.Data***REMOVED******REMOVED***</pre>
        <script type="text/javascript">
            (function() ***REMOVED***
                var data = document.getElementById("fileData");
                var conn = new WebSocket("ws://***REMOVED******REMOVED***.Host***REMOVED******REMOVED***/ws?lastMod=***REMOVED******REMOVED***.LastMod***REMOVED******REMOVED***");
                conn.onclose = function(evt) ***REMOVED***
                    data.textContent = 'Connection closed';
            ***REMOVED***
                conn.onmessage = function(evt) ***REMOVED***
                    console.log('file updated');
                    data.textContent = evt.data;
            ***REMOVED***
        ***REMOVED***)();
        </script>
    </body>
</html>
`
