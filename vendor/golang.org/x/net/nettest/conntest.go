// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package nettest provides utilities for network testing.
package nettest

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"testing"
	"time"
)

var (
	aLongTimeAgo = time.Unix(233431200, 0)
	neverTimeout = time.Time***REMOVED******REMOVED***
)

// MakePipe creates a connection between two endpoints and returns the pair
// as c1 and c2, such that anything written to c1 is read by c2 and vice-versa.
// The stop function closes all resources, including c1, c2, and the underlying
// net.Listener (if there is one), and should not be nil.
type MakePipe func() (c1, c2 net.Conn, stop func(), err error)

// TestConn tests that a net.Conn implementation properly satisfies the interface.
// The tests should not produce any false positives, but may experience
// false negatives. Thus, some issues may only be detected when the test is
// run multiple times. For maximal effectiveness, run the tests under the
// race detector.
func TestConn(t *testing.T, mp MakePipe) ***REMOVED***
	testConn(t, mp)
***REMOVED***

type connTester func(t *testing.T, c1, c2 net.Conn)

func timeoutWrapper(t *testing.T, mp MakePipe, f connTester) ***REMOVED***
	c1, c2, stop, err := mp()
	if err != nil ***REMOVED***
		t.Fatalf("unable to make pipe: %v", err)
	***REMOVED***
	var once sync.Once
	defer once.Do(func() ***REMOVED*** stop() ***REMOVED***)
	timer := time.AfterFunc(time.Minute, func() ***REMOVED***
		once.Do(func() ***REMOVED***
			t.Error("test timed out; terminating pipe")
			stop()
		***REMOVED***)
	***REMOVED***)
	defer timer.Stop()
	f(t, c1, c2)
***REMOVED***

// testBasicIO tests that the data sent on c1 is properly received on c2.
func testBasicIO(t *testing.T, c1, c2 net.Conn) ***REMOVED***
	want := make([]byte, 1<<20)
	rand.New(rand.NewSource(0)).Read(want)

	dataCh := make(chan []byte)
	go func() ***REMOVED***
		rd := bytes.NewReader(want)
		if err := chunkedCopy(c1, rd); err != nil ***REMOVED***
			t.Errorf("unexpected c1.Write error: %v", err)
		***REMOVED***
		if err := c1.Close(); err != nil ***REMOVED***
			t.Errorf("unexpected c1.Close error: %v", err)
		***REMOVED***
	***REMOVED***()

	go func() ***REMOVED***
		wr := new(bytes.Buffer)
		if err := chunkedCopy(wr, c2); err != nil ***REMOVED***
			t.Errorf("unexpected c2.Read error: %v", err)
		***REMOVED***
		if err := c2.Close(); err != nil ***REMOVED***
			t.Errorf("unexpected c2.Close error: %v", err)
		***REMOVED***
		dataCh <- wr.Bytes()
	***REMOVED***()

	if got := <-dataCh; !bytes.Equal(got, want) ***REMOVED***
		t.Errorf("transmitted data differs")
	***REMOVED***
***REMOVED***

// testPingPong tests that the two endpoints can synchronously send data to
// each other in a typical request-response pattern.
func testPingPong(t *testing.T, c1, c2 net.Conn) ***REMOVED***
	var wg sync.WaitGroup
	defer wg.Wait()

	pingPonger := func(c net.Conn) ***REMOVED***
		defer wg.Done()
		buf := make([]byte, 8)
		var prev uint64
		for ***REMOVED***
			if _, err := io.ReadFull(c, buf); err != nil ***REMOVED***
				if err == io.EOF ***REMOVED***
					break
				***REMOVED***
				t.Errorf("unexpected Read error: %v", err)
			***REMOVED***

			v := binary.LittleEndian.Uint64(buf)
			binary.LittleEndian.PutUint64(buf, v+1)
			if prev != 0 && prev+2 != v ***REMOVED***
				t.Errorf("mismatching value: got %d, want %d", v, prev+2)
			***REMOVED***
			prev = v
			if v == 1000 ***REMOVED***
				break
			***REMOVED***

			if _, err := c.Write(buf); err != nil ***REMOVED***
				t.Errorf("unexpected Write error: %v", err)
				break
			***REMOVED***
		***REMOVED***
		if err := c.Close(); err != nil ***REMOVED***
			t.Errorf("unexpected Close error: %v", err)
		***REMOVED***
	***REMOVED***

	wg.Add(2)
	go pingPonger(c1)
	go pingPonger(c2)

	// Start off the chain reaction.
	if _, err := c1.Write(make([]byte, 8)); err != nil ***REMOVED***
		t.Errorf("unexpected c1.Write error: %v", err)
	***REMOVED***
