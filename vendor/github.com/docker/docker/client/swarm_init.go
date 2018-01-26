package client

import (
	"encoding/json"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// SwarmInit initializes the swarm.
func (cli *Client) SwarmInit(ctx context.Context, req swarm.InitRequest) (string, error) ***REMOVED***
	serverResp, err := cli.post(ctx, "/swarm/init", nil, req, nil)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	var response string
	err = json.NewDecoder(serverResp.body).Decode(&response)
	ensureReaderClosed(serverResp)
	return response, err
***REMOVED***
