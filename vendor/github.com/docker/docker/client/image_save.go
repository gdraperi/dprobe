package client

import (
	"io"
	"net/url"

	"golang.org/x/net/context"
)

// ImageSave retrieves one or more images from the docker host as an io.ReadCloser.
// It's up to the caller to store the images and close the stream.
func (cli *Client) ImageSave(ctx context.Context, imageIDs []string) (io.ReadCloser, error) ***REMOVED***
	query := url.Values***REMOVED***
		"names": imageIDs,
	***REMOVED***

	resp, err := cli.get(ctx, "/images/get", query, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp.body, nil
***REMOVED***
