package dispatcher

import (
	"sync"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/manager/dispatcher/heartbeat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const rateLimitCount = 3

type registeredNode struct ***REMOVED***
	SessionID  string
	Heartbeat  *heartbeat.Heartbeat
	Registered time.Time
	Attempts   int
	Node       *api.Node
	Disconnect chan struct***REMOVED******REMOVED*** // signal to disconnect
	mu         sync.Mutex
***REMOVED***

// checkSessionID determines if the SessionID has changed and returns the
// appropriate GRPC error code.
//
// This may not belong here in the future.
func (rn *registeredNode) checkSessionID(sessionID string) error ***REMOVED***
	rn.mu.Lock()
	defer rn.mu.Unlock()

	// Before each message send, we need to check the nodes sessionID hasn't
	// changed. If it has, we will the stream and make the node
	// re-register.
	if sessionID == "" || rn.SessionID != sessionID ***REMOVED***
		return status.Errorf(codes.InvalidArgument, ErrSessionInvalid.Error())
	***REMOVED***

	return nil
***REMOVED***

type nodeStore struct ***REMOVED***
	periodChooser                *periodChooser
	gracePeriodMultiplierNormal  time.Duration
	gracePeriodMultiplierUnknown time.Duration
	rateLimitPeriod              time.Duration
	nodes                        map[string]*registeredNode
	mu                           sync.RWMutex
***REMOVED***

func newNodeStore(hbPeriod, hbEpsilon time.Duration, graceMultiplier int, rateLimitPeriod time.Duration) *nodeStore ***REMOVED***
	return &nodeStore***REMOVED***
		nodes:                        make(map[string]*registeredNode),
		periodChooser:                newPeriodChooser(hbPeriod, hbEpsilon),
		gracePeriodMultiplierNormal:  time.Duration(graceMultiplier),
		gracePeriodMultiplierUnknown: time.Duration(graceMultiplier) * 2,
		rateLimitPeriod:              rateLimitPeriod,
	***REMOVED***
***REMOVED***

func (s *nodeStore) updatePeriod(hbPeriod, hbEpsilon time.Duration, gracePeriodMultiplier int) ***REMOVED***
	s.mu.Lock()
	s.periodChooser = newPeriodChooser(hbPeriod, hbEpsilon)
	s.gracePeriodMultiplierNormal = time.Duration(gracePeriodMultiplier)
	s.gracePeriodMultiplierUnknown = s.gracePeriodMultiplierNormal * 2
	s.mu.Unlock()
***REMOVED***

func (s *nodeStore) Len() int ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.nodes)
***REMOVED***

func (s *nodeStore) AddUnknown(n *api.Node, expireFunc func()) error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	rn := &registeredNode***REMOVED***
		Node: n,
	***REMOVED***
	s.nodes[n.ID] = rn
	rn.Heartbeat = heartbeat.New(s.periodChooser.Choose()*s.gracePeriodMultiplierUnknown, expireFunc)
	return nil
***REMOVED***

// CheckRateLimit returns error if node with specified id is allowed to re-register
// again.
func (s *nodeStore) CheckRateLimit(id string) error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if existRn, ok := s.nodes[id]; ok ***REMOVED***
		if time.Since(existRn.Registered) > s.rateLimitPeriod ***REMOVED***
			existRn.Attempts = 0
		***REMOVED***
		existRn.Attempts++
		if existRn.Attempts > rateLimitCount ***REMOVED***
			return status.Errorf(codes.Unavailable, "node %s exceeded rate limit count of registrations", id)
		***REMOVED***
		existRn.Registered = time.Now()
	***REMOVED***
	return nil
***REMOVED***

// Add adds new node and returns it, it replaces existing without notification.
func (s *nodeStore) Add(n *api.Node, expireFunc func()) *registeredNode ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	var attempts int
	var registered time.Time
	if existRn, ok := s.nodes[n.ID]; ok ***REMOVED***
		attempts = existRn.Attempts
		registered = existRn.Registered
		existRn.Heartbeat.Stop()
		delete(s.nodes, n.ID)
	***REMOVED***
	if registered.IsZero() ***REMOVED***
		registered = time.Now()
	***REMOVED***
	rn := &registeredNode***REMOVED***
		SessionID:  identity.NewID(), // session ID is local to the dispatcher.
		Node:       n,
		Registered: registered,
		Attempts:   attempts,
		Disconnect: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
	s.nodes[n.ID] = rn
	rn.Heartbeat = heartbeat.New(s.periodChooser.Choose()*s.gracePeriodMultiplierNormal, expireFunc)
	return rn
***REMOVED***

func (s *nodeStore) Get(id string) (*registeredNode, error) ***REMOVED***
	s.mu.RLock()
	rn, ok := s.nodes[id]
	s.mu.RUnlock()
	if !ok ***REMOVED***
		return nil, status.Errorf(codes.NotFound, ErrNodeNotRegistered.Error())
	***REMOVED***
	return rn, nil
***REMOVED***

func (s *nodeStore) GetWithSession(id, sid string) (*registeredNode, error) ***REMOVED***
	s.mu.RLock()
	rn, ok := s.nodes[id]
	s.mu.RUnlock()
	if !ok ***REMOVED***
		return nil, status.Errorf(codes.NotFound, ErrNodeNotRegistered.Error())
	***REMOVED***
	return rn, rn.checkSessionID(sid)
***REMOVED***

func (s *nodeStore) Heartbeat(id, sid string) (time.Duration, error) ***REMOVED***
	rn, err := s.GetWithSession(id, sid)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	period := s.periodChooser.Choose() // base period for node
	grace := period * time.Duration(s.gracePeriodMultiplierNormal)
	rn.mu.Lock()
	rn.Heartbeat.Update(grace)
	rn.Heartbeat.Beat()
	rn.mu.Unlock()
	return period, nil
***REMOVED***

func (s *nodeStore) Delete(id string) *registeredNode ***REMOVED***
	s.mu.Lock()
	var node *registeredNode
	if rn, ok := s.nodes[id]; ok ***REMOVED***
		delete(s.nodes, id)
		rn.Heartbeat.Stop()
		node = rn
	***REMOVED***
	s.mu.Unlock()
	return node
***REMOVED***

func (s *nodeStore) Disconnect(id string) ***REMOVED***
	s.mu.Lock()
	if rn, ok := s.nodes[id]; ok ***REMOVED***
		close(rn.Disconnect)
		rn.Heartbeat.Stop()
	***REMOVED***
	s.mu.Unlock()
***REMOVED***

// Clean removes all nodes and stops their heartbeats.
// It's equivalent to invalidate all sessions.
func (s *nodeStore) Clean() ***REMOVED***
	s.mu.Lock()
	for _, rn := range s.nodes ***REMOVED***
		rn.Heartbeat.Stop()
	***REMOVED***
	s.nodes = make(map[string]*registeredNode)
	s.mu.Unlock()
***REMOVED***
