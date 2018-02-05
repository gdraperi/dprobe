// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"reflect"
	"testing"
	"testing/iotest"
	"time"
)

var _ net.Error = errWriteTimeout

type fakeNetConn struct ***REMOVED***
	io.Reader
	io.Writer
***REMOVED***

func (c fakeNetConn) Close() error                       ***REMOVED*** return nil ***REMOVED***
func (c fakeNetConn) LocalAddr() net.Addr                ***REMOVED*** return localAddr ***REMOVED***
func (c fakeNetConn) RemoteAddr() net.Addr               ***REMOVED*** return remoteAddr ***REMOVED***
func (c fakeNetConn) SetDeadline(t time.Time) error      ***REMOVED*** return nil ***REMOVED***
func (c fakeNetConn) SetReadDeadline(t time.Time) error  ***REMOVED*** return nil ***REMOVED***
func (c fakeNetConn) SetWriteDeadline(t time.Time) error ***REMOVED*** return nil ***REMOVED***

type fakeAddr int

var (
	localAddr  = fakeAddr(1)
	remoteAddr = fakeAddr(2)
)

func (a fakeAddr) Network() string ***REMOVED***
	return "net"
***REMOVED***

func (a fakeAddr) String() string ***REMOVED***
	return "str"
***REMOVED***

