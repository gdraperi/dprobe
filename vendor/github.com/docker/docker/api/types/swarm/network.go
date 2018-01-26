package swarm

import (
	"github.com/docker/docker/api/types/network"
)

// Endpoint represents an endpoint.
type Endpoint struct ***REMOVED***
	Spec       EndpointSpec        `json:",omitempty"`
	Ports      []PortConfig        `json:",omitempty"`
	VirtualIPs []EndpointVirtualIP `json:",omitempty"`
***REMOVED***

// EndpointSpec represents the spec of an endpoint.
type EndpointSpec struct ***REMOVED***
	Mode  ResolutionMode `json:",omitempty"`
	Ports []PortConfig   `json:",omitempty"`
***REMOVED***

// ResolutionMode represents a resolution mode.
type ResolutionMode string

const (
	// ResolutionModeVIP VIP
	ResolutionModeVIP ResolutionMode = "vip"
	// ResolutionModeDNSRR DNSRR
	ResolutionModeDNSRR ResolutionMode = "dnsrr"
)

// PortConfig represents the config of a port.
type PortConfig struct ***REMOVED***
	Name     string             `json:",omitempty"`
	Protocol PortConfigProtocol `json:",omitempty"`
	// TargetPort is the port inside the container
	TargetPort uint32 `json:",omitempty"`
	// PublishedPort is the port on the swarm hosts
	PublishedPort uint32 `json:",omitempty"`
	// PublishMode is the mode in which port is published
	PublishMode PortConfigPublishMode `json:",omitempty"`
***REMOVED***

// PortConfigPublishMode represents the mode in which the port is to
// be published.
type PortConfigPublishMode string

const (
	// PortConfigPublishModeIngress is used for ports published
	// for ingress load balancing using routing mesh.
	PortConfigPublishModeIngress PortConfigPublishMode = "ingress"
	// PortConfigPublishModeHost is used for ports published
	// for direct host level access on the host where the task is running.
	PortConfigPublishModeHost PortConfigPublishMode = "host"
)

// PortConfigProtocol represents the protocol of a port.
type PortConfigProtocol string

const (
	// TODO(stevvooe): These should be used generally, not just for PortConfig.

	// PortConfigProtocolTCP TCP
	PortConfigProtocolTCP PortConfigProtocol = "tcp"
	// PortConfigProtocolUDP UDP
	PortConfigProtocolUDP PortConfigProtocol = "udp"
)

// EndpointVirtualIP represents the virtual ip of a port.
type EndpointVirtualIP struct ***REMOVED***
	NetworkID string `json:",omitempty"`
	Addr      string `json:",omitempty"`
***REMOVED***

// Network represents a network.
type Network struct ***REMOVED***
	ID string
	Meta
	Spec        NetworkSpec  `json:",omitempty"`
	DriverState Driver       `json:",omitempty"`
	IPAMOptions *IPAMOptions `json:",omitempty"`
***REMOVED***

// NetworkSpec represents the spec of a network.
type NetworkSpec struct ***REMOVED***
	Annotations
	DriverConfiguration *Driver                  `json:",omitempty"`
	IPv6Enabled         bool                     `json:",omitempty"`
	Internal            bool                     `json:",omitempty"`
	Attachable          bool                     `json:",omitempty"`
	Ingress             bool                     `json:",omitempty"`
	IPAMOptions         *IPAMOptions             `json:",omitempty"`
	ConfigFrom          *network.ConfigReference `json:",omitempty"`
	Scope               string                   `json:",omitempty"`
***REMOVED***

// NetworkAttachmentConfig represents the configuration of a network attachment.
type NetworkAttachmentConfig struct ***REMOVED***
	Target     string            `json:",omitempty"`
	Aliases    []string          `json:",omitempty"`
	DriverOpts map[string]string `json:",omitempty"`
***REMOVED***

// NetworkAttachment represents a network attachment.
type NetworkAttachment struct ***REMOVED***
	Network   Network  `json:",omitempty"`
	Addresses []string `json:",omitempty"`
***REMOVED***

// IPAMOptions represents ipam options.
type IPAMOptions struct ***REMOVED***
	Driver  Driver       `json:",omitempty"`
	Configs []IPAMConfig `json:",omitempty"`
***REMOVED***

// IPAMConfig represents ipam configuration.
type IPAMConfig struct ***REMOVED***
	Subnet  string `json:",omitempty"`
	Range   string `json:",omitempty"`
	Gateway string `json:",omitempty"`
***REMOVED***
