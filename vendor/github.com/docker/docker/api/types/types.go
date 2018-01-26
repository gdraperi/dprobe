package types

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/go-connections/nat"
)

// RootFS returns Image's RootFS description including the layer IDs.
type RootFS struct ***REMOVED***
	Type      string
	Layers    []string `json:",omitempty"`
	BaseLayer string   `json:",omitempty"`
***REMOVED***

// ImageInspect contains response of Engine API:
// GET "/images/***REMOVED***name:.****REMOVED***/json"
type ImageInspect struct ***REMOVED***
	ID              string `json:"Id"`
	RepoTags        []string
	RepoDigests     []string
	Parent          string
	Comment         string
	Created         string
	Container       string
	ContainerConfig *container.Config
	DockerVersion   string
	Author          string
	Config          *container.Config
	Architecture    string
	Os              string
	OsVersion       string `json:",omitempty"`
	Size            int64
	VirtualSize     int64
	GraphDriver     GraphDriverData
	RootFS          RootFS
	Metadata        ImageMetadata
***REMOVED***

// ImageMetadata contains engine-local data about the image
type ImageMetadata struct ***REMOVED***
	LastTagTime time.Time `json:",omitempty"`
***REMOVED***

// Container contains response of Engine API:
// GET "/containers/json"
type Container struct ***REMOVED***
	ID         string `json:"Id"`
	Names      []string
	Image      string
	ImageID    string
	Command    string
	Created    int64
	Ports      []Port
	SizeRw     int64 `json:",omitempty"`
	SizeRootFs int64 `json:",omitempty"`
	Labels     map[string]string
	State      string
	Status     string
	HostConfig struct ***REMOVED***
		NetworkMode string `json:",omitempty"`
	***REMOVED***
	NetworkSettings *SummaryNetworkSettings
	Mounts          []MountPoint
***REMOVED***

// CopyConfig contains request body of Engine API:
// POST "/containers/"+containerID+"/copy"
type CopyConfig struct ***REMOVED***
	Resource string
***REMOVED***

// ContainerPathStat is used to encode the header from
// GET "/containers/***REMOVED***name:.****REMOVED***/archive"
// "Name" is the file or directory name.
type ContainerPathStat struct ***REMOVED***
	Name       string      `json:"name"`
	Size       int64       `json:"size"`
	Mode       os.FileMode `json:"mode"`
	Mtime      time.Time   `json:"mtime"`
	LinkTarget string      `json:"linkTarget"`
***REMOVED***

// ContainerStats contains response of Engine API:
// GET "/stats"
type ContainerStats struct ***REMOVED***
	Body   io.ReadCloser `json:"body"`
	OSType string        `json:"ostype"`
***REMOVED***

// Ping contains response of Engine API:
// GET "/_ping"
type Ping struct ***REMOVED***
	APIVersion   string
	OSType       string
	Experimental bool
***REMOVED***

// ComponentVersion describes the version information for a specific component.
type ComponentVersion struct ***REMOVED***
	Name    string
	Version string
	Details map[string]string `json:",omitempty"`
***REMOVED***

// Version contains response of Engine API:
// GET "/version"
type Version struct ***REMOVED***
	Platform   struct***REMOVED*** Name string ***REMOVED*** `json:",omitempty"`
	Components []ComponentVersion    `json:",omitempty"`

	// The following fields are deprecated, they relate to the Engine component and are kept for backwards compatibility

	Version       string
	APIVersion    string `json:"ApiVersion"`
	MinAPIVersion string `json:"MinAPIVersion,omitempty"`
	GitCommit     string
	GoVersion     string
	Os            string
	Arch          string
	KernelVersion string `json:",omitempty"`
	Experimental  bool   `json:",omitempty"`
	BuildTime     string `json:",omitempty"`
***REMOVED***

// Commit holds the Git-commit (SHA1) that a binary was built from, as reported
// in the version-string of external tools, such as containerd, or runC.
type Commit struct ***REMOVED***
	ID       string // ID is the actual commit ID of external tool.
	Expected string // Expected is the commit ID of external tool expected by dockerd as set at build time.
