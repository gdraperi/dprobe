package containerd

import (
	"context"
	"io"
	"path/filepath"
	"sync"

	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/linux/runctypes"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/libcontainerd"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// PluginNamespace is the name used for the plugins namespace
var PluginNamespace = "plugins.moby"

// ExitHandler represents an object that is called when the exit event is received from containerd
type ExitHandler interface ***REMOVED***
	HandleExitEvent(id string) error
***REMOVED***

// New creates a new containerd plugin executor
func New(rootDir string, remote libcontainerd.Remote, exitHandler ExitHandler) (*Executor, error) ***REMOVED***
	e := &Executor***REMOVED***
		rootDir:     rootDir,
		exitHandler: exitHandler,
	***REMOVED***
	client, err := remote.NewClient(PluginNamespace, e)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "error creating containerd exec client")
	***REMOVED***
	e.client = client
	return e, nil
***REMOVED***

// Executor is the containerd client implementation of a plugin executor
type Executor struct ***REMOVED***
	rootDir     string
	client      libcontainerd.Client
	exitHandler ExitHandler
***REMOVED***

// Create creates a new container
func (e *Executor) Create(id string, spec specs.Spec, stdout, stderr io.WriteCloser) error ***REMOVED***
	opts := runctypes.RuncOptions***REMOVED***
		RuntimeRoot: filepath.Join(e.rootDir, "runtime-root"),
	***REMOVED***
	ctx := context.Background()
	err := e.client.Create(ctx, id, &spec, &opts)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	_, err = e.client.Start(ctx, id, "", false, attachStreamsFunc(stdout, stderr))
	return err
***REMOVED***

// Restore restores a container
func (e *Executor) Restore(id string, stdout, stderr io.WriteCloser) error ***REMOVED***
	alive, _, err := e.client.Restore(context.Background(), id, attachStreamsFunc(stdout, stderr))
	if err != nil && !errdefs.IsNotFound(err) ***REMOVED***
		return err
	***REMOVED***
	if !alive ***REMOVED***
		_, _, err = e.client.DeleteTask(context.Background(), id)
		if err != nil && !errdefs.IsNotFound(err) ***REMOVED***
			logrus.WithError(err).Errorf("failed to delete container plugin %s task from containerd", id)
			return err
		***REMOVED***

		err = e.client.Delete(context.Background(), id)
		if err != nil && !errdefs.IsNotFound(err) ***REMOVED***
			logrus.WithError(err).Errorf("failed to delete container plugin %s from containerd", id)
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// IsRunning returns if the container with the given id is running
func (e *Executor) IsRunning(id string) (bool, error) ***REMOVED***
	status, err := e.client.Status(context.Background(), id)
	return status == libcontainerd.StatusRunning, err
***REMOVED***

// Signal sends the specified signal to the container
func (e *Executor) Signal(id string, signal int) error ***REMOVED***
	return e.client.SignalProcess(context.Background(), id, libcontainerd.InitProcessName, signal)
***REMOVED***

// ProcessEvent handles events from containerd
// All events are ignored except the exit event, which is sent of to the stored handler
func (e *Executor) ProcessEvent(id string, et libcontainerd.EventType, ei libcontainerd.EventInfo) error ***REMOVED***
	switch et ***REMOVED***
	case libcontainerd.EventExit:
		// delete task and container
		if _, _, err := e.client.DeleteTask(context.Background(), id); err != nil ***REMOVED***
			logrus.WithError(err).Errorf("failed to delete container plugin %s task from containerd", id)
		***REMOVED***

		if err := e.client.Delete(context.Background(), id); err != nil ***REMOVED***
			logrus.WithError(err).Errorf("failed to delete container plugin %s from containerd", id)
		***REMOVED***
		return e.exitHandler.HandleExitEvent(ei.ContainerID)
	***REMOVED***
	return nil
***REMOVED***

type rio struct ***REMOVED***
	cio.IO

	wg sync.WaitGroup
***REMOVED***

func (c *rio) Wait() ***REMOVED***
	c.wg.Wait()
	c.IO.Wait()
***REMOVED***

func attachStreamsFunc(stdout, stderr io.WriteCloser) libcontainerd.StdioCallback ***REMOVED***
	return func(iop *cio.DirectIO) (cio.IO, error) ***REMOVED***
		if iop.Stdin != nil ***REMOVED***
			iop.Stdin.Close()
			// closing stdin shouldn't be needed here, it should never be open
			panic("plugin stdin shouldn't have been created!")
		***REMOVED***

		rio := &rio***REMOVED***IO: iop***REMOVED***
		rio.wg.Add(2)
		go func() ***REMOVED***
			io.Copy(stdout, iop.Stdout)
			stdout.Close()
			rio.wg.Done()
		***REMOVED***()
		go func() ***REMOVED***
			io.Copy(stderr, iop.Stderr)
			stderr.Close()
			rio.wg.Done()
		***REMOVED***()
		return rio, nil
	***REMOVED***
***REMOVED***
