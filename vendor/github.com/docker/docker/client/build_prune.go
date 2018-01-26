package client

import (
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// BuildCachePrune requests the daemon to delete unused cache data
func (cli *Client) BuildCachePrune(ctx context.Context) (*types.BuildCachePruneReport, error) ***REMOVED***
	if err := cli.NewVersionError("1.31", "build prune"); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	report := types.BuildCachePruneReport***REMOVED******REMOVED***

	serverResp, err := cli.post(ctx, "/build/prune", nil, nil, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer ensureReaderClosed(serverResp)

	if err := json.NewDecoder(serverResp.body).Decode(&report); err != nil ***REMOVED***
		return nil, fmt.Errorf("Error retrieving disk usage: %v", err)
	***REMOVED***

	return &report, nil
***REMOVED***
