// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"bytes"
	"fmt"
	"sort"
	"testing"
)

func defaultPriorityWriteScheduler() *priorityWriteScheduler ***REMOVED***
	return NewPriorityWriteScheduler(nil).(*priorityWriteScheduler)
***REMOVED***

func checkPriorityWellFormed(ws *priorityWriteScheduler) error ***REMOVED***
	for id, n := range ws.nodes ***REMOVED***
		if id != n.id ***REMOVED***
			return fmt.Errorf("bad ws.nodes: ws.nodes[%d] = %d", id, n.id)
		***REMOVED***
		if n.parent == nil ***REMOVED***
			if n.next != nil || n.prev != nil ***REMOVED***
				return fmt.Errorf("bad node %d: nil parent but prev/next not nil", id)
			***REMOVED***
			continue
		***REMOVED***
		found := false
		for k := n.parent.kids; k != nil; k = k.next ***REMOVED***
			if k.id == id ***REMOVED***
				found = true
				break
			***REMOVED***
		***REMOVED***
		if !found ***REMOVED***
			return fmt.Errorf("bad node %d: not found in parent %d kids list", id, n.parent.id)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func fmtTree(ws *priorityWriteScheduler, fmtNode func(*priorityNode) string) string ***REMOVED***
	var ids []int
	for _, n := range ws.nodes ***REMOVED***
		ids = append(ids, int(n.id))
	***REMOVED***
	sort.Ints(ids)

	var buf bytes.Buffer
	for _, id := range ids ***REMOVED***
		if buf.Len() != 0 ***REMOVED***
			buf.WriteString(" ")
		***REMOVED***
		if id == 0 ***REMOVED***
			buf.WriteString(fmtNode(&ws.root))
		***REMOVED*** else ***REMOVED***
			buf.WriteString(fmtNode(ws.nodes[uint32(id)]))
		***REMOVED***
	***REMOVED***
	return buf.String()
***REMOVED***

func fmtNodeParentSkipRoot(n *priorityNode) string ***REMOVED***
	switch ***REMOVED***
	case n.id == 0:
		return ""
	case n.parent == nil:
		return fmt.Sprintf("%d***REMOVED***parent:nil***REMOVED***", n.id)
	default:
		return fmt.Sprintf("%d***REMOVED***parent:%d***REMOVED***", n.id, n.parent.id)
	***REMOVED***
***REMOVED***

func fmtNodeWeightParentSkipRoot(n *priorityNode) string ***REMOVED***
	switch ***REMOVED***
	case n.id == 0:
		return ""
	case n.parent == nil:
		return fmt.Sprintf("%d***REMOVED***weight:%d,parent:nil***REMOVED***", n.id, n.weight)
	default:
		return fmt.Sprintf("%d***REMOVED***weight:%d,parent:%d***REMOVED***", n.id, n.weight, n.parent.id)
	***REMOVED***
***REMOVED***

