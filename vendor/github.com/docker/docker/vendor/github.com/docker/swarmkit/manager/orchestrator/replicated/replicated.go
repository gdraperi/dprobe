package replicated

import (
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/orchestrator/restart"
	"github.com/docker/swarmkit/manager/orchestrator/update"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
)

// An Orchestrator runs a reconciliation loop to create and destroy
// tasks as necessary for the replicated services.
type Orchestrator struct ***REMOVED***
	store *store.MemoryStore

	reconcileServices map[string]*api.Service
	restartTasks      map[string]struct***REMOVED******REMOVED***

	// stopChan signals to the state machine to stop running.
	stopChan chan struct***REMOVED******REMOVED***
	// doneChan is closed when the state machine terminates.
	doneChan chan struct***REMOVED******REMOVED***

	updater  *update.Supervisor
	restarts *restart.Supervisor

	cluster *api.Cluster // local cluster instance
***REMOVED***

// NewReplicatedOrchestrator creates a new replicated Orchestrator.
func NewReplicatedOrchestrator(store *store.MemoryStore) *Orchestrator ***REMOVED***
	restartSupervisor := restart.NewSupervisor(store)
	updater := update.NewSupervisor(store, restartSupervisor)
	return &Orchestrator***REMOVED***
		store:             store,
		stopChan:          make(chan struct***REMOVED******REMOVED***),
		doneChan:          make(chan struct***REMOVED******REMOVED***),
		reconcileServices: make(map[string]*api.Service),
		restartTasks:      make(map[string]struct***REMOVED******REMOVED***),
		updater:           updater,
		restarts:          restartSupervisor,
	***REMOVED***
***REMOVED***

// Run contains the orchestrator event loop. It runs until Stop is called.
func (r *Orchestrator) Run(ctx context.Context) error ***REMOVED***
	defer close(r.doneChan)

	// Watch changes to services and tasks
	queue := r.store.WatchQueue()
	watcher, cancel := queue.Watch()
	defer cancel()

	// Balance existing services and drain initial tasks attached to invalid
	// nodes
	var err error
	r.store.View(func(readTx store.ReadTx) ***REMOVED***
		if err = r.initTasks(ctx, readTx); err != nil ***REMOVED***
			return
		***REMOVED***

		if err = r.initServices(readTx); err != nil ***REMOVED***
			return
		***REMOVED***

		if err = r.initCluster(readTx); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.tick(ctx)

	for ***REMOVED***
		select ***REMOVED***
		case event := <-watcher:
			// TODO(stevvooe): Use ctx to limit running time of operation.
			r.handleTaskEvent(ctx, event)
			r.handleServiceEvent(ctx, event)
			switch v := event.(type) ***REMOVED***
			case state.EventCommit:
				r.tick(ctx)
			case api.EventUpdateCluster:
				r.cluster = v.Cluster
			***REMOVED***
		case <-r.stopChan:
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stop stops the orchestrator.
func (r *Orchestrator) Stop() ***REMOVED***
	close(r.stopChan)
	<-r.doneChan
	r.updater.CancelAll()
	r.restarts.CancelAll()
***REMOVED***

func (r *Orchestrator) tick(ctx context.Context) ***REMOVED***
	// tickTasks must be called first, so we respond to task-level changes
	// before performing service reconciliation.
	r.tickTasks(ctx)
	r.tickServices(ctx)
***REMOVED***
