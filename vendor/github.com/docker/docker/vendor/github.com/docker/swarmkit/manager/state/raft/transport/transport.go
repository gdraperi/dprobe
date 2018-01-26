// Package transport provides grpc transport layer for raft.
// All methods are non-blocking.
package transport

import (
	"net"
	"sync"
	"time"

	"golang.org/x/net/context"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/docker/swarmkit/log"
	"github.com/pkg/errors"
)

// ErrIsNotFound indicates that peer was never added to transport.
var ErrIsNotFound = errors.New("peer not found")

// Raft is interface which represents Raft API for transport package.
type Raft interface ***REMOVED***
	ReportUnreachable(id uint64)
	ReportSnapshot(id uint64, status raft.SnapshotStatus)
	IsIDRemoved(id uint64) bool
	UpdateNode(id uint64, addr string)

	NodeRemoved()
***REMOVED***

// Config for Transport
type Config struct ***REMOVED***
	HeartbeatInterval time.Duration
	SendTimeout       time.Duration
	Credentials       credentials.TransportCredentials
	RaftID            string

	Raft
***REMOVED***

// Transport is structure which manages remote raft peers and sends messages
// to them.
type Transport struct ***REMOVED***
	config *Config

	unknownc chan raftpb.Message

	mu      sync.Mutex
	peers   map[uint64]*peer
	stopped bool

	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct***REMOVED******REMOVED***

	deferredConns map[*grpc.ClientConn]*time.Timer
***REMOVED***

// New returns new Transport with specified Config.
func New(cfg *Config) *Transport ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	if cfg.RaftID != "" ***REMOVED***
		ctx = log.WithField(ctx, "raft_id", cfg.RaftID)
	***REMOVED***
	t := &Transport***REMOVED***
		peers:    make(map[uint64]*peer),
		config:   cfg,
		unknownc: make(chan raftpb.Message),
		done:     make(chan struct***REMOVED******REMOVED***),
		ctx:      ctx,
		cancel:   cancel,

		deferredConns: make(map[*grpc.ClientConn]*time.Timer),
	***REMOVED***
	go t.run(ctx)
	return t
***REMOVED***

func (t *Transport) run(ctx context.Context) ***REMOVED***
	defer func() ***REMOVED***
		log.G(ctx).Debug("stop transport")
		t.mu.Lock()
		defer t.mu.Unlock()
		t.stopped = true
		for _, p := range t.peers ***REMOVED***
			p.stop()
			p.cc.Close()
		***REMOVED***
		for cc, timer := range t.deferredConns ***REMOVED***
			timer.Stop()
			cc.Close()
		***REMOVED***
		t.deferredConns = nil
		close(t.done)
	***REMOVED***()
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
		***REMOVED***

		select ***REMOVED***
		case m := <-t.unknownc:
			if err := t.sendUnknownMessage(ctx, m); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Warnf("ignored message %s to unknown peer %x", m.Type, m.To)
			***REMOVED***
		case <-ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stop stops transport and waits until it finished
func (t *Transport) Stop() ***REMOVED***
	t.cancel()
	<-t.done
***REMOVED***

// Send sends raft message to remote peers.
func (t *Transport) Send(m raftpb.Message) error ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.stopped ***REMOVED***
		return errors.New("transport stopped")
	***REMOVED***
	if t.config.IsIDRemoved(m.To) ***REMOVED***
		return errors.Errorf("refusing to send message %s to removed member %x", m.Type, m.To)
	***REMOVED***
	p, ok := t.peers[m.To]
	if !ok ***REMOVED***
		log.G(t.ctx).Warningf("sending message %s to an unrecognized member ID %x", m.Type, m.To)
		select ***REMOVED***
		// we need to process messages to unknown peers in separate goroutine
		// to not block sender
		case t.unknownc <- m:
		case <-t.ctx.Done():
			return t.ctx.Err()
		default:
			return errors.New("unknown messages queue is full")
		***REMOVED***
		return nil
	***REMOVED***
	if err := p.send(m); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to send message %x to %x", m.Type, m.To)
	***REMOVED***
	return nil
***REMOVED***

// AddPeer adds new peer with id and address addr to Transport.
// If there is already peer with such id in Transport it will return error if
// address is different (UpdatePeer should be used) or nil otherwise.
func (t *Transport) AddPeer(id uint64, addr string) error ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.stopped ***REMOVED***
		return errors.New("transport stopped")
	***REMOVED***
	if ep, ok := t.peers[id]; ok ***REMOVED***
		if ep.address() == addr ***REMOVED***
			return nil
		***REMOVED***
		return errors.Errorf("peer %x already added with addr %s", id, ep.addr)
	***REMOVED***
	log.G(t.ctx).Debugf("transport: add peer %x with address %s", id, addr)
	p, err := newPeer(id, addr, t)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to create peer %x with addr %s", id, addr)
	***REMOVED***
	t.peers[id] = p
	return nil
***REMOVED***

// RemovePeer removes peer from Transport and wait for it to stop.
func (t *Transport) RemovePeer(id uint64) error ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped ***REMOVED***
		return errors.New("transport stopped")
	***REMOVED***
	p, ok := t.peers[id]
	if !ok ***REMOVED***
		return ErrIsNotFound
	***REMOVED***
	delete(t.peers, id)
	cc := p.conn()
	p.stop()
	timer := time.AfterFunc(8*time.Second, func() ***REMOVED***
		t.mu.Lock()
		if !t.stopped ***REMOVED***
			delete(t.deferredConns, cc)
			cc.Close()
		***REMOVED***
		t.mu.Unlock()
	***REMOVED***)
	// store connection and timer for cleaning up on stop
	t.deferredConns[cc] = timer

	return nil
