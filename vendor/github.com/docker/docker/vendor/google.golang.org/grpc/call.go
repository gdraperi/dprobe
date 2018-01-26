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
	"bytes"
	"io"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/transport"
)

// recvResponse receives and parses an RPC response.
// On error, it returns the error and indicates whether the call should be retried.
//
// TODO(zhaoq): Check whether the received message sequence is valid.
// TODO ctx is used for stats collection and processing. It is the context passed from the application.
func recvResponse(ctx context.Context, dopts dialOptions, t transport.ClientTransport, c *callInfo, stream *transport.Stream, reply interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	// Try to acquire header metadata from the server if there is any.
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if _, ok := err.(transport.ConnectionError); !ok ***REMOVED***
				t.CloseStream(stream, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	c.headerMD, err = stream.Header()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	p := &parser***REMOVED***r: stream***REMOVED***
	var inPayload *stats.InPayload
	if dopts.copts.StatsHandler != nil ***REMOVED***
		inPayload = &stats.InPayload***REMOVED***
			Client: true,
		***REMOVED***
	***REMOVED***
	for ***REMOVED***
		if err = recv(p, dopts.codec, stream, dopts.dc, reply, dopts.maxMsgSize, inPayload); err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	if inPayload != nil && err == io.EOF && stream.Status().Code() == codes.OK ***REMOVED***
		// TODO in the current implementation, inTrailer may be handled before inPayload in some cases.
		// Fix the order if necessary.
		dopts.copts.StatsHandler.HandleRPC(ctx, inPayload)
	***REMOVED***
	c.trailerMD = stream.Trailer()
	if peer, ok := peer.FromContext(stream.Context()); ok ***REMOVED***
		c.peer = peer
	***REMOVED***
	return nil
***REMOVED***

// sendRequest writes out various information of an RPC such as Context and Message.
func sendRequest(ctx context.Context, dopts dialOptions, compressor Compressor, callHdr *transport.CallHdr, stream *transport.Stream, t transport.ClientTransport, args interface***REMOVED******REMOVED***, opts *transport.Options) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			// If err is connection error, t will be closed, no need to close stream here.
			if _, ok := err.(transport.ConnectionError); !ok ***REMOVED***
				t.CloseStream(stream, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	var (
		cbuf       *bytes.Buffer
		outPayload *stats.OutPayload
	)
	if compressor != nil ***REMOVED***
		cbuf = new(bytes.Buffer)
	***REMOVED***
	if dopts.copts.StatsHandler != nil ***REMOVED***
		outPayload = &stats.OutPayload***REMOVED***
			Client: true,
		***REMOVED***
	***REMOVED***
	outBuf, err := encode(dopts.codec, args, compressor, cbuf, outPayload)
	if err != nil ***REMOVED***
		return Errorf(codes.Internal, "grpc: %v", err)
	***REMOVED***
	err = t.Write(stream, outBuf, opts)
	if err == nil && outPayload != nil ***REMOVED***
		outPayload.SentTime = time.Now()
		dopts.copts.StatsHandler.HandleRPC(ctx, outPayload)
	***REMOVED***
	// t.NewStream(...) could lead to an early rejection of the RPC (e.g., the service/method
	// does not exist.) so that t.Write could get io.EOF from wait(...). Leave the following
	// recvResponse to get the final status.
	if err != nil && err != io.EOF ***REMOVED***
		return err
	***REMOVED***
	// Sent successfully.
	return nil
***REMOVED***

// Invoke sends the RPC request on the wire and returns after response is received.
// Invoke is called by generated code. Also users can call Invoke directly when it
// is really needed in their use cases.
func Invoke(ctx context.Context, method string, args, reply interface***REMOVED******REMOVED***, cc *ClientConn, opts ...CallOption) error ***REMOVED***
	if cc.dopts.unaryInt != nil ***REMOVED***
		return cc.dopts.unaryInt(ctx, method, args, reply, cc, invoke, opts...)
	***REMOVED***
	return invoke(ctx, method, args, reply, cc, opts...)
***REMOVED***

func invoke(ctx context.Context, method string, args, reply interface***REMOVED******REMOVED***, cc *ClientConn, opts ...CallOption) (e error) ***REMOVED***
	c := defaultCallInfo
	if mc, ok := cc.getMethodConfig(method); ok ***REMOVED***
		c.failFast = !mc.WaitForReady
		if mc.Timeout > 0 ***REMOVED***
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, mc.Timeout)
			defer cancel()
		***REMOVED***
	***REMOVED***
	for _, o := range opts ***REMOVED***
		if err := o.before(&c); err != nil ***REMOVED***
			return toRPCErr(err)
		***REMOVED***
	***REMOVED***
	defer func() ***REMOVED***
		for _, o := range opts ***REMOVED***
			o.after(&c)
		***REMOVED***
	***REMOVED***()
	if EnableTracing ***REMOVED***
		c.traceInfo.tr = trace.New("grpc.Sent."+methodFamily(method), method)
		defer c.traceInfo.tr.Finish()
		c.traceInfo.firstLine.client = true
		if deadline, ok := ctx.Deadline(); ok ***REMOVED***
			c.traceInfo.firstLine.deadline = deadline.Sub(time.Now())
		***REMOVED***
		c.traceInfo.tr.LazyLog(&c.traceInfo.firstLine, false)
		// TODO(dsymonds): Arrange for c.traceInfo.firstLine.remoteAddr to be set.
		defer func() ***REMOVED***
			if e != nil ***REMOVED***
				c.traceInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***e***REMOVED******REMOVED***, true)
				c.traceInfo.tr.SetError()
			***REMOVED***
		***REMOVED***()
	***REMOVED***
	ctx = newContextWithRPCInfo(ctx)
	sh := cc.dopts.copts.StatsHandler
	if sh != nil ***REMOVED***
		ctx = sh.TagRPC(ctx, &stats.RPCTagInfo***REMOVED***FullMethodName: method***REMOVED***)
		begin := &stats.Begin***REMOVED***
			Client:    true,
			BeginTime: time.Now(),
			FailFast:  c.failFast,
		***REMOVED***
		sh.HandleRPC(ctx, begin)
	***REMOVED***
	defer func() ***REMOVED***
		if sh != nil ***REMOVED***
			end := &stats.End***REMOVED***
				Client:  true,
				EndTime: time.Now(),
				Error:   e,
			***REMOVED***
			sh.HandleRPC(ctx, end)
		***REMOVED***
	***REMOVED***()
	topts := &transport.Options***REMOVED***
		Last:  true,
		Delay: false,
	***REMOVED***
	for ***REMOVED***
		var (
			err    error
			t      transport.ClientTransport
			stream *transport.Stream
			// Record the put handler from Balancer.Get(...). It is called once the
			// RPC has completed or failed.
			put func()
		)
		// TODO(zhaoq): Need a formal spec of fail-fast.
		callHdr := &transport.CallHdr***REMOVED***
			Host:   cc.authority,
			Method: method,
		***REMOVED***
		if cc.dopts.cp != nil ***REMOVED***
			callHdr.SendCompress = cc.dopts.cp.Type()
		***REMOVED***

		gopts := BalancerGetOptions***REMOVED***
			BlockingWait: !c.failFast,
		***REMOVED***
		t, put, err = cc.getTransport(ctx, gopts)
		if err != nil ***REMOVED***
			// TODO(zhaoq): Probably revisit the error handling.
			if _, ok := status.FromError(err); ok ***REMOVED***
				return err
			***REMOVED***
			if err == errConnClosing || err == errConnUnavailable ***REMOVED***
				if c.failFast ***REMOVED***
					return Errorf(codes.Unavailable, "%v", err)
				***REMOVED***
				continue
			***REMOVED***
			// All the other errors are treated as Internal errors.
			return Errorf(codes.Internal, "%v", err)
		***REMOVED***
		if c.traceInfo.tr != nil ***REMOVED***
			c.traceInfo.tr.LazyLog(&payload***REMOVED***sent: true, msg: args***REMOVED***, true)
		***REMOVED***
		stream, err = t.NewStream(ctx, callHdr)
		if err != nil ***REMOVED***
			if put != nil ***REMOVED***
				if _, ok := err.(transport.ConnectionError); ok ***REMOVED***
					// If error is connection error, transport was sending data on wire,
					// and we are not sure if anything has been sent on wire.
					// If error is not connection error, we are sure nothing has been sent.
					updateRPCInfoInContext(ctx, rpcInfo***REMOVED***bytesSent: true, bytesReceived: false***REMOVED***)
				***REMOVED***
				put()
			***REMOVED***
			if _, ok := err.(transport.ConnectionError); (ok || err == transport.ErrStreamDrain) && !c.failFast ***REMOVED***
				continue
			***REMOVED***
			return toRPCErr(err)
		***REMOVED***
		err = sendRequest(ctx, cc.dopts, cc.dopts.cp, callHdr, stream, t, args, topts)
		if err != nil ***REMOVED***
			if put != nil ***REMOVED***
				updateRPCInfoInContext(ctx, rpcInfo***REMOVED***
					bytesSent:     stream.BytesSent(),
					bytesReceived: stream.BytesReceived(),
				***REMOVED***)
				put()
			***REMOVED***
			// Retry a non-failfast RPC when
			// i) there is a connection error; or
			// ii) the server started to drain before this RPC was initiated.
			if _, ok := err.(transport.ConnectionError); (ok || err == transport.ErrStreamDrain) && !c.failFast ***REMOVED***
				continue
			***REMOVED***
			return toRPCErr(err)
		***REMOVED***
		err = recvResponse(ctx, cc.dopts, t, &c, stream, reply)
		if err != nil ***REMOVED***
			if put != nil ***REMOVED***
				updateRPCInfoInContext(ctx, rpcInfo***REMOVED***
					bytesSent:     stream.BytesSent(),
					bytesReceived: stream.BytesReceived(),
				***REMOVED***)
				put()
			***REMOVED***
			if _, ok := err.(transport.ConnectionError); (ok || err == transport.ErrStreamDrain) && !c.failFast ***REMOVED***
				continue
			***REMOVED***
			return toRPCErr(err)
		***REMOVED***
		if c.traceInfo.tr != nil ***REMOVED***
			c.traceInfo.tr.LazyLog(&payload***REMOVED***sent: false, msg: reply***REMOVED***, true)
		***REMOVED***
		t.CloseStream(stream, nil)
		if put != nil ***REMOVED***
			updateRPCInfoInContext(ctx, rpcInfo***REMOVED***
				bytesSent:     stream.BytesSent(),
				bytesReceived: stream.BytesReceived(),
			***REMOVED***)
			put()
		***REMOVED***
		return stream.Status().Err()
	***REMOVED***
***REMOVED***
