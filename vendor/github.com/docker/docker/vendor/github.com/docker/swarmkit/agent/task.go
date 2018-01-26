package agent

import (
	"sync"
	"time"

	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/equality"
	"github.com/docker/swarmkit/log"
	"golang.org/x/net/context"
)

// taskManager manages all aspects of task execution and reporting for an agent
// through state management.
type taskManager struct ***REMOVED***
	task     *api.Task
	ctlr     exec.Controller
	reporter StatusReporter

	updateq chan *api.Task

	shutdown     chan struct***REMOVED******REMOVED***
	shutdownOnce sync.Once
	closed       chan struct***REMOVED******REMOVED***
	closeOnce    sync.Once
***REMOVED***

func newTaskManager(ctx context.Context, task *api.Task, ctlr exec.Controller, reporter StatusReporter) *taskManager ***REMOVED***
	t := &taskManager***REMOVED***
		task:     task.Copy(),
		ctlr:     ctlr,
		reporter: reporter,
		updateq:  make(chan *api.Task),
		shutdown: make(chan struct***REMOVED******REMOVED***),
		closed:   make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
	go t.run(ctx)
	return t
***REMOVED***

// Update the task data.
func (tm *taskManager) Update(ctx context.Context, task *api.Task) error ***REMOVED***
	select ***REMOVED***
	case tm.updateq <- task:
		return nil
	case <-tm.closed:
		return ErrClosed
	case <-ctx.Done():
		return ctx.Err()
	***REMOVED***
***REMOVED***

// Close shuts down the task manager, blocking until it is closed.
func (tm *taskManager) Close() error ***REMOVED***
	tm.shutdownOnce.Do(func() ***REMOVED***
		close(tm.shutdown)
	***REMOVED***)

	<-tm.closed

	return nil
***REMOVED***

func (tm *taskManager) Logs(ctx context.Context, options api.LogSubscriptionOptions, publisher exec.LogPublisher) ***REMOVED***
	ctx = log.WithModule(ctx, "taskmanager")

	logCtlr, ok := tm.ctlr.(exec.ControllerLogs)
	if !ok ***REMOVED***
		return // no logs available
	***REMOVED***
	if err := logCtlr.Logs(ctx, publisher, options); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("logs call failed")
	***REMOVED***
***REMOVED***

