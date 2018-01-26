// +build !windows

package proc

import (
	"context"
	"sync"
	"syscall"

	"github.com/containerd/console"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/fifo"
	runc "github.com/containerd/go-runc"
	google_protobuf "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
)

type initState interface ***REMOVED***
	State

	Pause(context.Context) error
	Resume(context.Context) error
	Update(context.Context, *google_protobuf.Any) error
	Checkpoint(context.Context, *CheckpointConfig) error
***REMOVED***

type createdState struct ***REMOVED***
	p *Init
***REMOVED***

func (s *createdState) transition(name string) error ***REMOVED***
	switch name ***REMOVED***
	case "running":
		s.p.initState = &runningState***REMOVED***p: s.p***REMOVED***
	case "stopped":
		s.p.initState = &stoppedState***REMOVED***p: s.p***REMOVED***
	case "deleted":
		s.p.initState = &deletedState***REMOVED******REMOVED***
	default:
		return errors.Errorf("invalid state transition %q to %q", stateName(s), name)
	***REMOVED***
	return nil
***REMOVED***

func (s *createdState) Pause(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot pause task in created state")
***REMOVED***

func (s *createdState) Resume(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot resume task in created state")
***REMOVED***

func (s *createdState) Update(context context.Context, r *google_protobuf.Any) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.update(context, r)
***REMOVED***

func (s *createdState) Checkpoint(context context.Context, r *CheckpointConfig) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot checkpoint a task in created state")
***REMOVED***

func (s *createdState) Resize(ws console.WinSize) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.resize(ws)
***REMOVED***

func (s *createdState) Start(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()
	if err := s.p.start(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.transition("running")
***REMOVED***

func (s *createdState) Delete(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()
	if err := s.p.delete(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.transition("deleted")
***REMOVED***

func (s *createdState) Kill(ctx context.Context, sig uint32, all bool) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.kill(ctx, sig, all)
***REMOVED***

func (s *createdState) SetExited(status int) ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	s.p.setExited(status)

	if err := s.transition("stopped"); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

type createdCheckpointState struct ***REMOVED***
	p    *Init
	opts *runc.RestoreOpts
***REMOVED***

func (s *createdCheckpointState) transition(name string) error ***REMOVED***
	switch name ***REMOVED***
	case "running":
		s.p.initState = &runningState***REMOVED***p: s.p***REMOVED***
	case "stopped":
		s.p.initState = &stoppedState***REMOVED***p: s.p***REMOVED***
	case "deleted":
		s.p.initState = &deletedState***REMOVED******REMOVED***
	default:
		return errors.Errorf("invalid state transition %q to %q", stateName(s), name)
	***REMOVED***
	return nil
***REMOVED***

func (s *createdCheckpointState) Pause(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot pause task in created state")
***REMOVED***

func (s *createdCheckpointState) Resume(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot resume task in created state")
***REMOVED***

func (s *createdCheckpointState) Update(context context.Context, r *google_protobuf.Any) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.update(context, r)
***REMOVED***

func (s *createdCheckpointState) Checkpoint(context context.Context, r *CheckpointConfig) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot checkpoint a task in created state")
***REMOVED***

func (s *createdCheckpointState) Resize(ws console.WinSize) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.resize(ws)
***REMOVED***

func (s *createdCheckpointState) Start(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()
	p := s.p
	if _, err := s.p.runtime.Restore(ctx, p.id, p.bundle, s.opts); err != nil ***REMOVED***
		return p.runtimeError(err, "OCI runtime restore failed")
	***REMOVED***
	sio := p.stdio
	if sio.Stdin != "" ***REMOVED***
		sc, err := fifo.OpenFifo(ctx, sio.Stdin, syscall.O_WRONLY|syscall.O_NONBLOCK, 0)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to open stdin fifo %s", sio.Stdin)
		***REMOVED***
		p.stdin = sc
		p.closers = append(p.closers, sc)
	***REMOVED***
	var copyWaitGroup sync.WaitGroup
	if !sio.IsNull() ***REMOVED***
		if err := copyPipes(ctx, p.io, sio.Stdin, sio.Stdout, sio.Stderr, &p.wg, &copyWaitGroup); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to start io pipe copy")
		***REMOVED***
	***REMOVED***

	copyWaitGroup.Wait()
	pid, err := runc.ReadPidFile(s.opts.PidFile)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to retrieve OCI runtime container pid")
	***REMOVED***
	p.pid = pid

	return s.transition("running")
***REMOVED***

func (s *createdCheckpointState) Delete(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()
	if err := s.p.delete(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.transition("deleted")
***REMOVED***

func (s *createdCheckpointState) Kill(ctx context.Context, sig uint32, all bool) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.kill(ctx, sig, all)
***REMOVED***

func (s *createdCheckpointState) SetExited(status int) ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	s.p.setExited(status)

	if err := s.transition("stopped"); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

type runningState struct ***REMOVED***
	p *Init
***REMOVED***

func (s *runningState) transition(name string) error ***REMOVED***
	switch name ***REMOVED***
	case "stopped":
		s.p.initState = &stoppedState***REMOVED***p: s.p***REMOVED***
	case "paused":
		s.p.initState = &pausedState***REMOVED***p: s.p***REMOVED***
	default:
		return errors.Errorf("invalid state transition %q to %q", stateName(s), name)
	***REMOVED***
	return nil
***REMOVED***

func (s *runningState) Pause(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()
	if err := s.p.pause(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.transition("paused")
***REMOVED***

func (s *runningState) Resume(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot resume a running process")
***REMOVED***

func (s *runningState) Update(context context.Context, r *google_protobuf.Any) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.update(context, r)
***REMOVED***

func (s *runningState) Checkpoint(ctx context.Context, r *CheckpointConfig) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.checkpoint(ctx, r)
***REMOVED***

func (s *runningState) Resize(ws console.WinSize) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.resize(ws)
***REMOVED***

func (s *runningState) Start(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot start a running process")
***REMOVED***

func (s *runningState) Delete(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot delete a running process")
***REMOVED***

func (s *runningState) Kill(ctx context.Context, sig uint32, all bool) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.kill(ctx, sig, all)
***REMOVED***

func (s *runningState) SetExited(status int) ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	s.p.setExited(status)

	if err := s.transition("stopped"); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

type pausedState struct ***REMOVED***
	p *Init
***REMOVED***

func (s *pausedState) transition(name string) error ***REMOVED***
	switch name ***REMOVED***
	case "running":
		s.p.initState = &runningState***REMOVED***p: s.p***REMOVED***
	case "stopped":
		s.p.initState = &stoppedState***REMOVED***p: s.p***REMOVED***
	default:
		return errors.Errorf("invalid state transition %q to %q", stateName(s), name)
	***REMOVED***
	return nil
***REMOVED***

func (s *pausedState) Pause(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot pause a paused container")
***REMOVED***

func (s *pausedState) Resume(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	if err := s.p.resume(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.transition("running")
***REMOVED***

func (s *pausedState) Update(context context.Context, r *google_protobuf.Any) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.update(context, r)
***REMOVED***

func (s *pausedState) Checkpoint(ctx context.Context, r *CheckpointConfig) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.checkpoint(ctx, r)
***REMOVED***

func (s *pausedState) Resize(ws console.WinSize) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.resize(ws)
***REMOVED***

func (s *pausedState) Start(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot start a paused process")
***REMOVED***

func (s *pausedState) Delete(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot delete a paused process")
***REMOVED***

func (s *pausedState) Kill(ctx context.Context, sig uint32, all bool) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.kill(ctx, sig, all)
***REMOVED***

func (s *pausedState) SetExited(status int) ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	s.p.setExited(status)

	if err := s.transition("stopped"); err != nil ***REMOVED***
		panic(err)
	***REMOVED***

***REMOVED***

type stoppedState struct ***REMOVED***
	p *Init
***REMOVED***

func (s *stoppedState) transition(name string) error ***REMOVED***
	switch name ***REMOVED***
	case "deleted":
		s.p.initState = &deletedState***REMOVED******REMOVED***
	default:
		return errors.Errorf("invalid state transition %q to %q", stateName(s), name)
	***REMOVED***
	return nil
***REMOVED***

func (s *stoppedState) Pause(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot pause a stopped container")
***REMOVED***

func (s *stoppedState) Resume(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot resume a stopped container")
***REMOVED***

func (s *stoppedState) Update(context context.Context, r *google_protobuf.Any) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot update a stopped container")
***REMOVED***

func (s *stoppedState) Checkpoint(ctx context.Context, r *CheckpointConfig) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot checkpoint a stopped container")
***REMOVED***

func (s *stoppedState) Resize(ws console.WinSize) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot resize a stopped container")
***REMOVED***

func (s *stoppedState) Start(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot start a stopped process")
***REMOVED***

func (s *stoppedState) Delete(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()
	if err := s.p.delete(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.transition("deleted")
***REMOVED***

func (s *stoppedState) Kill(ctx context.Context, sig uint32, all bool) error ***REMOVED***
	return errdefs.ToGRPCf(errdefs.ErrNotFound, "process %s not found", s.p.id)
***REMOVED***

func (s *stoppedState) SetExited(status int) ***REMOVED***
	// no op
***REMOVED***
