package bridge

import (
	"fmt"
	"net"
)

// ErrConfigExists error is returned when driver already has a config applied.
type ErrConfigExists struct***REMOVED******REMOVED***

func (ece *ErrConfigExists) Error() string ***REMOVED***
	return "configuration already exists, bridge configuration can be applied only once"
***REMOVED***

// Forbidden denotes the type of this error
func (ece *ErrConfigExists) Forbidden() ***REMOVED******REMOVED***

// ErrInvalidDriverConfig error is returned when Bridge Driver is passed an invalid config
type ErrInvalidDriverConfig struct***REMOVED******REMOVED***

func (eidc *ErrInvalidDriverConfig) Error() string ***REMOVED***
	return "Invalid configuration passed to Bridge Driver"
***REMOVED***

// BadRequest denotes the type of this error
func (eidc *ErrInvalidDriverConfig) BadRequest() ***REMOVED******REMOVED***

// ErrInvalidNetworkConfig error is returned when a network is created on a driver without valid config.
type ErrInvalidNetworkConfig struct***REMOVED******REMOVED***

func (einc *ErrInvalidNetworkConfig) Error() string ***REMOVED***
	return "trying to create a network on a driver without valid config"
***REMOVED***

// Forbidden denotes the type of this error
func (einc *ErrInvalidNetworkConfig) Forbidden() ***REMOVED******REMOVED***

// ErrInvalidContainerConfig error is returned when an endpoint create is attempted with an invalid configuration.
type ErrInvalidContainerConfig struct***REMOVED******REMOVED***

func (eicc *ErrInvalidContainerConfig) Error() string ***REMOVED***
	return "Error in joining a container due to invalid configuration"
***REMOVED***

// BadRequest denotes the type of this error
func (eicc *ErrInvalidContainerConfig) BadRequest() ***REMOVED******REMOVED***

// ErrInvalidEndpointConfig error is returned when an endpoint create is attempted with an invalid endpoint configuration.
type ErrInvalidEndpointConfig struct***REMOVED******REMOVED***

func (eiec *ErrInvalidEndpointConfig) Error() string ***REMOVED***
	return "trying to create an endpoint with an invalid endpoint configuration"
***REMOVED***

// BadRequest denotes the type of this error
func (eiec *ErrInvalidEndpointConfig) BadRequest() ***REMOVED******REMOVED***

// ErrNetworkExists error is returned when a network already exists and another network is created.
type ErrNetworkExists struct***REMOVED******REMOVED***

func (ene *ErrNetworkExists) Error() string ***REMOVED***
	return "network already exists, bridge can only have one network"
***REMOVED***

// Forbidden denotes the type of this error
func (ene *ErrNetworkExists) Forbidden() ***REMOVED******REMOVED***

// ErrIfaceName error is returned when a new name could not be generated.
type ErrIfaceName struct***REMOVED******REMOVED***

func (ein *ErrIfaceName) Error() string ***REMOVED***
	return "failed to find name for new interface"
***REMOVED***

// InternalError denotes the type of this error
func (ein *ErrIfaceName) InternalError() ***REMOVED******REMOVED***

// ErrNoIPAddr error is returned when bridge has no IPv4 address configured.
type ErrNoIPAddr struct***REMOVED******REMOVED***

func (enip *ErrNoIPAddr) Error() string ***REMOVED***
	return "bridge has no IPv4 address configured"
***REMOVED***

// InternalError denotes the type of this error
func (enip *ErrNoIPAddr) InternalError() ***REMOVED******REMOVED***

// ErrInvalidGateway is returned when the user provided default gateway (v4/v6) is not not valid.
type ErrInvalidGateway struct***REMOVED******REMOVED***

func (eig *ErrInvalidGateway) Error() string ***REMOVED***
	return "default gateway ip must be part of the network"
***REMOVED***

// BadRequest denotes the type of this error
func (eig *ErrInvalidGateway) BadRequest() ***REMOVED******REMOVED***

// ErrInvalidContainerSubnet is returned when the container subnet (FixedCIDR) is not valid.
type ErrInvalidContainerSubnet struct***REMOVED******REMOVED***

func (eis *ErrInvalidContainerSubnet) Error() string ***REMOVED***
	return "container subnet must be a subset of bridge network"
***REMOVED***

// BadRequest denotes the type of this error
func (eis *ErrInvalidContainerSubnet) BadRequest() ***REMOVED******REMOVED***

// ErrInvalidMtu is returned when the user provided MTU is not valid.
type ErrInvalidMtu int

func (eim ErrInvalidMtu) Error() string ***REMOVED***
	return fmt.Sprintf("invalid MTU number: %d", int(eim))
***REMOVED***

// BadRequest denotes the type of this error
func (eim ErrInvalidMtu) BadRequest() ***REMOVED******REMOVED***

