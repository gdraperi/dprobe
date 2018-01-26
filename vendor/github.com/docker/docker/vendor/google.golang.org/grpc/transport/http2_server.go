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

package transport

import (
	"bytes"
	"errors"
	"io"
	"math"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)

// ErrIllegalHeaderWrite indicates that setting header is illegal because of
// the stream's state.
var ErrIllegalHeaderWrite = errors.New("transport: the stream is done or WriteHeader was already called")

// http2Server implements the ServerTransport interface with HTTP2.
type http2Server struct ***REMOVED***
	ctx         context.Context
	conn        net.Conn
	remoteAddr  net.Addr
	localAddr   net.Addr
	maxStreamID uint32               // max stream ID ever seen
	authInfo    credentials.AuthInfo // auth info about the connection
	inTapHandle tap.ServerInHandle
	// writableChan synchronizes write access to the transport.
	// A writer acquires the write lock by receiving a value on writableChan
	// and releases it by sending on writableChan.
	writableChan chan int
	// shutdownChan is closed when Close is called.
	// Blocking operations should select on shutdownChan to avoid
	// blocking forever after Close.
	shutdownChan chan struct***REMOVED******REMOVED***
	framer       *framer
	hBuf         *bytes.Buffer  // the buffer for HPACK encoding
	hEnc         *hpack.Encoder // HPACK encoder

	// The max number of concurrent streams.
	maxStreams uint32
	// controlBuf delivers all the control related tasks (e.g., window
	// updates, reset streams, and various settings) to the controller.
	controlBuf *recvBuffer
	fc         *inFlow
	// sendQuotaPool provides flow control to outbound message.
	sendQuotaPool *quotaPool

	stats stats.Handler

	// Flag to keep track of reading activity on transport.
	// 1 is true and 0 is false.
	activity uint32 // Accessed atomically.
	// Keepalive and max-age parameters for the server.
	kp keepalive.ServerParameters

	// Keepalive enforcement policy.
	kep keepalive.EnforcementPolicy
	// The time instance last ping was received.
	lastPingAt time.Time
	// Number of times the client has violated keepalive ping policy so far.
	pingStrikes uint8
	// Flag to signify that number of ping strikes should be reset to 0.
	// This is set whenever data or header frames are sent.
	// 1 means yes.
	resetPingStrikes uint32 // Accessed atomically.

	mu            sync.Mutex // guard the following
	state         transportState
	activeStreams map[uint32]*Stream
	// the per-stream outbound flow control window size set by the peer.
	streamSendQuota uint32
	// idle is the time instant when the connection went idle.
	// This is either the begining of the connection or when the number of
	// RPCs go down to 0.
	// When the connection is busy, this value is set to 0.
	idle time.Time
***REMOVED***

// newHTTP2Server constructs a ServerTransport based on HTTP2. ConnectionError is
// returned if something goes wrong.
func newHTTP2Server(conn net.Conn, config *ServerConfig) (_ ServerTransport, err error) ***REMOVED***
	framer := newFramer(conn)
	// Send initial settings as connection preface to client.
	var settings []http2.Setting
	// TODO(zhaoq): Have a better way to signal "no limit" because 0 is
	// permitted in the HTTP2 spec.
	maxStreams := config.MaxStreams
	if maxStreams == 0 ***REMOVED***
		maxStreams = math.MaxUint32
	***REMOVED*** else ***REMOVED***
		settings = append(settings, http2.Setting***REMOVED***
			ID:  http2.SettingMaxConcurrentStreams,
			Val: maxStreams,
		***REMOVED***)
	***REMOVED***
	if initialWindowSize != defaultWindowSize ***REMOVED***
		settings = append(settings, http2.Setting***REMOVED***
			ID:  http2.SettingInitialWindowSize,
			Val: uint32(initialWindowSize)***REMOVED***)
	***REMOVED***
	if err := framer.writeSettings(true, settings...); err != nil ***REMOVED***
		return nil, connectionErrorf(true, err, "transport: %v", err)
	***REMOVED***
	// Adjust the connection flow control window if needed.
	if delta := uint32(initialConnWindowSize - defaultWindowSize); delta > 0 ***REMOVED***
		if err := framer.writeWindowUpdate(true, 0, delta); err != nil ***REMOVED***
			return nil, connectionErrorf(true, err, "transport: %v", err)
		***REMOVED***
	***REMOVED***
	kp := config.KeepaliveParams
	if kp.MaxConnectionIdle == 0 ***REMOVED***
		kp.MaxConnectionIdle = defaultMaxConnectionIdle
	***REMOVED***
	if kp.MaxConnectionAge == 0 ***REMOVED***
		kp.MaxConnectionAge = defaultMaxConnectionAge
	***REMOVED***
	// Add a jitter to MaxConnectionAge.
	kp.MaxConnectionAge += getJitter(kp.MaxConnectionAge)
	if kp.MaxConnectionAgeGrace == 0 ***REMOVED***
		kp.MaxConnectionAgeGrace = defaultMaxConnectionAgeGrace
	***REMOVED***
	if kp.Time == 0 ***REMOVED***
		kp.Time = defaultServerKeepaliveTime
	***REMOVED***
	if kp.Timeout == 0 ***REMOVED***
		kp.Timeout = defaultServerKeepaliveTimeout
	***REMOVED***
	kep := config.KeepalivePolicy
	if kep.MinTime == 0 ***REMOVED***
		kep.MinTime = defaultKeepalivePolicyMinTime
	***REMOVED***
	var buf bytes.Buffer
	t := &http2Server***REMOVED***
		ctx:             context.Background(),
		conn:            conn,
		remoteAddr:      conn.RemoteAddr(),
		localAddr:       conn.LocalAddr(),
		authInfo:        config.AuthInfo,
		framer:          framer,
		hBuf:            &buf,
		hEnc:            hpack.NewEncoder(&buf),
		maxStreams:      maxStreams,
		inTapHandle:     config.InTapHandle,
		controlBuf:      newRecvBuffer(),
		fc:              &inFlow***REMOVED***limit: initialConnWindowSize***REMOVED***,
		sendQuotaPool:   newQuotaPool(defaultWindowSize),
		state:           reachable,
		writableChan:    make(chan int, 1),
		shutdownChan:    make(chan struct***REMOVED******REMOVED***),
		activeStreams:   make(map[uint32]*Stream),
		streamSendQuota: defaultWindowSize,
		stats:           config.StatsHandler,
		kp:              kp,
		idle:            time.Now(),
		kep:             kep,
	***REMOVED***
	if t.stats != nil ***REMOVED***
		t.ctx = t.stats.TagConn(t.ctx, &stats.ConnTagInfo***REMOVED***
			RemoteAddr: t.remoteAddr,
			LocalAddr:  t.localAddr,
		***REMOVED***)
		connBegin := &stats.ConnBegin***REMOVED******REMOVED***
		t.stats.HandleConn(t.ctx, connBegin)
	***REMOVED***
	go t.controller()
	go t.keepalive()
	t.writableChan <- 0
	return t, nil
