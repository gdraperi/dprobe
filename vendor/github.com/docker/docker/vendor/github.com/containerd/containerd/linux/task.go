// +build linux

package linux

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/containerd/cgroups"
	eventstypes "github.com/containerd/containerd/api/events"
	"github.com/containerd/containerd/api/types/task"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/events/exchange"
	"github.com/containerd/containerd/identifiers"
	"github.com/containerd/containerd/linux/shim/client"
	shim "github.com/containerd/containerd/linux/shim/v1"
	"github.com/containerd/containerd/runtime"
	runc "github.com/containerd/go-runc"
	"github.com/gogo/protobuf/types"
)

// Task on a linux based system
type Task struct ***REMOVED***
	mu        sync.Mutex
	id        string
	pid       int
	shim      *client.Client
	namespace string
	cg        cgroups.Cgroup
	monitor   runtime.TaskMonitor
	events    *exchange.Exchange
	runtime   *runc.Runc
***REMOVED***

func newTask(id, namespace string, pid int, shim *client.Client, monitor runtime.TaskMonitor, events *exchange.Exchange, runtime *runc.Runc) (*Task, error) ***REMOVED***
	var (
		err error
		cg  cgroups.Cgroup
	)
	if pid > 0 ***REMOVED***
		cg, err = cgroups.Load(cgroups.V1, cgroups.PidPath(pid))
		if err != nil && err != cgroups.ErrCgroupDeleted ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return &Task***REMOVED***
		id:        id,
		pid:       pid,
		shim:      shim,
		namespace: namespace,
		cg:        cg,
		monitor:   monitor,
		events:    events,
		runtime:   runtime,
	***REMOVED***, nil
***REMOVED***

// ID of the task
func (t *Task) ID() string ***REMOVED***
	return t.id
***REMOVED***

// Info returns task information about the runtime and namespace
func (t *Task) Info() runtime.TaskInfo ***REMOVED***
	return runtime.TaskInfo***REMOVED***
		ID:        t.id,
		Runtime:   pluginID,
		Namespace: t.namespace,
	***REMOVED***
***REMOVED***

