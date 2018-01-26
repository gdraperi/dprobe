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
	"errors"
	"io"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/transport"
)

// StreamHandler defines the handler called by gRPC server to complete the
// execution of a streaming RPC.
type StreamHandler func(srv interface***REMOVED******REMOVED***, stream ServerStream) error

// StreamDesc represents a streaming RPC service's method specification.
type StreamDesc struct ***REMOVED***
	StreamName string
	Handler    StreamHandler

	// At least one of these is true.
	ServerStreams bool
	ClientStreams bool
***REMOVED***

// Stream defines the common interface a client or server stream has to satisfy.
type Stream interface ***REMOVED***
	// Context returns the context for this stream.
	Context() context.Context
	// SendMsg blocks until it sends m, the stream is done or the stream
	// breaks.
	// On error, it aborts the stream and returns an RPC status on client
	// side. On server side, it simply returns the error to the caller.
	// SendMsg is called by generated code. Also Users can call SendMsg
	// directly when it is really needed in their use cases.
	SendMsg(m interface***REMOVED******REMOVED***) error
	// RecvMsg blocks until it receives a message or the stream is
	// done. On client side, it returns io.EOF when the stream is done. On
	// any other error, it aborts the stream and returns an RPC status. On
	// server side, it simply returns the error to the caller.
	RecvMsg(m interface***REMOVED******REMOVED***) error
***REMOVED***

// ClientStream defines the interface a client stream has to satisfy.
type ClientStream interface ***REMOVED***
	// Header returns the header metadata received from the server if there
	// is any. It blocks if the metadata is not ready to read.
	Header() (metadata.MD, error)
	// Trailer returns the trailer metadata from the server, if there is any.
	// It must only be called after stream.CloseAndRecv has returned, or
	// stream.Recv has returned a non-nil error (including io.EOF).
	Trailer() metadata.MD
	// CloseSend closes the send direction of the stream. It closes the stream
	// when non-nil error is met.
	CloseSend() error
	Stream
***REMOVED***

