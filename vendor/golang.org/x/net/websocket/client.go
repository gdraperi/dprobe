// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/url"
)

// DialError is an error that occurs while dialling a websocket server.
type DialError struct ***REMOVED***
	*Config
	Err error
***REMOVED***

func (e *DialError) Error() string ***REMOVED***
	return "websocket.Dial " + e.Config.Location.String() + ": " + e.Err.Error()
***REMOVED***

// NewConfig creates a new WebSocket config for client connection.
func NewConfig(server, origin string) (config *Config, err error) ***REMOVED***
	config = new(Config)
	config.Version = ProtocolVersionHybi13
	config.Location, err = url.ParseRequestURI(server)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	config.Origin, err = url.ParseRequestURI(origin)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	config.Header = http.Header(make(map[string][]string))
	return
***REMOVED***

// NewClient creates a new WebSocket client connection over rwc.
func NewClient(config *Config, rwc io.ReadWriteCloser) (ws *Conn, err error) ***REMOVED***
	br := bufio.NewReader(rwc)
	bw := bufio.NewWriter(rwc)
	err = hybiClientHandshake(config, br, bw)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	buf := bufio.NewReadWriter(br, bw)
	ws = newHybiClientConn(config, buf, rwc)
	return
***REMOVED***

// Dial opens a new client connection to a WebSocket.
func Dial(url_, protocol, origin string) (ws *Conn, err error) ***REMOVED***
	config, err := NewConfig(url_, origin)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if protocol != "" ***REMOVED***
		config.Protocol = []string***REMOVED***protocol***REMOVED***
	***REMOVED***
	return DialConfig(config)
***REMOVED***

var portMap = map[string]string***REMOVED***
	"ws":  "80",
	"wss": "443",
***REMOVED***

func parseAuthority(location *url.URL) string ***REMOVED***
	if _, ok := portMap[location.Scheme]; ok ***REMOVED***
		if _, _, err := net.SplitHostPort(location.Host); err != nil ***REMOVED***
			return net.JoinHostPort(location.Host, portMap[location.Scheme])
		***REMOVED***
	***REMOVED***
	return location.Host
***REMOVED***

// DialConfig opens a new client connection to a WebSocket with a config.
func DialConfig(config *Config) (ws *Conn, err error) ***REMOVED***
	var client net.Conn
	if config.Location == nil ***REMOVED***
		return nil, &DialError***REMOVED***config, ErrBadWebSocketLocation***REMOVED***
	***REMOVED***
	if config.Origin == nil ***REMOVED***
		return nil, &DialError***REMOVED***config, ErrBadWebSocketOrigin***REMOVED***
	***REMOVED***
	dialer := config.Dialer
	if dialer == nil ***REMOVED***
		dialer = &net.Dialer***REMOVED******REMOVED***
	***REMOVED***
	client, err = dialWithDialer(dialer, config)
	if err != nil ***REMOVED***
		goto Error
	***REMOVED***
	ws, err = NewClient(config, client)
	if err != nil ***REMOVED***
		client.Close()
		goto Error
	***REMOVED***
	return

Error:
	return nil, &DialError***REMOVED***config, err***REMOVED***
***REMOVED***