***REMOVED***

// operateHeader takes action on the decoded headers.
func (t *http2Server) operateHeaders(frame *http2.MetaHeadersFrame, handle func(*Stream), traceCtx func(context.Context, string) context.Context) (close bool) ***REMOVED***
	buf := newRecvBuffer()
	s := &Stream***REMOVED***
		id:  frame.Header().StreamID,
		st:  t,
		buf: buf,
		fc:  &inFlow***REMOVED***limit: initialWindowSize***REMOVED***,
	***REMOVED***

	var state decodeState
	for _, hf := range frame.Fields ***REMOVED***
		if err := state.processHeaderField(hf); err != nil ***REMOVED***
			if se, ok := err.(StreamError); ok ***REMOVED***
				t.controlBuf.put(&resetStream***REMOVED***s.id, statusCodeConvTab[se.Code]***REMOVED***)
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	if frame.StreamEnded() ***REMOVED***
		// s is just created by the caller. No lock needed.
		s.state = streamReadDone
	***REMOVED***
	s.recvCompress = state.encoding
	if state.timeoutSet ***REMOVED***
		s.ctx, s.cancel = context.WithTimeout(t.ctx, state.timeout)
	***REMOVED*** else ***REMOVED***
		s.ctx, s.cancel = context.WithCancel(t.ctx)
	***REMOVED***
	pr := &peer.Peer***REMOVED***
		Addr: t.remoteAddr,
	***REMOVED***
	// Attach Auth info if there is any.
	if t.authInfo != nil ***REMOVED***
		pr.AuthInfo = t.authInfo
	***REMOVED***
	s.ctx = peer.NewContext(s.ctx, pr)
	// Cache the current stream to the context so that the server application
	// can find out. Required when the server wants to send some metadata
	// back to the client (unary call only).
	s.ctx = newContextWithStream(s.ctx, s)
	// Attach the received metadata to the context.
	if len(state.mdata) > 0 ***REMOVED***
		s.ctx = metadata.NewIncomingContext(s.ctx, state.mdata)
	***REMOVED***

	s.dec = &recvBufferReader***REMOVED***
		ctx:  s.ctx,
		recv: s.buf,
	***REMOVED***
	s.recvCompress = state.encoding
	s.method = state.method
	if t.inTapHandle != nil ***REMOVED***
		var err error
		info := &tap.Info***REMOVED***
			FullMethodName: state.method,
		***REMOVED***
		s.ctx, err = t.inTapHandle(s.ctx, info)
		if err != nil ***REMOVED***
			// TODO: Log the real error.
			t.controlBuf.put(&resetStream***REMOVED***s.id, http2.ErrCodeRefusedStream***REMOVED***)
			return
		***REMOVED***
	***REMOVED***
	t.mu.Lock()
	if t.state != reachable ***REMOVED***
		t.mu.Unlock()
		return
	***REMOVED***
	if uint32(len(t.activeStreams)) >= t.maxStreams ***REMOVED***
		t.mu.Unlock()
		t.controlBuf.put(&resetStream***REMOVED***s.id, http2.ErrCodeRefusedStream***REMOVED***)
		return
	***REMOVED***
	if s.id%2 != 1 || s.id <= t.maxStreamID ***REMOVED***
		t.mu.Unlock()
		// illegal gRPC stream id.
		grpclog.Println("transport: http2Server.HandleStreams received an illegal stream id: ", s.id)
		return true
	***REMOVED***
	t.maxStreamID = s.id
	s.sendQuotaPool = newQuotaPool(int(t.streamSendQuota))
	t.activeStreams[s.id] = s
	if len(t.activeStreams) == 1 ***REMOVED***
		t.idle = time.Time***REMOVED******REMOVED***
	***REMOVED***
	t.mu.Unlock()
	s.windowHandler = func(n int) ***REMOVED***
		t.updateWindow(s, uint32(n))
	***REMOVED***
	s.ctx = traceCtx(s.ctx, s.method)
	if t.stats != nil ***REMOVED***
		s.ctx = t.stats.TagRPC(s.ctx, &stats.RPCTagInfo***REMOVED***FullMethodName: s.method***REMOVED***)
		inHeader := &stats.InHeader***REMOVED***
			FullMethod:  s.method,
			RemoteAddr:  t.remoteAddr,
			LocalAddr:   t.localAddr,
			Compression: s.recvCompress,
			WireLength:  int(frame.Header().Length),
		***REMOVED***
		t.stats.HandleRPC(s.ctx, inHeader)
	***REMOVED***
	handle(s)
	return
