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
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"golang.org/x/net/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
	"google.golang.org/grpc/transport"
)

type methodHandler func(srv interface***REMOVED******REMOVED***, ctx context.Context, dec func(interface***REMOVED******REMOVED***) error, interceptor UnaryServerInterceptor) (interface***REMOVED******REMOVED***, error)

// MethodDesc represents an RPC service's method specification.
type MethodDesc struct ***REMOVED***
	MethodName string
	Handler    methodHandler
***REMOVED***

// ServiceDesc represents an RPC service's specification.
type ServiceDesc struct ***REMOVED***
	ServiceName string
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	HandlerType interface***REMOVED******REMOVED***
	Methods     []MethodDesc
	Streams     []StreamDesc
	Metadata    interface***REMOVED******REMOVED***
***REMOVED***

// service consists of the information of the server serving this service and
// the methods in this service.
type service struct ***REMOVED***
	server interface***REMOVED******REMOVED*** // the server for service methods
	md     map[string]*MethodDesc
	sd     map[string]*StreamDesc
	mdata  interface***REMOVED******REMOVED***
***REMOVED***

// Server is a gRPC server to serve RPC requests.
type Server struct ***REMOVED***
	opts options

	mu     sync.Mutex // guards following
	lis    map[net.Listener]bool
	conns  map[io.Closer]bool
	drain  bool
	ctx    context.Context
	cancel context.CancelFunc
	// A CondVar to let GracefulStop() blocks until all the pending RPCs are finished
	// and all the transport goes away.
	cv     *sync.Cond
	m      map[string]*service // service name -> service info
	events trace.EventLog
***REMOVED***

type options struct ***REMOVED***
	creds                credentials.TransportCredentials
	codec                Codec
	cp                   Compressor
	dc                   Decompressor
	maxMsgSize           int
	unaryInt             UnaryServerInterceptor
	streamInt            StreamServerInterceptor
	inTapHandle          tap.ServerInHandle
	statsHandler         stats.Handler
	maxConcurrentStreams uint32
	useHandlerImpl       bool // use http.Handler-based server
	unknownStreamDesc    *StreamDesc
	keepaliveParams      keepalive.ServerParameters
	keepalivePolicy      keepalive.EnforcementPolicy
***REMOVED***

var defaultMaxMsgSize = 1024 * 1024 * 4 // use 4MB as the default message size limit

// A ServerOption sets options.
type ServerOption func(*options)

// KeepaliveParams returns a ServerOption that sets keepalive and max-age parameters for the server.
func KeepaliveParams(kp keepalive.ServerParameters) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		o.keepaliveParams = kp
	***REMOVED***
***REMOVED***

// KeepaliveEnforcementPolicy returns a ServerOption that sets keepalive enforcement policy for the server.
func KeepaliveEnforcementPolicy(kep keepalive.EnforcementPolicy) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		o.keepalivePolicy = kep
	***REMOVED***
***REMOVED***

// CustomCodec returns a ServerOption that sets a codec for message marshaling and unmarshaling.
func CustomCodec(codec Codec) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		o.codec = codec
	***REMOVED***
***REMOVED***

// RPCCompressor returns a ServerOption that sets a compressor for outbound messages.
func RPCCompressor(cp Compressor) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		o.cp = cp
	***REMOVED***
***REMOVED***

// RPCDecompressor returns a ServerOption that sets a decompressor for inbound messages.
func RPCDecompressor(dc Decompressor) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		o.dc = dc
	***REMOVED***
***REMOVED***

// MaxMsgSize returns a ServerOption to set the max message size in bytes for inbound mesages.
// If this is not set, gRPC uses the default 4MB.
func MaxMsgSize(m int) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		o.maxMsgSize = m
	***REMOVED***
***REMOVED***

// MaxConcurrentStreams returns a ServerOption that will apply a limit on the number
// of concurrent streams to each ServerTransport.
func MaxConcurrentStreams(n uint32) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		o.maxConcurrentStreams = n
	***REMOVED***
***REMOVED***

// Creds returns a ServerOption that sets credentials for server connections.
func Creds(c credentials.TransportCredentials) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		o.creds = c
	***REMOVED***
***REMOVED***

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the
// server. Only one unary interceptor can be installed. The construction of multiple
// interceptors (e.g., chaining) can be implemented at the caller.
func UnaryInterceptor(i UnaryServerInterceptor) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		if o.unaryInt != nil ***REMOVED***
			panic("The unary server interceptor has been set.")
		***REMOVED***
		o.unaryInt = i
	***REMOVED***
***REMOVED***

// StreamInterceptor returns a ServerOption that sets the StreamServerInterceptor for the
// server. Only one stream interceptor can be installed.
func StreamInterceptor(i StreamServerInterceptor) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		if o.streamInt != nil ***REMOVED***
			panic("The stream server interceptor has been set.")
		***REMOVED***
		o.streamInt = i
	***REMOVED***
***REMOVED***

