package libnetwork

import (
	"github.com/docker/libnetwork/drivers/null"
	"github.com/docker/libnetwork/drivers/remote"
)

func getInitializers(experimental bool) []initializer ***REMOVED***
	return []initializer***REMOVED***
		***REMOVED***null.Init, "null"***REMOVED***,
		***REMOVED***remote.Init, "remote"***REMOVED***,
	***REMOVED***
***REMOVED***
