package dockerfile

import (
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/builder"
	containerpkg "github.com/docker/docker/container"
	"github.com/docker/docker/pkg/stringid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type containerManager struct ***REMOVED***
	tmpContainers map[string]struct***REMOVED******REMOVED***
	backend       builder.ExecBackend
***REMOVED***

// newContainerManager creates a new container backend
func newContainerManager(docker builder.ExecBackend) *containerManager ***REMOVED***
	return &containerManager***REMOVED***
		backend:       docker,
		tmpContainers: make(map[string]struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// Create a container
func (c *containerManager) Create(runConfig *container.Config, hostConfig *container.HostConfig) (container.ContainerCreateCreatedBody, error) ***REMOVED***
	container, err := c.backend.ContainerCreate(types.ContainerCreateConfig***REMOVED***
		Config:     runConfig,
		HostConfig: hostConfig,
	***REMOVED***)
	if err != nil ***REMOVED***
		return container, err
	***REMOVED***
	c.tmpContainers[container.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return container, nil
***REMOVED***

var errCancelled = errors.New("build cancelled")

// Run a container by ID
func (c *containerManager) Run(ctx context.Context, cID string, stdout, stderr io.Writer) (err error) ***REMOVED***
	attached := make(chan struct***REMOVED******REMOVED***)
	errCh := make(chan error)
	go func() ***REMOVED***
		errCh <- c.backend.ContainerAttachRaw(cID, nil, stdout, stderr, true, attached)
	***REMOVED***()
	select ***REMOVED***
	case err := <-errCh:
		return err
	case <-attached:
	***REMOVED***

	finished := make(chan struct***REMOVED******REMOVED***)
	cancelErrCh := make(chan error, 1)
	go func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			logrus.Debugln("Build cancelled, killing and removing container:", cID)
			c.backend.ContainerKill(cID, 0)
			c.removeContainer(cID, stdout)
			cancelErrCh <- errCancelled
		case <-finished:
			cancelErrCh <- nil
		***REMOVED***
	***REMOVED***()

	if err := c.backend.ContainerStart(cID, nil, "", ""); err != nil ***REMOVED***
		close(finished)
		logCancellationError(cancelErrCh, "error from ContainerStart: "+err.Error())
		return err
	***REMOVED***

	// Block on reading output from container, stop on err or chan closed
	if err := <-errCh; err != nil ***REMOVED***
		close(finished)
		logCancellationError(cancelErrCh, "error from errCh: "+err.Error())
		return err
	***REMOVED***

	waitC, err := c.backend.ContainerWait(ctx, cID, containerpkg.WaitConditionNotRunning)
	if err != nil ***REMOVED***
		close(finished)
		logCancellationError(cancelErrCh, fmt.Sprintf("unable to begin ContainerWait: %s", err))
		return err
	***REMOVED***

	if status := <-waitC; status.ExitCode() != 0 ***REMOVED***
		close(finished)
		logCancellationError(cancelErrCh,
			fmt.Sprintf("a non-zero code from ContainerWait: %d", status.ExitCode()))
		return &statusCodeError***REMOVED***code: status.ExitCode(), err: err***REMOVED***
	***REMOVED***

	close(finished)
	return <-cancelErrCh
***REMOVED***

func logCancellationError(cancelErrCh chan error, msg string) ***REMOVED***
	if cancelErr := <-cancelErrCh; cancelErr != nil ***REMOVED***
		logrus.Debugf("Build cancelled (%v): %s", cancelErr, msg)
	***REMOVED***
***REMOVED***

type statusCodeError struct ***REMOVED***
	code int
	err  error
***REMOVED***

func (e *statusCodeError) Error() string ***REMOVED***
	return e.err.Error()
***REMOVED***

func (e *statusCodeError) StatusCode() int ***REMOVED***
	return e.code
***REMOVED***

func (c *containerManager) removeContainer(containerID string, stdout io.Writer) error ***REMOVED***
	rmConfig := &types.ContainerRmConfig***REMOVED***
		ForceRemove:  true,
		RemoveVolume: true,
	***REMOVED***
	if err := c.backend.ContainerRm(containerID, rmConfig); err != nil ***REMOVED***
		fmt.Fprintf(stdout, "Error removing intermediate container %s: %v\n", stringid.TruncateID(containerID), err)
		return err
	***REMOVED***
	return nil
***REMOVED***

// RemoveAll containers managed by this container manager
func (c *containerManager) RemoveAll(stdout io.Writer) ***REMOVED***
	for containerID := range c.tmpContainers ***REMOVED***
		if err := c.removeContainer(containerID, stdout); err != nil ***REMOVED***
			return
		***REMOVED***
		delete(c.tmpContainers, containerID)
		fmt.Fprintf(stdout, "Removing intermediate container %s\n", stringid.TruncateID(containerID))
	***REMOVED***
***REMOVED***