// InTapHandle returns a ServerOption that sets the tap handle for all the server
// transport to be created. Only one can be installed.
func InTapHandle(h tap.ServerInHandle) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		if o.inTapHandle != nil ***REMOVED***
			panic("The tap handle has been set.")
		***REMOVED***
		o.inTapHandle = h
	***REMOVED***
***REMOVED***

// StatsHandler returns a ServerOption that sets the stats handler for the server.
func StatsHandler(h stats.Handler) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		o.statsHandler = h
	***REMOVED***
***REMOVED***

// UnknownServiceHandler returns a ServerOption that allows for adding a custom
// unknown service handler. The provided method is a bidi-streaming RPC service
// handler that will be invoked instead of returning the the "unimplemented" gRPC
// error whenever a request is received for an unregistered service or method.
// The handling function has full access to the Context of the request and the
// stream, and the invocation passes through interceptors.
func UnknownServiceHandler(streamHandler StreamHandler) ServerOption ***REMOVED***
	return func(o *options) ***REMOVED***
		o.unknownStreamDesc = &StreamDesc***REMOVED***
			StreamName: "unknown_service_handler",
			Handler:    streamHandler,
			// We need to assume that the users of the streamHandler will want to use both.
			ClientStreams: true,
			ServerStreams: true,
		***REMOVED***
	***REMOVED***
***REMOVED***

// NewServer creates a gRPC server which has no service registered and has not
// started to accept requests yet.
func NewServer(opt ...ServerOption) *Server ***REMOVED***
	var opts options
	opts.maxMsgSize = defaultMaxMsgSize
	for _, o := range opt ***REMOVED***
		o(&opts)
	***REMOVED***
	if opts.codec == nil ***REMOVED***
		// Set the default codec.
		opts.codec = protoCodec***REMOVED******REMOVED***
	***REMOVED***
	s := &Server***REMOVED***
		lis:   make(map[net.Listener]bool),
		opts:  opts,
		conns: make(map[io.Closer]bool),
		m:     make(map[string]*service),
	***REMOVED***
	s.cv = sync.NewCond(&s.mu)
	s.ctx, s.cancel = context.WithCancel(context.Background())
	if EnableTracing ***REMOVED***
		_, file, line, _ := runtime.Caller(1)
		s.events = trace.NewEventLog("grpc.Server", fmt.Sprintf("%s:%d", file, line))
	***REMOVED***
	return s
***REMOVED***

