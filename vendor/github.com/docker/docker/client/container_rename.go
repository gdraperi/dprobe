package client

import (
	"net/url"

	"golang.org/x/net/context"
)

// ContainerRename changes the name of a given container.
func (cli *Client) ContainerRename(ctx context.Context, containerID, newContainerName string) error ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	query.Set("name", newContainerName)
	resp, err := cli.post(ctx, "/containers/"+containerID+"/rename", query, nil, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***
