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
)

// ErrStepLocalMsg is returned when try to step a local raft message
var ErrStepLocalMsg = errors.New("raft: cannot step raft local message")

// ErrStepPeerNotFound is returned when try to step a response message
// but there is no peer found in raft.prs for that node.
var ErrStepPeerNotFound = errors.New("raft: cannot step as peer not found")

// RawNode is a thread-unsafe Node.
// The methods of this struct correspond to the methods of Node and are described
// more fully there.
type RawNode struct ***REMOVED***
	raft       *raft
	prevSoftSt *SoftState
	prevHardSt pb.HardState
***REMOVED***

func (rn *RawNode) newReady() Ready ***REMOVED***
	return newReady(rn.raft, rn.prevSoftSt, rn.prevHardSt)
***REMOVED***

func (rn *RawNode) commitReady(rd Ready) ***REMOVED***
	if rd.SoftState != nil ***REMOVED***
		rn.prevSoftSt = rd.SoftState
	***REMOVED***
	if !IsEmptyHardState(rd.HardState) ***REMOVED***
		rn.prevHardSt = rd.HardState
	***REMOVED***
	if rn.prevHardSt.Commit != 0 ***REMOVED***
		// In most cases, prevHardSt and rd.HardState will be the same
		// because when there are new entries to apply we just sent a
		// HardState with an updated Commit value. However, on initial
		// startup the two are different because we don't send a HardState
		// until something changes, but we do send any un-applied but
		// committed entries (and previously-committed entries may be
		// incorporated into the snapshot, even if rd.CommittedEntries is
		// empty). Therefore we mark all committed entries as applied
		// whether they were included in rd.HardState or not.
		rn.raft.raftLog.appliedTo(rn.prevHardSt.Commit)
	***REMOVED***
	if len(rd.Entries) > 0 ***REMOVED***
		e := rd.Entries[len(rd.Entries)-1]
		rn.raft.raftLog.stableTo(e.Index, e.Term)
	***REMOVED***
	if !IsEmptySnap(rd.Snapshot) ***REMOVED***
		rn.raft.raftLog.stableSnapTo(rd.Snapshot.Metadata.Index)
	***REMOVED***
	if len(rd.ReadStates) != 0 ***REMOVED***
		rn.raft.readStates = nil
	***REMOVED***
***REMOVED***

// NewRawNode returns a new RawNode given configuration and a list of raft peers.
func NewRawNode(config *Config, peers []Peer) (*RawNode, error) ***REMOVED***
	if config.ID == 0 ***REMOVED***
		panic("config.ID must not be zero")
	***REMOVED***
	r := newRaft(config)
	rn := &RawNode***REMOVED***
		raft: r,
	***REMOVED***
	lastIndex, err := config.Storage.LastIndex()
	if err != nil ***REMOVED***
		panic(err) // TODO(bdarnell)
	***REMOVED***
	// If the log is empty, this is a new RawNode (like StartNode); otherwise it's
	// restoring an existing RawNode (like RestartNode).
	// TODO(bdarnell): rethink RawNode initialization and whether the application needs
	// to be able to tell us when it expects the RawNode to exist.
	if lastIndex == 0 ***REMOVED***
		r.becomeFollower(1, None)
		ents := make([]pb.Entry, len(peers))
		for i, peer := range peers ***REMOVED***
			cc := pb.ConfChange***REMOVED***Type: pb.ConfChangeAddNode, NodeID: peer.ID, Context: peer.Context***REMOVED***
			data, err := cc.Marshal()
			if err != nil ***REMOVED***
				panic("unexpected marshal error")
			***REMOVED***

			ents[i] = pb.Entry***REMOVED***Type: pb.EntryConfChange, Term: 1, Index: uint64(i + 1), Data: data***REMOVED***
		***REMOVED***
		r.raftLog.append(ents...)
		r.raftLog.committed = uint64(len(ents))
		for _, peer := range peers ***REMOVED***
			r.addNode(peer.ID)
		***REMOVED***
	***REMOVED***

	// Set the initial hard and soft states after performing all initialization.
	rn.prevSoftSt = r.softState()
	if lastIndex == 0 ***REMOVED***
		rn.prevHardSt = emptyState
	***REMOVED*** else ***REMOVED***
		rn.prevHardSt = r.hardState()
	***REMOVED***

	return rn, nil
***REMOVED***

// Tick advances the internal logical clock by a single tick.
func (rn *RawNode) Tick() ***REMOVED***
	rn.raft.tick()
***REMOVED***

// TickQuiesced advances the internal logical clock by a single tick without
// performing any other state machine processing. It allows the caller to avoid
// periodic heartbeats and elections when all of the peers in a Raft group are
// known to be at the same state. Expected usage is to periodically invoke Tick
// or TickQuiesced depending on whether the group is "active" or "quiesced".
//
// WARNING: Be very careful about using this method as it subverts the Raft
// state machine. You should probably be using Tick instead.
func (rn *RawNode) TickQuiesced() ***REMOVED***
	rn.raft.electionElapsed++
***REMOVED***

// Campaign causes this RawNode to transition to candidate state.
func (rn *RawNode) Campaign() error ***REMOVED***
	return rn.raft.Step(pb.Message***REMOVED***
		Type: pb.MsgHup,
	***REMOVED***)