// printf records an event in s's event log, unless s has been stopped.
// REQUIRES s.mu is held.
func (s *Server) printf(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	if s.events != nil ***REMOVED***
		s.events.Printf(format, a...)
	***REMOVED***
***REMOVED***

// errorf records an error in s's event log, unless s has been stopped.
// REQUIRES s.mu is held.
func (s *Server) errorf(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	if s.events != nil ***REMOVED***
		s.events.Errorf(format, a...)
	***REMOVED***
***REMOVED***

// RegisterService register a service and its implementation to the gRPC
// server. Called from the IDL generated code. This must be called before
// invoking Serve.
func (s *Server) RegisterService(sd *ServiceDesc, ss interface***REMOVED******REMOVED***) ***REMOVED***
	ht := reflect.TypeOf(sd.HandlerType).Elem()
	st := reflect.TypeOf(ss)
	if !st.Implements(ht) ***REMOVED***
		grpclog.Fatalf("grpc: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
	***REMOVED***
	s.register(sd, ss)
***REMOVED***

func (s *Server) register(sd *ServiceDesc, ss interface***REMOVED******REMOVED***) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	s.printf("RegisterService(%q)", sd.ServiceName)
	if _, ok := s.m[sd.ServiceName]; ok ***REMOVED***
		grpclog.Fatalf("grpc: Server.RegisterService found duplicate service registration for %q", sd.ServiceName)
	***REMOVED***
	srv := &service***REMOVED***
		server: ss,
		md:     make(map[string]*MethodDesc),
		sd:     make(map[string]*StreamDesc),
		mdata:  sd.Metadata,
	***REMOVED***
	for i := range sd.Methods ***REMOVED***
		d := &sd.Methods[i]
		srv.md[d.MethodName] = d
	***REMOVED***
	for i := range sd.Streams ***REMOVED***
		d := &sd.Streams[i]
		srv.sd[d.StreamName] = d
	***REMOVED***
	s.m[sd.ServiceName] = srv
***REMOVED***

// MethodInfo contains the information of an RPC including its method name and type.
type MethodInfo struct ***REMOVED***
	// Name is the method name only, without the service name or package name.
	Name string
	// IsClientStream indicates whether the RPC is a client streaming RPC.
	IsClientStream bool
	// IsServerStream indicates whether the RPC is a server streaming RPC.
	IsServerStream bool
***REMOVED***

// ServiceInfo contains unary RPC method info, streaming RPC methid info and metadata for a service.
type ServiceInfo struct ***REMOVED***
	Methods []MethodInfo
	// Metadata is the metadata specified in ServiceDesc when registering service.
	Metadata interface***REMOVED******REMOVED***
***REMOVED***

// GetServiceInfo returns a map from service names to ServiceInfo.
// Service names include the package names, in the form of <package>.<service>.
func (s *Server) GetServiceInfo() map[string]ServiceInfo ***REMOVED***
	ret := make(map[string]ServiceInfo)
	for n, srv := range s.m ***REMOVED***
		methods := make([]MethodInfo, 0, len(srv.md)+len(srv.sd))
		for m := range srv.md ***REMOVED***
			methods = append(methods, MethodInfo***REMOVED***
				Name:           m,
				IsClientStream: false,
				IsServerStream: false,
			***REMOVED***)
		***REMOVED***
		for m, d := range srv.sd ***REMOVED***
			methods = append(methods, MethodInfo***REMOVED***
				Name:           m,
				IsClientStream: d.ClientStreams,
				IsServerStream: d.ServerStreams,
			***REMOVED***)
		***REMOVED***

		ret[n] = ServiceInfo***REMOVED***
			Methods:  methods,
			Metadata: srv.mdata,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

var (
	// ErrServerStopped indicates that the operation is now illegal because of
	// the server being stopped.
	ErrServerStopped = errors.New("grpc: the server has been stopped")
)

func (s *Server) useTransportAuthenticator(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) ***REMOVED***
	if s.opts.creds == nil ***REMOVED***
		return rawConn, nil, nil
	***REMOVED***
	return s.opts.creds.ServerHandshake(rawConn)
***REMOVED***

// Serve accepts incoming connections on the listener lis, creating a new
// ServerTransport and service goroutine for each. The service goroutines
// read gRPC requests and then call the registered handlers to reply to them.
// Serve returns when lis.Accept fails with fatal errors.  lis will be closed when
// this method returns.
// Serve always returns non-nil error.
func (s *Server) Serve(lis net.Listener) error ***REMOVED***
	s.mu.Lock()
	s.printf("serving")
	if s.lis == nil ***REMOVED***
		s.mu.Unlock()
		lis.Close()
		return ErrServerStopped
	***REMOVED***
	s.lis[lis] = true
	s.mu.Unlock()
	defer func() ***REMOVED***
		s.mu.Lock()
		if s.lis != nil && s.lis[lis] ***REMOVED***
			lis.Close()
			delete(s.lis, lis)
		***REMOVED***
		s.mu.Unlock()
	***REMOVED***()

	var tempDelay time.Duration // how long to sleep on accept failure

	for ***REMOVED***
		rawConn, err := lis.Accept()
		if err != nil ***REMOVED***
			if ne, ok := err.(interface ***REMOVED***
				Temporary() bool
			***REMOVED***); ok && ne.Temporary() ***REMOVED***
				if tempDelay == 0 ***REMOVED***
					tempDelay = 5 * time.Millisecond
				***REMOVED*** else ***REMOVED***
					tempDelay *= 2
				***REMOVED***
				if max := 1 * time.Second; tempDelay > max ***REMOVED***
					tempDelay = max
				***REMOVED***
				s.mu.Lock()
				s.printf("Accept error: %v; retrying in %v", err, tempDelay)
				s.mu.Unlock()
				select ***REMOVED***
				case <-time.After(tempDelay):
				case <-s.ctx.Done():
				***REMOVED***
				continue
			***REMOVED***
			s.mu.Lock()
			s.printf("done serving; Accept = %v", err)
			s.mu.Unlock()
			return err
		***REMOVED***
		tempDelay = 0
		// Start a new goroutine to deal with rawConn
		// so we don't stall this Accept loop goroutine.
		go s.handleRawConn(rawConn)
	***REMOVED***
***REMOVED***

// handleRawConn is run in its own goroutine and handles a just-accepted
// connection that has not had any I/O performed on it yet.
func (s *Server) handleRawConn(rawConn net.Conn) ***REMOVED***
	conn, authInfo, err := s.useTransportAuthenticator(rawConn)
	if err != nil ***REMOVED***
		s.mu.Lock()
		s.errorf("ServerHandshake(%q) failed: %v", rawConn.RemoteAddr(), err)
		s.mu.Unlock()
		grpclog.Printf("grpc: Server.Serve failed to complete security handshake from %q: %v", rawConn.RemoteAddr(), err)
		// If serverHandShake returns ErrConnDispatched, keep rawConn open.
		if err != credentials.ErrConnDispatched ***REMOVED***
			rawConn.Close()
		***REMOVED***
		return
	***REMOVED***

	s.mu.Lock()
	if s.conns == nil ***REMOVED***
		s.mu.Unlock()
		conn.Close()
		return
	***REMOVED***
	s.mu.Unlock()

	if s.opts.useHandlerImpl ***REMOVED***
		s.serveUsingHandler(conn)
	***REMOVED*** else ***REMOVED***
		s.serveHTTP2Transport(conn, authInfo)
	***REMOVED***
***REMOVED***

// serveHTTP2Transport sets up a http/2 transport (using the
// gRPC http2 server transport in transport/http2_server.go) and
// serves streams on it.
// This is run in its own goroutine (it does network I/O in
// transport.NewServerTransport).
func (s *Server) serveHTTP2Transport(c net.Conn, authInfo credentials.AuthInfo) ***REMOVED***
	config := &transport.ServerConfig***REMOVED***
		MaxStreams:      s.opts.maxConcurrentStreams,
		AuthInfo:        authInfo,
		InTapHandle:     s.opts.inTapHandle,
		StatsHandler:    s.opts.statsHandler,
		KeepaliveParams: s.opts.keepaliveParams,
		KeepalivePolicy: s.opts.keepalivePolicy,
	***REMOVED***
	st, err := transport.NewServerTransport("http2", c, config)
	if err != nil ***REMOVED***
		s.mu.Lock()
		s.errorf("NewServerTransport(%q) failed: %v", c.RemoteAddr(), err)
		s.mu.Unlock()
		c.Close()
		grpclog.Println("grpc: Server.Serve failed to create ServerTransport: ", err)
		return
	***REMOVED***
	if !s.addConn(st) ***REMOVED***
		st.Close()
		return
	***REMOVED***
	s.serveStreams(st)
***REMOVED***

func (s *Server) serveStreams(st transport.ServerTransport) ***REMOVED***
	defer s.removeConn(st)
	defer st.Close()
	var wg sync.WaitGroup
	st.HandleStreams(func(stream *transport.Stream) ***REMOVED***
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			s.handleStream(st, stream, s.traceInfo(st, stream))
		***REMOVED***()
	***REMOVED***, func(ctx context.Context, method string) context.Context ***REMOVED***
		if !EnableTracing ***REMOVED***
			return ctx
		***REMOVED***
		tr := trace.New("grpc.Recv."+methodFamily(method), method)
		return trace.NewContext(ctx, tr)
	***REMOVED***)
	wg.Wait()
***REMOVED***

var _ http.Handler = (*Server)(nil)

// serveUsingHandler is called from handleRawConn when s is configured
// to handle requests via the http.Handler interface. It sets up a
// net/http.Server to handle the just-accepted conn. The http.Server
// is configured to route all incoming requests (all HTTP/2 streams)
// to ServeHTTP, which creates a new ServerTransport for each stream.
// serveUsingHandler blocks until conn closes.
//
// This codepath is only used when Server.TestingUseHandlerImpl has
// been configured. This lets the end2end tests exercise the ServeHTTP
// method as one of the environment types.
//
// conn is the *tls.Conn that's already been authenticated.
func (s *Server) serveUsingHandler(conn net.Conn) ***REMOVED***
	if !s.addConn(conn) ***REMOVED***
		conn.Close()
		return
	***REMOVED***
	defer s.removeConn(conn)
	h2s := &http2.Server***REMOVED***
		MaxConcurrentStreams: s.opts.maxConcurrentStreams,
	***REMOVED***
	h2s.ServeConn(conn, &http2.ServeConnOpts***REMOVED***
		Handler: s,
	***REMOVED***)
***REMOVED***

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	st, err := transport.NewServerHandlerTransport(w, r)
	if err != nil ***REMOVED***
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	if !s.addConn(st) ***REMOVED***
		st.Close()
		return
	***REMOVED***
	defer s.removeConn(st)
	s.serveStreams(st)
***REMOVED***

// traceInfo returns a traceInfo and associates it with stream, if tracing is enabled.
// If tracing is not enabled, it returns nil.
func (s *Server) traceInfo(st transport.ServerTransport, stream *transport.Stream) (trInfo *traceInfo) ***REMOVED***
	tr, ok := trace.FromContext(stream.Context())
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	trInfo = &traceInfo***REMOVED***
		tr: tr,
	***REMOVED***
	trInfo.firstLine.client = false
	trInfo.firstLine.remoteAddr = st.RemoteAddr()

	if dl, ok := stream.Context().Deadline(); ok ***REMOVED***
		trInfo.firstLine.deadline = dl.Sub(time.Now())
	***REMOVED***
	return trInfo
***REMOVED***

func (s *Server) addConn(c io.Closer) bool ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns == nil || s.drain ***REMOVED***
		return false
	***REMOVED***
	s.conns[c] = true
	return true
***REMOVED***

func (s *Server) removeConn(c io.Closer) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns != nil ***REMOVED***
		delete(s.conns, c)
		s.cv.Broadcast()
	***REMOVED***
***REMOVED***

func (s *Server) sendResponse(t transport.ServerTransport, stream *transport.Stream, msg interface***REMOVED******REMOVED***, cp Compressor, opts *transport.Options) error ***REMOVED***
	var (
		cbuf       *bytes.Buffer
		outPayload *stats.OutPayload
	)
	if cp != nil ***REMOVED***
		cbuf = new(bytes.Buffer)
	***REMOVED***
	if s.opts.statsHandler != nil ***REMOVED***
		outPayload = &stats.OutPayload***REMOVED******REMOVED***
	***REMOVED***
	p, err := encode(s.opts.codec, msg, cp, cbuf, outPayload)
	if err != nil ***REMOVED***
		// This typically indicates a fatal issue (e.g., memory
		// corruption or hardware faults) the application program
		// cannot handle.
		//
		// TODO(zhaoq): There exist other options also such as only closing the
		// faulty stream locally and remotely (Other streams can keep going). Find
		// the optimal option.
		grpclog.Fatalf("grpc: Server failed to encode response %v", err)
	***REMOVED***
	err = t.Write(stream, p, opts)
	if err == nil && outPayload != nil ***REMOVED***
		outPayload.SentTime = time.Now()
		s.opts.statsHandler.HandleRPC(stream.Context(), outPayload)
	***REMOVED***
	return err
***REMOVED***

func (s *Server) processUnaryRPC(t transport.ServerTransport, stream *transport.Stream, srv *service, md *MethodDesc, trInfo *traceInfo) (err error) ***REMOVED***
	sh := s.opts.statsHandler
	if sh != nil ***REMOVED***
		begin := &stats.Begin***REMOVED***
			BeginTime: time.Now(),
		***REMOVED***
		sh.HandleRPC(stream.Context(), begin)
	***REMOVED***
	defer func() ***REMOVED***
		if sh != nil ***REMOVED***
			end := &stats.End***REMOVED***
				EndTime: time.Now(),
			***REMOVED***
			if err != nil && err != io.EOF ***REMOVED***
				end.Error = toRPCErr(err)
			***REMOVED***
			sh.HandleRPC(stream.Context(), end)
		***REMOVED***
	***REMOVED***()
	if trInfo != nil ***REMOVED***
		defer trInfo.tr.Finish()
		trInfo.firstLine.client = false
		trInfo.tr.LazyLog(&trInfo.firstLine, false)
		defer func() ***REMOVED***
			if err != nil && err != io.EOF ***REMOVED***
				trInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***err***REMOVED******REMOVED***, true)
				trInfo.tr.SetError()
			***REMOVED***
		***REMOVED***()
	***REMOVED***
	if s.opts.cp != nil ***REMOVED***
		// NOTE: this needs to be ahead of all handling, https://github.com/grpc/grpc-go/issues/686.
		stream.SetSendCompress(s.opts.cp.Type())
	***REMOVED***
	p := &parser***REMOVED***r: stream***REMOVED***
	for ***REMOVED*** // TODO: delete
		pf, req, err := p.recvMsg(s.opts.maxMsgSize)
		if err == io.EOF ***REMOVED***
			// The entire stream is done (for unary RPC only).
			return err
		***REMOVED***
		if err == io.ErrUnexpectedEOF ***REMOVED***
			err = Errorf(codes.Internal, io.ErrUnexpectedEOF.Error())
		***REMOVED***
		if err != nil ***REMOVED***
			if st, ok := status.FromError(err); ok ***REMOVED***
				if e := t.WriteStatus(stream, st); e != nil ***REMOVED***
					grpclog.Printf("grpc: Server.processUnaryRPC failed to write status %v", e)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				switch st := err.(type) ***REMOVED***
				case transport.ConnectionError:
					// Nothing to do here.
				case transport.StreamError:
					if e := t.WriteStatus(stream, status.New(st.Code, st.Desc)); e != nil ***REMOVED***
						grpclog.Printf("grpc: Server.processUnaryRPC failed to write status %v", e)
					***REMOVED***
				default:
					panic(fmt.Sprintf("grpc: Unexpected error (%T) from recvMsg: %v", st, st))
				***REMOVED***
			***REMOVED***
			return err
		***REMOVED***

		if err := checkRecvPayload(pf, stream.RecvCompress(), s.opts.dc); err != nil ***REMOVED***
			if st, ok := status.FromError(err); ok ***REMOVED***
				if e := t.WriteStatus(stream, st); e != nil ***REMOVED***
					grpclog.Printf("grpc: Server.processUnaryRPC failed to write status %v", e)
				***REMOVED***
				return err
			***REMOVED***
			if e := t.WriteStatus(stream, status.New(codes.Internal, err.Error())); e != nil ***REMOVED***
				grpclog.Printf("grpc: Server.processUnaryRPC failed to write status %v", e)
			***REMOVED***

			// TODO checkRecvPayload always return RPC error. Add a return here if necessary.
		***REMOVED***
		var inPayload *stats.InPayload
		if sh != nil ***REMOVED***
			inPayload = &stats.InPayload***REMOVED***
				RecvTime: time.Now(),
			***REMOVED***
		***REMOVED***
		df := func(v interface***REMOVED******REMOVED***) error ***REMOVED***
			if inPayload != nil ***REMOVED***
				inPayload.WireLength = len(req)
			***REMOVED***
			if pf == compressionMade ***REMOVED***
				var err error
				req, err = s.opts.dc.Do(bytes.NewReader(req))
				if err != nil ***REMOVED***
					return Errorf(codes.Internal, err.Error())
				***REMOVED***
			***REMOVED***
			if len(req) > s.opts.maxMsgSize ***REMOVED***
				// TODO: Revisit the error code. Currently keep it consistent with
				// java implementation.
				return status.Errorf(codes.Internal, "grpc: server received a message of %d bytes exceeding %d limit", len(req), s.opts.maxMsgSize)
			***REMOVED***
			if err := s.opts.codec.Unmarshal(req, v); err != nil ***REMOVED***
				return status.Errorf(codes.Internal, "grpc: error unmarshalling request: %v", err)
			***REMOVED***
			if inPayload != nil ***REMOVED***
				inPayload.Payload = v
				inPayload.Data = req
				inPayload.Length = len(req)
				sh.HandleRPC(stream.Context(), inPayload)
			***REMOVED***
			if trInfo != nil ***REMOVED***
				trInfo.tr.LazyLog(&payload***REMOVED***sent: false, msg: v***REMOVED***, true)
			***REMOVED***
			return nil
		***REMOVED***
		reply, appErr := md.Handler(srv.server, stream.Context(), df, s.opts.unaryInt)
		if appErr != nil ***REMOVED***
			appStatus, ok := status.FromError(appErr)
			if !ok ***REMOVED***
				// Convert appErr if it is not a grpc status error.
				appErr = status.Error(convertCode(appErr), appErr.Error())
				appStatus, _ = status.FromError(appErr)
			***REMOVED***
			if trInfo != nil ***REMOVED***
				trInfo.tr.LazyLog(stringer(appStatus.Message()), true)
				trInfo.tr.SetError()
			***REMOVED***
			if e := t.WriteStatus(stream, appStatus); e != nil ***REMOVED***
				grpclog.Printf("grpc: Server.processUnaryRPC failed to write status: %v", e)
			***REMOVED***
			return appErr
		***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.LazyLog(stringer("OK"), false)
		***REMOVED***
		opts := &transport.Options***REMOVED***
			Last:  true,
			Delay: false,
		***REMOVED***
		if err := s.sendResponse(t, stream, reply, s.opts.cp, opts); err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				// The entire stream is done (for unary RPC only).
				return err
			***REMOVED***
			if s, ok := status.FromError(err); ok ***REMOVED***
				if e := t.WriteStatus(stream, s); e != nil ***REMOVED***
					grpclog.Printf("grpc: Server.processUnaryRPC failed to write status: %v", e)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				switch st := err.(type) ***REMOVED***
				case transport.ConnectionError:
					// Nothing to do here.
				case transport.StreamError:
					if e := t.WriteStatus(stream, status.New(st.Code, st.Desc)); e != nil ***REMOVED***
						grpclog.Printf("grpc: Server.processUnaryRPC failed to write status %v", e)
					***REMOVED***
				default:
					panic(fmt.Sprintf("grpc: Unexpected error (%T) from sendResponse: %v", st, st))
				***REMOVED***
			***REMOVED***
			return err
		***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.LazyLog(&payload***REMOVED***sent: true, msg: reply***REMOVED***, true)
		***REMOVED***
		// TODO: Should we be logging if writing status failed here, like above?
		// Should the logging be in WriteStatus?  Should we ignore the WriteStatus
		// error or allow the stats handler to see it?
		return t.WriteStatus(stream, status.New(codes.OK, ""))
	***REMOVED***
***REMOVED***

func (s *Server) processStreamingRPC(t transport.ServerTransport, stream *transport.Stream, srv *service, sd *StreamDesc, trInfo *traceInfo) (err error) ***REMOVED***
	sh := s.opts.statsHandler
	if sh != nil ***REMOVED***
		begin := &stats.Begin***REMOVED***
			BeginTime: time.Now(),
		***REMOVED***
		sh.HandleRPC(stream.Context(), begin)
	***REMOVED***
	defer func() ***REMOVED***
		if sh != nil ***REMOVED***
			end := &stats.End***REMOVED***
				EndTime: time.Now(),
			***REMOVED***
			if err != nil && err != io.EOF ***REMOVED***
				end.Error = toRPCErr(err)
			***REMOVED***
			sh.HandleRPC(stream.Context(), end)
		***REMOVED***
	***REMOVED***()
	if s.opts.cp != nil ***REMOVED***
		stream.SetSendCompress(s.opts.cp.Type())
	***REMOVED***
	ss := &serverStream***REMOVED***
		t:            t,
		s:            stream,
		p:            &parser***REMOVED***r: stream***REMOVED***,
		codec:        s.opts.codec,
		cp:           s.opts.cp,
		dc:           s.opts.dc,
		maxMsgSize:   s.opts.maxMsgSize,
		trInfo:       trInfo,
		statsHandler: sh,
	***REMOVED***
	if ss.cp != nil ***REMOVED***
		ss.cbuf = new(bytes.Buffer)
	***REMOVED***
	if trInfo != nil ***REMOVED***
		trInfo.tr.LazyLog(&trInfo.firstLine, false)
		defer func() ***REMOVED***
			ss.mu.Lock()
			if err != nil && err != io.EOF ***REMOVED***
				ss.trInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***err***REMOVED******REMOVED***, true)
				ss.trInfo.tr.SetError()
			***REMOVED***
			ss.trInfo.tr.Finish()
			ss.trInfo.tr = nil
			ss.mu.Unlock()
		***REMOVED***()
	***REMOVED***
	var appErr error
	var server interface***REMOVED******REMOVED***
	if srv != nil ***REMOVED***
		server = srv.server
	***REMOVED***
	if s.opts.streamInt == nil ***REMOVED***
		appErr = sd.Handler(server, ss)
	***REMOVED*** else ***REMOVED***
		info := &StreamServerInfo***REMOVED***
			FullMethod:     stream.Method(),
			IsClientStream: sd.ClientStreams,
			IsServerStream: sd.ServerStreams,
		***REMOVED***
		appErr = s.opts.streamInt(server, ss, info, sd.Handler)
	***REMOVED***
	if appErr != nil ***REMOVED***
		appStatus, ok := status.FromError(appErr)
		if !ok ***REMOVED***
			switch err := appErr.(type) ***REMOVED***
			case transport.StreamError:
				appStatus = status.New(err.Code, err.Desc)
			default:
				appStatus = status.New(convertCode(appErr), appErr.Error())
			***REMOVED***
			appErr = appStatus.Err()
		***REMOVED***
		if trInfo != nil ***REMOVED***
			ss.mu.Lock()
			ss.trInfo.tr.LazyLog(stringer(appStatus.Message()), true)
			ss.trInfo.tr.SetError()
			ss.mu.Unlock()
		***REMOVED***
		t.WriteStatus(ss.s, appStatus)
		// TODO: Should we log an error from WriteStatus here and below?
		return appErr
	***REMOVED***
	if trInfo != nil ***REMOVED***
		ss.mu.Lock()
		ss.trInfo.tr.LazyLog(stringer("OK"), false)
		ss.mu.Unlock()
	***REMOVED***
	return t.WriteStatus(ss.s, status.New(codes.OK, ""))

***REMOVED***

func (s *Server) handleStream(t transport.ServerTransport, stream *transport.Stream, trInfo *traceInfo) ***REMOVED***
	sm := stream.Method()
	if sm != "" && sm[0] == '/' ***REMOVED***
		sm = sm[1:]
	***REMOVED***
	pos := strings.LastIndex(sm, "/")
	if pos == -1 ***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.LazyLog(&fmtStringer***REMOVED***"Malformed method name %q", []interface***REMOVED******REMOVED******REMOVED***sm***REMOVED******REMOVED***, true)
			trInfo.tr.SetError()
		***REMOVED***
		errDesc := fmt.Sprintf("malformed method name: %q", stream.Method())
		if err := t.WriteStatus(stream, status.New(codes.InvalidArgument, errDesc)); err != nil ***REMOVED***
			if trInfo != nil ***REMOVED***
				trInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***err***REMOVED******REMOVED***, true)
				trInfo.tr.SetError()
			***REMOVED***
			grpclog.Printf("grpc: Server.handleStream failed to write status: %v", err)
		***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.Finish()
		***REMOVED***
		return
	***REMOVED***
	service := sm[:pos]
	method := sm[pos+1:]
	srv, ok := s.m[service]
	if !ok ***REMOVED***
		if unknownDesc := s.opts.unknownStreamDesc; unknownDesc != nil ***REMOVED***
			s.processStreamingRPC(t, stream, nil, unknownDesc, trInfo)
			return
		***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.LazyLog(&fmtStringer***REMOVED***"Unknown service %v", []interface***REMOVED******REMOVED******REMOVED***service***REMOVED******REMOVED***, true)
			trInfo.tr.SetError()
		***REMOVED***
		errDesc := fmt.Sprintf("unknown service %v", service)
		if err := t.WriteStatus(stream, status.New(codes.Unimplemented, errDesc)); err != nil ***REMOVED***
			if trInfo != nil ***REMOVED***
				trInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***err***REMOVED******REMOVED***, true)
				trInfo.tr.SetError()
			***REMOVED***
			grpclog.Printf("grpc: Server.handleStream failed to write status: %v", err)
		***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.Finish()
		***REMOVED***
		return
	***REMOVED***
	// Unary RPC or Streaming RPC?
	if md, ok := srv.md[method]; ok ***REMOVED***
		s.processUnaryRPC(t, stream, srv, md, trInfo)
		return
	***REMOVED***
	if sd, ok := srv.sd[method]; ok ***REMOVED***
		s.processStreamingRPC(t, stream, srv, sd, trInfo)
		return
	***REMOVED***
	if trInfo != nil ***REMOVED***
		trInfo.tr.LazyLog(&fmtStringer***REMOVED***"Unknown method %v", []interface***REMOVED******REMOVED******REMOVED***method***REMOVED******REMOVED***, true)
		trInfo.tr.SetError()
	***REMOVED***
	if unknownDesc := s.opts.unknownStreamDesc; unknownDesc != nil ***REMOVED***
		s.processStreamingRPC(t, stream, nil, unknownDesc, trInfo)
		return
	***REMOVED***
	errDesc := fmt.Sprintf("unknown method %v", method)
	if err := t.WriteStatus(stream, status.New(codes.Unimplemented, errDesc)); err != nil ***REMOVED***
		if trInfo != nil ***REMOVED***
			trInfo.tr.LazyLog(&fmtStringer***REMOVED***"%v", []interface***REMOVED******REMOVED******REMOVED***err***REMOVED******REMOVED***, true)
			trInfo.tr.SetError()
		***REMOVED***
		grpclog.Printf("grpc: Server.handleStream failed to write status: %v", err)
	***REMOVED***
	if trInfo != nil ***REMOVED***
		trInfo.tr.Finish()
	***REMOVED***
