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
	"errors"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/transport"
)

var (
	// ErrClientConnClosing indicates that the operation is illegal because
	// the ClientConn is closing.
	ErrClientConnClosing = errors.New("grpc: the client connection is closing")
	// ErrClientConnTimeout indicates that the ClientConn cannot establish the
	// underlying connections within the specified timeout.
	// DEPRECATED: Please use context.DeadlineExceeded instead. This error will be
	// removed in Q1 2017.
	ErrClientConnTimeout = errors.New("grpc: timed out when dialing")

	// errNoTransportSecurity indicates that there is no transport security
	// being set for ClientConn. Users should either set one or explicitly
	// call WithInsecure DialOption to disable security.
	errNoTransportSecurity = errors.New("grpc: no transport security set (use grpc.WithInsecure() explicitly or set credentials)")
	// errTransportCredentialsMissing indicates that users want to transmit security
	// information (e.g., oauth2 token) which requires secure connection on an insecure
	// connection.
	errTransportCredentialsMissing = errors.New("grpc: the credentials require transport level security (use grpc.WithTransportCredentials() to set)")
	// errCredentialsConflict indicates that grpc.WithTransportCredentials()
	// and grpc.WithInsecure() are both called for a connection.
	errCredentialsConflict = errors.New("grpc: transport credentials are set for an insecure connection (grpc.WithTransportCredentials() and grpc.WithInsecure() are both called)")
	// errNetworkIO indicates that the connection is down due to some network I/O error.
	errNetworkIO = errors.New("grpc: failed with network I/O error")
	// errConnDrain indicates that the connection starts to be drained and does not accept any new RPCs.
	errConnDrain = errors.New("grpc: the connection is drained")
	// errConnClosing indicates that the connection is closing.
	errConnClosing = errors.New("grpc: the connection is closing")
	// errConnUnavailable indicates that the connection is unavailable.
	errConnUnavailable = errors.New("grpc: the connection is unavailable")
	// minimum time to give a connection to complete
	minConnectTimeout = 20 * time.Second
)

// dialOptions configure a Dial call. dialOptions are set by the DialOption
// values passed to Dial.
type dialOptions struct ***REMOVED***
	unaryInt   UnaryClientInterceptor
	streamInt  StreamClientInterceptor
	codec      Codec
	cp         Compressor
	dc         Decompressor
	bs         backoffStrategy
	balancer   Balancer
	block      bool
	insecure   bool
	timeout    time.Duration
	scChan     <-chan ServiceConfig
	copts      transport.ConnectOptions
	maxMsgSize int
***REMOVED***

const defaultClientMaxMsgSize = math.MaxInt32

// DialOption configures how we set up the connection.
type DialOption func(*dialOptions)

// WithMaxMsgSize returns a DialOption which sets the maximum message size the client can receive.
func WithMaxMsgSize(s int) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.maxMsgSize = s
	***REMOVED***
***REMOVED***

// WithCodec returns a DialOption which sets a codec for message marshaling and unmarshaling.
func WithCodec(c Codec) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.codec = c
	***REMOVED***
***REMOVED***

// WithCompressor returns a DialOption which sets a CompressorGenerator for generating message
// compressor.
func WithCompressor(cp Compressor) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.cp = cp
	***REMOVED***
***REMOVED***

// WithDecompressor returns a DialOption which sets a DecompressorGenerator for generating
// message decompressor.
func WithDecompressor(dc Decompressor) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.dc = dc
	***REMOVED***
***REMOVED***

// WithBalancer returns a DialOption which sets a load balancer.
func WithBalancer(b Balancer) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.balancer = b
	***REMOVED***
***REMOVED***

// WithServiceConfig returns a DialOption which has a channel to read the service configuration.
func WithServiceConfig(c <-chan ServiceConfig) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.scChan = c
	***REMOVED***
***REMOVED***

// WithBackoffMaxDelay configures the dialer to use the provided maximum delay
// when backing off after failed connection attempts.
func WithBackoffMaxDelay(md time.Duration) DialOption ***REMOVED***
	return WithBackoffConfig(BackoffConfig***REMOVED***MaxDelay: md***REMOVED***)
***REMOVED***

