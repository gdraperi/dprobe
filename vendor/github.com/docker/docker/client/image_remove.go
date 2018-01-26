package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ImageRemove removes an image from the docker host.
func (cli *Client) ImageRemove(ctx context.Context, imageID string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***

	if options.Force ***REMOVED***
		query.Set("force", "1")
	***REMOVED***
	if !options.PruneChildren ***REMOVED***
		query.Set("noprune", "1")
	***REMOVED***

	var dels []types.ImageDeleteResponseItem
	resp, err := cli.delete(ctx, "/images/"+imageID, query, nil)
	if err != nil ***REMOVED***
		return dels, wrapResponseError(err, resp, "image", imageID)
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&dels)
	ensureReaderClosed(resp)
	return dels, err
***REMOVED***
