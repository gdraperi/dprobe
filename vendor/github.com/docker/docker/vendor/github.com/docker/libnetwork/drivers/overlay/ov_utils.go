package overlay

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/docker/libnetwork/netutils"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/osl"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

var soTimeout = ns.NetlinkSocketsTimeout

func validateID(nid, eid string) error ***REMOVED***
	if nid == "" ***REMOVED***
		return fmt.Errorf("invalid network id")
	***REMOVED***

	if eid == "" ***REMOVED***
		return fmt.Errorf("invalid endpoint id")
	***REMOVED***

	return nil
***REMOVED***

func createVethPair() (string, string, error) ***REMOVED***
	defer osl.InitOSContext()()
	nlh := ns.NlHandle()

	// Generate a name for what will be the host side pipe interface
	name1, err := netutils.GenerateIfaceName(nlh, vethPrefix, vethLen)
	if err != nil ***REMOVED***
		return "", "", fmt.Errorf("error generating veth name1: %v", err)
	***REMOVED***

	// Generate a name for what will be the sandbox side pipe interface
	name2, err := netutils.GenerateIfaceName(nlh, vethPrefix, vethLen)
	if err != nil ***REMOVED***
		return "", "", fmt.Errorf("error generating veth name2: %v", err)
	***REMOVED***

	// Generate and add the interface pipe host <-> sandbox
	veth := &netlink.Veth***REMOVED***
		LinkAttrs: netlink.LinkAttrs***REMOVED***Name: name1, TxQLen: 0***REMOVED***,
		PeerName:  name2***REMOVED***
	if err := nlh.LinkAdd(veth); err != nil ***REMOVED***
		return "", "", fmt.Errorf("error creating veth pair: %v", err)
	***REMOVED***

	return name1, name2, nil
***REMOVED***

func createVxlan(name string, vni uint32, mtu int) error ***REMOVED***
	defer osl.InitOSContext()()

	vxlan := &netlink.Vxlan***REMOVED***
		LinkAttrs: netlink.LinkAttrs***REMOVED***Name: name, MTU: mtu***REMOVED***,
		VxlanId:   int(vni),
		Learning:  true,
		Port:      vxlanPort,
		Proxy:     true,
		L3miss:    true,
		L2miss:    true,
	***REMOVED***

	if err := ns.NlHandle().LinkAdd(vxlan); err != nil ***REMOVED***
		return fmt.Errorf("error creating vxlan interface: %v", err)
	***REMOVED***

	return nil
***REMOVED***

func deleteInterfaceBySubnet(brPrefix string, s *subnet) error ***REMOVED***
	defer osl.InitOSContext()()

	nlh := ns.NlHandle()
	links, err := nlh.LinkList()
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to list interfaces while deleting bridge interface by subnet: %v", err)
	***REMOVED***

	for _, l := range links ***REMOVED***
		name := l.Attrs().Name
		if _, ok := l.(*netlink.Bridge); ok && strings.HasPrefix(name, brPrefix) ***REMOVED***
			addrList, err := nlh.AddrList(l, netlink.FAMILY_V4)
			if err != nil ***REMOVED***
				logrus.Errorf("error getting AddressList for bridge %s", name)
				continue
			***REMOVED***
			for _, addr := range addrList ***REMOVED***
				if netutils.NetworkOverlaps(addr.IPNet, s.subnetIP) ***REMOVED***
					err = nlh.LinkDel(l)
					if err != nil ***REMOVED***
						logrus.Errorf("error deleting bridge (%s) with subnet %v: %v", name, addr.IPNet, err)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil

***REMOVED***

func deleteInterface(name string) error ***REMOVED***
	defer osl.InitOSContext()()

	link, err := ns.NlHandle().LinkByName(name)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to find interface with name %s: %v", name, err)
	***REMOVED***

	if err := ns.NlHandle().LinkDel(link); err != nil ***REMOVED***
		return fmt.Errorf("error deleting interface with name %s: %v", name, err)
	***REMOVED***

	return nil
***REMOVED***

func deleteVxlanByVNI(path string, vni uint32) error ***REMOVED***
	defer osl.InitOSContext()()

	nlh := ns.NlHandle()
	if path != "" ***REMOVED***
		ns, err := netns.GetFromPath(path)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to get ns handle for %s: %v", path, err)
		***REMOVED***
		defer ns.Close()

		nlh, err = netlink.NewHandleAt(ns, syscall.NETLINK_ROUTE)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to get netlink handle for ns %s: %v", path, err)
		***REMOVED***
		defer nlh.Delete()
		err = nlh.SetSocketTimeout(soTimeout)
		if err != nil ***REMOVED***
			logrus.Warnf("Failed to set the timeout on the netlink handle sockets for vxlan deletion: %v", err)
		***REMOVED***
	***REMOVED***

	links, err := nlh.LinkList()
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to list interfaces while deleting vxlan interface by vni: %v", err)
	***REMOVED***

	for _, l := range links ***REMOVED***
		if l.Type() == "vxlan" && (vni == 0 || l.(*netlink.Vxlan).VxlanId == int(vni)) ***REMOVED***
			err = nlh.LinkDel(l)
			if err != nil ***REMOVED***
				return fmt.Errorf("error deleting vxlan interface with id %d: %v", vni, err)
			***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	return fmt.Errorf("could not find a vxlan interface to delete with id %d", vni)
***REMOVED***
