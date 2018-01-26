package agent

import (
	"reflect"
	"sync"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"golang.org/x/net/context"
)

// StatusReporter receives updates to task status. Method may be called
// concurrently, so implementations should be goroutine-safe.
type StatusReporter interface ***REMOVED***
	UpdateTaskStatus(ctx context.Context, taskID string, status *api.TaskStatus) error
***REMOVED***

type statusReporterFunc func(ctx context.Context, taskID string, status *api.TaskStatus) error

func (fn statusReporterFunc) UpdateTaskStatus(ctx context.Context, taskID string, status *api.TaskStatus) error ***REMOVED***
	return fn(ctx, taskID, status)
***REMOVED***

// statusReporter creates a reliable StatusReporter that will always succeed.
// It handles several tasks at once, ensuring all statuses are reported.
//
// The reporter will continue reporting the current status until it succeeds.
type statusReporter struct ***REMOVED***
	reporter StatusReporter
	statuses map[string]*api.TaskStatus
	mu       sync.Mutex
	cond     sync.Cond
	closed   bool
***REMOVED***

func newStatusReporter(ctx context.Context, upstream StatusReporter) *statusReporter ***REMOVED***
	r := &statusReporter***REMOVED***
		reporter: upstream,
		statuses: make(map[string]*api.TaskStatus),
	***REMOVED***

	r.cond.L = &r.mu

	go r.run(ctx)
	return r
***REMOVED***

func (sr *statusReporter) UpdateTaskStatus(ctx context.Context, taskID string, status *api.TaskStatus) error ***REMOVED***
	sr.mu.Lock()
	defer sr.mu.Unlock()

	current, ok := sr.statuses[taskID]
	if ok ***REMOVED***
		if reflect.DeepEqual(current, status) ***REMOVED***
			return nil
		***REMOVED***

		if current.State > status.State ***REMOVED***
			return nil // ignore old updates
		***REMOVED***
	***REMOVED***
	sr.statuses[taskID] = status
	sr.cond.Signal()

	return nil
***REMOVED***

func (sr *statusReporter) Close() error ***REMOVED***
	sr.mu.Lock()
	defer sr.mu.Unlock()

	sr.closed = true
	sr.cond.Signal()

	return nil
***REMOVED***

func (sr *statusReporter) run(ctx context.Context) ***REMOVED***
	done := make(chan struct***REMOVED******REMOVED***)
	defer close(done)

	sr.mu.Lock() // released during wait, below.
	defer sr.mu.Unlock()

	go func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			sr.Close()
		case <-done:
			return
		***REMOVED***
	***REMOVED***()

	for ***REMOVED***
		if len(sr.statuses) == 0 ***REMOVED***
			sr.cond.Wait()
		***REMOVED***

		if sr.closed ***REMOVED***
			// TODO(stevvooe): Add support here for waiting until all
			// statuses are flushed before shutting down.
			return
		***REMOVED***

		for taskID, status := range sr.statuses ***REMOVED***
			delete(sr.statuses, taskID) // delete the entry, while trying to send.

			sr.mu.Unlock()
			err := sr.reporter.UpdateTaskStatus(ctx, taskID, status)
			sr.mu.Lock()

			// reporter might be closed during UpdateTaskStatus call
			if sr.closed ***REMOVED***
				return
			***REMOVED***

			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("status reporter failed to report status to agent")

				// place it back in the map, if not there, allowing us to pick
				// the value if a new one came in when we were sending the last
				// update.
				if _, ok := sr.statuses[taskID]; !ok ***REMOVED***
					sr.statuses[taskID] = status
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
