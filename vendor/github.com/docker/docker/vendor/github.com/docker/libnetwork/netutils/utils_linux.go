// +build linux
// Network utility functions.

package netutils

import (
	"fmt"
	"net"
	"strings"

	"github.com/docker/libnetwork/ipamutils"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/osl"
	"github.com/docker/libnetwork/resolvconf"
	"github.com/docker/libnetwork/types"
	"github.com/vishvananda/netlink"
)

var (
	networkGetRoutesFct func(netlink.Link, int) ([]netlink.Route, error)
)

// CheckRouteOverlaps checks whether the passed network overlaps with any existing routes
func CheckRouteOverlaps(toCheck *net.IPNet) error ***REMOVED***
	if networkGetRoutesFct == nil ***REMOVED***
		networkGetRoutesFct = ns.NlHandle().RouteList
	***REMOVED***
	networks, err := networkGetRoutesFct(nil, netlink.FAMILY_V4)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, network := range networks ***REMOVED***
		if network.Dst != nil && NetworkOverlaps(toCheck, network.Dst) ***REMOVED***
			return ErrNetworkOverlaps
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// GenerateIfaceName returns an interface name using the passed in
// prefix and the length of random bytes. The api ensures that the
// there are is no interface which exists with that name.
func GenerateIfaceName(nlh *netlink.Handle, prefix string, len int) (string, error) ***REMOVED***
	linkByName := netlink.LinkByName
	if nlh != nil ***REMOVED***
		linkByName = nlh.LinkByName
	***REMOVED***
	for i := 0; i < 3; i++ ***REMOVED***
		name, err := GenerateRandomName(prefix, len)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		_, err = linkByName(name)
		if err != nil ***REMOVED***
			if strings.Contains(err.Error(), "not found") ***REMOVED***
				return name, nil
			***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***
	return "", types.InternalErrorf("could not generate interface name")
***REMOVED***

// ElectInterfaceAddresses looks for an interface on the OS with the
// specified name and returns returns all its IPv4 and IPv6 addresses in CIDR notation.
// If a failure in retrieving the addresses or no IPv4 address is found, an error is returned.
// If the interface does not exist, it chooses from a predefined
// list the first IPv4 address which does not conflict with other
// interfaces on the system.
func ElectInterfaceAddresses(name string) ([]*net.IPNet, []*net.IPNet, error) ***REMOVED***
	var (
		v4Nets []*net.IPNet
		v6Nets []*net.IPNet
	)

	defer osl.InitOSContext()()

	link, _ := ns.NlHandle().LinkByName(name)
	if link != nil ***REMOVED***
		v4addr, err := ns.NlHandle().AddrList(link, netlink.FAMILY_V4)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		v6addr, err := ns.NlHandle().AddrList(link, netlink.FAMILY_V6)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		for _, nlAddr := range v4addr ***REMOVED***
			v4Nets = append(v4Nets, nlAddr.IPNet)
		***REMOVED***
		for _, nlAddr := range v6addr ***REMOVED***
			v6Nets = append(v6Nets, nlAddr.IPNet)
		***REMOVED***
	***REMOVED***

	if link == nil || len(v4Nets) == 0 ***REMOVED***
		// Choose from predefined broad networks
		v4Net, err := FindAvailableNetwork(ipamutils.PredefinedBroadNetworks)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		v4Nets = append(v4Nets, v4Net)
	***REMOVED***

	return v4Nets, v6Nets, nil
***REMOVED***

// FindAvailableNetwork returns a network from the passed list which does not
// overlap with existing interfaces in the system
func FindAvailableNetwork(list []*net.IPNet) (*net.IPNet, error) ***REMOVED***
	// We don't check for an error here, because we don't really care if we
	// can't read /etc/resolv.conf. So instead we skip the append if resolvConf
	// is nil. It either doesn't exist, or we can't read it for some reason.
	var nameservers []string
	if rc, err := resolvconf.Get(); err == nil ***REMOVED***
		nameservers = resolvconf.GetNameserversAsCIDR(rc.Content)
	***REMOVED***
	for _, nw := range list ***REMOVED***
		if err := CheckNameserverOverlaps(nameservers, nw); err == nil ***REMOVED***
			if err := CheckRouteOverlaps(nw); err == nil ***REMOVED***
				return nw, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil, fmt.Errorf("no available network")
***REMOVED***
