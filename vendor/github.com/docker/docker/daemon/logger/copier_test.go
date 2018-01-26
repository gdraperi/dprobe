package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type TestLoggerJSON struct ***REMOVED***
	*json.Encoder
	mu    sync.Mutex
	delay time.Duration
***REMOVED***

func (l *TestLoggerJSON) Log(m *Message) error ***REMOVED***
	if l.delay > 0 ***REMOVED***
		time.Sleep(l.delay)
	***REMOVED***
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.Encode(m)
***REMOVED***

func (l *TestLoggerJSON) Close() error ***REMOVED*** return nil ***REMOVED***

func (l *TestLoggerJSON) Name() string ***REMOVED*** return "json" ***REMOVED***

type TestSizedLoggerJSON struct ***REMOVED***
	*json.Encoder
	mu sync.Mutex
***REMOVED***

func (l *TestSizedLoggerJSON) Log(m *Message) error ***REMOVED***
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.Encode(m)
***REMOVED***

func (*TestSizedLoggerJSON) Close() error ***REMOVED*** return nil ***REMOVED***

func (*TestSizedLoggerJSON) Name() string ***REMOVED*** return "sized-json" ***REMOVED***

func (*TestSizedLoggerJSON) BufSize() int ***REMOVED***
	return 32 * 1024
***REMOVED***

func TestCopier(t *testing.T) ***REMOVED***
	stdoutLine := "Line that thinks that it is log line from docker stdout"
	stderrLine := "Line that thinks that it is log line from docker stderr"
	stdoutTrailingLine := "stdout trailing line"
	stderrTrailingLine := "stderr trailing line"

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	for i := 0; i < 30; i++ ***REMOVED***
		if _, err := stdout.WriteString(stdoutLine + "\n"); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if _, err := stderr.WriteString(stderrLine + "\n"); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	// Test remaining lines without line-endings
	if _, err := stdout.WriteString(stdoutTrailingLine); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := stderr.WriteString(stderrTrailingLine); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	var jsonBuf bytes.Buffer

	jsonLog := &TestLoggerJSON***REMOVED***Encoder: json.NewEncoder(&jsonBuf)***REMOVED***

	c := NewCopier(
		map[string]io.Reader***REMOVED***
			"stdout": &stdout,
			"stderr": &stderr,
		***REMOVED***,
		jsonLog)
	c.Run()
	wait := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		c.Wait()
		close(wait)
	***REMOVED***()
	select ***REMOVED***
	case <-time.After(1 * time.Second):
		t.Fatal("Copier failed to do its work in 1 second")
	case <-wait:
	***REMOVED***
	dec := json.NewDecoder(&jsonBuf)
	for ***REMOVED***
		var msg Message
		if err := dec.Decode(&msg); err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if msg.Source != "stdout" && msg.Source != "stderr" ***REMOVED***
			t.Fatalf("Wrong Source: %q, should be %q or %q", msg.Source, "stdout", "stderr")
		***REMOVED***
		if msg.Source == "stdout" ***REMOVED***
			if string(msg.Line) != stdoutLine && string(msg.Line) != stdoutTrailingLine ***REMOVED***
				t.Fatalf("Wrong Line: %q, expected %q or %q", msg.Line, stdoutLine, stdoutTrailingLine)
			***REMOVED***
		***REMOVED***
		if msg.Source == "stderr" ***REMOVED***
			if string(msg.Line) != stderrLine && string(msg.Line) != stderrTrailingLine ***REMOVED***
				t.Fatalf("Wrong Line: %q, expected %q or %q", msg.Line, stderrLine, stderrTrailingLine)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestCopierLongLines tests long lines without line breaks
func TestCopierLongLines(t *testing.T) ***REMOVED***
	// Long lines (should be split at "defaultBufSize")
	stdoutLongLine := strings.Repeat("a", defaultBufSize)
	stderrLongLine := strings.Repeat("b", defaultBufSize)
	stdoutTrailingLine := "stdout trailing line"
	stderrTrailingLine := "stderr trailing line"

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	for i := 0; i < 3; i++ ***REMOVED***
		if _, err := stdout.WriteString(stdoutLongLine); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if _, err := stderr.WriteString(stderrLongLine); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	if _, err := stdout.WriteString(stdoutTrailingLine); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := stderr.WriteString(stderrTrailingLine); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	var jsonBuf bytes.Buffer

	jsonLog := &TestLoggerJSON***REMOVED***Encoder: json.NewEncoder(&jsonBuf)***REMOVED***

	c := NewCopier(
		map[string]io.Reader***REMOVED***
			"stdout": &stdout,
			"stderr": &stderr,
		***REMOVED***,
		jsonLog)
	c.Run()
	wait := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		c.Wait()
		close(wait)
	***REMOVED***()
	select ***REMOVED***
	case <-time.After(1 * time.Second):
		t.Fatal("Copier failed to do its work in 1 second")
	case <-wait:
	***REMOVED***
	dec := json.NewDecoder(&jsonBuf)
	for ***REMOVED***
		var msg Message
		if err := dec.Decode(&msg); err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if msg.Source != "stdout" && msg.Source != "stderr" ***REMOVED***
			t.Fatalf("Wrong Source: %q, should be %q or %q", msg.Source, "stdout", "stderr")
		***REMOVED***
		if msg.Source == "stdout" ***REMOVED***
			if string(msg.Line) != stdoutLongLine && string(msg.Line) != stdoutTrailingLine ***REMOVED***
				t.Fatalf("Wrong Line: %q, expected 'stdoutLongLine' or 'stdoutTrailingLine'", msg.Line)
			***REMOVED***
		***REMOVED***
		if msg.Source == "stderr" ***REMOVED***
			if string(msg.Line) != stderrLongLine && string(msg.Line) != stderrTrailingLine ***REMOVED***
				t.Fatalf("Wrong Line: %q, expected 'stderrLongLine' or 'stderrTrailingLine'", msg.Line)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCopierSlow(t *testing.T) ***REMOVED***
	stdoutLine := "Line that thinks that it is log line from docker stdout"
	var stdout bytes.Buffer
	for i := 0; i < 30; i++ ***REMOVED***
		if _, err := stdout.WriteString(stdoutLine + "\n"); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	var jsonBuf bytes.Buffer
	//encoder := &encodeCloser***REMOVED***Encoder: json.NewEncoder(&jsonBuf)***REMOVED***
	jsonLog := &TestLoggerJSON***REMOVED***Encoder: json.NewEncoder(&jsonBuf), delay: 100 * time.Millisecond***REMOVED***

	c := NewCopier(map[string]io.Reader***REMOVED***"stdout": &stdout***REMOVED***, jsonLog)
	c.Run()
	wait := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		c.Wait()
		close(wait)
	***REMOVED***()
	<-time.After(150 * time.Millisecond)
	c.Close()
	select ***REMOVED***
	case <-time.After(200 * time.Millisecond):
		t.Fatal("failed to exit in time after the copier is closed")
	case <-wait:
	***REMOVED***