// WithBackoffConfig configures the dialer to use the provided backoff
// parameters after connection failures.
//
// Use WithBackoffMaxDelay until more parameters on BackoffConfig are opened up
// for use.
func WithBackoffConfig(b BackoffConfig) DialOption ***REMOVED***
	// Set defaults to ensure that provided BackoffConfig is valid and
	// unexported fields get default values.
	setDefaults(&b)
	return withBackoff(b)
***REMOVED***

// withBackoff sets the backoff strategy used for retries after a
// failed connection attempt.
//
// This can be exported if arbitrary backoff strategies are allowed by gRPC.
func withBackoff(bs backoffStrategy) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.bs = bs
	***REMOVED***
***REMOVED***

// WithBlock returns a DialOption which makes caller of Dial blocks until the underlying
// connection is up. Without this, Dial returns immediately and connecting the server
// happens in background.
func WithBlock() DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.block = true
	***REMOVED***
***REMOVED***

// WithInsecure returns a DialOption which disables transport security for this ClientConn.
// Note that transport security is required unless WithInsecure is set.
func WithInsecure() DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.insecure = true
	***REMOVED***
***REMOVED***

// WithTransportCredentials returns a DialOption which configures a
// connection level security credentials (e.g., TLS/SSL).
func WithTransportCredentials(creds credentials.TransportCredentials) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.copts.TransportCredentials = creds
	***REMOVED***
***REMOVED***

// WithPerRPCCredentials returns a DialOption which sets
// credentials which will place auth state on each outbound RPC.
func WithPerRPCCredentials(creds credentials.PerRPCCredentials) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.copts.PerRPCCredentials = append(o.copts.PerRPCCredentials, creds)
	***REMOVED***
***REMOVED***

// WithTimeout returns a DialOption that configures a timeout for dialing a ClientConn
// initially. This is valid if and only if WithBlock() is present.
func WithTimeout(d time.Duration) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.timeout = d
	***REMOVED***
***REMOVED***

// WithDialer returns a DialOption that specifies a function to use for dialing network addresses.
// If FailOnNonTempDialError() is set to true, and an error is returned by f, gRPC checks the error's
// Temporary() method to decide if it should try to reconnect to the network address.
func WithDialer(f func(string, time.Duration) (net.Conn, error)) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.copts.Dialer = func(ctx context.Context, addr string) (net.Conn, error) ***REMOVED***
			if deadline, ok := ctx.Deadline(); ok ***REMOVED***
				return f(addr, deadline.Sub(time.Now()))
			***REMOVED***
			return f(addr, 0)
		***REMOVED***
	***REMOVED***
***REMOVED***

// WithStatsHandler returns a DialOption that specifies the stats handler
// for all the RPCs and underlying network connections in this ClientConn.
func WithStatsHandler(h stats.Handler) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.copts.StatsHandler = h
	***REMOVED***
***REMOVED***

// FailOnNonTempDialError returns a DialOption that specified if gRPC fails on non-temporary dial errors.
// If f is true, and dialer returns a non-temporary error, gRPC will fail the connection to the network
// address and won't try to reconnect.
// The default value of FailOnNonTempDialError is false.
// This is an EXPERIMENTAL API.
func FailOnNonTempDialError(f bool) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.copts.FailOnNonTempDialError = f
	***REMOVED***
***REMOVED***

// WithUserAgent returns a DialOption that specifies a user agent string for all the RPCs.
func WithUserAgent(s string) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.copts.UserAgent = s
	***REMOVED***
***REMOVED***

// WithKeepaliveParams returns a DialOption that specifies keepalive paramaters for the client transport.
func WithKeepaliveParams(kp keepalive.ClientParameters) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.copts.KeepaliveParams = kp
	***REMOVED***
***REMOVED***

// WithUnaryInterceptor returns a DialOption that specifies the interceptor for unary RPCs.
func WithUnaryInterceptor(f UnaryClientInterceptor) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.unaryInt = f
	***REMOVED***
***REMOVED***

// WithStreamInterceptor returns a DialOption that specifies the interceptor for streaming RPCs.
func WithStreamInterceptor(f StreamClientInterceptor) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.streamInt = f
	***REMOVED***
***REMOVED***

// WithAuthority returns a DialOption that specifies the value to be used as
// the :authority pseudo-header. This value only works with WithInsecure and
// has no effect if TransportCredentials are present.
func WithAuthority(a string) DialOption ***REMOVED***
	return func(o *dialOptions) ***REMOVED***
		o.copts.Authority = a
	***REMOVED***
