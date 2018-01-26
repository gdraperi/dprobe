package client

import (
	"io"
	"net/url"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
)

// ImageLoad loads an image in the docker host from the client host.
// It's up to the caller to close the io.ReadCloser in the
// ImageLoadResponse returned by this function.
func (cli *Client) ImageLoad(ctx context.Context, input io.Reader, quiet bool) (types.ImageLoadResponse, error) ***REMOVED***
	v := url.Values***REMOVED******REMOVED***
	v.Set("quiet", "0")
	if quiet ***REMOVED***
		v.Set("quiet", "1")
	***REMOVED***
	headers := map[string][]string***REMOVED***"Content-Type": ***REMOVED***"application/x-tar"***REMOVED******REMOVED***
	resp, err := cli.postRaw(ctx, "/images/load", v, input, headers)
	if err != nil ***REMOVED***
		return types.ImageLoadResponse***REMOVED******REMOVED***, err
	***REMOVED***
	return types.ImageLoadResponse***REMOVED***
		Body: resp.body,
		JSON: resp.header.Get("Content-Type") == "application/json",
	***REMOVED***, nil
***REMOVED***
