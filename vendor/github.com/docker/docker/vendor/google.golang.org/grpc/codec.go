/*
*
 * Copyright 2014, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
*/

package grpc

import (
	"math"
	"sync"

	"github.com/golang/protobuf/proto"
)

// Codec defines the interface gRPC uses to encode and decode messages.
// Note that implementations of this interface must be thread safe;
// a Codec's methods can be called from concurrent goroutines.
type Codec interface ***REMOVED***
	// Marshal returns the wire format of v.
	Marshal(v interface***REMOVED******REMOVED***) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface***REMOVED******REMOVED***) error
	// String returns the name of the Codec implementation. The returned
	// string will be used as part of content type in transmission.
	String() string
***REMOVED***

// protoCodec is a Codec implementation with protobuf. It is the default codec for gRPC.
type protoCodec struct ***REMOVED***
***REMOVED***

type cachedProtoBuffer struct ***REMOVED***
	lastMarshaledSize uint32
	proto.Buffer
***REMOVED***

func capToMaxInt32(val int) uint32 ***REMOVED***
	if val > math.MaxInt32 ***REMOVED***
		return uint32(math.MaxInt32)
	***REMOVED***
	return uint32(val)
***REMOVED***

func (p protoCodec) marshal(v interface***REMOVED******REMOVED***, cb *cachedProtoBuffer) ([]byte, error) ***REMOVED***
	protoMsg := v.(proto.Message)
	newSlice := make([]byte, 0, cb.lastMarshaledSize)

	cb.SetBuf(newSlice)
	cb.Reset()
	if err := cb.Marshal(protoMsg); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	out := cb.Bytes()
	cb.lastMarshaledSize = capToMaxInt32(len(out))
	return out, nil
***REMOVED***

func (p protoCodec) Marshal(v interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	cb := protoBufferPool.Get().(*cachedProtoBuffer)
	out, err := p.marshal(v, cb)

	// put back buffer and lose the ref to the slice
	cb.SetBuf(nil)
	protoBufferPool.Put(cb)
	return out, err
***REMOVED***

func (p protoCodec) Unmarshal(data []byte, v interface***REMOVED******REMOVED***) error ***REMOVED***
	cb := protoBufferPool.Get().(*cachedProtoBuffer)
	cb.SetBuf(data)
	err := cb.Unmarshal(v.(proto.Message))
	cb.SetBuf(nil)
	protoBufferPool.Put(cb)
	return err
***REMOVED***

func (protoCodec) String() string ***REMOVED***
	return "proto"
***REMOVED***

var (
	protoBufferPool = &sync.Pool***REMOVED***
		New: func() interface***REMOVED******REMOVED*** ***REMOVED***
			return &cachedProtoBuffer***REMOVED***
				Buffer:            proto.Buffer***REMOVED******REMOVED***,
				lastMarshaledSize: 16,
			***REMOVED***
		***REMOVED***,
	***REMOVED***
)
