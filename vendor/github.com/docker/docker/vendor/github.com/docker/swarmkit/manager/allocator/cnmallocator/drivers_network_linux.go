package cnmallocator

import (
	"github.com/docker/libnetwork/drivers/bridge/brmanager"
	"github.com/docker/libnetwork/drivers/host"
	"github.com/docker/libnetwork/drivers/ipvlan/ivmanager"
	"github.com/docker/libnetwork/drivers/macvlan/mvmanager"
	"github.com/docker/libnetwork/drivers/overlay/ovmanager"
	"github.com/docker/libnetwork/drivers/remote"
	"github.com/docker/swarmkit/manager/allocator/networkallocator"
)

var initializers = []initializer***REMOVED***
	***REMOVED***remote.Init, "remote"***REMOVED***,
	***REMOVED***ovmanager.Init, "overlay"***REMOVED***,
	***REMOVED***mvmanager.Init, "macvlan"***REMOVED***,
	***REMOVED***brmanager.Init, "bridge"***REMOVED***,
	***REMOVED***ivmanager.Init, "ipvlan"***REMOVED***,
	***REMOVED***host.Init, "host"***REMOVED***,
***REMOVED***

// PredefinedNetworks returns the list of predefined network structures
func PredefinedNetworks() []networkallocator.PredefinedNetworkData ***REMOVED***
	return []networkallocator.PredefinedNetworkData***REMOVED***
		***REMOVED***Name: "bridge", Driver: "bridge"***REMOVED***,
		***REMOVED***Name: "host", Driver: "host"***REMOVED***,
	***REMOVED***
***REMOVED***
