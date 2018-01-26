// Protocol Buffers for Go with Gadgets
//
// Copyright (c) 2013, The GoGo Authors. All rights reserved.
// http://github.com/gogo/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package io

import (
	"bufio"
	"encoding/binary"
	"errors"
	"github.com/gogo/protobuf/proto"
	"io"
)

var (
	errSmallBuffer = errors.New("Buffer Too Small")
	errLargeValue  = errors.New("Value is Larger than 64 bits")
)

func NewDelimitedWriter(w io.Writer) WriteCloser ***REMOVED***
	return &varintWriter***REMOVED***w, make([]byte, 10), nil***REMOVED***
***REMOVED***

type varintWriter struct ***REMOVED***
	w      io.Writer
	lenBuf []byte
	buffer []byte
***REMOVED***

func (this *varintWriter) WriteMsg(msg proto.Message) (err error) ***REMOVED***
	var data []byte
	if m, ok := msg.(marshaler); ok ***REMOVED***
		n, ok := getSize(m)
		if !ok ***REMOVED***
			data, err = proto.Marshal(msg)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if n >= len(this.buffer) ***REMOVED***
			this.buffer = make([]byte, n)
		***REMOVED***
		_, err = m.MarshalTo(this.buffer)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		data = this.buffer[:n]
	***REMOVED*** else ***REMOVED***
		data, err = proto.Marshal(msg)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	length := uint64(len(data))
	n := binary.PutUvarint(this.lenBuf, length)
	_, err = this.w.Write(this.lenBuf[:n])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = this.w.Write(data)
	return err
***REMOVED***

func (this *varintWriter) Close() error ***REMOVED***
	if closer, ok := this.w.(io.Closer); ok ***REMOVED***
		return closer.Close()
	***REMOVED***
	return nil
***REMOVED***

func NewDelimitedReader(r io.Reader, maxSize int) ReadCloser ***REMOVED***
	var closer io.Closer
	if c, ok := r.(io.Closer); ok ***REMOVED***
		closer = c
	***REMOVED***
	return &varintReader***REMOVED***bufio.NewReader(r), nil, maxSize, closer***REMOVED***
***REMOVED***

type varintReader struct ***REMOVED***
	r       *bufio.Reader
	buf     []byte
	maxSize int
	closer  io.Closer
***REMOVED***

func (this *varintReader) ReadMsg(msg proto.Message) error ***REMOVED***
	length64, err := binary.ReadUvarint(this.r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	length := int(length64)
	if length < 0 || length > this.maxSize ***REMOVED***
		return io.ErrShortBuffer
	***REMOVED***
	if len(this.buf) < length ***REMOVED***
		this.buf = make([]byte, length)
	***REMOVED***
	buf := this.buf[:length]
	if _, err := io.ReadFull(this.r, buf); err != nil ***REMOVED***
		return err
	***REMOVED***
	return proto.Unmarshal(buf, msg)
***REMOVED***

func (this *varintReader) Close() error ***REMOVED***
	if this.closer != nil ***REMOVED***
		return this.closer.Close()
	***REMOVED***
	return nil
***REMOVED***