***REMOVED***

// testRacyRead tests that it is safe to mutate the input Read buffer
// immediately after cancelation has occurred.
func testRacyRead(t *testing.T, c1, c2 net.Conn) ***REMOVED***
	go chunkedCopy(c2, rand.New(rand.NewSource(0)))

	var wg sync.WaitGroup
	defer wg.Wait()

	c1.SetReadDeadline(time.Now().Add(time.Millisecond))
	for i := 0; i < 10; i++ ***REMOVED***
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()

			b1 := make([]byte, 1024)
			b2 := make([]byte, 1024)
			for j := 0; j < 100; j++ ***REMOVED***
				_, err := c1.Read(b1)
				copy(b1, b2) // Mutate b1 to trigger potential race
				if err != nil ***REMOVED***
					checkForTimeoutError(t, err)
					c1.SetReadDeadline(time.Now().Add(time.Millisecond))
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***
***REMOVED***

// testRacyWrite tests that it is safe to mutate the input Write buffer
// immediately after cancelation has occurred.
func testRacyWrite(t *testing.T, c1, c2 net.Conn) ***REMOVED***
	go chunkedCopy(ioutil.Discard, c2)

	var wg sync.WaitGroup
	defer wg.Wait()

	c1.SetWriteDeadline(time.Now().Add(time.Millisecond))
	for i := 0; i < 10; i++ ***REMOVED***
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()

			b1 := make([]byte, 1024)
			b2 := make([]byte, 1024)
			for j := 0; j < 100; j++ ***REMOVED***
				_, err := c1.Write(b1)
				copy(b1, b2) // Mutate b1 to trigger potential race
				if err != nil ***REMOVED***
					checkForTimeoutError(t, err)
					c1.SetWriteDeadline(time.Now().Add(time.Millisecond))
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***
***REMOVED***

// testReadTimeout tests that Read timeouts do not affect Write.
func testReadTimeout(t *testing.T, c1, c2 net.Conn) ***REMOVED***
	go chunkedCopy(ioutil.Discard, c2)

	c1.SetReadDeadline(aLongTimeAgo)
	_, err := c1.Read(make([]byte, 1024))
	checkForTimeoutError(t, err)
	if _, err := c1.Write(make([]byte, 1024)); err != nil ***REMOVED***
		t.Errorf("unexpected Write error: %v", err)
	***REMOVED***
***REMOVED***

// testWriteTimeout tests that Write timeouts do not affect Read.
func testWriteTimeout(t *testing.T, c1, c2 net.Conn) ***REMOVED***
	go chunkedCopy(c2, rand.New(rand.NewSource(0)))

	c1.SetWriteDeadline(aLongTimeAgo)
	_, err := c1.Write(make([]byte, 1024))
	checkForTimeoutError(t, err)
	if _, err := c1.Read(make([]byte, 1024)); err != nil ***REMOVED***
		t.Errorf("unexpected Read error: %v", err)
	***REMOVED***
***REMOVED***

// testPastTimeout tests that a deadline set in the past immediately times out
// Read and Write requests.
func testPastTimeout(t *testing.T, c1, c2 net.Conn) ***REMOVED***
	go chunkedCopy(c2, c2)

	testRoundtrip(t, c1)

	c1.SetDeadline(aLongTimeAgo)
	n, err := c1.Write(make([]byte, 1024))
	if n != 0 ***REMOVED***
		t.Errorf("unexpected Write count: got %d, want 0", n)
	***REMOVED***
	checkForTimeoutError(t, err)
	n, err = c1.Read(make([]byte, 1024))
	if n != 0 ***REMOVED***
		t.Errorf("unexpected Read count: got %d, want 0", n)
	***REMOVED***
	checkForTimeoutError(t, err)

	testRoundtrip(t, c1)
