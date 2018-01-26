// +build !windows

package libnetwork

import "fmt"

func (ep *endpoint) DriverInfo() (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	ep, err := ep.retrieveFromStore()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if sb, ok := ep.getSandbox(); ok ***REMOVED***
		if gwep := sb.getEndpointInGWNetwork(); gwep != nil && gwep.ID() != ep.ID() ***REMOVED***
			return gwep.DriverInfo()
		***REMOVED***
	***REMOVED***

	n, err := ep.getNetworkFromStore()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("could not find network in store for driver info: %v", err)
	***REMOVED***

	driver, err := n.driver(true)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to get driver info: %v", err)
	***REMOVED***

	return driver.EndpointOperInfo(n.ID(), ep.ID())
***REMOVED***
