// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"net/http"
	"reflect"
	"testing"
)

var subprotocolTests = []struct ***REMOVED***
	h         string
	protocols []string
***REMOVED******REMOVED***
	***REMOVED***"", nil***REMOVED***,
	***REMOVED***"foo", []string***REMOVED***"foo"***REMOVED******REMOVED***,
	***REMOVED***"foo,bar", []string***REMOVED***"foo", "bar"***REMOVED******REMOVED***,
	***REMOVED***"foo, bar", []string***REMOVED***"foo", "bar"***REMOVED******REMOVED***,
	***REMOVED***" foo, bar", []string***REMOVED***"foo", "bar"***REMOVED******REMOVED***,
	***REMOVED***" foo, bar ", []string***REMOVED***"foo", "bar"***REMOVED******REMOVED***,
***REMOVED***

func TestSubprotocols(t *testing.T) ***REMOVED***
	for _, st := range subprotocolTests ***REMOVED***
		r := http.Request***REMOVED***Header: http.Header***REMOVED***"Sec-Websocket-Protocol": ***REMOVED***st.h***REMOVED******REMOVED******REMOVED***
		protocols := Subprotocols(&r)
		if !reflect.DeepEqual(st.protocols, protocols) ***REMOVED***
			t.Errorf("SubProtocols(%q) returned %#v, want %#v", st.h, protocols, st.protocols)
		***REMOVED***
	***REMOVED***
***REMOVED***

var isWebSocketUpgradeTests = []struct ***REMOVED***
	ok bool
	h  http.Header
***REMOVED******REMOVED***
	***REMOVED***false, http.Header***REMOVED***"Upgrade": ***REMOVED***"websocket"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***false, http.Header***REMOVED***"Connection": ***REMOVED***"upgrade"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***true, http.Header***REMOVED***"Connection": ***REMOVED***"upgRade"***REMOVED***, "Upgrade": ***REMOVED***"WebSocket"***REMOVED******REMOVED******REMOVED***,
***REMOVED***

func TestIsWebSocketUpgrade(t *testing.T) ***REMOVED***
	for _, tt := range isWebSocketUpgradeTests ***REMOVED***
		ok := IsWebSocketUpgrade(&http.Request***REMOVED***Header: tt.h***REMOVED***)
		if tt.ok != ok ***REMOVED***
			t.Errorf("IsWebSocketUpgrade(%v) returned %v, want %v", tt.h, ok, tt.ok)
		***REMOVED***
	***REMOVED***
***REMOVED***

var checkSameOriginTests = []struct ***REMOVED***
	ok bool
	r  *http.Request
***REMOVED******REMOVED***
	***REMOVED***false, &http.Request***REMOVED***Host: "example.org", Header: map[string][]string***REMOVED***"Origin": []string***REMOVED***"https://other.org"***REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***true, &http.Request***REMOVED***Host: "example.org", Header: map[string][]string***REMOVED***"Origin": []string***REMOVED***"https://example.org"***REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***true, &http.Request***REMOVED***Host: "Example.org", Header: map[string][]string***REMOVED***"Origin": []string***REMOVED***"https://example.org"***REMOVED******REMOVED******REMOVED******REMOVED***,
***REMOVED***

func TestCheckSameOrigin(t *testing.T) ***REMOVED***
	for _, tt := range checkSameOriginTests ***REMOVED***
		ok := checkSameOrigin(tt.r)
		if tt.ok != ok ***REMOVED***
			t.Errorf("checkSameOrigin(%+v) returned %v, want %v", tt.r, ok, tt.ok)
		***REMOVED***
	***REMOVED***
***REMOVED***
