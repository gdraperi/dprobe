// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package raft

import (
	"errors"

	pb "github.com/coreos/etcd/raft/raftpb"
	"golang.org/x/net/context"
)

type SnapshotStatus int

const (
	SnapshotFinish  SnapshotStatus = 1
	SnapshotFailure SnapshotStatus = 2
)

var (
	emptyState = pb.HardState***REMOVED******REMOVED***

	// ErrStopped is returned by methods on Nodes that have been stopped.
	ErrStopped = errors.New("raft: stopped")
)

// SoftState provides state that is useful for logging and debugging.
// The state is volatile and does not need to be persisted to the WAL.
type SoftState struct ***REMOVED***
	Lead      uint64 // must use atomic operations to access; keep 64-bit aligned.
	RaftState StateType
***REMOVED***

func (a *SoftState) equal(b *SoftState) bool ***REMOVED***
	return a.Lead == b.Lead && a.RaftState == b.RaftState
***REMOVED***

// Ready encapsulates the entries and messages that are ready to read,
// be saved to stable storage, committed or sent to other peers.
// All fields in Ready are read-only.
type Ready struct ***REMOVED***
	// The current volatile state of a Node.
	// SoftState will be nil if there is no update.
	// It is not required to consume or store SoftState.
	*SoftState

	// The current state of a Node to be saved to stable storage BEFORE
	// Messages are sent.
	// HardState will be equal to empty state if there is no update.
	pb.HardState

	// ReadStates can be used for node to serve linearizable read requests locally
	// when its applied index is greater than the index in ReadState.
	// Note that the readState will be returned when raft receives msgReadIndex.
	// The returned is only valid for the request that requested to read.
	ReadStates []ReadState

	// Entries specifies entries to be saved to stable storage BEFORE
	// Messages are sent.
	Entries []pb.Entry

	// Snapshot specifies the snapshot to be saved to stable storage.
	Snapshot pb.Snapshot

	// CommittedEntries specifies entries to be committed to a
	// store/state-machine. These have previously been committed to stable
	// store.
	CommittedEntries []pb.Entry

	// Messages specifies outbound messages to be sent AFTER Entries are
	// committed to stable storage.
	// If it contains a MsgSnap message, the application MUST report back to raft
	// when the snapshot has been received or has failed by calling ReportSnapshot.
	Messages []pb.Message

	// MustSync indicates whether the HardState and Entries must be synchronously
	// written to disk or if an asynchronous write is permissible.
	MustSync bool
***REMOVED***

func isHardStateEqual(a, b pb.HardState) bool ***REMOVED***
	return a.Term == b.Term && a.Vote == b.Vote && a.Commit == b.Commit
***REMOVED***

// IsEmptyHardState returns true if the given HardState is empty.
func IsEmptyHardState(st pb.HardState) bool ***REMOVED***
	return isHardStateEqual(st, emptyState)
***REMOVED***

// IsEmptySnap returns true if the given Snapshot is empty.
func IsEmptySnap(sp pb.Snapshot) bool ***REMOVED***
	return sp.Metadata.Index == 0
***REMOVED***

func (rd Ready) containsUpdates() bool ***REMOVED***
	return rd.SoftState != nil || !IsEmptyHardState(rd.HardState) ||
		!IsEmptySnap(rd.Snapshot) || len(rd.Entries) > 0 ||
		len(rd.CommittedEntries) > 0 || len(rd.Messages) > 0 || len(rd.ReadStates) != 0
***REMOVED***

