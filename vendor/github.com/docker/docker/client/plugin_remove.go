package client

import (
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// PluginRemove removes a plugin
func (cli *Client) PluginRemove(ctx context.Context, name string, options types.PluginRemoveOptions) error ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	if options.Force ***REMOVED***
		query.Set("force", "1")
	***REMOVED***

	resp, err := cli.delete(ctx, "/plugins/"+name, query, nil)
	ensureReaderClosed(resp)
	return wrapResponseError(err, resp, "plugin", name)
***REMOVED***