func TestPriorityTwoStreams(t *testing.T) ***REMOVED***
	ws := defaultPriorityWriteScheduler()
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED******REMOVED***)

	want := "1***REMOVED***weight:15,parent:0***REMOVED*** 2***REMOVED***weight:15,parent:0***REMOVED***"
	if got := fmtTree(ws, fmtNodeWeightParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After open\ngot  %q\nwant %q", got, want)
	***REMOVED***

	// Move 1's parent to 2.
	ws.AdjustStream(1, PriorityParam***REMOVED***
		StreamDep: 2,
		Weight:    32,
		Exclusive: false,
	***REMOVED***)
	want = "1***REMOVED***weight:32,parent:2***REMOVED*** 2***REMOVED***weight:15,parent:0***REMOVED***"
	if got := fmtTree(ws, fmtNodeWeightParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After adjust\ngot  %q\nwant %q", got, want)
	***REMOVED***

	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPriorityAdjustExclusiveZero(t *testing.T) ***REMOVED***
	// 1, 2, and 3 are all children of the 0 stream.
	// Exclusive reprioritization to any of the streams should bring
	// the rest of the streams under the reprioritized stream.
	ws := defaultPriorityWriteScheduler()
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(3, OpenStreamOptions***REMOVED******REMOVED***)

	want := "1***REMOVED***weight:15,parent:0***REMOVED*** 2***REMOVED***weight:15,parent:0***REMOVED*** 3***REMOVED***weight:15,parent:0***REMOVED***"
	if got := fmtTree(ws, fmtNodeWeightParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After open\ngot  %q\nwant %q", got, want)
	***REMOVED***

	ws.AdjustStream(2, PriorityParam***REMOVED***
		StreamDep: 0,
		Weight:    20,
		Exclusive: true,
	***REMOVED***)
	want = "1***REMOVED***weight:15,parent:2***REMOVED*** 2***REMOVED***weight:20,parent:0***REMOVED*** 3***REMOVED***weight:15,parent:2***REMOVED***"
	if got := fmtTree(ws, fmtNodeWeightParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After adjust\ngot  %q\nwant %q", got, want)
	***REMOVED***

	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPriorityAdjustOwnParent(t *testing.T) ***REMOVED***
	// Assigning a node as its own parent should have no effect.
	ws := defaultPriorityWriteScheduler()
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED******REMOVED***)
	ws.AdjustStream(2, PriorityParam***REMOVED***
		StreamDep: 2,
		Weight:    20,
		Exclusive: true,
	***REMOVED***)
	want := "1***REMOVED***weight:15,parent:0***REMOVED*** 2***REMOVED***weight:15,parent:0***REMOVED***"
	if got := fmtTree(ws, fmtNodeWeightParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After adjust\ngot  %q\nwant %q", got, want)
	***REMOVED***
	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPriorityClosedStreams(t *testing.T) ***REMOVED***
	ws := NewPriorityWriteScheduler(&PriorityWriteSchedulerConfig***REMOVED***MaxClosedNodesInTree: 2***REMOVED***).(*priorityWriteScheduler)
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(3, OpenStreamOptions***REMOVED***PusherID: 2***REMOVED***)
	ws.OpenStream(4, OpenStreamOptions***REMOVED***PusherID: 3***REMOVED***)

	// Close the first three streams. We lose 1, but keep 2 and 3.
	ws.CloseStream(1)
	ws.CloseStream(2)
	ws.CloseStream(3)

	want := "2***REMOVED***weight:15,parent:0***REMOVED*** 3***REMOVED***weight:15,parent:2***REMOVED*** 4***REMOVED***weight:15,parent:3***REMOVED***"
	if got := fmtTree(ws, fmtNodeWeightParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After close\ngot  %q\nwant %q", got, want)
	***REMOVED***
	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	// Adding a stream as an exclusive child of 1 gives it default
	// priorities, since 1 is gone.
	ws.OpenStream(5, OpenStreamOptions***REMOVED******REMOVED***)
	ws.AdjustStream(5, PriorityParam***REMOVED***StreamDep: 1, Weight: 15, Exclusive: true***REMOVED***)

	// Adding a stream as an exclusive child of 2 should work, since 2 is not gone.
	ws.OpenStream(6, OpenStreamOptions***REMOVED******REMOVED***)
	ws.AdjustStream(6, PriorityParam***REMOVED***StreamDep: 2, Weight: 15, Exclusive: true***REMOVED***)

	want = "2***REMOVED***weight:15,parent:0***REMOVED*** 3***REMOVED***weight:15,parent:6***REMOVED*** 4***REMOVED***weight:15,parent:3***REMOVED*** 5***REMOVED***weight:15,parent:0***REMOVED*** 6***REMOVED***weight:15,parent:2***REMOVED***"
	if got := fmtTree(ws, fmtNodeWeightParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After add streams\ngot  %q\nwant %q", got, want)
	***REMOVED***
	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPriorityClosedStreamsDisabled(t *testing.T) ***REMOVED***
	ws := NewPriorityWriteScheduler(&PriorityWriteSchedulerConfig***REMOVED******REMOVED***).(*priorityWriteScheduler)
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(3, OpenStreamOptions***REMOVED***PusherID: 2***REMOVED***)

	// Close the first two streams. We keep only 3.
	ws.CloseStream(1)
	ws.CloseStream(2)

	want := "3***REMOVED***weight:15,parent:0***REMOVED***"
	if got := fmtTree(ws, fmtNodeWeightParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After close\ngot  %q\nwant %q", got, want)
	***REMOVED***
	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPriorityIdleStreams(t *testing.T) ***REMOVED***
	ws := NewPriorityWriteScheduler(&PriorityWriteSchedulerConfig***REMOVED***MaxIdleNodesInTree: 2***REMOVED***).(*priorityWriteScheduler)
	ws.AdjustStream(1, PriorityParam***REMOVED***StreamDep: 0, Weight: 15***REMOVED***) // idle
	ws.AdjustStream(2, PriorityParam***REMOVED***StreamDep: 0, Weight: 15***REMOVED***) // idle
	ws.AdjustStream(3, PriorityParam***REMOVED***StreamDep: 2, Weight: 20***REMOVED***) // idle
	ws.OpenStream(4, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(5, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(6, OpenStreamOptions***REMOVED******REMOVED***)
	ws.AdjustStream(4, PriorityParam***REMOVED***StreamDep: 1, Weight: 15***REMOVED***)
	ws.AdjustStream(5, PriorityParam***REMOVED***StreamDep: 2, Weight: 15***REMOVED***)
	ws.AdjustStream(6, PriorityParam***REMOVED***StreamDep: 3, Weight: 15***REMOVED***)

	want := "2***REMOVED***weight:15,parent:0***REMOVED*** 3***REMOVED***weight:20,parent:2***REMOVED*** 4***REMOVED***weight:15,parent:0***REMOVED*** 5***REMOVED***weight:15,parent:2***REMOVED*** 6***REMOVED***weight:15,parent:3***REMOVED***"
	if got := fmtTree(ws, fmtNodeWeightParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After open\ngot  %q\nwant %q", got, want)
	***REMOVED***
	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPriorityIdleStreamsDisabled(t *testing.T) ***REMOVED***
	ws := NewPriorityWriteScheduler(&PriorityWriteSchedulerConfig***REMOVED******REMOVED***).(*priorityWriteScheduler)
	ws.AdjustStream(1, PriorityParam***REMOVED***StreamDep: 0, Weight: 15***REMOVED***) // idle
	ws.AdjustStream(2, PriorityParam***REMOVED***StreamDep: 0, Weight: 15***REMOVED***) // idle
	ws.AdjustStream(3, PriorityParam***REMOVED***StreamDep: 2, Weight: 20***REMOVED***) // idle
	ws.OpenStream(4, OpenStreamOptions***REMOVED******REMOVED***)

	want := "4***REMOVED***weight:15,parent:0***REMOVED***"
	if got := fmtTree(ws, fmtNodeWeightParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After open\ngot  %q\nwant %q", got, want)
	***REMOVED***
	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPrioritySection531NonExclusive(t *testing.T) ***REMOVED***
	// Example from RFC 7540 Section 5.3.1.
	// A,B,C,D = 1,2,3,4
	ws := defaultPriorityWriteScheduler()
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(3, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(4, OpenStreamOptions***REMOVED******REMOVED***)
	ws.AdjustStream(4, PriorityParam***REMOVED***
		StreamDep: 1,
		Weight:    15,
		Exclusive: false,
	***REMOVED***)
	want := "1***REMOVED***parent:0***REMOVED*** 2***REMOVED***parent:1***REMOVED*** 3***REMOVED***parent:1***REMOVED*** 4***REMOVED***parent:1***REMOVED***"
	if got := fmtTree(ws, fmtNodeParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After adjust\ngot  %q\nwant %q", got, want)
	***REMOVED***
	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPrioritySection531Exclusive(t *testing.T) ***REMOVED***
	// Example from RFC 7540 Section 5.3.1.
	// A,B,C,D = 1,2,3,4
	ws := defaultPriorityWriteScheduler()
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(3, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(4, OpenStreamOptions***REMOVED******REMOVED***)
	ws.AdjustStream(4, PriorityParam***REMOVED***
		StreamDep: 1,
		Weight:    15,
		Exclusive: true,
	***REMOVED***)
	want := "1***REMOVED***parent:0***REMOVED*** 2***REMOVED***parent:4***REMOVED*** 3***REMOVED***parent:4***REMOVED*** 4***REMOVED***parent:1***REMOVED***"
	if got := fmtTree(ws, fmtNodeParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After adjust\ngot  %q\nwant %q", got, want)
	***REMOVED***
	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func makeSection533Tree() *priorityWriteScheduler ***REMOVED***
	// Initial tree from RFC 7540 Section 5.3.3.
	// A,B,C,D,E,F = 1,2,3,4,5,6
	ws := defaultPriorityWriteScheduler()
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(3, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(4, OpenStreamOptions***REMOVED***PusherID: 3***REMOVED***)
	ws.OpenStream(5, OpenStreamOptions***REMOVED***PusherID: 3***REMOVED***)
	ws.OpenStream(6, OpenStreamOptions***REMOVED***PusherID: 4***REMOVED***)
	return ws
***REMOVED***

func TestPrioritySection533NonExclusive(t *testing.T) ***REMOVED***
	// Example from RFC 7540 Section 5.3.3.
	// A,B,C,D,E,F = 1,2,3,4,5,6
	ws := defaultPriorityWriteScheduler()
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(3, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(4, OpenStreamOptions***REMOVED***PusherID: 3***REMOVED***)
	ws.OpenStream(5, OpenStreamOptions***REMOVED***PusherID: 3***REMOVED***)
	ws.OpenStream(6, OpenStreamOptions***REMOVED***PusherID: 4***REMOVED***)
	ws.AdjustStream(1, PriorityParam***REMOVED***
		StreamDep: 4,
		Weight:    15,
		Exclusive: false,
	***REMOVED***)
	want := "1***REMOVED***parent:4***REMOVED*** 2***REMOVED***parent:1***REMOVED*** 3***REMOVED***parent:1***REMOVED*** 4***REMOVED***parent:0***REMOVED*** 5***REMOVED***parent:3***REMOVED*** 6***REMOVED***parent:4***REMOVED***"
	if got := fmtTree(ws, fmtNodeParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After adjust\ngot  %q\nwant %q", got, want)
	***REMOVED***
	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPrioritySection533Exclusive(t *testing.T) ***REMOVED***
	// Example from RFC 7540 Section 5.3.3.
	// A,B,C,D,E,F = 1,2,3,4,5,6
	ws := defaultPriorityWriteScheduler()
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(3, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(4, OpenStreamOptions***REMOVED***PusherID: 3***REMOVED***)
	ws.OpenStream(5, OpenStreamOptions***REMOVED***PusherID: 3***REMOVED***)
	ws.OpenStream(6, OpenStreamOptions***REMOVED***PusherID: 4***REMOVED***)
	ws.AdjustStream(1, PriorityParam***REMOVED***
		StreamDep: 4,
		Weight:    15,
		Exclusive: true,
	***REMOVED***)
	want := "1***REMOVED***parent:4***REMOVED*** 2***REMOVED***parent:1***REMOVED*** 3***REMOVED***parent:1***REMOVED*** 4***REMOVED***parent:0***REMOVED*** 5***REMOVED***parent:3***REMOVED*** 6***REMOVED***parent:1***REMOVED***"
	if got := fmtTree(ws, fmtNodeParentSkipRoot); got != want ***REMOVED***
		t.Errorf("After adjust\ngot  %q\nwant %q", got, want)
	***REMOVED***
	if err := checkPriorityWellFormed(ws); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func checkPopAll(ws WriteScheduler, order []uint32) error ***REMOVED***
	for k, id := range order ***REMOVED***
		wr, ok := ws.Pop()
		if !ok ***REMOVED***
			return fmt.Errorf("Pop[%d]: got ok=false, want %d (order=%v)", k, id, order)
		***REMOVED***
		if got := wr.StreamID(); got != id ***REMOVED***
			return fmt.Errorf("Pop[%d]: got %v, want %d (order=%v)", k, got, id, order)
		***REMOVED***
	***REMOVED***
	wr, ok := ws.Pop()
	if ok ***REMOVED***
		return fmt.Errorf("Pop[%d]: got %v, want ok=false (order=%v)", len(order), wr.StreamID(), order)
	***REMOVED***
	return nil
***REMOVED***

func TestPriorityPopFrom533Tree(t *testing.T) ***REMOVED***
	ws := makeSection533Tree()

	ws.Push(makeWriteHeadersRequest(3 /*C*/))
	ws.Push(makeWriteNonStreamRequest())
	ws.Push(makeWriteHeadersRequest(5 /*E*/))
	ws.Push(makeWriteHeadersRequest(1 /*A*/))
	t.Log("tree:", fmtTree(ws, fmtNodeParentSkipRoot))

	if err := checkPopAll(ws, []uint32***REMOVED***0 /*NonStream*/, 1, 3, 5***REMOVED***); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPriorityPopFromLinearTree(t *testing.T) ***REMOVED***
	ws := defaultPriorityWriteScheduler()
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)
	ws.OpenStream(3, OpenStreamOptions***REMOVED***PusherID: 2***REMOVED***)
	ws.OpenStream(4, OpenStreamOptions***REMOVED***PusherID: 3***REMOVED***)

	ws.Push(makeWriteHeadersRequest(3))
	ws.Push(makeWriteHeadersRequest(4))
	ws.Push(makeWriteHeadersRequest(1))
	ws.Push(makeWriteHeadersRequest(2))
	ws.Push(makeWriteNonStreamRequest())
	ws.Push(makeWriteNonStreamRequest())
	t.Log("tree:", fmtTree(ws, fmtNodeParentSkipRoot))

	if err := checkPopAll(ws, []uint32***REMOVED***0, 0 /*NonStreams*/, 1, 2, 3, 4***REMOVED***); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPriorityFlowControl(t *testing.T) ***REMOVED***
	ws := NewPriorityWriteScheduler(&PriorityWriteSchedulerConfig***REMOVED***ThrottleOutOfOrderWrites: false***REMOVED***)
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)

	sc := &serverConn***REMOVED***maxFrameSize: 16***REMOVED***
	st1 := &stream***REMOVED***id: 1, sc: sc***REMOVED***
	st2 := &stream***REMOVED***id: 2, sc: sc***REMOVED***

	ws.Push(FrameWriteRequest***REMOVED***&writeData***REMOVED***1, make([]byte, 16), false***REMOVED***, st1, nil***REMOVED***)
	ws.Push(FrameWriteRequest***REMOVED***&writeData***REMOVED***2, make([]byte, 16), false***REMOVED***, st2, nil***REMOVED***)
	ws.AdjustStream(2, PriorityParam***REMOVED***StreamDep: 1***REMOVED***)

	// No flow-control bytes available.
	if wr, ok := ws.Pop(); ok ***REMOVED***
		t.Fatalf("Pop(limited by flow control)=%v,true, want false", wr)
	***REMOVED***

	// Add enough flow-control bytes to write st2 in two Pop calls.
	// Should write data from st2 even though it's lower priority than st1.
	for i := 1; i <= 2; i++ ***REMOVED***
		st2.flow.add(8)
		wr, ok := ws.Pop()
		if !ok ***REMOVED***
			t.Fatalf("Pop(%d)=false, want true", i)
		***REMOVED***
		if got, want := wr.DataSize(), 8; got != want ***REMOVED***
			t.Fatalf("Pop(%d)=%d bytes, want %d bytes", i, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestPriorityThrottleOutOfOrderWrites(t *testing.T) ***REMOVED***
	ws := NewPriorityWriteScheduler(&PriorityWriteSchedulerConfig***REMOVED***ThrottleOutOfOrderWrites: true***REMOVED***)
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED***PusherID: 1***REMOVED***)

	sc := &serverConn***REMOVED***maxFrameSize: 4096***REMOVED***
	st1 := &stream***REMOVED***id: 1, sc: sc***REMOVED***
	st2 := &stream***REMOVED***id: 2, sc: sc***REMOVED***
	st1.flow.add(4096)
	st2.flow.add(4096)
	ws.Push(FrameWriteRequest***REMOVED***&writeData***REMOVED***2, make([]byte, 4096), false***REMOVED***, st2, nil***REMOVED***)
	ws.AdjustStream(2, PriorityParam***REMOVED***StreamDep: 1***REMOVED***)

	// We have enough flow-control bytes to write st2 in a single Pop call.
	// However, due to out-of-order write throttling, the first call should
	// only write 1KB.
	wr, ok := ws.Pop()
	if !ok ***REMOVED***
		t.Fatalf("Pop(st2.first)=false, want true")
	***REMOVED***
	if got, want := wr.StreamID(), uint32(2); got != want ***REMOVED***
		t.Fatalf("Pop(st2.first)=stream %d, want stream %d", got, want)
	***REMOVED***
	if got, want := wr.DataSize(), 1024; got != want ***REMOVED***
		t.Fatalf("Pop(st2.first)=%d bytes, want %d bytes", got, want)
	***REMOVED***

	// Now add data on st1. This should take precedence.
	ws.Push(FrameWriteRequest***REMOVED***&writeData***REMOVED***1, make([]byte, 4096), false***REMOVED***, st1, nil***REMOVED***)
	wr, ok = ws.Pop()
	if !ok ***REMOVED***
		t.Fatalf("Pop(st1)=false, want true")
	***REMOVED***
	if got, want := wr.StreamID(), uint32(1); got != want ***REMOVED***
		t.Fatalf("Pop(st1)=stream %d, want stream %d", got, want)
	***REMOVED***
	if got, want := wr.DataSize(), 4096; got != want ***REMOVED***
		t.Fatalf("Pop(st1)=%d bytes, want %d bytes", got, want)
	***REMOVED***

	// Should go back to writing 1KB from st2.
	wr, ok = ws.Pop()
	if !ok ***REMOVED***
		t.Fatalf("Pop(st2.last)=false, want true")
	***REMOVED***
	if got, want := wr.StreamID(), uint32(2); got != want ***REMOVED***
		t.Fatalf("Pop(st2.last)=stream %d, want stream %d", got, want)
	***REMOVED***
	if got, want := wr.DataSize(), 1024; got != want ***REMOVED***
		t.Fatalf("Pop(st2.last)=%d bytes, want %d bytes", got, want)
	***REMOVED***
***REMOVED***

func TestPriorityWeights(t *testing.T) ***REMOVED***
	ws := defaultPriorityWriteScheduler()
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.OpenStream(2, OpenStreamOptions***REMOVED******REMOVED***)

	sc := &serverConn***REMOVED***maxFrameSize: 8***REMOVED***
	st1 := &stream***REMOVED***id: 1, sc: sc***REMOVED***
	st2 := &stream***REMOVED***id: 2, sc: sc***REMOVED***
	st1.flow.add(40)
	st2.flow.add(40)

	ws.Push(FrameWriteRequest***REMOVED***&writeData***REMOVED***1, make([]byte, 40), false***REMOVED***, st1, nil***REMOVED***)
	ws.Push(FrameWriteRequest***REMOVED***&writeData***REMOVED***2, make([]byte, 40), false***REMOVED***, st2, nil***REMOVED***)
	ws.AdjustStream(1, PriorityParam***REMOVED***StreamDep: 0, Weight: 34***REMOVED***)
	ws.AdjustStream(2, PriorityParam***REMOVED***StreamDep: 0, Weight: 9***REMOVED***)

	// st1 gets 3.5x the bandwidth of st2 (3.5 = (34+1)/(9+1)).
	// The maximum frame size is 8 bytes. The write sequence should be:
	//   st1, total bytes so far is (st1=8,  st=0)
	//   st2, total bytes so far is (st1=8,  st=8)
	//   st1, total bytes so far is (st1=16, st=8)
	//   st1, total bytes so far is (st1=24, st=8)   // 3x bandwidth
	//   st1, total bytes so far is (st1=32, st=8)   // 4x bandwidth
	//   st2, total bytes so far is (st1=32, st=16)  // 2x bandwidth
	//   st1, total bytes so far is (st1=40, st=16)
	//   st2, total bytes so far is (st1=40, st=24)
	//   st2, total bytes so far is (st1=40, st=32)
	//   st2, total bytes so far is (st1=40, st=40)
	if err := checkPopAll(ws, []uint32***REMOVED***1, 2, 1, 1, 1, 2, 1, 2, 2, 2***REMOVED***); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestPriorityRstStreamOnNonOpenStreams(t *testing.T) ***REMOVED***
	ws := NewPriorityWriteScheduler(&PriorityWriteSchedulerConfig***REMOVED***
		MaxClosedNodesInTree: 0,
		MaxIdleNodesInTree:   0,
	***REMOVED***)
	ws.OpenStream(1, OpenStreamOptions***REMOVED******REMOVED***)
	ws.CloseStream(1)
	ws.Push(FrameWriteRequest***REMOVED***write: streamError(1, ErrCodeProtocol)***REMOVED***)
	ws.Push(FrameWriteRequest***REMOVED***write: streamError(2, ErrCodeProtocol)***REMOVED***)

	if err := checkPopAll(ws, []uint32***REMOVED***1, 2***REMOVED***); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***
