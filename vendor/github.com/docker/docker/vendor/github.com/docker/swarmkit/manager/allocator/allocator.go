package allocator

import (
	"sync"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/go-events"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
)

// Allocator controls how the allocation stage in the manager is handled.
type Allocator struct ***REMOVED***
	// The manager store.
	store *store.MemoryStore

	// the ballot used to synchronize across all allocators to ensure
	// all of them have completed their respective allocations so that the
	// task can be moved to ALLOCATED state.
	taskBallot *taskBallot

	// context for the network allocator that will be needed by
	// network allocator.
	netCtx *networkContext

	// stopChan signals to the allocator to stop running.
	stopChan chan struct***REMOVED******REMOVED***
	// doneChan is closed when the allocator is finished running.
	doneChan chan struct***REMOVED******REMOVED***

	// pluginGetter provides access to docker's plugin inventory.
	pluginGetter plugingetter.PluginGetter
***REMOVED***

// taskBallot controls how the voting for task allocation is
// coordinated b/w different allocators. This the only structure that
// will be written by all allocator goroutines concurrently. Hence the
// mutex.
type taskBallot struct ***REMOVED***
	sync.Mutex

	// List of registered voters who have to cast their vote to
	// indicate their allocation complete
	voters []string

	// List of votes collected for every task so far from different voters.
	votes map[string][]string
***REMOVED***

// allocActor controls the various phases in the lifecycle of one kind of allocator.
type allocActor struct ***REMOVED***
	// Task voter identity of the allocator.
	taskVoter string

	// Action routine which is called for every event that the
	// allocator received.
	action func(context.Context, events.Event)

	// Init routine which is called during the initialization of
	// the allocator.
	init func(ctx context.Context) error
***REMOVED***

// New returns a new instance of Allocator for use during allocation
// stage of the manager.
func New(store *store.MemoryStore, pg plugingetter.PluginGetter) (*Allocator, error) ***REMOVED***
	a := &Allocator***REMOVED***
		store: store,
		taskBallot: &taskBallot***REMOVED***
			votes: make(map[string][]string),
		***REMOVED***,
		stopChan:     make(chan struct***REMOVED******REMOVED***),
		doneChan:     make(chan struct***REMOVED******REMOVED***),
		pluginGetter: pg,
	***REMOVED***

	return a, nil
***REMOVED***

// Run starts all allocator go-routines and waits for Stop to be called.
func (a *Allocator) Run(ctx context.Context) error ***REMOVED***
	// Setup cancel context for all goroutines to use.
	ctx, cancel := context.WithCancel(ctx)
	var (
		wg     sync.WaitGroup
		actors []func() error
	)

	defer func() ***REMOVED***
		cancel()
		wg.Wait()
		close(a.doneChan)
	***REMOVED***()

	for _, aa := range []allocActor***REMOVED***
		***REMOVED***
			taskVoter: networkVoter,
			init:      a.doNetworkInit,
			action:    a.doNetworkAlloc,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		if aa.taskVoter != "" ***REMOVED***
			a.registerToVote(aa.taskVoter)
		***REMOVED***

		// Assign a pointer for variable capture
		aaPtr := &aa
		actor := func() error ***REMOVED***
			wg.Add(1)
			defer wg.Done()

			// init might return an allocator specific context
			// which is a child of the passed in context to hold
			// allocator specific state
			watch, watchCancel, err := a.init(ctx, aaPtr)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			wg.Add(1)
			go func(watch <-chan events.Event, watchCancel func()) ***REMOVED***
				defer func() ***REMOVED***
					wg.Done()
					watchCancel()
				***REMOVED***()
				a.run(ctx, *aaPtr, watch)
			***REMOVED***(watch, watchCancel)
			return nil
		***REMOVED***

		actors = append(actors, actor)
	***REMOVED***

	for _, actor := range actors ***REMOVED***
		if err := actor(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	<-a.stopChan
	return nil
***REMOVED***

// Stop stops the allocator
func (a *Allocator) Stop() ***REMOVED***
	close(a.stopChan)
	// Wait for all allocator goroutines to truly exit
	<-a.doneChan
***REMOVED***

func (a *Allocator) init(ctx context.Context, aa *allocActor) (<-chan events.Event, func(), error) ***REMOVED***
	watch, watchCancel := state.Watch(a.store.WatchQueue(),
		api.EventCreateNetwork***REMOVED******REMOVED***,
		api.EventDeleteNetwork***REMOVED******REMOVED***,
		api.EventCreateService***REMOVED******REMOVED***,
		api.EventUpdateService***REMOVED******REMOVED***,
		api.EventDeleteService***REMOVED******REMOVED***,
		api.EventCreateTask***REMOVED******REMOVED***,
		api.EventUpdateTask***REMOVED******REMOVED***,
		api.EventDeleteTask***REMOVED******REMOVED***,
		api.EventCreateNode***REMOVED******REMOVED***,
		api.EventUpdateNode***REMOVED******REMOVED***,
		api.EventDeleteNode***REMOVED******REMOVED***,
		state.EventCommit***REMOVED******REMOVED***,
	)

	if err := aa.init(ctx); err != nil ***REMOVED***
		watchCancel()
		return nil, nil, err
	***REMOVED***

	return watch, watchCancel, nil
***REMOVED***

func (a *Allocator) run(ctx context.Context, aa allocActor, watch <-chan events.Event) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case ev, ok := <-watch:
			if !ok ***REMOVED***
				return
			***REMOVED***

			aa.action(ctx, ev)
		case <-ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (a *Allocator) registerToVote(name string) ***REMOVED***
	a.taskBallot.Lock()
	defer a.taskBallot.Unlock()

	a.taskBallot.voters = append(a.taskBallot.voters, name)
***REMOVED***

func (a *Allocator) taskAllocateVote(voter string, id string) bool ***REMOVED***
	a.taskBallot.Lock()
	defer a.taskBallot.Unlock()

	// If voter has already voted, return false
	for _, v := range a.taskBallot.votes[id] ***REMOVED***
		// check if voter is in x
		if v == voter ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	a.taskBallot.votes[id] = append(a.taskBallot.votes[id], voter)

	// We haven't gotten enough votes yet
	if len(a.taskBallot.voters) > len(a.taskBallot.votes[id]) ***REMOVED***
		return false
	***REMOVED***

nextVoter:
	for _, voter := range a.taskBallot.voters ***REMOVED***
		for _, vote := range a.taskBallot.votes[id] ***REMOVED***
			if voter == vote ***REMOVED***
				continue nextVoter
			***REMOVED***
		***REMOVED***

		// Not every registered voter has registered a vote.
		return false
	***REMOVED***

	return true
***REMOVED***
