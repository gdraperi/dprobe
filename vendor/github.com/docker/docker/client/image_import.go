package client

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
)

// ImageImport creates a new image based in the source options.
// It returns the JSON content in the response body.
func (cli *Client) ImageImport(ctx context.Context, source types.ImageImportSource, ref string, options types.ImageImportOptions) (io.ReadCloser, error) ***REMOVED***
	if ref != "" ***REMOVED***
		//Check if the given image name can be resolved
		if _, err := reference.ParseNormalizedNamed(ref); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	query := url.Values***REMOVED******REMOVED***
	query.Set("fromSrc", source.SourceName)
	query.Set("repo", ref)
	query.Set("tag", options.Tag)
	query.Set("message", options.Message)
	if options.Platform != "" ***REMOVED***
		query.Set("platform", strings.ToLower(options.Platform))
	***REMOVED***
	for _, change := range options.Changes ***REMOVED***
		query.Add("changes", change)
	***REMOVED***

	resp, err := cli.postRaw(ctx, "/images/create", query, source.Source, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp.body, nil
***REMOVED***
