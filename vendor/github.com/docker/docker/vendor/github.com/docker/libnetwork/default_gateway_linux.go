package libnetwork

import (
	"fmt"
	"strconv"

	"github.com/docker/libnetwork/drivers/bridge"
)

const libnGWNetwork = "docker_gwbridge"

func getPlatformOption() EndpointOption ***REMOVED***
	return nil
***REMOVED***

func (c *controller) createGWNetwork() (Network, error) ***REMOVED***
	netOption := map[string]string***REMOVED***
		bridge.BridgeName:         libnGWNetwork,
		bridge.EnableICC:          strconv.FormatBool(false),
		bridge.EnableIPMasquerade: strconv.FormatBool(true),
	***REMOVED***

	n, err := c.NewNetwork("bridge", libnGWNetwork, "",
		NetworkOptionDriverOpts(netOption),
		NetworkOptionEnableIPv6(false),
	)

	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error creating external connectivity network: %v", err)
	***REMOVED***
	return n, err
***REMOVED***