***REMOVED***

// Info contains response of Engine API:
// GET "/info"
type Info struct ***REMOVED***
	ID                 string
	Containers         int
	ContainersRunning  int
	ContainersPaused   int
	ContainersStopped  int
	Images             int
	Driver             string
	DriverStatus       [][2]string
	SystemStatus       [][2]string
	Plugins            PluginsInfo
	MemoryLimit        bool
	SwapLimit          bool
	KernelMemory       bool
	CPUCfsPeriod       bool `json:"CpuCfsPeriod"`
	CPUCfsQuota        bool `json:"CpuCfsQuota"`
	CPUShares          bool
	CPUSet             bool
	IPv4Forwarding     bool
	BridgeNfIptables   bool
	BridgeNfIP6tables  bool `json:"BridgeNfIp6tables"`
	Debug              bool
	NFd                int
	OomKillDisable     bool
	NGoroutines        int
	SystemTime         string
	LoggingDriver      string
	CgroupDriver       string
	NEventsListener    int
	KernelVersion      string
	OperatingSystem    string
	OSType             string
	Architecture       string
	IndexServerAddress string
	RegistryConfig     *registry.ServiceConfig
	NCPU               int
	MemTotal           int64
	GenericResources   []swarm.GenericResource
	DockerRootDir      string
	HTTPProxy          string `json:"HttpProxy"`
	HTTPSProxy         string `json:"HttpsProxy"`
	NoProxy            string
	Name               string
	Labels             []string
	ExperimentalBuild  bool
	ServerVersion      string
	ClusterStore       string
	ClusterAdvertise   string
	Runtimes           map[string]Runtime
	DefaultRuntime     string
	Swarm              swarm.Info
	// LiveRestoreEnabled determines whether containers should be kept
	// running when the daemon is shutdown or upon daemon start if
	// running containers are detected
	LiveRestoreEnabled bool
	Isolation          container.Isolation
	InitBinary         string
	ContainerdCommit   Commit
	RuncCommit         Commit
	InitCommit         Commit
	SecurityOptions    []string
***REMOVED***

// KeyValue holds a key/value pair
type KeyValue struct ***REMOVED***
	Key, Value string
***REMOVED***

// SecurityOpt contains the name and options of a security option
type SecurityOpt struct ***REMOVED***
	Name    string
	Options []KeyValue
***REMOVED***

// DecodeSecurityOptions decodes a security options string slice to a type safe
// SecurityOpt
func DecodeSecurityOptions(opts []string) ([]SecurityOpt, error) ***REMOVED***
	so := []SecurityOpt***REMOVED******REMOVED***
	for _, opt := range opts ***REMOVED***
		// support output from a < 1.13 docker daemon
		if !strings.Contains(opt, "=") ***REMOVED***
			so = append(so, SecurityOpt***REMOVED***Name: opt***REMOVED***)
			continue
		***REMOVED***
		secopt := SecurityOpt***REMOVED******REMOVED***
		split := strings.Split(opt, ",")
		for _, s := range split ***REMOVED***
			kv := strings.SplitN(s, "=", 2)
			if len(kv) != 2 ***REMOVED***
				return nil, fmt.Errorf("invalid security option %q", s)
			***REMOVED***
			if kv[0] == "" || kv[1] == "" ***REMOVED***
				return nil, errors.New("invalid empty security option")
			***REMOVED***
			if kv[0] == "name" ***REMOVED***
				secopt.Name = kv[1]
				continue
			***REMOVED***
			secopt.Options = append(secopt.Options, KeyValue***REMOVED***Key: kv[0], Value: kv[1]***REMOVED***)
		***REMOVED***
		so = append(so, secopt)
	***REMOVED***
	return so, nil
***REMOVED***

// PluginsInfo is a temp struct holding Plugins name
// registered with docker daemon. It is used by Info struct
type PluginsInfo struct ***REMOVED***
	// List of Volume plugins registered
	Volume []string
	// List of Network plugins registered
	Network []string
	// List of Authorization plugins registered
	Authorization []string
	// List of Log plugins registered
	Log []string
