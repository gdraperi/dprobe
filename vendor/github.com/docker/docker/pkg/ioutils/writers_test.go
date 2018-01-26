package ioutils

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteCloserWrapperClose(t *testing.T) ***REMOVED***
	called := false
	writer := bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	wrapper := NewWriteCloserWrapper(writer, func() error ***REMOVED***
		called = true
		return nil
	***REMOVED***)
	if err := wrapper.Close(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !called ***REMOVED***
		t.Fatalf("writeCloserWrapper should have call the anonymous function.")
	***REMOVED***
***REMOVED***

func TestNopWriteCloser(t *testing.T) ***REMOVED***
	writer := bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	wrapper := NopWriteCloser(writer)
	if err := wrapper.Close(); err != nil ***REMOVED***
		t.Fatal("NopWriteCloser always return nil on Close.")
	***REMOVED***

***REMOVED***

func TestNopWriter(t *testing.T) ***REMOVED***
	nw := &NopWriter***REMOVED******REMOVED***
	l, err := nw.Write([]byte***REMOVED***'c'***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if l != 1 ***REMOVED***
		t.Fatalf("Expected 1 got %d", l)
	***REMOVED***
***REMOVED***

func TestWriteCounter(t *testing.T) ***REMOVED***
	dummy1 := "This is a dummy string."
	dummy2 := "This is another dummy string."
	totalLength := int64(len(dummy1) + len(dummy2))

	reader1 := strings.NewReader(dummy1)
	reader2 := strings.NewReader(dummy2)

	var buffer bytes.Buffer
	wc := NewWriteCounter(&buffer)

	reader1.WriteTo(wc)
	reader2.WriteTo(wc)

	if wc.Count != totalLength ***REMOVED***
		t.Errorf("Wrong count: %d vs. %d", wc.Count, totalLength)
	***REMOVED***

	if buffer.String() != dummy1+dummy2 ***REMOVED***
		t.Error("Wrong message written")
	***REMOVED***
***REMOVED***
