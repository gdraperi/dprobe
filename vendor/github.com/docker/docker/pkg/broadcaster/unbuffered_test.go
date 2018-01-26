package broadcaster

import (
	"bytes"
	"errors"
	"strings"

	"testing"
)

type dummyWriter struct ***REMOVED***
	buffer      bytes.Buffer
	failOnWrite bool
***REMOVED***

func (dw *dummyWriter) Write(p []byte) (n int, err error) ***REMOVED***
	if dw.failOnWrite ***REMOVED***
		return 0, errors.New("Fake fail")
	***REMOVED***
	return dw.buffer.Write(p)
***REMOVED***

func (dw *dummyWriter) String() string ***REMOVED***
	return dw.buffer.String()
***REMOVED***

func (dw *dummyWriter) Close() error ***REMOVED***
	return nil
***REMOVED***

func TestUnbuffered(t *testing.T) ***REMOVED***
	writer := new(Unbuffered)

	// Test 1: Both bufferA and bufferB should contain "foo"
	bufferA := &dummyWriter***REMOVED******REMOVED***
	writer.Add(bufferA)
	bufferB := &dummyWriter***REMOVED******REMOVED***
	writer.Add(bufferB)
	writer.Write([]byte("foo"))

	if bufferA.String() != "foo" ***REMOVED***
		t.Errorf("Buffer contains %v", bufferA.String())
	***REMOVED***

	if bufferB.String() != "foo" ***REMOVED***
		t.Errorf("Buffer contains %v", bufferB.String())
	***REMOVED***

	// Test2: bufferA and bufferB should contain "foobar",
	// while bufferC should only contain "bar"
	bufferC := &dummyWriter***REMOVED******REMOVED***
	writer.Add(bufferC)
	writer.Write([]byte("bar"))

	if bufferA.String() != "foobar" ***REMOVED***
		t.Errorf("Buffer contains %v", bufferA.String())
	***REMOVED***

	if bufferB.String() != "foobar" ***REMOVED***
		t.Errorf("Buffer contains %v", bufferB.String())
	***REMOVED***

	if bufferC.String() != "bar" ***REMOVED***
		t.Errorf("Buffer contains %v", bufferC.String())
	***REMOVED***

	// Test3: Test eviction on failure
	bufferA.failOnWrite = true
	writer.Write([]byte("fail"))
	if bufferA.String() != "foobar" ***REMOVED***
		t.Errorf("Buffer contains %v", bufferA.String())
	***REMOVED***
	if bufferC.String() != "barfail" ***REMOVED***
		t.Errorf("Buffer contains %v", bufferC.String())
	***REMOVED***
	// Even though we reset the flag, no more writes should go in there
	bufferA.failOnWrite = false
	writer.Write([]byte("test"))
	if bufferA.String() != "foobar" ***REMOVED***
		t.Errorf("Buffer contains %v", bufferA.String())
	***REMOVED***
	if bufferC.String() != "barfailtest" ***REMOVED***
		t.Errorf("Buffer contains %v", bufferC.String())
	***REMOVED***

	// Test4: Test eviction on multiple simultaneous failures
	bufferB.failOnWrite = true
	bufferC.failOnWrite = true
	bufferD := &dummyWriter***REMOVED******REMOVED***
	writer.Add(bufferD)
	writer.Write([]byte("yo"))
	writer.Write([]byte("ink"))
	if strings.Contains(bufferB.String(), "yoink") ***REMOVED***
		t.Errorf("bufferB received write. contents: %q", bufferB)
	***REMOVED***
	if strings.Contains(bufferC.String(), "yoink") ***REMOVED***
		t.Errorf("bufferC received write. contents: %q", bufferC)
	***REMOVED***
	if g, w := bufferD.String(), "yoink"; g != w ***REMOVED***
		t.Errorf("bufferD = %q, want %q", g, w)
	***REMOVED***

	writer.Clean()
***REMOVED***

type devNullCloser int

func (d devNullCloser) Close() error ***REMOVED***
	return nil
***REMOVED***

func (d devNullCloser) Write(buf []byte) (int, error) ***REMOVED***
	return len(buf), nil
***REMOVED***

// This test checks for races. It is only useful when run with the race detector.
func TestRaceUnbuffered(t *testing.T) ***REMOVED***
	writer := new(Unbuffered)
	c := make(chan bool)
	go func() ***REMOVED***
		writer.Add(devNullCloser(0))
		c <- true
	***REMOVED***()
	writer.Write([]byte("hello"))
	<-c
***REMOVED***

func BenchmarkUnbuffered(b *testing.B) ***REMOVED***
	writer := new(Unbuffered)
	setUpWriter := func() ***REMOVED***
		for i := 0; i < 100; i++ ***REMOVED***
			writer.Add(devNullCloser(0))
			writer.Add(devNullCloser(0))
			writer.Add(devNullCloser(0))
		***REMOVED***
	***REMOVED***
	testLine := "Line that thinks that it is log line from docker"
	var buf bytes.Buffer
	for i := 0; i < 100; i++ ***REMOVED***
		buf.Write([]byte(testLine + "\n"))
	***REMOVED***
	// line without eol
	buf.Write([]byte(testLine))
	testText := buf.Bytes()
	b.SetBytes(int64(5 * len(testText)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		b.StopTimer()
		setUpWriter()
		b.StartTimer()

		for j := 0; j < 5; j++ ***REMOVED***
			if _, err := writer.Write(testText); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***

		b.StopTimer()
		writer.Clean()
		b.StartTimer()
	***REMOVED***
***REMOVED***