// Node represents a node in a raft cluster.
type Node interface ***REMOVED***
	// Tick increments the internal logical clock for the Node by a single tick. Election
	// timeouts and heartbeat timeouts are in units of ticks.
	Tick()
	// Campaign causes the Node to transition to candidate state and start campaigning to become leader.
	Campaign(ctx context.Context) error
	// Propose proposes that data be appended to the log.
	Propose(ctx context.Context, data []byte) error
	// ProposeConfChange proposes config change.
	// At most one ConfChange can be in the process of going through consensus.
	// Application needs to call ApplyConfChange when applying EntryConfChange type entry.
	ProposeConfChange(ctx context.Context, cc pb.ConfChange) error
	// Step advances the state machine using the given message. ctx.Err() will be returned, if any.
	Step(ctx context.Context, msg pb.Message) error

	// Ready returns a channel that returns the current point-in-time state.
	// Users of the Node must call Advance after retrieving the state returned by Ready.
	//
	// NOTE: No committed entries from the next Ready may be applied until all committed entries
	// and snapshots from the previous one have finished.
	Ready() <-chan Ready

	// Advance notifies the Node that the application has saved progress up to the last Ready.
	// It prepares the node to return the next available Ready.
	//
	// The application should generally call Advance after it applies the entries in last Ready.
	//
	// However, as an optimization, the application may call Advance while it is applying the
	// commands. For example. when the last Ready contains a snapshot, the application might take
	// a long time to apply the snapshot data. To continue receiving Ready without blocking raft
	// progress, it can call Advance before finishing applying the last ready.
	Advance()
	// ApplyConfChange applies config change to the local node.
	// Returns an opaque ConfState protobuf which must be recorded
	// in snapshots. Will never return nil; it returns a pointer only
	// to match MemoryStorage.Compact.
	ApplyConfChange(cc pb.ConfChange) *pb.ConfState

	// TransferLeadership attempts to transfer leadership to the given transferee.
	TransferLeadership(ctx context.Context, lead, transferee uint64)

	// ReadIndex request a read state. The read state will be set in the ready.
	// Read state has a read index. Once the application advances further than the read
	// index, any linearizable read requests issued before the read request can be
	// processed safely. The read state will have the same rctx attached.
	ReadIndex(ctx context.Context, rctx []byte) error

	// Status returns the current status of the raft state machine.
	Status() Status
	// ReportUnreachable reports the given node is not reachable for the last send.
	ReportUnreachable(id uint64)
	// ReportSnapshot reports the status of the sent snapshot.
	ReportSnapshot(id uint64, status SnapshotStatus)
	// Stop performs any necessary termination of the Node.
	Stop()
***REMOVED***

type Peer struct ***REMOVED***
	ID      uint64
	Context []byte
***REMOVED***

// StartNode returns a new Node given configuration and a list of raft peers.
// It appends a ConfChangeAddNode entry for each given peer to the initial log.
func StartNode(c *Config, peers []Peer) Node ***REMOVED***
	r := newRaft(c)
	// become the follower at term 1 and apply initial configuration
	// entries of term 1
	r.becomeFollower(1, None)
	for _, peer := range peers ***REMOVED***
		cc := pb.ConfChange***REMOVED***Type: pb.ConfChangeAddNode, NodeID: peer.ID, Context: peer.Context***REMOVED***
		d, err := cc.Marshal()
		if err != nil ***REMOVED***
			panic("unexpected marshal error")
		***REMOVED***
		e := pb.Entry***REMOVED***Type: pb.EntryConfChange, Term: 1, Index: r.raftLog.lastIndex() + 1, Data: d***REMOVED***
		r.raftLog.append(e)
	***REMOVED***
	// Mark these initial entries as committed.
	// TODO(bdarnell): These entries are still unstable; do we need to preserve
	// the invariant that committed < unstable?
	r.raftLog.committed = r.raftLog.lastIndex()
	// Now apply them, mainly so that the application can call Campaign
	// immediately after StartNode in tests. Note that these nodes will
	// be added to raft twice: here and when the application's Ready
	// loop calls ApplyConfChange. The calls to addNode must come after
	// all calls to raftLog.append so progress.next is set after these
	// bootstrapping entries (it is an error if we try to append these
	// entries since they have already been committed).
	// We do not set raftLog.applied so the application will be able
	// to observe all conf changes via Ready.CommittedEntries.
	for _, peer := range peers ***REMOVED***
		r.addNode(peer.ID)
	***REMOVED***

	n := newNode()
	n.logger = c.Logger
	go n.run(r)
	return &n
***REMOVED***

