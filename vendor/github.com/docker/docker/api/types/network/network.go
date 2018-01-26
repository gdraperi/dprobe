package network

// Address represents an IP address
type Address struct ***REMOVED***
	Addr      string
	PrefixLen int
***REMOVED***

// IPAM represents IP Address Management
type IPAM struct ***REMOVED***
	Driver  string
	Options map[string]string //Per network IPAM driver options
	Config  []IPAMConfig
***REMOVED***

// IPAMConfig represents IPAM configurations
type IPAMConfig struct ***REMOVED***
	Subnet     string            `json:",omitempty"`
	IPRange    string            `json:",omitempty"`
	Gateway    string            `json:",omitempty"`
	AuxAddress map[string]string `json:"AuxiliaryAddresses,omitempty"`
***REMOVED***

// EndpointIPAMConfig represents IPAM configurations for the endpoint
type EndpointIPAMConfig struct ***REMOVED***
	IPv4Address  string   `json:",omitempty"`
	IPv6Address  string   `json:",omitempty"`
	LinkLocalIPs []string `json:",omitempty"`
***REMOVED***

// Copy makes a copy of the endpoint ipam config
func (cfg *EndpointIPAMConfig) Copy() *EndpointIPAMConfig ***REMOVED***
	cfgCopy := *cfg
	cfgCopy.LinkLocalIPs = make([]string, 0, len(cfg.LinkLocalIPs))
	cfgCopy.LinkLocalIPs = append(cfgCopy.LinkLocalIPs, cfg.LinkLocalIPs...)
	return &cfgCopy
***REMOVED***

// PeerInfo represents one peer of an overlay network
type PeerInfo struct ***REMOVED***
	Name string
	IP   string
***REMOVED***

// EndpointSettings stores the network endpoint details
type EndpointSettings struct ***REMOVED***
	// Configurations
	IPAMConfig *EndpointIPAMConfig
	Links      []string
	Aliases    []string
	// Operational data
	NetworkID           string
	EndpointID          string
	Gateway             string
	IPAddress           string
	IPPrefixLen         int
	IPv6Gateway         string
	GlobalIPv6Address   string
	GlobalIPv6PrefixLen int
	MacAddress          string
	DriverOpts          map[string]string
***REMOVED***

// Task carries the information about one backend task
type Task struct ***REMOVED***
	Name       string
	EndpointID string
	EndpointIP string
	Info       map[string]string
***REMOVED***

// ServiceInfo represents service parameters with the list of service's tasks
type ServiceInfo struct ***REMOVED***
	VIP          string
	Ports        []string
	LocalLBIndex int
	Tasks        []Task
***REMOVED***

// Copy makes a deep copy of `EndpointSettings`
func (es *EndpointSettings) Copy() *EndpointSettings ***REMOVED***
	epCopy := *es
	if es.IPAMConfig != nil ***REMOVED***
		epCopy.IPAMConfig = es.IPAMConfig.Copy()
	***REMOVED***

	if es.Links != nil ***REMOVED***
		links := make([]string, 0, len(es.Links))
		epCopy.Links = append(links, es.Links...)
	***REMOVED***

	if es.Aliases != nil ***REMOVED***
		aliases := make([]string, 0, len(es.Aliases))
		epCopy.Aliases = append(aliases, es.Aliases...)
	***REMOVED***
	return &epCopy
***REMOVED***

// NetworkingConfig represents the container's networking configuration for each of its interfaces
// Carries the networking configs specified in the `docker run` and `docker network connect` commands
type NetworkingConfig struct ***REMOVED***
	EndpointsConfig map[string]*EndpointSettings // Endpoint configs for each connecting network
***REMOVED***

// ConfigReference specifies the source which provides a network's configuration
type ConfigReference struct ***REMOVED***
	Network string
***REMOVED***
