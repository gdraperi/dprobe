package daemon

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/volume"
	volumestore "github.com/docker/docker/volume/store"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ContainerRm removes the container id from the filesystem. An error
// is returned if the container is not found, or if the remove
// fails. If the remove succeeds, the container name is released, and
// network links are removed.
func (daemon *Daemon) ContainerRm(name string, config *types.ContainerRmConfig) error ***REMOVED***
	start := time.Now()
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Container state RemovalInProgress should be used to avoid races.
	if inProgress := container.SetRemovalInProgress(); inProgress ***REMOVED***
		err := fmt.Errorf("removal of container %s is already in progress", name)
		return errdefs.Conflict(err)
	***REMOVED***
	defer container.ResetRemovalInProgress()

	// check if container wasn't deregistered by previous rm since Get
	if c := daemon.containers.Get(container.ID); c == nil ***REMOVED***
		return nil
	***REMOVED***

	if config.RemoveLink ***REMOVED***
		return daemon.rmLink(container, name)
	***REMOVED***

	err = daemon.cleanupContainer(container, config.ForceRemove, config.RemoveVolume)
	containerActions.WithValues("delete").UpdateSince(start)

	return err
***REMOVED***

func (daemon *Daemon) rmLink(container *container.Container, name string) error ***REMOVED***
	if name[0] != '/' ***REMOVED***
		name = "/" + name
	***REMOVED***
	parent, n := path.Split(name)
	if parent == "/" ***REMOVED***
		return fmt.Errorf("Conflict, cannot remove the default name of the container")
	***REMOVED***

	parent = strings.TrimSuffix(parent, "/")
	pe, err := daemon.containersReplica.Snapshot().GetID(parent)
	if err != nil ***REMOVED***
		return fmt.Errorf("Cannot get parent %s for name %s", parent, name)
	***REMOVED***

	daemon.releaseName(name)
	parentContainer, _ := daemon.GetContainer(pe)
	if parentContainer != nil ***REMOVED***
		daemon.linkIndex.unlink(name, container, parentContainer)
		if err := daemon.updateNetwork(parentContainer); err != nil ***REMOVED***
			logrus.Debugf("Could not update network to remove link %s: %v", n, err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// cleanupContainer unregisters a container from the daemon, stops stats
// collection and cleanly removes contents and metadata from the filesystem.
func (daemon *Daemon) cleanupContainer(container *container.Container, forceRemove, removeVolume bool) (err error) ***REMOVED***
	if container.IsRunning() ***REMOVED***
		if !forceRemove ***REMOVED***
			state := container.StateString()
			procedure := "Stop the container before attempting removal or force remove"
			if state == "paused" ***REMOVED***
				procedure = "Unpause and then " + strings.ToLower(procedure)
			***REMOVED***
			err := fmt.Errorf("You cannot remove a %s container %s. %s", state, container.ID, procedure)
			return errdefs.Conflict(err)
		***REMOVED***
		if err := daemon.Kill(container); err != nil ***REMOVED***
			return fmt.Errorf("Could not kill running container %s, cannot remove - %v", container.ID, err)
		***REMOVED***
	***REMOVED***
	if !system.IsOSSupported(container.OS) ***REMOVED***
		return fmt.Errorf("cannot remove %s: %s ", container.ID, system.ErrNotSupportedOperatingSystem)
	***REMOVED***

	// stop collection of stats for the container regardless
	// if stats are currently getting collected.
	daemon.statsCollector.StopCollection(container)

	if err = daemon.containerStop(container, 3); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Mark container dead. We don't want anybody to be restarting it.
	container.Lock()
	container.Dead = true

	// Save container state to disk. So that if error happens before
	// container meta file got removed from disk, then a restart of
	// docker should not make a dead container alive.
	if err := container.CheckpointTo(daemon.containersReplica); err != nil && !os.IsNotExist(err) ***REMOVED***
		logrus.Errorf("Error saving dying container to disk: %v", err)
	***REMOVED***
	container.Unlock()

	// When container creation fails and `RWLayer` has not been created yet, we
	// do not call `ReleaseRWLayer`
	if container.RWLayer != nil ***REMOVED***
		metadata, err := daemon.layerStores[container.OS].ReleaseRWLayer(container.RWLayer)
		layer.LogReleaseMetadata(metadata)
		if err != nil && err != layer.ErrMountDoesNotExist && !os.IsNotExist(errors.Cause(err)) ***REMOVED***
			e := errors.Wrapf(err, "driver %q failed to remove root filesystem for %s", daemon.GraphDriverName(container.OS), container.ID)
			container.SetRemovalError(e)
			return e
		***REMOVED***
	***REMOVED***

	if err := system.EnsureRemoveAll(container.Root); err != nil ***REMOVED***
		e := errors.Wrapf(err, "unable to remove filesystem for %s", container.ID)
		container.SetRemovalError(e)
		return e
	***REMOVED***

	linkNames := daemon.linkIndex.delete(container)
	selinuxFreeLxcContexts(container.ProcessLabel)
	daemon.idIndex.Delete(container.ID)
	daemon.containers.Delete(container.ID)
	daemon.containersReplica.Delete(container)
	if e := daemon.removeMountPoints(container, removeVolume); e != nil ***REMOVED***
		logrus.Error(e)
	***REMOVED***
	for _, name := range linkNames ***REMOVED***
		daemon.releaseName(name)
	***REMOVED***
	container.SetRemoved()
	stateCtr.del(container.ID)

	daemon.LogContainerEvent(container, "destroy")
	return nil
***REMOVED***

// VolumeRm removes the volume with the given name.
// If the volume is referenced by a container it is not removed
// This is called directly from the Engine API
func (daemon *Daemon) VolumeRm(name string, force bool) error ***REMOVED***
	v, err := daemon.volumes.Get(name)
	if err != nil ***REMOVED***
		if force && volumestore.IsNotExist(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	err = daemon.volumeRm(v)
	if err != nil && volumestore.IsInUse(err) ***REMOVED***
		return errdefs.Conflict(err)
	***REMOVED***

	if err == nil || force ***REMOVED***
		daemon.volumes.Purge(name)
		return nil
	***REMOVED***
	return err
***REMOVED***

func (daemon *Daemon) volumeRm(v volume.Volume) error ***REMOVED***
	if err := daemon.volumes.Remove(v); err != nil ***REMOVED***
		return errors.Wrap(err, "unable to remove volume")
	***REMOVED***
	daemon.LogVolumeEvent(v.Name(), "destroy", map[string]string***REMOVED***"driver": v.DriverName()***REMOVED***)
	return nil
***REMOVED***
