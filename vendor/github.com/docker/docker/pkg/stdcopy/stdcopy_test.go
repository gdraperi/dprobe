package stdcopy

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNewStdWriter(t *testing.T) ***REMOVED***
	writer := NewStdWriter(ioutil.Discard, Stdout)
	if writer == nil ***REMOVED***
		t.Fatalf("NewStdWriter with an invalid StdType should not return nil.")
	***REMOVED***
***REMOVED***

func TestWriteWithUninitializedStdWriter(t *testing.T) ***REMOVED***
	writer := stdWriter***REMOVED***
		Writer: nil,
		prefix: byte(Stdout),
	***REMOVED***
	n, err := writer.Write([]byte("Something here"))
	if n != 0 || err == nil ***REMOVED***
		t.Fatalf("Should fail when given an uncomplete or uninitialized StdWriter")
	***REMOVED***
***REMOVED***

func TestWriteWithNilBytes(t *testing.T) ***REMOVED***
	writer := NewStdWriter(ioutil.Discard, Stdout)
	n, err := writer.Write(nil)
	if err != nil ***REMOVED***
		t.Fatalf("Shouldn't have fail when given no data")
	***REMOVED***
	if n > 0 ***REMOVED***
		t.Fatalf("Write should have written 0 byte, but has written %d", n)
	***REMOVED***
***REMOVED***

func TestWrite(t *testing.T) ***REMOVED***
	writer := NewStdWriter(ioutil.Discard, Stdout)
	data := []byte("Test StdWrite.Write")
	n, err := writer.Write(data)
	if err != nil ***REMOVED***
		t.Fatalf("Error while writing with StdWrite")
	***REMOVED***
	if n != len(data) ***REMOVED***
		t.Fatalf("Write should have written %d byte but wrote %d.", len(data), n)
	***REMOVED***
***REMOVED***

type errWriter struct ***REMOVED***
	n   int
	err error
***REMOVED***

func (f *errWriter) Write(buf []byte) (int, error) ***REMOVED***
	return f.n, f.err
***REMOVED***

func TestWriteWithWriterError(t *testing.T) ***REMOVED***
	expectedError := errors.New("expected")
	expectedReturnedBytes := 10
	writer := NewStdWriter(&errWriter***REMOVED***
		n:   stdWriterPrefixLen + expectedReturnedBytes,
		err: expectedError***REMOVED***, Stdout)
	data := []byte("This won't get written, sigh")
	n, err := writer.Write(data)
	if err != expectedError ***REMOVED***
		t.Fatalf("Didn't get expected error.")
	***REMOVED***
	if n != expectedReturnedBytes ***REMOVED***
		t.Fatalf("Didn't get expected written bytes %d, got %d.",
			expectedReturnedBytes, n)
	***REMOVED***
***REMOVED***

func TestWriteDoesNotReturnNegativeWrittenBytes(t *testing.T) ***REMOVED***
	writer := NewStdWriter(&errWriter***REMOVED***n: -1***REMOVED***, Stdout)
	data := []byte("This won't get written, sigh")
	actual, _ := writer.Write(data)
	if actual != 0 ***REMOVED***
		t.Fatalf("Expected returned written bytes equal to 0, got %d", actual)
	***REMOVED***
***REMOVED***

func getSrcBuffer(stdOutBytes, stdErrBytes []byte) (buffer *bytes.Buffer, err error) ***REMOVED***
	buffer = new(bytes.Buffer)
	dstOut := NewStdWriter(buffer, Stdout)
	_, err = dstOut.Write(stdOutBytes)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	dstErr := NewStdWriter(buffer, Stderr)
	_, err = dstErr.Write(stdErrBytes)
	return
***REMOVED***

func TestStdCopyWriteAndRead(t *testing.T) ***REMOVED***
	stdOutBytes := []byte(strings.Repeat("o", startingBufLen))
	stdErrBytes := []byte(strings.Repeat("e", startingBufLen))
	buffer, err := getSrcBuffer(stdOutBytes, stdErrBytes)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	written, err := StdCopy(ioutil.Discard, ioutil.Discard, buffer)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expectedTotalWritten := len(stdOutBytes) + len(stdErrBytes)
	if written != int64(expectedTotalWritten) ***REMOVED***
		t.Fatalf("Expected to have total of %d bytes written, got %d", expectedTotalWritten, written)
	***REMOVED***
***REMOVED***

type customReader struct ***REMOVED***
	n            int
	err          error
	totalCalls   int
	correctCalls int
	src          *bytes.Buffer
***REMOVED***

func (f *customReader) Read(buf []byte) (int, error) ***REMOVED***
	f.totalCalls++
	if f.totalCalls <= f.correctCalls ***REMOVED***
		return f.src.Read(buf)
	***REMOVED***
	return f.n, f.err
***REMOVED***

func TestStdCopyReturnsErrorReadingHeader(t *testing.T) ***REMOVED***
	expectedError := errors.New("error")
	reader := &customReader***REMOVED***
		err: expectedError***REMOVED***
	written, err := StdCopy(ioutil.Discard, ioutil.Discard, reader)
	if written != 0 ***REMOVED***
		t.Fatalf("Expected 0 bytes read, got %d", written)
	***REMOVED***
	if err != expectedError ***REMOVED***
		t.Fatalf("Didn't get expected error")
	***REMOVED***
***REMOVED***

