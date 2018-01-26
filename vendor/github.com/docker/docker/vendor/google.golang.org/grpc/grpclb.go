/*
 *
 * Copyright 2016, Google Inc.
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
	"math/rand"
	"net"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	lbpb "google.golang.org/grpc/grpclb/grpc_lb_v1"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/naming"
)

// Client API for LoadBalancer service.
// Mostly copied from generated pb.go file.
// To avoid circular dependency.
type loadBalancerClient struct ***REMOVED***
	cc *ClientConn
***REMOVED***

func (c *loadBalancerClient) BalanceLoad(ctx context.Context, opts ...CallOption) (*balanceLoadClientStream, error) ***REMOVED***
	desc := &StreamDesc***REMOVED***
		StreamName:    "BalanceLoad",
		ServerStreams: true,
		ClientStreams: true,
	***REMOVED***
	stream, err := NewClientStream(ctx, desc, c.cc, "/grpc.lb.v1.LoadBalancer/BalanceLoad", opts...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	x := &balanceLoadClientStream***REMOVED***stream***REMOVED***
	return x, nil
***REMOVED***

type balanceLoadClientStream struct ***REMOVED***
	ClientStream
***REMOVED***

func (x *balanceLoadClientStream) Send(m *lbpb.LoadBalanceRequest) error ***REMOVED***
	return x.ClientStream.SendMsg(m)
***REMOVED***

func (x *balanceLoadClientStream) Recv() (*lbpb.LoadBalanceResponse, error) ***REMOVED***
	m := new(lbpb.LoadBalanceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return m, nil
***REMOVED***

// AddressType indicates the address type returned by name resolution.
type AddressType uint8

const (
	// Backend indicates the server is a backend server.
	Backend AddressType = iota
	// GRPCLB indicates the server is a grpclb load balancer.
	GRPCLB
)

// AddrMetadataGRPCLB contains the information the name resolution for grpclb should provide. The
// name resolver used by grpclb balancer is required to provide this type of metadata in
// its address updates.
type AddrMetadataGRPCLB struct ***REMOVED***
	// AddrType is the type of server (grpc load balancer or backend).
	AddrType AddressType
	// ServerName is the name of the grpc load balancer. Used for authentication.
	ServerName string
***REMOVED***

// NewGRPCLBBalancer creates a grpclb load balancer.
func NewGRPCLBBalancer(r naming.Resolver) Balancer ***REMOVED***
	return &balancer***REMOVED***
		r: r,
	***REMOVED***
***REMOVED***

type remoteBalancerInfo struct ***REMOVED***
	addr string
	// the server name used for authentication with the remote LB server.
	name string
***REMOVED***

// grpclbAddrInfo consists of the information of a backend server.
type grpclbAddrInfo struct ***REMOVED***
	addr      Address
	connected bool
	// dropForRateLimiting indicates whether this particular request should be
	// dropped by the client for rate limiting.
	dropForRateLimiting bool
	// dropForLoadBalancing indicates whether this particular request should be
	// dropped by the client for load balancing.
	dropForLoadBalancing bool
***REMOVED***

type balancer struct ***REMOVED***
	r        naming.Resolver
	target   string
	mu       sync.Mutex
	seq      int // a sequence number to make sure addrCh does not get stale addresses.
	w        naming.Watcher
	addrCh   chan []Address
	rbs      []remoteBalancerInfo
	addrs    []*grpclbAddrInfo
	next     int
	waitCh   chan struct***REMOVED******REMOVED***
	done     bool
	expTimer *time.Timer
	rand     *rand.Rand

	clientStats lbpb.ClientStats
***REMOVED***

func (b *balancer) watchAddrUpdates(w naming.Watcher, ch chan []remoteBalancerInfo) error ***REMOVED***
	updates, err := w.Next()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.done ***REMOVED***
		return ErrClientConnClosing
	***REMOVED***
	for _, update := range updates ***REMOVED***
		switch update.Op ***REMOVED***
		case naming.Add:
			var exist bool
			for _, v := range b.rbs ***REMOVED***
				// TODO: Is the same addr with different server name a different balancer?
				if update.Addr == v.addr ***REMOVED***
					exist = true
					break
				***REMOVED***
			***REMOVED***
			if exist ***REMOVED***
				continue
			***REMOVED***
			md, ok := update.Metadata.(*AddrMetadataGRPCLB)
			if !ok ***REMOVED***
				// TODO: Revisit the handling here and may introduce some fallback mechanism.
				grpclog.Printf("The name resolution contains unexpected metadata %v", update.Metadata)
				continue
			***REMOVED***
			switch md.AddrType ***REMOVED***
			case Backend:
				// TODO: Revisit the handling here and may introduce some fallback mechanism.
				grpclog.Printf("The name resolution does not give grpclb addresses")
				continue
			case GRPCLB:
				b.rbs = append(b.rbs, remoteBalancerInfo***REMOVED***
					addr: update.Addr,
					name: md.ServerName,
				***REMOVED***)
			default:
				grpclog.Printf("Received unknow address type %d", md.AddrType)
				continue
			***REMOVED***
		case naming.Delete:
			for i, v := range b.rbs ***REMOVED***
				if update.Addr == v.addr ***REMOVED***
					copy(b.rbs[i:], b.rbs[i+1:])
					b.rbs = b.rbs[:len(b.rbs)-1]
					break
				***REMOVED***
			***REMOVED***
		default:
			grpclog.Println("Unknown update.Op ", update.Op)
		***REMOVED***
	***REMOVED***
	// TODO: Fall back to the basic round-robin load balancing if the resulting address is
	// not a load balancer.
	select ***REMOVED***
	case <-ch:
	default:
	***REMOVED***
	ch <- b.rbs
	return nil
***REMOVED***

func (b *balancer) serverListExpire(seq int) ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()
	// TODO: gRPC interanls do not clear the connections when the server list is stale.
	// This means RPCs will keep using the existing server list until b receives new
	// server list even though the list is expired. Revisit this behavior later.
	if b.done || seq < b.seq ***REMOVED***
		return
	***REMOVED***
	b.next = 0
	b.addrs = nil
	// Ask grpc internals to close all the corresponding connections.
	b.addrCh <- nil
***REMOVED***

func convertDuration(d *lbpb.Duration) time.Duration ***REMOVED***
	if d == nil ***REMOVED***
		return 0
	***REMOVED***
	return time.Duration(d.Seconds)*time.Second + time.Duration(d.Nanos)*time.Nanosecond
***REMOVED***

func (b *balancer) processServerList(l *lbpb.ServerList, seq int) ***REMOVED***
	if l == nil ***REMOVED***
		return
	***REMOVED***
	servers := l.GetServers()
	expiration := convertDuration(l.GetExpirationInterval())
	var (
		sl    []*grpclbAddrInfo
		addrs []Address
	)
	for _, s := range servers ***REMOVED***
		md := metadata.Pairs("lb-token", s.LoadBalanceToken)
		addr := Address***REMOVED***
			Addr:     fmt.Sprintf("%s:%d", net.IP(s.IpAddress), s.Port),
			Metadata: &md,
		***REMOVED***
		sl = append(sl, &grpclbAddrInfo***REMOVED***
			addr:                 addr,
			dropForRateLimiting:  s.DropForRateLimiting,
			dropForLoadBalancing: s.DropForLoadBalancing,
		***REMOVED***)
		addrs = append(addrs, addr)
	***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.done || seq < b.seq ***REMOVED***
		return
	***REMOVED***
	if len(sl) > 0 ***REMOVED***
		// reset b.next to 0 when replacing the server list.
		b.next = 0
		b.addrs = sl
		b.addrCh <- addrs
		if b.expTimer != nil ***REMOVED***
			b.expTimer.Stop()
			b.expTimer = nil
		***REMOVED***
		if expiration > 0 ***REMOVED***
			b.expTimer = time.AfterFunc(expiration, func() ***REMOVED***
				b.serverListExpire(seq)
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (b *balancer) sendLoadReport(s *balanceLoadClientStream, interval time.Duration, done <-chan struct***REMOVED******REMOVED***) ***REMOVED***
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
		case <-done:
			return
		***REMOVED***
		b.mu.Lock()
		stats := b.clientStats
		b.clientStats = lbpb.ClientStats***REMOVED******REMOVED*** // Clear the stats.
		b.mu.Unlock()
		t := time.Now()
		stats.Timestamp = &lbpb.Timestamp***REMOVED***
			Seconds: t.Unix(),
			Nanos:   int32(t.Nanosecond()),
		***REMOVED***
		if err := s.Send(&lbpb.LoadBalanceRequest***REMOVED***
			LoadBalanceRequestType: &lbpb.LoadBalanceRequest_ClientStats***REMOVED***
				ClientStats: &stats,
			***REMOVED***,
		***REMOVED***); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (b *balancer) callRemoteBalancer(lbc *loadBalancerClient, seq int) (retry bool) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := lbc.BalanceLoad(ctx)
	if err != nil ***REMOVED***
		grpclog.Printf("Failed to perform RPC to the remote balancer %v", err)
		return
	***REMOVED***
	b.mu.Lock()
	if b.done ***REMOVED***
		b.mu.Unlock()
		return
	***REMOVED***
	b.mu.Unlock()
	initReq := &lbpb.LoadBalanceRequest***REMOVED***
		LoadBalanceRequestType: &lbpb.LoadBalanceRequest_InitialRequest***REMOVED***
			InitialRequest: &lbpb.InitialLoadBalanceRequest***REMOVED***
				Name: b.target,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	if err := stream.Send(initReq); err != nil ***REMOVED***
		// TODO: backoff on retry?
		return true
	***REMOVED***
	reply, err := stream.Recv()
	if err != nil ***REMOVED***
		// TODO: backoff on retry?
		return true
	***REMOVED***
	initResp := reply.GetInitialResponse()
	if initResp == nil ***REMOVED***
		grpclog.Println("Failed to receive the initial response from the remote balancer.")
		return
	***REMOVED***
	// TODO: Support delegation.
	if initResp.LoadBalancerDelegate != "" ***REMOVED***
		// delegation
		grpclog.Println("TODO: Delegation is not supported yet.")
		return
	***REMOVED***
	streamDone := make(chan struct***REMOVED******REMOVED***)
	defer close(streamDone)
	b.mu.Lock()
	b.clientStats = lbpb.ClientStats***REMOVED******REMOVED*** // Clear client stats.
	b.mu.Unlock()
	if d := convertDuration(initResp.ClientStatsReportInterval); d > 0 ***REMOVED***
		go b.sendLoadReport(stream, d, streamDone)
	***REMOVED***
	// Retrieve the server list.
	for ***REMOVED***
		reply, err := stream.Recv()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		b.mu.Lock()
		if b.done || seq < b.seq ***REMOVED***
			b.mu.Unlock()
			return
		***REMOVED***
		b.seq++ // tick when receiving a new list of servers.
		seq = b.seq
		b.mu.Unlock()
		if serverList := reply.GetServerList(); serverList != nil ***REMOVED***
			b.processServerList(serverList, seq)
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (b *balancer) Start(target string, config BalancerConfig) error ***REMOVED***
	b.rand = rand.New(rand.NewSource(time.Now().Unix()))
	// TODO: Fall back to the basic direct connection if there is no name resolver.
	if b.r == nil ***REMOVED***
		return errors.New("there is no name resolver installed")
	***REMOVED***
	b.target = target
	b.mu.Lock()
	if b.done ***REMOVED***
		b.mu.Unlock()
		return ErrClientConnClosing
	***REMOVED***
	b.addrCh = make(chan []Address)
	w, err := b.r.Resolve(target)
	if err != nil ***REMOVED***
		b.mu.Unlock()
		return err
	***REMOVED***
	b.w = w
	b.mu.Unlock()
	balancerAddrsCh := make(chan []remoteBalancerInfo, 1)
	// Spawn a goroutine to monitor the name resolution of remote load balancer.
	go func() ***REMOVED***
		for ***REMOVED***
			if err := b.watchAddrUpdates(w, balancerAddrsCh); err != nil ***REMOVED***
				grpclog.Printf("grpc: the naming watcher stops working due to %v.\n", err)
				close(balancerAddrsCh)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	// Spawn a goroutine to talk to the remote load balancer.
	go func() ***REMOVED***
		var (
			cc *ClientConn
			// ccError is closed when there is an error in the current cc.
			// A new rb should be picked from rbs and connected.
			ccError chan struct***REMOVED******REMOVED***
			rb      *remoteBalancerInfo
			rbs     []remoteBalancerInfo
			rbIdx   int
		)

		defer func() ***REMOVED***
			if ccError != nil ***REMOVED***
				select ***REMOVED***
				case <-ccError:
				default:
					close(ccError)
				***REMOVED***
			***REMOVED***
			if cc != nil ***REMOVED***
				cc.Close()
			***REMOVED***
		***REMOVED***()

		for ***REMOVED***
			var ok bool
			select ***REMOVED***
			case rbs, ok = <-balancerAddrsCh:
				if !ok ***REMOVED***
					return
				***REMOVED***
				foundIdx := -1
				if rb != nil ***REMOVED***
					for i, trb := range rbs ***REMOVED***
						if trb == *rb ***REMOVED***
							foundIdx = i
							break
						***REMOVED***
					***REMOVED***
				***REMOVED***
				if foundIdx >= 0 ***REMOVED***
					if foundIdx >= 1 ***REMOVED***
						// Move the address in use to the beginning of the list.
						b.rbs[0], b.rbs[foundIdx] = b.rbs[foundIdx], b.rbs[0]
						rbIdx = 0
					***REMOVED***
					continue // If found, don't dial new cc.
				***REMOVED*** else if len(rbs) > 0 ***REMOVED***
					// Pick a random one from the list, instead of always using the first one.
					if l := len(rbs); l > 1 && rb != nil ***REMOVED***
						tmpIdx := b.rand.Intn(l - 1)
						b.rbs[0], b.rbs[tmpIdx] = b.rbs[tmpIdx], b.rbs[0]
					***REMOVED***
					rbIdx = 0
					rb = &rbs[0]
				***REMOVED*** else ***REMOVED***
					// foundIdx < 0 && len(rbs) <= 0.
					rb = nil
				***REMOVED***
			case <-ccError:
				ccError = nil
				if rbIdx < len(rbs)-1 ***REMOVED***
					rbIdx++
					rb = &rbs[rbIdx]
				***REMOVED*** else ***REMOVED***
					rb = nil
				***REMOVED***
			***REMOVED***

			if rb == nil ***REMOVED***
				continue
			***REMOVED***

			if cc != nil ***REMOVED***
				cc.Close()
			***REMOVED***
			// Talk to the remote load balancer to get the server list.
			var err error
			creds := config.DialCreds
			ccError = make(chan struct***REMOVED******REMOVED***)
			if creds == nil ***REMOVED***
				cc, err = Dial(rb.addr, WithInsecure())
			***REMOVED*** else ***REMOVED***
				if rb.name != "" ***REMOVED***
					if err := creds.OverrideServerName(rb.name); err != nil ***REMOVED***
						grpclog.Printf("Failed to override the server name in the credentials: %v", err)
						continue
					***REMOVED***
				***REMOVED***
				cc, err = Dial(rb.addr, WithTransportCredentials(creds))
			***REMOVED***
			if err != nil ***REMOVED***
				grpclog.Printf("Failed to setup a connection to the remote balancer %v: %v", rb.addr, err)
				close(ccError)
				continue
			***REMOVED***
			b.mu.Lock()
			b.seq++ // tick when getting a new balancer address
			seq := b.seq
			b.next = 0
			b.mu.Unlock()
			go func(cc *ClientConn, ccError chan struct***REMOVED******REMOVED***) ***REMOVED***
				lbc := &loadBalancerClient***REMOVED***cc***REMOVED***
				b.callRemoteBalancer(lbc, seq)
				cc.Close()
				select ***REMOVED***
				case <-ccError:
				default:
					close(ccError)
				***REMOVED***
			***REMOVED***(cc, ccError)
		***REMOVED***
	***REMOVED***()
	return nil
***REMOVED***

func (b *balancer) down(addr Address, err error) ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, a := range b.addrs ***REMOVED***
		if addr == a.addr ***REMOVED***
			a.connected = false
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func (b *balancer) Up(addr Address) func(error) ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.done ***REMOVED***
		return nil
	***REMOVED***
	var cnt int
	for _, a := range b.addrs ***REMOVED***
		if a.addr == addr ***REMOVED***
			if a.connected ***REMOVED***
				return nil
			***REMOVED***
			a.connected = true
		***REMOVED***
		if a.connected && !a.dropForRateLimiting && !a.dropForLoadBalancing ***REMOVED***
			cnt++
		***REMOVED***
	***REMOVED***
	// addr is the only one which is connected. Notify the Get() callers who are blocking.
	if cnt == 1 && b.waitCh != nil ***REMOVED***
		close(b.waitCh)
		b.waitCh = nil
	***REMOVED***
	return func(err error) ***REMOVED***
		b.down(addr, err)
	***REMOVED***
***REMOVED***

func (b *balancer) Get(ctx context.Context, opts BalancerGetOptions) (addr Address, put func(), err error) ***REMOVED***
	var ch chan struct***REMOVED******REMOVED***
	b.mu.Lock()
	if b.done ***REMOVED***
		b.mu.Unlock()
		err = ErrClientConnClosing
		return
	***REMOVED***
	seq := b.seq

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			return
		***REMOVED***
		put = func() ***REMOVED***
			s, ok := rpcInfoFromContext(ctx)
			if !ok ***REMOVED***
				return
			***REMOVED***
			b.mu.Lock()
			defer b.mu.Unlock()
			if b.done || seq < b.seq ***REMOVED***
				return
			***REMOVED***
			b.clientStats.NumCallsFinished++
			if !s.bytesSent ***REMOVED***
				b.clientStats.NumCallsFinishedWithClientFailedToSend++
			***REMOVED*** else if s.bytesReceived ***REMOVED***
				b.clientStats.NumCallsFinishedKnownReceived++
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	b.clientStats.NumCallsStarted++
	if len(b.addrs) > 0 ***REMOVED***
		if b.next >= len(b.addrs) ***REMOVED***
			b.next = 0
		***REMOVED***
		next := b.next
		for ***REMOVED***
			a := b.addrs[next]
			next = (next + 1) % len(b.addrs)
			if a.connected ***REMOVED***
				if !a.dropForRateLimiting && !a.dropForLoadBalancing ***REMOVED***
					addr = a.addr
					b.next = next
					b.mu.Unlock()
					return
				***REMOVED***
				if !opts.BlockingWait ***REMOVED***
					b.next = next
					if a.dropForLoadBalancing ***REMOVED***
						b.clientStats.NumCallsFinished++
						b.clientStats.NumCallsFinishedWithDropForLoadBalancing++
					***REMOVED*** else if a.dropForRateLimiting ***REMOVED***
						b.clientStats.NumCallsFinished++
						b.clientStats.NumCallsFinishedWithDropForRateLimiting++
					***REMOVED***
					b.mu.Unlock()
					err = Errorf(codes.Unavailable, "%s drops requests", a.addr.Addr)
					return
				***REMOVED***
			***REMOVED***
			if next == b.next ***REMOVED***
				// Has iterated all the possible address but none is connected.
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if !opts.BlockingWait ***REMOVED***
		if len(b.addrs) == 0 ***REMOVED***
			b.clientStats.NumCallsFinished++
			b.clientStats.NumCallsFinishedWithClientFailedToSend++
			b.mu.Unlock()
			err = Errorf(codes.Unavailable, "there is no address available")
			return
		***REMOVED***
		// Returns the next addr on b.addrs for a failfast RPC.
		addr = b.addrs[b.next].addr
		b.next++
		b.mu.Unlock()
		return
	***REMOVED***
	// Wait on b.waitCh for non-failfast RPCs.
	if b.waitCh == nil ***REMOVED***
		ch = make(chan struct***REMOVED******REMOVED***)
		b.waitCh = ch
	***REMOVED*** else ***REMOVED***
		ch = b.waitCh
	***REMOVED***
	b.mu.Unlock()
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			b.mu.Lock()
			b.clientStats.NumCallsFinished++
			b.clientStats.NumCallsFinishedWithClientFailedToSend++
			b.mu.Unlock()
			err = ctx.Err()
			return
		case <-ch:
			b.mu.Lock()
			if b.done ***REMOVED***
				b.clientStats.NumCallsFinished++
				b.clientStats.NumCallsFinishedWithClientFailedToSend++
				b.mu.Unlock()
				err = ErrClientConnClosing
				return
			***REMOVED***

			if len(b.addrs) > 0 ***REMOVED***
				if b.next >= len(b.addrs) ***REMOVED***
					b.next = 0
				***REMOVED***
				next := b.next
				for ***REMOVED***
					a := b.addrs[next]
					next = (next + 1) % len(b.addrs)
					if a.connected ***REMOVED***
						if !a.dropForRateLimiting && !a.dropForLoadBalancing ***REMOVED***
							addr = a.addr
							b.next = next
							b.mu.Unlock()
							return
						***REMOVED***
						if !opts.BlockingWait ***REMOVED***
							b.next = next
							if a.dropForLoadBalancing ***REMOVED***
								b.clientStats.NumCallsFinished++
								b.clientStats.NumCallsFinishedWithDropForLoadBalancing++
							***REMOVED*** else if a.dropForRateLimiting ***REMOVED***
								b.clientStats.NumCallsFinished++
								b.clientStats.NumCallsFinishedWithDropForRateLimiting++
							***REMOVED***
							b.mu.Unlock()
							err = Errorf(codes.Unavailable, "drop requests for the addreess %s", a.addr.Addr)
							return
						***REMOVED***
					***REMOVED***
					if next == b.next ***REMOVED***
						// Has iterated all the possible address but none is connected.
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***
			// The newly added addr got removed by Down() again.
			if b.waitCh == nil ***REMOVED***
				ch = make(chan struct***REMOVED******REMOVED***)
				b.waitCh = ch
			***REMOVED*** else ***REMOVED***
				ch = b.waitCh
			***REMOVED***
			b.mu.Unlock()
		***REMOVED***
	***REMOVED***
***REMOVED***

func (b *balancer) Notify() <-chan []Address ***REMOVED***
	return b.addrCh
***REMOVED***

func (b *balancer) Close() error ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()
	b.done = true
	if b.expTimer != nil ***REMOVED***
		b.expTimer.Stop()
	***REMOVED***
	if b.waitCh != nil ***REMOVED***
		close(b.waitCh)
	***REMOVED***
	if b.addrCh != nil ***REMOVED***
		close(b.addrCh)
	***REMOVED***
	if b.w != nil ***REMOVED***
		b.w.Close()
	***REMOVED***
	return nil
***REMOVED***
