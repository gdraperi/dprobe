package daemon

import (
	"strings"

	dockercontainer "github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/libnetwork"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ContainerRename changes the name of a container, using the oldName
// to find the container. An error is returned if newName is already
// reserved.
func (daemon *Daemon) ContainerRename(oldName, newName string) error ***REMOVED***
	var (
		sid string
		sb  libnetwork.Sandbox
	)

	if oldName == "" || newName == "" ***REMOVED***
		return errdefs.InvalidParameter(errors.New("Neither old nor new names may be empty"))
	***REMOVED***

	if newName[0] != '/' ***REMOVED***
		newName = "/" + newName
	***REMOVED***

	container, err := daemon.GetContainer(oldName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	container.Lock()
	defer container.Unlock()

	oldName = container.Name
	oldIsAnonymousEndpoint := container.NetworkSettings.IsAnonymousEndpoint

	if oldName == newName ***REMOVED***
		return errdefs.InvalidParameter(errors.New("Renaming a container with the same name as its current name"))
	***REMOVED***

	links := map[string]*dockercontainer.Container***REMOVED******REMOVED***
	for k, v := range daemon.linkIndex.children(container) ***REMOVED***
		if !strings.HasPrefix(k, oldName) ***REMOVED***
			return errdefs.InvalidParameter(errors.Errorf("Linked container %s does not match parent %s", k, oldName))
		***REMOVED***
		links[strings.TrimPrefix(k, oldName)] = v
	***REMOVED***

	if newName, err = daemon.reserveName(container.ID, newName); err != nil ***REMOVED***
		return errors.Wrap(err, "Error when allocating new name")
	***REMOVED***

	for k, v := range links ***REMOVED***
		daemon.containersReplica.ReserveName(newName+k, v.ID)
		daemon.linkIndex.link(container, v, newName+k)
	***REMOVED***

	container.Name = newName
	container.NetworkSettings.IsAnonymousEndpoint = false

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			container.Name = oldName
			container.NetworkSettings.IsAnonymousEndpoint = oldIsAnonymousEndpoint
			daemon.reserveName(container.ID, oldName)
			for k, v := range links ***REMOVED***
				daemon.containersReplica.ReserveName(oldName+k, v.ID)
				daemon.linkIndex.link(container, v, oldName+k)
				daemon.linkIndex.unlink(newName+k, v, container)
				daemon.containersReplica.ReleaseName(newName + k)
			***REMOVED***
			daemon.releaseName(newName)
		***REMOVED***
	***REMOVED***()

	for k, v := range links ***REMOVED***
		daemon.linkIndex.unlink(oldName+k, v, container)
		daemon.containersReplica.ReleaseName(oldName + k)
	***REMOVED***
	daemon.releaseName(oldName)
	if err = container.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
		return err
	***REMOVED***

	attributes := map[string]string***REMOVED***
		"oldName": oldName,
	***REMOVED***

	if !container.Running ***REMOVED***
		daemon.LogContainerEventWithAttributes(container, "rename", attributes)
		return nil
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			container.Name = oldName
			container.NetworkSettings.IsAnonymousEndpoint = oldIsAnonymousEndpoint
			if e := container.CheckpointTo(daemon.containersReplica); e != nil ***REMOVED***
				logrus.Errorf("%s: Failed in writing to Disk on rename failure: %v", container.ID, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	sid = container.NetworkSettings.SandboxID
	if sid != "" && daemon.netController != nil ***REMOVED***
		sb, err = daemon.netController.SandboxByID(sid)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = sb.Rename(strings.TrimPrefix(container.Name, "/"))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	daemon.LogContainerEventWithAttributes(container, "rename", attributes)
	return nil
***REMOVED***
