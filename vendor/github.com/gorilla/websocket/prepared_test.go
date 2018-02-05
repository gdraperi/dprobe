// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bytes"
	"compress/flate"
	"math/rand"
	"testing"
)

var preparedMessageTests = []struct ***REMOVED***
	messageType            int
	isServer               bool
	enableWriteCompression bool
	compressionLevel       int
***REMOVED******REMOVED***
	// Server
	***REMOVED***TextMessage, true, false, flate.BestSpeed***REMOVED***,
	***REMOVED***TextMessage, true, true, flate.BestSpeed***REMOVED***,
	***REMOVED***TextMessage, true, true, flate.BestCompression***REMOVED***,
	***REMOVED***PingMessage, true, false, flate.BestSpeed***REMOVED***,
	***REMOVED***PingMessage, true, true, flate.BestSpeed***REMOVED***,

	// Client
	***REMOVED***TextMessage, false, false, flate.BestSpeed***REMOVED***,
	***REMOVED***TextMessage, false, true, flate.BestSpeed***REMOVED***,
	***REMOVED***TextMessage, false, true, flate.BestCompression***REMOVED***,
	***REMOVED***PingMessage, false, false, flate.BestSpeed***REMOVED***,
	***REMOVED***PingMessage, false, true, flate.BestSpeed***REMOVED***,
***REMOVED***

func TestPreparedMessage(t *testing.T) ***REMOVED***
	for _, tt := range preparedMessageTests ***REMOVED***
		var data = []byte("this is a test")
		var buf bytes.Buffer
		c := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: &buf***REMOVED***, tt.isServer, 1024, 1024)
		if tt.enableWriteCompression ***REMOVED***
			c.newCompressionWriter = compressNoContextTakeover
		***REMOVED***
		c.SetCompressionLevel(tt.compressionLevel)

		// Seed random number generator for consistent frame mask.
		rand.Seed(1234)

		if err := c.WriteMessage(tt.messageType, data); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		want := buf.String()

		pm, err := NewPreparedMessage(tt.messageType, data)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		// Scribble on data to ensure that NewPreparedMessage takes a snapshot.
		copy(data, "hello world")

		// Seed random number generator for consistent frame mask.
		rand.Seed(1234)

		buf.Reset()
		if err := c.WritePreparedMessage(pm); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		got := buf.String()

		if got != want ***REMOVED***
			t.Errorf("write message != prepared message for %+v", tt)
		***REMOVED***
	***REMOVED***
***REMOVED***
