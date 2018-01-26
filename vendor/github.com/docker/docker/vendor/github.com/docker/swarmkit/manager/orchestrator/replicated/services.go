package replicated

import (
	"sort"

	"github.com/docker/go-events"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/orchestrator"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
)

// This file provices service-level orchestration. It observes changes to
// services and creates and destroys tasks as necessary to match the service
// specifications. This is different from task-level orchestration, which
// responds to changes in individual tasks (or nodes which run them).

func (r *Orchestrator) initCluster(readTx store.ReadTx) error ***REMOVED***
	clusters, err := store.FindClusters(readTx, store.ByName(store.DefaultClusterName))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(clusters) != 1 ***REMOVED***
		// we'll just pick it when it is created.
		return nil
	***REMOVED***

	r.cluster = clusters[0]
	return nil
***REMOVED***

func (r *Orchestrator) initServices(readTx store.ReadTx) error ***REMOVED***
	services, err := store.FindServices(readTx, store.All)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, s := range services ***REMOVED***
		if orchestrator.IsReplicatedService(s) ***REMOVED***
			r.reconcileServices[s.ID] = s
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *Orchestrator) handleServiceEvent(ctx context.Context, event events.Event) ***REMOVED***
	switch v := event.(type) ***REMOVED***
	case api.EventDeleteService:
		if !orchestrator.IsReplicatedService(v.Service) ***REMOVED***
			return
		***REMOVED***
		orchestrator.SetServiceTasksRemove(ctx, r.store, v.Service)
		r.restarts.ClearServiceHistory(v.Service.ID)
		delete(r.reconcileServices, v.Service.ID)
	case api.EventCreateService:
		if !orchestrator.IsReplicatedService(v.Service) ***REMOVED***
			return
		***REMOVED***
		r.reconcileServices[v.Service.ID] = v.Service
	case api.EventUpdateService:
		if !orchestrator.IsReplicatedService(v.Service) ***REMOVED***
			return
		***REMOVED***
		r.reconcileServices[v.Service.ID] = v.Service
	***REMOVED***
***REMOVED***

func (r *Orchestrator) tickServices(ctx context.Context) ***REMOVED***
	if len(r.reconcileServices) > 0 ***REMOVED***
		for _, s := range r.reconcileServices ***REMOVED***
			r.reconcile(ctx, s)
		***REMOVED***
		r.reconcileServices = make(map[string]*api.Service)
	***REMOVED***
***REMOVED***

func (r *Orchestrator) resolveService(ctx context.Context, task *api.Task) *api.Service ***REMOVED***
	if task.ServiceID == "" ***REMOVED***
		return nil
	***REMOVED***
	var service *api.Service
	r.store.View(func(tx store.ReadTx) ***REMOVED***
		service = store.GetService(tx, task.ServiceID)
	***REMOVED***)
	return service
***REMOVED***

// reconcile decides what actions must be taken depending on the number of
// specificed slots and actual running slots. If the actual running slots are
// fewer than what is requested, it creates new tasks. If the actual running
// slots are more than requested, then it decides which slots must be removed
// and sets desired state of those tasks to REMOVE (the actual removal is handled
// by the task reaper, after the agent shuts the tasks down).
func (r *Orchestrator) reconcile(ctx context.Context, service *api.Service) ***REMOVED***
	runningSlots, deadSlots, err := r.updatableAndDeadSlots(ctx, service)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("reconcile failed finding tasks")
		return
	***REMOVED***

	numSlots := len(runningSlots)

	slotsSlice := make([]orchestrator.Slot, 0, numSlots)
	for _, slot := range runningSlots ***REMOVED***
		slotsSlice = append(slotsSlice, slot)
	***REMOVED***

	deploy := service.Spec.GetMode().(*api.ServiceSpec_Replicated)
	specifiedSlots := deploy.Replicated.Replicas

	switch ***REMOVED***
	case specifiedSlots > uint64(numSlots):
		log.G(ctx).Debugf("Service %s was scaled up from %d to %d instances", service.ID, numSlots, specifiedSlots)
		// Update all current tasks then add missing tasks
		r.updater.Update(ctx, r.cluster, service, slotsSlice)
		err = r.store.Batch(func(batch *store.Batch) error ***REMOVED***
			r.addTasks(ctx, batch, service, runningSlots, deadSlots, specifiedSlots-uint64(numSlots))
			r.deleteTasksMap(ctx, batch, deadSlots)
			return nil
		***REMOVED***)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("reconcile batch failed")
		***REMOVED***

	case specifiedSlots < uint64(numSlots):
		// Update up to N tasks then remove the extra
		log.G(ctx).Debugf("Service %s was scaled down from %d to %d instances", service.ID, numSlots, specifiedSlots)

		// Preferentially remove tasks on the nodes that have the most
		// copies of this service, to leave a more balanced result.

		// First sort tasks such that tasks which are currently running
		// (in terms of observed state) appear before non-running tasks.
		// This will cause us to prefer to remove non-running tasks, all
		// other things being equal in terms of node balance.

		sort.Sort(slotsByRunningState(slotsSlice))

		// Assign each task an index that counts it as the nth copy of
		// of the service on its node (1, 2, 3, ...), and sort the
		// tasks by this counter value.

		slotsByNode := make(map[string]int)
		slotsWithIndices := make(slotsByIndex, 0, numSlots)

		for _, slot := range slotsSlice ***REMOVED***
			if len(slot) == 1 && slot[0].NodeID != "" ***REMOVED***
				slotsByNode[slot[0].NodeID]++
				slotsWithIndices = append(slotsWithIndices, slotWithIndex***REMOVED***slot: slot, index: slotsByNode[slot[0].NodeID]***REMOVED***)
			***REMOVED*** else ***REMOVED***
				slotsWithIndices = append(slotsWithIndices, slotWithIndex***REMOVED***slot: slot, index: -1***REMOVED***)
			***REMOVED***
		***REMOVED***

		sort.Sort(slotsWithIndices)

		sortedSlots := make([]orchestrator.Slot, 0, numSlots)
		for _, slot := range slotsWithIndices ***REMOVED***
			sortedSlots = append(sortedSlots, slot.slot)
		***REMOVED***

		r.updater.Update(ctx, r.cluster, service, sortedSlots[:specifiedSlots])
		err = r.store.Batch(func(batch *store.Batch) error ***REMOVED***
			r.deleteTasksMap(ctx, batch, deadSlots)
			// for all slots that we are removing, we set the desired state of those tasks
			// to REMOVE. Then, the agent is responsible for shutting them down, and the
			// task reaper is responsible for actually removing them from the store after
			// shutdown.
			r.setTasksDesiredState(ctx, batch, sortedSlots[specifiedSlots:], api.TaskStateRemove)
			return nil
		***REMOVED***)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("reconcile batch failed")
		***REMOVED***

	case specifiedSlots == uint64(numSlots):
		err = r.store.Batch(func(batch *store.Batch) error ***REMOVED***
			r.deleteTasksMap(ctx, batch, deadSlots)
			return nil
		***REMOVED***)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("reconcile batch failed")
		***REMOVED***
		// Simple update, no scaling - update all tasks.
		r.updater.Update(ctx, r.cluster, service, slotsSlice)
	***REMOVED***
