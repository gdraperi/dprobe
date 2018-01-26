package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// ServiceList returns the list of services.
func (cli *Client) ServiceList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***

	if options.Filters.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToJSON(options.Filters)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		query.Set("filters", filterJSON)
	***REMOVED***

	resp, err := cli.get(ctx, "/services", query, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var services []swarm.Service
	err = json.NewDecoder(resp.body).Decode(&services)
	ensureReaderClosed(resp)
	return services, err
***REMOVED***