***REMOVED***

// HandleStreams receives incoming streams using the given handler. This is
// typically run in a separate goroutine.
// traceCtx attaches trace to ctx and returns the new context.
func (t *http2Server) HandleStreams(handle func(*Stream), traceCtx func(context.Context, string) context.Context) ***REMOVED***
	// Check the validity of client preface.
	preface := make([]byte, len(clientPreface))
	if _, err := io.ReadFull(t.conn, preface); err != nil ***REMOVED***
		grpclog.Printf("transport: http2Server.HandleStreams failed to receive the preface from client: %v", err)
		t.Close()
		return
	***REMOVED***
	if !bytes.Equal(preface, clientPreface) ***REMOVED***
		grpclog.Printf("transport: http2Server.HandleStreams received bogus greeting from client: %q", preface)
		t.Close()
		return
	***REMOVED***

	frame, err := t.framer.readFrame()
	if err == io.EOF || err == io.ErrUnexpectedEOF ***REMOVED***
		t.Close()
		return
	***REMOVED***
	if err != nil ***REMOVED***
		grpclog.Printf("transport: http2Server.HandleStreams failed to read frame: %v", err)
		t.Close()
		return
	***REMOVED***
	atomic.StoreUint32(&t.activity, 1)
	sf, ok := frame.(*http2.SettingsFrame)
	if !ok ***REMOVED***
		grpclog.Printf("transport: http2Server.HandleStreams saw invalid preface type %T from client", frame)
		t.Close()
		return
	***REMOVED***
	t.handleSettings(sf)

	for ***REMOVED***
		frame, err := t.framer.readFrame()
		atomic.StoreUint32(&t.activity, 1)
		if err != nil ***REMOVED***
			if se, ok := err.(http2.StreamError); ok ***REMOVED***
				t.mu.Lock()
				s := t.activeStreams[se.StreamID]
				t.mu.Unlock()
				if s != nil ***REMOVED***
					t.closeStream(s)
				***REMOVED***
				t.controlBuf.put(&resetStream***REMOVED***se.StreamID, se.Code***REMOVED***)
				continue
			***REMOVED***
			if err == io.EOF || err == io.ErrUnexpectedEOF ***REMOVED***
				t.Close()
				return
			***REMOVED***
			grpclog.Printf("transport: http2Server.HandleStreams failed to read frame: %v", err)
			t.Close()
			return
		***REMOVED***
		switch frame := frame.(type) ***REMOVED***
		case *http2.MetaHeadersFrame:
			if t.operateHeaders(frame, handle, traceCtx) ***REMOVED***
				t.Close()
				break
			***REMOVED***
		case *http2.DataFrame:
			t.handleData(frame)
		case *http2.RSTStreamFrame:
			t.handleRSTStream(frame)
		case *http2.SettingsFrame:
			t.handleSettings(frame)
		case *http2.PingFrame:
			t.handlePing(frame)
		case *http2.WindowUpdateFrame:
			t.handleWindowUpdate(frame)
		case *http2.GoAwayFrame:
			// TODO: Handle GoAway from the client appropriately.
		default:
			grpclog.Printf("transport: http2Server.HandleStreams found unhandled frame type %v.", frame)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (t *http2Server) getStream(f http2.Frame) (*Stream, bool) ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.activeStreams == nil ***REMOVED***
		// The transport is closing.
		return nil, false
	***REMOVED***
	s, ok := t.activeStreams[f.Header().StreamID]
	if !ok ***REMOVED***
		// The stream is already done.
		return nil, false
	***REMOVED***
	return s, true
***REMOVED***

