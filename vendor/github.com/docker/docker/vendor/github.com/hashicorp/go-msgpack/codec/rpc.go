// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"bufio"
	"io"
	"net/rpc"
	"sync"
)

// Rpc provides a rpc Server or Client Codec for rpc communication.
type Rpc interface ***REMOVED***
	ServerCodec(conn io.ReadWriteCloser, h Handle) rpc.ServerCodec
	ClientCodec(conn io.ReadWriteCloser, h Handle) rpc.ClientCodec
***REMOVED***

// RpcCodecBuffered allows access to the underlying bufio.Reader/Writer
// used by the rpc connection. It accomodates use-cases where the connection
// should be used by rpc and non-rpc functions, e.g. streaming a file after
// sending an rpc response.
type RpcCodecBuffered interface ***REMOVED***
	BufferedReader() *bufio.Reader
	BufferedWriter() *bufio.Writer
***REMOVED***

// -------------------------------------

// rpcCodec defines the struct members and common methods.
type rpcCodec struct ***REMOVED***
	rwc io.ReadWriteCloser
	dec *Decoder
	enc *Encoder
	bw  *bufio.Writer
	br  *bufio.Reader
	mu  sync.Mutex
	cls bool
***REMOVED***

func newRPCCodec(conn io.ReadWriteCloser, h Handle) rpcCodec ***REMOVED***
	bw := bufio.NewWriter(conn)
	br := bufio.NewReader(conn)
	return rpcCodec***REMOVED***
		rwc: conn,
		bw:  bw,
		br:  br,
		enc: NewEncoder(bw, h),
		dec: NewDecoder(br, h),
	***REMOVED***
***REMOVED***

func (c *rpcCodec) BufferedReader() *bufio.Reader ***REMOVED***
	return c.br
***REMOVED***

func (c *rpcCodec) BufferedWriter() *bufio.Writer ***REMOVED***
	return c.bw
***REMOVED***

func (c *rpcCodec) write(obj1, obj2 interface***REMOVED******REMOVED***, writeObj2, doFlush bool) (err error) ***REMOVED***
	if c.cls ***REMOVED***
		return io.EOF
	***REMOVED***
	if err = c.enc.Encode(obj1); err != nil ***REMOVED***
		return
	***REMOVED***
	if writeObj2 ***REMOVED***
		if err = c.enc.Encode(obj2); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	if doFlush && c.bw != nil ***REMOVED***
		return c.bw.Flush()
	***REMOVED***
	return
***REMOVED***

func (c *rpcCodec) read(obj interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	if c.cls ***REMOVED***
		return io.EOF
	***REMOVED***
	//If nil is passed in, we should still attempt to read content to nowhere.
	if obj == nil ***REMOVED***
		var obj2 interface***REMOVED******REMOVED***
		return c.dec.Decode(&obj2)
	***REMOVED***
	return c.dec.Decode(obj)
***REMOVED***

func (c *rpcCodec) Close() error ***REMOVED***
	if c.cls ***REMOVED***
		return io.EOF
	***REMOVED***
	c.cls = true
	return c.rwc.Close()
***REMOVED***

func (c *rpcCodec) ReadResponseBody(body interface***REMOVED******REMOVED***) error ***REMOVED***
	return c.read(body)
***REMOVED***

// -------------------------------------

type goRpcCodec struct ***REMOVED***
	rpcCodec
***REMOVED***

func (c *goRpcCodec) WriteRequest(r *rpc.Request, body interface***REMOVED******REMOVED***) error ***REMOVED***
	// Must protect for concurrent access as per API
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.write(r, body, true, true)
***REMOVED***

func (c *goRpcCodec) WriteResponse(r *rpc.Response, body interface***REMOVED******REMOVED***) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.write(r, body, true, true)
***REMOVED***

func (c *goRpcCodec) ReadResponseHeader(r *rpc.Response) error ***REMOVED***
	return c.read(r)
***REMOVED***

func (c *goRpcCodec) ReadRequestHeader(r *rpc.Request) error ***REMOVED***
	return c.read(r)
***REMOVED***

func (c *goRpcCodec) ReadRequestBody(body interface***REMOVED******REMOVED***) error ***REMOVED***
	return c.read(body)
***REMOVED***

// -------------------------------------

// goRpc is the implementation of Rpc that uses the communication protocol
// as defined in net/rpc package.
type goRpc struct***REMOVED******REMOVED***

// GoRpc implements Rpc using the communication protocol defined in net/rpc package.
// Its methods (ServerCodec and ClientCodec) return values that implement RpcCodecBuffered.
var GoRpc goRpc

func (x goRpc) ServerCodec(conn io.ReadWriteCloser, h Handle) rpc.ServerCodec ***REMOVED***
	return &goRpcCodec***REMOVED***newRPCCodec(conn, h)***REMOVED***
***REMOVED***

func (x goRpc) ClientCodec(conn io.ReadWriteCloser, h Handle) rpc.ClientCodec ***REMOVED***
	return &goRpcCodec***REMOVED***newRPCCodec(conn, h)***REMOVED***
***REMOVED***

var _ RpcCodecBuffered = (*rpcCodec)(nil) // ensure *rpcCodec implements RpcCodecBuffered
