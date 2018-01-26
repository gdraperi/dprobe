package client

import (
	"io"
	"net/url"

	"golang.org/x/net/context"
)

// ContainerExport retrieves the raw contents of a container
// and returns them as an io.ReadCloser. It's up to the caller
// to close the stream.
func (cli *Client) ContainerExport(ctx context.Context, containerID string) (io.ReadCloser, error) ***REMOVED***
	serverResp, err := cli.get(ctx, "/containers/"+containerID+"/export", url.Values***REMOVED******REMOVED***, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return serverResp.body, nil
***REMOVED***
