package client

import "golang.org/x/net/context"

// SecretRemove removes a Secret.
func (cli *Client) SecretRemove(ctx context.Context, id string) error ***REMOVED***
	if err := cli.NewVersionError("1.25", "secret remove"); err != nil ***REMOVED***
		return err
	***REMOVED***
	resp, err := cli.delete(ctx, "/secrets/"+id, nil, nil)
	ensureReaderClosed(resp)
	return wrapResponseError(err, resp, "secret", id)
***REMOVED***
