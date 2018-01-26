package update

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/go-events"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/defaults"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/orchestrator"
	"github.com/docker/swarmkit/manager/orchestrator/restart"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/protobuf/ptypes"
	"github.com/docker/swarmkit/watch"
	gogotypes "github.com/gogo/protobuf/types"
)

// Supervisor supervises a set of updates. It's responsible for keeping track of updates,
// shutting them down and replacing them.
type Supervisor struct ***REMOVED***
	store    *store.MemoryStore
	restarts *restart.Supervisor
	updates  map[string]*Updater
	l        sync.Mutex
***REMOVED***

// NewSupervisor creates a new UpdateSupervisor.
func NewSupervisor(store *store.MemoryStore, restartSupervisor *restart.Supervisor) *Supervisor ***REMOVED***
	return &Supervisor***REMOVED***
		store:    store,
		updates:  make(map[string]*Updater),
		restarts: restartSupervisor,
	***REMOVED***
***REMOVED***

// Update starts an Update of `slots` belonging to `service` in the background
// and returns immediately. Each slot contains a group of one or more tasks
// occupying the same slot (replicated service) or node (global service). There
// may be more than one task per slot in cases where an update is in progress
// and the new task was started before the old one was shut down. If an update
// for that service was already in progress, it will be cancelled before the
// new one starts.
func (u *Supervisor) Update(ctx context.Context, cluster *api.Cluster, service *api.Service, slots []orchestrator.Slot) ***REMOVED***
	u.l.Lock()
	defer u.l.Unlock()

	id := service.ID

	if update, ok := u.updates[id]; ok ***REMOVED***
		if reflect.DeepEqual(service.Spec, update.newService.Spec) ***REMOVED***
			// There's already an update working towards this goal.
			return
		***REMOVED***
		update.Cancel()
	***REMOVED***

	update := NewUpdater(u.store, u.restarts, cluster, service)
	u.updates[id] = update
	go func() ***REMOVED***
		update.Run(ctx, slots)
		u.l.Lock()
		if u.updates[id] == update ***REMOVED***
			delete(u.updates, id)
		***REMOVED***
		u.l.Unlock()
	***REMOVED***()
***REMOVED***

// CancelAll cancels all current updates.
func (u *Supervisor) CancelAll() ***REMOVED***
	u.l.Lock()
	defer u.l.Unlock()

	for _, update := range u.updates ***REMOVED***
		update.Cancel()
	***REMOVED***
***REMOVED***

// Updater updates a set of tasks to a new version.
type Updater struct ***REMOVED***
	store      *store.MemoryStore
	watchQueue *watch.Queue
	restarts   *restart.Supervisor

	cluster    *api.Cluster
	newService *api.Service

	updatedTasks   map[string]time.Time // task ID to creation time
	updatedTasksMu sync.Mutex

	// stopChan signals to the state machine to stop running.
	stopChan chan struct***REMOVED******REMOVED***
	// doneChan is closed when the state machine terminates.
	doneChan chan struct***REMOVED******REMOVED***
***REMOVED***

// NewUpdater creates a new Updater.
func NewUpdater(store *store.MemoryStore, restartSupervisor *restart.Supervisor, cluster *api.Cluster, newService *api.Service) *Updater ***REMOVED***
	return &Updater***REMOVED***
		store:        store,
		watchQueue:   store.WatchQueue(),
		restarts:     restartSupervisor,
		cluster:      cluster.Copy(),
		newService:   newService.Copy(),
		updatedTasks: make(map[string]time.Time),
		stopChan:     make(chan struct***REMOVED******REMOVED***),
		doneChan:     make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// Cancel cancels the current update immediately. It blocks until the cancellation is confirmed.
func (u *Updater) Cancel() ***REMOVED***
	close(u.stopChan)
	<-u.doneChan
***REMOVED***

