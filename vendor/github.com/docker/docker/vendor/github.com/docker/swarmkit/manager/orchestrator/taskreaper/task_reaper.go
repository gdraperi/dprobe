package taskreaper

import (
	"sort"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/orchestrator"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
)

const (
	// maxDirty is the size threshold for running a task pruning operation.
	maxDirty = 1000
	// reaperBatchingInterval is how often to prune old tasks.
	reaperBatchingInterval = 250 * time.Millisecond
)

// A TaskReaper deletes old tasks when more than TaskHistoryRetentionLimit tasks
// exist for the same service/instance or service/nodeid combination.
type TaskReaper struct ***REMOVED***
	store *store.MemoryStore

	// taskHistory is the number of tasks to keep
	taskHistory int64

	// List of slot tubles to be inspected for task history cleanup.
	dirty map[orchestrator.SlotTuple]struct***REMOVED******REMOVED***

	// List of tasks collected for cleanup, which includes two kinds of tasks
	// - serviceless orphaned tasks
	// - tasks with desired state REMOVE that have already been shut down
	cleanup  []string
	stopChan chan struct***REMOVED******REMOVED***
	doneChan chan struct***REMOVED******REMOVED***
***REMOVED***

// New creates a new TaskReaper.
func New(store *store.MemoryStore) *TaskReaper ***REMOVED***
	return &TaskReaper***REMOVED***
		store:    store,
		dirty:    make(map[orchestrator.SlotTuple]struct***REMOVED******REMOVED***),
		stopChan: make(chan struct***REMOVED******REMOVED***),
		doneChan: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// Run is the TaskReaper's watch loop which collects candidates for cleanup.
// Task history is mainly used in task restarts but is also available for administrative purposes.
// Note that the task history is stored per-slot-per-service for replicated services
// and per-node-per-service for global services. History does not apply to serviceless
// since they are not attached to a service. In addition, the TaskReaper watch loop is also
// responsible for cleaning up tasks associated with slots that were removed as part of
// service scale down or service removal.
func (tr *TaskReaper) Run(ctx context.Context) ***REMOVED***
	watcher, watchCancel := state.Watch(tr.store.WatchQueue(), api.EventCreateTask***REMOVED******REMOVED***, api.EventUpdateTask***REMOVED******REMOVED***, api.EventUpdateCluster***REMOVED******REMOVED***)

	defer func() ***REMOVED***
		close(tr.doneChan)
		watchCancel()
	***REMOVED***()

	var orphanedTasks []*api.Task
	var removeTasks []*api.Task
	tr.store.View(func(readTx store.ReadTx) ***REMOVED***
		var err error

		clusters, err := store.FindClusters(readTx, store.ByName(store.DefaultClusterName))
		if err == nil && len(clusters) == 1 ***REMOVED***
			tr.taskHistory = clusters[0].Spec.Orchestration.TaskHistoryRetentionLimit
		***REMOVED***

		// On startup, scan the entire store and inspect orphaned tasks from previous life.
		orphanedTasks, err = store.FindTasks(readTx, store.ByTaskState(api.TaskStateOrphaned))
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("failed to find Orphaned tasks in task reaper init")
		***REMOVED***
		removeTasks, err = store.FindTasks(readTx, store.ByDesiredState(api.TaskStateRemove))
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("failed to find tasks with desired state REMOVE in task reaper init")
		***REMOVED***
	***REMOVED***)

	if len(orphanedTasks)+len(removeTasks) > 0 ***REMOVED***
		for _, t := range orphanedTasks ***REMOVED***
			// Do not reap service tasks immediately.
			// Let them go through the regular history cleanup process
			// of checking TaskHistoryRetentionLimit.
			if t.ServiceID != "" ***REMOVED***
				continue
			***REMOVED***

			// Serviceless tasks can be cleaned up right away since they are not attached to a service.
			tr.cleanup = append(tr.cleanup, t.ID)
		***REMOVED***
		// tasks with desired state REMOVE that have progressed beyond COMPLETE can be cleaned up
		// right away
		for _, t := range removeTasks ***REMOVED***
			if t.Status.State >= api.TaskStateCompleted ***REMOVED***
				tr.cleanup = append(tr.cleanup, t.ID)
			***REMOVED***
		***REMOVED***
		// Clean up tasks in 'cleanup' right away
		if len(tr.cleanup) > 0 ***REMOVED***
			tr.tick()
		***REMOVED***
	***REMOVED***

	// Clean up when we hit TaskHistoryRetentionLimit or when the timer expires,
	// whichever happens first.
	timer := time.NewTimer(reaperBatchingInterval)

	// Watch for:
	// 1. EventCreateTask for cleaning slots, which is the best time to cleanup that node/slot.
	// 2. EventUpdateTask for cleaning
	//    - serviceless orphaned tasks (when orchestrator updates the task status to ORPHANED)
	//    - tasks which have desired state REMOVE and have been shut down by the agent
	//      (these are tasks which are associated with slots removed as part of service
	//       remove or scale down)
	// 3. EventUpdateCluster for TaskHistoryRetentionLimit update.
	for ***REMOVED***
		select ***REMOVED***
		case event := <-watcher:
			switch v := event.(type) ***REMOVED***
			case api.EventCreateTask:
				t := v.Task
				tr.dirty[orchestrator.SlotTuple***REMOVED***
					Slot:      t.Slot,
					ServiceID: t.ServiceID,
					NodeID:    t.NodeID,
				***REMOVED***] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			case api.EventUpdateTask:
				t := v.Task
				// add serviceless orphaned tasks
				if t.Status.State >= api.TaskStateOrphaned && t.ServiceID == "" ***REMOVED***
					tr.cleanup = append(tr.cleanup, t.ID)
				***REMOVED***
				// add tasks that have progressed beyond COMPLETE and have desired state REMOVE. These
				// tasks are associated with slots that were removed as part of a service scale down
				// or service removal.
				if t.DesiredState == api.TaskStateRemove && t.Status.State >= api.TaskStateCompleted ***REMOVED***
					tr.cleanup = append(tr.cleanup, t.ID)
				***REMOVED***
			case api.EventUpdateCluster:
				tr.taskHistory = v.Cluster.Spec.Orchestration.TaskHistoryRetentionLimit
			***REMOVED***

			if len(tr.dirty)+len(tr.cleanup) > maxDirty ***REMOVED***
				timer.Stop()
				tr.tick()
			***REMOVED*** else ***REMOVED***
				timer.Reset(reaperBatchingInterval)
			***REMOVED***
		case <-timer.C:
			timer.Stop()
			tr.tick()
		case <-tr.stopChan:
			timer.Stop()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// tick performs task history cleanup.