***REMOVED***

// Dial creates a client connection to the given target.
func Dial(target string, opts ...DialOption) (*ClientConn, error) ***REMOVED***
	return DialContext(context.Background(), target, opts...)
***REMOVED***

// DialContext creates a client connection to the given target. ctx can be used to
// cancel or expire the pending connecting. Once this function returns, the
// cancellation and expiration of ctx will be noop. Users should call ClientConn.Close
// to terminate all the pending operations after this function returns.
// This is the EXPERIMENTAL API.
func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) ***REMOVED***
	cc := &ClientConn***REMOVED***
		target: target,
		conns:  make(map[Address]*addrConn),
	***REMOVED***
	cc.ctx, cc.cancel = context.WithCancel(context.Background())
	cc.dopts.maxMsgSize = defaultClientMaxMsgSize
	for _, opt := range opts ***REMOVED***
		opt(&cc.dopts)
	***REMOVED***
	cc.mkp = cc.dopts.copts.KeepaliveParams

	if cc.dopts.copts.Dialer == nil ***REMOVED***
		cc.dopts.copts.Dialer = newProxyDialer(
			func(ctx context.Context, addr string) (net.Conn, error) ***REMOVED***
				return dialContext(ctx, "tcp", addr)
			***REMOVED***,
		)
	***REMOVED***

	if cc.dopts.copts.UserAgent != "" ***REMOVED***
		cc.dopts.copts.UserAgent += " " + grpcUA
	***REMOVED*** else ***REMOVED***
		cc.dopts.copts.UserAgent = grpcUA
	***REMOVED***

	if cc.dopts.timeout > 0 ***REMOVED***
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cc.dopts.timeout)
		defer cancel()
	***REMOVED***

	defer func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			conn, err = nil, ctx.Err()
		default:
		***REMOVED***

		if err != nil ***REMOVED***
			cc.Close()
		***REMOVED***
	***REMOVED***()

	if cc.dopts.scChan != nil ***REMOVED***
		// Wait for the initial service config.
		select ***REMOVED***
		case sc, ok := <-cc.dopts.scChan:
			if ok ***REMOVED***
				cc.sc = sc
			***REMOVED***
		case <-ctx.Done():
			return nil, ctx.Err()
		***REMOVED***
	***REMOVED***
	// Set defaults.
	if cc.dopts.codec == nil ***REMOVED***
		cc.dopts.codec = protoCodec***REMOVED******REMOVED***
	***REMOVED***
	if cc.dopts.bs == nil ***REMOVED***
		cc.dopts.bs = DefaultBackoffConfig
	***REMOVED***
	creds := cc.dopts.copts.TransportCredentials
	if creds != nil && creds.Info().ServerName != "" ***REMOVED***
		cc.authority = creds.Info().ServerName
	***REMOVED*** else if cc.dopts.insecure && cc.dopts.copts.Authority != "" ***REMOVED***
		cc.authority = cc.dopts.copts.Authority
	***REMOVED*** else ***REMOVED***
		cc.authority = target
	***REMOVED***
	waitC := make(chan error, 1)
	go func() ***REMOVED***
		defer close(waitC)
		if cc.dopts.balancer == nil && cc.sc.LB != nil ***REMOVED***
			cc.dopts.balancer = cc.sc.LB
		***REMOVED***
		if cc.dopts.balancer != nil ***REMOVED***
			var credsClone credentials.TransportCredentials
			if creds != nil ***REMOVED***
				credsClone = creds.Clone()
			***REMOVED***
			config := BalancerConfig***REMOVED***
				DialCreds: credsClone,
			***REMOVED***
			if err := cc.dopts.balancer.Start(target, config); err != nil ***REMOVED***
				waitC <- err
				return
			***REMOVED***
			ch := cc.dopts.balancer.Notify()
			if ch != nil ***REMOVED***
				if cc.dopts.block ***REMOVED***
					doneChan := make(chan struct***REMOVED******REMOVED***)
					go cc.lbWatcher(doneChan)
					<-doneChan
				***REMOVED*** else ***REMOVED***
					go cc.lbWatcher(nil)
				***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		// No balancer, or no resolver within the balancer.  Connect directly.
		if err := cc.resetAddrConn(Address***REMOVED***Addr: target***REMOVED***, cc.dopts.block, nil); err != nil ***REMOVED***
			waitC <- err
			return
		***REMOVED***
	***REMOVED***()
	select ***REMOVED***
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-waitC:
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if cc.dopts.scChan != nil ***REMOVED***
		go cc.scWatcher()
	***REMOVED***

	return cc, nil
