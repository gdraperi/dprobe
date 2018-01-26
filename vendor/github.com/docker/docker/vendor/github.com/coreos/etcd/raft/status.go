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

	pb "github.com/coreos/etcd/raft/raftpb"
)

type Status struct ***REMOVED***
	ID uint64

	pb.HardState
	SoftState

	Applied  uint64
	Progress map[uint64]Progress
***REMOVED***

// getStatus gets a copy of the current raft status.
func getStatus(r *raft) Status ***REMOVED***
	s := Status***REMOVED***ID: r.id***REMOVED***
	s.HardState = r.hardState()
	s.SoftState = *r.softState()

	s.Applied = r.raftLog.applied

	if s.RaftState == StateLeader ***REMOVED***
		s.Progress = make(map[uint64]Progress)
		for id, p := range r.prs ***REMOVED***
			s.Progress[id] = *p
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***

// MarshalJSON translates the raft status into JSON.
// TODO: try to simplify this by introducing ID type into raft
func (s Status) MarshalJSON() ([]byte, error) ***REMOVED***
	j := fmt.Sprintf(`***REMOVED***"id":"%x","term":%d,"vote":"%x","commit":%d,"lead":"%x","raftState":%q,"progress":***REMOVED***`,
		s.ID, s.Term, s.Vote, s.Commit, s.Lead, s.RaftState)

	if len(s.Progress) == 0 ***REMOVED***
		j += "***REMOVED******REMOVED***"
	***REMOVED*** else ***REMOVED***
		for k, v := range s.Progress ***REMOVED***
			subj := fmt.Sprintf(`"%x":***REMOVED***"match":%d,"next":%d,"state":%q***REMOVED***,`, k, v.Match, v.Next, v.State)
			j += subj
		***REMOVED***
		// remove the trailing ","
		j = j[:len(j)-1] + "***REMOVED******REMOVED***"
	***REMOVED***
	return []byte(j), nil
***REMOVED***

func (s Status) String() string ***REMOVED***
	b, err := s.MarshalJSON()
	if err != nil ***REMOVED***
		raftLogger.Panicf("unexpected error: %v", err)
	***REMOVED***
	return string(b)
***REMOVED***
