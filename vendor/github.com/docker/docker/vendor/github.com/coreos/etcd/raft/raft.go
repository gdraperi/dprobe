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
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	pb "github.com/coreos/etcd/raft/raftpb"
)

// None is a placeholder node ID used when there is no leader.
const None uint64 = 0
const noLimit = math.MaxUint64

// Possible values for StateType.
const (
	StateFollower StateType = iota
	StateCandidate
	StateLeader
	StatePreCandidate
	numStates
)

type ReadOnlyOption int

const (
	// ReadOnlySafe guarantees the linearizability of the read only request by
	// communicating with the quorum. It is the default and suggested option.
	ReadOnlySafe ReadOnlyOption = iota
	// ReadOnlyLeaseBased ensures linearizability of the read only request by
	// relying on the leader lease. It can be affected by clock drift.
	// If the clock drift is unbounded, leader might keep the lease longer than it
	// should (clock can move backward/pause without any bound). ReadIndex is not safe
	// in that case.
	ReadOnlyLeaseBased
)

// Possible values for CampaignType
const (
	// campaignPreElection represents the first phase of a normal election when
	// Config.PreVote is true.
	campaignPreElection CampaignType = "CampaignPreElection"
	// campaignElection represents a normal (time-based) election (the second phase
	// of the election when Config.PreVote is true).
	campaignElection CampaignType = "CampaignElection"
	// campaignTransfer represents the type of leader transfer
	campaignTransfer CampaignType = "CampaignTransfer"
)

// lockedRand is a small wrapper around rand.Rand to provide
// synchronization. Only the methods needed by the code are exposed
// (e.g. Intn).
type lockedRand struct ***REMOVED***
	mu   sync.Mutex
	rand *rand.Rand
***REMOVED***

func (r *lockedRand) Intn(n int) int ***REMOVED***
	r.mu.Lock()
	v := r.rand.Intn(n)
	r.mu.Unlock()
	return v
***REMOVED***

var globalRand = &lockedRand***REMOVED***
	rand: rand.New(rand.NewSource(time.Now().UnixNano())),
***REMOVED***

// CampaignType represents the type of campaigning
// the reason we use the type of string instead of uint64
// is because it's simpler to compare and fill in raft entries
type CampaignType string

// StateType represents the role of a node in a cluster.
type StateType uint64

var stmap = [...]string***REMOVED***
	"StateFollower",
	"StateCandidate",
	"StateLeader",
	"StatePreCandidate",
***REMOVED***

func (st StateType) String() string ***REMOVED***
	return stmap[uint64(st)]
***REMOVED***

// Config contains the parameters to start a raft.
type Config struct ***REMOVED***
	// ID is the identity of the local raft. ID cannot be 0.
	ID uint64

	// peers contains the IDs of all nodes (including self) in the raft cluster. It
	// should only be set when starting a new raft cluster. Restarting raft from
	// previous configuration will panic if peers is set. peer is private and only
	// used for testing right now.
	peers []uint64

	// ElectionTick is the number of Node.Tick invocations that must pass between
	// elections. That is, if a follower does not receive any message from the
	// leader of current term before ElectionTick has elapsed, it will become
	// candidate and start an election. ElectionTick must be greater than
	// HeartbeatTick. We suggest ElectionTick = 10 * HeartbeatTick to avoid
	// unnecessary leader switching.
	ElectionTick int
	// HeartbeatTick is the number of Node.Tick invocations that must pass between
	// heartbeats. That is, a leader sends heartbeat messages to maintain its
	// leadership every HeartbeatTick ticks.
	HeartbeatTick int

	// Storage is the storage for raft. raft generates entries and states to be
	// stored in storage. raft reads the persisted entries and states out of
	// Storage when it needs. raft reads out the previous state and configuration
	// out of storage when restarting.
	Storage Storage
	// Applied is the last applied index. It should only be set when restarting
	// raft. raft will not return entries to the application smaller or equal to
	// Applied. If Applied is unset when restarting, raft might return previous
	// applied entries. This is a very application dependent configuration.
	Applied uint64

	// MaxSizePerMsg limits the max size of each append message. Smaller value
	// lowers the raft recovery cost(initial probing and message lost during normal
	// operation). On the other side, it might affect the throughput during normal
	// replication. Note: math.MaxUint64 for unlimited, 0 for at most one entry per
	// message.
	MaxSizePerMsg uint64
	// MaxInflightMsgs limits the max number of in-flight append messages during
	// optimistic replication phase. The application transportation layer usually
	// has its own sending buffer over TCP/UDP. Setting MaxInflightMsgs to avoid
	// overflowing that sending buffer. TODO (xiangli): feedback to application to
	// limit the proposal rate?
	MaxInflightMsgs int

	// CheckQuorum specifies if the leader should check quorum activity. Leader
	// steps down when quorum is not active for an electionTimeout.
	CheckQuorum bool

	// PreVote enables the Pre-Vote algorithm described in raft thesis section
	// 9.6. This prevents disruption when a node that has been partitioned away
	// rejoins the cluster.
	PreVote bool

	// ReadOnlyOption specifies how the read only request is processed.
	//
	// ReadOnlySafe guarantees the linearizability of the read only request by
	// communicating with the quorum. It is the default and suggested option.
	//
	// ReadOnlyLeaseBased ensures linearizability of the read only request by
	// relying on the leader lease. It can be affected by clock drift.
	// If the clock drift is unbounded, leader might keep the lease longer than it
	// should (clock can move backward/pause without any bound). ReadIndex is not safe
	// in that case.
	ReadOnlyOption ReadOnlyOption

	// Logger is the logger used for raft log. For multinode which can host
	// multiple raft group, each raft group can have its own logger
	Logger Logger
***REMOVED***

func (c *Config) validate() error ***REMOVED***
	if c.ID == None ***REMOVED***
		return errors.New("cannot use none as id")
	***REMOVED***

	if c.HeartbeatTick <= 0 ***REMOVED***
		return errors.New("heartbeat tick must be greater than 0")
	***REMOVED***

	if c.ElectionTick <= c.HeartbeatTick ***REMOVED***
		return errors.New("election tick must be greater than heartbeat tick")
	***REMOVED***

	if c.Storage == nil ***REMOVED***
		return errors.New("storage cannot be nil")
	***REMOVED***

	if c.MaxInflightMsgs <= 0 ***REMOVED***
		return errors.New("max inflight messages must be greater than 0")
	***REMOVED***

	if c.Logger == nil ***REMOVED***
		c.Logger = raftLogger
	***REMOVED***

	return nil
