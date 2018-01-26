package macvlan

import (
	"fmt"
	"net"

	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/netutils"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/osl"
	"github.com/sirupsen/logrus"
)

// Join method is invoked when a Sandbox is attached to an endpoint.
func (d *driver) Join(nid, eid string, sboxKey string, jinfo driverapi.JoinInfo, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	defer osl.InitOSContext()()
	n, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	endpoint := n.endpoint(eid)
	if endpoint == nil ***REMOVED***
		return fmt.Errorf("could not find endpoint with id %s", eid)
	***REMOVED***
	// generate a name for the iface that will be renamed to eth0 in the sbox
	containerIfName, err := netutils.GenerateIfaceName(ns.NlHandle(), vethPrefix, vethLen)
	if err != nil ***REMOVED***
		return fmt.Errorf("error generating an interface name: %s", err)
	***REMOVED***
	// create the netlink macvlan interface
	vethName, err := createMacVlan(containerIfName, n.config.Parent, n.config.MacvlanMode)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// bind the generated iface name to the endpoint
	endpoint.srcName = vethName
	ep := n.endpoint(eid)
	if ep == nil ***REMOVED***
		return fmt.Errorf("could not find endpoint with id %s", eid)
	***REMOVED***
	// parse and match the endpoint address with the available v4 subnets
	if len(n.config.Ipv4Subnets) > 0 ***REMOVED***
		s := n.getSubnetforIPv4(ep.addr)
		if s == nil ***REMOVED***
			return fmt.Errorf("could not find a valid ipv4 subnet for endpoint %s", eid)
		***REMOVED***
		v4gw, _, err := net.ParseCIDR(s.GwIP)
		if err != nil ***REMOVED***
			return fmt.Errorf("gatway %s is not a valid ipv4 address: %v", s.GwIP, err)
		***REMOVED***
		err = jinfo.SetGateway(v4gw)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		logrus.Debugf("Macvlan Endpoint Joined with IPv4_Addr: %s, Gateway: %s, MacVlan_Mode: %s, Parent: %s",
			ep.addr.IP.String(), v4gw.String(), n.config.MacvlanMode, n.config.Parent)
	***REMOVED***
	// parse and match the endpoint address with the available v6 subnets
	if len(n.config.Ipv6Subnets) > 0 ***REMOVED***
		s := n.getSubnetforIPv6(ep.addrv6)
		if s == nil ***REMOVED***
			return fmt.Errorf("could not find a valid ipv6 subnet for endpoint %s", eid)
		***REMOVED***
		v6gw, _, err := net.ParseCIDR(s.GwIP)
		if err != nil ***REMOVED***
			return fmt.Errorf("gatway %s is not a valid ipv6 address: %v", s.GwIP, err)
		***REMOVED***
		err = jinfo.SetGatewayIPv6(v6gw)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		logrus.Debugf("Macvlan Endpoint Joined with IPv6_Addr: %s Gateway: %s MacVlan_Mode: %s, Parent: %s",
			ep.addrv6.IP.String(), v6gw.String(), n.config.MacvlanMode, n.config.Parent)
	***REMOVED***
	iNames := jinfo.InterfaceName()
	err = iNames.SetNames(vethName, containerVethPrefix)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := d.storeUpdate(ep); err != nil ***REMOVED***
		return fmt.Errorf("failed to save macvlan endpoint %s to store: %v", ep.id[0:7], err)
	***REMOVED***
	return nil
***REMOVED***

// Leave method is invoked when a Sandbox detaches from an endpoint.
func (d *driver) Leave(nid, eid string) error ***REMOVED***
	defer osl.InitOSContext()()
	network, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	endpoint, err := network.getEndpoint(eid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if endpoint == nil ***REMOVED***
		return fmt.Errorf("could not find endpoint with id %s", eid)
	***REMOVED***

	return nil
***REMOVED***

// getSubnetforIP returns the ipv4 subnet to which the given IP belongs
func (n *network) getSubnetforIPv4(ip *net.IPNet) *ipv4Subnet ***REMOVED***
	for _, s := range n.config.Ipv4Subnets ***REMOVED***
		_, snet, err := net.ParseCIDR(s.SubnetIP)
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		// first check if the mask lengths are the same
		i, _ := snet.Mask.Size()
		j, _ := ip.Mask.Size()
		if i != j ***REMOVED***
			continue
		***REMOVED***
		if snet.Contains(ip.IP) ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// getSubnetforIPv6 returns the ipv6 subnet to which the given IP belongs
func (n *network) getSubnetforIPv6(ip *net.IPNet) *ipv6Subnet ***REMOVED***
	for _, s := range n.config.Ipv6Subnets ***REMOVED***
		_, snet, err := net.ParseCIDR(s.SubnetIP)
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		// first check if the mask lengths are the same
		i, _ := snet.Mask.Size()
		j, _ := ip.Mask.Size()
		if i != j ***REMOVED***
			continue
		***REMOVED***
		if snet.Contains(ip.IP) ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
