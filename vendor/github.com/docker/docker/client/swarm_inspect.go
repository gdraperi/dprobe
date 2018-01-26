package client

import (
	"encoding/json"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// SwarmInspect inspects the swarm.
func (cli *Client) SwarmInspect(ctx context.Context) (swarm.Swarm, error) ***REMOVED***
	serverResp, err := cli.get(ctx, "/swarm", nil, nil)
	if err != nil ***REMOVED***
		return swarm.Swarm***REMOVED******REMOVED***, err
	***REMOVED***

	var response swarm.Swarm
	err = json.NewDecoder(serverResp.body).Decode(&response)
	ensureReaderClosed(serverResp)
	return response, err
***REMOVED***
