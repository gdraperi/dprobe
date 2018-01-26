package client

import (
	"net/url"
	"strconv"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ContainerResize changes the size of the tty for a container.
func (cli *Client) ContainerResize(ctx context.Context, containerID string, options types.ResizeOptions) error ***REMOVED***
	return cli.resize(ctx, "/containers/"+containerID, options.Height, options.Width)
***REMOVED***

// ContainerExecResize changes the size of the tty for an exec process running inside a container.
func (cli *Client) ContainerExecResize(ctx context.Context, execID string, options types.ResizeOptions) error ***REMOVED***
	return cli.resize(ctx, "/exec/"+execID, options.Height, options.Width)
***REMOVED***

func (cli *Client) resize(ctx context.Context, basePath string, height, width uint) error ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	query.Set("h", strconv.Itoa(int(height)))
	query.Set("w", strconv.Itoa(int(width)))

	resp, err := cli.post(ctx, basePath+"/resize", query, nil, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***