// RestartNode is similar to StartNode but does not take a list of peers.
// The current membership of the cluster will be restored from the Storage.
// If the caller has an existing state machine, pass in the last log index that
// has been applied to it; otherwise use zero.
func RestartNode(c *Config) Node ***REMOVED***
	r := newRaft(c)

	n := newNode()
	n.logger = c.Logger
	go n.run(r)
	return &n
***REMOVED***

// node is the canonical implementation of the Node interface
type node struct ***REMOVED***
	propc      chan pb.Message
	recvc      chan pb.Message
	confc      chan pb.ConfChange
	confstatec chan pb.ConfState
	readyc     chan Ready
	advancec   chan struct***REMOVED******REMOVED***
	tickc      chan struct***REMOVED******REMOVED***
	done       chan struct***REMOVED******REMOVED***
	stop       chan struct***REMOVED******REMOVED***
	status     chan chan Status

	logger Logger
***REMOVED***

func newNode() node ***REMOVED***
	return node***REMOVED***
		propc:      make(chan pb.Message),
		recvc:      make(chan pb.Message),
		confc:      make(chan pb.ConfChange),
		confstatec: make(chan pb.ConfState),
		readyc:     make(chan Ready),
		advancec:   make(chan struct***REMOVED******REMOVED***),
		// make tickc a buffered chan, so raft node can buffer some ticks when the node
		// is busy processing raft messages. Raft node will resume process buffered
		// ticks when it becomes idle.
		tickc:  make(chan struct***REMOVED******REMOVED***, 128),
		done:   make(chan struct***REMOVED******REMOVED***),
		stop:   make(chan struct***REMOVED******REMOVED***),
		status: make(chan chan Status),
	***REMOVED***
***REMOVED***

func (n *node) Stop() ***REMOVED***
	select ***REMOVED***
	case n.stop <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
		// Not already stopped, so trigger it
	case <-n.done:
		// Node has already been stopped - no need to do anything
		return
	***REMOVED***
	// Block until the stop has been acknowledged by run()
	<-n.done
***REMOVED***

