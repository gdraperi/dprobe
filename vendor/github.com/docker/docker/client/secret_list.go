package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// SecretList returns the list of secrets.
func (cli *Client) SecretList(ctx context.Context, options types.SecretListOptions) ([]swarm.Secret, error) ***REMOVED***
	if err := cli.NewVersionError("1.25", "secret list"); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	query := url.Values***REMOVED******REMOVED***

	if options.Filters.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToJSON(options.Filters)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		query.Set("filters", filterJSON)
	***REMOVED***

	resp, err := cli.get(ctx, "/secrets", query, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var secrets []swarm.Secret
	err = json.NewDecoder(resp.body).Decode(&secrets)
	ensureReaderClosed(resp)
	return secrets, err
***REMOVED***
