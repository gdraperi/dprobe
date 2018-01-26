package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/errdefs"
	"github.com/pkg/errors"
)

// ContainerUpdate updates configuration of the container
func (daemon *Daemon) ContainerUpdate(name string, hostConfig *container.HostConfig) (container.ContainerUpdateOKBody, error) ***REMOVED***
	var warnings []string

	c, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return container.ContainerUpdateOKBody***REMOVED***Warnings: warnings***REMOVED***, err
	***REMOVED***

	warnings, err = daemon.verifyContainerSettings(c.OS, hostConfig, nil, true)
	if err != nil ***REMOVED***
		return container.ContainerUpdateOKBody***REMOVED***Warnings: warnings***REMOVED***, errdefs.InvalidParameter(err)
	***REMOVED***

	if err := daemon.update(name, hostConfig); err != nil ***REMOVED***
		return container.ContainerUpdateOKBody***REMOVED***Warnings: warnings***REMOVED***, err
	***REMOVED***

	return container.ContainerUpdateOKBody***REMOVED***Warnings: warnings***REMOVED***, nil
***REMOVED***

func (daemon *Daemon) update(name string, hostConfig *container.HostConfig) error ***REMOVED***
	if hostConfig == nil ***REMOVED***
		return nil
	***REMOVED***

	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	restoreConfig := false
	backupHostConfig := *container.HostConfig
	defer func() ***REMOVED***
		if restoreConfig ***REMOVED***
			container.Lock()
			container.HostConfig = &backupHostConfig
			container.CheckpointTo(daemon.containersReplica)
			container.Unlock()
		***REMOVED***
	***REMOVED***()

	if container.RemovalInProgress || container.Dead ***REMOVED***
		return errCannotUpdate(container.ID, fmt.Errorf("container is marked for removal and cannot be \"update\""))
	***REMOVED***

	container.Lock()
	if err := container.UpdateContainer(hostConfig); err != nil ***REMOVED***
		restoreConfig = true
		container.Unlock()
		return errCannotUpdate(container.ID, err)
	***REMOVED***
	if err := container.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
		restoreConfig = true
		container.Unlock()
		return errCannotUpdate(container.ID, err)
	***REMOVED***
	container.Unlock()

	// if Restart Policy changed, we need to update container monitor
	if hostConfig.RestartPolicy.Name != "" ***REMOVED***
		container.UpdateMonitor(hostConfig.RestartPolicy)
	***REMOVED***

	// If container is not running, update hostConfig struct is enough,
	// resources will be updated when the container is started again.
	// If container is running (including paused), we need to update configs
	// to the real world.
	if container.IsRunning() && !container.IsRestarting() ***REMOVED***
		if err := daemon.containerd.UpdateResources(context.Background(), container.ID, toContainerdResources(hostConfig.Resources)); err != nil ***REMOVED***
			restoreConfig = true
			// TODO: it would be nice if containerd responded with better errors here so we can classify this better.
			return errCannotUpdate(container.ID, errdefs.System(err))
		***REMOVED***
	***REMOVED***

	daemon.LogContainerEvent(container, "update")

	return nil
***REMOVED***

func errCannotUpdate(containerID string, err error) error ***REMOVED***
	return errors.Wrap(err, "Cannot update container "+containerID)
***REMOVED***
