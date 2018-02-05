// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr    = flag.String("addr", "127.0.0.1:8080", "http service address")
	cmdPath string
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

func pumpStdin(ws *websocket.Conn, w io.Writer) ***REMOVED***
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error ***REMOVED*** ws.SetReadDeadline(time.Now().Add(pongWait)); return nil ***REMOVED***)
	for ***REMOVED***
		_, message, err := ws.ReadMessage()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		message = append(message, '\n')
		if _, err := w.Write(message); err != nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func pumpStdout(ws *websocket.Conn, r io.Reader, done chan struct***REMOVED******REMOVED***) ***REMOVED***
	defer func() ***REMOVED***
	***REMOVED***()
	s := bufio.NewScanner(r)
	for s.Scan() ***REMOVED***
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil ***REMOVED***
			ws.Close()
			break
		***REMOVED***
	***REMOVED***
	if s.Err() != nil ***REMOVED***
		log.Println("scan:", s.Err())
	***REMOVED***
	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
***REMOVED***

func ping(ws *websocket.Conn, done chan struct***REMOVED******REMOVED***) ***REMOVED***
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte***REMOVED******REMOVED***, time.Now().Add(writeWait)); err != nil ***REMOVED***
				log.Println("ping:", err)
			***REMOVED***
		case <-done:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func internalError(ws *websocket.Conn, msg string, err error) ***REMOVED***
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
***REMOVED***

var upgrader = websocket.Upgrader***REMOVED******REMOVED***

func serveWs(w http.ResponseWriter, r *http.Request) ***REMOVED***
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil ***REMOVED***
		log.Println("upgrade:", err)
		return
	***REMOVED***

	defer ws.Close()

	outr, outw, err := os.Pipe()
	if err != nil ***REMOVED***
		internalError(ws, "stdout:", err)
		return
	***REMOVED***
	defer outr.Close()
	defer outw.Close()

	inr, inw, err := os.Pipe()
	if err != nil ***REMOVED***
		internalError(ws, "stdin:", err)
		return
	***REMOVED***
	defer inr.Close()
	defer inw.Close()

	proc, err := os.StartProcess(cmdPath, flag.Args(), &os.ProcAttr***REMOVED***
		Files: []*os.File***REMOVED***inr, outw, outw***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		internalError(ws, "start:", err)
		return
	***REMOVED***

	inr.Close()
	outw.Close()

	stdoutDone := make(chan struct***REMOVED******REMOVED***)
	go pumpStdout(ws, outr, stdoutDone)
	go ping(ws, stdoutDone)

	pumpStdin(ws, inw)

	// Some commands will exit when stdin is closed.
	inw.Close()

	// Other commands need a bonk on the head.
	if err := proc.Signal(os.Interrupt); err != nil ***REMOVED***
		log.Println("inter:", err)
	***REMOVED***

	select ***REMOVED***
	case <-stdoutDone:
	case <-time.After(time.Second):
		// A bigger bonk on the head.
		if err := proc.Signal(os.Kill); err != nil ***REMOVED***
			log.Println("term:", err)
		***REMOVED***
		<-stdoutDone
	***REMOVED***

	if _, err := proc.Wait(); err != nil ***REMOVED***
		log.Println("wait:", err)
	***REMOVED***
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
	http.ServeFile(w, r, "home.html")
***REMOVED***

func main() ***REMOVED***
	flag.Parse()
	if len(flag.Args()) < 1 ***REMOVED***
		log.Fatal("must specify at least one argument")
	***REMOVED***
	var err error
	cmdPath, err = exec.LookPath(flag.Args()[0])
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(*addr, nil))
***REMOVED***
