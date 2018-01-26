package bridge

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"

	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func selectIPv4Address(addresses []netlink.Addr, selector *net.IPNet) (netlink.Addr, error) ***REMOVED***
	if len(addresses) == 0 ***REMOVED***
		return netlink.Addr***REMOVED******REMOVED***, errors.New("unable to select an address as the address pool is empty")
	***REMOVED***
	if selector != nil ***REMOVED***
		for _, addr := range addresses ***REMOVED***
			if selector.Contains(addr.IP) ***REMOVED***
				return addr, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return addresses[0], nil
***REMOVED***

func setupBridgeIPv4(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
	addrv4List, _, err := i.addresses()
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to retrieve bridge interface addresses: %v", err)
	***REMOVED***

	addrv4, _ := selectIPv4Address(addrv4List, config.AddressIPv4)

	if !types.CompareIPNet(addrv4.IPNet, config.AddressIPv4) ***REMOVED***
		if addrv4.IPNet != nil ***REMOVED***
			if err := i.nlh.AddrDel(i.Link, &addrv4); err != nil ***REMOVED***
				return fmt.Errorf("failed to remove current ip address from bridge: %v", err)
			***REMOVED***
		***REMOVED***
		logrus.Debugf("Assigning address to bridge interface %s: %s", config.BridgeName, config.AddressIPv4)
		if err := i.nlh.AddrAdd(i.Link, &netlink.Addr***REMOVED***IPNet: config.AddressIPv4***REMOVED***); err != nil ***REMOVED***
			return &IPv4AddrAddError***REMOVED***IP: config.AddressIPv4, Err: err***REMOVED***
		***REMOVED***
	***REMOVED***

	// Store bridge network and default gateway
	i.bridgeIPv4 = config.AddressIPv4
	i.gatewayIPv4 = config.AddressIPv4.IP

	return nil
***REMOVED***

func setupGatewayIPv4(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
	if !i.bridgeIPv4.Contains(config.DefaultGatewayIPv4) ***REMOVED***
		return &ErrInvalidGateway***REMOVED******REMOVED***
	***REMOVED***

	// Store requested default gateway
	i.gatewayIPv4 = config.DefaultGatewayIPv4

	return nil
***REMOVED***

func setupLoopbackAdressesRouting(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
	sysPath := filepath.Join("/proc/sys/net/ipv4/conf", config.BridgeName, "route_localnet")
	ipv4LoRoutingData, err := ioutil.ReadFile(sysPath)
	if err != nil ***REMOVED***
		return fmt.Errorf("Cannot read IPv4 local routing setup: %v", err)
	***REMOVED***
	// Enable loopback adresses routing only if it isn't already enabled
	if ipv4LoRoutingData[0] != '1' ***REMOVED***
		if err := ioutil.WriteFile(sysPath, []byte***REMOVED***'1', '\n'***REMOVED***, 0644); err != nil ***REMOVED***
			return fmt.Errorf("Unable to enable local routing for hairpin mode: %v", err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