***REMOVED***

// ExecStartCheck is a temp struct used by execStart
// Config fields is part of ExecConfig in runconfig package
type ExecStartCheck struct ***REMOVED***
	// ExecStart will first check if it's detached
	Detach bool
	// Check if there's a tty
	Tty bool
***REMOVED***

// HealthcheckResult stores information about a single run of a healthcheck probe
type HealthcheckResult struct ***REMOVED***
	Start    time.Time // Start is the time this check started
	End      time.Time // End is the time this check ended
	ExitCode int       // ExitCode meanings: 0=healthy, 1=unhealthy, 2=reserved (considered unhealthy), else=error running probe
	Output   string    // Output from last check
***REMOVED***

// Health states
const (
	NoHealthcheck = "none"      // Indicates there is no healthcheck
	Starting      = "starting"  // Starting indicates that the container is not yet ready
	Healthy       = "healthy"   // Healthy indicates that the container is running correctly
	Unhealthy     = "unhealthy" // Unhealthy indicates that the container has a problem
)

// Health stores information about the container's healthcheck results
type Health struct ***REMOVED***
	Status        string               // Status is one of Starting, Healthy or Unhealthy
	FailingStreak int                  // FailingStreak is the number of consecutive failures
	Log           []*HealthcheckResult // Log contains the last few results (oldest first)
***REMOVED***

// ContainerState stores container's running state
// it's part of ContainerJSONBase and will return by "inspect" command
type ContainerState struct ***REMOVED***
	Status     string // String representation of the container state. Can be one of "created", "running", "paused", "restarting", "removing", "exited", or "dead"
	Running    bool
	Paused     bool
	Restarting bool
	OOMKilled  bool
	Dead       bool
	Pid        int
	ExitCode   int
	Error      string
	StartedAt  string
	FinishedAt string
	Health     *Health `json:",omitempty"`
***REMOVED***

// ContainerNode stores information about the node that a container
// is running on.  It's only available in Docker Swarm
type ContainerNode struct ***REMOVED***
	ID        string
	IPAddress string `json:"IP"`
	Addr      string
	Name      string
	Cpus      int
	Memory    int64
	Labels    map[string]string
***REMOVED***

// ContainerJSONBase contains response of Engine API:
// GET "/containers/***REMOVED***name:.****REMOVED***/json"
type ContainerJSONBase struct ***REMOVED***
	ID              string `json:"Id"`
	Created         string
	Path            string
	Args            []string
	State           *ContainerState
	Image           string
	ResolvConfPath  string
	HostnamePath    string
	HostsPath       string
	LogPath         string
	Node            *ContainerNode `json:",omitempty"`
	Name            string
	RestartCount    int
	Driver          string
	Platform        string
	MountLabel      string
	ProcessLabel    string
	AppArmorProfile string
	ExecIDs         []string
	HostConfig      *container.HostConfig
	GraphDriver     GraphDriverData
	SizeRw          *int64 `json:",omitempty"`
	SizeRootFs      *int64 `json:",omitempty"`
***REMOVED***

// ContainerJSON is newly used struct along with MountPoint
type ContainerJSON struct ***REMOVED***
	*ContainerJSONBase
	Mounts          []MountPoint
	Config          *container.Config
	NetworkSettings *NetworkSettings
***REMOVED***

// NetworkSettings exposes the network settings in the api
type NetworkSettings struct ***REMOVED***
	NetworkSettingsBase
	DefaultNetworkSettings
	Networks map[string]*network.EndpointSettings
***REMOVED***

// SummaryNetworkSettings provides a summary of container's networks
// in /containers/json
type SummaryNetworkSettings struct ***REMOVED***
	Networks map[string]*network.EndpointSettings
***REMOVED***

