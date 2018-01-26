// +build !windows

package shim

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/containerd/console"
	eventstypes "github.com/containerd/containerd/api/events"
	"github.com/containerd/containerd/api/types/task"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/events"
	"github.com/containerd/containerd/linux/proc"
	"github.com/containerd/containerd/linux/runctypes"
	shimapi "github.com/containerd/containerd/linux/shim/v1"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/reaper"
	"github.com/containerd/containerd/runtime"
	runc "github.com/containerd/go-runc"
	"github.com/containerd/typeurl"
	ptypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var empty = &ptypes.Empty***REMOVED******REMOVED***

// Config contains shim specific configuration
type Config struct ***REMOVED***
	Path          string
	Namespace     string
	WorkDir       string
	Criu          string
	RuntimeRoot   string
	SystemdCgroup bool
***REMOVED***

// NewService returns a new shim service that can be used via GRPC
func NewService(config Config, publisher events.Publisher) (*Service, error) ***REMOVED***
	if config.Namespace == "" ***REMOVED***
		return nil, fmt.Errorf("shim namespace cannot be empty")
	***REMOVED***
	ctx := namespaces.WithNamespace(context.Background(), config.Namespace)
	ctx = log.WithLogger(ctx, logrus.WithFields(logrus.Fields***REMOVED***
		"namespace": config.Namespace,
		"path":      config.Path,
		"pid":       os.Getpid(),
	***REMOVED***))
	s := &Service***REMOVED***
		config:    config,
		context:   ctx,
		processes: make(map[string]proc.Process),
		events:    make(chan interface***REMOVED******REMOVED***, 128),
		ec:        reaper.Default.Subscribe(),
	***REMOVED***
	go s.processExits()
	if err := s.initPlatform(); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to initialized platform behavior")
	***REMOVED***
	go s.forward(publisher)
	return s, nil
***REMOVED***

// Service is the shim implementation of a remote shim over GRPC
type Service struct ***REMOVED***
	mu sync.Mutex

	config    Config
	context   context.Context
	processes map[string]proc.Process
	events    chan interface***REMOVED******REMOVED***
	platform  proc.Platform
	ec        chan runc.Exit

	// Filled by Create()
	id     string
	bundle string
***REMOVED***

