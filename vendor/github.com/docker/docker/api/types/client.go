package types

import (
	"bufio"
	"io"
	"net"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	units "github.com/docker/go-units"
)

// CheckpointCreateOptions holds parameters to create a checkpoint from a container
type CheckpointCreateOptions struct ***REMOVED***
	CheckpointID  string
	CheckpointDir string
	Exit          bool
***REMOVED***

// CheckpointListOptions holds parameters to list checkpoints for a container
type CheckpointListOptions struct ***REMOVED***
	CheckpointDir string
***REMOVED***

// CheckpointDeleteOptions holds parameters to delete a checkpoint from a container
type CheckpointDeleteOptions struct ***REMOVED***
	CheckpointID  string
	CheckpointDir string
***REMOVED***

// ContainerAttachOptions holds parameters to attach to a container.
type ContainerAttachOptions struct ***REMOVED***
	Stream     bool
	Stdin      bool
	Stdout     bool
	Stderr     bool
	DetachKeys string
	Logs       bool
***REMOVED***

// ContainerCommitOptions holds parameters to commit changes into a container.
type ContainerCommitOptions struct ***REMOVED***
	Reference string
	Comment   string
	Author    string
	Changes   []string
	Pause     bool
	Config    *container.Config
***REMOVED***

// ContainerExecInspect holds information returned by exec inspect.
type ContainerExecInspect struct ***REMOVED***
	ExecID      string
	ContainerID string
	Running     bool
	ExitCode    int
	Pid         int
***REMOVED***

// ContainerListOptions holds parameters to list containers with.
type ContainerListOptions struct ***REMOVED***
	Quiet   bool
	Size    bool
	All     bool
	Latest  bool
	Since   string
	Before  string
	Limit   int
	Filters filters.Args
***REMOVED***

// ContainerLogsOptions holds parameters to filter logs with.
type ContainerLogsOptions struct ***REMOVED***
	ShowStdout bool
	ShowStderr bool
	Since      string
	Until      string
	Timestamps bool
	Follow     bool
	Tail       string
	Details    bool
***REMOVED***

// ContainerRemoveOptions holds parameters to remove containers.
type ContainerRemoveOptions struct ***REMOVED***
	RemoveVolumes bool
	RemoveLinks   bool
	Force         bool
***REMOVED***

// ContainerStartOptions holds parameters to start containers.
type ContainerStartOptions struct ***REMOVED***
	CheckpointID  string
	CheckpointDir string
***REMOVED***

// CopyToContainerOptions holds information
// about files to copy into a container
type CopyToContainerOptions struct ***REMOVED***
	AllowOverwriteDirWithFile bool
	CopyUIDGID                bool
***REMOVED***

// EventsOptions holds parameters to filter events with.
type EventsOptions struct ***REMOVED***
	Since   string
	Until   string
	Filters filters.Args
***REMOVED***

// NetworkListOptions holds parameters to filter the list of networks with.
type NetworkListOptions struct ***REMOVED***
	Filters filters.Args
***REMOVED***

// HijackedResponse holds connection information for a hijacked request.
type HijackedResponse struct ***REMOVED***
	Conn   net.Conn
	Reader *bufio.Reader
***REMOVED***

// Close closes the hijacked connection and reader.
func (h *HijackedResponse) Close() ***REMOVED***
	h.Conn.Close()
***REMOVED***

// CloseWriter is an interface that implements structs
// that close input streams to prevent from writing.
type CloseWriter interface ***REMOVED***
	CloseWrite() error
***REMOVED***

// CloseWrite closes a readWriter for writing.
func (h *HijackedResponse) CloseWrite() error ***REMOVED***
	if conn, ok := h.Conn.(CloseWriter); ok ***REMOVED***
		return conn.CloseWrite()
	***REMOVED***
	return nil
***REMOVED***

