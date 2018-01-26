// +build !windows,!linux

package shim

import (
	"context"
	"io"
	"sync"
	"syscall"

	"github.com/containerd/console"
	"github.com/containerd/fifo"
)

type unixPlatform struct ***REMOVED***
***REMOVED***

func (p *unixPlatform) CopyConsole(ctx context.Context, console console.Console, stdin, stdout, stderr string, wg, cwg *sync.WaitGroup) (console.Console, error) ***REMOVED***
	if stdin != "" ***REMOVED***
		in, err := fifo.OpenFifo(ctx, stdin, syscall.O_RDONLY, 0)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cwg.Add(1)
		go func() ***REMOVED***
			cwg.Done()
			io.Copy(console, in)
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
		io.Copy(outw, console)
		console.Close()
		outr.Close()
		outw.Close()
		wg.Done()
	***REMOVED***()
	return console, nil
***REMOVED***

func (p *unixPlatform) ShutdownConsole(ctx context.Context, cons console.Console) error ***REMOVED***
	return nil
***REMOVED***

func (p *unixPlatform) Close() error ***REMOVED***
	return nil
***REMOVED***

func (s *Service) initPlatform() error ***REMOVED***
	s.platform = &unixPlatform***REMOVED******REMOVED***
	return nil
***REMOVED***
