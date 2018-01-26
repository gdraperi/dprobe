package swarm

// Node represents a node.
type Node struct ***REMOVED***
	ID string
	Meta
	// Spec defines the desired state of the node as specified by the user.
	// The system will honor this and will *never* modify it.
	Spec NodeSpec `json:",omitempty"`
	// Description encapsulates the properties of the Node as reported by the
	// agent.
	Description NodeDescription `json:",omitempty"`
	// Status provides the current status of the node, as seen by the manager.
	Status NodeStatus `json:",omitempty"`
	// ManagerStatus provides the current status of the node's manager
	// component, if the node is a manager.
	ManagerStatus *ManagerStatus `json:",omitempty"`
***REMOVED***

// NodeSpec represents the spec of a node.
type NodeSpec struct ***REMOVED***
	Annotations
	Role         NodeRole         `json:",omitempty"`
	Availability NodeAvailability `json:",omitempty"`
***REMOVED***

// NodeRole represents the role of a node.
type NodeRole string

const (
	// NodeRoleWorker WORKER
	NodeRoleWorker NodeRole = "worker"
	// NodeRoleManager MANAGER
	NodeRoleManager NodeRole = "manager"
)

// NodeAvailability represents the availability of a node.
type NodeAvailability string

const (
	// NodeAvailabilityActive ACTIVE
	NodeAvailabilityActive NodeAvailability = "active"
	// NodeAvailabilityPause PAUSE
	NodeAvailabilityPause NodeAvailability = "pause"
	// NodeAvailabilityDrain DRAIN
	NodeAvailabilityDrain NodeAvailability = "drain"
)

// NodeDescription represents the description of a node.
type NodeDescription struct ***REMOVED***
	Hostname  string            `json:",omitempty"`
	Platform  Platform          `json:",omitempty"`
	Resources Resources         `json:",omitempty"`
	Engine    EngineDescription `json:",omitempty"`
	TLSInfo   TLSInfo           `json:",omitempty"`
***REMOVED***

// Platform represents the platform (Arch/OS).
type Platform struct ***REMOVED***
	Architecture string `json:",omitempty"`
	OS           string `json:",omitempty"`
***REMOVED***

// EngineDescription represents the description of an engine.
type EngineDescription struct ***REMOVED***
	EngineVersion string              `json:",omitempty"`
	Labels        map[string]string   `json:",omitempty"`
	Plugins       []PluginDescription `json:",omitempty"`
***REMOVED***

// PluginDescription represents the description of an engine plugin.
type PluginDescription struct ***REMOVED***
	Type string `json:",omitempty"`
	Name string `json:",omitempty"`
***REMOVED***

// NodeStatus represents the status of a node.
type NodeStatus struct ***REMOVED***
	State   NodeState `json:",omitempty"`
	Message string    `json:",omitempty"`
	Addr    string    `json:",omitempty"`
***REMOVED***

// Reachability represents the reachability of a node.
type Reachability string

const (
	// ReachabilityUnknown UNKNOWN
	ReachabilityUnknown Reachability = "unknown"
	// ReachabilityUnreachable UNREACHABLE
	ReachabilityUnreachable Reachability = "unreachable"
	// ReachabilityReachable REACHABLE
	ReachabilityReachable Reachability = "reachable"
)

// ManagerStatus represents the status of a manager.
type ManagerStatus struct ***REMOVED***
	Leader       bool         `json:",omitempty"`
	Reachability Reachability `json:",omitempty"`
	Addr         string       `json:",omitempty"`
***REMOVED***

// NodeState represents the state of a node.
type NodeState string

const (
	// NodeStateUnknown UNKNOWN
	NodeStateUnknown NodeState = "unknown"
	// NodeStateDown DOWN
	NodeStateDown NodeState = "down"
	// NodeStateReady READY
	NodeStateReady NodeState = "ready"
	// NodeStateDisconnected DISCONNECTED
	NodeStateDisconnected NodeState = "disconnected"
)
