// +build !windows

package proc

import (
	"context"

	"github.com/containerd/console"
	"github.com/pkg/errors"
)

type execCreatedState struct ***REMOVED***
	p *execProcess
***REMOVED***

func (s *execCreatedState) transition(name string) error ***REMOVED***
	switch name ***REMOVED***
	case "running":
		s.p.State = &execRunningState***REMOVED***p: s.p***REMOVED***
	case "stopped":
		s.p.State = &execStoppedState***REMOVED***p: s.p***REMOVED***
	case "deleted":
		s.p.State = &deletedState***REMOVED******REMOVED***
	default:
		return errors.Errorf("invalid state transition %q to %q", stateName(s), name)
	***REMOVED***
	return nil
***REMOVED***

func (s *execCreatedState) Resize(ws console.WinSize) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.resize(ws)
***REMOVED***

func (s *execCreatedState) Start(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()
	if err := s.p.start(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.transition("running")
***REMOVED***

func (s *execCreatedState) Delete(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()
	if err := s.p.delete(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.transition("deleted")
***REMOVED***

func (s *execCreatedState) Kill(ctx context.Context, sig uint32, all bool) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.kill(ctx, sig, all)
***REMOVED***

func (s *execCreatedState) SetExited(status int) ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	s.p.setExited(status)

	if err := s.transition("stopped"); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

type execRunningState struct ***REMOVED***
	p *execProcess
***REMOVED***

func (s *execRunningState) transition(name string) error ***REMOVED***
	switch name ***REMOVED***
	case "stopped":
		s.p.State = &execStoppedState***REMOVED***p: s.p***REMOVED***
	default:
		return errors.Errorf("invalid state transition %q to %q", stateName(s), name)
	***REMOVED***
	return nil
***REMOVED***

func (s *execRunningState) Resize(ws console.WinSize) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.resize(ws)
***REMOVED***

func (s *execRunningState) Start(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot start a running process")
***REMOVED***

func (s *execRunningState) Delete(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot delete a running process")
***REMOVED***

func (s *execRunningState) Kill(ctx context.Context, sig uint32, all bool) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.kill(ctx, sig, all)
***REMOVED***

func (s *execRunningState) SetExited(status int) ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	s.p.setExited(status)

	if err := s.transition("stopped"); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

type execStoppedState struct ***REMOVED***
	p *execProcess
***REMOVED***

func (s *execStoppedState) transition(name string) error ***REMOVED***
	switch name ***REMOVED***
	case "deleted":
		s.p.State = &deletedState***REMOVED******REMOVED***
	default:
		return errors.Errorf("invalid state transition %q to %q", stateName(s), name)
	***REMOVED***
	return nil
***REMOVED***

func (s *execStoppedState) Resize(ws console.WinSize) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot resize a stopped container")
***REMOVED***

func (s *execStoppedState) Start(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return errors.Errorf("cannot start a stopped process")
***REMOVED***

func (s *execStoppedState) Delete(ctx context.Context) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()
	if err := s.p.delete(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.transition("deleted")
***REMOVED***

func (s *execStoppedState) Kill(ctx context.Context, sig uint32, all bool) error ***REMOVED***
	s.p.mu.Lock()
	defer s.p.mu.Unlock()

	return s.p.kill(ctx, sig, all)
***REMOVED***

func (s *execStoppedState) SetExited(status int) ***REMOVED***
	// no op
***REMOVED***