// ErrInvalidPort is returned when the container or host port specified in the port binding is not valid.
type ErrInvalidPort string

func (ip ErrInvalidPort) Error() string ***REMOVED***
	return fmt.Sprintf("invalid transport port: %s", string(ip))
***REMOVED***

// BadRequest denotes the type of this error
func (ip ErrInvalidPort) BadRequest() ***REMOVED******REMOVED***

// ErrUnsupportedAddressType is returned when the specified address type is not supported.
type ErrUnsupportedAddressType string

func (uat ErrUnsupportedAddressType) Error() string ***REMOVED***
	return fmt.Sprintf("unsupported address type: %s", string(uat))
***REMOVED***

// BadRequest denotes the type of this error
func (uat ErrUnsupportedAddressType) BadRequest() ***REMOVED******REMOVED***

// ErrInvalidAddressBinding is returned when the host address specified in the port binding is not valid.
type ErrInvalidAddressBinding string

func (iab ErrInvalidAddressBinding) Error() string ***REMOVED***
	return fmt.Sprintf("invalid host address in port binding: %s", string(iab))
***REMOVED***

// BadRequest denotes the type of this error
func (iab ErrInvalidAddressBinding) BadRequest() ***REMOVED******REMOVED***

// ActiveEndpointsError is returned when there are
// still active endpoints in the network being deleted.
type ActiveEndpointsError string

func (aee ActiveEndpointsError) Error() string ***REMOVED***
	return fmt.Sprintf("network %s has active endpoint", string(aee))
***REMOVED***

// Forbidden denotes the type of this error
func (aee ActiveEndpointsError) Forbidden() ***REMOVED******REMOVED***

// InvalidNetworkIDError is returned when the passed
// network id for an existing network is not a known id.
type InvalidNetworkIDError string

func (inie InvalidNetworkIDError) Error() string ***REMOVED***
	return fmt.Sprintf("invalid network id %s", string(inie))
***REMOVED***

// NotFound denotes the type of this error
func (inie InvalidNetworkIDError) NotFound() ***REMOVED******REMOVED***

// InvalidEndpointIDError is returned when the passed
// endpoint id is not valid.
type InvalidEndpointIDError string

func (ieie InvalidEndpointIDError) Error() string ***REMOVED***
	return fmt.Sprintf("invalid endpoint id: %s", string(ieie))
***REMOVED***

// BadRequest denotes the type of this error
func (ieie InvalidEndpointIDError) BadRequest() ***REMOVED******REMOVED***

// InvalidSandboxIDError is returned when the passed
// sandbox id is not valid.
type InvalidSandboxIDError string

func (isie InvalidSandboxIDError) Error() string ***REMOVED***
	return fmt.Sprintf("invalid sandbox id: %s", string(isie))
***REMOVED***

// BadRequest denotes the type of this error
func (isie InvalidSandboxIDError) BadRequest() ***REMOVED******REMOVED***

// EndpointNotFoundError is returned when the no endpoint
// with the passed endpoint id is found.
type EndpointNotFoundError string

func (enfe EndpointNotFoundError) Error() string ***REMOVED***
	return fmt.Sprintf("endpoint not found: %s", string(enfe))
***REMOVED***

// NotFound denotes the type of this error
func (enfe EndpointNotFoundError) NotFound() ***REMOVED******REMOVED***

// NonDefaultBridgeExistError is returned when a non-default
// bridge config is passed but it does not already exist.
type NonDefaultBridgeExistError string

func (ndbee NonDefaultBridgeExistError) Error() string ***REMOVED***
	return fmt.Sprintf("bridge device with non default name %s must be created manually", string(ndbee))
***REMOVED***

// Forbidden denotes the type of this error
func (ndbee NonDefaultBridgeExistError) Forbidden() ***REMOVED******REMOVED***

// NonDefaultBridgeNeedsIPError is returned when a non-default
// bridge config is passed but it has no ip configured
type NonDefaultBridgeNeedsIPError string

func (ndbee NonDefaultBridgeNeedsIPError) Error() string ***REMOVED***
	return fmt.Sprintf("bridge device with non default name %s must have a valid IP address", string(ndbee))
***REMOVED***

// Forbidden denotes the type of this error
func (ndbee NonDefaultBridgeNeedsIPError) Forbidden() ***REMOVED******REMOVED***

// FixedCIDRv4Error is returned when fixed-cidrv4 configuration
// failed.
type FixedCIDRv4Error struct ***REMOVED***
	Net    *net.IPNet
	Subnet *net.IPNet
	Err    error
***REMOVED***

func (fcv4 *FixedCIDRv4Error) Error() string ***REMOVED***
	return fmt.Sprintf("setup FixedCIDRv4 failed for subnet %s in %s: %v", fcv4.Subnet, fcv4.Net, fcv4.Err)
