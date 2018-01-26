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

type Writer interface ***REMOVED***
	WriteMsg(proto.Message) error
***REMOVED***

type WriteCloser interface ***REMOVED***
	Writer
	io.Closer
***REMOVED***

type Reader interface ***REMOVED***
	ReadMsg(msg proto.Message) error
***REMOVED***

type ReadCloser interface ***REMOVED***
	Reader
	io.Closer
***REMOVED***

type marshaler interface ***REMOVED***
	MarshalTo(data []byte) (n int, err error)
***REMOVED***

func getSize(v interface***REMOVED******REMOVED***) (int, bool) ***REMOVED***
	if sz, ok := v.(interface ***REMOVED***
		Size() (n int)
	***REMOVED***); ok ***REMOVED***
		return sz.Size(), true
	***REMOVED*** else if sz, ok := v.(interface ***REMOVED***
		ProtoSize() (n int)
	***REMOVED***); ok ***REMOVED***
		return sz.ProtoSize(), true
	***REMOVED*** else ***REMOVED***
		return 0, false
	***REMOVED***
***REMOVED***