***REMOVED***

// ConnectivityState indicates the state of a client connection.
type ConnectivityState int

const (
	// Idle indicates the ClientConn is idle.
	Idle ConnectivityState = iota
	// Connecting indicates the ClienConn is connecting.
	Connecting
	// Ready indicates the ClientConn is ready for work.
	Ready
	// TransientFailure indicates the ClientConn has seen a failure but expects to recover.
	TransientFailure
	// Shutdown indicates the ClientConn has started shutting down.
	Shutdown
)

func (s ConnectivityState) String() string ***REMOVED***
	switch s ***REMOVED***
	case Idle:
		return "IDLE"
	case Connecting:
		return "CONNECTING"
	case Ready:
		return "READY"
	case TransientFailure:
		return "TRANSIENT_FAILURE"
	case Shutdown:
		return "SHUTDOWN"
	default:
		panic(fmt.Sprintf("unknown connectivity state: %d", s))
	***REMOVED***
***REMOVED***

// ClientConn represents a client connection to an RPC server.
type ClientConn struct ***REMOVED***
	ctx    context.Context
	cancel context.CancelFunc

	target    string
	authority string
	dopts     dialOptions

	mu    sync.RWMutex
	sc    ServiceConfig
	conns map[Address]*addrConn
	// Keepalive parameter can be udated if a GoAway is received.
	mkp keepalive.ClientParameters
***REMOVED***

// lbWatcher watches the Notify channel of the balancer in cc and manages
// connections accordingly.  If doneChan is not nil, it is closed after the
// first successfull connection is made.
func (cc *ClientConn) lbWatcher(doneChan chan struct***REMOVED******REMOVED***) ***REMOVED***
	for addrs := range cc.dopts.balancer.Notify() ***REMOVED***
		var (
			add []Address   // Addresses need to setup connections.
			del []*addrConn // Connections need to tear down.
		)
		cc.mu.Lock()
		for _, a := range addrs ***REMOVED***
			if _, ok := cc.conns[a]; !ok ***REMOVED***
				add = append(add, a)
			***REMOVED***
		***REMOVED***
		for k, c := range cc.conns ***REMOVED***
			var keep bool
			for _, a := range addrs ***REMOVED***
				if k == a ***REMOVED***
					keep = true
					break
				***REMOVED***
			***REMOVED***
			if !keep ***REMOVED***
				del = append(del, c)
				delete(cc.conns, c.addr)
			***REMOVED***
		***REMOVED***
		cc.mu.Unlock()
		for _, a := range add ***REMOVED***
			if doneChan != nil ***REMOVED***
				err := cc.resetAddrConn(a, true, nil)
				if err == nil ***REMOVED***
					close(doneChan)
					doneChan = nil
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				cc.resetAddrConn(a, false, nil)
			***REMOVED***
		***REMOVED***
		for _, c := range del ***REMOVED***
			c.tearDown(errConnDrain)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (cc *ClientConn) scWatcher() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case sc, ok := <-cc.dopts.scChan:
			if !ok ***REMOVED***
				return
			***REMOVED***
			cc.mu.Lock()
			// TODO: load balance policy runtime change is ignored.
			// We may revist this decision in the future.
			cc.sc = sc
			cc.mu.Unlock()
		case <-cc.ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// resetAddrConn creates an addrConn for addr and adds it to cc.conns.
