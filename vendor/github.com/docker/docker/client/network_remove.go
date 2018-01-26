package client

import "golang.org/x/net/context"

// NetworkRemove removes an existent network from the docker host.
func (cli *Client) NetworkRemove(ctx context.Context, networkID string) error ***REMOVED***
	resp, err := cli.delete(ctx, "/networks/"+networkID, nil, nil)
	ensureReaderClosed(resp)
	return wrapResponseError(err, resp, "network", networkID)
***REMOVED***
