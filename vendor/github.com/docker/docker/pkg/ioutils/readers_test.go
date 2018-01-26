package ioutils

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

// Implement io.Reader
type errorReader struct***REMOVED******REMOVED***

func (r *errorReader) Read(p []byte) (int, error) ***REMOVED***
	return 0, fmt.Errorf("error reader always fail")
***REMOVED***

func TestReadCloserWrapperClose(t *testing.T) ***REMOVED***
	reader := strings.NewReader("A string reader")
	wrapper := NewReadCloserWrapper(reader, func() error ***REMOVED***
		return fmt.Errorf("This will be called when closing")
	***REMOVED***)
	err := wrapper.Close()
	if err == nil || !strings.Contains(err.Error(), "This will be called when closing") ***REMOVED***
		t.Fatalf("readCloserWrapper should have call the anonymous func and thus, fail.")
	***REMOVED***
***REMOVED***

func TestReaderErrWrapperReadOnError(t *testing.T) ***REMOVED***
	called := false
	reader := &errorReader***REMOVED******REMOVED***
	wrapper := NewReaderErrWrapper(reader, func() ***REMOVED***
		called = true
	***REMOVED***)
	_, err := wrapper.Read([]byte***REMOVED******REMOVED***)
	assert.EqualError(t, err, "error reader always fail")
	if !called ***REMOVED***
		t.Fatalf("readErrWrapper should have call the anonymous function on failure")
	***REMOVED***
***REMOVED***

func TestReaderErrWrapperRead(t *testing.T) ***REMOVED***
	reader := strings.NewReader("a string reader.")
	wrapper := NewReaderErrWrapper(reader, func() ***REMOVED***
		t.Fatalf("readErrWrapper should not have called the anonymous function")
	***REMOVED***)
	// Read 20 byte (should be ok with the string above)
	num, err := wrapper.Read(make([]byte, 20))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if num != 16 ***REMOVED***
		t.Fatalf("readerErrWrapper should have read 16 byte, but read %d", num)
	***REMOVED***
***REMOVED***

func TestHashData(t *testing.T) ***REMOVED***
	reader := strings.NewReader("hash-me")
	actual, err := HashData(reader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := "sha256:4d11186aed035cc624d553e10db358492c84a7cd6b9670d92123c144930450aa"
	if actual != expected ***REMOVED***
		t.Fatalf("Expecting %s, got %s", expected, actual)
	***REMOVED***
***REMOVED***

type perpetualReader struct***REMOVED******REMOVED***

func (p *perpetualReader) Read(buf []byte) (n int, err error) ***REMOVED***
	for i := 0; i != len(buf); i++ ***REMOVED***
		buf[i] = 'a'
	***REMOVED***
	return len(buf), nil
***REMOVED***

func TestCancelReadCloser(t *testing.T) ***REMOVED***
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	cancelReadCloser := NewCancelReadCloser(ctx, ioutil.NopCloser(&perpetualReader***REMOVED******REMOVED***))
	for ***REMOVED***
		var buf [128]byte
		_, err := cancelReadCloser.Read(buf[:])
		if err == context.DeadlineExceeded ***REMOVED***
			break
		***REMOVED*** else if err != nil ***REMOVED***
			t.Fatalf("got unexpected error: %v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***
