package replicated

import (
	"github.com/docker/go-events"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/orchestrator"
	"github.com/docker/swarmkit/manager/orchestrator/taskinit"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
)

// This file provides task-level orchestration. It observes changes to task
// and node state and kills/recreates tasks if necessary. This is distinct from
// service-level reconciliation, which observes changes to services and creates
// and/or kills tasks to match the service definition.

func (r *Orchestrator) initTasks(ctx context.Context, readTx store.ReadTx) error ***REMOVED***
	return taskinit.CheckTasks(ctx, r.store, readTx, r, r.restarts)
***REMOVED***

func (r *Orchestrator) handleTaskEvent(ctx context.Context, event events.Event) ***REMOVED***
	switch v := event.(type) ***REMOVED***
	case api.EventDeleteNode:
		r.restartTasksByNodeID(ctx, v.Node.ID)
	case api.EventCreateNode:
		r.handleNodeChange(ctx, v.Node)
	case api.EventUpdateNode:
		r.handleNodeChange(ctx, v.Node)
	case api.EventDeleteTask:
		if v.Task.DesiredState <= api.TaskStateRunning ***REMOVED***
			service := r.resolveService(ctx, v.Task)
			if !orchestrator.IsReplicatedService(service) ***REMOVED***
				return
			***REMOVED***
			r.reconcileServices[service.ID] = service
		***REMOVED***
		r.restarts.Cancel(v.Task.ID)
	case api.EventUpdateTask:
		r.handleTaskChange(ctx, v.Task)
	case api.EventCreateTask:
		r.handleTaskChange(ctx, v.Task)
	***REMOVED***
***REMOVED***

func (r *Orchestrator) tickTasks(ctx context.Context) ***REMOVED***
	if len(r.restartTasks) > 0 ***REMOVED***
		err := r.store.Batch(func(batch *store.Batch) error ***REMOVED***
			for taskID := range r.restartTasks ***REMOVED***
				err := batch.Update(func(tx store.Tx) error ***REMOVED***
					// TODO(aaronl): optimistic update?
					t := store.GetTask(tx, taskID)
					if t != nil ***REMOVED***
						if t.DesiredState > api.TaskStateRunning ***REMOVED***
							return nil
						***REMOVED***

						service := store.GetService(tx, t.ServiceID)
						if !orchestrator.IsReplicatedService(service) ***REMOVED***
							return nil
						***REMOVED***

						// Restart task if applicable
						if err := r.restarts.Restart(ctx, tx, r.cluster, service, *t); err != nil ***REMOVED***
							return err
						***REMOVED***
					***REMOVED***
					return nil
				***REMOVED***)
				if err != nil ***REMOVED***
					log.G(ctx).WithError(err).Errorf("Orchestrator task reaping transaction failed")
				***REMOVED***
			***REMOVED***
			return nil
		***REMOVED***)

		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("orchestrator task removal batch failed")
		***REMOVED***

		r.restartTasks = make(map[string]struct***REMOVED******REMOVED***)
	***REMOVED***
***REMOVED***

func (r *Orchestrator) restartTasksByNodeID(ctx context.Context, nodeID string) ***REMOVED***
	var err error
	r.store.View(func(tx store.ReadTx) ***REMOVED***
		var tasks []*api.Task
		tasks, err = store.FindTasks(tx, store.ByNodeID(nodeID))
		if err != nil ***REMOVED***
			return
		***REMOVED***

		for _, t := range tasks ***REMOVED***
			if t.DesiredState > api.TaskStateRunning ***REMOVED***
				continue
			***REMOVED***
			service := store.GetService(tx, t.ServiceID)
			if orchestrator.IsReplicatedService(service) ***REMOVED***
				r.restartTasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("failed to list tasks to remove")
	***REMOVED***
***REMOVED***

func (r *Orchestrator) handleNodeChange(ctx context.Context, n *api.Node) ***REMOVED***
	if !orchestrator.InvalidNode(n) ***REMOVED***
		return
	***REMOVED***

	r.restartTasksByNodeID(ctx, n.ID)
***REMOVED***

// handleTaskChange defines what orchestrator does when a task is updated by agent.
func (r *Orchestrator) handleTaskChange(ctx context.Context, t *api.Task) ***REMOVED***
	// If we already set the desired state past TaskStateRunning, there is no
	// further action necessary.
	if t.DesiredState > api.TaskStateRunning ***REMOVED***
		return
	***REMOVED***

	var (
		n       *api.Node
		service *api.Service
	)
	r.store.View(func(tx store.ReadTx) ***REMOVED***
		if t.NodeID != "" ***REMOVED***
			n = store.GetNode(tx, t.NodeID)
		***REMOVED***
		if t.ServiceID != "" ***REMOVED***
			service = store.GetService(tx, t.ServiceID)
		***REMOVED***
	***REMOVED***)

	if !orchestrator.IsReplicatedService(service) ***REMOVED***
		return
	***REMOVED***

	if t.Status.State > api.TaskStateRunning ||
		(t.NodeID != "" && orchestrator.InvalidNode(n)) ***REMOVED***
		r.restartTasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// FixTask validates a task with the current cluster settings, and takes
// action to make it conformant. it's called at orchestrator initialization.
func (r *Orchestrator) FixTask(ctx context.Context, batch *store.Batch, t *api.Task) ***REMOVED***
	// If we already set the desired state past TaskStateRunning, there is no
	// further action necessary.
	if t.DesiredState > api.TaskStateRunning ***REMOVED***
		return
	***REMOVED***

	var (
		n       *api.Node
		service *api.Service
	)
	batch.Update(func(tx store.Tx) error ***REMOVED***
		if t.NodeID != "" ***REMOVED***
			n = store.GetNode(tx, t.NodeID)
		***REMOVED***
		if t.ServiceID != "" ***REMOVED***
			service = store.GetService(tx, t.ServiceID)
		***REMOVED***
		return nil
	***REMOVED***)

	if !orchestrator.IsReplicatedService(service) ***REMOVED***
		return
	***REMOVED***

	if t.Status.State > api.TaskStateRunning ||
		(t.NodeID != "" && orchestrator.InvalidNode(n)) ***REMOVED***
		r.restartTasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		return
	***REMOVED***
***REMOVED***
