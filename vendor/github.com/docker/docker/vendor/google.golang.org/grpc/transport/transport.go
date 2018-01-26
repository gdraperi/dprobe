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

/*
Package transport defines and implements message oriented communication channel
to complete various transactions (e.g., an RPC).
*/
package transport // import "google.golang.org/grpc/transport"

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"

	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)

// recvMsg represents the received msg from the transport. All transport
// protocol specific info has been removed.
type recvMsg struct ***REMOVED***
	data []byte
	// nil: received some data
	// io.EOF: stream is completed. data is nil.
	// other non-nil error: transport failure. data is nil.
	err error
***REMOVED***

func (*recvMsg) item() ***REMOVED******REMOVED***

// All items in an out of a recvBuffer should be the same type.
type item interface ***REMOVED***
	item()
***REMOVED***

// recvBuffer is an unbounded channel of item.
type recvBuffer struct ***REMOVED***
	c       chan item
	mu      sync.Mutex
	backlog []item
***REMOVED***

func newRecvBuffer() *recvBuffer ***REMOVED***
	b := &recvBuffer***REMOVED***
		c: make(chan item, 1),
	***REMOVED***
	return b
***REMOVED***

func (b *recvBuffer) put(r item) ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.backlog) == 0 ***REMOVED***
		select ***REMOVED***
		case b.c <- r:
			return
		default:
		***REMOVED***
	***REMOVED***
	b.backlog = append(b.backlog, r)
***REMOVED***

func (b *recvBuffer) load() ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.backlog) > 0 ***REMOVED***
		select ***REMOVED***
		case b.c <- b.backlog[0]:
			b.backlog = b.backlog[1:]
		default:
		***REMOVED***
	***REMOVED***
***REMOVED***

// get returns the channel that receives an item in the buffer.
//
// Upon receipt of an item, the caller should call load to send another
// item onto the channel if there is any.
func (b *recvBuffer) get() <-chan item ***REMOVED***
	return b.c
***REMOVED***

// recvBufferReader implements io.Reader interface to read the data from
// recvBuffer.
type recvBufferReader struct ***REMOVED***
	ctx    context.Context
	goAway chan struct***REMOVED******REMOVED***
	recv   *recvBuffer
	last   *bytes.Reader // Stores the remaining data in the previous calls.
	err    error
***REMOVED***

// Read reads the next len(p) bytes from last. If last is drained, it tries to
// read additional data from recv. It blocks if there no additional data available
// in recv. If Read returns any non-nil error, it will continue to return that error.
func (r *recvBufferReader) Read(p []byte) (n int, err error) ***REMOVED***
	if r.err != nil ***REMOVED***
		return 0, r.err
	***REMOVED***
	defer func() ***REMOVED*** r.err = err ***REMOVED***()
	if r.last != nil && r.last.Len() > 0 ***REMOVED***
		// Read remaining data left in last call.
		return r.last.Read(p)
	***REMOVED***
	select ***REMOVED***
	case <-r.ctx.Done():
		return 0, ContextErr(r.ctx.Err())
	case <-r.goAway:
		return 0, ErrStreamDrain
	case i := <-r.recv.get():
		r.recv.load()
		m := i.(*recvMsg)
		if m.err != nil ***REMOVED***
			return 0, m.err
		***REMOVED***
		r.last = bytes.NewReader(m.data)
		return r.last.Read(p)
	***REMOVED***
***REMOVED***

type streamState uint8

const (
	streamActive    streamState = iota
	streamWriteDone             // EndStream sent
	streamReadDone              // EndStream received
	streamDone                  // the entire stream is finished.
)