// ImageBuildOptions holds the information
// necessary to build images.
type ImageBuildOptions struct ***REMOVED***
	Tags           []string
	SuppressOutput bool
	RemoteContext  string
	NoCache        bool
	Remove         bool
	ForceRemove    bool
	PullParent     bool
	Isolation      container.Isolation
	CPUSetCPUs     string
	CPUSetMems     string
	CPUShares      int64
	CPUQuota       int64
	CPUPeriod      int64
	Memory         int64
	MemorySwap     int64
	CgroupParent   string
	NetworkMode    string
	ShmSize        int64
	Dockerfile     string
	Ulimits        []*units.Ulimit
	// BuildArgs needs to be a *string instead of just a string so that
	// we can tell the difference between "" (empty string) and no value
	// at all (nil). See the parsing of buildArgs in
	// api/server/router/build/build_routes.go for even more info.
	BuildArgs   map[string]*string
	AuthConfigs map[string]AuthConfig
	Context     io.Reader
	Labels      map[string]string
	// squash the resulting image's layers to the parent
	// preserves the original image and creates a new one from the parent with all
	// the changes applied to a single layer
	Squash bool
	// CacheFrom specifies images that are used for matching cache. Images
	// specified here do not need to have a valid parent chain to match cache.
	CacheFrom   []string
	SecurityOpt []string
	ExtraHosts  []string // List of extra hosts
	Target      string
	SessionID   string
	Platform    string
***REMOVED***

// ImageBuildResponse holds information
// returned by a server after building
// an image.
type ImageBuildResponse struct ***REMOVED***
	Body   io.ReadCloser
	OSType string
***REMOVED***

// ImageCreateOptions holds information to create images.
type ImageCreateOptions struct ***REMOVED***
	RegistryAuth string // RegistryAuth is the base64 encoded credentials for the registry.
	Platform     string // Platform is the target platform of the image if it needs to be pulled from the registry.
***REMOVED***

// ImageImportSource holds source information for ImageImport
type ImageImportSource struct ***REMOVED***
	Source     io.Reader // Source is the data to send to the server to create this image from. You must set SourceName to "-" to leverage this.
	SourceName string    // SourceName is the name of the image to pull. Set to "-" to leverage the Source attribute.
***REMOVED***

// ImageImportOptions holds information to import images from the client host.
type ImageImportOptions struct ***REMOVED***
	Tag      string   // Tag is the name to tag this image with. This attribute is deprecated.
	Message  string   // Message is the message to tag the image with
	Changes  []string // Changes are the raw changes to apply to this image
	Platform string   // Platform is the target platform of the image
***REMOVED***

// ImageListOptions holds parameters to filter the list of images with.
type ImageListOptions struct ***REMOVED***
	All     bool
	Filters filters.Args
***REMOVED***

// ImageLoadResponse returns information to the client about a load process.
type ImageLoadResponse struct ***REMOVED***
	// Body must be closed to avoid a resource leak
	Body io.ReadCloser
	JSON bool
***REMOVED***

// ImagePullOptions holds information to pull images.
type ImagePullOptions struct ***REMOVED***
	All           bool
	RegistryAuth  string // RegistryAuth is the base64 encoded credentials for the registry
	PrivilegeFunc RequestPrivilegeFunc
	Platform      string
***REMOVED***

// RequestPrivilegeFunc is a function interface that
// clients can supply to retry operations after
// getting an authorization error.
// This function returns the registry authentication
// header value in base 64 format, or an error
// if the privilege request fails.
type RequestPrivilegeFunc func() (string, error)

//ImagePushOptions holds information to push images.
type ImagePushOptions ImagePullOptions

// ImageRemoveOptions holds parameters to remove images.
type ImageRemoveOptions struct ***REMOVED***
	Force         bool
	PruneChildren bool
***REMOVED***

// ImageSearchOptions holds parameters to search images with.
type ImageSearchOptions struct ***REMOVED***
	RegistryAuth  string
	PrivilegeFunc RequestPrivilegeFunc
	Filters       filters.Args
	Limit         int
***REMOVED***

// ResizeOptions holds parameters to resize a tty.
// It can be used to resize container ttys and
// exec process ttys too.
type ResizeOptions struct ***REMOVED***
	Height uint
	Width  uint
***REMOVED***

// NodeListOptions holds parameters to list nodes with.
type NodeListOptions struct ***REMOVED***
	Filters filters.Args
***REMOVED***

