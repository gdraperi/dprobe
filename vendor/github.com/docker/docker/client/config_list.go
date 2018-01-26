package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// ConfigList returns the list of configs.
func (cli *Client) ConfigList(ctx context.Context, options types.ConfigListOptions) ([]swarm.Config, error) ***REMOVED***
	if err := cli.NewVersionError("1.30", "config list"); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	query := url.Values***REMOVED******REMOVED***

	if options.Filters.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToJSON(options.Filters)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		query.Set("filters", filterJSON)
	***REMOVED***

	resp, err := cli.get(ctx, "/configs", query, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var configs []swarm.Config
	err = json.NewDecoder(resp.body).Decode(&configs)
	ensureReaderClosed(resp)
	return configs, err
***REMOVED***