// Stream represents an RPC in the transport layer.
type Stream struct ***REMOVED***
	id uint32
	// nil for client side Stream.
	st ServerTransport
	// clientStatsCtx keeps the user context for stats handling.
	// It's only valid on client side. Server side stats context is same as s.ctx.
	// All client side stats collection should use the clientStatsCtx (instead of the stream context)
	// so that all the generated stats for a particular RPC can be associated in the processing phase.
	clientStatsCtx context.Context
	// ctx is the associated context of the stream.
	ctx context.Context
	// cancel is always nil for client side Stream.
	cancel context.CancelFunc
	// done is closed when the final status arrives.
	done chan struct***REMOVED******REMOVED***
	// goAway is closed when the server sent GoAways signal before this stream was initiated.
	goAway chan struct***REMOVED******REMOVED***
	// method records the associated RPC method of the stream.
	method       string
	recvCompress string
	sendCompress string
	buf          *recvBuffer
	dec          io.Reader
	fc           *inFlow
	recvQuota    uint32
	// The accumulated inbound quota pending for window update.
	updateQuota uint32
	// The handler to control the window update procedure for both this
	// particular stream and the associated transport.
	windowHandler func(int)

	sendQuotaPool *quotaPool
	// Close headerChan to indicate the end of reception of header metadata.
	headerChan chan struct***REMOVED******REMOVED***
	// header caches the received header metadata.
	header metadata.MD
	// The key-value map of trailer metadata.
	trailer metadata.MD

	mu sync.RWMutex // guard the following
	// headerOK becomes true from the first header is about to send.
	headerOk bool
	state    streamState
	// true iff headerChan is closed. Used to avoid closing headerChan
	// multiple times.
	headerDone bool
	// the status error received from the server.
	status *status.Status
	// rstStream indicates whether a RST_STREAM frame needs to be sent
	// to the server to signify that this stream is closing.
	rstStream bool
	// rstError is the error that needs to be sent along with the RST_STREAM frame.
	rstError http2.ErrCode
	// bytesSent and bytesReceived indicates whether any bytes have been sent or
	// received on this stream.
	bytesSent     bool
	bytesReceived bool
***REMOVED***

// RecvCompress returns the compression algorithm applied to the inbound
// message. It is empty string if there is no compression applied.
func (s *Stream) RecvCompress() string ***REMOVED***
	return s.recvCompress
***REMOVED***

// SetSendCompress sets the compression algorithm to the stream.
func (s *Stream) SetSendCompress(str string) ***REMOVED***
	s.sendCompress = str
***REMOVED***

// Done returns a chanel which is closed when it receives the final status
// from the server.
func (s *Stream) Done() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return s.done
***REMOVED***

// GoAway returns a channel which is closed when the server sent GoAways signal
// before this stream was initiated.
func (s *Stream) GoAway() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return s.goAway
***REMOVED***

// Header acquires the key-value pairs of header metadata once it
// is available. It blocks until i) the metadata is ready or ii) there is no
// header metadata or iii) the stream is cancelled/expired.
func (s *Stream) Header() (metadata.MD, error) ***REMOVED***
	select ***REMOVED***
	case <-s.ctx.Done():
		return nil, ContextErr(s.ctx.Err())
	case <-s.goAway:
		return nil, ErrStreamDrain
	case <-s.headerChan:
		return s.header.Copy(), nil
	***REMOVED***
***REMOVED***

// Trailer returns the cached trailer metedata. Note that if it is not called
// after the entire stream is done, it could return an empty MD. Client
// side only.
func (s *Stream) Trailer() metadata.MD ***REMOVED***
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.trailer.Copy()
***REMOVED***

// ServerTransport returns the underlying ServerTransport for the stream.
// The client side stream always returns nil.
func (s *Stream) ServerTransport() ServerTransport ***REMOVED***
	return s.st
***REMOVED***

// Context returns the context of the stream.
func (s *Stream) Context() context.Context ***REMOVED***
	return s.ctx
***REMOVED***

// Method returns the method for the stream.
func (s *Stream) Method() string ***REMOVED***
	return s.method
***REMOVED***

// Status returns the status received from the server.
func (s *Stream) Status() *status.Status ***REMOVED***
	return s.status
***REMOVED***

// SetHeader sets the header metadata. This can be called multiple times.
// Server side only.
func (s *Stream) SetHeader(md metadata.MD) error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.headerOk || s.state == streamDone ***REMOVED***
		return ErrIllegalHeaderWrite
	***REMOVED***
	if md.Len() == 0 ***REMOVED***
		return nil
	***REMOVED***
	s.header = metadata.Join(s.header, md)
	return nil
***REMOVED***

