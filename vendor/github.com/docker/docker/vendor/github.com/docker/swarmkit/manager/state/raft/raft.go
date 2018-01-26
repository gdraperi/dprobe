package raft

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/pkg/idutil"
	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/docker/docker/pkg/signal"
	"github.com/docker/go-events"
	"github.com/docker/go-metrics"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/raftselector"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/raft/membership"
	"github.com/docker/swarmkit/manager/state/raft/storage"
	"github.com/docker/swarmkit/manager/state/raft/transport"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/watch"
	"github.com/gogo/protobuf/proto"
	"github.com/pivotal-golang/clock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	// ErrNoRaftMember is thrown when the node is not yet part of a raft cluster
	ErrNoRaftMember = errors.New("raft: node is not yet part of a raft cluster")
	// ErrConfChangeRefused is returned when there is an issue with the configuration change
	ErrConfChangeRefused = errors.New("raft: propose configuration change refused")
	// ErrApplyNotSpecified is returned during the creation of a raft node when no apply method was provided
	ErrApplyNotSpecified = errors.New("raft: apply method was not specified")
	// ErrSetHardState is returned when the node fails to set the hard state
	ErrSetHardState = errors.New("raft: failed to set the hard state for log append entry")
	// ErrStopped is returned when an operation was submitted but the node was stopped in the meantime
	ErrStopped = errors.New("raft: failed to process the request: node is stopped")
	// ErrLostLeadership is returned when an operation was submitted but the node lost leader status before it became committed
	ErrLostLeadership = errors.New("raft: failed to process the request: node lost leader status")
	// ErrRequestTooLarge is returned when a raft internal message is too large to be sent
	ErrRequestTooLarge = errors.New("raft: raft message is too large and can't be sent")
	// ErrCannotRemoveMember is thrown when we try to remove a member from the cluster but this would result in a loss of quorum
	ErrCannotRemoveMember = errors.New("raft: member cannot be removed, because removing it may result in loss of quorum")
	// ErrNoClusterLeader is thrown when the cluster has no elected leader
	ErrNoClusterLeader = errors.New("raft: no elected cluster leader")
	// ErrMemberUnknown is sent in response to a message from an
	// unrecognized peer.
	ErrMemberUnknown = errors.New("raft: member unknown")

	// work around lint
	lostQuorumMessage = "The swarm does not have a leader. It's possible that too few managers are online. Make sure more than half of the managers are online."
	errLostQuorum     = errors.New(lostQuorumMessage)

	// Timer to capture ProposeValue() latency.
	proposeLatencyTimer metrics.Timer
)

// LeadershipState indicates whether the node is a leader or follower.
type LeadershipState int

const (
	// IsLeader indicates that the node is a raft leader.
	IsLeader LeadershipState = iota
	// IsFollower indicates that the node is a raft follower.
	IsFollower

	// lostQuorumTimeout is the number of ticks that can elapse with no
	// leader before LeaderConn starts returning an error right away.
	lostQuorumTimeout = 10
)

// EncryptionKeys are the current and, if necessary, pending DEKs with which to
// encrypt raft data
type EncryptionKeys struct ***REMOVED***
	CurrentDEK []byte
	PendingDEK []byte
***REMOVED***

// EncryptionKeyRotator is an interface to find out if any keys need rotating.
type EncryptionKeyRotator interface ***REMOVED***
	GetKeys() EncryptionKeys
	UpdateKeys(EncryptionKeys) error
	NeedsRotation() bool
	RotationNotify() chan struct***REMOVED******REMOVED***
***REMOVED***

// Node represents the Raft Node useful
// configuration.
type Node struct ***REMOVED***
	raftNode  raft.Node
	cluster   *membership.Cluster
	transport *transport.Transport

	raftStore           *raft.MemoryStorage
	memoryStore         *store.MemoryStore
	Config              *raft.Config
	opts                NodeOptions
	reqIDGen            *idutil.Generator
	wait                *wait
	campaignWhenAble    bool
	signalledLeadership uint32
	isMember            uint32
	bootstrapMembers    []*api.RaftMember

	// waitProp waits for all the proposals to be terminated before
	// shutting down the node.
	waitProp sync.WaitGroup

	confState       raftpb.ConfState
	appliedIndex    uint64
	snapshotMeta    raftpb.SnapshotMetadata
	writtenWALIndex uint64

	ticker clock.Ticker
	doneCh chan struct***REMOVED******REMOVED***
	// RemovedFromRaft notifies about node deletion from raft cluster
	RemovedFromRaft chan struct***REMOVED******REMOVED***
	cancelFunc      func()
	// removeRaftCh notifies about node deletion from raft cluster
	removeRaftCh        chan struct***REMOVED******REMOVED***
	removeRaftOnce      sync.Once
	leadershipBroadcast *watch.Queue

	// used to coordinate shutdown
	// Lock should be used only in stop(), all other functions should use RLock.
	stopMu sync.RWMutex
	// used for membership management checks
	membershipLock sync.Mutex
	// synchronizes access to n.opts.Addr, and makes sure the address is not
	// updated concurrently with JoinAndStart.
	addrLock sync.Mutex

	snapshotInProgress chan raftpb.SnapshotMetadata
	asyncTasks         sync.WaitGroup

	// stopped chan is used for notifying grpc handlers that raft node going
	// to stop.
	stopped chan struct***REMOVED******REMOVED***

	raftLogger          *storage.EncryptedRaftLogger
	keyRotator          EncryptionKeyRotator
	rotationQueued      bool
	clearData           bool
	waitForAppliedIndex uint64
	ticksWithNoLeader   uint32
***REMOVED***

// NodeOptions provides node-level options.
type NodeOptions struct ***REMOVED***
	// ID is the node's ID, from its certificate's CN field.
	ID string
	// Addr is the address of this node's listener
	Addr string
	// ForceNewCluster defines if we have to force a new cluster
	// because we are recovering from a backup data directory.
	ForceNewCluster bool
	// JoinAddr is the cluster to join. May be an empty string to create
	// a standalone cluster.
	JoinAddr string
	// ForceJoin tells us to join even if already part of a cluster.
	ForceJoin bool
	// Config is the raft config.
	Config *raft.Config
	// StateDir is the directory to store durable state.
	StateDir string
	// TickInterval interval is the time interval between raft ticks.
	TickInterval time.Duration
	// ClockSource is a Clock interface to use as a time base.
	// Leave this nil except for tests that are designed not to run in real
	// time.
	ClockSource clock.Clock
	// SendTimeout is the timeout on the sending messages to other raft
	// nodes. Leave this as 0 to get the default value.
	SendTimeout    time.Duration
	TLSCredentials credentials.TransportCredentials
	KeyRotator     EncryptionKeyRotator
	// DisableStackDump prevents Run from dumping goroutine stacks when the
	// store becomes stuck.
	DisableStackDump bool
***REMOVED***

func init() ***REMOVED***
	rand.Seed(time.Now().UnixNano())
	ns := metrics.NewNamespace("swarm", "raft", nil)
	proposeLatencyTimer = ns.NewTimer("transaction_latency", "Raft transaction latency.")
	metrics.Register(ns)
***REMOVED***

// NewNode generates a new Raft node
func NewNode(opts NodeOptions) *Node ***REMOVED***
	cfg := opts.Config
	if cfg == nil ***REMOVED***
		cfg = DefaultNodeConfig()
	***REMOVED***
	if opts.TickInterval == 0 ***REMOVED***
		opts.TickInterval = time.Second
	***REMOVED***
	if opts.SendTimeout == 0 ***REMOVED***
		opts.SendTimeout = 2 * time.Second
	***REMOVED***

	raftStore := raft.NewMemoryStorage()

	n := &Node***REMOVED***
		cluster:   membership.NewCluster(),
		raftStore: raftStore,
		opts:      opts,
		Config: &raft.Config***REMOVED***
			ElectionTick:    cfg.ElectionTick,
			HeartbeatTick:   cfg.HeartbeatTick,
			Storage:         raftStore,
			MaxSizePerMsg:   cfg.MaxSizePerMsg,
			MaxInflightMsgs: cfg.MaxInflightMsgs,
			Logger:          cfg.Logger,
			CheckQuorum:     cfg.CheckQuorum,
		***REMOVED***,
		doneCh:              make(chan struct***REMOVED******REMOVED***),
		RemovedFromRaft:     make(chan struct***REMOVED******REMOVED***),
		stopped:             make(chan struct***REMOVED******REMOVED***),
		leadershipBroadcast: watch.NewQueue(),
		keyRotator:          opts.KeyRotator,
	***REMOVED***
	n.memoryStore = store.NewMemoryStore(n)

	if opts.ClockSource == nil ***REMOVED***
		n.ticker = clock.NewClock().NewTicker(opts.TickInterval)
	***REMOVED*** else ***REMOVED***
		n.ticker = opts.ClockSource.NewTicker(opts.TickInterval)
	***REMOVED***

	n.reqIDGen = idutil.NewGenerator(uint16(n.Config.ID), time.Now())
	n.wait = newWait()

	n.cancelFunc = func(n *Node) func() ***REMOVED***
		var cancelOnce sync.Once
		return func() ***REMOVED***
			cancelOnce.Do(func() ***REMOVED***
				close(n.stopped)
			***REMOVED***)
		***REMOVED***
	***REMOVED***(n)

	return n
***REMOVED***

// IsIDRemoved reports if member with id was removed from cluster.
// Part of transport.Raft interface.
func (n *Node) IsIDRemoved(id uint64) bool ***REMOVED***
	return n.cluster.IsIDRemoved(id)