***REMOVED***

// Propose proposes data be appended to the raft log.
func (rn *RawNode) Propose(data []byte) error ***REMOVED***
	return rn.raft.Step(pb.Message***REMOVED***
		Type: pb.MsgProp,
		From: rn.raft.id,
		Entries: []pb.Entry***REMOVED***
			***REMOVED***Data: data***REMOVED***,
		***REMOVED******REMOVED***)
***REMOVED***

// ProposeConfChange proposes a config change.
func (rn *RawNode) ProposeConfChange(cc pb.ConfChange) error ***REMOVED***
	data, err := cc.Marshal()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return rn.raft.Step(pb.Message***REMOVED***
		Type: pb.MsgProp,
		Entries: []pb.Entry***REMOVED***
			***REMOVED***Type: pb.EntryConfChange, Data: data***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// ApplyConfChange applies a config change to the local node.
func (rn *RawNode) ApplyConfChange(cc pb.ConfChange) *pb.ConfState ***REMOVED***
	if cc.NodeID == None ***REMOVED***
		rn.raft.resetPendingConf()
		return &pb.ConfState***REMOVED***Nodes: rn.raft.nodes()***REMOVED***
	***REMOVED***
	switch cc.Type ***REMOVED***
	case pb.ConfChangeAddNode:
		rn.raft.addNode(cc.NodeID)
	case pb.ConfChangeRemoveNode:
		rn.raft.removeNode(cc.NodeID)
	case pb.ConfChangeUpdateNode:
		rn.raft.resetPendingConf()
	default:
		panic("unexpected conf type")
	***REMOVED***
	return &pb.ConfState***REMOVED***Nodes: rn.raft.nodes()***REMOVED***
***REMOVED***

// Step advances the state machine using the given message.
func (rn *RawNode) Step(m pb.Message) error ***REMOVED***
	// ignore unexpected local messages receiving over network
	if IsLocalMsg(m.Type) ***REMOVED***
		return ErrStepLocalMsg
	***REMOVED***
	if _, ok := rn.raft.prs[m.From]; ok || !IsResponseMsg(m.Type) ***REMOVED***
		return rn.raft.Step(m)
	***REMOVED***
	return ErrStepPeerNotFound
***REMOVED***

// Ready returns the current point-in-time state of this RawNode.
func (rn *RawNode) Ready() Ready ***REMOVED***
	rd := rn.newReady()
	rn.raft.msgs = nil
	return rd
***REMOVED***

// HasReady called when RawNode user need to check if any Ready pending.
// Checking logic in this method should be consistent with Ready.containsUpdates().
func (rn *RawNode) HasReady() bool ***REMOVED***
	r := rn.raft
	if !r.softState().equal(rn.prevSoftSt) ***REMOVED***
		return true
	***REMOVED***
	if hardSt := r.hardState(); !IsEmptyHardState(hardSt) && !isHardStateEqual(hardSt, rn.prevHardSt) ***REMOVED***
		return true
	***REMOVED***
	if r.raftLog.unstable.snapshot != nil && !IsEmptySnap(*r.raftLog.unstable.snapshot) ***REMOVED***
		return true
	***REMOVED***
	if len(r.msgs) > 0 || len(r.raftLog.unstableEntries()) > 0 || r.raftLog.hasNextEnts() ***REMOVED***
		return true
	***REMOVED***
	if len(r.readStates) != 0 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// Advance notifies the RawNode that the application has applied and saved progress in the
// last Ready results.
func (rn *RawNode) Advance(rd Ready) ***REMOVED***
	rn.commitReady(rd)
***REMOVED***

// Status returns the current status of the given group.
func (rn *RawNode) Status() *Status ***REMOVED***
	status := getStatus(rn.raft)
	return &status
***REMOVED***

// ReportUnreachable reports the given node is not reachable for the last send.
func (rn *RawNode) ReportUnreachable(id uint64) ***REMOVED***
	_ = rn.raft.Step(pb.Message***REMOVED***Type: pb.MsgUnreachable, From: id***REMOVED***)
***REMOVED***

// ReportSnapshot reports the status of the sent snapshot.
func (rn *RawNode) ReportSnapshot(id uint64, status SnapshotStatus) ***REMOVED***
	rej := status == SnapshotFailure

	_ = rn.raft.Step(pb.Message***REMOVED***Type: pb.MsgSnapStatus, From: id, Reject: rej***REMOVED***)
***REMOVED***

// TransferLeader tries to transfer leadership to the given transferee.
func (rn *RawNode) TransferLeader(transferee uint64) ***REMOVED***
	_ = rn.raft.Step(pb.Message***REMOVED***Type: pb.MsgTransferLeader, From: transferee***REMOVED***)
***REMOVED***

// ReadIndex requests a read state. The read state will be set in ready.
// Read State has a read index. Once the application advances further than the read
// index, any linearizable read requests issued before the read request can be
// processed safely. The read state will have the same rctx attached.
func (rn *RawNode) ReadIndex(rctx []byte) ***REMOVED***
	_ = rn.raft.Step(pb.Message***REMOVED***Type: pb.MsgReadIndex, Entries: []pb.Entry***REMOVED******REMOVED***Data: rctx***REMOVED******REMOVED******REMOVED***)
***REMOVED***