***REMOVED***

func (r *Orchestrator) addTasks(ctx context.Context, batch *store.Batch, service *api.Service, runningSlots map[uint64]orchestrator.Slot, deadSlots map[uint64]orchestrator.Slot, count uint64) ***REMOVED***
	slot := uint64(0)
	for i := uint64(0); i < count; i++ ***REMOVED***
		// Find a slot number that is missing a running task
		for ***REMOVED***
			slot++
			if _, ok := runningSlots[slot]; !ok ***REMOVED***
				break
			***REMOVED***
		***REMOVED***

		delete(deadSlots, slot)
		err := batch.Update(func(tx store.Tx) error ***REMOVED***
			return store.CreateTask(tx, orchestrator.NewTask(r.cluster, service, slot, ""))
		***REMOVED***)
		if err != nil ***REMOVED***
			log.G(ctx).Errorf("Failed to create task: %v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// setTasksDesiredState sets the desired state for all tasks for the given slots to the
// requested state
func (r *Orchestrator) setTasksDesiredState(ctx context.Context, batch *store.Batch, slots []orchestrator.Slot, newDesiredState api.TaskState) ***REMOVED***
	for _, slot := range slots ***REMOVED***
		for _, t := range slot ***REMOVED***
			err := batch.Update(func(tx store.Tx) error ***REMOVED***
				// time travel is not allowed. if the current desired state is
				// above the one we're trying to go to we can't go backwards.
				// we have nothing to do and we should skip to the next task
				if t.DesiredState > newDesiredState ***REMOVED***
					// log a warning, though. we shouln't be trying to rewrite
					// a state to an earlier state
					log.G(ctx).Warnf(
						"cannot update task %v in desired state %v to an earlier desired state %v",
						t.ID, t.DesiredState, newDesiredState,
					)
					return nil
				***REMOVED***
				// update desired state
				t.DesiredState = newDesiredState

				return store.UpdateTask(tx, t)
			***REMOVED***)

			// log an error if we get one
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Errorf("failed to update task to %v", newDesiredState.String())
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *Orchestrator) deleteTasksMap(ctx context.Context, batch *store.Batch, slots map[uint64]orchestrator.Slot) ***REMOVED***
	for _, slot := range slots ***REMOVED***
		for _, t := range slot ***REMOVED***
			r.deleteTask(ctx, batch, t)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *Orchestrator) deleteTask(ctx context.Context, batch *store.Batch, t *api.Task) ***REMOVED***
	err := batch.Update(func(tx store.Tx) error ***REMOVED***
		return store.DeleteTask(tx, t.ID)
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("deleting task %s failed", t.ID)
	***REMOVED***
***REMOVED***

// IsRelatedService returns true if the service should be governed by this orchestrator
func (r *Orchestrator) IsRelatedService(service *api.Service) bool ***REMOVED***
	return orchestrator.IsReplicatedService(service)
***REMOVED***
