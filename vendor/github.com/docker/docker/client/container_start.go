package client

import (
	"net/url"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
)

// ContainerStart sends a request to the docker daemon to start a container.
func (cli *Client) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	if len(options.CheckpointID) != 0 ***REMOVED***
		query.Set("checkpoint", options.CheckpointID)
	***REMOVED***
	if len(options.CheckpointDir) != 0 ***REMOVED***
		query.Set("checkpoint-dir", options.CheckpointDir)
	***REMOVED***

	resp, err := cli.post(ctx, "/containers/"+containerID+"/start", query, nil, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***
