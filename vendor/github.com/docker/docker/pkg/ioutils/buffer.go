package ioutils

import (
	"errors"
	"io"
)

var errBufferFull = errors.New("buffer is full")

type fixedBuffer struct ***REMOVED***
	buf      []byte
	pos      int
	lastRead int
***REMOVED***

func (b *fixedBuffer) Write(p []byte) (int, error) ***REMOVED***
	n := copy(b.buf[b.pos:cap(b.buf)], p)
	b.pos += n

	if n < len(p) ***REMOVED***
		if b.pos == cap(b.buf) ***REMOVED***
			return n, errBufferFull
		***REMOVED***
		return n, io.ErrShortWrite
	***REMOVED***
	return n, nil
***REMOVED***

func (b *fixedBuffer) Read(p []byte) (int, error) ***REMOVED***
	n := copy(p, b.buf[b.lastRead:b.pos])
	b.lastRead += n
	return n, nil
***REMOVED***

func (b *fixedBuffer) Len() int ***REMOVED***
	return b.pos - b.lastRead
***REMOVED***

func (b *fixedBuffer) Cap() int ***REMOVED***
	return cap(b.buf)
***REMOVED***

func (b *fixedBuffer) Reset() ***REMOVED***
	b.pos = 0
	b.lastRead = 0
	b.buf = b.buf[:0]
***REMOVED***

func (b *fixedBuffer) String() string ***REMOVED***
	return string(b.buf[b.lastRead:b.pos])
***REMOVED***