// updateWindow adjusts the inbound quota for the stream and the transport.
// Window updates will deliver to the controller for sending when
// the cumulative quota exceeds the corresponding threshold.
func (t *http2Server) updateWindow(s *Stream, n uint32) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state == streamDone ***REMOVED***
		return
	***REMOVED***
	if w := t.fc.onRead(n); w > 0 ***REMOVED***
		t.controlBuf.put(&windowUpdate***REMOVED***0, w***REMOVED***)
	***REMOVED***
	if w := s.fc.onRead(n); w > 0 ***REMOVED***
		t.controlBuf.put(&windowUpdate***REMOVED***s.id, w***REMOVED***)
	***REMOVED***
***REMOVED***

func (t *http2Server) handleData(f *http2.DataFrame) ***REMOVED***
	size := f.Header().Length
	if err := t.fc.onData(uint32(size)); err != nil ***REMOVED***
		grpclog.Printf("transport: http2Server %v", err)
		t.Close()
		return
	***REMOVED***
	// Select the right stream to dispatch.
	s, ok := t.getStream(f)
	if !ok ***REMOVED***
		if w := t.fc.onRead(uint32(size)); w > 0 ***REMOVED***
			t.controlBuf.put(&windowUpdate***REMOVED***0, w***REMOVED***)
		***REMOVED***
		return
	***REMOVED***
	if size > 0 ***REMOVED***
		if f.Header().Flags.Has(http2.FlagDataPadded) ***REMOVED***
			if w := t.fc.onRead(uint32(size) - uint32(len(f.Data()))); w > 0 ***REMOVED***
				t.controlBuf.put(&windowUpdate***REMOVED***0, w***REMOVED***)
			***REMOVED***
		***REMOVED***
		s.mu.Lock()
		if s.state == streamDone ***REMOVED***
			s.mu.Unlock()
			// The stream has been closed. Release the corresponding quota.
			if w := t.fc.onRead(uint32(size)); w > 0 ***REMOVED***
				t.controlBuf.put(&windowUpdate***REMOVED***0, w***REMOVED***)
			***REMOVED***
			return
		***REMOVED***
		if err := s.fc.onData(uint32(size)); err != nil ***REMOVED***
			s.mu.Unlock()
			t.closeStream(s)
			t.controlBuf.put(&resetStream***REMOVED***s.id, http2.ErrCodeFlowControl***REMOVED***)
			return
		***REMOVED***
		if f.Header().Flags.Has(http2.FlagDataPadded) ***REMOVED***
			if w := s.fc.onRead(uint32(size) - uint32(len(f.Data()))); w > 0 ***REMOVED***
				t.controlBuf.put(&windowUpdate***REMOVED***s.id, w***REMOVED***)
			***REMOVED***
		***REMOVED***
		s.mu.Unlock()
		// TODO(bradfitz, zhaoq): A copy is required here because there is no
		// guarantee f.Data() is consumed before the arrival of next frame.
		// Can this copy be eliminated?
		if len(f.Data()) > 0 ***REMOVED***
			data := make([]byte, len(f.Data()))
			copy(data, f.Data())
			s.write(recvMsg***REMOVED***data: data***REMOVED***)
		***REMOVED***
	***REMOVED***
	if f.Header().Flags.Has(http2.FlagDataEndStream) ***REMOVED***
		// Received the end of stream from the client.
		s.mu.Lock()
		if s.state != streamDone ***REMOVED***
			s.state = streamReadDone
		***REMOVED***
		s.mu.Unlock()
		s.write(recvMsg***REMOVED***err: io.EOF***REMOVED***)
	***REMOVED***
***REMOVED***

func (t *http2Server) handleRSTStream(f *http2.RSTStreamFrame) ***REMOVED***
	s, ok := t.getStream(f)
	if !ok ***REMOVED***
		return
	***REMOVED***
	t.closeStream(s)
***REMOVED***

func (t *http2Server) handleSettings(f *http2.SettingsFrame) ***REMOVED***
	if f.IsAck() ***REMOVED***
		return
	***REMOVED***
	var ss []http2.Setting
	f.ForeachSetting(func(s http2.Setting) error ***REMOVED***
		ss = append(ss, s)
		return nil
	***REMOVED***)
	// The settings will be applied once the ack is sent.
	t.controlBuf.put(&settings***REMOVED***ack: true, ss: ss***REMOVED***)
***REMOVED***

const (
	maxPingStrikes     = 2
	defaultPingTimeout = 2 * time.Hour
)

