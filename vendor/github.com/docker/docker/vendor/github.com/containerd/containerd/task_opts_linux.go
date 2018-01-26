package containerd

import (
	"context"

	"github.com/opencontainers/runtime-spec/specs-go"
)

// WithResources sets the provided resources for task updates
func WithResources(resources *specs.LinuxResources) UpdateTaskOpts ***REMOVED***
	return func(ctx context.Context, client *Client, r *UpdateTaskInfo) error ***REMOVED***
		r.Resources = resources
		return nil
	***REMOVED***
***REMOVED***
