// +build !linux

package cluster

import "net"

func (c *Cluster) resolveSystemAddr() (net.IP, error) ***REMOVED***
	return c.resolveSystemAddrViaSubnetCheck()
***REMOVED***
