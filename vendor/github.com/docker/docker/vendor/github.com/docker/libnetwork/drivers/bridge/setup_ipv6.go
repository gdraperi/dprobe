package bridge

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

var bridgeIPv6 *net.IPNet

const (
	bridgeIPv6Str          = "fe80::1/64"
	ipv6ForwardConfPerm    = 0644
	ipv6ForwardConfDefault = "/proc/sys/net/ipv6/conf/default/forwarding"
	ipv6ForwardConfAll     = "/proc/sys/net/ipv6/conf/all/forwarding"
)

func init() ***REMOVED***
	// We allow ourselves to panic in this special case because we indicate a
	// failure to parse a compile-time define constant.
	var err error
	if bridgeIPv6, err = types.ParseCIDR(bridgeIPv6Str); err != nil ***REMOVED***
		panic(fmt.Sprintf("Cannot parse default bridge IPv6 address %q: %v", bridgeIPv6Str, err))
	***REMOVED***
***REMOVED***

func setupBridgeIPv6(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
	procFile := "/proc/sys/net/ipv6/conf/" + config.BridgeName + "/disable_ipv6"
	ipv6BridgeData, err := ioutil.ReadFile(procFile)
	if err != nil ***REMOVED***
		return fmt.Errorf("Cannot read IPv6 setup for bridge %v: %v", config.BridgeName, err)
	***REMOVED***
	// Enable IPv6 on the bridge only if it isn't already enabled
	if ipv6BridgeData[0] != '0' ***REMOVED***
		if err := ioutil.WriteFile(procFile, []byte***REMOVED***'0', '\n'***REMOVED***, ipv6ForwardConfPerm); err != nil ***REMOVED***
			return fmt.Errorf("Unable to enable IPv6 addresses on bridge: %v", err)
		***REMOVED***
	***REMOVED***

	// Store bridge network and default gateway
	i.bridgeIPv6 = bridgeIPv6
	i.gatewayIPv6 = i.bridgeIPv6.IP

	if err := i.programIPv6Address(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if config.AddressIPv6 == nil ***REMOVED***
		return nil
	***REMOVED***

	// Store the user specified bridge network and network gateway and program it
	i.bridgeIPv6 = config.AddressIPv6
	i.gatewayIPv6 = config.AddressIPv6.IP

	if err := i.programIPv6Address(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Setting route to global IPv6 subnet
	logrus.Debugf("Adding route to IPv6 network %s via device %s", config.AddressIPv6.String(), config.BridgeName)
	err = i.nlh.RouteAdd(&netlink.Route***REMOVED***
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: i.Link.Attrs().Index,
		Dst:       config.AddressIPv6,
	***REMOVED***)
	if err != nil && !os.IsExist(err) ***REMOVED***
		logrus.Errorf("Could not add route to IPv6 network %s via device %s", config.AddressIPv6.String(), config.BridgeName)
	***REMOVED***

	return nil
***REMOVED***

func setupGatewayIPv6(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
	if config.AddressIPv6 == nil ***REMOVED***
		return &ErrInvalidContainerSubnet***REMOVED******REMOVED***
	***REMOVED***
	if !config.AddressIPv6.Contains(config.DefaultGatewayIPv6) ***REMOVED***
		return &ErrInvalidGateway***REMOVED******REMOVED***
	***REMOVED***

	// Store requested default gateway
	i.gatewayIPv6 = config.DefaultGatewayIPv6

	return nil
***REMOVED***

func setupIPv6Forwarding(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
	// Get current IPv6 default forwarding setup
	ipv6ForwardDataDefault, err := ioutil.ReadFile(ipv6ForwardConfDefault)
	if err != nil ***REMOVED***
		return fmt.Errorf("Cannot read IPv6 default forwarding setup: %v", err)
	***REMOVED***
	// Enable IPv6 default forwarding only if it is not already enabled
	if ipv6ForwardDataDefault[0] != '1' ***REMOVED***
		if err := ioutil.WriteFile(ipv6ForwardConfDefault, []byte***REMOVED***'1', '\n'***REMOVED***, ipv6ForwardConfPerm); err != nil ***REMOVED***
			logrus.Warnf("Unable to enable IPv6 default forwarding: %v", err)
		***REMOVED***
	***REMOVED***

	// Get current IPv6 all forwarding setup
	ipv6ForwardDataAll, err := ioutil.ReadFile(ipv6ForwardConfAll)
	if err != nil ***REMOVED***
		return fmt.Errorf("Cannot read IPv6 all forwarding setup: %v", err)
	***REMOVED***
	// Enable IPv6 all forwarding only if it is not already enabled
	if ipv6ForwardDataAll[0] != '1' ***REMOVED***
		if err := ioutil.WriteFile(ipv6ForwardConfAll, []byte***REMOVED***'1', '\n'***REMOVED***, ipv6ForwardConfPerm); err != nil ***REMOVED***
			logrus.Warnf("Unable to enable IPv6 all forwarding: %v", err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
