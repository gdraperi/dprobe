package client

import (
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// PluginDisable disables a plugin
func (cli *Client) PluginDisable(ctx context.Context, name string, options types.PluginDisableOptions) error ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	if options.Force ***REMOVED***
		query.Set("force", "1")
	***REMOVED***
	resp, err := cli.post(ctx, "/plugins/"+name+"/disable", query, nil, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***
