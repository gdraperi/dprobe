package client

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
)

// ImageCreate creates a new image based in the parent options.
// It returns the JSON content in the response body.
func (cli *Client) ImageCreate(ctx context.Context, parentReference string, options types.ImageCreateOptions) (io.ReadCloser, error) ***REMOVED***
	ref, err := reference.ParseNormalizedNamed(parentReference)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	query := url.Values***REMOVED******REMOVED***
	query.Set("fromImage", reference.FamiliarName(ref))
	query.Set("tag", getAPITagFromNamedRef(ref))
	if options.Platform != "" ***REMOVED***
		query.Set("platform", strings.ToLower(options.Platform))
	***REMOVED***
	resp, err := cli.tryImageCreate(ctx, query, options.RegistryAuth)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp.body, nil
***REMOVED***

func (cli *Client) tryImageCreate(ctx context.Context, query url.Values, registryAuth string) (serverResponse, error) ***REMOVED***
	headers := map[string][]string***REMOVED***"X-Registry-Auth": ***REMOVED***registryAuth***REMOVED******REMOVED***
	return cli.post(ctx, "/images/create", query, nil, headers)
***REMOVED***