// SetTrailer sets the trailer metadata which will be sent with the RPC status
// by the server. This can be called multiple times. Server side only.
func (s *Stream) SetTrailer(md metadata.MD) error ***REMOVED***
	if md.Len() == 0 ***REMOVED***
		return nil
	***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	s.trailer = metadata.Join(s.trailer, md)
	return nil
***REMOVED***

func (s *Stream) write(m recvMsg) ***REMOVED***
	s.buf.put(&m)
***REMOVED***

// Read reads all the data available for this Stream from the transport and
// passes them into the decoder, which converts them into a gRPC message stream.
// The error is io.EOF when the stream is done or another non-nil error if
// the stream broke.
func (s *Stream) Read(p []byte) (n int, err error) ***REMOVED***
	n, err = s.dec.Read(p)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	s.windowHandler(n)
	return
***REMOVED***

// finish sets the stream's state and status, and closes the done channel.
// s.mu must be held by the caller.  st must always be non-nil.
func (s *Stream) finish(st *status.Status) ***REMOVED***
	s.status = st
	s.state = streamDone
	close(s.done)
***REMOVED***

// BytesSent indicates whether any bytes have been sent on this stream.
func (s *Stream) BytesSent() bool ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.bytesSent
***REMOVED***

// BytesReceived indicates whether any bytes have been received on this stream.
func (s *Stream) BytesReceived() bool ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.bytesReceived
***REMOVED***

// GoString is implemented by Stream so context.String() won't
// race when printing %#v.
func (s *Stream) GoString() string ***REMOVED***
	return fmt.Sprintf("<stream: %p, %v>", s, s.method)
***REMOVED***

// The key to save transport.Stream in the context.
type streamKey struct***REMOVED******REMOVED***

// newContextWithStream creates a new context from ctx and attaches stream
// to it.
func newContextWithStream(ctx context.Context, stream *Stream) context.Context ***REMOVED***
	return context.WithValue(ctx, streamKey***REMOVED******REMOVED***, stream)
***REMOVED***

// StreamFromContext returns the stream saved in ctx.
func StreamFromContext(ctx context.Context) (s *Stream, ok bool) ***REMOVED***
	s, ok = ctx.Value(streamKey***REMOVED******REMOVED***).(*Stream)
	return
***REMOVED***

// state of transport
type transportState int

const (
	reachable transportState = iota
	unreachable
	closing
	draining
)

// ServerConfig consists of all the configurations to establish a server transport.
type ServerConfig struct ***REMOVED***
	MaxStreams      uint32
	AuthInfo        credentials.AuthInfo
	InTapHandle     tap.ServerInHandle
	StatsHandler    stats.Handler
	KeepaliveParams keepalive.ServerParameters
	KeepalivePolicy keepalive.EnforcementPolicy
***REMOVED***

// NewServerTransport creates a ServerTransport with conn or non-nil error
// if it fails.
func NewServerTransport(protocol string, conn net.Conn, config *ServerConfig) (ServerTransport, error) ***REMOVED***
	return newHTTP2Server(conn, config)
***REMOVED***

// ConnectOptions covers all relevant options for communicating with the server.
type ConnectOptions struct ***REMOVED***
	// UserAgent is the application user agent.
	UserAgent string
	// Authority is the :authority pseudo-header to use. This field has no effect if
	// TransportCredentials is set.
	Authority string
	// Dialer specifies how to dial a network address.
	Dialer func(context.Context, string) (net.Conn, error)
	// FailOnNonTempDialError specifies if gRPC fails on non-temporary dial errors.
	FailOnNonTempDialError bool
	// PerRPCCredentials stores the PerRPCCredentials required to issue RPCs.
	PerRPCCredentials []credentials.PerRPCCredentials
	// TransportCredentials stores the Authenticator required to setup a client connection.
	TransportCredentials credentials.TransportCredentials
	// KeepaliveParams stores the keepalive parameters.
	KeepaliveParams keepalive.ClientParameters
	// StatsHandler stores the handler for stats.
	StatsHandler stats.Handler
***REMOVED***

// TargetInfo contains the information of the target such as network address and metadata.
type TargetInfo struct ***REMOVED***
	Addr     string
	Metadata interface***REMOVED******REMOVED***
