package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

// PluginList returns the installed plugins
func (cli *Client) PluginList(ctx context.Context, filter filters.Args) (types.PluginsListResponse, error) ***REMOVED***
	var plugins types.PluginsListResponse
	query := url.Values***REMOVED******REMOVED***

	if filter.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToParamWithVersion(cli.version, filter)
		if err != nil ***REMOVED***
			return plugins, err
		***REMOVED***
		query.Set("filters", filterJSON)
	***REMOVED***
	resp, err := cli.get(ctx, "/plugins", query, nil)
	if err != nil ***REMOVED***
		return plugins, wrapResponseError(err, resp, "plugin", "")
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&plugins)
	ensureReaderClosed(resp)
	return plugins, err
***REMOVED***
