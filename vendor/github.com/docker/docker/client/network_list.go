package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

// NetworkList returns the list of networks configured in the docker host.
func (cli *Client) NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	if options.Filters.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToParamWithVersion(cli.version, options.Filters)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		query.Set("filters", filterJSON)
	***REMOVED***
	var networkResources []types.NetworkResource
	resp, err := cli.get(ctx, "/networks", query, nil)
	if err != nil ***REMOVED***
		return networkResources, err
	***REMOVED***
	err = json.NewDecoder(resp.body).Decode(&networkResources)
	ensureReaderClosed(resp)
	return networkResources, err
***REMOVED***
