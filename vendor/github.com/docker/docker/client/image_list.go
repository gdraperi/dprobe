package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/versions"
	"golang.org/x/net/context"
)

// ImageList returns a list of images in the docker host.
func (cli *Client) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) ***REMOVED***
	var images []types.ImageSummary
	query := url.Values***REMOVED******REMOVED***

	optionFilters := options.Filters
	referenceFilters := optionFilters.Get("reference")
	if versions.LessThan(cli.version, "1.25") && len(referenceFilters) > 0 ***REMOVED***
		query.Set("filter", referenceFilters[0])
		for _, filterValue := range referenceFilters ***REMOVED***
			optionFilters.Del("reference", filterValue)
		***REMOVED***
	***REMOVED***
	if optionFilters.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToParamWithVersion(cli.version, optionFilters)
		if err != nil ***REMOVED***
			return images, err
		***REMOVED***
		query.Set("filters", filterJSON)
	***REMOVED***
	if options.All ***REMOVED***
		query.Set("all", "1")
	***REMOVED***

	serverResp, err := cli.get(ctx, "/images/json", query, nil)
	if err != nil ***REMOVED***
		return images, err
	***REMOVED***

	err = json.NewDecoder(serverResp.body).Decode(&images)
	ensureReaderClosed(serverResp)
	return images, err
***REMOVED***
