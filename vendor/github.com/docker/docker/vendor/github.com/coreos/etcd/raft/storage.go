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
	"sync"

	pb "github.com/coreos/etcd/raft/raftpb"
)

// ErrCompacted is returned by Storage.Entries/Compact when a requested
// index is unavailable because it predates the last snapshot.
var ErrCompacted = errors.New("requested index is unavailable due to compaction")

// ErrSnapOutOfDate is returned by Storage.CreateSnapshot when a requested
// index is older than the existing snapshot.
var ErrSnapOutOfDate = errors.New("requested index is older than the existing snapshot")

// ErrUnavailable is returned by Storage interface when the requested log entries
// are unavailable.
var ErrUnavailable = errors.New("requested entry at index is unavailable")

// ErrSnapshotTemporarilyUnavailable is returned by the Storage interface when the required
// snapshot is temporarily unavailable.
var ErrSnapshotTemporarilyUnavailable = errors.New("snapshot is temporarily unavailable")

// Storage is an interface that may be implemented by the application
// to retrieve log entries from storage.
//
// If any Storage method returns an error, the raft instance will
// become inoperable and refuse to participate in elections; the
// application is responsible for cleanup and recovery in this case.
type Storage interface ***REMOVED***
	// InitialState returns the saved HardState and ConfState information.
	InitialState() (pb.HardState, pb.ConfState, error)
	// Entries returns a slice of log entries in the range [lo,hi).
	// MaxSize limits the total size of the log entries returned, but
	// Entries returns at least one entry if any.
	Entries(lo, hi, maxSize uint64) ([]pb.Entry, error)
	// Term returns the term of entry i, which must be in the range
	// [FirstIndex()-1, LastIndex()]. The term of the entry before
	// FirstIndex is retained for matching purposes even though the
	// rest of that entry may not be available.
	Term(i uint64) (uint64, error)
	// LastIndex returns the index of the last entry in the log.
	LastIndex() (uint64, error)
	// FirstIndex returns the index of the first log entry that is
	// possibly available via Entries (older entries have been incorporated
	// into the latest Snapshot; if storage only contains the dummy entry the
	// first log entry is not available).
	FirstIndex() (uint64, error)
	// Snapshot returns the most recent snapshot.
	// If snapshot is temporarily unavailable, it should return ErrSnapshotTemporarilyUnavailable,
	// so raft state machine could know that Storage needs some time to prepare
	// snapshot and call Snapshot later.
	Snapshot() (pb.Snapshot, error)
***REMOVED***

// MemoryStorage implements the Storage interface backed by an
// in-memory array.
type MemoryStorage struct ***REMOVED***
	// Protects access to all fields. Most methods of MemoryStorage are
	// run on the raft goroutine, but Append() is run on an application
	// goroutine.
	sync.Mutex

	hardState pb.HardState
	snapshot  pb.Snapshot
	// ents[i] has raft log position i+snapshot.Metadata.Index
	ents []pb.Entry
***REMOVED***

// NewMemoryStorage creates an empty MemoryStorage.
func NewMemoryStorage() *MemoryStorage ***REMOVED***
	return &MemoryStorage***REMOVED***
		// When starting from scratch populate the list with a dummy entry at term zero.
		ents: make([]pb.Entry, 1),
	***REMOVED***
***REMOVED***

// InitialState implements the Storage interface.
func (ms *MemoryStorage) InitialState() (pb.HardState, pb.ConfState, error) ***REMOVED***
	return ms.hardState, ms.snapshot.Metadata.ConfState, nil
***REMOVED***

// SetHardState saves the current HardState.
func (ms *MemoryStorage) SetHardState(st pb.HardState) error ***REMOVED***
	ms.Lock()
	defer ms.Unlock()
	ms.hardState = st
	return nil
***REMOVED***

// Entries implements the Storage interface.
func (ms *MemoryStorage) Entries(lo, hi, maxSize uint64) ([]pb.Entry, error) ***REMOVED***
	ms.Lock()
	defer ms.Unlock()
	offset := ms.ents[0].Index
	if lo <= offset ***REMOVED***
		return nil, ErrCompacted
	***REMOVED***
	if hi > ms.lastIndex()+1 ***REMOVED***
		raftLogger.Panicf("entries' hi(%d) is out of bound lastindex(%d)", hi, ms.lastIndex())
	***REMOVED***
	// only contains dummy entries.
	if len(ms.ents) == 1 ***REMOVED***
		return nil, ErrUnavailable
	***REMOVED***

	ents := ms.ents[lo-offset : hi-offset]
	return limitSize(ents, maxSize), nil
***REMOVED***

// Term implements the Storage interface.
func (ms *MemoryStorage) Term(i uint64) (uint64, error) ***REMOVED***
	ms.Lock()
	defer ms.Unlock()
	offset := ms.ents[0].Index
	if i < offset ***REMOVED***
		return 0, ErrCompacted
	***REMOVED***
	if int(i-offset) >= len(ms.ents) ***REMOVED***
		return 0, ErrUnavailable
	***REMOVED***
	return ms.ents[i-offset].Term, nil
***REMOVED***

