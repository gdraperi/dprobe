package client

import (
	"encoding/json"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// SwarmGetUnlockKey retrieves the swarm's unlock key.
func (cli *Client) SwarmGetUnlockKey(ctx context.Context) (types.SwarmUnlockKeyResponse, error) ***REMOVED***
	serverResp, err := cli.get(ctx, "/swarm/unlockkey", nil, nil)
	if err != nil ***REMOVED***
		return types.SwarmUnlockKeyResponse***REMOVED******REMOVED***, err
	***REMOVED***

	var response types.SwarmUnlockKeyResponse
	err = json.NewDecoder(serverResp.body).Decode(&response)
	ensureReaderClosed(serverResp)
	return response, err
***REMOVED***
