package client

import (
	"encoding/json"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// NetworkCreate creates a new network in the docker host.
func (cli *Client) NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error) ***REMOVED***
	networkCreateRequest := types.NetworkCreateRequest***REMOVED***
		NetworkCreate: options,
		Name:          name,
	***REMOVED***
	var response types.NetworkCreateResponse
	serverResp, err := cli.post(ctx, "/networks/create", nil, networkCreateRequest, nil)
	if err != nil ***REMOVED***
		return response, err
	***REMOVED***

	json.NewDecoder(serverResp.body).Decode(&response)
	ensureReaderClosed(serverResp)
	return response, err
***REMOVED***
