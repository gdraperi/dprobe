package client

import "golang.org/x/net/context"

// ConfigRemove removes a Config.
func (cli *Client) ConfigRemove(ctx context.Context, id string) error ***REMOVED***
	if err := cli.NewVersionError("1.30", "config remove"); err != nil ***REMOVED***
		return err
	***REMOVED***
	resp, err := cli.delete(ctx, "/configs/"+id, nil, nil)
	ensureReaderClosed(resp)
	return wrapResponseError(err, resp, "config", id)
***REMOVED***