// If there is an old addrConn for addr, it will be torn down, using tearDownErr as the reason.
// If tearDownErr is nil, errConnDrain will be used instead.
func (cc *ClientConn) resetAddrConn(addr Address, block bool, tearDownErr error) error ***REMOVED***
	ac := &addrConn***REMOVED***
		cc:    cc,
		addr:  addr,
		dopts: cc.dopts,
	***REMOVED***
	cc.mu.RLock()
	ac.dopts.copts.KeepaliveParams = cc.mkp
	cc.mu.RUnlock()
	ac.ctx, ac.cancel = context.WithCancel(cc.ctx)
	ac.stateCV = sync.NewCond(&ac.mu)
	if EnableTracing ***REMOVED***
		ac.events = trace.NewEventLog("grpc.ClientConn", ac.addr.Addr)
	***REMOVED***
	if !ac.dopts.insecure ***REMOVED***
		if ac.dopts.copts.TransportCredentials == nil ***REMOVED***
			return errNoTransportSecurity
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if ac.dopts.copts.TransportCredentials != nil ***REMOVED***
			return errCredentialsConflict
		***REMOVED***
		for _, cd := range ac.dopts.copts.PerRPCCredentials ***REMOVED***
			if cd.RequireTransportSecurity() ***REMOVED***
				return errTransportCredentialsMissing
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Track ac in cc. This needs to be done before any getTransport(...) is called.
	cc.mu.Lock()
	if cc.conns == nil ***REMOVED***
		cc.mu.Unlock()
		return ErrClientConnClosing
	***REMOVED***
	stale := cc.conns[ac.addr]
	cc.conns[ac.addr] = ac
	cc.mu.Unlock()
	if stale != nil ***REMOVED***
		// There is an addrConn alive on ac.addr already. This could be due to
		// 1) a buggy Balancer notifies duplicated Addresses;
		// 2) goaway was received, a new ac will replace the old ac.
		//    The old ac should be deleted from cc.conns, but the
		//    underlying transport should drain rather than close.
		if tearDownErr == nil ***REMOVED***
			// tearDownErr is nil if resetAddrConn is called by
			// 1) Dial
			// 2) lbWatcher
			// In both cases, the stale ac should drain, not close.
			stale.tearDown(errConnDrain)
		***REMOVED*** else ***REMOVED***
			stale.tearDown(tearDownErr)
		***REMOVED***
	***REMOVED***
	if block ***REMOVED***
		if err := ac.resetTransport(false); err != nil ***REMOVED***
			if err != errConnClosing ***REMOVED***
				// Tear down ac and delete it from cc.conns.
				cc.mu.Lock()
				delete(cc.conns, ac.addr)
				cc.mu.Unlock()
				ac.tearDown(err)
			***REMOVED***
			if e, ok := err.(transport.ConnectionError); ok && !e.Temporary() ***REMOVED***
				return e.Origin()
			***REMOVED***
			return err
		***REMOVED***
		// Start to monitor the error status of transport.
		go ac.transportMonitor()
	***REMOVED*** else ***REMOVED***
		// Start a goroutine connecting to the server asynchronously.
		go func() ***REMOVED***
			if err := ac.resetTransport(false); err != nil ***REMOVED***
				grpclog.Printf("Failed to dial %s: %v; please retry.", ac.addr.Addr, err)
				if err != errConnClosing ***REMOVED***
					// Keep this ac in cc.conns, to get the reason it's torn down.
					ac.tearDown(err)
				***REMOVED***
				return
			***REMOVED***
			ac.transportMonitor()
		***REMOVED***()
	***REMOVED***
	return nil
***REMOVED***

// TODO: Avoid the locking here.
func (cc *ClientConn) getMethodConfig(method string) (m MethodConfig, ok bool) ***REMOVED***
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	m, ok = cc.sc.Methods[method]
	return
***REMOVED***

