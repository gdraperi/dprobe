package runconfig

import (
	"github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
)

// ContainerConfigWrapper is a Config wrapper that holds the container Config (portable)
// and the corresponding HostConfig (non-portable).
type ContainerConfigWrapper struct ***REMOVED***
	*container.Config
	HostConfig       *container.HostConfig          `json:"HostConfig,omitempty"`
	NetworkingConfig *networktypes.NetworkingConfig `json:"NetworkingConfig,omitempty"`
***REMOVED***

// getHostConfig gets the HostConfig of the Config.
func (w *ContainerConfigWrapper) getHostConfig() *container.HostConfig ***REMOVED***
	return w.HostConfig
***REMOVED***