***REMOVED***

type raft struct ***REMOVED***
	id uint64

	Term uint64
	Vote uint64

	readStates []ReadState

	// the log
	raftLog *raftLog

	maxInflight int
	maxMsgSize  uint64
	prs         map[uint64]*Progress

	state StateType

	votes map[uint64]bool

	msgs []pb.Message

	// the leader id
	lead uint64
	// leadTransferee is id of the leader transfer target when its value is not zero.
	// Follow the procedure defined in raft thesis 3.10.
	leadTransferee uint64
	// New configuration is ignored if there exists unapplied configuration.
	pendingConf bool

	readOnly *readOnly

	// number of ticks since it reached last electionTimeout when it is leader
	// or candidate.
	// number of ticks since it reached last electionTimeout or received a
	// valid message from current leader when it is a follower.
	electionElapsed int

	// number of ticks since it reached last heartbeatTimeout.
	// only leader keeps heartbeatElapsed.
	heartbeatElapsed int

	checkQuorum bool
	preVote     bool

	heartbeatTimeout int
	electionTimeout  int
	// randomizedElectionTimeout is a random number between
	// [electiontimeout, 2 * electiontimeout - 1]. It gets reset
	// when raft changes its state to follower or candidate.
	randomizedElectionTimeout int

	tick func()
	step stepFunc

	logger Logger
***REMOVED***

func newRaft(c *Config) *raft ***REMOVED***
	if err := c.validate(); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
	raftlog := newLog(c.Storage, c.Logger)
	hs, cs, err := c.Storage.InitialState()
	if err != nil ***REMOVED***
		panic(err) // TODO(bdarnell)
	***REMOVED***
	peers := c.peers
	if len(cs.Nodes) > 0 ***REMOVED***
		if len(peers) > 0 ***REMOVED***
			// TODO(bdarnell): the peers argument is always nil except in
			// tests; the argument should be removed and these tests should be
			// updated to specify their nodes through a snapshot.
			panic("cannot specify both newRaft(peers) and ConfState.Nodes)")
		***REMOVED***
		peers = cs.Nodes
	***REMOVED***
	r := &raft***REMOVED***
		id:               c.ID,
		lead:             None,
		raftLog:          raftlog,
		maxMsgSize:       c.MaxSizePerMsg,
		maxInflight:      c.MaxInflightMsgs,
		prs:              make(map[uint64]*Progress),
		electionTimeout:  c.ElectionTick,
		heartbeatTimeout: c.HeartbeatTick,
		logger:           c.Logger,
		checkQuorum:      c.CheckQuorum,
		preVote:          c.PreVote,
		readOnly:         newReadOnly(c.ReadOnlyOption),
	***REMOVED***
	for _, p := range peers ***REMOVED***
		r.prs[p] = &Progress***REMOVED***Next: 1, ins: newInflights(r.maxInflight)***REMOVED***
	***REMOVED***
	if !isHardStateEqual(hs, emptyState) ***REMOVED***
		r.loadState(hs)
	***REMOVED***
	if c.Applied > 0 ***REMOVED***
		raftlog.appliedTo(c.Applied)
	***REMOVED***
	r.becomeFollower(r.Term, None)

	var nodesStrs []string
	for _, n := range r.nodes() ***REMOVED***
		nodesStrs = append(nodesStrs, fmt.Sprintf("%x", n))
	***REMOVED***

	r.logger.Infof("newRaft %x [peers: [%s], term: %d, commit: %d, applied: %d, lastindex: %d, lastterm: %d]",
		r.id, strings.Join(nodesStrs, ","), r.Term, r.raftLog.committed, r.raftLog.applied, r.raftLog.lastIndex(), r.raftLog.lastTerm())
	return r
***REMOVED***

func (r *raft) hasLeader() bool ***REMOVED*** return r.lead != None ***REMOVED***

func (r *raft) softState() *SoftState ***REMOVED*** return &SoftState***REMOVED***Lead: r.lead, RaftState: r.state***REMOVED*** ***REMOVED***

func (r *raft) hardState() pb.HardState ***REMOVED***
	return pb.HardState***REMOVED***
		Term:   r.Term,
		Vote:   r.Vote,
		Commit: r.raftLog.committed,
	***REMOVED***
***REMOVED***

func (r *raft) quorum() int ***REMOVED*** return len(r.prs)/2 + 1 ***REMOVED***

func (r *raft) nodes() []uint64 ***REMOVED***
	nodes := make([]uint64, 0, len(r.prs))
	for id := range r.prs ***REMOVED***
		nodes = append(nodes, id)
	***REMOVED***
	sort.Sort(uint64Slice(nodes))
	return nodes
***REMOVED***

// send persists state to stable storage and then sends to its mailbox.
func (r *raft) send(m pb.Message) ***REMOVED***
	m.From = r.id
	if m.Type == pb.MsgVote || m.Type == pb.MsgPreVote ***REMOVED***
		if m.Term == 0 ***REMOVED***
			// PreVote RPCs are sent at a term other than our actual term, so the code
			// that sends these messages is responsible for setting the term.
			panic(fmt.Sprintf("term should be set when sending %s", m.Type))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if m.Term != 0 ***REMOVED***
			panic(fmt.Sprintf("term should not be set when sending %s (was %d)", m.Type, m.Term))
		***REMOVED***
		// do not attach term to MsgProp, MsgReadIndex
		// proposals are a way to forward to the leader and
		// should be treated as local message.
		// MsgReadIndex is also forwarded to leader.
		if m.Type != pb.MsgProp && m.Type != pb.MsgReadIndex ***REMOVED***
			m.Term = r.Term
		***REMOVED***
	***REMOVED***
	r.msgs = append(r.msgs, m)
***REMOVED***

