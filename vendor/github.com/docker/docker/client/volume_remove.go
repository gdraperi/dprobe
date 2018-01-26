package client

import (
	"net/url"

	"github.com/docker/docker/api/types/versions"
	"golang.org/x/net/context"
)

// VolumeRemove removes a volume from the docker host.
func (cli *Client) VolumeRemove(ctx context.Context, volumeID string, force bool) error ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	if versions.GreaterThanOrEqualTo(cli.version, "1.25") ***REMOVED***
		if force ***REMOVED***
			query.Set("force", "1")
		***REMOVED***
	***REMOVED***
	resp, err := cli.delete(ctx, "/volumes/"+volumeID, query, nil)
	ensureReaderClosed(resp)
	return wrapResponseError(err, resp, "volume", volumeID)
***REMOVED***
