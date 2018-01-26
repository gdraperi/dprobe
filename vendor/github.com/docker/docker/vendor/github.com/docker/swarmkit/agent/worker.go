package agent

import (
	"sync"

	"github.com/boltdb/bolt"
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/watch"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// Worker implements the core task management logic and persistence. It
// coordinates the set of assignments with the executor.
type Worker interface ***REMOVED***
	// Init prepares the worker for task assignment.
	Init(ctx context.Context) error

	// Close performs worker cleanup when no longer needed.
	//
	// It is not safe to call any worker function after that.
	Close()

	// Assign assigns a complete set of tasks and configs/secrets to a
	// worker. Any items not included in this set will be removed.
	Assign(ctx context.Context, assignments []*api.AssignmentChange) error

	// Updates updates an incremental set of tasks or configs/secrets of
	// the worker. Any items not included either in added or removed will
	// remain untouched.
	Update(ctx context.Context, assignments []*api.AssignmentChange) error

	// Listen to updates about tasks controlled by the worker. When first
	// called, the reporter will receive all updates for all tasks controlled
	// by the worker.
	//
	// The listener will be removed if the context is cancelled.
	Listen(ctx context.Context, reporter StatusReporter)

	// Subscribe to log messages matching the subscription.
	Subscribe(ctx context.Context, subscription *api.SubscriptionMessage) error

	// Wait blocks until all task managers have closed
	Wait(ctx context.Context) error
***REMOVED***

// statusReporterKey protects removal map from panic.
type statusReporterKey struct ***REMOVED***
	StatusReporter
***REMOVED***

type worker struct ***REMOVED***
	db                *bolt.DB
	executor          exec.Executor
	publisher         exec.LogPublisher
	listeners         map[*statusReporterKey]struct***REMOVED******REMOVED***
	taskevents        *watch.Queue
	publisherProvider exec.LogPublisherProvider

	taskManagers map[string]*taskManager
	mu           sync.RWMutex

	closed  bool
	closers sync.WaitGroup // keeps track of active closers
***REMOVED***

func newWorker(db *bolt.DB, executor exec.Executor, publisherProvider exec.LogPublisherProvider) *worker ***REMOVED***
	return &worker***REMOVED***
		db:                db,
		executor:          executor,
		publisherProvider: publisherProvider,
		taskevents:        watch.NewQueue(),
		listeners:         make(map[*statusReporterKey]struct***REMOVED******REMOVED***),
		taskManagers:      make(map[string]*taskManager),
	***REMOVED***
***REMOVED***

// Init prepares the worker for assignments.
func (w *worker) Init(ctx context.Context) error ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()

	ctx = log.WithModule(ctx, "worker")

	// TODO(stevvooe): Start task cleanup process.

	// read the tasks from the database and start any task managers that may be needed.
	return w.db.Update(func(tx *bolt.Tx) error ***REMOVED***
		return WalkTasks(tx, func(task *api.Task) error ***REMOVED***
			if !TaskAssigned(tx, task.ID) ***REMOVED***
				// NOTE(stevvooe): If tasks can survive worker restart, we need
				// to startup the controller and ensure they are removed. For
				// now, we can simply remove them from the database.
				if err := DeleteTask(tx, task.ID); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Errorf("error removing task %v", task.ID)
				***REMOVED***
				return nil
			***REMOVED***

			status, err := GetTaskStatus(tx, task.ID)
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("unable to read tasks status")
				return nil
			***REMOVED***

			task.Status = *status // merges the status into the task, ensuring we start at the right point.
			return w.startTask(ctx, tx, task)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

// Close performs worker cleanup when no longer needed.
func (w *worker) Close() ***REMOVED***
	w.mu.Lock()
	w.closed = true
	w.mu.Unlock()

	w.taskevents.Close()
***REMOVED***

// Assign assigns a full set of tasks, configs, and secrets to the worker.
// Any tasks not previously known will be started. Any tasks that are in the task set
// and already running will be updated, if possible. Any tasks currently running on
// the worker outside the task set will be terminated.
// Anything not in the set of assignments will be removed.
func (w *worker) Assign(ctx context.Context, assignments []*api.AssignmentChange) error ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed ***REMOVED***
		return ErrClosed
	***REMOVED***

	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"len(assignments)": len(assignments),
	***REMOVED***).Debug("(*worker).Assign")

	// Need to update dependencies before tasks

	err := reconcileSecrets(ctx, w, assignments, true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = reconcileConfigs(ctx, w, assignments, true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return reconcileTaskState(ctx, w, assignments, true)
***REMOVED***

// Update updates the set of tasks, configs, and secrets for the worker.
// Tasks in the added set will be added to the worker, and tasks in the removed set
// will be removed from the worker
// Secrets in the added set will be added to the worker, and secrets in the removed set
// will be removed from the worker.
// Configs in the added set will be added to the worker, and configs in the removed set
// will be removed from the worker.
func (w *worker) Update(ctx context.Context, assignments []*api.AssignmentChange) error ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed ***REMOVED***
		return ErrClosed
	***REMOVED***

	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"len(assignments)": len(assignments),
	***REMOVED***).Debug("(*worker).Update")

	err := reconcileSecrets(ctx, w, assignments, false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = reconcileConfigs(ctx, w, assignments, false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return reconcileTaskState(ctx, w, assignments, false)
***REMOVED***

func reconcileTaskState(ctx context.Context, w *worker, assignments []*api.AssignmentChange, fullSnapshot bool) error ***REMOVED***
	var (
		updatedTasks []*api.Task
		removedTasks []*api.Task
	)
	for _, a := range assignments ***REMOVED***
		if t := a.Assignment.GetTask(); t != nil ***REMOVED***
			switch a.Action ***REMOVED***
			case api.AssignmentChange_AssignmentActionUpdate:
				updatedTasks = append(updatedTasks, t)
			case api.AssignmentChange_AssignmentActionRemove:
				removedTasks = append(removedTasks, t)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"len(updatedTasks)": len(updatedTasks),
		"len(removedTasks)": len(removedTasks),
	***REMOVED***).Debug("(*worker).reconcileTaskState")

	tx, err := w.db.Begin(true)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("failed starting transaction against task database")
		return err
	***REMOVED***
	defer tx.Rollback()

	assigned := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***

	for _, task := range updatedTasks ***REMOVED***
		log.G(ctx).WithFields(
			logrus.Fields***REMOVED***
				"task.id":           task.ID,
				"task.desiredstate": task.DesiredState***REMOVED***).Debug("assigned")
		if err := PutTask(tx, task); err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := SetTaskAssignment(tx, task.ID, true); err != nil ***REMOVED***
			return err
		***REMOVED***

		if mgr, ok := w.taskManagers[task.ID]; ok ***REMOVED***
			if err := mgr.Update(ctx, task); err != nil && err != ErrClosed ***REMOVED***
				log.G(ctx).WithError(err).Error("failed updating assigned task")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// we may have still seen the task, let's grab the status from
			// storage and replace it with our status, if we have it.
			status, err := GetTaskStatus(tx, task.ID)
			if err != nil ***REMOVED***
				if err != errTaskUnknown ***REMOVED***
					return err
				***REMOVED***

				// never seen before, register the provided status
				if err := PutTaskStatus(tx, task.ID, &task.Status); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				task.Status = *status
			***REMOVED***
			w.startTask(ctx, tx, task)
		***REMOVED***

		assigned[task.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	closeManager := func(tm *taskManager) ***REMOVED***
		go func(tm *taskManager) ***REMOVED***
			defer w.closers.Done()
			// when a task is no longer assigned, we shutdown the task manager
			if err := tm.Close(); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("error closing task manager")
			***REMOVED***
		***REMOVED***(tm)

		// make an attempt at removing. this is best effort. any errors will be
		// retried by the reaper later.
		if err := tm.ctlr.Remove(ctx); err != nil ***REMOVED***
			log.G(ctx).WithError(err).WithField("task.id", tm.task.ID).Error("remove task failed")
		***REMOVED***

		if err := tm.ctlr.Close(); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("error closing controller")
		***REMOVED***
	***REMOVED***

	removeTaskAssignment := func(taskID string) error ***REMOVED***
		ctx := log.WithLogger(ctx, log.G(ctx).WithField("task.id", taskID))
		if err := SetTaskAssignment(tx, taskID, false); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("error setting task assignment in database")
		***REMOVED***
		return err
	***REMOVED***

	// If this was a complete set of assignments, we're going to remove all the remaining
	// tasks.
	if fullSnapshot ***REMOVED***
		for id, tm := range w.taskManagers ***REMOVED***
			if _, ok := assigned[id]; ok ***REMOVED***
				continue
			***REMOVED***

			err := removeTaskAssignment(id)
			if err == nil ***REMOVED***
				delete(w.taskManagers, id)
				go closeManager(tm)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// If this was an incremental set of assignments, we're going to remove only the tasks
		// in the removed set
		for _, task := range removedTasks ***REMOVED***
			err := removeTaskAssignment(task.ID)
			if err != nil ***REMOVED***
				continue
			***REMOVED***

			tm, ok := w.taskManagers[task.ID]
			if ok ***REMOVED***
				delete(w.taskManagers, task.ID)
				go closeManager(tm)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return tx.Commit()
***REMOVED***

func reconcileSecrets(ctx context.Context, w *worker, assignments []*api.AssignmentChange, fullSnapshot bool) error ***REMOVED***
	var (
		updatedSecrets []api.Secret
		removedSecrets []string
	)
	for _, a := range assignments ***REMOVED***
		if s := a.Assignment.GetSecret(); s != nil ***REMOVED***
			switch a.Action ***REMOVED***
			case api.AssignmentChange_AssignmentActionUpdate:
				updatedSecrets = append(updatedSecrets, *s)
			case api.AssignmentChange_AssignmentActionRemove:
				removedSecrets = append(removedSecrets, s.ID)
			***REMOVED***

		***REMOVED***
	***REMOVED***

	secretsProvider, ok := w.executor.(exec.SecretsProvider)
	if !ok ***REMOVED***
		if len(updatedSecrets) != 0 || len(removedSecrets) != 0 ***REMOVED***
			log.G(ctx).Warn("secrets update ignored; executor does not support secrets")
		***REMOVED***
		return nil
	***REMOVED***

	secrets := secretsProvider.Secrets()

	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"len(updatedSecrets)": len(updatedSecrets),
		"len(removedSecrets)": len(removedSecrets),
	***REMOVED***).Debug("(*worker).reconcileSecrets")

	// If this was a complete set of secrets, we're going to clear the secrets map and add all of them
	if fullSnapshot ***REMOVED***
		secrets.Reset()
	***REMOVED*** else ***REMOVED***
		secrets.Remove(removedSecrets)
	***REMOVED***
	secrets.Add(updatedSecrets...)

	return nil
***REMOVED***

func reconcileConfigs(ctx context.Context, w *worker, assignments []*api.AssignmentChange, fullSnapshot bool) error ***REMOVED***
	var (
		updatedConfigs []api.Config
		removedConfigs []string
	)
	for _, a := range assignments ***REMOVED***
		if r := a.Assignment.GetConfig(); r != nil ***REMOVED***
			switch a.Action ***REMOVED***
			case api.AssignmentChange_AssignmentActionUpdate:
				updatedConfigs = append(updatedConfigs, *r)
			case api.AssignmentChange_AssignmentActionRemove:
				removedConfigs = append(removedConfigs, r.ID)
			***REMOVED***

		***REMOVED***
	***REMOVED***

	configsProvider, ok := w.executor.(exec.ConfigsProvider)
	if !ok ***REMOVED***
		if len(updatedConfigs) != 0 || len(removedConfigs) != 0 ***REMOVED***
			log.G(ctx).Warn("configs update ignored; executor does not support configs")
		***REMOVED***
		return nil
	***REMOVED***

	configs := configsProvider.Configs()

	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"len(updatedConfigs)": len(updatedConfigs),
		"len(removedConfigs)": len(removedConfigs),
	***REMOVED***).Debug("(*worker).reconcileConfigs")

	// If this was a complete set of configs, we're going to clear the configs map and add all of them
	if fullSnapshot ***REMOVED***
		configs.Reset()
	***REMOVED*** else ***REMOVED***
		configs.Remove(removedConfigs)
	***REMOVED***
	configs.Add(updatedConfigs...)

	return nil
***REMOVED***

func (w *worker) Listen(ctx context.Context, reporter StatusReporter) ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()

	key := &statusReporterKey***REMOVED***reporter***REMOVED***
	w.listeners[key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	go func() ***REMOVED***
		<-ctx.Done()
		w.mu.Lock()
		defer w.mu.Unlock()
		delete(w.listeners, key) // remove the listener if the context is closed.
	***REMOVED***()

	// report the current statuses to the new listener
	if err := w.db.View(func(tx *bolt.Tx) error ***REMOVED***
		return WalkTaskStatus(tx, func(id string, status *api.TaskStatus) error ***REMOVED***
			return reporter.UpdateTaskStatus(ctx, id, status)
		***REMOVED***)
	***REMOVED***); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("failed reporting initial statuses to registered listener %v", reporter)
	***REMOVED***
***REMOVED***

func (w *worker) startTask(ctx context.Context, tx *bolt.Tx, task *api.Task) error ***REMOVED***
	_, err := w.taskManager(ctx, tx, task) // side-effect taskManager creation.

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("failed to start taskManager")
		// we ignore this error: it gets reported in the taskStatus within
		// `newTaskManager`. We log it here and move on. If their is an
		// attempted restart, the lack of taskManager will have this retry
		// again.
		return nil
	***REMOVED***

	// only publish if controller resolution was successful.
	w.taskevents.Publish(task.Copy())
	return nil
***REMOVED***

func (w *worker) taskManager(ctx context.Context, tx *bolt.Tx, task *api.Task) (*taskManager, error) ***REMOVED***
	if tm, ok := w.taskManagers[task.ID]; ok ***REMOVED***
		return tm, nil
	***REMOVED***

	tm, err := w.newTaskManager(ctx, tx, task)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	w.taskManagers[task.ID] = tm
	// keep track of active tasks
	w.closers.Add(1)
	return tm, nil
***REMOVED***

func (w *worker) newTaskManager(ctx context.Context, tx *bolt.Tx, task *api.Task) (*taskManager, error) ***REMOVED***
	ctx = log.WithLogger(ctx, log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"task.id":    task.ID,
		"service.id": task.ServiceID,
	***REMOVED***))

	ctlr, status, err := exec.Resolve(ctx, task, w.executor)
	if err := w.updateTaskStatus(ctx, tx, task.ID, status); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("error updating task status after controller resolution")
	***REMOVED***

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("controller resolution failed")
		return nil, err
	***REMOVED***

	return newTaskManager(ctx, task, ctlr, statusReporterFunc(func(ctx context.Context, taskID string, status *api.TaskStatus) error ***REMOVED***
		w.mu.RLock()
		defer w.mu.RUnlock()

		return w.db.Update(func(tx *bolt.Tx) error ***REMOVED***
			return w.updateTaskStatus(ctx, tx, taskID, status)
		***REMOVED***)
	***REMOVED***)), nil
***REMOVED***

// updateTaskStatus reports statuses to listeners, read lock must be held.
func (w *worker) updateTaskStatus(ctx context.Context, tx *bolt.Tx, taskID string, status *api.TaskStatus) error ***REMOVED***
	if err := PutTaskStatus(tx, taskID, status); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("failed writing status to disk")
		return err
	***REMOVED***

	// broadcast the task status out.
	for key := range w.listeners ***REMOVED***
		if err := key.StatusReporter.UpdateTaskStatus(ctx, taskID, status); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed updating status for reporter %v", key.StatusReporter)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Subscribe to log messages matching the subscription.
func (w *worker) Subscribe(ctx context.Context, subscription *api.SubscriptionMessage) error ***REMOVED***
	log.G(ctx).Debugf("Received subscription %s (selector: %v)", subscription.ID, subscription.Selector)

	publisher, cancel, err := w.publisherProvider.Publisher(ctx, subscription.ID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// Send a close once we're done
	defer cancel()

	match := func(t *api.Task) bool ***REMOVED***
		// TODO(aluzzardi): Consider using maps to limit the iterations.
		for _, tid := range subscription.Selector.TaskIDs ***REMOVED***
			if t.ID == tid ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***

		for _, sid := range subscription.Selector.ServiceIDs ***REMOVED***
			if t.ServiceID == sid ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***

		for _, nid := range subscription.Selector.NodeIDs ***REMOVED***
			if t.NodeID == nid ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***

		return false
	***REMOVED***

	wg := sync.WaitGroup***REMOVED******REMOVED***
	w.mu.Lock()
	for _, tm := range w.taskManagers ***REMOVED***
		if match(tm.task) ***REMOVED***
			wg.Add(1)
			go func(tm *taskManager) ***REMOVED***
				defer wg.Done()
				tm.Logs(ctx, *subscription.Options, publisher)
			***REMOVED***(tm)
		***REMOVED***
	***REMOVED***
	w.mu.Unlock()

	// If follow mode is disabled, wait for the current set of matched tasks
	// to finish publishing logs, then close the subscription by returning.
	if subscription.Options == nil || !subscription.Options.Follow ***REMOVED***
		waitCh := make(chan struct***REMOVED******REMOVED***)
		go func() ***REMOVED***
			defer close(waitCh)
			wg.Wait()
		***REMOVED***()

		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case <-waitCh:
			return nil
		***REMOVED***
	***REMOVED***

	// In follow mode, watch for new tasks. Don't close the subscription
	// until it's cancelled.
	ch, cancel := w.taskevents.Watch()
	defer cancel()
	for ***REMOVED***
		select ***REMOVED***
		case v := <-ch:
			task := v.(*api.Task)
			if match(task) ***REMOVED***
				w.mu.RLock()
				tm, ok := w.taskManagers[task.ID]
				w.mu.RUnlock()
				if !ok ***REMOVED***
					continue
				***REMOVED***

				go tm.Logs(ctx, *subscription.Options, publisher)
			***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		***REMOVED***
	***REMOVED***
***REMOVED***

func (w *worker) Wait(ctx context.Context) error ***REMOVED***
	ch := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		w.closers.Wait()
		close(ch)
	***REMOVED***()

	select ***REMOVED***
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	***REMOVED***
***REMOVED***
