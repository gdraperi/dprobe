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
	"fmt"
	"log"

	pb "github.com/coreos/etcd/raft/raftpb"
)

type raftLog struct ***REMOVED***
	// storage contains all stable entries since the last snapshot.
	storage Storage

	// unstable contains all unstable entries and snapshot.
	// they will be saved into storage.
	unstable unstable

	// committed is the highest log position that is known to be in
	// stable storage on a quorum of nodes.
	committed uint64
	// applied is the highest log position that the application has
	// been instructed to apply to its state machine.
	// Invariant: applied <= committed
	applied uint64

	logger Logger
***REMOVED***

// newLog returns log using the given storage. It recovers the log to the state
// that it just commits and applies the latest snapshot.
func newLog(storage Storage, logger Logger) *raftLog ***REMOVED***
	if storage == nil ***REMOVED***
		log.Panic("storage must not be nil")
	***REMOVED***
	log := &raftLog***REMOVED***
		storage: storage,
		logger:  logger,
	***REMOVED***
	firstIndex, err := storage.FirstIndex()
	if err != nil ***REMOVED***
		panic(err) // TODO(bdarnell)
	***REMOVED***
	lastIndex, err := storage.LastIndex()
	if err != nil ***REMOVED***
		panic(err) // TODO(bdarnell)
	***REMOVED***
	log.unstable.offset = lastIndex + 1
	log.unstable.logger = logger
	// Initialize our committed and applied pointers to the time of the last compaction.
	log.committed = firstIndex - 1
	log.applied = firstIndex - 1

	return log
***REMOVED***

func (l *raftLog) String() string ***REMOVED***
	return fmt.Sprintf("committed=%d, applied=%d, unstable.offset=%d, len(unstable.Entries)=%d", l.committed, l.applied, l.unstable.offset, len(l.unstable.entries))
***REMOVED***

// maybeAppend returns (0, false) if the entries cannot be appended. Otherwise,
// it returns (last index of new entries, true).
func (l *raftLog) maybeAppend(index, logTerm, committed uint64, ents ...pb.Entry) (lastnewi uint64, ok bool) ***REMOVED***
	if l.matchTerm(index, logTerm) ***REMOVED***
		lastnewi = index + uint64(len(ents))
		ci := l.findConflict(ents)
		switch ***REMOVED***
		case ci == 0:
		case ci <= l.committed:
			l.logger.Panicf("entry %d conflict with committed entry [committed(%d)]", ci, l.committed)
		default:
			offset := index + 1
			l.append(ents[ci-offset:]...)
		***REMOVED***
		l.commitTo(min(committed, lastnewi))
		return lastnewi, true
	***REMOVED***
	return 0, false
***REMOVED***

func (l *raftLog) append(ents ...pb.Entry) uint64 ***REMOVED***
	if len(ents) == 0 ***REMOVED***
		return l.lastIndex()
	***REMOVED***
	if after := ents[0].Index - 1; after < l.committed ***REMOVED***
		l.logger.Panicf("after(%d) is out of range [committed(%d)]", after, l.committed)
	***REMOVED***
	l.unstable.truncateAndAppend(ents)
	return l.lastIndex()
***REMOVED***

// findConflict finds the index of the conflict.
// It returns the first pair of conflicting entries between the existing
// entries and the given entries, if there are any.
// If there is no conflicting entries, and the existing entries contains
// all the given entries, zero will be returned.
// If there is no conflicting entries, but the given entries contains new
// entries, the index of the first new entry will be returned.
// An entry is considered to be conflicting if it has the same index but
// a different term.
// The first entry MUST have an index equal to the argument 'from'.
// The index of the given entries MUST be continuously increasing.
func (l *raftLog) findConflict(ents []pb.Entry) uint64 ***REMOVED***
	for _, ne := range ents ***REMOVED***
		if !l.matchTerm(ne.Index, ne.Term) ***REMOVED***
			if ne.Index <= l.lastIndex() ***REMOVED***
				l.logger.Infof("found conflict at index %d [existing term: %d, conflicting term: %d]",
					ne.Index, l.zeroTermOnErrCompacted(l.term(ne.Index)), ne.Term)
			***REMOVED***
			return ne.Index
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func (l *raftLog) unstableEntries() []pb.Entry ***REMOVED***
	if len(l.unstable.entries) == 0 ***REMOVED***
		return nil
	***REMOVED***
	return l.unstable.entries
***REMOVED***

