package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/container"
	"github.com/sirupsen/logrus"
)

// ContainerPause pauses a container
func (daemon *Daemon) ContainerPause(name string) error ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return daemon.containerPause(container)
***REMOVED***

// containerPause pauses the container execution without stopping the process.
// The execution can be resumed by calling containerUnpause.
func (daemon *Daemon) containerPause(container *container.Container) error ***REMOVED***
	container.Lock()
	defer container.Unlock()

	// We cannot Pause the container which is not running
	if !container.Running ***REMOVED***
		return errNotRunning(container.ID)
	***REMOVED***

	// We cannot Pause the container which is already paused
	if container.Paused ***REMOVED***
		return errNotPaused(container.ID)
	***REMOVED***

	// We cannot Pause the container which is restarting
	if container.Restarting ***REMOVED***
		return errContainerIsRestarting(container.ID)
	***REMOVED***

	if err := daemon.containerd.Pause(context.Background(), container.ID); err != nil ***REMOVED***
		return fmt.Errorf("Cannot pause container %s: %s", container.ID, err)
	***REMOVED***

	container.Paused = true
	daemon.setStateCounter(container)
	daemon.updateHealthMonitor(container)
	daemon.LogContainerEvent(container, "pause")

	if err := container.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
		logrus.WithError(err).Warn("could not save container to disk")
	***REMOVED***

	return nil
***REMOVED***
