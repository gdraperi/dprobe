/*
Package api represents all requests and responses suitable for conversation
with a remote driver.
*/
package api

import (
	"net"

	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
)

// Response is the basic response structure used in all responses.
type Response struct ***REMOVED***
	Err string
***REMOVED***

// GetError returns the error from the response, if any.
func (r *Response) GetError() string ***REMOVED***
	return r.Err
***REMOVED***

// GetCapabilityResponse is the response of GetCapability request
type GetCapabilityResponse struct ***REMOVED***
	Response
	Scope             string
	ConnectivityScope string
***REMOVED***

// AllocateNetworkRequest requests allocation of new network by manager
type AllocateNetworkRequest struct ***REMOVED***
	// A network ID that remote plugins are expected to store for future
	// reference.
	NetworkID string

	// A free form map->object interface for communication of options.
	Options map[string]string

	// IPAMData contains the address pool information for this network
	IPv4Data, IPv6Data []driverapi.IPAMData
***REMOVED***

// AllocateNetworkResponse is the response to the AllocateNetworkRequest.
type AllocateNetworkResponse struct ***REMOVED***
	Response
	// A free form plugin specific string->string object to be sent in
	// CreateNetworkRequest call in the libnetwork agents
	Options map[string]string
***REMOVED***

// FreeNetworkRequest is the request to free allocated network in the manager
type FreeNetworkRequest struct ***REMOVED***
	// The ID of the network to be freed.
	NetworkID string
***REMOVED***

// FreeNetworkResponse is the response to a request for freeing a network.
type FreeNetworkResponse struct ***REMOVED***
	Response
***REMOVED***

// CreateNetworkRequest requests a new network.
type CreateNetworkRequest struct ***REMOVED***
	// A network ID that remote plugins are expected to store for future
	// reference.
	NetworkID string

	// A free form map->object interface for communication of options.
	Options map[string]interface***REMOVED******REMOVED***

	// IPAMData contains the address pool information for this network
	IPv4Data, IPv6Data []driverapi.IPAMData
***REMOVED***

// CreateNetworkResponse is the response to the CreateNetworkRequest.
type CreateNetworkResponse struct ***REMOVED***
	Response
***REMOVED***

// DeleteNetworkRequest is the request to delete an existing network.
type DeleteNetworkRequest struct ***REMOVED***
	// The ID of the network to delete.
	NetworkID string
***REMOVED***

// DeleteNetworkResponse is the response to a request for deleting a network.
type DeleteNetworkResponse struct ***REMOVED***
	Response
***REMOVED***

// CreateEndpointRequest is the request to create an endpoint within a network.
type CreateEndpointRequest struct ***REMOVED***
	// Provided at create time, this will be the network id referenced.
	NetworkID string
	// The ID of the endpoint for later reference.
	EndpointID string
	Interface  *EndpointInterface
	Options    map[string]interface***REMOVED******REMOVED***
***REMOVED***

// EndpointInterface represents an interface endpoint.
type EndpointInterface struct ***REMOVED***
	Address     string
	AddressIPv6 string
	MacAddress  string
***REMOVED***

// CreateEndpointResponse is the response to the CreateEndpoint action.
type CreateEndpointResponse struct ***REMOVED***
	Response
	Interface *EndpointInterface
***REMOVED***

// Interface is the representation of a linux interface.
type Interface struct ***REMOVED***
	Address     *net.IPNet
	AddressIPv6 *net.IPNet
	MacAddress  net.HardwareAddr
***REMOVED***

// DeleteEndpointRequest describes the API for deleting an endpoint.
type DeleteEndpointRequest struct ***REMOVED***
	NetworkID  string
	EndpointID string
***REMOVED***

// DeleteEndpointResponse is the response to the DeleteEndpoint action.
type DeleteEndpointResponse struct ***REMOVED***
	Response
***REMOVED***

// EndpointInfoRequest retrieves information about the endpoint from the network driver.
type EndpointInfoRequest struct ***REMOVED***
	NetworkID  string
	EndpointID string
***REMOVED***

// EndpointInfoResponse is the response to an EndpointInfoRequest.
type EndpointInfoResponse struct ***REMOVED***
	Response
	Value map[string]interface***REMOVED******REMOVED***
***REMOVED***

// JoinRequest describes the API for joining an endpoint to a sandbox.
type JoinRequest struct ***REMOVED***
	NetworkID  string
	EndpointID string
	SandboxKey string
	Options    map[string]interface***REMOVED******REMOVED***
***REMOVED***

// InterfaceName is the struct represetation of a pair of devices with source
// and destination, for the purposes of putting an endpoint into a container.
type InterfaceName struct ***REMOVED***
	SrcName   string
	DstName   string
	DstPrefix string
***REMOVED***

// StaticRoute is the plain JSON representation of a static route.
type StaticRoute struct ***REMOVED***
	Destination string
	RouteType   int
	NextHop     string
***REMOVED***

// JoinResponse is the response to a JoinRequest.
type JoinResponse struct ***REMOVED***
	Response
	InterfaceName         *InterfaceName
	Gateway               string
	GatewayIPv6           string
	StaticRoutes          []StaticRoute
	DisableGatewayService bool
***REMOVED***

// LeaveRequest describes the API for detaching an endpoint from a sandbox.
type LeaveRequest struct ***REMOVED***
	NetworkID  string
	EndpointID string
***REMOVED***

// LeaveResponse is the answer to LeaveRequest.
type LeaveResponse struct ***REMOVED***
	Response
***REMOVED***

// ProgramExternalConnectivityRequest describes the API for programming the external connectivity for the given endpoint.
type ProgramExternalConnectivityRequest struct ***REMOVED***
	NetworkID  string
	EndpointID string
	Options    map[string]interface***REMOVED******REMOVED***
***REMOVED***

// ProgramExternalConnectivityResponse is the answer to ProgramExternalConnectivityRequest.
type ProgramExternalConnectivityResponse struct ***REMOVED***
	Response
***REMOVED***

// RevokeExternalConnectivityRequest describes the API for revoking the external connectivity for the given endpoint.
type RevokeExternalConnectivityRequest struct ***REMOVED***
	NetworkID  string
	EndpointID string
***REMOVED***

// RevokeExternalConnectivityResponse is the answer to RevokeExternalConnectivityRequest.
type RevokeExternalConnectivityResponse struct ***REMOVED***
	Response
***REMOVED***

// DiscoveryNotification represents a discovery notification
type DiscoveryNotification struct ***REMOVED***
	DiscoveryType discoverapi.DiscoveryType
	DiscoveryData interface***REMOVED******REMOVED***
***REMOVED***

// DiscoveryResponse is used by libnetwork to log any plugin error processing the discovery notifications
type DiscoveryResponse struct ***REMOVED***
	Response
***REMOVED***