***REMOVED***

// testPresentTimeout tests that a deadline set while there are pending
// Read and Write operations immediately times out those operations.
func testPresentTimeout(t *testing.T, c1, c2 net.Conn) ***REMOVED***
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(3)

	deadlineSet := make(chan bool, 1)
	go func() ***REMOVED***
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		deadlineSet <- true
		c1.SetReadDeadline(aLongTimeAgo)
		c1.SetWriteDeadline(aLongTimeAgo)
	***REMOVED***()
	go func() ***REMOVED***
		defer wg.Done()
		n, err := c1.Read(make([]byte, 1024))
		if n != 0 ***REMOVED***
			t.Errorf("unexpected Read count: got %d, want 0", n)
		***REMOVED***
		checkForTimeoutError(t, err)
		if len(deadlineSet) == 0 ***REMOVED***
			t.Error("Read timed out before deadline is set")
		***REMOVED***
	***REMOVED***()
	go func() ***REMOVED***
		defer wg.Done()
		var err error
		for err == nil ***REMOVED***
			_, err = c1.Write(make([]byte, 1024))
		***REMOVED***
		checkForTimeoutError(t, err)
		if len(deadlineSet) == 0 ***REMOVED***
			t.Error("Write timed out before deadline is set")
		***REMOVED***
	***REMOVED***()
***REMOVED***

// testFutureTimeout tests that a future deadline will eventually time out
// Read and Write operations.
func testFutureTimeout(t *testing.T, c1, c2 net.Conn) ***REMOVED***
	var wg sync.WaitGroup
	wg.Add(2)

	c1.SetDeadline(time.Now().Add(100 * time.Millisecond))
	go func() ***REMOVED***
		defer wg.Done()
		_, err := c1.Read(make([]byte, 1024))
		checkForTimeoutError(t, err)
	***REMOVED***()
	go func() ***REMOVED***
		defer wg.Done()
		var err error
		for err == nil ***REMOVED***
			_, err = c1.Write(make([]byte, 1024))
		***REMOVED***
		checkForTimeoutError(t, err)
	***REMOVED***()
	wg.Wait()

	go chunkedCopy(c2, c2)
	resyncConn(t, c1)
	testRoundtrip(t, c1)
***REMOVED***

// testCloseTimeout tests that calling Close immediately times out pending
// Read and Write operations.
func testCloseTimeout(t *testing.T, c1, c2 net.Conn) ***REMOVED***
	go chunkedCopy(c2, c2)

	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(3)

	// Test for cancelation upon connection closure.
	c1.SetDeadline(neverTimeout)
	go func() ***REMOVED***
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		c1.Close()
	***REMOVED***()
	go func() ***REMOVED***
		defer wg.Done()
		var err error
		buf := make([]byte, 1024)
		for err == nil ***REMOVED***
			_, err = c1.Read(buf)
		***REMOVED***
	***REMOVED***()
	go func() ***REMOVED***
		defer wg.Done()
		var err error
		buf := make([]byte, 1024)
		for err == nil ***REMOVED***
			_, err = c1.Write(buf)
		***REMOVED***
	***REMOVED***()
***REMOVED***

