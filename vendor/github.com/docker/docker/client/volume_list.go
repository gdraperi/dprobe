package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"golang.org/x/net/context"
)

// VolumeList returns the volumes configured in the docker host.
func (cli *Client) VolumeList(ctx context.Context, filter filters.Args) (volumetypes.VolumesListOKBody, error) ***REMOVED***
	var volumes volumetypes.VolumesListOKBody
	query := url.Values***REMOVED******REMOVED***

	if filter.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToParamWithVersion(cli.version, filter)
		if err != nil ***REMOVED***
			return volumes, err
		***REMOVED***
		query.Set("filters", filterJSON)
	***REMOVED***
	resp, err := cli.get(ctx, "/volumes", query, nil)
	if err != nil ***REMOVED***
		return volumes, err
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&volumes)
	ensureReaderClosed(resp)
	return volumes, err
***REMOVED***
