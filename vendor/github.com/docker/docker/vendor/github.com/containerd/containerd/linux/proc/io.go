// +build !windows

package proc

import (
	"context"
	"fmt"
	"io"
	"sync"
	"syscall"

	"github.com/containerd/fifo"
	runc "github.com/containerd/go-runc"
)

func copyPipes(ctx context.Context, rio runc.IO, stdin, stdout, stderr string, wg, cwg *sync.WaitGroup) error ***REMOVED***
	for name, dest := range map[string]func(wc io.WriteCloser, rc io.Closer)***REMOVED***
		stdout: func(wc io.WriteCloser, rc io.Closer) ***REMOVED***
			wg.Add(1)
			cwg.Add(1)
			go func() ***REMOVED***
				cwg.Done()
				io.Copy(wc, rio.Stdout())
				wg.Done()
				wc.Close()
				rc.Close()
			***REMOVED***()
		***REMOVED***,
		stderr: func(wc io.WriteCloser, rc io.Closer) ***REMOVED***
			wg.Add(1)
			cwg.Add(1)
			go func() ***REMOVED***
				cwg.Done()
				io.Copy(wc, rio.Stderr())
				wg.Done()
				wc.Close()
				rc.Close()
			***REMOVED***()
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		fw, err := fifo.OpenFifo(ctx, name, syscall.O_WRONLY, 0)
		if err != nil ***REMOVED***
			return fmt.Errorf("containerd-shim: opening %s failed: %s", name, err)
		***REMOVED***
		fr, err := fifo.OpenFifo(ctx, name, syscall.O_RDONLY, 0)
		if err != nil ***REMOVED***
			return fmt.Errorf("containerd-shim: opening %s failed: %s", name, err)
		***REMOVED***
		dest(fw, fr)
	***REMOVED***
	if stdin == "" ***REMOVED***
		rio.Stdin().Close()
		return nil
	***REMOVED***
	f, err := fifo.OpenFifo(ctx, stdin, syscall.O_RDONLY, 0)
	if err != nil ***REMOVED***
		return fmt.Errorf("containerd-shim: opening %s failed: %s", stdin, err)
	***REMOVED***
	cwg.Add(1)
	go func() ***REMOVED***
		cwg.Done()
		io.Copy(rio.Stdin(), f)
		rio.Stdin().Close()
		f.Close()
	***REMOVED***()
	return nil
***REMOVED***