// Create a new initial process and container with the underlying OCI runtime
func (s *Service) Create(ctx context.Context, r *shimapi.CreateTaskRequest) (*shimapi.CreateTaskResponse, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	process, err := proc.New(
		ctx,
		s.config.Path,
		s.config.WorkDir,
		s.config.RuntimeRoot,
		s.config.Namespace,
		s.config.Criu,
		s.config.SystemdCgroup,
		s.platform,
		&proc.CreateConfig***REMOVED***
			ID:               r.ID,
			Bundle:           r.Bundle,
			Runtime:          r.Runtime,
			Rootfs:           r.Rootfs,
			Terminal:         r.Terminal,
			Stdin:            r.Stdin,
			Stdout:           r.Stdout,
			Stderr:           r.Stderr,
			Checkpoint:       r.Checkpoint,
			ParentCheckpoint: r.ParentCheckpoint,
			Options:          r.Options,
		***REMOVED***,
	)
	if err != nil ***REMOVED***
		return nil, errdefs.ToGRPC(err)
	***REMOVED***
	// save the main task id and bundle to the shim for additional requests
	s.id = r.ID
	s.bundle = r.Bundle
	pid := process.Pid()
	s.processes[r.ID] = process
	return &shimapi.CreateTaskResponse***REMOVED***
		Pid: uint32(pid),
	***REMOVED***, nil
***REMOVED***

// Start a process
func (s *Service) Start(ctx context.Context, r *shimapi.StartRequest) (*shimapi.StartResponse, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.processes[r.ID]
	if p == nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrNotFound, "process %s not found", r.ID)
	***REMOVED***
	if err := p.Start(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &shimapi.StartResponse***REMOVED***
		ID:  p.ID(),
		Pid: uint32(p.Pid()),
	***REMOVED***, nil
***REMOVED***

// Delete the initial process and container
func (s *Service) Delete(ctx context.Context, r *ptypes.Empty) (*shimapi.DeleteResponse, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.processes[s.id]
	if p == nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrFailedPrecondition, "container must be created")
	***REMOVED***
	if err := p.Delete(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	delete(s.processes, s.id)
	s.platform.Close()
	return &shimapi.DeleteResponse***REMOVED***
		ExitStatus: uint32(p.ExitStatus()),
		ExitedAt:   p.ExitedAt(),
		Pid:        uint32(p.Pid()),
	***REMOVED***, nil
***REMOVED***

// DeleteProcess deletes an exec'd process
func (s *Service) DeleteProcess(ctx context.Context, r *shimapi.DeleteProcessRequest) (*shimapi.DeleteResponse, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if r.ID == s.id ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "cannot delete init process with DeleteProcess")
	***REMOVED***
	p := s.processes[r.ID]
	if p == nil ***REMOVED***
		return nil, errors.Wrapf(errdefs.ErrNotFound, "process %s", r.ID)
	***REMOVED***
	if err := p.Delete(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	delete(s.processes, r.ID)
	return &shimapi.DeleteResponse***REMOVED***
		ExitStatus: uint32(p.ExitStatus()),
		ExitedAt:   p.ExitedAt(),
		Pid:        uint32(p.Pid()),
	***REMOVED***, nil
***REMOVED***

// Exec an additional process inside the container
func (s *Service) Exec(ctx context.Context, r *shimapi.ExecProcessRequest) (*ptypes.Empty, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	if p := s.processes[r.ID]; p != nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrAlreadyExists, "id %s", r.ID)
	***REMOVED***

	p := s.processes[s.id]
	if p == nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrFailedPrecondition, "container must be created")
	***REMOVED***

	process, err := p.(*proc.Init).Exec(ctx, s.config.Path, &proc.ExecConfig***REMOVED***
		ID:       r.ID,
		Terminal: r.Terminal,
		Stdin:    r.Stdin,
		Stdout:   r.Stdout,
		Stderr:   r.Stderr,
		Spec:     r.Spec,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.ToGRPC(err)
	***REMOVED***
	s.processes[r.ID] = process
	return empty, nil
***REMOVED***

// ResizePty of a process
func (s *Service) ResizePty(ctx context.Context, r *shimapi.ResizePtyRequest) (*ptypes.Empty, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if r.ID == "" ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrInvalidArgument, "id not provided")
	***REMOVED***
	ws := console.WinSize***REMOVED***
		Width:  uint16(r.Width),
		Height: uint16(r.Height),
	***REMOVED***
	p := s.processes[r.ID]
	if p == nil ***REMOVED***
		return nil, errors.Errorf("process does not exist %s", r.ID)
	***REMOVED***
	if err := p.Resize(ws); err != nil ***REMOVED***
		return nil, errdefs.ToGRPC(err)
	***REMOVED***
	return empty, nil
***REMOVED***

// State returns runtime state information for a process
func (s *Service) State(ctx context.Context, r *shimapi.StateRequest) (*shimapi.StateResponse, error) ***REMOVED***
	s.mu.Lock()
	p := s.processes[r.ID]
	s.mu.Unlock()
	if p == nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrNotFound, "process id %s not found", r.ID)
	***REMOVED***
	st, err := p.Status(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	status := task.StatusUnknown
	switch st ***REMOVED***
	case "created":
		status = task.StatusCreated
	case "running":
		status = task.StatusRunning
	case "stopped":
		status = task.StatusStopped
	case "paused":
		status = task.StatusPaused
	case "pausing":
		status = task.StatusPausing
	***REMOVED***
	sio := p.Stdio()
	return &shimapi.StateResponse***REMOVED***
		ID:         p.ID(),
		Bundle:     s.bundle,
		Pid:        uint32(p.Pid()),
		Status:     status,
		Stdin:      sio.Stdin,
		Stdout:     sio.Stdout,
		Stderr:     sio.Stderr,
		Terminal:   sio.Terminal,
		ExitStatus: uint32(p.ExitStatus()),
		ExitedAt:   p.ExitedAt(),
	***REMOVED***, nil
***REMOVED***

// Pause the container
func (s *Service) Pause(ctx context.Context, r *ptypes.Empty) (*ptypes.Empty, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.processes[s.id]
	if p == nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrFailedPrecondition, "container must be created")
	***REMOVED***
	if err := p.(*proc.Init).Pause(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return empty, nil
***REMOVED***

// Resume the container
func (s *Service) Resume(ctx context.Context, r *ptypes.Empty) (*ptypes.Empty, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.processes[s.id]
	if p == nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrFailedPrecondition, "container must be created")
	***REMOVED***
	if err := p.(*proc.Init).Resume(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return empty, nil
***REMOVED***

// Kill a process with the provided signal
func (s *Service) Kill(ctx context.Context, r *shimapi.KillRequest) (*ptypes.Empty, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if r.ID == "" ***REMOVED***
		p := s.processes[s.id]
		if p == nil ***REMOVED***
			return nil, errdefs.ToGRPCf(errdefs.ErrFailedPrecondition, "container must be created")
		***REMOVED***
		if err := p.Kill(ctx, r.Signal, r.All); err != nil ***REMOVED***
			return nil, errdefs.ToGRPC(err)
		***REMOVED***
		return empty, nil
	***REMOVED***

	p := s.processes[r.ID]
	if p == nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrNotFound, "process id %s not found", r.ID)
	***REMOVED***
	if err := p.Kill(ctx, r.Signal, r.All); err != nil ***REMOVED***
		return nil, errdefs.ToGRPC(err)
	***REMOVED***
	return empty, nil
***REMOVED***

// ListPids returns all pids inside the container
func (s *Service) ListPids(ctx context.Context, r *shimapi.ListPidsRequest) (*shimapi.ListPidsResponse, error) ***REMOVED***
	pids, err := s.getContainerPids(ctx, r.ID)
	if err != nil ***REMOVED***
		return nil, errdefs.ToGRPC(err)
	***REMOVED***
	var processes []*task.ProcessInfo
	for _, pid := range pids ***REMOVED***
		pInfo := task.ProcessInfo***REMOVED***
			Pid: pid,
		***REMOVED***
		for _, p := range s.processes ***REMOVED***
			if p.Pid() == int(pid) ***REMOVED***
				d := &runctypes.ProcessDetails***REMOVED***
					ExecID: p.ID(),
				***REMOVED***
				a, err := typeurl.MarshalAny(d)
				if err != nil ***REMOVED***
					return nil, errors.Wrapf(err, "failed to marshal process %d info", pid)
				***REMOVED***
				pInfo.Info = a
				break
			***REMOVED***
		***REMOVED***
		processes = append(processes, &pInfo)
	***REMOVED***
	return &shimapi.ListPidsResponse***REMOVED***
		Processes: processes,
	***REMOVED***, nil
***REMOVED***

// CloseIO of a process
func (s *Service) CloseIO(ctx context.Context, r *shimapi.CloseIORequest) (*ptypes.Empty, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.processes[r.ID]
	if p == nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrNotFound, "process does not exist %s", r.ID)
	***REMOVED***
	if stdin := p.Stdin(); stdin != nil ***REMOVED***
		if err := stdin.Close(); err != nil ***REMOVED***
			return nil, errors.Wrap(err, "close stdin")
		***REMOVED***
	***REMOVED***
	return empty, nil
***REMOVED***

// Checkpoint the container
func (s *Service) Checkpoint(ctx context.Context, r *shimapi.CheckpointTaskRequest) (*ptypes.Empty, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.processes[s.id]
	if p == nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrFailedPrecondition, "container must be created")
	***REMOVED***
	if err := p.(*proc.Init).Checkpoint(ctx, &proc.CheckpointConfig***REMOVED***
		Path:    r.Path,
		Options: r.Options,
	***REMOVED***); err != nil ***REMOVED***
		return nil, errdefs.ToGRPC(err)
	***REMOVED***
	return empty, nil
***REMOVED***

// ShimInfo returns shim information such as the shim's pid
func (s *Service) ShimInfo(ctx context.Context, r *ptypes.Empty) (*shimapi.ShimInfoResponse, error) ***REMOVED***
	return &shimapi.ShimInfoResponse***REMOVED***
		ShimPid: uint32(os.Getpid()),
	***REMOVED***, nil
***REMOVED***

// Update a running container
func (s *Service) Update(ctx context.Context, r *shimapi.UpdateTaskRequest) (*ptypes.Empty, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.processes[s.id]
	if p == nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrFailedPrecondition, "container must be created")
	***REMOVED***
	if err := p.(*proc.Init).Update(ctx, r.Resources); err != nil ***REMOVED***
		return nil, errdefs.ToGRPC(err)
	***REMOVED***
	return empty, nil
***REMOVED***

// Wait for a process to exit
func (s *Service) Wait(ctx context.Context, r *shimapi.WaitRequest) (*shimapi.WaitResponse, error) ***REMOVED***
	s.mu.Lock()
	p := s.processes[r.ID]
	s.mu.Unlock()
	if p == nil ***REMOVED***
		return nil, errdefs.ToGRPCf(errdefs.ErrFailedPrecondition, "container must be created")
	***REMOVED***
	p.Wait()

	return &shimapi.WaitResponse***REMOVED***
		ExitStatus: uint32(p.ExitStatus()),
		ExitedAt:   p.ExitedAt(),
	***REMOVED***, nil
***REMOVED***

func (s *Service) processExits() ***REMOVED***
	for e := range s.ec ***REMOVED***
		s.checkProcesses(e)
	***REMOVED***
***REMOVED***

func (s *Service) checkProcesses(e runc.Exit) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.processes ***REMOVED***
		if p.Pid() == e.Pid ***REMOVED***
			if ip, ok := p.(*proc.Init); ok ***REMOVED***
				// Ensure all children are killed
				if err := ip.KillAll(s.context); err != nil ***REMOVED***
					log.G(s.context).WithError(err).WithField("id", ip.ID()).
						Error("failed to kill init's children")
				***REMOVED***
			***REMOVED***
			p.SetExited(e.Status)
			s.events <- &eventstypes.TaskExit***REMOVED***
				ContainerID: s.id,
				ID:          p.ID(),
				Pid:         uint32(e.Pid),
				ExitStatus:  uint32(e.Status),
				ExitedAt:    p.ExitedAt(),
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *Service) getContainerPids(ctx context.Context, id string) ([]uint32, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.processes[s.id]
	if p == nil ***REMOVED***
		return nil, errors.Wrapf(errdefs.ErrFailedPrecondition, "container must be created")
	***REMOVED***

	ps, err := p.(*proc.Init).Runtime().Ps(ctx, id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	pids := make([]uint32, 0, len(ps))
	for _, pid := range ps ***REMOVED***
		pids = append(pids, uint32(pid))
	***REMOVED***
	return pids, nil
***REMOVED***

func (s *Service) forward(publisher events.Publisher) ***REMOVED***
	for e := range s.events ***REMOVED***
		if err := publisher.Publish(s.context, getTopic(s.context, e), e); err != nil ***REMOVED***
			log.G(s.context).WithError(err).Error("post event")
		***REMOVED***
	***REMOVED***
***REMOVED***

func getTopic(ctx context.Context, e interface***REMOVED******REMOVED***) string ***REMOVED***
	switch e.(type) ***REMOVED***
	case *eventstypes.TaskCreate:
		return runtime.TaskCreateEventTopic
	case *eventstypes.TaskStart:
		return runtime.TaskStartEventTopic
	case *eventstypes.TaskOOM:
		return runtime.TaskOOMEventTopic
	case *eventstypes.TaskExit:
		return runtime.TaskExitEventTopic
	case *eventstypes.TaskDelete:
		return runtime.TaskDeleteEventTopic
	case *eventstypes.TaskExecAdded:
		return runtime.TaskExecAddedEventTopic
	case *eventstypes.TaskExecStarted:
		return runtime.TaskExecStartedEventTopic
	case *eventstypes.TaskPaused:
		return runtime.TaskPausedEventTopic
	case *eventstypes.TaskResumed:
		return runtime.TaskResumedEventTopic
	case *eventstypes.TaskCheckpointed:
		return runtime.TaskCheckpointedEventTopic
	default:
		logrus.Warnf("no topic for type %#v", e)
	***REMOVED***
	return runtime.TaskUnknownTopic
***REMOVED***