***REMOVED***

// NewClientTransport establishes the transport with the required ConnectOptions
// and returns it to the caller.
func NewClientTransport(ctx context.Context, target TargetInfo, opts ConnectOptions) (ClientTransport, error) ***REMOVED***
	return newHTTP2Client(ctx, target, opts)
***REMOVED***

// Options provides additional hints and information for message
// transmission.
type Options struct ***REMOVED***
	// Last indicates whether this write is the last piece for
	// this stream.
	Last bool

	// Delay is a hint to the transport implementation for whether
	// the data could be buffered for a batching write. The
	// Transport implementation may ignore the hint.
	Delay bool
***REMOVED***

// CallHdr carries the information of a particular RPC.
type CallHdr struct ***REMOVED***
	// Host specifies the peer's host.
	Host string

	// Method specifies the operation to perform.
	Method string

	// RecvCompress specifies the compression algorithm applied on
	// inbound messages.
	RecvCompress string

	// SendCompress specifies the compression algorithm applied on
	// outbound message.
	SendCompress string

	// Flush indicates whether a new stream command should be sent
	// to the peer without waiting for the first data. This is
	// only a hint. The transport may modify the flush decision
	// for performance purposes.
	Flush bool
***REMOVED***

// ClientTransport is the common interface for all gRPC client-side transport
// implementations.
type ClientTransport interface ***REMOVED***
	// Close tears down this transport. Once it returns, the transport
	// should not be accessed any more. The caller must make sure this
	// is called only once.
	Close() error

	// GracefulClose starts to tear down the transport. It stops accepting
	// new RPCs and wait the completion of the pending RPCs.
	GracefulClose() error

	// Write sends the data for the given stream. A nil stream indicates
	// the write is to be performed on the transport as a whole.
	Write(s *Stream, data []byte, opts *Options) error

	// NewStream creates a Stream for an RPC.
	NewStream(ctx context.Context, callHdr *CallHdr) (*Stream, error)

	// CloseStream clears the footprint of a stream when the stream is
	// not needed any more. The err indicates the error incurred when
	// CloseStream is called. Must be called when a stream is finished
	// unless the associated transport is closing.
	CloseStream(stream *Stream, err error)

	// Error returns a channel that is closed when some I/O error
	// happens. Typically the caller should have a goroutine to monitor
	// this in order to take action (e.g., close the current transport
	// and create a new one) in error case. It should not return nil
	// once the transport is initiated.
	Error() <-chan struct***REMOVED******REMOVED***

	// GoAway returns a channel that is closed when ClientTranspor
	// receives the draining signal from the server (e.g., GOAWAY frame in
	// HTTP/2).
	GoAway() <-chan struct***REMOVED******REMOVED***

	// GetGoAwayReason returns the reason why GoAway frame was received.
	GetGoAwayReason() GoAwayReason
***REMOVED***

// ServerTransport is the common interface for all gRPC server-side transport
// implementations.
//
// Methods may be called concurrently from multiple goroutines, but
// Write methods for a given Stream will be called serially.
type ServerTransport interface ***REMOVED***
	// HandleStreams receives incoming streams using the given handler.
	HandleStreams(func(*Stream), func(context.Context, string) context.Context)

	// WriteHeader sends the header metadata for the given stream.
	// WriteHeader may not be called on all streams.
	WriteHeader(s *Stream, md metadata.MD) error

	// Write sends the data for the given stream.
	// Write may not be called on all streams.
	Write(s *Stream, data []byte, opts *Options) error

	// WriteStatus sends the status of a stream to the client.  WriteStatus is
	// the final call made on a stream and always occurs.
	WriteStatus(s *Stream, st *status.Status) error

	// Close tears down the transport. Once it is called, the transport
	// should not be accessed any more. All the pending streams and their
	// handlers will be terminated asynchronously.
	Close() error

	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr

	// Drain notifies the client this ServerTransport stops accepting new RPCs.
	Drain()
***REMOVED***

// streamErrorf creates an StreamError with the specified error code and description.
func streamErrorf(c codes.Code, format string, a ...interface***REMOVED******REMOVED***) StreamError ***REMOVED***
	return StreamError***REMOVED***
		Code: c,
		Desc: fmt.Sprintf(format, a...),
	***REMOVED***
