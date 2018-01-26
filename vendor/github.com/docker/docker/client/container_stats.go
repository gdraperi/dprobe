package client

import (
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ContainerStats returns near realtime stats for a given container.
// It's up to the caller to close the io.ReadCloser returned.
func (cli *Client) ContainerStats(ctx context.Context, containerID string, stream bool) (types.ContainerStats, error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	query.Set("stream", "0")
	if stream ***REMOVED***
		query.Set("stream", "1")
	***REMOVED***

	resp, err := cli.get(ctx, "/containers/"+containerID+"/stats", query, nil)
	if err != nil ***REMOVED***
		return types.ContainerStats***REMOVED******REMOVED***, err
	***REMOVED***

	osType := getDockerOS(resp.header.Get("Server"))
	return types.ContainerStats***REMOVED***Body: resp.body, OSType: osType***REMOVED***, err
***REMOVED***
