package replicated

import (
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/orchestrator"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
)

type slotsByRunningState []orchestrator.Slot

func (is slotsByRunningState) Len() int      ***REMOVED*** return len(is) ***REMOVED***
func (is slotsByRunningState) Swap(i, j int) ***REMOVED*** is[i], is[j] = is[j], is[i] ***REMOVED***

// Less returns true if the first task should be preferred over the second task,
// all other things being equal in terms of node balance.
func (is slotsByRunningState) Less(i, j int) bool ***REMOVED***
	iRunning := false
	jRunning := false

	for _, ii := range is[i] ***REMOVED***
		if ii.Status.State == api.TaskStateRunning ***REMOVED***
			iRunning = true
			break
		***REMOVED***
	***REMOVED***
	for _, ij := range is[j] ***REMOVED***
		if ij.Status.State == api.TaskStateRunning ***REMOVED***
			jRunning = true
			break
		***REMOVED***
	***REMOVED***

	if iRunning && !jRunning ***REMOVED***
		return true
	***REMOVED***

	if !iRunning && jRunning ***REMOVED***
		return false
	***REMOVED***

	// Use Slot number as a tie-breaker to prefer to remove tasks in reverse
	// order of Slot number. This would help us avoid unnecessary master
	// migration when scaling down a stateful service because the master
	// task of a stateful service is usually in a low numbered Slot.
	return is[i][0].Slot < is[j][0].Slot
***REMOVED***

type slotWithIndex struct ***REMOVED***
	slot orchestrator.Slot

	// index is a counter that counts this task as the nth instance of
	// the service on its node. This is used for sorting the tasks so that
	// when scaling down we leave tasks more evenly balanced.
	index int
***REMOVED***

type slotsByIndex []slotWithIndex

func (is slotsByIndex) Len() int      ***REMOVED*** return len(is) ***REMOVED***
func (is slotsByIndex) Swap(i, j int) ***REMOVED*** is[i], is[j] = is[j], is[i] ***REMOVED***

func (is slotsByIndex) Less(i, j int) bool ***REMOVED***
	if is[i].index < 0 && is[j].index >= 0 ***REMOVED***
		return false
	***REMOVED***
	if is[j].index < 0 && is[i].index >= 0 ***REMOVED***
		return true
	***REMOVED***
	return is[i].index < is[j].index
***REMOVED***

// updatableAndDeadSlots returns two maps of slots. The first contains slots
// that have at least one task with a desired state above NEW and lesser or
// equal to RUNNING, or a task that shouldn't be restarted. The second contains
// all other slots with at least one task.
func (r *Orchestrator) updatableAndDeadSlots(ctx context.Context, service *api.Service) (map[uint64]orchestrator.Slot, map[uint64]orchestrator.Slot, error) ***REMOVED***
	var (
		tasks []*api.Task
		err   error
	)
	r.store.View(func(tx store.ReadTx) ***REMOVED***
		tasks, err = store.FindTasks(tx, store.ByServiceID(service.ID))
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	updatableSlots := make(map[uint64]orchestrator.Slot)
	for _, t := range tasks ***REMOVED***
		updatableSlots[t.Slot] = append(updatableSlots[t.Slot], t)
	***REMOVED***

	deadSlots := make(map[uint64]orchestrator.Slot)
	for slotID, slot := range updatableSlots ***REMOVED***
		updatable := r.restarts.UpdatableTasksInSlot(ctx, slot, service)
		if len(updatable) != 0 ***REMOVED***
			updatableSlots[slotID] = updatable
		***REMOVED*** else ***REMOVED***
			delete(updatableSlots, slotID)
			deadSlots[slotID] = slot
		***REMOVED***
	***REMOVED***

	return updatableSlots, deadSlots, nil
***REMOVED***

// SlotTuple returns a slot tuple for the replicated service task.
func (r *Orchestrator) SlotTuple(t *api.Task) orchestrator.SlotTuple ***REMOVED***
	return orchestrator.SlotTuple***REMOVED***
		ServiceID: t.ServiceID,
		Slot:      t.Slot,
	***REMOVED***
***REMOVED***
