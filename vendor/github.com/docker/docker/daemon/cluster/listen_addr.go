package cluster

import (
	"fmt"
	"net"
)

const (
	errNoSuchInterface         configError = "no such interface"
	errNoIP                    configError = "could not find the system's IP address"
	errMustSpecifyListenAddr   configError = "must specify a listening address because the address to advertise is not recognized as a system address, and a system's IP address to use could not be uniquely identified"
	errBadNetworkIdentifier    configError = "must specify a valid IP address or interface name"
	errBadListenAddr           configError = "listen address must be an IP address or network interface (with optional port number)"
	errBadAdvertiseAddr        configError = "advertise address must be a non-zero IP address or network interface (with optional port number)"
	errBadDataPathAddr         configError = "data path address must be a non-zero IP address or network interface (without a port number)"
	errBadDefaultAdvertiseAddr configError = "default advertise address must be a non-zero IP address or network interface (without a port number)"
)

func resolveListenAddr(specifiedAddr string) (string, string, error) ***REMOVED***
	specifiedHost, specifiedPort, err := net.SplitHostPort(specifiedAddr)
	if err != nil ***REMOVED***
		return "", "", fmt.Errorf("could not parse listen address %s", specifiedAddr)
	***REMOVED***
	// Does the host component match any of the interface names on the
	// system? If so, use the address from that interface.
	specifiedIP, err := resolveInputIPAddr(specifiedHost, true)
	if err != nil ***REMOVED***
		if err == errBadNetworkIdentifier ***REMOVED***
			err = errBadListenAddr
		***REMOVED***
		return "", "", err
	***REMOVED***

	return specifiedIP.String(), specifiedPort, nil
***REMOVED***

func (c *Cluster) resolveAdvertiseAddr(advertiseAddr, listenAddrPort string) (string, string, error) ***REMOVED***
	// Approach:
	// - If an advertise address is specified, use that. Resolve the
	//   interface's address if an interface was specified in
	//   advertiseAddr. Fill in the port from listenAddrPort if necessary.
	// - If DefaultAdvertiseAddr is not empty, use that with the port from
	//   listenAddrPort. Resolve the interface's address from
	//   if an interface name was specified in DefaultAdvertiseAddr.
	// - Otherwise, try to autodetect the system's address. Use the port in
	//   listenAddrPort with this address if autodetection succeeds.

	if advertiseAddr != "" ***REMOVED***
		advertiseHost, advertisePort, err := net.SplitHostPort(advertiseAddr)
		if err != nil ***REMOVED***
			// Not a host:port specification
			advertiseHost = advertiseAddr
			advertisePort = listenAddrPort
		***REMOVED***
		// Does the host component match any of the interface names on the
		// system? If so, use the address from that interface.
		advertiseIP, err := resolveInputIPAddr(advertiseHost, false)
		if err != nil ***REMOVED***
			if err == errBadNetworkIdentifier ***REMOVED***
				err = errBadAdvertiseAddr
			***REMOVED***
			return "", "", err
		***REMOVED***

		return advertiseIP.String(), advertisePort, nil
	***REMOVED***

	if c.config.DefaultAdvertiseAddr != "" ***REMOVED***
		// Does the default advertise address component match any of the
		// interface names on the system? If so, use the address from
		// that interface.
		defaultAdvertiseIP, err := resolveInputIPAddr(c.config.DefaultAdvertiseAddr, false)
		if err != nil ***REMOVED***
			if err == errBadNetworkIdentifier ***REMOVED***
				err = errBadDefaultAdvertiseAddr
			***REMOVED***
			return "", "", err
		***REMOVED***

		return defaultAdvertiseIP.String(), listenAddrPort, nil
	***REMOVED***

	systemAddr, err := c.resolveSystemAddr()
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	return systemAddr.String(), listenAddrPort, nil
***REMOVED***

