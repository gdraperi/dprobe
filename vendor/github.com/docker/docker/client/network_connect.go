package client

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"golang.org/x/net/context"
)

// NetworkConnect connects a container to an existent network in the docker host.
func (cli *Client) NetworkConnect(ctx context.Context, networkID, containerID string, config *network.EndpointSettings) error ***REMOVED***
	nc := types.NetworkConnect***REMOVED***
		Container:      containerID,
		EndpointConfig: config,
	***REMOVED***
	resp, err := cli.post(ctx, "/networks/"+networkID+"/connect", nil, nc, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***