func TestFraming(t *testing.T) ***REMOVED***
	frameSizes := []int***REMOVED***0, 1, 2, 124, 125, 126, 127, 128, 129, 65534, 65535, 65536, 65537***REMOVED***
	var readChunkers = []struct ***REMOVED***
		name string
		f    func(io.Reader) io.Reader
	***REMOVED******REMOVED***
		***REMOVED***"half", iotest.HalfReader***REMOVED***,
		***REMOVED***"one", iotest.OneByteReader***REMOVED***,
		***REMOVED***"asis", func(r io.Reader) io.Reader ***REMOVED*** return r ***REMOVED******REMOVED***,
	***REMOVED***
	writeBuf := make([]byte, 65537)
	for i := range writeBuf ***REMOVED***
		writeBuf[i] = byte(i)
	***REMOVED***
	var writers = []struct ***REMOVED***
		name string
		f    func(w io.Writer, n int) (int, error)
	***REMOVED******REMOVED***
		***REMOVED***"iocopy", func(w io.Writer, n int) (int, error) ***REMOVED***
			nn, err := io.Copy(w, bytes.NewReader(writeBuf[:n]))
			return int(nn), err
		***REMOVED******REMOVED***,
		***REMOVED***"write", func(w io.Writer, n int) (int, error) ***REMOVED***
			return w.Write(writeBuf[:n])
		***REMOVED******REMOVED***,
		***REMOVED***"string", func(w io.Writer, n int) (int, error) ***REMOVED***
			return io.WriteString(w, string(writeBuf[:n]))
		***REMOVED******REMOVED***,
	***REMOVED***

	for _, compress := range []bool***REMOVED***false, true***REMOVED*** ***REMOVED***
		for _, isServer := range []bool***REMOVED***true, false***REMOVED*** ***REMOVED***
			for _, chunker := range readChunkers ***REMOVED***

				var connBuf bytes.Buffer
				wc := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: &connBuf***REMOVED***, isServer, 1024, 1024)
				rc := newConn(fakeNetConn***REMOVED***Reader: chunker.f(&connBuf), Writer: nil***REMOVED***, !isServer, 1024, 1024)
				if compress ***REMOVED***
					wc.newCompressionWriter = compressNoContextTakeover
					rc.newDecompressionReader = decompressNoContextTakeover
				***REMOVED***
				for _, n := range frameSizes ***REMOVED***
					for _, writer := range writers ***REMOVED***
						name := fmt.Sprintf("z:%v, s:%v, r:%s, n:%d w:%s", compress, isServer, chunker.name, n, writer.name)

						w, err := wc.NextWriter(TextMessage)
						if err != nil ***REMOVED***
							t.Errorf("%s: wc.NextWriter() returned %v", name, err)
							continue
						***REMOVED***
						nn, err := writer.f(w, n)
						if err != nil || nn != n ***REMOVED***
							t.Errorf("%s: w.Write(writeBuf[:n]) returned %d, %v", name, nn, err)
							continue
						***REMOVED***
						err = w.Close()
						if err != nil ***REMOVED***
							t.Errorf("%s: w.Close() returned %v", name, err)
							continue
						***REMOVED***

						opCode, r, err := rc.NextReader()
						if err != nil || opCode != TextMessage ***REMOVED***
							t.Errorf("%s: NextReader() returned %d, r, %v", name, opCode, err)
							continue
						***REMOVED***
						rbuf, err := ioutil.ReadAll(r)
						if err != nil ***REMOVED***
							t.Errorf("%s: ReadFull() returned rbuf, %v", name, err)
							continue
						***REMOVED***

						if len(rbuf) != n ***REMOVED***
							t.Errorf("%s: len(rbuf) is %d, want %d", name, len(rbuf), n)
							continue
						***REMOVED***

						for i, b := range rbuf ***REMOVED***
							if byte(i) != b ***REMOVED***
								t.Errorf("%s: bad byte at offset %d", name, i)
								break
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestControl(t *testing.T) ***REMOVED***
	const message = "this is a ping/pong messsage"
	for _, isServer := range []bool***REMOVED***true, false***REMOVED*** ***REMOVED***
		for _, isWriteControl := range []bool***REMOVED***true, false***REMOVED*** ***REMOVED***
			name := fmt.Sprintf("s:%v, wc:%v", isServer, isWriteControl)
			var connBuf bytes.Buffer
			wc := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: &connBuf***REMOVED***, isServer, 1024, 1024)
			rc := newConn(fakeNetConn***REMOVED***Reader: &connBuf, Writer: nil***REMOVED***, !isServer, 1024, 1024)
			if isWriteControl ***REMOVED***
				wc.WriteControl(PongMessage, []byte(message), time.Now().Add(time.Second))
			***REMOVED*** else ***REMOVED***
				w, err := wc.NextWriter(PongMessage)
				if err != nil ***REMOVED***
					t.Errorf("%s: wc.NextWriter() returned %v", name, err)
					continue
				***REMOVED***
				if _, err := w.Write([]byte(message)); err != nil ***REMOVED***
					t.Errorf("%s: w.Write() returned %v", name, err)
					continue
				***REMOVED***
				if err := w.Close(); err != nil ***REMOVED***
					t.Errorf("%s: w.Close() returned %v", name, err)
					continue
				***REMOVED***
				var actualMessage string
				rc.SetPongHandler(func(s string) error ***REMOVED*** actualMessage = s; return nil ***REMOVED***)
				rc.NextReader()
				if actualMessage != message ***REMOVED***
					t.Errorf("%s: pong=%q, want %q", name, actualMessage, message)
					continue
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCloseFrameBeforeFinalMessageFrame(t *testing.T) ***REMOVED***
	const bufSize = 512

	expectedErr := &CloseError***REMOVED***Code: CloseNormalClosure, Text: "hello"***REMOVED***

	var b1, b2 bytes.Buffer
	wc := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: &b1***REMOVED***, false, 1024, bufSize)
	rc := newConn(fakeNetConn***REMOVED***Reader: &b1, Writer: &b2***REMOVED***, true, 1024, 1024)

	w, _ := wc.NextWriter(BinaryMessage)
	w.Write(make([]byte, bufSize+bufSize/2))
	wc.WriteControl(CloseMessage, FormatCloseMessage(expectedErr.Code, expectedErr.Text), time.Now().Add(10*time.Second))
	w.Close()

	op, r, err := rc.NextReader()
	if op != BinaryMessage || err != nil ***REMOVED***
		t.Fatalf("NextReader() returned %d, %v", op, err)
	***REMOVED***
	_, err = io.Copy(ioutil.Discard, r)
	if !reflect.DeepEqual(err, expectedErr) ***REMOVED***
		t.Fatalf("io.Copy() returned %v, want %v", err, expectedErr)
	***REMOVED***
	_, _, err = rc.NextReader()
	if !reflect.DeepEqual(err, expectedErr) ***REMOVED***
		t.Fatalf("NextReader() returned %v, want %v", err, expectedErr)
	***REMOVED***
***REMOVED***

