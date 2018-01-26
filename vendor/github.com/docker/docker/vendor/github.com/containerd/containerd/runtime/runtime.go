package runtime

import (
	"context"
	"time"

	"github.com/containerd/containerd/mount"
	"github.com/gogo/protobuf/types"
)

// IO holds process IO information
type IO struct ***REMOVED***
	Stdin    string
	Stdout   string
	Stderr   string
	Terminal bool
***REMOVED***

// CreateOpts contains task creation data
type CreateOpts struct ***REMOVED***
	// Spec is the OCI runtime spec
	Spec *types.Any
	// Rootfs mounts to perform to gain access to the container's filesystem
	Rootfs []mount.Mount
	// IO for the container's main process
	IO IO
	// Checkpoint digest to restore container state
	Checkpoint string
	// Options for the runtime and container
	Options *types.Any
***REMOVED***

// Exit information for a process
type Exit struct ***REMOVED***
	Pid       uint32
	Status    uint32
	Timestamp time.Time
***REMOVED***

// Runtime is responsible for the creation of containers for a certain platform,
// arch, or custom usage.
type Runtime interface ***REMOVED***
	// ID of the runtime
	ID() string
	// Create creates a task with the provided id and options.
	Create(ctx context.Context, id string, opts CreateOpts) (Task, error)
	// Get returns a task.
	Get(context.Context, string) (Task, error)
	// Tasks returns all the current tasks for the runtime.
	// Any container runs at most one task at a time.
	Tasks(context.Context) ([]Task, error)
	// Delete removes the task in the runtime.
	Delete(context.Context, Task) (*Exit, error)
***REMOVED***
