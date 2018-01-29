// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"io"
	"sync"
)

// buffer provides a linked list buffer for data exchange
// between producer and consumer. Theoretically the buffer is
// of unlimited capacity as it does no allocation of its own.
type buffer struct ***REMOVED***
	// protects concurrent access to head, tail and closed
	*sync.Cond

	head *element // the buffer that will be read first
	tail *element // the buffer that will be read last

	closed bool
***REMOVED***

// An element represents a single link in a linked list.
type element struct ***REMOVED***
	buf  []byte
	next *element
***REMOVED***

// newBuffer returns an empty buffer that is not closed.
func newBuffer() *buffer ***REMOVED***
	e := new(element)
	b := &buffer***REMOVED***
		Cond: newCond(),
		head: e,
		tail: e,
	***REMOVED***
	return b
***REMOVED***

// write makes buf available for Read to receive.
// buf must not be modified after the call to write.
func (b *buffer) write(buf []byte) ***REMOVED***
	b.Cond.L.Lock()
	e := &element***REMOVED***buf: buf***REMOVED***
	b.tail.next = e
	b.tail = e
	b.Cond.Signal()
	b.Cond.L.Unlock()
***REMOVED***

// eof closes the buffer. Reads from the buffer once all
// the data has been consumed will receive io.EOF.
func (b *buffer) eof() ***REMOVED***
	b.Cond.L.Lock()
	b.closed = true
	b.Cond.Signal()
	b.Cond.L.Unlock()
***REMOVED***

// Read reads data from the internal buffer in buf.  Reads will block
// if no data is available, or until the buffer is closed.
func (b *buffer) Read(buf []byte) (n int, err error) ***REMOVED***
	b.Cond.L.Lock()
	defer b.Cond.L.Unlock()

	for len(buf) > 0 ***REMOVED***
		// if there is data in b.head, copy it
		if len(b.head.buf) > 0 ***REMOVED***
			r := copy(buf, b.head.buf)
			buf, b.head.buf = buf[r:], b.head.buf[r:]
			n += r
			continue
		***REMOVED***
		// if there is a next buffer, make it the head
		if len(b.head.buf) == 0 && b.head != b.tail ***REMOVED***
			b.head = b.head.next
			continue
		***REMOVED***

		// if at least one byte has been copied, return
		if n > 0 ***REMOVED***
			break
		***REMOVED***

		// if nothing was read, and there is nothing outstanding
		// check to see if the buffer is closed.
		if b.closed ***REMOVED***
			err = io.EOF
			break
		***REMOVED***
		// out of buffers, wait for producer
		b.Cond.Wait()
	***REMOVED***
	return
***REMOVED***