func TestEOFWithinFrame(t *testing.T) ***REMOVED***
	const bufSize = 64

	for n := 0; ; n++ ***REMOVED***
		var b bytes.Buffer
		wc := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: &b***REMOVED***, false, 1024, 1024)
		rc := newConn(fakeNetConn***REMOVED***Reader: &b, Writer: nil***REMOVED***, true, 1024, 1024)

		w, _ := wc.NextWriter(BinaryMessage)
		w.Write(make([]byte, bufSize))
		w.Close()

		if n >= b.Len() ***REMOVED***
			break
		***REMOVED***
		b.Truncate(n)

		op, r, err := rc.NextReader()
		if err == errUnexpectedEOF ***REMOVED***
			continue
		***REMOVED***
		if op != BinaryMessage || err != nil ***REMOVED***
			t.Fatalf("%d: NextReader() returned %d, %v", n, op, err)
		***REMOVED***
		_, err = io.Copy(ioutil.Discard, r)
		if err != errUnexpectedEOF ***REMOVED***
			t.Fatalf("%d: io.Copy() returned %v, want %v", n, err, errUnexpectedEOF)
		***REMOVED***
		_, _, err = rc.NextReader()
		if err != errUnexpectedEOF ***REMOVED***
			t.Fatalf("%d: NextReader() returned %v, want %v", n, err, errUnexpectedEOF)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEOFBeforeFinalFrame(t *testing.T) ***REMOVED***
	const bufSize = 512

	var b1, b2 bytes.Buffer
	wc := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: &b1***REMOVED***, false, 1024, bufSize)
	rc := newConn(fakeNetConn***REMOVED***Reader: &b1, Writer: &b2***REMOVED***, true, 1024, 1024)

	w, _ := wc.NextWriter(BinaryMessage)
	w.Write(make([]byte, bufSize+bufSize/2))

	op, r, err := rc.NextReader()
	if op != BinaryMessage || err != nil ***REMOVED***
		t.Fatalf("NextReader() returned %d, %v", op, err)
	***REMOVED***
	_, err = io.Copy(ioutil.Discard, r)
	if err != errUnexpectedEOF ***REMOVED***
		t.Fatalf("io.Copy() returned %v, want %v", err, errUnexpectedEOF)
	***REMOVED***
	_, _, err = rc.NextReader()
	if err != errUnexpectedEOF ***REMOVED***
		t.Fatalf("NextReader() returned %v, want %v", err, errUnexpectedEOF)
	***REMOVED***
***REMOVED***

func TestWriteAfterMessageWriterClose(t *testing.T) ***REMOVED***
	wc := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: &bytes.Buffer***REMOVED******REMOVED******REMOVED***, false, 1024, 1024)
	w, _ := wc.NextWriter(BinaryMessage)
	io.WriteString(w, "hello")
	if err := w.Close(); err != nil ***REMOVED***
		t.Fatalf("unxpected error closing message writer, %v", err)
	***REMOVED***

	if _, err := io.WriteString(w, "world"); err == nil ***REMOVED***
		t.Fatalf("no error writing after close")
	***REMOVED***

	w, _ = wc.NextWriter(BinaryMessage)
	io.WriteString(w, "hello")

	// close w by getting next writer
	_, err := wc.NextWriter(BinaryMessage)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error getting next writer, %v", err)
	***REMOVED***

	if _, err := io.WriteString(w, "world"); err == nil ***REMOVED***
		t.Fatalf("no error writing after close")
	***REMOVED***
***REMOVED***

func TestReadLimit(t *testing.T) ***REMOVED***

	const readLimit = 512
	message := make([]byte, readLimit+1)

	var b1, b2 bytes.Buffer
	wc := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: &b1***REMOVED***, false, 1024, readLimit-2)
	rc := newConn(fakeNetConn***REMOVED***Reader: &b1, Writer: &b2***REMOVED***, true, 1024, 1024)
	rc.SetReadLimit(readLimit)

	// Send message at the limit with interleaved pong.
	w, _ := wc.NextWriter(BinaryMessage)
	w.Write(message[:readLimit-1])
	wc.WriteControl(PongMessage, []byte("this is a pong"), time.Now().Add(10*time.Second))
	w.Write(message[:1])
	w.Close()

	// Send message larger than the limit.
	wc.WriteMessage(BinaryMessage, message[:readLimit+1])

	op, _, err := rc.NextReader()
	if op != BinaryMessage || err != nil ***REMOVED***
		t.Fatalf("1: NextReader() returned %d, %v", op, err)
	***REMOVED***
	op, r, err := rc.NextReader()
	if op != BinaryMessage || err != nil ***REMOVED***
		t.Fatalf("2: NextReader() returned %d, %v", op, err)
	***REMOVED***
	_, err = io.Copy(ioutil.Discard, r)
	if err != ErrReadLimit ***REMOVED***
		t.Fatalf("io.Copy() returned %v", err)
	***REMOVED***
