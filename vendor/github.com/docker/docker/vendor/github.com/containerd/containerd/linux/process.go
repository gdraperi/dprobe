// +build linux

package linux

import (
	"context"

	eventstypes "github.com/containerd/containerd/api/events"
	"github.com/containerd/containerd/api/types/task"
	"github.com/containerd/containerd/errdefs"
	shim "github.com/containerd/containerd/linux/shim/v1"
	"github.com/containerd/containerd/runtime"
)

// Process implements a linux process
type Process struct ***REMOVED***
	id string
	t  *Task
***REMOVED***

// ID of the process
func (p *Process) ID() string ***REMOVED***
	return p.id
***REMOVED***

// Kill sends the provided signal to the underlying process
//
// Unable to kill all processes in the task using this method on a process
func (p *Process) Kill(ctx context.Context, signal uint32, _ bool) error ***REMOVED***
	_, err := p.t.shim.Kill(ctx, &shim.KillRequest***REMOVED***
		Signal: signal,
		ID:     p.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	return err
***REMOVED***

// State of process
func (p *Process) State(ctx context.Context) (runtime.State, error) ***REMOVED***
	// use the container status for the status of the process
	response, err := p.t.shim.State(ctx, &shim.StateRequest***REMOVED***
		ID: p.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return runtime.State***REMOVED******REMOVED***, errdefs.FromGRPC(err)
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
	***REMOVED***, nil
***REMOVED***

// ResizePty changes the side of the process's PTY to the provided width and height
func (p *Process) ResizePty(ctx context.Context, size runtime.ConsoleSize) error ***REMOVED***
	_, err := p.t.shim.ResizePty(ctx, &shim.ResizePtyRequest***REMOVED***
		ID:     p.id,
		Width:  size.Width,
		Height: size.Height,
	***REMOVED***)
	if err != nil ***REMOVED***
		err = errdefs.FromGRPC(err)
	***REMOVED***
	return err
***REMOVED***

// CloseIO closes the provided IO pipe for the process
func (p *Process) CloseIO(ctx context.Context) error ***REMOVED***
	_, err := p.t.shim.CloseIO(ctx, &shim.CloseIORequest***REMOVED***
		ID:    p.id,
		Stdin: true,
	***REMOVED***)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	return nil
***REMOVED***

// Start the process
func (p *Process) Start(ctx context.Context) error ***REMOVED***
	r, err := p.t.shim.Start(ctx, &shim.StartRequest***REMOVED***
		ID: p.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	p.t.events.Publish(ctx, runtime.TaskExecStartedEventTopic, &eventstypes.TaskExecStarted***REMOVED***
		ContainerID: p.t.id,
		Pid:         r.Pid,
		ExecID:      p.id,
	***REMOVED***)
	return nil
***REMOVED***

// Wait on the process to exit and return the exit status and timestamp
func (p *Process) Wait(ctx context.Context) (*runtime.Exit, error) ***REMOVED***
	r, err := p.t.shim.Wait(ctx, &shim.WaitRequest***REMOVED***
		ID: p.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &runtime.Exit***REMOVED***
		Timestamp: r.ExitedAt,
		Status:    r.ExitStatus,
	***REMOVED***, nil
***REMOVED***
