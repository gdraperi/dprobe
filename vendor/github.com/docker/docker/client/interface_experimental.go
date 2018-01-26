package client

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

type apiClientExperimental interface ***REMOVED***
	CheckpointAPIClient
***REMOVED***

// CheckpointAPIClient defines API client methods for the checkpoints
type CheckpointAPIClient interface ***REMOVED***
	CheckpointCreate(ctx context.Context, container string, options types.CheckpointCreateOptions) error
	CheckpointDelete(ctx context.Context, container string, options types.CheckpointDeleteOptions) error
	CheckpointList(ctx context.Context, container string, options types.CheckpointListOptions) ([]types.Checkpoint, error)
***REMOVED***
