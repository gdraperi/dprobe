package daemon

import (
	swarmtypes "github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
)

// SetContainerSecretReferences sets the container secret references needed
func (daemon *Daemon) SetContainerSecretReferences(name string, refs []*swarmtypes.SecretReference) error ***REMOVED***
	if !secretsSupported() && len(refs) > 0 ***REMOVED***
		logrus.Warn("secrets are not supported on this platform")
		return nil
	***REMOVED***

	c, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.SecretReferences = refs

	return nil
***REMOVED***