// LastIndex implements the Storage interface.
func (ms *MemoryStorage) LastIndex() (uint64, error) ***REMOVED***
	ms.Lock()
	defer ms.Unlock()
	return ms.lastIndex(), nil
***REMOVED***

func (ms *MemoryStorage) lastIndex() uint64 ***REMOVED***
	return ms.ents[0].Index + uint64(len(ms.ents)) - 1
***REMOVED***

// FirstIndex implements the Storage interface.
func (ms *MemoryStorage) FirstIndex() (uint64, error) ***REMOVED***
	ms.Lock()
	defer ms.Unlock()
	return ms.firstIndex(), nil
***REMOVED***

func (ms *MemoryStorage) firstIndex() uint64 ***REMOVED***
	return ms.ents[0].Index + 1
***REMOVED***

// Snapshot implements the Storage interface.
func (ms *MemoryStorage) Snapshot() (pb.Snapshot, error) ***REMOVED***
	ms.Lock()
	defer ms.Unlock()
	return ms.snapshot, nil
***REMOVED***

// ApplySnapshot overwrites the contents of this Storage object with
// those of the given snapshot.
func (ms *MemoryStorage) ApplySnapshot(snap pb.Snapshot) error ***REMOVED***
	ms.Lock()
	defer ms.Unlock()

	//handle check for old snapshot being applied
	msIndex := ms.snapshot.Metadata.Index
	snapIndex := snap.Metadata.Index
	if msIndex >= snapIndex ***REMOVED***
		return ErrSnapOutOfDate
	***REMOVED***

	ms.snapshot = snap
	ms.ents = []pb.Entry***REMOVED******REMOVED***Term: snap.Metadata.Term, Index: snap.Metadata.Index***REMOVED******REMOVED***
	return nil
***REMOVED***

// CreateSnapshot makes a snapshot which can be retrieved with Snapshot() and
// can be used to reconstruct the state at that point.
// If any configuration changes have been made since the last compaction,
// the result of the last ApplyConfChange must be passed in.
func (ms *MemoryStorage) CreateSnapshot(i uint64, cs *pb.ConfState, data []byte) (pb.Snapshot, error) ***REMOVED***
	ms.Lock()
	defer ms.Unlock()
	if i <= ms.snapshot.Metadata.Index ***REMOVED***
		return pb.Snapshot***REMOVED******REMOVED***, ErrSnapOutOfDate
	***REMOVED***

	offset := ms.ents[0].Index
	if i > ms.lastIndex() ***REMOVED***
		raftLogger.Panicf("snapshot %d is out of bound lastindex(%d)", i, ms.lastIndex())
	***REMOVED***

	ms.snapshot.Metadata.Index = i
	ms.snapshot.Metadata.Term = ms.ents[i-offset].Term
	if cs != nil ***REMOVED***
		ms.snapshot.Metadata.ConfState = *cs
	***REMOVED***
	ms.snapshot.Data = data
	return ms.snapshot, nil
***REMOVED***

// Compact discards all log entries prior to compactIndex.
// It is the application's responsibility to not attempt to compact an index
// greater than raftLog.applied.
func (ms *MemoryStorage) Compact(compactIndex uint64) error ***REMOVED***
	ms.Lock()
	defer ms.Unlock()
	offset := ms.ents[0].Index
	if compactIndex <= offset ***REMOVED***
		return ErrCompacted
	***REMOVED***
	if compactIndex > ms.lastIndex() ***REMOVED***
		raftLogger.Panicf("compact %d is out of bound lastindex(%d)", compactIndex, ms.lastIndex())
	***REMOVED***

	i := compactIndex - offset
	ents := make([]pb.Entry, 1, 1+uint64(len(ms.ents))-i)
	ents[0].Index = ms.ents[i].Index
	ents[0].Term = ms.ents[i].Term
	ents = append(ents, ms.ents[i+1:]...)
	ms.ents = ents
	return nil
***REMOVED***

// Append the new entries to storage.
// TODO (xiangli): ensure the entries are continuous and
// entries[0].Index > ms.entries[0].Index
func (ms *MemoryStorage) Append(entries []pb.Entry) error ***REMOVED***
	if len(entries) == 0 ***REMOVED***
		return nil
	***REMOVED***

	ms.Lock()
	defer ms.Unlock()

	first := ms.firstIndex()
	last := entries[0].Index + uint64(len(entries)) - 1

	// shortcut if there is no new entry.
	if last < first ***REMOVED***
		return nil
	***REMOVED***
	// truncate compacted entries
	if first > entries[0].Index ***REMOVED***
		entries = entries[first-entries[0].Index:]
	***REMOVED***

	offset := entries[0].Index - ms.ents[0].Index
	switch ***REMOVED***
	case uint64(len(ms.ents)) > offset:
		ms.ents = append([]pb.Entry***REMOVED******REMOVED***, ms.ents[:offset]...)
		ms.ents = append(ms.ents, entries...)
	case uint64(len(ms.ents)) == offset:
		ms.ents = append(ms.ents, entries...)
	default:
		raftLogger.Panicf("missing log entry [last: %d, append at: %d]",
			ms.lastIndex(), entries[0].Index)
	***REMOVED***
	return nil
***REMOVED***
