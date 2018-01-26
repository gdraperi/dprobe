package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// CheckpointList returns the checkpoints of the given container in the docker host
func (cli *Client) CheckpointList(ctx context.Context, container string, options types.CheckpointListOptions) ([]types.Checkpoint, error) ***REMOVED***
	var checkpoints []types.Checkpoint

	query := url.Values***REMOVED******REMOVED***
	if options.CheckpointDir != "" ***REMOVED***
		query.Set("dir", options.CheckpointDir)
	***REMOVED***

	resp, err := cli.get(ctx, "/containers/"+container+"/checkpoints", query, nil)
	if err != nil ***REMOVED***
		return checkpoints, wrapResponseError(err, resp, "container", container)
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&checkpoints)
	ensureReaderClosed(resp)
	return checkpoints, err
***REMOVED***