***REMOVED***

func TestAddrs(t *testing.T) ***REMOVED***
	c := newConn(&fakeNetConn***REMOVED******REMOVED***, true, 1024, 1024)
	if c.LocalAddr() != localAddr ***REMOVED***
		t.Errorf("LocalAddr = %v, want %v", c.LocalAddr(), localAddr)
	***REMOVED***
	if c.RemoteAddr() != remoteAddr ***REMOVED***
		t.Errorf("RemoteAddr = %v, want %v", c.RemoteAddr(), remoteAddr)
	***REMOVED***
***REMOVED***

func TestUnderlyingConn(t *testing.T) ***REMOVED***
	var b1, b2 bytes.Buffer
	fc := fakeNetConn***REMOVED***Reader: &b1, Writer: &b2***REMOVED***
	c := newConn(fc, true, 1024, 1024)
	ul := c.UnderlyingConn()
	if ul != fc ***REMOVED***
		t.Fatalf("Underlying conn is not what it should be.")
	***REMOVED***
***REMOVED***

func TestBufioReadBytes(t *testing.T) ***REMOVED***
	// Test calling bufio.ReadBytes for value longer than read buffer size.

	m := make([]byte, 512)
	m[len(m)-1] = '\n'

	var b1, b2 bytes.Buffer
	wc := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: &b1***REMOVED***, false, len(m)+64, len(m)+64)
	rc := newConn(fakeNetConn***REMOVED***Reader: &b1, Writer: &b2***REMOVED***, true, len(m)-64, len(m)-64)

	w, _ := wc.NextWriter(BinaryMessage)
	w.Write(m)
	w.Close()

	op, r, err := rc.NextReader()
	if op != BinaryMessage || err != nil ***REMOVED***
		t.Fatalf("NextReader() returned %d, %v", op, err)
	***REMOVED***

	br := bufio.NewReader(r)
	p, err := br.ReadBytes('\n')
	if err != nil ***REMOVED***
		t.Fatalf("ReadBytes() returned %v", err)
	***REMOVED***
	if len(p) != len(m) ***REMOVED***
		t.Fatalf("read returned %d bytes, want %d bytes", len(p), len(m))
	***REMOVED***
***REMOVED***

var closeErrorTests = []struct ***REMOVED***
	err   error
	codes []int
	ok    bool
***REMOVED******REMOVED***
	***REMOVED***&CloseError***REMOVED***Code: CloseNormalClosure***REMOVED***, []int***REMOVED***CloseNormalClosure***REMOVED***, true***REMOVED***,
	***REMOVED***&CloseError***REMOVED***Code: CloseNormalClosure***REMOVED***, []int***REMOVED***CloseNoStatusReceived***REMOVED***, false***REMOVED***,
	***REMOVED***&CloseError***REMOVED***Code: CloseNormalClosure***REMOVED***, []int***REMOVED***CloseNoStatusReceived, CloseNormalClosure***REMOVED***, true***REMOVED***,
	***REMOVED***errors.New("hello"), []int***REMOVED***CloseNormalClosure***REMOVED***, false***REMOVED***,
***REMOVED***

func TestCloseError(t *testing.T) ***REMOVED***
	for _, tt := range closeErrorTests ***REMOVED***
		ok := IsCloseError(tt.err, tt.codes...)
		if ok != tt.ok ***REMOVED***
			t.Errorf("IsCloseError(%#v, %#v) returned %v, want %v", tt.err, tt.codes, ok, tt.ok)
		***REMOVED***
	***REMOVED***
***REMOVED***

var unexpectedCloseErrorTests = []struct ***REMOVED***
	err   error
	codes []int
	ok    bool
