package daemon

import (
	swarmtypes "github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
)

// SetContainerConfigReferences sets the container config references needed
func (daemon *Daemon) SetContainerConfigReferences(name string, refs []*swarmtypes.ConfigReference) error ***REMOVED***
	if !configsSupported() && len(refs) > 0 ***REMOVED***
		logrus.Warn("configs are not supported on this platform")
		return nil
	***REMOVED***

	c, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.ConfigReferences = refs

	return nil
***REMOVED***