// Run starts the update and returns only once its complete or cancelled.
func (u *Updater) Run(ctx context.Context, slots []orchestrator.Slot) ***REMOVED***
	defer close(u.doneChan)

	service := u.newService

	// If the update is in a PAUSED state, we should not do anything.
	if service.UpdateStatus != nil &&
		(service.UpdateStatus.State == api.UpdateStatus_PAUSED ||
			service.UpdateStatus.State == api.UpdateStatus_ROLLBACK_PAUSED) ***REMOVED***
		return
	***REMOVED***

	var dirtySlots []orchestrator.Slot
	for _, slot := range slots ***REMOVED***
		if u.isSlotDirty(slot) ***REMOVED***
			dirtySlots = append(dirtySlots, slot)
		***REMOVED***
	***REMOVED***
	// Abort immediately if all tasks are clean.
	if len(dirtySlots) == 0 ***REMOVED***
		if service.UpdateStatus != nil &&
			(service.UpdateStatus.State == api.UpdateStatus_UPDATING ||
				service.UpdateStatus.State == api.UpdateStatus_ROLLBACK_STARTED) ***REMOVED***
			u.completeUpdate(ctx, service.ID)
		***REMOVED***
		return
	***REMOVED***

	// If there's no update in progress, we are starting one.
	if service.UpdateStatus == nil ***REMOVED***
		u.startUpdate(ctx, service.ID)
	***REMOVED***

	var (
		monitoringPeriod time.Duration
		updateConfig     *api.UpdateConfig
	)

	if service.UpdateStatus != nil && service.UpdateStatus.State == api.UpdateStatus_ROLLBACK_STARTED ***REMOVED***
		monitoringPeriod, _ = gogotypes.DurationFromProto(defaults.Service.Rollback.Monitor)
		updateConfig = service.Spec.Rollback
		if updateConfig == nil ***REMOVED***
			updateConfig = defaults.Service.Rollback
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		monitoringPeriod, _ = gogotypes.DurationFromProto(defaults.Service.Update.Monitor)
		updateConfig = service.Spec.Update
		if updateConfig == nil ***REMOVED***
			updateConfig = defaults.Service.Update
		***REMOVED***
	***REMOVED***

	parallelism := int(updateConfig.Parallelism)
	if updateConfig.Monitor != nil ***REMOVED***
		newMonitoringPeriod, err := gogotypes.DurationFromProto(updateConfig.Monitor)
		if err == nil ***REMOVED***
			monitoringPeriod = newMonitoringPeriod
		***REMOVED***
	***REMOVED***

	if parallelism == 0 ***REMOVED***
		// TODO(aluzzardi): We could try to optimize unlimited parallelism by performing updates in a single
		// goroutine using a batch transaction.
		parallelism = len(dirtySlots)
	***REMOVED***

	// Start the workers.
	slotQueue := make(chan orchestrator.Slot)
	wg := sync.WaitGroup***REMOVED******REMOVED***
	wg.Add(parallelism)
	for i := 0; i < parallelism; i++ ***REMOVED***
		go func() ***REMOVED***
			u.worker(ctx, slotQueue, updateConfig)
			wg.Done()
		***REMOVED***()
	***REMOVED***

	var failedTaskWatch chan events.Event

	if updateConfig.FailureAction != api.UpdateConfig_CONTINUE ***REMOVED***
		var cancelWatch func()
		failedTaskWatch, cancelWatch = state.Watch(
			u.store.WatchQueue(),
			api.EventUpdateTask***REMOVED***
				Task:   &api.Task***REMOVED***ServiceID: service.ID, Status: api.TaskStatus***REMOVED***State: api.TaskStateRunning***REMOVED******REMOVED***,
				Checks: []api.TaskCheckFunc***REMOVED***api.TaskCheckServiceID, state.TaskCheckStateGreaterThan***REMOVED***,
			***REMOVED***,
		)
		defer cancelWatch()
	***REMOVED***

	stopped := false
	failedTasks := make(map[string]struct***REMOVED******REMOVED***)
	totalFailures := 0

	failureTriggersAction := func(failedTask *api.Task) bool ***REMOVED***
		// Ignore tasks we have already seen as failures.
		if _, found := failedTasks[failedTask.ID]; found ***REMOVED***
			return false
		***REMOVED***

		// If this failed/completed task is one that we
		// created as part of this update, we should
		// follow the failure action.
		u.updatedTasksMu.Lock()
		startedAt, found := u.updatedTasks[failedTask.ID]
		u.updatedTasksMu.Unlock()

		if found && (startedAt.IsZero() || time.Since(startedAt) <= monitoringPeriod) ***REMOVED***
			failedTasks[failedTask.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			totalFailures++
			if float32(totalFailures)/float32(len(dirtySlots)) > updateConfig.MaxFailureRatio ***REMOVED***
				switch updateConfig.FailureAction ***REMOVED***
				case api.UpdateConfig_PAUSE:
					stopped = true
					message := fmt.Sprintf("update paused due to failure or early termination of task %s", failedTask.ID)
					u.pauseUpdate(ctx, service.ID, message)
					return true
				case api.UpdateConfig_ROLLBACK:
					// Never roll back a rollback
					if service.UpdateStatus != nil && service.UpdateStatus.State == api.UpdateStatus_ROLLBACK_STARTED ***REMOVED***
						message := fmt.Sprintf("rollback paused due to failure or early termination of task %s", failedTask.ID)
						u.pauseUpdate(ctx, service.ID, message)
						return true
					***REMOVED***
					stopped = true
					message := fmt.Sprintf("update rolled back due to failure or early termination of task %s", failedTask.ID)
					u.rollbackUpdate(ctx, service.ID, message)
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED***

		return false
	***REMOVED***

slotsLoop:
	for _, slot := range dirtySlots ***REMOVED***
	retryLoop:
		for ***REMOVED***
			// Wait for a worker to pick up the task or abort the update, whichever comes first.
			select ***REMOVED***
			case <-u.stopChan:
				stopped = true
				break slotsLoop
			case ev := <-failedTaskWatch:
				if failureTriggersAction(ev.(api.EventUpdateTask).Task) ***REMOVED***
					break slotsLoop
				***REMOVED***
			case slotQueue <- slot:
				break retryLoop
			***REMOVED***
		***REMOVED***
	***REMOVED***

	close(slotQueue)
	wg.Wait()

	if !stopped ***REMOVED***
		// Keep watching for task failures for one more monitoringPeriod,
		// before declaring the update complete.
		doneMonitoring := time.After(monitoringPeriod)
	monitorLoop:
		for ***REMOVED***
			select ***REMOVED***
			case <-u.stopChan:
				stopped = true
				break monitorLoop
			case <-doneMonitoring:
				break monitorLoop
			case ev := <-failedTaskWatch:
				if failureTriggersAction(ev.(api.EventUpdateTask).Task) ***REMOVED***
					break monitorLoop
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// TODO(aaronl): Potentially roll back the service if not enough tasks
	// have reached RUNNING by this point.

	if !stopped ***REMOVED***
		u.completeUpdate(ctx, service.ID)
	***REMOVED***
***REMOVED***

func (u *Updater) worker(ctx context.Context, queue <-chan orchestrator.Slot, updateConfig *api.UpdateConfig) ***REMOVED***
	for slot := range queue ***REMOVED***
		// Do we have a task with the new spec in desired state = RUNNING?
		// If so, all we have to do to complete the update is remove the
		// other tasks. Or if we have a task with the new spec that has
		// desired state < RUNNING, advance it to running and remove the
		// other tasks.
		var (
			runningTask *api.Task
			cleanTask   *api.Task
		)
		for _, t := range slot ***REMOVED***
			if !u.isTaskDirty(t) ***REMOVED***
				if t.DesiredState == api.TaskStateRunning ***REMOVED***
					runningTask = t
					break
				***REMOVED***
				if t.DesiredState < api.TaskStateRunning ***REMOVED***
					cleanTask = t
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if runningTask != nil ***REMOVED***
			if err := u.useExistingTask(ctx, slot, runningTask); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("update failed")
			***REMOVED***
		***REMOVED*** else if cleanTask != nil ***REMOVED***
			if err := u.useExistingTask(ctx, slot, cleanTask); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("update failed")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			updated := orchestrator.NewTask(u.cluster, u.newService, slot[0].Slot, "")
			if orchestrator.IsGlobalService(u.newService) ***REMOVED***
				updated = orchestrator.NewTask(u.cluster, u.newService, slot[0].Slot, slot[0].NodeID)
			***REMOVED***
			updated.DesiredState = api.TaskStateReady

			if err := u.updateTask(ctx, slot, updated, updateConfig.Order); err != nil ***REMOVED***
				log.G(ctx).WithError(err).WithField("task.id", updated.ID).Error("update failed")
			***REMOVED***
		***REMOVED***

		if updateConfig.Delay != 0 ***REMOVED***
			select ***REMOVED***
			case <-time.After(updateConfig.Delay):
			case <-u.stopChan:
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (u *Updater) updateTask(ctx context.Context, slot orchestrator.Slot, updated *api.Task, order api.UpdateConfig_UpdateOrder) error ***REMOVED***
	// Kick off the watch before even creating the updated task. This is in order to avoid missing any event.
	taskUpdates, cancel := state.Watch(u.watchQueue, api.EventUpdateTask***REMOVED***
		Task:   &api.Task***REMOVED***ID: updated.ID***REMOVED***,
		Checks: []api.TaskCheckFunc***REMOVED***api.TaskCheckID***REMOVED***,
	***REMOVED***)
	defer cancel()

	// Create an empty entry for this task, so the updater knows a failure
	// should count towards the failure count. The timestamp is added
	// if/when the task reaches RUNNING.
	u.updatedTasksMu.Lock()
	u.updatedTasks[updated.ID] = time.Time***REMOVED******REMOVED***
	u.updatedTasksMu.Unlock()

	startThenStop := false
	var delayStartCh <-chan struct***REMOVED******REMOVED***
	// Atomically create the updated task and bring down the old one.
	err := u.store.Batch(func(batch *store.Batch) error ***REMOVED***
		err := batch.Update(func(tx store.Tx) error ***REMOVED***
			if store.GetService(tx, updated.ServiceID) == nil ***REMOVED***
				return errors.New("service was deleted")
			***REMOVED***

			return store.CreateTask(tx, updated)
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if order == api.UpdateConfig_START_FIRST ***REMOVED***
			delayStartCh = u.restarts.DelayStart(ctx, nil, nil, updated.ID, 0, false)
			startThenStop = true
		***REMOVED*** else ***REMOVED***
			oldTask, err := u.removeOldTasks(ctx, batch, slot)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			delayStartCh = u.restarts.DelayStart(ctx, nil, oldTask, updated.ID, 0, true)
		***REMOVED***

		return nil

	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if delayStartCh != nil ***REMOVED***
		select ***REMOVED***
		case <-delayStartCh:
		case <-u.stopChan:
			return nil
		***REMOVED***
	***REMOVED***

	// Wait for the new task to come up.
	// TODO(aluzzardi): Consider adding a timeout here.
	for ***REMOVED***
		select ***REMOVED***
		case e := <-taskUpdates:
			updated = e.(api.EventUpdateTask).Task
			if updated.Status.State >= api.TaskStateRunning ***REMOVED***
				u.updatedTasksMu.Lock()
				u.updatedTasks[updated.ID] = time.Now()
				u.updatedTasksMu.Unlock()

				if startThenStop && updated.Status.State == api.TaskStateRunning ***REMOVED***
					err := u.store.Batch(func(batch *store.Batch) error ***REMOVED***
						_, err := u.removeOldTasks(ctx, batch, slot)
						if err != nil ***REMOVED***
							log.G(ctx).WithError(err).WithField("task.id", updated.ID).Warning("failed to remove old task after starting replacement")
						***REMOVED***
						return nil
					***REMOVED***)
					return err
				***REMOVED***
				return nil
			***REMOVED***
		case <-u.stopChan:
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (u *Updater) useExistingTask(ctx context.Context, slot orchestrator.Slot, existing *api.Task) error ***REMOVED***
	var removeTasks []*api.Task
	for _, t := range slot ***REMOVED***
		if t != existing ***REMOVED***
			removeTasks = append(removeTasks, t)
		***REMOVED***
	***REMOVED***
	if len(removeTasks) != 0 || existing.DesiredState != api.TaskStateRunning ***REMOVED***
		var delayStartCh <-chan struct***REMOVED******REMOVED***
		err := u.store.Batch(func(batch *store.Batch) error ***REMOVED***
			var oldTask *api.Task
			if len(removeTasks) != 0 ***REMOVED***
				var err error
				oldTask, err = u.removeOldTasks(ctx, batch, removeTasks)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***

			if existing.DesiredState != api.TaskStateRunning ***REMOVED***
				delayStartCh = u.restarts.DelayStart(ctx, nil, oldTask, existing.ID, 0, true)
			***REMOVED***
			return nil
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if delayStartCh != nil ***REMOVED***
			select ***REMOVED***
			case <-delayStartCh:
			case <-u.stopChan:
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// removeOldTasks shuts down the given tasks and returns one of the tasks that
// was shut down, or an error.
func (u *Updater) removeOldTasks(ctx context.Context, batch *store.Batch, removeTasks []*api.Task) (*api.Task, error) ***REMOVED***
	var (
		lastErr     error
		removedTask *api.Task
	)
	for _, original := range removeTasks ***REMOVED***
		if original.DesiredState > api.TaskStateRunning ***REMOVED***
			continue
		***REMOVED***
		err := batch.Update(func(tx store.Tx) error ***REMOVED***
			t := store.GetTask(tx, original.ID)
			if t == nil ***REMOVED***
				return fmt.Errorf("task %s not found while trying to shut it down", original.ID)
			***REMOVED***
			if t.DesiredState > api.TaskStateRunning ***REMOVED***
				return fmt.Errorf("task %s was already shut down when reached by updater", original.ID)
			***REMOVED***
			t.DesiredState = api.TaskStateShutdown
			return store.UpdateTask(tx, t)
		***REMOVED***)
		if err != nil ***REMOVED***
			lastErr = err
		***REMOVED*** else ***REMOVED***
			removedTask = original
		***REMOVED***
	***REMOVED***

	if removedTask == nil ***REMOVED***
		return nil, lastErr
	***REMOVED***
	return removedTask, nil
***REMOVED***

func (u *Updater) isTaskDirty(t *api.Task) bool ***REMOVED***
	return orchestrator.IsTaskDirty(u.newService, t)
***REMOVED***

func (u *Updater) isSlotDirty(slot orchestrator.Slot) bool ***REMOVED***
	return len(slot) > 1 || (len(slot) == 1 && u.isTaskDirty(slot[0]))
***REMOVED***

func (u *Updater) startUpdate(ctx context.Context, serviceID string) ***REMOVED***
	err := u.store.Update(func(tx store.Tx) error ***REMOVED***
		service := store.GetService(tx, serviceID)
		if service == nil ***REMOVED***
			return nil
		***REMOVED***
		if service.UpdateStatus != nil ***REMOVED***
			return nil
		***REMOVED***

		service.UpdateStatus = &api.UpdateStatus***REMOVED***
			State:     api.UpdateStatus_UPDATING,
			Message:   "update in progress",
			StartedAt: ptypes.MustTimestampProto(time.Now()),
		***REMOVED***

		return store.UpdateService(tx, service)
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("failed to mark update of service %s in progress", serviceID)
	***REMOVED***
***REMOVED***

func (u *Updater) pauseUpdate(ctx context.Context, serviceID, message string) ***REMOVED***
	log.G(ctx).Debugf("pausing update of service %s", serviceID)

	err := u.store.Update(func(tx store.Tx) error ***REMOVED***
		service := store.GetService(tx, serviceID)
		if service == nil ***REMOVED***
			return nil
		***REMOVED***
		if service.UpdateStatus == nil ***REMOVED***
			// The service was updated since we started this update
			return nil
		***REMOVED***

		if service.UpdateStatus.State == api.UpdateStatus_ROLLBACK_STARTED ***REMOVED***
			service.UpdateStatus.State = api.UpdateStatus_ROLLBACK_PAUSED
		***REMOVED*** else ***REMOVED***
			service.UpdateStatus.State = api.UpdateStatus_PAUSED
		***REMOVED***
		service.UpdateStatus.Message = message

		return store.UpdateService(tx, service)
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("failed to pause update of service %s", serviceID)
	***REMOVED***
***REMOVED***

func (u *Updater) rollbackUpdate(ctx context.Context, serviceID, message string) ***REMOVED***
	log.G(ctx).Debugf("starting rollback of service %s", serviceID)

	err := u.store.Update(func(tx store.Tx) error ***REMOVED***
		service := store.GetService(tx, serviceID)
		if service == nil ***REMOVED***
			return nil
		***REMOVED***
		if service.UpdateStatus == nil ***REMOVED***
			// The service was updated since we started this update
			return nil
		***REMOVED***

		service.UpdateStatus.State = api.UpdateStatus_ROLLBACK_STARTED
		service.UpdateStatus.Message = message

		if service.PreviousSpec == nil ***REMOVED***
			return errors.New("cannot roll back service because no previous spec is available")
		***REMOVED***
		service.Spec = *service.PreviousSpec
		service.SpecVersion = service.PreviousSpecVersion.Copy()
		service.PreviousSpec = nil
		service.PreviousSpecVersion = nil

		return store.UpdateService(tx, service)
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("failed to start rollback of service %s", serviceID)
		return
	***REMOVED***
***REMOVED***

func (u *Updater) completeUpdate(ctx context.Context, serviceID string) ***REMOVED***
	log.G(ctx).Debugf("update of service %s complete", serviceID)

	err := u.store.Update(func(tx store.Tx) error ***REMOVED***
		service := store.GetService(tx, serviceID)
		if service == nil ***REMOVED***
			return nil
		***REMOVED***
		if service.UpdateStatus == nil ***REMOVED***
			// The service was changed since we started this update
			return nil
		***REMOVED***
		if service.UpdateStatus.State == api.UpdateStatus_ROLLBACK_STARTED ***REMOVED***
			service.UpdateStatus.State = api.UpdateStatus_ROLLBACK_COMPLETED
			service.UpdateStatus.Message = "rollback completed"
		***REMOVED*** else ***REMOVED***
			service.UpdateStatus.State = api.UpdateStatus_COMPLETED
			service.UpdateStatus.Message = "update completed"
		***REMOVED***
		service.UpdateStatus.CompletedAt = ptypes.MustTimestampProto(time.Now())

		return store.UpdateService(tx, service)
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("failed to mark update of service %s complete", serviceID)
	***REMOVED***
***REMOVED***
