package logger

import (
	"context"
	"strconv"
	"testing"
	"time"
)

type mockLogger struct***REMOVED*** c chan *Message ***REMOVED***

func (l *mockLogger) Log(msg *Message) error ***REMOVED***
	l.c <- msg
	return nil
***REMOVED***

func (l *mockLogger) Name() string ***REMOVED***
	return "mock"
***REMOVED***

func (l *mockLogger) Close() error ***REMOVED***
	return nil
***REMOVED***

func TestRingLogger(t *testing.T) ***REMOVED***
	mockLog := &mockLogger***REMOVED***make(chan *Message)***REMOVED*** // no buffer on this channel
	ring := newRingLogger(mockLog, Info***REMOVED******REMOVED***, 1)
	defer ring.setClosed()

	// this should never block
	ring.Log(&Message***REMOVED***Line: []byte("1")***REMOVED***)
	ring.Log(&Message***REMOVED***Line: []byte("2")***REMOVED***)
	ring.Log(&Message***REMOVED***Line: []byte("3")***REMOVED***)

	select ***REMOVED***
	case msg := <-mockLog.c:
		if string(msg.Line) != "1" ***REMOVED***
			t.Fatalf("got unexpected msg: %q", string(msg.Line))
		***REMOVED***
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout reading log message")
	***REMOVED***

	select ***REMOVED***
	case msg := <-mockLog.c:
		t.Fatalf("expected no more messages in the queue, got: %q", string(msg.Line))
	default:
	***REMOVED***
***REMOVED***

