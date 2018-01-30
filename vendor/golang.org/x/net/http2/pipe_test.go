// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"
)

func TestPipeClose(t *testing.T) ***REMOVED***
	var p pipe
	p.b = new(bytes.Buffer)
	a := errors.New("a")
	b := errors.New("b")
	p.CloseWithError(a)
	p.CloseWithError(b)
	_, err := p.Read(make([]byte, 1))
	if err != a ***REMOVED***
		t.Errorf("err = %v want %v", err, a)
	***REMOVED***
***REMOVED***

func TestPipeDoneChan(t *testing.T) ***REMOVED***
	var p pipe
	done := p.Done()
	select ***REMOVED***
	case <-done:
		t.Fatal("done too soon")
	default:
	***REMOVED***
	p.CloseWithError(io.EOF)
	select ***REMOVED***
	case <-done:
	default:
		t.Fatal("should be done")
	***REMOVED***
***REMOVED***

func TestPipeDoneChan_ErrFirst(t *testing.T) ***REMOVED***
	var p pipe
	p.CloseWithError(io.EOF)
	done := p.Done()
	select ***REMOVED***
	case <-done:
	default:
		t.Fatal("should be done")
	***REMOVED***
***REMOVED***

func TestPipeDoneChan_Break(t *testing.T) ***REMOVED***
	var p pipe
	done := p.Done()
	select ***REMOVED***
	case <-done:
		t.Fatal("done too soon")
	default:
	***REMOVED***
	p.BreakWithError(io.EOF)
	select ***REMOVED***
	case <-done:
	default:
		t.Fatal("should be done")
	***REMOVED***
***REMOVED***

func TestPipeDoneChan_Break_ErrFirst(t *testing.T) ***REMOVED***
	var p pipe
	p.BreakWithError(io.EOF)
	done := p.Done()
	select ***REMOVED***
	case <-done:
	default:
		t.Fatal("should be done")
	***REMOVED***
***REMOVED***

func TestPipeCloseWithError(t *testing.T) ***REMOVED***
	p := &pipe***REMOVED***b: new(bytes.Buffer)***REMOVED***
	const body = "foo"
	io.WriteString(p, body)
	a := errors.New("test error")
	p.CloseWithError(a)
	all, err := ioutil.ReadAll(p)
	if string(all) != body ***REMOVED***
		t.Errorf("read bytes = %q; want %q", all, body)
	***REMOVED***
	if err != a ***REMOVED***
		t.Logf("read error = %v, %v", err, a)
	***REMOVED***
	// Read and Write should fail.
	if n, err := p.Write([]byte("abc")); err != errClosedPipeWrite || n != 0 ***REMOVED***
		t.Errorf("Write(abc) after close\ngot %v, %v\nwant 0, %v", n, err, errClosedPipeWrite)
	***REMOVED***
	if n, err := p.Read(make([]byte, 1)); err == nil || n != 0 ***REMOVED***
		t.Errorf("Read() after close\ngot %v, nil\nwant 0, %v", n, errClosedPipeWrite)
	***REMOVED***
***REMOVED***

func TestPipeBreakWithError(t *testing.T) ***REMOVED***
	p := &pipe***REMOVED***b: new(bytes.Buffer)***REMOVED***
	io.WriteString(p, "foo")
	a := errors.New("test err")
	p.BreakWithError(a)
	all, err := ioutil.ReadAll(p)
	if string(all) != "" ***REMOVED***
		t.Errorf("read bytes = %q; want empty string", all)
	***REMOVED***
	if err != a ***REMOVED***
		t.Logf("read error = %v, %v", err, a)
	***REMOVED***
	if p.b != nil ***REMOVED***
		t.Errorf("buffer should be nil after BreakWithError")
	***REMOVED***
	// Write should succeed silently.
	if n, err := p.Write([]byte("abc")); err != nil || n != 3 ***REMOVED***
		t.Errorf("Write(abc) after break\ngot %v, %v\nwant 0, nil", n, err)
	***REMOVED***
	if p.b != nil ***REMOVED***
		t.Errorf("buffer should be nil after Write")
	***REMOVED***
	// Read should fail.
	if n, err := p.Read(make([]byte, 1)); err == nil || n != 0 ***REMOVED***
		t.Errorf("Read() after close\ngot %v, nil\nwant 0, not nil", n)
	***REMOVED***
***REMOVED***
