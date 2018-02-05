// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command server is a test server for the Autobahn WebSockets Test Suite.
package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader***REMOVED***
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool ***REMOVED***
		return true
	***REMOVED***,
***REMOVED***

// echoCopy echoes messages from the client using io.Copy.
func echoCopy(w http.ResponseWriter, r *http.Request, writerOnly bool) ***REMOVED***
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil ***REMOVED***
		log.Println("Upgrade:", err)
		return
	***REMOVED***
	defer conn.Close()
	for ***REMOVED***
		mt, r, err := conn.NextReader()
		if err != nil ***REMOVED***
			if err != io.EOF ***REMOVED***
				log.Println("NextReader:", err)
			***REMOVED***
			return
		***REMOVED***
		if mt == websocket.TextMessage ***REMOVED***
			r = &validator***REMOVED***r: r***REMOVED***
		***REMOVED***
		w, err := conn.NextWriter(mt)
		if err != nil ***REMOVED***
			log.Println("NextWriter:", err)
			return
		***REMOVED***
		if mt == websocket.TextMessage ***REMOVED***
			r = &validator***REMOVED***r: r***REMOVED***
		***REMOVED***
		if writerOnly ***REMOVED***
			_, err = io.Copy(struct***REMOVED*** io.Writer ***REMOVED******REMOVED***w***REMOVED***, r)
		***REMOVED*** else ***REMOVED***
			_, err = io.Copy(w, r)
		***REMOVED***
		if err != nil ***REMOVED***
			if err == errInvalidUTF8 ***REMOVED***
				conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, ""),
					time.Time***REMOVED******REMOVED***)
			***REMOVED***
			log.Println("Copy:", err)
			return
		***REMOVED***
		err = w.Close()
		if err != nil ***REMOVED***
			log.Println("Close:", err)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func echoCopyWriterOnly(w http.ResponseWriter, r *http.Request) ***REMOVED***
	echoCopy(w, r, true)
***REMOVED***

func echoCopyFull(w http.ResponseWriter, r *http.Request) ***REMOVED***
	echoCopy(w, r, false)
***REMOVED***

// echoReadAll echoes messages from the client by reading the entire message
// with ioutil.ReadAll.
func echoReadAll(w http.ResponseWriter, r *http.Request, writeMessage, writePrepared bool) ***REMOVED***
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil ***REMOVED***
		log.Println("Upgrade:", err)
		return
	***REMOVED***
	defer conn.Close()
	for ***REMOVED***
		mt, b, err := conn.ReadMessage()
		if err != nil ***REMOVED***
			if err != io.EOF ***REMOVED***
				log.Println("NextReader:", err)
			***REMOVED***
			return
		***REMOVED***
		if mt == websocket.TextMessage ***REMOVED***
			if !utf8.Valid(b) ***REMOVED***
				conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, ""),
					time.Time***REMOVED******REMOVED***)
				log.Println("ReadAll: invalid utf8")
			***REMOVED***
		***REMOVED***
		if writeMessage ***REMOVED***
			if !writePrepared ***REMOVED***
				err = conn.WriteMessage(mt, b)
				if err != nil ***REMOVED***
					log.Println("WriteMessage:", err)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				pm, err := websocket.NewPreparedMessage(mt, b)
				if err != nil ***REMOVED***
					log.Println("NewPreparedMessage:", err)
					return
				***REMOVED***
				err = conn.WritePreparedMessage(pm)
				if err != nil ***REMOVED***
					log.Println("WritePreparedMessage:", err)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			w, err := conn.NextWriter(mt)
			if err != nil ***REMOVED***
				log.Println("NextWriter:", err)
				return
			***REMOVED***
			if _, err := w.Write(b); err != nil ***REMOVED***
				log.Println("Writer:", err)
				return
			***REMOVED***
			if err := w.Close(); err != nil ***REMOVED***
				log.Println("Close:", err)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func echoReadAllWriter(w http.ResponseWriter, r *http.Request) ***REMOVED***
	echoReadAll(w, r, false, false)
***REMOVED***

func echoReadAllWriteMessage(w http.ResponseWriter, r *http.Request) ***REMOVED***
	echoReadAll(w, r, true, false)
***REMOVED***

func echoReadAllWritePreparedMessage(w http.ResponseWriter, r *http.Request) ***REMOVED***
	echoReadAll(w, r, true, true)
***REMOVED***

func serveHome(w http.ResponseWriter, r *http.Request) ***REMOVED***
	if r.URL.Path != "/" ***REMOVED***
		http.Error(w, "Not found.", 404)
		return
	***REMOVED***
	if r.Method != "GET" ***REMOVED***
		http.Error(w, "Method not allowed", 405)
		return
	***REMOVED***
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, "<html><body>Echo Server</body></html>")
***REMOVED***

var addr = flag.String("addr", ":9000", "http service address")

func main() ***REMOVED***
	flag.Parse()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/c", echoCopyWriterOnly)
	http.HandleFunc("/f", echoCopyFull)
	http.HandleFunc("/r", echoReadAllWriter)
	http.HandleFunc("/m", echoReadAllWriteMessage)
	http.HandleFunc("/p", echoReadAllWritePreparedMessage)
	err := http.ListenAndServe(*addr, nil)
	if err != nil ***REMOVED***
		log.Fatal("ListenAndServe: ", err)
	***REMOVED***
***REMOVED***

type validator struct ***REMOVED***
	state int
	x     rune
	r     io.Reader
***REMOVED***

var errInvalidUTF8 = errors.New("invalid utf8")

func (r *validator) Read(p []byte) (int, error) ***REMOVED***
	n, err := r.r.Read(p)
	state := r.state
	x := r.x
	for _, b := range p[:n] ***REMOVED***
		state, x = decode(state, x, b)
		if state == utf8Reject ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	r.state = state
	r.x = x
	if state == utf8Reject || (err == io.EOF && state != utf8Accept) ***REMOVED***
		return n, errInvalidUTF8
	***REMOVED***
	return n, err
***REMOVED***

// UTF-8 decoder from http://bjoern.hoehrmann.de/utf-8/decoder/dfa/
//
// Copyright (c) 2008-2009 Bjoern Hoehrmann <bjoern@hoehrmann.de>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.
var utf8d = [...]byte***REMOVED***
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 00..1f
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 20..3f
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 40..5f
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 60..7f
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, // 80..9f
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, // a0..bf
	8, 8, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, // c0..df
	0xa, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x4, 0x3, 0x3, // e0..ef
	0xb, 0x6, 0x6, 0x6, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, // f0..ff
	0x0, 0x1, 0x2, 0x3, 0x5, 0x8, 0x7, 0x1, 0x1, 0x1, 0x4, 0x6, 0x1, 0x1, 0x1, 0x1, // s0..s0
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, // s1..s2
	1, 2, 1, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, // s3..s4
	1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 1, 1, // s5..s6
	1, 3, 1, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 1, 1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // s7..s8
***REMOVED***

const (
	utf8Accept = 0
	utf8Reject = 1
)

func decode(state int, x rune, b byte) (int, rune) ***REMOVED***
	t := utf8d[b]
	if state != utf8Accept ***REMOVED***
		x = rune(b&0x3f) | (x << 6)
	***REMOVED*** else ***REMOVED***
		x = rune((0xff >> t) & b)
	***REMOVED***
	state = int(utf8d[256+state*16+int(t)])
	return state, x
***REMOVED***
