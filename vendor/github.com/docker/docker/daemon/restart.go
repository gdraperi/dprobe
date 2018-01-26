package daemon

import (
	"fmt"

	"github.com/docker/docker/container"
	"github.com/sirupsen/logrus"
)

// ContainerRestart stops and starts a container. It attempts to
// gracefully stop the container within the given timeout, forcefully
// stopping it if the timeout is exceeded. If given a negative
// timeout, ContainerRestart will wait forever until a graceful
// stop. Returns an error if the container cannot be found, or if
// there is an underlying error at any stage of the restart.
func (daemon *Daemon) ContainerRestart(name string, seconds *int) error ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if seconds == nil ***REMOVED***
		stopTimeout := container.StopTimeout()
		seconds = &stopTimeout
	***REMOVED***
	if err := daemon.containerRestart(container, *seconds); err != nil ***REMOVED***
		return fmt.Errorf("Cannot restart container %s: %v", name, err)
	***REMOVED***
	return nil

***REMOVED***

// containerRestart attempts to gracefully stop and then start the
// container. When stopping, wait for the given duration in seconds to
// gracefully stop, before forcefully terminating the container. If
// given a negative duration, wait forever for a graceful stop.
func (daemon *Daemon) containerRestart(container *container.Container, seconds int) error ***REMOVED***
	// Avoid unnecessarily unmounting and then directly mounting
	// the container when the container stops and then starts
	// again
	if err := daemon.Mount(container); err == nil ***REMOVED***
		defer daemon.Unmount(container)
	***REMOVED***

	if container.IsRunning() ***REMOVED***
		// set AutoRemove flag to false before stop so the container won't be
		// removed during restart process
		autoRemove := container.HostConfig.AutoRemove

		container.HostConfig.AutoRemove = false
		err := daemon.containerStop(container, seconds)
		// restore AutoRemove irrespective of whether the stop worked or not
		container.HostConfig.AutoRemove = autoRemove
		// containerStop will write HostConfig to disk, we shall restore AutoRemove
		// in disk too
		if toDiskErr := daemon.checkpointAndSave(container); toDiskErr != nil ***REMOVED***
			logrus.Errorf("Write container to disk error: %v", toDiskErr)
		***REMOVED***

		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if err := daemon.containerStart(container, "", "", true); err != nil ***REMOVED***
		return err
	***REMOVED***

	daemon.LogContainerEvent(container, "restart")
	return nil
***REMOVED***
