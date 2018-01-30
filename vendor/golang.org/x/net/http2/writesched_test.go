// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func makeWriteNonStreamRequest() FrameWriteRequest ***REMOVED***
	return FrameWriteRequest***REMOVED***writeSettingsAck***REMOVED******REMOVED***, nil, nil***REMOVED***
***REMOVED***

func makeWriteHeadersRequest(streamID uint32) FrameWriteRequest ***REMOVED***
	st := &stream***REMOVED***id: streamID***REMOVED***
	return FrameWriteRequest***REMOVED***&writeResHeaders***REMOVED***streamID: streamID, httpResCode: 200***REMOVED***, st, nil***REMOVED***
***REMOVED***

func checkConsume(wr FrameWriteRequest, nbytes int32, want []FrameWriteRequest) error ***REMOVED***
	consumed, rest, n := wr.Consume(nbytes)
	var wantConsumed, wantRest FrameWriteRequest
	switch len(want) ***REMOVED***
	case 0:
	case 1:
		wantConsumed = want[0]
	case 2:
		wantConsumed = want[0]
		wantRest = want[1]
	***REMOVED***
	if !reflect.DeepEqual(consumed, wantConsumed) || !reflect.DeepEqual(rest, wantRest) || n != len(want) ***REMOVED***
		return fmt.Errorf("got %v, %v, %v\nwant %v, %v, %v", consumed, rest, n, wantConsumed, wantRest, len(want))
	***REMOVED***
	return nil
***REMOVED***

func TestFrameWriteRequestNonData(t *testing.T) ***REMOVED***
	wr := makeWriteNonStreamRequest()
	if got, want := wr.DataSize(), 0; got != want ***REMOVED***
		t.Errorf("DataSize: got %v, want %v", got, want)
	***REMOVED***

	// Non-DATA frames are always consumed whole.
	if err := checkConsume(wr, 0, []FrameWriteRequest***REMOVED***wr***REMOVED***); err != nil ***REMOVED***
		t.Errorf("Consume:\n%v", err)
	***REMOVED***
***REMOVED***

func TestFrameWriteRequestData(t *testing.T) ***REMOVED***
	st := &stream***REMOVED***
		id: 1,
		sc: &serverConn***REMOVED***maxFrameSize: 16***REMOVED***,
	***REMOVED***
	const size = 32
	wr := FrameWriteRequest***REMOVED***&writeData***REMOVED***st.id, make([]byte, size), true***REMOVED***, st, make(chan error)***REMOVED***
	if got, want := wr.DataSize(), size; got != want ***REMOVED***
		t.Errorf("DataSize: got %v, want %v", got, want)
	***REMOVED***

	// No flow-control bytes available: cannot consume anything.
	if err := checkConsume(wr, math.MaxInt32, []FrameWriteRequest***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Errorf("Consume(limited by flow control):\n%v", err)
	***REMOVED***

	// Add enough flow-control bytes to consume the entire frame,
	// but we're now restricted by st.sc.maxFrameSize.
	st.flow.add(size)
	want := []FrameWriteRequest***REMOVED***
		***REMOVED***
			write:  &writeData***REMOVED***st.id, make([]byte, st.sc.maxFrameSize), false***REMOVED***,
			stream: st,
			done:   nil,
		***REMOVED***,
		***REMOVED***
			write:  &writeData***REMOVED***st.id, make([]byte, size-st.sc.maxFrameSize), true***REMOVED***,
			stream: st,
			done:   wr.done,
		***REMOVED***,
	***REMOVED***
	if err := checkConsume(wr, math.MaxInt32, want); err != nil ***REMOVED***
		t.Errorf("Consume(limited by maxFrameSize):\n%v", err)
	***REMOVED***
	rest := want[1]

	// Consume 8 bytes from the remaining frame.
	want = []FrameWriteRequest***REMOVED***
		***REMOVED***
			write:  &writeData***REMOVED***st.id, make([]byte, 8), false***REMOVED***,
			stream: st,
			done:   nil,
		***REMOVED***,
		***REMOVED***
			write:  &writeData***REMOVED***st.id, make([]byte, size-st.sc.maxFrameSize-8), true***REMOVED***,
			stream: st,
			done:   wr.done,
		***REMOVED***,
	***REMOVED***
	if err := checkConsume(rest, 8, want); err != nil ***REMOVED***
		t.Errorf("Consume(8):\n%v", err)
	***REMOVED***
	rest = want[1]

	// Consume all remaining bytes.
	want = []FrameWriteRequest***REMOVED***
		***REMOVED***
			write:  &writeData***REMOVED***st.id, make([]byte, size-st.sc.maxFrameSize-8), true***REMOVED***,
			stream: st,
			done:   wr.done,
		***REMOVED***,
	***REMOVED***
	if err := checkConsume(rest, math.MaxInt32, want); err != nil ***REMOVED***
		t.Errorf("Consume(remainder):\n%v", err)
	***REMOVED***
***REMOVED***

func TestFrameWriteRequest_StreamID(t *testing.T) ***REMOVED***
	const streamID = 123
	wr := FrameWriteRequest***REMOVED***write: streamError(streamID, ErrCodeNo)***REMOVED***
	if got := wr.StreamID(); got != streamID ***REMOVED***
		t.Errorf("FrameWriteRequest(StreamError) = %v; want %v", got, streamID)
	***REMOVED***
***REMOVED***
