// Copyright 2014 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"net/url"
	"testing"
)

var hostPortNoPortTests = []struct ***REMOVED***
	u                    *url.URL
	hostPort, hostNoPort string
***REMOVED******REMOVED***
	***REMOVED***&url.URL***REMOVED***Scheme: "ws", Host: "example.com"***REMOVED***, "example.com:80", "example.com"***REMOVED***,
	***REMOVED***&url.URL***REMOVED***Scheme: "wss", Host: "example.com"***REMOVED***, "example.com:443", "example.com"***REMOVED***,
	***REMOVED***&url.URL***REMOVED***Scheme: "ws", Host: "example.com:7777"***REMOVED***, "example.com:7777", "example.com"***REMOVED***,
	***REMOVED***&url.URL***REMOVED***Scheme: "wss", Host: "example.com:7777"***REMOVED***, "example.com:7777", "example.com"***REMOVED***,
***REMOVED***

func TestHostPortNoPort(t *testing.T) ***REMOVED***
	for _, tt := range hostPortNoPortTests ***REMOVED***
		hostPort, hostNoPort := hostPortNoPort(tt.u)
		if hostPort != tt.hostPort ***REMOVED***
			t.Errorf("hostPortNoPort(%v) returned hostPort %q, want %q", tt.u, hostPort, tt.hostPort)
		***REMOVED***
		if hostNoPort != tt.hostNoPort ***REMOVED***
			t.Errorf("hostPortNoPort(%v) returned hostNoPort %q, want %q", tt.u, hostNoPort, tt.hostNoPort)
		***REMOVED***
	***REMOVED***
***REMOVED***
