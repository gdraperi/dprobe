package bridge

import (
	"fmt"
	"strings"

	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func setupVerifyAndReconcile(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
	// Fetch a slice of IPv4 addresses and a slice of IPv6 addresses from the bridge.
	addrsv4, addrsv6, err := i.addresses()
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to verify ip addresses: %v", err)
	***REMOVED***

	addrv4, _ := selectIPv4Address(addrsv4, config.AddressIPv4)

	// Verify that the bridge does have an IPv4 address.
	if addrv4.IPNet == nil ***REMOVED***
		return &ErrNoIPAddr***REMOVED******REMOVED***
	***REMOVED***

	// Verify that the bridge IPv4 address matches the requested configuration.
	if config.AddressIPv4 != nil && !addrv4.IP.Equal(config.AddressIPv4.IP) ***REMOVED***
		return &IPv4AddrNoMatchError***REMOVED***IP: addrv4.IP, CfgIP: config.AddressIPv4.IP***REMOVED***
	***REMOVED***

	// Verify that one of the bridge IPv6 addresses matches the requested
	// configuration.
	if config.EnableIPv6 && !findIPv6Address(netlink.Addr***REMOVED***IPNet: bridgeIPv6***REMOVED***, addrsv6) ***REMOVED***
		return (*IPv6AddrNoMatchError)(bridgeIPv6)
	***REMOVED***

	// Release any residual IPv6 address that might be there because of older daemon instances
	for _, addrv6 := range addrsv6 ***REMOVED***
		if addrv6.IP.IsGlobalUnicast() && !types.CompareIPNet(addrv6.IPNet, i.bridgeIPv6) ***REMOVED***
			if err := i.nlh.AddrDel(i.Link, &addrv6); err != nil ***REMOVED***
				logrus.Warnf("Failed to remove residual IPv6 address %s from bridge: %v", addrv6.IPNet, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func findIPv6Address(addr netlink.Addr, addresses []netlink.Addr) bool ***REMOVED***
	for _, addrv6 := range addresses ***REMOVED***
		if addrv6.String() == addr.String() ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func bridgeInterfaceExists(name string) (bool, error) ***REMOVED***
	nlh := ns.NlHandle()
	link, err := nlh.LinkByName(name)
	if err != nil ***REMOVED***
		if strings.Contains(err.Error(), "Link not found") ***REMOVED***
			return false, nil
		***REMOVED***
		return false, fmt.Errorf("failed to check bridge interface existence: %v", err)
	***REMOVED***

	if link.Type() == "bridge" ***REMOVED***
		return true, nil
	***REMOVED***
	return false, fmt.Errorf("existing interface %s is not a bridge", name)
***REMOVED***