func TestStdCopyReturnsErrorReadingFrame(t *testing.T) ***REMOVED***
	expectedError := errors.New("error")
	stdOutBytes := []byte(strings.Repeat("o", startingBufLen))
	stdErrBytes := []byte(strings.Repeat("e", startingBufLen))
	buffer, err := getSrcBuffer(stdOutBytes, stdErrBytes)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	reader := &customReader***REMOVED***
		correctCalls: 1,
		n:            stdWriterPrefixLen + 1,
		err:          expectedError,
		src:          buffer***REMOVED***
	written, err := StdCopy(ioutil.Discard, ioutil.Discard, reader)
	if written != 0 ***REMOVED***
		t.Fatalf("Expected 0 bytes read, got %d", written)
	***REMOVED***
	if err != expectedError ***REMOVED***
		t.Fatalf("Didn't get expected error")
	***REMOVED***
***REMOVED***

func TestStdCopyDetectsCorruptedFrame(t *testing.T) ***REMOVED***
	stdOutBytes := []byte(strings.Repeat("o", startingBufLen))
	stdErrBytes := []byte(strings.Repeat("e", startingBufLen))
	buffer, err := getSrcBuffer(stdOutBytes, stdErrBytes)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	reader := &customReader***REMOVED***
		correctCalls: 1,
		n:            stdWriterPrefixLen + 1,
		err:          io.EOF,
		src:          buffer***REMOVED***
	written, err := StdCopy(ioutil.Discard, ioutil.Discard, reader)
	if written != startingBufLen ***REMOVED***
		t.Fatalf("Expected %d bytes read, got %d", startingBufLen, written)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Fatal("Didn't get nil error")
	***REMOVED***
***REMOVED***

func TestStdCopyWithInvalidInputHeader(t *testing.T) ***REMOVED***
	dstOut := NewStdWriter(ioutil.Discard, Stdout)
	dstErr := NewStdWriter(ioutil.Discard, Stderr)
	src := strings.NewReader("Invalid input")
	_, err := StdCopy(dstOut, dstErr, src)
	if err == nil ***REMOVED***
		t.Fatal("StdCopy with invalid input header should fail.")
	***REMOVED***
***REMOVED***

func TestStdCopyWithCorruptedPrefix(t *testing.T) ***REMOVED***
	data := []byte***REMOVED***0x01, 0x02, 0x03***REMOVED***
	src := bytes.NewReader(data)
	written, err := StdCopy(nil, nil, src)
	if err != nil ***REMOVED***
		t.Fatalf("StdCopy should not return an error with corrupted prefix.")
	***REMOVED***
	if written != 0 ***REMOVED***
		t.Fatalf("StdCopy should have written 0, but has written %d", written)
	***REMOVED***
***REMOVED***

func TestStdCopyReturnsWriteErrors(t *testing.T) ***REMOVED***
	stdOutBytes := []byte(strings.Repeat("o", startingBufLen))
	stdErrBytes := []byte(strings.Repeat("e", startingBufLen))
	buffer, err := getSrcBuffer(stdOutBytes, stdErrBytes)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expectedError := errors.New("expected")

	dstOut := &errWriter***REMOVED***err: expectedError***REMOVED***

	written, err := StdCopy(dstOut, ioutil.Discard, buffer)
	if written != 0 ***REMOVED***
		t.Fatalf("StdCopy should have written 0, but has written %d", written)
	***REMOVED***
	if err != expectedError ***REMOVED***
		t.Fatalf("Didn't get expected error, got %v", err)
	***REMOVED***
***REMOVED***

func TestStdCopyDetectsNotFullyWrittenFrames(t *testing.T) ***REMOVED***
	stdOutBytes := []byte(strings.Repeat("o", startingBufLen))
	stdErrBytes := []byte(strings.Repeat("e", startingBufLen))
	buffer, err := getSrcBuffer(stdOutBytes, stdErrBytes)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	dstOut := &errWriter***REMOVED***n: startingBufLen - 10***REMOVED***

	written, err := StdCopy(dstOut, ioutil.Discard, buffer)
	if written != 0 ***REMOVED***
		t.Fatalf("StdCopy should have return 0 written bytes, but returned %d", written)
	***REMOVED***
	if err != io.ErrShortWrite ***REMOVED***
		t.Fatalf("Didn't get expected io.ErrShortWrite error")
	***REMOVED***
***REMOVED***

// TestStdCopyReturnsErrorFromSystem tests that StdCopy correctly returns an
// error, when that error is muxed into the Systemerr stream.
func TestStdCopyReturnsErrorFromSystem(t *testing.T) ***REMOVED***
	// write in the basic messages, just so there's some fluff in there
	stdOutBytes := []byte(strings.Repeat("o", startingBufLen))
	stdErrBytes := []byte(strings.Repeat("e", startingBufLen))
	buffer, err := getSrcBuffer(stdOutBytes, stdErrBytes)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// add in an error message on the Systemerr stream
	systemErrBytes := []byte(strings.Repeat("S", startingBufLen))
	systemWriter := NewStdWriter(buffer, Systemerr)
	_, err = systemWriter.Write(systemErrBytes)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// now copy and demux. we should expect an error containing the string we
	// wrote out
	_, err = StdCopy(ioutil.Discard, ioutil.Discard, buffer)
	if err == nil ***REMOVED***
		t.Fatal("expected error, got none")
	***REMOVED***
	if !strings.Contains(err.Error(), string(systemErrBytes)) ***REMOVED***
		t.Fatal("expected error to contain message")
	***REMOVED***
***REMOVED***

func BenchmarkWrite(b *testing.B) ***REMOVED***
	w := NewStdWriter(ioutil.Discard, Stdout)
	data := []byte("Test line for testing stdwriter performance\n")
	data = bytes.Repeat(data, 100)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if _, err := w.Write(data); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
