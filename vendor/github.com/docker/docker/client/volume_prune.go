package client

import (
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

// VolumesPrune requests the daemon to delete unused data
func (cli *Client) VolumesPrune(ctx context.Context, pruneFilters filters.Args) (types.VolumesPruneReport, error) ***REMOVED***
	var report types.VolumesPruneReport

	if err := cli.NewVersionError("1.25", "volume prune"); err != nil ***REMOVED***
		return report, err
	***REMOVED***

	query, err := getFiltersQuery(pruneFilters)
	if err != nil ***REMOVED***
		return report, err
	***REMOVED***

	serverResp, err := cli.post(ctx, "/volumes/prune", query, nil, nil)
	if err != nil ***REMOVED***
		return report, err
	***REMOVED***
	defer ensureReaderClosed(serverResp)

	if err := json.NewDecoder(serverResp.body).Decode(&report); err != nil ***REMOVED***
		return report, fmt.Errorf("Error retrieving volume prune report: %v", err)
	***REMOVED***

	return report, nil
***REMOVED***