***REMOVED***

// Stop stops the gRPC server. It immediately closes all open
// connections and listeners.
// It cancels all active RPCs on the server side and the corresponding
// pending RPCs on the client side will get notified by connection
// errors.
func (s *Server) Stop() ***REMOVED***
	s.mu.Lock()
	listeners := s.lis
	s.lis = nil
	st := s.conns
	s.conns = nil
	// interrupt GracefulStop if Stop and GracefulStop are called concurrently.
	s.cv.Broadcast()
	s.mu.Unlock()

	for lis := range listeners ***REMOVED***
		lis.Close()
	***REMOVED***
	for c := range st ***REMOVED***
		c.Close()
	***REMOVED***

	s.mu.Lock()
	s.cancel()
	if s.events != nil ***REMOVED***
		s.events.Finish()
		s.events = nil
	***REMOVED***
	s.mu.Unlock()
***REMOVED***

// GracefulStop stops the gRPC server gracefully. It stops the server to accept new
// connections and RPCs and blocks until all the pending RPCs are finished.
func (s *Server) GracefulStop() ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns == nil ***REMOVED***
		return
	***REMOVED***
	for lis := range s.lis ***REMOVED***
		lis.Close()
	***REMOVED***
	s.lis = nil
	s.cancel()
	if !s.drain ***REMOVED***
		for c := range s.conns ***REMOVED***
			c.(transport.ServerTransport).Drain()
		***REMOVED***
		s.drain = true
	***REMOVED***
	for len(s.conns) != 0 ***REMOVED***
		s.cv.Wait()
	***REMOVED***
	s.conns = nil
	if s.events != nil ***REMOVED***
		s.events.Finish()
		s.events = nil
	***REMOVED***
