package restart

import (
	"container/list"
	"errors"
	"sync"
	"time"

	"github.com/docker/go-events"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/defaults"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/orchestrator"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	gogotypes "github.com/gogo/protobuf/types"
	"golang.org/x/net/context"
)

const defaultOldTaskTimeout = time.Minute

type restartedInstance struct ***REMOVED***
	timestamp time.Time
***REMOVED***

type instanceRestartInfo struct ***REMOVED***
	// counter of restarts for this instance.
	totalRestarts uint64
	// Linked list of restartedInstance structs. Only used when
	// Restart.MaxAttempts and Restart.Window are both
	// nonzero.
	restartedInstances *list.List
	// Why is specVersion in this structure and not in the map key? While
	// putting it in the key would be a very simple solution, it wouldn't
	// be easy to clean up map entries corresponding to old specVersions.
	// Making the key version-agnostic and clearing the value whenever the
	// version changes avoids the issue of stale map entries for old
	// versions.
	specVersion api.Version
***REMOVED***

type delayedStart struct ***REMOVED***
	// cancel is called to cancel the delayed start.
	cancel func()
	doneCh chan struct***REMOVED******REMOVED***

	// waiter is set to true if the next restart is waiting for this delay
	// to complete.
	waiter bool
***REMOVED***

// Supervisor initiates and manages restarts. It's responsible for
// delaying restarts when applicable.
type Supervisor struct ***REMOVED***
	mu               sync.Mutex
	store            *store.MemoryStore
	delays           map[string]*delayedStart
	historyByService map[string]map[orchestrator.SlotTuple]*instanceRestartInfo
	TaskTimeout      time.Duration
***REMOVED***

// NewSupervisor creates a new RestartSupervisor.
func NewSupervisor(store *store.MemoryStore) *Supervisor ***REMOVED***
	return &Supervisor***REMOVED***
		store:            store,
		delays:           make(map[string]*delayedStart),
		historyByService: make(map[string]map[orchestrator.SlotTuple]*instanceRestartInfo),
		TaskTimeout:      defaultOldTaskTimeout,
	***REMOVED***
***REMOVED***

