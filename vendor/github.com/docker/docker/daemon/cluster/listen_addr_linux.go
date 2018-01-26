package cluster

import (
	"net"

	"github.com/vishvananda/netlink"
)

func (c *Cluster) resolveSystemAddr() (net.IP, error) ***REMOVED***
	// Use the system's only device IP address, or fail if there are
	// multiple addresses to choose from.
	interfaces, err := netlink.LinkList()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var (
		systemAddr      net.IP
		systemInterface string
		deviceFound     bool
	)

	for _, intf := range interfaces ***REMOVED***
		// Skip non device or inactive interfaces
		if intf.Type() != "device" || intf.Attrs().Flags&net.FlagUp == 0 ***REMOVED***
			continue
		***REMOVED***

		addrs, err := netlink.AddrList(intf, netlink.FAMILY_ALL)
		if err != nil ***REMOVED***
			continue
		***REMOVED***

		var interfaceAddr4, interfaceAddr6 net.IP

		for _, addr := range addrs ***REMOVED***
			ipAddr := addr.IPNet.IP

			// Skip loopback and link-local addresses
			if !ipAddr.IsGlobalUnicast() ***REMOVED***
				continue
			***REMOVED***

			// At least one non-loopback device is found and it is administratively up
			deviceFound = true

			if ipAddr.To4() != nil ***REMOVED***
				if interfaceAddr4 != nil ***REMOVED***
					return nil, errMultipleIPs(intf.Attrs().Name, intf.Attrs().Name, interfaceAddr4, ipAddr)
				***REMOVED***
				interfaceAddr4 = ipAddr
			***REMOVED*** else ***REMOVED***
				if interfaceAddr6 != nil ***REMOVED***
					return nil, errMultipleIPs(intf.Attrs().Name, intf.Attrs().Name, interfaceAddr6, ipAddr)
				***REMOVED***
				interfaceAddr6 = ipAddr
			***REMOVED***
		***REMOVED***

		// In the case that this interface has exactly one IPv4 address
		// and exactly one IPv6 address, favor IPv4 over IPv6.
		if interfaceAddr4 != nil ***REMOVED***
			if systemAddr != nil ***REMOVED***
				return nil, errMultipleIPs(systemInterface, intf.Attrs().Name, systemAddr, interfaceAddr4)
			***REMOVED***
			systemAddr = interfaceAddr4
			systemInterface = intf.Attrs().Name
		***REMOVED*** else if interfaceAddr6 != nil ***REMOVED***
			if systemAddr != nil ***REMOVED***
				return nil, errMultipleIPs(systemInterface, intf.Attrs().Name, systemAddr, interfaceAddr6)
			***REMOVED***
			systemAddr = interfaceAddr6
			systemInterface = intf.Attrs().Name
		***REMOVED***
	***REMOVED***

	if systemAddr == nil ***REMOVED***
		if !deviceFound ***REMOVED***
			// If no non-loopback device type interface is found,
			// fall back to the regular auto-detection mechanism.
			// This is to cover the case where docker is running
			// inside a container (eths are in fact veths).
			return c.resolveSystemAddrViaSubnetCheck()
		***REMOVED***
		return nil, errNoIP
	***REMOVED***

	return systemAddr, nil
***REMOVED***
