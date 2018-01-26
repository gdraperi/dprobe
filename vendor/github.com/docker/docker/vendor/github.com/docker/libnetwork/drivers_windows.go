package libnetwork

import (
	"github.com/docker/libnetwork/drivers/null"
	"github.com/docker/libnetwork/drivers/remote"
	"github.com/docker/libnetwork/drivers/windows"
	"github.com/docker/libnetwork/drivers/windows/overlay"
)

func getInitializers(experimental bool) []initializer ***REMOVED***
	return []initializer***REMOVED***
		***REMOVED***null.Init, "null"***REMOVED***,
		***REMOVED***overlay.Init, "overlay"***REMOVED***,
		***REMOVED***remote.Init, "remote"***REMOVED***,
		***REMOVED***windows.GetInit("transparent"), "transparent"***REMOVED***,
		***REMOVED***windows.GetInit("l2bridge"), "l2bridge"***REMOVED***,
		***REMOVED***windows.GetInit("l2tunnel"), "l2tunnel"***REMOVED***,
		***REMOVED***windows.GetInit("nat"), "nat"***REMOVED***,
		***REMOVED***windows.GetInit("ics"), "ics"***REMOVED***,
	***REMOVED***
***REMOVED***