func (r *Supervisor) waitRestart(ctx context.Context, oldDelay *delayedStart, cluster *api.Cluster, taskID string) ***REMOVED***
	// Wait for the last restart delay to elapse.
	select ***REMOVED***
	case <-oldDelay.doneCh:
	case <-ctx.Done():
		return
	***REMOVED***

	// Start the next restart
	err := r.store.Update(func(tx store.Tx) error ***REMOVED***
		t := store.GetTask(tx, taskID)
		if t == nil ***REMOVED***
			return nil
		***REMOVED***
		if t.DesiredState > api.TaskStateRunning ***REMOVED***
			return nil
		***REMOVED***
		service := store.GetService(tx, t.ServiceID)
		if service == nil ***REMOVED***
			return nil
		***REMOVED***
		return r.Restart(ctx, tx, cluster, service, *t)
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("failed to restart task after waiting for previous restart")
	***REMOVED***
***REMOVED***

// Restart initiates a new task to replace t if appropriate under the service's
// restart policy.
func (r *Supervisor) Restart(ctx context.Context, tx store.Tx, cluster *api.Cluster, service *api.Service, t api.Task) error ***REMOVED***
	// TODO(aluzzardi): This function should not depend on `service`.

	// Is the old task still in the process of restarting? If so, wait for
	// its restart delay to elapse, to avoid tight restart loops (for
	// example, when the image doesn't exist).
	r.mu.Lock()
	oldDelay, ok := r.delays[t.ID]
	if ok ***REMOVED***
		if !oldDelay.waiter ***REMOVED***
			oldDelay.waiter = true
			go r.waitRestart(ctx, oldDelay, cluster, t.ID)
		***REMOVED***
		r.mu.Unlock()
		return nil
	***REMOVED***
	r.mu.Unlock()

	// Sanity check: was the task shut down already by a separate call to
	// Restart? If so, we must avoid restarting it, because this will create
	// an extra task. This should never happen unless there is a bug.
	if t.DesiredState > api.TaskStateRunning ***REMOVED***
		return errors.New("Restart called on task that was already shut down")
	***REMOVED***

	t.DesiredState = api.TaskStateShutdown
	err := store.UpdateTask(tx, &t)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("failed to set task desired state to dead")
		return err
	***REMOVED***

	if !r.shouldRestart(ctx, &t, service) ***REMOVED***
		return nil
	***REMOVED***

	var restartTask *api.Task

	if orchestrator.IsReplicatedService(service) ***REMOVED***
		restartTask = orchestrator.NewTask(cluster, service, t.Slot, "")
	***REMOVED*** else if orchestrator.IsGlobalService(service) ***REMOVED***
		restartTask = orchestrator.NewTask(cluster, service, 0, t.NodeID)
	***REMOVED*** else ***REMOVED***
		log.G(ctx).Error("service not supported by restart supervisor")
		return nil
	***REMOVED***

	n := store.GetNode(tx, t.NodeID)

	restartTask.DesiredState = api.TaskStateReady

	var restartDelay time.Duration
	// Restart delay is not applied to drained nodes
	if n == nil || n.Spec.Availability != api.NodeAvailabilityDrain ***REMOVED***
		if t.Spec.Restart != nil && t.Spec.Restart.Delay != nil ***REMOVED***
			var err error
			restartDelay, err = gogotypes.DurationFromProto(t.Spec.Restart.Delay)
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("invalid restart delay; using default")
				restartDelay, _ = gogotypes.DurationFromProto(defaults.Service.Task.Restart.Delay)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			restartDelay, _ = gogotypes.DurationFromProto(defaults.Service.Task.Restart.Delay)
		***REMOVED***
	***REMOVED***

	waitStop := true

	// Normally we wait for the old task to stop running, but we skip this
	// if the old task is already dead or the node it's assigned to is down.
	if (n != nil && n.Status.State == api.NodeStatus_DOWN) || t.Status.State > api.TaskStateRunning ***REMOVED***
		waitStop = false
	***REMOVED***

	if err := store.CreateTask(tx, restartTask); err != nil ***REMOVED***
		log.G(ctx).WithError(err).WithField("task.id", restartTask.ID).Error("task create failed")
		return err
	***REMOVED***

	tuple := orchestrator.SlotTuple***REMOVED***
		Slot:      restartTask.Slot,
		ServiceID: restartTask.ServiceID,
		NodeID:    restartTask.NodeID,
	***REMOVED***
	r.RecordRestartHistory(tuple, restartTask)

	r.DelayStart(ctx, tx, &t, restartTask.ID, restartDelay, waitStop)
	return nil
***REMOVED***

// shouldRestart returns true if a task should be restarted according to the
// restart policy.
func (r *Supervisor) shouldRestart(ctx context.Context, t *api.Task, service *api.Service) bool ***REMOVED***
	// TODO(aluzzardi): This function should not depend on `service`.
	condition := orchestrator.RestartCondition(t)

	if condition != api.RestartOnAny &&
		(condition != api.RestartOnFailure || t.Status.State == api.TaskStateCompleted) ***REMOVED***
		return false
	***REMOVED***

	if t.Spec.Restart == nil || t.Spec.Restart.MaxAttempts == 0 ***REMOVED***
		return true
	***REMOVED***

	instanceTuple := orchestrator.SlotTuple***REMOVED***
		Slot:      t.Slot,
		ServiceID: t.ServiceID,
	***REMOVED***

	// Slot is not meaningful for "global" tasks, so they need to be
	// indexed by NodeID.
	if orchestrator.IsGlobalService(service) ***REMOVED***
		instanceTuple.NodeID = t.NodeID
	***REMOVED***

	r.mu.Lock()
	defer r.mu.Unlock()

	restartInfo := r.historyByService[t.ServiceID][instanceTuple]
	if restartInfo == nil || (t.SpecVersion != nil && *t.SpecVersion != restartInfo.specVersion) ***REMOVED***
		return true
	***REMOVED***

	if t.Spec.Restart.Window == nil || (t.Spec.Restart.Window.Seconds == 0 && t.Spec.Restart.Window.Nanos == 0) ***REMOVED***
		return restartInfo.totalRestarts < t.Spec.Restart.MaxAttempts
	***REMOVED***

	if restartInfo.restartedInstances == nil ***REMOVED***
		return true
	***REMOVED***

	window, err := gogotypes.DurationFromProto(t.Spec.Restart.Window)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("invalid restart lookback window")
		return restartInfo.totalRestarts < t.Spec.Restart.MaxAttempts
	***REMOVED***

	var timestamp time.Time
	// Prefer the manager's timestamp over the agent's, since manager
	// clocks are more trustworthy.
	if t.Status.AppliedAt != nil ***REMOVED***
		timestamp, err = gogotypes.TimestampFromProto(t.Status.AppliedAt)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("invalid task status AppliedAt timestamp")
			return restartInfo.totalRestarts < t.Spec.Restart.MaxAttempts
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// It's safe to call TimestampFromProto with a nil timestamp
		timestamp, err = gogotypes.TimestampFromProto(t.Status.Timestamp)
		if t.Status.Timestamp == nil || err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("invalid task completion timestamp")
			return restartInfo.totalRestarts < t.Spec.Restart.MaxAttempts
		***REMOVED***
	***REMOVED***
	lookback := timestamp.Add(-window)

	numRestarts := uint64(restartInfo.restartedInstances.Len())

	// Disregard any restarts that happened before the lookback window,
	// and remove them from the linked list since they will no longer
	// be relevant to figuring out if tasks should be restarted going
	// forward.
	var next *list.Element
	for e := restartInfo.restartedInstances.Front(); e != nil; e = next ***REMOVED***
		next = e.Next()

		if e.Value.(restartedInstance).timestamp.After(lookback) ***REMOVED***
			break
		***REMOVED***
		restartInfo.restartedInstances.Remove(e)
		numRestarts--
	***REMOVED***

	// Ignore restarts that didn't happen before the task we're looking at.
	for e2 := restartInfo.restartedInstances.Back(); e2 != nil; e2 = e2.Prev() ***REMOVED***
		if e2.Value.(restartedInstance).timestamp.Before(timestamp) ***REMOVED***
			break
		***REMOVED***
		numRestarts--
	***REMOVED***

	if restartInfo.restartedInstances.Len() == 0 ***REMOVED***
		restartInfo.restartedInstances = nil
	***REMOVED***

	return numRestarts < t.Spec.Restart.MaxAttempts
***REMOVED***

// UpdatableTasksInSlot returns the set of tasks that should be passed to the
// updater from this slot, or an empty slice if none should be.  An updatable
// slot has either at least one task that with desired state <= RUNNING, or its
// most recent task has stopped running and should not be restarted. The latter
// case is for making sure that tasks that shouldn't normally be restarted will
// still be handled by rolling updates when they become outdated.  There is a
// special case for rollbacks to make sure that a rollback always takes the
// service to a converged state, instead of ignoring tasks with the original
// spec that stopped running and shouldn't be restarted according to the
// restart policy.
func (r *Supervisor) UpdatableTasksInSlot(ctx context.Context, slot orchestrator.Slot, service *api.Service) orchestrator.Slot ***REMOVED***
	if len(slot) < 1 ***REMOVED***
		return nil
	***REMOVED***

	var updatable orchestrator.Slot
	for _, t := range slot ***REMOVED***
		if t.DesiredState <= api.TaskStateRunning ***REMOVED***
			updatable = append(updatable, t)
		***REMOVED***
	***REMOVED***
	if len(updatable) > 0 ***REMOVED***
		return updatable
	***REMOVED***

	if service.UpdateStatus != nil && service.UpdateStatus.State == api.UpdateStatus_ROLLBACK_STARTED ***REMOVED***
		return nil
	***REMOVED***

	// Find most recent task
	byTimestamp := orchestrator.TasksByTimestamp(slot)
	newestIndex := 0
	for i := 1; i != len(slot); i++ ***REMOVED***
		if byTimestamp.Less(newestIndex, i) ***REMOVED***
			newestIndex = i
		***REMOVED***
	***REMOVED***

	if !r.shouldRestart(ctx, slot[newestIndex], service) ***REMOVED***
		return orchestrator.Slot***REMOVED***slot[newestIndex]***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// RecordRestartHistory updates the historyByService map to reflect the restart
// of restartedTask.
func (r *Supervisor) RecordRestartHistory(tuple orchestrator.SlotTuple, replacementTask *api.Task) ***REMOVED***
	if replacementTask.Spec.Restart == nil || replacementTask.Spec.Restart.MaxAttempts == 0 ***REMOVED***
		// No limit on the number of restarts, so no need to record
		// history.
		return
	***REMOVED***

	r.mu.Lock()
	defer r.mu.Unlock()

	serviceID := replacementTask.ServiceID
	if r.historyByService[serviceID] == nil ***REMOVED***
		r.historyByService[serviceID] = make(map[orchestrator.SlotTuple]*instanceRestartInfo)
	***REMOVED***
	if r.historyByService[serviceID][tuple] == nil ***REMOVED***
		r.historyByService[serviceID][tuple] = &instanceRestartInfo***REMOVED******REMOVED***
	***REMOVED***

	restartInfo := r.historyByService[serviceID][tuple]

	if replacementTask.SpecVersion != nil && *replacementTask.SpecVersion != restartInfo.specVersion ***REMOVED***
		// This task has a different SpecVersion from the one we're
		// tracking. Most likely, the service was updated. Past failures
		// shouldn't count against the new service definition, so clear
		// the history for this instance.
		*restartInfo = instanceRestartInfo***REMOVED***
			specVersion: *replacementTask.SpecVersion,
		***REMOVED***
	***REMOVED***

	restartInfo.totalRestarts++

	if replacementTask.Spec.Restart.Window != nil && (replacementTask.Spec.Restart.Window.Seconds != 0 || replacementTask.Spec.Restart.Window.Nanos != 0) ***REMOVED***
		if restartInfo.restartedInstances == nil ***REMOVED***
			restartInfo.restartedInstances = list.New()
		***REMOVED***

		// it's okay to call TimestampFromProto with a nil argument
		timestamp, err := gogotypes.TimestampFromProto(replacementTask.Meta.CreatedAt)
		if replacementTask.Meta.CreatedAt == nil || err != nil ***REMOVED***
			timestamp = time.Now()
		***REMOVED***

		restartedInstance := restartedInstance***REMOVED***
			timestamp: timestamp,
		***REMOVED***

		restartInfo.restartedInstances.PushBack(restartedInstance)
	***REMOVED***
***REMOVED***

// DelayStart starts a timer that moves the task from READY to RUNNING once:
// - The restart delay has elapsed (if applicable)
// - The old task that it's replacing has stopped running (or this times out)
// It must be called during an Update transaction to ensure that it does not
// miss events. The purpose of the store.Tx argument is to avoid accidental
// calls outside an Update transaction.
func (r *Supervisor) DelayStart(ctx context.Context, _ store.Tx, oldTask *api.Task, newTaskID string, delay time.Duration, waitStop bool) <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	doneCh := make(chan struct***REMOVED******REMOVED***)

	r.mu.Lock()
	for ***REMOVED***
		oldDelay, ok := r.delays[newTaskID]
		if !ok ***REMOVED***
			break
		***REMOVED***
		oldDelay.cancel()
		r.mu.Unlock()
		// Note that this channel read should only block for a very
		// short time, because we cancelled the existing delay and
		// that should cause it to stop immediately.
		<-oldDelay.doneCh
		r.mu.Lock()
	***REMOVED***
	r.delays[newTaskID] = &delayedStart***REMOVED***cancel: cancel, doneCh: doneCh***REMOVED***
	r.mu.Unlock()

	var watch chan events.Event
	cancelWatch := func() ***REMOVED******REMOVED***

	waitForTask := waitStop && oldTask != nil && oldTask.Status.State <= api.TaskStateRunning

	if waitForTask ***REMOVED***
		// Wait for either the old task to complete, or the old task's
		// node to become unavailable.
		watch, cancelWatch = state.Watch(
			r.store.WatchQueue(),
			api.EventUpdateTask***REMOVED***
				Task:   &api.Task***REMOVED***ID: oldTask.ID, Status: api.TaskStatus***REMOVED***State: api.TaskStateRunning***REMOVED******REMOVED***,
				Checks: []api.TaskCheckFunc***REMOVED***api.TaskCheckID, state.TaskCheckStateGreaterThan***REMOVED***,
			***REMOVED***,
			api.EventUpdateNode***REMOVED***
				Node:   &api.Node***REMOVED***ID: oldTask.NodeID, Status: api.NodeStatus***REMOVED***State: api.NodeStatus_DOWN***REMOVED******REMOVED***,
				Checks: []api.NodeCheckFunc***REMOVED***api.NodeCheckID, state.NodeCheckState***REMOVED***,
			***REMOVED***,
			api.EventDeleteNode***REMOVED***
				Node:   &api.Node***REMOVED***ID: oldTask.NodeID***REMOVED***,
				Checks: []api.NodeCheckFunc***REMOVED***api.NodeCheckID***REMOVED***,
			***REMOVED***,
		)
	***REMOVED***

	go func() ***REMOVED***
		defer func() ***REMOVED***
			cancelWatch()
			r.mu.Lock()
			delete(r.delays, newTaskID)
			r.mu.Unlock()
			close(doneCh)
		***REMOVED***()

		oldTaskTimer := time.NewTimer(r.TaskTimeout)
		defer oldTaskTimer.Stop()

		// Wait for the delay to elapse, if one is specified.
		if delay != 0 ***REMOVED***
			select ***REMOVED***
			case <-time.After(delay):
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***

		if waitForTask ***REMOVED***
			select ***REMOVED***
			case <-watch:
			case <-oldTaskTimer.C:
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***

		err := r.store.Update(func(tx store.Tx) error ***REMOVED***
			err := r.StartNow(tx, newTaskID)
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).WithField("task.id", newTaskID).Error("moving task out of delayed state failed")
			***REMOVED***
			return nil
		***REMOVED***)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).WithField("task.id", newTaskID).Error("task restart transaction failed")
		***REMOVED***
	***REMOVED***()

	return doneCh