func (n *node) run(r *raft) ***REMOVED***
	var propc chan pb.Message
	var readyc chan Ready
	var advancec chan struct***REMOVED******REMOVED***
	var prevLastUnstablei, prevLastUnstablet uint64
	var havePrevLastUnstablei bool
	var prevSnapi uint64
	var rd Ready

	lead := None
	prevSoftSt := r.softState()
	prevHardSt := emptyState

	for ***REMOVED***
		if advancec != nil ***REMOVED***
			readyc = nil
		***REMOVED*** else ***REMOVED***
			rd = newReady(r, prevSoftSt, prevHardSt)
			if rd.containsUpdates() ***REMOVED***
				readyc = n.readyc
			***REMOVED*** else ***REMOVED***
				readyc = nil
			***REMOVED***
		***REMOVED***

		if lead != r.lead ***REMOVED***
			if r.hasLeader() ***REMOVED***
				if lead == None ***REMOVED***
					r.logger.Infof("raft.node: %x elected leader %x at term %d", r.id, r.lead, r.Term)
				***REMOVED*** else ***REMOVED***
					r.logger.Infof("raft.node: %x changed leader from %x to %x at term %d", r.id, lead, r.lead, r.Term)
				***REMOVED***
				propc = n.propc
			***REMOVED*** else ***REMOVED***
				r.logger.Infof("raft.node: %x lost leader %x at term %d", r.id, lead, r.Term)
				propc = nil
			***REMOVED***
			lead = r.lead
		***REMOVED***

		select ***REMOVED***
		// TODO: maybe buffer the config propose if there exists one (the way
		// described in raft dissertation)
		// Currently it is dropped in Step silently.
		case m := <-propc:
			m.From = r.id
			r.Step(m)
		case m := <-n.recvc:
			// filter out response message from unknown From.
			if _, ok := r.prs[m.From]; ok || !IsResponseMsg(m.Type) ***REMOVED***
				r.Step(m) // raft never returns an error
			***REMOVED***
		case cc := <-n.confc:
			if cc.NodeID == None ***REMOVED***
				r.resetPendingConf()
				select ***REMOVED***
				case n.confstatec <- pb.ConfState***REMOVED***Nodes: r.nodes()***REMOVED***:
				case <-n.done:
				***REMOVED***
				break
			***REMOVED***
			switch cc.Type ***REMOVED***
			case pb.ConfChangeAddNode:
				r.addNode(cc.NodeID)
			case pb.ConfChangeRemoveNode:
				// block incoming proposal when local node is
				// removed
				if cc.NodeID == r.id ***REMOVED***
					propc = nil
				***REMOVED***
				r.removeNode(cc.NodeID)
			case pb.ConfChangeUpdateNode:
				r.resetPendingConf()
			default:
				panic("unexpected conf type")
			***REMOVED***
			select ***REMOVED***
			case n.confstatec <- pb.ConfState***REMOVED***Nodes: r.nodes()***REMOVED***:
			case <-n.done:
			***REMOVED***
		case <-n.tickc:
			r.tick()
		case readyc <- rd:
			if rd.SoftState != nil ***REMOVED***
				prevSoftSt = rd.SoftState
			***REMOVED***
			if len(rd.Entries) > 0 ***REMOVED***
				prevLastUnstablei = rd.Entries[len(rd.Entries)-1].Index
				prevLastUnstablet = rd.Entries[len(rd.Entries)-1].Term
				havePrevLastUnstablei = true
			***REMOVED***
			if !IsEmptyHardState(rd.HardState) ***REMOVED***
				prevHardSt = rd.HardState
			***REMOVED***
			if !IsEmptySnap(rd.Snapshot) ***REMOVED***
				prevSnapi = rd.Snapshot.Metadata.Index
			***REMOVED***

			r.msgs = nil
			r.readStates = nil
			advancec = n.advancec
		case <-advancec:
			if prevHardSt.Commit != 0 ***REMOVED***
				r.raftLog.appliedTo(prevHardSt.Commit)
			***REMOVED***
			if havePrevLastUnstablei ***REMOVED***
				r.raftLog.stableTo(prevLastUnstablei, prevLastUnstablet)
				havePrevLastUnstablei = false
			***REMOVED***
			r.raftLog.stableSnapTo(prevSnapi)
			advancec = nil
		case c := <-n.status:
			c <- getStatus(r)
		case <-n.stop:
			close(n.done)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Tick increments the internal logical clock for this Node. Election timeouts
// and heartbeat timeouts are in units of ticks.
func (n *node) Tick() ***REMOVED***
	select ***REMOVED***
	case n.tickc <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
	case <-n.done:
	default:
		n.logger.Warningf("A tick missed to fire. Node blocks too long!")
	***REMOVED***
***REMOVED***

func (n *node) Campaign(ctx context.Context) error ***REMOVED*** return n.step(ctx, pb.Message***REMOVED***Type: pb.MsgHup***REMOVED***) ***REMOVED***

func (n *node) Propose(ctx context.Context, data []byte) error ***REMOVED***
	return n.step(ctx, pb.Message***REMOVED***Type: pb.MsgProp, Entries: []pb.Entry***REMOVED******REMOVED***Data: data***REMOVED******REMOVED******REMOVED***)
***REMOVED***

func (n *node) Step(ctx context.Context, m pb.Message) error ***REMOVED***
	// ignore unexpected local messages receiving over network
	if IsLocalMsg(m.Type) ***REMOVED***
		// TODO: return an error?
		return nil
	***REMOVED***
	return n.step(ctx, m)
***REMOVED***

func (n *node) ProposeConfChange(ctx context.Context, cc pb.ConfChange) error ***REMOVED***
	data, err := cc.Marshal()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return n.Step(ctx, pb.Message***REMOVED***Type: pb.MsgProp, Entries: []pb.Entry***REMOVED******REMOVED***Type: pb.EntryConfChange, Data: data***REMOVED******REMOVED******REMOVED***)
***REMOVED***

