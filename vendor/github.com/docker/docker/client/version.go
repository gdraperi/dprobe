package client

import (
	"encoding/json"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ServerVersion returns information of the docker client and server host.
func (cli *Client) ServerVersion(ctx context.Context) (types.Version, error) ***REMOVED***
	resp, err := cli.get(ctx, "/version", nil, nil)
	if err != nil ***REMOVED***
		return types.Version***REMOVED******REMOVED***, err
	***REMOVED***

	var server types.Version
	err = json.NewDecoder(resp.body).Decode(&server)
	ensureReaderClosed(resp)
	return server, err
***REMOVED***
