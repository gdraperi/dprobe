package cnmallocator

import (
	"github.com/docker/libnetwork/drivers/overlay/ovmanager"
	"github.com/docker/libnetwork/drivers/remote"
	"github.com/docker/swarmkit/manager/allocator/networkallocator"
)

var initializers = []initializer***REMOVED***
	***REMOVED***remote.Init, "remote"***REMOVED***,
	***REMOVED***ovmanager.Init, "overlay"***REMOVED***,
***REMOVED***

// PredefinedNetworks returns the list of predefined network structures
func PredefinedNetworks() []networkallocator.PredefinedNetworkData ***REMOVED***
	return nil
***REMOVED***
