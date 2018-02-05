// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader***REMOVED******REMOVED*** // use default options

func echo(w http.ResponseWriter, r *http.Request) ***REMOVED***
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil ***REMOVED***
		log.Print("upgrade:", err)
		return
	***REMOVED***
	defer c.Close()
	for ***REMOVED***
		mt, message, err := c.ReadMessage()
		if err != nil ***REMOVED***
			log.Println("read:", err)
			break
		***REMOVED***
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil ***REMOVED***
			log.Println("write:", err)
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func home(w http.ResponseWriter, r *http.Request) ***REMOVED***
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
***REMOVED***

func main() ***REMOVED***
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
***REMOVED***

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) ***REMOVED***

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) ***REMOVED***
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
***REMOVED***;

    document.getElementById("open").onclick = function(evt) ***REMOVED***
        if (ws) ***REMOVED***
            return false;
    ***REMOVED***
        ws = new WebSocket("***REMOVED******REMOVED***.***REMOVED******REMOVED***");
        ws.onopen = function(evt) ***REMOVED***
            print("OPEN");
    ***REMOVED***
        ws.onclose = function(evt) ***REMOVED***
            print("CLOSE");
            ws = null;
    ***REMOVED***
        ws.onmessage = function(evt) ***REMOVED***
            print("RESPONSE: " + evt.data);
    ***REMOVED***
        ws.onerror = function(evt) ***REMOVED***
            print("ERROR: " + evt.data);
    ***REMOVED***
        return false;
***REMOVED***;

    document.getElementById("send").onclick = function(evt) ***REMOVED***
        if (!ws) ***REMOVED***
            return false;
    ***REMOVED***
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
***REMOVED***;

    document.getElementById("close").onclick = function(evt) ***REMOVED***
        if (!ws) ***REMOVED***
            return false;
    ***REMOVED***
        ws.close();
        return false;
***REMOVED***;

***REMOVED***);
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
