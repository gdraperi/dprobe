package client

import (
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// SecretCreate creates a new Secret.
func (cli *Client) SecretCreate(ctx context.Context, secret swarm.SecretSpec) (types.SecretCreateResponse, error) ***REMOVED***
	var response types.SecretCreateResponse
	if err := cli.NewVersionError("1.25", "secret create"); err != nil ***REMOVED***
		return response, err
	***REMOVED***
	resp, err := cli.post(ctx, "/secrets/create", nil, secret, nil)
	if err != nil ***REMOVED***
		return response, err
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
***REMOVED***