***REMOVED***

// NodeRemoved signals that node was removed from cluster and should stop.
// Part of transport.Raft interface.
func (n *Node) NodeRemoved() ***REMOVED***
	n.removeRaftOnce.Do(func() ***REMOVED***
		atomic.StoreUint32(&n.isMember, 0)
		close(n.RemovedFromRaft)
	***REMOVED***)
***REMOVED***

// ReportSnapshot reports snapshot status to underlying raft node.
// Part of transport.Raft interface.
func (n *Node) ReportSnapshot(id uint64, status raft.SnapshotStatus) ***REMOVED***
	n.raftNode.ReportSnapshot(id, status)
***REMOVED***

// ReportUnreachable reports to underlying raft node that member with id is
// unreachable.
// Part of transport.Raft interface.
func (n *Node) ReportUnreachable(id uint64) ***REMOVED***
	n.raftNode.ReportUnreachable(id)
***REMOVED***

// SetAddr provides the raft node's address. This can be used in cases where
// opts.Addr was not provided to NewNode, for example when a port was not bound
// until after the raft node was created.
func (n *Node) SetAddr(ctx context.Context, addr string) error ***REMOVED***
	n.addrLock.Lock()
	defer n.addrLock.Unlock()

	n.opts.Addr = addr

	if !n.IsMember() ***REMOVED***
		return nil
	***REMOVED***

	newRaftMember := &api.RaftMember***REMOVED***
		RaftID: n.Config.ID,
		NodeID: n.opts.ID,
		Addr:   addr,
	***REMOVED***
	if err := n.cluster.UpdateMember(n.Config.ID, newRaftMember); err != nil ***REMOVED***
		return err
	***REMOVED***

	// If the raft node is running, submit a configuration change
	// with the new address.

	// TODO(aaronl): Currently, this node must be the leader to
	// submit this configuration change. This works for the initial
	// use cases (single-node cluster late binding ports, or calling
	// SetAddr before joining a cluster). In the future, we may want
	// to support having a follower proactively change its remote
	// address.

	leadershipCh, cancelWatch := n.SubscribeLeadership()
	defer cancelWatch()

	ctx, cancelCtx := n.WithContext(ctx)
	defer cancelCtx()

	isLeader := atomic.LoadUint32(&n.signalledLeadership) == 1
	for !isLeader ***REMOVED***
		select ***REMOVED***
		case leadershipChange := <-leadershipCh:
			if leadershipChange == IsLeader ***REMOVED***
				isLeader = true
			***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		***REMOVED***
	***REMOVED***

	return n.updateNodeBlocking(ctx, n.Config.ID, addr)
***REMOVED***

// WithContext returns context which is cancelled when parent context cancelled
// or node is stopped.
func (n *Node) WithContext(ctx context.Context) (context.Context, context.CancelFunc) ***REMOVED***
	ctx, cancel := context.WithCancel(ctx)

	go func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
		case <-n.stopped:
			cancel()
		***REMOVED***
	***REMOVED***()
	return ctx, cancel
***REMOVED***

func (n *Node) initTransport() ***REMOVED***
	transportConfig := &transport.Config***REMOVED***
		HeartbeatInterval: time.Duration(n.Config.ElectionTick) * n.opts.TickInterval,
		SendTimeout:       n.opts.SendTimeout,
		Credentials:       n.opts.TLSCredentials,
		Raft:              n,
	***REMOVED***
	n.transport = transport.New(transportConfig)
***REMOVED***

