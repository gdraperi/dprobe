package macvlan

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/libnetwork/ns"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

const (
	dummyPrefix      = "dm-" // macvlan prefix for dummy parent interface
	macvlanKernelVer = 3     // minimum macvlan kernel support
	macvlanMajorVer  = 9     // minimum macvlan major kernel support
)

// Create the macvlan slave specifying the source name
func createMacVlan(containerIfName, parent, macvlanMode string) (string, error) ***REMOVED***
	// Set the macvlan mode. Default is bridge mode
	mode, err := setMacVlanMode(macvlanMode)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("Unsupported %s macvlan mode: %v", macvlanMode, err)
	***REMOVED***
	// verify the Docker host interface acting as the macvlan parent iface exists
	if !parentExists(parent) ***REMOVED***
		return "", fmt.Errorf("the requested parent interface %s was not found on the Docker host", parent)
	***REMOVED***
	// Get the link for the master index (Example: the docker host eth iface)
	parentLink, err := ns.NlHandle().LinkByName(parent)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("error occoured looking up the %s parent iface %s error: %s", macvlanType, parent, err)
	***REMOVED***
	// Create a macvlan link
	macvlan := &netlink.Macvlan***REMOVED***
		LinkAttrs: netlink.LinkAttrs***REMOVED***
			Name:        containerIfName,
			ParentIndex: parentLink.Attrs().Index,
		***REMOVED***,
		Mode: mode,
	***REMOVED***
	if err := ns.NlHandle().LinkAdd(macvlan); err != nil ***REMOVED***
		// If a user creates a macvlan and ipvlan on same parent, only one slave iface can be active at a time.
		return "", fmt.Errorf("failed to create the %s port: %v", macvlanType, err)
	***REMOVED***

	return macvlan.Attrs().Name, nil
***REMOVED***

// setMacVlanMode setter for one of the four macvlan port types
func setMacVlanMode(mode string) (netlink.MacvlanMode, error) ***REMOVED***
	switch mode ***REMOVED***
	case modePrivate:
		return netlink.MACVLAN_MODE_PRIVATE, nil
	case modeVepa:
		return netlink.MACVLAN_MODE_VEPA, nil
	case modeBridge:
		return netlink.MACVLAN_MODE_BRIDGE, nil
	case modePassthru:
		return netlink.MACVLAN_MODE_PASSTHRU, nil
	default:
		return 0, fmt.Errorf("unknown macvlan mode: %s", mode)
	***REMOVED***
***REMOVED***

// parentExists checks if the specified interface exists in the default namespace
func parentExists(ifaceStr string) bool ***REMOVED***
	_, err := ns.NlHandle().LinkByName(ifaceStr)
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***

// createVlanLink parses sub-interfaces and vlan id for creation
func createVlanLink(parentName string) error ***REMOVED***
	if strings.Contains(parentName, ".") ***REMOVED***
		parent, vidInt, err := parseVlan(parentName)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// VLAN identifier or VID is a 12-bit field specifying the VLAN to which the frame belongs
		if vidInt > 4094 || vidInt < 1 ***REMOVED***
			return fmt.Errorf("vlan id must be between 1-4094, received: %d", vidInt)
		***REMOVED***
		// get the parent link to attach a vlan subinterface
		parentLink, err := ns.NlHandle().LinkByName(parent)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to find master interface %s on the Docker host: %v", parent, err)
		***REMOVED***
		vlanLink := &netlink.Vlan***REMOVED***
			LinkAttrs: netlink.LinkAttrs***REMOVED***
				Name:        parentName,
				ParentIndex: parentLink.Attrs().Index,
			***REMOVED***,
			VlanId: vidInt,
		***REMOVED***
		// create the subinterface
		if err := ns.NlHandle().LinkAdd(vlanLink); err != nil ***REMOVED***
			return fmt.Errorf("failed to create %s vlan link: %v", vlanLink.Name, err)
		***REMOVED***
		// Bring the new netlink iface up
		if err := ns.NlHandle().LinkSetUp(vlanLink); err != nil ***REMOVED***
			return fmt.Errorf("failed to enable %s the macvlan parent link %v", vlanLink.Name, err)
		***REMOVED***
		logrus.Debugf("Added a vlan tagged netlink subinterface: %s with a vlan id: %d", parentName, vidInt)
		return nil
	***REMOVED***

	return fmt.Errorf("invalid subinterface vlan name %s, example formatting is eth0.10", parentName)