// sendAppend sends RPC, with entries to the given peer.
func (r *raft) sendAppend(to uint64) ***REMOVED***
	pr := r.prs[to]
	if pr.IsPaused() ***REMOVED***
		return
	***REMOVED***
	m := pb.Message***REMOVED******REMOVED***
	m.To = to

	term, errt := r.raftLog.term(pr.Next - 1)
	ents, erre := r.raftLog.entries(pr.Next, r.maxMsgSize)

	if errt != nil || erre != nil ***REMOVED*** // send snapshot if we failed to get term or entries
		if !pr.RecentActive ***REMOVED***
			r.logger.Debugf("ignore sending snapshot to %x since it is not recently active", to)
			return
		***REMOVED***

		m.Type = pb.MsgSnap
		snapshot, err := r.raftLog.snapshot()
		if err != nil ***REMOVED***
			if err == ErrSnapshotTemporarilyUnavailable ***REMOVED***
				r.logger.Debugf("%x failed to send snapshot to %x because snapshot is temporarily unavailable", r.id, to)
				return
			***REMOVED***
			panic(err) // TODO(bdarnell)
		***REMOVED***
		if IsEmptySnap(snapshot) ***REMOVED***
			panic("need non-empty snapshot")
		***REMOVED***
		m.Snapshot = snapshot
		sindex, sterm := snapshot.Metadata.Index, snapshot.Metadata.Term
		r.logger.Debugf("%x [firstindex: %d, commit: %d] sent snapshot[index: %d, term: %d] to %x [%s]",
			r.id, r.raftLog.firstIndex(), r.raftLog.committed, sindex, sterm, to, pr)
		pr.becomeSnapshot(sindex)
		r.logger.Debugf("%x paused sending replication messages to %x [%s]", r.id, to, pr)
	***REMOVED*** else ***REMOVED***
		m.Type = pb.MsgApp
		m.Index = pr.Next - 1
		m.LogTerm = term
		m.Entries = ents
		m.Commit = r.raftLog.committed
		if n := len(m.Entries); n != 0 ***REMOVED***
			switch pr.State ***REMOVED***
			// optimistically increase the next when in ProgressStateReplicate
			case ProgressStateReplicate:
				last := m.Entries[n-1].Index
				pr.optimisticUpdate(last)
				pr.ins.add(last)
			case ProgressStateProbe:
				pr.pause()
			default:
				r.logger.Panicf("%x is sending append in unhandled state %s", r.id, pr.State)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	r.send(m)
***REMOVED***

// sendHeartbeat sends an empty MsgApp
func (r *raft) sendHeartbeat(to uint64, ctx []byte) ***REMOVED***
	// Attach the commit as min(to.matched, r.committed).
	// When the leader sends out heartbeat message,
	// the receiver(follower) might not be matched with the leader
	// or it might not have all the committed entries.
	// The leader MUST NOT forward the follower's commit to
	// an unmatched index.
	commit := min(r.prs[to].Match, r.raftLog.committed)
	m := pb.Message***REMOVED***
		To:      to,
		Type:    pb.MsgHeartbeat,
		Commit:  commit,
		Context: ctx,
	***REMOVED***

	r.send(m)
***REMOVED***

// bcastAppend sends RPC, with entries to all peers that are not up-to-date
// according to the progress recorded in r.prs.
func (r *raft) bcastAppend() ***REMOVED***
	for id := range r.prs ***REMOVED***
		if id == r.id ***REMOVED***
			continue
		***REMOVED***
		r.sendAppend(id)
	***REMOVED***
***REMOVED***

// bcastHeartbeat sends RPC, without entries to all the peers.
func (r *raft) bcastHeartbeat() ***REMOVED***
	lastCtx := r.readOnly.lastPendingRequestCtx()
	if len(lastCtx) == 0 ***REMOVED***
		r.bcastHeartbeatWithCtx(nil)
	***REMOVED*** else ***REMOVED***
		r.bcastHeartbeatWithCtx([]byte(lastCtx))
	***REMOVED***
***REMOVED***

func (r *raft) bcastHeartbeatWithCtx(ctx []byte) ***REMOVED***
	for id := range r.prs ***REMOVED***
		if id == r.id ***REMOVED***
			continue
		***REMOVED***
		r.sendHeartbeat(id, ctx)
	***REMOVED***
***REMOVED***

// maybeCommit attempts to advance the commit index. Returns true if
// the commit index changed (in which case the caller should call
// r.bcastAppend).
func (r *raft) maybeCommit() bool ***REMOVED***
	// TODO(bmizerany): optimize.. Currently naive
	mis := make(uint64Slice, 0, len(r.prs))
	for id := range r.prs ***REMOVED***
		mis = append(mis, r.prs[id].Match)
	***REMOVED***
	sort.Sort(sort.Reverse(mis))
	mci := mis[r.quorum()-1]
	return r.raftLog.maybeCommit(mci, r.Term)
***REMOVED***

func (r *raft) reset(term uint64) ***REMOVED***
	if r.Term != term ***REMOVED***
		r.Term = term
		r.Vote = None
	***REMOVED***
	r.lead = None

	r.electionElapsed = 0
	r.heartbeatElapsed = 0
	r.resetRandomizedElectionTimeout()

	r.abortLeaderTransfer()

	r.votes = make(map[uint64]bool)
	for id := range r.prs ***REMOVED***
		r.prs[id] = &Progress***REMOVED***Next: r.raftLog.lastIndex() + 1, ins: newInflights(r.maxInflight)***REMOVED***
		if id == r.id ***REMOVED***
			r.prs[id].Match = r.raftLog.lastIndex()
		***REMOVED***
	***REMOVED***
	r.pendingConf = false
	r.readOnly = newReadOnly(r.readOnly.option)
***REMOVED***

func (r *raft) appendEntry(es ...pb.Entry) ***REMOVED***
	li := r.raftLog.lastIndex()
	for i := range es ***REMOVED***
		es[i].Term = r.Term
		es[i].Index = li + 1 + uint64(i)
	***REMOVED***
	r.raftLog.append(es...)
	r.prs[r.id].maybeUpdate(r.raftLog.lastIndex())
	// Regardless of maybeCommit's return, our caller will call bcastAppend.
	r.maybeCommit()
***REMOVED***

// tickElection is run by followers and candidates after r.electionTimeout.
func (r *raft) tickElection() ***REMOVED***
	r.electionElapsed++

	if r.promotable() && r.pastElectionTimeout() ***REMOVED***
		r.electionElapsed = 0
		r.Step(pb.Message***REMOVED***From: r.id, Type: pb.MsgHup***REMOVED***)
	***REMOVED***
***REMOVED***

// tickHeartbeat is run by leaders to send a MsgBeat after r.heartbeatTimeout.
func (r *raft) tickHeartbeat() ***REMOVED***
	r.heartbeatElapsed++
	r.electionElapsed++

	if r.electionElapsed >= r.electionTimeout ***REMOVED***
		r.electionElapsed = 0
		if r.checkQuorum ***REMOVED***
			r.Step(pb.Message***REMOVED***From: r.id, Type: pb.MsgCheckQuorum***REMOVED***)
		***REMOVED***
		// If current leader cannot transfer leadership in electionTimeout, it becomes leader again.
		if r.state == StateLeader && r.leadTransferee != None ***REMOVED***
			r.abortLeaderTransfer()
		***REMOVED***
	***REMOVED***

	if r.state != StateLeader ***REMOVED***
		return
	***REMOVED***

	if r.heartbeatElapsed >= r.heartbeatTimeout ***REMOVED***
		r.heartbeatElapsed = 0
		r.Step(pb.Message***REMOVED***From: r.id, Type: pb.MsgBeat***REMOVED***)
	***REMOVED***
