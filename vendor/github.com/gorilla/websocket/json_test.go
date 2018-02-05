// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"
)

func TestJSON(t *testing.T) ***REMOVED***
	var buf bytes.Buffer
	c := fakeNetConn***REMOVED***&buf, &buf***REMOVED***
	wc := newConn(c, true, 1024, 1024)
	rc := newConn(c, false, 1024, 1024)

	var actual, expect struct ***REMOVED***
		A int
		B string
	***REMOVED***
	expect.A = 1
	expect.B = "hello"

	if err := wc.WriteJSON(&expect); err != nil ***REMOVED***
		t.Fatal("write", err)
	***REMOVED***

	if err := rc.ReadJSON(&actual); err != nil ***REMOVED***
		t.Fatal("read", err)
	***REMOVED***

	if !reflect.DeepEqual(&actual, &expect) ***REMOVED***
		t.Fatal("equal", actual, expect)
	***REMOVED***
***REMOVED***

func TestPartialJSONRead(t *testing.T) ***REMOVED***
	var buf bytes.Buffer
	c := fakeNetConn***REMOVED***&buf, &buf***REMOVED***
	wc := newConn(c, true, 1024, 1024)
	rc := newConn(c, false, 1024, 1024)

	var v struct ***REMOVED***
		A int
		B string
	***REMOVED***
	v.A = 1
	v.B = "hello"

	messageCount := 0

	// Partial JSON values.

	data, err := json.Marshal(v)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	for i := len(data) - 1; i >= 0; i-- ***REMOVED***
		if err := wc.WriteMessage(TextMessage, data[:i]); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		messageCount++
	***REMOVED***

	// Whitespace.

	if err := wc.WriteMessage(TextMessage, []byte(" ")); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	messageCount++

	// Close.

	if err := wc.WriteMessage(CloseMessage, FormatCloseMessage(CloseNormalClosure, "")); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	for i := 0; i < messageCount; i++ ***REMOVED***
		err := rc.ReadJSON(&v)
		if err != io.ErrUnexpectedEOF ***REMOVED***
			t.Error("read", i, err)
		***REMOVED***
	***REMOVED***

	err = rc.ReadJSON(&v)
	if _, ok := err.(*CloseError); !ok ***REMOVED***
		t.Error("final", err)
	***REMOVED***
***REMOVED***

func TestDeprecatedJSON(t *testing.T) ***REMOVED***
	var buf bytes.Buffer
	c := fakeNetConn***REMOVED***&buf, &buf***REMOVED***
	wc := newConn(c, true, 1024, 1024)
	rc := newConn(c, false, 1024, 1024)

	var actual, expect struct ***REMOVED***
		A int
		B string
	***REMOVED***
	expect.A = 1
	expect.B = "hello"

	if err := WriteJSON(wc, &expect); err != nil ***REMOVED***
		t.Fatal("write", err)
	***REMOVED***

	if err := ReadJSON(rc, &actual); err != nil ***REMOVED***
		t.Fatal("read", err)
	***REMOVED***

	if !reflect.DeepEqual(&actual, &expect) ***REMOVED***
		t.Fatal("equal", actual, expect)
	***REMOVED***
***REMOVED***