func resolveDataPathAddr(dataPathAddr string) (string, error) ***REMOVED***
	if dataPathAddr == "" ***REMOVED***
		// dataPathAddr is not defined
		return "", nil
	***REMOVED***
	// If a data path flag is specified try to resolve the IP address.
	dataPathIP, err := resolveInputIPAddr(dataPathAddr, false)
	if err != nil ***REMOVED***
		if err == errBadNetworkIdentifier ***REMOVED***
			err = errBadDataPathAddr
		***REMOVED***
		return "", err
	***REMOVED***
	return dataPathIP.String(), nil
***REMOVED***

func resolveInterfaceAddr(specifiedInterface string) (net.IP, error) ***REMOVED***
	// Use a specific interface's IP address.
	intf, err := net.InterfaceByName(specifiedInterface)
	if err != nil ***REMOVED***
		return nil, errNoSuchInterface
	***REMOVED***

	addrs, err := intf.Addrs()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var interfaceAddr4, interfaceAddr6 net.IP

	for _, addr := range addrs ***REMOVED***
		ipAddr, ok := addr.(*net.IPNet)

		if ok ***REMOVED***
			if ipAddr.IP.To4() != nil ***REMOVED***
				// IPv4
				if interfaceAddr4 != nil ***REMOVED***
					return nil, configError(fmt.Sprintf("interface %s has more than one IPv4 address (%s and %s)", specifiedInterface, interfaceAddr4, ipAddr.IP))
				***REMOVED***
				interfaceAddr4 = ipAddr.IP
			***REMOVED*** else ***REMOVED***
				// IPv6
				if interfaceAddr6 != nil ***REMOVED***
					return nil, configError(fmt.Sprintf("interface %s has more than one IPv6 address (%s and %s)", specifiedInterface, interfaceAddr6, ipAddr.IP))
				***REMOVED***
				interfaceAddr6 = ipAddr.IP
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if interfaceAddr4 == nil && interfaceAddr6 == nil ***REMOVED***
		return nil, configError(fmt.Sprintf("interface %s has no usable IPv4 or IPv6 address", specifiedInterface))
	***REMOVED***

	// In the case that there's exactly one IPv4 address
	// and exactly one IPv6 address, favor IPv4 over IPv6.
	if interfaceAddr4 != nil ***REMOVED***
		return interfaceAddr4, nil
	***REMOVED***
	return interfaceAddr6, nil
***REMOVED***

// resolveInputIPAddr tries to resolve the IP address from the string passed as input
// - tries to match the string as an interface name, if so returns the IP address associated with it
// - on failure of previous step tries to parse the string as an IP address itself
//	 if succeeds returns the IP address
func resolveInputIPAddr(input string, isUnspecifiedValid bool) (net.IP, error) ***REMOVED***
	// Try to see if it is an interface name
	interfaceAddr, err := resolveInterfaceAddr(input)
	if err == nil ***REMOVED***
		return interfaceAddr, nil
	***REMOVED***
	// String matched interface but there is a potential ambiguity to be resolved
	if err != errNoSuchInterface ***REMOVED***
		return nil, err
	***REMOVED***

	// String is not an interface check if it is a valid IP
	if ip := net.ParseIP(input); ip != nil && (isUnspecifiedValid || !ip.IsUnspecified()) ***REMOVED***
		return ip, nil
	***REMOVED***

	// Not valid IP found
	return nil, errBadNetworkIdentifier
***REMOVED***

func (c *Cluster) resolveSystemAddrViaSubnetCheck() (net.IP, error) ***REMOVED***
	// Use the system's only IP address, or fail if there are
	// multiple addresses to choose from. Skip interfaces which
	// are managed by docker via subnet check.
	interfaces, err := net.Interfaces()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var systemAddr net.IP
	var systemInterface string

	// List Docker-managed subnets
	v4Subnets, v6Subnets := c.config.NetworkSubnetsProvider.Subnets()

