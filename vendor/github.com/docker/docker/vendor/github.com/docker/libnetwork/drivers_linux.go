package libnetwork

import (
	"github.com/docker/libnetwork/drivers/bridge"
	"github.com/docker/libnetwork/drivers/host"
	"github.com/docker/libnetwork/drivers/macvlan"
	"github.com/docker/libnetwork/drivers/null"
	"github.com/docker/libnetwork/drivers/overlay"
	"github.com/docker/libnetwork/drivers/remote"
)

func getInitializers(experimental bool) []initializer ***REMOVED***
	in := []initializer***REMOVED***
		***REMOVED***bridge.Init, "bridge"***REMOVED***,
		***REMOVED***host.Init, "host"***REMOVED***,
		***REMOVED***macvlan.Init, "macvlan"***REMOVED***,
		***REMOVED***null.Init, "null"***REMOVED***,
		***REMOVED***remote.Init, "remote"***REMOVED***,
		***REMOVED***overlay.Init, "overlay"***REMOVED***,
	***REMOVED***

	if experimental ***REMOVED***
		in = append(in, additionalDrivers()...)
	***REMOVED***
	return in
***REMOVED***