***REMOVED***

func (r *raft) becomeFollower(term uint64, lead uint64) ***REMOVED***
	r.step = stepFollower
	r.reset(term)
	r.tick = r.tickElection
	r.lead = lead
	r.state = StateFollower
	r.logger.Infof("%x became follower at term %d", r.id, r.Term)
***REMOVED***

func (r *raft) becomeCandidate() ***REMOVED***
	// TODO(xiangli) remove the panic when the raft implementation is stable
	if r.state == StateLeader ***REMOVED***
		panic("invalid transition [leader -> candidate]")
	***REMOVED***
	r.step = stepCandidate
	r.reset(r.Term + 1)
	r.tick = r.tickElection
	r.Vote = r.id
	r.state = StateCandidate
	r.logger.Infof("%x became candidate at term %d", r.id, r.Term)
***REMOVED***

func (r *raft) becomePreCandidate() ***REMOVED***
	// TODO(xiangli) remove the panic when the raft implementation is stable
	if r.state == StateLeader ***REMOVED***
		panic("invalid transition [leader -> pre-candidate]")
	***REMOVED***
	// Becoming a pre-candidate changes our step functions and state,
	// but doesn't change anything else. In particular it does not increase
	// r.Term or change r.Vote.
	r.step = stepCandidate
	r.tick = r.tickElection
	r.state = StatePreCandidate
	r.logger.Infof("%x became pre-candidate at term %d", r.id, r.Term)
***REMOVED***

func (r *raft) becomeLeader() ***REMOVED***
	// TODO(xiangli) remove the panic when the raft implementation is stable
	if r.state == StateFollower ***REMOVED***
		panic("invalid transition [follower -> leader]")
	***REMOVED***
	r.step = stepLeader
	r.reset(r.Term)
	r.tick = r.tickHeartbeat
	r.lead = r.id
	r.state = StateLeader
	ents, err := r.raftLog.entries(r.raftLog.committed+1, noLimit)
	if err != nil ***REMOVED***
		r.logger.Panicf("unexpected error getting uncommitted entries (%v)", err)
	***REMOVED***

	nconf := numOfPendingConf(ents)
	if nconf > 1 ***REMOVED***
		panic("unexpected multiple uncommitted config entry")
	***REMOVED***
	if nconf == 1 ***REMOVED***
		r.pendingConf = true
	***REMOVED***

	r.appendEntry(pb.Entry***REMOVED***Data: nil***REMOVED***)
	r.logger.Infof("%x became leader at term %d", r.id, r.Term)
***REMOVED***

func (r *raft) campaign(t CampaignType) ***REMOVED***
	var term uint64
	var voteMsg pb.MessageType
	if t == campaignPreElection ***REMOVED***
		r.becomePreCandidate()
		voteMsg = pb.MsgPreVote
		// PreVote RPCs are sent for the next term before we've incremented r.Term.
		term = r.Term + 1
	***REMOVED*** else ***REMOVED***
		r.becomeCandidate()
		voteMsg = pb.MsgVote
		term = r.Term
	***REMOVED***
	if r.quorum() == r.poll(r.id, voteRespMsgType(voteMsg), true) ***REMOVED***
		// We won the election after voting for ourselves (which must mean that
		// this is a single-node cluster). Advance to the next state.
		if t == campaignPreElection ***REMOVED***
			r.campaign(campaignElection)
		***REMOVED*** else ***REMOVED***
			r.becomeLeader()
		***REMOVED***
		return
	***REMOVED***
	for id := range r.prs ***REMOVED***
		if id == r.id ***REMOVED***
			continue
		***REMOVED***
		r.logger.Infof("%x [logterm: %d, index: %d] sent %s request to %x at term %d",
			r.id, r.raftLog.lastTerm(), r.raftLog.lastIndex(), voteMsg, id, r.Term)

		var ctx []byte
		if t == campaignTransfer ***REMOVED***
			ctx = []byte(t)
		***REMOVED***
		r.send(pb.Message***REMOVED***Term: term, To: id, Type: voteMsg, Index: r.raftLog.lastIndex(), LogTerm: r.raftLog.lastTerm(), Context: ctx***REMOVED***)
	***REMOVED***
***REMOVED***

func (r *raft) poll(id uint64, t pb.MessageType, v bool) (granted int) ***REMOVED***
	if v ***REMOVED***
		r.logger.Infof("%x received %s from %x at term %d", r.id, t, id, r.Term)
	***REMOVED*** else ***REMOVED***
		r.logger.Infof("%x received %s rejection from %x at term %d", r.id, t, id, r.Term)
	***REMOVED***
	if _, ok := r.votes[id]; !ok ***REMOVED***
		r.votes[id] = v
	***REMOVED***
	for _, vv := range r.votes ***REMOVED***
		if vv ***REMOVED***
			granted++
		***REMOVED***
	***REMOVED***
	return granted
***REMOVED***