// NewClientStream creates a new Stream for the client side. This is called
// by generated code.
func NewClientStream(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (_ ClientStream, err error) ***REMOVED***
	if cc.dopts.streamInt != nil ***REMOVED***
		return cc.dopts.streamInt(ctx, desc, cc, method, newClientStream, opts...)
	***REMOVED***
	return newClientStream(ctx, desc, cc, method, opts...)
***REMOVED***

func newClientStream(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (_ ClientStream, err error) ***REMOVED***
	var (
		t      transport.ClientTransport
		s      *transport.Stream
		put    func()
		cancel context.CancelFunc
	)
	c := defaultCallInfo
	if mc, ok := cc.getMethodConfig(method); ok ***REMOVED***
		c.failFast = !mc.WaitForReady
		if mc.Timeout > 0 ***REMOVED***
			ctx, cancel = context.WithTimeout(ctx, mc.Timeout)
		***REMOVED***
	***REMOVED***
	for _, o := range opts ***REMOVED***
		if err := o.before(&c); err != nil ***REMOVED***
			return nil, toRPCErr(err)
		***REMOVED***
	***REMOVED***
	callHdr := &transport.CallHdr***REMOVED***
		Host:   cc.authority,
		Method: method,
		Flush:  desc.ServerStreams && desc.ClientStreams,
	***REMOVED***
	if cc.dopts.cp != nil ***REMOVED***
		callHdr.SendCompress = cc.dopts.cp.Type()
	***REMOVED***
	var trInfo traceInfo
	if EnableTracing ***REMOVED***
		trInfo.tr = trace.New("grpc.Sent."+methodFamily(method), method)
		trInfo.firstLine.client = true
		if deadline, ok := ctx.Deadline(); ok ***REMOVED***
			trInfo.firstLine.deadline = deadline.Sub(time.Now())
		***REMOVED***
		trInfo.tr.LazyLog(&trInfo.firstLine, false)
		ctx = trace.NewContext(ctx, trInfo.tr)
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				// Need to call tr.finish() if error is returned.
				// Because tr will not be returned to caller.
				trInfo.tr.LazyPrintf("RPC: [%v]", err)
				trInfo.tr.SetError()
				trInfo.tr.Finish()
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
		if err != nil && sh != nil ***REMOVED***
			// Only handle end stats if err != nil.
			end := &stats.End***REMOVED***
				Client: true,
				Error:  err,
			***REMOVED***
			sh.HandleRPC(ctx, end)
		***REMOVED***
	***REMOVED***()
	gopts := BalancerGetOptions***REMOVED***
		BlockingWait: !c.failFast,
	***REMOVED***
	for ***REMOVED***
		t, put, err = cc.getTransport(ctx, gopts)
		if err != nil ***REMOVED***
			// TODO(zhaoq): Probably revisit the error handling.
			if _, ok := status.FromError(err); ok ***REMOVED***
				return nil, err
			***REMOVED***
			if err == errConnClosing || err == errConnUnavailable ***REMOVED***
				if c.failFast ***REMOVED***
					return nil, Errorf(codes.Unavailable, "%v", err)
				***REMOVED***
				continue
			***REMOVED***
			// All the other errors are treated as Internal errors.
			return nil, Errorf(codes.Internal, "%v", err)
		***REMOVED***

		s, err = t.NewStream(ctx, callHdr)
		if err != nil ***REMOVED***
			if _, ok := err.(transport.ConnectionError); ok && put != nil ***REMOVED***
				// If error is connection error, transport was sending data on wire,
				// and we are not sure if anything has been sent on wire.
				// If error is not connection error, we are sure nothing has been sent.
				updateRPCInfoInContext(ctx, rpcInfo***REMOVED***bytesSent: true, bytesReceived: false***REMOVED***)
			***REMOVED***
			if put != nil ***REMOVED***
				put()
				put = nil
			***REMOVED***
			if _, ok := err.(transport.ConnectionError); (ok || err == transport.ErrStreamDrain) && !c.failFast ***REMOVED***
				continue
			***REMOVED***
			return nil, toRPCErr(err)
		***REMOVED***
		break
	***REMOVED***
	cs := &clientStream***REMOVED***
		opts:       opts,
		c:          c,
		desc:       desc,
		codec:      cc.dopts.codec,
		cp:         cc.dopts.cp,
		dc:         cc.dopts.dc,
		maxMsgSize: cc.dopts.maxMsgSize,
		cancel:     cancel,

		put: put,
		t:   t,
		s:   s,
		p:   &parser***REMOVED***r: s***REMOVED***,

		tracing: EnableTracing,
		trInfo:  trInfo,

		statsCtx:     ctx,
		statsHandler: cc.dopts.copts.StatsHandler,
	***REMOVED***
	if cc.dopts.cp != nil ***REMOVED***
		cs.cbuf = new(bytes.Buffer)
	***REMOVED***
	// Listen on ctx.Done() to detect cancellation and s.Done() to detect normal termination
	// when there is no pending I/O operations on this stream.
	go func() ***REMOVED***
		select ***REMOVED***
		case <-t.Error():
			// Incur transport error, simply exit.
		case <-cc.ctx.Done():
			cs.finish(ErrClientConnClosing)
			cs.closeTransportStream(ErrClientConnClosing)
		case <-s.Done():
			// TODO: The trace of the RPC is terminated here when there is no pending
			// I/O, which is probably not the optimal solution.
			cs.finish(s.Status().Err())
			cs.closeTransportStream(nil)
		case <-s.GoAway():
			cs.finish(errConnDrain)
			cs.closeTransportStream(errConnDrain)
		case <-s.Context().Done():
			err := s.Context().Err()
			cs.finish(err)
			cs.closeTransportStream(transport.ContextErr(err))
		***REMOVED***
	***REMOVED***()
	return cs, nil
***REMOVED***

// clientStream implements a client side Stream.
type clientStream struct ***REMOVED***
	opts       []CallOption
	c          callInfo
	t          transport.ClientTransport
	s          *transport.Stream
	p          *parser
	desc       *StreamDesc
	codec      Codec
	cp         Compressor
	cbuf       *bytes.Buffer
	dc         Decompressor
	maxMsgSize int
	cancel     context.CancelFunc

	tracing bool // set to EnableTracing when the clientStream is created.

	mu       sync.Mutex
	put      func()
	closed   bool
	finished bool
	// trInfo.tr is set when the clientStream is created (if EnableTracing is true),
	// and is set to nil when the clientStream's finish method is called.
	trInfo traceInfo

	// statsCtx keeps the user context for stats handling.
	// All stats collection should use the statsCtx (instead of the stream context)
	// so that all the generated stats for a particular RPC can be associated in the processing phase.
	statsCtx     context.Context
	statsHandler stats.Handler
***REMOVED***

func (cs *clientStream) Context() context.Context ***REMOVED***
	return cs.s.Context()
***REMOVED***

func (cs *clientStream) Header() (metadata.MD, error) ***REMOVED***
	m, err := cs.s.Header()
	if err != nil ***REMOVED***
		if _, ok := err.(transport.ConnectionError); !ok ***REMOVED***
			cs.closeTransportStream(err)
		***REMOVED***
	***REMOVED***
	return m, err
***REMOVED***

func (cs *clientStream) Trailer() metadata.MD ***REMOVED***
	return cs.s.Trailer()
***REMOVED***

func (cs *clientStream) SendMsg(m interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	if cs.tracing ***REMOVED***
		cs.mu.Lock()
		if cs.trInfo.tr != nil ***REMOVED***
			cs.trInfo.tr.LazyLog(&payload***REMOVED***sent: true, msg: m***REMOVED***, true)
		***REMOVED***
		cs.mu.Unlock()
	***REMOVED***
	// TODO Investigate how to signal the stats handling party.
	// generate error stats if err != nil && err != io.EOF?
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			cs.finish(err)
		***REMOVED***
		if err == nil ***REMOVED***
			return
		***REMOVED***
		if err == io.EOF ***REMOVED***
			// Specialize the process for server streaming. SendMesg is only called
			// once when creating the stream object. io.EOF needs to be skipped when
			// the rpc is early finished (before the stream object is created.).
			// TODO: It is probably better to move this into the generated code.
			if !cs.desc.ClientStreams && cs.desc.ServerStreams ***REMOVED***
				err = nil
			***REMOVED***
			return
		***REMOVED***
		if _, ok := err.(transport.ConnectionError); !ok ***REMOVED***
			cs.closeTransportStream(err)
		***REMOVED***
		err = toRPCErr(err)
	***REMOVED***()
	var outPayload *stats.OutPayload
	if cs.statsHandler != nil ***REMOVED***
		outPayload = &stats.OutPayload***REMOVED***
			Client: true,
		***REMOVED***
	***REMOVED***
	out, err := encode(cs.codec, m, cs.cp, cs.cbuf, outPayload)
	defer func() ***REMOVED***
		if cs.cbuf != nil ***REMOVED***
			cs.cbuf.Reset()
		***REMOVED***
	***REMOVED***()
	if err != nil ***REMOVED***
		return Errorf(codes.Internal, "grpc: %v", err)
	***REMOVED***
	err = cs.t.Write(cs.s, out, &transport.Options***REMOVED***Last: false***REMOVED***)
	if err == nil && outPayload != nil ***REMOVED***
		outPayload.SentTime = time.Now()
		cs.statsHandler.HandleRPC(cs.statsCtx, outPayload)
	***REMOVED***
	return err
***REMOVED***

func (cs *clientStream) RecvMsg(m interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	var inPayload *stats.InPayload
	if cs.statsHandler != nil ***REMOVED***
		inPayload = &stats.InPayload***REMOVED***
			Client: true,
		***REMOVED***
	***REMOVED***
	err = recv(cs.p, cs.codec, cs.s, cs.dc, m, cs.maxMsgSize, inPayload)
	defer func() ***REMOVED***
		// err != nil indicates the termination of the stream.
		if err != nil ***REMOVED***
			cs.finish(err)
		***REMOVED***
	***REMOVED***()
	if err == nil ***REMOVED***
		if cs.tracing ***REMOVED***
			cs.mu.Lock()
			if cs.trInfo.tr != nil ***REMOVED***
				cs.trInfo.tr.LazyLog(&payload***REMOVED***sent: false, msg: m***REMOVED***, true)
			***REMOVED***
			cs.mu.Unlock()
		***REMOVED***
		if inPayload != nil ***REMOVED***
			cs.statsHandler.HandleRPC(cs.statsCtx, inPayload)
		***REMOVED***
		if !cs.desc.ClientStreams || cs.desc.ServerStreams ***REMOVED***
			return
		***REMOVED***
		// Special handling for client streaming rpc.
		// This recv expects EOF or errors, so we don't collect inPayload.
		err = recv(cs.p, cs.codec, cs.s, cs.dc, m, cs.maxMsgSize, nil)
		cs.closeTransportStream(err)
		if err == nil ***REMOVED***
			return toRPCErr(errors.New("grpc: client streaming protocol violation: get <nil>, want <EOF>"))
		***REMOVED***
		if err == io.EOF ***REMOVED***
			if se := cs.s.Status().Err(); se != nil ***REMOVED***
				return se
			***REMOVED***
			cs.finish(err)
			return nil
		***REMOVED***
		return toRPCErr(err)
	***REMOVED***
	if _, ok := err.(transport.ConnectionError); !ok ***REMOVED***
		cs.closeTransportStream(err)
	***REMOVED***
	if err == io.EOF ***REMOVED***
		if statusErr := cs.s.Status().Err(); statusErr != nil ***REMOVED***
			return statusErr
		***REMOVED***
		// Returns io.EOF to indicate the end of the stream.
		return
	***REMOVED***
	return toRPCErr(err)
***REMOVED***

func (cs *clientStream) CloseSend() (err error) ***REMOVED***
	err = cs.t.Write(cs.s, nil, &transport.Options***REMOVED***Last: true***REMOVED***)
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			cs.finish(err)
		***REMOVED***
	***REMOVED***()
	if err == nil || err == io.EOF ***REMOVED***
		return nil
	***REMOVED***
	if _, ok := err.(transport.ConnectionError); !ok ***REMOVED***
		cs.closeTransportStream(err)
	***REMOVED***
	err = toRPCErr(err)
	return
***REMOVED***

func (cs *clientStream) closeTransportStream(err error) ***REMOVED***
	cs.mu.Lock()
	if cs.closed ***REMOVED***
		cs.mu.Unlock()
		return
	***REMOVED***
	cs.closed = true
	cs.mu.Unlock()
	cs.t.CloseStream(cs.s, err)
***REMOVED***

func (cs *clientStream) finish(err error) ***REMOVED***
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if cs.finished ***REMOVED***
		return
	***REMOVED***
	cs.finished = true
	defer func() ***REMOVED***
		if cs.cancel != nil ***REMOVED***
			cs.cancel()
		***REMOVED***
	***REMOVED***()
	for _, o := range cs.opts ***REMOVED***
		o.after(&cs.c)
	***REMOVED***
	if cs.put != nil ***REMOVED***
		updateRPCInfoInContext(cs.s.Context(), rpcInfo***REMOVED***
			bytesSent:     cs.s.BytesSent(),
			bytesReceived: cs.s.BytesReceived(),
		***REMOVED***)
		cs.put()
		cs.put = nil
	***REMOVED***
	if cs.statsHandler != nil ***REMOVED***
		end := &stats.End***REMOVED***
			Client:  true,
			EndTime: time.Now(),
		***REMOVED***
		if err != io.EOF ***REMOVED***
			// end.Error is nil if the RPC finished successfully.
			end.Error = toRPCErr(err)
		***REMOVED***
		cs.statsHandler.HandleRPC(cs.statsCtx, end)
	***REMOVED***
	if !cs.tracing ***REMOVED***
		return
	***REMOVED***
	if cs.trInfo.tr != nil ***REMOVED***
		if err == nil || err == io.EOF ***REMOVED***
			cs.trInfo.tr.LazyPrintf("RPC: [OK]")
		***REMOVED*** else ***REMOVED***
			cs.trInfo.tr.LazyPrintf("RPC: [%v]", err)
			cs.trInfo.tr.SetError()
		***REMOVED***
		cs.trInfo.tr.Finish()
		cs.trInfo.tr = nil
	***REMOVED***
***REMOVED***

// ServerStream defines the interface a server stream has to satisfy.
type ServerStream interface ***REMOVED***
	// SetHeader sets the header metadata. It may be called multiple times.
	// When call multiple times, all the provided metadata will be merged.
	// All the metadata will be sent out when one of the following happens:
	//  - ServerStream.SendHeader() is called;
	//  - The first response is sent out;
	//  - An RPC status is sent out (error or success).
	SetHeader(metadata.MD) error
	// SendHeader sends the header metadata.
	// The provided md and headers set by SetHeader() will be sent.
	// It fails if called multiple times.
	SendHeader(metadata.MD) error
	// SetTrailer sets the trailer metadata which will be sent with the RPC status.
	// When called more than once, all the provided metadata will be merged.
	SetTrailer(metadata.MD)
	Stream
***REMOVED***

// serverStream implements a server side Stream.
type serverStream struct ***REMOVED***
	t          transport.ServerTransport
	s          *transport.Stream
	p          *parser
	codec      Codec
	cp         Compressor
	dc         Decompressor
	cbuf       *bytes.Buffer
	maxMsgSize int
	trInfo     *traceInfo

	statsHandler stats.Handler

	mu sync.Mutex // protects trInfo.tr after the service handler runs.
***REMOVED***

func (ss *serverStream) Context() context.Context ***REMOVED***
	return ss.s.Context()
***REMOVED***

func (ss *serverStream) SetHeader(md metadata.MD) error ***REMOVED***
	if md.Len() == 0 ***REMOVED***
		return nil
	***REMOVED***
	return ss.s.SetHeader(md)
***REMOVED***

func (ss *serverStream) SendHeader(md metadata.MD) error ***REMOVED***
	return ss.t.WriteHeader(ss.s, md)
***REMOVED***

func (ss *serverStream) SetTrailer(md metadata.MD) ***REMOVED***
	if md.Len() == 0 ***REMOVED***
		return
	***REMOVED***
	ss.s.SetTrailer(md)
	return
***REMOVED***

func (ss *serverStream) SendMsg(m interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if ss.trInfo != nil ***REMOVED***
			ss.mu.Lock()
			if ss.trInfo.tr != nil ***REMOVED***
				if err == nil ***REMOVED***
					ss.trInfo.tr.LazyLog(&payload***REMOVED***sent: true, msg: m***REMOVED***, true)
				***REMOVED*** else ***REMOVED***
					ss.trInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***err***REMOVED******REMOVED***, true)
					ss.trInfo.tr.SetError()
				***REMOVED***
			***REMOVED***
			ss.mu.Unlock()
		***REMOVED***
	***REMOVED***()
	var outPayload *stats.OutPayload
	if ss.statsHandler != nil ***REMOVED***
		outPayload = &stats.OutPayload***REMOVED******REMOVED***
	***REMOVED***
	out, err := encode(ss.codec, m, ss.cp, ss.cbuf, outPayload)
	defer func() ***REMOVED***
		if ss.cbuf != nil ***REMOVED***
			ss.cbuf.Reset()
		***REMOVED***
	***REMOVED***()
	if err != nil ***REMOVED***
		err = Errorf(codes.Internal, "grpc: %v", err)
		return err
	***REMOVED***
	if err := ss.t.Write(ss.s, out, &transport.Options***REMOVED***Last: false***REMOVED***); err != nil ***REMOVED***
		return toRPCErr(err)
	***REMOVED***
	if outPayload != nil ***REMOVED***
		outPayload.SentTime = time.Now()
		ss.statsHandler.HandleRPC(ss.s.Context(), outPayload)
	***REMOVED***
	return nil
***REMOVED***

func (ss *serverStream) RecvMsg(m interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if ss.trInfo != nil ***REMOVED***
			ss.mu.Lock()
			if ss.trInfo.tr != nil ***REMOVED***
				if err == nil ***REMOVED***
					ss.trInfo.tr.LazyLog(&payload***REMOVED***sent: false, msg: m***REMOVED***, true)
				***REMOVED*** else if err != io.EOF ***REMOVED***
					ss.trInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***err***REMOVED******REMOVED***, true)
					ss.trInfo.tr.SetError()
				***REMOVED***
			***REMOVED***
			ss.mu.Unlock()
		***REMOVED***
	***REMOVED***()
	var inPayload *stats.InPayload
	if ss.statsHandler != nil ***REMOVED***
		inPayload = &stats.InPayload***REMOVED******REMOVED***
	***REMOVED***
	if err := recv(ss.p, ss.codec, ss.s, ss.dc, m, ss.maxMsgSize, inPayload); err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			return err
		***REMOVED***
		if err == io.ErrUnexpectedEOF ***REMOVED***
			err = Errorf(codes.Internal, io.ErrUnexpectedEOF.Error())
		***REMOVED***
		return toRPCErr(err)
	***REMOVED***
	if inPayload != nil ***REMOVED***
		ss.statsHandler.HandleRPC(ss.s.Context(), inPayload)
	***REMOVED***
	return nil
***REMOVED***
