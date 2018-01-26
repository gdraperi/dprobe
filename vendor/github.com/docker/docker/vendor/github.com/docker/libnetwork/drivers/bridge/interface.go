package bridge

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

const (
	// DefaultBridgeName is the default name for the bridge interface managed
	// by the driver when unspecified by the caller.
	DefaultBridgeName = "docker0"
)

// Interface models the bridge network device.
type bridgeInterface struct ***REMOVED***
	Link        netlink.Link
	bridgeIPv4  *net.IPNet
	bridgeIPv6  *net.IPNet
	gatewayIPv4 net.IP
	gatewayIPv6 net.IP
	nlh         *netlink.Handle
***REMOVED***

// newInterface creates a new bridge interface structure. It attempts to find
// an already existing device identified by the configuration BridgeName field,
// or the default bridge name when unspecified, but doesn't attempt to create
// one when missing
func newInterface(nlh *netlink.Handle, config *networkConfiguration) (*bridgeInterface, error) ***REMOVED***
	var err error
	i := &bridgeInterface***REMOVED***nlh: nlh***REMOVED***

	// Initialize the bridge name to the default if unspecified.
	if config.BridgeName == "" ***REMOVED***
		config.BridgeName = DefaultBridgeName
	***REMOVED***

	// Attempt to find an existing bridge named with the specified name.
	i.Link, err = nlh.LinkByName(config.BridgeName)
	if err != nil ***REMOVED***
		logrus.Debugf("Did not find any interface with name %s: %v", config.BridgeName, err)
	***REMOVED*** else if _, ok := i.Link.(*netlink.Bridge); !ok ***REMOVED***
		return nil, fmt.Errorf("existing interface %s is not a bridge", i.Link.Attrs().Name)
	***REMOVED***
	return i, nil
***REMOVED***

// exists indicates if the existing bridge interface exists on the system.
func (i *bridgeInterface) exists() bool ***REMOVED***
	return i.Link != nil
***REMOVED***

// addresses returns all IPv4 addresses and all IPv6 addresses for the bridge interface.
func (i *bridgeInterface) addresses() ([]netlink.Addr, []netlink.Addr, error) ***REMOVED***
	v4addr, err := i.nlh.AddrList(i.Link, netlink.FAMILY_V4)
	if err != nil ***REMOVED***
		return nil, nil, fmt.Errorf("Failed to retrieve V4 addresses: %v", err)
	***REMOVED***

	v6addr, err := i.nlh.AddrList(i.Link, netlink.FAMILY_V6)
	if err != nil ***REMOVED***
		return nil, nil, fmt.Errorf("Failed to retrieve V6 addresses: %v", err)
	***REMOVED***

	if len(v4addr) == 0 ***REMOVED***
		return nil, v6addr, nil
	***REMOVED***
	return v4addr, v6addr, nil
***REMOVED***

func (i *bridgeInterface) programIPv6Address() error ***REMOVED***
	_, nlAddressList, err := i.addresses()
	if err != nil ***REMOVED***
		return &IPv6AddrAddError***REMOVED***IP: i.bridgeIPv6, Err: fmt.Errorf("failed to retrieve address list: %v", err)***REMOVED***
	***REMOVED***
	nlAddr := netlink.Addr***REMOVED***IPNet: i.bridgeIPv6***REMOVED***
	if findIPv6Address(nlAddr, nlAddressList) ***REMOVED***
		return nil
	***REMOVED***
	if err := i.nlh.AddrAdd(i.Link, &nlAddr); err != nil ***REMOVED***
		return &IPv6AddrAddError***REMOVED***IP: i.bridgeIPv6, Err: err***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