func (r *raft) Step(m pb.Message) error ***REMOVED***
	// Handle the message term, which may result in our stepping down to a follower.
	switch ***REMOVED***
	case m.Term == 0:
		// local message
	case m.Term > r.Term:
		lead := m.From
		if m.Type == pb.MsgVote || m.Type == pb.MsgPreVote ***REMOVED***
			force := bytes.Equal(m.Context, []byte(campaignTransfer))
			inLease := r.checkQuorum && r.lead != None && r.electionElapsed < r.electionTimeout
			if !force && inLease ***REMOVED***
				// If a server receives a RequestVote request within the minimum election timeout
				// of hearing from a current leader, it does not update its term or grant its vote
				r.logger.Infof("%x [logterm: %d, index: %d, vote: %x] ignored %s from %x [logterm: %d, index: %d] at term %d: lease is not expired (remaining ticks: %d)",
					r.id, r.raftLog.lastTerm(), r.raftLog.lastIndex(), r.Vote, m.Type, m.From, m.LogTerm, m.Index, r.Term, r.electionTimeout-r.electionElapsed)
				return nil
			***REMOVED***
			lead = None
		***REMOVED***
		switch ***REMOVED***
		case m.Type == pb.MsgPreVote:
			// Never change our term in response to a PreVote
		case m.Type == pb.MsgPreVoteResp && !m.Reject:
			// We send pre-vote requests with a term in our future. If the
			// pre-vote is granted, we will increment our term when we get a
			// quorum. If it is not, the term comes from the node that
			// rejected our vote so we should become a follower at the new
			// term.
		default:
			r.logger.Infof("%x [term: %d] received a %s message with higher term from %x [term: %d]",
				r.id, r.Term, m.Type, m.From, m.Term)
			r.becomeFollower(m.Term, lead)
		***REMOVED***

	case m.Term < r.Term:
		if r.checkQuorum && (m.Type == pb.MsgHeartbeat || m.Type == pb.MsgApp) ***REMOVED***
			// We have received messages from a leader at a lower term. It is possible
			// that these messages were simply delayed in the network, but this could
			// also mean that this node has advanced its term number during a network
			// partition, and it is now unable to either win an election or to rejoin
			// the majority on the old term. If checkQuorum is false, this will be
			// handled by incrementing term numbers in response to MsgVote with a
			// higher term, but if checkQuorum is true we may not advance the term on
			// MsgVote and must generate other messages to advance the term. The net
			// result of these two features is to minimize the disruption caused by
			// nodes that have been removed from the cluster's configuration: a
			// removed node will send MsgVotes (or MsgPreVotes) which will be ignored,
			// but it will not receive MsgApp or MsgHeartbeat, so it will not create
			// disruptive term increases
			r.send(pb.Message***REMOVED***To: m.From, Type: pb.MsgAppResp***REMOVED***)
		***REMOVED*** else ***REMOVED***
			// ignore other cases
			r.logger.Infof("%x [term: %d] ignored a %s message with lower term from %x [term: %d]",
				r.id, r.Term, m.Type, m.From, m.Term)
		***REMOVED***
		return nil
	***REMOVED***

	switch m.Type ***REMOVED***
	case pb.MsgHup:
		if r.state != StateLeader ***REMOVED***
			ents, err := r.raftLog.slice(r.raftLog.applied+1, r.raftLog.committed+1, noLimit)
			if err != nil ***REMOVED***
				r.logger.Panicf("unexpected error getting unapplied entries (%v)", err)
			***REMOVED***
			if n := numOfPendingConf(ents); n != 0 && r.raftLog.committed > r.raftLog.applied ***REMOVED***
				r.logger.Warningf("%x cannot campaign at term %d since there are still %d pending configuration changes to apply", r.id, r.Term, n)
				return nil
			***REMOVED***

			r.logger.Infof("%x is starting a new election at term %d", r.id, r.Term)
			if r.preVote ***REMOVED***
				r.campaign(campaignPreElection)
			***REMOVED*** else ***REMOVED***
				r.campaign(campaignElection)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			r.logger.Debugf("%x ignoring MsgHup because already leader", r.id)
		***REMOVED***

	case pb.MsgVote, pb.MsgPreVote:
		// The m.Term > r.Term clause is for MsgPreVote. For MsgVote m.Term should
		// always equal r.Term.
		if (r.Vote == None || m.Term > r.Term || r.Vote == m.From) && r.raftLog.isUpToDate(m.Index, m.LogTerm) ***REMOVED***
			r.logger.Infof("%x [logterm: %d, index: %d, vote: %x] cast %s for %x [logterm: %d, index: %d] at term %d",
				r.id, r.raftLog.lastTerm(), r.raftLog.lastIndex(), r.Vote, m.Type, m.From, m.LogTerm, m.Index, r.Term)
			r.send(pb.Message***REMOVED***To: m.From, Type: voteRespMsgType(m.Type)***REMOVED***)
			if m.Type == pb.MsgVote ***REMOVED***
				// Only record real votes.
				r.electionElapsed = 0
				r.Vote = m.From
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			r.logger.Infof("%x [logterm: %d, index: %d, vote: %x] rejected %s from %x [logterm: %d, index: %d] at term %d",
				r.id, r.raftLog.lastTerm(), r.raftLog.lastIndex(), r.Vote, m.Type, m.From, m.LogTerm, m.Index, r.Term)
			r.send(pb.Message***REMOVED***To: m.From, Type: voteRespMsgType(m.Type), Reject: true***REMOVED***)
		***REMOVED***

	default:
		r.step(r, m)
	***REMOVED***
	return nil
***REMOVED***

type stepFunc func(r *raft, m pb.Message)

