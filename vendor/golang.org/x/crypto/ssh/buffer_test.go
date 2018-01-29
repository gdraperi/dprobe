// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"io"
	"testing"
)

var alphabet = []byte("abcdefghijklmnopqrstuvwxyz")

func TestBufferReadwrite(t *testing.T) ***REMOVED***
	b := newBuffer()
	b.write(alphabet[:10])
	r, _ := b.Read(make([]byte, 10))
	if r != 10 ***REMOVED***
		t.Fatalf("Expected written == read == 10, written: 10, read %d", r)
	***REMOVED***

	b = newBuffer()
	b.write(alphabet[:5])
	r, _ = b.Read(make([]byte, 10))
	if r != 5 ***REMOVED***
		t.Fatalf("Expected written == read == 5, written: 5, read %d", r)
	***REMOVED***

	b = newBuffer()
	b.write(alphabet[:10])
	r, _ = b.Read(make([]byte, 5))
	if r != 5 ***REMOVED***
		t.Fatalf("Expected written == 10, read == 5, written: 10, read %d", r)
	***REMOVED***

	b = newBuffer()
	b.write(alphabet[:5])
	b.write(alphabet[5:15])
	r, _ = b.Read(make([]byte, 10))
	r2, _ := b.Read(make([]byte, 10))
	if r != 10 || r2 != 5 || 15 != r+r2 ***REMOVED***
		t.Fatal("Expected written == read == 15")
	***REMOVED***
***REMOVED***

func TestBufferClose(t *testing.T) ***REMOVED***
	b := newBuffer()
	b.write(alphabet[:10])
	b.eof()
	_, err := b.Read(make([]byte, 5))
	if err != nil ***REMOVED***
		t.Fatal("expected read of 5 to not return EOF")
	***REMOVED***
	b = newBuffer()
	b.write(alphabet[:10])
	b.eof()
	r, err := b.Read(make([]byte, 5))
	r2, err2 := b.Read(make([]byte, 10))
	if r != 5 || r2 != 5 || err != nil || err2 != nil ***REMOVED***
		t.Fatal("expected reads of 5 and 5")
	***REMOVED***

	b = newBuffer()
	b.write(alphabet[:10])
	b.eof()
	r, err = b.Read(make([]byte, 5))
	r2, err2 = b.Read(make([]byte, 10))
	r3, err3 := b.Read(make([]byte, 10))
	if r != 5 || r2 != 5 || r3 != 0 || err != nil || err2 != nil || err3 != io.EOF ***REMOVED***
		t.Fatal("expected reads of 5 and 5 and 0, with EOF")
	***REMOVED***

	b = newBuffer()
	b.write(make([]byte, 5))
	b.write(make([]byte, 10))
	b.eof()
	r, err = b.Read(make([]byte, 9))
	r2, err2 = b.Read(make([]byte, 3))
	r3, err3 = b.Read(make([]byte, 3))
	r4, err4 := b.Read(make([]byte, 10))
	if err != nil || err2 != nil || err3 != nil || err4 != io.EOF ***REMOVED***
		t.Fatalf("Expected EOF on forth read only, err=%v, err2=%v, err3=%v, err4=%v", err, err2, err3, err4)
	***REMOVED***
	if r != 9 || r2 != 3 || r3 != 3 || r4 != 0 ***REMOVED***
		t.Fatal("Expected written == read == 15", r, r2, r3, r4)
	***REMOVED***
***REMOVED***