***REMOVED***

// StartNow moves the task into the RUNNING state so it will proceed to start
// up.
func (r *Supervisor) StartNow(tx store.Tx, taskID string) error ***REMOVED***
	t := store.GetTask(tx, taskID)
	if t == nil || t.DesiredState >= api.TaskStateRunning ***REMOVED***
		return nil
	***REMOVED***
	t.DesiredState = api.TaskStateRunning
	return store.UpdateTask(tx, t)
***REMOVED***

// Cancel cancels a pending restart.
func (r *Supervisor) Cancel(taskID string) ***REMOVED***
	r.mu.Lock()
	delay, ok := r.delays[taskID]
	r.mu.Unlock()

	if !ok ***REMOVED***
		return
	***REMOVED***

	delay.cancel()
	<-delay.doneCh
***REMOVED***

// CancelAll aborts all pending restarts and waits for any instances of
// StartNow that have already triggered to complete.
func (r *Supervisor) CancelAll() ***REMOVED***
	var cancelled []delayedStart

	r.mu.Lock()
	for _, delay := range r.delays ***REMOVED***
		delay.cancel()
	***REMOVED***
	r.mu.Unlock()

	for _, delay := range cancelled ***REMOVED***
		<-delay.doneCh
	***REMOVED***
***REMOVED***

// ClearServiceHistory forgets restart history related to a given service ID.
func (r *Supervisor) ClearServiceHistory(serviceID string) ***REMOVED***
	r.mu.Lock()
	delete(r.historyByService, serviceID)
	r.mu.Unlock()
***REMOVED***
