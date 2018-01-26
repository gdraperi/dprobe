package daemon

import (
	"github.com/docker/swarmkit/agent/exec"
)

// SetContainerDependencyStore sets the dependency store backend for the container
func (daemon *Daemon) SetContainerDependencyStore(name string, store exec.DependencyGetter) error ***REMOVED***
	c, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.DependencyStore = store

	return nil
***REMOVED***
