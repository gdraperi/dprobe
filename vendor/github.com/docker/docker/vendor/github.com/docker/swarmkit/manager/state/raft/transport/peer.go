package transport

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state/raft/membership"
	"github.com/pkg/errors"
)

const (
	// GRPCMaxMsgSize is the max allowed gRPC message size for raft messages.
	GRPCMaxMsgSize = 4 << 20
)

type peer struct ***REMOVED***
	id uint64

	tr *Transport

	msgc chan raftpb.Message

	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct***REMOVED******REMOVED***

	mu      sync.Mutex
	cc      *grpc.ClientConn
	addr    string
	newAddr string

	active       bool
	becameActive time.Time
***REMOVED***

func newPeer(id uint64, addr string, tr *Transport) (*peer, error) ***REMOVED***
	cc, err := tr.dial(addr)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to create conn for %x with addr %s", id, addr)
	***REMOVED***
	ctx, cancel := context.WithCancel(tr.ctx)
	ctx = log.WithField(ctx, "peer_id", fmt.Sprintf("%x", id))
	p := &peer***REMOVED***
		id:     id,
		addr:   addr,
		cc:     cc,
		tr:     tr,
		ctx:    ctx,
		cancel: cancel,
		msgc:   make(chan raftpb.Message, 4096),
		done:   make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
	go p.run(ctx)
	return p, nil
***REMOVED***

func (p *peer) send(m raftpb.Message) (err error) ***REMOVED***
	p.mu.Lock()
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			p.active = false
			p.becameActive = time.Time***REMOVED******REMOVED***
		***REMOVED***
		p.mu.Unlock()
	***REMOVED***()
	select ***REMOVED***
	case <-p.ctx.Done():
		return p.ctx.Err()
	default:
	***REMOVED***
	select ***REMOVED***
	case p.msgc <- m:
	case <-p.ctx.Done():
		return p.ctx.Err()
	default:
		p.tr.config.ReportUnreachable(p.id)
		return errors.Errorf("peer is unreachable")
	***REMOVED***
	return nil
***REMOVED***

