package taskinit

import (
	"sort"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/defaults"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/orchestrator"
	"github.com/docker/swarmkit/manager/orchestrator/restart"
	"github.com/docker/swarmkit/manager/state/store"
	gogotypes "github.com/gogo/protobuf/types"
	"golang.org/x/net/context"
)

// InitHandler defines orchestrator's action to fix tasks at start.
type InitHandler interface ***REMOVED***
	IsRelatedService(service *api.Service) bool
	FixTask(ctx context.Context, batch *store.Batch, t *api.Task)
	SlotTuple(t *api.Task) orchestrator.SlotTuple
***REMOVED***

// CheckTasks fixes tasks in the store before orchestrator runs. The previous leader might
// not have finished processing their updates and left them in an inconsistent state.
func CheckTasks(ctx context.Context, s *store.MemoryStore, readTx store.ReadTx, initHandler InitHandler, startSupervisor *restart.Supervisor) error ***REMOVED***
	instances := make(map[orchestrator.SlotTuple][]*api.Task)
	err := s.Batch(func(batch *store.Batch) error ***REMOVED***
		tasks, err := store.FindTasks(readTx, store.All)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, t := range tasks ***REMOVED***
			if t.ServiceID == "" ***REMOVED***
				continue
			***REMOVED***

			// TODO(aluzzardi): We should NOT retrieve the service here.
			service := store.GetService(readTx, t.ServiceID)
			if service == nil ***REMOVED***
				// Service was deleted
				err := batch.Update(func(tx store.Tx) error ***REMOVED***
					return store.DeleteTask(tx, t.ID)
				***REMOVED***)
				if err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("failed to delete task")
				***REMOVED***
				continue
			***REMOVED***
			if !initHandler.IsRelatedService(service) ***REMOVED***
				continue
			***REMOVED***

			tuple := initHandler.SlotTuple(t)
			instances[tuple] = append(instances[tuple], t)

			// handle task updates from agent which should have been triggered by task update events
			initHandler.FixTask(ctx, batch, t)

			// desired state ready is a transient state that it should be started.
			// however previous leader may not have started it, retry start here
			if t.DesiredState != api.TaskStateReady || t.Status.State > api.TaskStateRunning ***REMOVED***
				continue
			***REMOVED***
			restartDelay, _ := gogotypes.DurationFromProto(defaults.Service.Task.Restart.Delay)
			if t.Spec.Restart != nil && t.Spec.Restart.Delay != nil ***REMOVED***
				var err error
				restartDelay, err = gogotypes.DurationFromProto(t.Spec.Restart.Delay)
				if err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("invalid restart delay")
					restartDelay, _ = gogotypes.DurationFromProto(defaults.Service.Task.Restart.Delay)
				***REMOVED***
			***REMOVED***
			if restartDelay != 0 ***REMOVED***
				var timestamp time.Time
				if t.Status.AppliedAt != nil ***REMOVED***
					timestamp, err = gogotypes.TimestampFromProto(t.Status.AppliedAt)
				***REMOVED*** else ***REMOVED***
					timestamp, err = gogotypes.TimestampFromProto(t.Status.Timestamp)
				***REMOVED***
				if err == nil ***REMOVED***
					restartTime := timestamp.Add(restartDelay)
					calculatedRestartDelay := restartTime.Sub(time.Now())
					if calculatedRestartDelay < restartDelay ***REMOVED***
						restartDelay = calculatedRestartDelay
					***REMOVED***
					if restartDelay > 0 ***REMOVED***
						_ = batch.Update(func(tx store.Tx) error ***REMOVED***
							t := store.GetTask(tx, t.ID)
							// TODO(aluzzardi): This is shady as well. We should have a more generic condition.
							if t == nil || t.DesiredState != api.TaskStateReady ***REMOVED***
								return nil
							***REMOVED***
							startSupervisor.DelayStart(ctx, tx, nil, t.ID, restartDelay, true)
							return nil
						***REMOVED***)
						continue
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					log.G(ctx).WithError(err).Error("invalid status timestamp")
				***REMOVED***
			***REMOVED***

			// Start now
			err := batch.Update(func(tx store.Tx) error ***REMOVED***
				return startSupervisor.StartNow(tx, t.ID)
			***REMOVED***)
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).WithField("task.id", t.ID).Error("moving task out of delayed state failed")
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for tuple, instance := range instances ***REMOVED***
		// Find the most current spec version. That's the only one
		// we care about for the purpose of reconstructing restart
		// history.
		maxVersion := uint64(0)
		for _, t := range instance ***REMOVED***
			if t.SpecVersion != nil && t.SpecVersion.Index > maxVersion ***REMOVED***
				maxVersion = t.SpecVersion.Index
			***REMOVED***
		***REMOVED***

		// Create a new slice with just the current spec version tasks.
		var upToDate []*api.Task
		for _, t := range instance ***REMOVED***
			if t.SpecVersion != nil && t.SpecVersion.Index == maxVersion ***REMOVED***
				upToDate = append(upToDate, t)
			***REMOVED***
		***REMOVED***

		// Sort by creation timestamp
		sort.Sort(tasksByCreationTimestamp(upToDate))

		// All up-to-date tasks in this instance except the first one
		// should be considered restarted.
		if len(upToDate) < 2 ***REMOVED***
			continue
		***REMOVED***
		for _, t := range upToDate[1:] ***REMOVED***
			startSupervisor.RecordRestartHistory(tuple, t)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type tasksByCreationTimestamp []*api.Task

func (t tasksByCreationTimestamp) Len() int ***REMOVED***
	return len(t)
***REMOVED***
func (t tasksByCreationTimestamp) Swap(i, j int) ***REMOVED***
	t[i], t[j] = t[j], t[i]
***REMOVED***
func (t tasksByCreationTimestamp) Less(i, j int) bool ***REMOVED***
	if t[i].Meta.CreatedAt == nil ***REMOVED***
		return true
	***REMOVED***
	if t[j].Meta.CreatedAt == nil ***REMOVED***
		return false
	***REMOVED***
	if t[i].Meta.CreatedAt.Seconds < t[j].Meta.CreatedAt.Seconds ***REMOVED***
		return true
	***REMOVED***
	if t[i].Meta.CreatedAt.Seconds > t[j].Meta.CreatedAt.Seconds ***REMOVED***
		return false
	***REMOVED***
	return t[i].Meta.CreatedAt.Nanos < t[j].Meta.CreatedAt.Nanos
***REMOVED***
