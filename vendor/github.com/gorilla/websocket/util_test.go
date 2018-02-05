// Copyright 2014 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"net/http"
	"reflect"
	"testing"
)

var equalASCIIFoldTests = []struct ***REMOVED***
	t, s string
	eq   bool
***REMOVED******REMOVED***
	***REMOVED***"WebSocket", "websocket", true***REMOVED***,
	***REMOVED***"websocket", "WebSocket", true***REMOVED***,
	***REMOVED***"Öyster", "öyster", false***REMOVED***,
***REMOVED***

func TestEqualASCIIFold(t *testing.T) ***REMOVED***
	for _, tt := range equalASCIIFoldTests ***REMOVED***
		eq := equalASCIIFold(tt.s, tt.t)
		if eq != tt.eq ***REMOVED***
			t.Errorf("equalASCIIFold(%q, %q) = %v, want %v", tt.s, tt.t, eq, tt.eq)
		***REMOVED***
	***REMOVED***
***REMOVED***

var tokenListContainsValueTests = []struct ***REMOVED***
	value string
	ok    bool
***REMOVED******REMOVED***
	***REMOVED***"WebSocket", true***REMOVED***,
	***REMOVED***"WEBSOCKET", true***REMOVED***,
	***REMOVED***"websocket", true***REMOVED***,
	***REMOVED***"websockets", false***REMOVED***,
	***REMOVED***"x websocket", false***REMOVED***,
	***REMOVED***"websocket x", false***REMOVED***,
	***REMOVED***"other,websocket,more", true***REMOVED***,
	***REMOVED***"other, websocket, more", true***REMOVED***,
***REMOVED***

func TestTokenListContainsValue(t *testing.T) ***REMOVED***
	for _, tt := range tokenListContainsValueTests ***REMOVED***
		h := http.Header***REMOVED***"Upgrade": ***REMOVED***tt.value***REMOVED******REMOVED***
		ok := tokenListContainsValue(h, "Upgrade", "websocket")
		if ok != tt.ok ***REMOVED***
			t.Errorf("tokenListContainsValue(h, n, %q) = %v, want %v", tt.value, ok, tt.ok)
		***REMOVED***
	***REMOVED***
***REMOVED***

var parseExtensionTests = []struct ***REMOVED***
	value      string
	extensions []map[string]string
***REMOVED******REMOVED***
	***REMOVED***`foo`, []map[string]string***REMOVED******REMOVED***"": "foo"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`foo, bar; baz=2`, []map[string]string***REMOVED***
		***REMOVED***"": "foo"***REMOVED***,
		***REMOVED***"": "bar", "baz": "2"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`foo; bar="b,a;z"`, []map[string]string***REMOVED***
		***REMOVED***"": "foo", "bar": "b,a;z"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`foo , bar; baz = 2`, []map[string]string***REMOVED***
		***REMOVED***"": "foo"***REMOVED***,
		***REMOVED***"": "bar", "baz": "2"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`foo, bar; baz=2 junk`, []map[string]string***REMOVED***
		***REMOVED***"": "foo"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`foo junk, bar; baz=2 junk`, nil***REMOVED***,
	***REMOVED***`mux; max-channels=4; flow-control, deflate-stream`, []map[string]string***REMOVED***
		***REMOVED***"": "mux", "max-channels": "4", "flow-control": ""***REMOVED***,
		***REMOVED***"": "deflate-stream"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`permessage-foo; x="10"`, []map[string]string***REMOVED***
		***REMOVED***"": "permessage-foo", "x": "10"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`permessage-foo; use_y, permessage-foo`, []map[string]string***REMOVED***
		***REMOVED***"": "permessage-foo", "use_y": ""***REMOVED***,
		***REMOVED***"": "permessage-foo"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`permessage-deflate; client_max_window_bits; server_max_window_bits=10 , permessage-deflate; client_max_window_bits`, []map[string]string***REMOVED***
		***REMOVED***"": "permessage-deflate", "client_max_window_bits": "", "server_max_window_bits": "10"***REMOVED***,
		***REMOVED***"": "permessage-deflate", "client_max_window_bits": ""***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"permessage-deflate; server_no_context_takeover; client_max_window_bits=15", []map[string]string***REMOVED***
		***REMOVED***"": "permessage-deflate", "server_no_context_takeover": "", "client_max_window_bits": "15"***REMOVED***,
	***REMOVED******REMOVED***,
***REMOVED***

func TestParseExtensions(t *testing.T) ***REMOVED***
	for _, tt := range parseExtensionTests ***REMOVED***
		h := http.Header***REMOVED***http.CanonicalHeaderKey("Sec-WebSocket-Extensions"): ***REMOVED***tt.value***REMOVED******REMOVED***
		extensions := parseExtensions(h)
		if !reflect.DeepEqual(extensions, tt.extensions) ***REMOVED***
			t.Errorf("parseExtensions(%q)\n    = %v,\nwant %v", tt.value, extensions, tt.extensions)
		***REMOVED***
	***REMOVED***
***REMOVED***