func (tr *TaskReaper) tick() ***REMOVED***
	if len(tr.dirty) == 0 && len(tr.cleanup) == 0 ***REMOVED***
		return
	***REMOVED***

	defer func() ***REMOVED***
		tr.cleanup = nil
	***REMOVED***()

	deleteTasks := make(map[string]struct***REMOVED******REMOVED***)
	for _, tID := range tr.cleanup ***REMOVED***
		deleteTasks[tID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	// Check history of dirty tasks for cleanup.
	tr.store.View(func(tx store.ReadTx) ***REMOVED***
		for dirty := range tr.dirty ***REMOVED***
			service := store.GetService(tx, dirty.ServiceID)
			if service == nil ***REMOVED***
				continue
			***REMOVED***

			taskHistory := tr.taskHistory

			// If MaxAttempts is set, keep at least one more than
			// that number of tasks (this overrides TaskHistoryRetentionLimit).
			// This is necessary to reconstruct restart history when the orchestrator starts up.
			// TODO(aaronl): Consider hiding tasks beyond the normal
			// retention limit in the UI.
			// TODO(aaronl): There are some ways to cut down the
			// number of retained tasks at the cost of more
			// complexity:
			//   - Don't force retention of tasks with an older spec
			//     version.
			//   - Don't force retention of tasks outside of the
			//     time window configured for restart lookback.
			if service.Spec.Task.Restart != nil && service.Spec.Task.Restart.MaxAttempts > 0 ***REMOVED***
				taskHistory = int64(service.Spec.Task.Restart.MaxAttempts) + 1
			***REMOVED***

			// Negative value for TaskHistoryRetentionLimit is an indication to never clean up task history.
			if taskHistory < 0 ***REMOVED***
				continue
			***REMOVED***

			var historicTasks []*api.Task

			switch service.Spec.GetMode().(type) ***REMOVED***
			case *api.ServiceSpec_Replicated:
				// Clean out the slot for which we received EventCreateTask.
				var err error
				historicTasks, err = store.FindTasks(tx, store.BySlot(dirty.ServiceID, dirty.Slot))
				if err != nil ***REMOVED***
					continue
				***REMOVED***

			case *api.ServiceSpec_Global:
				// Clean out the node history in case of global services.
				tasksByNode, err := store.FindTasks(tx, store.ByNodeID(dirty.NodeID))
				if err != nil ***REMOVED***
					continue
				***REMOVED***

				for _, t := range tasksByNode ***REMOVED***
					if t.ServiceID == dirty.ServiceID ***REMOVED***
						historicTasks = append(historicTasks, t)
					***REMOVED***
				***REMOVED***
			***REMOVED***

			if int64(len(historicTasks)) <= taskHistory ***REMOVED***
				continue
			***REMOVED***

			// TODO(aaronl): This could filter for non-running tasks and use quickselect
			// instead of sorting the whole slice.
			// TODO(aaronl): This sort should really use lamport time instead of wall
			// clock time. We should store a Version in the Status field.
			sort.Sort(orchestrator.TasksByTimestamp(historicTasks))

			runningTasks := 0
			for _, t := range historicTasks ***REMOVED***
				if t.DesiredState <= api.TaskStateRunning || t.Status.State <= api.TaskStateRunning ***REMOVED***
					// Don't delete running tasks
					runningTasks++
					continue
				***REMOVED***

				deleteTasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

				taskHistory++
				if int64(len(historicTasks)) <= taskHistory ***REMOVED***
					break
				***REMOVED***
			***REMOVED***

			if runningTasks <= 1 ***REMOVED***
				delete(tr.dirty, dirty)
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	// Perform cleanup.
	if len(deleteTasks) > 0 ***REMOVED***
		tr.store.Batch(func(batch *store.Batch) error ***REMOVED***
			for taskID := range deleteTasks ***REMOVED***
				batch.Update(func(tx store.Tx) error ***REMOVED***
					return store.DeleteTask(tx, taskID)
				***REMOVED***)
			***REMOVED***
			return nil
		***REMOVED***)
	***REMOVED***
***REMOVED***

// Stop stops the TaskReaper and waits for the main loop to exit.
func (tr *TaskReaper) Stop() ***REMOVED***
	// TODO(dperny) calling stop on the task reaper twice will cause a panic
	// because we try to close a channel that will already have been closed.
	close(tr.stopChan)
	<-tr.doneChan
***REMOVED***
