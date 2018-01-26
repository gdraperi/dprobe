// Copyright 2016 The etcd Authors
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

import pb "github.com/coreos/etcd/raft/raftpb"

// ReadState provides state for read only query.
// It's caller's responsibility to call ReadIndex first before getting
// this state from ready, It's also caller's duty to differentiate if this
// state is what it requests through RequestCtx, eg. given a unique id as
// RequestCtx
type ReadState struct ***REMOVED***
	Index      uint64
	RequestCtx []byte
***REMOVED***

type readIndexStatus struct ***REMOVED***
	req   pb.Message
	index uint64
	acks  map[uint64]struct***REMOVED******REMOVED***
***REMOVED***

type readOnly struct ***REMOVED***
	option           ReadOnlyOption
	pendingReadIndex map[string]*readIndexStatus
	readIndexQueue   []string
***REMOVED***

func newReadOnly(option ReadOnlyOption) *readOnly ***REMOVED***
	return &readOnly***REMOVED***
		option:           option,
		pendingReadIndex: make(map[string]*readIndexStatus),
	***REMOVED***
***REMOVED***

// addRequest adds a read only reuqest into readonly struct.
// `index` is the commit index of the raft state machine when it received
// the read only request.
// `m` is the original read only request message from the local or remote node.
func (ro *readOnly) addRequest(index uint64, m pb.Message) ***REMOVED***
	ctx := string(m.Entries[0].Data)
	if _, ok := ro.pendingReadIndex[ctx]; ok ***REMOVED***
		return
	***REMOVED***
	ro.pendingReadIndex[ctx] = &readIndexStatus***REMOVED***index: index, req: m, acks: make(map[uint64]struct***REMOVED******REMOVED***)***REMOVED***
	ro.readIndexQueue = append(ro.readIndexQueue, ctx)
***REMOVED***

// recvAck notifies the readonly struct that the raft state machine received
// an acknowledgment of the heartbeat that attached with the read only request
// context.
func (ro *readOnly) recvAck(m pb.Message) int ***REMOVED***
	rs, ok := ro.pendingReadIndex[string(m.Context)]
	if !ok ***REMOVED***
		return 0
	***REMOVED***

	rs.acks[m.From] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	// add one to include an ack from local node
	return len(rs.acks) + 1
***REMOVED***

// advance advances the read only request queue kept by the readonly struct.
// It dequeues the requests until it finds the read only request that has
// the same context as the given `m`.
func (ro *readOnly) advance(m pb.Message) []*readIndexStatus ***REMOVED***
	var (
		i     int
		found bool
	)

	ctx := string(m.Context)
	rss := []*readIndexStatus***REMOVED******REMOVED***

	for _, okctx := range ro.readIndexQueue ***REMOVED***
		i++
		rs, ok := ro.pendingReadIndex[okctx]
		if !ok ***REMOVED***
			panic("cannot find corresponding read state from pending map")
		***REMOVED***
		rss = append(rss, rs)
		if okctx == ctx ***REMOVED***
			found = true
			break
		***REMOVED***
	***REMOVED***

	if found ***REMOVED***
		ro.readIndexQueue = ro.readIndexQueue[i:]
		for _, rs := range rss ***REMOVED***
			delete(ro.pendingReadIndex, string(rs.req.Entries[0].Data))
		***REMOVED***
		return rss
	***REMOVED***

	return nil
***REMOVED***

// lastPendingRequestCtx returns the context of the last pending read only
// request in readonly struct.
func (ro *readOnly) lastPendingRequestCtx() string ***REMOVED***
	if len(ro.readIndexQueue) == 0 ***REMOVED***
		return ""
	***REMOVED***
	return ro.readIndexQueue[len(ro.readIndexQueue)-1]
***REMOVED***