// testConcurrentMethods tests that the methods of net.Conn can safely
// be called concurrently.
func testConcurrentMethods(t *testing.T, c1, c2 net.Conn) ***REMOVED***
	if runtime.GOOS == "plan9" ***REMOVED***
		t.Skip("skipping on plan9; see https://golang.org/issue/20489")
	***REMOVED***
	go chunkedCopy(c2, c2)

	// The results of the calls may be nonsensical, but this should
	// not trigger a race detector warning.
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ ***REMOVED***
		wg.Add(7)
		go func() ***REMOVED***
			defer wg.Done()
			c1.Read(make([]byte, 1024))
		***REMOVED***()
		go func() ***REMOVED***
			defer wg.Done()
			c1.Write(make([]byte, 1024))
		***REMOVED***()
		go func() ***REMOVED***
			defer wg.Done()
			c1.SetDeadline(time.Now().Add(10 * time.Millisecond))
		***REMOVED***()
		go func() ***REMOVED***
			defer wg.Done()
			c1.SetReadDeadline(aLongTimeAgo)
		***REMOVED***()
		go func() ***REMOVED***
			defer wg.Done()
			c1.SetWriteDeadline(aLongTimeAgo)
		***REMOVED***()
		go func() ***REMOVED***
			defer wg.Done()
			c1.LocalAddr()
		***REMOVED***()
		go func() ***REMOVED***
			defer wg.Done()
			c1.RemoteAddr()
		***REMOVED***()
	***REMOVED***
	wg.Wait() // At worst, the deadline is set 10ms into the future

	resyncConn(t, c1)
	testRoundtrip(t, c1)
***REMOVED***

// checkForTimeoutError checks that the error satisfies the Error interface
// and that Timeout returns true.
func checkForTimeoutError(t *testing.T, err error) ***REMOVED***
	if nerr, ok := err.(net.Error); ok ***REMOVED***
		if !nerr.Timeout() ***REMOVED***
			t.Errorf("err.Timeout() = false, want true")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		t.Errorf("got %T, want net.Error", err)
	***REMOVED***
***REMOVED***

// testRoundtrip writes something into c and reads it back.
// It assumes that everything written into c is echoed back to itself.
func testRoundtrip(t *testing.T, c net.Conn) ***REMOVED***
	if err := c.SetDeadline(neverTimeout); err != nil ***REMOVED***
		t.Errorf("roundtrip SetDeadline error: %v", err)
	***REMOVED***

	const s = "Hello, world!"
	buf := []byte(s)
	if _, err := c.Write(buf); err != nil ***REMOVED***
		t.Errorf("roundtrip Write error: %v", err)
	***REMOVED***
	if _, err := io.ReadFull(c, buf); err != nil ***REMOVED***
		t.Errorf("roundtrip Read error: %v", err)
	***REMOVED***
	if string(buf) != s ***REMOVED***
		t.Errorf("roundtrip data mismatch: got %q, want %q", buf, s)
	***REMOVED***
***REMOVED***

// resyncConn resynchronizes the connection into a sane state.
// It assumes that everything written into c is echoed back to itself.
// It assumes that 0xff is not currently on the wire or in the read buffer.
func resyncConn(t *testing.T, c net.Conn) ***REMOVED***
	c.SetDeadline(neverTimeout)
	errCh := make(chan error)
	go func() ***REMOVED***
		_, err := c.Write([]byte***REMOVED***0xff***REMOVED***)
		errCh <- err
	***REMOVED***()
	buf := make([]byte, 1024)
	for ***REMOVED***
		n, err := c.Read(buf)
		if n > 0 && bytes.IndexByte(buf[:n], 0xff) == n-1 ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			t.Errorf("unexpected Read error: %v", err)
			break
		***REMOVED***
	***REMOVED***
	if err := <-errCh; err != nil ***REMOVED***
		t.Errorf("unexpected Write error: %v", err)
	***REMOVED***
***REMOVED***

// chunkedCopy copies from r to w in fixed-width chunks to avoid
// causing a Write that exceeds the maximum packet size for packet-based
// connections like "unixpacket".
// We assume that the maximum packet size is at least 1024.
func chunkedCopy(w io.Writer, r io.Reader) error ***REMOVED***
	b := make([]byte, 1024)
	_, err := io.CopyBuffer(struct***REMOVED*** io.Writer ***REMOVED******REMOVED***w***REMOVED***, struct***REMOVED*** io.Reader ***REMOVED******REMOVED***r***REMOVED***, b)
	return err
***REMOVED***
