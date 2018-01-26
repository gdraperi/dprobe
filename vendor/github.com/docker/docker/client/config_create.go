package client

import (
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// ConfigCreate creates a new Config.
func (cli *Client) ConfigCreate(ctx context.Context, config swarm.ConfigSpec) (types.ConfigCreateResponse, error) ***REMOVED***
	var response types.ConfigCreateResponse
	if err := cli.NewVersionError("1.30", "config create"); err != nil ***REMOVED***
		return response, err
	***REMOVED***
	resp, err := cli.post(ctx, "/configs/create", nil, config, nil)
	if err != nil ***REMOVED***
		return response, err
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
***REMOVED***
