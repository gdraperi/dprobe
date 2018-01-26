package client

import (
	"path"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// Ping pings the server and returns the value of the "Docker-Experimental", "OS-Type" & "API-Version" headers
func (cli *Client) Ping(ctx context.Context) (types.Ping, error) ***REMOVED***
	var ping types.Ping
	req, err := cli.buildRequest("GET", path.Join(cli.basePath, "/_ping"), nil, nil)
	if err != nil ***REMOVED***
		return ping, err
	***REMOVED***
	serverResp, err := cli.doRequest(ctx, req)
	if err != nil ***REMOVED***
		return ping, err
	***REMOVED***
	defer ensureReaderClosed(serverResp)

	if serverResp.header != nil ***REMOVED***
		ping.APIVersion = serverResp.header.Get("API-Version")

		if serverResp.header.Get("Docker-Experimental") == "true" ***REMOVED***
			ping.Experimental = true
		***REMOVED***
		ping.OSType = serverResp.header.Get("OSType")
	***REMOVED***
	return ping, cli.checkResponseErr(serverResp)
***REMOVED***
