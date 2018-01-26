package containerd

import (
	"context"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// WithResources sets the provided resources on the spec for task updates
func WithResources(resources *specs.WindowsResources) UpdateTaskOpts ***REMOVED***
	return func(ctx context.Context, client *Client, r *UpdateTaskInfo) error ***REMOVED***
		r.Resources = resources
		return nil
	***REMOVED***
***REMOVED***
