package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// NodeList returns the list of nodes.
func (cli *Client) NodeList(ctx context.Context, options types.NodeListOptions) ([]swarm.Node, error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***

	if options.Filters.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToJSON(options.Filters)

		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		query.Set("filters", filterJSON)
	***REMOVED***

	resp, err := cli.get(ctx, "/nodes", query, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var nodes []swarm.Node
	err = json.NewDecoder(resp.body).Decode(&nodes)
	ensureReaderClosed(resp)
	return nodes, err
***REMOVED***
