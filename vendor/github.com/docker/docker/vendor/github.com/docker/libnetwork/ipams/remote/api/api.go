// Package api defines the data structure to be used in the request/response
// messages between libnetwork and the remote ipam plugin
package api

import "github.com/docker/libnetwork/ipamapi"

// Response is the basic response structure used in all responses
type Response struct ***REMOVED***
	Error string
***REMOVED***

// IsSuccess returns whether the plugin response is successful
func (r *Response) IsSuccess() bool ***REMOVED***
	return r.Error == ""
***REMOVED***

// GetError returns the error from the response, if any.
func (r *Response) GetError() string ***REMOVED***
	return r.Error
***REMOVED***

// GetCapabilityResponse is the response of GetCapability request
type GetCapabilityResponse struct ***REMOVED***
	Response
	RequiresMACAddress    bool
	RequiresRequestReplay bool
***REMOVED***

// ToCapability converts the capability response into the internal ipam driver capability structure
func (capRes GetCapabilityResponse) ToCapability() *ipamapi.Capability ***REMOVED***
	return &ipamapi.Capability***REMOVED***
		RequiresMACAddress:    capRes.RequiresMACAddress,
		RequiresRequestReplay: capRes.RequiresRequestReplay,
	***REMOVED***
***REMOVED***

// GetAddressSpacesResponse is the response to the ``get default address spaces`` request message
type GetAddressSpacesResponse struct ***REMOVED***
	Response
	LocalDefaultAddressSpace  string
	GlobalDefaultAddressSpace string
***REMOVED***

// RequestPoolRequest represents the expected data in a ``request address pool`` request message
type RequestPoolRequest struct ***REMOVED***
	AddressSpace string
	Pool         string
	SubPool      string
	Options      map[string]string
	V6           bool
***REMOVED***

// RequestPoolResponse represents the response message to a ``request address pool`` request
type RequestPoolResponse struct ***REMOVED***
	Response
	PoolID string
	Pool   string // CIDR format
	Data   map[string]string
***REMOVED***

// ReleasePoolRequest represents the expected data in a ``release address pool`` request message
type ReleasePoolRequest struct ***REMOVED***
	PoolID string
***REMOVED***

// ReleasePoolResponse represents the response message to a ``release address pool`` request
type ReleasePoolResponse struct ***REMOVED***
	Response
***REMOVED***

// RequestAddressRequest represents the expected data in a ``request address`` request message
type RequestAddressRequest struct ***REMOVED***
	PoolID  string
	Address string
	Options map[string]string
***REMOVED***

// RequestAddressResponse represents the expected data in the response message to a ``request address`` request
type RequestAddressResponse struct ***REMOVED***
	Response
	Address string // in CIDR format
	Data    map[string]string
***REMOVED***

// ReleaseAddressRequest represents the expected data in a ``release address`` request message
type ReleaseAddressRequest struct ***REMOVED***
	PoolID  string
	Address string
***REMOVED***

// ReleaseAddressResponse represents the response message to a ``release address`` request
type ReleaseAddressResponse struct ***REMOVED***
	Response
***REMOVED***