func TestRingCap(t *testing.T) ***REMOVED***
	r := newRing(5)
	for i := 0; i < 10; i++ ***REMOVED***
		// queue messages with "0" to "10"
		// the "5" to "10" messages should be dropped since we only allow 5 bytes in the buffer
		if err := r.Enqueue(&Message***REMOVED***Line: []byte(strconv.Itoa(i))***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	// should have messages in the queue for "5" to "10"
	for i := 0; i < 5; i++ ***REMOVED***
		m, err := r.Dequeue()
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if string(m.Line) != strconv.Itoa(i) ***REMOVED***
			t.Fatalf("got unexpected message for iter %d: %s", i, string(m.Line))
		***REMOVED***
	***REMOVED***

	// queue a message that's bigger than the buffer cap
	if err := r.Enqueue(&Message***REMOVED***Line: []byte("hello world")***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// queue another message that's bigger than the buffer cap
	if err := r.Enqueue(&Message***REMOVED***Line: []byte("eat a banana")***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	m, err := r.Dequeue()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(m.Line) != "hello world" ***REMOVED***
		t.Fatalf("got unexpected message: %s", string(m.Line))
	***REMOVED***
	if len(r.queue) != 0 ***REMOVED***
		t.Fatalf("expected queue to be empty, got: %d", len(r.queue))
	***REMOVED***
***REMOVED***

func TestRingClose(t *testing.T) ***REMOVED***
	r := newRing(1)
	if err := r.Enqueue(&Message***REMOVED***Line: []byte("hello")***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	r.Close()
	if err := r.Enqueue(&Message***REMOVED******REMOVED***); err != errClosed ***REMOVED***
		t.Fatalf("expected errClosed, got: %v", err)
	***REMOVED***
	if len(r.queue) != 1 ***REMOVED***
		t.Fatal("expected empty queue")
	***REMOVED***
	if m, err := r.Dequeue(); err == nil || m != nil ***REMOVED***
		t.Fatal("expected err on Dequeue after close")
	***REMOVED***

	ls := r.Drain()
	if len(ls) != 1 ***REMOVED***
		t.Fatalf("expected one message: %v", ls)
	***REMOVED***
	if string(ls[0].Line) != "hello" ***REMOVED***
		t.Fatalf("got unexpected message: %s", string(ls[0].Line))
	***REMOVED***
***REMOVED***

func TestRingDrain(t *testing.T) ***REMOVED***
	r := newRing(5)
	for i := 0; i < 5; i++ ***REMOVED***
		if err := r.Enqueue(&Message***REMOVED***Line: []byte(strconv.Itoa(i))***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	ls := r.Drain()
	if len(ls) != 5 ***REMOVED***
		t.Fatal("got unexpected length after drain")
	***REMOVED***

	for i := 0; i < 5; i++ ***REMOVED***
		if string(ls[i].Line) != strconv.Itoa(i) ***REMOVED***
			t.Fatalf("got unexpected message at position %d: %s", i, string(ls[i].Line))
		***REMOVED***
	***REMOVED***
	if r.sizeBytes != 0 ***REMOVED***
		t.Fatalf("expected buffer size to be 0 after drain, got: %d", r.sizeBytes)
	***REMOVED***

	ls = r.Drain()
	if len(ls) != 0 ***REMOVED***
		t.Fatalf("expected 0 messages on 2nd drain: %v", ls)
	***REMOVED***

***REMOVED***

type nopLogger struct***REMOVED******REMOVED***

func (nopLogger) Name() string       ***REMOVED*** return "nopLogger" ***REMOVED***
func (nopLogger) Close() error       ***REMOVED*** return nil ***REMOVED***
func (nopLogger) Log(*Message) error ***REMOVED*** return nil ***REMOVED***

func BenchmarkRingLoggerThroughputNoReceiver(b *testing.B) ***REMOVED***
	mockLog := &mockLogger***REMOVED***make(chan *Message)***REMOVED***
	defer mockLog.Close()
	l := NewRingLogger(mockLog, Info***REMOVED******REMOVED***, -1)
	msg := &Message***REMOVED***Line: []byte("hello humans and everyone else!")***REMOVED***
	b.SetBytes(int64(len(msg.Line)))

	for i := 0; i < b.N; i++ ***REMOVED***
		if err := l.Log(msg); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkRingLoggerThroughputWithReceiverDelay0(b *testing.B) ***REMOVED***
	l := NewRingLogger(nopLogger***REMOVED******REMOVED***, Info***REMOVED******REMOVED***, -1)
	msg := &Message***REMOVED***Line: []byte("hello humans and everyone else!")***REMOVED***
	b.SetBytes(int64(len(msg.Line)))

	for i := 0; i < b.N; i++ ***REMOVED***
		if err := l.Log(msg); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func consumeWithDelay(delay time.Duration, c <-chan *Message) (cancel func()) ***REMOVED***
	started := make(chan struct***REMOVED******REMOVED***)
	ctx, cancel := context.WithCancel(context.Background())
	go func() ***REMOVED***
		close(started)
		ticker := time.NewTicker(delay)
		for range ticker.C ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-c:
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	<-started
	return cancel
***REMOVED***

func BenchmarkRingLoggerThroughputConsumeDelay1(b *testing.B) ***REMOVED***
	mockLog := &mockLogger***REMOVED***make(chan *Message)***REMOVED***
	defer mockLog.Close()
	l := NewRingLogger(mockLog, Info***REMOVED******REMOVED***, -1)
	msg := &Message***REMOVED***Line: []byte("hello humans and everyone else!")***REMOVED***
	b.SetBytes(int64(len(msg.Line)))

	cancel := consumeWithDelay(1*time.Millisecond, mockLog.c)
	defer cancel()

	for i := 0; i < b.N; i++ ***REMOVED***
		if err := l.Log(msg); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkRingLoggerThroughputConsumeDelay10(b *testing.B) ***REMOVED***
	mockLog := &mockLogger***REMOVED***make(chan *Message)***REMOVED***
	defer mockLog.Close()
	l := NewRingLogger(mockLog, Info***REMOVED******REMOVED***, -1)
	msg := &Message***REMOVED***Line: []byte("hello humans and everyone else!")***REMOVED***
	b.SetBytes(int64(len(msg.Line)))

	cancel := consumeWithDelay(10*time.Millisecond, mockLog.c)
	defer cancel()

	for i := 0; i < b.N; i++ ***REMOVED***
		if err := l.Log(msg); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkRingLoggerThroughputConsumeDelay50(b *testing.B) ***REMOVED***
	mockLog := &mockLogger***REMOVED***make(chan *Message)***REMOVED***
	defer mockLog.Close()
	l := NewRingLogger(mockLog, Info***REMOVED******REMOVED***, -1)
	msg := &Message***REMOVED***Line: []byte("hello humans and everyone else!")***REMOVED***
	b.SetBytes(int64(len(msg.Line)))

	cancel := consumeWithDelay(50*time.Millisecond, mockLog.c)
	defer cancel()

	for i := 0; i < b.N; i++ ***REMOVED***
		if err := l.Log(msg); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkRingLoggerThroughputConsumeDelay100(b *testing.B) ***REMOVED***
	mockLog := &mockLogger***REMOVED***make(chan *Message)***REMOVED***
	defer mockLog.Close()
	l := NewRingLogger(mockLog, Info***REMOVED******REMOVED***, -1)
	msg := &Message***REMOVED***Line: []byte("hello humans and everyone else!")***REMOVED***
	b.SetBytes(int64(len(msg.Line)))

	cancel := consumeWithDelay(100*time.Millisecond, mockLog.c)
	defer cancel()

	for i := 0; i < b.N; i++ ***REMOVED***
		if err := l.Log(msg); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkRingLoggerThroughputConsumeDelay300(b *testing.B) ***REMOVED***
	mockLog := &mockLogger***REMOVED***make(chan *Message)***REMOVED***
	defer mockLog.Close()
	l := NewRingLogger(mockLog, Info***REMOVED******REMOVED***, -1)
	msg := &Message***REMOVED***Line: []byte("hello humans and everyone else!")***REMOVED***
	b.SetBytes(int64(len(msg.Line)))

	cancel := consumeWithDelay(300*time.Millisecond, mockLog.c)
	defer cancel()

	for i := 0; i < b.N; i++ ***REMOVED***
		if err := l.Log(msg); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkRingLoggerThroughputConsumeDelay500(b *testing.B) ***REMOVED***
	mockLog := &mockLogger***REMOVED***make(chan *Message)***REMOVED***
	defer mockLog.Close()
	l := NewRingLogger(mockLog, Info***REMOVED******REMOVED***, -1)
	msg := &Message***REMOVED***Line: []byte("hello humans and everyone else!")***REMOVED***
	b.SetBytes(int64(len(msg.Line)))

	cancel := consumeWithDelay(500*time.Millisecond, mockLog.c)
	defer cancel()

	for i := 0; i < b.N; i++ ***REMOVED***
		if err := l.Log(msg); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