func (cc *ClientConn) getTransport(ctx context.Context, opts BalancerGetOptions) (transport.ClientTransport, func(), error) ***REMOVED***
	var (
		ac  *addrConn
		ok  bool
		put func()
	)
	if cc.dopts.balancer == nil ***REMOVED***
		// If balancer is nil, there should be only one addrConn available.
		cc.mu.RLock()
		if cc.conns == nil ***REMOVED***
			cc.mu.RUnlock()
			return nil, nil, toRPCErr(ErrClientConnClosing)
		***REMOVED***
		for _, ac = range cc.conns ***REMOVED***
			// Break after the first iteration to get the first addrConn.
			ok = true
			break
		***REMOVED***
		cc.mu.RUnlock()
	***REMOVED*** else ***REMOVED***
		var (
			addr Address
			err  error
		)
		addr, put, err = cc.dopts.balancer.Get(ctx, opts)
		if err != nil ***REMOVED***
			return nil, nil, toRPCErr(err)
		***REMOVED***
		cc.mu.RLock()
		if cc.conns == nil ***REMOVED***
			cc.mu.RUnlock()
			return nil, nil, toRPCErr(ErrClientConnClosing)
		***REMOVED***
		ac, ok = cc.conns[addr]
		cc.mu.RUnlock()
	***REMOVED***
	if !ok ***REMOVED***
		if put != nil ***REMOVED***
			updateRPCInfoInContext(ctx, rpcInfo***REMOVED***bytesSent: false, bytesReceived: false***REMOVED***)
			put()
		***REMOVED***
		return nil, nil, errConnClosing
	***REMOVED***
	t, err := ac.wait(ctx, cc.dopts.balancer != nil, !opts.BlockingWait)
	if err != nil ***REMOVED***
		if put != nil ***REMOVED***
			updateRPCInfoInContext(ctx, rpcInfo***REMOVED***bytesSent: false, bytesReceived: false***REMOVED***)
			put()
		***REMOVED***
		return nil, nil, err
	***REMOVED***
	return t, put, nil
***REMOVED***

// Close tears down the ClientConn and all underlying connections.
func (cc *ClientConn) Close() error ***REMOVED***
	cc.cancel()

	cc.mu.Lock()
	if cc.conns == nil ***REMOVED***
		cc.mu.Unlock()
		return ErrClientConnClosing
	***REMOVED***
	conns := cc.conns
	cc.conns = nil
	cc.mu.Unlock()
	if cc.dopts.balancer != nil ***REMOVED***
		cc.dopts.balancer.Close()
	***REMOVED***
	for _, ac := range conns ***REMOVED***
		ac.tearDown(ErrClientConnClosing)
	***REMOVED***
	return nil
***REMOVED***

// addrConn is a network connection to a given address.
type addrConn struct ***REMOVED***
	ctx    context.Context
	cancel context.CancelFunc

	cc     *ClientConn
	addr   Address
	dopts  dialOptions
	events trace.EventLog

	mu      sync.Mutex
	state   ConnectivityState
	stateCV *sync.Cond
	down    func(error) // the handler called when a connection is down.
	// ready is closed and becomes nil when a new transport is up or failed
	// due to timeout.
	ready     chan struct***REMOVED******REMOVED***
	transport transport.ClientTransport

	// The reason this addrConn is torn down.
	tearDownErr error
***REMOVED***

// adjustParams updates parameters used to create transports upon
// receiving a GoAway.
func (ac *addrConn) adjustParams(r transport.GoAwayReason) ***REMOVED***
	switch r ***REMOVED***
	case transport.TooManyPings:
		v := 2 * ac.dopts.copts.KeepaliveParams.Time
		ac.cc.mu.Lock()
		if v > ac.cc.mkp.Time ***REMOVED***
			ac.cc.mkp.Time = v
		***REMOVED***
		ac.cc.mu.Unlock()
	***REMOVED***
***REMOVED***

