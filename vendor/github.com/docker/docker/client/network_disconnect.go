package client

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// NetworkDisconnect disconnects a container from an existent network in the docker host.
func (cli *Client) NetworkDisconnect(ctx context.Context, networkID, containerID string, force bool) error ***REMOVED***
	nd := types.NetworkDisconnect***REMOVED***Container: containerID, Force: force***REMOVED***
	resp, err := cli.post(ctx, "/networks/"+networkID+"/disconnect", nil, nd, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***
