package daemon

import (
	"context"
	"runtime"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/mount"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ContainerStart starts a container.
func (daemon *Daemon) ContainerStart(name string, hostConfig *containertypes.HostConfig, checkpoint string, checkpointDir string) error ***REMOVED***
	if checkpoint != "" && !daemon.HasExperimental() ***REMOVED***
		return errdefs.InvalidParameter(errors.New("checkpoint is only supported in experimental mode"))
	***REMOVED***

	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	validateState := func() error ***REMOVED***
		container.Lock()
		defer container.Unlock()

		if container.Paused ***REMOVED***
			return errdefs.Conflict(errors.New("cannot start a paused container, try unpause instead"))
		***REMOVED***

		if container.Running ***REMOVED***
			return containerNotModifiedError***REMOVED***running: true***REMOVED***
		***REMOVED***

		if container.RemovalInProgress || container.Dead ***REMOVED***
			return errdefs.Conflict(errors.New("container is marked for removal and cannot be started"))
		***REMOVED***
		return nil
	***REMOVED***

	if err := validateState(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Windows does not have the backwards compatibility issue here.
	if runtime.GOOS != "windows" ***REMOVED***
		// This is kept for backward compatibility - hostconfig should be passed when
		// creating a container, not during start.
		if hostConfig != nil ***REMOVED***
			logrus.Warn("DEPRECATED: Setting host configuration options when the container starts is deprecated and has been removed in Docker 1.12")
			oldNetworkMode := container.HostConfig.NetworkMode
			if err := daemon.setSecurityOptions(container, hostConfig); err != nil ***REMOVED***
				return errdefs.InvalidParameter(err)
			***REMOVED***
			if err := daemon.mergeAndVerifyLogConfig(&hostConfig.LogConfig); err != nil ***REMOVED***
				return errdefs.InvalidParameter(err)
			***REMOVED***
			if err := daemon.setHostConfig(container, hostConfig); err != nil ***REMOVED***
				return errdefs.InvalidParameter(err)
			***REMOVED***
			newNetworkMode := container.HostConfig.NetworkMode
			if string(oldNetworkMode) != string(newNetworkMode) ***REMOVED***
				// if user has change the network mode on starting, clean up the
				// old networks. It is a deprecated feature and has been removed in Docker 1.12
				container.NetworkSettings.Networks = nil
				if err := container.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
					return errdefs.System(err)
				***REMOVED***
			***REMOVED***
			container.InitDNSHostConfig()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if hostConfig != nil ***REMOVED***
			return errdefs.InvalidParameter(errors.New("Supplying a hostconfig on start is not supported. It should be supplied on create"))
		***REMOVED***
	***REMOVED***

	// check if hostConfig is in line with the current system settings.
	// It may happen cgroups are umounted or the like.
	if _, err = daemon.verifyContainerSettings(container.OS, container.HostConfig, nil, false); err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***
	// Adapt for old containers in case we have updates in this function and
	// old containers never have chance to call the new function in create stage.
	if hostConfig != nil ***REMOVED***
		if err := daemon.adaptContainerSettings(container.HostConfig, false); err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***
	***REMOVED***
	return daemon.containerStart(container, checkpoint, checkpointDir, true)
***REMOVED***

// containerStart prepares the container to run by setting up everything the
// container needs, such as storage and networking, as well as links
// between containers. The container is left waiting for a signal to
// begin running.
func (daemon *Daemon) containerStart(container *container.Container, checkpoint string, checkpointDir string, resetRestartManager bool) (err error) ***REMOVED***
	start := time.Now()
	container.Lock()
	defer container.Unlock()

	if resetRestartManager && container.Running ***REMOVED*** // skip this check if already in restarting step and resetRestartManager==false
		return nil
	***REMOVED***

	if container.RemovalInProgress || container.Dead ***REMOVED***
		return errdefs.Conflict(errors.New("container is marked for removal and cannot be started"))
	***REMOVED***

	if checkpointDir != "" ***REMOVED***
		// TODO(mlaventure): how would we support that?
		return errdefs.Forbidden(errors.New("custom checkpointdir is not supported"))
	***REMOVED***

	// if we encounter an error during start we need to ensure that any other
	// setup has been cleaned up properly
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			container.SetError(err)
			// if no one else has set it, make sure we don't leave it at zero
			if container.ExitCode() == 0 ***REMOVED***
				container.SetExitCode(128)
			***REMOVED***
			if err := container.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
				logrus.Errorf("%s: failed saving state on start failure: %v", container.ID, err)
			***REMOVED***
			container.Reset(false)

			daemon.Cleanup(container)
			// if containers AutoRemove flag is set, remove it after clean up
			if container.HostConfig.AutoRemove ***REMOVED***
				container.Unlock()
				if err := daemon.ContainerRm(container.ID, &types.ContainerRmConfig***REMOVED***ForceRemove: true, RemoveVolume: true***REMOVED***); err != nil ***REMOVED***
					logrus.Errorf("can't remove container %s: %v", container.ID, err)
				***REMOVED***
				container.Lock()
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if err := daemon.conditionalMountOnStart(container); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := daemon.initializeNetworking(container); err != nil ***REMOVED***
		return err
	***REMOVED***

	spec, err := daemon.createSpec(container)
	if err != nil ***REMOVED***
		return errdefs.System(err)
	***REMOVED***

	if resetRestartManager ***REMOVED***
		container.ResetRestartManager(true)
	***REMOVED***

	if daemon.saveApparmorConfig(container); err != nil ***REMOVED***
		return err
	***REMOVED***

	if checkpoint != "" ***REMOVED***
		checkpointDir, err = getCheckpointDir(checkpointDir, checkpoint, container.Name, container.ID, container.CheckpointDir(), false)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	createOptions, err := daemon.getLibcontainerdCreateOptions(container)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = daemon.containerd.Create(context.Background(), container.ID, spec, createOptions)
	if err != nil ***REMOVED***
		return translateContainerdStartErr(container.Path, container.SetExitCode, err)
	***REMOVED***

	// TODO(mlaventure): we need to specify checkpoint options here
	pid, err := daemon.containerd.Start(context.Background(), container.ID, checkpointDir,
		container.StreamConfig.Stdin() != nil || container.Config.Tty,
		container.InitializeStdio)
	if err != nil ***REMOVED***
		if err := daemon.containerd.Delete(context.Background(), container.ID); err != nil ***REMOVED***
			logrus.WithError(err).WithField("container", container.ID).
				Error("failed to delete failed start container")
		***REMOVED***
		return translateContainerdStartErr(container.Path, container.SetExitCode, err)
	***REMOVED***

	container.SetRunning(pid, true)
	container.HasBeenManuallyStopped = false
	container.HasBeenStartedBefore = true
	daemon.setStateCounter(container)

	daemon.initHealthMonitor(container)

	if err := container.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
		logrus.WithError(err).WithField("container", container.ID).
			Errorf("failed to store container")
	***REMOVED***

	daemon.LogContainerEvent(container, "start")
	containerActions.WithValues("start").UpdateSince(start)

	return nil
***REMOVED***

// Cleanup releases any network resources allocated to the container along with any rules
// around how containers are linked together.  It also unmounts the container's root filesystem.
func (daemon *Daemon) Cleanup(container *container.Container) ***REMOVED***
	daemon.releaseNetwork(container)

	if err := container.UnmountIpcMount(detachMounted); err != nil ***REMOVED***
		logrus.Warnf("%s cleanup: failed to unmount IPC: %s", container.ID, err)
	***REMOVED***

	if err := daemon.conditionalUnmountOnCleanup(container); err != nil ***REMOVED***
		// FIXME: remove once reference counting for graphdrivers has been refactored
		// Ensure that all the mounts are gone
		if mountid, err := daemon.layerStores[container.OS].GetMountID(container.ID); err == nil ***REMOVED***
			daemon.cleanupMountsByID(mountid)
		***REMOVED***
	***REMOVED***

	if err := container.UnmountSecrets(); err != nil ***REMOVED***
		logrus.Warnf("%s cleanup: failed to unmount secrets: %s", container.ID, err)
	***REMOVED***

	if err := mount.RecursiveUnmount(container.Root); err != nil ***REMOVED***
		logrus.WithError(err).WithField("container", container.ID).Warn("Error while cleaning up container resource mounts.")
	***REMOVED***

	for _, eConfig := range container.ExecCommands.Commands() ***REMOVED***
		daemon.unregisterExecCommand(container, eConfig)
	***REMOVED***

	if container.BaseFS != nil && container.BaseFS.Path() != "" ***REMOVED***
		if err := container.UnmountVolumes(daemon.LogVolumeEvent); err != nil ***REMOVED***
			logrus.Warnf("%s cleanup: Failed to umount volumes: %v", container.ID, err)
		***REMOVED***
	***REMOVED***

	container.CancelAttachContext()

	if err := daemon.containerd.Delete(context.Background(), container.ID); err != nil ***REMOVED***
		logrus.Errorf("%s cleanup: failed to delete container from containerd: %v", container.ID, err)
	***REMOVED***
***REMOVED***
