package client

import (
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

// ImagesPrune requests the daemon to delete unused data
func (cli *Client) ImagesPrune(ctx context.Context, pruneFilters filters.Args) (types.ImagesPruneReport, error) ***REMOVED***
	var report types.ImagesPruneReport

	if err := cli.NewVersionError("1.25", "image prune"); err != nil ***REMOVED***
		return report, err
	***REMOVED***

	query, err := getFiltersQuery(pruneFilters)
	if err != nil ***REMOVED***
		return report, err
	***REMOVED***

	serverResp, err := cli.post(ctx, "/images/prune", query, nil, nil)
	if err != nil ***REMOVED***
		return report, err
	***REMOVED***
	defer ensureReaderClosed(serverResp)

	if err := json.NewDecoder(serverResp.body).Decode(&report); err != nil ***REMOVED***
		return report, fmt.Errorf("Error retrieving disk usage: %v", err)
	***REMOVED***

	return report, nil
***REMOVED***