***REMOVED***

func init() ***REMOVED***
	internal.TestingCloseConns = func(arg interface***REMOVED******REMOVED***) ***REMOVED***
		arg.(*Server).testingCloseConns()
	***REMOVED***
	internal.TestingUseHandlerImpl = func(arg interface***REMOVED******REMOVED***) ***REMOVED***
		arg.(*Server).opts.useHandlerImpl = true
	***REMOVED***
***REMOVED***

// testingCloseConns closes all existing transports but keeps s.lis
// accepting new connections.
func (s *Server) testingCloseConns() ***REMOVED***
	s.mu.Lock()
	for c := range s.conns ***REMOVED***
		c.Close()
		delete(s.conns, c)
	***REMOVED***
	s.mu.Unlock()
***REMOVED***

// SetHeader sets the header metadata.
// When called multiple times, all the provided metadata will be merged.
// All the metadata will be sent out when one of the following happens:
//  - grpc.SendHeader() is called;
//  - The first response is sent out;
//  - An RPC status is sent out (error or success).
func SetHeader(ctx context.Context, md metadata.MD) error ***REMOVED***
	if md.Len() == 0 ***REMOVED***
		return nil
	***REMOVED***
	stream, ok := transport.StreamFromContext(ctx)
	if !ok ***REMOVED***
		return Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	***REMOVED***
	return stream.SetHeader(md)
***REMOVED***

// SendHeader sends header metadata. It may be called at most once.
// The provided md and headers set by SetHeader() will be sent.
func SendHeader(ctx context.Context, md metadata.MD) error ***REMOVED***
	stream, ok := transport.StreamFromContext(ctx)
	if !ok ***REMOVED***
		return Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	***REMOVED***
	t := stream.ServerTransport()
	if t == nil ***REMOVED***
		grpclog.Fatalf("grpc: SendHeader: %v has no ServerTransport to send header metadata.", stream)
	***REMOVED***
	if err := t.WriteHeader(stream, md); err != nil ***REMOVED***
		return toRPCErr(err)
	***REMOVED***
	return nil
***REMOVED***

// SetTrailer sets the trailer metadata that will be sent when an RPC returns.
// When called more than once, all the provided metadata will be merged.
func SetTrailer(ctx context.Context, md metadata.MD) error ***REMOVED***
	if md.Len() == 0 ***REMOVED***
		return nil
	***REMOVED***
	stream, ok := transport.StreamFromContext(ctx)
	if !ok ***REMOVED***
		return Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	***REMOVED***
	return stream.SetTrailer(md)
***REMOVED***