***REMOVED***

// UpdatePeer updates peer with new address. It replaces connection immediately.
func (t *Transport) UpdatePeer(id uint64, addr string) error ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped ***REMOVED***
		return errors.New("transport stopped")
	***REMOVED***
	p, ok := t.peers[id]
	if !ok ***REMOVED***
		return ErrIsNotFound
	***REMOVED***
	if err := p.update(addr); err != nil ***REMOVED***
		return err
	***REMOVED***
	log.G(t.ctx).Debugf("peer %x updated to address %s", id, addr)
	return nil
***REMOVED***

// UpdatePeerAddr updates peer with new address, but delays connection creation.
// New address won't be used until first failure on old address.
func (t *Transport) UpdatePeerAddr(id uint64, addr string) error ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped ***REMOVED***
		return errors.New("transport stopped")
	***REMOVED***
	p, ok := t.peers[id]
	if !ok ***REMOVED***
		return ErrIsNotFound
	***REMOVED***
	return p.updateAddr(addr)
***REMOVED***

// PeerConn returns raw grpc connection to peer.
func (t *Transport) PeerConn(id uint64) (*grpc.ClientConn, error) ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	p, ok := t.peers[id]
	if !ok ***REMOVED***
		return nil, ErrIsNotFound
	***REMOVED***
	p.mu.Lock()
	active := p.active
	p.mu.Unlock()
	if !active ***REMOVED***
		return nil, errors.New("peer is inactive")
	***REMOVED***
	return p.conn(), nil
***REMOVED***

// PeerAddr returns address of peer with id.
func (t *Transport) PeerAddr(id uint64) (string, error) ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	p, ok := t.peers[id]
	if !ok ***REMOVED***
		return "", ErrIsNotFound
	***REMOVED***
	return p.address(), nil
***REMOVED***

// HealthCheck checks health of particular peer.
func (t *Transport) HealthCheck(ctx context.Context, id uint64) error ***REMOVED***
	t.mu.Lock()
	p, ok := t.peers[id]
	t.mu.Unlock()
	if !ok ***REMOVED***
		return ErrIsNotFound
	***REMOVED***
	ctx, cancel := t.withContext(ctx)
	defer cancel()
	return p.healthCheck(ctx)
***REMOVED***

// Active returns true if node was recently active and false otherwise.
func (t *Transport) Active(id uint64) bool ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	p, ok := t.peers[id]
	if !ok ***REMOVED***
		return false
	***REMOVED***
	p.mu.Lock()
	active := p.active
	p.mu.Unlock()
	return active
***REMOVED***

// LongestActive returns the ID of the peer that has been active for the longest
// length of time.
func (t *Transport) LongestActive() (uint64, error) ***REMOVED***
	p, err := t.longestActive()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return p.id, nil
***REMOVED***

// longestActive returns the peer that has been active for the longest length of
// time.
func (t *Transport) longestActive() (*peer, error) ***REMOVED***
	var longest *peer
	var longestTime time.Time
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, p := range t.peers ***REMOVED***
		becameActive := p.activeTime()
		if becameActive.IsZero() ***REMOVED***
			continue
		***REMOVED***
		if longest == nil ***REMOVED***
			longest = p
			continue
		***REMOVED***
		if becameActive.Before(longestTime) ***REMOVED***
			longest = p
			longestTime = becameActive
		***REMOVED***
	***REMOVED***
	if longest == nil ***REMOVED***
		return nil, errors.New("failed to find longest active peer")
	***REMOVED***
	return longest, nil
***REMOVED***

func (t *Transport) dial(addr string) (*grpc.ClientConn, error) ***REMOVED***
	grpcOptions := []grpc.DialOption***REMOVED***
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		grpc.WithBackoffMaxDelay(8 * time.Second),
	***REMOVED***
	if t.config.Credentials != nil ***REMOVED***
		grpcOptions = append(grpcOptions, grpc.WithTransportCredentials(t.config.Credentials))
	***REMOVED*** else ***REMOVED***
		grpcOptions = append(grpcOptions, grpc.WithInsecure())
	***REMOVED***

	if t.config.SendTimeout > 0 ***REMOVED***
		grpcOptions = append(grpcOptions, grpc.WithTimeout(t.config.SendTimeout))
	***REMOVED***

	// gRPC dialer connects to proxy first. Provide a custom dialer here avoid that.
	// TODO(anshul) Add an option to configure this.
	grpcOptions = append(grpcOptions,
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) ***REMOVED***
			return net.DialTimeout("tcp", addr, timeout)
		***REMOVED***))

	cc, err := grpc.Dial(addr, grpcOptions...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return cc, nil
***REMOVED***

func (t *Transport) withContext(ctx context.Context) (context.Context, context.CancelFunc) ***REMOVED***
	ctx, cancel := context.WithCancel(ctx)

	go func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
		case <-t.ctx.Done():
			cancel()
		***REMOVED***
	***REMOVED***()
	return ctx, cancel
***REMOVED***

func (t *Transport) resolvePeer(ctx context.Context, id uint64) (*peer, error) ***REMOVED***
	longestActive, err := t.longestActive()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ctx, cancel := context.WithTimeout(ctx, t.config.SendTimeout)
	defer cancel()
	addr, err := longestActive.resolveAddr(ctx, id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return newPeer(id, addr, t)
***REMOVED***

func (t *Transport) sendUnknownMessage(ctx context.Context, m raftpb.Message) error ***REMOVED***
	p, err := t.resolvePeer(ctx, m.To)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to resolve peer")
	***REMOVED***
	defer p.cancel()
	if err := p.sendProcessMessage(ctx, m); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to send message")
	***REMOVED***
	return nil
***REMOVED***
