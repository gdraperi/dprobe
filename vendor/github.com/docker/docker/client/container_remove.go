package client

import (
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ContainerRemove kills and removes a container from the docker host.
func (cli *Client) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	if options.RemoveVolumes ***REMOVED***
		query.Set("v", "1")
	***REMOVED***
	if options.RemoveLinks ***REMOVED***
		query.Set("link", "1")
	***REMOVED***

	if options.Force ***REMOVED***
		query.Set("force", "1")
	***REMOVED***

	resp, err := cli.delete(ctx, "/containers/"+containerID, query, nil)
	ensureReaderClosed(resp)
	return wrapResponseError(err, resp, "container", containerID)
***REMOVED***
