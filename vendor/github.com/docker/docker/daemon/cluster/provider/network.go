package provider

import "github.com/docker/docker/api/types"

// NetworkCreateRequest is a request when creating a network.
type NetworkCreateRequest struct ***REMOVED***
	ID string
	types.NetworkCreateRequest
***REMOVED***

// NetworkCreateResponse is a response when creating a network.
type NetworkCreateResponse struct ***REMOVED***
	ID string `json:"Id"`
***REMOVED***

// VirtualAddress represents a virtual address.
type VirtualAddress struct ***REMOVED***
	IPv4 string
	IPv6 string
***REMOVED***

// PortConfig represents a port configuration.
type PortConfig struct ***REMOVED***
	Name          string
	Protocol      int32
	TargetPort    uint32
	PublishedPort uint32
***REMOVED***

// ServiceConfig represents a service configuration.
type ServiceConfig struct ***REMOVED***
	ID               string
	Name             string
	Aliases          map[string][]string
	VirtualAddresses map[string]*VirtualAddress
	ExposedPorts     []*PortConfig
***REMOVED***