***REMOVED******REMOVED***
	***REMOVED***&CloseError***REMOVED***Code: CloseNormalClosure***REMOVED***, []int***REMOVED***CloseNormalClosure***REMOVED***, false***REMOVED***,
	***REMOVED***&CloseError***REMOVED***Code: CloseNormalClosure***REMOVED***, []int***REMOVED***CloseNoStatusReceived***REMOVED***, true***REMOVED***,
	***REMOVED***&CloseError***REMOVED***Code: CloseNormalClosure***REMOVED***, []int***REMOVED***CloseNoStatusReceived, CloseNormalClosure***REMOVED***, false***REMOVED***,
	***REMOVED***errors.New("hello"), []int***REMOVED***CloseNormalClosure***REMOVED***, false***REMOVED***,
***REMOVED***

func TestUnexpectedCloseErrors(t *testing.T) ***REMOVED***
	for _, tt := range unexpectedCloseErrorTests ***REMOVED***
		ok := IsUnexpectedCloseError(tt.err, tt.codes...)
		if ok != tt.ok ***REMOVED***
			t.Errorf("IsUnexpectedCloseError(%#v, %#v) returned %v, want %v", tt.err, tt.codes, ok, tt.ok)
		***REMOVED***
	***REMOVED***
***REMOVED***

type blockingWriter struct ***REMOVED***
	c1, c2 chan struct***REMOVED******REMOVED***
***REMOVED***

func (w blockingWriter) Write(p []byte) (int, error) ***REMOVED***
	// Allow main to continue
	close(w.c1)
	// Wait for panic in main
	<-w.c2
	return len(p), nil
***REMOVED***

func TestConcurrentWritePanic(t *testing.T) ***REMOVED***
	w := blockingWriter***REMOVED***make(chan struct***REMOVED******REMOVED***), make(chan struct***REMOVED******REMOVED***)***REMOVED***
	c := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: w***REMOVED***, false, 1024, 1024)
	go func() ***REMOVED***
		c.WriteMessage(TextMessage, []byte***REMOVED******REMOVED***)
	***REMOVED***()

	// wait for goroutine to block in write.
	<-w.c1

	defer func() ***REMOVED***
		close(w.c2)
		if v := recover(); v != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***()

	c.WriteMessage(TextMessage, []byte***REMOVED******REMOVED***)
	t.Fatal("should not get here")
***REMOVED***

type failingReader struct***REMOVED******REMOVED***

func (r failingReader) Read(p []byte) (int, error) ***REMOVED***
	return 0, io.EOF
***REMOVED***

func TestFailedConnectionReadPanic(t *testing.T) ***REMOVED***
	c := newConn(fakeNetConn***REMOVED***Reader: failingReader***REMOVED******REMOVED***, Writer: nil***REMOVED***, false, 1024, 1024)

	defer func() ***REMOVED***
		if v := recover(); v != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***()

	for i := 0; i < 20000; i++ ***REMOVED***
		c.ReadMessage()
	***REMOVED***
	t.Fatal("should not get here")
***REMOVED***

func TestBufioReuse(t *testing.T) ***REMOVED***
	brw := bufio.NewReadWriter(bufio.NewReader(nil), bufio.NewWriter(nil))
	c := newConnBRW(nil, false, 0, 0, brw)

	if c.br != brw.Reader ***REMOVED***
		t.Error("connection did not reuse bufio.Reader")
	***REMOVED***

	var wh writeHook
	brw.Writer.Reset(&wh)
	brw.WriteByte(0)
	brw.Flush()
	if &c.writeBuf[0] != &wh.p[0] ***REMOVED***
		t.Error("connection did not reuse bufio.Writer")
	***REMOVED***

	brw = bufio.NewReadWriter(bufio.NewReaderSize(nil, 0), bufio.NewWriterSize(nil, 0))
	c = newConnBRW(nil, false, 0, 0, brw)

	if c.br == brw.Reader ***REMOVED***
		t.Error("connection used bufio.Reader with small size")
	***REMOVED***

	brw.Writer.Reset(&wh)
	brw.WriteByte(0)
	brw.Flush()
	if &c.writeBuf[0] != &wh.p[0] ***REMOVED***
		t.Error("connection used bufio.Writer with small size")
	***REMOVED***

***REMOVED***
