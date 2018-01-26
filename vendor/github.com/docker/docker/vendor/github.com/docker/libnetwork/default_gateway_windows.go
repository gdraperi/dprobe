package libnetwork

import (
	windriver "github.com/docker/libnetwork/drivers/windows"
	"github.com/docker/libnetwork/options"
	"github.com/docker/libnetwork/types"
)

const libnGWNetwork = "nat"

func getPlatformOption() EndpointOption ***REMOVED***

	epOption := options.Generic***REMOVED***
		windriver.DisableICC: true,
		windriver.DisableDNS: true,
	***REMOVED***
	return EndpointOptionGeneric(epOption)
***REMOVED***

func (c *controller) createGWNetwork() (Network, error) ***REMOVED***
	return nil, types.NotImplementedErrorf("default gateway functionality is not implemented in windows")
***REMOVED***