// nextEnts returns all the available entries for execution.
// If applied is smaller than the index of snapshot, it returns all committed
// entries after the index of snapshot.
func (l *raftLog) nextEnts() (ents []pb.Entry) ***REMOVED***
	off := max(l.applied+1, l.firstIndex())
	if l.committed+1 > off ***REMOVED***
		ents, err := l.slice(off, l.committed+1, noLimit)
		if err != nil ***REMOVED***
			l.logger.Panicf("unexpected error when getting unapplied entries (%v)", err)
		***REMOVED***
		return ents
	***REMOVED***
	return nil
***REMOVED***

// hasNextEnts returns if there is any available entries for execution. This
// is a fast check without heavy raftLog.slice() in raftLog.nextEnts().
func (l *raftLog) hasNextEnts() bool ***REMOVED***
	off := max(l.applied+1, l.firstIndex())
	return l.committed+1 > off
***REMOVED***

func (l *raftLog) snapshot() (pb.Snapshot, error) ***REMOVED***
	if l.unstable.snapshot != nil ***REMOVED***
		return *l.unstable.snapshot, nil
	***REMOVED***
	return l.storage.Snapshot()
***REMOVED***

func (l *raftLog) firstIndex() uint64 ***REMOVED***
	if i, ok := l.unstable.maybeFirstIndex(); ok ***REMOVED***
		return i
	***REMOVED***
	index, err := l.storage.FirstIndex()
	if err != nil ***REMOVED***
		panic(err) // TODO(bdarnell)
	***REMOVED***
	return index
***REMOVED***

func (l *raftLog) lastIndex() uint64 ***REMOVED***
	if i, ok := l.unstable.maybeLastIndex(); ok ***REMOVED***
		return i
	***REMOVED***
	i, err := l.storage.LastIndex()
	if err != nil ***REMOVED***
		panic(err) // TODO(bdarnell)
	***REMOVED***
	return i
***REMOVED***

func (l *raftLog) commitTo(tocommit uint64) ***REMOVED***
	// never decrease commit
	if l.committed < tocommit ***REMOVED***
		if l.lastIndex() < tocommit ***REMOVED***
			l.logger.Panicf("tocommit(%d) is out of range [lastIndex(%d)]. Was the raft log corrupted, truncated, or lost?", tocommit, l.lastIndex())
		***REMOVED***
		l.committed = tocommit
	***REMOVED***
***REMOVED***

func (l *raftLog) appliedTo(i uint64) ***REMOVED***
	if i == 0 ***REMOVED***
		return
	***REMOVED***
	if l.committed < i || i < l.applied ***REMOVED***
		l.logger.Panicf("applied(%d) is out of range [prevApplied(%d), committed(%d)]", i, l.applied, l.committed)
	***REMOVED***
	l.applied = i
***REMOVED***

func (l *raftLog) stableTo(i, t uint64) ***REMOVED*** l.unstable.stableTo(i, t) ***REMOVED***

func (l *raftLog) stableSnapTo(i uint64) ***REMOVED*** l.unstable.stableSnapTo(i) ***REMOVED***

func (l *raftLog) lastTerm() uint64 ***REMOVED***
	t, err := l.term(l.lastIndex())
	if err != nil ***REMOVED***
		l.logger.Panicf("unexpected error when getting the last term (%v)", err)
	***REMOVED***
	return t
***REMOVED***

func (l *raftLog) term(i uint64) (uint64, error) ***REMOVED***
	// the valid term range is [index of dummy entry, last index]
	dummyIndex := l.firstIndex() - 1
	if i < dummyIndex || i > l.lastIndex() ***REMOVED***
		// TODO: return an error instead?
		return 0, nil
	***REMOVED***

	if t, ok := l.unstable.maybeTerm(i); ok ***REMOVED***
		return t, nil
	***REMOVED***

	t, err := l.storage.Term(i)
	if err == nil ***REMOVED***
		return t, nil
	***REMOVED***
	if err == ErrCompacted || err == ErrUnavailable ***REMOVED***
		return 0, err
	***REMOVED***
	panic(err) // TODO(bdarnell)
***REMOVED***

func (l *raftLog) entries(i, maxsize uint64) ([]pb.Entry, error) ***REMOVED***
	if i > l.lastIndex() ***REMOVED***
		return nil, nil
	***REMOVED***
	return l.slice(i, l.lastIndex()+1, maxsize)
***REMOVED***