func (t *http2Server) handlePing(f *http2.PingFrame) ***REMOVED***
	if f.IsAck() ***REMOVED*** // Do nothing.
		return
	***REMOVED***
	pingAck := &ping***REMOVED***ack: true***REMOVED***
	copy(pingAck.data[:], f.Data[:])
	t.controlBuf.put(pingAck)

	now := time.Now()
	defer func() ***REMOVED***
		t.lastPingAt = now
	***REMOVED***()
	// A reset ping strikes means that we don't need to check for policy
	// violation for this ping and the pingStrikes counter should be set
	// to 0.
	if atomic.CompareAndSwapUint32(&t.resetPingStrikes, 1, 0) ***REMOVED***
		t.pingStrikes = 0
		return
	***REMOVED***
	t.mu.Lock()
	ns := len(t.activeStreams)
	t.mu.Unlock()
	if ns < 1 && !t.kep.PermitWithoutStream ***REMOVED***
		// Keepalive shouldn't be active thus, this new ping should
		// have come after atleast defaultPingTimeout.
		if t.lastPingAt.Add(defaultPingTimeout).After(now) ***REMOVED***
			t.pingStrikes++
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Check if keepalive policy is respected.
		if t.lastPingAt.Add(t.kep.MinTime).After(now) ***REMOVED***
			t.pingStrikes++
		***REMOVED***
	***REMOVED***

	if t.pingStrikes > maxPingStrikes ***REMOVED***
		// Send goaway and close the connection.
		t.controlBuf.put(&goAway***REMOVED***code: http2.ErrCodeEnhanceYourCalm, debugData: []byte("too_many_pings")***REMOVED***)
	***REMOVED***
***REMOVED***

func (t *http2Server) handleWindowUpdate(f *http2.WindowUpdateFrame) ***REMOVED***
	id := f.Header().StreamID
	incr := f.Increment
	if id == 0 ***REMOVED***
		t.sendQuotaPool.add(int(incr))
		return
	***REMOVED***
	if s, ok := t.getStream(f); ok ***REMOVED***
		s.sendQuotaPool.add(int(incr))
	***REMOVED***
***REMOVED***

