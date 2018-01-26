package constraintenforcer

import (
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/genericresource"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/constraint"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/protobuf/ptypes"
)

// ConstraintEnforcer watches for updates to nodes and shuts down tasks that no
// longer satisfy scheduling constraints or resource limits.
type ConstraintEnforcer struct ***REMOVED***
	store    *store.MemoryStore
	stopChan chan struct***REMOVED******REMOVED***
	doneChan chan struct***REMOVED******REMOVED***
***REMOVED***

// New creates a new ConstraintEnforcer.
func New(store *store.MemoryStore) *ConstraintEnforcer ***REMOVED***
	return &ConstraintEnforcer***REMOVED***
		store:    store,
		stopChan: make(chan struct***REMOVED******REMOVED***),
		doneChan: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// Run is the ConstraintEnforcer's main loop.
func (ce *ConstraintEnforcer) Run() ***REMOVED***
	defer close(ce.doneChan)

	watcher, cancelWatch := state.Watch(ce.store.WatchQueue(), api.EventUpdateNode***REMOVED******REMOVED***)
	defer cancelWatch()

	var (
		nodes []*api.Node
		err   error
	)
	ce.store.View(func(readTx store.ReadTx) ***REMOVED***
		nodes, err = store.FindNodes(readTx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		log.L.WithError(err).Error("failed to check nodes for noncompliant tasks")
	***REMOVED*** else ***REMOVED***
		for _, node := range nodes ***REMOVED***
			ce.rejectNoncompliantTasks(node)
		***REMOVED***
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case event := <-watcher:
			node := event.(api.EventUpdateNode).Node
			ce.rejectNoncompliantTasks(node)
		case <-ce.stopChan:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ce *ConstraintEnforcer) rejectNoncompliantTasks(node *api.Node) ***REMOVED***
	// If the availability is "drain", the orchestrator will
	// shut down all tasks.
	// If the availability is "pause", we shouldn't touch
	// the tasks on this node.
	if node.Spec.Availability != api.NodeAvailabilityActive ***REMOVED***
		return
	***REMOVED***

	var (
		tasks []*api.Task
		err   error
	)

	ce.store.View(func(tx store.ReadTx) ***REMOVED***
		tasks, err = store.FindTasks(tx, store.ByNodeID(node.ID))
	***REMOVED***)

	if err != nil ***REMOVED***
		log.L.WithError(err).Errorf("failed to list tasks for node ID %s", node.ID)
	***REMOVED***

	available := &api.Resources***REMOVED******REMOVED***
	var fakeStore []*api.GenericResource

	if node.Description != nil && node.Description.Resources != nil ***REMOVED***
		available = node.Description.Resources.Copy()
	***REMOVED***

	removeTasks := make(map[string]*api.Task)

	// TODO(aaronl): The set of tasks removed will be
	// nondeterministic because it depends on the order of
	// the slice returned from FindTasks. We could do
	// a separate pass over the tasks for each type of
	// resource, and sort by the size of the reservation
	// to remove the most resource-intensive tasks.
loop:
	for _, t := range tasks ***REMOVED***
		if t.DesiredState < api.TaskStateAssigned || t.DesiredState > api.TaskStateRunning ***REMOVED***
			continue
		***REMOVED***

		// Ensure that the task still meets scheduling
		// constraints.
		if t.Spec.Placement != nil && len(t.Spec.Placement.Constraints) != 0 ***REMOVED***
			constraints, _ := constraint.Parse(t.Spec.Placement.Constraints)
			if !constraint.NodeMatches(constraints, node) ***REMOVED***
				removeTasks[t.ID] = t
				continue
			***REMOVED***
		***REMOVED***

		// Ensure that the task assigned to the node
		// still satisfies the resource limits.
		if t.Spec.Resources != nil && t.Spec.Resources.Reservations != nil ***REMOVED***
			if t.Spec.Resources.Reservations.MemoryBytes > available.MemoryBytes ***REMOVED***
				removeTasks[t.ID] = t
				continue
			***REMOVED***
			if t.Spec.Resources.Reservations.NanoCPUs > available.NanoCPUs ***REMOVED***
				removeTasks[t.ID] = t
				continue
			***REMOVED***
			for _, ta := range t.AssignedGenericResources ***REMOVED***
				// Type change or no longer available
				if genericresource.HasResource(ta, available.Generic) ***REMOVED***
					removeTasks[t.ID] = t
					break loop
				***REMOVED***
			***REMOVED***

			available.MemoryBytes -= t.Spec.Resources.Reservations.MemoryBytes
			available.NanoCPUs -= t.Spec.Resources.Reservations.NanoCPUs

			genericresource.ClaimResources(&available.Generic,
				&fakeStore, t.AssignedGenericResources)
		***REMOVED***
	***REMOVED***

	if len(removeTasks) != 0 ***REMOVED***
		err := ce.store.Batch(func(batch *store.Batch) error ***REMOVED***
			for _, t := range removeTasks ***REMOVED***
				err := batch.Update(func(tx store.Tx) error ***REMOVED***
					t = store.GetTask(tx, t.ID)
					if t == nil || t.DesiredState > api.TaskStateRunning ***REMOVED***
						return nil
					***REMOVED***

					// We set the observed state to
					// REJECTED, rather than the desired
					// state. Desired state is owned by the
					// orchestrator, and setting it directly
					// will bypass actions such as
					// restarting the task on another node
					// (if applicable).
					t.Status.State = api.TaskStateRejected
					t.Status.Message = "task rejected by constraint enforcer"
					t.Status.Err = "assigned node no longer meets constraints"
					t.Status.Timestamp = ptypes.MustTimestampProto(time.Now())
					return store.UpdateTask(tx, t)
				***REMOVED***)
				if err != nil ***REMOVED***
					log.L.WithError(err).Errorf("failed to shut down task %s", t.ID)
				***REMOVED***
			***REMOVED***
			return nil
		***REMOVED***)

		if err != nil ***REMOVED***
			log.L.WithError(err).Errorf("failed to shut down tasks")
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stop stops the ConstraintEnforcer and waits for the main loop to exit.
func (ce *ConstraintEnforcer) Stop() ***REMOVED***
	close(ce.stopChan)
	<-ce.doneChan
***REMOVED***