// printf records an event in ac's event log, unless ac has been closed.
// REQUIRES ac.mu is held.
func (ac *addrConn) printf(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	if ac.events != nil ***REMOVED***
		ac.events.Printf(format, a...)
	***REMOVED***
***REMOVED***

// errorf records an error in ac's event log, unless ac has been closed.
// REQUIRES ac.mu is held.
func (ac *addrConn) errorf(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	if ac.events != nil ***REMOVED***
		ac.events.Errorf(format, a...)
	***REMOVED***
***REMOVED***

// getState returns the connectivity state of the Conn
func (ac *addrConn) getState() ConnectivityState ***REMOVED***
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.state
***REMOVED***

// waitForStateChange blocks until the state changes to something other than the sourceState.
func (ac *addrConn) waitForStateChange(ctx context.Context, sourceState ConnectivityState) (ConnectivityState, error) ***REMOVED***
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if sourceState != ac.state ***REMOVED***
		return ac.state, nil
	***REMOVED***
	done := make(chan struct***REMOVED******REMOVED***)
	var err error
	go func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			ac.mu.Lock()
			err = ctx.Err()
			ac.stateCV.Broadcast()
			ac.mu.Unlock()
		case <-done:
		***REMOVED***
	***REMOVED***()
	defer close(done)
	for sourceState == ac.state ***REMOVED***
		ac.stateCV.Wait()
		if err != nil ***REMOVED***
			return ac.state, err
		***REMOVED***
	***REMOVED***
	return ac.state, nil
***REMOVED***

func (ac *addrConn) resetTransport(closeTransport bool) error ***REMOVED***
	for retries := 0; ; retries++ ***REMOVED***
		ac.mu.Lock()
		ac.printf("connecting")
		if ac.state == Shutdown ***REMOVED***
			// ac.tearDown(...) has been invoked.
			ac.mu.Unlock()
			return errConnClosing
		***REMOVED***
		if ac.down != nil ***REMOVED***
			ac.down(downErrorf(false, true, "%v", errNetworkIO))
			ac.down = nil
		***REMOVED***
		ac.state = Connecting
		ac.stateCV.Broadcast()
		t := ac.transport
		ac.mu.Unlock()
		if closeTransport && t != nil ***REMOVED***
			t.Close()
		***REMOVED***
		sleepTime := ac.dopts.bs.backoff(retries)
		timeout := minConnectTimeout
		if timeout < sleepTime ***REMOVED***
			timeout = sleepTime
		***REMOVED***
		ctx, cancel := context.WithTimeout(ac.ctx, timeout)
		connectTime := time.Now()
		sinfo := transport.TargetInfo***REMOVED***
			Addr:     ac.addr.Addr,
			Metadata: ac.addr.Metadata,
		***REMOVED***
		newTransport, err := transport.NewClientTransport(ctx, sinfo, ac.dopts.copts)
		// Don't call cancel in success path due to a race in Go 1.6:
		// https://github.com/golang/go/issues/15078.
		if err != nil ***REMOVED***
			cancel()

			if e, ok := err.(transport.ConnectionError); ok && !e.Temporary() ***REMOVED***
				return err
			***REMOVED***
			grpclog.Printf("grpc: addrConn.resetTransport failed to create client transport: %v; Reconnecting to %v", err, ac.addr)
			ac.mu.Lock()
			if ac.state == Shutdown ***REMOVED***
				// ac.tearDown(...) has been invoked.
				ac.mu.Unlock()
				return errConnClosing
			***REMOVED***
			ac.errorf("transient failure: %v", err)
			ac.state = TransientFailure
			ac.stateCV.Broadcast()
			if ac.ready != nil ***REMOVED***
				close(ac.ready)
				ac.ready = nil
			***REMOVED***
			ac.mu.Unlock()
			closeTransport = false
			select ***REMOVED***
			case <-time.After(sleepTime - time.Since(connectTime)):
			case <-ac.ctx.Done():
				return ac.ctx.Err()
			***REMOVED***
			continue
		***REMOVED***
		ac.mu.Lock()
		ac.printf("ready")
		if ac.state == Shutdown ***REMOVED***
			// ac.tearDown(...) has been invoked.
			ac.mu.Unlock()
			newTransport.Close()
			return errConnClosing
		***REMOVED***
		ac.state = Ready
		ac.stateCV.Broadcast()
		ac.transport = newTransport
		if ac.ready != nil ***REMOVED***
			close(ac.ready)
			ac.ready = nil
		***REMOVED***
		if ac.cc.dopts.balancer != nil ***REMOVED***
			ac.down = ac.cc.dopts.balancer.Up(ac.addr)
		***REMOVED***
		ac.mu.Unlock()
		return nil
	***REMOVED***
***REMOVED***

// Run in a goroutine to track the error in transport and create the
// new transport if an error happens. It returns when the channel is closing.
func (ac *addrConn) transportMonitor() ***REMOVED***
	for ***REMOVED***
		ac.mu.Lock()
		t := ac.transport
		ac.mu.Unlock()
		select ***REMOVED***
		// This is needed to detect the teardown when
		// the addrConn is idle (i.e., no RPC in flight).
		case <-ac.ctx.Done():
			select ***REMOVED***
			case <-t.Error():
				t.Close()
			default:
			***REMOVED***
			return
		case <-t.GoAway():
			ac.adjustParams(t.GetGoAwayReason())
			// If GoAway happens without any network I/O error, ac is closed without shutting down the
			// underlying transport (the transport will be closed when all the pending RPCs finished or
			// failed.).
			// If GoAway and some network I/O error happen concurrently, ac and its underlying transport
			// are closed.
			// In both cases, a new ac is created.
			select ***REMOVED***
			case <-t.Error():
				ac.cc.resetAddrConn(ac.addr, false, errNetworkIO)
			default:
				ac.cc.resetAddrConn(ac.addr, false, errConnDrain)
			***REMOVED***
			return
		case <-t.Error():
			select ***REMOVED***
			case <-ac.ctx.Done():
				t.Close()
				return
			case <-t.GoAway():
				ac.adjustParams(t.GetGoAwayReason())
				ac.cc.resetAddrConn(ac.addr, false, errNetworkIO)
				return
			default:
			***REMOVED***
			ac.mu.Lock()
			if ac.state == Shutdown ***REMOVED***
				// ac has been shutdown.
				ac.mu.Unlock()
				return
			***REMOVED***
			ac.state = TransientFailure
			ac.stateCV.Broadcast()
			ac.mu.Unlock()
			if err := ac.resetTransport(true); err != nil ***REMOVED***
				ac.mu.Lock()
				ac.printf("transport exiting: %v", err)
				ac.mu.Unlock()
				grpclog.Printf("grpc: addrConn.transportMonitor exits due to: %v", err)
				if err != errConnClosing ***REMOVED***
					// Keep this ac in cc.conns, to get the reason it's torn down.
					ac.tearDown(err)
				***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// wait blocks until i) the new transport is up or ii) ctx is done or iii) ac is closed or
// iv) transport is in TransientFailure and there is a balancer/failfast is true.
func (ac *addrConn) wait(ctx context.Context, hasBalancer, failfast bool) (transport.ClientTransport, error) ***REMOVED***
	for ***REMOVED***
		ac.mu.Lock()
		switch ***REMOVED***
		case ac.state == Shutdown:
			if failfast || !hasBalancer ***REMOVED***
				// RPC is failfast or balancer is nil. This RPC should fail with ac.tearDownErr.
				err := ac.tearDownErr
				ac.mu.Unlock()
				return nil, err
			***REMOVED***
			ac.mu.Unlock()
			return nil, errConnClosing
		case ac.state == Ready:
			ct := ac.transport
			ac.mu.Unlock()
			return ct, nil
		case ac.state == TransientFailure:
			if failfast || hasBalancer ***REMOVED***
				ac.mu.Unlock()
				return nil, errConnUnavailable
			***REMOVED***
		***REMOVED***
		ready := ac.ready
		if ready == nil ***REMOVED***
			ready = make(chan struct***REMOVED******REMOVED***)
			ac.ready = ready
		***REMOVED***
		ac.mu.Unlock()
		select ***REMOVED***
		case <-ctx.Done():
			return nil, toRPCErr(ctx.Err())
		// Wait until the new transport is ready or failed.
		case <-ready:
		***REMOVED***
	***REMOVED***
***REMOVED***

// tearDown starts to tear down the addrConn.
// TODO(zhaoq): Make this synchronous to avoid unbounded memory consumption in
// some edge cases (e.g., the caller opens and closes many addrConn's in a
// tight loop.
// tearDown doesn't remove ac from ac.cc.conns.
func (ac *addrConn) tearDown(err error) ***REMOVED***
	ac.cancel()

	ac.mu.Lock()
	defer ac.mu.Unlock()
	if ac.down != nil ***REMOVED***
		ac.down(downErrorf(false, false, "%v", err))
		ac.down = nil
	***REMOVED***
	if err == errConnDrain && ac.transport != nil ***REMOVED***
		// GracefulClose(...) may be executed multiple times when
		// i) receiving multiple GoAway frames from the server; or
		// ii) there are concurrent name resolver/Balancer triggered
		// address removal and GoAway.
		ac.transport.GracefulClose()
	***REMOVED***
	if ac.state == Shutdown ***REMOVED***
		return
	***REMOVED***
	ac.state = Shutdown
	ac.tearDownErr = err
	ac.stateCV.Broadcast()
	if ac.events != nil ***REMOVED***
		ac.events.Finish()
		ac.events = nil
	***REMOVED***
	if ac.ready != nil ***REMOVED***
		close(ac.ready)
		ac.ready = nil
	***REMOVED***
	if ac.transport != nil && err != errConnDrain ***REMOVED***
		ac.transport.Close()
	***REMOVED***
	return
***REMOVED***