// Start the task
func (t *Task) Start(ctx context.Context) error ***REMOVED***
	t.mu.Lock()
	hasCgroup := t.cg != nil
	t.mu.Unlock()
	r, err := t.shim.Start(ctx, &shim.StartRequest***REMOVED***
		ID: t.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	t.pid = int(r.Pid)
	if !hasCgroup ***REMOVED***
		cg, err := cgroups.Load(cgroups.V1, cgroups.PidPath(t.pid))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		t.mu.Lock()
		t.cg = cg
		t.mu.Unlock()
		if err := t.monitor.Monitor(t); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	t.events.Publish(ctx, runtime.TaskStartEventTopic, &eventstypes.TaskStart***REMOVED***
		ContainerID: t.id,
		Pid:         uint32(t.pid),
	***REMOVED***)
	return nil
***REMOVED***

// State returns runtime information for the task
func (t *Task) State(ctx context.Context) (runtime.State, error) ***REMOVED***
	response, err := t.shim.State(ctx, &shim.StateRequest***REMOVED***
		ID: t.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		if err != grpc.ErrServerStopped ***REMOVED***
			return runtime.State***REMOVED******REMOVED***, errdefs.FromGRPC(err)
		***REMOVED***
		return runtime.State***REMOVED******REMOVED***, errdefs.ErrNotFound
	***REMOVED***
	var status runtime.Status
	switch response.Status ***REMOVED***
	case task.StatusCreated:
		status = runtime.CreatedStatus
	case task.StatusRunning:
		status = runtime.RunningStatus
	case task.StatusStopped:
		status = runtime.StoppedStatus
	case task.StatusPaused:
		status = runtime.PausedStatus
	case task.StatusPausing:
		status = runtime.PausingStatus
	***REMOVED***
	return runtime.State***REMOVED***
		Pid:        response.Pid,
		Status:     status,
		Stdin:      response.Stdin,
		Stdout:     response.Stdout,
		Stderr:     response.Stderr,
		Terminal:   response.Terminal,
		ExitStatus: response.ExitStatus,
		ExitedAt:   response.ExitedAt,
	***REMOVED***, nil
***REMOVED***

// Pause the task and all processes
func (t *Task) Pause(ctx context.Context) error ***REMOVED***
	if _, err := t.shim.Pause(ctx, empty); err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	t.events.Publish(ctx, runtime.TaskPausedEventTopic, &eventstypes.TaskPaused***REMOVED***
		ContainerID: t.id,
	***REMOVED***)
	return nil
***REMOVED***

// Resume the task and all processes
func (t *Task) Resume(ctx context.Context) error ***REMOVED***
	if _, err := t.shim.Resume(ctx, empty); err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	t.events.Publish(ctx, runtime.TaskResumedEventTopic, &eventstypes.TaskResumed***REMOVED***
		ContainerID: t.id,
	***REMOVED***)
	return nil
***REMOVED***

// Kill the task using the provided signal
//
// Optionally send the signal to all processes that are a child of the task
func (t *Task) Kill(ctx context.Context, signal uint32, all bool) error ***REMOVED***
	if _, err := t.shim.Kill(ctx, &shim.KillRequest***REMOVED***
		ID:     t.id,
		Signal: signal,
		All:    all,
	***REMOVED***); err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	return nil
***REMOVED***

// Exec creates a new process inside the task
func (t *Task) Exec(ctx context.Context, id string, opts runtime.ExecOpts) (runtime.Process, error) ***REMOVED***
	if err := identifiers.Validate(id); err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "invalid exec id")
	***REMOVED***
	request := &shim.ExecProcessRequest***REMOVED***
		ID:       id,
		Stdin:    opts.IO.Stdin,
		Stdout:   opts.IO.Stdout,
		Stderr:   opts.IO.Stderr,
		Terminal: opts.IO.Terminal,
		Spec:     opts.Spec,
	***REMOVED***
	if _, err := t.shim.Exec(ctx, request); err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	return &Process***REMOVED***
		id: id,
		t:  t,
	***REMOVED***, nil
***REMOVED***

// Pids returns all system level process ids running inside the task
func (t *Task) Pids(ctx context.Context) ([]runtime.ProcessInfo, error) ***REMOVED***
	resp, err := t.shim.ListPids(ctx, &shim.ListPidsRequest***REMOVED***
		ID: t.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	var processList []runtime.ProcessInfo
	for _, p := range resp.Processes ***REMOVED***
		processList = append(processList, runtime.ProcessInfo***REMOVED***
			Pid:  p.Pid,
			Info: p.Info,
		***REMOVED***)
	***REMOVED***
	return processList, nil
***REMOVED***

// ResizePty changes the side of the task's PTY to the provided width and height
func (t *Task) ResizePty(ctx context.Context, size runtime.ConsoleSize) error ***REMOVED***
	_, err := t.shim.ResizePty(ctx, &shim.ResizePtyRequest***REMOVED***
		ID:     t.id,
		Width:  size.Width,
		Height: size.Height,
	***REMOVED***)
	if err != nil ***REMOVED***
		err = errdefs.FromGRPC(err)
	***REMOVED***
	return err
***REMOVED***

// CloseIO closes the provided IO on the task
func (t *Task) CloseIO(ctx context.Context) error ***REMOVED***
	_, err := t.shim.CloseIO(ctx, &shim.CloseIORequest***REMOVED***
		ID:    t.id,
		Stdin: true,
	***REMOVED***)
	if err != nil ***REMOVED***
		err = errdefs.FromGRPC(err)
	***REMOVED***
	return err
***REMOVED***

// Checkpoint creates a system level dump of the task and process information that can be later restored
func (t *Task) Checkpoint(ctx context.Context, path string, options *types.Any) error ***REMOVED***
	r := &shim.CheckpointTaskRequest***REMOVED***
		Path:    path,
		Options: options,
	***REMOVED***
	if _, err := t.shim.Checkpoint(ctx, r); err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	t.events.Publish(ctx, runtime.TaskCheckpointedEventTopic, &eventstypes.TaskCheckpointed***REMOVED***
		ContainerID: t.id,
	***REMOVED***)
	return nil
***REMOVED***

// DeleteProcess removes the provided process from the task and deletes all on disk state
func (t *Task) DeleteProcess(ctx context.Context, id string) (*runtime.Exit, error) ***REMOVED***
	r, err := t.shim.DeleteProcess(ctx, &shim.DeleteProcessRequest***REMOVED***
		ID: id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	return &runtime.Exit***REMOVED***
		Status:    r.ExitStatus,
		Timestamp: r.ExitedAt,
		Pid:       r.Pid,
	***REMOVED***, nil
***REMOVED***

// Update changes runtime information of a running task
func (t *Task) Update(ctx context.Context, resources *types.Any) error ***REMOVED***
	if _, err := t.shim.Update(ctx, &shim.UpdateTaskRequest***REMOVED***
		Resources: resources,
	***REMOVED***); err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	return nil
***REMOVED***

// Process returns a specific process inside the task by the process id
func (t *Task) Process(ctx context.Context, id string) (runtime.Process, error) ***REMOVED***
	p := &Process***REMOVED***
		id: id,
		t:  t,
	***REMOVED***
	if _, err := p.State(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return p, nil
***REMOVED***

// Metrics returns runtime specific system level metric information for the task
func (t *Task) Metrics(ctx context.Context) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.cg == nil ***REMOVED***
		return nil, errors.Wrap(errdefs.ErrNotFound, "cgroup does not exist")
	***REMOVED***
	stats, err := t.cg.Stat(cgroups.IgnoreNotExist)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return stats, nil
***REMOVED***

// Cgroup returns the underlying cgroup for a linux task
func (t *Task) Cgroup() (cgroups.Cgroup, error) ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.cg == nil ***REMOVED***
		return nil, errors.Wrap(errdefs.ErrNotFound, "cgroup does not exist")
	***REMOVED***
	return t.cg, nil
***REMOVED***

// Wait for the task to exit returning the status and timestamp
func (t *Task) Wait(ctx context.Context) (*runtime.Exit, error) ***REMOVED***
	r, err := t.shim.Wait(ctx, &shim.WaitRequest***REMOVED***
		ID: t.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &runtime.Exit***REMOVED***
		Timestamp: r.ExitedAt,
		Status:    r.ExitStatus,
	***REMOVED***, nil
***REMOVED***
