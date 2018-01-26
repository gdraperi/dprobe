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

import "fmt"

const (
	ProgressStateProbe ProgressStateType = iota
	ProgressStateReplicate
	ProgressStateSnapshot
)

type ProgressStateType uint64

var prstmap = [...]string***REMOVED***
	"ProgressStateProbe",
	"ProgressStateReplicate",
	"ProgressStateSnapshot",
***REMOVED***

func (st ProgressStateType) String() string ***REMOVED*** return prstmap[uint64(st)] ***REMOVED***

// Progress represents a follower’s progress in the view of the leader. Leader maintains
// progresses of all followers, and sends entries to the follower based on its progress.
type Progress struct ***REMOVED***
	Match, Next uint64
	// State defines how the leader should interact with the follower.
	//
	// When in ProgressStateProbe, leader sends at most one replication message
	// per heartbeat interval. It also probes actual progress of the follower.
	//
	// When in ProgressStateReplicate, leader optimistically increases next
	// to the latest entry sent after sending replication message. This is
	// an optimized state for fast replicating log entries to the follower.
	//
	// When in ProgressStateSnapshot, leader should have sent out snapshot
	// before and stops sending any replication message.
	State ProgressStateType
	// Paused is used in ProgressStateProbe.
	// When Paused is true, raft should pause sending replication message to this peer.
	Paused bool
	// PendingSnapshot is used in ProgressStateSnapshot.
	// If there is a pending snapshot, the pendingSnapshot will be set to the
	// index of the snapshot. If pendingSnapshot is set, the replication process of
	// this Progress will be paused. raft will not resend snapshot until the pending one
	// is reported to be failed.
	PendingSnapshot uint64

	// RecentActive is true if the progress is recently active. Receiving any messages
	// from the corresponding follower indicates the progress is active.
	// RecentActive can be reset to false after an election timeout.
	RecentActive bool

	// inflights is a sliding window for the inflight messages.
	// Each inflight message contains one or more log entries.
	// The max number of entries per message is defined in raft config as MaxSizePerMsg.
	// Thus inflight effectively limits both the number of inflight messages
	// and the bandwidth each Progress can use.
	// When inflights is full, no more message should be sent.
	// When a leader sends out a message, the index of the last
	// entry should be added to inflights. The index MUST be added
	// into inflights in order.
	// When a leader receives a reply, the previous inflights should
	// be freed by calling inflights.freeTo with the index of the last
	// received entry.
	ins *inflights
***REMOVED***

func (pr *Progress) resetState(state ProgressStateType) ***REMOVED***
	pr.Paused = false
	pr.PendingSnapshot = 0
	pr.State = state
	pr.ins.reset()
***REMOVED***

func (pr *Progress) becomeProbe() ***REMOVED***
	// If the original state is ProgressStateSnapshot, progress knows that
	// the pending snapshot has been sent to this peer successfully, then
	// probes from pendingSnapshot + 1.
	if pr.State == ProgressStateSnapshot ***REMOVED***
		pendingSnapshot := pr.PendingSnapshot
		pr.resetState(ProgressStateProbe)
		pr.Next = max(pr.Match+1, pendingSnapshot+1)
	***REMOVED*** else ***REMOVED***
		pr.resetState(ProgressStateProbe)
		pr.Next = pr.Match + 1
	***REMOVED***
***REMOVED***

func (pr *Progress) becomeReplicate() ***REMOVED***
	pr.resetState(ProgressStateReplicate)
	pr.Next = pr.Match + 1
***REMOVED***

func (pr *Progress) becomeSnapshot(snapshoti uint64) ***REMOVED***
	pr.resetState(ProgressStateSnapshot)
	pr.PendingSnapshot = snapshoti
***REMOVED***

// maybeUpdate returns false if the given n index comes from an outdated message.
// Otherwise it updates the progress and returns true.
func (pr *Progress) maybeUpdate(n uint64) bool ***REMOVED***
	var updated bool
	if pr.Match < n ***REMOVED***
		pr.Match = n
		updated = true
		pr.resume()
	***REMOVED***
	if pr.Next < n+1 ***REMOVED***
		pr.Next = n + 1
	***REMOVED***
	return updated
***REMOVED***

func (pr *Progress) optimisticUpdate(n uint64) ***REMOVED*** pr.Next = n + 1 ***REMOVED***

