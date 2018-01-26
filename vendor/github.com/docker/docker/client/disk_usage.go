package client

import (
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// DiskUsage requests the current data usage from the daemon
func (cli *Client) DiskUsage(ctx context.Context) (types.DiskUsage, error) ***REMOVED***
	var du types.DiskUsage

	serverResp, err := cli.get(ctx, "/system/df", nil, nil)
	if err != nil ***REMOVED***
		return du, err
	***REMOVED***
	defer ensureReaderClosed(serverResp)

	if err := json.NewDecoder(serverResp.body).Decode(&du); err != nil ***REMOVED***
		return du, fmt.Errorf("Error retrieving disk usage: %v", err)
	***REMOVED***

	return du, nil
***REMOVED***
