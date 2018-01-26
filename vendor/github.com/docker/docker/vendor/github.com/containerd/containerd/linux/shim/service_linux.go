package shim

import (
	"context"
	"io"
	"sync"
	"syscall"

	"github.com/containerd/console"
	"github.com/containerd/fifo"
	"github.com/pkg/errors"
)

type linuxPlatform struct ***REMOVED***
	epoller *console.Epoller
***REMOVED***

func (p *linuxPlatform) CopyConsole(ctx context.Context, console console.Console, stdin, stdout, stderr string, wg, cwg *sync.WaitGroup) (console.Console, error) ***REMOVED***
	if p.epoller == nil ***REMOVED***
		return nil, errors.New("uninitialized epoller")
	***REMOVED***

	epollConsole, err := p.epoller.Add(console)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if stdin != "" ***REMOVED***
		in, err := fifo.OpenFifo(ctx, stdin, syscall.O_RDONLY, 0)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cwg.Add(1)
		go func() ***REMOVED***
			cwg.Done()
			io.Copy(epollConsole, in)
		***REMOVED***()
	***REMOVED***

	outw, err := fifo.OpenFifo(ctx, stdout, syscall.O_WRONLY, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	outr, err := fifo.OpenFifo(ctx, stdout, syscall.O_RDONLY, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	wg.Add(1)
	cwg.Add(1)
	go func() ***REMOVED***
		cwg.Done()
		io.Copy(outw, epollConsole)
		epollConsole.Close()
		outr.Close()
		outw.Close()
		wg.Done()
	***REMOVED***()
	return epollConsole, nil
***REMOVED***

func (p *linuxPlatform) ShutdownConsole(ctx context.Context, cons console.Console) error ***REMOVED***
	if p.epoller == nil ***REMOVED***
		return errors.New("uninitialized epoller")
	***REMOVED***
	epollConsole, ok := cons.(*console.EpollConsole)
	if !ok ***REMOVED***
		return errors.Errorf("expected EpollConsole, got %#v", cons)
	***REMOVED***
	return epollConsole.Shutdown(p.epoller.CloseConsole)
***REMOVED***

func (p *linuxPlatform) Close() error ***REMOVED***
	return p.epoller.Close()
***REMOVED***

// initialize a single epoll fd to manage our consoles. `initPlatform` should
// only be called once.
func (s *Service) initPlatform() error ***REMOVED***
	if s.platform != nil ***REMOVED***
		return nil
	***REMOVED***
	epoller, err := console.NewEpoller()
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to initialize epoller")
	***REMOVED***
	s.platform = &linuxPlatform***REMOVED***
		epoller: epoller,
	***REMOVED***
	go epoller.Wait()
	return nil
***REMOVED***