func stepLeader(r *raft, m pb.Message) ***REMOVED***
	// These message types do not require any progress for m.From.
	switch m.Type ***REMOVED***
	case pb.MsgBeat:
		r.bcastHeartbeat()
		return
	case pb.MsgCheckQuorum:
		if !r.checkQuorumActive() ***REMOVED***
			r.logger.Warningf("%x stepped down to follower since quorum is not active", r.id)
			r.becomeFollower(r.Term, None)
		***REMOVED***
		return
	case pb.MsgProp:
		if len(m.Entries) == 0 ***REMOVED***
			r.logger.Panicf("%x stepped empty MsgProp", r.id)
		***REMOVED***
		if _, ok := r.prs[r.id]; !ok ***REMOVED***
			// If we are not currently a member of the range (i.e. this node
			// was removed from the configuration while serving as leader),
			// drop any new proposals.
			return
		***REMOVED***
		if r.leadTransferee != None ***REMOVED***
			r.logger.Debugf("%x [term %d] transfer leadership to %x is in progress; dropping proposal", r.id, r.Term, r.leadTransferee)
			return
		***REMOVED***

		for i, e := range m.Entries ***REMOVED***
			if e.Type == pb.EntryConfChange ***REMOVED***
				if r.pendingConf ***REMOVED***
					r.logger.Infof("propose conf %s ignored since pending unapplied configuration", e.String())
					m.Entries[i] = pb.Entry***REMOVED***Type: pb.EntryNormal***REMOVED***
				***REMOVED***
				r.pendingConf = true
			***REMOVED***
		***REMOVED***
		r.appendEntry(m.Entries...)
		r.bcastAppend()
		return
	case pb.MsgReadIndex:
		if r.quorum() > 1 ***REMOVED***
			if r.raftLog.zeroTermOnErrCompacted(r.raftLog.term(r.raftLog.committed)) != r.Term ***REMOVED***
				// Reject read only request when this leader has not committed any log entry at its term.
				return
			***REMOVED***

			// thinking: use an interally defined context instead of the user given context.
			// We can express this in terms of the term and index instead of a user-supplied value.
			// This would allow multiple reads to piggyback on the same message.
			switch r.readOnly.option ***REMOVED***
			case ReadOnlySafe:
				r.readOnly.addRequest(r.raftLog.committed, m)
				r.bcastHeartbeatWithCtx(m.Entries[0].Data)
			case ReadOnlyLeaseBased:
				var ri uint64
				if r.checkQuorum ***REMOVED***
					ri = r.raftLog.committed
				***REMOVED***
				if m.From == None || m.From == r.id ***REMOVED*** // from local member
					r.readStates = append(r.readStates, ReadState***REMOVED***Index: r.raftLog.committed, RequestCtx: m.Entries[0].Data***REMOVED***)
				***REMOVED*** else ***REMOVED***
					r.send(pb.Message***REMOVED***To: m.From, Type: pb.MsgReadIndexResp, Index: ri, Entries: m.Entries***REMOVED***)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			r.readStates = append(r.readStates, ReadState***REMOVED***Index: r.raftLog.committed, RequestCtx: m.Entries[0].Data***REMOVED***)
		***REMOVED***

		return
	***REMOVED***

	// All other message types require a progress for m.From (pr).
	pr, prOk := r.prs[m.From]
	if !prOk ***REMOVED***
		r.logger.Debugf("%x no progress available for %x", r.id, m.From)
		return
	***REMOVED***
	switch m.Type ***REMOVED***
	case pb.MsgAppResp:
		pr.RecentActive = true

		if m.Reject ***REMOVED***
			r.logger.Debugf("%x received msgApp rejection(lastindex: %d) from %x for index %d",
				r.id, m.RejectHint, m.From, m.Index)
			if pr.maybeDecrTo(m.Index, m.RejectHint) ***REMOVED***
				r.logger.Debugf("%x decreased progress of %x to [%s]", r.id, m.From, pr)
				if pr.State == ProgressStateReplicate ***REMOVED***
					pr.becomeProbe()
				***REMOVED***
				r.sendAppend(m.From)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			oldPaused := pr.IsPaused()
			if pr.maybeUpdate(m.Index) ***REMOVED***
				switch ***REMOVED***
				case pr.State == ProgressStateProbe:
					pr.becomeReplicate()
				case pr.State == ProgressStateSnapshot && pr.needSnapshotAbort():
					r.logger.Debugf("%x snapshot aborted, resumed sending replication messages to %x [%s]", r.id, m.From, pr)
					pr.becomeProbe()
				case pr.State == ProgressStateReplicate:
					pr.ins.freeTo(m.Index)
				***REMOVED***

				if r.maybeCommit() ***REMOVED***
					r.bcastAppend()
				***REMOVED*** else if oldPaused ***REMOVED***
					// update() reset the wait state on this node. If we had delayed sending
					// an update before, send it now.
					r.sendAppend(m.From)
				***REMOVED***
				// Transfer leadership is in progress.
				if m.From == r.leadTransferee && pr.Match == r.raftLog.lastIndex() ***REMOVED***
					r.logger.Infof("%x sent MsgTimeoutNow to %x after received MsgAppResp", r.id, m.From)
					r.sendTimeoutNow(m.From)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	case pb.MsgHeartbeatResp:
		pr.RecentActive = true
		pr.resume()

		// free one slot for the full inflights window to allow progress.
		if pr.State == ProgressStateReplicate && pr.ins.full() ***REMOVED***
			pr.ins.freeFirstOne()
		***REMOVED***
		if pr.Match < r.raftLog.lastIndex() ***REMOVED***
			r.sendAppend(m.From)
		***REMOVED***

		if r.readOnly.option != ReadOnlySafe || len(m.Context) == 0 ***REMOVED***
			return
		***REMOVED***

		ackCount := r.readOnly.recvAck(m)
		if ackCount < r.quorum() ***REMOVED***
			return
		***REMOVED***

		rss := r.readOnly.advance(m)
		for _, rs := range rss ***REMOVED***
			req := rs.req
			if req.From == None || req.From == r.id ***REMOVED*** // from local member
				r.readStates = append(r.readStates, ReadState***REMOVED***Index: rs.index, RequestCtx: req.Entries[0].Data***REMOVED***)
			***REMOVED*** else ***REMOVED***
				r.send(pb.Message***REMOVED***To: req.From, Type: pb.MsgReadIndexResp, Index: rs.index, Entries: req.Entries***REMOVED***)
			***REMOVED***
		***REMOVED***
	case pb.MsgSnapStatus:
		if pr.State != ProgressStateSnapshot ***REMOVED***
			return
		***REMOVED***
		if !m.Reject ***REMOVED***
			pr.becomeProbe()
			r.logger.Debugf("%x snapshot succeeded, resumed sending replication messages to %x [%s]", r.id, m.From, pr)
		***REMOVED*** else ***REMOVED***
			pr.snapshotFailure()
			pr.becomeProbe()
			r.logger.Debugf("%x snapshot failed, resumed sending replication messages to %x [%s]", r.id, m.From, pr)
		***REMOVED***
		// If snapshot finish, wait for the msgAppResp from the remote node before sending
		// out the next msgApp.
		// If snapshot failure, wait for a heartbeat interval before next try
		pr.pause()
	case pb.MsgUnreachable:
		// During optimistic replication, if the remote becomes unreachable,
		// there is huge probability that a MsgApp is lost.
		if pr.State == ProgressStateReplicate ***REMOVED***
			pr.becomeProbe()
		***REMOVED***
		r.logger.Debugf("%x failed to send message to %x because it is unreachable [%s]", r.id, m.From, pr)
	case pb.MsgTransferLeader:
		leadTransferee := m.From
		lastLeadTransferee := r.leadTransferee
		if lastLeadTransferee != None ***REMOVED***
			if lastLeadTransferee == leadTransferee ***REMOVED***
				r.logger.Infof("%x [term %d] transfer leadership to %x is in progress, ignores request to same node %x",
					r.id, r.Term, leadTransferee, leadTransferee)
				return
			***REMOVED***
			r.abortLeaderTransfer()
			r.logger.Infof("%x [term %d] abort previous transferring leadership to %x", r.id, r.Term, lastLeadTransferee)
		***REMOVED***
		if leadTransferee == r.id ***REMOVED***
			r.logger.Debugf("%x is already leader. Ignored transferring leadership to self", r.id)
			return
		***REMOVED***
		// Transfer leadership to third party.
		r.logger.Infof("%x [term %d] starts to transfer leadership to %x", r.id, r.Term, leadTransferee)
		// Transfer leadership should be finished in one electionTimeout, so reset r.electionElapsed.
		r.electionElapsed = 0
		r.leadTransferee = leadTransferee
		if pr.Match == r.raftLog.lastIndex() ***REMOVED***
			r.sendTimeoutNow(leadTransferee)
			r.logger.Infof("%x sends MsgTimeoutNow to %x immediately as %x already has up-to-date log", r.id, leadTransferee, leadTransferee)
		***REMOVED*** else ***REMOVED***
			r.sendAppend(leadTransferee)
		***REMOVED***
	***REMOVED***