func (tm *taskManager) run(ctx context.Context) ***REMOVED***
	ctx, cancelAll := context.WithCancel(ctx)
	defer cancelAll() // cancel all child operations on exit.

	ctx = log.WithModule(ctx, "taskmanager")

	var (
		opctx    context.Context
		cancel   context.CancelFunc
		run      = make(chan struct***REMOVED******REMOVED***, 1)
		statusq  = make(chan *api.TaskStatus)
		errs     = make(chan error)
		shutdown = tm.shutdown
		updated  bool // true if the task was updated.
	)

	defer func() ***REMOVED***
		// closure  picks up current value of cancel.
		if cancel != nil ***REMOVED***
			cancel()
		***REMOVED***
	***REMOVED***()

	run <- struct***REMOVED******REMOVED******REMOVED******REMOVED*** // prime the pump
	for ***REMOVED***
		select ***REMOVED***
		case <-run:
			// always check for shutdown before running.
			select ***REMOVED***
			case <-tm.shutdown:
				shutdown = tm.shutdown // a little questionable
				continue               // ignore run request and handle shutdown
			case <-tm.closed:
				continue
			default:
			***REMOVED***

			opctx, cancel = context.WithCancel(ctx)

			// Several variables need to be snapshotted for the closure below.
			opcancel := cancel        // fork for the closure
			running := tm.task.Copy() // clone the task before dispatch
			statusqLocal := statusq
			updatedLocal := updated // capture state of update for goroutine
			updated = false
			go runctx(ctx, tm.closed, errs, func(ctx context.Context) error ***REMOVED***
				defer opcancel()

				if updatedLocal ***REMOVED***
					// before we do anything, update the task for the controller.
					// always update the controller before running.
					if err := tm.ctlr.Update(opctx, running); err != nil ***REMOVED***
						log.G(ctx).WithError(err).Error("updating task controller failed")
						return err
					***REMOVED***
				***REMOVED***

				status, err := exec.Do(opctx, running, tm.ctlr)
				if status != nil ***REMOVED***
					// always report the status if we get one back. This
					// returns to the manager loop, then reports the status
					// upstream.
					select ***REMOVED***
					case statusqLocal <- status:
					case <-ctx.Done(): // not opctx, since that may have been cancelled.
					***REMOVED***

					if err := tm.reporter.UpdateTaskStatus(ctx, running.ID, status); err != nil ***REMOVED***
						log.G(ctx).WithError(err).Error("task manager failed to report status to agent")
					***REMOVED***
				***REMOVED***

				return err
			***REMOVED***)
		case err := <-errs:
			// This branch is always executed when an operations completes. The
			// goal is to decide whether or not we re-dispatch the operation.
			cancel = nil

			select ***REMOVED***
			case <-tm.shutdown:
				shutdown = tm.shutdown // re-enable the shutdown branch
				continue               // no dispatch if we are in shutdown.
			default:
			***REMOVED***

			switch err ***REMOVED***
			case exec.ErrTaskNoop:
				if !updated ***REMOVED***
					continue // wait till getting pumped via update.
				***REMOVED***
			case exec.ErrTaskRetry:
				// TODO(stevvooe): Add exponential backoff with random jitter
				// here. For now, this backoff is enough to keep the task
				// manager from running away with the CPU.
				time.AfterFunc(time.Second, func() ***REMOVED***
					errs <- nil // repump this branch, with no err
				***REMOVED***)
				continue
			case nil, context.Canceled, context.DeadlineExceeded:
				// no log in this case
			default:
				log.G(ctx).WithError(err).Error("task operation failed")
			***REMOVED***

			select ***REMOVED***
			case run <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
			default:
			***REMOVED***
		case status := <-statusq:
			tm.task.Status = *status
		case task := <-tm.updateq:
			if equality.TasksEqualStable(task, tm.task) ***REMOVED***
				continue // ignore the update
			***REMOVED***

			if task.ID != tm.task.ID ***REMOVED***
				log.G(ctx).WithField("task.update.id", task.ID).Error("received update for incorrect task")
				continue
			***REMOVED***

			if task.DesiredState < tm.task.DesiredState ***REMOVED***
				log.G(ctx).WithField("task.update.desiredstate", task.DesiredState).
					Error("ignoring task update with invalid desired state")
				continue
			***REMOVED***

			task = task.Copy()
			task.Status = tm.task.Status // overwrite our status, as it is canonical.
			tm.task = task
			updated = true

			// we have accepted the task update
			if cancel != nil ***REMOVED***
				cancel() // cancel outstanding if necessary.
			***REMOVED*** else ***REMOVED***
				// If this channel op fails, it means there is already a
				// message on the run queue.
				select ***REMOVED***
				case run <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
				default:
				***REMOVED***
			***REMOVED***
		case <-shutdown:
			if cancel != nil ***REMOVED***
				// cancel outstanding operation.
				cancel()

				// subtle: after a cancellation, we want to avoid busy wait
				// here. this gets renabled in the errs branch and we'll come
				// back around and try shutdown again.
				shutdown = nil // turn off this branch until op proceeds
				continue       // wait until operation actually exits.
			***REMOVED***

			// disable everything, and prepare for closing.
			statusq = nil
			errs = nil
			shutdown = nil
			tm.closeOnce.Do(func() ***REMOVED***
				close(tm.closed)
			***REMOVED***)
		case <-tm.closed:
			return
		case <-ctx.Done():
			tm.closeOnce.Do(func() ***REMOVED***
				close(tm.closed)
			***REMOVED***)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***
