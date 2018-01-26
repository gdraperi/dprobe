package client

import (
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

// NetworksPrune requests the daemon to delete unused networks
func (cli *Client) NetworksPrune(ctx context.Context, pruneFilters filters.Args) (types.NetworksPruneReport, error) ***REMOVED***
	var report types.NetworksPruneReport

	if err := cli.NewVersionError("1.25", "network prune"); err != nil ***REMOVED***
		return report, err
	***REMOVED***

	query, err := getFiltersQuery(pruneFilters)
	if err != nil ***REMOVED***
		return report, err
	***REMOVED***

	serverResp, err := cli.post(ctx, "/networks/prune", query, nil, nil)
	if err != nil ***REMOVED***
		return report, err
	***REMOVED***
	defer ensureReaderClosed(serverResp)

	if err := json.NewDecoder(serverResp.body).Decode(&report); err != nil ***REMOVED***
		return report, fmt.Errorf("Error retrieving network prune report: %v", err)
	***REMOVED***

	return report, nil
***REMOVED***
