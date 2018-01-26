package swarm

import (
	"time"

	"github.com/docker/docker/api/types/swarm/runtime"
)

// TaskState represents the state of a task.
type TaskState string

const (
	// TaskStateNew NEW
	TaskStateNew TaskState = "new"
	// TaskStateAllocated ALLOCATED
	TaskStateAllocated TaskState = "allocated"
	// TaskStatePending PENDING
	TaskStatePending TaskState = "pending"
	// TaskStateAssigned ASSIGNED
	TaskStateAssigned TaskState = "assigned"
	// TaskStateAccepted ACCEPTED
	TaskStateAccepted TaskState = "accepted"
	// TaskStatePreparing PREPARING
	TaskStatePreparing TaskState = "preparing"
	// TaskStateReady READY
	TaskStateReady TaskState = "ready"
	// TaskStateStarting STARTING
	TaskStateStarting TaskState = "starting"
	// TaskStateRunning RUNNING
	TaskStateRunning TaskState = "running"
	// TaskStateComplete COMPLETE
	TaskStateComplete TaskState = "complete"
	// TaskStateShutdown SHUTDOWN
	TaskStateShutdown TaskState = "shutdown"
	// TaskStateFailed FAILED
	TaskStateFailed TaskState = "failed"
	// TaskStateRejected REJECTED
	TaskStateRejected TaskState = "rejected"
)

// Task represents a task.
type Task struct ***REMOVED***
	ID string
	Meta
	Annotations

	Spec                TaskSpec            `json:",omitempty"`
	ServiceID           string              `json:",omitempty"`
	Slot                int                 `json:",omitempty"`
	NodeID              string              `json:",omitempty"`
	Status              TaskStatus          `json:",omitempty"`
	DesiredState        TaskState           `json:",omitempty"`
	NetworksAttachments []NetworkAttachment `json:",omitempty"`
	GenericResources    []GenericResource   `json:",omitempty"`
***REMOVED***

// TaskSpec represents the spec of a task.
type TaskSpec struct ***REMOVED***
	// ContainerSpec and PluginSpec are mutually exclusive.
	// PluginSpec will only be used when the `Runtime` field is set to `plugin`
	ContainerSpec *ContainerSpec      `json:",omitempty"`
	PluginSpec    *runtime.PluginSpec `json:",omitempty"`

	Resources     *ResourceRequirements     `json:",omitempty"`
	RestartPolicy *RestartPolicy            `json:",omitempty"`
	Placement     *Placement                `json:",omitempty"`
	Networks      []NetworkAttachmentConfig `json:",omitempty"`

	// LogDriver specifies the LogDriver to use for tasks created from this
	// spec. If not present, the one on cluster default on swarm.Spec will be
	// used, finally falling back to the engine default if not specified.
	LogDriver *Driver `json:",omitempty"`

	// ForceUpdate is a counter that triggers an update even if no relevant
	// parameters have been changed.
	ForceUpdate uint64

	Runtime RuntimeType `json:",omitempty"`
***REMOVED***

// Resources represents resources (CPU/Memory).
type Resources struct ***REMOVED***
	NanoCPUs         int64             `json:",omitempty"`
	MemoryBytes      int64             `json:",omitempty"`
	GenericResources []GenericResource `json:",omitempty"`
***REMOVED***

// GenericResource represents a "user defined" resource which can
// be either an integer (e.g: SSD=3) or a string (e.g: SSD=sda1)
type GenericResource struct ***REMOVED***
	NamedResourceSpec    *NamedGenericResource    `json:",omitempty"`
	DiscreteResourceSpec *DiscreteGenericResource `json:",omitempty"`
***REMOVED***

// NamedGenericResource represents a "user defined" resource which is defined
// as a string.
// "Kind" is used to describe the Kind of a resource (e.g: "GPU", "FPGA", "SSD", ...)
// Value is used to identify the resource (GPU="UUID-1", FPGA="/dev/sdb5", ...)
type NamedGenericResource struct ***REMOVED***
	Kind  string `json:",omitempty"`
	Value string `json:",omitempty"`
***REMOVED***

// DiscreteGenericResource represents a "user defined" resource which is defined
// as an integer
// "Kind" is used to describe the Kind of a resource (e.g: "GPU", "FPGA", "SSD", ...)
// Value is used to count the resource (SSD=5, HDD=3, ...)
type DiscreteGenericResource struct ***REMOVED***
	Kind  string `json:",omitempty"`
	Value int64  `json:",omitempty"`
***REMOVED***

// ResourceRequirements represents resources requirements.
type ResourceRequirements struct ***REMOVED***
	Limits       *Resources `json:",omitempty"`
	Reservations *Resources `json:",omitempty"`
***REMOVED***

// Placement represents orchestration parameters.
type Placement struct ***REMOVED***
	Constraints []string              `json:",omitempty"`
	Preferences []PlacementPreference `json:",omitempty"`

	// Platforms stores all the platforms that the image can run on.
	// This field is used in the platform filter for scheduling. If empty,
	// then the platform filter is off, meaning there are no scheduling restrictions.
	Platforms []Platform `json:",omitempty"`
***REMOVED***

// PlacementPreference provides a way to make the scheduler aware of factors
// such as topology.
type PlacementPreference struct ***REMOVED***
	Spread *SpreadOver
***REMOVED***

// SpreadOver is a scheduling preference that instructs the scheduler to spread
// tasks evenly over groups of nodes identified by labels.
type SpreadOver struct ***REMOVED***
	// label descriptor, such as engine.labels.az
	SpreadDescriptor string
***REMOVED***

// RestartPolicy represents the restart policy.
type RestartPolicy struct ***REMOVED***
	Condition   RestartPolicyCondition `json:",omitempty"`
	Delay       *time.Duration         `json:",omitempty"`
	MaxAttempts *uint64                `json:",omitempty"`
	Window      *time.Duration         `json:",omitempty"`
***REMOVED***

// RestartPolicyCondition represents when to restart.
type RestartPolicyCondition string

const (
	// RestartPolicyConditionNone NONE
	RestartPolicyConditionNone RestartPolicyCondition = "none"
	// RestartPolicyConditionOnFailure ON_FAILURE
	RestartPolicyConditionOnFailure RestartPolicyCondition = "on-failure"
	// RestartPolicyConditionAny ANY
	RestartPolicyConditionAny RestartPolicyCondition = "any"
)

// TaskStatus represents the status of a task.
type TaskStatus struct ***REMOVED***
	Timestamp       time.Time       `json:",omitempty"`
	State           TaskState       `json:",omitempty"`
	Message         string          `json:",omitempty"`
	Err             string          `json:",omitempty"`
	ContainerStatus ContainerStatus `json:",omitempty"`
	PortStatus      PortStatus      `json:",omitempty"`
***REMOVED***

// ContainerStatus represents the status of a container.
type ContainerStatus struct ***REMOVED***
	ContainerID string `json:",omitempty"`
	PID         int    `json:",omitempty"`
	ExitCode    int    `json:",omitempty"`
***REMOVED***

// PortStatus represents the port status of a task's host ports whose
// service has published host ports
type PortStatus struct ***REMOVED***
	Ports []PortConfig `json:",omitempty"`
***REMOVED***
