package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/container"
	"github.com/sirupsen/logrus"
)

// ContainerUnpause unpauses a container
func (daemon *Daemon) ContainerUnpause(name string) error ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return daemon.containerUnpause(container)
***REMOVED***

// containerUnpause resumes the container execution after the container is paused.
func (daemon *Daemon) containerUnpause(container *container.Container) error ***REMOVED***
	container.Lock()
	defer container.Unlock()

	// We cannot unpause the container which is not paused
	if !container.Paused ***REMOVED***
		return fmt.Errorf("Container %s is not paused", container.ID)
	***REMOVED***

	if err := daemon.containerd.Resume(context.Background(), container.ID); err != nil ***REMOVED***
		return fmt.Errorf("Cannot unpause container %s: %s", container.ID, err)
	***REMOVED***

	container.Paused = false
	daemon.setStateCounter(container)
	daemon.updateHealthMonitor(container)
	daemon.LogContainerEvent(container, "unpause")

	if err := container.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
		logrus.WithError(err).Warnf("could not save container to disk")
	***REMOVED***

	return nil
***REMOVED***
