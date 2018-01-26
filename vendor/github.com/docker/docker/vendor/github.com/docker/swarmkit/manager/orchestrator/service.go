package orchestrator

import (
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
)

// IsReplicatedService checks if a service is a replicated service.
func IsReplicatedService(service *api.Service) bool ***REMOVED***
	// service nil validation is required as there are scenarios
	// where service is removed from store
	if service == nil ***REMOVED***
		return false
	***REMOVED***
	_, ok := service.Spec.GetMode().(*api.ServiceSpec_Replicated)
	return ok
***REMOVED***

// IsGlobalService checks if the service is a global service.
func IsGlobalService(service *api.Service) bool ***REMOVED***
	if service == nil ***REMOVED***
		return false
	***REMOVED***
	_, ok := service.Spec.GetMode().(*api.ServiceSpec_Global)
	return ok
***REMOVED***

// SetServiceTasksRemove sets the desired state of tasks associated with a service
// to REMOVE, so that they can be properly shut down by the agent and later removed
// by the task reaper.
func SetServiceTasksRemove(ctx context.Context, s *store.MemoryStore, service *api.Service) ***REMOVED***
	var (
		tasks []*api.Task
		err   error
	)
	s.View(func(tx store.ReadTx) ***REMOVED***
		tasks, err = store.FindTasks(tx, store.ByServiceID(service.ID))
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("failed to list tasks")
		return
	***REMOVED***

	err = s.Batch(func(batch *store.Batch) error ***REMOVED***
		for _, t := range tasks ***REMOVED***
			err := batch.Update(func(tx store.Tx) error ***REMOVED***
				// time travel is not allowed. if the current desired state is
				// above the one we're trying to go to we can't go backwards.
				// we have nothing to do and we should skip to the next task
				if t.DesiredState > api.TaskStateRemove ***REMOVED***
					// log a warning, though. we shouln't be trying to rewrite
					// a state to an earlier state
					log.G(ctx).Warnf(
						"cannot update task %v in desired state %v to an earlier desired state %v",
						t.ID, t.DesiredState, api.TaskStateRemove,
					)
					return nil
				***REMOVED***
				// update desired state to REMOVE
				t.DesiredState = api.TaskStateRemove

				if err := store.UpdateTask(tx, t); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Errorf("failed transaction: update task desired state to REMOVE")
				***REMOVED***
				return nil
			***REMOVED***)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("task search transaction failed")
	***REMOVED***
***REMOVED***
