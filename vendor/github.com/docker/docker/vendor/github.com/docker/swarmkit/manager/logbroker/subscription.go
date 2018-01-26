package logbroker

import (
	"fmt"
	"strings"
	"sync"

	events "github.com/docker/go-events"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/watch"
	"golang.org/x/net/context"
)

type subscription struct ***REMOVED***
	mu sync.RWMutex
	wg sync.WaitGroup

	store   *store.MemoryStore
	message *api.SubscriptionMessage
	changed *watch.Queue

	ctx    context.Context
	cancel context.CancelFunc

	errors       []error
	nodes        map[string]struct***REMOVED******REMOVED***
	pendingTasks map[string]struct***REMOVED******REMOVED***
***REMOVED***

func newSubscription(store *store.MemoryStore, message *api.SubscriptionMessage, changed *watch.Queue) *subscription ***REMOVED***
	return &subscription***REMOVED***
		store:        store,
		message:      message,
		changed:      changed,
		nodes:        make(map[string]struct***REMOVED******REMOVED***),
		pendingTasks: make(map[string]struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

func (s *subscription) follow() bool ***REMOVED***
	return s.message.Options != nil && s.message.Options.Follow
***REMOVED***

func (s *subscription) Contains(nodeID string) bool ***REMOVED***
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.nodes[nodeID]
	return ok
***REMOVED***

func (s *subscription) Nodes() []string ***REMOVED***
	s.mu.RLock()
	defer s.mu.RUnlock()

	nodes := make([]string, 0, len(s.nodes))
	for node := range s.nodes ***REMOVED***
		nodes = append(nodes, node)
	***REMOVED***
	return nodes
***REMOVED***

func (s *subscription) Run(ctx context.Context) ***REMOVED***
	s.ctx, s.cancel = context.WithCancel(ctx)

	if s.follow() ***REMOVED***
		wq := s.store.WatchQueue()
		ch, cancel := state.Watch(wq, api.EventCreateTask***REMOVED******REMOVED***, api.EventUpdateTask***REMOVED******REMOVED***)
		go func() ***REMOVED***
			defer cancel()
			s.watch(ch)
		***REMOVED***()
	***REMOVED***

	s.match()
***REMOVED***

func (s *subscription) Stop() ***REMOVED***
	if s.cancel != nil ***REMOVED***
		s.cancel()
	***REMOVED***
***REMOVED***

func (s *subscription) Wait(ctx context.Context) <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	// Follow subscriptions never end
	if s.follow() ***REMOVED***
		return nil
	***REMOVED***

	ch := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		defer close(ch)
		s.wg.Wait()
	***REMOVED***()
	return ch
***REMOVED***

func (s *subscription) Done(nodeID string, err error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	if err != nil ***REMOVED***
		s.errors = append(s.errors, err)
	***REMOVED***

	if s.follow() ***REMOVED***
		return
	***REMOVED***

	if _, ok := s.nodes[nodeID]; !ok ***REMOVED***
		return
	***REMOVED***

	delete(s.nodes, nodeID)
	s.wg.Done()
***REMOVED***

func (s *subscription) Err() error ***REMOVED***
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.errors) == 0 && len(s.pendingTasks) == 0 ***REMOVED***
		return nil
	***REMOVED***

	messages := make([]string, 0, len(s.errors))
	for _, err := range s.errors ***REMOVED***
		messages = append(messages, err.Error())
	***REMOVED***
	for t := range s.pendingTasks ***REMOVED***
		messages = append(messages, fmt.Sprintf("task %s has not been scheduled", t))
	***REMOVED***

	return fmt.Errorf("warning: incomplete log stream. some logs could not be retrieved for the following reasons: %s", strings.Join(messages, ", "))
***REMOVED***

func (s *subscription) Close() ***REMOVED***
	s.mu.Lock()
	s.message.Close = true
	s.mu.Unlock()
***REMOVED***

func (s *subscription) Closed() bool ***REMOVED***
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.message.Close
***REMOVED***

func (s *subscription) match() ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	add := func(t *api.Task) ***REMOVED***
		if t.NodeID == "" ***REMOVED***
			s.pendingTasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			return
		***REMOVED***
		if _, ok := s.nodes[t.NodeID]; !ok ***REMOVED***
			s.nodes[t.NodeID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			s.wg.Add(1)
		***REMOVED***
	***REMOVED***

	s.store.View(func(tx store.ReadTx) ***REMOVED***
		for _, nid := range s.message.Selector.NodeIDs ***REMOVED***
			s.nodes[nid] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***

		for _, tid := range s.message.Selector.TaskIDs ***REMOVED***
			if task := store.GetTask(tx, tid); task != nil ***REMOVED***
				add(task)
			***REMOVED***
		***REMOVED***

		for _, sid := range s.message.Selector.ServiceIDs ***REMOVED***
			tasks, err := store.FindTasks(tx, store.ByServiceID(sid))
			if err != nil ***REMOVED***
				log.L.Warning(err)
				continue
			***REMOVED***
			for _, task := range tasks ***REMOVED***
				// if we're not following, don't add tasks that aren't running yet
				if !s.follow() && task.Status.State < api.TaskStateRunning ***REMOVED***
					continue
				***REMOVED***
				add(task)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func (s *subscription) watch(ch <-chan events.Event) error ***REMOVED***
	matchTasks := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, tid := range s.message.Selector.TaskIDs ***REMOVED***
		matchTasks[tid] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	matchServices := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, sid := range s.message.Selector.ServiceIDs ***REMOVED***
		matchServices[sid] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	add := func(t *api.Task) ***REMOVED***
		s.mu.Lock()
		defer s.mu.Unlock()

		// Un-allocated task.
		if t.NodeID == "" ***REMOVED***
			s.pendingTasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			return
		***REMOVED***

		delete(s.pendingTasks, t.ID)
		if _, ok := s.nodes[t.NodeID]; !ok ***REMOVED***
			s.nodes[t.NodeID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			s.changed.Publish(s)
		***REMOVED***
	***REMOVED***

	for ***REMOVED***
		var t *api.Task
		select ***REMOVED***
		case <-s.ctx.Done():
			return s.ctx.Err()
		case event := <-ch:
			switch v := event.(type) ***REMOVED***
			case api.EventCreateTask:
				t = v.Task
			case api.EventUpdateTask:
				t = v.Task
			***REMOVED***
		***REMOVED***

		if t == nil ***REMOVED***
			panic("received invalid task from the watch queue")
		***REMOVED***

		if _, ok := matchTasks[t.ID]; ok ***REMOVED***
			add(t)
		***REMOVED***
		if _, ok := matchServices[t.ServiceID]; ok ***REMOVED***
			add(t)
		***REMOVED***
	***REMOVED***
***REMOVED***
