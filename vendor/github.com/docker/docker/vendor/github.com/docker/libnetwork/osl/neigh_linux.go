package osl

import (
	"bytes"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

// NeighborSearchError indicates that the neighbor is already present
type NeighborSearchError struct ***REMOVED***
	ip      net.IP
	mac     net.HardwareAddr
	present bool
***REMOVED***

func (n NeighborSearchError) Error() string ***REMOVED***
	return fmt.Sprintf("Search neighbor failed for IP %v, mac %v, present in db:%t", n.ip, n.mac, n.present)
***REMOVED***

// NeighOption is a function option type to set interface options
type NeighOption func(nh *neigh)

type neigh struct ***REMOVED***
	dstIP    net.IP
	dstMac   net.HardwareAddr
	linkName string
	linkDst  string
	family   int
***REMOVED***

func (n *networkNamespace) findNeighbor(dstIP net.IP, dstMac net.HardwareAddr) *neigh ***REMOVED***
	n.Lock()
	defer n.Unlock()

	for _, nh := range n.neighbors ***REMOVED***
		if nh.dstIP.Equal(dstIP) && bytes.Equal(nh.dstMac, dstMac) ***REMOVED***
			return nh
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (n *networkNamespace) DeleteNeighbor(dstIP net.IP, dstMac net.HardwareAddr, osDelete bool) error ***REMOVED***
	var (
		iface netlink.Link
		err   error
	)

	nh := n.findNeighbor(dstIP, dstMac)
	if nh == nil ***REMOVED***
		return NeighborSearchError***REMOVED***dstIP, dstMac, false***REMOVED***
	***REMOVED***

	if osDelete ***REMOVED***
		n.Lock()
		nlh := n.nlHandle
		n.Unlock()

		if nh.linkDst != "" ***REMOVED***
			iface, err = nlh.LinkByName(nh.linkDst)
			if err != nil ***REMOVED***
				return fmt.Errorf("could not find interface with destination name %s: %v",
					nh.linkDst, err)
			***REMOVED***
		***REMOVED***

		nlnh := &netlink.Neigh***REMOVED***
			IP:     dstIP,
			State:  netlink.NUD_PERMANENT,
			Family: nh.family,
		***REMOVED***

		if nlnh.Family > 0 ***REMOVED***
			nlnh.HardwareAddr = dstMac
			nlnh.Flags = netlink.NTF_SELF
		***REMOVED***

		if nh.linkDst != "" ***REMOVED***
			nlnh.LinkIndex = iface.Attrs().Index
		***REMOVED***

		// If the kernel deletion fails for the neighbor entry still remote it
		// from the namespace cache. Otherwise if the neighbor moves back to the
		// same host again, kernel update can fail.
		if err := nlh.NeighDel(nlnh); err != nil ***REMOVED***
			logrus.Warnf("Deleting neighbor IP %s, mac %s failed, %v", dstIP, dstMac, err)
		***REMOVED***

		// Delete the dynamic entry in the bridge
		if nlnh.Family > 0 ***REMOVED***
			nlnh := &netlink.Neigh***REMOVED***
				IP:     dstIP,
				Family: nh.family,
			***REMOVED***

			nlnh.HardwareAddr = dstMac
			nlnh.Flags = netlink.NTF_MASTER
			if nh.linkDst != "" ***REMOVED***
				nlnh.LinkIndex = iface.Attrs().Index
			***REMOVED***
			nlh.NeighDel(nlnh)
		***REMOVED***
	***REMOVED***

	n.Lock()
	for i, nh := range n.neighbors ***REMOVED***
		if nh.dstIP.Equal(dstIP) && bytes.Equal(nh.dstMac, dstMac) ***REMOVED***
			n.neighbors = append(n.neighbors[:i], n.neighbors[i+1:]...)
			break
		***REMOVED***
	***REMOVED***
	n.Unlock()
	logrus.Debugf("Neighbor entry deleted for IP %v, mac %v osDelete:%t", dstIP, dstMac, osDelete)

	return nil
***REMOVED***

func (n *networkNamespace) AddNeighbor(dstIP net.IP, dstMac net.HardwareAddr, force bool, options ...NeighOption) error ***REMOVED***
	var (
		iface                  netlink.Link
		err                    error
		neighborAlreadyPresent bool
	)

	// If the namespace already has the neighbor entry but the AddNeighbor is called
	// because of a miss notification (force flag) program the kernel anyway.
	nh := n.findNeighbor(dstIP, dstMac)
	if nh != nil ***REMOVED***
		neighborAlreadyPresent = true
		logrus.Warnf("Neighbor entry already present for IP %v, mac %v neighbor:%+v forceUpdate:%t", dstIP, dstMac, nh, force)
		if !force ***REMOVED***
			return NeighborSearchError***REMOVED***dstIP, dstMac, true***REMOVED***
		***REMOVED***
	***REMOVED***

	nh = &neigh***REMOVED***
		dstIP:  dstIP,
		dstMac: dstMac,
	***REMOVED***

	nh.processNeighOptions(options...)

	if nh.linkName != "" ***REMOVED***
		nh.linkDst = n.findDst(nh.linkName, false)
		if nh.linkDst == "" ***REMOVED***
			return fmt.Errorf("could not find the interface with name %s", nh.linkName)
		***REMOVED***
	***REMOVED***

	n.Lock()
	nlh := n.nlHandle
	n.Unlock()

	if nh.linkDst != "" ***REMOVED***
		iface, err = nlh.LinkByName(nh.linkDst)
		if err != nil ***REMOVED***
			return fmt.Errorf("could not find interface with destination name %s: %v", nh.linkDst, err)
		***REMOVED***
	***REMOVED***

	nlnh := &netlink.Neigh***REMOVED***
		IP:           dstIP,
		HardwareAddr: dstMac,
		State:        netlink.NUD_PERMANENT,
		Family:       nh.family,
	***REMOVED***

	if nlnh.Family > 0 ***REMOVED***
		nlnh.Flags = netlink.NTF_SELF
	***REMOVED***

	if nh.linkDst != "" ***REMOVED***
		nlnh.LinkIndex = iface.Attrs().Index
	***REMOVED***

	if err := nlh.NeighSet(nlnh); err != nil ***REMOVED***
		return fmt.Errorf("could not add neighbor entry:%+v error:%v", nlnh, err)
	***REMOVED***

	if neighborAlreadyPresent ***REMOVED***
		return nil
	***REMOVED***

	n.Lock()
	n.neighbors = append(n.neighbors, nh)
	n.Unlock()
	logrus.Debugf("Neighbor entry added for IP:%v, mac:%v on ifc:%s", dstIP, dstMac, nh.linkName)

	return nil
***REMOVED***