func (p *peer) update(addr string) error ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.addr == addr ***REMOVED***
		return nil
	***REMOVED***
	cc, err := p.tr.dial(addr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	p.cc.Close()
	p.cc = cc
	p.addr = addr
	return nil
***REMOVED***

func (p *peer) updateAddr(addr string) error ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.addr == addr ***REMOVED***
		return nil
	***REMOVED***
	log.G(p.ctx).Debugf("peer %x updated to address %s, it will be used if old failed", p.id, addr)
	p.newAddr = addr
	return nil
***REMOVED***

func (p *peer) conn() *grpc.ClientConn ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.cc
***REMOVED***

func (p *peer) address() string ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.addr
***REMOVED***

func (p *peer) resolveAddr(ctx context.Context, id uint64) (string, error) ***REMOVED***
	resp, err := api.NewRaftClient(p.conn()).ResolveAddress(ctx, &api.ResolveAddressRequest***REMOVED***RaftID: id***REMOVED***)
	if err != nil ***REMOVED***
		return "", errors.Wrap(err, "failed to resolve address")
	***REMOVED***
	return resp.Addr, nil
***REMOVED***

// Returns the raft message struct size (not including the payload size) for the given raftpb.Message.
// The payload is typically the snapshot or append entries.
func raftMessageStructSize(m *raftpb.Message) int ***REMOVED***
	return (&api.ProcessRaftMessageRequest***REMOVED***Message: m***REMOVED***).Size() - len(m.Snapshot.Data)
***REMOVED***

// Returns the max allowable payload based on MaxRaftMsgSize and
// the struct size for the given raftpb.Message.
func raftMessagePayloadSize(m *raftpb.Message) int ***REMOVED***
	return GRPCMaxMsgSize - raftMessageStructSize(m)
***REMOVED***

// Split a large raft message into smaller messages.
// Currently this means splitting the []Snapshot.Data into chunks whose size
// is dictacted by MaxRaftMsgSize.
func splitSnapshotData(ctx context.Context, m *raftpb.Message) []api.StreamRaftMessageRequest ***REMOVED***
	var messages []api.StreamRaftMessageRequest
	if m.Type != raftpb.MsgSnap ***REMOVED***
		return messages
	***REMOVED***

	// get the size of the data to be split.
	size := len(m.Snapshot.Data)

	// Get the max payload size.
	payloadSize := raftMessagePayloadSize(m)

	// split the snapshot into smaller messages.
	for snapDataIndex := 0; snapDataIndex < size; ***REMOVED***
		chunkSize := size - snapDataIndex
		if chunkSize > payloadSize ***REMOVED***
			chunkSize = payloadSize
		***REMOVED***

		raftMsg := *m

		// sub-slice for this snapshot chunk.
		raftMsg.Snapshot.Data = m.Snapshot.Data[snapDataIndex : snapDataIndex+chunkSize]

		snapDataIndex += chunkSize

		// add message to the list of messages to be sent.
		msg := api.StreamRaftMessageRequest***REMOVED***Message: &raftMsg***REMOVED***
		messages = append(messages, msg)
	***REMOVED***

	return messages
***REMOVED***

// Function to check if this message needs to be split to be streamed
// (because it is larger than GRPCMaxMsgSize).
// Returns true if the message type is MsgSnap
// and size larger than MaxRaftMsgSize.
func needsSplitting(m *raftpb.Message) bool ***REMOVED***
	raftMsg := api.ProcessRaftMessageRequest***REMOVED***Message: m***REMOVED***
	return m.Type == raftpb.MsgSnap && raftMsg.Size() > GRPCMaxMsgSize
***REMOVED***

func (p *peer) sendProcessMessage(ctx context.Context, m raftpb.Message) error ***REMOVED***
	ctx, cancel := context.WithTimeout(ctx, p.tr.config.SendTimeout)
	defer cancel()

	var err error
	var stream api.Raft_StreamRaftMessageClient
	stream, err = api.NewRaftClient(p.conn()).StreamRaftMessage(ctx)

	if err == nil ***REMOVED***
		// Split the message if needed.
		// Currently only supported for MsgSnap.
		var msgs []api.StreamRaftMessageRequest
		if needsSplitting(&m) ***REMOVED***
			msgs = splitSnapshotData(ctx, &m)
		***REMOVED*** else ***REMOVED***
			raftMsg := api.StreamRaftMessageRequest***REMOVED***Message: &m***REMOVED***
			msgs = append(msgs, raftMsg)
		***REMOVED***

		// Stream
		for _, msg := range msgs ***REMOVED***
			err = stream.Send(&msg)
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("error streaming message to peer")
				stream.CloseAndRecv()
				break
			***REMOVED***
		***REMOVED***

		// Finished sending all the messages.
		// Close and receive response.
		if err == nil ***REMOVED***
			_, err = stream.CloseAndRecv()

			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("error receiving response")
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		log.G(ctx).WithError(err).Error("error sending message to peer")
	***REMOVED***

	// Try doing a regular rpc if the receiver doesn't support streaming.
	if grpc.Code(err) == codes.Unimplemented ***REMOVED***
		_, err = api.NewRaftClient(p.conn()).ProcessRaftMessage(ctx, &api.ProcessRaftMessageRequest***REMOVED***Message: &m***REMOVED***)
	***REMOVED***

	// Handle errors.
	if grpc.Code(err) == codes.NotFound && grpc.ErrorDesc(err) == membership.ErrMemberRemoved.Error() ***REMOVED***
		p.tr.config.NodeRemoved()
	***REMOVED***
	if m.Type == raftpb.MsgSnap ***REMOVED***
		if err != nil ***REMOVED***
			p.tr.config.ReportSnapshot(m.To, raft.SnapshotFailure)
		***REMOVED*** else ***REMOVED***
			p.tr.config.ReportSnapshot(m.To, raft.SnapshotFinish)
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		p.tr.config.ReportUnreachable(m.To)
		return err
	***REMOVED***
	return nil
***REMOVED***

func healthCheckConn(ctx context.Context, cc *grpc.ClientConn) error ***REMOVED***
	resp, err := api.NewHealthClient(cc).Check(ctx, &api.HealthCheckRequest***REMOVED***Service: "Raft"***REMOVED***)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to check health")
	***REMOVED***
	if resp.Status != api.HealthCheckResponse_SERVING ***REMOVED***
		return errors.Errorf("health check returned status %s", resp.Status)
	***REMOVED***
	return nil
***REMOVED***

func (p *peer) healthCheck(ctx context.Context) error ***REMOVED***
	ctx, cancel := context.WithTimeout(ctx, p.tr.config.SendTimeout)
	defer cancel()
	return healthCheckConn(ctx, p.conn())
***REMOVED***

func (p *peer) setActive() ***REMOVED***
	p.mu.Lock()
	if !p.active ***REMOVED***
		p.active = true
		p.becameActive = time.Now()
	***REMOVED***
	p.mu.Unlock()
***REMOVED***

func (p *peer) setInactive() ***REMOVED***
	p.mu.Lock()
	p.active = false
	p.becameActive = time.Time***REMOVED******REMOVED***
	p.mu.Unlock()
***REMOVED***

func (p *peer) activeTime() time.Time ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.becameActive
***REMOVED***

func (p *peer) drain() error ***REMOVED***
	ctx, cancel := context.WithTimeout(context.Background(), 16*time.Second)
	defer cancel()
	for ***REMOVED***
		select ***REMOVED***
		case m, ok := <-p.msgc:
			if !ok ***REMOVED***
				// all messages proceeded
				return nil
			***REMOVED***
			if err := p.sendProcessMessage(ctx, m); err != nil ***REMOVED***
				return errors.Wrap(err, "send drain message")
			***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *peer) handleAddressChange(ctx context.Context) error ***REMOVED***
	p.mu.Lock()
	newAddr := p.newAddr
	p.newAddr = ""
	p.mu.Unlock()
	if newAddr == "" ***REMOVED***
		return nil
	***REMOVED***
	cc, err := p.tr.dial(newAddr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	ctx, cancel := context.WithTimeout(ctx, p.tr.config.SendTimeout)
	defer cancel()
	if err := healthCheckConn(ctx, cc); err != nil ***REMOVED***
		cc.Close()
		return err
	***REMOVED***
	// there is possibility of race if host changing address too fast, but
	// it's unlikely and eventually thing should be settled
	p.mu.Lock()
	p.cc.Close()
	p.cc = cc
	p.addr = newAddr
	p.tr.config.UpdateNode(p.id, p.addr)
	p.mu.Unlock()
	return nil
***REMOVED***

func (p *peer) run(ctx context.Context) ***REMOVED***
	defer func() ***REMOVED***
		p.mu.Lock()
		p.active = false
		p.becameActive = time.Time***REMOVED******REMOVED***
		// at this point we can be sure that nobody will write to msgc
		if p.msgc != nil ***REMOVED***
			close(p.msgc)
		***REMOVED***
		p.mu.Unlock()
		if err := p.drain(); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("failed to drain message queue")
		***REMOVED***
		close(p.done)
	***REMOVED***()
	if err := p.healthCheck(ctx); err == nil ***REMOVED***
		p.setActive()
	***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
		***REMOVED***

		select ***REMOVED***
		case m := <-p.msgc:
			// we do not propagate context here, because this operation should be finished
			// or timed out for correct raft work.
			err := p.sendProcessMessage(context.Background(), m)
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Debugf("failed to send message %s", m.Type)
				p.setInactive()
				if err := p.handleAddressChange(ctx); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("failed to change address after failure")
				***REMOVED***
				continue
			***REMOVED***
			p.setActive()
		case <-ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *peer) stop() ***REMOVED***
	p.cancel()
	<-p.done
***REMOVED***
