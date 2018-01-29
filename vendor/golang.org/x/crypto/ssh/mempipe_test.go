// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"io"
	"sync"
	"testing"
)

// An in-memory packetConn. It is safe to call Close and writePacket
// from different goroutines.
type memTransport struct ***REMOVED***
	eof     bool
	pending [][]byte
	write   *memTransport
	sync.Mutex
	*sync.Cond
***REMOVED***

func (t *memTransport) readPacket() ([]byte, error) ***REMOVED***
	t.Lock()
	defer t.Unlock()
	for ***REMOVED***
		if len(t.pending) > 0 ***REMOVED***
			r := t.pending[0]
			t.pending = t.pending[1:]
			return r, nil
		***REMOVED***
		if t.eof ***REMOVED***
			return nil, io.EOF
		***REMOVED***
		t.Cond.Wait()
	***REMOVED***
***REMOVED***

func (t *memTransport) closeSelf() error ***REMOVED***
	t.Lock()
	defer t.Unlock()
	if t.eof ***REMOVED***
		return io.EOF
	***REMOVED***
	t.eof = true
	t.Cond.Broadcast()
	return nil
***REMOVED***

func (t *memTransport) Close() error ***REMOVED***
	err := t.write.closeSelf()
	t.closeSelf()
	return err
***REMOVED***

func (t *memTransport) writePacket(p []byte) error ***REMOVED***
	t.write.Lock()
	defer t.write.Unlock()
	if t.write.eof ***REMOVED***
		return io.EOF
	***REMOVED***
	c := make([]byte, len(p))
	copy(c, p)
	t.write.pending = append(t.write.pending, c)
	t.write.Cond.Signal()
	return nil
***REMOVED***

func memPipe() (a, b packetConn) ***REMOVED***
	t1 := memTransport***REMOVED******REMOVED***
	t2 := memTransport***REMOVED******REMOVED***
	t1.write = &t2
	t2.write = &t1
	t1.Cond = sync.NewCond(&t1.Mutex)
	t2.Cond = sync.NewCond(&t2.Mutex)
	return &t1, &t2
***REMOVED***

func TestMemPipe(t *testing.T) ***REMOVED***
	a, b := memPipe()
	if err := a.writePacket([]byte***REMOVED***42***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("writePacket: %v", err)
	***REMOVED***
	if err := a.Close(); err != nil ***REMOVED***
		t.Fatal("Close: ", err)
	***REMOVED***
	p, err := b.readPacket()
	if err != nil ***REMOVED***
		t.Fatal("readPacket: ", err)
	***REMOVED***
	if len(p) != 1 || p[0] != 42 ***REMOVED***
		t.Fatalf("got %v, want ***REMOVED***42***REMOVED***", p)
	***REMOVED***
	p, err = b.readPacket()
	if err != io.EOF ***REMOVED***
		t.Fatalf("got %v, %v, want EOF", p, err)
	***REMOVED***
***REMOVED***

func TestDoubleClose(t *testing.T) ***REMOVED***
	a, _ := memPipe()
	err := a.Close()
	if err != nil ***REMOVED***
		t.Errorf("Close: %v", err)
	***REMOVED***
	err = a.Close()
	if err != io.EOF ***REMOVED***
		t.Errorf("expect EOF on double close.")
	***REMOVED***
***REMOVED***