// NetworkSettingsBase holds basic information about networks
type NetworkSettingsBase struct ***REMOVED***
	Bridge                 string      // Bridge is the Bridge name the network uses(e.g. `docker0`)
	SandboxID              string      // SandboxID uniquely represents a container's network stack
	HairpinMode            bool        // HairpinMode specifies if hairpin NAT should be enabled on the virtual interface
	LinkLocalIPv6Address   string      // LinkLocalIPv6Address is an IPv6 unicast address using the link-local prefix
	LinkLocalIPv6PrefixLen int         // LinkLocalIPv6PrefixLen is the prefix length of an IPv6 unicast address
	Ports                  nat.PortMap // Ports is a collection of PortBinding indexed by Port
	SandboxKey             string      // SandboxKey identifies the sandbox
	SecondaryIPAddresses   []network.Address
	SecondaryIPv6Addresses []network.Address
***REMOVED***

// DefaultNetworkSettings holds network information
// during the 2 release deprecation period.
// It will be removed in Docker 1.11.
type DefaultNetworkSettings struct ***REMOVED***
	EndpointID          string // EndpointID uniquely represents a service endpoint in a Sandbox
	Gateway             string // Gateway holds the gateway address for the network
	GlobalIPv6Address   string // GlobalIPv6Address holds network's global IPv6 address
	GlobalIPv6PrefixLen int    // GlobalIPv6PrefixLen represents mask length of network's global IPv6 address
	IPAddress           string // IPAddress holds the IPv4 address for the network
	IPPrefixLen         int    // IPPrefixLen represents mask length of network's IPv4 address
	IPv6Gateway         string // IPv6Gateway holds gateway address specific for IPv6
	MacAddress          string // MacAddress holds the MAC address for the network
***REMOVED***

// MountPoint represents a mount point configuration inside the container.
// This is used for reporting the mountpoints in use by a container.
type MountPoint struct ***REMOVED***
	Type        mount.Type `json:",omitempty"`
	Name        string     `json:",omitempty"`
	Source      string
	Destination string
	Driver      string `json:",omitempty"`
	Mode        string
	RW          bool
	Propagation mount.Propagation
***REMOVED***

// NetworkResource is the body of the "get network" http response message
type NetworkResource struct ***REMOVED***
	Name       string                         // Name is the requested name of the network
	ID         string                         `json:"Id"` // ID uniquely identifies a network on a single machine
	Created    time.Time                      // Created is the time the network created
	Scope      string                         // Scope describes the level at which the network exists (e.g. `swarm` for cluster-wide or `local` for machine level)
	Driver     string                         // Driver is the Driver name used to create the network (e.g. `bridge`, `overlay`)
	EnableIPv6 bool                           // EnableIPv6 represents whether to enable IPv6
	IPAM       network.IPAM                   // IPAM is the network's IP Address Management
	Internal   bool                           // Internal represents if the network is used internal only
	Attachable bool                           // Attachable represents if the global scope is manually attachable by regular containers from workers in swarm mode.
	Ingress    bool                           // Ingress indicates the network is providing the routing-mesh for the swarm cluster.
	ConfigFrom network.ConfigReference        // ConfigFrom specifies the source which will provide the configuration for this network.
	ConfigOnly bool                           // ConfigOnly networks are place-holder networks for network configurations to be used by other networks. ConfigOnly networks cannot be used directly to run containers or services.
	Containers map[string]EndpointResource    // Containers contains endpoints belonging to the network
	Options    map[string]string              // Options holds the network specific options to use for when creating the network
	Labels     map[string]string              // Labels holds metadata specific to the network being created
	Peers      []network.PeerInfo             `json:",omitempty"` // List of peer nodes for an overlay network
	Services   map[string]network.ServiceInfo `json:",omitempty"`
***REMOVED***

// EndpointResource contains network resources allocated and used for a container in a network
type EndpointResource struct ***REMOVED***
	Name        string
	EndpointID  string
	MacAddress  string
	IPv4Address string
	IPv6Address string
***REMOVED***

// NetworkCreate is the expected body of the "create network" http request message
type NetworkCreate struct ***REMOVED***
	// Check for networks with duplicate names.
	// Network is primarily keyed based on a random ID and not on the name.
	// Network name is strictly a user-friendly alias to the network
	// which is uniquely identified using ID.
	// And there is no guaranteed way to check for duplicates.
	// Option CheckDuplicate is there to provide a best effort checking of any networks
	// which has the same name but it is not guaranteed to catch all name collisions.
	CheckDuplicate bool
	Driver         string
	Scope          string
	EnableIPv6     bool
	IPAM           *network.IPAM
	Internal       bool
	Attachable     bool
	Ingress        bool
	ConfigOnly     bool
	ConfigFrom     *network.ConfigReference
	Options        map[string]string
	Labels         map[string]string