// allEntries returns all entries in the log.
func (l *raftLog) allEntries() []pb.Entry ***REMOVED***
	ents, err := l.entries(l.firstIndex(), noLimit)
	if err == nil ***REMOVED***
		return ents
	***REMOVED***
	if err == ErrCompacted ***REMOVED*** // try again if there was a racing compaction
		return l.allEntries()
	***REMOVED***
	// TODO (xiangli): handle error?
	panic(err)
***REMOVED***

// isUpToDate determines if the given (lastIndex,term) log is more up-to-date
// by comparing the index and term of the last entries in the existing logs.
// If the logs have last entries with different terms, then the log with the
// later term is more up-to-date. If the logs end with the same term, then
// whichever log has the larger lastIndex is more up-to-date. If the logs are
// the same, the given log is up-to-date.
func (l *raftLog) isUpToDate(lasti, term uint64) bool ***REMOVED***
	return term > l.lastTerm() || (term == l.lastTerm() && lasti >= l.lastIndex())
***REMOVED***

func (l *raftLog) matchTerm(i, term uint64) bool ***REMOVED***
	t, err := l.term(i)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return t == term
***REMOVED***

func (l *raftLog) maybeCommit(maxIndex, term uint64) bool ***REMOVED***
	if maxIndex > l.committed && l.zeroTermOnErrCompacted(l.term(maxIndex)) == term ***REMOVED***
		l.commitTo(maxIndex)
		return true
	***REMOVED***
	return false
***REMOVED***

func (l *raftLog) restore(s pb.Snapshot) ***REMOVED***
	l.logger.Infof("log [%s] starts to restore snapshot [index: %d, term: %d]", l, s.Metadata.Index, s.Metadata.Term)
	l.committed = s.Metadata.Index
	l.unstable.restore(s)
***REMOVED***

// slice returns a slice of log entries from lo through hi-1, inclusive.
func (l *raftLog) slice(lo, hi, maxSize uint64) ([]pb.Entry, error) ***REMOVED***
	err := l.mustCheckOutOfBounds(lo, hi)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if lo == hi ***REMOVED***
		return nil, nil
	***REMOVED***
	var ents []pb.Entry
	if lo < l.unstable.offset ***REMOVED***
		storedEnts, err := l.storage.Entries(lo, min(hi, l.unstable.offset), maxSize)
		if err == ErrCompacted ***REMOVED***
			return nil, err
		***REMOVED*** else if err == ErrUnavailable ***REMOVED***
			l.logger.Panicf("entries[%d:%d) is unavailable from storage", lo, min(hi, l.unstable.offset))
		***REMOVED*** else if err != nil ***REMOVED***
			panic(err) // TODO(bdarnell)
		***REMOVED***

		// check if ents has reached the size limitation
		if uint64(len(storedEnts)) < min(hi, l.unstable.offset)-lo ***REMOVED***
			return storedEnts, nil
		***REMOVED***

		ents = storedEnts
	***REMOVED***
	if hi > l.unstable.offset ***REMOVED***
		unstable := l.unstable.slice(max(lo, l.unstable.offset), hi)
		if len(ents) > 0 ***REMOVED***
			ents = append([]pb.Entry***REMOVED******REMOVED***, ents...)
			ents = append(ents, unstable...)
		***REMOVED*** else ***REMOVED***
			ents = unstable
		***REMOVED***
	***REMOVED***
	return limitSize(ents, maxSize), nil
***REMOVED***

// l.firstIndex <= lo <= hi <= l.firstIndex + len(l.entries)
func (l *raftLog) mustCheckOutOfBounds(lo, hi uint64) error ***REMOVED***
	if lo > hi ***REMOVED***
		l.logger.Panicf("invalid slice %d > %d", lo, hi)
	***REMOVED***
	fi := l.firstIndex()
	if lo < fi ***REMOVED***
		return ErrCompacted
	***REMOVED***

	length := l.lastIndex() + 1 - fi
	if lo < fi || hi > fi+length ***REMOVED***
		l.logger.Panicf("slice[%d,%d) out of bound [%d,%d]", lo, hi, fi, l.lastIndex())
	***REMOVED***
	return nil
***REMOVED***

func (l *raftLog) zeroTermOnErrCompacted(t uint64, err error) uint64 ***REMOVED***
	if err == nil ***REMOVED***
		return t
	***REMOVED***
	if err == ErrCompacted ***REMOVED***
		return 0
	***REMOVED***
	l.logger.Panicf("unexpected error (%v)", err)
	return 0
***REMOVED***