***REMOVED***

func TestCopierWithSized(t *testing.T) ***REMOVED***
	var jsonBuf bytes.Buffer
	expectedMsgs := 2
	sizedLogger := &TestSizedLoggerJSON***REMOVED***Encoder: json.NewEncoder(&jsonBuf)***REMOVED***
	logbuf := bytes.NewBufferString(strings.Repeat(".", sizedLogger.BufSize()*expectedMsgs))
	c := NewCopier(map[string]io.Reader***REMOVED***"stdout": logbuf***REMOVED***, sizedLogger)

	c.Run()
	// Wait for Copier to finish writing to the buffered logger.
	c.Wait()
	c.Close()

	recvdMsgs := 0
	dec := json.NewDecoder(&jsonBuf)
	for ***REMOVED***
		var msg Message
		if err := dec.Decode(&msg); err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if msg.Source != "stdout" ***REMOVED***
			t.Fatalf("Wrong Source: %q, should be %q", msg.Source, "stdout")
		***REMOVED***
		if len(msg.Line) != sizedLogger.BufSize() ***REMOVED***
			t.Fatalf("Line was not of expected max length %d, was %d", sizedLogger.BufSize(), len(msg.Line))
		***REMOVED***
		recvdMsgs++
	***REMOVED***
	if recvdMsgs != expectedMsgs ***REMOVED***
		t.Fatalf("expected to receive %d messages, actually received %d", expectedMsgs, recvdMsgs)
	***REMOVED***
***REMOVED***

type BenchmarkLoggerDummy struct ***REMOVED***
***REMOVED***

func (l *BenchmarkLoggerDummy) Log(m *Message) error ***REMOVED*** PutMessage(m); return nil ***REMOVED***

func (l *BenchmarkLoggerDummy) Close() error ***REMOVED*** return nil ***REMOVED***

func (l *BenchmarkLoggerDummy) Name() string ***REMOVED*** return "dummy" ***REMOVED***

func BenchmarkCopier64(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<6)
***REMOVED***
func BenchmarkCopier128(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<7)
***REMOVED***
func BenchmarkCopier256(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<8)
***REMOVED***
func BenchmarkCopier512(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<9)
***REMOVED***
func BenchmarkCopier1K(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<10)
***REMOVED***
func BenchmarkCopier2K(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<11)
***REMOVED***
func BenchmarkCopier4K(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<12)
***REMOVED***
func BenchmarkCopier8K(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<13)
***REMOVED***
func BenchmarkCopier16K(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<14)
***REMOVED***
func BenchmarkCopier32K(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<15)
***REMOVED***
func BenchmarkCopier64K(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<16)
***REMOVED***
func BenchmarkCopier128K(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<17)
***REMOVED***
func BenchmarkCopier256K(b *testing.B) ***REMOVED***
	benchmarkCopier(b, 1<<18)
***REMOVED***

func piped(b *testing.B, iterations int, delay time.Duration, buf []byte) io.Reader ***REMOVED***
	r, w, err := os.Pipe()
	if err != nil ***REMOVED***
		b.Fatal(err)
		return nil
	***REMOVED***
	go func() ***REMOVED***
		for i := 0; i < iterations; i++ ***REMOVED***
			time.Sleep(delay)
			if n, err := w.Write(buf); err != nil || n != len(buf) ***REMOVED***
				if err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
				b.Fatal(fmt.Errorf("short write"))
			***REMOVED***
		***REMOVED***
		w.Close()
	***REMOVED***()
	return r
***REMOVED***

func benchmarkCopier(b *testing.B, length int) ***REMOVED***
	b.StopTimer()
	buf := []byte***REMOVED***'A'***REMOVED***
	for len(buf) < length ***REMOVED***
		buf = append(buf, buf...)
	***REMOVED***
	buf = append(buf[:length-1], []byte***REMOVED***'\n'***REMOVED***...)
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		c := NewCopier(
			map[string]io.Reader***REMOVED***
				"buffer": piped(b, 10, time.Nanosecond, buf),
			***REMOVED***,
			&BenchmarkLoggerDummy***REMOVED******REMOVED***)
		c.Run()
		c.Wait()
		c.Close()
	***REMOVED***
***REMOVED***
