// +build !linux,!windows

package libnetwork

import (
	"fmt"
	"net"
)

func (c *controller) cleanupServiceBindings(nid string) ***REMOVED***
***REMOVED***

func (c *controller) addServiceBinding(name, sid, nid, eid string, vip net.IP, ingressPorts []*PortConfig, aliases []string, ip net.IP) error ***REMOVED***
	return fmt.Errorf("not supported")
***REMOVED***

func (c *controller) rmServiceBinding(name, sid, nid, eid string, vip net.IP, ingressPorts []*PortConfig, aliases []string, ip net.IP) error ***REMOVED***
	return fmt.Errorf("not supported")
***REMOVED***

func (sb *sandbox) populateLoadbalancers(ep *endpoint) ***REMOVED***
***REMOVED***

func arrangeIngressFilterRule() ***REMOVED***
***REMOVED***
