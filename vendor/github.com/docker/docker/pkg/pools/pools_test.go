package pools

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBufioReaderPoolGetWithNoReaderShouldCreateOne(t *testing.T) ***REMOVED***
	reader := BufioReader32KPool.Get(nil)
	if reader == nil ***REMOVED***
		t.Fatalf("BufioReaderPool should have create a bufio.Reader but did not.")
	***REMOVED***
***REMOVED***

func TestBufioReaderPoolPutAndGet(t *testing.T) ***REMOVED***
	sr := bufio.NewReader(strings.NewReader("foobar"))
	reader := BufioReader32KPool.Get(sr)
	if reader == nil ***REMOVED***
		t.Fatalf("BufioReaderPool should not return a nil reader.")
	***REMOVED***
	// verify the first 3 byte
	buf1 := make([]byte, 3)
	_, err := reader.Read(buf1)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if actual := string(buf1); actual != "foo" ***REMOVED***
		t.Fatalf("The first letter should have been 'foo' but was %v", actual)
	***REMOVED***
	BufioReader32KPool.Put(reader)
	// Try to read the next 3 bytes
	_, err = sr.Read(make([]byte, 3))
	if err == nil || err != io.EOF ***REMOVED***
		t.Fatalf("The buffer should have been empty, issue an EOF error.")
	***REMOVED***
***REMOVED***

type simpleReaderCloser struct ***REMOVED***
	io.Reader
	closed bool
***REMOVED***

func (r *simpleReaderCloser) Close() error ***REMOVED***
	r.closed = true
	return nil
***REMOVED***

func TestNewReadCloserWrapperWithAReadCloser(t *testing.T) ***REMOVED***
	br := bufio.NewReader(strings.NewReader(""))
	sr := &simpleReaderCloser***REMOVED***
		Reader: strings.NewReader("foobar"),
		closed: false,
	***REMOVED***
	reader := BufioReader32KPool.NewReadCloserWrapper(br, sr)
	if reader == nil ***REMOVED***
		t.Fatalf("NewReadCloserWrapper should not return a nil reader.")
	***REMOVED***
	// Verify the content of reader
	buf := make([]byte, 3)
	_, err := reader.Read(buf)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if actual := string(buf); actual != "foo" ***REMOVED***
		t.Fatalf("The first 3 letter should have been 'foo' but were %v", actual)
	***REMOVED***
	reader.Close()
	// Read 3 more bytes "bar"
	_, err = reader.Read(buf)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if actual := string(buf); actual != "bar" ***REMOVED***
		t.Fatalf("The first 3 letter should have been 'bar' but were %v", actual)
	***REMOVED***
	if !sr.closed ***REMOVED***
		t.Fatalf("The ReaderCloser should have been closed, it is not.")
	***REMOVED***
***REMOVED***

func TestBufioWriterPoolGetWithNoReaderShouldCreateOne(t *testing.T) ***REMOVED***
	writer := BufioWriter32KPool.Get(nil)
	if writer == nil ***REMOVED***
		t.Fatalf("BufioWriterPool should have create a bufio.Writer but did not.")
	***REMOVED***
***REMOVED***

func TestBufioWriterPoolPutAndGet(t *testing.T) ***REMOVED***
	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)
	writer := BufioWriter32KPool.Get(bw)
	require.NotNil(t, writer)

	written, err := writer.Write([]byte("foobar"))
	require.NoError(t, err)
	assert.Equal(t, 6, written)

	// Make sure we Flush all the way ?
	writer.Flush()
	bw.Flush()
	assert.Len(t, buf.Bytes(), 6)
	// Reset the buffer
	buf.Reset()
	BufioWriter32KPool.Put(writer)
	// Try to write something
	if _, err = writer.Write([]byte("barfoo")); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// If we now try to flush it, it should panic (the writer is nil)
	// recover it
	defer func() ***REMOVED***
		if r := recover(); r == nil ***REMOVED***
			t.Fatal("Trying to flush the writter should have 'paniced', did not.")
		***REMOVED***
	***REMOVED***()
	writer.Flush()
***REMOVED***

type simpleWriterCloser struct ***REMOVED***
	io.Writer
	closed bool
***REMOVED***

func (r *simpleWriterCloser) Close() error ***REMOVED***
	r.closed = true
	return nil
***REMOVED***

func TestNewWriteCloserWrapperWithAWriteCloser(t *testing.T) ***REMOVED***
	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)
	sw := &simpleWriterCloser***REMOVED***
		Writer: new(bytes.Buffer),
		closed: false,
	***REMOVED***
	bw.Flush()
	writer := BufioWriter32KPool.NewWriteCloserWrapper(bw, sw)
	if writer == nil ***REMOVED***
		t.Fatalf("BufioReaderPool should not return a nil writer.")
	***REMOVED***
	written, err := writer.Write([]byte("foobar"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if written != 6 ***REMOVED***
		t.Fatalf("Should have written 6 bytes, but wrote %v bytes", written)
	***REMOVED***
	writer.Close()
	if !sw.closed ***REMOVED***
		t.Fatalf("The ReaderCloser should have been closed, it is not.")
	***REMOVED***
***REMOVED***

func TestBufferPoolPutAndGet(t *testing.T) ***REMOVED***
	buf := buffer32KPool.Get()
	buffer32KPool.Put(buf)
***REMOVED***