// JoinAndStart joins and starts the raft server
func (n *Node) JoinAndStart(ctx context.Context) (err error) ***REMOVED***
	ctx, cancel := n.WithContext(ctx)
	defer func() ***REMOVED***
		cancel()
		if err != nil ***REMOVED***
			n.stopMu.Lock()
			// to shutdown transport
			n.cancelFunc()
			n.stopMu.Unlock()
			n.done()
		***REMOVED*** else ***REMOVED***
			atomic.StoreUint32(&n.isMember, 1)
		***REMOVED***
	***REMOVED***()

	loadAndStartErr := n.loadAndStart(ctx, n.opts.ForceNewCluster)
	if loadAndStartErr != nil && loadAndStartErr != storage.ErrNoWAL ***REMOVED***
		return loadAndStartErr
	***REMOVED***

	snapshot, err := n.raftStore.Snapshot()
	// Snapshot never returns an error
	if err != nil ***REMOVED***
		panic("could not get snapshot of raft store")
	***REMOVED***

	n.confState = snapshot.Metadata.ConfState
	n.appliedIndex = snapshot.Metadata.Index
	n.snapshotMeta = snapshot.Metadata
	n.writtenWALIndex, _ = n.raftStore.LastIndex() // lastIndex always returns nil as an error

	n.addrLock.Lock()
	defer n.addrLock.Unlock()

	// override the module field entirely, since etcd/raft is not exactly a submodule
	n.Config.Logger = log.G(ctx).WithField("module", "raft")

	// restore from snapshot
	if loadAndStartErr == nil ***REMOVED***
		if n.opts.JoinAddr != "" && n.opts.ForceJoin ***REMOVED***
			if err := n.joinCluster(ctx); err != nil ***REMOVED***
				return errors.Wrap(err, "failed to rejoin cluster")
			***REMOVED***
		***REMOVED***
		n.campaignWhenAble = true
		n.initTransport()
		n.raftNode = raft.RestartNode(n.Config)
		return nil
	***REMOVED***

	if n.opts.JoinAddr == "" ***REMOVED***
		// First member in the cluster, self-assign ID
		n.Config.ID = uint64(rand.Int63()) + 1
		peer, err := n.newRaftLogs(n.opts.ID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		n.campaignWhenAble = true
		n.initTransport()
		n.raftNode = raft.StartNode(n.Config, []raft.Peer***REMOVED***peer***REMOVED***)
		return nil
	***REMOVED***

	// join to existing cluster

	if err := n.joinCluster(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***

	if _, err := n.newRaftLogs(n.opts.ID); err != nil ***REMOVED***
		return err
	***REMOVED***

	n.initTransport()
	n.raftNode = raft.StartNode(n.Config, nil)

	return nil
***REMOVED***

func (n *Node) joinCluster(ctx context.Context) error ***REMOVED***
	if n.opts.Addr == "" ***REMOVED***
		return errors.New("attempted to join raft cluster without knowing own address")
	***REMOVED***

	conn, err := dial(n.opts.JoinAddr, "tcp", n.opts.TLSCredentials, 10*time.Second)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer conn.Close()
	client := api.NewRaftMembershipClient(conn)

	joinCtx, joinCancel := context.WithTimeout(ctx, n.reqTimeout())
	defer joinCancel()
	resp, err := client.Join(joinCtx, &api.JoinRequest***REMOVED***
		Addr: n.opts.Addr,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	n.Config.ID = resp.RaftID
	n.bootstrapMembers = resp.Members
	return nil
***REMOVED***

// DefaultNodeConfig returns the default config for a
// raft node that can be modified and customized
func DefaultNodeConfig() *raft.Config ***REMOVED***
	return &raft.Config***REMOVED***
		HeartbeatTick:   1,
		ElectionTick:    3,
		MaxSizePerMsg:   math.MaxUint16,
		MaxInflightMsgs: 256,
		Logger:          log.L,
		CheckQuorum:     true,
	***REMOVED***
***REMOVED***

// DefaultRaftConfig returns a default api.RaftConfig.
func DefaultRaftConfig() api.RaftConfig ***REMOVED***
	return api.RaftConfig***REMOVED***
		KeepOldSnapshots:           0,
		SnapshotInterval:           10000,
		LogEntriesForSlowFollowers: 500,
		ElectionTick:               3,
		HeartbeatTick:              1,
	***REMOVED***
***REMOVED***

// MemoryStore returns the memory store that is kept in sync with the raft log.
func (n *Node) MemoryStore() *store.MemoryStore ***REMOVED***
	return n.memoryStore
***REMOVED***

func (n *Node) done() ***REMOVED***
	n.cluster.Clear()

	n.ticker.Stop()
	n.leadershipBroadcast.Close()
	n.cluster.PeersBroadcast.Close()
	n.memoryStore.Close()
	if n.transport != nil ***REMOVED***
		n.transport.Stop()
	***REMOVED***

	close(n.doneCh)
***REMOVED***

// ClearData tells the raft node to delete its WALs, snapshots, and keys on
// shutdown.
func (n *Node) ClearData() ***REMOVED***
	n.clearData = true
***REMOVED***

// Run is the main loop for a Raft node, it goes along the state machine,
// acting on the messages received from other Raft nodes in the cluster.
//
// Before running the main loop, it first starts the raft node based on saved
// cluster state. If no saved state exists, it starts a single-node cluster.
func (n *Node) Run(ctx context.Context) error ***REMOVED***
	ctx = log.WithLogger(ctx, logrus.WithField("raft_id", fmt.Sprintf("%x", n.Config.ID)))
	ctx, cancel := context.WithCancel(ctx)

	for _, node := range n.bootstrapMembers ***REMOVED***
		if err := n.registerNode(node); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed to register member %x", node.RaftID)
		***REMOVED***
	***REMOVED***

	defer func() ***REMOVED***
		cancel()
		n.stop(ctx)
		if n.clearData ***REMOVED***
			// Delete WAL and snapshots, since they are no longer
			// usable.
			if err := n.raftLogger.Clear(ctx); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("failed to move wal after node removal")
			***REMOVED***
			// clear out the DEKs
			if err := n.keyRotator.UpdateKeys(EncryptionKeys***REMOVED******REMOVED***); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("could not remove DEKs")
			***REMOVED***
		***REMOVED***
		n.done()
	***REMOVED***()

	// Flag that indicates if this manager node is *currently* the raft leader.
	wasLeader := false
	transferLeadershipLimit := rate.NewLimiter(rate.Every(time.Minute), 1)

	for ***REMOVED***
		select ***REMOVED***
		case <-n.ticker.C():
			n.raftNode.Tick()

			if n.leader() == raft.None ***REMOVED***
				atomic.AddUint32(&n.ticksWithNoLeader, 1)
			***REMOVED*** else ***REMOVED***
				atomic.StoreUint32(&n.ticksWithNoLeader, 0)
			***REMOVED***
		case rd := <-n.raftNode.Ready():
			raftConfig := n.getCurrentRaftConfig()

			// Save entries to storage
			if err := n.saveToStorage(ctx, &raftConfig, rd.HardState, rd.Entries, rd.Snapshot); err != nil ***REMOVED***
				return errors.Wrap(err, "failed to save entries to storage")
			***REMOVED***

			// If the memory store lock has been held for too long,
			// transferring leadership is an easy way to break out of it.
			if wasLeader &&
				(rd.SoftState == nil || rd.SoftState.RaftState == raft.StateLeader) &&
				n.memoryStore.Wedged() &&
				transferLeadershipLimit.Allow() ***REMOVED***
				log.G(ctx).Error("Attempting to transfer leadership")
				if !n.opts.DisableStackDump ***REMOVED***
					signal.DumpStacks("")
				***REMOVED***
				transferee, err := n.transport.LongestActive()
				if err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("failed to get longest-active member")
				***REMOVED*** else ***REMOVED***
					log.G(ctx).Error("data store lock held too long - transferring leadership")
					n.raftNode.TransferLeadership(ctx, n.Config.ID, transferee)
				***REMOVED***
			***REMOVED***

			for _, msg := range rd.Messages ***REMOVED***
				// Send raft messages to peers
				if err := n.transport.Send(msg); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("failed to send message to member")
				***REMOVED***
			***REMOVED***

			// Apply snapshot to memory store. The snapshot
			// was applied to the raft store in
			// saveToStorage.
			if !raft.IsEmptySnap(rd.Snapshot) ***REMOVED***
				// Load the snapshot data into the store
				if err := n.restoreFromSnapshot(ctx, rd.Snapshot.Data); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("failed to restore cluster from snapshot")
				***REMOVED***
				n.appliedIndex = rd.Snapshot.Metadata.Index
				n.snapshotMeta = rd.Snapshot.Metadata
				n.confState = rd.Snapshot.Metadata.ConfState
			***REMOVED***

			// If we cease to be the leader, we must cancel any
			// proposals that are currently waiting for a quorum to
			// acknowledge them. It is still possible for these to
			// become committed, but if that happens we will apply
			// them as any follower would.

			// It is important that we cancel these proposals before
			// calling processCommitted, so processCommitted does
			// not deadlock.

			if rd.SoftState != nil ***REMOVED***
				if wasLeader && rd.SoftState.RaftState != raft.StateLeader ***REMOVED***
					wasLeader = false
					log.G(ctx).Error("soft state changed, node no longer a leader, resetting and cancelling all waits")

					if atomic.LoadUint32(&n.signalledLeadership) == 1 ***REMOVED***
						atomic.StoreUint32(&n.signalledLeadership, 0)
						n.leadershipBroadcast.Publish(IsFollower)
					***REMOVED***

					// It is important that we set n.signalledLeadership to 0
					// before calling n.wait.cancelAll. When a new raft
					// request is registered, it checks n.signalledLeadership
					// afterwards, and cancels the registration if it is 0.
					// If cancelAll was called first, this call might run
					// before the new request registers, but
					// signalledLeadership would be set after the check.
					// Setting signalledLeadership before calling cancelAll
					// ensures that if a new request is registered during
					// this transition, it will either be cancelled by
					// cancelAll, or by its own check of signalledLeadership.
					n.wait.cancelAll()
				***REMOVED*** else if !wasLeader && rd.SoftState.RaftState == raft.StateLeader ***REMOVED***
					// Node just became a leader.
					wasLeader = true
				***REMOVED***
			***REMOVED***

			// Process committed entries
			for _, entry := range rd.CommittedEntries ***REMOVED***
				if err := n.processCommitted(ctx, entry); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("failed to process committed entries")
				***REMOVED***
			***REMOVED***

			// in case the previous attempt to update the key failed
			n.maybeMarkRotationFinished(ctx)

			// Trigger a snapshot every once in awhile
			if n.snapshotInProgress == nil &&
				(n.needsSnapshot(ctx) || raftConfig.SnapshotInterval > 0 &&
					n.appliedIndex-n.snapshotMeta.Index >= raftConfig.SnapshotInterval) ***REMOVED***
				n.triggerSnapshot(ctx, raftConfig)
			***REMOVED***

			if wasLeader && atomic.LoadUint32(&n.signalledLeadership) != 1 ***REMOVED***
				// If all the entries in the log have become
				// committed, broadcast our leadership status.
				if n.caughtUp() ***REMOVED***
					atomic.StoreUint32(&n.signalledLeadership, 1)
					n.leadershipBroadcast.Publish(IsLeader)
				***REMOVED***
			***REMOVED***

			// Advance the state machine
			n.raftNode.Advance()

			// On the first startup, or if we are the only
			// registered member after restoring from the state,
			// campaign to be the leader.
			if n.campaignWhenAble ***REMOVED***
				members := n.cluster.Members()
				if len(members) >= 1 ***REMOVED***
					n.campaignWhenAble = false
				***REMOVED***
				if len(members) == 1 && members[n.Config.ID] != nil ***REMOVED***
					n.raftNode.Campaign(ctx)
				***REMOVED***
			***REMOVED***

		case snapshotMeta := <-n.snapshotInProgress:
			raftConfig := n.getCurrentRaftConfig()
			if snapshotMeta.Index > n.snapshotMeta.Index ***REMOVED***
				n.snapshotMeta = snapshotMeta
				if err := n.raftLogger.GC(snapshotMeta.Index, snapshotMeta.Term, raftConfig.KeepOldSnapshots); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("failed to clean up old snapshots and WALs")
				***REMOVED***
			***REMOVED***
			n.snapshotInProgress = nil
			n.maybeMarkRotationFinished(ctx)
			if n.rotationQueued && n.needsSnapshot(ctx) ***REMOVED***
				// there was a key rotation that took place before while the snapshot
				// was in progress - we have to take another snapshot and encrypt with the new key
				n.rotationQueued = false
				n.triggerSnapshot(ctx, raftConfig)
			***REMOVED***
		case <-n.keyRotator.RotationNotify():
			// There are 2 separate checks:  rotationQueued, and n.needsSnapshot().
			// We set rotationQueued so that when we are notified of a rotation, we try to
			// do a snapshot as soon as possible.  However, if there is an error while doing
			// the snapshot, we don't want to hammer the node attempting to do snapshots over
			// and over.  So if doing a snapshot fails, wait until the next entry comes in to
			// try again.
			switch ***REMOVED***
			case n.snapshotInProgress != nil:
				n.rotationQueued = true
			case n.needsSnapshot(ctx):
				n.triggerSnapshot(ctx, n.getCurrentRaftConfig())
			***REMOVED***
		case <-ctx.Done():
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (n *Node) restoreFromSnapshot(ctx context.Context, data []byte) error ***REMOVED***
	snapCluster, err := n.clusterSnapshot(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	oldMembers := n.cluster.Members()

	for _, member := range snapCluster.Members ***REMOVED***
		delete(oldMembers, member.RaftID)
	***REMOVED***

	for _, removedMember := range snapCluster.Removed ***REMOVED***
		n.cluster.RemoveMember(removedMember)
		n.transport.RemovePeer(removedMember)
		delete(oldMembers, removedMember)
	***REMOVED***

	for id, member := range oldMembers ***REMOVED***
		n.cluster.ClearMember(id)
		if err := n.transport.RemovePeer(member.RaftID); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed to remove peer %x from transport", member.RaftID)
		***REMOVED***
	***REMOVED***
	for _, node := range snapCluster.Members ***REMOVED***
		if err := n.registerNode(&api.RaftMember***REMOVED***RaftID: node.RaftID, NodeID: node.NodeID, Addr: node.Addr***REMOVED***); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("failed to register node from snapshot")
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (n *Node) needsSnapshot(ctx context.Context) bool ***REMOVED***
	if n.waitForAppliedIndex == 0 && n.keyRotator.NeedsRotation() ***REMOVED***
		keys := n.keyRotator.GetKeys()
		if keys.PendingDEK != nil ***REMOVED***
			n.raftLogger.RotateEncryptionKey(keys.PendingDEK)
			// we want to wait for the last index written with the old DEK to be committed, else a snapshot taken
			// may have an index less than the index of a WAL written with an old DEK.  We want the next snapshot
			// written with the new key to supercede any WAL written with an old DEK.
			n.waitForAppliedIndex = n.writtenWALIndex
			// if there is already a snapshot at this index or higher, bump the wait index up to 1 higher than the current
			// snapshot index, because the rotation cannot be completed until the next snapshot
			if n.waitForAppliedIndex <= n.snapshotMeta.Index ***REMOVED***
				n.waitForAppliedIndex = n.snapshotMeta.Index + 1
			***REMOVED***
			log.G(ctx).Debugf(
				"beginning raft DEK rotation - last indices written with the old key are (snapshot: %d, WAL: %d) - waiting for snapshot of index %d to be written before rotation can be completed", n.snapshotMeta.Index, n.writtenWALIndex, n.waitForAppliedIndex)
		***REMOVED***
	***REMOVED***

	result := n.waitForAppliedIndex > 0 && n.waitForAppliedIndex <= n.appliedIndex
	if result ***REMOVED***
		log.G(ctx).Debugf(
			"a snapshot at index %d is needed in order to complete raft DEK rotation - a snapshot with index >= %d can now be triggered",
			n.waitForAppliedIndex, n.appliedIndex)
	***REMOVED***
	return result
***REMOVED***

func (n *Node) maybeMarkRotationFinished(ctx context.Context) ***REMOVED***
	if n.waitForAppliedIndex > 0 && n.waitForAppliedIndex <= n.snapshotMeta.Index ***REMOVED***
		// this means we tried to rotate - so finish the rotation
		if err := n.keyRotator.UpdateKeys(EncryptionKeys***REMOVED***CurrentDEK: n.raftLogger.EncryptionKey***REMOVED***); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("failed to update encryption keys after a successful rotation")
		***REMOVED*** else ***REMOVED***
			log.G(ctx).Debugf(
				"a snapshot with index %d is available, which completes the DEK rotation requiring a snapshot of at least index %d - throwing away DEK and older snapshots encrypted with the old key",
				n.snapshotMeta.Index, n.waitForAppliedIndex)
			n.waitForAppliedIndex = 0

			if err := n.raftLogger.GC(n.snapshotMeta.Index, n.snapshotMeta.Term, 0); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("failed to remove old snapshots and WALs that were written with the previous raft DEK")
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (n *Node) getCurrentRaftConfig() api.RaftConfig ***REMOVED***
	raftConfig := DefaultRaftConfig()
	n.memoryStore.View(func(readTx store.ReadTx) ***REMOVED***
		clusters, err := store.FindClusters(readTx, store.ByName(store.DefaultClusterName))
		if err == nil && len(clusters) == 1 ***REMOVED***
			raftConfig = clusters[0].Spec.Raft
		***REMOVED***
	***REMOVED***)
	return raftConfig
***REMOVED***

// Cancel interrupts all ongoing proposals, and prevents new ones from
// starting. This is useful for the shutdown sequence because it allows
// the manager to shut down raft-dependent services that might otherwise
// block on shutdown if quorum isn't met. Then the raft node can be completely
// shut down once no more code is using it.
func (n *Node) Cancel() ***REMOVED***
	n.cancelFunc()
***REMOVED***

// Done returns channel which is closed when raft node is fully stopped.
func (n *Node) Done() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return n.doneCh
***REMOVED***

func (n *Node) stop(ctx context.Context) ***REMOVED***
	n.stopMu.Lock()
	defer n.stopMu.Unlock()

	n.Cancel()
	n.waitProp.Wait()
	n.asyncTasks.Wait()

	n.raftNode.Stop()
	n.ticker.Stop()
	n.raftLogger.Close(ctx)
	atomic.StoreUint32(&n.isMember, 0)
	// TODO(stevvooe): Handle ctx.Done()
***REMOVED***

// isLeader checks if we are the leader or not, without the protection of lock
func (n *Node) isLeader() bool ***REMOVED***
	if !n.IsMember() ***REMOVED***
		return false
	***REMOVED***

	if n.Status().Lead == n.Config.ID ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// IsLeader checks if we are the leader or not, with the protection of lock
func (n *Node) IsLeader() bool ***REMOVED***
	n.stopMu.RLock()
	defer n.stopMu.RUnlock()

	return n.isLeader()
***REMOVED***

// leader returns the id of the leader, without the protection of lock and
// membership check, so it's caller task.
func (n *Node) leader() uint64 ***REMOVED***
	return n.Status().Lead
***REMOVED***

// Leader returns the id of the leader, with the protection of lock
func (n *Node) Leader() (uint64, error) ***REMOVED***
	n.stopMu.RLock()
	defer n.stopMu.RUnlock()

	if !n.IsMember() ***REMOVED***
		return raft.None, ErrNoRaftMember
	***REMOVED***
	leader := n.leader()
	if leader == raft.None ***REMOVED***
		return raft.None, ErrNoClusterLeader
	***REMOVED***

	return leader, nil
***REMOVED***

// ReadyForProposals returns true if the node has broadcasted a message
// saying that it has become the leader. This means it is ready to accept
// proposals.
func (n *Node) ReadyForProposals() bool ***REMOVED***
	return atomic.LoadUint32(&n.signalledLeadership) == 1
***REMOVED***

func (n *Node) caughtUp() bool ***REMOVED***
	// obnoxious function that always returns a nil error
	lastIndex, _ := n.raftStore.LastIndex()
	return n.appliedIndex >= lastIndex
***REMOVED***

// Join asks to a member of the raft to propose
// a configuration change and add us as a member thus
// beginning the log replication process. This method
// is called from an aspiring member to an existing member
func (n *Node) Join(ctx context.Context, req *api.JoinRequest) (*api.JoinResponse, error) ***REMOVED***
	nodeInfo, err := ca.RemoteNode(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	fields := logrus.Fields***REMOVED***
		"node.id": nodeInfo.NodeID,
		"method":  "(*Node).Join",
		"raft_id": fmt.Sprintf("%x", n.Config.ID),
	***REMOVED***
	if nodeInfo.ForwardedBy != nil ***REMOVED***
		fields["forwarder.id"] = nodeInfo.ForwardedBy.NodeID
	***REMOVED***
	log := log.G(ctx).WithFields(fields)
	log.Debug("")

	// can't stop the raft node while an async RPC is in progress
	n.stopMu.RLock()
	defer n.stopMu.RUnlock()

	n.membershipLock.Lock()
	defer n.membershipLock.Unlock()

	if !n.IsMember() ***REMOVED***
		return nil, status.Errorf(codes.FailedPrecondition, "%s", ErrNoRaftMember.Error())
	***REMOVED***

	if !n.isLeader() ***REMOVED***
		return nil, status.Errorf(codes.FailedPrecondition, "%s", ErrLostLeadership.Error())
	***REMOVED***

	remoteAddr := req.Addr

	// If the joining node sent an address like 0.0.0.0:4242, automatically
	// determine its actual address based on the GRPC connection. This
	// avoids the need for a prospective member to know its own address.

	requestHost, requestPort, err := net.SplitHostPort(remoteAddr)
	if err != nil ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "invalid address %s in raft join request", remoteAddr)
	***REMOVED***

	requestIP := net.ParseIP(requestHost)
	if requestIP != nil && requestIP.IsUnspecified() ***REMOVED***
		remoteHost, _, err := net.SplitHostPort(nodeInfo.RemoteAddr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		remoteAddr = net.JoinHostPort(remoteHost, requestPort)
	***REMOVED***

	// We do not bother submitting a configuration change for the
	// new member if we can't contact it back using its address
	if err := n.checkHealth(ctx, remoteAddr, 5*time.Second); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// If the peer is already a member of the cluster, we will only update
	// its information, not add it as a new member. Adding it again would
	// cause the quorum to be computed incorrectly.
	for _, m := range n.cluster.Members() ***REMOVED***
		if m.NodeID == nodeInfo.NodeID ***REMOVED***
			if remoteAddr == m.Addr ***REMOVED***
				return n.joinResponse(m.RaftID), nil
			***REMOVED***
			updatedRaftMember := &api.RaftMember***REMOVED***
				RaftID: m.RaftID,
				NodeID: m.NodeID,
				Addr:   remoteAddr,
			***REMOVED***
			if err := n.cluster.UpdateMember(m.RaftID, updatedRaftMember); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if err := n.updateNodeBlocking(ctx, m.RaftID, remoteAddr); err != nil ***REMOVED***
				log.WithError(err).Error("failed to update node address")
				return nil, err
			***REMOVED***

			log.Info("updated node address")
			return n.joinResponse(m.RaftID), nil
		***REMOVED***
	***REMOVED***

	// Find a unique ID for the joining member.
	var raftID uint64
	for ***REMOVED***
		raftID = uint64(rand.Int63()) + 1
		if n.cluster.GetMember(raftID) == nil && !n.cluster.IsIDRemoved(raftID) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	err = n.addMember(ctx, remoteAddr, raftID, nodeInfo.NodeID)
	if err != nil ***REMOVED***
		log.WithError(err).Errorf("failed to add member %x", raftID)
		return nil, err
	***REMOVED***

	log.Debug("node joined")

	return n.joinResponse(raftID), nil
***REMOVED***

func (n *Node) joinResponse(raftID uint64) *api.JoinResponse ***REMOVED***
	var nodes []*api.RaftMember
	for _, node := range n.cluster.Members() ***REMOVED***
		nodes = append(nodes, &api.RaftMember***REMOVED***
			RaftID: node.RaftID,
			NodeID: node.NodeID,
			Addr:   node.Addr,
		***REMOVED***)
	***REMOVED***

	return &api.JoinResponse***REMOVED***Members: nodes, RaftID: raftID***REMOVED***
***REMOVED***

// checkHealth tries to contact an aspiring member through its advertised address
// and checks if its raft server is running.
func (n *Node) checkHealth(ctx context.Context, addr string, timeout time.Duration) error ***REMOVED***
	conn, err := dial(addr, "tcp", n.opts.TLSCredentials, timeout)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	defer conn.Close()

	if timeout != 0 ***REMOVED***
		tctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		ctx = tctx
	***REMOVED***

	healthClient := api.NewHealthClient(conn)
	resp, err := healthClient.Check(ctx, &api.HealthCheckRequest***REMOVED***Service: "Raft"***REMOVED***)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "could not connect to prospective new cluster member using its advertised address")
	***REMOVED***
	if resp.Status != api.HealthCheckResponse_SERVING ***REMOVED***
		return fmt.Errorf("health check returned status %s", resp.Status.String())
	***REMOVED***

	return nil
***REMOVED***

// addMember submits a configuration change to add a new member on the raft cluster.
func (n *Node) addMember(ctx context.Context, addr string, raftID uint64, nodeID string) error ***REMOVED***
	node := api.RaftMember***REMOVED***
		RaftID: raftID,
		NodeID: nodeID,
		Addr:   addr,
	***REMOVED***

	meta, err := node.Marshal()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cc := raftpb.ConfChange***REMOVED***
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  raftID,
		Context: meta,
	***REMOVED***

	// Wait for a raft round to process the configuration change
	return n.configure(ctx, cc)
***REMOVED***

// updateNodeBlocking runs synchronous job to update node address in whole cluster.
func (n *Node) updateNodeBlocking(ctx context.Context, id uint64, addr string) error ***REMOVED***
	m := n.cluster.GetMember(id)
	if m == nil ***REMOVED***
		return errors.Errorf("member %x is not found for update", id)
	***REMOVED***
	node := api.RaftMember***REMOVED***
		RaftID: m.RaftID,
		NodeID: m.NodeID,
		Addr:   addr,
	***REMOVED***

	meta, err := node.Marshal()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cc := raftpb.ConfChange***REMOVED***
		Type:    raftpb.ConfChangeUpdateNode,
		NodeID:  id,
		Context: meta,
	***REMOVED***

	// Wait for a raft round to process the configuration change
	return n.configure(ctx, cc)
***REMOVED***

// UpdateNode submits a configuration change to change a member's address.
func (n *Node) UpdateNode(id uint64, addr string) ***REMOVED***
	ctx, cancel := n.WithContext(context.Background())
	defer cancel()
	// spawn updating info in raft in background to unblock transport
	go func() ***REMOVED***
		if err := n.updateNodeBlocking(ctx, id, addr); err != nil ***REMOVED***
			log.G(ctx).WithFields(logrus.Fields***REMOVED***"raft_id": n.Config.ID, "update_id": id***REMOVED***).WithError(err).Error("failed to update member address in cluster")
		***REMOVED***
	***REMOVED***()
***REMOVED***

// Leave asks to a member of the raft to remove
// us from the raft cluster. This method is called
// from a member who is willing to leave its raft
// membership to an active member of the raft
func (n *Node) Leave(ctx context.Context, req *api.LeaveRequest) (*api.LeaveResponse, error) ***REMOVED***
	if req.Node == nil ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "no node information provided")
	***REMOVED***

	nodeInfo, err := ca.RemoteNode(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ctx, cancel := n.WithContext(ctx)
	defer cancel()

	fields := logrus.Fields***REMOVED***
		"node.id": nodeInfo.NodeID,
		"method":  "(*Node).Leave",
		"raft_id": fmt.Sprintf("%x", n.Config.ID),
	***REMOVED***
	if nodeInfo.ForwardedBy != nil ***REMOVED***
		fields["forwarder.id"] = nodeInfo.ForwardedBy.NodeID
	***REMOVED***
	log.G(ctx).WithFields(fields).Debug("")

	if err := n.removeMember(ctx, req.Node.RaftID); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &api.LeaveResponse***REMOVED******REMOVED***, nil
***REMOVED***

// CanRemoveMember checks if a member can be removed from
// the context of the current node.
func (n *Node) CanRemoveMember(id uint64) bool ***REMOVED***
	members := n.cluster.Members()
	nreachable := 0 // reachable managers after removal

	for _, m := range members ***REMOVED***
		if m.RaftID == id ***REMOVED***
			continue
		***REMOVED***

		// Local node from where the remove is issued
		if m.RaftID == n.Config.ID ***REMOVED***
			nreachable++
			continue
		***REMOVED***

		if n.transport.Active(m.RaftID) ***REMOVED***
			nreachable++
		***REMOVED***
	***REMOVED***

	nquorum := (len(members)-1)/2 + 1
	if nreachable < nquorum ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***

func (n *Node) removeMember(ctx context.Context, id uint64) error ***REMOVED***
	// can't stop the raft node while an async RPC is in progress
	n.stopMu.RLock()
	defer n.stopMu.RUnlock()

	if !n.IsMember() ***REMOVED***
		return ErrNoRaftMember
	***REMOVED***

	if !n.isLeader() ***REMOVED***
		return ErrLostLeadership
	***REMOVED***

	n.membershipLock.Lock()
	defer n.membershipLock.Unlock()
	if !n.CanRemoveMember(id) ***REMOVED***
		return ErrCannotRemoveMember
	***REMOVED***

	cc := raftpb.ConfChange***REMOVED***
		ID:      id,
		Type:    raftpb.ConfChangeRemoveNode,
		NodeID:  id,
		Context: []byte(""),
	***REMOVED***
	return n.configure(ctx, cc)
***REMOVED***

// TransferLeadership attempts to transfer leadership to a different node,
// and wait for the transfer to happen.
func (n *Node) TransferLeadership(ctx context.Context) error ***REMOVED***
	ctx, cancelTransfer := context.WithTimeout(ctx, n.reqTimeout())
	defer cancelTransfer()

	n.stopMu.RLock()
	defer n.stopMu.RUnlock()

	if !n.IsMember() ***REMOVED***
		return ErrNoRaftMember
	***REMOVED***

	if !n.isLeader() ***REMOVED***
		return ErrLostLeadership
	***REMOVED***

	transferee, err := n.transport.LongestActive()
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to get longest-active member")
	***REMOVED***
	start := time.Now()
	n.raftNode.TransferLeadership(ctx, n.Config.ID, transferee)
	ticker := time.NewTicker(n.opts.TickInterval / 10)
	defer ticker.Stop()
	var leader uint64
	for ***REMOVED***
		leader = n.leader()
		if leader != raft.None && leader != n.Config.ID ***REMOVED***
			break
		***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		***REMOVED***
	***REMOVED***
	log.G(ctx).Infof("raft: transfer leadership %x -> %x finished in %v", n.Config.ID, leader, time.Since(start))
	return nil
***REMOVED***

// RemoveMember submits a configuration change to remove a member from the raft cluster
// after checking if the operation would not result in a loss of quorum.
func (n *Node) RemoveMember(ctx context.Context, id uint64) error ***REMOVED***
	ctx, cancel := n.WithContext(ctx)
	defer cancel()
	return n.removeMember(ctx, id)
***REMOVED***

// processRaftMessageLogger is used to lazily create a logger for
// ProcessRaftMessage. Usually nothing will be logged, so it is useful to avoid
// formatting strings and allocating a logger when it won't be used.
func (n *Node) processRaftMessageLogger(ctx context.Context, msg *api.ProcessRaftMessageRequest) *logrus.Entry ***REMOVED***
	fields := logrus.Fields***REMOVED***
		"method": "(*Node).ProcessRaftMessage",
	***REMOVED***

	if n.IsMember() ***REMOVED***
		fields["raft_id"] = fmt.Sprintf("%x", n.Config.ID)
	***REMOVED***

	if msg != nil && msg.Message != nil ***REMOVED***
		fields["from"] = fmt.Sprintf("%x", msg.Message.From)
	***REMOVED***

	return log.G(ctx).WithFields(fields)
***REMOVED***

func (n *Node) reportNewAddress(ctx context.Context, id uint64) error ***REMOVED***
	// too early
	if !n.IsMember() ***REMOVED***
		return nil
	***REMOVED***
	p, ok := peer.FromContext(ctx)
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	oldAddr, err := n.transport.PeerAddr(id)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if oldAddr == "" ***REMOVED***
		// Don't know the address of the peer yet, so can't report an
		// update.
		return nil
	***REMOVED***
	newHost, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, officialPort, err := net.SplitHostPort(oldAddr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	newAddr := net.JoinHostPort(newHost, officialPort)
	return n.transport.UpdatePeerAddr(id, newAddr)
***REMOVED***

// StreamRaftMessage is the server endpoint for streaming Raft messages.
// It accepts a stream of raft messages to be processed on this raft member,
// returning a StreamRaftMessageResponse when processing of the streamed
// messages is complete.
// It is called from the Raft leader, which uses it to stream messages
// to this raft member.
// A single stream corresponds to a single raft message,
// which may be disassembled and streamed by the sender
// as individual messages. Therefore, each of the messages
// received by the stream will have the same raft message type and index.
// Currently, only messages of type raftpb.MsgSnap can be disassembled, sent
// and received on the stream.
func (n *Node) StreamRaftMessage(stream api.Raft_StreamRaftMessageServer) error ***REMOVED***
	// recvdMsg is the current messasge received from the stream.
	// assembledMessage is where the data from recvdMsg is appended to.
	var recvdMsg, assembledMessage *api.StreamRaftMessageRequest
	var err error

	// First message index.
	var raftMsgIndex uint64

	for ***REMOVED***
		recvdMsg, err = stream.Recv()
		if err == io.EOF ***REMOVED***
			break
		***REMOVED*** else if err != nil ***REMOVED***
			log.G(stream.Context()).WithError(err).Error("error while reading from stream")
			return err
		***REMOVED***

		// Initialized the message to be used for assembling
		// the raft message.
		if assembledMessage == nil ***REMOVED***
			// For all message types except raftpb.MsgSnap,
			// we don't expect more than a single message
			// on the stream so we'll get an EOF on the next Recv()
			// and go on to process the received message.
			assembledMessage = recvdMsg
			raftMsgIndex = recvdMsg.Message.Index
			continue
		***REMOVED***

		// Verify raft message index.
		if recvdMsg.Message.Index != raftMsgIndex ***REMOVED***
			errMsg := fmt.Sprintf("Raft message chunk with index %d is different from the previously received raft message index %d",
				recvdMsg.Message.Index, raftMsgIndex)
			log.G(stream.Context()).Errorf(errMsg)
			return status.Errorf(codes.InvalidArgument, "%s", errMsg)
		***REMOVED***

		// Verify that multiple message received on a stream
		// can only be of type raftpb.MsgSnap.
		if recvdMsg.Message.Type != raftpb.MsgSnap ***REMOVED***
			errMsg := fmt.Sprintf("Raft message chunk is not of type %d",
				raftpb.MsgSnap)
			log.G(stream.Context()).Errorf(errMsg)
			return status.Errorf(codes.InvalidArgument, "%s", errMsg)
		***REMOVED***

		// Append the received snapshot data.
		assembledMessage.Message.Snapshot.Data = append(assembledMessage.Message.Snapshot.Data, recvdMsg.Message.Snapshot.Data...)
	***REMOVED***

	// We should have the complete snapshot. Verify and process.
	if err == io.EOF ***REMOVED***
		_, err = n.ProcessRaftMessage(stream.Context(), &api.ProcessRaftMessageRequest***REMOVED***Message: assembledMessage.Message***REMOVED***)
		if err == nil ***REMOVED***
			// Translate the response of ProcessRaftMessage() from
			// ProcessRaftMessageResponse to StreamRaftMessageResponse if needed.
			return stream.SendAndClose(&api.StreamRaftMessageResponse***REMOVED******REMOVED***)
		***REMOVED***
	***REMOVED***

	return err
***REMOVED***

// ProcessRaftMessage calls 'Step' which advances the
// raft state machine with the provided message on the
// receiving node
func (n *Node) ProcessRaftMessage(ctx context.Context, msg *api.ProcessRaftMessageRequest) (*api.ProcessRaftMessageResponse, error) ***REMOVED***
	if msg == nil || msg.Message == nil ***REMOVED***
		n.processRaftMessageLogger(ctx, msg).Debug("received empty message")
		return &api.ProcessRaftMessageResponse***REMOVED******REMOVED***, nil
	***REMOVED***

	// Don't process the message if this comes from
	// a node in the remove set
	if n.cluster.IsIDRemoved(msg.Message.From) ***REMOVED***
		n.processRaftMessageLogger(ctx, msg).Debug("received message from removed member")
		return nil, status.Errorf(codes.NotFound, "%s", membership.ErrMemberRemoved.Error())
	***REMOVED***

	ctx, cancel := n.WithContext(ctx)
	defer cancel()

	// TODO(aaronl): Address changes are temporarily disabled.
	// See https://github.com/docker/docker/issues/30455.
	// This should be reenabled in the future with additional
	// safeguards (perhaps storing multiple addresses per node).
	//if err := n.reportNewAddress(ctx, msg.Message.From); err != nil ***REMOVED***
	//	log.G(ctx).WithError(err).Errorf("failed to report new address of %x to transport", msg.Message.From)
	//***REMOVED***

	// Reject vote requests from unreachable peers
	if msg.Message.Type == raftpb.MsgVote ***REMOVED***
		member := n.cluster.GetMember(msg.Message.From)
		if member == nil ***REMOVED***
			n.processRaftMessageLogger(ctx, msg).Debug("received message from unknown member")
			return &api.ProcessRaftMessageResponse***REMOVED******REMOVED***, nil
		***REMOVED***

		if err := n.transport.HealthCheck(ctx, msg.Message.From); err != nil ***REMOVED***
			n.processRaftMessageLogger(ctx, msg).WithError(err).Debug("member which sent vote request failed health check")
			return &api.ProcessRaftMessageResponse***REMOVED******REMOVED***, nil
		***REMOVED***
	***REMOVED***

	if msg.Message.Type == raftpb.MsgProp ***REMOVED***
		// We don't accept forwarded proposals. Our
		// current architecture depends on only the leader
		// making proposals, so in-flight proposals can be
		// guaranteed not to conflict.
		n.processRaftMessageLogger(ctx, msg).Debug("dropped forwarded proposal")
		return &api.ProcessRaftMessageResponse***REMOVED******REMOVED***, nil
	***REMOVED***

	// can't stop the raft node while an async RPC is in progress
	n.stopMu.RLock()
	defer n.stopMu.RUnlock()

	if n.IsMember() ***REMOVED***
		if msg.Message.To != n.Config.ID ***REMOVED***
			n.processRaftMessageLogger(ctx, msg).Errorf("received message intended for raft_id %x", msg.Message.To)
			return &api.ProcessRaftMessageResponse***REMOVED******REMOVED***, nil
		***REMOVED***

		if err := n.raftNode.Step(ctx, *msg.Message); err != nil ***REMOVED***
			n.processRaftMessageLogger(ctx, msg).WithError(err).Debug("raft Step failed")
		***REMOVED***
	***REMOVED***

	return &api.ProcessRaftMessageResponse***REMOVED******REMOVED***, nil
***REMOVED***

// ResolveAddress returns the address reaching for a given node ID.
func (n *Node) ResolveAddress(ctx context.Context, msg *api.ResolveAddressRequest) (*api.ResolveAddressResponse, error) ***REMOVED***
	if !n.IsMember() ***REMOVED***
		return nil, ErrNoRaftMember
	***REMOVED***

	nodeInfo, err := ca.RemoteNode(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	fields := logrus.Fields***REMOVED***
		"node.id": nodeInfo.NodeID,
		"method":  "(*Node).ResolveAddress",
		"raft_id": fmt.Sprintf("%x", n.Config.ID),
	***REMOVED***
	if nodeInfo.ForwardedBy != nil ***REMOVED***
		fields["forwarder.id"] = nodeInfo.ForwardedBy.NodeID
	***REMOVED***
	log.G(ctx).WithFields(fields).Debug("")

	member := n.cluster.GetMember(msg.RaftID)
	if member == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "member %x not found", msg.RaftID)
	***REMOVED***
	return &api.ResolveAddressResponse***REMOVED***Addr: member.Addr***REMOVED***, nil
***REMOVED***

func (n *Node) getLeaderConn() (*grpc.ClientConn, error) ***REMOVED***
	leader, err := n.Leader()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if leader == n.Config.ID ***REMOVED***
		return nil, raftselector.ErrIsLeader
	***REMOVED***
	conn, err := n.transport.PeerConn(leader)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to get connection to leader")
	***REMOVED***
	return conn, nil
***REMOVED***

// LeaderConn returns current connection to cluster leader or raftselector.ErrIsLeader
// if current machine is leader.
func (n *Node) LeaderConn(ctx context.Context) (*grpc.ClientConn, error) ***REMOVED***
	cc, err := n.getLeaderConn()
	if err == nil ***REMOVED***
		return cc, nil
	***REMOVED***
	if err == raftselector.ErrIsLeader ***REMOVED***
		return nil, err
	***REMOVED***
	if atomic.LoadUint32(&n.ticksWithNoLeader) > lostQuorumTimeout ***REMOVED***
		return nil, errLostQuorum
	***REMOVED***

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			cc, err := n.getLeaderConn()
			if err == nil ***REMOVED***
				return cc, nil
			***REMOVED***
			if err == raftselector.ErrIsLeader ***REMOVED***
				return nil, err
			***REMOVED***
		case <-ctx.Done():
			return nil, ctx.Err()
		***REMOVED***
	***REMOVED***
***REMOVED***

// registerNode registers a new node on the cluster memberlist
func (n *Node) registerNode(node *api.RaftMember) error ***REMOVED***
	if n.cluster.IsIDRemoved(node.RaftID) ***REMOVED***
		return nil
	***REMOVED***

	member := &membership.Member***REMOVED******REMOVED***

	existingMember := n.cluster.GetMember(node.RaftID)
	if existingMember != nil ***REMOVED***
		// Member already exists

		// If the address is different from what we thought it was,
		// update it. This can happen if we just joined a cluster
		// and are adding ourself now with the remotely-reachable
		// address.
		if existingMember.Addr != node.Addr ***REMOVED***
			if node.RaftID != n.Config.ID ***REMOVED***
				if err := n.transport.UpdatePeer(node.RaftID, node.Addr); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			member.RaftMember = node
			n.cluster.AddMember(member)
		***REMOVED***

		return nil
	***REMOVED***

	// Avoid opening a connection to the local node
	if node.RaftID != n.Config.ID ***REMOVED***
		if err := n.transport.AddPeer(node.RaftID, node.Addr); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	member.RaftMember = node
	err := n.cluster.AddMember(member)
	if err != nil ***REMOVED***
		if rerr := n.transport.RemovePeer(node.RaftID); rerr != nil ***REMOVED***
			return errors.Wrapf(rerr, "failed to remove peer after error %v", err)
		***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// ProposeValue calls Propose on the underlying raft library(etcd/raft) and waits
// on the commit log action before returning a result
func (n *Node) ProposeValue(ctx context.Context, storeAction []api.StoreAction, cb func()) error ***REMOVED***
	defer metrics.StartTimer(proposeLatencyTimer)()
	ctx, cancel := n.WithContext(ctx)
	defer cancel()
	_, err := n.processInternalRaftRequest(ctx, &api.InternalRaftRequest***REMOVED***Action: storeAction***REMOVED***, cb)

	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// GetVersion returns the sequence information for the current raft round.
func (n *Node) GetVersion() *api.Version ***REMOVED***
	n.stopMu.RLock()
	defer n.stopMu.RUnlock()

	if !n.IsMember() ***REMOVED***
		return nil
	***REMOVED***

	status := n.Status()
	return &api.Version***REMOVED***Index: status.Commit***REMOVED***
***REMOVED***

// ChangesBetween returns the changes starting after "from", up to and
// including "to". If these changes are not available because the log
// has been compacted, an error will be returned.
func (n *Node) ChangesBetween(from, to api.Version) ([]state.Change, error) ***REMOVED***
	n.stopMu.RLock()
	defer n.stopMu.RUnlock()

	if from.Index > to.Index ***REMOVED***
		return nil, errors.New("versions are out of order")
	***REMOVED***

	if !n.IsMember() ***REMOVED***
		return nil, ErrNoRaftMember
	***REMOVED***

	// never returns error
	last, _ := n.raftStore.LastIndex()

	if to.Index > last ***REMOVED***
		return nil, errors.New("last version is out of bounds")
	***REMOVED***

	pbs, err := n.raftStore.Entries(from.Index+1, to.Index+1, math.MaxUint64)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var changes []state.Change
	for _, pb := range pbs ***REMOVED***
		if pb.Type != raftpb.EntryNormal || pb.Data == nil ***REMOVED***
			continue
		***REMOVED***
		r := &api.InternalRaftRequest***REMOVED******REMOVED***
		err := proto.Unmarshal(pb.Data, r)
		if err != nil ***REMOVED***
			return nil, errors.Wrap(err, "error umarshalling internal raft request")
		***REMOVED***

		if r.Action != nil ***REMOVED***
			changes = append(changes, state.Change***REMOVED***StoreActions: r.Action, Version: api.Version***REMOVED***Index: pb.Index***REMOVED******REMOVED***)
		***REMOVED***
	***REMOVED***

	return changes, nil
***REMOVED***

// SubscribePeers subscribes to peer updates in cluster. It sends always full
// list of peers.
func (n *Node) SubscribePeers() (q chan events.Event, cancel func()) ***REMOVED***
	return n.cluster.PeersBroadcast.Watch()
***REMOVED***

// GetMemberlist returns the current list of raft members in the cluster.
func (n *Node) GetMemberlist() map[uint64]*api.RaftMember ***REMOVED***
	memberlist := make(map[uint64]*api.RaftMember)
	members := n.cluster.Members()
	leaderID, err := n.Leader()
	if err != nil ***REMOVED***
		leaderID = raft.None
	***REMOVED***

	for id, member := range members ***REMOVED***
		reachability := api.RaftMemberStatus_REACHABLE
		leader := false

		if member.RaftID != n.Config.ID ***REMOVED***
			if !n.transport.Active(member.RaftID) ***REMOVED***
				reachability = api.RaftMemberStatus_UNREACHABLE
			***REMOVED***
		***REMOVED***

		if member.RaftID == leaderID ***REMOVED***
			leader = true
		***REMOVED***

		memberlist[id] = &api.RaftMember***REMOVED***
			RaftID: member.RaftID,
			NodeID: member.NodeID,
			Addr:   member.Addr,
			Status: api.RaftMemberStatus***REMOVED***
				Leader:       leader,
				Reachability: reachability,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	return memberlist
***REMOVED***

// Status returns status of underlying etcd.Node.
func (n *Node) Status() raft.Status ***REMOVED***
	return n.raftNode.Status()
***REMOVED***

// GetMemberByNodeID returns member information based
// on its generic Node ID.
func (n *Node) GetMemberByNodeID(nodeID string) *membership.Member ***REMOVED***
	members := n.cluster.Members()
	for _, member := range members ***REMOVED***
		if member.NodeID == nodeID ***REMOVED***
			return member
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// IsMember checks if the raft node has effectively joined
// a cluster of existing members.
func (n *Node) IsMember() bool ***REMOVED***
	return atomic.LoadUint32(&n.isMember) == 1
***REMOVED***

// Saves a log entry to our Store
func (n *Node) saveToStorage(
	ctx context.Context,
	raftConfig *api.RaftConfig,
	hardState raftpb.HardState,
	entries []raftpb.Entry,
	snapshot raftpb.Snapshot,
) (err error) ***REMOVED***

	if !raft.IsEmptySnap(snapshot) ***REMOVED***
		if err := n.raftLogger.SaveSnapshot(snapshot); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to save snapshot")
		***REMOVED***
		if err := n.raftLogger.GC(snapshot.Metadata.Index, snapshot.Metadata.Term, raftConfig.KeepOldSnapshots); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("unable to clean old snapshots and WALs")
		***REMOVED***
		if err = n.raftStore.ApplySnapshot(snapshot); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to apply snapshot on raft node")
		***REMOVED***
	***REMOVED***

	if err := n.raftLogger.SaveEntries(hardState, entries); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to save raft log entries")
	***REMOVED***

	if len(entries) > 0 ***REMOVED***
		lastIndex := entries[len(entries)-1].Index
		if lastIndex > n.writtenWALIndex ***REMOVED***
			n.writtenWALIndex = lastIndex
		***REMOVED***
	***REMOVED***

	if err = n.raftStore.Append(entries); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to append raft log entries")
	***REMOVED***

	return nil
***REMOVED***

// processInternalRaftRequest proposes a value to be appended to the raft log.
// It calls Propose() on etcd/raft, which calls back into the raft FSM,
// which then sends a message to each of the participating nodes
// in the raft group to apply a log entry and then waits for it to be applied
// on this node. It will block until the this node:
// 1. Gets the necessary replies back from the participating nodes and also performs the commit itself, or
// 2. There is an error, or
// 3. Until the raft node finalizes all the proposals on node shutdown.
func (n *Node) processInternalRaftRequest(ctx context.Context, r *api.InternalRaftRequest, cb func()) (proto.Message, error) ***REMOVED***
	n.stopMu.RLock()
	if !n.IsMember() ***REMOVED***
		n.stopMu.RUnlock()
		return nil, ErrStopped
	***REMOVED***
	n.waitProp.Add(1)
	defer n.waitProp.Done()
	n.stopMu.RUnlock()

	r.ID = n.reqIDGen.Next()

	// This must be derived from the context which is cancelled by stop()
	// to avoid a deadlock on shutdown.
	waitCtx, cancel := context.WithCancel(ctx)

	ch := n.wait.register(r.ID, cb, cancel)

	// Do this check after calling register to avoid a race.
	if atomic.LoadUint32(&n.signalledLeadership) != 1 ***REMOVED***
		log.G(ctx).Error("node is no longer leader, aborting propose")
		n.wait.cancel(r.ID)
		return nil, ErrLostLeadership
	***REMOVED***

	data, err := r.Marshal()
	if err != nil ***REMOVED***
		n.wait.cancel(r.ID)
		return nil, err
	***REMOVED***

	if len(data) > store.MaxTransactionBytes ***REMOVED***
		n.wait.cancel(r.ID)
		return nil, ErrRequestTooLarge
	***REMOVED***

	err = n.raftNode.Propose(waitCtx, data)
	if err != nil ***REMOVED***
		n.wait.cancel(r.ID)
		return nil, err
	***REMOVED***

	select ***REMOVED***
	case x, ok := <-ch:
		if !ok ***REMOVED***
			// Wait notification channel was closed. This should only happen if the wait was cancelled.
			log.G(ctx).Error("wait cancelled")
			if atomic.LoadUint32(&n.signalledLeadership) == 1 ***REMOVED***
				log.G(ctx).Error("wait cancelled but node is still a leader")
			***REMOVED***
			return nil, ErrLostLeadership
		***REMOVED***
		return x.(proto.Message), nil
	case <-waitCtx.Done():
		n.wait.cancel(r.ID)
		// If we can read from the channel, wait item was triggered. Otherwise it was cancelled.
		x, ok := <-ch
		if !ok ***REMOVED***
			log.G(ctx).WithError(waitCtx.Err()).Error("wait context cancelled")
			if atomic.LoadUint32(&n.signalledLeadership) == 1 ***REMOVED***
				log.G(ctx).Error("wait context cancelled but node is still a leader")
			***REMOVED***
			return nil, ErrLostLeadership
		***REMOVED***
		return x.(proto.Message), nil
	case <-ctx.Done():
		n.wait.cancel(r.ID)
		// if channel is closed, wait item was canceled, otherwise it was triggered
		x, ok := <-ch
		if !ok ***REMOVED***
			return nil, ctx.Err()
		***REMOVED***
		return x.(proto.Message), nil
	***REMOVED***
***REMOVED***

// configure sends a configuration change through consensus and
// then waits for it to be applied to the server. It will block
// until the change is performed or there is an error.
func (n *Node) configure(ctx context.Context, cc raftpb.ConfChange) error ***REMOVED***
	cc.ID = n.reqIDGen.Next()

	ctx, cancel := context.WithCancel(ctx)
	ch := n.wait.register(cc.ID, nil, cancel)

	if err := n.raftNode.ProposeConfChange(ctx, cc); err != nil ***REMOVED***
		n.wait.cancel(cc.ID)
		return err
	***REMOVED***

	select ***REMOVED***
	case x := <-ch:
		if err, ok := x.(error); ok ***REMOVED***
			return err
		***REMOVED***
		if x != nil ***REMOVED***
			log.G(ctx).Panic("raft: configuration change error, return type should always be error")
		***REMOVED***
		return nil
	case <-ctx.Done():
		n.wait.cancel(cc.ID)
		return ctx.Err()
	***REMOVED***
***REMOVED***

func (n *Node) processCommitted(ctx context.Context, entry raftpb.Entry) error ***REMOVED***
	// Process a normal entry
	if entry.Type == raftpb.EntryNormal && entry.Data != nil ***REMOVED***
		if err := n.processEntry(ctx, entry); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Process a configuration change (add/remove node)
	if entry.Type == raftpb.EntryConfChange ***REMOVED***
		n.processConfChange(ctx, entry)
	***REMOVED***

	n.appliedIndex = entry.Index
	return nil
***REMOVED***

func (n *Node) processEntry(ctx context.Context, entry raftpb.Entry) error ***REMOVED***
	r := &api.InternalRaftRequest***REMOVED******REMOVED***
	err := proto.Unmarshal(entry.Data, r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !n.wait.trigger(r.ID, r) ***REMOVED***
		// There was no wait on this ID, meaning we don't have a
		// transaction in progress that would be committed to the
		// memory store by the "trigger" call. This could mean that:
		// 1. Startup is in progress, and the raft WAL is being parsed,
		// processed and applied to the store, or
		// 2. Either a different node wrote this to raft,
		// or we wrote it before losing the leader
		// position and cancelling the transaction. This entry still needs
		// to be committed since other nodes have already committed it.
		// Create a new transaction to commit this entry.

		// It should not be possible for processInternalRaftRequest
		// to be running in this situation, but out of caution we
		// cancel any current invocations to avoid a deadlock.
		// TODO(anshul) This call is likely redundant, remove after consideration.
		n.wait.cancelAll()

		err := n.memoryStore.ApplyStoreActions(r.Action)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("failed to apply actions from raft")
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (n *Node) processConfChange(ctx context.Context, entry raftpb.Entry) ***REMOVED***
	var (
		err error
		cc  raftpb.ConfChange
	)

	if err := proto.Unmarshal(entry.Data, &cc); err != nil ***REMOVED***
		n.wait.trigger(cc.ID, err)
	***REMOVED***

	if err := n.cluster.ValidateConfigurationChange(cc); err != nil ***REMOVED***
		n.wait.trigger(cc.ID, err)
	***REMOVED***

	switch cc.Type ***REMOVED***
	case raftpb.ConfChangeAddNode:
		err = n.applyAddNode(cc)
	case raftpb.ConfChangeUpdateNode:
		err = n.applyUpdateNode(ctx, cc)
	case raftpb.ConfChangeRemoveNode:
		err = n.applyRemoveNode(ctx, cc)
	***REMOVED***

	if err != nil ***REMOVED***
		n.wait.trigger(cc.ID, err)
	***REMOVED***

	n.confState = *n.raftNode.ApplyConfChange(cc)
	n.wait.trigger(cc.ID, nil)
***REMOVED***

// applyAddNode is called when we receive a ConfChange
// from a member in the raft cluster, this adds a new
// node to the existing raft cluster
func (n *Node) applyAddNode(cc raftpb.ConfChange) error ***REMOVED***
	member := &api.RaftMember***REMOVED******REMOVED***
	err := proto.Unmarshal(cc.Context, member)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// ID must be non zero
	if member.RaftID == 0 ***REMOVED***
		return nil
	***REMOVED***

	return n.registerNode(member)
***REMOVED***

// applyUpdateNode is called when we receive a ConfChange from a member in the
// raft cluster which update the address of an existing node.
func (n *Node) applyUpdateNode(ctx context.Context, cc raftpb.ConfChange) error ***REMOVED***
	newMember := &api.RaftMember***REMOVED******REMOVED***
	err := proto.Unmarshal(cc.Context, newMember)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if newMember.RaftID == n.Config.ID ***REMOVED***
		return nil
	***REMOVED***
	if err := n.transport.UpdatePeer(newMember.RaftID, newMember.Addr); err != nil ***REMOVED***
		return err
	***REMOVED***
	return n.cluster.UpdateMember(newMember.RaftID, newMember)
***REMOVED***

// applyRemoveNode is called when we receive a ConfChange
// from a member in the raft cluster, this removes a node
// from the existing raft cluster
func (n *Node) applyRemoveNode(ctx context.Context, cc raftpb.ConfChange) (err error) ***REMOVED***
	// If the node from where the remove is issued is
	// a follower and the leader steps down, Campaign
	// to be the leader.

	if cc.NodeID == n.leader() && !n.isLeader() ***REMOVED***
		if err = n.raftNode.Campaign(ctx); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if cc.NodeID == n.Config.ID ***REMOVED***
		// wait for the commit ack to be sent before closing connection
		n.asyncTasks.Wait()

		n.NodeRemoved()
	***REMOVED*** else if err := n.transport.RemovePeer(cc.NodeID); err != nil ***REMOVED***
		return err
	***REMOVED***

	return n.cluster.RemoveMember(cc.NodeID)
***REMOVED***

// SubscribeLeadership returns channel to which events about leadership change
// will be sent in form of raft.LeadershipState. Also cancel func is returned -
// it should be called when listener is no longer interested in events.
func (n *Node) SubscribeLeadership() (q chan events.Event, cancel func()) ***REMOVED***
	return n.leadershipBroadcast.Watch()
***REMOVED***

// createConfigChangeEnts creates a series of Raft entries (i.e.
// EntryConfChange) to remove the set of given IDs from the cluster. The ID
// `self` is _not_ removed, even if present in the set.
// If `self` is not inside the given ids, it creates a Raft entry to add a
// default member with the given `self`.
func createConfigChangeEnts(ids []uint64, self uint64, term, index uint64) []raftpb.Entry ***REMOVED***
	var ents []raftpb.Entry
	next := index + 1
	found := false
	for _, id := range ids ***REMOVED***
		if id == self ***REMOVED***
			found = true
			continue
		***REMOVED***
		cc := &raftpb.ConfChange***REMOVED***
			Type:   raftpb.ConfChangeRemoveNode,
			NodeID: id,
		***REMOVED***
		data, err := cc.Marshal()
		if err != nil ***REMOVED***
			log.L.WithError(err).Panic("marshal configuration change should never fail")
		***REMOVED***
		e := raftpb.Entry***REMOVED***
			Type:  raftpb.EntryConfChange,
			Data:  data,
			Term:  term,
			Index: next,
		***REMOVED***
		ents = append(ents, e)
		next++
	***REMOVED***
	if !found ***REMOVED***
		node := &api.RaftMember***REMOVED***RaftID: self***REMOVED***
		meta, err := node.Marshal()
		if err != nil ***REMOVED***
			log.L.WithError(err).Panic("marshal member should never fail")
		***REMOVED***
		cc := &raftpb.ConfChange***REMOVED***
			Type:    raftpb.ConfChangeAddNode,
			NodeID:  self,
			Context: meta,
		***REMOVED***
		data, err := cc.Marshal()
		if err != nil ***REMOVED***
			log.L.WithError(err).Panic("marshal configuration change should never fail")
		***REMOVED***
		e := raftpb.Entry***REMOVED***
			Type:  raftpb.EntryConfChange,
			Data:  data,
			Term:  term,
			Index: next,
		***REMOVED***
		ents = append(ents, e)
	***REMOVED***
	return ents
***REMOVED***

// getIDs returns an ordered set of IDs included in the given snapshot and
// the entries. The given snapshot/entries can contain two kinds of
// ID-related entry:
// - ConfChangeAddNode, in which case the contained ID will be added into the set.
// - ConfChangeRemoveNode, in which case the contained ID will be removed from the set.
func getIDs(snap *raftpb.Snapshot, ents []raftpb.Entry) []uint64 ***REMOVED***
	ids := make(map[uint64]struct***REMOVED******REMOVED***)
	if snap != nil ***REMOVED***
		for _, id := range snap.Metadata.ConfState.Nodes ***REMOVED***
			ids[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	for _, e := range ents ***REMOVED***
		if e.Type != raftpb.EntryConfChange ***REMOVED***
			continue
		***REMOVED***
		if snap != nil && e.Index < snap.Metadata.Index ***REMOVED***
			continue
		***REMOVED***
		var cc raftpb.ConfChange
		if err := cc.Unmarshal(e.Data); err != nil ***REMOVED***
			log.L.WithError(err).Panic("unmarshal configuration change should never fail")
		***REMOVED***
		switch cc.Type ***REMOVED***
		case raftpb.ConfChangeAddNode:
			ids[cc.NodeID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		case raftpb.ConfChangeRemoveNode:
			delete(ids, cc.NodeID)
		case raftpb.ConfChangeUpdateNode:
			// do nothing
		default:
			log.L.Panic("ConfChange Type should be either ConfChangeAddNode, or ConfChangeRemoveNode, or ConfChangeUpdateNode!")
		***REMOVED***
	***REMOVED***
	var sids []uint64
	for id := range ids ***REMOVED***
		sids = append(sids, id)
	***REMOVED***
	return sids
***REMOVED***

func (n *Node) reqTimeout() time.Duration ***REMOVED***
	return 5*time.Second + 2*time.Duration(n.Config.ElectionTick)*n.opts.TickInterval
***REMOVED***