***REMOVED***

// InternalError denotes the type of this error
func (fcv4 *FixedCIDRv4Error) InternalError() ***REMOVED******REMOVED***

// FixedCIDRv6Error is returned when fixed-cidrv6 configuration
// failed.
type FixedCIDRv6Error struct ***REMOVED***
	Net *net.IPNet
	Err error
***REMOVED***

func (fcv6 *FixedCIDRv6Error) Error() string ***REMOVED***
	return fmt.Sprintf("setup FixedCIDRv6 failed for subnet %s in %s: %v", fcv6.Net, fcv6.Net, fcv6.Err)
***REMOVED***

// InternalError denotes the type of this error
func (fcv6 *FixedCIDRv6Error) InternalError() ***REMOVED******REMOVED***

// IPTableCfgError is returned when an unexpected ip tables configuration is entered
type IPTableCfgError string

func (name IPTableCfgError) Error() string ***REMOVED***
	return fmt.Sprintf("unexpected request to set IP tables for interface: %s", string(name))
***REMOVED***

// BadRequest denotes the type of this error
func (name IPTableCfgError) BadRequest() ***REMOVED******REMOVED***

// InvalidIPTablesCfgError is returned when an invalid ip tables configuration is entered
type InvalidIPTablesCfgError string

func (action InvalidIPTablesCfgError) Error() string ***REMOVED***
	return fmt.Sprintf("Invalid IPTables action '%s'", string(action))
***REMOVED***

// BadRequest denotes the type of this error
func (action InvalidIPTablesCfgError) BadRequest() ***REMOVED******REMOVED***

// IPv4AddrRangeError is returned when a valid IP address range couldn't be found.
type IPv4AddrRangeError string

func (name IPv4AddrRangeError) Error() string ***REMOVED***
	return fmt.Sprintf("can't find an address range for interface %q", string(name))
***REMOVED***

// BadRequest denotes the type of this error
func (name IPv4AddrRangeError) BadRequest() ***REMOVED******REMOVED***

// IPv4AddrAddError is returned when IPv4 address could not be added to the bridge.
type IPv4AddrAddError struct ***REMOVED***
	IP  *net.IPNet
	Err error
***REMOVED***

func (ipv4 *IPv4AddrAddError) Error() string ***REMOVED***
	return fmt.Sprintf("failed to add IPv4 address %s to bridge: %v", ipv4.IP, ipv4.Err)
***REMOVED***

// InternalError denotes the type of this error
func (ipv4 *IPv4AddrAddError) InternalError() ***REMOVED******REMOVED***

// IPv6AddrAddError is returned when IPv6 address could not be added to the bridge.
type IPv6AddrAddError struct ***REMOVED***
	IP  *net.IPNet
	Err error
***REMOVED***

func (ipv6 *IPv6AddrAddError) Error() string ***REMOVED***
	return fmt.Sprintf("failed to add IPv6 address %s to bridge: %v", ipv6.IP, ipv6.Err)
***REMOVED***

// InternalError denotes the type of this error
func (ipv6 *IPv6AddrAddError) InternalError() ***REMOVED******REMOVED***

// IPv4AddrNoMatchError is returned when the bridge's IPv4 address does not match configured.
type IPv4AddrNoMatchError struct ***REMOVED***
	IP    net.IP
	CfgIP net.IP
***REMOVED***

func (ipv4 *IPv4AddrNoMatchError) Error() string ***REMOVED***
	return fmt.Sprintf("bridge IPv4 (%s) does not match requested configuration %s", ipv4.IP, ipv4.CfgIP)
***REMOVED***

// BadRequest denotes the type of this error
func (ipv4 *IPv4AddrNoMatchError) BadRequest() ***REMOVED******REMOVED***

// IPv6AddrNoMatchError is returned when the bridge's IPv6 address does not match configured.
type IPv6AddrNoMatchError net.IPNet

func (ipv6 *IPv6AddrNoMatchError) Error() string ***REMOVED***
	return fmt.Sprintf("bridge IPv6 addresses do not match the expected bridge configuration %s", (*net.IPNet)(ipv6).String())
***REMOVED***

// BadRequest denotes the type of this error
func (ipv6 *IPv6AddrNoMatchError) BadRequest() ***REMOVED******REMOVED***

// InvalidLinkIPAddrError is returned when a link is configured to a container with an invalid ip address
type InvalidLinkIPAddrError string

func (address InvalidLinkIPAddrError) Error() string ***REMOVED***
	return fmt.Sprintf("Cannot link to a container with Invalid IP Address '%s'", string(address))
***REMOVED***

// BadRequest denotes the type of this error
func (address InvalidLinkIPAddrError) BadRequest() ***REMOVED******REMOVED***
