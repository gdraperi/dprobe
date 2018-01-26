package containerd

import (
	"context"
	"strings"
	"syscall"
	"time"

	"github.com/containerd/containerd/api/services/tasks/v1"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/errdefs"
	"github.com/pkg/errors"
)

// Process represents a system process
type Process interface ***REMOVED***
	// Pid is the system specific process id
	Pid() uint32
	// Start starts the process executing the user's defined binary
	Start(context.Context) error
	// Delete removes the process and any resources allocated returning the exit status
	Delete(context.Context, ...ProcessDeleteOpts) (*ExitStatus, error)
	// Kill sends the provided signal to the process
	Kill(context.Context, syscall.Signal, ...KillOpts) error
	// Wait asynchronously waits for the process to exit, and sends the exit code to the returned channel
	Wait(context.Context) (<-chan ExitStatus, error)
	// CloseIO allows various pipes to be closed on the process
	CloseIO(context.Context, ...IOCloserOpts) error
	// Resize changes the width and heigh of the process's terminal
	Resize(ctx context.Context, w, h uint32) error
	// IO returns the io set for the process
	IO() cio.IO
	// Status returns the executing status of the process
	Status(context.Context) (Status, error)
***REMOVED***

// ExitStatus encapsulates a process' exit status.
// It is used by `Wait()` to return either a process exit code or an error
type ExitStatus struct ***REMOVED***
	code     uint32
	exitedAt time.Time
	err      error
***REMOVED***

// Result returns the exit code and time of the exit status.
// An error may be returned here to which indicates there was an error
//   at some point while waiting for the exit status. It does not signify
//   an error with the process itself.
// If an error is returned, the process may still be running.
func (s ExitStatus) Result() (uint32, time.Time, error) ***REMOVED***
	return s.code, s.exitedAt, s.err
***REMOVED***

// ExitCode returns the exit code of the process.
// This is only valid is Error() returns nil
func (s ExitStatus) ExitCode() uint32 ***REMOVED***
	return s.code
***REMOVED***

// ExitTime returns the exit time of the process
// This is only valid is Error() returns nil
func (s ExitStatus) ExitTime() time.Time ***REMOVED***
	return s.exitedAt
***REMOVED***

// Error returns the error, if any, that occured while waiting for the
// process.
func (s ExitStatus) Error() error ***REMOVED***
	return s.err
***REMOVED***

type process struct ***REMOVED***
	id   string
	task *task
	pid  uint32
	io   cio.IO
***REMOVED***

func (p *process) ID() string ***REMOVED***
	return p.id
***REMOVED***

// Pid returns the pid of the process
// The pid is not set until start is called and returns
func (p *process) Pid() uint32 ***REMOVED***
	return p.pid
***REMOVED***

// Start starts the exec process
func (p *process) Start(ctx context.Context) error ***REMOVED***
	r, err := p.task.client.TaskService().Start(ctx, &tasks.StartRequest***REMOVED***
		ContainerID: p.task.id,
		ExecID:      p.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		p.io.Cancel()
		p.io.Wait()
		p.io.Close()
		return errdefs.FromGRPC(err)
	***REMOVED***
	p.pid = r.Pid
	return nil
***REMOVED***

func (p *process) Kill(ctx context.Context, s syscall.Signal, opts ...KillOpts) error ***REMOVED***
	var i KillInfo
	for _, o := range opts ***REMOVED***
		if err := o(ctx, &i); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	_, err := p.task.client.TaskService().Kill(ctx, &tasks.KillRequest***REMOVED***
		Signal:      uint32(s),
		ContainerID: p.task.id,
		ExecID:      p.id,
		All:         i.All,
	***REMOVED***)
	return errdefs.FromGRPC(err)
***REMOVED***

func (p *process) Wait(ctx context.Context) (<-chan ExitStatus, error) ***REMOVED***
	c := make(chan ExitStatus, 1)
	go func() ***REMOVED***
		defer close(c)
		r, err := p.task.client.TaskService().Wait(ctx, &tasks.WaitRequest***REMOVED***
			ContainerID: p.task.id,
			ExecID:      p.id,
		***REMOVED***)
		if err != nil ***REMOVED***
			c <- ExitStatus***REMOVED***
				code: UnknownExitStatus,
				err:  err,
			***REMOVED***
			return
		***REMOVED***
		c <- ExitStatus***REMOVED***
			code:     r.ExitStatus,
			exitedAt: r.ExitedAt,
		***REMOVED***
	***REMOVED***()
	return c, nil
***REMOVED***

func (p *process) CloseIO(ctx context.Context, opts ...IOCloserOpts) error ***REMOVED***
	r := &tasks.CloseIORequest***REMOVED***
		ContainerID: p.task.id,
		ExecID:      p.id,
	***REMOVED***
	var i IOCloseInfo
	for _, o := range opts ***REMOVED***
		o(&i)
	***REMOVED***
	r.Stdin = i.Stdin
	_, err := p.task.client.TaskService().CloseIO(ctx, r)
	return errdefs.FromGRPC(err)
***REMOVED***

func (p *process) IO() cio.IO ***REMOVED***
	return p.io
***REMOVED***

func (p *process) Resize(ctx context.Context, w, h uint32) error ***REMOVED***
	_, err := p.task.client.TaskService().ResizePty(ctx, &tasks.ResizePtyRequest***REMOVED***
		ContainerID: p.task.id,
		Width:       w,
		Height:      h,
		ExecID:      p.id,
	***REMOVED***)
	return errdefs.FromGRPC(err)
***REMOVED***

func (p *process) Delete(ctx context.Context, opts ...ProcessDeleteOpts) (*ExitStatus, error) ***REMOVED***
	for _, o := range opts ***REMOVED***
		if err := o(ctx, p); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	status, err := p.Status(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	switch status.Status ***REMOVED***
	case Running, Paused, Pausing:
		return nil, errors.Wrapf(errdefs.ErrFailedPrecondition, "process must be stopped before deletion")
	***REMOVED***
	r, err := p.task.client.TaskService().DeleteProcess(ctx, &tasks.DeleteProcessRequest***REMOVED***
		ContainerID: p.task.id,
		ExecID:      p.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	if p.io != nil ***REMOVED***
		p.io.Cancel()
		p.io.Wait()
		p.io.Close()
	***REMOVED***
	return &ExitStatus***REMOVED***code: r.ExitStatus, exitedAt: r.ExitedAt***REMOVED***, nil
***REMOVED***

func (p *process) Status(ctx context.Context) (Status, error) ***REMOVED***
	r, err := p.task.client.TaskService().Get(ctx, &tasks.GetRequest***REMOVED***
		ContainerID: p.task.id,
		ExecID:      p.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return Status***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***
	return Status***REMOVED***
		Status:     ProcessStatus(strings.ToLower(r.Process.Status.String())),
		ExitStatus: r.Process.ExitStatus,
	***REMOVED***, nil
***REMOVED***