// Step advances the state machine using msgs. The ctx.Err() will be returned,
// if any.
func (n *node) step(ctx context.Context, m pb.Message) error ***REMOVED***
	ch := n.recvc
	if m.Type == pb.MsgProp ***REMOVED***
		ch = n.propc
	***REMOVED***

	select ***REMOVED***
	case ch <- m:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-n.done:
		return ErrStopped
	***REMOVED***
***REMOVED***

func (n *node) Ready() <-chan Ready ***REMOVED*** return n.readyc ***REMOVED***

func (n *node) Advance() ***REMOVED***
	select ***REMOVED***
	case n.advancec <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
	case <-n.done:
	***REMOVED***
***REMOVED***

func (n *node) ApplyConfChange(cc pb.ConfChange) *pb.ConfState ***REMOVED***
	var cs pb.ConfState
	select ***REMOVED***
	case n.confc <- cc:
	case <-n.done:
	***REMOVED***
	select ***REMOVED***
	case cs = <-n.confstatec:
	case <-n.done:
	***REMOVED***
	return &cs
***REMOVED***

func (n *node) Status() Status ***REMOVED***
	c := make(chan Status)
	select ***REMOVED***
	case n.status <- c:
		return <-c
	case <-n.done:
		return Status***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

func (n *node) ReportUnreachable(id uint64) ***REMOVED***
	select ***REMOVED***
	case n.recvc <- pb.Message***REMOVED***Type: pb.MsgUnreachable, From: id***REMOVED***:
	case <-n.done:
	***REMOVED***
***REMOVED***

func (n *node) ReportSnapshot(id uint64, status SnapshotStatus) ***REMOVED***
	rej := status == SnapshotFailure

	select ***REMOVED***
	case n.recvc <- pb.Message***REMOVED***Type: pb.MsgSnapStatus, From: id, Reject: rej***REMOVED***:
	case <-n.done:
	***REMOVED***
***REMOVED***

func (n *node) TransferLeadership(ctx context.Context, lead, transferee uint64) ***REMOVED***
	select ***REMOVED***
	// manually set 'from' and 'to', so that leader can voluntarily transfers its leadership
	case n.recvc <- pb.Message***REMOVED***Type: pb.MsgTransferLeader, From: transferee, To: lead***REMOVED***:
	case <-n.done:
	case <-ctx.Done():
	***REMOVED***
***REMOVED***

func (n *node) ReadIndex(ctx context.Context, rctx []byte) error ***REMOVED***
	return n.step(ctx, pb.Message***REMOVED***Type: pb.MsgReadIndex, Entries: []pb.Entry***REMOVED******REMOVED***Data: rctx***REMOVED******REMOVED******REMOVED***)
***REMOVED***

func newReady(r *raft, prevSoftSt *SoftState, prevHardSt pb.HardState) Ready ***REMOVED***
	rd := Ready***REMOVED***
		Entries:          r.raftLog.unstableEntries(),
		CommittedEntries: r.raftLog.nextEnts(),
		Messages:         r.msgs,
	***REMOVED***
	if softSt := r.softState(); !softSt.equal(prevSoftSt) ***REMOVED***
		rd.SoftState = softSt
	***REMOVED***
	if hardSt := r.hardState(); !isHardStateEqual(hardSt, prevHardSt) ***REMOVED***
		rd.HardState = hardSt
	***REMOVED***
	if r.raftLog.unstable.snapshot != nil ***REMOVED***
		rd.Snapshot = *r.raftLog.unstable.snapshot
	***REMOVED***
	if len(r.readStates) != 0 ***REMOVED***
		rd.ReadStates = r.readStates
	***REMOVED***
	rd.MustSync = MustSync(rd.HardState, prevHardSt, len(rd.Entries))
	return rd
***REMOVED***

// MustSync returns true if the hard state and count of Raft entries indicate
// that a synchronous write to persistent storage is required.
func MustSync(st, prevst pb.HardState, entsnum int) bool ***REMOVED***
	// Persistent state on all servers:
	// (Updated on stable storage before responding to RPCs)
	// currentTerm
	// votedFor
	// log entries[]
	return entsnum != 0 || st.Vote != prevst.Vote || st.Term != prevst.Term
***REMOVED***
