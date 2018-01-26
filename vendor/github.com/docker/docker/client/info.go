package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// Info returns information about the docker server.
func (cli *Client) Info(ctx context.Context) (types.Info, error) ***REMOVED***
	var info types.Info
	serverResp, err := cli.get(ctx, "/info", url.Values***REMOVED******REMOVED***, nil)
	if err != nil ***REMOVED***
		return info, err
	***REMOVED***
	defer ensureReaderClosed(serverResp)

	if err := json.NewDecoder(serverResp.body).Decode(&info); err != nil ***REMOVED***
		return info, fmt.Errorf("Error reading remote info: %v", err)
	***REMOVED***

	return info, nil
***REMOVED***