func (t *http2Server) writeHeaders(s *Stream, b *bytes.Buffer, endStream bool) error ***REMOVED***
	first := true
	endHeaders := false
	var err error
	defer func() ***REMOVED***
		if err == nil ***REMOVED***
			// Reset ping strikes when seding headers since that might cause the
			// peer to send ping.
			atomic.StoreUint32(&t.resetPingStrikes, 1)
		***REMOVED***
	***REMOVED***()
	// Sends the headers in a single batch.
	for !endHeaders ***REMOVED***
		size := t.hBuf.Len()
		if size > http2MaxFrameLen ***REMOVED***
			size = http2MaxFrameLen
		***REMOVED*** else ***REMOVED***
			endHeaders = true
		***REMOVED***
		if first ***REMOVED***
			p := http2.HeadersFrameParam***REMOVED***
				StreamID:      s.id,
				BlockFragment: b.Next(size),
				EndStream:     endStream,
				EndHeaders:    endHeaders,
			***REMOVED***
			err = t.framer.writeHeaders(endHeaders, p)
			first = false
		***REMOVED*** else ***REMOVED***
			err = t.framer.writeContinuation(endHeaders, s.id, endHeaders, b.Next(size))
		***REMOVED***
		if err != nil ***REMOVED***
			t.Close()
			return connectionErrorf(true, err, "transport: %v", err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// WriteHeader sends the header metedata md back to the client.
func (t *http2Server) WriteHeader(s *Stream, md metadata.MD) error ***REMOVED***
	s.mu.Lock()
	if s.headerOk || s.state == streamDone ***REMOVED***
		s.mu.Unlock()
		return ErrIllegalHeaderWrite
	***REMOVED***
	s.headerOk = true
	if md.Len() > 0 ***REMOVED***
		if s.header.Len() > 0 ***REMOVED***
			s.header = metadata.Join(s.header, md)
		***REMOVED*** else ***REMOVED***
			s.header = md
		***REMOVED***
	***REMOVED***
	md = s.header
	s.mu.Unlock()
	if _, err := wait(s.ctx, nil, nil, t.shutdownChan, t.writableChan); err != nil ***REMOVED***
		return err
	***REMOVED***
	t.hBuf.Reset()
	t.hEnc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
	t.hEnc.WriteField(hpack.HeaderField***REMOVED***Name: "content-type", Value: "application/grpc"***REMOVED***)
	if s.sendCompress != "" ***REMOVED***
		t.hEnc.WriteField(hpack.HeaderField***REMOVED***Name: "grpc-encoding", Value: s.sendCompress***REMOVED***)
	***REMOVED***
	for k, vv := range md ***REMOVED***
		if isReservedHeader(k) ***REMOVED***
			// Clients don't tolerate reading restricted headers after some non restricted ones were sent.
			continue
		***REMOVED***
		for _, v := range vv ***REMOVED***
			t.hEnc.WriteField(hpack.HeaderField***REMOVED***Name: k, Value: encodeMetadataHeader(k, v)***REMOVED***)
		***REMOVED***
	***REMOVED***
	bufLen := t.hBuf.Len()
	if err := t.writeHeaders(s, t.hBuf, false); err != nil ***REMOVED***
		return err
	***REMOVED***
	if t.stats != nil ***REMOVED***
		outHeader := &stats.OutHeader***REMOVED***
			WireLength: bufLen,
		***REMOVED***
		t.stats.HandleRPC(s.Context(), outHeader)
	***REMOVED***
	t.writableChan <- 0
	return nil
***REMOVED***

// WriteStatus sends stream status to the client and terminates the stream.
// There is no further I/O operations being able to perform on this stream.
// TODO(zhaoq): Now it indicates the end of entire stream. Revisit if early
// OK is adopted.
func (t *http2Server) WriteStatus(s *Stream, st *status.Status) error ***REMOVED***
	var headersSent, hasHeader bool
	s.mu.Lock()
	if s.state == streamDone ***REMOVED***
		s.mu.Unlock()
		return nil
	***REMOVED***
	if s.headerOk ***REMOVED***
		headersSent = true
	***REMOVED***
	if s.header.Len() > 0 ***REMOVED***
		hasHeader = true
	***REMOVED***
	s.mu.Unlock()

	if !headersSent && hasHeader ***REMOVED***
		t.WriteHeader(s, nil)
		headersSent = true
	***REMOVED***

	if _, err := wait(s.ctx, nil, nil, t.shutdownChan, t.writableChan); err != nil ***REMOVED***
		return err
	***REMOVED***
	t.hBuf.Reset()
	if !headersSent ***REMOVED***
		t.hEnc.WriteField(hpack.HeaderField***REMOVED***Name: ":status", Value: "200"***REMOVED***)
		t.hEnc.WriteField(hpack.HeaderField***REMOVED***Name: "content-type", Value: "application/grpc"***REMOVED***)
	***REMOVED***
	t.hEnc.WriteField(
		hpack.HeaderField***REMOVED***
			Name:  "grpc-status",
			Value: strconv.Itoa(int(st.Code())),
		***REMOVED***)
	t.hEnc.WriteField(hpack.HeaderField***REMOVED***Name: "grpc-message", Value: encodeGrpcMessage(st.Message())***REMOVED***)

	if p := st.Proto(); p != nil && len(p.Details) > 0 ***REMOVED***
		stBytes, err := proto.Marshal(p)
		if err != nil ***REMOVED***
			// TODO: return error instead, when callers are able to handle it.
			panic(err)
		***REMOVED***

		t.hEnc.WriteField(hpack.HeaderField***REMOVED***Name: "grpc-status-details-bin", Value: encodeBinHeader(stBytes)***REMOVED***)
	***REMOVED***

	// Attach the trailer metadata.
	for k, vv := range s.trailer ***REMOVED***
		// Clients don't tolerate reading restricted headers after some non restricted ones were sent.
		if isReservedHeader(k) ***REMOVED***
			continue
		***REMOVED***
		for _, v := range vv ***REMOVED***
			t.hEnc.WriteField(hpack.HeaderField***REMOVED***Name: k, Value: encodeMetadataHeader(k, v)***REMOVED***)
		***REMOVED***
	***REMOVED***
	bufLen := t.hBuf.Len()
	if err := t.writeHeaders(s, t.hBuf, true); err != nil ***REMOVED***
		t.Close()
		return err
	***REMOVED***
	if t.stats != nil ***REMOVED***
		outTrailer := &stats.OutTrailer***REMOVED***
			WireLength: bufLen,
		***REMOVED***
		t.stats.HandleRPC(s.Context(), outTrailer)
	***REMOVED***
	t.closeStream(s)
	t.writableChan <- 0
	return nil
***REMOVED***

// Write converts the data into HTTP2 data frame and sends it out. Non-nil error
// is returns if it fails (e.g., framing error, transport error).
func (t *http2Server) Write(s *Stream, data []byte, opts *Options) (err error) ***REMOVED***
	// TODO(zhaoq): Support multi-writers for a single stream.
	var writeHeaderFrame bool
	s.mu.Lock()
	if s.state == streamDone ***REMOVED***
		s.mu.Unlock()
		return streamErrorf(codes.Unknown, "the stream has been done")
	***REMOVED***
	if !s.headerOk ***REMOVED***
		writeHeaderFrame = true
	***REMOVED***
	s.mu.Unlock()
	if writeHeaderFrame ***REMOVED***
		t.WriteHeader(s, nil)
	***REMOVED***
	defer func() ***REMOVED***
		if err == nil ***REMOVED***
			// Reset ping strikes when sending data since this might cause
			// the peer to send ping.
			atomic.StoreUint32(&t.resetPingStrikes, 1)
		***REMOVED***
	***REMOVED***()
	r := bytes.NewBuffer(data)
	for ***REMOVED***
		if r.Len() == 0 ***REMOVED***
			return nil
		***REMOVED***
		size := http2MaxFrameLen
		// Wait until the stream has some quota to send the data.
		sq, err := wait(s.ctx, nil, nil, t.shutdownChan, s.sendQuotaPool.acquire())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// Wait until the transport has some quota to send the data.
		tq, err := wait(s.ctx, nil, nil, t.shutdownChan, t.sendQuotaPool.acquire())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if sq < size ***REMOVED***
			size = sq
		***REMOVED***
		if tq < size ***REMOVED***
			size = tq
		***REMOVED***
		p := r.Next(size)
		ps := len(p)
		if ps < sq ***REMOVED***
			// Overbooked stream quota. Return it back.
			s.sendQuotaPool.add(sq - ps)
		***REMOVED***
		if ps < tq ***REMOVED***
			// Overbooked transport quota. Return it back.
			t.sendQuotaPool.add(tq - ps)
		***REMOVED***
		t.framer.adjustNumWriters(1)
		// Got some quota. Try to acquire writing privilege on the
		// transport.
		if _, err := wait(s.ctx, nil, nil, t.shutdownChan, t.writableChan); err != nil ***REMOVED***
			if _, ok := err.(StreamError); ok ***REMOVED***
				// Return the connection quota back.
				t.sendQuotaPool.add(ps)
			***REMOVED***
			if t.framer.adjustNumWriters(-1) == 0 ***REMOVED***
				// This writer is the last one in this batch and has the
				// responsibility to flush the buffered frames. It queues
				// a flush request to controlBuf instead of flushing directly
				// in order to avoid the race with other writing or flushing.
				t.controlBuf.put(&flushIO***REMOVED******REMOVED***)
			***REMOVED***
			return err
		***REMOVED***
		select ***REMOVED***
		case <-s.ctx.Done():
			t.sendQuotaPool.add(ps)
			if t.framer.adjustNumWriters(-1) == 0 ***REMOVED***
				t.controlBuf.put(&flushIO***REMOVED******REMOVED***)
			***REMOVED***
			t.writableChan <- 0
			return ContextErr(s.ctx.Err())
		default:
		***REMOVED***
		var forceFlush bool
		if r.Len() == 0 && t.framer.adjustNumWriters(0) == 1 && !opts.Last ***REMOVED***
			forceFlush = true
		***REMOVED***
		if err := t.framer.writeData(forceFlush, s.id, false, p); err != nil ***REMOVED***
			t.Close()
			return connectionErrorf(true, err, "transport: %v", err)
		***REMOVED***
		if t.framer.adjustNumWriters(-1) == 0 ***REMOVED***
			t.framer.flushWrite()
		***REMOVED***
		t.writableChan <- 0
	***REMOVED***

***REMOVED***

func (t *http2Server) applySettings(ss []http2.Setting) ***REMOVED***
	for _, s := range ss ***REMOVED***
		if s.ID == http2.SettingInitialWindowSize ***REMOVED***
			t.mu.Lock()
			defer t.mu.Unlock()
			for _, stream := range t.activeStreams ***REMOVED***
				stream.sendQuotaPool.add(int(s.Val - t.streamSendQuota))
			***REMOVED***
			t.streamSendQuota = s.Val
		***REMOVED***

	***REMOVED***
***REMOVED***

// keepalive running in a separate goroutine does the following:
// 1. Gracefully closes an idle connection after a duration of keepalive.MaxConnectionIdle.
// 2. Gracefully closes any connection after a duration of keepalive.MaxConnectionAge.
// 3. Forcibly closes a connection after an additive period of keepalive.MaxConnectionAgeGrace over keepalive.MaxConnectionAge.
// 4. Makes sure a connection is alive by sending pings with a frequency of keepalive.Time and closes a non-resposive connection
// after an additional duration of keepalive.Timeout.
func (t *http2Server) keepalive() ***REMOVED***
	p := &ping***REMOVED******REMOVED***
	var pingSent bool
	maxIdle := time.NewTimer(t.kp.MaxConnectionIdle)
	maxAge := time.NewTimer(t.kp.MaxConnectionAge)
	keepalive := time.NewTimer(t.kp.Time)
	// NOTE: All exit paths of this function should reset their
	// respecitve timers. A failure to do so will cause the
	// following clean-up to deadlock and eventually leak.
	defer func() ***REMOVED***
		if !maxIdle.Stop() ***REMOVED***
			<-maxIdle.C
		***REMOVED***
		if !maxAge.Stop() ***REMOVED***
			<-maxAge.C
		***REMOVED***
		if !keepalive.Stop() ***REMOVED***
			<-keepalive.C
		***REMOVED***
	***REMOVED***()
	for ***REMOVED***
		select ***REMOVED***
		case <-maxIdle.C:
			t.mu.Lock()
			idle := t.idle
			if idle.IsZero() ***REMOVED*** // The connection is non-idle.
				t.mu.Unlock()
				maxIdle.Reset(t.kp.MaxConnectionIdle)
				continue
			***REMOVED***
			val := t.kp.MaxConnectionIdle - time.Since(idle)
			if val <= 0 ***REMOVED***
				// The connection has been idle for a duration of keepalive.MaxConnectionIdle or more.
				// Gracefully close the connection.
				t.state = draining
				t.mu.Unlock()
				t.Drain()
				// Reseting the timer so that the clean-up doesn't deadlock.
				maxIdle.Reset(infinity)
				return
			***REMOVED***
			t.mu.Unlock()
			maxIdle.Reset(val)
		case <-maxAge.C:
			t.mu.Lock()
			t.state = draining
			t.mu.Unlock()
			t.Drain()
			maxAge.Reset(t.kp.MaxConnectionAgeGrace)
			select ***REMOVED***
			case <-maxAge.C:
				// Close the connection after grace period.
				t.Close()
				// Reseting the timer so that the clean-up doesn't deadlock.
				maxAge.Reset(infinity)
			case <-t.shutdownChan:
			***REMOVED***
			return
		case <-keepalive.C:
			if atomic.CompareAndSwapUint32(&t.activity, 1, 0) ***REMOVED***
				pingSent = false
				keepalive.Reset(t.kp.Time)
				continue
			***REMOVED***
			if pingSent ***REMOVED***
				t.Close()
				// Reseting the timer so that the clean-up doesn't deadlock.
				keepalive.Reset(infinity)
				return
			***REMOVED***
			pingSent = true
			t.controlBuf.put(p)
			keepalive.Reset(t.kp.Timeout)
		case <-t.shutdownChan:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// controller running in a separate goroutine takes charge of sending control
// frames (e.g., window update, reset stream, setting, etc.) to the server.
func (t *http2Server) controller() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case i := <-t.controlBuf.get():
			t.controlBuf.load()
			select ***REMOVED***
			case <-t.writableChan:
				switch i := i.(type) ***REMOVED***
				case *windowUpdate:
					t.framer.writeWindowUpdate(true, i.streamID, i.increment)
				case *settings:
					if i.ack ***REMOVED***
						t.framer.writeSettingsAck(true)
						t.applySettings(i.ss)
					***REMOVED*** else ***REMOVED***
						t.framer.writeSettings(true, i.ss...)
					***REMOVED***
				case *resetStream:
					t.framer.writeRSTStream(true, i.streamID, i.code)
				case *goAway:
					t.mu.Lock()
					if t.state == closing ***REMOVED***
						t.mu.Unlock()
						// The transport is closing.
						return
					***REMOVED***
					sid := t.maxStreamID
					t.state = draining
					t.mu.Unlock()
					t.framer.writeGoAway(true, sid, i.code, i.debugData)
					if i.code == http2.ErrCodeEnhanceYourCalm ***REMOVED***
						t.Close()
					***REMOVED***
				case *flushIO:
					t.framer.flushWrite()
				case *ping:
					t.framer.writePing(true, i.ack, i.data)
				default:
					grpclog.Printf("transport: http2Server.controller got unexpected item type %v\n", i)
				***REMOVED***
				t.writableChan <- 0
				continue
			case <-t.shutdownChan:
				return
			***REMOVED***
		case <-t.shutdownChan:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Close starts shutting down the http2Server transport.
// TODO(zhaoq): Now the destruction is not blocked on any pending streams. This
// could cause some resource issue. Revisit this later.
func (t *http2Server) Close() (err error) ***REMOVED***
	t.mu.Lock()
	if t.state == closing ***REMOVED***
		t.mu.Unlock()
		return errors.New("transport: Close() was already called")
	***REMOVED***
	t.state = closing
	streams := t.activeStreams
	t.activeStreams = nil
	t.mu.Unlock()
	close(t.shutdownChan)
	err = t.conn.Close()
	// Cancel all active streams.
	for _, s := range streams ***REMOVED***
		s.cancel()
	***REMOVED***
	if t.stats != nil ***REMOVED***
		connEnd := &stats.ConnEnd***REMOVED******REMOVED***
		t.stats.HandleConn(t.ctx, connEnd)
	***REMOVED***
	return
***REMOVED***

// closeStream clears the footprint of a stream when the stream is not needed
// any more.
func (t *http2Server) closeStream(s *Stream) ***REMOVED***
	t.mu.Lock()
	delete(t.activeStreams, s.id)
	if len(t.activeStreams) == 0 ***REMOVED***
		t.idle = time.Now()
	***REMOVED***
	if t.state == draining && len(t.activeStreams) == 0 ***REMOVED***
		defer t.Close()
	***REMOVED***
	t.mu.Unlock()
	// In case stream sending and receiving are invoked in separate
	// goroutines (e.g., bi-directional streaming), cancel needs to be
	// called to interrupt the potential blocking on other goroutines.
	s.cancel()
	s.mu.Lock()
	if q := s.fc.resetPendingData(); q > 0 ***REMOVED***
		if w := t.fc.onRead(q); w > 0 ***REMOVED***
			t.controlBuf.put(&windowUpdate***REMOVED***0, w***REMOVED***)
		***REMOVED***
	***REMOVED***
	if s.state == streamDone ***REMOVED***
		s.mu.Unlock()
		return
	***REMOVED***
	s.state = streamDone
	s.mu.Unlock()
***REMOVED***

func (t *http2Server) RemoteAddr() net.Addr ***REMOVED***
	return t.remoteAddr
***REMOVED***

func (t *http2Server) Drain() ***REMOVED***
	t.controlBuf.put(&goAway***REMOVED***code: http2.ErrCodeNo***REMOVED***)
***REMOVED***

var rgen = rand.New(rand.NewSource(time.Now().UnixNano()))

func getJitter(v time.Duration) time.Duration ***REMOVED***
	if v == infinity ***REMOVED***
		return 0
	***REMOVED***
	// Generate a jitter between +/- 10% of the value.
	r := int64(v / 10)
	j := rgen.Int63n(2*r) - r
	return time.Duration(j)
***REMOVED***
