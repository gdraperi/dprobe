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
	"github.com/gogo/protobuf/proto"
	"io"
)

func NewFullWriter(w io.Writer) WriteCloser ***REMOVED***
	return &fullWriter***REMOVED***w, nil***REMOVED***
***REMOVED***

type fullWriter struct ***REMOVED***
	w      io.Writer
	buffer []byte
***REMOVED***

func (this *fullWriter) WriteMsg(msg proto.Message) (err error) ***REMOVED***
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
	_, err = this.w.Write(data)
	return err
***REMOVED***

func (this *fullWriter) Close() error ***REMOVED***
	if closer, ok := this.w.(io.Closer); ok ***REMOVED***
		return closer.Close()
	***REMOVED***
	return nil
***REMOVED***

type fullReader struct ***REMOVED***
	r   io.Reader
	buf []byte
***REMOVED***

func NewFullReader(r io.Reader, maxSize int) ReadCloser ***REMOVED***
	return &fullReader***REMOVED***r, make([]byte, maxSize)***REMOVED***
***REMOVED***

func (this *fullReader) ReadMsg(msg proto.Message) error ***REMOVED***
	length, err := this.r.Read(this.buf)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return proto.Unmarshal(this.buf[:length], msg)
***REMOVED***

func (this *fullReader) Close() error ***REMOVED***
	if closer, ok := this.r.(io.Closer); ok ***REMOVED***
		return closer.Close()
	***REMOVED***
	return nil
***REMOVED***
