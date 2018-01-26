package client

import (
	"net/url"
	"strconv"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// NodeUpdate updates a Node.
func (cli *Client) NodeUpdate(ctx context.Context, nodeID string, version swarm.Version, node swarm.NodeSpec) error ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	query.Set("version", strconv.FormatUint(version.Index, 10))
	resp, err := cli.post(ctx, "/nodes/"+nodeID+"/update", query, node, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***
