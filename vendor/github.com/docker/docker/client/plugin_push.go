package client

import (
	"io"

	"golang.org/x/net/context"
)

// PluginPush pushes a plugin to a registry
func (cli *Client) PluginPush(ctx context.Context, name string, registryAuth string) (io.ReadCloser, error) ***REMOVED***
	headers := map[string][]string***REMOVED***"X-Registry-Auth": ***REMOVED***registryAuth***REMOVED******REMOVED***
	resp, err := cli.post(ctx, "/plugins/"+name+"/push", nil, nil, headers)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp.body, nil
***REMOVED***
