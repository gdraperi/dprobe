package libnetwork

import "github.com/docker/libnetwork/types"

const libnGWNetwork = "docker_gwbridge"

func getPlatformOption() EndpointOption ***REMOVED***
	return nil
***REMOVED***

func (c *controller) createGWNetwork() (Network, error) ***REMOVED***
	return nil, types.NotImplementedErrorf("default gateway functionality is not implemented in freebsd")
***REMOVED***
