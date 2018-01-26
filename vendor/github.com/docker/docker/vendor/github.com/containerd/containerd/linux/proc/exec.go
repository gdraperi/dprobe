// +build !windows

package proc

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"

	"github.com/containerd/console"
	"github.com/containerd/fifo"
	runc "github.com/containerd/go-runc"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

type execProcess struct ***REMOVED***
	wg sync.WaitGroup

	State

	mu      sync.Mutex
	id      string
	console console.Console
	io      runc.IO
	status  int
	exited  time.Time
	pid     int
	closers []io.Closer
	stdin   io.Closer
	stdio   Stdio
	path    string
	spec    specs.Process

	parent    *Init
	waitBlock chan struct***REMOVED******REMOVED***
***REMOVED***

func (e *execProcess) Wait() ***REMOVED***
	<-e.waitBlock
***REMOVED***

func (e *execProcess) ID() string ***REMOVED***
	return e.id
***REMOVED***

func (e *execProcess) Pid() int ***REMOVED***
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.pid
***REMOVED***

func (e *execProcess) ExitStatus() int ***REMOVED***
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.status
***REMOVED***

func (e *execProcess) ExitedAt() time.Time ***REMOVED***
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.exited
***REMOVED***

func (e *execProcess) setExited(status int) ***REMOVED***
	e.status = status
	e.exited = time.Now()
	e.parent.platform.ShutdownConsole(context.Background(), e.console)
	close(e.waitBlock)
***REMOVED***

func (e *execProcess) delete(ctx context.Context) error ***REMOVED***
	e.wg.Wait()
	if e.io != nil ***REMOVED***
		for _, c := range e.closers ***REMOVED***
			c.Close()
		***REMOVED***
		e.io.Close()
	***REMOVED***
	pidfile := filepath.Join(e.path, fmt.Sprintf("%s.pid", e.id))
	// silently ignore error
	os.Remove(pidfile)
	return nil
***REMOVED***

func (e *execProcess) resize(ws console.WinSize) error ***REMOVED***
	if e.console == nil ***REMOVED***
		return nil
	***REMOVED***
	return e.console.Resize(ws)
***REMOVED***

func (e *execProcess) kill(ctx context.Context, sig uint32, _ bool) error ***REMOVED***
	pid := e.pid
	if pid != 0 ***REMOVED***
		if err := unix.Kill(pid, syscall.Signal(sig)); err != nil ***REMOVED***
			return errors.Wrapf(checkKillError(err), "exec kill error")
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (e *execProcess) Stdin() io.Closer ***REMOVED***
	return e.stdin
***REMOVED***

func (e *execProcess) Stdio() Stdio ***REMOVED***
	return e.stdio
***REMOVED***

func (e *execProcess) start(ctx context.Context) (err error) ***REMOVED***
	var (
		socket  *runc.Socket
		pidfile = filepath.Join(e.path, fmt.Sprintf("%s.pid", e.id))
	)
	if e.stdio.Terminal ***REMOVED***
		if socket, err = runc.NewTempConsoleSocket(); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to create runc console socket")
		***REMOVED***
		defer socket.Close()
	***REMOVED*** else if e.stdio.IsNull() ***REMOVED***
		if e.io, err = runc.NewNullIO(); err != nil ***REMOVED***
			return errors.Wrap(err, "creating new NULL IO")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if e.io, err = runc.NewPipeIO(e.parent.IoUID, e.parent.IoGID); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to create runc io pipes")
		***REMOVED***
	***REMOVED***
	opts := &runc.ExecOpts***REMOVED***
		PidFile: pidfile,
		IO:      e.io,
		Detach:  true,
	***REMOVED***
	if socket != nil ***REMOVED***
		opts.ConsoleSocket = socket
	***REMOVED***
	if err := e.parent.runtime.Exec(ctx, e.parent.id, e.spec, opts); err != nil ***REMOVED***
		close(e.waitBlock)
		return e.parent.runtimeError(err, "OCI runtime exec failed")
	***REMOVED***
	if e.stdio.Stdin != "" ***REMOVED***
		sc, err := fifo.OpenFifo(ctx, e.stdio.Stdin, syscall.O_WRONLY|syscall.O_NONBLOCK, 0)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to open stdin fifo %s", e.stdio.Stdin)
		***REMOVED***
		e.closers = append(e.closers, sc)
		e.stdin = sc
	***REMOVED***
	var copyWaitGroup sync.WaitGroup
	if socket != nil ***REMOVED***
		console, err := socket.ReceiveMaster()
		if err != nil ***REMOVED***
			return errors.Wrap(err, "failed to retrieve console master")
		***REMOVED***
		if e.console, err = e.parent.platform.CopyConsole(ctx, console, e.stdio.Stdin, e.stdio.Stdout, e.stdio.Stderr, &e.wg, &copyWaitGroup); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to start console copy")
		***REMOVED***
	***REMOVED*** else if !e.stdio.IsNull() ***REMOVED***
		if err := copyPipes(ctx, e.io, e.stdio.Stdin, e.stdio.Stdout, e.stdio.Stderr, &e.wg, &copyWaitGroup); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to start io pipe copy")
		***REMOVED***
	***REMOVED***
	copyWaitGroup.Wait()
	pid, err := runc.ReadPidFile(opts.PidFile)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to retrieve OCI runtime exec pid")
	***REMOVED***
	e.pid = pid
	return nil
***REMOVED***

func (e *execProcess) Status(ctx context.Context) (string, error) ***REMOVED***
	s, err := e.parent.Status(ctx)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	// if the container as a whole is in the pausing/paused state, so are all
	// other processes inside the container, use container state here
	switch s ***REMOVED***
	case "paused", "pausing":
		return s, nil
	***REMOVED***
	e.mu.Lock()
	defer e.mu.Unlock()
	// if we don't have a pid then the exec process has just been created
	if e.pid == 0 ***REMOVED***
		return "created", nil
	***REMOVED***
	// if we have a pid and it can be signaled, the process is running
	if err := unix.Kill(e.pid, 0); err == nil ***REMOVED***
		return "running", nil
	***REMOVED***
	// else if we have a pid but it can nolonger be signaled, it has stopped
	return "stopped", nil
***REMOVED***