***REMOVED***

// stepCandidate is shared by StateCandidate and StatePreCandidate; the difference is
// whether they respond to MsgVoteResp or MsgPreVoteResp.
func stepCandidate(r *raft, m pb.Message) ***REMOVED***
	// Only handle vote responses corresponding to our candidacy (while in
	// StateCandidate, we may get stale MsgPreVoteResp messages in this term from
	// our pre-candidate state).
	var myVoteRespType pb.MessageType
	if r.state == StatePreCandidate ***REMOVED***
		myVoteRespType = pb.MsgPreVoteResp
	***REMOVED*** else ***REMOVED***
		myVoteRespType = pb.MsgVoteResp
	***REMOVED***
	switch m.Type ***REMOVED***
	case pb.MsgProp:
		r.logger.Infof("%x no leader at term %d; dropping proposal", r.id, r.Term)
		return
	case pb.MsgApp:
		r.becomeFollower(r.Term, m.From)
		r.handleAppendEntries(m)
	case pb.MsgHeartbeat:
		r.becomeFollower(r.Term, m.From)
		r.handleHeartbeat(m)
	case pb.MsgSnap:
		r.becomeFollower(m.Term, m.From)
		r.handleSnapshot(m)
	case myVoteRespType:
		gr := r.poll(m.From, m.Type, !m.Reject)
		r.logger.Infof("%x [quorum:%d] has received %d %s votes and %d vote rejections", r.id, r.quorum(), gr, m.Type, len(r.votes)-gr)
		switch r.quorum() ***REMOVED***
		case gr:
			if r.state == StatePreCandidate ***REMOVED***
				r.campaign(campaignElection)
			***REMOVED*** else ***REMOVED***
				r.becomeLeader()
				r.bcastAppend()
			***REMOVED***
		case len(r.votes) - gr:
			r.becomeFollower(r.Term, None)
		***REMOVED***
	case pb.MsgTimeoutNow:
		r.logger.Debugf("%x [term %d state %v] ignored MsgTimeoutNow from %x", r.id, r.Term, r.state, m.From)
	***REMOVED***
***REMOVED***

func stepFollower(r *raft, m pb.Message) ***REMOVED***
	switch m.Type ***REMOVED***
	case pb.MsgProp:
		if r.lead == None ***REMOVED***
			r.logger.Infof("%x no leader at term %d; dropping proposal", r.id, r.Term)
			return
		***REMOVED***
		m.To = r.lead
		r.send(m)
	case pb.MsgApp:
		r.electionElapsed = 0
		r.lead = m.From
		r.handleAppendEntries(m)
	case pb.MsgHeartbeat:
		r.electionElapsed = 0
		r.lead = m.From
		r.handleHeartbeat(m)
	case pb.MsgSnap:
		r.electionElapsed = 0
		r.lead = m.From
		r.handleSnapshot(m)
	case pb.MsgTransferLeader:
		if r.lead == None ***REMOVED***
			r.logger.Infof("%x no leader at term %d; dropping leader transfer msg", r.id, r.Term)
			return
		***REMOVED***
		m.To = r.lead
		r.send(m)
	case pb.MsgTimeoutNow:
		if r.promotable() ***REMOVED***
			r.logger.Infof("%x [term %d] received MsgTimeoutNow from %x and starts an election to get leadership.", r.id, r.Term, m.From)
			// Leadership transfers never use pre-vote even if r.preVote is true; we
			// know we are not recovering from a partition so there is no need for the
			// extra round trip.
			r.campaign(campaignTransfer)
		***REMOVED*** else ***REMOVED***
			r.logger.Infof("%x received MsgTimeoutNow from %x but is not promotable", r.id, m.From)
		***REMOVED***
	case pb.MsgReadIndex:
		if r.lead == None ***REMOVED***
			r.logger.Infof("%x no leader at term %d; dropping index reading msg", r.id, r.Term)
			return
		***REMOVED***
		m.To = r.lead
		r.send(m)
	case pb.MsgReadIndexResp:
		if len(m.Entries) != 1 ***REMOVED***
			r.logger.Errorf("%x invalid format of MsgReadIndexResp from %x, entries count: %d", r.id, m.From, len(m.Entries))
			return
		***REMOVED***
		r.readStates = append(r.readStates, ReadState***REMOVED***Index: m.Index, RequestCtx: m.Entries[0].Data***REMOVED***)
	***REMOVED***
***REMOVED***

func (r *raft) handleAppendEntries(m pb.Message) ***REMOVED***
	if m.Index < r.raftLog.committed ***REMOVED***
		r.send(pb.Message***REMOVED***To: m.From, Type: pb.MsgAppResp, Index: r.raftLog.committed***REMOVED***)
		return
	***REMOVED***

	if mlastIndex, ok := r.raftLog.maybeAppend(m.Index, m.LogTerm, m.Commit, m.Entries...); ok ***REMOVED***
		r.send(pb.Message***REMOVED***To: m.From, Type: pb.MsgAppResp, Index: mlastIndex***REMOVED***)
	***REMOVED*** else ***REMOVED***
		r.logger.Debugf("%x [logterm: %d, index: %d] rejected msgApp [logterm: %d, index: %d] from %x",
			r.id, r.raftLog.zeroTermOnErrCompacted(r.raftLog.term(m.Index)), m.Index, m.LogTerm, m.Index, m.From)
		r.send(pb.Message***REMOVED***To: m.From, Type: pb.MsgAppResp, Index: m.Index, Reject: true, RejectHint: r.raftLog.lastIndex()***REMOVED***)
	***REMOVED***
***REMOVED***

func (r *raft) handleHeartbeat(m pb.Message) ***REMOVED***
	r.raftLog.commitTo(m.Commit)
	r.send(pb.Message***REMOVED***To: m.From, Type: pb.MsgHeartbeatResp, Context: m.Context***REMOVED***)
***REMOVED***