***REMOVED***

// connectionErrorf creates an ConnectionError with the specified error description.
func connectionErrorf(temp bool, e error, format string, a ...interface***REMOVED******REMOVED***) ConnectionError ***REMOVED***
	return ConnectionError***REMOVED***
		Desc: fmt.Sprintf(format, a...),
		temp: temp,
		err:  e,
	***REMOVED***
***REMOVED***

// ConnectionError is an error that results in the termination of the
// entire connection and the retry of all the active streams.
type ConnectionError struct ***REMOVED***
	Desc string
	temp bool
	err  error
***REMOVED***

func (e ConnectionError) Error() string ***REMOVED***
	return fmt.Sprintf("connection error: desc = %q", e.Desc)
***REMOVED***

// Temporary indicates if this connection error is temporary or fatal.
func (e ConnectionError) Temporary() bool ***REMOVED***
	return e.temp
***REMOVED***

// Origin returns the original error of this connection error.
func (e ConnectionError) Origin() error ***REMOVED***
	// Never return nil error here.
	// If the original error is nil, return itself.
	if e.err == nil ***REMOVED***
		return e
	***REMOVED***
	return e.err
***REMOVED***

var (
	// ErrConnClosing indicates that the transport is closing.
	ErrConnClosing = connectionErrorf(true, nil, "transport is closing")
	// ErrStreamDrain indicates that the stream is rejected by the server because
	// the server stops accepting new RPCs.
	ErrStreamDrain = streamErrorf(codes.Unavailable, "the server stops accepting new RPCs")
)

// TODO: See if we can replace StreamError with status package errors.

// StreamError is an error that only affects one stream within a connection.
type StreamError struct ***REMOVED***
	Code codes.Code
	Desc string
***REMOVED***

func (e StreamError) Error() string ***REMOVED***
	return fmt.Sprintf("stream error: code = %s desc = %q", e.Code, e.Desc)
***REMOVED***

// ContextErr converts the error from context package into a StreamError.
func ContextErr(err error) StreamError ***REMOVED***
	switch err ***REMOVED***
	case context.DeadlineExceeded:
		return streamErrorf(codes.DeadlineExceeded, "%v", err)
	case context.Canceled:
		return streamErrorf(codes.Canceled, "%v", err)
	***REMOVED***
	panic(fmt.Sprintf("Unexpected error from context packet: %v", err))
***REMOVED***

// wait blocks until it can receive from ctx.Done, closing, or proceed.
// If it receives from ctx.Done, it returns 0, the StreamError for ctx.Err.
// If it receives from done, it returns 0, io.EOF if ctx is not done; otherwise
// it return the StreamError for ctx.Err.
// If it receives from goAway, it returns 0, ErrStreamDrain.
// If it receives from closing, it returns 0, ErrConnClosing.
// If it receives from proceed, it returns the received integer, nil.
func wait(ctx context.Context, done, goAway, closing <-chan struct***REMOVED******REMOVED***, proceed <-chan int) (int, error) ***REMOVED***
	select ***REMOVED***
	case <-ctx.Done():
		return 0, ContextErr(ctx.Err())
	case <-done:
		// User cancellation has precedence.
		select ***REMOVED***
		case <-ctx.Done():
			return 0, ContextErr(ctx.Err())
		default:
		***REMOVED***
		return 0, io.EOF
	case <-goAway:
		return 0, ErrStreamDrain
	case <-closing:
		return 0, ErrConnClosing
	case i := <-proceed:
		return i, nil
	***REMOVED***
***REMOVED***

// GoAwayReason contains the reason for the GoAway frame received.
type GoAwayReason uint8

const (
	// Invalid indicates that no GoAway frame is received.
	Invalid GoAwayReason = 0
	// NoReason is the default value when GoAway frame is received.
	NoReason GoAwayReason = 1
	// TooManyPings indicates that a GoAway frame with ErrCodeEnhanceYourCalm
	// was recieved and that the debug data said "too_many_pings".
	TooManyPings GoAwayReason = 2
)