// maybeDecrTo returns false if the given to index comes from an out of order message.
// Otherwise it decreases the progress next index to min(rejected, last) and returns true.
func (pr *Progress) maybeDecrTo(rejected, last uint64) bool ***REMOVED***
	if pr.State == ProgressStateReplicate ***REMOVED***
		// the rejection must be stale if the progress has matched and "rejected"
		// is smaller than "match".
		if rejected <= pr.Match ***REMOVED***
			return false
		***REMOVED***
		// directly decrease next to match + 1
		pr.Next = pr.Match + 1
		return true
	***REMOVED***

	// the rejection must be stale if "rejected" does not match next - 1
	if pr.Next-1 != rejected ***REMOVED***
		return false
	***REMOVED***

	if pr.Next = min(rejected, last+1); pr.Next < 1 ***REMOVED***
		pr.Next = 1
	***REMOVED***
	pr.resume()
	return true
***REMOVED***

func (pr *Progress) pause()  ***REMOVED*** pr.Paused = true ***REMOVED***
func (pr *Progress) resume() ***REMOVED*** pr.Paused = false ***REMOVED***

// IsPaused returns whether sending log entries to this node has been
// paused. A node may be paused because it has rejected recent
// MsgApps, is currently waiting for a snapshot, or has reached the
// MaxInflightMsgs limit.
func (pr *Progress) IsPaused() bool ***REMOVED***
	switch pr.State ***REMOVED***
	case ProgressStateProbe:
		return pr.Paused
	case ProgressStateReplicate:
		return pr.ins.full()
	case ProgressStateSnapshot:
		return true
	default:
		panic("unexpected state")
	***REMOVED***
***REMOVED***

func (pr *Progress) snapshotFailure() ***REMOVED*** pr.PendingSnapshot = 0 ***REMOVED***

// needSnapshotAbort returns true if snapshot progress's Match
// is equal or higher than the pendingSnapshot.
func (pr *Progress) needSnapshotAbort() bool ***REMOVED***
	return pr.State == ProgressStateSnapshot && pr.Match >= pr.PendingSnapshot
***REMOVED***

func (pr *Progress) String() string ***REMOVED***
	return fmt.Sprintf("next = %d, match = %d, state = %s, waiting = %v, pendingSnapshot = %d", pr.Next, pr.Match, pr.State, pr.IsPaused(), pr.PendingSnapshot)
***REMOVED***

type inflights struct ***REMOVED***
	// the starting index in the buffer
	start int
	// number of inflights in the buffer
	count int

	// the size of the buffer
	size int

	// buffer contains the index of the last entry
	// inside one message.
	buffer []uint64
***REMOVED***

func newInflights(size int) *inflights ***REMOVED***
	return &inflights***REMOVED***
		size: size,
	***REMOVED***
***REMOVED***

// add adds an inflight into inflights
func (in *inflights) add(inflight uint64) ***REMOVED***
	if in.full() ***REMOVED***
		panic("cannot add into a full inflights")
	***REMOVED***
	next := in.start + in.count
	size := in.size
	if next >= size ***REMOVED***
		next -= size
	***REMOVED***
	if next >= len(in.buffer) ***REMOVED***
		in.growBuf()
	***REMOVED***
	in.buffer[next] = inflight
	in.count++
***REMOVED***

// grow the inflight buffer by doubling up to inflights.size. We grow on demand
// instead of preallocating to inflights.size to handle systems which have
// thousands of Raft groups per process.
func (in *inflights) growBuf() ***REMOVED***
	newSize := len(in.buffer) * 2
	if newSize == 0 ***REMOVED***
		newSize = 1
	***REMOVED*** else if newSize > in.size ***REMOVED***
		newSize = in.size
	***REMOVED***
	newBuffer := make([]uint64, newSize)
	copy(newBuffer, in.buffer)
	in.buffer = newBuffer
***REMOVED***

// freeTo frees the inflights smaller or equal to the given `to` flight.
func (in *inflights) freeTo(to uint64) ***REMOVED***
	if in.count == 0 || to < in.buffer[in.start] ***REMOVED***
		// out of the left side of the window
		return
	***REMOVED***

	i, idx := 0, in.start
	for i = 0; i < in.count; i++ ***REMOVED***
		if to < in.buffer[idx] ***REMOVED*** // found the first large inflight
			break
		***REMOVED***

		// increase index and maybe rotate
		size := in.size
		if idx++; idx >= size ***REMOVED***
			idx -= size
		***REMOVED***
	***REMOVED***
	// free i inflights and set new start index
	in.count -= i
	in.start = idx
	if in.count == 0 ***REMOVED***
		// inflights is empty, reset the start index so that we don't grow the
		// buffer unnecessarily.
		in.start = 0
	***REMOVED***
***REMOVED***

func (in *inflights) freeFirstOne() ***REMOVED*** in.freeTo(in.buffer[in.start]) ***REMOVED***

// full returns true if the inflights is full.
func (in *inflights) full() bool ***REMOVED***
	return in.count == in.size
***REMOVED***

// resets frees all inflights.
func (in *inflights) reset() ***REMOVED***
	in.count = 0
	in.start = 0
***REMOVED***
