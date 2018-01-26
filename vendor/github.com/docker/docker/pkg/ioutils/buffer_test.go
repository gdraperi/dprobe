package ioutils

import (
	"bytes"
	"testing"
)

func TestFixedBufferCap(t *testing.T) ***REMOVED***
	buf := &fixedBuffer***REMOVED***buf: make([]byte, 0, 5)***REMOVED***

	n := buf.Cap()
	if n != 5 ***REMOVED***
		t.Fatalf("expected buffer capacity to be 5 bytes, got %d", n)
	***REMOVED***
***REMOVED***

func TestFixedBufferLen(t *testing.T) ***REMOVED***
	buf := &fixedBuffer***REMOVED***buf: make([]byte, 0, 10)***REMOVED***

	buf.Write([]byte("hello"))
	l := buf.Len()
	if l != 5 ***REMOVED***
		t.Fatalf("expected buffer length to be 5 bytes, got %d", l)
	***REMOVED***

	buf.Write([]byte("world"))
	l = buf.Len()
	if l != 10 ***REMOVED***
		t.Fatalf("expected buffer length to be 10 bytes, got %d", l)
	***REMOVED***

	// read 5 bytes
	b := make([]byte, 5)
	buf.Read(b)

	l = buf.Len()
	if l != 5 ***REMOVED***
		t.Fatalf("expected buffer length to be 5 bytes, got %d", l)
	***REMOVED***

	n, err := buf.Write([]byte("i-wont-fit"))
	if n != 0 ***REMOVED***
		t.Fatalf("expected no bytes to be written to buffer, got %d", n)
	***REMOVED***
	if err != errBufferFull ***REMOVED***
		t.Fatalf("expected errBufferFull, got %v", err)
	***REMOVED***

	l = buf.Len()
	if l != 5 ***REMOVED***
		t.Fatalf("expected buffer length to still be 5 bytes, got %d", l)
	***REMOVED***

	buf.Reset()
	l = buf.Len()
	if l != 0 ***REMOVED***
		t.Fatalf("expected buffer length to still be 0 bytes, got %d", l)
	***REMOVED***
***REMOVED***

func TestFixedBufferString(t *testing.T) ***REMOVED***
	buf := &fixedBuffer***REMOVED***buf: make([]byte, 0, 10)***REMOVED***

	buf.Write([]byte("hello"))
	buf.Write([]byte("world"))

	out := buf.String()
	if out != "helloworld" ***REMOVED***
		t.Fatalf("expected output to be \"helloworld\", got %q", out)
	***REMOVED***

	// read 5 bytes
	b := make([]byte, 5)
	buf.Read(b)

	// test that fixedBuffer.String() only returns the part that hasn't been read
	out = buf.String()
	if out != "world" ***REMOVED***
		t.Fatalf("expected output to be \"world\", got %q", out)
	***REMOVED***
***REMOVED***

func TestFixedBufferWrite(t *testing.T) ***REMOVED***
	buf := &fixedBuffer***REMOVED***buf: make([]byte, 0, 64)***REMOVED***
	n, err := buf.Write([]byte("hello"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if n != 5 ***REMOVED***
		t.Fatalf("expected 5 bytes written, got %d", n)
	***REMOVED***

	if string(buf.buf[:5]) != "hello" ***REMOVED***
		t.Fatalf("expected \"hello\", got %q", string(buf.buf[:5]))
	***REMOVED***

	n, err = buf.Write(bytes.Repeat([]byte***REMOVED***1***REMOVED***, 64))
	if n != 59 ***REMOVED***
		t.Fatalf("expected 59 bytes written before buffer is full, got %d", n)
	***REMOVED***
	if err != errBufferFull ***REMOVED***
		t.Fatalf("expected errBufferFull, got %v - %v", err, buf.buf[:64])
	***REMOVED***
***REMOVED***

func TestFixedBufferRead(t *testing.T) ***REMOVED***
	buf := &fixedBuffer***REMOVED***buf: make([]byte, 0, 64)***REMOVED***
	if _, err := buf.Write([]byte("hello world")); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	b := make([]byte, 5)
	n, err := buf.Read(b)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if n != 5 ***REMOVED***
		t.Fatalf("expected 5 bytes read, got %d - %s", n, buf.String())
	***REMOVED***

	if string(b) != "hello" ***REMOVED***
		t.Fatalf("expected \"hello\", got %q", string(b))
	***REMOVED***

	n, err = buf.Read(b)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if n != 5 ***REMOVED***
		t.Fatalf("expected 5 bytes read, got %d", n)
	***REMOVED***

	if string(b) != " worl" ***REMOVED***
		t.Fatalf("expected \" worl\", got %s", string(b))
	***REMOVED***

	b = b[:1]
	n, err = buf.Read(b)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if n != 1 ***REMOVED***
		t.Fatalf("expected 1 byte read, got %d - %s", n, buf.String())
	***REMOVED***

	if string(b) != "d" ***REMOVED***
		t.Fatalf("expected \"d\", got %s", string(b))
	***REMOVED***
***REMOVED***
