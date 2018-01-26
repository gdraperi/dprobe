package membership

import (
	"errors"
	"sync"

	"github.com/coreos/etcd/raft/raftpb"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/watch"
	"github.com/gogo/protobuf/proto"
)

var (
	// ErrIDExists is thrown when a node wants to join the existing cluster but its ID already exists
	ErrIDExists = errors.New("membership: can't add node to cluster, node id is a duplicate")
	// ErrIDRemoved is thrown when a node tries to perform an operation on an existing cluster but was removed
	ErrIDRemoved = errors.New("membership: node was removed during cluster lifetime")
	// ErrIDNotFound is thrown when we try an operation on a member that does not exist in the cluster list
	ErrIDNotFound = errors.New("membership: member not found in cluster list")
	// ErrConfigChangeInvalid is thrown when a configuration change we received looks invalid in form
	ErrConfigChangeInvalid = errors.New("membership: ConfChange type should be either AddNode, RemoveNode or UpdateNode")
	// ErrCannotUnmarshalConfig is thrown when a node cannot unmarshal a configuration change
	ErrCannotUnmarshalConfig = errors.New("membership: cannot unmarshal configuration change")
	// ErrMemberRemoved is thrown when a node was removed from the cluster
	ErrMemberRemoved = errors.New("raft: member was removed from the cluster")
)

// Cluster represents a set of active
// raft Members
type Cluster struct ***REMOVED***
	mu      sync.RWMutex
	members map[uint64]*Member

	// removed contains the list of removed Members,
	// those ids cannot be reused
	removed map[uint64]bool

	PeersBroadcast *watch.Queue
***REMOVED***

// Member represents a raft Cluster Member
type Member struct ***REMOVED***
	*api.RaftMember
***REMOVED***

// NewCluster creates a new Cluster neighbors list for a raft Member.
func NewCluster() *Cluster ***REMOVED***
	// TODO(abronan): generate Cluster ID for federation

	return &Cluster***REMOVED***
		members:        make(map[uint64]*Member),
		removed:        make(map[uint64]bool),
		PeersBroadcast: watch.NewQueue(),
	***REMOVED***
***REMOVED***

// Members returns the list of raft Members in the Cluster.
func (c *Cluster) Members() map[uint64]*Member ***REMOVED***
	members := make(map[uint64]*Member)
	c.mu.RLock()
	for k, v := range c.members ***REMOVED***
		members[k] = v
	***REMOVED***
	c.mu.RUnlock()
	return members
***REMOVED***

// Removed returns the list of raft Members removed from the Cluster.
func (c *Cluster) Removed() []uint64 ***REMOVED***
	c.mu.RLock()
	removed := make([]uint64, 0, len(c.removed))
	for k := range c.removed ***REMOVED***
		removed = append(removed, k)
	***REMOVED***
	c.mu.RUnlock()
	return removed
***REMOVED***

// GetMember returns informations on a given Member.
func (c *Cluster) GetMember(id uint64) *Member ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.members[id]
***REMOVED***

func (c *Cluster) broadcastUpdate() ***REMOVED***
	peers := make([]*api.Peer, 0, len(c.members))
	for _, m := range c.members ***REMOVED***
		peers = append(peers, &api.Peer***REMOVED***
			NodeID: m.NodeID,
			Addr:   m.Addr,
		***REMOVED***)
	***REMOVED***
	c.PeersBroadcast.Publish(peers)
***REMOVED***

// AddMember adds a node to the Cluster Memberlist.
func (c *Cluster) AddMember(member *Member) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.removed[member.RaftID] ***REMOVED***
		return ErrIDRemoved
	***REMOVED***

	c.members[member.RaftID] = member

	c.broadcastUpdate()
	return nil
***REMOVED***

// RemoveMember removes a node from the Cluster Memberlist, and adds it to
// the removed list.
func (c *Cluster) RemoveMember(id uint64) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	c.removed[id] = true

	return c.clearMember(id)
***REMOVED***

// UpdateMember updates member address.
func (c *Cluster) UpdateMember(id uint64, m *api.RaftMember) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.removed[id] ***REMOVED***
		return ErrIDRemoved
	***REMOVED***

	oldMember, ok := c.members[id]
	if !ok ***REMOVED***
		return ErrIDNotFound
	***REMOVED***

	if oldMember.NodeID != m.NodeID ***REMOVED***
		// Should never happen; this is a sanity check
		return errors.New("node ID mismatch match on node update")
	***REMOVED***

	if oldMember.Addr == m.Addr ***REMOVED***
		// nothing to do
		return nil
	***REMOVED***
	oldMember.RaftMember = m
	c.broadcastUpdate()
	return nil
***REMOVED***

// ClearMember removes a node from the Cluster Memberlist, but does NOT add it
// to the removed list.
func (c *Cluster) ClearMember(id uint64) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.clearMember(id)
***REMOVED***

func (c *Cluster) clearMember(id uint64) error ***REMOVED***
	if _, ok := c.members[id]; ok ***REMOVED***
		delete(c.members, id)
		c.broadcastUpdate()
	***REMOVED***
	return nil
***REMOVED***

// IsIDRemoved checks if a Member is in the remove set.
func (c *Cluster) IsIDRemoved(id uint64) bool ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.removed[id]
***REMOVED***

// Clear resets the list of active Members and removed Members.
func (c *Cluster) Clear() ***REMOVED***
	c.mu.Lock()

	c.members = make(map[uint64]*Member)
	c.removed = make(map[uint64]bool)
	c.mu.Unlock()
***REMOVED***

// ValidateConfigurationChange takes a proposed ConfChange and
// ensures that it is valid.
func (c *Cluster) ValidateConfigurationChange(cc raftpb.ConfChange) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.removed[cc.NodeID] ***REMOVED***
		return ErrIDRemoved
	***REMOVED***
	switch cc.Type ***REMOVED***
	case raftpb.ConfChangeAddNode:
		if c.members[cc.NodeID] != nil ***REMOVED***
			return ErrIDExists
		***REMOVED***
	case raftpb.ConfChangeRemoveNode:
		if c.members[cc.NodeID] == nil ***REMOVED***
			return ErrIDNotFound
		***REMOVED***
	case raftpb.ConfChangeUpdateNode:
		if c.members[cc.NodeID] == nil ***REMOVED***
			return ErrIDNotFound
		***REMOVED***
	default:
		return ErrConfigChangeInvalid
	***REMOVED***
	m := &api.RaftMember***REMOVED******REMOVED***
	if err := proto.Unmarshal(cc.Context, m); err != nil ***REMOVED***
		return ErrCannotUnmarshalConfig
	***REMOVED***
	return nil
***REMOVED***
