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
	"encoding/binary"
	"github.com/gogo/protobuf/proto"
	"io"
)

func NewUint32DelimitedWriter(w io.Writer, byteOrder binary.ByteOrder) WriteCloser ***REMOVED***
	return &uint32Writer***REMOVED***w, byteOrder, nil***REMOVED***
***REMOVED***

func NewSizeUint32DelimitedWriter(w io.Writer, byteOrder binary.ByteOrder, size int) WriteCloser ***REMOVED***
	return &uint32Writer***REMOVED***w, byteOrder, make([]byte, size)***REMOVED***
***REMOVED***

type uint32Writer struct ***REMOVED***
	w         io.Writer
	byteOrder binary.ByteOrder
	buffer    []byte
***REMOVED***

func (this *uint32Writer) WriteMsg(msg proto.Message) (err error) ***REMOVED***
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
	length := uint32(len(data))
	if err = binary.Write(this.w, this.byteOrder, &length); err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = this.w.Write(data)
	return err
***REMOVED***

func (this *uint32Writer) Close() error ***REMOVED***
	if closer, ok := this.w.(io.Closer); ok ***REMOVED***
		return closer.Close()
	***REMOVED***
	return nil
***REMOVED***

type uint32Reader struct ***REMOVED***
	r         io.Reader
	byteOrder binary.ByteOrder
	lenBuf    []byte
	buf       []byte
	maxSize   int
***REMOVED***

func NewUint32DelimitedReader(r io.Reader, byteOrder binary.ByteOrder, maxSize int) ReadCloser ***REMOVED***
	return &uint32Reader***REMOVED***r, byteOrder, make([]byte, 4), nil, maxSize***REMOVED***
***REMOVED***

func (this *uint32Reader) ReadMsg(msg proto.Message) error ***REMOVED***
	if _, err := io.ReadFull(this.r, this.lenBuf); err != nil ***REMOVED***
		return err
	***REMOVED***
	length32 := this.byteOrder.Uint32(this.lenBuf)
	length := int(length32)
	if length < 0 || length > this.maxSize ***REMOVED***
		return io.ErrShortBuffer
	***REMOVED***
	if length >= len(this.buf) ***REMOVED***
		this.buf = make([]byte, length)
	***REMOVED***
	_, err := io.ReadFull(this.r, this.buf[:length])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return proto.Unmarshal(this.buf[:length], msg)
***REMOVED***

func (this *uint32Reader) Close() error ***REMOVED***
	if closer, ok := this.r.(io.Closer); ok ***REMOVED***
		return closer.Close()
	***REMOVED***
	return nil
***REMOVED***