func (r *raft) handleSnapshot(m pb.Message) ***REMOVED***
	sindex, sterm := m.Snapshot.Metadata.Index, m.Snapshot.Metadata.Term
	if r.restore(m.Snapshot) ***REMOVED***
		r.logger.Infof("%x [commit: %d] restored snapshot [index: %d, term: %d]",
			r.id, r.raftLog.committed, sindex, sterm)
		r.send(pb.Message***REMOVED***To: m.From, Type: pb.MsgAppResp, Index: r.raftLog.lastIndex()***REMOVED***)
	***REMOVED*** else ***REMOVED***
		r.logger.Infof("%x [commit: %d] ignored snapshot [index: %d, term: %d]",
			r.id, r.raftLog.committed, sindex, sterm)
		r.send(pb.Message***REMOVED***To: m.From, Type: pb.MsgAppResp, Index: r.raftLog.committed***REMOVED***)
	***REMOVED***
***REMOVED***

// restore recovers the state machine from a snapshot. It restores the log and the
// configuration of state machine.
func (r *raft) restore(s pb.Snapshot) bool ***REMOVED***
	if s.Metadata.Index <= r.raftLog.committed ***REMOVED***
		return false
	***REMOVED***
	if r.raftLog.matchTerm(s.Metadata.Index, s.Metadata.Term) ***REMOVED***
		r.logger.Infof("%x [commit: %d, lastindex: %d, lastterm: %d] fast-forwarded commit to snapshot [index: %d, term: %d]",
			r.id, r.raftLog.committed, r.raftLog.lastIndex(), r.raftLog.lastTerm(), s.Metadata.Index, s.Metadata.Term)
		r.raftLog.commitTo(s.Metadata.Index)
		return false
	***REMOVED***

	r.logger.Infof("%x [commit: %d, lastindex: %d, lastterm: %d] starts to restore snapshot [index: %d, term: %d]",
		r.id, r.raftLog.committed, r.raftLog.lastIndex(), r.raftLog.lastTerm(), s.Metadata.Index, s.Metadata.Term)

	r.raftLog.restore(s)
	r.prs = make(map[uint64]*Progress)
	for _, n := range s.Metadata.ConfState.Nodes ***REMOVED***
		match, next := uint64(0), r.raftLog.lastIndex()+1
		if n == r.id ***REMOVED***
			match = next - 1
		***REMOVED***
		r.setProgress(n, match, next)
		r.logger.Infof("%x restored progress of %x [%s]", r.id, n, r.prs[n])
	***REMOVED***
	return true
***REMOVED***

// promotable indicates whether state machine can be promoted to leader,
// which is true when its own id is in progress list.
func (r *raft) promotable() bool ***REMOVED***
	_, ok := r.prs[r.id]
	return ok
***REMOVED***

func (r *raft) addNode(id uint64) ***REMOVED***
	r.pendingConf = false
	if _, ok := r.prs[id]; ok ***REMOVED***
		// Ignore any redundant addNode calls (which can happen because the
		// initial bootstrapping entries are applied twice).
		return
	***REMOVED***

	r.setProgress(id, 0, r.raftLog.lastIndex()+1)
	// When a node is first added, we should mark it as recently active.
	// Otherwise, CheckQuorum may cause us to step down if it is invoked
	// before the added node has a chance to communicate with us.
	r.prs[id].RecentActive = true
***REMOVED***

func (r *raft) removeNode(id uint64) ***REMOVED***
	r.delProgress(id)
	r.pendingConf = false

	// do not try to commit or abort transferring if there is no nodes in the cluster.
	if len(r.prs) == 0 ***REMOVED***
		return
	***REMOVED***

	// The quorum size is now smaller, so see if any pending entries can
	// be committed.
	if r.maybeCommit() ***REMOVED***
		r.bcastAppend()
	***REMOVED***
	// If the removed node is the leadTransferee, then abort the leadership transferring.
	if r.state == StateLeader && r.leadTransferee == id ***REMOVED***
		r.abortLeaderTransfer()
	***REMOVED***
***REMOVED***

func (r *raft) resetPendingConf() ***REMOVED*** r.pendingConf = false ***REMOVED***

func (r *raft) setProgress(id, match, next uint64) ***REMOVED***
	r.prs[id] = &Progress***REMOVED***Next: next, Match: match, ins: newInflights(r.maxInflight)***REMOVED***
***REMOVED***

func (r *raft) delProgress(id uint64) ***REMOVED***
	delete(r.prs, id)
***REMOVED***

func (r *raft) loadState(state pb.HardState) ***REMOVED***
	if state.Commit < r.raftLog.committed || state.Commit > r.raftLog.lastIndex() ***REMOVED***
		r.logger.Panicf("%x state.commit %d is out of range [%d, %d]", r.id, state.Commit, r.raftLog.committed, r.raftLog.lastIndex())
	***REMOVED***
	r.raftLog.committed = state.Commit
	r.Term = state.Term
	r.Vote = state.Vote
***REMOVED***

// pastElectionTimeout returns true iff r.electionElapsed is greater
// than or equal to the randomized election timeout in
// [electiontimeout, 2 * electiontimeout - 1].
func (r *raft) pastElectionTimeout() bool ***REMOVED***
	return r.electionElapsed >= r.randomizedElectionTimeout
***REMOVED***

func (r *raft) resetRandomizedElectionTimeout() ***REMOVED***
	r.randomizedElectionTimeout = r.electionTimeout + globalRand.Intn(r.electionTimeout)
***REMOVED***

// checkQuorumActive returns true if the quorum is active from
// the view of the local raft state machine. Otherwise, it returns
// false.
// checkQuorumActive also resets all RecentActive to false.
func (r *raft) checkQuorumActive() bool ***REMOVED***
	var act int

	for id := range r.prs ***REMOVED***
		if id == r.id ***REMOVED*** // self is always active
			act++
			continue
		***REMOVED***

		if r.prs[id].RecentActive ***REMOVED***
			act++
		***REMOVED***

		r.prs[id].RecentActive = false
	***REMOVED***

	return act >= r.quorum()
***REMOVED***

func (r *raft) sendTimeoutNow(to uint64) ***REMOVED***
	r.send(pb.Message***REMOVED***To: to, Type: pb.MsgTimeoutNow***REMOVED***)
***REMOVED***

func (r *raft) abortLeaderTransfer() ***REMOVED***
	r.leadTransferee = None
***REMOVED***

func numOfPendingConf(ents []pb.Entry) int ***REMOVED***
	n := 0
	for i := range ents ***REMOVED***
		if ents[i].Type == pb.EntryConfChange ***REMOVED***
			n++
		***REMOVED***
	***REMOVED***
	return n
***REMOVED***