// NodeRemoveOptions holds parameters to remove nodes with.
type NodeRemoveOptions struct ***REMOVED***
	Force bool
***REMOVED***

// ServiceCreateOptions contains the options to use when creating a service.
type ServiceCreateOptions struct ***REMOVED***
	// EncodedRegistryAuth is the encoded registry authorization credentials to
	// use when updating the service.
	//
	// This field follows the format of the X-Registry-Auth header.
	EncodedRegistryAuth string

	// QueryRegistry indicates whether the service update requires
	// contacting a registry. A registry may be contacted to retrieve
	// the image digest and manifest, which in turn can be used to update
	// platform or other information about the service.
	QueryRegistry bool
***REMOVED***

// ServiceCreateResponse contains the information returned to a client
// on the creation of a new service.
type ServiceCreateResponse struct ***REMOVED***
	// ID is the ID of the created service.
	ID string
	// Warnings is a set of non-fatal warning messages to pass on to the user.
	Warnings []string `json:",omitempty"`
***REMOVED***

// Values for RegistryAuthFrom in ServiceUpdateOptions
const (
	RegistryAuthFromSpec         = "spec"
	RegistryAuthFromPreviousSpec = "previous-spec"
)

// ServiceUpdateOptions contains the options to be used for updating services.
type ServiceUpdateOptions struct ***REMOVED***
	// EncodedRegistryAuth is the encoded registry authorization credentials to
	// use when updating the service.
	//
	// This field follows the format of the X-Registry-Auth header.
	EncodedRegistryAuth string

	// TODO(stevvooe): Consider moving the version parameter of ServiceUpdate
	// into this field. While it does open API users up to racy writes, most
	// users may not need that level of consistency in practice.

	// RegistryAuthFrom specifies where to find the registry authorization
	// credentials if they are not given in EncodedRegistryAuth. Valid
	// values are "spec" and "previous-spec".
	RegistryAuthFrom string

	// Rollback indicates whether a server-side rollback should be
	// performed. When this is set, the provided spec will be ignored.
	// The valid values are "previous" and "none". An empty value is the
	// same as "none".
	Rollback string

	// QueryRegistry indicates whether the service update requires
	// contacting a registry. A registry may be contacted to retrieve
	// the image digest and manifest, which in turn can be used to update
	// platform or other information about the service.
	QueryRegistry bool
***REMOVED***

// ServiceListOptions holds parameters to list services with.
type ServiceListOptions struct ***REMOVED***
	Filters filters.Args
***REMOVED***

// ServiceInspectOptions holds parameters related to the "service inspect"
// operation.
type ServiceInspectOptions struct ***REMOVED***
	InsertDefaults bool
***REMOVED***

// TaskListOptions holds parameters to list tasks with.
type TaskListOptions struct ***REMOVED***
	Filters filters.Args
***REMOVED***

// PluginRemoveOptions holds parameters to remove plugins.
type PluginRemoveOptions struct ***REMOVED***
	Force bool
***REMOVED***

// PluginEnableOptions holds parameters to enable plugins.
type PluginEnableOptions struct ***REMOVED***
	Timeout int
***REMOVED***

// PluginDisableOptions holds parameters to disable plugins.
type PluginDisableOptions struct ***REMOVED***
	Force bool
***REMOVED***

// PluginInstallOptions holds parameters to install a plugin.
type PluginInstallOptions struct ***REMOVED***
	Disabled              bool
	AcceptAllPermissions  bool
	RegistryAuth          string // RegistryAuth is the base64 encoded credentials for the registry
	RemoteRef             string // RemoteRef is the plugin name on the registry
	PrivilegeFunc         RequestPrivilegeFunc
	AcceptPermissionsFunc func(PluginPrivileges) (bool, error)
	Args                  []string
***REMOVED***

// SwarmUnlockKeyResponse contains the response for Engine API:
// GET /swarm/unlockkey
type SwarmUnlockKeyResponse struct ***REMOVED***
	// UnlockKey is the unlock key in ASCII-armored format.
	UnlockKey string
***REMOVED***

// PluginCreateOptions hold all options to plugin create.
type PluginCreateOptions struct ***REMOVED***
	RepoName string
***REMOVED***