***REMOVED***

// delVlanLink verifies only sub-interfaces with a vlan id get deleted
func delVlanLink(linkName string) error ***REMOVED***
	if strings.Contains(linkName, ".") ***REMOVED***
		_, _, err := parseVlan(linkName)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// delete the vlan subinterface
		vlanLink, err := ns.NlHandle().LinkByName(linkName)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to find interface %s on the Docker host : %v", linkName, err)
		***REMOVED***
		// verify a parent interface isn't being deleted
		if vlanLink.Attrs().ParentIndex == 0 ***REMOVED***
			return fmt.Errorf("interface %s does not appear to be a slave device: %v", linkName, err)
		***REMOVED***
		// delete the macvlan slave device
		if err := ns.NlHandle().LinkDel(vlanLink); err != nil ***REMOVED***
			return fmt.Errorf("failed to delete  %s link: %v", linkName, err)
		***REMOVED***
		logrus.Debugf("Deleted a vlan tagged netlink subinterface: %s", linkName)
	***REMOVED***
	// if the subinterface doesn't parse to iface.vlan_id leave the interface in
	// place since it could be a user specified name not created by the driver.
	return nil
***REMOVED***

// parseVlan parses and verifies a slave interface name: -o parent=eth0.10
func parseVlan(linkName string) (string, int, error) ***REMOVED***
	// parse -o parent=eth0.10
	splitName := strings.Split(linkName, ".")
	if len(splitName) != 2 ***REMOVED***
		return "", 0, fmt.Errorf("required interface name format is: name.vlan_id, ex. eth0.10 for vlan 10, instead received %s", linkName)
	***REMOVED***
	parent, vidStr := splitName[0], splitName[1]
	// validate type and convert vlan id to int
	vidInt, err := strconv.Atoi(vidStr)
	if err != nil ***REMOVED***
		return "", 0, fmt.Errorf("unable to parse a valid vlan id from: %s (ex. eth0.10 for vlan 10)", vidStr)
	***REMOVED***
	// Check if the interface exists
	if !parentExists(parent) ***REMOVED***
		return "", 0, fmt.Errorf("-o parent interface does was not found on the host: %s", parent)
	***REMOVED***

	return parent, vidInt, nil
***REMOVED***

// createDummyLink creates a dummy0 parent link
func createDummyLink(dummyName, truncNetID string) error ***REMOVED***
	// create a parent interface since one was not specified
	parent := &netlink.Dummy***REMOVED***
		LinkAttrs: netlink.LinkAttrs***REMOVED***
			Name: dummyName,
		***REMOVED***,
	***REMOVED***
	if err := ns.NlHandle().LinkAdd(parent); err != nil ***REMOVED***
		return err
	***REMOVED***
	parentDummyLink, err := ns.NlHandle().LinkByName(dummyName)
	if err != nil ***REMOVED***
		return fmt.Errorf("error occoured looking up the %s parent iface %s error: %s", macvlanType, dummyName, err)
	***REMOVED***
	// bring the new netlink iface up
	if err := ns.NlHandle().LinkSetUp(parentDummyLink); err != nil ***REMOVED***
		return fmt.Errorf("failed to enable %s the macvlan parent link: %v", dummyName, err)
	***REMOVED***

	return nil
***REMOVED***

// delDummyLink deletes the link type dummy used when -o parent is not passed
func delDummyLink(linkName string) error ***REMOVED***
	// delete the vlan subinterface
	dummyLink, err := ns.NlHandle().LinkByName(linkName)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to find link %s on the Docker host : %v", linkName, err)
	***REMOVED***
	// verify a parent interface is being deleted
	if dummyLink.Attrs().ParentIndex != 0 ***REMOVED***
		return fmt.Errorf("link %s is not a parent dummy interface", linkName)
	***REMOVED***
	// delete the macvlan dummy device
	if err := ns.NlHandle().LinkDel(dummyLink); err != nil ***REMOVED***
		return fmt.Errorf("failed to delete the dummy %s link: %v", linkName, err)
	***REMOVED***
	logrus.Debugf("Deleted a dummy parent link: %s", linkName)

	return nil
***REMOVED***

// getDummyName returns the name of a dummy parent with truncated net ID and driver prefix
func getDummyName(netID string) string ***REMOVED***
	return dummyPrefix + netID
***REMOVED***
