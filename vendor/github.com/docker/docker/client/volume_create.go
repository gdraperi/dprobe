package client

import (
	"encoding/json"

	"github.com/docker/docker/api/types"
	volumetypes "github.com/docker/docker/api/types/volume"
	"golang.org/x/net/context"
)

// VolumeCreate creates a volume in the docker host.
func (cli *Client) VolumeCreate(ctx context.Context, options volumetypes.VolumesCreateBody) (types.Volume, error) ***REMOVED***
	var volume types.Volume
	resp, err := cli.post(ctx, "/volumes/create", nil, options, nil)
	if err != nil ***REMOVED***
		return volume, err
	***REMOVED***
	err = json.NewDecoder(resp.body).Decode(&volume)
	ensureReaderClosed(resp)
	return volume, err
***REMOVED***
