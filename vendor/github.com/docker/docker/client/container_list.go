package client

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

// ContainerList returns the list of containers in the docker host.
func (cli *Client) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***

	if options.All ***REMOVED***
		query.Set("all", "1")
	***REMOVED***

	if options.Limit != -1 ***REMOVED***
		query.Set("limit", strconv.Itoa(options.Limit))
	***REMOVED***

	if options.Since != "" ***REMOVED***
		query.Set("since", options.Since)
	***REMOVED***

	if options.Before != "" ***REMOVED***
		query.Set("before", options.Before)
	***REMOVED***

	if options.Size ***REMOVED***
		query.Set("size", "1")
	***REMOVED***

	if options.Filters.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToParamWithVersion(cli.version, options.Filters)

		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		query.Set("filters", filterJSON)
	***REMOVED***

	resp, err := cli.get(ctx, "/containers/json", query, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var containers []types.Container
	err = json.NewDecoder(resp.body).Decode(&containers)
	ensureReaderClosed(resp)
	return containers, err
***REMOVED***