ifaceLoop:
	for _, intf := range interfaces ***REMOVED***
		// Skip inactive interfaces and loopback interfaces
		if (intf.Flags&net.FlagUp == 0) || (intf.Flags&net.FlagLoopback) != 0 ***REMOVED***
			continue
		***REMOVED***

		addrs, err := intf.Addrs()
		if err != nil ***REMOVED***
			continue
		***REMOVED***

		var interfaceAddr4, interfaceAddr6 net.IP

		for _, addr := range addrs ***REMOVED***
			ipAddr, ok := addr.(*net.IPNet)

			// Skip loopback and link-local addresses
			if !ok || !ipAddr.IP.IsGlobalUnicast() ***REMOVED***
				continue
			***REMOVED***

			if ipAddr.IP.To4() != nil ***REMOVED***
				// IPv4

				// Ignore addresses in subnets that are managed by Docker.
				for _, subnet := range v4Subnets ***REMOVED***
					if subnet.Contains(ipAddr.IP) ***REMOVED***
						continue ifaceLoop
					***REMOVED***
				***REMOVED***

				if interfaceAddr4 != nil ***REMOVED***
					return nil, errMultipleIPs(intf.Name, intf.Name, interfaceAddr4, ipAddr.IP)
				***REMOVED***

				interfaceAddr4 = ipAddr.IP
			***REMOVED*** else ***REMOVED***
				// IPv6

				// Ignore addresses in subnets that are managed by Docker.
				for _, subnet := range v6Subnets ***REMOVED***
					if subnet.Contains(ipAddr.IP) ***REMOVED***
						continue ifaceLoop
					***REMOVED***
				***REMOVED***

				if interfaceAddr6 != nil ***REMOVED***
					return nil, errMultipleIPs(intf.Name, intf.Name, interfaceAddr6, ipAddr.IP)
				***REMOVED***

				interfaceAddr6 = ipAddr.IP
			***REMOVED***
		***REMOVED***

		// In the case that this interface has exactly one IPv4 address
		// and exactly one IPv6 address, favor IPv4 over IPv6.
		if interfaceAddr4 != nil ***REMOVED***
			if systemAddr != nil ***REMOVED***
				return nil, errMultipleIPs(systemInterface, intf.Name, systemAddr, interfaceAddr4)
			***REMOVED***
			systemAddr = interfaceAddr4
			systemInterface = intf.Name
		***REMOVED*** else if interfaceAddr6 != nil ***REMOVED***
			if systemAddr != nil ***REMOVED***
				return nil, errMultipleIPs(systemInterface, intf.Name, systemAddr, interfaceAddr6)
			***REMOVED***
			systemAddr = interfaceAddr6
			systemInterface = intf.Name
		***REMOVED***
	***REMOVED***

	if systemAddr == nil ***REMOVED***
		return nil, errNoIP
	***REMOVED***

	return systemAddr, nil
***REMOVED***

func listSystemIPs() []net.IP ***REMOVED***
	interfaces, err := net.Interfaces()
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	var systemAddrs []net.IP

	for _, intf := range interfaces ***REMOVED***
		addrs, err := intf.Addrs()
		if err != nil ***REMOVED***
			continue
		***REMOVED***

		for _, addr := range addrs ***REMOVED***
			ipAddr, ok := addr.(*net.IPNet)

			if ok ***REMOVED***
				systemAddrs = append(systemAddrs, ipAddr.IP)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return systemAddrs
***REMOVED***

func errMultipleIPs(interfaceA, interfaceB string, addrA, addrB net.IP) error ***REMOVED***
	if interfaceA == interfaceB ***REMOVED***
		return configError(fmt.Sprintf("could not choose an IP address to advertise since this system has multiple addresses on interface %s (%s and %s)", interfaceA, addrA, addrB))
	***REMOVED***
	return configError(fmt.Sprintf("could not choose an IP address to advertise since this system has multiple addresses on different interfaces (%s on %s and %s on %s)", addrA, interfaceA, addrB, interfaceB))
***REMOVED***
