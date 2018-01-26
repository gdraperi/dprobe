package bridge

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Enumeration type saying which versions of IP protocol to process.
type ipVersion int

const (
	ipvnone ipVersion = iota
	ipv4
	ipv6
	ipvboth
)

//Gets the IP version in use ( [ipv4], [ipv6] or [ipv4 and ipv6] )
func getIPVersion(config *networkConfiguration) ipVersion ***REMOVED***
	ipVersion := ipv4
	if config.AddressIPv6 != nil || config.EnableIPv6 ***REMOVED***
		ipVersion |= ipv6
	***REMOVED***
	return ipVersion
***REMOVED***

func setupBridgeNetFiltering(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
	err := checkBridgeNetFiltering(config, i)
	if err != nil ***REMOVED***
		if ptherr, ok := err.(*os.PathError); ok ***REMOVED***
			if errno, ok := ptherr.Err.(syscall.Errno); ok && errno == syscall.ENOENT ***REMOVED***
				if isRunningInContainer() ***REMOVED***
					logrus.Warnf("running inside docker container, ignoring missing kernel params: %v", err)
					err = nil
				***REMOVED*** else ***REMOVED***
					err = errors.New("please ensure that br_netfilter kernel module is loaded")
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			return fmt.Errorf("cannot restrict inter-container communication: %v", err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

//Enable bridge net filtering if ip forwarding is enabled. See github issue #11404
func checkBridgeNetFiltering(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
	ipVer := getIPVersion(config)
	iface := config.BridgeName
	doEnable := func(ipVer ipVersion) error ***REMOVED***
		var ipVerName string
		if ipVer == ipv4 ***REMOVED***
			ipVerName = "IPv4"
		***REMOVED*** else ***REMOVED***
			ipVerName = "IPv6"
		***REMOVED***
		enabled, err := isPacketForwardingEnabled(ipVer, iface)
		if err != nil ***REMOVED***
			logrus.Warnf("failed to check %s forwarding: %v", ipVerName, err)
		***REMOVED*** else if enabled ***REMOVED***
			enabled, err := getKernelBoolParam(getBridgeNFKernelParam(ipVer))
			if err != nil || enabled ***REMOVED***
				return err
			***REMOVED***
			return setKernelBoolParam(getBridgeNFKernelParam(ipVer), true)
		***REMOVED***
		return nil
	***REMOVED***

	switch ipVer ***REMOVED***
	case ipv4, ipv6:
		return doEnable(ipVer)
	case ipvboth:
		v4err := doEnable(ipv4)
		v6err := doEnable(ipv6)
		if v4err == nil ***REMOVED***
			return v6err
		***REMOVED***
		return v4err
	default:
		return nil
	***REMOVED***
***REMOVED***

// Get kernel param path saying whether IPv$***REMOVED***ipVer***REMOVED*** traffic is being forwarded
// on particular interface. Interface may be specified for IPv6 only. If
// `iface` is empty, `default` will be assumed, which represents default value
// for new interfaces.
func getForwardingKernelParam(ipVer ipVersion, iface string) string ***REMOVED***
	switch ipVer ***REMOVED***
	case ipv4:
		return "/proc/sys/net/ipv4/ip_forward"
	case ipv6:
		if iface == "" ***REMOVED***
			iface = "default"
		***REMOVED***
		return fmt.Sprintf("/proc/sys/net/ipv6/conf/%s/forwarding", iface)
	default:
		return ""
	***REMOVED***
***REMOVED***

// Get kernel param path saying whether bridged IPv$***REMOVED***ipVer***REMOVED*** traffic shall be
// passed to ip$***REMOVED***ipVer***REMOVED***tables' chains.
func getBridgeNFKernelParam(ipVer ipVersion) string ***REMOVED***
	switch ipVer ***REMOVED***
	case ipv4:
		return "/proc/sys/net/bridge/bridge-nf-call-iptables"
	case ipv6:
		return "/proc/sys/net/bridge/bridge-nf-call-ip6tables"
	default:
		return ""
	***REMOVED***
***REMOVED***

//Gets the value of the kernel parameters located at the given path
func getKernelBoolParam(path string) (bool, error) ***REMOVED***
	enabled := false
	line, err := ioutil.ReadFile(path)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	if len(line) > 0 ***REMOVED***
		enabled = line[0] == '1'
	***REMOVED***
	return enabled, err
***REMOVED***

//Sets the value of the kernel parameter located at the given path
func setKernelBoolParam(path string, on bool) error ***REMOVED***
	value := byte('0')
	if on ***REMOVED***
		value = byte('1')
	***REMOVED***
	return ioutil.WriteFile(path, []byte***REMOVED***value, '\n'***REMOVED***, 0644)
***REMOVED***

//Checks to see if packet forwarding is enabled
func isPacketForwardingEnabled(ipVer ipVersion, iface string) (bool, error) ***REMOVED***
	switch ipVer ***REMOVED***
	case ipv4, ipv6:
		return getKernelBoolParam(getForwardingKernelParam(ipVer, iface))
	case ipvboth:
		enabled, err := getKernelBoolParam(getForwardingKernelParam(ipv4, ""))
		if err != nil || !enabled ***REMOVED***
			return enabled, err
		***REMOVED***
		return getKernelBoolParam(getForwardingKernelParam(ipv6, iface))
	default:
		return true, nil
	***REMOVED***
***REMOVED***

func isRunningInContainer() bool ***REMOVED***
	_, err := os.Stat("/.dockerenv")
	return !os.IsNotExist(err)
***REMOVED***