***REMOVED***

// NetworkCreateRequest is the request message sent to the server for network create call.
type NetworkCreateRequest struct ***REMOVED***
	NetworkCreate
	Name string
***REMOVED***

// NetworkCreateResponse is the response message sent by the server for network create call
type NetworkCreateResponse struct ***REMOVED***
	ID      string `json:"Id"`
	Warning string
***REMOVED***

// NetworkConnect represents the data to be used to connect a container to the network
type NetworkConnect struct ***REMOVED***
	Container      string
	EndpointConfig *network.EndpointSettings `json:",omitempty"`
***REMOVED***

// NetworkDisconnect represents the data to be used to disconnect a container from the network
type NetworkDisconnect struct ***REMOVED***
	Container string
	Force     bool
***REMOVED***

// NetworkInspectOptions holds parameters to inspect network
type NetworkInspectOptions struct ***REMOVED***
	Scope   string
	Verbose bool
***REMOVED***

// Checkpoint represents the details of a checkpoint
type Checkpoint struct ***REMOVED***
	Name string // Name is the name of the checkpoint
***REMOVED***

// Runtime describes an OCI runtime
type Runtime struct ***REMOVED***
	Path string   `json:"path"`
	Args []string `json:"runtimeArgs,omitempty"`
***REMOVED***

// DiskUsage contains response of Engine API:
// GET "/system/df"
type DiskUsage struct ***REMOVED***
	LayersSize  int64
	Images      []*ImageSummary
	Containers  []*Container
	Volumes     []*Volume
	BuilderSize int64
***REMOVED***

// ContainersPruneReport contains the response for Engine API:
// POST "/containers/prune"
type ContainersPruneReport struct ***REMOVED***
	ContainersDeleted []string
	SpaceReclaimed    uint64
***REMOVED***

// VolumesPruneReport contains the response for Engine API:
// POST "/volumes/prune"
type VolumesPruneReport struct ***REMOVED***
	VolumesDeleted []string
	SpaceReclaimed uint64
***REMOVED***

// ImagesPruneReport contains the response for Engine API:
// POST "/images/prune"
type ImagesPruneReport struct ***REMOVED***
	ImagesDeleted  []ImageDeleteResponseItem
	SpaceReclaimed uint64
***REMOVED***

// BuildCachePruneReport contains the response for Engine API:
// POST "/build/prune"
type BuildCachePruneReport struct ***REMOVED***
	SpaceReclaimed uint64
***REMOVED***

// NetworksPruneReport contains the response for Engine API:
// POST "/networks/prune"
type NetworksPruneReport struct ***REMOVED***
	NetworksDeleted []string
***REMOVED***

// SecretCreateResponse contains the information returned to a client
// on the creation of a new secret.
type SecretCreateResponse struct ***REMOVED***
	// ID is the id of the created secret.
	ID string
***REMOVED***

// SecretListOptions holds parameters to list secrets
type SecretListOptions struct ***REMOVED***
	Filters filters.Args
***REMOVED***

// ConfigCreateResponse contains the information returned to a client
// on the creation of a new config.
type ConfigCreateResponse struct ***REMOVED***
	// ID is the id of the created config.
	ID string
***REMOVED***

// ConfigListOptions holds parameters to list configs
type ConfigListOptions struct ***REMOVED***
	Filters filters.Args
***REMOVED***

// PushResult contains the tag, manifest digest, and manifest size from the
// push. It's used to signal this information to the trust code in the client
// so it can sign the manifest if necessary.
type PushResult struct ***REMOVED***
	Tag    string
	Digest string
	Size   int
***REMOVED***

// BuildResult contains the image id of a successful build
type BuildResult struct ***REMOVED***
	ID string
***REMOVED***
